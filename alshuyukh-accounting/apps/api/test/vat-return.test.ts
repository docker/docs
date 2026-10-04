import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { commerce, type Commerce } from './commerce-helpers.js';
import { chart, cr, dr, fiscalYears, register, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;
let c: Commerce;
let range: string;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
  c = await commerce(t.app, owner);
  const [year] = await fiscalYears(t.app, owner.token);
  const p = year!.periods[2]!; // the month that contains c.date
  range = `dateFrom=${p.startDate}&dateTo=${p.endDate}`;

  const cust = (await c.customer()).id;
  const sup = (await c.supplier()).id;
  const std = (await c.product({ salePrice: '200' })).id;
  const zero = (await c.product({ salePrice: '200', vatCategory: 'Z' })).id;
  const exempt = (await c.product({ salePrice: '50', vatCategory: 'E' })).id;
  const goods = (await c.goods({ purchasePrice: '60' })).id;
  const cats = Object.fromEntries((await c.api.get('/api/expense-categories')).json().data.map((x: { code: string; id: string }) => [x.code, x.id]));

  // Sales: standard 1000 (VAT 150), zero-rated 200, exempt 50.
  const inv = await c.postInvoice(cust, [{ productId: std, quantity: '5' }, { productId: zero, quantity: '1' }, { productId: exempt, quantity: '1' }]);
  // Sales return of one standard unit: −200 / −30.
  const r = (await c.api.post('/api/sales-returns', { originalInvoiceId: inv.id, docDate: c.date, reason: 'إرجاع', lines: [{ sourceItemId: inv.lines[0].id, quantity: '1' }] })).json();
  await c.api.post(`/api/sales-returns/${r.id}/post`);
  // A cancelled invoice nets to zero.
  const cancelled = await c.postInvoice(cust, [{ productId: std, quantity: '1' }]);
  await c.api.post(`/api/invoices/${cancelled.id}/cancel`, { reason: 'خطأ' });
  // Purchases: goods 600 (VAT 90), marketing expense 1000 (150), exempt rent 5000.
  const bill = await c.purchase(sup, [{ productId: goods, quantity: '10' }]);
  await c.api.post('/api/expenses', { expenseDate: c.date, paymentType: 'CASH', methodId: await c.method('CASH'), post: true, lines: [{ categoryId: cats.MARKETING, amount: '1000' }] });
  await c.api.post('/api/expenses', { expenseDate: c.date, paymentType: 'CASH', methodId: await c.method('CASH'), post: true, lines: [{ categoryId: cats.RENT, amount: '5000', vatCategory: 'E' }] });
  // Purchase return of goods worth 100: −100 / −15.
  const pr = (await c.api.post('/api/purchase-returns', { originalInvoiceId: bill.id, docDate: c.date, reason: 'معيب', lines: [{ sourceItemId: bill.lines[0].id, quantity: '1' }] })).json();
  await c.api.post(`/api/purchase-returns/${pr.id}/post`);
});
afterAll(async () => { await t.close(); });

describe('VAT return', () => {
  it('fills the ZATCA boxes from posted documents', async () => {
    const res = await c.api.get(`/api/reports/vat-return?${range}`);
    expect(res.statusCode).toBe(200);
    const v = res.json();
    expect(v.sales['1_standardRated']).toEqual({ amount: '1000.00', adjustment: '-200.00', vat: '120.00' });
    expect(v.sales['3_zeroRatedDomestic']).toEqual({ amount: '200.00', adjustment: '0.00', vat: '0.00' });
    expect(v.sales['5_exempt']).toEqual({ amount: '50.00', adjustment: '0.00', vat: '0.00' });
    expect(v.sales['6_total']).toEqual({ amount: '1250.00', adjustment: '-200.00', vat: '120.00' });
    // Purchases: goods 600 + expense 1000 (standard), return −60 / −9, rent 5000 exempt.
    expect(v.purchases['7_standardRatedDomestic']).toEqual({ amount: '1600.00', adjustment: '-60.00', vat: '231.00' });
    expect(v.purchases['11_exempt']).toEqual({ amount: '5000.00', adjustment: '0.00', vat: '0.00' });
    expect(v.purchases['12_total']).toEqual({ amount: '6600.00', adjustment: '-60.00', vat: '231.00' });
    expect(v['13_totalVatDue']).toBe('-111.00');
    expect(v['16_netVatDue']).toBe('-111.00');
    expect(v.notSupported).toContain('4_exports');
  });

  it('reconciles with the VAT accounts in the ledger', async () => {
    const v = (await c.api.get(`/api/reports/vat-return?${range}`)).json();
    expect(v.ledger).toMatchObject({ vatOutput: '120.00', vatInput: '231.00', reconciled: true });
  });

  it('flags manual journal entries to VAT accounts', async () => {
    const s = await register(t.app);
    const sc = await commerce(t.app, s);
    const ch = await chart(t.app, s.token);
    await sc.api.post('/api/journal-entries', { entryDate: sc.date, description: 'قيد يدوي على الضريبة', post: true,
      lines: [dr(ch.byKey.get('CASH')!, '10.00'), cr(ch.byKey.get('VAT_OUTPUT')!, '10.00')] });
    const [year] = await fiscalYears(t.app, s.token);
    const v = (await sc.api.get(`/api/reports/vat-return?dateFrom=${year!.startDate}&dateTo=${year!.endDate}`)).json();
    expect(v.ledger).toMatchObject({ outputDifference: '10.00', reconciled: false });
  });

  it('only counts transactions inside the period', async () => {
    const [year] = await fiscalYears(t.app, owner.token);
    const v = (await c.api.get(`/api/reports/vat-return?dateFrom=${year!.periods[5]!.startDate}&dateTo=${year!.periods[5]!.endDate}`)).json();
    expect(v.sales['6_total']).toEqual({ amount: '0.00', adjustment: '0.00', vat: '0.00' });
  });

  it('lists the transactions behind the return', async () => {
    const rows = (await c.api.get(`/api/reports/vat-transactions?${range}&direction=OUTPUT&vatCategory=S`)).json().data;
    expect(rows.map((r: { sourceType: string; taxAmount: string }) => [r.sourceType, r.taxAmount])).toEqual([
      ['SALES_INVOICE', '150.00'], ['SALES_RETURN', '-30.00'], ['SALES_INVOICE', '30.00'], ['SALES_INVOICE', '-30.00'],
    ]);
  });
});
