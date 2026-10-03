import type { Db } from '../../db/tx.js';
import { addDays, addMonths, todayIn } from '../../lib/dates.js';
import { badRequest, conflict } from '../../lib/errors.js';
import { DEFAULT_ACCOUNTS, DEFAULT_GROUPS } from './chart-template.js';

export interface CompanyRef { tenantId: string; companyId: string; userId: string }

/** Creates the default account groups and chart of accounts. Skips if the company already has accounts. */
export async function seedDefaultChart(db: Db, ref: CompanyRef): Promise<boolean> {
  const existing = await db.query(`SELECT 1 FROM accounts WHERE company_id = $1 AND deleted_at IS NULL LIMIT 1`, [ref.companyId]);
  if (existing.rowCount) return false;

  const groupIds = new Map<string, string>();
  for (const g of DEFAULT_GROUPS) {
    const { rows: [row] } = await db.query<{ id: string }>(
      `INSERT INTO account_groups (tenant_id, company_id, code, name_ar, name_en, account_type, sort_order, is_system)
       VALUES ($1, $2, $3, $4, $5, $6, $7, true)
       ON CONFLICT (company_id, code) DO UPDATE SET code = EXCLUDED.code
       RETURNING id`,
      [ref.tenantId, ref.companyId, g.code, g.ar, g.en, g.type, g.sort]);
    groupIds.set(g.code, row!.id);
  }

  const accountIds = new Map<string, string>();
  for (const a of DEFAULT_ACCOUNTS) {
    const { rows: [row] } = await db.query<{ id: string }>(
      `INSERT INTO accounts (tenant_id, company_id, code, name_ar, name_en, account_type, parent_id, group_id,
                             is_postable, is_system, system_key, created_by, updated_by)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true, $10, $11, $11) RETURNING id`,
      [ref.tenantId, ref.companyId, a.code, a.ar, a.en, a.type,
       a.parent ? accountIds.get(a.parent)! : null, a.group ? groupIds.get(a.group)! : null,
       a.postable, a.systemKey ?? null, ref.userId]);
    accountIds.set(a.code, row!.id);
  }
  return true;
}

const MONTHS_AR = ['يناير', 'فبراير', 'مارس', 'أبريل', 'مايو', 'يونيو', 'يوليو', 'أغسطس', 'سبتمبر', 'أكتوبر', 'نوفمبر', 'ديسمبر'];

/**
 * Creates a fiscal year with monthly periods. `startDate` must be the 1st of
 * a month. Without `endDate` the year is 12 months. The last period may be
 * shorter or longer only through an explicit endDate (first/last years).
 */
export async function createFiscalYear(db: Db, ref: CompanyRef, startDate: string, endDate?: string, name?: string) {
  if (!startDate.endsWith('-01')) throw badRequest('INVALID_FISCAL_START', 'A fiscal year must start on the first day of a month');
  const end = endDate ?? addDays(addMonths(startDate, 12), -1);
  if (end <= startDate) throw badRequest('INVALID_FISCAL_END', 'Fiscal year end must be after its start');

  const overlap = await db.query(
    `SELECT 1 FROM fiscal_years WHERE company_id = $1 AND daterange(start_date, end_date, '[]') && daterange($2::date, $3::date, '[]')`,
    [ref.companyId, startDate, end]);
  if (overlap.rowCount) throw conflict('FISCAL_YEAR_OVERLAP', 'This fiscal year overlaps an existing one');

  const yearName = name ?? (startDate.slice(0, 4) === end.slice(0, 4) ? startDate.slice(0, 4) : `${startDate.slice(0, 4)}/${end.slice(0, 4)}`);
  const { rows: [year] } = await db.query<{ id: string }>(
    `INSERT INTO fiscal_years (tenant_id, company_id, name, start_date, end_date, created_by)
     VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
    [ref.tenantId, ref.companyId, yearName, startDate, end, ref.userId]);

  let periodStart = startDate;
  let n = 1;
  while (periodStart <= end) {
    const monthEnd = addDays(addMonths(periodStart, 1), -1);
    const periodEnd = monthEnd > end ? end : monthEnd;
    const month = Number(periodStart.slice(5, 7));
    await db.query(
      `INSERT INTO fiscal_periods (tenant_id, company_id, fiscal_year_id, period_number, name, start_date, end_date)
       VALUES ($1, $2, $3, $4, $5, $6, $7)`,
      [ref.tenantId, ref.companyId, year!.id, n, `${MONTHS_AR[month - 1]} ${periodStart.slice(0, 4)}`, periodStart, periodEnd]);
    periodStart = addDays(periodEnd, 1);
    n += 1;
    if (n > 18) throw badRequest('FISCAL_YEAR_TOO_LONG', 'A fiscal year can have at most 18 periods');
  }
  await db.query(`INSERT INTO journal_sequences (fiscal_year_id, tenant_id, company_id) VALUES ($1, $2, $3)`,
    [year!.id, ref.tenantId, ref.companyId]);
  return year!.id;
}

/** Start date (1st of month) of the fiscal year that contains `today`. */
export function currentFiscalYearStart(today: string, startMonth: number): string {
  const year = Number(today.slice(0, 4));
  const month = Number(today.slice(5, 7));
  const startYear = month >= startMonth ? year : year - 1;
  return `${startYear}-${String(startMonth).padStart(2, '0')}-01`;
}

// UN/ECE Recommendation 20 codes, used later in ZATCA XML.
const DEFAULT_UNITS: readonly [string, string, string][] = [
  ['PCE', 'حبة', 'Piece'], ['BX', 'كرتون', 'Box'], ['KGM', 'كيلوجرام', 'Kilogram'], ['LTR', 'لتر', 'Litre'],
  ['MTR', 'متر', 'Metre'], ['SET', 'طقم', 'Set'], ['HUR', 'ساعة', 'Hour'], ['DAY', 'يوم', 'Day'],
];

/** Creates the default units of measure. Existing units are kept. */
export async function seedDefaultUnits(db: Db, ref: CompanyRef): Promise<void> {
  for (const [code, ar, en] of DEFAULT_UNITS) {
    await db.query(
      `INSERT INTO units (tenant_id, company_id, code, name_ar, name_en) VALUES ($1, $2, $3, $4, $5)
       ON CONFLICT (company_id, code) DO NOTHING`,
      [ref.tenantId, ref.companyId, code, ar, en]);
  }
}

/** Standard VAT rate (configurable later in tax settings) and default payment methods. */
export async function seedTaxAndPayments(db: Db, ref: CompanyRef): Promise<void> {
  const hasRate = await db.query(`SELECT 1 FROM tax_rates WHERE company_id = $1 LIMIT 1`, [ref.companyId]);
  if (!hasRate.rowCount) {
    await db.query(
      `INSERT INTO tax_rates (tenant_id, company_id, name_ar, rate, effective_from, created_by)
       VALUES ($1, $2, 'ضريبة القيمة المضافة - النسبة الأساسية', 0.15, DATE '2020-07-01', $3)`,
      [ref.tenantId, ref.companyId, ref.userId]);
  }
  const methods: [string, string, string, string, number][] = [
    ['CASH', 'نقدًا', 'CASH', 'CASH', 1], ['BANK', 'تحويل بنكي', 'BANK', 'BANK', 2],
    ['CARD', 'بطاقة مدى / ائتمان', 'CARD', 'BANK', 3], ['STC_PAY', 'STC Pay', 'STC_PAY', 'BANK', 4],
    ['TAMARA', 'تمارا', 'TAMARA', 'BANK', 5],
  ];
  for (const [code, name, type, accountKey, sort] of methods) {
    await db.query(
      `INSERT INTO payment_methods (tenant_id, company_id, code, name_ar, method_type, account_id, sort_order)
       SELECT $1, $2, $3, $4, $5, a.id, $7 FROM accounts a
        WHERE a.company_id = $2 AND a.system_key = $6 AND a.deleted_at IS NULL
       ON CONFLICT (company_id, code) DO NOTHING`,
      [ref.tenantId, ref.companyId, code, name, type, accountKey, sort]);
  }
}

/**
 * Prepares a new company for accounting: default chart and the current
 * fiscal year. Safe to call again; existing setup is left as is.
 */
export async function setupCompanyAccounting(db: Db, ref: CompanyRef, timezone = 'Asia/Riyadh') {
  const chartCreated = await seedDefaultChart(db, ref);
  await seedDefaultUnits(db, ref);
  await seedTaxAndPayments(db, ref);
  let fiscalYearId: string | null = null;
  const hasYear = await db.query(`SELECT 1 FROM fiscal_years WHERE company_id = $1 LIMIT 1`, [ref.companyId]);
  if (!hasYear.rowCount) {
    const { rows: [s] } = await db.query<{ fiscal_year_start_month: number }>(
      `SELECT fiscal_year_start_month FROM tenant_settings WHERE tenant_id = $1`, [ref.tenantId]);
    fiscalYearId = await createFiscalYear(db, ref, currentFiscalYearStart(todayIn(timezone), s?.fiscal_year_start_month ?? 1));
  }
  return { chartCreated, fiscalYearCreated: fiscalYearId !== null };
}
