import type pg from 'pg';
import type { Db } from '../../db/tx.js';
import { AppError, badRequest, notFound } from '../../lib/errors.js';
import { riyadhNow } from '../zatca/service.js';

/**
 * Subscription state and plan limits. Prices and limits live in the plans
 * table (edited by the platform administrator); NULL means unlimited.
 */

export const LIMIT_KEYS = ['max_users', 'max_companies', 'max_branches', 'max_warehouses', 'max_products',
  'max_invoices_per_month', 'max_storage_mb', 'max_api_calls_per_month'] as const;
export type LimitKey = (typeof LIMIT_KEYS)[number];
export type Limits = Record<LimitKey, number | null>;

/** TRIALING / ACTIVE: within the period. GRACE: period ended, still within grace days. */
export type SubscriptionState = 'TRIALING' | 'ACTIVE' | 'GRACE' | 'EXPIRED' | 'CANCELLED' | 'NONE';

export interface SubscriptionView {
  id: string | null; planId: string | null; planCode: string | null; planName: string | null; status: string | null;
  state: SubscriptionState; writable: boolean; billingCycle: string | null;
  periodStart: Date | null; periodEnd: Date | null; graceEnd: Date | null; limits: Limits; limitOverrides: Partial<Limits>;
}

const SUB_SQL = `
  SELECT s.id, s.plan_id, s.status, s.billing_cycle, s.current_period_start, s.current_period_end, s.limit_overrides,
         p.code AS plan_code, p.name_ar AS plan_name, p.grace_days, ${LIMIT_KEYS.map((k) => `p.${k}`).join(', ')}
    FROM subscriptions s JOIN plans p ON p.id = s.plan_id
   WHERE s.tenant_id = $1 AND s.status <> 'CANCELLED'`;

type SubRow = Record<string, any>;

export function viewOf(row: SubRow | undefined, now = new Date()): SubscriptionView {
  const emptyLimits = Object.fromEntries(LIMIT_KEYS.map((k) => [k, null])) as Limits;
  if (!row) {
    return { id: null, planId: null, planCode: null, planName: null, status: null, state: 'NONE', writable: false, billingCycle: null,
      periodStart: null, periodEnd: null, graceEnd: null, limits: emptyLimits, limitOverrides: {} };
  }
  const overrides = (row.limit_overrides ?? {}) as Partial<Limits>;
  const limits = Object.fromEntries(LIMIT_KEYS.map((k) => [k, k in overrides ? overrides[k] ?? null : row[k] ?? null])) as Limits;
  const end = new Date(row.current_period_end);
  const graceEnd = new Date(end.getTime() + Number(row.grace_days) * 86_400_000);
  const state: SubscriptionState = now <= end ? row.status : now <= graceEnd ? 'GRACE' : 'EXPIRED';
  return {
    id: row.id, planId: row.plan_id, planCode: row.plan_code, planName: row.plan_name, status: row.status, state,
    writable: state !== 'EXPIRED', billingCycle: row.billing_cycle, periodStart: new Date(row.current_period_start),
    periodEnd: end, graceEnd, limits, limitOverrides: overrides,
  };
}

export async function currentSubscription(db: Db | pg.Pool, tenantId: string, lock = false): Promise<SubscriptionView> {
  const { rows: [row] } = await db.query<SubRow>(`${SUB_SQL}${lock ? ' FOR UPDATE OF s' : ''}`, [tenantId]);
  return viewOf(row);
}

/** First day of the current month in Riyadh, as YYYY-MM-01. */
export const monthStart = () => `${riyadhNow().date.slice(0, 7)}-01`;

/** What each limit counts, for one tenant. */
export const USAGE_SQL: Record<Exclude<LimitKey, 'max_storage_mb'>, string> = {
  max_users: `SELECT count(*) FROM user_tenants WHERE tenant_id = $1 AND status = 'ACTIVE'`,
  max_companies: `SELECT count(*) FROM companies WHERE tenant_id = $1 AND deleted_at IS NULL`,
  max_branches: `SELECT count(*) FROM branches WHERE tenant_id = $1 AND deleted_at IS NULL`,
  max_warehouses: `SELECT count(*) FROM warehouses WHERE tenant_id = $1 AND deleted_at IS NULL`,
  max_products: `SELECT count(*) FROM products WHERE tenant_id = $1 AND deleted_at IS NULL`,
  max_invoices_per_month: `SELECT count(*) FROM sales_invoices WHERE tenant_id = $1 AND status <> 'DRAFT'
                             AND issued_at >= ($2::date::timestamp AT TIME ZONE 'Asia/Riyadh')`,
  max_api_calls_per_month: `SELECT COALESCE(sum(quantity), 0) FROM usage_records WHERE tenant_id = $1 AND metric = 'API_CALLS' AND period = $2::date`,
};

export async function usageOf(db: Db, tenantId: string): Promise<Record<LimitKey, number>> {
  const month = monthStart();
  const out = { max_storage_mb: 0 } as Record<LimitKey, number>;
  for (const [k, sql] of Object.entries(USAGE_SQL)) {
    const { rows: [r] } = await db.query<{ count?: string; coalesce?: string }>(sql, sql.includes('$2') ? [tenantId, month] : [tenantId]);
    out[k as LimitKey] = Number(Object.values(r!)[0]);
  }
  out.max_api_calls_per_month += pendingApiCalls(tenantId);
  return out;
}

const LABEL_AR: Record<LimitKey, string> = {
  max_users: 'المستخدمين', max_companies: 'الشركات', max_branches: 'الفروع', max_warehouses: 'المستودعات', max_products: 'المنتجات',
  max_invoices_per_month: 'الفواتير الشهرية', max_storage_mb: 'التخزين', max_api_calls_per_month: 'طلبات API الشهرية',
};

/**
 * Throws 402 PLAN_LIMIT_REACHED when creating one more item would exceed the
 * plan. The subscription row is locked, so concurrent creations in the same
 * tenant are serialised and cannot overshoot the limit together.
 */
export async function assertWithinLimit(db: Db, tenantId: string, key: Exclude<LimitKey, 'max_storage_mb' | 'max_api_calls_per_month'>) {
  const sub = await currentSubscription(db, tenantId, true);
  const limit = sub.limits[key];
  if (limit === null) return;
  const sql = USAGE_SQL[key];
  const { rows: [r] } = await db.query(sql, sql.includes('$2') ? [tenantId, monthStart()] : [tenantId]);
  const used = Number(Object.values(r!)[0]);
  if (used >= limit) {
    throw new AppError(402, 'PLAN_LIMIT_REACHED', `بلغت الحد الأقصى لـ${LABEL_AR[key]} في باقتك (${limit}). رقِّ الباقة لإضافة المزيد.`, { limit: key, max: limit, used });
  }
}

// --- Starting and changing subscriptions -----------------------------------------------

export function addPeriod(from: Date, cycle: 'MONTHLY' | 'YEARLY', count = 1): Date {
  const d = new Date(from);
  d.setUTCMonth(d.getUTCMonth() + (cycle === 'YEARLY' ? 12 : 1) * count);
  return d;
}

async function planRow(db: Db, planId: string | null) {
  const { rows: [p] } = await db.query<{ id: string; code: string; name_ar: string; price_monthly: string; price_yearly: string; currency: string; trial_days: number; is_active: boolean }>(
    `SELECT id, code, name_ar, price_monthly::text, price_yearly::text, currency, trial_days, is_active FROM plans WHERE ${planId ? 'id = $1' : 'is_default'}`,
    planId ? [planId] : []);
  if (!p) throw planId ? notFound('Plan') : new AppError(500, 'NO_DEFAULT_PLAN', 'No default plan is configured');
  if (!p.is_active) throw badRequest('PLAN_INACTIVE', 'This plan is no longer offered');
  return p;
}

export async function replaceItems(db: Db, tenantId: string, subscriptionId: string, plan: Awaited<ReturnType<typeof planRow>>, cycle: 'MONTHLY' | 'YEARLY') {
  await db.query(`DELETE FROM subscription_items WHERE subscription_id = $1`, [subscriptionId]);
  await db.query(
    `INSERT INTO subscription_items (tenant_id, subscription_id, item_type, code, description, unit_price, currency) VALUES ($1, $2, 'PLAN', $3, $4, $5, $6)`,
    [tenantId, subscriptionId, plan.code, `${plan.name_ar} — ${cycle === 'YEARLY' ? 'سنوي' : 'شهري'}`, cycle === 'YEARLY' ? plan.price_yearly : plan.price_monthly, plan.currency]);
}

/** New tenant: the given plan (or the default one), as a trial when the plan has trial days. */
export async function startSubscription(db: Db, tenantId: string, userId: string | null, planId: string | null = null) {
  const plan = await planRow(db, planId);
  const now = new Date();
  const trial = plan.trial_days > 0;
  const end = trial ? new Date(now.getTime() + plan.trial_days * 86_400_000) : addPeriod(now, 'MONTHLY');
  const { rows: [s] } = await db.query<{ id: string }>(
    `INSERT INTO subscriptions (tenant_id, plan_id, status, current_period_start, current_period_end) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
    [tenantId, plan.id, trial ? 'TRIALING' : 'ACTIVE', now, end]);
  await replaceItems(db, tenantId, s!.id, plan, 'MONTHLY');
  await billingEvent(db, { tenantId, subscriptionId: s!.id, type: trial ? 'TRIAL_STARTED' : 'PLAN_CHANGED', userId, details: { plan: plan.code, periodEnd: end } });
  return s!.id;
}

export async function changePlan(db: Db, tenantId: string, userId: string, planId: string, cycle: 'MONTHLY' | 'YEARLY') {
  const sub = await currentSubscription(db, tenantId, true);
  if (!sub.id) throw notFound('Subscription');
  const plan = await planRow(db, planId);
  await db.query(`UPDATE subscriptions SET plan_id = $2, billing_cycle = $3 WHERE id = $1`, [sub.id, plan.id, cycle]);
  await replaceItems(db, tenantId, sub.id, plan, cycle);
  await billingEvent(db, { tenantId, subscriptionId: sub.id, type: 'PLAN_CHANGED', userId, details: { from: sub.planCode, to: plan.code, billingCycle: cycle } });
  return currentSubscription(db, tenantId);
}

export type BillingEventType = 'TRIAL_STARTED' | 'PLAN_CHANGED' | 'PERIOD_EXTENDED' | 'PAYMENT_RECORDED' | 'LIMITS_CHANGED'
  | 'PLAN_CHANGE_REQUESTED' | 'SUBSCRIPTION_CANCELLED' | 'TENANT_SUSPENDED' | 'TENANT_ACTIVATED';

export async function billingEvent(db: Db, e: { tenantId: string; subscriptionId?: string | null; type: BillingEventType; userId: string | null; amount?: string | null; currency?: string | null; reference?: string | null; details?: unknown }) {
  await db.query(
    `INSERT INTO billing_events (tenant_id, subscription_id, event_type, amount, currency, reference, details, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
    [e.tenantId, e.subscriptionId ?? null, e.type, e.amount ?? null, e.currency ?? null, e.reference ?? null, JSON.stringify(e.details ?? {}), e.userId]);
}

// --- API call metering -------------------------------------------------------------------

const pending = new Map<string, number>();
export const pendingApiCalls = (tenantId: string) => pending.get(tenantId) ?? 0;
export const countApiCall = (tenantId: string) => pending.set(tenantId, pendingApiCalls(tenantId) + 1);

/** Writes the in-memory API call counters (one statement per tenant). */
export async function flushUsage(pool: pg.Pool) {
  if (!pending.size) return;
  const batch = [...pending.entries()];
  pending.clear();
  const month = monthStart();
  for (const [tenantId, n] of batch) {
    try {
      await pool.query(`SELECT usage_increment($1, 'API_CALLS', $2, $3)`, [tenantId, month, n]);
    } catch {
      pending.set(tenantId, pendingApiCalls(tenantId) + n); // try again on the next flush
    }
  }
}

/** Effective feature flags for a tenant: the per-tenant override, else the global default. */
export async function featuresOf(db: Db, tenantId: string): Promise<Record<string, boolean>> {
  const { rows } = await db.query<{ key: string; enabled: boolean }>(
    `SELECT f.key, COALESCE(t.enabled, f.enabled) AS enabled FROM feature_flags f
       LEFT JOIN tenant_feature_flags t ON t.flag_key = f.key AND t.tenant_id = $1`, [tenantId]);
  return Object.fromEntries(rows.map((r) => [r.key, r.enabled]));
}
