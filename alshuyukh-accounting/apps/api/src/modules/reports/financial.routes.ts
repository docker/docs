import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { isoDate } from '../../lib/dates.js';
import { badRequest, notFound } from '../../lib/errors.js';
import { Decimal, sum, toMoney } from '../../lib/money.js';
import { parse } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { DEBIT_NORMAL, type AccountType } from '../accounting/chart-template.js';
import { resolveCompanyId } from '../accounting/company-context.js';
import { CASH_ACCOUNTS, checkRange, DOC_NUMBERS, GL, NOT_CLOSING, rangeFields } from './common.js';

interface AccountRow { accountId: string; code: string; nameAr: string; type: AccountType; groupCode: string | null; groupName: string | null; amount: string }

const SECTION_AR: Record<string, string> = {
  REVENUE: 'الإيرادات', COST_OF_GOODS_SOLD: 'تكلفة المبيعات', EXPENSE: 'المصروفات',
  ASSET: 'الأصول', LIABILITY: 'الالتزامات', EQUITY: 'حقوق الملكية',
};

const section = (type: AccountType, rows: AccountRow[]) => {
  const accounts = rows.filter((r) => r.type === type);
  return { type, title: SECTION_AR[type], accounts, total: toMoney(sum(accounts.map((a) => a.amount))) };
};

export default async function financialReportRoutes(app: FastifyInstance) {
  const canView = requirePermission(app, 'financial_report.view');

  /** Income statement: revenue − cost of sales = gross profit; − expenses = net profit. */
  app.get('/reports/profit-loss', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ ...rangeFields, branchId: z.uuid().optional(), costCenterId: z.uuid().optional() }), req.query);
    checkRange(q);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query<AccountRow>(
        `WITH g AS (${GL})
         SELECT a.id AS "accountId", a.code, a.name_ar AS "nameAr", a.account_type AS type, ag.code AS "groupCode", ag.name_ar AS "groupName",
                (CASE WHEN a.account_type = 'REVENUE' THEN sum(g.credit - g.debit) ELSE sum(g.debit - g.credit) END)::numeric(18,2)::text AS amount
           FROM g JOIN accounts a ON a.id = g.account_id LEFT JOIN account_groups ag ON ag.id = a.group_id
          WHERE a.account_type IN ('REVENUE', 'COST_OF_GOODS_SOLD', 'EXPENSE') AND ${NOT_CLOSING}
            AND g.entry_date BETWEEN $3 AND $4 AND ($5::uuid IS NULL OR g.branch_id = $5) AND ($6::uuid IS NULL OR g.cost_center_id = $6)
          GROUP BY a.id, ag.code, ag.name_ar HAVING sum(g.debit) <> 0 OR sum(g.credit) <> 0
          ORDER BY a.code`,
        [req.auth!.tenantId, companyId, q.dateFrom, q.dateTo, q.branchId ?? null, q.costCenterId ?? null]);
      const revenue = section('REVENUE', rows);
      const cogs = section('COST_OF_GOODS_SOLD', rows);
      const expenses = section('EXPENSE', rows);
      const grossProfit = new Decimal(revenue.total).minus(cogs.total);
      return {
        companyId, dateFrom: q.dateFrom, dateTo: q.dateTo,
        sections: [revenue, cogs, expenses],
        revenue: revenue.total, costOfSales: cogs.total, grossProfit: toMoney(grossProfit),
        expenses: expenses.total, netProfit: toMoney(grossProfit.minus(expenses.total)),
      };
    });
  });

  /**
   * Balance sheet at a date. Profit not yet closed into retained earnings is
   * shown as current-period earnings, so assets = liabilities + equity.
   */
  app.get('/reports/balance-sheet', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional(), asOf: isoDate }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query<AccountRow>(
        `WITH g AS (${GL})
         SELECT a.id AS "accountId", a.code, a.name_ar AS "nameAr", a.account_type AS type, ag.code AS "groupCode", ag.name_ar AS "groupName",
                (CASE WHEN a.account_type IN ('ASSET', 'EXPENSE', 'COST_OF_GOODS_SOLD') THEN sum(g.debit - g.credit) ELSE sum(g.credit - g.debit) END)::numeric(18,2)::text AS amount
           FROM g JOIN accounts a ON a.id = g.account_id LEFT JOIN account_groups ag ON ag.id = a.group_id
          WHERE g.entry_date <= $3
          GROUP BY a.id, ag.code, ag.name_ar HAVING sum(g.debit - g.credit) <> 0
          ORDER BY a.code`,
        [req.auth!.tenantId, companyId, q.asOf]);
      const pl = rows.filter((r) => !['ASSET', 'LIABILITY', 'EQUITY'].includes(r.type));
      const earnings = pl.reduce((s, r) => (r.type === 'REVENUE' ? s.plus(r.amount) : s.minus(r.amount)), new Decimal(0));
      const assets = section('ASSET', rows);
      const liabilities = section('LIABILITY', rows);
      const equity = section('EQUITY', rows);
      const equityTotal = new Decimal(equity.total).plus(earnings);
      const liabilitiesAndEquity = equityTotal.plus(liabilities.total);
      return {
        companyId, asOf: q.asOf,
        assets, liabilities, equity: { ...equity, currentEarnings: toMoney(earnings), total: toMoney(equityTotal) },
        totalAssets: assets.total, totalLiabilitiesAndEquity: toMoney(liabilitiesAndEquity),
        balanced: liabilitiesAndEquity.equals(assets.total),
      };
    });
  });

  /**
   * Cash flow statement, direct method. For every entry that moves cash, the
   * other lines say where the cash came from or went: their credit − debit
   * equals the cash movement exactly, so opening + net flow = closing.
   * Classification: equity and non-current liabilities are financing,
   * non-current assets are investing, everything else is operating.
   */
  app.get('/reports/cash-flow', { preHandler: canView }, async (req) => {
    const q = parse(z.object(rangeFields), req.query);
    checkRange(q);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const params = [req.auth!.tenantId, companyId, q.dateFrom, q.dateTo];
      const { rows: [bal] } = await db.query<{ opening: string; closing: string }>(
        `WITH g AS (${GL})
         SELECT COALESCE(sum(CASE WHEN g.entry_date < $3 THEN g.debit - g.credit END), 0)::numeric(18,2)::text AS opening,
                COALESCE(sum(g.debit - g.credit), 0)::numeric(18,2)::text AS closing
           FROM g WHERE g.entry_date <= $4 AND g.account_id IN (${CASH_ACCOUNTS})`, params);
      const { rows } = await db.query<{ category: string; accountId: string; code: string; nameAr: string; inflow: string; outflow: string }>(
        `WITH g AS (${GL}),
              cash_entries AS (SELECT DISTINCT entry_id FROM g WHERE g.entry_date BETWEEN $3 AND $4 AND g.account_id IN (${CASH_ACCOUNTS}))
         SELECT CASE WHEN a.account_type = 'EQUITY' OR ag.code = 'NON_CURRENT_LIABILITIES' THEN 'FINANCING'
                     WHEN ag.code = 'NON_CURRENT_ASSETS' THEN 'INVESTING' ELSE 'OPERATING' END AS category,
                a.id AS "accountId", a.code, a.name_ar AS "nameAr",
                sum(g.credit)::numeric(18,2)::text AS inflow, sum(g.debit)::numeric(18,2)::text AS outflow
           FROM g JOIN cash_entries c ON c.entry_id = g.entry_id
           JOIN accounts a ON a.id = g.account_id LEFT JOIN account_groups ag ON ag.id = a.group_id
          WHERE g.account_id NOT IN (${CASH_ACCOUNTS})
          GROUP BY 1, a.id, a.code, a.name_ar ORDER BY a.code`, params);
      const sections = (['OPERATING', 'INVESTING', 'FINANCING'] as const).map((category) => {
        const lines = rows.filter((r) => r.category === category).map((r) => ({ ...r, net: toMoney(new Decimal(r.inflow).minus(r.outflow)) }));
        return { category, lines, total: toMoney(sum(lines.map((l) => l.net))) };
      });
      const net = sum(sections.map((s) => s.total));
      return {
        companyId, dateFrom: q.dateFrom, dateTo: q.dateTo, sections,
        openingCash: bal!.opening, netChange: toMoney(net), closingCash: bal!.closing,
        reconciled: new Decimal(bal!.opening).plus(net).equals(bal!.closing),
      };
    });
  });

  /** General ledger for one account with opening balance and running balance. */
  app.get('/reports/general-ledger', { preHandler: canView }, async (req) => {
    const q = parse(z.object({
      ...rangeFields, accountId: z.uuid(), branchId: z.uuid().optional(), costCenterId: z.uuid().optional(),
      customerId: z.uuid().optional(), supplierId: z.uuid().optional(),
    }), req.query);
    checkRange(q);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows: [account] } = await db.query<{ id: string; code: string; nameAr: string; type: AccountType; isPostable: boolean }>(
        `SELECT id, code, name_ar AS "nameAr", account_type AS type, is_postable AS "isPostable" FROM accounts
          WHERE tenant_id = $1 AND company_id = $2 AND id = $3 AND deleted_at IS NULL`, [req.auth!.tenantId, companyId, q.accountId]);
      if (!account) throw notFound('Account');
      if (!account.isPostable) throw badRequest('HEADER_ACCOUNT', 'Choose a postable account; header accounts have no lines');
      const filters = `g.account_id = $3 AND ($6::uuid IS NULL OR g.branch_id = $6) AND ($7::uuid IS NULL OR g.cost_center_id = $7)
                       AND ($8::uuid IS NULL OR g.customer_id = $8) AND ($9::uuid IS NULL OR g.supplier_id = $9)`;
      const params = [req.auth!.tenantId, companyId, q.accountId, q.dateFrom, q.dateTo, q.branchId ?? null, q.costCenterId ?? null, q.customerId ?? null, q.supplierId ?? null];
      // Same parameter list as the lines query; $5 (dateTo) is referenced only so PostgreSQL can type it.
      const { rows: [o] } = await db.query<{ opening: string }>(
        `WITH g AS (${GL}) SELECT COALESCE(sum(g.debit - g.credit), 0)::text AS opening FROM g WHERE ${filters} AND g.entry_date < $4 AND $5::date IS NOT NULL`, params);
      const LIMIT = 5000;
      const { rows } = await db.query<{ entryId: string; entryNumber: string; date: string; description: string; referenceType: string; documentNumber: string | null; debit: string; credit: string }>(
        `WITH g AS (${GL}), dn AS (${DOC_NUMBERS})
         SELECT g.entry_id AS "entryId", g.entry_number AS "entryNumber", g.entry_date AS date,
                COALESCE(g.line_description, g.description) AS description, g.origin_type AS "referenceType", g.origin_id AS "referenceId",
                dn.number AS "documentNumber", g.debit::text, g.credit::text
           FROM g LEFT JOIN dn ON dn.id = g.origin_id WHERE ${filters} AND g.entry_date BETWEEN $4 AND $5
          ORDER BY g.entry_date, g.posted_at, g.entry_number, g.line_no LIMIT ${LIMIT + 1}`, params);
      const truncated = rows.length > LIMIT;
      return { companyId, dateFrom: q.dateFrom, dateTo: q.dateTo, account, ...runningLines(o!.opening, rows.slice(0, LIMIT), DEBIT_NORMAL.has(account.type)), truncated };
    });
  });

  /** Journal report: posted entries in a period with their lines. */
  app.get('/reports/journal', { preHandler: canView }, async (req) => {
    const q = parse(z.object({
      ...rangeFields, referenceType: z.string().regex(/^[A-Z_]+$/).optional(),
      limit: z.coerce.number().int().min(1).max(500).default(100), offset: z.coerce.number().int().min(0).default(0),
    }), req.query);
    checkRange(q);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const params = [req.auth!.tenantId, companyId, q.dateFrom, q.dateTo, q.referenceType ?? null];
      const where = `e.tenant_id = $1 AND e.company_id = $2 AND e.status IN ('POSTED', 'REVERSED') AND e.entry_date BETWEEN $3 AND $4
                     AND ($5::text IS NULL OR e.reference_type = $5)`;
      const { rows: [t] } = await db.query<{ count: string; debit: string; credit: string }>(
        `SELECT count(*)::text, COALESCE(sum(total_debit), 0)::text AS debit, COALESCE(sum(total_credit), 0)::text AS credit FROM journal_entries e WHERE ${where}`, params);
      const { rows: entries } = await db.query<{ id: string }>(
        `SELECT e.id, e.entry_number AS "entryNumber", e.entry_date AS date, e.description, e.reference_type AS "referenceType",
                e.status, e.total_debit::text AS "totalDebit", e.total_credit::text AS "totalCredit"
           FROM journal_entries e WHERE ${where}
          ORDER BY e.entry_date, e.posted_at, e.entry_number LIMIT $6 OFFSET $7`, [...params, q.limit, q.offset]);
      const { rows: lines } = await db.query<{ entryId: string }>(
        `SELECT l.journal_entry_id AS "entryId", l.line_no AS "lineNo", a.code AS "accountCode", a.name_ar AS "accountName",
                l.description, l.debit::text, l.credit::text
           FROM journal_entry_lines l JOIN accounts a ON a.id = l.account_id
          WHERE l.journal_entry_id = ANY($1) ORDER BY l.journal_entry_id, l.line_no`, [entries.map((e) => e.id)]);
      return {
        companyId, dateFrom: q.dateFrom, dateTo: q.dateTo,
        data: entries.map((e) => ({ ...e, lines: lines.filter((l) => l.entryId === e.id) })),
        total: Number(t!.count), totals: { debit: toMoney(t!.debit), credit: toMoney(t!.credit) },
      };
    });
  });
}

/** Adds a running balance. Debit-normal accounts show debit − credit; others credit − debit. */
export function runningLines<T extends { debit: string; credit: string }>(opening: string, rows: T[], debitNormal: boolean) {
  const sign = debitNormal ? 1 : -1;
  let balance = new Decimal(opening).times(sign);
  const lines = rows.map((r) => {
    balance = balance.plus(new Decimal(r.debit).minus(r.credit).times(sign));
    return { ...r, balance: toMoney(balance) };
  });
  return {
    openingBalance: toMoney(new Decimal(opening).times(sign)),
    lines,
    totals: { debit: toMoney(sum(rows.map((r) => r.debit))), credit: toMoney(sum(rows.map((r) => r.credit))) },
    closingBalance: toMoney(balance),
  };
}
