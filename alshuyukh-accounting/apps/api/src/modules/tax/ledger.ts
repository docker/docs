import type { Db } from '../../db/tx.js';
import { Decimal, toMoney } from '../../lib/money.js';
import type { VatCategory } from '../documents/calc.js';

/**
 * VAT sub-ledger. Every taxable document writes one row per (VAT category,
 * rate) when it posts; cancelling writes the negated rows. The VAT return is
 * built only from these rows.
 */

export type TaxSource = 'SALES_INVOICE' | 'SALES_RETURN' | 'PURCHASE_INVOICE' | 'PURCHASE_RETURN' | 'EXPENSE';

export interface TaxGroup { vatCategory: VatCategory; vatRate: string; taxableAmount: string; taxAmount: string }

/**
 * Groups document lines by (category, rate). The VAT of each group is
 * round(taxable × rate); if the document's VAT differs (a capped final
 * return), the difference goes to the largest standard-rated group so the
 * groups always add up to the posted VAT.
 */
export function taxGroups(lines: { netAmount: string; vatCategory: VatCategory; vatRate: string }[], documentTax: string): TaxGroup[] {
  const map = new Map<string, { vatCategory: VatCategory; vatRate: Decimal; taxable: Decimal }>();
  for (const l of lines) {
    const rate = new Decimal(l.vatRate);
    const key = `${l.vatCategory}|${rate.toString()}`;
    const g = map.get(key) ?? { vatCategory: l.vatCategory, vatRate: rate, taxable: new Decimal(0) };
    g.taxable = g.taxable.plus(l.netAmount);
    map.set(key, g);
  }
  const groups = [...map.values()].map((g) => ({
    vatCategory: g.vatCategory, vatRate: g.vatRate.toString(), taxable: g.taxable,
    tax: g.taxable.times(g.vatRate).toDecimalPlaces(2, Decimal.ROUND_HALF_UP),
  }));
  const diff = new Decimal(documentTax).minus(groups.reduce((s, g) => s.plus(g.tax), new Decimal(0)));
  if (!diff.isZero()) {
    const target = groups.filter((g) => g.vatCategory === 'S').sort((a, b) => b.taxable.comparedTo(a.taxable))[0];
    if (target) target.tax = target.tax.plus(diff);
  }
  return groups.map((g) => ({ vatCategory: g.vatCategory, vatRate: g.vatRate, taxableAmount: toMoney(g.taxable), taxAmount: toMoney(g.tax) }));
}

export interface RecordInput {
  tenantId: string;
  companyId: string;
  sourceType: TaxSource;
  sourceId: string;
  sourceNumber: string | null;
  journalEntryId: string;
  date: string;
  groups: TaxGroup[];
  partyName?: string | null;
  partyVatNumber?: string | null;
}

const DIRECTION: Record<TaxSource, 'OUTPUT' | 'INPUT'> = {
  SALES_INVOICE: 'OUTPUT', SALES_RETURN: 'OUTPUT', PURCHASE_INVOICE: 'INPUT', PURCHASE_RETURN: 'INPUT', EXPENSE: 'INPUT',
};
const IS_RETURN = (s: TaxSource) => s === 'SALES_RETURN' || s === 'PURCHASE_RETURN';

export async function recordTax(db: Db, input: RecordInput): Promise<void> {
  const sign = IS_RETURN(input.sourceType) ? -1 : 1;
  for (const g of input.groups) {
    await db.query(
      `INSERT INTO tax_transactions (tenant_id, company_id, direction, source_type, source_id, source_number, journal_entry_id,
         transaction_date, vat_category, vat_rate, taxable_amount, tax_amount, is_adjustment, party_name, party_vat_number)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
      [input.tenantId, input.companyId, DIRECTION[input.sourceType], input.sourceType, input.sourceId, input.sourceNumber,
       input.journalEntryId, input.date, g.vatCategory, g.vatRate,
       toMoney(new Decimal(g.taxableAmount).times(sign)), toMoney(new Decimal(g.taxAmount).times(sign)),
       IS_RETURN(input.sourceType), input.partyName ?? null, input.partyVatNumber ?? null]);
  }
}

/** Writes the negated rows of a cancelled document, dated like the originals. */
export async function reverseTax(db: Db, sourceType: TaxSource, sourceId: string): Promise<void> {
  await db.query(
    `INSERT INTO tax_transactions (tenant_id, company_id, direction, source_type, source_id, source_number, journal_entry_id,
       transaction_date, vat_category, vat_rate, taxable_amount, tax_amount, is_adjustment, reverses_id, party_name, party_vat_number)
     SELECT t.tenant_id, t.company_id, t.direction, t.source_type, t.source_id, t.source_number, t.journal_entry_id,
            t.transaction_date, t.vat_category, t.vat_rate, -t.taxable_amount, -t.tax_amount, t.is_adjustment, t.id,
            t.party_name, t.party_vat_number
       FROM tax_transactions t
      WHERE t.source_type = $1 AND t.source_id = $2 AND t.reverses_id IS NULL
        AND NOT EXISTS (SELECT 1 FROM tax_transactions r WHERE r.reverses_id = t.id)`,
    [sourceType, sourceId]);
}
