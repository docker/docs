import { formatAmount } from '../../money';

export interface TotalsData { subtotal: string; discountTotal: string; taxableAmount: string; taxAmount: string; total: string }

export default function Totals({ t, extra }: { t: TotalsData; extra?: [string, string | undefined][] }) {
  const rows: [string, string | undefined, boolean?][] = [
    ['الإجمالي قبل الخصم', t.subtotal], ['الخصم', t.discountTotal], ['المبلغ الخاضع للضريبة', t.taxableAmount],
    ['ضريبة القيمة المضافة', t.taxAmount], ['الإجمالي شامل الضريبة', t.total, true],
    ...(extra ?? []).map(([a, b]) => [a, b] as [string, string | undefined]),
  ];
  return (
    <dl className="totals">
      {rows.filter(([, v]) => v !== undefined).map(([label, value, strong]) => (
        <div key={label} className={strong ? 'strong' : ''}><dt>{label}</dt><dd dir="ltr">{formatAmount(value)}</dd></div>
      ))}
    </dl>
  );
}
