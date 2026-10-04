import { randomBytes, randomUUID } from 'node:crypto';
import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { badRequest, notFound } from '../../lib/errors.js';
import { amountString, Decimal, toMoney } from '../../lib/money.js';
import { hashPassword } from '../../lib/password.js';
import { parse, uuidParam } from '../../lib/validation.js';
import { writeAudit } from '../audit/audit.service.js';
import { provisionTenant } from '../auth/auth.service.js';
import {
  addPeriod, billingEvent, changePlan, currentSubscription, featuresOf, LIMIT_KEYS, startSubscription, usageOf,
} from '../subscriptions/service.js';
import { platformAudit, platformTx, requirePlatformAdmin, STATE_SQL } from './guard.js';

const limitOverrides = z.object(Object.fromEntries(LIMIT_KEYS.map((k) => [k, z.number().int().positive().nullable().optional()]))).strict();

/** Organizations (tenants) and their subscriptions, for the platform administrator. */
export default async function adminTenantRoutes(app: FastifyInstance) {
  const guard = requirePlatformAdmin(app);

  app.get('/admin/tenants', { preHandler: guard }, async (req) => {
    const q = parse(z.object({
      search: z.string().trim().max(100).optional(), status: z.enum(['ACTIVE', 'SUSPENDED', 'CANCELLED']).optional(),
      limit: z.coerce.number().int().min(1).max(200).default(50), offset: z.coerce.number().int().min(0).default(0),
    }), req.query);
    return platformTx(app, req, async (db) => {
      const { rows } = await db.query<{ total: string }>(
        `SELECT count(*) OVER () AS total, t.id, t.name, t.slug, t.status, t.created_at AS "createdAt",
                (SELECT u.email FROM user_tenants ut JOIN users u ON u.id = ut.user_id WHERE ut.tenant_id = t.id AND ut.is_owner ORDER BY ut.created_at LIMIT 1) AS "ownerEmail",
                p.name_ar AS "planName", ${STATE_SQL} AS "subscriptionState", s.current_period_end AS "periodEnd",
                (SELECT count(*) FROM user_tenants ut WHERE ut.tenant_id = t.id AND ut.status = 'ACTIVE')::int AS users,
                (SELECT count(*) FROM companies c WHERE c.tenant_id = t.id AND c.deleted_at IS NULL)::int AS companies
           FROM tenants t
           LEFT JOIN subscriptions s ON s.tenant_id = t.id AND s.status <> 'CANCELLED'
           LEFT JOIN plans p ON p.id = s.plan_id
          WHERE t.deleted_at IS NULL AND ($1::text IS NULL OR t.name ILIKE '%' || $1 || '%' OR t.slug ILIKE '%' || $1 || '%'
                 OR EXISTS (SELECT 1 FROM user_tenants ut JOIN users u ON u.id = ut.user_id WHERE ut.tenant_id = t.id AND u.email ILIKE '%' || $1 || '%'))
            AND ($2::text IS NULL OR t.status = $2)
          ORDER BY t.created_at DESC LIMIT $3 OFFSET $4`,
        [q.search || null, q.status ?? null, q.limit, q.offset]);
      return { data: rows.map(({ total: _t, ...r }) => r), total: Number(rows[0]?.total ?? 0) };
    });
  });

  app.get('/admin/tenants/:id', { preHandler: guard }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return platformTx(app, req, async (db) => {
      const { rows: [t] } = await db.query(`SELECT id, name, slug, status, created_at AS "createdAt" FROM tenants WHERE id = $1 AND deleted_at IS NULL`, [id]);
      if (!t) throw notFound('Tenant');
      const sub = await currentSubscription(db, id);
      const [items, events, members, companies, overrides] = await Promise.all([
        db.query(`SELECT code, description, quantity, unit_price::text AS "unitPrice", currency FROM subscription_items WHERE subscription_id = $1`, [sub.id]),
        db.query(`SELECT b.id, b.event_type AS "eventType", b.amount::text, b.currency, b.reference, b.details, b.created_at AS "createdAt", u.email AS "byEmail"
                    FROM billing_events b LEFT JOIN users u ON u.id = b.created_by WHERE b.tenant_id = $1 ORDER BY b.created_at DESC LIMIT 50`, [id]),
        db.query(`SELECT u.id, u.email, u.full_name AS "fullName", u.status AS "userStatus", ut.status AS "memberStatus", ut.is_owner AS "isOwner", u.last_login_at AS "lastLoginAt"
                    FROM user_tenants ut JOIN users u ON u.id = ut.user_id WHERE ut.tenant_id = $1 ORDER BY ut.is_owner DESC, u.email`, [id]),
        db.query(`SELECT id, name, vat_number AS "vatNumber", commercial_registration AS "commercialRegistration", city FROM companies WHERE tenant_id = $1 AND deleted_at IS NULL ORDER BY created_at`, [id]),
        db.query(`SELECT flag_key AS key, enabled FROM tenant_feature_flags WHERE tenant_id = $1`, [id]),
      ]);
      return {
        ...t, subscription: { ...sub, items: items.rows }, usage: await usageOf(db, id), billingEvents: events.rows, members: members.rows,
        companies: companies.rows, features: await featuresOf(db, id), featureOverrides: overrides.rows,
      };
    });
  });

  /** Creates an organization with its owner; the temporary password is shown once. */
  app.post('/admin/tenants', { preHandler: guard }, async (req, reply) => {
    const body = parse(z.object({
      tenantName: z.string().trim().min(2).max(200), companyName: z.string().trim().min(2).max(200),
      ownerName: z.string().trim().min(2).max(200), ownerEmail: z.string().trim().toLowerCase().pipe(z.email().max(254)),
      planId: z.uuid().nullish(),
    }), req.body);
    const tenantId = randomUUID();
    const temporaryPassword = `${randomBytes(9).toString('base64url')}A1!`;
    const passwordHash = await hashPassword(temporaryPassword);
    await platformTx(app, req, async (db) => {
      await provisionTenant(db, {
        tenantId, passwordHash, slug: `org-${tenantId.slice(0, 8)}`, tenantName: body.tenantName, companyName: body.companyName,
        fullName: body.ownerName, email: body.ownerEmail, password: '', mustChangePassword: true, planId: body.planId ?? null, createdBy: req.auth!.userId,
      }, req.auditMeta());
      await platformAudit(db, req, { action: 'CREATE', entityType: 'tenant', entityId: tenantId, targetTenantId: tenantId, newValues: { ...body } });
    }, tenantId);
    reply.code(201);
    return { tenantId, ownerEmail: body.ownerEmail, temporaryPassword };
  });

  for (const [action, status, event] of [['suspend', 'SUSPENDED', 'TENANT_SUSPENDED'], ['activate', 'ACTIVE', 'TENANT_ACTIVATED']] as const) {
    app.post(`/admin/tenants/:id/${action}`, { preHandler: guard }, async (req) => {
      const { id } = parse(uuidParam, req.params);
      const { reason } = parse(z.object({ reason: z.string().trim().min(3).max(500) }), req.body);
      if (action === 'suspend' && id === req.auth!.tenantId) throw badRequest('OWN_TENANT', 'You cannot suspend the organization you are signed in to');
      return platformTx(app, req, async (db) => {
        const { rows: [t] } = await db.query<{ status: string }>(`SELECT status FROM tenants WHERE id = $1 AND deleted_at IS NULL`, [id]);
        if (!t) throw notFound('Tenant');
        await db.query(`UPDATE tenants SET status = $2 WHERE id = $1`, [id, status]);
        const sub = await currentSubscription(db, id);
        await billingEvent(db, { tenantId: id, subscriptionId: sub.id, type: event, userId: req.auth!.userId, details: { reason } });
        await writeAudit(db, { tenantId: id, userId: req.auth!.userId, action: 'SETTINGS_CHANGE', entityType: 'tenant', entityId: id, oldValues: { status: t.status }, newValues: { status, reason, by: 'platform' } }, req.auditMeta());
        await platformAudit(db, req, { action: action.toUpperCase(), entityType: 'tenant', entityId: id, targetTenantId: id, oldValues: { status: t.status }, newValues: { status, reason } });
        return { id, status };
      }, id);
    });
  }

  /** Change plan (or start a new subscription when none is current). */
  app.post('/admin/tenants/:id/subscription/plan', { preHandler: guard }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ planId: z.uuid(), billingCycle: z.enum(['MONTHLY', 'YEARLY']).default('MONTHLY') }), req.body);
    return platformTx(app, req, async (db) => {
      const before = await currentSubscription(db, id);
      const after = before.id ? await changePlan(db, id, req.auth!.userId, body.planId, body.billingCycle)
        : (await startSubscription(db, id, req.auth!.userId, body.planId), await currentSubscription(db, id));
      await platformAudit(db, req, { action: 'PLAN_CHANGE', entityType: 'subscription', entityId: after.id, targetTenantId: id, oldValues: { plan: before.planCode }, newValues: { plan: after.planCode, billingCycle: body.billingCycle } });
      return after;
    }, id);
  });

  /** Per-tenant limits that replace the plan's (null = unlimited); the whole object is replaced. */
  app.post('/admin/tenants/:id/subscription/limits', { preHandler: guard }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const { overrides } = parse(z.object({ overrides: limitOverrides }), req.body);
    return platformTx(app, req, async (db) => {
      const sub = await currentSubscription(db, id, true);
      if (!sub.id) throw notFound('Subscription');
      await db.query(`UPDATE subscriptions SET limit_overrides = $2 WHERE id = $1`, [sub.id, JSON.stringify(overrides)]);
      await billingEvent(db, { tenantId: id, subscriptionId: sub.id, type: 'LIMITS_CHANGED', userId: req.auth!.userId, details: { from: sub.limitOverrides, to: overrides } });
      await platformAudit(db, req, { action: 'LIMITS_CHANGE', entityType: 'subscription', entityId: sub.id, targetTenantId: id, oldValues: sub.limitOverrides, newValues: overrides });
      return currentSubscription(db, id);
    }, id);
  });

  /**
   * Records a payment received outside the system (bank transfer…) and extends
   * the period by whole billing cycles. Online payment is not integrated.
   */
  app.post('/admin/tenants/:id/subscription/payment', { preHandler: guard }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ amount: amountString, reference: z.string().trim().min(2).max(100), cycles: z.number().int().min(1).max(36).default(1) }), req.body);
    if (!new Decimal(body.amount).greaterThan(0)) throw badRequest('INVALID_AMOUNT', 'Amount must be greater than zero');
    return platformTx(app, req, async (db) => {
      const sub = await currentSubscription(db, id, true);
      if (!sub.id) throw badRequest('NO_SUBSCRIPTION', 'Choose a plan for this organization first');
      // An expired or trial subscription starts paying from today; an active one
      // (or one in its grace days) continues from the end of its period.
      const restart = sub.state === 'EXPIRED' || sub.status === 'TRIALING';
      const end = addPeriod(restart ? new Date() : sub.periodEnd!, sub.billingCycle as 'MONTHLY' | 'YEARLY', body.cycles);
      await db.query(
        `UPDATE subscriptions SET status = 'ACTIVE', current_period_start = CASE WHEN $3 THEN now() ELSE current_period_start END, current_period_end = $2 WHERE id = $1`,
        [sub.id, end, restart]);
      const { rows: [p] } = await db.query<{ currency: string }>(`SELECT currency FROM plans WHERE id = $1`, [sub.planId]);
      await billingEvent(db, { tenantId: id, subscriptionId: sub.id, type: 'PAYMENT_RECORDED', userId: req.auth!.userId, amount: toMoney(body.amount),
        currency: p!.currency, reference: body.reference, details: { cycles: body.cycles, billingCycle: sub.billingCycle, periodEnd: end } });
      await platformAudit(db, req, { action: 'PAYMENT', entityType: 'subscription', entityId: sub.id, targetTenantId: id, newValues: { ...body, periodEnd: end } });
      return currentSubscription(db, id);
    }, id);
  });

  /** Extends the current period (e.g. a longer trial) without a payment. */
  app.post('/admin/tenants/:id/subscription/extend', { preHandler: guard }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ days: z.number().int().min(1).max(365), reason: z.string().trim().min(3).max(500) }), req.body);
    return platformTx(app, req, async (db) => {
      const sub = await currentSubscription(db, id, true);
      if (!sub.id) throw badRequest('NO_SUBSCRIPTION', 'Choose a plan for this organization first');
      const base = sub.periodEnd! > new Date() ? sub.periodEnd! : new Date();
      const end = new Date(base.getTime() + body.days * 86_400_000);
      await db.query(`UPDATE subscriptions SET current_period_end = $2 WHERE id = $1`, [sub.id, end]);
      await billingEvent(db, { tenantId: id, subscriptionId: sub.id, type: 'PERIOD_EXTENDED', userId: req.auth!.userId, details: { days: body.days, reason: body.reason, periodEnd: end } });
      await platformAudit(db, req, { action: 'EXTEND', entityType: 'subscription', entityId: sub.id, targetTenantId: id, newValues: { ...body, periodEnd: end } });
      return currentSubscription(db, id);
    }, id);
  });

  app.post('/admin/tenants/:id/subscription/cancel', { preHandler: guard }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const { reason } = parse(z.object({ reason: z.string().trim().min(3).max(500) }), req.body);
    return platformTx(app, req, async (db) => {
      const sub = await currentSubscription(db, id, true);
      if (!sub.id) throw notFound('Subscription');
      await db.query(`UPDATE subscriptions SET status = 'CANCELLED', cancelled_at = now() WHERE id = $1`, [sub.id]);
      await billingEvent(db, { tenantId: id, subscriptionId: sub.id, type: 'SUBSCRIPTION_CANCELLED', userId: req.auth!.userId, details: { reason } });
      await platformAudit(db, req, { action: 'CANCEL', entityType: 'subscription', entityId: sub.id, targetTenantId: id, newValues: { reason } });
      return currentSubscription(db, id);
    }, id);
  });

  /** Per-tenant feature flag; enabled: null removes the override. */
  app.put('/admin/tenants/:id/features/:key', { preHandler: guard }, async (req) => {
    const { id, key } = parse(z.object({ id: z.uuid(), key: z.string().regex(/^[a-z][a-z0-9_]{1,62}$/) }), req.params);
    const { enabled } = parse(z.object({ enabled: z.boolean().nullable() }), req.body);
    return platformTx(app, req, async (db) => {
      const { rowCount } = await db.query(`SELECT 1 FROM feature_flags WHERE key = $1`, [key]);
      if (!rowCount) throw notFound('Feature flag');
      if (enabled === null) await db.query(`DELETE FROM tenant_feature_flags WHERE tenant_id = $1 AND flag_key = $2`, [id, key]);
      else await db.query(
        `INSERT INTO tenant_feature_flags (tenant_id, flag_key, enabled) VALUES ($1, $2, $3)
         ON CONFLICT (tenant_id, flag_key) DO UPDATE SET enabled = EXCLUDED.enabled, updated_at = now()`, [id, key, enabled]);
      await platformAudit(db, req, { action: 'FEATURE_OVERRIDE', entityType: 'feature_flag', entityId: key, targetTenantId: id, newValues: { enabled } });
      return { features: await featuresOf(db, id) };
    }, id);
  });
}
