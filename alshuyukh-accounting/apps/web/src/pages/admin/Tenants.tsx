import { useState, type FormEvent } from 'react';
import { Link, useParams } from 'react-router-dom';
import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, formatDateTime, useLoad } from '../../ui';
import { EVENT_AR, LIMIT_AR, STATE_AR, UsageRow } from '../subscription/Subscription';

interface TenantRow { id: string; name: string; slug: string; status: string; createdAt: string; ownerEmail: string | null; planName: string | null; subscriptionState: string; periodEnd: string | null; users: number; companies: number }
interface PlanOpt { id: string; nameAr: string; isActive: boolean; priceMonthly: string }
const TENANT_STATUS: Record<string, [string, string]> = { ACTIVE: ['نشطة', 'status-posted'], SUSPENDED: ['موقوفة', 'status-reversed'], CANCELLED: ['ملغاة', 'status-reversed'] };

export function Tenants() {
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');
  const [creating, setCreating] = useState(false);
  const list = useLoad(() => api<{ data: TenantRow[]; total: number }>('GET', `/api/admin/tenants?limit=100${search ? `&search=${encodeURIComponent(search)}` : ''}${status ? `&status=${status}` : ''}`), [search, status]);
  return (
    <div className="card">
      <div className="toolbar">
        <input type="search" className="grow" placeholder="بحث بالاسم أو بريد المالك" onKeyDown={(e) => { if (e.key === 'Enter') setSearch(e.currentTarget.value); }} />
        <select aria-label="الحالة" value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="">كل الحالات</option>{Object.entries(TENANT_STATUS).map(([k, [l]]) => <option key={k} value={k}>{l}</option>)}
        </select>
        <button className="btn btn-primary" onClick={() => setCreating(true)}>منشأة جديدة</button>
      </div>
      {creating && <NewTenant onDone={() => { setCreating(false); list.reload(); }} />}
      <ErrorBox error={list.error} />
      <div className="table-wrap">
        <table className="table">
          <thead><tr><th>المنشأة</th><th>المالك</th><th>الباقة</th><th>الاشتراك</th><th>نهاية الفترة</th><th className="num">مستخدمون</th><th>الحالة</th></tr></thead>
          <tbody>
            {list.data?.data.map((t) => {
              const [sl, sc] = STATE_AR[t.subscriptionState] ?? [t.subscriptionState, ''];
              const [tl, tc] = TENANT_STATUS[t.status] ?? [t.status, ''];
              return (
                <tr key={t.id}>
                  <td><Link to={`/admin/tenants/${t.id}`}>{t.name}</Link><div className="muted small" dir="ltr">{t.slug}</div></td>
                  <td dir="ltr" className="small">{t.ownerEmail}</td><td>{t.planName ?? '—'}</td>
                  <td><span className={`tag ${sc}`}>{sl}</span></td><td>{formatDateTime(t.periodEnd)}</td>
                  <td className="num">{t.users}</td><td><span className={`tag ${tc}`}>{tl}</span></td>
                </tr>
              );
            })}
            {list.data?.data.length === 0 && <tr><td colSpan={7} className="muted empty-row">لا توجد منشآت</td></tr>}
          </tbody>
        </table>
      </div>
      {list.data && <p className="muted small">{list.data.total} منشأة</p>}
    </div>
  );
}

function NewTenant({ onDone }: { onDone: () => void }) {
  const plans = useLoad(() => api<{ data: PlanOpt[] }>('GET', '/api/admin/plans'));
  const [error, setError] = useState<unknown>(null);
  const [created, setCreated] = useState<{ ownerEmail: string; temporaryPassword: string; tenantId: string } | null>(null);
  async function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = Object.fromEntries(new FormData(e.currentTarget)) as Record<string, string>;
    setError(null);
    try { setCreated(await api('POST', '/api/admin/tenants', { ...f, planId: f.planId || null })); } catch (err) { setError(err); }
  }
  if (created) {
    return (
      <div className="alert alert-ok">
        أُنشئت المنشأة. كلمة المرور المؤقتة للمالك <strong dir="ltr">{created.ownerEmail}</strong>: <code dir="ltr">{created.temporaryPassword}</code>
        <div className="small">تظهر مرة واحدة فقط، وسيُطلب منه تغييرها عند أول دخول.</div>
        <div className="form-actions"><Link className="btn btn-small" to={`/admin/tenants/${created.tenantId}`}>فتح المنشأة</Link><button className="btn btn-small" onClick={onDone}>إغلاق</button></div>
      </div>
    );
  }
  return (
    <form className="form sub-form" onSubmit={submit}>
      <ErrorBox error={error} />
      <div className="grid-3">
        <label>اسم المنشأة<input name="tenantName" required minLength={2} /></label>
        <label>اسم الشركة<input name="companyName" required minLength={2} /></label>
        <label>الباقة<select name="planId"><option value="">الافتراضية</option>{plans.data?.data.filter((p) => p.isActive).map((p) => <option key={p.id} value={p.id}>{p.nameAr}</option>)}</select></label>
        <label>اسم المالك<input name="ownerName" required minLength={2} /></label>
        <label>بريد المالك<input name="ownerEmail" type="email" required dir="ltr" /></label>
      </div>
      <div className="form-actions"><button className="btn btn-primary">إنشاء</button><button type="button" className="btn btn-ghost" onClick={onDone}>إلغاء</button></div>
    </form>
  );
}

interface Detail {
  id: string; name: string; slug: string; status: string; createdAt: string;
  subscription: { planId: string | null; planName: string | null; state: string; billingCycle: string | null; periodEnd: string | null; limits: Record<string, number | null>; limitOverrides: Record<string, number | null>; items: { description: string; unitPrice: string; currency: string }[] };
  usage: Record<string, number>;
  billingEvents: { id: string; eventType: string; amount: string | null; currency: string | null; reference: string | null; details: Record<string, unknown>; createdAt: string; byEmail: string | null }[];
  members: { id: string; email: string; fullName: string; userStatus: string; memberStatus: string; isOwner: boolean; lastLoginAt: string | null }[];
  companies: { id: string; name: string; vatNumber: string | null }[];
  features: Record<string, boolean>; featureOverrides: { key: string; enabled: boolean }[];
}

export function TenantDetail() {
  const { id } = useParams();
  const d = useLoad(() => api<Detail>('GET', `/api/admin/tenants/${id}`), [id]);
  const plans = useLoad(() => api<{ data: PlanOpt[] }>('GET', '/api/admin/plans'));
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState(false);
  async function act(path: string, body: unknown) {
    setError(null); setBusy(true);
    try { await api(path.startsWith('PUT ') ? 'PUT' : 'POST', `/api/admin/tenants/${id}/${path.replace('PUT ', '')}`, body); d.reload(); } catch (e) { setError(e); } finally { setBusy(false); }
  }
  const form = (fn: (f: Record<string, string>) => void) => (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    fn(Object.fromEntries(new FormData(e.currentTarget)) as Record<string, string>);
    e.currentTarget.reset();
  };
  const t = d.data;
  if (!t) return <ErrorBox error={d.error} />;
  const [sl, sc] = STATE_AR[t.subscription.state] ?? [t.subscription.state, ''];
  return (
    <>
      <ErrorBox error={error} />
      <div className="card">
        <div className="entry-head">
          <div><h2>{t.name}</h2><p className="muted small" dir="ltr">{t.slug} · {t.id}</p></div>
          <span className={`tag ${TENANT_STATUS[t.status]?.[1] ?? ''}`}>{TENANT_STATUS[t.status]?.[0] ?? t.status}</span>
        </div>
        <form className="inline-form" onSubmit={form((f) => act(t.status === 'ACTIVE' ? 'suspend' : 'activate', { reason: f.reason }))}>
          <input name="reason" required minLength={3} placeholder="السبب" aria-label="سبب الإيقاف أو التفعيل" />
          <button className={`btn ${t.status === 'ACTIVE' ? 'btn-danger' : 'btn-primary'}`} disabled={busy}>{t.status === 'ACTIVE' ? 'إيقاف المنشأة' : 'تفعيل المنشأة'}</button>
        </form>
      </div>
      <div className="grid-cards">
        <div className="card">
          <div className="entry-head"><h2>الاشتراك: {t.subscription.planName ?? '—'}</h2><span className={`tag ${sc}`}>{sl}</span></div>
          <dl className="meta">
            <div><dt>نهاية الفترة</dt><dd>{formatDateTime(t.subscription.periodEnd)}</dd></div>
            <div><dt>الدورة</dt><dd>{t.subscription.billingCycle === 'YEARLY' ? 'سنوية' : 'شهرية'}</dd></div>
            {t.subscription.items.map((i) => <div key={i.description}><dt>{i.description}</dt><dd dir="ltr">{formatAmount(i.unitPrice)} {i.currency}</dd></div>)}
          </dl>
          <form key={`${t.subscription.planId}-${t.subscription.billingCycle}`} className="inline-form" onSubmit={form((f) => act('subscription/plan', { planId: f.planId, billingCycle: f.billingCycle }))}>
            <select name="planId" required aria-label="الباقة" defaultValue={t.subscription.planId ?? ''}>
              <option value="">— الباقة —</option>{plans.data?.data.filter((p) => p.isActive).map((p) => <option key={p.id} value={p.id}>{p.nameAr}</option>)}
            </select>
            <select name="billingCycle" aria-label="الدورة" defaultValue={t.subscription.billingCycle ?? 'MONTHLY'}><option value="MONTHLY">شهري</option><option value="YEARLY">سنوي</option></select>
            <button className="btn" disabled={busy}>تغيير الباقة</button>
          </form>
          <form className="inline-form" onSubmit={form((f) => act('subscription/payment', { amount: f.amount, reference: f.reference, cycles: Number(f.cycles || 1) }))}>
            <input name="amount" required dir="ltr" inputMode="decimal" pattern="[0-9]{1,14}(\.[0-9]{1,2})?" placeholder="المبلغ" aria-label="المبلغ" />
            <input name="reference" required minLength={2} placeholder="مرجع التحويل" aria-label="المرجع" />
            <input name="cycles" type="number" min={1} max={36} defaultValue={1} aria-label="عدد الدورات" />
            <button className="btn btn-primary" disabled={busy}>تسجيل دفعة وتمديد</button>
          </form>
          <form className="inline-form" onSubmit={form((f) => act('subscription/extend', { days: Number(f.days), reason: f.reason }))}>
            <input name="days" type="number" required min={1} max={365} placeholder="أيام" aria-label="أيام التمديد" />
            <input name="reason" required minLength={3} placeholder="سبب التمديد" aria-label="سبب التمديد" />
            <button className="btn" disabled={busy}>تمديد دون دفع</button>
          </form>
        </div>
        <div className="card">
          <h2>الاستخدام والحدود</h2>
          {Object.keys(LIMIT_AR).filter((k) => k !== 'max_storage_mb').map((k) => <UsageRow key={k} label={LIMIT_AR[k]!} used={t.usage[k] ?? 0} max={t.subscription.limits[k] ?? null} />)}
          <LimitsForm overrides={t.subscription.limitOverrides} onSave={(o) => act('subscription/limits', { overrides: o })} busy={busy} />
        </div>
      </div>
      <div className="card">
        <h2>خصائص التشغيل</h2>
        <div className="table-wrap"><table className="table">
          <thead><tr><th>الخاصية</th><th>الفعّال</th><th>تخصيص المنشأة</th></tr></thead>
          <tbody>{Object.entries(t.features).map(([k, v]) => {
            const o = t.featureOverrides.find((x) => x.key === k);
            return (
              <tr key={k}><td dir="ltr" className="code">{k}</td><td>{v ? 'مفعّلة' : 'معطّلة'}</td>
                <td><select aria-label={`تخصيص ${k}`} value={o ? String(o.enabled) : ''} disabled={busy} onChange={(e) => act(`PUT features/${k}`, { enabled: e.target.value === '' ? null : e.target.value === 'true' })}>
                  <option value="">الافتراضي العام</option><option value="true">تفعيل لهذه المنشأة</option><option value="false">تعطيل لهذه المنشأة</option>
                </select></td></tr>
            );
          })}</tbody>
        </table></div>
      </div>
      <div className="card">
        <h2>المستخدمون</h2>
        <table className="table"><thead><tr><th>البريد</th><th>الاسم</th><th>الحالة</th><th>آخر دخول</th></tr></thead>
          <tbody>{t.members.map((m) => <tr key={m.id}><td dir="ltr">{m.email}{m.isOwner && <span className="tag"> مالك</span>}</td><td>{m.fullName}</td><td>{m.userStatus === 'ACTIVE' && m.memberStatus === 'ACTIVE' ? 'نشط' : 'معطل'}</td><td>{formatDateTime(m.lastLoginAt)}</td></tr>)}</tbody>
        </table>
      </div>
      <div className="card">
        <h2>سجل الفوترة</h2>
        <table className="table"><thead><tr><th>التاريخ</th><th>الحدث</th><th className="num">المبلغ</th><th>المرجع</th><th>بواسطة</th></tr></thead>
          <tbody>{t.billingEvents.map((e) => <tr key={e.id}><td>{formatDateTime(e.createdAt)}</td><td>{EVENT_AR[e.eventType] ?? e.eventType}{typeof e.details?.reason === 'string' ? ` — ${e.details.reason}` : ''}</td>
            <td className="num" dir="ltr">{e.amount ? `${formatAmount(e.amount)} ${e.currency}` : ''}</td><td dir="ltr">{e.reference ?? ''}</td><td dir="ltr" className="small">{e.byEmail ?? ''}</td></tr>)}</tbody>
        </table>
      </div>
    </>
  );
}

function LimitsForm({ overrides, onSave, busy }: { overrides: Record<string, number | null>; onSave: (o: Record<string, number | null>) => void; busy: boolean }) {
  function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const out: Record<string, number | null> = {};
    for (const k of Object.keys(LIMIT_AR)) {
      const v = String(f.get(k) ?? '').trim();
      if (v === '') continue;            // follow the plan
      out[k] = v === '∞' ? null : Number(v); // explicit unlimited or a number
    }
    onSave(out);
  }
  return (
    <details className="limits-form">
      <summary className="small">تخصيص الحدود لهذه المنشأة</summary>
      <form className="form" onSubmit={submit}>
        <p className="muted small">اترك الحقل فارغًا لاتباع الباقة، أو اكتب ∞ لغير محدود.</p>
        <div className="grid-2">
          {Object.entries(LIMIT_AR).map(([k, l]) => (
            <label key={k}>{l}<input name={k} dir="ltr" pattern="[0-9]{1,9}|∞" defaultValue={k in overrides ? (overrides[k] === null ? '∞' : String(overrides[k])) : ''} /></label>
          ))}
        </div>
        <button className="btn" disabled={busy}>حفظ الحدود</button>
      </form>
    </details>
  );
}
