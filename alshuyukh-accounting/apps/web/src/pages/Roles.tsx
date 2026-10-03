import { useState, type FormEvent } from 'react';
import { api } from '../api';
import { useAuth } from '../auth';
import { ErrorBox, PageHeader, useLoad } from '../ui';

interface Role { id: string; code: string; nameAr: string; nameEn: string; isSystem: boolean; permissions: string[] }
interface Permission { code: string; module: string; descriptionAr: string }

export default function Roles() {
  const { can } = useAuth();
  const roles = useLoad(() => api<{ data: Role[] }>('GET', '/api/roles'));
  const perms = useLoad(() => api<{ data: Permission[] }>('GET', '/api/permissions'));
  const [open, setOpen] = useState<string | null>(null);
  const [error, setError] = useState<unknown>(null);
  const manage = can('role.manage');

  const byModule = (perms.data?.data ?? []).reduce<Record<string, Permission[]>>((acc, p) => {
    (acc[p.module] ??= []).push(p);
    return acc;
  }, {});
  const label = (code: string) => perms.data?.data.find((p) => p.code === code)?.descriptionAr ?? code;

  async function create(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const f = new FormData(form);
    setError(null);
    try {
      await api('POST', '/api/roles', {
        code: String(f.get('code')), nameAr: String(f.get('nameAr')), nameEn: String(f.get('nameEn')),
        permissions: f.getAll('permissions').map(String),
      });
      form.reset();
      roles.reload();
    } catch (err) { setError(err); }
  }

  async function remove(id: string) {
    setError(null);
    try { await api('DELETE', `/api/roles/${id}`); roles.reload(); } catch (err) { setError(err); }
  }

  return (
    <>
      <PageHeader title="الأدوار والصلاحيات" />
      <ErrorBox error={error ?? roles.error} />
      <div className="card">
        {roles.data?.data.map((r) => (
          <div key={r.id} className="role-row">
            <button className="role-head" onClick={() => setOpen(open === r.id ? null : r.id)} aria-expanded={open === r.id}>
              <span>{r.nameAr}</span>
              <span className="muted small" dir="ltr">{r.code}</span>
              <span className={`tag ${r.isSystem ? '' : 'tag-ok'}`}>{r.isSystem ? 'نظامي' : 'مخصص'}</span>
              <span className="muted small">{r.permissions.length} صلاحية</span>
            </button>
            {open === r.id && (
              <div className="role-body">
                <ul className="perm-list">{r.permissions.map((p) => <li key={p}>{label(p)}</li>)}</ul>
                {manage && !r.isSystem && <button className="btn btn-danger btn-small" onClick={() => remove(r.id)}>حذف الدور</button>}
              </div>
            )}
          </div>
        ))}
      </div>

      {manage && (
        <form className="card form" onSubmit={create}>
          <h2>دور مخصص جديد</h2>
          <div className="grid-2">
            <label>الرمز<input name="code" required pattern="[A-Za-z][A-Za-z0-9_]{1,62}" dir="ltr" placeholder="BRANCH_AUDITOR" /></label>
            <label>الاسم بالعربية<input name="nameAr" required minLength={2} /></label>
            <label>الاسم بالإنجليزية<input name="nameEn" required minLength={2} dir="ltr" /></label>
          </div>
          <p className="muted small">لا يمكنك منح صلاحيات لا تملكها.</p>
          <div className="perm-grid">
            {Object.entries(byModule).map(([mod, list]) => (
              <fieldset key={mod}>
                <legend dir="ltr">{mod}</legend>
                {list.map((p) => <label key={p.code} className="check"><input type="checkbox" name="permissions" value={p.code} disabled={!can(p.code)} />{p.descriptionAr}</label>)}
              </fieldset>
            ))}
          </div>
          <button className="btn btn-primary">إنشاء الدور</button>
        </form>
      )}
    </>
  );
}
