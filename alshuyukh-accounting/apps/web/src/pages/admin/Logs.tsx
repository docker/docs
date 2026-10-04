import { Fragment, useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { ErrorBox, formatDateTime, useLoad } from '../../ui';

export function Logs() {
  const [tab, setTab] = useState<'platform' | 'audit' | 'errors'>('platform');
  return (
    <>
      <nav className="tabs">
        {([['platform', 'إجراءات مديري المنصة'], ['audit', 'سجل تدقيق المنشآت'], ['errors', 'أخطاء النظام']] as const).map(([k, l]) =>
          <button key={k} className={`tab ${tab === k ? 'active' : ''}`} onClick={() => setTab(k)}>{l}</button>)}
      </nav>
      {tab === 'platform' && <PlatformLogs />}
      {tab === 'audit' && <TenantAudit />}
      {tab === 'errors' && <Errors />}
    </>
  );
}

function PlatformLogs() {
  const l = useLoad(() => api<{ data: { id: string; adminEmail: string; action: string; entityType: string; entityId: string | null; tenantId: string | null; tenantName: string | null; newValues: unknown; createdAt: string }[] }>('GET', '/api/admin/platform-logs'));
  return (
    <div className="card"><ErrorBox error={l.error} />
      <div className="table-wrap"><table className="table">
        <thead><tr><th>الوقت</th><th>المدير</th><th>الإجراء</th><th>الهدف</th><th>التفاصيل</th></tr></thead>
        <tbody>{l.data?.data.map((x) => (
          <tr key={x.id}><td>{formatDateTime(x.createdAt)}</td><td dir="ltr" className="small">{x.adminEmail}</td><td className="code" dir="ltr">{x.action}</td>
            <td>{x.entityType}{x.tenantId && <> — <Link to={`/admin/tenants/${x.tenantId}`}>{x.tenantName}</Link></>}</td>
            <td dir="ltr" className="small code">{x.newValues ? JSON.stringify(x.newValues).slice(0, 160) : ''}</td></tr>
        ))}</tbody>
      </table></div>
    </div>
  );
}

function TenantAudit() {
  const l = useLoad(() => api<{ data: { id: string; tenantId: string; tenantName: string; userEmail: string | null; action: string; entityType: string; ipAddress: string | null; createdAt: string }[] }>('GET', '/api/admin/audit-logs?limit=200'));
  return (
    <div className="card"><ErrorBox error={l.error} />
      <div className="table-wrap"><table className="table">
        <thead><tr><th>الوقت</th><th>المنشأة</th><th>المستخدم</th><th>الإجراء</th><th>الكيان</th><th>IP</th></tr></thead>
        <tbody>{l.data?.data.map((x) => (
          <tr key={x.id}><td>{formatDateTime(x.createdAt)}</td><td><Link to={`/admin/tenants/${x.tenantId}`}>{x.tenantName}</Link></td>
            <td dir="ltr" className="small">{x.userEmail ?? '—'}</td><td className="code" dir="ltr">{x.action}</td><td>{x.entityType}</td><td dir="ltr" className="small">{x.ipAddress}</td></tr>
        ))}</tbody>
      </table></div>
    </div>
  );
}

function Errors() {
  const l = useLoad(() => api<{ data: { id: string; tenantName: string | null; requestId: string | null; method: string | null; path: string | null; errorCode: string | null; message: string; stack: string | null; createdAt: string }[] }>('GET', '/api/admin/errors'));
  const [open, setOpen] = useState<string | null>(null);
  return (
    <div className="card"><ErrorBox error={l.error} />
      <div className="table-wrap"><table className="table">
        <thead><tr><th>الوقت</th><th>الطلب</th><th>المنشأة</th><th>الرسالة</th></tr></thead>
        <tbody>{l.data?.data.map((x) => (
          <Fragment key={x.id}>
            <tr onClick={() => setOpen(open === x.id ? null : x.id)} className="clickable">
              <td>{formatDateTime(x.createdAt)}</td><td dir="ltr" className="small code">{x.method} {x.path}</td><td>{x.tenantName ?? '—'}</td>
              <td dir="ltr" className="small">{x.errorCode ? `[${x.errorCode}] ` : ''}{x.message}</td>
            </tr>
            {open === x.id && x.stack && <tr className="sub-detail"><td colSpan={4}><pre dir="ltr" className="stack">{x.stack}</pre><div className="muted small" dir="ltr">request {x.requestId}</div></td></tr>}
          </Fragment>
        ))}{l.data?.data.length === 0 && <tr><td colSpan={4} className="muted empty-row">لا توجد أخطاء مسجلة</td></tr>}</tbody>
      </table></div>
    </div>
  );
}
