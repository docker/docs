import { useState, type FormEvent } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../auth';
import { ErrorBox } from '../ui';

export default function Login() {
  const { login } = useAuth();
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState(false);

  async function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    setBusy(true);
    setError(null);
    try {
      await login(String(f.get('email')), String(f.get('password')));
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
        <h1>تسجيل الدخول</h1>
        <ErrorBox error={error} />
        <label>البريد الإلكتروني<input name="email" type="email" required autoComplete="email" dir="ltr" /></label>
        <label>كلمة المرور<input name="password" type="password" required autoComplete="current-password" dir="ltr" /></label>
        <button className="btn btn-primary" disabled={busy}>{busy ? 'جارٍ الدخول…' : 'دخول'}</button>
        <p className="muted">ليس لديك حساب؟ <Link to="/register">أنشئ منشأة جديدة</Link></p>
      </form>
    </div>
  );
}
