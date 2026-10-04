import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { withTx } from '../src/db/tx.js';
import { commerce, type Commerce } from './commerce-helpers.js';
import { addMember, chart, client, fiscalYears, register, setupApp, type Chart, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;
let c: Commerce;
let ch: Chart;
let main: string;
let second: string;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
  c = await commerce(t.app, owner);
  ch = await chart(t.app, owner.token);
  const company = (await c.api.get('/api/companies')).json().data[0];
  main = (await c.api.get(`/api/companies/${company.id}/warehouses`)).json().data[0].id;
  second = (await c.api.post(`/api/companies/${company.id}/warehouses`, { code: 'JED', name: 'مستودع جدة' })).json().id;
});
afterAll(async () => { await t.close(); });

const balance = async (productId: string, warehouseId = main) => {
  const rows = (await c.api.get(`/api/inventory/balances?productId=${productId}&warehouseId=${warehouseId}&includeZero=true`)).json().data;
  return rows[0] ? { quantity: rows[0].quantity, value: rows[0].value } : { quantity: '0', value: '0' };
};
const reconciled = async () => {
  const v = (await c.api.get('/api/inventory/valuation')).json();
  expect(v.difference, `stock ${v.stockValue} vs ledger ${v.ledgerBalance}`).toBe('0.00');
  return v;
};
const adjustIn = (productId: string, quantity: string, unitCost: string, warehouseId = main, offsetAccountId?: string) =>
  c.api.post('/api/stock-adjustments', { warehouseId, date: c.date, reason: 'رصيد افتتاحي', offsetAccountId, lines: [{ productId, direction: 'IN', quantity, unitCost }] });
const sell = (customerId: string, productId: string, quantity: string, extra = {}) =>
  c.postInvoice(customerId, [{ productId, quantity, unitPrice: '100' }], extra);
const entryLines = async (journalEntryId: string) => (await c.journal(journalEntryId)).map((l) => [l.accountCode, l.debit, l.credit]);

describe('receipts, sales and weighted average cost', () => {
  let product: string;
  let customer: string;

  beforeAll(async () => {
    product = (await c.goods()).id;
    customer = (await c.customer()).id;
  });

  it('records opening stock with an adjustment against capital', async () => {
    const res = await adjustIn(product, '10', '50', main, ch.byKey.get('CAPITAL'));
    expect(res.statusCode).toBe(201);
    expect(res.json()).toMatchObject({ totalIncrease: '500.00', offsetAccountCode: '3100' });
    expect(await entryLines(res.json().journalEntryId)).toEqual([['1400', '500.00', '0.00'], ['3100', '0.00', '500.00']]);
    expect(await balance(product)).toEqual({ quantity: '10.0000', value: '500.00' });
    await reconciled();
  });

  it('adds purchases at the net purchase price and moves the average', async () => {
    const sup = await c.supplier();
    await c.purchase(sup.id, [{ productId: product, quantity: '10', unitPrice: '60' }]);
    expect(await balance(product)).toEqual({ quantity: '20.0000', value: '1100.00' });
    await reconciled();
  });

  it('posts cost of goods sold with the sale: Dr COGS / Cr Inventory at average cost', async () => {
    const inv = await sell(customer, product, '4');
    expect(await entryLines(inv.journalEntryId)).toEqual([
      ['1300', '460.00', '0.00'], ['4100', '0.00', '400.00'], ['2210', '0.00', '60.00'],
      ['5100', '220.00', '0.00'], ['1400', '0.00', '220.00'],
    ]);
    expect(await balance(product)).toEqual({ quantity: '16.0000', value: '880.00' });
    const card = (await c.api.get(`/api/inventory/movements?productId=${product}`)).json().data;
    expect(card[0]).toMatchObject({ type: 'SALE', direction: 'OUT', quantity: '4.0000', totalCost: '220.00', unitCost: '55.000000', balanceQuantity: '16.0000', balanceValue: '880.00' });
    await reconciled();
  });

  it('refuses to sell more than is in stock and changes nothing', async () => {
    const draft = await c.invoice(customer, [{ productId: product, quantity: '100', unitPrice: '100' }]);
    const res = await c.api.post(`/api/invoices/${draft.id}/post`);
    expect(res.statusCode).toBe(409);
    expect(res.json().error.code).toBe('INSUFFICIENT_STOCK');
    expect((await c.api.get(`/api/invoices/${draft.id}`)).json().status).toBe('DRAFT');
    expect(await balance(product)).toEqual({ quantity: '16.0000', value: '880.00' });
  });

  it('brings returned goods back at the cost they left at, exactly', async () => {
    const inv = await sell(customer, product, '3'); // 3 × 55 = 165
    const ret = async (q: string) => {
      const r = (await c.api.post('/api/sales-returns', { originalInvoiceId: inv.id, docDate: c.date, reason: 'إرجاع', lines: [{ sourceItemId: inv.lines[0].id, quantity: q }] })).json();
      return (await c.api.post(`/api/sales-returns/${r.id}/post`)).json();
    };
    const r1 = await ret('1');
    expect((await entryLines(r1.journalEntryId)).slice(-2)).toEqual([['5100', '0.00', '55.00'], ['1400', '55.00', '0.00']]);
    const r2 = await ret('2');
    expect((await entryLines(r2.journalEntryId)).slice(-2)).toEqual([['5100', '0.00', '110.00'], ['1400', '110.00', '0.00']]);
    expect(await balance(product)).toEqual({ quantity: '16.0000', value: '880.00' });
    await reconciled();
  });

  it('restores stock when a sale is cancelled', async () => {
    const inv = await sell(customer, product, '2');
    expect(await balance(product)).toEqual({ quantity: '14.0000', value: '770.00' });
    await c.api.post(`/api/invoices/${inv.id}/cancel`, { reason: 'خطأ' });
    expect(await balance(product)).toEqual({ quantity: '16.0000', value: '880.00' });
    const card = (await c.api.get(`/api/inventory/movements?productId=${product}`)).json().data;
    expect(card[0]).toMatchObject({ type: 'CANCELLATION_IN', totalCost: '110.00' });
    await reconciled();
  });
});

describe('rounding', () => {
  it('issues 100.00 over 3 units at the running average and leaves nothing behind', async () => {
    const product = (await c.goods()).id;
    const cust = (await c.customer()).id;
    await adjustIn(product, '3', '33.333333');
    expect(await balance(product)).toEqual({ quantity: '3.0000', value: '100.00' });
    const costs = [];
    for (let i = 0; i < 3; i++) {
      const inv = await sell(cust, product, '1');
      costs.push((await c.journal(inv.journalEntryId)).find((l) => l.accountCode === '5100')!.debit);
    }
    // 100.00 / 3 = 33.33; then 66.67 / 2 = 33.335 → 33.34; the last unit takes the remaining 33.33.
    expect(costs).toEqual(['33.33', '33.34', '33.33']);
    expect(await balance(product)).toEqual({ quantity: '0.0000', value: '0.00' });
    await reconciled();
  });
});

describe('purchase returns and cancellations', () => {
  it('returns goods at average cost and books the price difference to COGS', async () => {
    const product = (await c.goods()).id;
    const sup = (await c.supplier()).id;
    await c.purchase(sup, [{ productId: product, quantity: '10', unitPrice: '10' }]);
    const expensive = await c.purchase(sup, [{ productId: product, quantity: '10', unitPrice: '20' }]); // average 15
    const r = (await c.api.post('/api/purchase-returns', { originalInvoiceId: expensive.id, docDate: c.date, reason: 'معيب', lines: [{ sourceItemId: expensive.lines[0].id, quantity: '5' }] })).json();
    const ret = (await c.api.post(`/api/purchase-returns/${r.id}/post`)).json();
    expect(await entryLines(ret.journalEntryId)).toEqual([
      ['2100', '115.00', '0.00'], ['1400', '0.00', '75.00'], ['5100', '0.00', '25.00'], ['2220', '0.00', '15.00'],
    ]);
    expect(await balance(product)).toEqual({ quantity: '15.0000', value: '225.00' });
    await reconciled();
  });

  it('refuses to cancel a purchase whose goods were already used', async () => {
    const product = (await c.goods()).id;
    const bill = await c.purchase((await c.supplier()).id, [{ productId: product, quantity: '5', unitPrice: '10' }]);
    await sell((await c.customer()).id, product, '3');
    const res = await c.api.post(`/api/purchase-invoices/${bill.id}/cancel`, { reason: 'إلغاء' });
    expect(res.json().error.code).toBe('INSUFFICIENT_STOCK');
  });

  it('expenses leftover value when cancelling a purchase empties the warehouse', async () => {
    const product = (await c.goods()).id;
    const sup = (await c.supplier()).id;
    await c.purchase(sup, [{ productId: product, quantity: '10', unitPrice: '20' }]);
    const cheap = await c.purchase(sup, [{ productId: product, quantity: '10', unitPrice: '10' }]); // 20 units, 300.00
    await sell((await c.customer()).id, product, '10'); // 150.00 out; 10 left worth 150.00
    const res = await c.api.post(`/api/purchase-invoices/${cheap.id}/cancel`, { reason: 'فاتورة مكررة' });
    expect(res.json().status).toBe('CANCELLED');
    expect(await balance(product)).toEqual({ quantity: '0.0000', value: '0.00' });
    const residual = (await c.api.get('/api/journal-entries?referenceType=INVENTORY_RESIDUAL')).json().data;
    expect(residual[0].totalDebit).toBe('50.00');
    await reconciled();
  });
});

describe('transfers', () => {
  it('moves stock between warehouses at average cost without a journal entry', async () => {
    const product = (await c.goods()).id;
    await adjustIn(product, '10', '12');
    const entriesBefore = (await c.api.get('/api/journal-entries')).json().total;
    const res = await c.api.post('/api/stock-transfers', { fromWarehouseId: main, toWarehouseId: second, date: c.date, lines: [{ productId: product, quantity: '4' }] });
    expect(res.statusCode).toBe(201);
    expect(res.json()).toMatchObject({ totalCost: '48.00', fromWarehouseId: main, toWarehouseId: second });
    expect(res.json().number).toMatch(/^TRF-\d{6}$/);
    expect(await balance(product, main)).toEqual({ quantity: '6.0000', value: '72.00' });
    expect(await balance(product, second)).toEqual({ quantity: '4.0000', value: '48.00' });
    expect((await c.api.get('/api/journal-entries')).json().total).toBe(entriesBefore);
    await reconciled();

    expect((await c.api.post('/api/stock-transfers', { fromWarehouseId: second, toWarehouseId: main, date: c.date, lines: [{ productId: product, quantity: '5' }] })).json().error.code).toBe('INSUFFICIENT_STOCK');
    expect((await c.api.post('/api/stock-transfers', { fromWarehouseId: main, toWarehouseId: main, date: c.date, lines: [{ productId: product, quantity: '1' }] })).statusCode).toBe(400);
  });

  it('sells from the warehouse chosen on the invoice', async () => {
    const product = (await c.goods()).id;
    await adjustIn(product, '5', '10', second);
    const res = await c.api.post('/api/invoices', { partyId: (await c.customer()).id, docDate: c.date, warehouseId: second, lines: [{ productId: product, quantity: '2', unitPrice: '50' }] });
    const inv = (await c.api.post(`/api/invoices/${res.json().id}/post`)).json();
    expect(inv.status).toBe('ISSUED');
    expect(await balance(product, second)).toEqual({ quantity: '3.0000', value: '30.00' });
  });
});

describe('stock counts and adjustments', () => {
  it('adjusts to the counted quantity at average cost', async () => {
    const product = (await c.goods()).id;
    await adjustIn(product, '10', '8');
    const res = await c.api.post('/api/stock-adjustments', { warehouseId: main, date: c.date, reason: 'جرد نهاية الشهر', lines: [{ productId: product, countedQuantity: '7' }] });
    expect(res.statusCode).toBe(201);
    expect(res.json().lines[0]).toMatchObject({ direction: 'OUT', quantity: '3.0000', countedQuantity: '7.0000', systemQuantity: '10.0000', totalCost: '24.00' });
    expect(await entryLines(res.json().journalEntryId)).toEqual([['5200', '24.00', '0.00'], ['1400', '0.00', '24.00']]);
    const same = await c.api.post('/api/stock-adjustments', { warehouseId: main, date: c.date, reason: 'جرد ثان', lines: [{ productId: product, countedQuantity: '7' }] });
    expect(same.json().error.code).toBe('NOTHING_TO_ADJUST');
    await reconciled();
  });

  it('validates adjustments', async () => {
    const product = (await c.goods()).id;
    expect((await c.api.post('/api/stock-adjustments', { warehouseId: main, date: c.date, reason: 'بلا تكلفة', lines: [{ productId: product, direction: 'IN', quantity: '1' }] })).json().error.code).toBe('UNIT_COST_REQUIRED');
    expect((await c.api.post('/api/stock-adjustments', { warehouseId: main, date: c.date, reason: 'نقص', lines: [{ productId: product, direction: 'OUT', quantity: '1' }] })).json().error.code).toBe('INSUFFICIENT_STOCK');
    expect((await adjustIn(product, '1', '5', main, ch.byKey.get('INVENTORY'))).json().error.code).toBe('INVALID_OFFSET_ACCOUNT');
    const service = (await c.product()).id;
    expect((await adjustIn(service, '1', '5')).json().error.code).toBe('NOT_STOCKED');
  });
});

describe('product rules', () => {
  it('locks type and tracking, and deletion, once a product has stock history', async () => {
    const product = (await c.goods()).id;
    await adjustIn(product, '1', '5');
    expect((await c.api.patch(`/api/products/${product}`, { productType: 'SERVICE' })).json().error.code).toBe('PRODUCT_HAS_STOCK_HISTORY');
    expect((await c.api.patch(`/api/products/${product}`, { trackInventory: false })).json().error.code).toBe('PRODUCT_HAS_STOCK_HISTORY');
    expect((await c.api.del(`/api/products/${product}`)).json().error.code).toBe('PRODUCT_IN_USE');
    const listed = (await c.api.get(`/api/products?search=${encodeURIComponent((await c.api.get(`/api/products/${product}`)).json().sku)}`)).json().data[0];
    expect(listed.onHand).toBe('1.0000');
  });
});

describe('integrity, permissions and isolation', () => {
  it('keeps movements append-only and balances non-negative in the database', async () => {
    const product = (await c.goods()).id;
    await adjustIn(product, '2', '5');
    const ctx = { tenantId: owner.tenantId, userId: owner.userId };
    await expect(withTx(t.pool, ctx, (db) => db.query(`UPDATE stock_movements SET quantity = 99 WHERE product_id = $1`, [product]))).rejects.toThrow(/permission denied|append-only/);
    await expect(withTx(t.pool, ctx, (db) => db.query(`DELETE FROM stock_movements WHERE product_id = $1`, [product]))).rejects.toThrow(/permission denied|append-only/);
    await expect(withTx(t.pool, ctx, (db) => db.query(`UPDATE inventory_balances SET quantity = -1 WHERE product_id = $1`, [product]))).rejects.toThrow(/check constraint/);
    await expect(withTx(t.pool, ctx, (db) => db.query(`UPDATE inventory_balances SET quantity = 0 WHERE product_id = $1`, [product]))).rejects.toThrow(/empty_ck/);
  });

  it('applies inventory permissions', async () => {
    const product = (await c.goods()).id;
    await adjustIn(product, '5', '5');
    const emp = await addMember(t.app, owner, 'WAREHOUSE_EMPLOYEE');
    const eApi = client(t.app, emp.token);
    expect((await eApi.post('/api/stock-transfers', { fromWarehouseId: main, toWarehouseId: second, date: c.date, lines: [{ productId: product, quantity: '1' }] })).statusCode).toBe(201);
    expect((await eApi.post('/api/stock-adjustments', { warehouseId: main, date: c.date, reason: 'تجربة', lines: [{ productId: product, countedQuantity: '1' }] })).statusCode).toBe(403);
    const sales = await addMember(t.app, owner, 'SALES_EMPLOYEE');
    expect((await client(t.app, sales.token).get('/api/inventory/balances')).statusCode).toBe(403);
  });

  it('isolates stock between tenants', async () => {
    const product = (await c.goods()).id;
    await adjustIn(product, '5', '5');
    const other = await register(t.app);
    const oc = await commerce(t.app, other);
    expect((await oc.api.get('/api/inventory/balances')).json().data).toEqual([]);
    expect((await oc.api.get(`/api/inventory/movements?productId=${product}`)).json().total).toBe(0);
    const res = await oc.api.post('/api/stock-transfers', { fromWarehouseId: main, toWarehouseId: second, date: oc.date, lines: [{ productId: product, quantity: '1' }] });
    expect(res.statusCode).toBe(400);
  });

  it('ends with stock value equal to the ledger and a balanced trial balance', async () => {
    const v = await reconciled();
    expect(v.reconciled).toBe(true);
    const [year] = await fiscalYears(t.app, owner.token);
    expect((await c.api.get(`/api/reports/trial-balance?dateFrom=${year!.startDate}&dateTo=${year!.endDate}`)).json().balanced).toBe(true);
  });
});
