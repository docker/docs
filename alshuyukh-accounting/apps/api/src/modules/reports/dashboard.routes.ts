import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { addMonths, isoDate } from '../../lib/dates.js';
import { Decimal, toMoney } from '../../lib/money.js';
import { parse } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from '../accounting/company-context.js';
import { CASH_ACCOUNTS, checkRange, companyToday, GL, NOT_CLOSING } from './common.js';

/**
 * Company dashboard. Every number comes from posted journal lines:
 * profit from revenue/cost/expense accounts (never from invoices),
 * balances from the control accounts' sub-ledgers and cash accounts.
 */
export default async function dashboardRoutes(app: FastifyInstance) {
  app.get('/dashboard', { preHandler: requirePermission(app, 'financial_report.view') }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional(), dateFrom: isoDate.optional(), dateTo: isoDate.optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const today = await companyToday(db, companyId);
      // Default period: the fiscal year that contains today (or the calendar year), up to today.
      let dateFrom = q.dateFrom;
      const dateTo = q.dateTo ?? today;
      if (!dateFrom) {
        const { rows: [fy] } = await db.query<{ start_date: string }>(
          `SELECT start_date FROM fiscal_years WHERE company_id = $1 AND $2::date BETWEEN start_date AND end_date`, [companyId, dateTo]);
        dateFrom = fy?.start_date ?? `${dateTo.slice(0, 4)}-01-01`;
      }
      checkRange({ dateFrom, dateTo });
      const params = [req.auth!.tenantId, companyId, dateFrom, dateTo];

      const { rows: [k] } = await db.query<Record<string, string>>(
        `WITH g AS (${GL}), cash AS (${CASH_ACCOUNTS}),
              cash_only AS (SELECT a.id FROM accounts a WHERE a.company_id = $2 AND (a.system_key = 'CASH'
                             OR a.id IN (SELECT account_id FROM payment_methods WHERE company_id = $2 AND method_type = 'CASH')))
         SELECT
           COALESCE(sum(g.credit - g.debit) FILTER (WHERE a.account_type = 'REVENUE' AND in_period AND ${NOT_CLOSING}), 0) AS sales,
           COALESCE(sum(g.debit - g.credit) FILTER (WHERE a.account_type = 'COST_OF_GOODS_SOLD' AND in_period AND ${NOT_CLOSING}), 0) AS cogs,
           COALESCE(sum(g.debit - g.credit) FILTER (WHERE a.account_type = 'EXPENSE' AND in_period AND ${NOT_CLOSING}), 0) AS expenses,
           -- Purchases net of VAT: what purchase invoices and returns put into stock, expense and cost accounts.
           COALESCE(sum(g.debit - g.credit) FILTER (WHERE in_period AND g.origin_type IN ('PURCHASE_INVOICE', 'PURCHASE_RETURN')
                    AND g.supplier_id IS NULL AND a.system_key IS DISTINCT FROM 'VAT_INPUT'), 0) AS purchases,
           COALESCE(sum(g.debit - g.credit) FILTER (WHERE g.customer_id IS NOT NULL), 0) AS receivables,
           COALESCE(sum(g.credit - g.debit) FILTER (WHERE g.supplier_id IS NOT NULL), 0) AS payables,
           COALESCE(sum(g.debit - g.credit) FILTER (WHERE g.account_id IN (SELECT id FROM cash_only)), 0) AS cash,
           COALESCE(sum(g.debit - g.credit) FILTER (WHERE g.account_id IN (SELECT id FROM cash) AND g.account_id NOT IN (SELECT id FROM cash_only)), 0) AS bank,
           COALESCE(sum(g.debit - g.credit) FILTER (WHERE a.system_key = 'INVENTORY'), 0) AS inventory,
           COALESCE(sum(g.credit - g.debit) FILTER (WHERE a.system_key IN ('VAT_OUTPUT', 'VAT_INPUT')), 0) AS vat
         FROM (SELECT g.*, g.entry_date >= $3 AS in_period FROM g WHERE g.entry_date <= $4) g JOIN accounts a ON a.id = g.account_id`, params);

      // Twelve months ending with the month of dateTo.
      const firstMonth = addMonths(`${dateTo.slice(0, 7)}-01`, -11);
      const { rows: monthly } = await db.query<{ month: string; revenue: string; costs: string; cash: string }>(
        `WITH g AS (${GL})
         SELECT to_char(g.entry_date, 'YYYY-MM') AS month,
                COALESCE(sum(g.credit - g.debit) FILTER (WHERE a.account_type = 'REVENUE' AND ${NOT_CLOSING}), 0)::text AS revenue,
                COALESCE(sum(g.debit - g.credit) FILTER (WHERE a.account_type IN ('EXPENSE', 'COST_OF_GOODS_SOLD') AND ${NOT_CLOSING}), 0)::text AS costs,
                COALESCE(sum(g.debit - g.credit) FILTER (WHERE g.account_id IN (${CASH_ACCOUNTS})), 0)::text AS cash
           FROM g JOIN accounts a ON a.id = g.account_id
          WHERE g.entry_date BETWEEN $3 AND $4 GROUP BY 1`, [req.auth!.tenantId, companyId, firstMonth, dateTo]);
      const { rows: [open] } = await db.query<{ cash: string }>(
        `WITH g AS (${GL}) SELECT COALESCE(sum(g.debit - g.credit), 0)::text AS cash FROM g WHERE g.entry_date < $3 AND g.account_id IN (${CASH_ACCOUNTS})`,
        [req.auth!.tenantId, companyId, firstMonth]);
      let cash = new Decimal(open!.cash);
      const months = Array.from({ length: 12 }, (_, i) => addMonths(firstMonth, i).slice(0, 7)).map((month) => {
        const m = monthly.find((r) => r.month === month);
        cash = cash.plus(m?.cash ?? 0);
        const revenue = new Decimal(m?.revenue ?? 0);
        const costs = new Decimal(m?.costs ?? 0);
        return { month, revenue: toMoney(revenue), costs: toMoney(costs), netProfit: toMoney(revenue.minus(costs)), cashBalance: toMoney(cash) };
      });

      const money = (key: string) => toMoney(k![key] ?? 0);
      return {
        companyId, dateFrom, dateTo,
        kpis: {
          totalSales: money('sales'), totalPurchases: money('purchases'), costOfSales: money('cogs'), expenses: money('expenses'),
          netProfit: toMoney(new Decimal(k!.sales!).minus(k!.cogs!).minus(k!.expenses!)),
          receivables: money('receivables'), payables: money('payables'), cash: money('cash'), bank: money('bank'),
          inventoryValue: money('inventory'), vatPayable: money('vat'),
        },
        monthly: months,
      };
    });
  });
}
