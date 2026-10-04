import { readdirSync } from 'node:fs';
import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { withTx } from '../src/db/tx.js';
import { addMember, client, login, register, setupApp, type Session, type TestContext } from './helpers.js';

const latestMigration = readdirSync(new URL('../src/db/migrations/', import.meta.url)).filter((f) => f.endsWith('.sql')).sort().at(-1);
let t: TestContext;
let adminSession: Session;
let admin: ReturnType<typeof client>;
let tenant: Session;

beforeAll(async () => {
  t = await setupApp();
  adminSession = await register(t.app, { tenantName: 'مشغل المنصة' });
  await t.ownerPool.query(`UPDATE users SET is_platform_admin = true WHERE id = $1`, [adminSession.userId]);
  admin = client(t.app, adminSession.token);
  tenant = await register(t.app, { tenantName: 'شركة العميل الأولى' });
});
afterAll(async () => { await t.close(); });

const planBody = (over: Record<string, unknown> = {}) => ({
  code: `PRO_${Math.floor(Math.random() * 1e6)}`, nameAr: 'الاحترافية', priceMonthly: '299.00', priceYearly: '2990.00', trialDays: 7, graceDays: 5,
  max_users: 10, max_companies: 3, max_branches: 10, max_warehouses: 10, max_products: null, max_invoices_per_month: null, max_storage_mb: 5120,
  max_api_calls_per_month: null, isPublic: true, isActive: true, isDefault: false, sortOrder: 20, ...over,
});

describe('access', () => {
  it('is for platform administrators only', async () => {
    const owner = client(t.app, tenant.token);
    for (const path of ['/api/admin/overview', '/api/admin/tenants', '/api/admin/users', '/api/admin/plans', '/api/admin/errors', '/api/admin/health']) {
      expect((await owner.get(path)).statusCode).toBe(403);
    }
    expect((await owner.post('/api/admin/plans', planBody())).statusCode).toBe(403);
  });

  it('keeps the database closed to tenants: plans and platform tables', async () => {
    const asTenant = <T>(fn: Parameters<typeof withTx<T>>[2]) => withTx(t.pool, { tenantId: tenant.tenantId, userId: tenant.userId }, fn);
    await expect(asTenant((db) => db.query(`INSERT INTO plans (code, name_ar) VALUES ('HACK', 'اختراق')`))).rejects.toThrow(/row-level security/);
    expect((await asTenant((db) => db.query(`UPDATE plans SET price_monthly = 0`))).rowCount).toBe(0);
    await asTenant((db) => db.query(`INSERT INTO system_errors (message) VALUES ('visible only to admins')`));
    expect((await asTenant((db) => db.query(`SELECT * FROM system_errors`))).rowCount).toBe(0);
    // Without the platform flag a tenant still sees only itself.
    expect((await asTenant((db) => db.query(`SELECT * FROM subscriptions`))).rows.every((r) => r.tenant_id === tenant.tenantId)).toBe(true);
  });

  it('stops a tenant from changing its own subscription, status or platform access in SQL', async () => {
    const asTenant = <T>(fn: Parameters<typeof withTx<T>>[2]) => withTx(t.pool, { tenantId: tenant.tenantId, userId: tenant.userId }, fn);
    expect((await asTenant((db) => db.query(`UPDATE subscriptions SET current_period_end = now() + interval '10 years' WHERE tenant_id = $1`, [tenant.tenantId]))).rowCount).toBe(0);
    await expect(asTenant((db) => db.query(`INSERT INTO subscriptions (tenant_id, plan_id, status, current_period_start, current_period_end)
      SELECT $1, id, 'ACTIVE', now(), now() + interval '10 years' FROM plans LIMIT 1`, [tenant.tenantId]))).rejects.toThrow(/row-level security/);
    await expect(asTenant((db) => db.query(`INSERT INTO billing_events (tenant_id, event_type) VALUES ($1, 'PAYMENT_RECORDED')`, [tenant.tenantId]))).rejects.toThrow(/row-level security/);
    await expect(asTenant((db) => db.query(`INSERT INTO tenant_feature_flags (tenant_id, flag_key, enabled) VALUES ($1, 'zatca_einvoicing', true)`, [tenant.tenantId]))).rejects.toThrow(/row-level security/);
    await expect(asTenant((db) => db.query(`UPDATE tenants SET status = 'SUSPENDED' WHERE id = $1`, [tenant.tenantId]))).rejects.toThrow(/platform administrator/);
    await expect(asTenant((db) => db.query(`UPDATE users SET is_platform_admin = true WHERE id = $1`, [tenant.userId]))).rejects.toThrow(/platform administrator/);
    // Ordinary profile edits still work.
    expect((await asTenant((db) => db.query(`UPDATE tenants SET name = name WHERE id = $1`, [tenant.tenantId]))).rowCount).toBe(1);
  });
});

describe('organizations', () => {
  it('lists and searches every organization with its subscription', async () => {
    const list = (await admin.get('/api/admin/tenants?search=العميل الأولى')).json();
    expect(list.data).toEqual([expect.objectContaining({ id: tenant.tenantId, status: 'ACTIVE', subscriptionState: 'TRIALING', ownerEmail: tenant.email, users: 1, companies: 1 })]);
    const byEmail = (await admin.get(`/api/admin/tenants?search=${encodeURIComponent(tenant.email)}`)).json();
    expect(byEmail.total).toBe(1);
    const detail = (await admin.get(`/api/admin/tenants/${tenant.tenantId}`)).json();
    expect(detail).toMatchObject({ name: 'شركة العميل الأولى', subscription: { planCode: 'TRIAL', state: 'TRIALING' }, usage: { max_users: 1 } });
    expect(detail.members).toEqual([expect.objectContaining({ email: tenant.email, isOwner: true })]);
    expect(detail.billingEvents.map((e: { eventType: string }) => e.eventType)).toEqual(['TRIAL_STARTED']);
  });

  it('creates an organization whose owner must change the temporary password', async () => {
    const plan = (await admin.post('/api/admin/plans', planBody())).json();
    const email = `new-owner-${Date.now()}@example.test`;
    const res = await admin.post('/api/admin/tenants', { tenantName: 'مؤسسة جديدة', companyName: 'مؤسسة جديدة للتجارة', ownerName: 'سعد', ownerEmail: email, planId: plan.id });
    expect(res.statusCode).toBe(201);
    const { tenantId, temporaryPassword } = res.json();
    const owner = await login(t.app, email, temporaryPassword);
    const me = (await client(t.app, owner.token).get('/api/auth/me')).json();
    expect(me).toMatchObject({ user: { mustChangePassword: true }, tenant: { id: tenantId, isOwner: true }, subscription: { state: 'TRIALING', planName: 'الاحترافية' } });
    expect((await admin.post('/api/admin/tenants', { tenantName: 'مكرر', companyName: 'مكرر', ownerName: 'سعد', ownerEmail: email })).json().error.code).toBe('EMAIL_TAKEN');
  });

  it('suspends and reactivates an organization', async () => {
    const s = await register(t.app);
    const member = await addMember(t.app, s, 'ACCOUNTANT');
    expect((await admin.post(`/api/admin/tenants/${s.tenantId}/suspend`, { reason: 'عدم السداد' })).json()).toMatchObject({ status: 'SUSPENDED' });
    const blocked = await client(t.app, member.token).get('/api/companies');
    expect(blocked.statusCode).toBe(403);
    expect(blocked.json().error.code).toBe('TENANT_SUSPENDED');
    await admin.post(`/api/admin/tenants/${s.tenantId}/activate`, { reason: 'تم السداد' });
    expect((await client(t.app, member.token).get('/api/companies')).statusCode).toBe(200);
    // The organization sees what happened in its own audit log.
    const audit = (await client(t.app, s.token).get('/api/audit-logs?limit=50')).json().data;
    expect(audit.filter((a: { action: string; entityType: string }) => a.entityType === 'tenant' && a.action === 'SETTINGS_CHANGE')).toHaveLength(2);
    expect((await admin.post(`/api/admin/tenants/${adminSession.tenantId}/suspend`, { reason: 'خطأ' })).json().error.code).toBe('OWN_TENANT');
  });

  it('changes plan, limits and period, and records every step', async () => {
    const s = await register(t.app);
    const plan = (await admin.post('/api/admin/plans', planBody({ priceMonthly: '150.00', priceYearly: '1500.00' }))).json();
    const yearly = (await admin.post(`/api/admin/tenants/${s.tenantId}/subscription/plan`, { planId: plan.id, billingCycle: 'YEARLY' })).json();
    expect(yearly).toMatchObject({ planCode: plan.code, billingCycle: 'YEARLY', limits: { max_users: 10, max_products: null } });
    const limited = (await admin.post(`/api/admin/tenants/${s.tenantId}/subscription/limits`, { overrides: { max_users: 25, max_products: 50 } })).json();
    expect(limited.limits).toMatchObject({ max_users: 25, max_products: 50, max_companies: 3 });
    expect((await admin.post(`/api/admin/tenants/${s.tenantId}/subscription/limits`, { overrides: { max_unknown: 1 } })).statusCode).toBe(400);
    const before = new Date(limited.periodEnd).getTime();
    const extended = (await admin.post(`/api/admin/tenants/${s.tenantId}/subscription/extend`, { days: 10, reason: 'تمديد التجربة' })).json();
    expect(new Date(extended.periodEnd).getTime() - before).toBe(10 * 86_400_000);
    const paid = (await admin.post(`/api/admin/tenants/${s.tenantId}/subscription/payment`, { amount: '1500.00', reference: 'INV-2026-001' })).json();
    expect(paid).toMatchObject({ status: 'ACTIVE', state: 'ACTIVE' });
    expect(new Date(paid.periodEnd).getTime()).toBeGreaterThan(Date.now() + 360 * 86_400_000); // one yearly cycle
    const detail = (await admin.get(`/api/admin/tenants/${s.tenantId}`)).json();
    expect(detail.subscription.items).toEqual([expect.objectContaining({ unitPrice: '1500.00', code: plan.code })]);
    expect(detail.billingEvents.map((e: { eventType: string }) => e.eventType)).toEqual(['PAYMENT_RECORDED', 'PERIOD_EXTENDED', 'LIMITS_CHANGED', 'PLAN_CHANGED', 'TRIAL_STARTED']);

    const cancelled = (await admin.post(`/api/admin/tenants/${s.tenantId}/subscription/cancel`, { reason: 'طلب العميل' })).json();
    expect(cancelled).toMatchObject({ state: 'NONE', writable: false });
    expect((await client(t.app, s.token).post('/api/products', { nameAr: 'بعد الإلغاء' })).json().error.code).toBe('SUBSCRIPTION_INACTIVE');
    const restarted = (await admin.post(`/api/admin/tenants/${s.tenantId}/subscription/plan`, { planId: plan.id })).json();
    expect(restarted).toMatchObject({ state: 'TRIALING', writable: true });
  });
});

describe('plans', () => {
  it('keeps exactly one active default plan', async () => {
    const plans = (await admin.get('/api/admin/plans')).json().data;
    const trial = plans.find((p: { code: string }) => p.code === 'TRIAL');
    expect(trial).toMatchObject({ isDefault: true, priceMonthly: '0.00' });
    expect((await admin.patch(`/api/admin/plans/${trial.id}`, { isActive: false })).json().error.code).toBe('DEFAULT_PLAN');
    expect((await admin.patch(`/api/admin/plans/${trial.id}`, { isDefault: false })).json().error.code).toBe('DEFAULT_PLAN');
    const next = (await admin.post('/api/admin/plans', planBody({ isDefault: true }))).json();
    const after = (await admin.get('/api/admin/plans')).json().data;
    expect(after.filter((p: { isDefault: boolean }) => p.isDefault).map((p: { id: string }) => p.id)).toEqual([next.id]);
    await admin.patch(`/api/admin/plans/${trial.id}`, { isDefault: true });
    expect((await admin.post('/api/admin/plans', planBody({ priceMonthly: '-1' }))).statusCode).toBe(400);
    expect((await admin.post('/api/admin/plans', planBody({ code: next.code }))).statusCode).toBe(409);
    // Only public, active plans are offered to organizations.
    await admin.patch(`/api/admin/plans/${next.id}`, { isPublic: false });
    const offered = (await client(t.app, tenant.token).get('/api/subscription/plans')).json().data.map((p: { id: string }) => p.id);
    expect(offered).not.toContain(next.id);
  });
});

describe('users', () => {
  it('disables, unlocks and grants platform access', async () => {
    const s = await register(t.app);
    const u = (await admin.get(`/api/admin/users?search=${encodeURIComponent(s.email)}`)).json();
    expect(u.data[0]).toMatchObject({ email: s.email, status: 'ACTIVE', isPlatformAdmin: false, tenants: [expect.objectContaining({ isOwner: true })] });
    await admin.post(`/api/admin/users/${s.userId}/status`, { status: 'DISABLED' });
    expect((await client(t.app, s.token).get('/api/companies')).statusCode).toBe(401);
    await admin.post(`/api/admin/users/${s.userId}/status`, { status: 'ACTIVE' });
    expect((await client(t.app, s.token).get('/api/companies')).statusCode).toBe(200);
    await t.ownerPool.query(`INSERT INTO login_failures (user_id, ip, failures, locked_until) VALUES ($1, '127.0.0.1', 9, now() + interval '1 hour')`, [s.userId]);
    await expect(login(t.app, s.email)).rejects.toThrow();
    expect((await admin.get(`/api/admin/users?search=${encodeURIComponent(s.email)}`)).json().data[0]).toMatchObject({ failedLogins: 9, lockedUntil: expect.any(String) });
    await admin.post(`/api/admin/users/${s.userId}/unlock`);
    expect((await login(t.app, s.email)).token).toBeTruthy();
    await admin.post(`/api/admin/users/${s.userId}/platform-admin`, { grant: true });
    expect((await client(t.app, s.token).get('/api/admin/overview')).statusCode).toBe(200);
    await admin.post(`/api/admin/users/${s.userId}/platform-admin`, { grant: false });
    expect((await client(t.app, s.token).get('/api/admin/overview')).statusCode).toBe(403);
    expect((await admin.post(`/api/admin/users/${adminSession.userId}/status`, { status: 'DISABLED' })).json().error.code).toBe('SELF');
    expect((await admin.post(`/api/admin/users/${adminSession.userId}/platform-admin`, { grant: false })).json().error.code).toBe('SELF');
  });
});

describe('feature flags', () => {
  it('combines the global default with per-organization overrides', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    expect((await admin.post('/api/admin/feature-flags', { key: 'beta_dashboard', nameAr: 'لوحة تجريبية', enabled: false })).statusCode).toBe(201);
    expect((await api.get('/api/auth/me')).json().features).toMatchObject({ beta_dashboard: false, zatca_einvoicing: true });
    await admin.put(`/api/admin/tenants/${s.tenantId}/features/beta_dashboard`, { enabled: true });
    await admin.put(`/api/admin/tenants/${s.tenantId}/features/zatca_einvoicing`, { enabled: false });
    expect((await api.get('/api/auth/me')).json().features).toMatchObject({ beta_dashboard: true, zatca_einvoicing: false });
    const blocked = await api.post('/api/zatca/devices', { name: 'وحدة', environment: 'DEVELOPER', businessCategory: 'تجارة' });
    expect(blocked.statusCode).toBe(403);
    await admin.put(`/api/admin/tenants/${s.tenantId}/features/zatca_einvoicing`, { enabled: null });
    expect((await api.get('/api/auth/me')).json().features.zatca_einvoicing).toBe(true);
    const flags = (await admin.get('/api/admin/feature-flags')).json().data;
    expect(flags.find((f: { key: string }) => f.key === 'beta_dashboard')).toMatchObject({ overrides: 1 });
  });
});

describe('monitoring', () => {
  it('reports revenue, usage, logs, errors and health', async () => {
    const plan = (await admin.post('/api/admin/plans', planBody({ priceMonthly: '299.00' }))).json();
    await admin.post(`/api/admin/tenants/${tenant.tenantId}/subscription/plan`, { planId: plan.id, billingCycle: 'MONTHLY' });
    await admin.post(`/api/admin/tenants/${tenant.tenantId}/subscription/payment`, { amount: '299.00', reference: 'TRF-7781' });
    const overview = (await admin.get('/api/admin/overview')).json();
    expect(overview.tenants).toBeGreaterThanOrEqual(5);
    expect(Number(overview.revenueThisMonth)).toBeGreaterThanOrEqual(1799);
    expect(Number(overview.mrr)).toBeGreaterThanOrEqual(299);
    expect(overview.subscriptions.ACTIVE).toBeGreaterThanOrEqual(1);
    const revenue = (await admin.get('/api/admin/revenue')).json();
    expect(revenue.payments[0]).toMatchObject({ amount: '299.00', reference: 'TRF-7781' });
    expect(revenue.payments.map((x: { reference: string }) => x.reference)).toContain('INV-2026-001');
    expect(revenue.monthly.at(-1).payments).toBeGreaterThanOrEqual(2);
    expect(revenue.byPlan.length).toBeGreaterThan(0);
    const usage = (await admin.get('/api/admin/usage')).json().data;
    expect(usage.find((u: { tenantId: string }) => u.tenantId === tenant.tenantId)).toMatchObject({ users: 1, companies: 1 });
    const subs = (await admin.get('/api/admin/subscriptions?state=ACTIVE')).json().data;
    expect(subs.every((x: { state: string }) => x.state === 'ACTIVE')).toBe(true);
    const plog = (await admin.get('/api/admin/platform-logs')).json().data;
    expect(plog.map((l: { action: string }) => l.action)).toEqual(expect.arrayContaining(['SUSPEND', 'ACTIVATE', 'PAYMENT', 'PLAN_CHANGE', 'LIMITS_CHANGE', 'USER_STATUS']));
    const audit = (await admin.get(`/api/admin/audit-logs?tenantId=${tenant.tenantId}`)).json().data;
    expect(audit.length).toBeGreaterThan(0);
    expect(audit.every((a: { tenantId: string }) => a.tenantId === tenant.tenantId)).toBe(true);
    const errors = (await admin.get('/api/admin/errors')).json().data;
    expect(errors.some((e: { message: string }) => e.message === 'visible only to admins')).toBe(true);
    const health = (await admin.get('/api/admin/health')).json();
    expect(health).toMatchObject({ status: 'ok', database: { migrations: expect.any(Number), lastMigration: latestMigration }, zatca: { pending: expect.any(Number) } });
  });
});
