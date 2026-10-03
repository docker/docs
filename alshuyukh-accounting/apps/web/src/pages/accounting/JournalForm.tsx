import { useState, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { formatAmount, fromHalalas, toHalalas } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import type { Account } from './ChartOfAccounts';
import type { Entry } from './JournalDetail';

interface Line { accountId: string; debit: string; credit: string; description: string }

const emptyLine = (): Line => ({ accountId: '', debit: '', credit: '', description: '' });

export default function JournalForm({ entry, onSaved }: { entry?: Entry; onSaved?: () => void }) {
  const { can } = useAuth();
  const navigate = useNavigate();
  const accounts = useLoad(() => api<{ data: Account[] }>('GET', '/api/accounts?postableOnly=true'));
  const [lines, setLines] = useState<Line[]>(() =>
    entry ? entry.lines.map((l) => ({
      accountId: l.accountId, debit: l.debit === '0.00' ? '' : l.debit, credit: l.credit === '0.00' ? '' : l.credit, description: l.description ?? '',
    })) : [emptyLine(), emptyLine()]);
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState(false);

  const update = (i: number, patch: Partial<Line>) => setLines((ls) => ls.map((l, j) => (j === i ? { ...l, ...patch } : l)));
  const debit = lines.reduce((s, l) => s + (toHalalas(l.debit) ?? 0), 0);
  const credit = lines.reduce((s, l) => s + (toHalalas(l.credit) ?? 0), 0);
  const diff = debit - credit;

  async function submit(e: FormEvent<HTMLFormElement>, postNow: boolean) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const payload = {
      entryDate: String(f.get('entryDate')),
      description: String(f.get('description')),
      lines: lines.filter((l) => l.accountId).map((l) => ({
        accountId: l.accountId,
        ...(l.debit.trim() ? { debit: l.debit.trim() } : {}),
        ...(l.credit.trim() ? { credit: l.credit.trim() } : {}),
        description: l.description || null,
      })),
    };
    setBusy(true);
    setError(null);
    try {
      if (entry) {
        await api('PATCH', `/api/journal-entries/${entry.id}`, payload);
        if (postNow) await api('POST', `/api/journal-entries/${entry.id}/post`);
        onSaved?.();
      } else {
        const created = await api<{ id: string }>('POST', '/api/journal-entries', { ...payload, post: postNow });
        navigate(`/accounting/journal/${created.id}`);
      }
    } catch (err) {
      setError(err);
    } finally {
      setBusy(false);
    }
  }

  const today = new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh' }).format(new Date());

  return (
    <form className="card form" onSubmit={(e) => submit(e, (e.nativeEvent as SubmitEvent).submitter?.getAttribute('value') === 'post')}>
      <h2>{entry ? 'تعديل مسودة القيد' : 'قيد يومية جديد'}</h2>
      <ErrorBox error={error ?? accounts.error} />
      <div className="grid-2">
        <label>التاريخ<input name="entryDate" type="date" required defaultValue={entry?.entryDate ?? today} /></label>
        <label>البيان<input name="description" required maxLength={1000} defaultValue={entry?.description ?? ''} /></label>
      </div>

      <div className="table-wrap">
        <table className="table lines-table">
          <thead><tr><th>الحساب</th><th>مدين</th><th>دائن</th><th>البيان</th><th /></tr></thead>
          <tbody>
            {lines.map((l, i) => (
              <tr key={i}>
                <td>
                  <select value={l.accountId} onChange={(e) => update(i, { accountId: e.target.value })} aria-label={`حساب السطر ${i + 1}`}>
                    <option value="">— اختر الحساب —</option>
                    {accounts.data?.data.map((a) => <option key={a.id} value={a.id}>{a.code} — {a.nameAr}</option>)}
                  </select>
                </td>
                <td><input inputMode="decimal" dir="ltr" value={l.debit} aria-label={`مدين السطر ${i + 1}`}
                  aria-invalid={toHalalas(l.debit) === null} onChange={(e) => update(i, { debit: e.target.value, credit: e.target.value ? '' : l.credit })} /></td>
                <td><input inputMode="decimal" dir="ltr" value={l.credit} aria-label={`دائن السطر ${i + 1}`}
                  aria-invalid={toHalalas(l.credit) === null} onChange={(e) => update(i, { credit: e.target.value, debit: e.target.value ? '' : l.debit })} /></td>
                <td><input value={l.description} onChange={(e) => update(i, { description: e.target.value })} aria-label={`بيان السطر ${i + 1}`} /></td>
                <td>{lines.length > 2 && <button type="button" className="icon-btn" aria-label="حذف السطر" onClick={() => setLines((ls) => ls.filter((_, j) => j !== i))}>×</button>}</td>
              </tr>
            ))}
          </tbody>
          <tfoot>
            <tr>
              <td><button type="button" className="btn btn-small" onClick={() => setLines((ls) => [...ls, emptyLine()])}>+ سطر</button></td>
              <td className="num" dir="ltr"><strong>{formatAmount(fromHalalas(debit))}</strong></td>
              <td className="num" dir="ltr"><strong>{formatAmount(fromHalalas(credit))}</strong></td>
              <td colSpan={2}>
                {diff === 0 && debit > 0
                  ? <span className="tag tag-ok">متوازن</span>
                  : <span className="tag tag-off">الفرق: <span dir="ltr">{formatAmount(fromHalalas(Math.abs(diff)))}</span></span>}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
      <p className="muted small">يتحقق الخادم من توازن القيد وصحة الحسابات والفترة المالية عند الترحيل.</p>
      <div className="form-actions">
        <button className="btn" disabled={busy} value="draft">حفظ كمسودة</button>
        {can('journal.post') && <button className="btn btn-primary" disabled={busy} value="post">حفظ وترحيل</button>}
      </div>
    </form>
  );
}
