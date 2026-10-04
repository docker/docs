import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { flushUsage } from '../src/modules/subscriptions/service.js';
import { commerce } from './commerce-helpers.js';
import { addMember, client, register, roleId, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let admin: ReturnType<typeof client>;

/** A platform administrator (the flag is set directly; the API never lets a tenant grant it). */
async function makeAdmin(): Promise<Session> {
  const s = await register(t.app);
  await t.ownerPool.query(`UPDATE users SET is_platform_admin = true WHERE id = $1`, [s.userId]);
  return s;
}
const setLimits = async (tenantId: string, overrides: Record<string, number | null>) => {
  const r = await admin.post(`/api/admin/tenants/${tenantId}/subscription/limits`, { overrides });
  expect(r.statusCode).toBe(200);
};
/** Moves the subscription period into the past (the API only extends periods). */
const endPeriod = (s: Session, daysAgo: number) => t.ownerPool.query(`UPDATE subscriptions SET current_period_start = now() - interval '60 days', current_period_end = now() - make_interval(days => $2)
             WHERE tenant_id = $1 AND status <> 'CANCELLED'`, [s.tenantId, daysAgo]);

beforeAll(async () => {
  t = await setupApp();
  admin = client(t.app, (await makeAdmin()).token);
});
afterAll(async () => { await t.close(); });

describe('trial on sign-up', () => {
  it('starts every new organization on the default plan as a trial', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    const sub = (await api.get('/api/subscription')).json();
    expect(sub).toMatchObject({ planCode: 'TRIAL', status: 'TRIALING', state: 'TRIALING', writable: true, billingCycle: 'MONTHLY' });
    expect(new Date(sub.periodEnd).getTime() - Date.now()).toBeGreaterThan(13 * 86_400_000);
    expect(sub.items).toEqual([expect.objectContaining({ code: 'TRIAL', unitPrice: '0.00' })]);
    expect(sub.usage).toMatchObject({ max_users: 1, max_companies: 1, max_branches: 1, max_warehouses: 1 });
    const me = (await api.get('/api/auth/me')).json();
    expect(me.subscription).toMatchObject({ state: 'TRIALING', writable: true, planName: 'الباقة التجريبية' });
    expect(me.features).toMatchObject({ zatca_einvoicing: true });
    const events = (await api.get('/api/subscription/billing-events')).json().data;
    expect(events.map((e: { eventType: string }) => e.eventType)).toEqual(['TRIAL_STARTED']);
  });
});

describe('plan limits', () => {
  it('counts reactivated members and accepted invitations against the user limit', async () => {
    const s = await register(t.app);
    const other = await register(t.app);
    const api = client(t.app, s.token);
    const viewer = await roleId(t.app, s.token, 'VIEWER');
    const member = (await api.post('/api/users', { email: `m-${Date.now()}@example.test`, fullName: 'عضو', initialPassword: 'Str0ng-Passw0rd!', roleIds: [viewer] })).json();
    expect((await api.post('/api/users', { email: other.email, fullName: 'مدعو', roleIds: [viewer] })).statusCode).toBe(201);
    expect((await api.patch(`/api/users/${member.id}`, { status: 'DISABLED' })).statusCode).toBe(200);
    await setLimits(s.tenantId, { max_users: 1 });
    const reactivate = await api.patch(`/api/users/${member.id}`, { status: 'ACTIVE' });
    expect(reactivate.json().error).toMatchObject({ code: 'PLAN_LIMIT_REACHED', details: { limit: 'max_users' } });
    const accept = await client(t.app, other.token).post(`/api/auth/invitations/${s.tenantId}/accept`);
    expect(accept.json().error.code).toBe('PLAN_LIMIT_REACHED');
    await setLimits(s.tenantId, { max_users: null });
    expect((await client(t.app, other.token).post(`/api/auth/invitations/${s.tenantId}/accept`)).statusCode).toBe(204);
  });

  it('stops at each limit with 402 and leaves nothing half-created', async () => {
    const s = await register(t.app);
    const c = await commerce(t.app, s);
    const companyId = (await c.api.get('/api/companies')).json().data[0].id;
    await setLimits(s.tenantId, { max_users: 1, max_companies: 1, max_branches: 1, max_warehouses: 1, max_products: 2, max_invoices_per_month: 1 });

    const member = await c.api.post('/api/users', { email: `x-${Date.now()}@example.test`, fullName: 'عضو', initialPassword: 'Str0ng-Passw0rd!', roleIds: [await roleId(t.app, s.token, 'VIEWER')] });
    expect(member.statusCode).toBe(402);
    expect(member.json().error).toMatchObject({ code: 'PLAN_LIMIT_REACHED', details: { limit: 'max_users', max: 1, used: 1 } });
    expect((await c.api.post('/api/companies', { name: 'شركة ثانية' })).json().error.code).toBe('PLAN_LIMIT_REACHED');
    expect((await c.api.post(`/api/companies/${companyId}/branches`, { code: 'B2', name: 'فرع جدة' })).json().error.code).toBe('PLAN_LIMIT_REACHED');
    expect((await c.api.post(`/api/companies/${companyId}/warehouses`, { code: 'W2', name: 'مستودع جدة' })).json().error.code).toBe('PLAN_LIMIT_REACHED');

    await c.product();
    await c.product();
    expect((await c.api.post('/api/products', { nameAr: 'منتج ثالث' })).statusCode).toBe(402);

    const cust = (await c.customer()).id;
    const svc = (await c.api.get('/api/products')).json().data[0].id;
    await c.postInvoice(cust, [{ productId: svc, quantity: '1' }]);
    const second = await c.invoice(cust, [{ productId: svc, quantity: '1' }]);
    expect((await c.api.post(`/api/invoices/${second.id}/post`)).json().error.code).toBe('PLAN_LIMIT_REACHED');
    expect((await c.api.get(`/api/invoices/${second.id}`)).json()).toMatchObject({ status: 'DRAFT', journalEntryId: null });

    // Raising a limit (null = unlimited) applies immediately.
    await setLimits(s.tenantId, { max_products: null });
    expect((await c.api.post('/api/products', { nameAr: 'منتج ثالث' })).statusCode).toBe(201);
  });

  it('does not overshoot a limit under concurrent requests', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    await setLimits(s.tenantId, { max_products: 3 });
    const results = await Promise.all(Array.from({ length: 6 }, (_, i) => api.post('/api/products', { nameAr: `منتج متزامن ${i}` })));
    expect(results.filter((r) => r.statusCode === 201)).toHaveLength(3);
    expect(results.filter((r) => r.statusCode === 402)).toHaveLength(3);
  });

  it('meters API calls and refuses them past the monthly allowance', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    await flushUsage(t.pool);
    await setLimits(s.tenantId, { max_api_calls_per_month: 5 });
    const codes: number[] = [];
    for (let i = 0; i < 8; i++) codes.push((await api.get('/api/companies')).statusCode);
    expect(codes.filter((x) => x === 429).length).toBeGreaterThan(0);
    expect(codes.indexOf(429)).toBeLessThanOrEqual(5);
    // Sign-in and the subscription page stay reachable.
    expect((await api.get('/api/subscription')).statusCode).toBe(200);
    await flushUsage(t.pool);
    expect((await api.get('/api/subscription')).json().usage.max_api_calls_per_month).toBeGreaterThanOrEqual(5);
    await setLimits(s.tenantId, {});
    expect((await api.get('/api/companies')).statusCode).toBe(200);
  });
});

describe('expiry', () => {
  it('keeps full access during the grace days, then becomes read-only until payment', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    await endPeriod(s, 2); // within the 7 grace days
    expect((await api.get('/api/auth/me')).json().subscription).toMatchObject({ state: 'GRACE', writable: true });
    expect((await api.post('/api/products', { nameAr: 'أثناء المهلة' })).statusCode).toBe(201);

    await endPeriod(s, 30);
    expect((await api.get('/api/auth/me')).json().subscription).toMatchObject({ state: 'EXPIRED', writable: false });
    expect((await api.get('/api/products')).statusCode).toBe(200);
    const write = await api.post('/api/products', { nameAr: 'بعد الانتهاء' });
    expect(write.statusCode).toBe(402);
    expect(write.json().error.code).toBe('SUBSCRIPTION_INACTIVE');
    // The owner can still ask for a plan and sign out.
    const plans = (await admin.get('/api/admin/plans')).json().data;
    const paid = (await admin.post('/api/admin/plans', {
      code: `BASIC_${Date.now() % 100000}`, nameAr: 'الأساسية', priceMonthly: '99.00', priceYearly: '990.00', trialDays: 0, graceDays: 7,
      max_users: 5, max_companies: 1, max_branches: 3, max_warehouses: 3, max_products: 1000, max_invoices_per_month: 500, max_storage_mb: 1024,
      max_api_calls_per_month: null, isPublic: true, isActive: true, isDefault: false, sortOrder: 10,
    })).json();
    expect(plans.length).toBeGreaterThan(0);
    expect((await api.post('/api/subscription/request-change', { planId: paid.id, billingCycle: 'MONTHLY' })).json()).toEqual({ requested: true });

    // The administrator applies the plan and records the payment: access returns.
    await admin.post(`/api/admin/tenants/${s.tenantId}/subscription/plan`, { planId: paid.id, billingCycle: 'MONTHLY' });
    const paidSub = (await admin.post(`/api/admin/tenants/${s.tenantId}/subscription/payment`, { amount: '99.00', reference: 'TRF-1001' })).json();
    expect(paidSub).toMatchObject({ status: 'ACTIVE', state: 'ACTIVE', writable: true, planName: 'الأساسية' });
    expect(new Date(paidSub.periodEnd).getTime()).toBeGreaterThan(Date.now() + 27 * 86_400_000);
    expect((await api.post('/api/products', { nameAr: 'بعد التجديد' })).statusCode).toBe(201);
    const events = (await api.get('/api/subscription/billing-events')).json().data.map((e: { eventType: string }) => e.eventType);
    expect(events).toEqual(['PAYMENT_RECORDED', 'PLAN_CHANGED', 'PLAN_CHANGE_REQUESTED', 'TRIAL_STARTED']);
    expect((await api.get('/api/subscription')).json().items).toEqual([expect.objectContaining({ unitPrice: '99.00' })]);
  });

  it('limits plan change requests to subscription managers', async () => {
    const s = await register(t.app);
    const accountant = client(t.app, (await addMember(t.app, s, 'ACCOUNTANT')).token);
    const plan = (await client(t.app, s.token).get('/api/subscription/plans')).json().data[0];
    expect((await accountant.post('/api/subscription/request-change', { planId: plan.id, billingCycle: 'MONTHLY' })).statusCode).toBe(403);
    expect((await accountant.get('/api/subscription')).statusCode).toBe(200);
  });
});
