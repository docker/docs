import { useMemo, useState, type FormEvent } from 'react';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { ACCOUNT_TYPE_AR, formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';

export interface Account {
  id: string; code: string; nameAr: string; nameEn: string | null; type: string; parentId: string | null;
  level: number; isPostable: boolean; isActive: boolean; isSystem: boolean; systemKey: string | null; balance: string;
}

export default function ChartOfAccounts() {
  const { can } = useAuth();
  const [showInactive, setShowInactive] = useState(false);
  const accounts = useLoad(() => api<{ data: Account[] }>('GET', `/api/accounts?includeInactive=${showInactive}`), [showInactive]);
  const [error, setError] = useState<unknown>(null);
  const [editing, setEditing] = useState<string | null>(null);
  const manage = can('account.manage');

  // Depth-first order so children follow their parent.
  const ordered = useMemo(() => {
    const list = accounts.data?.data ?? [];
    const children = new Map<string | null, Account[]>();
    for (const a of list) children.set(a.parentId, [...(children.get(a.parentId) ?? []), a]);
    const out: Account[] = [];
    const walk = (parent: string | null) => (children.get(parent) ?? []).sort((x, y) => x.code.localeCompare(y.code)).forEach((a) => { out.push(a); walk(a.id); });
    walk(null);
    // Accounts whose parent is filtered out (inactive) still appear.
    for (const a of list) if (!out.includes(a)) out.push(a);
    return out;
  }, [accounts.data]);

  async function run(fn: () => Promise<unknown>) {
    setError(null);
    try { await fn(); accounts.reload(); return true; } catch (e) { setError(e); return false; }
  }

  async function add(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const f = Object.fromEntries(new FormData(form)) as Record<string, string>;
    const ok = await run(() => api('POST', '/api/accounts', {
      code: f.code, nameAr: f.nameAr, nameEn: f.nameEn || null, parentId: f.parentId || null,
      type: f.parentId ? undefined : f.type, isPostable: f.kind === 'postable',
    }));
    if (ok) form.reset();
  }

  async function saveEdit(e: FormEvent<HTMLFormElement>, a: Account) {
    e.preventDefault();
    const f = Object.fromEntries(new FormData(e.currentTarget)) as Record<string, string>;
    if (await run(() => api('PATCH', `/api/accounts/${a.id}`, { code: f.code, nameAr: f.nameAr }))) setEditing(null);
  }

  const headers = ordered.filter((a) => !a.isPostable && a.isActive);

  return (
    <>
      <ErrorBox error={error ?? accounts.error} />
      <div className="card">
        <div className="toolbar">
          <label className="check"><input type="checkbox" checked={showInactive} onChange={(e) => setShowInactive(e.target.checked)} />إظهار الحسابات غير النشطة</label>
        </div>
        <div className="table-wrap">
          <table className="table">
            <thead><tr><th>الرمز</th><th>اسم الحساب</th><th>النوع</th><th className="num">الرصيد</th><th /></tr></thead>
            <tbody>
              {ordered.map((a) => (
                <tr key={a.id} className={`${a.isPostable ? '' : 'row-header'} ${a.isActive ? '' : 'row-muted'}`}>
                  {editing === a.id ? (
                    <td colSpan={5}>
                      <form className="inline-form" onSubmit={(e) => saveEdit(e, a)}>
                        <input name="code" defaultValue={a.code} required dir="ltr" />
                        <input name="nameAr" defaultValue={a.nameAr} required />
                        <button className="btn btn-primary">حفظ</button>
                        <button type="button" className="btn btn-ghost" onClick={() => setEditing(null)}>إلغاء</button>
                      </form>
                    </td>
                  ) : (
                    <>
                      <td dir="ltr" className="code">{a.code}</td>
                      <td style={{ paddingInlineStart: `${(a.level - 1) * 22 + 8}px` }}>
                        {a.nameAr}
                        {a.isSystem && <span className="tag tag-soft" title={a.systemKey ?? ''}>نظامي</span>}
                        {!a.isActive && <span className="tag tag-off">غير نشط</span>}
                      </td>
                      <td className="muted small">{ACCOUNT_TYPE_AR[a.type]}</td>
                      <td className="num" dir="ltr">{a.isPostable ? formatAmount(a.balance) : ''}</td>
                      <td className="row-actions">
                        {manage && <button className="btn btn-small" onClick={() => setEditing(a.id)}>تعديل</button>}
                        {manage && !a.isSystem && a.isPostable && (
                          <button className="btn btn-small" onClick={() => run(() => api('PATCH', `/api/accounts/${a.id}`, { isActive: !a.isActive }))}>
                            {a.isActive ? 'إيقاف' : 'تفعيل'}
                          </button>
                        )}
                      </td>
                    </>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <p className="muted small">الرصيد = المدين − الدائن من القيود المرحّلة. الرصيد السالب يعني رصيدًا دائنًا.</p>
      </div>

      {manage && (
        <form className="card form" onSubmit={add}>
          <h2>إضافة حساب</h2>
          <div className="grid-2">
            <label>الحساب الرئيسي
              <select name="parentId" defaultValue="">
                <option value="">— حساب رئيسي جديد —</option>
                {headers.map((h) => <option key={h.id} value={h.id}>{h.code} — {h.nameAr}</option>)}
              </select>
            </label>
            <label>النوع (للحسابات الرئيسية فقط)
              <select name="type" defaultValue="EXPENSE">
                {Object.entries(ACCOUNT_TYPE_AR).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
              </select>
            </label>
            <label>الرمز<input name="code" required pattern="[0-9A-Za-z.\-]{1,20}" dir="ltr" /></label>
            <label>الاسم بالعربية<input name="nameAr" required /></label>
            <label>الاسم بالإنجليزية<input name="nameEn" dir="ltr" /></label>
            <label>طبيعة الحساب
              <select name="kind" defaultValue="postable">
                <option value="postable">حساب فرعي (يقبل القيود)</option>
                <option value="header">حساب تجميعي (لا يقبل القيود)</option>
              </select>
            </label>
          </div>
          <button className="btn btn-primary">إضافة</button>
        </form>
      )}
    </>
  );
}
