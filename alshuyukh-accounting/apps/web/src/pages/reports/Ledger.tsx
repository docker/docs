import { useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { Amount, docLink, RangeBar, REF_AR, ReportFrame, preset } from './shared';

interface Line { entryId: string; entryNumber: string; date: string; description: string; referenceType: string; referenceId?: string; documentNumber: string | null; debit: string; credit: string; balance: string }
export interface LedgerData { openingBalance: string; lines: Line[]; totals: { debit: string; credit: string }; closingBalance: string }

/** Range state kept in the URL so links from other reports (and reloads) keep their filters. */
export function useUrlRange() {
  const [params, setParams] = useSearchParams();
  const year = preset('year');
  const range = { from: params.get('dateFrom') ?? year.from, to: params.get('dateTo') ?? year.to };
  const set = (patch: Record<string, string>) => setParams((p) => { for (const [k, v] of Object.entries(patch)) { if (v) p.set(k, v); else p.delete(k); } return p; }, { replace: true });
  return { params, range, setRange: (r: { from: string; to: string }) => set({ dateFrom: r.from, dateTo: r.to }), set, query: `dateFrom=${range.from}&dateTo=${range.to}` };
}

/** Ledger lines with opening, running and closing balance; shared by the general ledger and party statements. */
export function LedgerTable({ d }: { d: LedgerData }) {
  return (
    <div className="table-wrap">
      <table className="table report-table">
        <thead><tr><th>التاريخ</th><th>القيد</th><th>المستند</th><th>البيان</th><th className="num">مدين</th><th className="num">دائن</th><th className="num">الرصيد</th></tr></thead>
        <tbody>
          <tr className="total-row"><td colSpan={6}>الرصيد الافتتاحي</td><Amount v={d.openingBalance} /></tr>
          {d.lines.map((l, i) => {
            const link = docLink(l.referenceType, l.referenceId);
            return (
              <tr key={`${l.entryId}-${i}`}>
                <td dir="ltr">{l.date}</td>
                <td className="code" dir="ltr"><Link to={`/accounting/journal/${l.entryId}`}>{l.entryNumber}</Link></td>
                <td>{REF_AR[l.referenceType] ?? l.referenceType} {l.documentNumber && (link ? <Link to={link} dir="ltr" className="code">{l.documentNumber}</Link> : <span dir="ltr" className="code">{l.documentNumber}</span>)}</td>
                <td>{l.description}</td>
                <td className="num" dir="ltr">{l.debit === '0.00' ? '' : formatAmount(l.debit)}</td>
                <td className="num" dir="ltr">{l.credit === '0.00' ? '' : formatAmount(l.credit)}</td>
                <Amount v={l.balance} />
              </tr>
            );
          })}
          {d.lines.length === 0 && <tr><td colSpan={7} className="muted empty-row">لا توجد حركات في الفترة</td></tr>}
          <tr className="total-row emphasis"><td colSpan={4}>الإجمالي والرصيد الختامي</td><Amount v={d.totals.debit} /><Amount v={d.totals.credit} /><Amount v={d.closingBalance} strong /></tr>
        </tbody>
      </table>
    </div>
  );
}

export const ledgerCsv = (d: LedgerData) => [
  ['التاريخ', 'القيد', 'النوع', 'المستند', 'البيان', 'مدين', 'دائن', 'الرصيد'],
  ['', '', '', '', 'الرصيد الافتتاحي', '', '', d.openingBalance],
  ...d.lines.map((l) => [l.date, l.entryNumber, REF_AR[l.referenceType] ?? l.referenceType, l.documentNumber, l.description, l.debit, l.credit, l.balance]),
  ['', '', '', '', 'الإجمالي', d.totals.debit, d.totals.credit, d.closingBalance],
];

export function GeneralLedger() {
  const { params, range, setRange, set, query } = useUrlRange();
  const accountId = params.get('accountId') ?? '';
  const accounts = useLoad(() => api<{ data: { id: string; code: string; nameAr: string; isPostable: boolean }[] }>('GET', '/api/accounts?postableOnly=true'));
  const r = useLoad(() => (accountId
    ? api<LedgerData & { account: { code: string; nameAr: string }; truncated: boolean }>('GET', `/api/reports/general-ledger?${query}&accountId=${accountId}`)
    : Promise.resolve(null)), [accountId, query]);
  const d = r.data;
  return (
    <ReportFrame title={d ? `دفتر الأستاذ — ${d.account.code} ${d.account.nameAr}` : 'دفتر الأستاذ العام'} subtitle={`من ${range.from} إلى ${range.to}`} csv={d ? () => ledgerCsv(d) : undefined}>
      <div className="card">
        <RangeBar range={range} setRange={setRange}>
          <select aria-label="الحساب" value={accountId} onChange={(e) => set({ accountId: e.target.value })}>
            <option value="">— اختر الحساب —</option>
            {accounts.data?.data.map((a) => <option key={a.id} value={a.id}>{a.code} — {a.nameAr}</option>)}
          </select>
        </RangeBar>
        <ErrorBox error={r.error ?? accounts.error} />
        {!accountId && <p className="muted">اختر حسابًا لعرض حركاته.</p>}
        {d && <>
          {d.truncated && <div className="alert alert-warn">عُرضت أول 5000 حركة فقط؛ ضيّق الفترة.</div>}
          <LedgerTable d={d} />
        </>}
      </div>
    </ReportFrame>
  );
}

interface Entry { id: string; entryNumber: string; date: string; description: string; referenceType: string; status: string; totalDebit: string; lines: { lineNo: number; accountCode: string; accountName: string; description: string | null; debit: string; credit: string }[] }

export function JournalReport() {
  const { range, setRange, query } = useUrlRange();
  const [type, setType] = useState('');
  const [offset, setOffset] = useState(0);
  const PAGE = 100;
  const r = useLoad(() => api<{ data: Entry[]; total: number; totals: { debit: string; credit: string } }>(
    'GET', `/api/reports/journal?${query}&limit=${PAGE}&offset=${offset}${type ? `&referenceType=${type}` : ''}`), [query, type, offset]);
  const d = r.data;
  return (
    <ReportFrame title="تقرير اليومية" subtitle={`من ${range.from} إلى ${range.to}`} csv={d ? () => [
      ['القيد', 'التاريخ', 'النوع', 'البيان', 'الحساب', 'مدين', 'دائن'],
      ...d.data.flatMap((e) => e.lines.map((l) => [e.entryNumber, e.date, REF_AR[e.referenceType] ?? e.referenceType, l.description ?? e.description, `${l.accountCode} ${l.accountName}`, l.debit, l.credit])),
    ] : undefined}>
      <div className="card">
        <RangeBar range={range} setRange={(x) => { setOffset(0); setRange(x); }}>
          <select aria-label="نوع القيد" value={type} onChange={(e) => { setType(e.target.value); setOffset(0); }}>
            <option value="">كل الأنواع</option>
            {Object.entries(REF_AR).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
          </select>
        </RangeBar>
        <ErrorBox error={r.error} />
        {d && (
          <>
            <p className="muted small">{d.total} قيدًا — إجمالي المدين <span dir="ltr">{formatAmount(d.totals.debit)}</span> = إجمالي الدائن <span dir="ltr">{formatAmount(d.totals.credit)}</span></p>
            <div className="table-wrap">
              <table className="table report-table">
                <thead><tr><th>القيد / الحساب</th><th>البيان</th><th className="num">مدين</th><th className="num">دائن</th></tr></thead>
                <tbody>
                  {d.data.map((e) => [
                    <tr key={e.id} className="group-row">
                      <th><Link to={`/accounting/journal/${e.id}`} dir="ltr" className="code">{e.entryNumber}</Link> <span dir="ltr">{e.date}</span></th>
                      <th>{REF_AR[e.referenceType] ?? e.referenceType} — {e.description}{e.status === 'REVERSED' && <span className="muted small"> (معكوس)</span>}</th><th /><th />
                    </tr>,
                    ...e.lines.map((l) => (
                      <tr key={`${e.id}-${l.lineNo}`}>
                        <td><span className="code" dir="ltr">{l.accountCode}</span> {l.accountName}</td><td>{l.description ?? ''}</td>
                        <td className="num" dir="ltr">{l.debit === '0.00' ? '' : formatAmount(l.debit)}</td><td className="num" dir="ltr">{l.credit === '0.00' ? '' : formatAmount(l.credit)}</td>
                      </tr>
                    )),
                  ])}
                  {d.data.length === 0 && <tr><td colSpan={4} className="muted empty-row">لا توجد قيود في الفترة</td></tr>}
                </tbody>
              </table>
            </div>
            {d.total > PAGE && (
              <div className="pager no-print">
                <span className="muted small">{offset + 1}–{Math.min(offset + PAGE, d.total)} من {d.total}</span>
                <button className="btn" disabled={offset === 0} onClick={() => setOffset(offset - PAGE)}>السابق</button>
                <button className="btn" disabled={offset + PAGE >= d.total} onClick={() => setOffset(offset + PAGE)}>التالي</button>
              </div>
            )}
          </>
        )}
      </div>
    </ReportFrame>
  );
}
