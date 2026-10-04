import { useState, type FormEvent } from 'react';
import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { LIMIT_AR } from '../subscription/Subscription';

interface Plan { id: string; code: string; nameAr: string; description: string | null; currency: string; priceMonthly: string; priceYearly: string; trialDays: number; graceDays: number; isPublic: boolean; isActive: boolean; isDefault: boolean; sortOrder: number; subscribers: number; [k: string]: unknown }

/** Plans, prices and limits are data: edited here, never in code. */
export function Plans() {
  const plans = useLoad(() => api<{ data: Plan[] }>('GET', '/api/admin/plans'));
  const [editing, setEditing] = useState<Plan | 'new' | null>(null);
  const [error, setError] = useState<unknown>(null);
  async function save(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const num = (k: string) => { const v = String(f.get(k) ?? '').trim(); return v === '' ? null : Number(v); };
    const body = {
      code: String(f.get('code')), nameAr: String(f.get('nameAr')), description: String(f.get('description') || '') || null,
      priceMonthly: String(f.get('priceMonthly')), priceYearly: String(f.get('priceYearly')),
      trialDays: Number(f.get('trialDays')), graceDays: Number(f.get('graceDays')), sortOrder: Number(f.get('sortOrder')),
      isPublic: f.get('isPublic') === 'on', isActive: f.get('isActive') === 'on', isDefault: f.get('isDefault') === 'on',
      ...Object.fromEntries(Object.keys(LIMIT_AR).map((k) => [k, num(k)])),
    };
    setError(null);
    try {
      if (editing === 'new') await api('POST', '/api/admin/plans', body);
      else if (editing) await api('PATCH', `/api/admin/plans/${editing.id}`, body);
      setEditing(null); plans.reload();
    } catch (err) { setError(err); }
  }
  const p = editing && editing !== 'new' ? editing : null;
  return (
    <>
      <ErrorBox error={plans.error} />
      {editing ? (
        <form className="card form" onSubmit={save} key={p?.id ?? 'new'}>
          <h2>{p ? `تعديل ${p.nameAr}` : 'باقة جديدة'}</h2>
          <ErrorBox error={error} />
          <div className="grid-3">
            <label>الرمز<input name="code" required dir="ltr" pattern="[A-Za-z][A-Za-z0-9_]{1,30}" defaultValue={p?.code ?? ''} /></label>
            <label>الاسم<input name="nameAr" required minLength={2} defaultValue={p?.nameAr ?? ''} /></label>
            <label>الترتيب<input name="sortOrder" type="number" min={0} defaultValue={p?.sortOrder ?? 10} /></label>
            <label>السعر الشهري<input name="priceMonthly" required dir="ltr" pattern="[0-9]{1,14}(\.[0-9]{1,2})?" defaultValue={p?.priceMonthly ?? ''} /></label>
            <label>السعر السنوي<input name="priceYearly" required dir="ltr" pattern="[0-9]{1,14}(\.[0-9]{1,2})?" defaultValue={p?.priceYearly ?? ''} /></label>
            <label>أيام التجربة<input name="trialDays" type="number" min={0} max={365} defaultValue={p?.trialDays ?? 0} /></label>
            <label>أيام السماح بعد الانتهاء<input name="graceDays" type="number" min={0} max={90} defaultValue={p?.graceDays ?? 7} /></label>
          </div>
          <label>الوصف<input name="description" defaultValue={p?.description ?? ''} /></label>
          <h3>الحدود <span className="muted small">(فارغ = غير محدود)</span></h3>
          <div className="grid-3">{Object.entries(LIMIT_AR).map(([k, l]) => <label key={k}>{l}<input name={k} type="number" min={1} dir="ltr" defaultValue={p?.[k] == null ? '' : String(p[k])} /></label>)}</div>
          <label className="check"><input type="checkbox" name="isPublic" defaultChecked={p?.isPublic ?? true} />معروضة للمنشآت</label>
          <label className="check"><input type="checkbox" name="isActive" defaultChecked={p?.isActive ?? true} />متاحة</label>
          <label className="check"><input type="checkbox" name="isDefault" defaultChecked={p?.isDefault ?? false} />الباقة الافتراضية للتسجيل الجديد</label>
          <p className="muted small">تغيير السعر يطبق على الاشتراكات الجديدة وتغييرات الباقة؛ الاشتراكات الحالية تحتفظ بسعرها.</p>
          <div className="form-actions"><button className="btn btn-primary">حفظ</button><button type="button" className="btn btn-ghost" onClick={() => setEditing(null)}>إلغاء</button></div>
        </form>
      ) : <div className="toolbar"><button className="btn btn-primary" onClick={() => setEditing('new')}>باقة جديدة</button></div>}
      <div className="plan-grid">
        {plans.data?.data.map((x) => (
          <div key={x.id} className={`plan-card ${x.isActive ? '' : 'muted'}`}>
            <div className="card-head"><h3>{x.nameAr}</h3><span dir="ltr" className="code small">{x.code}</span></div>
            <div className="plan-price" dir="ltr">{formatAmount(x.priceMonthly)} <small>{x.currency} / شهر</small></div>
            <div className="chips">
              {x.isDefault && <span className="tag status-posted">افتراضية</span>}
              {!x.isPublic && <span className="tag">مخفية</span>}
              {!x.isActive && <span className="tag status-reversed">موقوفة</span>}
              {x.trialDays > 0 && <span className="tag">تجربة {x.trialDays} يومًا</span>}
              <span className="tag">{x.subscribers} مشترك</span>
            </div>
            <ul className="plan-limits">{Object.entries(LIMIT_AR).map(([k, l]) => <li key={k}>{l}: <strong>{x[k] == null ? 'غير محدود' : String(x[k])}</strong></li>)}</ul>
            <button className="btn btn-small" onClick={() => { setError(null); setEditing(x); }}>تعديل</button>
          </div>
        ))}
      </div>
    </>
  );
}
