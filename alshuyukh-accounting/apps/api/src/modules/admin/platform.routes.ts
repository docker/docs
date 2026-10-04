import { readFileSync } from 'node:fs';
import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { badRequest, notFound } from '../../lib/errors.js';
import { amountString } from '../../lib/money.js';
import { optionalText, parse, uuidParam } from '../../lib/validation.js';
import { PLAN_PUBLIC_SELECT } from '../subscriptions/routes.js';
import { LIMIT_KEYS, monthStart } from '../subscriptions/service.js';
import { platformAudit, platformTx, requirePlatformAdmin, STATE_SQL } from './guard.js';

const started = Date.now();
const version = (() => { try { return JSON.parse(readFileSync(new URL('../../../package.json', import.meta.url), 'utf8')).version as string; } catch { return 'unknown'; } })();

/** Effective limit of a subscription row in SQL: an override key replaces the plan value (null = unlimited). */
const limitSql = (k: string) => `CASE WHEN s.limit_overrides ? '${k}' THEN (s.limit_overrides->>'${k}')::int ELSE p.${k} END`;
/** Monthly value of a subscription: its items, yearly prices spread over 12 months. */
const MONTHLY_VALUE = `COALESCE((SELECT sum(i.unit_price * i.quantity) FROM subscription_items i WHERE i.subscription_id = s.id), 0)
  / CASE WHEN s.billing_cycle = 'YEARLY' THEN 12 ELSE 1 END`;

const planFields = {
  code: z.string().trim().toUpperCase().regex(/^[A-Z][A-Z0-9_]{1,30}$/),
  nameAr: z.string().trim().min(2).max(100), nameEn: optionalText(100), description: optionalText(1000),
  currency: z.string().regex(/^[A-Z]{3}$/).default('SAR'),
  priceMonthly: amountString, priceYearly: amountString,
  trialDays: z.number().int().min(0).max(365), graceDays: z.number().int().min(0).max(90),
  ...Object.fromEntries(LIMIT_KEYS.map((k) => [k, z.number().int().positive().nullable()])) as Record<(typeof LIMIT_KEYS)[number], z.ZodNullable<z.ZodNumber>>,
  isPublic: z.boolean(), isActive: z.boolean(), isDefault: z.boolean(), sortOrder: z.number().int().min(0).max(1000),
};
const PLAN_COLUMNS: Record<string, string> = {
  code: 'code', nameAr: 'name_ar', nameEn: 'name_en', description: 'description', currency: 'currency', priceMonthly: 'price_monthly',
  priceYearly: 'price_yearly', trialDays: 'trial_days', graceDays: 'grace_days', isPublic: 'is_public', isActive: 'is_active',
  isDefault: 'is_default', sortOrder: 'sort_order', ...Object.fromEntries(LIMIT_KEYS.map((k) => [k, k])),
};

export default async function adminPlatformRoutes(app: FastifyInstance) {
  const guard = requirePlatformAdmin(app);

  app.get('/admin/overview', { preHandler: guard }, async (req) => platformTx(app, req, async (db) => {
    const { rows: [o] } = await db.query<Record<string, string>>(
      `SELECT (SELECT count(*) FROM tenants WHERE deleted_at IS NULL) AS tenants,
              (SELECT count(*) FROM tenants WHERE status = 'SUSPENDED' AND deleted_at IS NULL) AS suspended,
              (SELECT count(*) FROM users WHERE deleted_at IS NULL) AS users,
              (SELECT count(*) FROM tenants WHERE created_at >= now() - interval '30 days') AS "newTenants30d",
              (SELECT count(*) FROM system_errors WHERE created_at >= now() - interval '24 hours') AS "errors24h",
              (SELECT COALESCE(sum(amount), 0) FROM billing_events WHERE event_type = 'PAYMENT_RECORDED' AND created_at >= ($1::date::timestamp AT TIME ZONE 'Asia/Riyadh')) AS "revenueThisMonth"`,
      [monthStart()]);
    const { rows: states } = await db.query<{ state: string; count: string; mrr: string }>(
      `SELECT ${STATE_SQL} AS state, count(*) AS count, COALESCE(sum(${MONTHLY_VALUE}), 0)::numeric(18,2)::text AS mrr
         FROM subscriptions s JOIN plans p ON p.id = s.plan_id WHERE s.status <> 'CANCELLED' GROUP BY 1`);
    const mrr = states.filter((s) => s.state === 'ACTIVE' || s.state === 'GRACE').reduce((a, s) => a + Number(s.mrr), 0);
    return {
      ...Object.fromEntries(Object.entries(o!).map(([k, v]) => [k, k === 'revenueThisMonth' ? Number(v).toFixed(2) : Number(v)])),
      subscriptions: Object.fromEntries(states.map((s) => [s.state, Number(s.count)])), mrr: mrr.toFixed(2),
    };
  }));

  // Users --------------------------------------------------------------------------------
  app.get('/admin/users', { preHandler: guard }, async (req) => {
    const q = parse(z.object({ search: z.string().trim().max(100).optional(), limit: z.coerce.number().int().min(1).max(200).default(50), offset: z.coerce.number().int().min(0).default(0) }), req.query);
    return platformTx(app, req, async (db) => {
      const { rows } = await db.query<{ total: string }>(
        `SELECT count(*) OVER () AS total, u.id, u.email, u.full_name AS "fullName", u.status, u.is_platform_admin AS "isPlatformAdmin",
                u.locked_until AS "lockedUntil", u.failed_login_attempts AS "failedLogins", u.last_login_at AS "lastLoginAt", u.created_at AS "createdAt",
                COALESCE((SELECT json_agg(json_build_object('tenantId', t.id, 'name', t.name, 'isOwner', ut.is_owner) ORDER BY t.name)
                            FROM user_tenants ut JOIN tenants t ON t.id = ut.tenant_id WHERE ut.user_id = u.id), '[]') AS tenants
           FROM users u WHERE u.deleted_at IS NULL AND ($1::text IS NULL OR u.email ILIKE '%' || $1 || '%' OR u.full_name ILIKE '%' || $1 || '%')
          ORDER BY u.created_at DESC LIMIT $2 OFFSET $3`, [q.search || null, q.limit, q.offset]);
      return { data: rows.map(({ total: _t, ...r }) => r), total: Number(rows[0]?.total ?? 0) };
    });
  });

  const userAction = (path: string, body: z.ZodTypeAny, apply: (db: any, id: string, b: any, self: boolean) => Promise<Record<string, unknown>>) =>
    app.post(`/admin/users/:id/${path}`, { preHandler: guard }, async (req) => {
      const { id } = parse(uuidParam, req.params);
      const b = parse(body, req.body ?? {});
      return platformTx(app, req, async (db) => {
        const { rows: [u] } = await db.query(`SELECT id, status, is_platform_admin FROM users WHERE id = $1 AND deleted_at IS NULL`, [id]);
        if (!u) throw notFound('User');
        const changes = await apply(db, id, b, id === req.auth!.userId);
        await platformAudit(db, req, { action: `USER_${path.toUpperCase().replace('-', '_')}`, entityType: 'user', entityId: id, oldValues: u, newValues: changes });
        return { id, ...changes };
      });
    });
  userAction('status', z.object({ status: z.enum(['ACTIVE', 'DISABLED']) }), async (db, id, b, self) => {
    if (self && b.status === 'DISABLED') throw badRequest('SELF', 'You cannot disable your own account');
    await db.query(`UPDATE users SET status = $2 WHERE id = $1`, [id, b.status]);
    return { status: b.status };
  });
  userAction('unlock', z.object({}), async (db, id) => {
    await db.query(`UPDATE users SET failed_login_attempts = 0, locked_until = NULL WHERE id = $1`, [id]);
    return { unlocked: true };
  });
  userAction('platform-admin', z.object({ grant: z.boolean() }), async (db, id, b, self) => {
    if (self && !b.grant) throw badRequest('SELF', 'You cannot remove your own platform access');
    await db.query(`UPDATE users SET is_platform_admin = $2 WHERE id = $1`, [id, b.grant]);
    return { isPlatformAdmin: b.grant };
  });

  // Plans --------------------------------------------------------------------------------
  app.get('/admin/plans', { preHandler: guard }, async (req) => platformTx(app, req, async (db) => ({
    data: (await db.query(`${PLAN_PUBLIC_SELECT.replace('FROM plans', `, (SELECT count(*) FROM subscriptions s WHERE s.plan_id = plans.id AND s.status <> 'CANCELLED')::int AS subscribers FROM plans`)}
      ORDER BY sort_order, price_monthly`)).rows,
  })));

  const savePlan = async (req: Parameters<typeof platformTx>[1], id: string | null, body: Record<string, unknown>) => platformTx(app, req, async (db) => {
    const before = id ? (await db.query(`SELECT * FROM plans WHERE id = $1`, [id])).rows[0] : null;
    if (id && !before) throw notFound('Plan');
    const isDefault = body.isDefault ?? before?.is_default;
    if (isDefault && (body.isActive ?? before?.is_active) === false) throw badRequest('DEFAULT_PLAN', 'The default plan must stay active');
    if (before?.is_default && body.isDefault === false) throw badRequest('DEFAULT_PLAN', 'Choose another default plan instead');
    if (body.isDefault === true) await db.query(`UPDATE plans SET is_default = false WHERE is_default AND id IS DISTINCT FROM $1`, [id]);
    const entries = Object.entries(body).filter(([k, v]) => PLAN_COLUMNS[k] && v !== undefined);
    let planId = id;
    if (id) {
      if (entries.length) await db.query(`UPDATE plans SET ${entries.map(([k], i) => `${PLAN_COLUMNS[k]} = $${i + 2}`).join(', ')} WHERE id = $1`, [id, ...entries.map(([, v]) => v)]);
    } else {
      const { rows: [p] } = await db.query<{ id: string }>(
        `INSERT INTO plans (${entries.map(([k]) => PLAN_COLUMNS[k]).join(', ')}) VALUES (${entries.map((_, i) => `$${i + 1}`).join(', ')}) RETURNING id`,
        entries.map(([, v]) => v));
      planId = p!.id;
    }
    const after = (await db.query(`${PLAN_PUBLIC_SELECT} WHERE id = $1`, [planId])).rows[0];
    await platformAudit(db, req, { action: id ? 'UPDATE' : 'CREATE', entityType: 'plan', entityId: planId, oldValues: before ?? undefined, newValues: after });
    return after;
  });
  app.post('/admin/plans', { preHandler: guard }, async (req, reply) => {
    const body = parse(z.object(planFields), req.body);
    reply.code(201);
    return savePlan(req, null, body);
  });
  app.patch('/admin/plans/:id', { preHandler: guard }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    // Prices changed here apply to new subscriptions and plan changes; current ones keep their snapshot.
    return savePlan(req, id, parse(z.object(planFields).partial(), req.body));
  });

  // Subscriptions, revenue, usage ---------------------------------------------------------
  app.get('/admin/subscriptions', { preHandler: guard }, async (req) => {
    const q = parse(z.object({ state: z.enum(['TRIALING', 'ACTIVE', 'GRACE', 'EXPIRED']).optional() }), req.query);
    return platformTx(app, req, async (db) => ({
      data: (await db.query(
        `SELECT * FROM (SELECT s.id, s.tenant_id AS "tenantId", t.name AS "tenantName", t.status AS "tenantStatus", p.name_ar AS "planName",
                s.billing_cycle AS "billingCycle", ${STATE_SQL} AS state, s.current_period_end AS "periodEnd",
                (${MONTHLY_VALUE})::numeric(18,2)::text AS "monthlyValue"
           FROM subscriptions s JOIN plans p ON p.id = s.plan_id JOIN tenants t ON t.id = s.tenant_id
          WHERE s.status <> 'CANCELLED' AND t.deleted_at IS NULL) x
          WHERE ($1::text IS NULL OR state = $1) ORDER BY "periodEnd"`, [q.state ?? null])).rows,
    }));
  });

  app.get('/admin/revenue', { preHandler: guard }, async (req) => platformTx(app, req, async (db) => {
    const { rows: monthly } = await db.query(
      `SELECT to_char(created_at AT TIME ZONE 'Asia/Riyadh', 'YYYY-MM') AS month, sum(amount)::numeric(18,2)::text AS amount, count(*)::int AS payments
         FROM billing_events WHERE event_type = 'PAYMENT_RECORDED' AND created_at >= now() - interval '12 months' GROUP BY 1 ORDER BY 1`);
    const { rows: byPlan } = await db.query(
      `SELECT p.name_ar AS "planName", count(*)::int AS subscriptions, sum(${MONTHLY_VALUE})::numeric(18,2)::text AS mrr
         FROM subscriptions s JOIN plans p ON p.id = s.plan_id
        WHERE s.status <> 'CANCELLED' AND ${STATE_SQL} IN ('ACTIVE', 'GRACE') GROUP BY p.name_ar ORDER BY 3 DESC`);
    const { rows: payments } = await db.query(
      `SELECT b.id, b.tenant_id AS "tenantId", t.name AS "tenantName", b.amount::text, b.currency, b.reference, b.created_at AS "createdAt"
         FROM billing_events b JOIN tenants t ON t.id = b.tenant_id WHERE b.event_type = 'PAYMENT_RECORDED' ORDER BY b.created_at DESC LIMIT 50`);
    return { monthly, byPlan, payments };
  }));

  app.get('/admin/usage', { preHandler: guard }, async (req) => platformTx(app, req, async (db) => ({
    data: (await db.query(
      `SELECT t.id AS "tenantId", t.name AS "tenantName", p.name_ar AS "planName",
              (SELECT count(*) FROM user_tenants ut WHERE ut.tenant_id = t.id AND ut.status = 'ACTIVE')::int AS users, ${limitSql('max_users')} AS "maxUsers",
              (SELECT count(*) FROM companies c WHERE c.tenant_id = t.id AND c.deleted_at IS NULL)::int AS companies, ${limitSql('max_companies')} AS "maxCompanies",
              (SELECT count(*) FROM products x WHERE x.tenant_id = t.id AND x.deleted_at IS NULL)::int AS products, ${limitSql('max_products')} AS "maxProducts",
              (SELECT count(*) FROM sales_invoices i WHERE i.tenant_id = t.id AND i.status <> 'DRAFT' AND i.issued_at >= ($1::date::timestamp AT TIME ZONE 'Asia/Riyadh'))::int AS invoices,
              ${limitSql('max_invoices_per_month')} AS "maxInvoices",
              COALESCE((SELECT quantity FROM usage_records u WHERE u.tenant_id = t.id AND u.metric = 'API_CALLS' AND u.period = $1::date), 0)::int AS "apiCalls",
              ${limitSql('max_api_calls_per_month')} AS "maxApiCalls"
         FROM tenants t LEFT JOIN subscriptions s ON s.tenant_id = t.id AND s.status <> 'CANCELLED' LEFT JOIN plans p ON p.id = s.plan_id
        WHERE t.deleted_at IS NULL ORDER BY "apiCalls" DESC, t.name LIMIT 500`, [monthStart()])).rows,
  })));

  // Logs and errors ------------------------------------------------------------------------
  app.get('/admin/audit-logs', { preHandler: guard }, async (req) => {
    const q = parse(z.object({ tenantId: z.uuid().optional(), action: z.string().regex(/^[A-Z_]+$/).optional(), limit: z.coerce.number().int().min(1).max(200).default(100), offset: z.coerce.number().int().min(0).default(0) }), req.query);
    return platformTx(app, req, async (db) => ({
      data: (await db.query(
        `SELECT a.id, a.tenant_id AS "tenantId", t.name AS "tenantName", u.email AS "userEmail", a.action, a.entity_type AS "entityType",
                a.entity_id AS "entityId", a.ip_address AS "ipAddress", a.created_at AS "createdAt"
           FROM audit_logs a JOIN tenants t ON t.id = a.tenant_id LEFT JOIN users u ON u.id = a.user_id
          WHERE ($1::uuid IS NULL OR a.tenant_id = $1) AND ($2::text IS NULL OR a.action = $2)
          ORDER BY a.created_at DESC LIMIT $3 OFFSET $4`, [q.tenantId ?? null, q.action ?? null, q.limit, q.offset])).rows,
    }));
  });

  app.get('/admin/platform-logs', { preHandler: guard }, async (req) => platformTx(app, req, async (db) => ({
    data: (await db.query(
      `SELECT l.id, u.email AS "adminEmail", l.action, l.entity_type AS "entityType", l.entity_id AS "entityId", l.target_tenant_id AS "tenantId",
              t.name AS "tenantName", l.new_values AS "newValues", l.created_at AS "createdAt"
         FROM platform_audit_logs l JOIN users u ON u.id = l.admin_user_id LEFT JOIN tenants t ON t.id = l.target_tenant_id
        ORDER BY l.created_at DESC LIMIT 200`)).rows,
  })));

  app.get('/admin/errors', { preHandler: guard }, async (req) => platformTx(app, req, async (db) => ({
    data: (await db.query(
      `SELECT e.id, e.tenant_id AS "tenantId", t.name AS "tenantName", e.request_id AS "requestId", e.method, e.path, e.status_code AS "statusCode",
              e.error_code AS "errorCode", e.message, e.stack, e.created_at AS "createdAt"
         FROM system_errors e LEFT JOIN tenants t ON t.id = e.tenant_id ORDER BY e.created_at DESC LIMIT 200`)).rows,
  })));

  // Feature flags --------------------------------------------------------------------------
  app.get('/admin/feature-flags', { preHandler: guard }, async (req) => platformTx(app, req, async (db) => ({
    data: (await db.query(
      `SELECT f.key, f.name_ar AS "nameAr", f.description, f.enabled,
              (SELECT count(*) FROM tenant_feature_flags t WHERE t.flag_key = f.key)::int AS overrides
         FROM feature_flags f ORDER BY f.key`)).rows,
  })));
  app.post('/admin/feature-flags', { preHandler: guard }, async (req, reply) => {
    const body = parse(z.object({ key: z.string().regex(/^[a-z][a-z0-9_]{1,62}$/), nameAr: z.string().trim().min(2).max(100), description: optionalText(500), enabled: z.boolean().default(false) }), req.body);
    const flag = await platformTx(app, req, async (db) => {
      await db.query(`INSERT INTO feature_flags (key, name_ar, description, enabled) VALUES ($1, $2, $3, $4)`, [body.key, body.nameAr, body.description ?? null, body.enabled]);
      await platformAudit(db, req, { action: 'CREATE', entityType: 'feature_flag', entityId: body.key, newValues: body });
      return body;
    });
    reply.code(201);
    return flag;
  });
  app.patch('/admin/feature-flags/:key', { preHandler: guard }, async (req) => {
    const { key } = parse(z.object({ key: z.string().regex(/^[a-z][a-z0-9_]{1,62}$/) }), req.params);
    const body = parse(z.object({ nameAr: z.string().trim().min(2).max(100), description: optionalText(500), enabled: z.boolean() }).partial(), req.body);
    return platformTx(app, req, async (db) => {
      const { rows: [before] } = await db.query(`SELECT key, name_ar, description, enabled FROM feature_flags WHERE key = $1`, [key]);
      if (!before) throw notFound('Feature flag');
      await db.query(`UPDATE feature_flags SET name_ar = COALESCE($2, name_ar), description = COALESCE($3, description), enabled = COALESCE($4, enabled) WHERE key = $1`,
        [key, body.nameAr ?? null, body.description ?? null, body.enabled ?? null]);
      await platformAudit(db, req, { action: 'UPDATE', entityType: 'feature_flag', entityId: key, oldValues: before, newValues: body });
      return { key, ...body };
    });
  });

  // System health --------------------------------------------------------------------------
  app.get('/admin/health', { preHandler: guard }, async (req) => {
    const t0 = performance.now();
    await app.deps.pool.query('SELECT 1');
    const dbLatencyMs = Math.round((performance.now() - t0) * 10) / 10;
    return platformTx(app, req, async (db) => {
      const { rows: [d] } = await db.query<Record<string, string>>(
        `SELECT version() AS "postgres", pg_database_size(current_database()) AS "databaseBytes",
                (SELECT count(*) FROM schema_migrations) AS migrations, (SELECT max(version) FROM schema_migrations) AS "lastMigration",
                (SELECT count(*) FROM zatca_invoices WHERE status = 'PENDING') AS "zatcaPending",
                (SELECT count(*) FROM zatca_invoices WHERE status = 'PENDING' AND invoice_kind = 'SIMPLIFIED' AND created_at < now() - interval '24 hours') AS "zatcaOverdue",
                (SELECT count(*) FROM system_errors WHERE created_at >= now() - interval '24 hours') AS "errors24h"`);
      const pool = app.deps.pool;
      const mem = process.memoryUsage();
      return {
        status: 'ok', version, node: process.version, uptimeSeconds: Math.round((Date.now() - started) / 1000),
        memoryMb: { rss: Math.round(mem.rss / 1048576), heapUsed: Math.round(mem.heapUsed / 1048576) },
        database: { latencyMs: dbLatencyMs, postgres: d!.postgres!.split(' on ')[0], sizeMb: Math.round(Number(d!.databaseBytes) / 1048576),
          migrations: Number(d!.migrations), lastMigration: d!.lastMigration, pool: { total: pool.totalCount, idle: pool.idleCount, waiting: pool.waitingCount } },
        zatca: { pending: Number(d!.zatcaPending), overdue: Number(d!.zatcaOverdue), worker: app.deps.env.ZATCA_WORKER },
        errors24h: Number(d!.errors24h),
      };
    });
  });
}
