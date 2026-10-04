import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { qty } from './shared';

interface Result { data: { productId: string; sku: string; productName: string; quantity: string; value: string }[]; stockValue: string; ledgerBalance: string; difference: string; reconciled: boolean }

export default function Valuation() {
  const v = useLoad(() => api<Result>('GET', '/api/inventory/valuation'));
  const d = v.data;
  return (
    <div className="card">
      <ErrorBox error={v.error} />
      {d && (
        <>
          <dl className="meta">
            <div><dt>قيمة المخزون</dt><dd dir="ltr">{formatAmount(d.stockValue)}</dd></div>
            <div><dt>رصيد حساب المخزون في الدفتر</dt><dd dir="ltr">{formatAmount(d.ledgerBalance)}</dd></div>
            <div><dt>المطابقة</dt><dd>{d.reconciled ? <span className="tag tag-ok">مطابق</span> : <span className="tag tag-off">فرق <span dir="ltr">{formatAmount(d.difference)}</span></span>}</dd></div>
          </dl>
          <table className="table">
            <thead><tr><th>الرمز</th><th>الصنف</th><th className="num">الكمية</th><th className="num">القيمة</th></tr></thead>
            <tbody>{d.data.map((r) => (
              <tr key={r.productId}><td dir="ltr" className="code">{r.sku}</td><td>{r.productName}</td><td className="num" dir="ltr">{qty(r.quantity)}</td><td className="num" dir="ltr">{formatAmount(r.value)}</td></tr>
            ))}</tbody>
          </table>
        </>
      )}
      <p className="muted small">التقييم بطريقة المتوسط المرجح لكل مستودع.</p>
    </div>
  );
}
