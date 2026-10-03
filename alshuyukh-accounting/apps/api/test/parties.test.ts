import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { addMember, chart, client, cr, dateInYear, dr, register, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
});
afterAll(async () => { await t.close(); });

const api = () => client(t.app, owner.token);

const address = {
  addressType: 'BILLING', buildingNumber: '1234', street: 'طريق الملك فهد', district: 'العليا',
  city: 'الرياض', postalCode: '12345', additionalNumber: '6789',
};

describe.each([
  ['customers', 'CUS', 'ACCOUNTS_RECEIVABLE', 'customer'],
  ['suppliers', 'SUP', 'ACCOUNTS_PAYABLE', 'supplier'],
] as const)('%s', (path, prefix, controlKey, entity) => {
  it('creates a party with an automatic code and a national address', async () => {
    const res = await api().post(`/api/${path}`, {
      nameAr: 'شركة النخبة للتجارة', vatNumber: '300000000000003', commercialRegistration: '1010101010',
      phone: '+966 50 123 4567', email: 'Info@Nukhba.SA', paymentTermsDays: 30, creditLimit: '50000.00',
      addresses: [address],
    });
    expect(res.statusCode).toBe(201);
    const p = res.json();
    expect(p.code).toMatch(new RegExp(`^${prefix}-\\d{5}$`));
    expect(p).toMatchObject({ partyType: 'BUSINESS', email: 'info@nukhba.sa', creditLimit: '50000.00', balance: '0.00', isActive: true });
    expect(p.addresses).toEqual([expect.objectContaining({ ...address, isDefault: true, country: 'SA' })]);

    const next = (await api().post(`/api/${path}`, { nameAr: 'عميل ثان' })).json();
    expect(Number(next.code.slice(-5))).toBe(Number(p.code.slice(-5)) + 1);
  });

  it('skips generated codes that were already taken manually', async () => {
    const taken = (await api().post(`/api/${path}`, { nameAr: 'يدوي', code: `${prefix}-99999` })).json();
    expect(taken.code).toBe(`${prefix}-99999`);
    expect((await api().post(`/api/${path}`, { nameAr: 'مكرر', code: `${prefix}-99999` })).statusCode).toBe(409);
  });

  it('validates Saudi identifiers and address formats', async () => {
    const bad = [
      { nameAr: 'خطأ', vatNumber: '123456789012345' },
      { nameAr: 'خطأ', commercialRegistration: '123' },
      { nameAr: 'خطأ', nationalId: '3123456789' },
      { nameAr: 'خطأ', creditLimit: 100 },
      { nameAr: 'خطأ', addresses: [{ ...address, postalCode: '123' }] },
      { nameAr: 'خطأ', addresses: [{ ...address, buildingNumber: '12' }] },
      { nameAr: 'خطأ', addresses: [{ ...address, isDefault: true }, { ...address, isDefault: true }] },
      { nameAr: 'خ' },
    ];
    for (const body of bad) expect((await api().post(`/api/${path}`, body)).statusCode, JSON.stringify(body)).toBe(400);
    expect((await api().post(`/api/${path}`, { nameAr: 'تكرار ضريبي', vatNumber: '300000000000003' })).statusCode).toBe(409);
  });

  it('searches by name, code, VAT number and phone', async () => {
    await api().post(`/api/${path}`, { nameAr: 'مؤسسة البحث الفريدة', phone: '0555000111', vatNumber: '311111111111113' });
    for (const term of ['البحث الفريدة', '0555000111', '311111111111113']) {
      const res = (await api().get(`/api/${path}?search=${encodeURIComponent(term)}`)).json();
      expect(res.data.map((p: { nameAr: string }) => p.nameAr), term).toEqual(['مؤسسة البحث الفريدة']);
    }
  });

  it('updates, replaces addresses, and deactivates', async () => {
    const p = (await api().post(`/api/${path}`, { nameAr: 'للتعديل', addresses: [address] })).json();
    const updated = (await api().patch(`/api/${path}/${p.id}`, { nameAr: 'بعد التعديل', paymentTermsDays: 60, phone: null })).json();
    expect(updated).toMatchObject({ nameAr: 'بعد التعديل', paymentTermsDays: 60, phone: null });

    const addrs = (await api().put(`/api/${path}/${p.id}/addresses`, {
      addresses: [{ ...address, city: 'جدة' }, { ...address, addressType: 'SHIPPING', city: 'الدمام' }],
    })).json().addresses;
    expect(addrs.map((a: { city: string; isDefault: boolean }) => [a.city, a.isDefault])).toEqual([['جدة', true], ['الدمام', true]]);

    expect((await api().patch(`/api/${path}/${p.id}`, { isActive: false })).json().isActive).toBe(false);
    const active = (await api().get(`/api/${path}?search=${encodeURIComponent('بعد التعديل')}`)).json().data;
    expect(active).toEqual([]);
    const all = (await api().get(`/api/${path}?status=all&search=${encodeURIComponent('بعد التعديل')}`)).json().data;
    expect(all).toHaveLength(1);
  });

  it('accepts only a matching control account', async () => {
    const c = await chart(t.app, owner.token);
    expect((await api().post(`/api/${path}`, { nameAr: 'حساب خاطئ', controlAccountId: c.byKey.get('SALES') })).statusCode).toBe(400);
    const ok = await api().post(`/api/${path}`, { nameAr: 'حساب صحيح', controlAccountId: c.byKey.get(controlKey) });
    expect(ok.statusCode).toBe(201);
    expect(ok.json().controlAccountId).toBe(c.byKey.get(controlKey));
  });

  it('takes its balance from tagged journal lines and cannot be deleted afterwards', async () => {
    const c = await chart(t.app, owner.token);
    const p = (await api().post(`/api/${path}`, { nameAr: 'رصيد افتتاحي' })).json();
    const tag = entity === 'customer' ? { customerId: p.id } : { supplierId: p.id };
    const control = c.byKey.get(controlKey)!;
    const capital = c.byKey.get('CAPITAL')!;
    const lines = entity === 'customer'
      ? [{ ...dr(control, '2500.00'), ...tag }, cr(capital, '2500.00')]
      : [dr(capital, '2500.00'), { ...cr(control, '2500.00'), ...tag }];
    const entry = await api().post('/api/journal-entries', { entryDate: await dateInYear(t.app, owner.token), description: 'رصيد افتتاحي', post: true, lines });
    expect(entry.statusCode).toBe(201);
    const tagged = entry.json().lines.filter((l: { customerId: string | null; supplierId: string | null }) => (l.customerId ?? l.supplierId) === p.id);
    expect(tagged).toHaveLength(1);

    const after = (await api().get(`/api/${path}/${p.id}`)).json();
    expect(after.balance).toBe(entity === 'customer' ? '2500.00' : '-2500.00');
    expect((await api().del(`/api/${path}/${p.id}`)).json().error.code).toBe('PARTY_IN_USE');

    const unused = (await api().post(`/api/${path}`, { nameAr: 'للحذف' })).json();
    expect((await api().del(`/api/${path}/${unused.id}`)).statusCode).toBe(204);
    expect((await api().get(`/api/${path}/${unused.id}`)).statusCode).toBe(404);
  });

  it('isolates parties between tenants, including journal tagging', async () => {
    const p = (await api().post(`/api/${path}`, { nameAr: 'خاص بالمنشأة' })).json();
    const other = await register(t.app);
    const oApi = client(t.app, other.token);
    expect((await oApi.get(`/api/${path}/${p.id}`)).statusCode).toBe(404);
    expect((await oApi.patch(`/api/${path}/${p.id}`, { nameAr: 'اختراق' })).statusCode).toBe(404);
    expect((await oApi.get(`/api/${path}`)).json().total).toBe(0);
    const oc = await chart(t.app, other.token);
    const tag = entity === 'customer' ? { customerId: p.id } : { supplierId: p.id };
    const res = await oApi.post('/api/journal-entries', {
      entryDate: await dateInYear(t.app, other.token), description: 'حقن',
      lines: [{ ...dr(oc.byKey.get(controlKey)!, '1.00'), ...tag }, cr(oc.byKey.get('CAPITAL')!, '1.00')],
    });
    expect(res.statusCode).toBe(400);
  });
});

describe('party permissions', () => {
  it('lets each role manage only its own parties', async () => {
    const sales = await addMember(t.app, owner, 'SALES_MANAGER');
    const sApi = client(t.app, sales.token);
    expect((await sApi.post('/api/customers', { nameAr: 'من المبيعات' })).statusCode).toBe(201);
    expect((await sApi.post('/api/suppliers', { nameAr: 'من المبيعات' })).statusCode).toBe(403);
    expect((await sApi.get('/api/suppliers')).statusCode).toBe(403);

    const purchase = await addMember(t.app, owner, 'PURCHASE_MANAGER');
    const pApi = client(t.app, purchase.token);
    expect((await pApi.post('/api/suppliers', { nameAr: 'من المشتريات' })).statusCode).toBe(201);
    expect((await pApi.post('/api/customers', { nameAr: 'من المشتريات' })).statusCode).toBe(403);

    const employee = await addMember(t.app, owner, 'SALES_EMPLOYEE');
    const eApi = client(t.app, employee.token);
    expect((await eApi.get('/api/customers')).statusCode).toBe(200);
    expect((await eApi.post('/api/customers', { nameAr: 'موظف' })).statusCode).toBe(403);

    const accountant = await addMember(t.app, owner, 'ACCOUNTANT');
    expect((await client(t.app, accountant.token).post('/api/suppliers', { nameAr: 'من المحاسب' })).statusCode).toBe(201);
  });

  it('a party cannot be tagged on both sides of one line', async () => {
    const c = await chart(t.app, owner.token);
    const cust = (await api().post('/api/customers', { nameAr: 'عميل الوسم' })).json();
    const supp = (await api().post('/api/suppliers', { nameAr: 'مورد الوسم' })).json();
    const res = await api().post('/api/journal-entries', {
      entryDate: await dateInYear(t.app, owner.token), description: 'خطأ',
      lines: [{ ...dr(c.byKey.get('CASH')!, '1.00'), customerId: cust.id, supplierId: supp.id }, cr(c.byKey.get('CAPITAL')!, '1.00')],
    });
    expect(res.statusCode).toBe(400);
  });

  it('records audit entries for party changes', async () => {
    const actions = (await api().get('/api/audit-logs?entityType=customer&limit=200')).json().data.map((l: { action: string }) => l.action);
    expect(actions).toEqual(expect.arrayContaining(['CREATE', 'UPDATE', 'DELETE']));
  });
});
