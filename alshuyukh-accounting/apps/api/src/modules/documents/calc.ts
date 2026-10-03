import { Decimal, toMoney } from '../../lib/money.js';

/**
 * Line and document totals for every commercial document (quotes, invoices,
 * returns, purchase orders and bills). Pure functions: no database access.
 *
 * Rounding (ZATCA-consistent):
 *  - Each line is rounded to 2 decimals: gross, discount, taxable (net).
 *  - Document VAT is computed once per (category, rate) group:
 *      round(sum of line taxable amounts × rate, 2)
 *    Line VAT amounts are informative and may differ from the document VAT
 *    by a halala; the document VAT is what is posted and reported.
 *  - When prices include VAT, the line taxable amount is
 *      round((qty × price − discount) / (1 + rate), 2)
 *    so all stored line amounts are VAT-exclusive.
 */

export type VatCategory = 'S' | 'Z' | 'E' | 'O';

export interface CalcLineInput {
  quantity: string;          // up to 4 decimals
  unitPrice: string;         // up to 4 decimals, in the price basis of the document
  discountAmount?: string;   // in the price basis, up to 2 decimals
  discountPercent?: string;  // 0–100, up to 2 decimals
  vatCategory: VatCategory;
  vatRate: string;           // fraction, e.g. "0.15"
}

export interface CalcLine {
  quantity: string;
  unitPrice: string;
  grossAmount: string;       // VAT-exclusive, before discount
  discountAmount: string;    // VAT-exclusive
  netAmount: string;         // taxable amount
  vatCategory: VatCategory;
  vatRate: string;
  vatAmount: string;         // informative line VAT
  totalAmount: string;       // net + line VAT
}

export interface CalcTotals {
  subtotal: string;
  discountTotal: string;
  taxableAmount: string;
  taxAmount: string;
  total: string;
  taxBreakdown: { vatCategory: VatCategory; vatRate: string; taxableAmount: string; taxAmount: string }[];
}

export class CalcError extends Error {
  constructor(public readonly line: number, message: string) {
    super(`Line ${line}: ${message}`);
  }
}

const round2 = (d: Decimal) => d.toDecimalPlaces(2, Decimal.ROUND_HALF_UP);

export function calculateLine(input: CalcLineInput, pricesIncludeVat: boolean, lineNo: number): CalcLine {
  const qty = new Decimal(input.quantity);
  const price = new Decimal(input.unitPrice);
  const rate = new Decimal(input.vatRate);
  if (!qty.greaterThan(0)) throw new CalcError(lineNo, 'quantity must be greater than zero');
  if (price.isNegative()) throw new CalcError(lineNo, 'unit price cannot be negative');
  if (rate.isNegative() || rate.greaterThan(1)) throw new CalcError(lineNo, 'invalid VAT rate');
  if (input.vatCategory !== 'S' && !rate.isZero()) throw new CalcError(lineNo, `VAT category ${input.vatCategory} must have a zero rate`);
  if (input.vatCategory === 'S' && rate.isZero()) throw new CalcError(lineNo, 'standard-rated lines need a positive VAT rate');

  const grossBasis = round2(qty.times(price));
  let discountBasis = new Decimal(0);
  if (input.discountAmount && input.discountPercent) throw new CalcError(lineNo, 'use either a discount amount or a discount percent');
  if (input.discountAmount) discountBasis = new Decimal(input.discountAmount);
  if (input.discountPercent) {
    const pct = new Decimal(input.discountPercent);
    if (pct.isNegative() || pct.greaterThan(100)) throw new CalcError(lineNo, 'discount percent must be between 0 and 100');
    discountBasis = round2(grossBasis.times(pct).dividedBy(100));
  }
  if (discountBasis.isNegative()) throw new CalcError(lineNo, 'discount cannot be negative');
  if (discountBasis.greaterThan(grossBasis)) throw new CalcError(lineNo, 'discount cannot exceed the line amount');

  let gross: Decimal;
  let discount: Decimal;
  let net: Decimal;
  if (pricesIncludeVat) {
    const divisor = rate.plus(1);
    net = round2(grossBasis.minus(discountBasis).dividedBy(divisor));
    discount = round2(discountBasis.dividedBy(divisor));
    gross = net.plus(discount);
  } else {
    gross = grossBasis;
    discount = discountBasis;
    net = gross.minus(discount);
  }
  const vat = round2(net.times(rate));
  return {
    quantity: qty.toString(),
    unitPrice: price.toString(),
    grossAmount: toMoney(gross),
    discountAmount: toMoney(discount),
    netAmount: toMoney(net),
    vatCategory: input.vatCategory,
    vatRate: rate.toString(),
    vatAmount: toMoney(vat),
    totalAmount: toMoney(net.plus(vat)),
  };
}

/** Document totals from already-calculated lines (also used for returns). */
export function totalsOf(lines: Pick<CalcLine, 'grossAmount' | 'discountAmount' | 'netAmount' | 'vatCategory' | 'vatRate'>[]): CalcTotals {
  const groups = new Map<string, { vatCategory: VatCategory; vatRate: Decimal; taxable: Decimal }>();
  let subtotal = new Decimal(0);
  let discount = new Decimal(0);
  for (const l of lines) {
    subtotal = subtotal.plus(l.grossAmount);
    discount = discount.plus(l.discountAmount);
    const rate = new Decimal(l.vatRate);
    const key = `${l.vatCategory}|${rate.toString()}`;
    const g = groups.get(key) ?? { vatCategory: l.vatCategory, vatRate: rate, taxable: new Decimal(0) };
    g.taxable = g.taxable.plus(l.netAmount);
    groups.set(key, g);
  }
  const taxBreakdown = [...groups.values()].map((g) => ({
    vatCategory: g.vatCategory,
    vatRate: g.vatRate.toString(),
    taxableAmount: toMoney(g.taxable),
    taxAmount: toMoney(round2(g.taxable.times(g.vatRate))),
  }));
  const taxable = subtotal.minus(discount);
  const tax = taxBreakdown.reduce((s, g) => s.plus(g.taxAmount), new Decimal(0));
  return {
    subtotal: toMoney(subtotal),
    discountTotal: toMoney(discount),
    taxableAmount: toMoney(taxable),
    taxAmount: toMoney(tax),
    total: toMoney(taxable.plus(tax)),
    taxBreakdown,
  };
}

export function calculateDocument(lines: CalcLineInput[], pricesIncludeVat: boolean) {
  const calculated = lines.map((l, i) => calculateLine(l, pricesIncludeVat, i + 1));
  return { lines: calculated, totals: totalsOf(calculated) };
}
