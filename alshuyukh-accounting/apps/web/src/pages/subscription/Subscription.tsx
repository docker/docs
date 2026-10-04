import { useState } from 'react';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { formatAmount } from '../../money';
import { ErrorBox, formatDateTime, PageHeader, useLoad } from '../../ui';

export const LIMIT_AR: Record<string, string> = {
  max_users: 'المستخدمون', max_companies: 'الشركات', max_branches: 'الفروع', max_warehouses: 'المستودعات', max_products: 'المنتجات والخدمات',
  max_invoices_per_month: 'فواتير المبيعات هذا الشهر', max_storage_mb: 'التخزين (MB)', max_api_calls_per_month: 'طلبات API هذا الشهر',
};
export const STATE_AR: Record<string, [string, string]> = {
  TRIALING: ['فترة تجريبية', 'status-partial'], ACTIVE: ['نشط', 'status-posted'], GRACE: ['فترة سماح', 'status-partial'],
  EXPIRED: ['منتهٍ — للاطلاع فقط', 'status-reversed'], CANCELLED: ['ملغى', 'status-reversed'], NONE: ['بلا اشتراك', 'status-reversed'],
};
const EVENT_AR: Record<string, string> = {
  TRIAL_STARTED: 'بدء الفترة التجريبية', PLAN_CHANGED: 'تغيير الباقة', PERIOD_EXTENDED: 'تمديد الفترة', PAYMENT_RECORDED: 'تسجيل دفعة',
  LIMITS_CHANGED: 'تعديل الحدود', PLAN_CHANGE_REQUESTED: 'طلب تغيير الباقة', SUBSCRIPTION_CANCELLED: 'إلغاء الاشتراك',
  TENANT_SUSPENDED: 'إيقاف المنشأة', TENANT_ACTIVATED: 'تفعيل المنشأة',
};
export { EVENT_AR };

interface Sub {
  planName: string | null; state: string; billingCycle: string | null; periodEnd: string | null; graceEnd: string | null;
  limits: Record<string, number | null>; usage: Record<string, number>; items: { description: string; unitPrice: string; currency: string; quantity: number }[];
}
interface Plan { id: string; nameAr: string; description: string | null; priceMonthly: string; priceYearly: string; currency: string; [k: string]: unknown }

/** Usage against a limit; a bar only when the limit is finite. */
export function UsageRow({ label, used, max }: { label: string; used: number; max: number | null }) {
  const pct = max ? Math.min(100, Math.round((used / max) * 100)) : 0;
  return (
    <div className="usage-row">
      <div className="usage-label"><span>{label}</span><span dir="ltr">{used}{max === null ? ' / ∞' : ` / ${max}`}</span></div>
      {max !== null && <div className="usage-bar" role="progressbar" aria-valuenow={pct} aria-valuemin={0} aria-valuemax={100} aria-label={label}><i style={{ width: `${pct}%` }} className={pct >= 90 ? 'high' : ''} /></div>}
    </div>
  );
}

export default function Subscription() {
  const { can } = useAuth();
  const sub = useLoad(() => api<Sub>('GET', '/api/subscription'));
  const plans = useLoad(() => api<{ data: Plan[] }>('GET', '/api/subscription/plans'));
  const events = useLoad(() => (can('subscription.manage') ? api<{ data: { id: string; eventType: string; amount: string | null; currency: string | null; reference: string | null; createdAt: string }[] }>('GET', '/api/subscription/billing-events') : Promise.resolve({ data: [] })));
  const [cycle, setCycle] = useState<'MONTHLY' | 'YEARLY'>('MONTHLY');
  const [sent, setSent] = useState<string | null>(null);
  const [error, setError] = useState<unknown>(null);
  async function request(planId: string) {
    setError(null);
    try { await api('POST', '/api/subscription/request-change', { planId, billingCycle: cycle }); setSent(planId); events.reload(); } catch (e) { setError(e); }
  }
  const s = sub.data;
  const [label, cls] = s ? STATE_AR[s.state] ?? [s.state, ''] : ['', ''];
  return (
    <>
      <PageHeader title="الاشتراك والباقة" />
      <ErrorBox error={sub.error ?? error} />
      {s && (
        <div className="grid-cards">
          <div className="card">
            <div className="entry-head"><h2>{s.planName ?? 'بلا باقة'}</h2><span className={`tag ${cls}`}>{label}</span></div>
            <dl className="meta">
              <div><dt>{s.state === 'TRIALING' ? 'تنتهي التجربة' : 'نهاية الفترة'}</dt><dd>{formatDateTime(s.periodEnd)}</dd></div>
              {s.billingCycle && <div><dt>دورة الفوترة</dt><dd>{s.billingCycle === 'YEARLY' ? 'سنوية' : 'شهرية'}</dd></div>}
              {s.items.map((i) => <div key={i.description}><dt>{i.description}</dt><dd dir="ltr">{formatAmount(i.unitPrice)} {i.currency}</dd></div>)}
            </dl>
          </div>
          <div className="card">
            <h2>الاستخدام والحدود</h2>
            {Object.keys(LIMIT_AR).filter((k) => k !== 'max_storage_mb').map((k) => <UsageRow key={k} label={LIMIT_AR[k]!} used={s.usage[k] ?? 0} max={s.limits[k] ?? null} />)}
          </div>
        </div>
      )}
      {can('subscription.manage') && plans.data && (
        <div className="card">
          <div className="card-head">
            <h2>الباقات المتاحة</h2>
            <select aria-label="دورة الفوترة" value={cycle} onChange={(e) => setCycle(e.target.value as 'MONTHLY' | 'YEARLY')}>
              <option value="MONTHLY">شهري</option><option value="YEARLY">سنوي</option>
            </select>
          </div>
          <p className="muted small">الدفع الإلكتروني غير مفعّل بعد: اطلب الباقة وسيتواصل معك فريق المنصة لإتمام الدفع وتفعيلها.</p>
          <div className="plan-grid">
            {plans.data.data.map((p) => (
              <div key={p.id} className="plan-card">
                <h3>{p.nameAr}</h3>
                <div className="plan-price" dir="ltr">{formatAmount(cycle === 'YEARLY' ? p.priceYearly : p.priceMonthly)} <small>{p.currency} / {cycle === 'YEARLY' ? 'سنة' : 'شهر'}</small></div>
                {p.description && <p className="muted small">{p.description}</p>}
                <ul className="plan-limits">
                  {Object.keys(LIMIT_AR).map((k) => <li key={k}>{LIMIT_AR[k]}: <strong>{p[k] === null ? 'غير محدود' : String(p[k])}</strong></li>)}
                </ul>
                {sent === p.id ? <div className="alert alert-ok small">أُرسل الطلب</div> : <button className="btn btn-primary" onClick={() => request(p.id)}>طلب هذه الباقة</button>}
              </div>
            ))}
            {plans.data.data.length === 0 && <p className="muted">لا توجد باقات معروضة حاليًا.</p>}
          </div>
        </div>
      )}
      {(events.data?.data.length ?? 0) > 0 && (
        <div className="card">
          <h2>سجل الفوترة</h2>
          <table className="table">
            <thead><tr><th>التاريخ</th><th>الحدث</th><th className="num">المبلغ</th><th>المرجع</th></tr></thead>
            <tbody>{events.data!.data.map((e) => (
              <tr key={e.id}><td>{formatDateTime(e.createdAt)}</td><td>{EVENT_AR[e.eventType] ?? e.eventType}</td>
                <td className="num" dir="ltr">{e.amount ? `${formatAmount(e.amount)} ${e.currency}` : ''}</td><td dir="ltr">{e.reference ?? ''}</td></tr>
            ))}</tbody>
          </table>
        </div>
      )}
    </>
  );
}
