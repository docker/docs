import { useState, type FormEvent } from 'react';
import { api } from '../api';
import { useAuth } from '../auth';
import { ErrorBox, PageHeader, formatDateTime, useLoad } from '../ui';

interface Member { id: string; email: string; fullName: string; status: string; isOwner: boolean; lastLoginAt: string | null; roles: { id: string; code: string; nameAr: string }[] }
interface Role { id: string; code: string; nameAr: string; isSystem: boolean }

export default function Users() {
  const { me, can } = useAuth();
  const users = useLoad(() => api<{ data: Member[] }>('GET', '/api/users'));
  const roles = useLoad(() => (can('role.view') ? api<{ data: Role[] }>('GET', '/api/roles') : Promise.resolve({ data: [] })));
  const [error, setError] = useState<unknown>(null);
  const [editing, setEditing] = useState<string | null>(null);
  const manage = can('user.manage');

  async function run(fn: () => Promise<unknown>) {
    setError(null);
    try { await fn(); users.reload(); return true; } catch (err) { setError(err); return false; }
  }

  async function add(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const f = new FormData(form);
    const ok = await run(() => api('POST', '/api/users', {
      email: String(f.get('email')), fullName: String(f.get('fullName')),
      initialPassword: String(f.get('initialPassword')) || undefined, roleIds: [String(f.get('roleId'))],
    }));
    if (ok) form.reset();
  }

  async function saveRoles(e: FormEvent<HTMLFormElement>, userId: string) {
    e.preventDefault();
    const roleIds = new FormData(e.currentTarget).getAll('roleIds').map(String);
    if (await run(() => api('PUT', `/api/users/${userId}/roles`, { roleIds }))) setEditing(null);
  }

  return (
    <>
      <PageHeader title="المستخدمون" />
      <ErrorBox error={error ?? users.error} />
      <div className="card">
        <div className="table-wrap">
          <table className="table">
            <thead><tr><th>الاسم</th><th>البريد</th><th>الأدوار</th><th>الحالة</th><th>آخر دخول</th><th /></tr></thead>
            <tbody>
              {users.data?.data.map((u) => (
                <tr key={u.id}>
                  <td>{u.fullName} {u.isOwner && <span className="tag">المالك</span>}</td>
                  <td dir="ltr">{u.email}</td>
                  <td>
                    {editing === u.id ? (
                      <form onSubmit={(e) => saveRoles(e, u.id)} className="role-picker">
                        {roles.data?.data.map((r) => (
                          <label key={r.id} className="check"><input type="checkbox" name="roleIds" value={r.id} defaultChecked={u.roles.some((x) => x.id === r.id)} />{r.nameAr}</label>
                        ))}
                        <button className="btn btn-primary">حفظ</button>
                        <button type="button" className="btn btn-ghost" onClick={() => setEditing(null)}>إلغاء</button>
                      </form>
                    ) : u.roles.map((r) => r.nameAr).join('، ')}
                  </td>
                  <td><span className={`tag ${u.status === 'ACTIVE' ? 'tag-ok' : 'tag-off'}`}>{u.status === 'ACTIVE' ? 'نشط' : 'معطّل'}</span></td>
                  <td>{formatDateTime(u.lastLoginAt)}</td>
                  <td className="row-actions">
                    {manage && u.id !== me?.user.id && !u.isOwner && (
                      <>
                        <button className="btn btn-small" onClick={() => setEditing(u.id)}>الأدوار</button>
                        <button className="btn btn-small" onClick={() => run(() => api('PATCH', `/api/users/${u.id}`, { status: u.status === 'ACTIVE' ? 'DISABLED' : 'ACTIVE' }))}>
                          {u.status === 'ACTIVE' ? 'تعطيل' : 'تفعيل'}
                        </button>
                      </>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {can('user.invite') && manage && (
        <form className="card form" onSubmit={add}>
          <h2>إضافة مستخدم</h2>
          <div className="grid-2">
            <label>الاسم الكامل<input name="fullName" required minLength={2} /></label>
            <label>البريد الإلكتروني<input name="email" type="email" required dir="ltr" /></label>
            <label>كلمة مرور مؤقتة<input name="initialPassword" type="password" minLength={10} dir="ltr" />
              <span className="hint">مطلوبة للحسابات الجديدة. سيُطلب من المستخدم تغييرها.</span></label>
            <label>الدور
              <select name="roleId" required>
                {roles.data?.data.filter((r) => r.code !== 'TENANT_OWNER').map((r) => <option key={r.id} value={r.id}>{r.nameAr}</option>)}
              </select>
            </label>
          </div>
          <button className="btn btn-primary">إضافة</button>
        </form>
      )}
    </>
  );
}
