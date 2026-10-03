import { useState, type FormEvent } from 'react';
import { api } from '../../api';
import { ACCOUNT_TYPE_AR, formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';

interface Row { accountId: string; code: string; nameAr: string; type: string; openingDebit: string; openingCredit: string; periodDebit: string; periodCredit: string; closingDebit: string; closingCredit: string }
interface Result { data: Row[]; totals: Omit<Row, 'accountId' | 'code' | 'nameAr' | 'type'>; balanced: boolean }

const year = new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh', year: 'numeric' }).format(new Date());

export default function TrialBalance() {
  const [range, setRange] = useState({ from: `${year}-01-01`, to: `${year}-12-31` });
  const tb = useLoad(() => api<Result>('GET', `/api/reports/trial-balance?dateFrom=${range.from}&dateTo=${range.to}`), [range.from, range.to]);

  function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    setRange({ from: String(f.get('from')), to: String(f.get('to')) });
  }

  const cell = (v: string) => <td className="num" dir="ltr">{v === '0.00' ? '' : formatAmount(v)}</td>;
  return (
    <div className="card">
      <form className="toolbar" onSubmit={submit}>
        <label className="inline">من<input name="from" type="date" defaultValue={range.from} required /></label>
        <label className="inline">إلى<input name="to" type="date" defaultValue={range.to} required /></label>
        <button className="btn">عرض</button>
        <span className="spacer" />
        {tb.data && <span className={`tag ${tb.data.balanced ? 'tag-ok' : 'tag-off'}`}>{tb.data.balanced ? 'الميزان متوازن' : 'الميزان غير متوازن'}</span>}
      </form>
      <ErrorBox error={tb.error} />
      <div className="table-wrap">
        <table className="table tb-table">
          <thead>
            <tr><th rowSpan={2}>الرمز</th><th rowSpan={2}>الحساب</th><th colSpan={2}>الرصيد الافتتاحي</th><th colSpan={2}>حركة الفترة</th><th colSpan={2}>الرصيد الختامي</th></tr>
            <tr><th className="num">مدين</th><th className="num">دائن</th><th className="num">مدين</th><th className="num">دائن</th><th className="num">مدين</th><th className="num">دائن</th></tr>
          </thead>
          <tbody>
            {tb.data?.data.map((r) => (
              <tr key={r.accountId}>
                <td dir="ltr" className="code">{r.code}</td>
                <td>{r.nameAr} <span className="muted small">{ACCOUNT_TYPE_AR[r.type]}</span></td>
                {cell(r.openingDebit)}{cell(r.openingCredit)}{cell(r.periodDebit)}{cell(r.periodCredit)}{cell(r.closingDebit)}{cell(r.closingCredit)}
              </tr>
            ))}
            {tb.data?.data.length === 0 && <tr><td colSpan={8} className="muted empty-row">لا توجد حركات في هذه الفترة</td></tr>}
          </tbody>
          {tb.data && (
            <tfoot>
              <tr><td colSpan={2}><strong>الإجمالي</strong></td>
                {(['openingDebit', 'openingCredit', 'periodDebit', 'periodCredit', 'closingDebit', 'closingCredit'] as const).map((k) =>
                  <td key={k} className="num" dir="ltr"><strong>{formatAmount(tb.data!.totals[k])}</strong></td>)}
              </tr>
            </tfoot>
          )}
        </table>
      </div>
    </div>
  );
}
