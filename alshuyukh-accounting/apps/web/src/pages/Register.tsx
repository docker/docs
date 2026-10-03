import { useState, type FormEvent } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../auth';
import { ErrorBox } from '../ui';

export default function Register() {
  const { register } = useAuth();
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState(false);

  async function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = Object.fromEntries(new FormData(e.currentTarget)) as Record<string, string>;
    setBusy(true);
    setError(null);
    try {
      await register({
        fullName: f.fullName, email: f.email, password: f.password,
        tenantName: f.companyName, companyName: f.companyName,
        vatNumber: f.vatNumber || null, commercialRegistration: f.commercialRegistration || null,
      });
    } catch (err) {
      setError(err);
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="auth-page">
      <form className="card auth-card" onSubmit={submit}>
        <div className="brand brand-lg"><span className="brand-mark">ش</span> الشيوخ للمحاسبة</div>
        <h1>إنشاء منشأة جديدة</h1>
        <ErrorBox error={error} />
        <label>الاسم الكامل<input name="fullName" required minLength={2} /></label>
        <label>البريد الإلكتروني<input name="email" type="email" required dir="ltr" /></label>
        <label>كلمة المرور<input name="password" type="password" required minLength={10} dir="ltr" />
          <span className="hint">10 أحرف على الأقل، وتحتوي على حروف وأرقام</span></label>
        <label>اسم الشركة<input name="companyName" required minLength={2} /></label>
        <div className="grid-2">
          <label>الرقم الضريبي (اختياري)<input name="vatNumber" pattern="3[0-9]{13}3" dir="ltr" placeholder="3XXXXXXXXXXXXX3" /></label>
          <label>السجل التجاري (اختياري)<input name="commercialRegistration" pattern="[0-9]{10}" dir="ltr" /></label>
        </div>
        <button className="btn btn-primary" disabled={busy}>{busy ? 'جارٍ الإنشاء…' : 'إنشاء الحساب'}</button>
        <p className="muted">لديك حساب؟ <Link to="/">تسجيل الدخول</Link></p>
      </form>
    </div>
  );
}
