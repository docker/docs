/**
 * Display-only helpers. Amounts arrive from the API as strings ("1150.00").
 * Totals shown while typing are computed in integer halalas, never floats;
 * the server recomputes and validates everything.
 */
export function formatAmount(value: string | null | undefined): string {
  if (value === null || value === undefined || value === '') return '—';
  const negative = value.startsWith('-');
  const [int, frac = ''] = value.replace('-', '').split('.');
  const grouped = int!.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  return `${negative ? '-' : ''}${grouped}.${frac.padEnd(2, '0').slice(0, 2)}`;
}

/** "12.5" → 1250 halalas. Returns null for invalid input. */
export function toHalalas(value: string): number | null {
  const v = value.trim();
  if (v === '') return 0;
  if (!/^\d{1,13}(\.\d{1,2})?$/.test(v)) return null;
  const [int, frac = ''] = v.split('.');
  return Number(int) * 100 + Number(frac.padEnd(2, '0'));
}

export const fromHalalas = (h: number) =>
  `${h < 0 ? '-' : ''}${Math.floor(Math.abs(h) / 100)}.${String(Math.abs(h) % 100).padStart(2, '0')}`;

export const ACCOUNT_TYPE_AR: Record<string, string> = {
  ASSET: 'أصول', LIABILITY: 'التزامات', EQUITY: 'حقوق ملكية', REVENUE: 'إيرادات', EXPENSE: 'مصروفات', COST_OF_GOODS_SOLD: 'تكلفة المبيعات',
};
export const ENTRY_STATUS_AR: Record<string, string> = { DRAFT: 'مسودة', POSTED: 'مرحّل', REVERSED: 'معكوس' };
export const REFERENCE_AR: Record<string, string> = { MANUAL: 'يدوي', REVERSAL: 'قيد عكسي', YEAR_CLOSING: 'إقفال سنة', SALES_INVOICE: 'فاتورة مبيعات' };
