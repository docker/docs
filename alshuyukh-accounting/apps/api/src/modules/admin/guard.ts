import type { FastifyInstance, FastifyRequest, preHandlerAsyncHookHandler } from 'fastify';
import { withTx, type Db } from '../../db/tx.js';
import { forbidden } from '../../lib/errors.js';

/**
 * Platform administration (the SaaS operator, not a tenant role). The flag is
 * re-read from the database on every request by authenticate(); only then do
 * admin transactions turn on app.platform_admin, which lets RLS show every
 * tenant's rows read-only. Writes run in the target tenant's own context.
 */
export function requirePlatformAdmin(app: FastifyInstance): preHandlerAsyncHookHandler {
  return async (req) => {
    await app.authenticate(req);
    if (!req.auth!.isPlatformAdmin) throw forbidden('Platform administrators only');
  };
}

export function platformTx<T>(app: FastifyInstance, req: FastifyRequest, fn: (db: Db) => Promise<T>, tenantId: string | null = null): Promise<T> {
  if (!req.auth?.isPlatformAdmin) throw forbidden('Platform administrators only');
  return withTx(app.deps.pool, { tenantId, userId: req.auth.userId }, async (db) => {
    await db.query(`SELECT set_config('app.platform_admin', 'on', true)`);
    return fn(db);
  });
}

export async function platformAudit(db: Db, req: FastifyRequest, e: {
  action: string; entityType: string; entityId?: string | null; targetTenantId?: string | null; oldValues?: unknown; newValues?: unknown;
}) {
  await db.query(
    `INSERT INTO platform_audit_logs (admin_user_id, action, entity_type, entity_id, target_tenant_id, old_values, new_values, ip_address, user_agent)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
    [req.auth!.userId, e.action, e.entityType, e.entityId ?? null, e.targetTenantId ?? null,
     e.oldValues === undefined ? null : JSON.stringify(e.oldValues), e.newValues === undefined ? null : JSON.stringify(e.newValues),
     req.ip, req.headers['user-agent'] ?? null]);
}

/** Subscription state in SQL, matching viewOf() in the subscriptions service. */
export const STATE_SQL = `CASE WHEN s.id IS NULL THEN 'NONE'
  WHEN now() <= s.current_period_end THEN s.status
  WHEN now() <= s.current_period_end + make_interval(days => p.grace_days) THEN 'GRACE' ELSE 'EXPIRED' END`;
