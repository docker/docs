import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { balanceOf, commerce, type Commerce } from './commerce-helpers.js';
import { addMember, chart, client, fiscalYears, register, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;
let c: Commerce;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
  c = await commerce(t.app, owner);
});
afterAll(async () => { await t.close(); });

describe('purchase invoices', () => {
  it('posts goods to inventory: Dr Inventory + Dr VAT input / Cr AP', async () => {
    const sup = await c.supplier({ vatNumber: '310000000000003' });
    const prod = await c.goods({ purchasePrice: '60' });
    const bill = await c.purchase(sup.id, [{ productId: prod.id, quantity: '10' }], { supplierInvoiceNumber: 'S-9001' });
    expect(bill).toMatchObject({ status: 'POSTED', taxableAmount: '600.00', taxAmount: '90.00', total: '690.00', supplierInvoiceNumber: 'S-9001' });
    expect(bill.number).toMatch(/^PINV-\d{6}$/);
    const lines = await c.journal(bill.journalEntryId);
    expect(lines.map((l) => [l.accountCode, l.debit, l.credit])).toEqual([
      ['2100', '0.00', '690.00'], ['1400', '600.00', '0.00'], ['2220', '90.00', '0.00'],
    ]);
    expect(lines[0]!.supplierId).toBe(sup.id);
    expect(await balanceOf(c.api, 'suppliers', sup.id)).toBe('-690.00');
  });

  it('records the supplier invoice number once per supplier', async () => {
    const sup = await c.supplier();
    const prod = await c.goods();
    await c.purchase(sup.id, [{ productId: prod.id, quantity: '1' }], { supplierInvoiceNumber: 'DUP-1' });
    const dup = await c.api.post('/api/purchase-invoices', { partyId: sup.id, docDate: c.date, supplierInvoiceNumber: 'DUP-1', lines: [{ productId: prod.id, quantity: '1' }] });
    expect(dup.statusCode).toBe(409);
  });

  it('posts expense lines to their account and requires one for untracked services', async () => {
    const sup = await c.supplier();
    const ch = await chart(t.app, owner.token);
    const service = await c.product({ purchasePrice: '1000' });
    const draft = (await c.api.post('/api/purchase-invoices', { partyId: sup.id, docDate: c.date, lines: [{ productId: service.id, quantity: '1' }] })).json();
    expect((await c.api.post(`/api/purchase-invoices/${draft.id}/post`)).json().error.code).toBe('ACCOUNT_REQUIRED');

    const rent = await c.purchase(sup.id, [{ accountId: ch.byCode.get('6100'), description: 'إيجار مارس', quantity: '1', unitPrice: '5000' }]);
    expect((await c.journal(rent.journalEntryId)).map((l) => [l.accountCode, l.debit, l.credit])).toEqual([
      ['2100', '0.00', '5750.00'], ['6100', '5000.00', '0.00'], ['2220', '750.00', '0.00'],
    ]);
    // Exempt line: no VAT.
    const exempt = await c.purchase(sup.id, [{ accountId: ch.byCode.get('6400'), quantity: '1', unitPrice: '300', vatCategory: 'E' }]);
    expect(exempt).toMatchObject({ taxAmount: '0.00', total: '300.00' });
    // Sales lines cannot use bare accounts.
    const cust = await c.customer();
    expect((await c.api.post('/api/invoices', { partyId: cust.id, docDate: c.date, lines: [{ accountId: ch.byCode.get('6100'), quantity: '1', unitPrice: '1' }] })).statusCode).toBe(400);
  });

  it('returns goods to the supplier: Dr AP / Cr Inventory + Cr VAT input', async () => {
    const sup = await c.supplier();
    const bill = await c.purchase(sup.id, [{ productId: (await c.goods()).id, quantity: '4', unitPrice: '50' }]); // 230
    const r = (await c.api.post('/api/purchase-returns', { originalInvoiceId: bill.id, docDate: c.date, reason: 'معيب', lines: [{ sourceItemId: bill.lines[0].id, quantity: '1' }] })).json();
    const ret = (await c.api.post(`/api/purchase-returns/${r.id}/post`)).json();
    expect(ret.number).toMatch(/^DN-\d{6}$/);
    expect((await c.journal(ret.journalEntryId)).map((l) => [l.accountCode, l.debit, l.credit])).toEqual([
      ['2100', '57.50', '0.00'], ['1400', '0.00', '50.00'], ['2220', '0.00', '7.50'],
    ]);
    expect((await c.api.get(`/api/purchase-invoices/${bill.id}`)).json()).toMatchObject({ returnedAmount: '57.50', remainingAmount: '172.50' });
    expect(await balanceOf(c.api, 'suppliers', sup.id)).toBe('-172.50');
  });

  it('cancels an unpaid purchase invoice', async () => {
    const sup = await c.supplier();
    const bill = await c.purchase(sup.id, [{ productId: (await c.goods()).id, quantity: '1' }]);
    expect((await c.api.post(`/api/purchase-invoices/${bill.id}/cancel`, { reason: 'مكرر' })).json().status).toBe('CANCELLED');
    expect(await balanceOf(c.api, 'suppliers', sup.id)).toBe('0.00');
  });
});

describe('purchase orders', () => {
  it('approves an order and converts it to a draft bill', async () => {
    const sup = await c.supplier();
    const po = (await c.api.post('/api/purchase-orders', { partyId: sup.id, docDate: c.date, lines: [{ productId: (await c.goods()).id, quantity: '12' }] })).json();
    expect(po.number).toMatch(/^PO-\d{6}$/);
    expect((await c.api.post(`/api/purchase-orders/${po.id}/convert`, {})).statusCode).toBe(409);
    expect((await c.api.post(`/api/purchase-orders/${po.id}/status`, { status: 'APPROVED' })).json().status).toBe('APPROVED');
    const bill = (await c.api.post(`/api/purchase-orders/${po.id}/convert`, { docDate: c.date })).json();
    expect(bill).toMatchObject({ status: 'DRAFT', sourceOrderId: po.id, total: po.total });
    expect((await c.api.post(`/api/purchase-invoices/${bill.id}/post`)).json().status).toBe('POSTED');
  });
});

describe('purchase permissions', () => {
  it('lets purchase managers post bills and sales staff see nothing', async () => {
    const sup = await c.supplier();
    const prod = await c.goods();
    const pm = await addMember(t.app, owner, 'PURCHASE_MANAGER');
    const pApi = client(t.app, pm.token);
    const draft = await pApi.post('/api/purchase-invoices', { partyId: sup.id, docDate: c.date, lines: [{ productId: prod.id, quantity: '1' }] });
    expect((await pApi.post(`/api/purchase-invoices/${draft.json().id}/post`)).statusCode).toBe(200);
    const sales = await addMember(t.app, owner, 'SALES_EMPLOYEE');
    expect((await client(t.app, sales.token).get('/api/purchase-invoices')).statusCode).toBe(403);
  });

  it('keeps the ledger balanced', async () => {
    const [year] = await fiscalYears(t.app, owner.token);
    expect((await c.api.get(`/api/reports/trial-balance?dateFrom=${year!.startDate}&dateTo=${year!.endDate}`)).json().balanced).toBe(true);
  });
});
