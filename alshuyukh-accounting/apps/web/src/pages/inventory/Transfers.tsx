import { useState, type FormEvent } from 'react';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { loadStockedProducts, loadWarehouses } from './shared';

interface Transfer { id: string; number: string; date: string; fromWarehouseName: string; toWarehouseName: string; totalCost: string; notes: string | null }
interface Line { productId: string; quantity: string }

const today = () => new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh' }).format(new Date());

export default function Transfers() {
  const { can } = useAuth();
  const list = useLoad(() => api<{ data: Transfer[] }>('GET', '/api/stock-transfers'));
  const warehouses = useLoad(loadWarehouses);
  const products = useLoad(loadStockedProducts);
  const [creating, setCreating] = useState(false);
  const [lines, setLines] = useState<Line[]>([{ productId: '', quantity: '' }]);
  const [error, setError] = useState<unknown>(null);

  async function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    setError(null);
    try {
      await api('POST', '/api/stock-transfers', {
        fromWarehouseId: f.get('from'), toWarehouseId: f.get('to'), date: f.get('date'), notes: f.get('notes') || null,
        lines: lines.filter((l) => l.productId && l.quantity.trim()).map((l) => ({ productId: l.productId, quantity: l.quantity.trim() })),
      });
      setCreating(false);
      setLines([{ productId: '', quantity: '' }]);
      list.reload();
    } catch (err) { setError(err); }
  }

  return (
    <>
      {creating && warehouses.data && products.data && (
        <form className="card form" onSubmit={submit}>
          <h2>تحويل مخزون بين المستودعات</h2>
          <ErrorBox error={error} />
          <div className="grid-3">
            <label>من مستودع<select name="from" required>{warehouses.data?.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}</select></label>
            <label>إلى مستودع<select name="to" required defaultValue={warehouses.data?.[1]?.id}>{warehouses.data?.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}</select></label>
            <label>التاريخ<input name="date" type="date" required defaultValue={today()} /></label>
          </div>
          <table className="table">
            <thead><tr><th>الصنف</th><th>الكمية</th><th /></tr></thead>
            <tbody>{lines.map((l, i) => (
              <tr key={i}>
                <td><select value={l.productId} aria-label={`صنف ${i + 1}`} onChange={(e) => setLines((ls) => ls.map((x, j) => (j === i ? { ...x, productId: e.target.value } : x)))}>
                  <option value="">— اختر —</option>
                  {products.data?.map((p) => <option key={p.id} value={p.id}>{p.sku} — {p.nameAr}</option>)}
                </select></td>
                <td><input dir="ltr" inputMode="decimal" value={l.quantity} aria-label={`كمية ${i + 1}`} onChange={(e) => setLines((ls) => ls.map((x, j) => (j === i ? { ...x, quantity: e.target.value } : x)))} /></td>
                <td>{lines.length > 1 && <button type="button" className="icon-btn" onClick={() => setLines((ls) => ls.filter((_, j) => j !== i))}>×</button>}</td>
              </tr>
            ))}</tbody>
          </table>
          <button type="button" className="btn btn-small" onClick={() => setLines((ls) => [...ls, { productId: '', quantity: '' }])}>+ سطر</button>
          <label>ملاحظات<input name="notes" /></label>
          <p className="muted small">ينتقل المخزون بمتوسط تكلفته في المستودع المصدر. لا يُنشأ قيد لأن القيمة تبقى في حساب المخزون.</p>
          <div className="form-actions"><button className="btn btn-primary">تنفيذ التحويل</button><button type="button" className="btn btn-ghost" onClick={() => setCreating(false)}>إلغاء</button></div>
        </form>
      )}
      <div className="card">
        <div className="toolbar"><span className="spacer" />{can('inventory.transfer') && !creating && <button className="btn btn-primary" onClick={() => setCreating(true)}>تحويل جديد</button>}</div>
        <ErrorBox error={list.error} />
        <table className="table">
          <thead><tr><th>الرقم</th><th>التاريخ</th><th>من</th><th>إلى</th><th className="num">التكلفة</th></tr></thead>
          <tbody>
            {list.data?.data.map((t) => <tr key={t.id}><td dir="ltr" className="code">{t.number}</td><td dir="ltr">{t.date}</td><td>{t.fromWarehouseName}</td><td>{t.toWarehouseName}</td><td className="num" dir="ltr">{formatAmount(t.totalCost)}</td></tr>)}
            {list.data?.data.length === 0 && <tr><td colSpan={5} className="muted empty-row">لا توجد تحويلات</td></tr>}
          </tbody>
        </table>
      </div>
    </>
  );
}
