import { useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { MonthlyBars } from '../../components/Charts';
import { formatAmount } from '../../money';
import { ErrorBox, formatDateTime, useLoad } from '../../ui';
import { STATE_AR } from '../subscription/Subscription';

export function Subscriptions() {
  const [state, setState] = useState('');
  const list = useLoad(() => api<{ data: { id: string; tenantId: string; tenantName: string; tenantStatus: string; planName: string; billingCycle: string; state: string; periodEnd: string; monthlyValue: string }[] }>(
    'GET', `/api/admin/subscriptions${state ? `?state=${state}` : ''}`), [state]);
  return (
    <div className="card">
      <div className="toolbar">
        <select aria-label="الحالة" value={state} onChange={(e) => setState(e.target.value)}>
          <option value="">كل الحالات</option>{['TRIALING', 'ACTIVE', 'GRACE', 'EXPIRED'].map((k) => <option key={k} value={k}>{STATE_AR[k]![0]}</option>)}
        </select>
      </div>
      <ErrorBox error={list.error} />
      <div className="table-wrap"><table className="table">
        <thead><tr><th>المنشأة</th><th>الباقة</th><th>الدورة</th><th>الحالة</th><th>نهاية الفترة</th><th className="num">القيمة الشهرية</th></tr></thead>
        <tbody>{list.data?.data.map((s) => (
          <tr key={s.id}><td><Link to={`/admin/tenants/${s.tenantId}`}>{s.tenantName}</Link></td><td>{s.planName}</td><td>{s.billingCycle === 'YEARLY' ? 'سنوية' : 'شهرية'}</td>
            <td><span className={`tag ${STATE_AR[s.state]?.[1] ?? ''}`}>{STATE_AR[s.state]?.[0] ?? s.state}</span></td><td>{formatDateTime(s.periodEnd)}</td>
            <td className="num" dir="ltr">{formatAmount(s.monthlyValue)}</td></tr>
        ))}{list.data?.data.length === 0 && <tr><td colSpan={6} className="muted empty-row">لا توجد اشتراكات</td></tr>}</tbody>
      </table></div>
    </div>
  );
}

export function Revenue() {
  const r = useLoad(() => api<{ monthly: { month: string; amount: string; payments: number }[]; byPlan: { planName: string; subscriptions: number; mrr: string }[]; payments: { id: string; tenantId: string; tenantName: string; amount: string; currency: string; reference: string | null; createdAt: string }[] }>('GET', '/api/admin/revenue'));
  const d = r.data;
  // Fill the last 12 months so the chart has a continuous axis.
  const months = Array.from({ length: 12 }, (_, i) => { const x = new Date(); x.setUTCDate(1); x.setUTCMonth(x.getUTCMonth() - 11 + i); return x.toISOString().slice(0, 7); });
  const series = months.map((m) => ({ month: m, amount: d?.monthly.find((x) => x.month === m)?.amount ?? '0' }));
  return (
    <>
      <ErrorBox error={r.error} />
      <p className="muted small">الإيراد من الدفعات المسجلة يدويًا (الدفع الإلكتروني غير مفعّل بعد).</p>
      {d && (
        <div className="grid-cards charts">
          <div className="card"><h2>المحصّل شهريًا</h2><MonthlyBars data={series} series={[{ key: 'amount', label: 'المحصّل', color: 'var(--series-1)' }]} /></div>
          <div className="card">
            <h2>الإيراد الشهري المتكرر حسب الباقة</h2>
            <table className="table"><thead><tr><th>الباقة</th><th className="num">اشتراكات</th><th className="num">MRR</th></tr></thead>
              <tbody>{d.byPlan.map((p) => <tr key={p.planName}><td>{p.planName}</td><td className="num">{p.subscriptions}</td><td className="num" dir="ltr">{formatAmount(p.mrr)}</td></tr>)}
                {d.byPlan.length === 0 && <tr><td colSpan={3} className="muted empty-row">لا توجد اشتراكات مدفوعة نشطة</td></tr>}</tbody>
            </table>
          </div>
        </div>
      )}
      {d && (
        <div className="card">
          <h2>آخر الدفعات</h2>
          <table className="table"><thead><tr><th>التاريخ</th><th>المنشأة</th><th className="num">المبلغ</th><th>المرجع</th></tr></thead>
            <tbody>{d.payments.map((p) => <tr key={p.id}><td>{formatDateTime(p.createdAt)}</td><td><Link to={`/admin/tenants/${p.tenantId}`}>{p.tenantName}</Link></td>
              <td className="num" dir="ltr">{formatAmount(p.amount)} {p.currency}</td><td dir="ltr">{p.reference}</td></tr>)}</tbody>
          </table>
        </div>
      )}
    </>
  );
}

export function Usage() {
  const u = useLoad(() => api<{ data: Record<string, any>[] }>('GET', '/api/admin/usage'));
  const cell = (used: number, max: number | null) => {
    const high = max !== null && used >= max * 0.9;
    return <td className={`num ${high ? 'negative' : ''}`} dir="ltr">{used}{max === null ? '' : ` / ${max}`}</td>;
  };
  return (
    <div className="card">
      <ErrorBox error={u.error} />
      <div className="table-wrap"><table className="table">
        <thead><tr><th>المنشأة</th><th>الباقة</th><th className="num">مستخدمون</th><th className="num">شركات</th><th className="num">منتجات</th><th className="num">فواتير الشهر</th><th className="num">طلبات API الشهر</th></tr></thead>
        <tbody>{u.data?.data.map((r) => (
          <tr key={r.tenantId}><td><Link to={`/admin/tenants/${r.tenantId}`}>{r.tenantName}</Link></td><td>{r.planName ?? '—'}</td>
            {cell(r.users, r.maxUsers)}{cell(r.companies, r.maxCompanies)}{cell(r.products, r.maxProducts)}{cell(r.invoices, r.maxInvoices)}{cell(r.apiCalls, r.maxApiCalls)}</tr>
        ))}</tbody>
      </table></div>
      <p className="muted small">يتجاوز 90% من الحد يظهر بالأحمر. طلبات API تُكتب كل 15 ثانية.</p>
    </div>
  );
}
