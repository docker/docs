import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { sum, toMoney } from '../../lib/money.js';
import { parse } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from '../accounting/company-context.js';
import { checkRange, GL, NOT_CLOSING, rangeFields } from './common.js';

/**
 * Operational sales and purchase reports, read from issued documents.
 * Returns subtract from the period of the return. Cancelled documents and
 * drafts are excluded. Financial totals (revenue, cost) belong to the P&L;
 * these reports explain them by customer, supplier, product or month.
 */
const SIDES = {
  sales: {
    party: 'customer', doc: 'sales_invoices', items: 'sales_invoice_items', ret: 'sales_returns', retItems: 'sales_return_items',
    out: 'SALE', back: 'SALE_RETURN', partyTable: 'customers',
  },
  purchases: {
    party: 'supplier', doc: 'purchase_invoices', items: 'purchase_invoice_items', ret: 'purchase_returns', retItems: 'purchase_return_items',
    out: 'PURCHASE', back: 'PURCHASE_RETURN', partyTable: 'suppliers',
  },
} as const;

const GROUP_BY = ['party', 'product', 'month', 'document'] as const;

export default async function commercialReportRoutes(app: FastifyInstance) {
  const canView = requirePermission(app, 'report.view');

  for (const [name, s] of Object.entries(SIDES)) {
    app.get(`/reports/${name}`, { preHandler: canView }, async (req) => {
      const q = parse(z.object({
        ...rangeFields, groupBy: z.enum(GROUP_BY).default('party'),
        partyId: z.uuid().optional(), productId: z.uuid().optional(), branchId: z.uuid().optional(),
      }), req.query);
      checkRange(q);
      return req.tenantTx(async (db) => {
        const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
        // For sales the cost is what the stock movements took out (COGS); returns bring it back.
        const lines = (table: string, items: string, sign: 1 | -1, movement: string) => `
          SELECT d.id AS doc_id, d.doc_number, d.doc_date, d.${s.party}_id AS party_id, i.product_id, ${sign} AS sign,
                 i.quantity, i.net_amount, i.vat_amount, i.total_amount,
                 COALESCE((SELECT sum(m.total_cost) FROM stock_movements m WHERE m.reference_line_id = i.id AND m.movement_type = '${movement}'), 0) AS cost
            FROM ${table} d JOIN ${items} i ON i.document_id = d.id
           WHERE d.tenant_id = $1 AND d.company_id = $2 AND d.status NOT IN ('DRAFT', 'CANCELLED') AND d.doc_date BETWEEN $3 AND $4
             AND ($5::uuid IS NULL OR d.${s.party}_id = $5) AND ($6::uuid IS NULL OR i.product_id = $6) AND ($7::uuid IS NULL OR d.branch_id = $7)`;
        const key = {
          party: { key: 'l.party_id::text', label: 'pt.name_ar', code: 'pt.code' },
          product: { key: `COALESCE(l.product_id::text, '')`, label: `COALESCE(p.name_ar, 'بنود بدون صنف (حسابات)')`, code: 'p.sku' },
          month: { key: `to_char(l.doc_date, 'YYYY-MM')`, label: `to_char(l.doc_date, 'YYYY-MM')`, code: 'NULL::text' },
          document: { key: 'l.doc_id::text', label: 'l.doc_number', code: 'l.doc_date::text' },
        }[q.groupBy];
        const { rows } = await db.query<{ key: string; label: string; code: string | null; quantity: string; netAmount: string; vatAmount: string; totalAmount: string; cost: string; documents: string; returns: string }>(
          `WITH l AS (${lines(s.doc, s.items, 1, s.out)} UNION ALL ${lines(s.ret, s.retItems, -1, s.back)})
           SELECT ${key.key} AS key, ${key.label} AS label, ${key.code} AS code,
                  sum(l.sign * l.quantity)::numeric(18,4)::text AS quantity,
                  sum(l.sign * l.net_amount)::numeric(18,2)::text AS "netAmount",
                  sum(l.sign * l.vat_amount)::numeric(18,2)::text AS "vatAmount",
                  sum(l.sign * l.total_amount)::numeric(18,2)::text AS "totalAmount",
                  sum(l.sign * l.cost)::numeric(18,2)::text AS cost,
                  count(DISTINCT l.doc_id) FILTER (WHERE l.sign = 1)::text AS documents,
                  count(DISTINCT l.doc_id) FILTER (WHERE l.sign = -1)::text AS returns
             FROM l LEFT JOIN ${s.partyTable} pt ON pt.id = l.party_id LEFT JOIN products p ON p.id = l.product_id
            GROUP BY 1, 2, 3 ORDER BY ${q.groupBy === 'month' || q.groupBy === 'document' ? '3, 2' : `sum(l.sign * l.net_amount) DESC`}`,
          [req.auth!.tenantId, companyId, q.dateFrom, q.dateTo, q.partyId ?? null, q.productId ?? null, q.branchId ?? null]);
        const isSales = name === 'sales';
        const data = rows.map((r) => ({
          ...r, documents: Number(r.documents), returns: Number(r.returns),
          // Gross profit only makes sense for sales; for purchases "cost" is the stock value received.
          ...(isSales ? { grossProfit: toMoney(sum([r.netAmount]).minus(r.cost)) } : {}),
        }));
        const total = (k: 'netAmount' | 'vatAmount' | 'totalAmount' | 'cost') => toMoney(sum(data.map((r) => r[k])));
        const totals = { netAmount: total('netAmount'), vatAmount: total('vatAmount'), totalAmount: total('totalAmount'), cost: total('cost') };
        return {
          companyId, dateFrom: q.dateFrom, dateTo: q.dateTo, groupBy: q.groupBy, data,
          totals: isSales ? { ...totals, grossProfit: toMoney(sum([totals.netAmount]).minus(totals.cost)) } : totals,
        };
      });
    });
  }

  /** Expenses from the ledger (vouchers, purchase-invoice expense lines and manual entries alike). */
  app.get('/reports/expenses', { preHandler: requirePermission(app, 'financial_report.view') }, async (req) => {
    const q = parse(z.object({
      ...rangeFields, groupBy: z.enum(['account', 'costCenter', 'month']).default('account'),
      branchId: z.uuid().optional(), costCenterId: z.uuid().optional(), accountId: z.uuid().optional(),
    }), req.query);
    checkRange(q);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const key = {
        account: { key: 'a.id::text', label: 'a.name_ar', code: 'a.code' },
        costCenter: { key: `COALESCE(cc.id::text, '')`, label: `COALESCE(cc.name, 'بدون مركز تكلفة')`, code: 'cc.code' },
        month: { key: `to_char(g.entry_date, 'YYYY-MM')`, label: `to_char(g.entry_date, 'YYYY-MM')`, code: 'NULL::text' },
      }[q.groupBy];
      const { rows } = await db.query<{ key: string; label: string; code: string | null; amount: string; entries: string }>(
        `WITH g AS (${GL})
         SELECT ${key.key} AS key, ${key.label} AS label, ${key.code} AS code,
                sum(g.debit - g.credit)::numeric(18,2)::text AS amount, count(DISTINCT g.entry_id)::text AS entries
           FROM g JOIN accounts a ON a.id = g.account_id LEFT JOIN cost_centers cc ON cc.id = g.cost_center_id
          WHERE a.account_type = 'EXPENSE' AND ${NOT_CLOSING} AND g.entry_date BETWEEN $3 AND $4
            AND ($5::uuid IS NULL OR g.branch_id = $5) AND ($6::uuid IS NULL OR g.cost_center_id = $6) AND ($7::uuid IS NULL OR g.account_id = $7)
          GROUP BY 1, 2, 3 HAVING sum(g.debit - g.credit) <> 0
          ORDER BY ${q.groupBy === 'month' ? '2' : 'sum(g.debit - g.credit) DESC'}`,
        [req.auth!.tenantId, companyId, q.dateFrom, q.dateTo, q.branchId ?? null, q.costCenterId ?? null, q.accountId ?? null]);
      const data = rows.map((r) => ({ ...r, entries: Number(r.entries) }));
      return { companyId, dateFrom: q.dateFrom, dateTo: q.dateTo, groupBy: q.groupBy, data, total: toMoney(sum(data.map((r) => r.amount))) };
    });
  });
}
