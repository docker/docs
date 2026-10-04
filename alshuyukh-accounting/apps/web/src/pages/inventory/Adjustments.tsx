import { useState, type FormEvent } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import type { Account } from '../accounting/ChartOfAccounts';
import { loadStockedProducts, loadWarehouses, qty } from './shared';

interface Adjustment { id: string; number: string; date: string; warehouseName: string; reason: string; offsetAccountName: string; totalIncrease: string; totalDecrease: string; journalEntryId: string | null }
interface Line { productId: string; counted: string; direction: 'IN' | 'OUT'; quantity: string; unitCost: string }

const today = () => new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh' }).format(new Date());
const empty = (): Line => ({ productId: '', counted: '', direction: 'IN', quantity: '', unitCost: '' });

export default function Adjustments() {
  const { can } = useAuth();
  const list = useLoad(() => api<{ data: Adjustment[] }>('GET', '/api/stock-adjustments'));
  const warehouses = useLoad(loadWarehouses);
  const products = useLoad(loadStockedProducts);
  const accounts = useLoad(() => api<{ data: Account[] }>('GET', '/api/accounts?postableOnly=true'));
  const [mode, setMode] = useState<'count' | 'manual' | null>(null);
  const [lines, setLines] = useState<Line[]>([empty()]);
  const [error, setError] = useState<unknown>(null);
  const update = (i: number, patch: Partial<Line>) => setLines((ls) => ls.map((l, j) => (j === i ? { ...l, ...patch } : l)));

  async function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    setError(null);
    try {
      await api('POST', '/api/stock-adjustments', {
        warehouseId: f.get('warehouseId'), date: f.get('date'), reason: f.get('reason'),
        ...(f.get('offsetAccountId') ? { offsetAccountId: f.get('offsetAccountId') } : {}),
        lines: lines.filter((l) => l.productId).map((l) => mode === 'count'
          ? { productId: l.productId, countedQuantity: l.counted.trim() }
          : { productId: l.productId, direction: l.direction, quantity: l.quantity.trim(), ...(l.unitCost.trim() ? { unitCost: l.unitCost.trim() } : {}) }),
      });
      setMode(null);
      setLines([empty()]);
      list.reload();
    } catch (err) { setError(err); }
  }

  const offsetAccounts = accounts.data?.data.filter((a) => ['EXPENSE', 'COST_OF_GOODS_SOLD', 'EQUITY', 'REVENUE'].includes(a.type)) ?? [];
  return (
    <>
      {mode && warehouses.data && products.data && accounts.data && (
        <form className="card form" onSubmit={submit}>
          <h2>{mode === 'count' ? 'جرد المخزون' : 'تسوية مخزون'}</h2>
          <ErrorBox error={error} />
          <div className="grid-3">
            <label>المستودع<select name="warehouseId" required>{warehouses.data?.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}</select></label>
            <label>التاريخ<input name="date" type="date" required defaultValue={today()} /></label>
            <label>السبب<input name="reason" required minLength={3} defaultValue={mode === 'count' ? 'جرد دوري' : ''} /></label>
            <label>الحساب المقابل
              <select name="offsetAccountId" defaultValue="">
                <option value="">فروقات وتسويات المخزون (افتراضي)</option>
                {offsetAccounts.map((a) => <option key={a.id} value={a.id}>{a.code} — {a.nameAr}</option>)}
              </select>
            </label>
          </div>
          <table className="table">
            <thead><tr><th>الصنف</th>{mode === 'count' ? <><th className="num">الرصيد الحالي</th><th>الكمية المعدودة</th></> : <><th>النوع</th><th>الكمية</th><th>تكلفة الوحدة (للزيادة)</th></>}<th /></tr></thead>
            <tbody>{lines.map((l, i) => (
              <tr key={i}>
                <td><select value={l.productId} aria-label={`صنف ${i + 1}`} onChange={(e) => update(i, { productId: e.target.value })}>
                  <option value="">— اختر —</option>
                  {products.data?.map((p) => <option key={p.id} value={p.id}>{p.sku} — {p.nameAr}</option>)}
                </select></td>
                {mode === 'count' ? (
                  <>
                    <td className="num" dir="ltr">{qty(products.data?.find((p) => p.id === l.productId)?.onHand)}</td>
                    <td><input dir="ltr" inputMode="decimal" value={l.counted} aria-label={`المعدود ${i + 1}`} onChange={(e) => update(i, { counted: e.target.value })} /></td>
                  </>
                ) : (
                  <>
                    <td><select value={l.direction} onChange={(e) => update(i, { direction: e.target.value as 'IN' | 'OUT' })}><option value="IN">زيادة</option><option value="OUT">نقص</option></select></td>
                    <td><input dir="ltr" inputMode="decimal" value={l.quantity} aria-label={`كمية ${i + 1}`} onChange={(e) => update(i, { quantity: e.target.value })} /></td>
                    <td><input dir="ltr" inputMode="decimal" value={l.unitCost} disabled={l.direction === 'OUT'} placeholder="متوسط التكلفة" aria-label={`تكلفة ${i + 1}`} onChange={(e) => update(i, { unitCost: e.target.value })} /></td>
                  </>
                )}
                <td>{lines.length > 1 && <button type="button" className="icon-btn" onClick={() => setLines((ls) => ls.filter((_, j) => j !== i))}>×</button>}</td>
              </tr>
            ))}</tbody>
          </table>
          <button type="button" className="btn btn-small" onClick={() => setLines((ls) => [...ls, empty()])}>+ سطر</button>
          <p className="muted small">
            {mode === 'count' ? 'الرصيد الحالي هو مجموع كل المستودعات للعرض فقط؛ الفرق يُحسب في الخادم لرصيد المستودع المختار. ' : ''}
            النقص يُقيَّم بمتوسط التكلفة. للرصيد الافتتاحي اختر حساب رأس المال كحساب مقابل.
          </p>
          <div className="form-actions"><button className="btn btn-primary">ترحيل</button><button type="button" className="btn btn-ghost" onClick={() => setMode(null)}>إلغاء</button></div>
        </form>
      )}
      <div className="card">
        <div className="toolbar">
          <span className="spacer" />
          {can('inventory.adjust') && !mode && (
            <>
              <button className="btn" onClick={() => { setLines([empty()]); setMode('manual'); }}>تسوية جديدة</button>
              <button className="btn btn-primary" onClick={() => { setLines([empty()]); setMode('count'); }}>جرد جديد</button>
            </>
          )}
        </div>
        <ErrorBox error={list.error} />
        <table className="table">
          <thead><tr><th>الرقم</th><th>التاريخ</th><th>المستودع</th><th>السبب</th><th className="num">زيادة</th><th className="num">نقص</th><th /></tr></thead>
          <tbody>
            {list.data?.data.map((a) => (
              <tr key={a.id}>
                <td dir="ltr" className="code">{a.number}</td><td dir="ltr">{a.date}</td><td>{a.warehouseName}</td><td>{a.reason}</td>
                <td className="num" dir="ltr">{a.totalIncrease === '0.00' ? '' : formatAmount(a.totalIncrease)}</td>
                <td className="num" dir="ltr">{a.totalDecrease === '0.00' ? '' : formatAmount(a.totalDecrease)}</td>
                <td>{a.journalEntryId && <Link to={`/accounting/journal/${a.journalEntryId}`}>القيد</Link>}</td>
              </tr>
            ))}
            {list.data?.data.length === 0 && <tr><td colSpan={7} className="muted empty-row">لا توجد تسويات</td></tr>}
          </tbody>
        </table>
      </div>
    </>
  );
}
