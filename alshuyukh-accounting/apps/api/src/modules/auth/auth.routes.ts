import { assertWithinLimit, currentSubscription, featuresOf } from '../subscriptions/service.js';
import type { FastifyInstance, FastifyReply, FastifyRequest } from 'fastify';
import { z } from 'zod';
import { withTx } from '../../db/tx.js';
import { badRequest, notFound, unauthorized } from '../../lib/errors.js';
import { writeAudit } from '../audit/audit.service.js';
import { parse, password } from '../../lib/validation.js';
import { requireAuth } from '../../plugins/auth.js';
import { AuthService, type IssuedTokens } from './auth.service.js';

export const REFRESH_COOKIE = 'ash_rt';
const CSRF_HEADER = 'x-csrf-protection';

const email = z.string().trim().toLowerCase().pipe(z.email().max(254));

const registerBody = z.object({
  fullName: z.string().trim().min(2).max(200),
  email,
  password,
  tenantName: z.string().trim().min(2).max(200),
  companyName: z.string().trim().min(2).max(200),
  vatNumber: z.string().regex(/^3\d{13}3$/, 'Saudi VAT number must be 15 digits starting and ending with 3').nullish(),
  commercialRegistration: z.string().regex(/^\d{10}$/, 'Commercial registration must be 10 digits').nullish(),
});

const loginBody = z.object({
  email,
  password: z.string().min(1).max(200),
  tenantId: z.uuid().optional(),
});

/**
 * The refresh token travels only in an HttpOnly, SameSite=Strict cookie
 * scoped to /api/auth. Endpoints that read it also require a custom header,
 * which a cross-site form cannot send and which triggers a CORS preflight
 * for scripts on other origins (CSRF protection).
 */
function requireCsrfHeader(req: FastifyRequest) {
  if (req.headers[CSRF_HEADER] !== '1') throw badRequest('CSRF_HEADER_MISSING', `Missing ${CSRF_HEADER} header`);
}

export default async function authRoutes(app: FastifyInstance) {
  const service = new AuthService(app.deps);
  const { env } = app.deps;
  const isTest = env.NODE_ENV === 'test';
  // Credential endpoints: tight per-IP limit against password guessing.
  const strictLimit = { rateLimit: { max: isTest ? 1000 : 10, timeWindow: '1 minute' } };
  // Refresh runs on every page load, and an office may share one IP.
  const refreshLimit = { rateLimit: { max: isTest ? 1000 : 120, timeWindow: '1 minute' } };

  const setRefreshCookie = (reply: FastifyReply, token: string) =>
    reply.setCookie(REFRESH_COOKIE, token, {
      httpOnly: true,
      secure: env.NODE_ENV === 'production',
      sameSite: 'strict',
      path: '/api/auth',
      maxAge: env.REFRESH_TOKEN_TTL_DAYS * 86400,
    });

  const respond = async (reply: FastifyReply, t: IssuedTokens) => {
    setRefreshCookie(reply, t.refreshToken);
    return { accessToken: t.accessToken, expiresIn: t.expiresIn, tenantId: t.tenantId, userId: t.userId };
  };

  app.post('/register', { config: strictLimit }, async (req, reply) => {
    const body = parse(registerBody, req.body);
    reply.code(201);
    return respond(reply, await service.register(body, req.auditMeta()));
  });

  app.post('/login', { config: strictLimit }, async (req, reply) => {
    const body = parse(loginBody, req.body);
    return respond(reply, await service.login(body.email, body.password, body.tenantId, req.auditMeta()));
  });

  app.post('/refresh', { config: refreshLimit }, async (req, reply) => {
    requireCsrfHeader(req);
    const token = req.cookies[REFRESH_COOKIE];
    if (!token) throw unauthorized('Missing refresh token');
    try {
      return await respond(reply, await service.refresh(token, req.auditMeta()));
    } catch (err) {
      reply.clearCookie(REFRESH_COOKIE, { path: '/api/auth' });
      throw err;
    }
  });

  app.post('/logout', { preHandler: requireAuth(app) }, async (req, reply) => {
    requireCsrfHeader(req);
    const a = req.auth!;
    await service.logout(a.sessionId, a.userId, a.tenantId, req.auditMeta());
    reply.clearCookie(REFRESH_COOKIE, { path: '/api/auth' });
    return reply.code(204).send();
  });

  app.post('/switch-tenant', { preHandler: requireAuth(app) }, async (req) => {
    const { tenantId } = parse(z.object({ tenantId: z.uuid() }), req.body);
    const a = req.auth!;
    return service.switchTenant(a.sessionId, a.userId, a.tenantId, tenantId, req.auditMeta());
  });

  /** Organizations that invited the signed-in user (an existing account is never added without consent). */
  app.get('/invitations', { preHandler: requireAuth(app) }, async (req) => {
    const a = req.auth!;
    const ids = await withTx(app.deps.pool, { tenantId: a.tenantId, userId: a.userId }, async (db) =>
      (await db.query<{ tenant_id: string; invited_at: Date }>(
        `SELECT tenant_id, joined_at AS invited_at FROM user_tenants WHERE user_id = $1 AND status = 'INVITED'`, [a.userId])).rows);
    const data = [];
    // The tenant's name is visible only in that tenant's own context.
    for (const inv of ids) {
      const name = await withTx(app.deps.pool, { tenantId: inv.tenant_id, userId: a.userId }, async (db) =>
        (await db.query<{ name: string }>(`SELECT name FROM tenants WHERE id = $1 AND status = 'ACTIVE' AND deleted_at IS NULL`, [inv.tenant_id])).rows[0]?.name);
      if (name) data.push({ tenantId: inv.tenant_id, tenantName: name, invitedAt: inv.invited_at });
    }
    return { data };
  });

  for (const action of ['accept', 'decline'] as const) {
    app.post(`/invitations/:tenantId/${action}`, { preHandler: requireAuth(app) }, async (req, reply) => {
      const { tenantId } = parse(z.object({ tenantId: z.uuid() }), req.params);
      const a = req.auth!;
      await withTx(app.deps.pool, { tenantId, userId: a.userId }, async (db) => {
        const { rows: [m] } = await db.query(`SELECT 1 FROM user_tenants WHERE tenant_id = $1 AND user_id = $2 AND status = 'INVITED' FOR UPDATE`, [tenantId, a.userId]);
        if (!m) throw notFound('Invitation');
        if (action === 'accept') await assertWithinLimit(db, tenantId, 'max_users');
        await db.query(`UPDATE user_tenants SET status = $3 WHERE tenant_id = $1 AND user_id = $2`, [tenantId, a.userId, action === 'accept' ? 'ACTIVE' : 'DISABLED']);
        await writeAudit(db, { tenantId, userId: a.userId, action: 'UPDATE', entityType: 'user', entityId: a.userId, newValues: { invitation: action === 'accept' ? 'ACCEPTED' : 'DECLINED' } }, req.auditMeta());
      });
      return reply.code(204).send();
    });
  }

  app.post('/change-password', { preHandler: requireAuth(app), config: strictLimit }, async (req, reply) => {
    const body = parse(z.object({ currentPassword: z.string().min(1).max(200), newPassword: password }), req.body);
    const a = req.auth!;
    await service.changePassword(a.userId, a.tenantId, a.sessionId, body.currentPassword, body.newPassword, req.auditMeta());
    return reply.code(204).send();
  });

  app.get('/me', { preHandler: requireAuth(app) }, async (req) => {
    const a = req.auth!;
    const profile = await withTx(app.deps.pool, { tenantId: a.tenantId, userId: a.userId }, async (db) => {
      const { rows: [u] } = await db.query<{ id: string; email: string; full_name: string; must_change_password: boolean }>(
        `SELECT id, email, full_name, must_change_password FROM users WHERE id = $1`, [a.userId]);
      const { rows: roles } = await db.query<{ id: string; code: string; name_ar: string; name_en: string }>(
        `SELECT r.id, r.code, r.name_ar, r.name_en FROM user_roles ur JOIN roles r ON r.id = ur.role_id
          WHERE ur.tenant_id = $1 AND ur.user_id = $2 AND r.deleted_at IS NULL ORDER BY r.code`,
        [a.tenantId, a.userId]);
      const { rows: [tenant] } = await db.query<{ id: string; name: string; status: string }>(
        `SELECT id, name, status FROM tenants WHERE id = $1`, [a.tenantId]);
      const sub = await currentSubscription(db, a.tenantId);
      return { u: u!, roles, tenant: tenant!, sub, features: await featuresOf(db, a.tenantId) };
    });
    return {
      user: {
        id: profile.u.id, email: profile.u.email, fullName: profile.u.full_name,
        mustChangePassword: profile.u.must_change_password, isPlatformAdmin: a.isPlatformAdmin,
      },
      tenant: { ...profile.tenant, isOwner: a.isOwner },
      roles: profile.roles.map((r) => ({ id: r.id, code: r.code, nameAr: r.name_ar, nameEn: r.name_en })),
      permissions: [...a.permissions].sort(),
      memberships: await service.listMemberships(a.userId),
      subscription: {
        state: profile.sub.state, writable: profile.sub.writable, planName: profile.sub.planName,
        periodEnd: profile.sub.periodEnd, graceEnd: profile.sub.graceEnd,
      },
      features: profile.features,
    };
  });
}
