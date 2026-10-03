import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { withTx } from '../src/db/tx.js';
import { balanceOf, commerce, type Commerce } from './commerce-helpers.js';
import { addMember, client, fiscalYears, register, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;
let c: Commerce;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
  c = await commerce(t.app, owner);
});
afterAll(async () => { await t.close(); });

const address = { buildingNumber: '1234', street: 'طريق الملك فهد', district: 'العليا', city: 'الرياض', postalCode: '12345', additionalNumber: '6789' };

describe('sales invoice posting', () => {
  it('posts the reference example: 1000 + 150 VAT → Dr AR 1150 / Cr Sales 1000 / Cr VAT 150', async () => {
    const customer = await c.customer({ vatNumber: '300000000000003', paymentTermsDays: 30, addresses: [address] });
    const product = await c.product({ salePrice: '1000' });
    const draft = await c.invoice(customer.id, [{ productId: product.id, quantity: '1' }]);
    expect(draft).toMatchObject({ status: 'DRAFT', number: null, subtotal: '1000.00', taxableAmount: '1000.00', taxAmount: '150.00', total: '1150.00' });

    const inv = (await c.api.post(`/api/invoices/${draft.id}/post`)).json();
    expect(inv).toMatchObject({ status: 'ISSUED', total: '1150.00', paidAmount: '0.00', remainingAmount: '1150.00', invoiceKind: 'STANDARD' });
    expect(inv.number).toMatch(/^INV-\d{6}$/);
    expect(inv.partySnapshot).toMatchObject({ vatNumber: '300000000000003', address: expect.objectContaining({ postalCode: '12345' }) });
    const due = new Date(`${c.date}T00:00:00Z`); due.setUTCDate(due.getUTCDate() + 30);
    expect(inv.dueDate).toBe(due.toISOString().slice(0, 10));

    const lines = await c.journal(inv.journalEntryId);
    expect(lines.map((l) => [l.accountCode, l.debit, l.credit])).toEqual([
      ['1300', '1150.00', '0.00'], ['4100', '0.00', '1000.00'], ['2210', '0.00', '150.00'],
    ]);
    expect(lines[0]!.customerId).toBe(customer.id);
    const entry = (await c.api.get(`/api/journal-entries/${inv.journalEntryId}`)).json();
    expect(entry).toMatchObject({ status: 'POSTED', source: 'SYSTEM', referenceType: 'SALES_INVOICE', referenceId: inv.id, totalDebit: '1150.00', totalCredit: '1150.00' });
    expect(await balanceOf(c.api, 'customers', customer.id)).toBe('1150.00');
  });

  it('marks invoices to customers without a VAT number as simplified (B2C)', async () => {
    const inv = await c.postInvoice((await c.customer()).id, [{ productId: (await c.product()).id, quantity: '1' }]);
    expect(inv.invoiceKind).toBe('SIMPLIFIED');
  });

  it('numbers issued invoices sequentially; drafts have no number', async () => {
    const cust = await c.customer();
    const prod = await c.product();
    const a = await c.postInvoice(cust.id, [{ productId: prod.id, quantity: '1' }]);
    const draft = await c.invoice(cust.id, [{ productId: prod.id, quantity: '1' }]);
    const b = await c.postInvoice(cust.id, [{ productId: prod.id, quantity: '1' }]);
    expect(draft.number).toBeNull();
    expect(Number(b.number.slice(-6))).toBe(Number(a.number.slice(-6)) + 1);
  });

  it('calculates on the server without saving', async () => {
    const cust = await c.customer();
    const prod = await c.product({ salePrice: '200' });
    const before = (await c.api.get('/api/invoices')).json().total;
    const res = await c.api.post('/api/invoices/calculate', { partyId: cust.id, docDate: c.date, lines: [{ productId: prod.id, quantity: '3', discountPercent: '10' }] });
    expect(res.json().totals).toMatchObject({ subtotal: '600.00', discountTotal: '60.00', taxableAmount: '540.00', taxAmount: '81.00', total: '621.00' });
    expect((await c.api.get('/api/invoices')).json().total).toBe(before);
  });

  it('converts VAT-inclusive list prices to the document basis', async () => {
    const prod = await c.product({ salePrice: '115', salePriceIncludesVat: true });
    const inv = await c.invoice((await c.customer()).id, [{ productId: prod.id, quantity: '2' }]);
    expect(inv.lines[0]).toMatchObject({ unitPrice: '100.0000', netAmount: '200.00' });
    expect(inv.total).toBe('230.00');
    const incl = await c.invoice((await c.customer()).id, [{ productId: prod.id, quantity: '2' }], { pricesIncludeVat: true });
    expect(incl).toMatchObject({ taxableAmount: '200.00', taxAmount: '30.00', total: '230.00' });
  });

  it('edits and deletes drafts, then refuses changes once issued', async () => {
    const cust = await c.customer();
    const prod = await c.product();
    const draft = await c.invoice(cust.id, [{ productId: prod.id, quantity: '1' }]);
    const edited = (await c.api.patch(`/api/invoices/${draft.id}`, { lines: [{ productId: prod.id, quantity: '5', unitPrice: '10' }] })).json();
    expect(edited.total).toBe('57.50');
    const other = await c.invoice(cust.id, [{ productId: prod.id, quantity: '1' }]);
    expect((await c.api.del(`/api/invoices/${other.id}`)).statusCode).toBe(204);

    const issued = (await c.api.post(`/api/invoices/${draft.id}/post`)).json();
    expect((await c.api.patch(`/api/invoices/${issued.id}`, { notes: 'x' })).statusCode).toBe(409);
    expect((await c.api.del(`/api/invoices/${issued.id}`)).statusCode).toBe(409);
    expect((await c.api.post(`/api/invoices/${issued.id}/post`)).statusCode).toBe(409);
    const ctx = { tenantId: owner.tenantId, userId: owner.userId };
    await expect(withTx(t.pool, ctx, (db) => db.query(`UPDATE sales_invoices SET total = 1 WHERE id = $1`, [issued.id]))).rejects.toThrow(/cannot be modified/);
    await expect(withTx(t.pool, ctx, (db) => db.query(`UPDATE sales_invoice_items SET quantity = 1 WHERE document_id = $1`, [issued.id]))).rejects.toThrow(/cannot be changed/);
    await expect(withTx(t.pool, ctx, (db) => db.query(`UPDATE sales_invoices SET status = 'DRAFT' WHERE id = $1`, [issued.id]))).rejects.toThrow(/draft/);
  });

  it('enforces the credit limit from the ledger balance', async () => {
    const cust = await c.customer({ creditLimit: '1000.00' });
    const prod = await c.product({ salePrice: '500' });
    await c.postInvoice(cust.id, [{ productId: prod.id, quantity: '1' }]); // 575
    const second = await c.invoice(cust.id, [{ productId: prod.id, quantity: '1' }]);
    const res = await c.api.post(`/api/invoices/${second.id}/post`);
    expect(res.json().error.code).toBe('CREDIT_LIMIT_EXCEEDED');
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe('575.00');
  });

  it('uses the VAT rate effective on the invoice date', async () => {
    const s = await register(t.app);
    const sc = await commerce(t.app, s);
    const [year] = await fiscalYears(t.app, s.token);
    const july = year!.periods[6]!.startDate;
    expect((await sc.api.post('/api/tax-rates', { nameAr: 'نسبة جديدة', rate: '0.2', effectiveFrom: july })).statusCode).toBe(201);
    const cust = await sc.customer();
    const prod = await sc.product({ salePrice: '100' });
    const before = await sc.invoice(cust.id, [{ productId: prod.id, quantity: '1' }], { docDate: year!.periods[5]!.startDate });
    const after = await sc.invoice(cust.id, [{ productId: prod.id, quantity: '1' }], { docDate: july });
    expect(before.taxAmount).toBe('15.00');
    expect(after.taxAmount).toBe('20.00');
    const rates = (await sc.api.get('/api/tax-rates')).json().data;
    expect(rates.find((r: { rate: string }) => r.rate === '0.1500').effectiveTo).toBe(year!.periods[5]!.endDate);
  });
});

describe('cancellation', () => {
  it('cancels an unpaid invoice by reversing its entry', async () => {
    const cust = await c.customer();
    const inv = await c.postInvoice(cust.id, [{ productId: (await c.product()).id, quantity: '2' }]);
    const res = await c.api.post(`/api/invoices/${inv.id}/cancel`, { reason: 'خطأ في العميل' });
    expect(res.json()).toMatchObject({ status: 'CANCELLED', cancelReason: 'خطأ في العميل', number: inv.number });
    expect((await c.api.get(`/api/journal-entries/${inv.journalEntryId}`)).json().status).toBe('REVERSED');
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe('0.00');
    expect((await c.api.post(`/api/invoices/${inv.id}/cancel`, { reason: 'مرة ثانية' })).statusCode).toBe(409);
  });

  it('refuses to cancel a paid invoice', async () => {
    const cust = await c.customer();
    const inv = await c.postInvoice(cust.id, [{ productId: (await c.product()).id, quantity: '1' }]);
    await c.api.post('/api/payments', {
      direction: 'RECEIPT', customerId: cust.id, paymentDate: c.date, methodId: await c.method('CASH'), amount: '50.00',
      allocations: [{ documentType: 'SALES_INVOICE', documentId: inv.id, amount: '50.00' }],
    });
    expect((await c.api.post(`/api/invoices/${inv.id}/cancel`, { reason: 'إلغاء' })).json().error.code).toBe('HAS_PAYMENTS');
  });
});

describe('sales returns (credit notes)', () => {
  it('returns part of an invoice, then the rest, and nets exactly to zero', async () => {
    const cust = await c.customer();
    const prod = await c.product({ salePrice: '33.33' });
    const inv = await c.postInvoice(cust.id, [{ productId: prod.id, quantity: '3' }]);
    expect(inv).toMatchObject({ taxableAmount: '99.99', taxAmount: '15.00', total: '114.99' });
    const itemId = inv.lines[0].id;

    const r1 = (await c.api.post('/api/sales-returns', { originalInvoiceId: inv.id, docDate: c.date, reason: 'تالف', lines: [{ sourceItemId: itemId, quantity: '1' }] })).json();
    expect(r1).toMatchObject({ status: 'DRAFT', taxableAmount: '33.33', taxAmount: '5.00', total: '38.33' });
    const p1 = (await c.api.post(`/api/sales-returns/${r1.id}/post`)).json();
    expect(p1.number).toMatch(/^CN-\d{6}$/);
    expect(p1.appliedAmount).toBe('38.33');
    const lines = await c.journal(p1.journalEntryId);
    expect(lines.map((l) => [l.accountCode, l.debit, l.credit])).toEqual([
      ['1300', '0.00', '38.33'], ['4100', '33.33', '0.00'], ['2210', '5.00', '0.00'],
    ]);
    const mid = (await c.api.get(`/api/invoices/${inv.id}`)).json();
    expect(mid).toMatchObject({ status: 'ISSUED', returnedAmount: '38.33', remainingAmount: '76.66' });

    const over = await c.api.post('/api/sales-returns', { originalInvoiceId: inv.id, docDate: c.date, reason: 'زيادة', lines: [{ sourceItemId: itemId, quantity: '3' }] });
    expect(over.json().error.code).toBe('RETURN_EXCEEDS_INVOICE');

    const r2 = (await c.api.post('/api/sales-returns', { originalInvoiceId: inv.id, docDate: c.date, reason: 'الباقي', lines: [{ sourceItemId: itemId, quantity: '2' }] })).json();
    const p2 = (await c.api.post(`/api/sales-returns/${r2.id}/post`)).json();
    expect(p2).toMatchObject({ taxableAmount: '66.66', taxAmount: '10.00', total: '76.66' });
    const final = (await c.api.get(`/api/invoices/${inv.id}`)).json();
    expect(final).toMatchObject({ status: 'RETURNED', returnedAmount: '114.99', remainingAmount: '0.00' });
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe('0.00');
  });

  it('cancelling a return restores the invoice', async () => {
    const cust = await c.customer();
    const inv = await c.postInvoice(cust.id, [{ productId: (await c.product()).id, quantity: '2' }]);
    const r = (await c.api.post('/api/sales-returns', { originalInvoiceId: inv.id, docDate: c.date, reason: 'خطأ', lines: [{ sourceItemId: inv.lines[0].id, quantity: '2' }] })).json();
    await c.api.post(`/api/sales-returns/${r.id}/post`);
    expect((await c.api.get(`/api/invoices/${inv.id}`)).json().status).toBe('RETURNED');
    expect((await c.api.post(`/api/invoices/${inv.id}/cancel`, { reason: 'إلغاء' })).json().error.code).toBe('HAS_RETURNS');
    await c.api.post(`/api/sales-returns/${r.id}/cancel`, { reason: 'أُدخل خطأ' });
    expect((await c.api.get(`/api/invoices/${inv.id}`)).json()).toMatchObject({ status: 'ISSUED', returnedAmount: '0.00', remainingAmount: inv.total });
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe(inv.total);
  });

  it('refuses returns of drafts and of lines from another invoice', async () => {
    const cust = await c.customer();
    const prod = await c.product();
    const draft = await c.invoice(cust.id, [{ productId: prod.id, quantity: '1' }]);
    expect((await c.api.post('/api/sales-returns', { originalInvoiceId: draft.id, docDate: c.date, reason: 'مسودة', lines: [{ sourceItemId: draft.lines[0].id, quantity: '1' }] })).json().error.code).toBe('INVOICE_NOT_ISSUED');
    const a = await c.postInvoice(cust.id, [{ productId: prod.id, quantity: '1' }]);
    const b = await c.postInvoice(cust.id, [{ productId: prod.id, quantity: '1' }]);
    expect((await c.api.post('/api/sales-returns', { originalInvoiceId: a.id, docDate: c.date, reason: 'خلط', lines: [{ sourceItemId: b.lines[0].id, quantity: '1' }] })).statusCode).toBe(400);
  });
});

describe('quotes', () => {
  it('sends, accepts and converts a quote into a draft invoice', async () => {
    const cust = await c.customer();
    const prod = await c.product({ salePrice: '80' });
    const q = (await c.api.post('/api/sales-quotes', { partyId: cust.id, docDate: c.date, lines: [{ productId: prod.id, quantity: '5', discountAmount: '20' }] })).json();
    expect(q.number).toMatch(/^QT-\d{6}$/);
    expect((await c.api.post(`/api/sales-quotes/${q.id}/status`, { status: 'SENT' })).json().status).toBe('SENT');
    expect((await c.api.post(`/api/sales-quotes/${q.id}/status`, { status: 'ACCEPTED' })).json().status).toBe('ACCEPTED');
    const inv = (await c.api.post(`/api/sales-quotes/${q.id}/convert`, { docDate: c.date })).json();
    expect(inv).toMatchObject({ status: 'DRAFT', sourceQuoteId: q.id, total: q.total });
    expect((await c.api.get(`/api/sales-quotes/${q.id}`)).json()).toMatchObject({ status: 'CONVERTED', convertedInvoiceId: inv.id });
    expect((await c.api.post(`/api/sales-quotes/${q.id}/convert`, {})).statusCode).toBe(409);
    expect((await c.api.post(`/api/sales-quotes/${q.id}/post`)).statusCode).toBe(404);
  });
});

describe('permissions and isolation', () => {
  it('lets sales employees draft but not issue invoices', async () => {
    const cust = await c.customer();
    const prod = await c.product();
    const emp = await addMember(t.app, owner, 'SALES_EMPLOYEE');
    const eApi = client(t.app, emp.token);
    const draft = await eApi.post('/api/invoices', { partyId: cust.id, docDate: c.date, lines: [{ productId: prod.id, quantity: '1' }] });
    expect(draft.statusCode).toBe(201);
    expect((await eApi.post(`/api/invoices/${draft.json().id}/post`)).statusCode).toBe(403);
    const purchase = await addMember(t.app, owner, 'PURCHASE_MANAGER');
    expect((await client(t.app, purchase.token).get('/api/invoices')).statusCode).toBe(403);
  });

  it('hides invoices and products of other tenants', async () => {
    const inv = await c.postInvoice((await c.customer()).id, [{ productId: (await c.product()).id, quantity: '1' }]);
    const other = await register(t.app);
    const oc = await commerce(t.app, other);
    expect((await oc.api.get(`/api/invoices/${inv.id}`)).statusCode).toBe(404);
    expect((await oc.api.post(`/api/invoices/${inv.id}/cancel`, { reason: 'اختراق' })).statusCode).toBe(404);
    const foreign = await oc.api.post('/api/invoices', { partyId: (await oc.customer()).id, docDate: oc.date, lines: [{ productId: inv.lines[0].productId, quantity: '1' }] });
    expect(foreign.json().error.code).toBe('INVALID_PRODUCT');
  });

  it('refuses to delete customers and products used in documents', async () => {
    const cust = await c.customer();
    const prod = await c.product();
    await c.invoice(cust.id, [{ productId: prod.id, quantity: '1' }]);
    expect((await c.api.del(`/api/customers/${cust.id}`)).json().error.code).toBe('PARTY_IN_USE');
    expect((await c.api.del(`/api/products/${prod.id}`)).json().error.code).toBe('PRODUCT_IN_USE');
  });

  it('keeps the ledger balanced', async () => {
    const [year] = await fiscalYears(t.app, owner.token);
    const tb = (await c.api.get(`/api/reports/trial-balance?dateFrom=${year!.startDate}&dateTo=${year!.endDate}`)).json();
    expect(tb.balanced).toBe(true);
  });
});
