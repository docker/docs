import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { notFound } from '../../lib/errors.js';
import { parse } from '../../lib/validation.js';
import { requireAuth, requirePermission } from '../../plugins/auth.js';
import { writeAudit } from '../audit/audit.service.js';
import { billingEvent, currentSubscription, featuresOf, LIMIT_KEYS, usageOf } from './service.js';

export const PLAN_PUBLIC_SELECT = `SELECT id, code, name_ar AS "nameAr", name_en AS "nameEn", description, currency,
  price_monthly::text AS "priceMonthly", price_yearly::text AS "priceYearly", trial_days AS "trialDays", grace_days AS "graceDays",
  ${LIMIT_KEYS.map((k) => `${k} AS "${k}"`).join(', ')}, is_public AS "isPublic", is_default AS "isDefault", is_active AS "isActive", sort_order AS "sortOrder"
  FROM plans`;

/** The organization's own subscription: state, limits and usage, and plan change requests. */
export default async function subscriptionRoutes(app: FastifyInstance) {
  app.get('/subscription', { preHandler: requireAuth(app) }, async (req) =>
    req.tenantTx(async (db) => {
      const sub = await currentSubscription(db, req.auth!.tenantId);
      const { rows: items } = await db.query(
        `SELECT code, description, quantity, unit_price::text AS "unitPrice", currency FROM subscription_items WHERE subscription_id = $1`, [sub.id]);
      return { ...sub, items, usage: await usageOf(db, req.auth!.tenantId), features: await featuresOf(db, req.auth!.tenantId) };
    }));

  app.get('/subscription/plans', { preHandler: requireAuth(app) }, async (req) =>
    req.tenantTx(async (db) => ({ data: (await db.query(`${PLAN_PUBLIC_SELECT} WHERE is_active AND is_public ORDER BY sort_order, price_monthly`)).rows })));

  app.get('/subscription/billing-events', { preHandler: requirePermission(app, 'subscription.manage') }, async (req) =>
    req.tenantTx(async (db) => ({
      data: (await db.query(
        `SELECT id, event_type AS "eventType", amount::text, currency, reference, details, created_at AS "createdAt"
           FROM billing_events WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT 100`, [req.auth!.tenantId])).rows,
    })));

  /**
   * Online payment is not integrated yet: the owner asks for a plan change and
   * the platform administrator applies it once payment is received.
   */
  app.post('/subscription/request-change', { preHandler: requirePermission(app, 'subscription.manage') }, async (req) => {
    const body = parse(z.object({ planId: z.uuid(), billingCycle: z.enum(['MONTHLY', 'YEARLY']), note: z.string().trim().max(500).nullish() }), req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const { rows: [plan] } = await db.query<{ code: string }>(`SELECT code FROM plans WHERE id = $1 AND is_active AND is_public`, [body.planId]);
      if (!plan) throw notFound('Plan');
      const sub = await currentSubscription(db, a.tenantId);
      await billingEvent(db, { tenantId: a.tenantId, subscriptionId: sub.id, type: 'PLAN_CHANGE_REQUESTED', userId: a.userId, details: { plan: plan.code, billingCycle: body.billingCycle, note: body.note ?? null } });
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'subscription', entityId: sub.id, newValues: { requestedPlan: plan.code, billingCycle: body.billingCycle } }, req.auditMeta());
      return { requested: true };
    });
  });
}
