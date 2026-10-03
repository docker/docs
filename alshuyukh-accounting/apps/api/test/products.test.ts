import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { addMember, chart, client, register, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
});
afterAll(async () => { await t.close(); });

const api = () => client(t.app, owner.token);
const unitId = async (code: string) => (await api().get('/api/units')).json().data.find((u: { code: string }) => u.code === code).id as string;

describe('units and categories', () => {
  it('seeds UN/ECE units for every company', async () => {
    const codes = (await api().get('/api/units')).json().data.map((u: { code: string }) => u.code);
    expect(codes).toEqual(expect.arrayContaining(['PCE', 'BX', 'KGM', 'LTR', 'HUR']));
  });

  it('adds units and nested categories', async () => {
    expect((await api().post('/api/units', { code: 'pk', nameAr: 'باكيت' })).json().code).toBe('PK');
    expect((await api().post('/api/units', { code: 'PK', nameAr: 'مكرر' })).statusCode).toBe(409);
    const parent = (await api().post('/api/product-categories', { nameAr: 'مواد غذائية' })).json();
    const child = await api().post('/api/product-categories', { nameAr: 'ألبان', parentId: parent.id });
    expect(child.json().parentId).toBe(parent.id);
    expect((await api().post('/api/product-categories', { nameAr: 'ألبان', parentId: parent.id })).statusCode).toBe(409);
  });
});

describe('products', () => {
  it('creates goods with an automatic SKU and exact prices', async () => {
    const res = await api().post('/api/products', {
      nameAr: 'حليب طازج 1 لتر', barcode: '6281007000017', unitId: await unitId('LTR'),
      salePrice: '6.5000', purchasePrice: '4.1250', salePriceIncludesVat: true,
    });
    expect(res.statusCode).toBe(201);
    expect(res.json()).toMatchObject({
      productType: 'GOODS', trackInventory: true, vatCategory: 'S', unitCode: 'LTR',
      salePrice: '6.5000', purchasePrice: '4.1250', salePriceIncludesVat: true,
    });
    expect(res.json().sku).toMatch(/^PRD-\d{5}$/);
  });

  it('creates services without stock and defaults to hours', async () => {
    const res = await api().post('/api/products', { nameAr: 'خدمة صيانة', productType: 'SERVICE', salePrice: '250' });
    expect(res.json()).toMatchObject({ productType: 'SERVICE', trackInventory: false, unitCode: 'HUR' });
    expect((await api().post('/api/products', { nameAr: 'خدمة خاطئة', productType: 'SERVICE', trackInventory: true })).statusCode).toBe(400);
    const goods = (await api().post('/api/products', { nameAr: 'سلعة تتحول' })).json();
    const changed = (await api().patch(`/api/products/${goods.id}`, { productType: 'SERVICE' })).json();
    expect(changed).toMatchObject({ productType: 'SERVICE', trackInventory: false });
  });

  it('validates prices, codes and uniqueness', async () => {
    const bad = [
      { nameAr: 'سعر رقمي', salePrice: 10 },
      { nameAr: 'خمس خانات', salePrice: '1.12345' },
      { nameAr: 'سالب', salePrice: '-1' },
      { nameAr: 'رمز خاطئ', sku: 'a b' },
      { nameAr: 'فئة ضريبية', vatCategory: 'X' },
    ];
    for (const body of bad) expect((await api().post('/api/products', body)).statusCode, JSON.stringify(body)).toBe(400);
    await api().post('/api/products', { nameAr: 'منتج أ', sku: 'SKU-1', barcode: '1111222233334' });
    expect((await api().post('/api/products', { nameAr: 'منتج ب', sku: 'SKU-1' })).statusCode).toBe(409);
    expect((await api().post('/api/products', { nameAr: 'منتج ج', barcode: '1111222233334' })).statusCode).toBe(409);
  });

  it('accepts account overrides of the right type only', async () => {
    const c = await chart(t.app, owner.token);
    expect((await api().post('/api/products', { nameAr: 'حساب خاطئ', salesAccountId: c.byKey.get('CASH') })).statusCode).toBe(400);
    expect((await api().post('/api/products', { nameAr: 'حساب تجميعي', salesAccountId: c.byCode.get('4000') })).statusCode).toBe(400);
    const ok = await api().post('/api/products', { nameAr: 'حساب صحيح', salesAccountId: c.byKey.get('SALES'), purchaseAccountId: c.byKey.get('INVENTORY') });
    expect(ok.statusCode).toBe(201);
  });

  it('searches by name, SKU and exact barcode, and filters by category', async () => {
    const cat = (await api().post('/api/product-categories', { nameAr: 'منظفات' })).json();
    await api().post('/api/products', { nameAr: 'منظف أرضيات', sku: 'CLN-1', barcode: '9990001112223', categoryId: cat.id });
    for (const term of ['أرضيات', 'CLN-1', '9990001112223']) {
      expect((await api().get(`/api/products?search=${encodeURIComponent(term)}`)).json().data.map((p: { sku: string }) => p.sku), term).toEqual(['CLN-1']);
    }
    expect((await api().get(`/api/products?categoryId=${cat.id}`)).json().total).toBe(1);
  });

  it('soft-deletes and frees the SKU', async () => {
    const p = (await api().post('/api/products', { nameAr: 'مؤقت', sku: 'TMP-1' })).json();
    expect((await api().del(`/api/products/${p.id}`)).statusCode).toBe(204);
    expect((await api().get(`/api/products/${p.id}`)).statusCode).toBe(404);
    expect((await api().post('/api/products', { nameAr: 'بديل', sku: 'TMP-1' })).statusCode).toBe(201);
  });

  it('refuses units and categories from another company or tenant', async () => {
    const other = await register(t.app);
    const foreignUnit = (await client(t.app, other.token).get('/api/units')).json().data[0].id;
    expect((await api().post('/api/products', { nameAr: 'وحدة أجنبية', unitId: foreignUnit })).json().error.code).toBe('INVALID_UNIT');
    const p = (await api().post('/api/products', { nameAr: 'سري' })).json();
    expect((await client(t.app, other.token).get(`/api/products/${p.id}`)).statusCode).toBe(404);
    expect((await client(t.app, other.token).get('/api/products')).json().total).toBe(0);
  });

  it('applies product permissions per role', async () => {
    const wh = await addMember(t.app, owner, 'WAREHOUSE_MANAGER');
    expect((await client(t.app, wh.token).post('/api/products', { nameAr: 'من المستودع' })).statusCode).toBe(201);
    const sales = await addMember(t.app, owner, 'SALES_EMPLOYEE');
    expect((await client(t.app, sales.token).get('/api/products')).statusCode).toBe(200);
    expect((await client(t.app, sales.token).post('/api/products', { nameAr: 'من المبيعات' })).statusCode).toBe(403);
    expect((await client(t.app, sales.token).post('/api/units', { code: 'X', nameAr: 'س' })).statusCode).toBe(403);
  });
});
