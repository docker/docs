import { Fragment, useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { LedgerTable, ledgerCsv, useUrlRange, type LedgerData } from './Ledger';
import { Amount, docLink, RangeBar, ReportFrame } from './shared';

type PartyType = 'customer' | 'supplier';
const BUCKETS: [string, string][] = [['current', 'غير مستحق'], ['days1to30', '1–30 يومًا'], ['days31to60', '31–60'], ['days61to90', '61–90'], ['over90', 'أكثر من 90']];
interface AgingRow extends Record<string, unknown> {
  partyId: string; code: string; nameAr: string; unapplied: string; balance: string;
  documents: { type: string; id: string; number: string; date: string; dueDate: string; total: string; remaining: string; daysOverdue: number }[];
}

export function Aging({ party }: { party: PartyType }) {
  const r = useLoad(() => api<{ asOf: string; data: AgingRow[]; totals: Record<string, string> }>('GET', `/api/reports/${party === 'customer' ? 'receivables' : 'payables'}-aging`), [party]);
  const [open, setOpen] = useState<string | null>(null);
  const d = r.data;
  const title = party === 'customer' ? 'أعمار الذمم المدينة (العملاء)' : 'أعمار الذمم الدائنة (الموردون)';
  return (
    <ReportFrame title={title} subtitle={d ? `كما في ${d.asOf}` : undefined} csv={d ? () => [
      ['الرمز', 'الاسم', ...BUCKETS.map(([, l]) => l), 'أرصدة غير مخصصة', 'الرصيد'],
      ...d.data.map((p) => [p.code, p.nameAr, ...BUCKETS.map(([k]) => p[k] as string), p.unapplied, p.balance]),
      ['', 'الإجمالي', ...BUCKETS.map(([k]) => d.totals[k]), d.totals.unapplied, d.totals.balance],
    ] : undefined}>
      <div className="card">
        <ErrorBox error={r.error} />
        {d && (
          <>
            <p className="muted small">كما في <span dir="ltr">{d.asOf}</span>، حسب تاريخ الاستحقاق. الرصيد من دفتر الأستاذ؛ «غير مخصص» يشمل الدفعات المقدمة والإشعارات غير المطبقة.</p>
            <div className="table-wrap">
              <table className="table report-table">
                <thead><tr><th>{party === 'customer' ? 'العميل' : 'المورد'}</th>{BUCKETS.map(([k, l]) => <th key={k} className="num">{l}</th>)}<th className="num">غير مخصص</th><th className="num">الرصيد</th></tr></thead>
                <tbody>
                  {d.data.map((p) => (
                    <Fragment key={p.partyId}>
                      <tr>
                        <td>
                          {p.documents.length > 0 && <button className="icon-btn no-print" aria-label="عرض المستندات" onClick={() => setOpen(open === p.partyId ? null : p.partyId)}>{open === p.partyId ? '−' : '+'}</button>}
                          <Link className="plain-link" to={`/reports/${party}-statement?partyId=${p.partyId}`}>{p.code} — {p.nameAr}</Link>
                        </td>
                        {BUCKETS.map(([k]) => <Amount key={k} v={p[k] as string} />)}
                        <Amount v={p.unapplied} /><Amount v={p.balance} strong />
                      </tr>
                      {open === p.partyId && p.documents.map((doc) => (
                        <tr key={doc.id} className="sub-detail">
                          <td><Link to={docLink(doc.type, doc.id) ?? '#'} dir="ltr" className="code">{doc.number}</Link> <span className="muted small">استحقاق <span dir="ltr">{doc.dueDate}</span>{doc.daysOverdue > 0 ? ` — متأخر ${doc.daysOverdue} يومًا` : ''}</span></td>
                          <td colSpan={BUCKETS.length + 1} className="muted small">الإجمالي <span dir="ltr">{formatAmount(doc.total)}</span></td>
                          <Amount v={doc.remaining} />
                        </tr>
                      ))}
                    </Fragment>
                  ))}
                  {d.data.length === 0 && <tr><td colSpan={BUCKETS.length + 3} className="muted empty-row">لا توجد أرصدة</td></tr>}
                  <tr className="total-row emphasis"><td>الإجمالي</td>{BUCKETS.map(([k]) => <Amount key={k} v={d.totals[k]} />)}<Amount v={d.totals.unapplied} /><Amount v={d.totals.balance} strong /></tr>
                </tbody>
              </table>
            </div>
          </>
        )}
      </div>
    </ReportFrame>
  );
}

export function Statement({ party }: { party: PartyType }) {
  const { params, range, setRange, set, query } = useUrlRange();
  const partyId = params.get('partyId') ?? '';
  const parties = useLoad(() => api<{ data: { id: string; code: string; nameAr: string }[] }>('GET', `/api/${party}s?limit=200`), [party]);
  const r = useLoad(() => (partyId
    ? api<LedgerData & { party: { code: string; nameAr: string; vatNumber: string | null } }>('GET', `/api/reports/${party}-statement?${query}&partyId=${partyId}`)
    : Promise.resolve(null)), [party, partyId, query]);
  const d = r.data;
  const label = party === 'customer' ? 'العميل' : 'المورد';
  return (
    <ReportFrame title={d ? `كشف حساب ${label} — ${d.party.nameAr}` : `كشف حساب ${label}`} subtitle={`من ${range.from} إلى ${range.to}`} csv={d ? () => ledgerCsv(d) : undefined}>
      <div className="card">
        <RangeBar range={range} setRange={setRange}>
          <select aria-label={label} value={partyId} onChange={(e) => set({ partyId: e.target.value })}>
            <option value="">— اختر {label} —</option>
            {parties.data?.data.map((p) => <option key={p.id} value={p.id}>{p.code} — {p.nameAr}</option>)}
          </select>
        </RangeBar>
        <ErrorBox error={r.error ?? parties.error} />
        {!partyId && <p className="muted">اختر {label} لعرض كشف الحساب.</p>}
        {d && <>
          <dl className="meta">
            <div><dt>{label}</dt><dd>{d.party.code} — {d.party.nameAr}</dd></div>
            {d.party.vatNumber && <div><dt>الرقم الضريبي</dt><dd dir="ltr">{d.party.vatNumber}</dd></div>}
            <div><dt>الرصيد الختامي</dt><dd dir="ltr">{formatAmount(d.closingBalance)}</dd></div>
          </dl>
          <LedgerTable d={d} />
          <p className="muted small">{party === 'customer' ? 'الرصيد الموجب مستحق على العميل.' : 'الرصيد الموجب مستحق للمورد.'}</p>
        </>}
      </div>
    </ReportFrame>
  );
}
