import type { FastifyInstance, FastifyRequest, preHandlerAsyncHookHandler } from 'fastify';
import fp from 'fastify-plugin';
import { withTx, type Db } from '../db/tx.js';
import { AppError, forbidden, unauthorized } from '../lib/errors.js';
import type { AuthContext } from '../types.js';
import { countApiCall, currentSubscription, monthStart, pendingApiCalls } from '../modules/subscriptions/service.js';

const READS = new Set(['GET', 'HEAD', 'OPTIONS']);
/** Paths that stay usable when the subscription has expired (sign in/out, renewing). */
const ALWAYS_ALLOWED = /^\/api\/(auth|subscription)(\/|$|\?)/;
// Not counted or limited: signing in, the subscription page itself, and platform administration.
const UNMETERED = /^\/api\/(auth|admin|subscription)(\/|$|\?)/;

export async function loadPermissions(db: Db, tenantId: string, userId: string): Promise<Set<string>> {
  const { rows } = await db.query<{ code: string }>(
    `SELECT DISTINCT p.code
       FROM user_roles ur
       JOIN roles r ON r.id = ur.role_id AND r.deleted_at IS NULL
       JOIN role_permissions rp ON rp.role_id = r.id
       JOIN permissions p ON p.id = rp.permission_id
      WHERE ur.tenant_id = $1 AND ur.user_id = $2`,
    [tenantId, userId],
  );
  return new Set(rows.map((r) => r.code));
}

/**
 * Verifies the access token, then re-checks the database on every request:
 * the session is not revoked, the user is active, the membership in the
 * token's tenant is active, and the tenant is not suspended. Permissions are
 * loaded fresh, so role changes apply immediately.
 */
async function authenticate(this: FastifyInstance, req: FastifyRequest): Promise<void> {
  if (req.auth) return;
  const header = req.headers.authorization;
  if (!header?.startsWith('Bearer ')) throw unauthorized();
  const claims = await this.deps.tokens.verifyAccess(header.slice(7));

  req.auth = await withTx(this.deps.pool, { tenantId: claims.tid, userId: claims.sub }, async (db) => {
    const { rows } = await db.query<{
      user_status: string; user_deleted: Date | null; is_platform_admin: boolean;
      member_status: string; is_owner: boolean; tenant_status: string;
      session_revoked: Date | null; session_expires: Date; session_tenant: string;
    }>(
      `SELECT u.status AS user_status, u.deleted_at AS user_deleted, u.is_platform_admin,
              ut.status AS member_status, ut.is_owner, t.status AS tenant_status,
              s.revoked_at AS session_revoked, s.expires_at AS session_expires, s.tenant_id AS session_tenant
         FROM users u
         JOIN user_tenants ut ON ut.user_id = u.id AND ut.tenant_id = $2
         JOIN tenants t ON t.id = ut.tenant_id
         JOIN user_sessions s ON s.id = $3 AND s.user_id = u.id
        WHERE u.id = $1`,
      [claims.sub, claims.tid, claims.sid],
    );
    const r = rows[0];
    if (!r || r.session_revoked || r.session_expires < new Date() || r.session_tenant !== claims.tid) {
      throw unauthorized('Session is no longer valid');
    }
    if (r.user_status !== 'ACTIVE' || r.user_deleted) throw unauthorized('Account is disabled');
    if (r.member_status !== 'ACTIVE') throw forbidden('Your access to this organization is disabled');
    if (r.tenant_status !== 'ACTIVE') throw new AppError(403, 'TENANT_SUSPENDED', 'This organization is suspended');

    // Subscription: an expired one leaves the data readable but refuses changes.
    const sub = await currentSubscription(db, claims.tid);
    const url = req.url;
    const adminRequest = r.is_platform_admin && url.startsWith('/api/admin');
    if (!sub.writable && !READS.has(req.method) && !ALWAYS_ALLOWED.test(url) && !adminRequest) {
      throw new AppError(402, 'SUBSCRIPTION_INACTIVE', 'انتهى الاشتراك؛ البيانات متاحة للاطلاع فقط حتى التجديد');
    }
    // Monthly API call allowance (counted in memory, flushed periodically).
    if (!UNMETERED.test(url)) {
      const limit = sub.limits.max_api_calls_per_month;
      if (limit !== null) {
        const { rows: [u] } = await db.query<{ quantity: string }>(
          `SELECT quantity::text FROM usage_records WHERE tenant_id = $1 AND metric = 'API_CALLS' AND period = $2::date`, [claims.tid, monthStart()]);
        if (Number(u?.quantity ?? 0) + pendingApiCalls(claims.tid) >= limit) {
          throw new AppError(429, 'PLAN_LIMIT_REACHED', 'بلغت المنشأة الحد الشهري لطلبات API في باقتها');
        }
      }
      countApiCall(claims.tid);
    }

    const ctx: AuthContext = {
      userId: claims.sub,
      tenantId: claims.tid,
      sessionId: claims.sid,
      isOwner: r.is_owner,
      isPlatformAdmin: r.is_platform_admin,
      permissions: await loadPermissions(db, claims.tid, claims.sub),
      subscription: { state: sub.state, writable: sub.writable },
    };
    return ctx;
  });
}

export default fp(async (app) => {
  app.decorate('authenticate', authenticate);
  app.decorateRequest('auth', null);
  app.decorateRequest('tenantTx', function <T>(this: FastifyRequest, fn: (db: Db) => Promise<T>) {
    if (!this.auth) throw unauthorized();
    return withTx(app.deps.pool, { tenantId: this.auth.tenantId, userId: this.auth.userId }, fn);
  });
  app.decorateRequest('auditMeta', function (this: FastifyRequest) {
    return { ip: this.ip, userAgent: this.headers['user-agent'] ?? null, requestId: String(this.id) };
  });
});

/** Route guard: authenticated and holding every listed permission. */
export function requirePermission(app: FastifyInstance, ...codes: string[]): preHandlerAsyncHookHandler {
  return async (req) => {
    await app.authenticate(req);
    const missing = codes.filter((c) => !req.auth!.permissions.has(c));
    if (missing.length) throw forbidden(`Missing permission: ${missing.join(', ')}`);
  };
}

export function requireAuth(app: FastifyInstance): preHandlerAsyncHookHandler {
  return async (req) => app.authenticate(req);
}
