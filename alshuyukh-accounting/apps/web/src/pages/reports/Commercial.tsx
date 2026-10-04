import { useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { ErrorBox, useLoad } from '../../ui';
import { Amount, RangeBar, ReportFrame, useRange } from './shared';

interface Row { key: string; label: string; code: string | null; quantity: string; netAmount: string; vatAmount: string; totalAmount: string; cost: string; grossProfit?: string; documents: number; returns: number }

export function SalesPurchases({ kind }: { kind: 'sales' | 'purchases' }) {
  const { range, setRange, query } = useRange('year');
  const isSales = kind === 'sales';
  const [groupBy, setGroupBy] = useState('party');
  const r = useLoad(() => api<{ data: Row[]; totals: Record<string, string> }>('GET', `/api/reports/${kind}?${query}&groupBy=${groupBy}`), [kind, query, groupBy]);
  const d = r.data;
  const groups: [string, string][] = [['party', isSales ? 'حسب العميل' : 'حسب المورد'], ['product', 'حسب الصنف'], ['month', 'حسب الشهر'], ['document', 'حسب المستند']];
  const showQty = groupBy === 'product';
  const title = isSales ? 'تقرير المبيعات' : 'تقرير المشتريات';
  return (
    <ReportFrame title={title} subtitle={`من ${range.from} إلى ${range.to}`} csv={d ? () => [
      ['البند', 'الرمز', 'الكمية', 'الصافي', 'الضريبة', 'الإجمالي', ...(isSales ? ['التكلفة', 'مجمل الربح'] : []), 'المستندات', 'المرتجعات'],
      ...d.data.map((x) => [x.label, x.code, x.quantity, x.netAmount, x.vatAmount, x.totalAmount, ...(isSales ? [x.cost, x.grossProfit] : []), x.documents, x.returns]),
      ['الإجمالي', '', '', d.totals.netAmount, d.totals.vatAmount, d.totals.totalAmount, ...(isSales ? [d.totals.cost, d.totals.grossProfit] : [])],
    ] : undefined}>
      <div className="card">
        <RangeBar range={range} setRange={setRange}>
          <select aria-label="التجميع" value={groupBy} onChange={(e) => setGroupBy(e.target.value)}>
            {groups.map(([k, l]) => <option key={k} value={k}>{l}</option>)}
          </select>
        </RangeBar>
        <ErrorBox error={r.error} />
        {d && (
          <div className="table-wrap">
            <table className="table report-table">
              <thead>
                <tr>
                  <th>{groups.find(([k]) => k === groupBy)![1].replace('حسب ', '')}</th>
                  {showQty && <th className="num">الكمية</th>}
                  <th className="num">الصافي</th><th className="num">الضريبة</th><th className="num">الإجمالي</th>
                  {isSales && <><th className="num">التكلفة</th><th className="num">مجمل الربح</th></>}
                  <th className="num">مستندات</th><th className="num">مرتجعات</th>
                </tr>
              </thead>
              <tbody>
                {d.data.map((x) => (
                  <tr key={x.key}>
                    <td>{x.code && groupBy !== 'document' && <span className="code" dir="ltr">{x.code} </span>}<span dir={groupBy === 'month' || groupBy === 'document' ? 'ltr' : undefined}>{x.label}</span>{groupBy === 'document' && <span className="muted small" dir="ltr"> {x.code}</span>}</td>
                    {showQty && <td className="num" dir="ltr">{x.quantity.replace(/\.?0+$/, '')}</td>}
                    <Amount v={x.netAmount} /><Amount v={x.vatAmount} /><Amount v={x.totalAmount} />
                    {isSales && <><Amount v={x.cost} /><Amount v={x.grossProfit} /></>}
                    <td className="num">{x.documents}</td><td className="num">{x.returns || ''}</td>
                  </tr>
                ))}
                {d.data.length === 0 && <tr><td colSpan={9} className="muted empty-row">لا توجد مستندات في الفترة</td></tr>}
                <tr className="total-row emphasis">
                  <td>الإجمالي</td>{showQty && <td />}
                  <Amount v={d.totals.netAmount} /><Amount v={d.totals.vatAmount} /><Amount v={d.totals.totalAmount} strong />
                  {isSales && <><Amount v={d.totals.cost} /><Amount v={d.totals.grossProfit} /></>}
                  <td /><td />
                </tr>
              </tbody>
            </table>
          </div>
        )}
        <p className="muted small">
          من المستندات المُصدرة؛ المرتجعات تُطرح والملغاة مستبعدة.
          {isSales ? ' التكلفة من حركات المخزون الفعلية.' : ''} الأرقام المالية الرسمية في <Link to="/reports/profit-loss">قائمة الدخل</Link>.
        </p>
      </div>
    </ReportFrame>
  );
}

export function ExpenseReport() {
  const { range, setRange, query } = useRange('year');
  const [groupBy, setGroupBy] = useState('account');
  const r = useLoad(() => api<{ data: { key: string; label: string; code: string | null; amount: string; entries: number }[]; total: string }>(
    'GET', `/api/reports/expenses?${query}&groupBy=${groupBy}`), [query, groupBy]);
  const d = r.data;
  return (
    <ReportFrame title="تقرير المصروفات" subtitle={`من ${range.from} إلى ${range.to}`} csv={d ? () => [
      ['البند', 'الرمز', 'المبلغ', 'القيود'], ...d.data.map((x) => [x.label, x.code, x.amount, x.entries]), ['الإجمالي', '', d.total],
    ] : undefined}>
      <div className="card">
        <RangeBar range={range} setRange={setRange}>
          <select aria-label="التجميع" value={groupBy} onChange={(e) => setGroupBy(e.target.value)}>
            <option value="account">حسب الحساب</option><option value="costCenter">حسب مركز التكلفة</option><option value="month">حسب الشهر</option>
          </select>
        </RangeBar>
        <ErrorBox error={r.error} />
        {d && (
          <div className="table-wrap">
            <table className="table report-table">
              <thead><tr><th>البند</th><th className="num">المبلغ</th><th className="num">النسبة</th><th className="num">القيود</th></tr></thead>
              <tbody>
                {d.data.map((x) => (
                  <tr key={x.key}>
                    <td>{x.code && <span className="code" dir="ltr">{x.code} </span>}<span dir={groupBy === 'month' ? 'ltr' : undefined}>{x.label}</span></td>
                    <Amount v={x.amount} />
                    <td className="num" dir="ltr">{d.total !== '0.00' ? `${((Number(x.amount) / Number(d.total)) * 100).toFixed(1)}%` : ''}</td>
                    <td className="num">{x.entries}</td>
                  </tr>
                ))}
                {d.data.length === 0 && <tr><td colSpan={4} className="muted empty-row">لا توجد مصروفات في الفترة</td></tr>}
                <tr className="total-row emphasis"><td>الإجمالي</td><Amount v={d.total} strong /><td /><td /></tr>
              </tbody>
            </table>
          </div>
        )}
        <p className="muted small">من حسابات المصروفات في دفتر الأستاذ: سندات المصروفات، وبنود المصروفات في فواتير المشتريات، والقيود اليدوية.</p>
      </div>
    </ReportFrame>
  );
}
