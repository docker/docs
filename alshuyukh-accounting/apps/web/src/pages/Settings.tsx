import { useState, type FormEvent } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api';
import { useAuth } from '../auth';
import { ErrorBox, PageHeader, useLoad } from '../ui';
import TaxRates from './TaxRates';

interface TenantSettings { tenantName: string; defaultCurrency: string; timezone: string; locale: string; fiscalYearStartMonth: number; dateFormat: string }

const MONTHS = ['يناير', 'فبراير', 'مارس', 'أبريل', 'مايو', 'يونيو', 'يوليو', 'أغسطس', 'سبتمبر', 'أكتوبر', 'نوفمبر', 'ديسمبر'];

export default function Settings() {
  const { can, reload } = useAuth();
  const settings = useLoad(() => api<TenantSettings>('GET', '/api/settings/tenant'));
  const [error, setError] = useState<unknown>(null);
  const [saved, setSaved] = useState(false);
  const [pwError, setPwError] = useState<unknown>(null);
  const [pwSaved, setPwSaved] = useState(false);
  const editable = can('settings.manage');

  async function save(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = Object.fromEntries(new FormData(e.currentTarget)) as Record<string, string>;
    setError(null); setSaved(false);
    try {
      await api('PATCH', '/api/settings/tenant', {
        tenantName: f.tenantName, defaultCurrency: f.defaultCurrency, timezone: f.timezone,
        fiscalYearStartMonth: Number(f.fiscalYearStartMonth), dateFormat: f.dateFormat,
      });
      setSaved(true);
      await reload();
    } catch (err) { setError(err); }
  }

  async function changePassword(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const f = Object.fromEntries(new FormData(form)) as Record<string, string>;
    setPwError(null); setPwSaved(false);
    try {
      await api('POST', '/api/auth/change-password', { currentPassword: f.currentPassword, newPassword: f.newPassword });
      setPwSaved(true);
      form.reset();
      await reload();
    } catch (err) { setPwError(err); }
  }

  const s = settings.data;
  return (
    <>
      <PageHeader title="الإعدادات" />
      <div className="settings-links">
        {can('company.view') && <Link className="card link-card" to="/settings/companies">الشركات والفروع والمستودعات</Link>}
        {can('user.view') && <Link className="card link-card" to="/settings/users">المستخدمون</Link>}
        {can('role.view') && <Link className="card link-card" to="/settings/roles">الأدوار والصلاحيات</Link>}
        {can('audit.view') && <Link className="card link-card" to="/settings/audit">سجل التدقيق</Link>}
        <Link className="card link-card" to="/settings/subscription">الاشتراك والباقة</Link>
      </div>

      {s && (
        <form className="card form" onSubmit={save}>
          <h2>إعدادات المنشأة</h2>
          <ErrorBox error={error} />
          {saved && <div className="alert alert-ok">تم الحفظ</div>}
          <fieldset disabled={!editable}>
            <label>اسم المنشأة<input name="tenantName" defaultValue={s.tenantName} required /></label>
            <div className="grid-2">
              <label>العملة الافتراضية<input name="defaultCurrency" defaultValue={s.defaultCurrency} pattern="[A-Z]{3}" dir="ltr" /></label>
              <label>المنطقة الزمنية<input name="timezone" defaultValue={s.timezone} dir="ltr" /></label>
              <label>بداية السنة المالية
                <select name="fiscalYearStartMonth" defaultValue={s.fiscalYearStartMonth}>
                  {MONTHS.map((m, i) => <option key={m} value={i + 1}>{m}</option>)}
                </select>
              </label>
              <label>تنسيق التاريخ
                <select name="dateFormat" defaultValue={s.dateFormat}>
                  <option>YYYY-MM-DD</option><option>DD/MM/YYYY</option><option>MM/DD/YYYY</option>
                </select>
              </label>
            </div>
            {editable && <button className="btn btn-primary">حفظ</button>}
          </fieldset>
        </form>
      )}

      {can('invoice.view') && <TaxRates />}

      <form className="card form" onSubmit={changePassword}>
        <h2>تغيير كلمة المرور</h2>
        <ErrorBox error={pwError} />
        {pwSaved && <div className="alert alert-ok">تم تغيير كلمة المرور وتسجيل الخروج من الأجهزة الأخرى</div>}
        <div className="grid-2">
          <label>كلمة المرور الحالية<input name="currentPassword" type="password" required dir="ltr" autoComplete="current-password" /></label>
          <label>كلمة المرور الجديدة<input name="newPassword" type="password" required minLength={10} dir="ltr" autoComplete="new-password" /></label>
        </div>
        <button className="btn btn-primary">تغيير</button>
      </form>
    </>
  );
}
