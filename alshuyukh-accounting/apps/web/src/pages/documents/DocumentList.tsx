import { useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { STATUS_AR, statusClass, type DocConfig, type Doc } from './config';

const PAGE = 50;

export default function DocumentList({ config }: { config: DocConfig }) {
  const { can } = useAuth();
  const [status, setStatus] = useState('');
  const [search, setSearch] = useState('');
  const [offset, setOffset] = useState(0);
  const list = useLoad(() => {
    const q = new URLSearchParams({ limit: String(PAGE), offset: String(offset) });
    if (status === 'OPEN') q.set('open', 'true');
    else if (status) q.set('status', status);
    if (search) q.set('search', search);
    return api<{ data: Doc[]; total: number }>('GET', `/api/${config.api}?${q}`);
  }, [config.api, status, search, offset]);
  const statuses = Object.keys(STATUS_AR).filter((s) => s !== 'VOIDED');

  return (
    <div className="card">
      <div className="toolbar">
        <input type="search" className="grow" placeholder="بحث بالرقم أو الاسم"
          onKeyDown={(e) => { if (e.key === 'Enter') { setSearch(e.currentTarget.value); setOffset(0); } }} />
        <select value={status} onChange={(e) => { setStatus(e.target.value); setOffset(0); }} aria-label="الحالة">
          <option value="">كل الحالات</option>
          {config.legal && !config.isReturn && <option value="OPEN">غير مسددة</option>}
          {statuses.map((s) => <option key={s} value={s}>{STATUS_AR[s]}</option>)}
        </select>
        {can(config.perm.create) && !config.isReturn && <Link className="btn btn-primary" to="new">{config.newLabel}</Link>}
      </div>
      <ErrorBox error={list.error} />
      <div className="table-wrap">
        <table className="table">
          <thead>
            <tr>
              <th>الرقم</th><th>التاريخ</th><th>{config.party === 'customer' ? 'العميل' : 'المورد'}</th>
              <th className="num">الإجمالي</th>{config.legal && <th className="num">المتبقي</th>}<th>الحالة</th>
            </tr>
          </thead>
          <tbody>
            {!list.data && !list.error && <tr><td colSpan={6} className="muted empty-row">جارٍ التحميل…</td></tr>}
            {list.data?.data.map((d) => (
              <tr key={d.id}>
                <td dir="ltr" className="code"><Link to={d.id}>{d.number ?? 'مسودة'}</Link></td>
                <td dir="ltr">{d.date}</td>
                <td><Link to={d.id} className="plain-link">{d.partyName}</Link></td>
                <td className="num" dir="ltr">{formatAmount(d.total)}</td>
                {config.legal && <td className="num" dir="ltr">{d.status === 'DRAFT' || d.status === 'CANCELLED' ? '' : formatAmount(d.remainingAmount)}</td>}
                <td><span className={`tag ${statusClass(d.status)}`}>{STATUS_AR[d.status] ?? d.status}</span></td>
              </tr>
            ))}
            {list.data?.data.length === 0 && <tr><td colSpan={6} className="muted empty-row">لا توجد مستندات</td></tr>}
          </tbody>
        </table>
      </div>
      {list.data && list.data.total > PAGE && (
        <div className="pager">
          <span className="muted small">{offset + 1}–{Math.min(offset + PAGE, list.data.total)} من {list.data.total}</span>
          <button className="btn" disabled={offset === 0} onClick={() => setOffset(offset - PAGE)}>السابق</button>
          <button className="btn" disabled={offset + PAGE >= list.data.total} onClick={() => setOffset(offset + PAGE)}>التالي</button>
        </div>
      )}
    </div>
  );
}
