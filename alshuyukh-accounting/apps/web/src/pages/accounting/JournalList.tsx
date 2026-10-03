import { useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { ENTRY_STATUS_AR, REFERENCE_AR, formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';

interface Row { id: string; entryNumber: string | null; entryDate: string; description: string; referenceType: string; status: string; totalDebit: string }

const PAGE = 50;

export default function JournalList() {
  const { can } = useAuth();
  const [status, setStatus] = useState('');
  const [search, setSearch] = useState('');
  const [offset, setOffset] = useState(0);
  const list = useLoad(() => {
    const q = new URLSearchParams({ limit: String(PAGE), offset: String(offset) });
    if (status) q.set('status', status);
    if (search) q.set('search', search);
    return api<{ data: Row[]; total: number }>('GET', `/api/journal-entries?${q}`);
  }, [status, search, offset]);

  return (
    <div className="card">
      <div className="toolbar">
        <select value={status} onChange={(e) => { setStatus(e.target.value); setOffset(0); }} aria-label="الحالة">
          <option value="">كل الحالات</option>
          {Object.entries(ENTRY_STATUS_AR).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
        </select>
        <input type="search" placeholder="بحث بالوصف أو الرقم" defaultValue={search}
          onKeyDown={(e) => { if (e.key === 'Enter') { setSearch(e.currentTarget.value); setOffset(0); } }} />
        <span className="spacer" />
        {can('journal.create') && <Link className="btn btn-primary" to="new">قيد جديد</Link>}
      </div>
      <ErrorBox error={list.error} />
      <div className="table-wrap">
        <table className="table">
          <thead><tr><th>الرقم</th><th>التاريخ</th><th>الوصف</th><th>المصدر</th><th className="num">المبلغ</th><th>الحالة</th></tr></thead>
          <tbody>
            {list.data?.data.map((e) => (
              <tr key={e.id}>
                <td dir="ltr" className="code"><Link to={e.id}>{e.entryNumber ?? 'مسودة'}</Link></td>
                <td dir="ltr">{e.entryDate}</td>
                <td><Link to={e.id} className="plain-link">{e.description}</Link></td>
                <td className="muted small">{REFERENCE_AR[e.referenceType] ?? e.referenceType}</td>
                <td className="num" dir="ltr">{formatAmount(e.totalDebit)}</td>
                <td><span className={`tag status-${e.status.toLowerCase()}`}>{ENTRY_STATUS_AR[e.status]}</span></td>
              </tr>
            ))}
            {list.data?.data.length === 0 && <tr><td colSpan={6} className="muted empty-row">لا توجد قيود</td></tr>}
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
