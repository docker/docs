import { useState, type FormEvent } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { ErrorBox, formatDateTime, useLoad } from '../../ui';

interface U { id: string; email: string; fullName: string; status: string; isPlatformAdmin: boolean; lockedUntil: string | null; failedLogins: number; lastLoginAt: string | null; tenants: { tenantId: string; name: string; isOwner: boolean }[] }

export function Users() {
  const { me } = useAuth();
  const [search, setSearch] = useState('');
  const list = useLoad(() => api<{ data: U[]; total: number }>('GET', `/api/admin/users?limit=100${search ? `&search=${encodeURIComponent(search)}` : ''}`), [search]);
  const [error, setError] = useState<unknown>(null);
  async function act(id: string, path: string, body: unknown = {}) {
    setError(null);
    try { await api('POST', `/api/admin/users/${id}/${path}`, body); list.reload(); } catch (e) { setError(e); }
  }
  return (
    <div className="card">
      <div className="toolbar"><input type="search" className="grow" placeholder="بحث بالبريد أو الاسم" onKeyDown={(e) => { if (e.key === 'Enter') setSearch(e.currentTarget.value); }} /></div>
      <ErrorBox error={error ?? list.error} />
      <div className="table-wrap"><table className="table">
        <thead><tr><th>المستخدم</th><th>المنشآت</th><th>آخر دخول</th><th>الحالة</th><th /></tr></thead>
        <tbody>{list.data?.data.map((u) => {
          const self = u.id === me?.user.id;
          const locked = u.lockedUntil && new Date(u.lockedUntil) > new Date();
          return (
            <tr key={u.id}>
              <td><div dir="ltr">{u.email}</div><div className="muted small">{u.fullName}{u.isPlatformAdmin && <span className="tag"> مدير منصة</span>}</div></td>
              <td className="small">{u.tenants.map((t) => <div key={t.tenantId}><Link to={`/admin/tenants/${t.tenantId}`}>{t.name}</Link>{t.isOwner ? ' (مالك)' : ''}</div>)}</td>
              <td>{formatDateTime(u.lastLoginAt)}</td>
              <td>{u.status === 'ACTIVE' ? 'نشط' : <span className="negative">معطل</span>}{locked && <div className="negative small">مقفل مؤقتًا</div>}</td>
              <td className="actions-cell">
                {!self && <button className="btn btn-small" onClick={() => act(u.id, 'status', { status: u.status === 'ACTIVE' ? 'DISABLED' : 'ACTIVE' })}>{u.status === 'ACTIVE' ? 'تعطيل' : 'تفعيل'}</button>}
                {locked && <button className="btn btn-small" onClick={() => act(u.id, 'unlock')}>فك القفل</button>}
                {!self && <button className="btn btn-small" onClick={() => act(u.id, 'platform-admin', { grant: !u.isPlatformAdmin })}>{u.isPlatformAdmin ? 'سحب صلاحية المنصة' : 'منح صلاحية المنصة'}</button>}
              </td>
            </tr>
          );
        })}</tbody>
      </table></div>
    </div>
  );
}

export function Flags() {
  const list = useLoad(() => api<{ data: { key: string; nameAr: string; description: string | null; enabled: boolean; overrides: number }[] }>('GET', '/api/admin/feature-flags'));
  const [error, setError] = useState<unknown>(null);
  async function toggle(key: string, enabled: boolean) {
    setError(null);
    try { await api('PATCH', `/api/admin/feature-flags/${key}`, { enabled }); list.reload(); } catch (e) { setError(e); }
  }
  async function add(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const f = Object.fromEntries(new FormData(form)) as Record<string, string>;
    setError(null);
    try { await api('POST', '/api/admin/feature-flags', { key: f.key, nameAr: f.nameAr, description: f.description || null, enabled: false }); form.reset(); list.reload(); } catch (err) { setError(err); }
  }
  return (
    <div className="card">
      <ErrorBox error={error ?? list.error} />
      <p className="muted small">القيمة العامة تطبق على كل المنشآت إلا من خُصّصت له قيمة من صفحة المنشأة.</p>
      <table className="table"><thead><tr><th>المفتاح</th><th>الاسم</th><th>تخصيصات</th><th>عام</th></tr></thead>
        <tbody>{list.data?.data.map((f) => (
          <tr key={f.key}><td dir="ltr" className="code">{f.key}</td><td>{f.nameAr}<div className="muted small">{f.description}</div></td><td>{f.overrides}</td>
            <td><label className="check"><input type="checkbox" checked={f.enabled} onChange={(e) => toggle(f.key, e.target.checked)} />{f.enabled ? 'مفعّلة' : 'معطّلة'}</label></td></tr>
        ))}</tbody>
      </table>
      <form className="inline-form" onSubmit={add}>
        <input name="key" required dir="ltr" pattern="[a-z][a-z0-9_]{1,62}" placeholder="new_feature" aria-label="المفتاح" />
        <input name="nameAr" required minLength={2} placeholder="الاسم" aria-label="الاسم" />
        <input name="description" placeholder="الوصف" aria-label="الوصف" />
        <button className="btn">إضافة خاصية</button>
      </form>
    </div>
  );
}
