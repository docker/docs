import { describe, expect, it } from 'vitest';
import { CalcError, calculateDocument, calculateLine, totalsOf } from '../src/modules/documents/calc.js';

const S = (quantity: string, unitPrice: string, extra = {}) => ({ quantity, unitPrice, vatCategory: 'S' as const, vatRate: '0.15', ...extra });

describe('VAT calculation', () => {
  it('computes the reference example: 1000 SAR + 15% VAT = 1150', () => {
    const { totals } = calculateDocument([S('1', '1000')], false);
    expect(totals).toMatchObject({ subtotal: '1000.00', discountTotal: '0.00', taxableAmount: '1000.00', taxAmount: '150.00', total: '1150.00' });
  });

  it('applies line discounts before VAT', () => {
    const { lines, totals } = calculateDocument([S('4', '250', { discountPercent: '10' }), S('2', '50', { discountAmount: '5' })], false);
    expect(lines[0]).toMatchObject({ grossAmount: '1000.00', discountAmount: '100.00', netAmount: '900.00', vatAmount: '135.00' });
    expect(totals).toMatchObject({ subtotal: '1100.00', discountTotal: '105.00', taxableAmount: '995.00', taxAmount: '149.25', total: '1144.25' });
  });

  it('rounds VAT once per category, not per line', () => {
    // Three lines of 0.10: per-line VAT rounds 0.015 → 0.02 each (0.06), document VAT is round(0.30 × 0.15) = 0.05.
    const { lines, totals } = calculateDocument([S('1', '0.10'), S('1', '0.10'), S('1', '0.10')], false);
    expect(lines.map((l) => l.vatAmount)).toEqual(['0.02', '0.02', '0.02']);
    expect(totals.taxAmount).toBe('0.05');
    expect(totals.total).toBe('0.35');
  });

  it('extracts VAT from VAT-inclusive prices', () => {
    const { lines, totals } = calculateDocument([S('1', '115'), S('3', '10')], true);
    expect(lines[0]).toMatchObject({ netAmount: '100.00', vatAmount: '15.00' });
    expect(lines[1]).toMatchObject({ netAmount: '26.09' });
    expect(totals).toMatchObject({ taxableAmount: '126.09', taxAmount: '18.91', total: '145.00' });
  });

  it('handles quantities and prices with 4 decimals exactly', () => {
    const line = calculateLine(S('2.5', '4.1250'), false, 1);
    expect(line).toMatchObject({ grossAmount: '10.31', netAmount: '10.31', vatAmount: '1.55' });
  });

  it('separates zero-rated, exempt and standard lines', () => {
    const { totals } = calculateDocument([
      S('1', '100'),
      { quantity: '1', unitPrice: '200', vatCategory: 'Z', vatRate: '0' },
      { quantity: '1', unitPrice: '50', vatCategory: 'E', vatRate: '0' },
    ], false);
    expect(totals.taxBreakdown).toEqual([
      { vatCategory: 'S', vatRate: '0.15', taxableAmount: '100.00', taxAmount: '15.00' },
      { vatCategory: 'Z', vatRate: '0', taxableAmount: '200.00', taxAmount: '0.00' },
      { vatCategory: 'E', vatRate: '0', taxableAmount: '50.00', taxAmount: '0.00' },
    ]);
    expect(totals.total).toBe('365.00');
  });

  it('rejects invalid lines', () => {
    const bad = [
      S('0', '10'), S('-1', '10'), S('1', '-5'), S('1', '10', { discountAmount: '11' }),
      S('1', '10', { discountPercent: '101' }), S('1', '10', { discountAmount: '1', discountPercent: '5' }),
      { quantity: '1', unitPrice: '10', vatCategory: 'Z' as const, vatRate: '0.15' },
      { quantity: '1', unitPrice: '10', vatCategory: 'S' as const, vatRate: '0' },
    ];
    for (const l of bad) expect(() => calculateLine(l, false, 1), JSON.stringify(l)).toThrow(CalcError);
  });

  it('totals arbitrary stored lines (used for returns)', () => {
    const t = totalsOf([{ grossAmount: '500.00', discountAmount: '0.00', netAmount: '500.00', vatCategory: 'S', vatRate: '0.15' }]);
    expect(t.total).toBe('575.00');
  });
});
