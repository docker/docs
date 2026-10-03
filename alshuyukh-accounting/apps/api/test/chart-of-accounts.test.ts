import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { DEFAULT_ACCOUNTS } from '../src/modules/accounting/chart-template.js';
import { chart, client, cr, dateInYear, dr, fiscalYears, register, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
});
afterAll(async () => { await t.close(); });

describe('default setup for a new company', () => {
  it('creates the default Saudi chart of accounts with system keys', async () => {
    const accounts = (await client(t.app, owner.token).get('/api/accounts')).json().data;
    expect(accounts.map((a: { code: string }) => a.code)).toEqual(DEFAULT_ACCOUNTS.map((a) => a.code).sort());
    const byCode = Object.fromEntries(accounts.map((a: { code: string }) => [a.code, a]));
    expect(byCode['1000']).toMatchObject({ nameAr: 'الأصول', type: 'ASSET', isPostable: false, level: 1 });
    expect(byCode['1300']).toMatchObject({ nameAr: 'العملاء', systemKey: 'ACCOUNTS_RECEIVABLE', isPostable: true, level: 2, balance: '0.00' });
    expect(byCode['2210']).toMatchObject({ systemKey: 'VAT_OUTPUT', level: 3, type: 'LIABILITY' });
    expect(byCode['3200']).toMatchObject({ systemKey: 'RETAINED_EARNINGS', type: 'EQUITY' });
    expect(byCode['4100']).toMatchObject({ systemKey: 'SALES', type: 'REVENUE' });
    expect(byCode['5100']).toMatchObject({ systemKey: 'COGS', type: 'COST_OF_GOODS_SOLD' });
  });

  it('creates the current fiscal year with 12 open monthly periods', async () => {
    const [year] = await fiscalYears(t.app, owner.token);
    expect(year!.status).toBe('OPEN');
    expect(year!.periods).toHaveLength(12);
    expect(year!.startDate).toMatch(/-01-01$/);
    expect(year!.endDate).toMatch(/-12-31$/);
    expect(year!.periods[1]).toMatchObject({ startDate: `${year!.startDate.slice(0, 4)}-02-01`, status: 'OPEN' });
    expect(year!.periods.every((p) => p.status === 'OPEN')).toBe(true);
  });

  it('sets up accounting for companies created later, and the setup call is idempotent', async () => {
    const api = client(t.app, owner.token);
    const company = (await api.post('/api/companies', { name: `Second ${Date.now()}` })).json();
    const accounts = (await api.get(`/api/accounts?companyId=${company.id}`)).json().data;
    expect(accounts.length).toBe(DEFAULT_ACCOUNTS.length);
    const again = await api.post('/api/accounting/setup', { companyId: company.id });
    expect(again.json()).toMatchObject({ chartCreated: false, fiscalYearCreated: false });
    // With two companies, list endpoints require companyId.
    expect((await api.get('/api/accounts')).statusCode).toBe(400);
  });
});

describe('managing accounts', () => {
  let session: Session;
  beforeAll(async () => { session = await register(t.app); });

  it('creates a sub-account that inherits the parent type', async () => {
    const api = client(t.app, session.token);
    const c = await chart(t.app, session.token);
    const res = await api.post('/api/accounts', { code: '6500', nameAr: 'مصروفات الصيانة', parentId: c.byCode.get('6000') });
    expect(res.statusCode).toBe(201);
    expect(res.json()).toMatchObject({ type: 'EXPENSE', level: 2, isPostable: true, isSystem: false });
  });

  it('turns an unused leaf into a header when a child is added under it', async () => {
    const api = client(t.app, session.token);
    const c = await chart(t.app, session.token);
    const res = await api.post('/api/accounts', { code: '1210', nameAr: 'بنك الراجحي', parentId: c.byCode.get('1200') });
    expect(res.statusCode).toBe(201);
    expect(res.json().level).toBe(3);
    expect((await api.get(`/api/accounts/${c.byCode.get('1200')}`)).json().isPostable).toBe(false);
  });

  it('rejects invalid structures', async () => {
    const api = client(t.app, session.token);
    const c = await chart(t.app, session.token);
    expect((await api.post('/api/accounts', { code: '6100', nameAr: 'مكرر', parentId: c.byCode.get('6000') })).statusCode).toBe(409);
    expect((await api.post('/api/accounts', { code: '6600', nameAr: 'نوع خاطئ', type: 'ASSET', parentId: c.byCode.get('6000') })).statusCode).toBe(400);
    expect((await api.post('/api/accounts', { code: '9000', nameAr: 'بلا نوع' })).statusCode).toBe(400);
    expect((await api.post('/api/accounts', { code: 'bad code!', nameAr: 'رمز', type: 'ASSET' })).statusCode).toBe(400);

    const h1 = (await api.post('/api/accounts', { code: '7000', nameAr: 'رئيسي', type: 'EXPENSE', isPostable: false })).json();
    const h2 = (await api.post('/api/accounts', { code: '7100', nameAr: 'فرعي', parentId: h1.id, isPostable: false })).json();
    const cycle = await api.patch(`/api/accounts/${h1.id}`, { parentId: h2.id });
    expect(cycle.statusCode).toBe(400);
    expect(cycle.json().error.message).toMatch(/descendant/);
  });

  it('protects system accounts and accounts with postings', async () => {
    const api = client(t.app, session.token);
    const c = await chart(t.app, session.token);
    expect((await api.patch(`/api/accounts/${c.byKey.get('CASH')}`, { isActive: false })).statusCode).toBe(409);
    expect((await api.del(`/api/accounts/${c.byKey.get('CASH')}`)).statusCode).toBe(409);
    // Renaming a system account is allowed.
    expect((await api.patch(`/api/accounts/${c.byKey.get('CASH')}`, { nameAr: 'الصندوق' })).statusCode).toBe(200);

    const rent = c.byCode.get('6100')!;
    const date = await dateInYear(t.app, session.token);
    const posted = await api.post('/api/journal-entries', { entryDate: date, description: 'إيجار', post: true, lines: [dr(rent, '500.00'), cr(c.byKey.get('CASH')!, '500.00')] });
    expect(posted.statusCode).toBe(201);
    expect((await api.get(`/api/accounts/${rent}`)).json().balance).toBe('500.00');
    expect((await api.patch(`/api/accounts/${rent}`, { isActive: false })).statusCode).toBe(409);
    expect((await api.del(`/api/accounts/${rent}`)).statusCode).toBe(409);
    const toHeader = await api.patch(`/api/accounts/${rent}`, { isPostable: false });
    expect(toHeader.statusCode).toBe(400);
  });

  it('deactivates and soft-deletes unused accounts', async () => {
    const api = client(t.app, session.token);
    const c = await chart(t.app, session.token);
    const acc = (await api.post('/api/accounts', { code: '6900', nameAr: 'مؤقت', parentId: c.byCode.get('6000') })).json();
    expect((await api.patch(`/api/accounts/${acc.id}`, { isActive: false })).json().isActive).toBe(false);
    const listed = (await api.get('/api/accounts')).json().data.map((a: { id: string }) => a.id);
    expect(listed).not.toContain(acc.id);
    expect((await api.del(`/api/accounts/${acc.id}`)).statusCode).toBe(204);
    expect((await api.get(`/api/accounts/${acc.id}`)).statusCode).toBe(404);
    // The code can be reused after deletion.
    expect((await api.post('/api/accounts', { code: '6900', nameAr: 'جديد', parentId: c.byCode.get('6000') })).statusCode).toBe(201);
  });

  it('manages cost centers', async () => {
    const api = client(t.app, session.token);
    const cc = await api.post('/api/cost-centers', { code: 'RUH', name: 'مركز الرياض' });
    expect(cc.statusCode).toBe(201);
    expect((await api.post('/api/cost-centers', { code: 'RUH', name: 'مكرر' })).statusCode).toBe(409);
    expect((await api.patch(`/api/cost-centers/${cc.json().id}`, { isActive: false })).json().isActive).toBe(false);
  });

  it('isolates charts between tenants', async () => {
    const other = await register(t.app);
    const c = await chart(t.app, session.token);
    const api = client(t.app, other.token);
    expect((await api.get(`/api/accounts/${c.byKey.get('CASH')}`)).statusCode).toBe(404);
    expect((await api.patch(`/api/accounts/${c.byKey.get('CASH')}`, { nameAr: 'اختراق' })).statusCode).toBe(404);
    expect((await api.post('/api/accounts', { code: '6700', nameAr: 'حقن', parentId: c.byCode.get('6000') })).statusCode).toBe(400);
  });
});
