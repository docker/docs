import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { withTx } from '../src/db/tx.js';
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

const receipt = async (customerId: string, amount: string, allocations: { documentId: string; amount: string }[] = [], method = 'CASH') =>
  c.api.post('/api/payments', {
    direction: 'RECEIPT', customerId, paymentDate: c.date, methodId: await c.method(method), amount,
    allocations: allocations.map((a) => ({ documentType: 'SALES_INVOICE', ...a })),
  });

const newInvoice = async (price = '1000') => {
  const cust = await c.customer();
  const inv = await c.postInvoice(cust.id, [{ productId: (await c.product({ salePrice: price })).id, quantity: '1' }]);
  return { cust, inv };
};

describe('customer receipts', () => {
  it('collects an invoice in full: Dr Cash / Cr AR, invoice PAID', async () => {
    const { cust, inv } = await newInvoice();
    const res = await receipt(cust.id, '1150.00', [{ documentId: inv.id, amount: '1150.00' }]);
    expect(res.statusCode).toBe(201);
    const pay = res.json();
    expect(pay).toMatchObject({ status: 'POSTED', amount: '1150.00', allocatedAmount: '1150.00', unallocatedAmount: '0.00' });
    expect(pay.number).toMatch(/^RCPT-\d{6}$/);
    const lines = await c.journal(pay.journalEntryId);
    expect(lines.map((l) => [l.accountCode, l.debit, l.credit])).toEqual([['1100', '1150.00', '0.00'], ['1300', '0.00', '1150.00']]);
    expect(lines[1]!.customerId).toBe(cust.id);
    expect((await c.api.get(`/api/invoices/${inv.id}`)).json()).toMatchObject({ status: 'PAID', paidAmount: '1150.00', remainingAmount: '0.00' });
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe('0.00');
  });

  it('records partial payments and settles several invoices with one receipt', async () => {
    const cust = await c.customer();
    const prod = await c.product({ salePrice: '100' });
    const a = await c.postInvoice(cust.id, [{ productId: prod.id, quantity: '1' }]); // 115
    const b = await c.postInvoice(cust.id, [{ productId: prod.id, quantity: '2' }]); // 230
    await receipt(cust.id, '50.00', [{ documentId: a.id, amount: '50.00' }]);
    expect((await c.api.get(`/api/invoices/${a.id}`)).json()).toMatchObject({ status: 'PARTIALLY_PAID', remainingAmount: '65.00' });
    const res = await receipt(cust.id, '295.00', [{ documentId: a.id, amount: '65.00' }, { documentId: b.id, amount: '230.00' }], 'BANK');
    expect(res.statusCode).toBe(201);
    expect((await c.api.get(`/api/invoices/${a.id}`)).json().status).toBe('PAID');
    expect((await c.api.get(`/api/invoices/${b.id}`)).json().status).toBe('PAID');
    const open = (await c.api.get(`/api/invoices?partyId=${cust.id}&open=true`)).json();
    expect(open.total).toBe(0);
  });

  it('keeps an advance on account and allocates it later', async () => {
    const cust = await c.customer();
    const adv = (await receipt(cust.id, '500.00')).json();
    expect(adv.unallocatedAmount).toBe('500.00');
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe('-500.00');
    const inv = await c.postInvoice(cust.id, [{ productId: (await c.product({ salePrice: '200' })).id, quantity: '1' }]); // 230
    const res = await c.api.post(`/api/payments/${adv.id}/allocations`, { allocations: [{ documentType: 'SALES_INVOICE', documentId: inv.id, amount: '230.00' }] });
    expect(res.json()).toMatchObject({ allocatedAmount: '230.00', unallocatedAmount: '270.00' });
    expect((await c.api.get(`/api/invoices/${inv.id}`)).json().status).toBe('PAID');
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe('-270.00');
  });

  it('rejects invalid allocations', async () => {
    const { cust, inv } = await newInvoice('100'); // 115
    const other = await newInvoice('100');
    const cases: [unknown[], string, string][] = [
      [[{ documentId: inv.id, amount: '200.00' }], '300.00', 'OVER_ALLOCATED'],             // more than the invoice
      [[{ documentId: inv.id, amount: '100.00' }], '50.00', 'OVER_ALLOCATED'],              // more than the payment
      [[{ documentId: other.inv.id, amount: '10.00' }], '50.00', 'INVALID_ALLOCATION'],     // another customer's invoice
      [[{ documentId: inv.id, amount: '0.00' }], '50.00', 'INVALID_ALLOCATION'],            // zero
    ];
    for (const [allocations, amount, code] of cases) {
      const res = await receipt(cust.id, amount, allocations as never);
      expect(res.json().error?.code, JSON.stringify(allocations)).toBe(code);
    }
    const wrongType = await c.api.post('/api/payments', {
      direction: 'RECEIPT', customerId: cust.id, paymentDate: c.date, methodId: await c.method('CASH'), amount: '10.00',
      allocations: [{ documentType: 'SALES_RETURN', documentId: inv.id, amount: '10.00' }],
    });
    expect(wrongType.json().error.code).toBe('INVALID_ALLOCATION');
    // Nothing was recorded by the failed attempts.
    expect((await c.api.get(`/api/payments?customerId=${cust.id}`)).json().total).toBe(0);
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe('115.00');
  });

  it('voids a payment: reverses the entry and reopens the invoices', async () => {
    const { cust, inv } = await newInvoice('100');
    const pay = (await receipt(cust.id, '115.00', [{ documentId: inv.id, amount: '115.00' }])).json();
    const res = await c.api.post(`/api/payments/${pay.id}/void`, { reason: 'شيك مرتجع' });
    expect(res.json()).toMatchObject({ status: 'VOIDED', allocatedAmount: '0.00', voidReason: 'شيك مرتجع' });
    expect(res.json().allocations[0].reversedAt).not.toBeNull();
    expect((await c.api.get(`/api/invoices/${inv.id}`)).json()).toMatchObject({ status: 'ISSUED', paidAmount: '0.00' });
    expect((await c.api.get(`/api/journal-entries/${pay.journalEntryId}`)).json().status).toBe('REVERSED');
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe('115.00');
    expect((await c.api.post(`/api/payments/${pay.id}/void`, { reason: 'مرة ثانية' })).statusCode).toBe(409);
  });

  it('refunds a customer after a return of a paid invoice', async () => {
    const { cust, inv } = await newInvoice('100');
    await receipt(cust.id, '115.00', [{ documentId: inv.id, amount: '115.00' }]);
    const r = (await c.api.post('/api/sales-returns', { originalInvoiceId: inv.id, docDate: c.date, reason: 'إرجاع', lines: [{ sourceItemId: inv.lines[0].id, quantity: '1' }] })).json();
    const ret = (await c.api.post(`/api/sales-returns/${r.id}/post`)).json();
    // Paid in full, so nothing is applied to the invoice: the full credit is owed to the customer.
    expect(ret).toMatchObject({ appliedAmount: '0.00', remainingAmount: '115.00' });
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe('-115.00');
    const refund = await c.api.post('/api/payments', {
      direction: 'DISBURSEMENT', customerId: cust.id, paymentDate: c.date, methodId: await c.method('CASH'), amount: '115.00',
      allocations: [{ documentType: 'SALES_RETURN', documentId: ret.id, amount: '115.00' }],
    });
    expect(refund.statusCode).toBe(201);
    expect((await c.journal(refund.json().journalEntryId)).map((l) => [l.accountCode, l.debit, l.credit])).toEqual([['1100', '0.00', '115.00'], ['1300', '115.00', '0.00']]);
    expect((await c.api.get(`/api/sales-returns/${ret.id}`)).json()).toMatchObject({ refundedAmount: '115.00', remainingAmount: '0.00' });
    expect(await balanceOf(c.api, 'customers', cust.id)).toBe('0.00');
  });

  it('cannot change a payment directly in the database', async () => {
    const { cust } = await newInvoice('100');
    const pay = (await receipt(cust.id, '10.00')).json();
    await expect(withTx(t.pool, { tenantId: owner.tenantId, userId: owner.userId }, (db) =>
      db.query(`UPDATE payments SET amount = 1 WHERE id = $1`, [pay.id]))).rejects.toThrow(/cannot be modified/);
  });
});

describe('supplier payments', () => {
  it('pays a purchase invoice: Dr AP / Cr Bank', async () => {
    const sup = await c.supplier();
    const bill = await c.purchase(sup.id, [{ productId: (await c.product({ purchasePrice: '400' })).id, quantity: '1' }]); // 460
    const res = await c.api.post('/api/payments', {
      direction: 'DISBURSEMENT', supplierId: sup.id, paymentDate: c.date, methodId: await c.method('BANK'), amount: '460.00',
      allocations: [{ documentType: 'PURCHASE_INVOICE', documentId: bill.id, amount: '460.00' }],
    });
    expect(res.statusCode).toBe(201);
    expect(res.json().number).toMatch(/^PAY-\d{6}$/);
    expect((await c.journal(res.json().journalEntryId)).map((l) => [l.accountCode, l.debit, l.credit])).toEqual([['1200', '0.00', '460.00'], ['2100', '460.00', '0.00']]);
    expect((await c.api.get(`/api/purchase-invoices/${bill.id}`)).json().status).toBe('PAID');
    expect(await balanceOf(c.api, 'suppliers', sup.id)).toBe('0.00');
  });
});

describe('payment methods', () => {
  it('seeds the Saudi methods and accepts custom ones on asset accounts only', async () => {
    const methods = (await c.api.get('/api/payment-methods')).json().data.map((m: { code: string }) => m.code);
    expect(methods).toEqual(['CASH', 'BANK', 'CARD', 'STC_PAY', 'TAMARA']);
    const ch = await chart(t.app, owner.token);
    const clearing = (await c.api.post('/api/accounts', { code: '1250', nameAr: 'مستحقات تمارا', parentId: ch.byCode.get('1000') })).json();
    const ok = await c.api.post('/api/payment-methods', { code: 'TAMARA_CLR', nameAr: 'تمارا (تسوية)', methodType: 'TAMARA', accountId: clearing.id });
    expect(ok.statusCode).toBe(201);
    expect((await c.api.post('/api/payment-methods', { code: 'BAD', nameAr: 'خطأ', methodType: 'OTHER', accountId: ch.byKey.get('SALES') })).statusCode).toBe(400);
    const { cust, inv } = await newInvoice('100');
    const res = await c.api.post('/api/payments', {
      direction: 'RECEIPT', customerId: cust.id, paymentDate: c.date, methodId: ok.json().id, amount: '115.00',
      allocations: [{ documentType: 'SALES_INVOICE', documentId: inv.id, amount: '115.00' }],
    });
    expect((await c.journal(res.json().journalEntryId))[0]!.accountCode).toBe('1250');
  });
});

describe('payment permissions and isolation', () => {
  it('applies payment permissions', async () => {
    const { cust } = await newInvoice('100');
    const sales = await addMember(t.app, owner, 'SALES_MANAGER');
    const sApi = client(t.app, sales.token);
    const pay = await sApi.post('/api/payments', { direction: 'RECEIPT', customerId: cust.id, paymentDate: c.date, methodId: await c.method('CASH'), amount: '10.00' });
    expect(pay.statusCode).toBe(201);
    expect((await sApi.post(`/api/payments/${pay.json().id}/void`, { reason: 'test' })).statusCode).toBe(403);
    const viewer = await addMember(t.app, owner, 'VIEWER');
    expect((await client(t.app, viewer.token).post('/api/payments', { direction: 'RECEIPT', customerId: cust.id, paymentDate: c.date, methodId: await c.method('CASH'), amount: '10.00' })).statusCode).toBe(403);
  });

  it('cannot pay another tenant\'s invoice', async () => {
    const { inv } = await newInvoice('100');
    const other = await register(t.app);
    const oc = await commerce(t.app, other);
    const res = await oc.api.post('/api/payments', {
      direction: 'RECEIPT', customerId: (await oc.customer()).id, paymentDate: oc.date, methodId: await oc.method('CASH'), amount: '10.00',
      allocations: [{ documentType: 'SALES_INVOICE', documentId: inv.id, amount: '10.00' }],
    });
    expect(res.json().error.code).toBe('INVALID_ALLOCATION');
  });

  it('keeps the ledger balanced', async () => {
    const [year] = await fiscalYears(t.app, owner.token);
    expect((await c.api.get(`/api/reports/trial-balance?dateFrom=${year!.startDate}&dateTo=${year!.endDate}`)).json().balanced).toBe(true);
  });
});
