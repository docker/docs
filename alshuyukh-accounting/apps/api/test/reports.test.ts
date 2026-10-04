import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { commerce, type Commerce } from './commerce-helpers.js';
import { addMember, chart, client, cr, dr, fiscalYears, register, setupApp, type Chart, type Session, type TestContext } from './helpers.js';

/**
 * One company, one realistic month. Every report must agree with the others:
 *   capital 15,000 (cash 10,000 + bank 5,000)
 *   purchase 10 goods × 60 on credit (VAT 90)                 → AP 690
 *   invoice to A: 4 goods × 100 + 1 service × 100 (VAT 75)    → AR 575, COGS 240
 *   sales return of 1 good (−100, VAT −15, COGS −60)
 *   receipt from A 300 (cash), advance from B 50 (cash)
 *   cash expense: marketing 200 + VAT 30
 *   supplier payment 400 (bank)
 *   an invoice to B that is cancelled (must vanish everywhere)
 * Expected: revenue 400, COGS 180, gross 220, expenses 200, net profit 20.
 */
let t: TestContext;
let owner: Session;
let c: Commerce;
let ch: Chart;
let y: { from: string; to: string };
let A: string, B: string, S: string, goods: string, service: string, invoice: Record<string, any>;
const get = async (path: string) => {
  const r = await c.api.get(path);
  if (r.statusCode !== 200) throw new Error(`${path}: ${r.statusCode} ${r.body}`);
  return r.json();
};
const range = () => `dateFrom=${y.from}&dateTo=${y.to}`;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
  c = await commerce(t.app, owner);
  ch = await chart(t.app, owner.token);
  const [year] = await fiscalYears(t.app, owner.token);
  y = { from: year!.startDate, to: year!.endDate };

  await c.api.post('/api/journal-entries', { entryDate: c.date, description: 'رأس المال', post: true,
    lines: [dr(ch.byKey.get('CASH')!, '10000.00'), dr(ch.byKey.get('BANK')!, '5000.00'), cr(ch.byKey.get('CAPITAL')!, '15000.00')] });
  A = (await c.customer()).id;
  B = (await c.customer()).id;
  S = (await c.supplier()).id;
  goods = (await c.goods({ salePrice: '100', purchasePrice: '60' })).id;
  service = (await c.product({ salePrice: '100' })).id;
  const bill = await c.purchase(S, [{ productId: goods, quantity: '10' }]);
  invoice = await c.postInvoice(A, [{ productId: goods, quantity: '4' }, { productId: service, quantity: '1' }], { dueDate: c.date });
  const ret = (await c.api.post('/api/sales-returns', { originalInvoiceId: invoice.id, docDate: c.date, reason: 'تالف', lines: [{ sourceItemId: invoice.lines[0].id, quantity: '1' }] })).json();
  await c.api.post(`/api/sales-returns/${ret.id}/post`);
  const cash = await c.method('CASH');
  await c.api.post('/api/payments', { direction: 'RECEIPT', customerId: A, paymentDate: c.date, methodId: cash, amount: '300.00',
    allocations: [{ documentType: 'SALES_INVOICE', documentId: invoice.id, amount: '300.00' }] });
  await c.api.post('/api/payments', { direction: 'RECEIPT', customerId: B, paymentDate: c.date, methodId: cash, amount: '50.00', allocations: [] });
  const cats = (await c.api.get('/api/expense-categories')).json().data;
  await c.api.post('/api/expenses', { expenseDate: c.date, paymentType: 'CASH', methodId: cash, post: true,
    lines: [{ categoryId: cats.find((x: { code: string }) => x.code === 'MARKETING').id, amount: '200' }] });
  await c.api.post('/api/payments', { direction: 'DISBURSEMENT', supplierId: S, paymentDate: c.date, methodId: await c.method('BANK'), amount: '400.00',
    allocations: [{ documentType: 'PURCHASE_INVOICE', documentId: bill.id, amount: '400.00' }] });
  const wrong = await c.postInvoice(B, [{ productId: service, quantity: '3' }]);
  await c.api.post(`/api/invoices/${wrong.id}/cancel`, { reason: 'خطأ في العميل' });
});
afterAll(async () => { await t.close(); });

describe('financial statements', () => {
  it('builds the income statement from the ledger', async () => {
    const pl = await get(`/api/reports/profit-loss?${range()}`);
    expect(pl).toMatchObject({ revenue: '400.00', costOfSales: '180.00', grossProfit: '220.00', expenses: '200.00', netProfit: '20.00' });
    expect(pl.sections[0].accounts.map((a: { code: string; amount: string }) => [a.code, a.amount])).toEqual([['4100', '400.00']]);
    expect(pl.sections[2].accounts.map((a: { code: string; amount: string }) => [a.code, a.amount])).toEqual([['6300', '200.00']]);
  });

  it('balances the balance sheet with current earnings', async () => {
    const bs = await get(`/api/reports/balance-sheet?asOf=${y.to}`);
    const amounts = (s: { accounts: { code: string; amount: string }[] }) => Object.fromEntries(s.accounts.map((a) => [a.code, a.amount]));
    expect(amounts(bs.assets)).toEqual({ 1100: '10120.00', 1200: '4600.00', 1300: '110.00', 1400: '420.00' });
    expect(amounts(bs.liabilities)).toEqual({ 2100: '290.00', 2210: '60.00', 2220: '-120.00' });
    expect(bs.equity).toMatchObject({ currentEarnings: '20.00', total: '15020.00' });
    expect(bs).toMatchObject({ totalAssets: '15250.00', totalLiabilitiesAndEquity: '15250.00', balanced: true });
    const before = await get(`/api/reports/balance-sheet?asOf=${y.from}`);
    expect(before).toMatchObject({ totalAssets: '0.00', balanced: true });
  });

  it('explains every riyal of cash movement (direct method)', async () => {
    const cf = await get(`/api/reports/cash-flow?${range()}`);
    expect(cf).toMatchObject({ openingCash: '0.00', netChange: '14720.00', closingCash: '14720.00', reconciled: true });
    const [operating, investing, financing] = cf.sections;
    expect(financing.total).toBe('15000.00');
    expect(investing.total).toBe('0.00');
    expect(operating.total).toBe('-280.00');
    expect(operating.lines.map((l: { code: string; net: string }) => [l.code, l.net])).toEqual([['1300', '350.00'], ['2100', '-400.00'], ['2220', '-30.00'], ['6300', '-200.00']]);
  });

  it('lists an account with opening and running balances', async () => {
    const gl = await get(`/api/reports/general-ledger?${range()}&accountId=${ch.byKey.get('CASH')}`);
    expect(gl.openingBalance).toBe('0.00');
    expect(gl.lines.map((l: { debit: string; credit: string; balance: string }) => [l.debit, l.credit, l.balance])).toEqual([
      ['10000.00', '0.00', '10000.00'], ['300.00', '0.00', '10300.00'], ['50.00', '0.00', '10350.00'], ['0.00', '230.00', '10120.00'],
    ]);
    expect(gl.lines[1].documentNumber).toMatch(/^RCPT-\d{6}$/);
    expect(gl.lines[3].documentNumber).toMatch(/^EXP-\d{6}$/);
    expect(gl).toMatchObject({ closingBalance: '10120.00', totals: { debit: '10350.00', credit: '230.00' }, truncated: false });
    // A later window starts from the closing balance.
    const later = await get(`/api/reports/general-ledger?dateFrom=${y.to}&dateTo=${y.to}&accountId=${ch.byKey.get('CASH')}`);
    expect(later).toMatchObject({ openingBalance: '10120.00', lines: [], closingBalance: '10120.00' });
    // Credit-normal accounts show credit − debit.
    const ap = await get(`/api/reports/general-ledger?${range()}&accountId=${ch.byKey.get('ACCOUNTS_PAYABLE')}`);
    expect(ap.closingBalance).toBe('290.00');
    const header = await c.api.get(`/api/reports/general-ledger?${range()}&accountId=${ch.byCode.get('1000')}`);
    expect(header.json().error.code).toBe('HEADER_ACCOUNT');
  });

  it('lists the journal with balanced totals', async () => {
    const j = await get(`/api/reports/journal?${range()}`);
    expect(j.total).toBe(j.data.length);
    expect(j.totals.debit).toBe(j.totals.credit);
    expect(j.data.every((e: { totalDebit: string; lines: { debit: string }[] }) => e.lines.length >= 2)).toBe(true);
    const manual = await get(`/api/reports/journal?${range()}&referenceType=MANUAL`);
    expect(manual.total).toBe(1);
  });

  it('rejects an inverted range', async () => {
    expect((await c.api.get(`/api/reports/profit-loss?dateFrom=${y.to}&dateTo=${y.from}`)).json().error.code).toBe('INVALID_RANGE');
  });
});

describe('receivables, payables and statements', () => {
  it('ages receivables and reconciles each customer to the ledger', async () => {
    const r = await get('/api/reports/receivables-aging');
    const a = r.data.find((x: { partyId: string }) => x.partyId === A);
    const b = r.data.find((x: { partyId: string }) => x.partyId === B);
    expect(a.documents).toHaveLength(1);
    expect(a.documents[0]).toMatchObject({ number: invoice.number, remaining: '160.00' });
    const days = Math.round((Date.parse(r.asOf) - Date.parse(c.date)) / 86_400_000);
    const bucket = days <= 0 ? 'current' : days <= 30 ? 'days1to30' : days <= 60 ? 'days31to60' : days <= 90 ? 'days61to90' : 'over90';
    expect(a[bucket]).toBe('160.00');
    expect(a).toMatchObject({ unapplied: '0.00', balance: '160.00' });
    // B paid in advance and its invoice was cancelled: only the credit remains.
    expect(b).toMatchObject({ documents: [], unapplied: '-50.00', balance: '-50.00' });
    expect(r.totals.balance).toBe('110.00'); // equals AR in the balance sheet
  });

  it('ages payables', async () => {
    const r = await get('/api/reports/payables-aging');
    expect(r.data).toHaveLength(1);
    expect(r.data[0]).toMatchObject({ partyId: S, unapplied: '0.00', balance: '290.00' });
    expect(r.totals.balance).toBe('290.00');
  });

  it('produces customer and supplier statements with running balances', async () => {
    const st = await get(`/api/reports/customer-statement?${range()}&partyId=${A}`);
    expect(st.lines.map((l: { referenceType: string; debit: string; credit: string; balance: string }) => [l.referenceType, l.debit, l.credit, l.balance])).toEqual([
      ['SALES_INVOICE', '575.00', '0.00', '575.00'], ['SALES_RETURN', '0.00', '115.00', '460.00'], ['PAYMENT_RECEIPT', '0.00', '300.00', '160.00'],
    ]);
    expect(st.lines[0].documentNumber).toBe(invoice.number);
    expect(st.closingBalance).toBe('160.00');
    // The cancelled invoice and its reversal both appear and cancel out.
    const sb = await get(`/api/reports/customer-statement?${range()}&partyId=${B}`);
    expect(sb.lines.map((l: { isReversal: boolean; debit: string; credit: string }) => [l.isReversal, l.debit, l.credit])).toEqual([
      [false, '0.00', '50.00'], [false, '345.00', '0.00'], [true, '0.00', '345.00'],
    ]);
    expect(sb.closingBalance).toBe('-50.00');
    const ss = await get(`/api/reports/supplier-statement?${range()}&partyId=${S}`);
    expect(ss).toMatchObject({ openingBalance: '0.00', closingBalance: '290.00' });
    expect((await c.api.get(`/api/reports/supplier-statement?${range()}&partyId=${A}`)).statusCode).toBe(404);
  });
});

describe('operational reports', () => {
  it('reports sales by product with cost and gross profit, net of returns and without cancelled invoices', async () => {
    const r = await get(`/api/reports/sales?${range()}&groupBy=product`);
    const row = (id: string) => r.data.find((x: { key: string }) => x.key === id);
    expect(row(goods)).toMatchObject({ quantity: '3.0000', netAmount: '300.00', cost: '180.00', grossProfit: '120.00', documents: 1, returns: 1 });
    expect(row(service)).toMatchObject({ quantity: '1.0000', netAmount: '100.00', cost: '0.00', documents: 1 });
    // Agrees with the income statement.
    expect(r.totals).toMatchObject({ netAmount: '400.00', vatAmount: '60.00', totalAmount: '460.00', cost: '180.00', grossProfit: '220.00' });
    const byCustomer = await get(`/api/reports/sales?${range()}&groupBy=party`);
    expect(byCustomer.data.map((x: { key: string; totalAmount: string }) => [x.key, x.totalAmount])).toEqual([[A, '460.00']]);
    const byMonth = await get(`/api/reports/sales?${range()}&groupBy=month`);
    expect(byMonth.data).toEqual([expect.objectContaining({ key: c.date.slice(0, 7), netAmount: '400.00' })]);
    const filtered = await get(`/api/reports/sales?${range()}&groupBy=document&productId=${service}`);
    expect(filtered.totals.netAmount).toBe('100.00');
  });

  it('reports purchases by supplier', async () => {
    const r = await get(`/api/reports/purchases?${range()}`);
    expect(r.data).toEqual([expect.objectContaining({ key: S, netAmount: '600.00', vatAmount: '90.00', totalAmount: '690.00' })]);
  });

  it('reports expenses from the ledger', async () => {
    const r = await get(`/api/reports/expenses?${range()}`);
    expect(r.data.map((x: { code: string; amount: string }) => [x.code, x.amount])).toEqual([['6300', '200.00']]);
    expect(r.total).toBe('200.00');
  });
});

describe('dashboard', () => {
  it('shows ledger-based figures and a monthly series', async () => {
    const d = await get(`/api/dashboard?${range()}`);
    expect(d.kpis).toEqual({
      totalSales: '400.00', totalPurchases: '600.00', costOfSales: '180.00', expenses: '200.00', netProfit: '20.00',
      receivables: '110.00', payables: '290.00', cash: '10120.00', bank: '4600.00', inventoryValue: '420.00', vatPayable: '-60.00',
    });
    expect(d.monthly).toHaveLength(12);
    const m = d.monthly.find((x: { month: string }) => x.month === c.date.slice(0, 7));
    expect(m).toMatchObject({ revenue: '400.00', costs: '380.00', netProfit: '20.00', cashBalance: '14720.00' });
    expect(d.monthly.at(-1).cashBalance).toBe('14720.00');
  });

  it('defaults to the fiscal year to date', async () => {
    const d = await get('/api/dashboard');
    expect(d.dateFrom <= d.dateTo).toBe(true);
  });
});

describe('year closing', () => {
  it('keeps the closed year\'s profit in the income statement and moves it to retained earnings', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    const k = await chart(t.app, s.token);
    const [year] = await fiscalYears(t.app, s.token);
    const date = year!.periods[1]!.startDate;
    await api.post('/api/journal-entries', { entryDate: date, description: 'بيع نقدي', post: true, lines: [dr(k.byKey.get('CASH')!, '1000.00'), cr(k.byKey.get('SALES')!, '1000.00')] });
    expect((await api.post(`/api/fiscal-years/${year!.id}/close`)).statusCode).toBe(200);
    const pl = (await api.get(`/api/reports/profit-loss?dateFrom=${year!.startDate}&dateTo=${year!.endDate}`)).json();
    expect(pl.netProfit).toBe('1000.00');
    const bs = (await api.get(`/api/reports/balance-sheet?asOf=${year!.endDate}`)).json();
    expect(bs.equity).toMatchObject({ currentEarnings: '0.00', total: '1000.00' });
    expect(bs.equity.accounts.map((a: { code: string; amount: string }) => [a.code, a.amount])).toEqual([['3200', '1000.00']]);
    expect(bs.balanced).toBe(true);
    const d = (await api.get(`/api/dashboard?dateFrom=${year!.startDate}&dateTo=${year!.endDate}`)).json();
    expect(d.kpis.netProfit).toBe('1000.00');
  });
});

describe('permissions and isolation', () => {
  it('separates financial statements from operational reports', async () => {
    const sales = await addMember(t.app, owner, 'SALES_MANAGER');
    const sApi = client(t.app, sales.token);
    expect((await sApi.get(`/api/reports/sales?${range()}`)).statusCode).toBe(200);
    expect((await sApi.get('/api/reports/receivables-aging')).statusCode).toBe(200);
    for (const path of [`/api/reports/profit-loss?${range()}`, `/api/reports/balance-sheet?asOf=${y.to}`, '/api/dashboard',
      `/api/reports/trial-balance?${range()}`, `/api/reports/general-ledger?${range()}&accountId=${ch.byKey.get('CASH')}`]) {
      expect((await sApi.get(path)).statusCode).toBe(403);
    }
    const accountant = await addMember(t.app, owner, 'ACCOUNTANT');
    expect((await client(t.app, accountant.token).get(`/api/reports/profit-loss?${range()}`)).statusCode).toBe(200);
    const viewer = await addMember(t.app, owner, 'VIEWER');
    expect((await client(t.app, viewer.token).get('/api/dashboard')).statusCode).toBe(200);
    const clerk = await addMember(t.app, owner, 'SALES_EMPLOYEE');
    expect((await client(t.app, clerk.token).get(`/api/reports/sales?${range()}`)).statusCode).toBe(403);
  });

  it('never shows another tenant\'s figures', async () => {
    const other = await register(t.app);
    const o = client(t.app, other.token);
    expect((await o.get(`/api/reports/general-ledger?${range()}&accountId=${ch.byKey.get('CASH')}`)).statusCode).toBe(404);
    expect((await o.get(`/api/reports/customer-statement?${range()}&partyId=${A}`)).statusCode).toBe(404);
    expect((await o.get(`/api/reports/sales?${range()}&partyId=${A}`)).json().data).toEqual([]);
    const pl = (await o.get(`/api/reports/profit-loss?${range()}`)).json();
    expect(pl.netProfit).toBe('0.00');
    expect((await o.get('/api/reports/receivables-aging')).json().data).toEqual([]);
    expect((await o.get('/api/dashboard')).json().kpis.cash).toBe('0.00');
  });
});
