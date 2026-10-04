/**
 * Performance check. Builds one large organization through the real API code
 * paths (invoices, payments, purchases, expenses), then times the heavy
 * endpoints. Runs only against a database whose name ends in _perf, which it
 * rebuilds from scratch:
 *
 *   PERF_DATABASE_URL=postgres://alshuyukh_app:…/alshuyukh_perf \
 *   PERF_DATABASE_URL_MIGRATE=postgres://alshuyukh_owner:…/alshuyukh_perf \
 *   PERF_INVOICES=3000 npm run perf -w apps/api
 *
 * PERF_REUSE=1 skips the rebuild and only measures the data already there.
 */
import pg from 'pg';
import { buildApp } from '../src/app.js';
import { loadEnv } from '../src/config/env.js';
import { createPool } from '../src/db/pool.js';
import { migrate } from '../src/db/migrate.js';

const appUrl = process.env.PERF_DATABASE_URL;
const ownerUrl = process.env.PERF_DATABASE_URL_MIGRATE;
if (!appUrl || !ownerUrl || !/_perf$/.test(new URL(ownerUrl).pathname)) {
  console.error('Set PERF_DATABASE_URL and PERF_DATABASE_URL_MIGRATE to a database whose name ends in _perf');
  process.exit(1);
}
const N = Number(process.env.PERF_INVOICES ?? 3000);
const REUSE = process.env.PERF_REUSE === '1';
const EMAIL = 'perf@example.test';

if (!REUSE) {
  const reset = new pg.Client({ connectionString: ownerUrl });
  await reset.connect();
  await reset.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  await reset.end();
}
await migrate(ownerUrl);
// Unlimited plan for the load (plan limits are not what is measured here).
const owner = new pg.Client({ connectionString: ownerUrl });
await owner.connect();
await owner.query(`UPDATE plans SET max_users = NULL, max_companies = NULL, max_branches = NULL, max_warehouses = NULL, max_products = NULL,
  max_invoices_per_month = NULL, max_api_calls_per_month = NULL WHERE is_default`);

const env = loadEnv({ ...process.env, DATABASE_URL: appUrl, NODE_ENV: 'test' });
const pool = createPool(appUrl);
const app = await buildApp({ env, pool, logger: false });
await app.ready();

let token = '';
const call = async (method: string, url: string, payload?: unknown) => {
  const res = await app.inject({ method: method as 'GET', url, payload: payload as object, headers: { authorization: `Bearer ${token}` } });
  if (res.statusCode >= 400) throw new Error(`${method} ${url}: ${res.statusCode} ${res.body.slice(0, 300)}`);
  return res.json();
};

const auth = REUSE
  ? (await app.inject({ method: 'POST', url: '/api/auth/login', payload: { email: EMAIL, password: 'Str0ng-Passw0rd!' } })).json()
  : (await app.inject({ method: 'POST', url: '/api/auth/register', payload: { fullName: 'Perf', email: EMAIL, password: 'Str0ng-Passw0rd!', tenantName: 'Perf Org', companyName: 'Perf Co' } })).json();
token = auth.accessToken;
const [year] = (await call('GET', '/api/fiscal-years')).data;
const customers: string[] = (await call('GET', '/api/customers?limit=1')).data.map((x: { id: string }) => x.id);
if (!REUSE) {
  customers.length = 0;
  const day = (i: number) => {
    const d = new Date(`${year.startDate}T00:00:00Z`);
    d.setUTCDate(d.getUTCDate() + (i % 330));
    return d.toISOString().slice(0, 10);
  };
  const t0 = Date.now();
  const methods = (await call('GET', '/api/payment-methods')).data;
  const cash = methods.find((m: { code: string }) => m.code === 'CASH').id;
  for (let i = 0; i < 200; i++) customers.push((await call('POST', '/api/customers', { nameAr: `عميل ${i}`, paymentTermsDays: 30 })).id);
  const suppliers = [];
  for (let i = 0; i < 50; i++) suppliers.push((await call('POST', '/api/suppliers', { nameAr: `مورد ${i}` })).id);
  const services = [];
  for (let i = 0; i < 150; i++) services.push((await call('POST', '/api/products', { nameAr: `خدمة ${i}`, productType: 'SERVICE', salePrice: String(50 + i) })).id);
  const goods = [];
  for (let i = 0; i < 50; i++) goods.push((await call('POST', '/api/products', { nameAr: `صنف ${i}`, productType: 'GOODS', salePrice: '100', purchasePrice: '60' })).id);
  // Stock first so goods can be sold.
  for (let i = 0; i < 100; i++) {
    const bill = await call('POST', '/api/purchase-invoices', { partyId: suppliers[i % 50], docDate: day(0), lines: goods.slice(0, 10).map((g) => ({ productId: g, quantity: '200' })) });
    await call('POST', `/api/purchase-invoices/${bill.id}/post`);
  }
  console.log(`master data + purchases: ${((Date.now() - t0) / 1000).toFixed(1)}s`);
  const t1 = Date.now();
  const cats = (await call('GET', '/api/expense-categories')).data;
  for (let i = 0; i < N; i++) {
    const inv = await call('POST', '/api/invoices', { partyId: customers[i % 200], docDate: day(i), lines: [
      { productId: services[i % 150], quantity: '2' }, { productId: services[(i * 7) % 150], quantity: '1' }, { productId: goods[i % 10], quantity: '1' },
    ] });
    const issued = await call('POST', `/api/invoices/${inv.id}/post`);
    if (i % 2 === 0) {
      await call('POST', '/api/payments', { direction: 'RECEIPT', customerId: customers[i % 200], paymentDate: day(i), methodId: cash, amount: issued.total,
        allocations: [{ documentType: 'SALES_INVOICE', documentId: inv.id, amount: issued.total }] });
    }
    if (i % 5 === 0) {
      await call('POST', '/api/expenses', { expenseDate: day(i), paymentType: 'CASH', methodId: cash, post: true, lines: [{ categoryId: cats[i % cats.length].id, amount: String(100 + (i % 50)) }] });
    }
  }
  const seconds = (Date.now() - t1) / 1000;
  const { rows: [c] } = await owner.query(`SELECT (SELECT count(*) FROM journal_entries) AS entries, (SELECT count(*) FROM journal_entry_lines) AS lines,
    (SELECT count(*) FROM sales_invoices) AS invoices, (SELECT count(*) FROM stock_movements) AS movements`);
  console.log(`${N} invoices (+ payments, expenses) in ${seconds.toFixed(1)}s → ${(N / seconds).toFixed(1)} invoices/s`);
  console.log(`data: ${c.entries} journal entries, ${c.lines} lines, ${c.invoices} invoices, ${c.movements} stock movements`);

}
const range = `dateFrom=${year.startDate}&dateTo=${year.endDate}`;
const cashAccount = (await call('GET', '/api/accounts?postableOnly=true')).data.find((a: { code: string }) => a.code === '1100').id;
const endpoints = [
  '/api/dashboard', `/api/reports/profit-loss?${range}`, `/api/reports/balance-sheet?asOf=${year.endDate}`, `/api/reports/cash-flow?${range}`,
  `/api/reports/trial-balance?${range}`, `/api/reports/general-ledger?${range}&accountId=${cashAccount}`, `/api/reports/journal?${range}&limit=100`,
  '/api/reports/receivables-aging', `/api/reports/customer-statement?${range}&partyId=${customers[0]}`, `/api/reports/sales?${range}&groupBy=product`,
  `/api/reports/expenses?${range}`, `/api/reports/vat-return?${range}`, '/api/inventory/valuation', '/api/invoices?limit=50', '/api/journal-entries?limit=50',
  '/api/customers?limit=50', `/api/customers/${customers[0]}`, '/api/auth/me',
];
console.log('\nendpoint                                              p50 ms   p95 ms');
for (const url of endpoints) {
  const times: number[] = [];
  for (let i = 0; i < 12; i++) {
    const s = performance.now();
    await call('GET', url);
    times.push(performance.now() - s);
  }
  times.sort((a, b) => a - b);
  console.log(`${url.split('?')[0]!.padEnd(52)} ${times[5]!.toFixed(1).padStart(7)} ${times[11]!.toFixed(1).padStart(8)}`);
}
await app.close();
await pool.end();
await owner.end();
