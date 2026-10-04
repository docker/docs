import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { withTx } from '../src/db/tx.js';
import { balanceOf, commerce, type Commerce } from './commerce-helpers.js';
import { addMember, client, register, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;
let c: Commerce;
let cat: Record<string, string>;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
  c = await commerce(t.app, owner);
  const cats = (await c.api.get('/api/expense-categories')).json().data as { id: string; code: string }[];
  cat = Object.fromEntries(cats.map((x) => [x.code, x.id]));
});
afterAll(async () => { await t.close(); });

const expense = async (body: Record<string, unknown>, post = true) =>
  c.api.post('/api/expenses', { expenseDate: c.date, paymentType: 'CASH', methodId: await c.method('CASH'), post, ...body });
const lines = async (id: string) => (await c.journal(id)).map((l) => [l.accountCode, l.debit, l.credit]);

describe('expense categories', () => {
  it('seeds categories with their accounts and default VAT treatment', async () => {
    const cats = (await c.api.get('/api/expense-categories')).json().data;
    expect(cats.map((x: { code: string; accountCode: string; vatCategory: string }) => [x.code, x.accountCode, x.vatCategory]).sort()).toEqual([
      ['MARKETING', '6300', 'S'], ['RENT', '6100', 'S'], ['SALARIES', '6200', 'O'], ['UTILITIES', '6400', 'S'],
    ]);
  });
});

describe('posting expenses', () => {
  it('posts a cash expense: Dr Expense + Dr VAT input / Cr Cash', async () => {
    const res = await expense({ payeeName: 'وكالة إعلانات', reference: 'A-77', vendorVatNumber: '311111111111113', lines: [{ categoryId: cat.MARKETING, amount: '1000' }] });
    expect(res.statusCode).toBe(201);
    const e = res.json();
    expect(e).toMatchObject({ status: 'PAID', subtotal: '1000.00', taxAmount: '150.00', total: '1150.00', remainingAmount: '0.00' });
    expect(e.number).toMatch(/^EXP-\d{6}$/);
    expect(await lines(e.journalEntryId)).toEqual([['6300', '1000.00', '0.00'], ['2220', '150.00', '0.00'], ['1100', '0.00', '1150.00']]);
    const tx = (await c.api.get(`/api/reports/vat-transactions?dateFrom=${c.date}&dateTo=${c.date}&direction=INPUT`)).json().data;
    expect(tx.find((x: { sourceId: string }) => x.sourceId === e.id)).toMatchObject({ sourceType: 'EXPENSE', vatCategory: 'S', taxableAmount: '1000.00', taxAmount: '150.00', partyVatNumber: '311111111111113' });
  });

  it('extracts VAT from inclusive amounts and honours a per-line exempt override', async () => {
    const incl = (await expense({ pricesIncludeVat: true, lines: [{ categoryId: cat.UTILITIES, amount: '115' }] })).json();
    expect(incl).toMatchObject({ subtotal: '100.00', taxAmount: '15.00', total: '115.00' });
    const rent = (await expense({ paymentType: 'BANK', methodId: await c.method('BANK'), lines: [{ categoryId: cat.RENT, amount: '5000', vatCategory: 'E', description: 'سكن موظفين' }] })).json();
    expect(rent).toMatchObject({ taxAmount: '0.00', total: '5000.00' });
    expect(await lines(rent.journalEntryId)).toEqual([['6100', '5000.00', '0.00'], ['1200', '0.00', '5000.00']]);
  });

  it('groups lines per account and cost center', async () => {
    const cc = (await c.api.post('/api/cost-centers', { code: 'MKT', name: 'التسويق' })).json();
    const e = (await expense({ lines: [
      { categoryId: cat.MARKETING, amount: '100', costCenterId: cc.id },
      { categoryId: cat.MARKETING, amount: '50', costCenterId: cc.id },
      { categoryId: cat.MARKETING, amount: '30' },
    ] })).json();
    const entry = (await c.api.get(`/api/journal-entries/${e.journalEntryId}`)).json();
    expect(entry.lines.map((l: { accountCode: string; debit: string; costCenterId: string | null }) => [l.accountCode, l.debit, l.costCenterId])).toEqual([
      ['6300', '150.00', cc.id], ['6300', '30.00', null], ['2220', '27.00', null], ['1100', '0.00', null],
    ]);
  });

  it('records a credit expense against the supplier and settles it with a supplier payment', async () => {
    expect((await expense({ paymentType: 'CREDIT', methodId: null, lines: [{ categoryId: cat.UTILITIES, amount: '200' }] })).json().error.code).toBe('SUPPLIER_REQUIRED');
    const sup = await c.supplier({ vatNumber: '322222222222223' });
    const e = (await expense({ paymentType: 'CREDIT', methodId: null, supplierId: sup.id, lines: [{ categoryId: cat.UTILITIES, amount: '200' }] })).json();
    expect(e).toMatchObject({ status: 'POSTED', total: '230.00', remainingAmount: '230.00', supplierName: expect.any(String) });
    expect(await lines(e.journalEntryId)).toEqual([['6400', '200.00', '0.00'], ['2220', '30.00', '0.00'], ['2100', '0.00', '230.00']]);
    expect(await balanceOf(c.api, 'suppliers', sup.id)).toBe('-230.00');

    const open = (await c.api.get(`/api/expenses?supplierId=${sup.id}&open=true`)).json();
    expect(open.total).toBe(1);
    const pay = await c.api.post('/api/payments', {
      direction: 'DISBURSEMENT', supplierId: sup.id, paymentDate: c.date, methodId: await c.method('BANK'), amount: '230.00',
      allocations: [{ documentType: 'EXPENSE', documentId: e.id, amount: '230.00' }],
    });
    expect(pay.statusCode).toBe(201);
    expect((await c.api.get(`/api/expenses/${e.id}`)).json()).toMatchObject({ status: 'PAID', paidAmount: '230.00' });
    expect(await balanceOf(c.api, 'suppliers', sup.id)).toBe('0.00');
    expect((await c.api.post(`/api/expenses/${e.id}/cancel`, { reason: 'خطأ' })).json().error.code).toBe('HAS_PAYMENTS');

    await c.api.post(`/api/payments/${pay.json().id}/void`, { reason: 'تحويل مرتجع' });
    expect((await c.api.get(`/api/expenses/${e.id}`)).json()).toMatchObject({ status: 'POSTED', paidAmount: '0.00' });
  });

  it('cancels a posted expense by reversing its entry and VAT', async () => {
    const e = (await expense({ lines: [{ categoryId: cat.MARKETING, amount: '400' }] })).json();
    const res = await c.api.post(`/api/expenses/${e.id}/cancel`, { reason: 'أُدخل مرتين' });
    expect(res.json().status).toBe('CANCELLED');
    expect((await c.api.get(`/api/journal-entries/${e.journalEntryId}`)).json().status).toBe('REVERSED');
    const tx = (await c.api.get(`/api/reports/vat-transactions?dateFrom=${c.date}&dateTo=${c.date}`)).json().data
      .filter((x: { sourceId: string }) => x.sourceId === e.id);
    expect(tx.map((x: { taxAmount: string }) => x.taxAmount).sort()).toEqual(['-60.00', '60.00']);
  });
});

describe('drafts, validation and protection', () => {
  it('edits and deletes drafts; refuses changes after posting', async () => {
    const d = (await expense({ lines: [{ categoryId: cat.MARKETING, amount: '10' }] }, false)).json();
    expect(d).toMatchObject({ status: 'DRAFT', number: null });
    const edited = (await c.api.patch(`/api/expenses/${d.id}`, { lines: [{ categoryId: cat.UTILITIES, amount: '20' }] })).json();
    expect(edited).toMatchObject({ total: '23.00' });
    const posted = (await c.api.post(`/api/expenses/${d.id}/post`)).json();
    expect(posted.status).toBe('PAID');
    expect((await c.api.patch(`/api/expenses/${d.id}`, { notes: 'x' })).statusCode).toBe(409);
    expect((await c.api.del(`/api/expenses/${d.id}`)).statusCode).toBe(409);
    await expect(withTx(t.pool, { tenantId: owner.tenantId, userId: owner.userId }, (db) =>
      db.query(`UPDATE expenses SET total = 1 WHERE id = $1`, [d.id]))).rejects.toThrow(/cannot be modified/);
    const other = (await expense({ lines: [{ categoryId: cat.MARKETING, amount: '10' }] }, false)).json();
    expect((await c.api.del(`/api/expenses/${other.id}`)).statusCode).toBe(204);
  });

  it('validates lines and references', async () => {
    expect((await expense({ lines: [{ categoryId: cat.MARKETING, amount: '0' }] })).statusCode).toBe(400);
    expect((await expense({ lines: [{ categoryId: cat.MARKETING, amount: '-5' }] })).statusCode).toBe(400);
    expect((await expense({ lines: [{ categoryId: cat.MARKETING, amount: 5 }] })).statusCode).toBe(400);
    expect((await expense({ methodId: null, lines: [{ categoryId: cat.MARKETING, amount: '5' }] })).json().error.code).toBe('METHOD_REQUIRED');
    const other = await register(t.app);
    const foreignCat = (await client(t.app, other.token).get('/api/expense-categories')).json().data[0].id;
    expect((await expense({ lines: [{ categoryId: foreignCat, amount: '5' }] })).json().error.code).toBe('INVALID_CATEGORY');
  });
});

describe('permissions and isolation', () => {
  it('applies expense permissions', async () => {
    const sales = await addMember(t.app, owner, 'SALES_EMPLOYEE');
    expect((await client(t.app, sales.token).get('/api/expenses')).statusCode).toBe(403);
    const viewer = await addMember(t.app, owner, 'VIEWER');
    expect((await client(t.app, viewer.token).get('/api/expenses')).statusCode).toBe(200);
    expect((await client(t.app, viewer.token).post('/api/expenses', { expenseDate: c.date, paymentType: 'CASH', methodId: await c.method('CASH'), lines: [{ categoryId: cat.MARKETING, amount: '1' }] })).statusCode).toBe(403);
    const accountant = await addMember(t.app, owner, 'ACCOUNTANT');
    const aApi = client(t.app, accountant.token);
    const e = await aApi.post('/api/expenses', { expenseDate: c.date, paymentType: 'CASH', methodId: await c.method('CASH'), post: true, lines: [{ categoryId: cat.MARKETING, amount: '1' }] });
    expect(e.statusCode).toBe(201);
    expect((await aApi.post(`/api/expenses/${e.json().id}/cancel`, { reason: 'تجربة' })).statusCode).toBe(200);
  });

  it('hides expenses from other tenants', async () => {
    const e = (await expense({ lines: [{ categoryId: cat.MARKETING, amount: '1' }] })).json();
    const other = await register(t.app);
    const oApi = client(t.app, other.token);
    expect((await oApi.get(`/api/expenses/${e.id}`)).statusCode).toBe(404);
    expect((await oApi.post(`/api/expenses/${e.id}/cancel`, { reason: 'اختراق' })).statusCode).toBe(404);
    expect((await oApi.get('/api/expenses')).json().total).toBe(0);
  });
});
