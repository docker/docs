import { useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { ENTRY_STATUS_AR, REFERENCE_AR, formatAmount } from '../../money';
import { ErrorBox, formatDateTime, useLoad } from '../../ui';
import JournalForm from './JournalForm';

export interface Entry {
  id: string; entryNumber: string | null; entryDate: string; description: string; referenceType: string; referenceId: string | null;
  source: string; status: string; currency: string; totalDebit: string; totalCredit: string;
  reversalOfId: string | null; reversedByEntryId: string | null; correctionOfId: string | null; postedAt: string | null; createdAt: string;
  lines: { id: string; lineNo: number; accountId: string; accountCode: string; accountName: string; debit: string; credit: string; description: string | null }[];
}

export default function JournalDetail() {
  const { id } = useParams();
  const { can } = useAuth();
  const navigate = useNavigate();
  const entry = useLoad(() => api<Entry>('GET', `/api/journal-entries/${id}`), [id]);
  const [error, setError] = useState<unknown>(null);
  const [editing, setEditing] = useState(false);
  const [reversing, setReversing] = useState(false);

  async function run(fn: () => Promise<unknown>) {
    setError(null);
    try { await fn(); entry.reload(); } catch (e) { setError(e); }
  }

  const e = entry.data;
  if (!e) return <ErrorBox error={entry.error} />;
  if (editing) return <JournalForm entry={e} onSaved={() => { setEditing(false); entry.reload(); }} />;

  const manual = e.source === 'MANUAL';
  return (
    <div className="card">
      <ErrorBox error={error} />
      <div className="entry-head">
        <div>
          <h2 dir="ltr" className="code">{e.entryNumber ?? 'مسودة'}</h2>
          <p>{e.description}</p>
        </div>
        <span className={`tag status-${e.status.toLowerCase()}`}>{ENTRY_STATUS_AR[e.status]}</span>
      </div>
      <dl className="meta">
        <div><dt>التاريخ</dt><dd dir="ltr">{e.entryDate}</dd></div>
        <div><dt>المصدر</dt><dd>{REFERENCE_AR[e.referenceType] ?? e.referenceType}</dd></div>
        <div><dt>العملة</dt><dd>{e.currency}</dd></div>
        <div><dt>تاريخ الترحيل</dt><dd>{formatDateTime(e.postedAt)}</dd></div>
        {e.reversalOfId && <div><dt>يعكس القيد</dt><dd><Link to={`/accounting/journal/${e.reversalOfId}`}>عرض الأصل</Link></dd></div>}
        {e.reversedByEntryId && <div><dt>عُكس بالقيد</dt><dd><Link to={`/accounting/journal/${e.reversedByEntryId}`}>عرض القيد العكسي</Link></dd></div>}
        {e.correctionOfId && <div><dt>تصحيح للقيد</dt><dd><Link to={`/accounting/journal/${e.correctionOfId}`}>عرض الأصل</Link></dd></div>}
      </dl>

      <div className="table-wrap">
        <table className="table">
          <thead><tr><th>#</th><th>الحساب</th><th>البيان</th><th className="num">مدين</th><th className="num">دائن</th></tr></thead>
          <tbody>
            {e.lines.map((l) => (
              <tr key={l.id}>
                <td>{l.lineNo}</td>
                <td><span dir="ltr" className="code">{l.accountCode}</span> {l.accountName}</td>
                <td className="muted">{l.description ?? ''}</td>
                <td className="num" dir="ltr">{l.debit === '0.00' ? '' : formatAmount(l.debit)}</td>
                <td className="num" dir="ltr">{l.credit === '0.00' ? '' : formatAmount(l.credit)}</td>
              </tr>
            ))}
          </tbody>
          <tfoot>
            <tr><td colSpan={3}><strong>الإجمالي</strong></td>
              <td className="num" dir="ltr"><strong>{formatAmount(e.totalDebit)}</strong></td>
              <td className="num" dir="ltr"><strong>{formatAmount(e.totalCredit)}</strong></td></tr>
          </tfoot>
        </table>
      </div>

      <div className="form-actions">
        {e.status === 'DRAFT' && manual && can('journal.create') && <button className="btn" onClick={() => setEditing(true)}>تعديل</button>}
        {e.status === 'DRAFT' && can('journal.post') && <button className="btn btn-primary" onClick={() => run(() => api('POST', `/api/journal-entries/${e.id}/post`))}>ترحيل</button>}
        {e.status === 'DRAFT' && manual && can('journal.create') && (
          <button className="btn btn-danger" onClick={() => { if (confirm('حذف المسودة؟')) void run(async () => { await api('DELETE', `/api/journal-entries/${e.id}`); navigate('/accounting/journal'); }); }}>حذف المسودة</button>
        )}
        {e.status === 'POSTED' && manual && e.referenceType !== 'REVERSAL' && can('journal.reverse') && !reversing && (
          <button className="btn btn-danger" onClick={() => setReversing(true)}>عكس القيد</button>
        )}
      </div>
      {e.status === 'POSTED' && !manual && <p className="muted small">هذا القيد صادر عن مستند. يُصحَّح بإلغاء المستند أو إرجاعه.</p>}

      {reversing && (
        <form className="reverse-box" onSubmit={(ev) => {
          ev.preventDefault();
          const f = new FormData(ev.currentTarget);
          void run(async () => {
            const r = await api<{ reversal: { id: string } }>('POST', `/api/journal-entries/${e.id}/reverse`, {
              reason: String(f.get('reason')), ...(f.get('date') ? { date: String(f.get('date')) } : {}),
            });
            setReversing(false);
            navigate(`/accounting/journal/${r.reversal.id}`);
          });
        }}>
          <p>لا يُعدَّل القيد المرحّل. سيُنشأ قيد عكسي مرحّل بنفس المبالغ معكوسة، ثم يمكنك إنشاء قيد تصحيح جديد.</p>
          <div className="grid-2">
            <label>سبب العكس<input name="reason" required minLength={3} /></label>
            <label>تاريخ القيد العكسي (اختياري)<input name="date" type="date" /></label>
          </div>
          <div className="form-actions">
            <button className="btn btn-danger">تأكيد العكس</button>
            <button type="button" className="btn btn-ghost" onClick={() => setReversing(false)}>إلغاء</button>
          </div>
        </form>
      )}
    </div>
  );
}
