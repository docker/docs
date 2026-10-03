import { useState } from 'react';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { ErrorBox, useLoad } from '../../ui';

interface Year {
  id: string; name: string; startDate: string; endDate: string; status: string; closingEntryId: string | null;
  periods: { id: string; periodNumber: number; name: string; startDate: string; endDate: string; status: string }[];
}

export default function FiscalYears() {
  const { can } = useAuth();
  const years = useLoad(() => api<{ data: Year[] }>('GET', '/api/fiscal-years'));
  const [error, setError] = useState<unknown>(null);
  const [open, setOpen] = useState<string | null>(null);
  const manage = can('fiscal.manage');

  async function run(fn: () => Promise<unknown>) {
    setError(null);
    try { await fn(); years.reload(); } catch (e) { setError(e); }
  }

  const latest = years.data?.data[0];
  const nextStart = latest ? new Date(`${latest.endDate}T00:00:00Z`) : null;
  nextStart?.setUTCDate(nextStart.getUTCDate() + 1);

  return (
    <>
      <ErrorBox error={error ?? years.error} />
      {manage && nextStart && (
        <div className="toolbar">
          <button className="btn" onClick={() => run(() => api('POST', '/api/fiscal-years', { startDate: nextStart.toISOString().slice(0, 10) }))}>
            إنشاء السنة المالية التالية (تبدأ <span dir="ltr">{nextStart.toISOString().slice(0, 10)}</span>)
          </button>
        </div>
      )}
      {years.data?.data.map((y) => (
        <div key={y.id} className="card">
          <div className="entry-head">
            <button className="role-head" onClick={() => setOpen(open === y.id ? null : y.id)} aria-expanded={open === y.id}>
              <span>السنة المالية {y.name}</span>
              <span className="muted small" dir="ltr">{y.startDate} → {y.endDate}</span>
              <span className={`tag ${y.status === 'OPEN' ? 'tag-ok' : 'tag-off'}`}>{y.status === 'OPEN' ? 'مفتوحة' : 'مقفلة'}</span>
            </button>
            {manage && y.status === 'OPEN' && (
              <button className="btn btn-danger btn-small" onClick={() => {
                if (confirm('سيتم ترحيل قيد الإقفال إلى الأرباح المحتجزة وإقفال جميع الفترات. لا يمكن التراجع. متابعة؟')) void run(() => api('POST', `/api/fiscal-years/${y.id}/close`));
              }}>إقفال السنة</button>
            )}
          </div>
          {open === y.id && (
            <table className="table">
              <thead><tr><th>#</th><th>الفترة</th><th>من</th><th>إلى</th><th>الحالة</th><th /></tr></thead>
              <tbody>
                {y.periods.map((p) => (
                  <tr key={p.id}>
                    <td>{p.periodNumber}</td><td>{p.name}</td>
                    <td dir="ltr">{p.startDate}</td><td dir="ltr">{p.endDate}</td>
                    <td><span className={`tag ${p.status === 'OPEN' ? 'tag-ok' : 'tag-off'}`}>{p.status === 'OPEN' ? 'مفتوحة' : 'مقفلة'}</span></td>
                    <td>
                      {manage && y.status === 'OPEN' && (
                        <button className="btn btn-small" onClick={() => run(() => api('POST', `/api/fiscal-periods/${p.id}/${p.status === 'OPEN' ? 'close' : 'reopen'}`))}>
                          {p.status === 'OPEN' ? 'إقفال' : 'إعادة فتح'}
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      ))}
    </>
  );
}
