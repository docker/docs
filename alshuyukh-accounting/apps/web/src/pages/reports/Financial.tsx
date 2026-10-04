import { Fragment, useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { ErrorBox, useLoad } from '../../ui';
import { Amount, RangeBar, ReportFrame, today, useRange } from './shared';

interface AccountRow { accountId: string; code: string; nameAr: string; groupName: string | null; amount: string }
interface Section { type: string; title: string; accounts: AccountRow[]; total: string }

const glLink = (accountId: string, query: string) => `/reports/general-ledger?accountId=${accountId}&${query}`;

function SectionRows({ s, query, emphasis }: { s: Section; query: string; emphasis?: boolean }) {
  return (
    <>
      <tr className="group-row"><th colSpan={3}>{s.title}</th></tr>
      {s.accounts.map((a) => (
        <tr key={a.accountId}>
          <td className="code" dir="ltr">{a.code}</td>
          <td><Link className="plain-link" to={glLink(a.accountId, query)}>{a.nameAr}</Link></td>
          <Amount v={a.amount} />
        </tr>
      ))}
      {s.accounts.length === 0 && <tr><td /><td className="muted">لا توجد حركات</td><td /></tr>}
      <tr className={`total-row ${emphasis ? 'emphasis' : ''}`}><td /><td>إجمالي {s.title}</td><Amount v={s.total} strong={emphasis} /></tr>
    </>
  );
}

export function ProfitLoss() {
  const { range, setRange, query } = useRange('year');
  const [costCenterId, setCostCenterId] = useState('');
  const centers = useLoad(() => api<{ data: { id: string; code: string; name: string }[] }>('GET', '/api/cost-centers').catch(() => ({ data: [] })));
  const r = useLoad(() => api<{ sections: Section[]; revenue: string; costOfSales: string; grossProfit: string; expenses: string; netProfit: string }>(
    'GET', `/api/reports/profit-loss?${query}${costCenterId ? `&costCenterId=${costCenterId}` : ''}`), [query, costCenterId]);
  const d = r.data;
  return (
    <ReportFrame title="قائمة الدخل" subtitle={`من ${range.from} إلى ${range.to}`} csv={d ? () => [
      ['الرمز', 'الحساب', 'المبلغ'],
      ...d.sections.flatMap((s) => [[s.title], ...s.accounts.map((a) => [a.code, a.nameAr, a.amount]), ['', `إجمالي ${s.title}`, s.total]]),
      ['', 'مجمل الربح', d.grossProfit], ['', 'صافي الربح', d.netProfit],
    ] : undefined}>
      <div className="card">
        <RangeBar range={range} setRange={setRange}>
          {(centers.data?.data.length ?? 0) > 0 && (
            <select aria-label="مركز التكلفة" value={costCenterId} onChange={(e) => setCostCenterId(e.target.value)}>
              <option value="">كل مراكز التكلفة</option>
              {centers.data!.data.map((c) => <option key={c.id} value={c.id}>{c.code} — {c.name}</option>)}
            </select>
          )}
        </RangeBar>
        <ErrorBox error={r.error} />
        {d && (
          <div className="table-wrap">
            <table className="table report-table">
              <tbody>
                <SectionRows s={d.sections[0]!} query={query} />
                <SectionRows s={d.sections[1]!} query={query} />
                <tr className="total-row emphasis"><td /><td>مجمل الربح</td><Amount v={d.grossProfit} /></tr>
                <SectionRows s={d.sections[2]!} query={query} />
                <tr className="total-row emphasis"><td /><td>{d.netProfit.startsWith('-') ? 'صافي الخسارة' : 'صافي الربح'}</td><Amount v={d.netProfit} strong /></tr>
              </tbody>
            </table>
          </div>
        )}
        <p className="muted small">من القيود المرحّلة فقط، دون قيود إقفال السنة.</p>
      </div>
    </ReportFrame>
  );
}

export function BalanceSheet() {
  const [asOf, setAsOf] = useState(today);
  const r = useLoad(() => api<{
    assets: Section; liabilities: Section; equity: Section & { currentEarnings: string };
    totalAssets: string; totalLiabilitiesAndEquity: string; balanced: boolean;
  }>('GET', `/api/reports/balance-sheet?asOf=${asOf}`), [asOf]);
  const d = r.data;
  const query = `dateFrom=${asOf.slice(0, 4)}-01-01&dateTo=${asOf}`;
  return (
    <ReportFrame title="الميزانية العمومية" subtitle={`كما في ${asOf}`} csv={d ? () => [
      ['الرمز', 'الحساب', 'المبلغ'],
      ...[d.assets, d.liabilities, d.equity].flatMap((s) => [[s.title], ...s.accounts.map((a) => [a.code, a.nameAr, a.amount]), ['', `إجمالي ${s.title}`, s.total]]),
      ['', 'أرباح الفترة غير المقفلة', d.equity.currentEarnings], ['', 'إجمالي الالتزامات وحقوق الملكية', d.totalLiabilitiesAndEquity],
    ] : undefined}>
      <div className="card">
        <div className="toolbar no-print"><label>كما في<input type="date" value={asOf} onChange={(e) => setAsOf(e.target.value)} /></label></div>
        <ErrorBox error={r.error} />
        {d && (
          <>
            {!d.balanced && <div className="alert alert-error">الميزانية غير متوازنة — راجع القيود.</div>}
            <div className="table-wrap">
              <table className="table report-table">
                <tbody>
                  <SectionRows s={d.assets} query={query} emphasis />
                  <SectionRows s={d.liabilities} query={query} />
                  <tr className="group-row"><th colSpan={3}>{d.equity.title}</th></tr>
                  {d.equity.accounts.map((a) => (
                    <tr key={a.accountId}><td className="code" dir="ltr">{a.code}</td><td><Link className="plain-link" to={glLink(a.accountId, query)}>{a.nameAr}</Link></td><Amount v={a.amount} /></tr>
                  ))}
                  <tr><td /><td><Link className="plain-link" to="/reports/profit-loss">أرباح الفترة غير المقفلة</Link></td><Amount v={d.equity.currentEarnings} /></tr>
                  <tr className="total-row"><td /><td>إجمالي حقوق الملكية</td><Amount v={d.equity.total} /></tr>
                  <tr className="total-row emphasis"><td /><td>إجمالي الالتزامات وحقوق الملكية</td><Amount v={d.totalLiabilitiesAndEquity} strong /></tr>
                </tbody>
              </table>
            </div>
            <p className="muted small">رصيد ضريبة المدخلات يظهر سالبًا ضمن الالتزامات لأنه يُخصم من ضريبة المخرجات.</p>
          </>
        )}
      </div>
    </ReportFrame>
  );
}

const CF_AR: Record<string, string> = { OPERATING: 'الأنشطة التشغيلية', INVESTING: 'الأنشطة الاستثمارية', FINANCING: 'الأنشطة التمويلية' };

export function CashFlow() {
  const { range, setRange, query } = useRange('year');
  const r = useLoad(() => api<{
    sections: { category: string; total: string; lines: { accountId: string; code: string; nameAr: string; inflow: string; outflow: string; net: string }[] }[];
    openingCash: string; netChange: string; closingCash: string; reconciled: boolean;
  }>('GET', `/api/reports/cash-flow?${query}`), [query]);
  const d = r.data;
  return (
    <ReportFrame title="قائمة التدفقات النقدية" subtitle={`من ${range.from} إلى ${range.to}`} csv={d ? () => [
      ['النشاط', 'الرمز', 'الحساب', 'داخل', 'خارج', 'الصافي'],
      ...d.sections.flatMap((s) => s.lines.map((l) => [CF_AR[s.category], l.code, l.nameAr, l.inflow, l.outflow, l.net])),
      ['', '', 'النقدية أول الفترة', '', '', d.openingCash], ['', '', 'صافي التغير', '', '', d.netChange], ['', '', 'النقدية آخر الفترة', '', '', d.closingCash],
    ] : undefined}>
      <div className="card">
        <RangeBar range={range} setRange={setRange} />
        <ErrorBox error={r.error} />
        {d && (
          <div className="table-wrap">
            <table className="table report-table">
              <thead><tr><th>الحساب المقابل</th><th className="num">داخل</th><th className="num">خارج</th><th className="num">الصافي</th></tr></thead>
              <tbody>
                {d.sections.map((s) => (
                  <Fragment key={s.category}>
                    <tr className="group-row"><th colSpan={4}>{CF_AR[s.category]}</th></tr>
                    {s.lines.map((l) => (
                      <tr key={l.accountId}><td><span className="code" dir="ltr">{l.code}</span> {l.nameAr}</td><Amount v={l.inflow} /><Amount v={l.outflow} /><Amount v={l.net} /></tr>
                    ))}
                    <tr className="total-row"><td>صافي {CF_AR[s.category]}</td><td /><td /><Amount v={s.total} /></tr>
                  </Fragment>
                ))}
                <tr className="total-row"><td>النقدية وما في حكمها أول الفترة</td><td /><td /><Amount v={d.openingCash} /></tr>
                <tr className="total-row"><td>صافي التغير في النقدية</td><td /><td /><Amount v={d.netChange} /></tr>
                <tr className="total-row emphasis"><td>النقدية وما في حكمها آخر الفترة</td><td /><td /><Amount v={d.closingCash} strong /></tr>
              </tbody>
            </table>
          </div>
        )}
        <p className="muted small">الطريقة المباشرة: كل حركة على الصندوق أو البنك أو حسابات طرق الدفع تُنسب إلى الحسابات المقابلة في القيد.</p>
      </div>
    </ReportFrame>
  );
}
