import { z } from 'zod';
import type { Db } from '../../db/tx.js';
import { isoDate, todayIn } from '../../lib/dates.js';
import { badRequest } from '../../lib/errors.js';

/**
 * Shared building blocks for the reporting engine. Every financial figure is
 * read from posted journal lines; documents are only used for operational
 * reports (sales and purchases by customer or product) and for aging.
 */

export const rangeFields = { companyId: z.uuid().optional(), dateFrom: isoDate, dateTo: isoDate };

/** Longest period a report may cover. */
export const MAX_RANGE_YEARS = 5;

export function checkRange(q: { dateFrom: string; dateTo: string }) {
  if (q.dateFrom > q.dateTo) throw badRequest('INVALID_RANGE', 'dateFrom must be on or before dateTo');
  const limit = new Date(`${q.dateFrom}T00:00:00Z`);
  limit.setUTCFullYear(limit.getUTCFullYear() + MAX_RANGE_YEARS);
  if (new Date(`${q.dateTo}T00:00:00Z`) > limit) throw badRequest('RANGE_TOO_LONG', `A report can cover at most ${MAX_RANGE_YEARS} years`);
}

/**
 * Posted ledger lines of one company ($1 tenant, $2 company). A reversal
 * inherits the origin of the entry it reverses, so a cancelled document and
 * its reversal can be recognised and netted together. Entries with status
 * REVERSED stay in the ledger alongside their reversal; both count.
 */
export const GL = `
  SELECT l.id AS line_id, l.account_id, l.debit, l.credit, l.branch_id, l.cost_center_id, l.customer_id, l.supplier_id,
         l.line_no, l.description AS line_description,
         e.id AS entry_id, e.entry_number, e.entry_date, e.posted_at, e.description, e.reference_type, e.reference_id,
         COALESCE(o.reference_type, e.reference_type) AS origin_type, COALESCE(o.reference_id, e.reference_id) AS origin_id
    FROM journal_entry_lines l
    JOIN journal_entries e ON e.id = l.journal_entry_id
    LEFT JOIN journal_entries o ON o.id = e.reversal_of_id
   WHERE e.tenant_id = $1 AND e.company_id = $2 AND e.status IN ('POSTED', 'REVERSED')`;

/** Year-closing entries move P&L balances into retained earnings; income reports exclude them. */
export const NOT_CLOSING = `g.origin_type <> 'YEAR_CLOSING'`;

/**
 * Document numbers of the company ($1 tenant, $2 company), keyed by id, to be
 * LEFT JOINed on a ledger line's origin_id. A join rather than a correlated
 * subquery per line: under RLS those subqueries cannot use the primary-key
 * index and turned into a sequential scan per ledger line.
 */
export const DOC_NUMBERS = `
  SELECT id, doc_number AS number FROM sales_invoices WHERE tenant_id = $1 AND company_id = $2
  UNION ALL SELECT id, doc_number FROM sales_returns WHERE tenant_id = $1 AND company_id = $2
  UNION ALL SELECT id, doc_number FROM purchase_invoices WHERE tenant_id = $1 AND company_id = $2
  UNION ALL SELECT id, doc_number FROM purchase_returns WHERE tenant_id = $1 AND company_id = $2
  UNION ALL SELECT id, payment_number FROM payments WHERE tenant_id = $1 AND company_id = $2
  UNION ALL SELECT id, expense_number FROM expenses WHERE tenant_id = $1 AND company_id = $2
  UNION ALL SELECT id, adjustment_number FROM stock_adjustments WHERE tenant_id = $1 AND company_id = $2`;

/**
 * Cash and cash equivalents: the system cash and bank accounts plus every
 * account behind a payment method (cards, wallets and other clearing accounts).
 */
export const CASH_ACCOUNTS = `
  SELECT a.id FROM accounts a
   WHERE a.company_id = $2 AND (a.system_key IN ('CASH', 'BANK') OR a.id IN (SELECT account_id FROM payment_methods WHERE company_id = $2))`;

/** Today in the company's time zone (Asia/Riyadh by default). */
export async function companyToday(db: Db, companyId: string): Promise<string> {
  const { rows: [t] } = await db.query<{ timezone: string }>(`SELECT timezone FROM companies WHERE id = $1`, [companyId]);
  return todayIn(t?.timezone ?? 'Asia/Riyadh');
}
