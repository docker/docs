import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { isoDate } from '../../lib/dates.js';
import { badRequest } from '../../lib/errors.js';
import { Decimal, sum, toMoney } from '../../lib/money.js';
import { parse } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from './company-context.js';

/**
 * Trial balance from posted journal lines only (reversed entries and their
 * reversals both count, so they cancel out). Other reports arrive in Phase 7.
 */
export default async function trialBalanceRoutes(app: FastifyInstance) {
  app.get('/reports/trial-balance', { preHandler: requirePermission(app, 'report.view') }, async (req) => {
    const q = parse(z.object({
      companyId: z.uuid().optional(),
      dateFrom: isoDate,
      dateTo: isoDate,
      includeZero: z.enum(['true', 'false']).optional(),
      branchId: z.uuid().optional(),
    }), req.query);
    if (q.dateFrom > q.dateTo) throw badRequest('INVALID_RANGE', 'dateFrom must be on or before dateTo');
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query<{
        accountId: string; code: string; nameAr: string; type: string;
        openingDebit: string; openingCredit: string; periodDebit: string; periodCredit: string;
      }>(
        `WITH movements AS (
           SELECT l.account_id,
                  sum(CASE WHEN e.entry_date < $3 THEN l.debit - l.credit ELSE 0 END) AS opening,
                  sum(CASE WHEN e.entry_date >= $3 THEN l.debit ELSE 0 END) AS period_debit,
                  sum(CASE WHEN e.entry_date >= $3 THEN l.credit ELSE 0 END) AS period_credit
             FROM journal_entry_lines l JOIN journal_entries e ON e.id = l.journal_entry_id
            WHERE e.tenant_id = $1 AND e.company_id = $2 AND e.status IN ('POSTED', 'REVERSED')
              AND e.entry_date <= $4 AND ($6::uuid IS NULL OR l.branch_id = $6)
            GROUP BY l.account_id)
         SELECT a.id AS "accountId", a.code, a.name_ar AS "nameAr", a.account_type AS type,
                GREATEST(COALESCE(m.opening, 0), 0)::numeric(18,2)::text AS "openingDebit",
                GREATEST(-COALESCE(m.opening, 0), 0)::numeric(18,2)::text AS "openingCredit",
                COALESCE(m.period_debit, 0)::numeric(18,2)::text AS "periodDebit",
                COALESCE(m.period_credit, 0)::numeric(18,2)::text AS "periodCredit"
           FROM accounts a LEFT JOIN movements m ON m.account_id = a.id
          WHERE a.tenant_id = $1 AND a.company_id = $2 AND a.is_postable AND a.deleted_at IS NULL
            AND ($5::boolean OR m.account_id IS NOT NULL)
          ORDER BY a.code`,
        [req.auth!.tenantId, companyId, q.dateFrom, q.dateTo, q.includeZero === 'true', q.branchId ?? null]);

      const data = rows.map((r) => {
        const closing = new Decimal(r.openingDebit).minus(r.openingCredit).plus(r.periodDebit).minus(r.periodCredit);
        return {
          ...r,
          closingDebit: toMoney(Decimal.max(closing, 0)),
          closingCredit: toMoney(Decimal.max(closing.neg(), 0)),
        };
      });
      const total = (k: keyof (typeof data)[number]) => toMoney(sum(data.map((r) => r[k] as string)));
      const totals = {
        openingDebit: total('openingDebit'), openingCredit: total('openingCredit'),
        periodDebit: total('periodDebit'), periodCredit: total('periodCredit'),
        closingDebit: total('closingDebit'), closingCredit: total('closingCredit'),
      };
      return {
        companyId, dateFrom: q.dateFrom, dateTo: q.dateTo, data, totals,
        balanced: totals.openingDebit === totals.openingCredit && totals.periodDebit === totals.periodCredit && totals.closingDebit === totals.closingCredit,
      };
    });
  });
}
