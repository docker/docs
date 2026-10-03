import { Fragment, useState } from 'react';
import { api } from '../api';
import { ErrorBox, PageHeader, formatDateTime, useLoad } from '../ui';

interface Entry { id: string; action: string; entityType: string | null; entityId: string | null; userName: string | null; ipAddress: string | null; createdAt: string; oldValues: unknown; newValues: unknown }

const ACTIONS: Record<string, string> = {
  REGISTER: 'تسجيل منشأة', LOGIN: 'تسجيل دخول', LOGOUT: 'تسجيل خروج', LOGIN_FAILED: 'دخول فاشل',
  TOKEN_REUSE_DETECTED: 'إعادة استخدام رمز', TENANT_SWITCH: 'تبديل المنشأة', PASSWORD_CHANGE: 'تغيير كلمة المرور',
  CREATE: 'إنشاء', UPDATE: 'تعديل', DELETE: 'حذف', SETTINGS_CHANGE: 'تغيير الإعدادات', PERMISSION_CHANGE: 'تغيير الصلاحيات',
};

export default function AuditLog() {
  const [action, setAction] = useState('');
  const [cursor, setCursor] = useState<string | null>(null);
  const [open, setOpen] = useState<string | null>(null);
  const page = useLoad(() => {
    const q = new URLSearchParams({ limit: '50' });
    if (action) q.set('action', action);
    if (cursor) q.set('cursor', cursor);
    return api<{ data: Entry[]; nextCursor: string | null }>('GET', `/api/audit-logs?${q}`);
  }, [action, cursor]);

  return (
    <>
      <PageHeader title="سجل التدقيق">
        <select value={action} onChange={(e) => { setAction(e.target.value); setCursor(null); }} aria-label="نوع العملية">
          <option value="">كل العمليات</option>
          {Object.entries(ACTIONS).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
        </select>
      </PageHeader>
      <ErrorBox error={page.error} />
      <div className="card">
        <div className="table-wrap">
          <table className="table">
            <thead><tr><th>الوقت</th><th>المستخدم</th><th>العملية</th><th>الكيان</th><th>IP</th></tr></thead>
            <tbody>
              {page.data?.data.map((e) => (
                <Fragment key={e.id}>
                  <tr className="clickable" onClick={() => setOpen(open === e.id ? null : e.id)}>
                    <td>{formatDateTime(e.createdAt)}</td>
                    <td>{e.userName ?? '—'}</td>
                    <td>{ACTIONS[e.action] ?? e.action}</td>
                    <td dir="ltr">{e.entityType ?? '—'}</td>
                    <td dir="ltr">{e.ipAddress ?? '—'}</td>
                  </tr>
                  {open === e.id && (
                    <tr><td colSpan={5}>
                      <div className="diff">
                        <div><strong>قبل</strong><pre dir="ltr">{JSON.stringify(e.oldValues, null, 2)}</pre></div>
                        <div><strong>بعد</strong><pre dir="ltr">{JSON.stringify(e.newValues, null, 2)}</pre></div>
                      </div>
                    </td></tr>
                  )}
                </Fragment>
              ))}
            </tbody>
          </table>
        </div>
        <div className="pager">
          {cursor && <button className="btn" onClick={() => setCursor(null)}>الأحدث</button>}
          {page.data?.nextCursor && <button className="btn" onClick={() => setCursor(page.data!.nextCursor)}>الأقدم ←</button>}
        </div>
      </div>
    </>
  );
}
