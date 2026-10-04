// End-to-end smoke test in a real browser against a running deployment.
//   docker compose -f docker-compose.prod.yml --env-file deploy/.env up -d
//   E2E_BASE_URL=http://localhost:8080 npm run e2e
// It registers a new organization, so point it at a staging or local stack.
// CHROMIUM_PATH selects the browser binary (defaults to Playwright's own lookup).
import { chromium } from 'playwright-core';

const BASE = (process.env.E2E_BASE_URL ?? 'http://localhost:8080').replace(/\/$/, '');
const browser = await chromium.launch({ executablePath: process.env.CHROMIUM_PATH || undefined });
const problems = [];
const page = await (await browser.newContext({ viewport: { width: 1280, height: 900 }, locale: 'ar-SA' })).newPage();
page.on('pageerror', (e) => problems.push(`page error: ${e.message}`));
page.on('console', (m) => {
  // Expected 4xx answers show up as console errors; anything else (CSP violations, crashes) fails the run.
  if (m.type() === 'error' && !/status of 4\d\d/.test(m.text())) problems.push(`console: ${m.text()}`);
});

let current = '';
const step = (name) => { current = name; console.log(`• ${name}`); };
const run = Date.now().toString(36);

try {
  step('register a new organization');
  await page.goto(`${BASE}/register`);
  await page.fill('input[name=fullName]', 'مستخدم اختبار');
  await page.fill('input[name=email]', `smoke-${run}@example.test`);
  await page.fill('input[name=password]', 'Str0ng-Passw0rd!');
  await page.fill('input[name=companyName]', `شركة الاختبار ${run}`);
  await page.click('button:has-text("إنشاء الحساب")');
  await page.waitForSelector('text=مرحبًا');

  step('customer and product');
  await page.goto(`${BASE}/customers`);
  await page.click('button:has-text("عميل جديد")');
  await page.fill('input[name=nameAr]', 'شركة النخبة للتجارة');
  await page.fill('input[name=vatNumber]', '300000000000003');
  await page.click('button:has-text("حفظ")');
  await page.waitForSelector('td:has-text("CUS-00001")');
  await page.goto(`${BASE}/inventory`);
  await page.click('button:has-text("منتج جديد")');
  await page.selectOption('select[name=productType]', 'SERVICE');
  await page.fill('input[name=nameAr]', 'خدمة استشارية');
  await page.fill('input[name=salePrice]', '1000');
  await page.fill('input[name=purchasePrice]', '600');
  await page.click('button:has-text("حفظ")');
  await page.waitForSelector('td:has-text("PRD-00001")');

  step('sales invoice: 1000 SAR + 15% VAT');
  await page.goto(`${BASE}/sales`);
  await page.click('a:has-text("فاتورة جديدة")');
  await page.locator('form select').first().selectOption({ label: 'CUS-00001 — شركة النخبة للتجارة' });
  await page.locator('[aria-label="صنف السطر 1"]').selectOption({ label: 'PRD-00001 — خدمة استشارية' });
  await page.waitForSelector('.totals .strong dd:has-text("1,150.00")');
  await page.click('button:has-text("حفظ كمسودة")');
  await page.waitForSelector('h2:has-text("مسودة")');
  await page.click('button:has-text("إصدار")');
  await page.waitForSelector('h2:has-text("INV-000001")');

  step('journal entry: Dr receivables 1150 / Cr sales 1000 / Cr VAT 150');
  await page.click('a:has-text("عرض القيد المحاسبي")');
  await page.waitForSelector('text=فاتورة مبيعات INV-000001');
  const body = await page.textContent('main');
  for (const amount of ['1,150.00', '1,000.00', '150.00']) {
    if (!body.includes(amount)) throw new Error(`journal entry is missing ${amount}`);
  }

  step('customer payment');
  await page.goBack();
  await page.waitForSelector('h2:has-text("INV-000001")');
  await page.click('button:has-text("تسجيل دفعة")');
  await page.fill('input[name=amount]', '1150.00');
  await page.click('.reverse-box button:has-text("حفظ")');
  await page.waitForSelector('.tag:has-text("مدفوعة")');

  step('trial balance is balanced');
  await page.goto(`${BASE}/accounting/trial-balance`);
  await page.waitForSelector('text=الميزان متوازن');

  step('reports and dashboard load');
  await page.goto(`${BASE}/reports`);
  await page.waitForSelector('main h2, main h1');
  await page.goto(`${BASE}/`);
  await page.waitForSelector('text=مرحبًا');

  if (problems.length) throw new Error(`browser reported problems:\n  ${problems.join('\n  ')}`);
  console.log('smoke test passed');
} catch (err) {
  const shot = `e2e-failure-${run}.png`;
  await page.screenshot({ path: shot, fullPage: true }).catch(() => undefined);
  console.error(`failed at "${current}": ${err.message}\nscreenshot: ${shot}`);
  process.exitCode = 1;
} finally {
  await browser.close();
}
