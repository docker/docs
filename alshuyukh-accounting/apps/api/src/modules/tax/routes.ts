import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { isoDate } from '../../lib/dates.js';
import { badRequest } from '../../lib/errors.js';
import { Decimal, toMoney } from '../../lib/money.js';
import { parse } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from '../accounting/company-context.js';

interface Box { amount: string; adjustment: string; vat: string }

/**
 * VAT return in the layout of the ZATCA form. Built from tax_transactions
 * and reconciled with the VAT accounts in the general ledger: any manual
 * journal entry to a VAT account shows up as a difference.
 */
export default async function taxRoutes(app: FastifyInstance) {
  const range = z.object({ companyId: z.uuid().optional(), dateFrom: isoDate, dateTo: isoDate });

  app.get('/reports/vat-return', { preHandler: requirePermission(app, 'report.view') }, async (req) => {
    const q = parse(range, req.query);
    if (q.dateFrom > q.dateTo) throw badRequest('INVALID_RANGE', 'dateFrom must be on or before dateTo');
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query<{ direction: string; vat_category: string; is_adjustment: boolean; taxable: string; tax: string }>(
        `SELECT direction, vat_category, is_adjustment, sum(taxable_amount)::text AS taxable, sum(tax_amount)::text AS tax
           FROM tax_transactions
          WHERE tenant_id = $1 AND company_id = $2 AND transaction_date BETWEEN $3 AND $4
          GROUP BY direction, vat_category, is_adjustment`,
        [req.auth!.tenantId, companyId, q.dateFrom, q.dateTo]);
      const box = (direction: string, category: string): Box => {
        const pick = (adj: boolean) => rows.find((r) => r.direction === direction && r.vat_category === category && r.is_adjustment === adj);
        const sale = pick(false);
        const adj = pick(true);
        return {
          amount: toMoney(sale?.taxable ?? 0),
          adjustment: toMoney(adj?.taxable ?? 0),
          vat: toMoney(new Decimal(sale?.tax ?? 0).plus(adj?.tax ?? 0)),
        };
      };
      const sum = (boxes: Box[]): Box => ({
        amount: toMoney(boxes.reduce((s, b) => s.plus(b.amount), new Decimal(0))),
        adjustment: toMoney(boxes.reduce((s, b) => s.plus(b.adjustment), new Decimal(0))),
        vat: toMoney(boxes.reduce((s, b) => s.plus(b.vat), new Decimal(0))),
      });
      const zero: Box = { amount: '0.00', adjustment: '0.00', vat: '0.00' };

      const sales = { standardRated: box('OUTPUT', 'S'), zeroRated: box('OUTPUT', 'Z'), exempt: box('OUTPUT', 'E') };
      const purchases = { standardRated: box('INPUT', 'S'), zeroRated: box('INPUT', 'Z'), exempt: box('INPUT', 'E') };
      const salesTotal = sum(Object.values(sales));
      const purchasesTotal = sum(Object.values(purchases));
      const netVat = new Decimal(salesTotal.vat).minus(purchasesTotal.vat);

      const { rows: [gl] } = await db.query<{ output: string; input: string }>(
        `SELECT COALESCE(sum(CASE WHEN a.system_key = 'VAT_OUTPUT' THEN l.credit - l.debit END), 0)::text AS output,
                COALESCE(sum(CASE WHEN a.system_key = 'VAT_INPUT' THEN l.debit - l.credit END), 0)::text AS input
           FROM journal_entry_lines l JOIN journal_entries e ON e.id = l.journal_entry_id JOIN accounts a ON a.id = l.account_id
          WHERE e.company_id = $1 AND e.status IN ('POSTED', 'REVERSED') AND e.entry_date BETWEEN $2 AND $3
            AND a.system_key IN ('VAT_OUTPUT', 'VAT_INPUT')`, [companyId, q.dateFrom, q.dateTo]);
      const outputDiff = new Decimal(gl!.output).minus(salesTotal.vat);
      const inputDiff = new Decimal(gl!.input).minus(purchasesTotal.vat);

      return {
        companyId, dateFrom: q.dateFrom, dateTo: q.dateTo,
        // Box numbers follow the ZATCA VAT return. Boxes this system cannot
        // yet separate (private health/education, exports, imports, reverse
        // charge) are reported as zero and listed in `notSupported`.
        sales: {
          '1_standardRated': sales.standardRated, '2_citizenHealthEducation': zero, '3_zeroRatedDomestic': sales.zeroRated,
          '4_exports': zero, '5_exempt': sales.exempt, '6_total': salesTotal,
        },
        purchases: {
          '7_standardRatedDomestic': purchases.standardRated, '8_importsPaidAtCustoms': zero, '9_importsReverseCharge': zero,
          '10_zeroRated': purchases.zeroRated, '11_exempt': purchases.exempt, '12_total': purchasesTotal,
        },
        '13_totalVatDue': toMoney(netVat),
        '14_previousPeriodCorrections': '0.00',
        '15_creditCarriedForward': '0.00',
        '16_netVatDue': toMoney(netVat),
        notSupported: ['2_citizenHealthEducation', '4_exports', '8_importsPaidAtCustoms', '9_importsReverseCharge', '14_previousPeriodCorrections', '15_creditCarriedForward'],
        ledger: {
          vatOutput: toMoney(gl!.output), vatInput: toMoney(gl!.input),
          outputDifference: toMoney(outputDiff), inputDifference: toMoney(inputDiff),
          reconciled: outputDiff.isZero() && inputDiff.isZero(),
        },
      };
    });
  });

  /** The tax transactions behind the return, for review and audit. */
  app.get('/reports/vat-transactions', { preHandler: requirePermission(app, 'report.view') }, async (req) => {
    const q = parse(range.extend({ direction: z.enum(['OUTPUT', 'INPUT']).optional(), vatCategory: z.enum(['S', 'Z', 'E', 'O']).optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query(
        `SELECT id, direction, source_type AS "sourceType", source_id AS "sourceId", source_number AS "sourceNumber",
                transaction_date AS date, vat_category AS "vatCategory", vat_rate::text AS "vatRate",
                taxable_amount::text AS "taxableAmount", tax_amount::text AS "taxAmount", is_adjustment AS "isAdjustment",
                reverses_id AS "reversesId", party_name AS "partyName", party_vat_number AS "partyVatNumber"
           FROM tax_transactions
          WHERE tenant_id = $1 AND company_id = $2 AND transaction_date BETWEEN $3 AND $4
            AND ($5::text IS NULL OR direction = $5) AND ($6::text IS NULL OR vat_category = $6)
          ORDER BY transaction_date, created_at`,
        [req.auth!.tenantId, companyId, q.dateFrom, q.dateTo, q.direction ?? null, q.vatCategory ?? null]);
      return { data: rows };
    });
  });
}
