import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { notFound } from '../../lib/errors.js';
import { Decimal, sum, toMoney } from '../../lib/money.js';
import { parse } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from '../accounting/company-context.js';
import { checkRange, companyToday, DOCUMENT_NUMBER, GL, rangeFields } from './common.js';
import { runningLines } from './financial.routes.js';

type Party = 'customer' | 'supplier';
const BUCKETS = ['current', 'days1to30', 'days31to60', 'days61to90', 'over90'] as const;
type Bucket = (typeof BUCKETS)[number];

const bucketOf = (daysOverdue: number): Bucket =>
  daysOverdue <= 0 ? 'current' : daysOverdue <= 30 ? 'days1to30' : daysOverdue <= 60 ? 'days31to60' : daysOverdue <= 90 ? 'days61to90' : 'over90';

/** Open documents that make up a party's balance, by party type. */
const OPEN_DOCS: Record<Party, string> = {
  customer: `
    SELECT 'SALES_INVOICE' AS type, id, doc_number AS number, doc_date AS date, COALESCE(due_date, doc_date) AS due, customer_id AS party_id, total, remaining_amount AS remaining
      FROM sales_invoices WHERE tenant_id = $1 AND company_id = $2 AND status IN ('ISSUED', 'PARTIALLY_PAID') AND remaining_amount > 0`,
  supplier: `
    SELECT 'PURCHASE_INVOICE' AS type, id, doc_number AS number, doc_date AS date, COALESCE(due_date, doc_date) AS due, supplier_id AS party_id, total, remaining_amount AS remaining
      FROM purchase_invoices WHERE tenant_id = $1 AND company_id = $2 AND status IN ('POSTED', 'PARTIALLY_PAID') AND remaining_amount > 0
    UNION ALL
    SELECT 'EXPENSE', id, expense_number, expense_date, expense_date, supplier_id, total, remaining_amount
      FROM expenses WHERE tenant_id = $1 AND company_id = $2 AND status IN ('POSTED', 'PARTIALLY_PAID') AND remaining_amount > 0
       AND supplier_id IS NOT NULL AND deleted_at IS NULL`,
};

export default async function partyReportRoutes(app: FastifyInstance) {
  const canView = requirePermission(app, 'report.view');

  /**
   * Receivables / payables aging at today's date. Open documents are bucketed
   * by days past due. The party balance comes from the ledger; whatever the
   * open documents do not explain (advance payments, unapplied credit notes,
   * manual entries) appears as "unapplied", so every row reconciles to the
   * ledger and the grand total equals the control accounts' sub-ledger.
   */
  for (const party of ['customer', 'supplier'] as const) {
    const path = party === 'customer' ? '/reports/receivables-aging' : '/reports/payables-aging';
    app.get(path, { preHandler: canView }, async (req) => {
      const q = parse(z.object({ companyId: z.uuid().optional(), partyId: z.uuid().optional() }), req.query);
      return req.tenantTx(async (db) => {
        const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
        const asOf = await companyToday(db, companyId);
        const params = [req.auth!.tenantId, companyId, q.partyId ?? null];
        const { rows: docs } = await db.query<{ type: string; id: string; number: string; date: string; due: string; party_id: string; total: string; remaining: string; overdue: number }>(
          `SELECT d.*, ($4::date - d.due) AS overdue FROM (${OPEN_DOCS[party]}) d
            WHERE ($3::uuid IS NULL OR d.party_id = $3) ORDER BY d.due, d.number`, [...params, asOf]);
        const col = party === 'customer' ? 'customer_id' : 'supplier_id';
        const { rows: balances } = await db.query<{ party_id: string; balance: string }>(
          `WITH g AS (${GL})
           SELECT g.${col} AS party_id, ${party === 'customer' ? 'sum(g.debit - g.credit)' : 'sum(g.credit - g.debit)'}::text AS balance
             FROM g WHERE g.${col} IS NOT NULL AND ($3::uuid IS NULL OR g.${col} = $3) GROUP BY g.${col}`, params);
        const { rows: parties } = await db.query<{ id: string; code: string; nameAr: string }>(
          `SELECT id, code, name_ar AS "nameAr" FROM ${party}s WHERE tenant_id = $1 AND company_id = $2 AND ($3::uuid IS NULL OR id = $3)`, params);

        const rows = parties.map((p) => {
          const own = docs.filter((d) => d.party_id === p.id);
          const buckets = Object.fromEntries(BUCKETS.map((b) => [b, new Decimal(0)])) as Record<Bucket, Decimal>;
          for (const d of own) buckets[bucketOf(d.overdue)] = buckets[bucketOf(d.overdue)].plus(d.remaining);
          const openTotal = sum(own.map((d) => d.remaining));
          const balance = new Decimal(balances.find((b) => b.party_id === p.id)?.balance ?? 0);
          return {
            partyId: p.id, code: p.code, nameAr: p.nameAr,
            ...Object.fromEntries(BUCKETS.map((b) => [b, toMoney(buckets[b])])) as Record<Bucket, string>,
            unapplied: toMoney(balance.minus(openTotal)),
            balance: toMoney(balance),
            documents: own.map((d) => ({ type: d.type, id: d.id, number: d.number, date: d.date, dueDate: d.due, total: d.total, remaining: d.remaining, daysOverdue: Math.max(d.overdue, 0) })),
          };
        }).filter((r) => r.documents.length || r.balance !== '0.00').sort((a, b) => a.code.localeCompare(b.code));

        const totals = Object.fromEntries([...BUCKETS, 'unapplied', 'balance'].map((k) =>
          [k, toMoney(sum(rows.map((r) => r[k as Bucket | 'unapplied' | 'balance'])))])) as Record<Bucket | 'unapplied' | 'balance', string>;
        return { companyId, asOf, data: rows, totals };
      });
    });
  }

  /** Customer or supplier statement: opening balance, every ledger movement, running balance. */
  for (const party of ['customer', 'supplier'] as const) {
    app.get(`/reports/${party}-statement`, { preHandler: canView }, async (req) => {
      const q = parse(z.object({ ...rangeFields, partyId: z.uuid() }), req.query);
      checkRange(q);
      return req.tenantTx(async (db) => {
        const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
        const { rows: [p] } = await db.query<{ id: string; code: string; nameAr: string; vatNumber: string | null }>(
          `SELECT id, code, name_ar AS "nameAr", vat_number AS "vatNumber" FROM ${party}s WHERE tenant_id = $1 AND company_id = $2 AND id = $3`,
          [req.auth!.tenantId, companyId, q.partyId]);
        if (!p) throw notFound(party === 'customer' ? 'Customer' : 'Supplier');
        const col = party === 'customer' ? 'customer_id' : 'supplier_id';
        const params = [req.auth!.tenantId, companyId, q.partyId, q.dateFrom, q.dateTo];
        const { rows: [o] } = await db.query<{ opening: string }>(
          `WITH g AS (${GL}) SELECT COALESCE(sum(g.debit - g.credit), 0)::text AS opening FROM g WHERE g.${col} = $3 AND g.entry_date < $4`, params.slice(0, 4));
        const { rows } = await db.query<{ debit: string; credit: string }>(
          `WITH g AS (${GL})
           SELECT g.entry_id AS "entryId", g.entry_number AS "entryNumber", g.entry_date AS date, g.description,
                  g.origin_type AS "referenceType", g.origin_id AS "referenceId", (g.reference_type = 'REVERSAL') AS "isReversal",
                  ${DOCUMENT_NUMBER} AS "documentNumber", g.debit::text, g.credit::text
             FROM g WHERE g.${col} = $3 AND g.entry_date BETWEEN $4 AND $5
            ORDER BY g.entry_date, g.posted_at, g.entry_number, g.line_no`, params);
        return { companyId, dateFrom: q.dateFrom, dateTo: q.dateTo, party: p, ...runningLines(o!.opening, rows, party === 'customer') };
      });
    });
  }
}
