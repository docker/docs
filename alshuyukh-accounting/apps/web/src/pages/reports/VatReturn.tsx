import { useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { ReportFrame } from './shared';

interface Box { amount: string; adjustment: string; vat: string }
interface VatReturnData {
  dateFrom: string; dateTo: string;
  sales: Record<string, Box>; purchases: Record<string, Box>;
  '13_totalVatDue': string; '14_previousPeriodCorrections': string; '15_creditCarriedForward': string; '16_netVatDue': string;
  notSupported: string[];
  ledger: { vatOutput: string; vatInput: string; outputDifference: string; inputDifference: string; reconciled: boolean };
}
interface TaxTx { id: string; direction: string; sourceType: string; sourceId: string; sourceNumber: string | null; date: string; vatCategory: string; taxableAmount: string; taxAmount: string; isAdjustment: boolean; reversesId: string | null; partyName: string | null; partyVatNumber: string | null }

const LABELS: Record<string, string> = {
  '1_standardRated': 'المبيعات الخاضعة للنسبة الأساسية', '2_citizenHealthEducation': 'مبيعات المواطنين (خدمات صحية خاصة وتعليم أهلي)',
  '3_zeroRatedDomestic': 'المبيعات المحلية الخاضعة للنسبة الصفرية', '4_exports': 'الصادرات', '5_exempt': 'المبيعات المعفاة', '6_total': 'إجمالي المبيعات',
  '7_standardRatedDomestic': 'المشتريات الخاضعة للنسبة الأساسية', '8_importsPaidAtCustoms': 'الاستيرادات الخاضعة للضريبة المدفوعة في الجمارك',
  '9_importsReverseCharge': 'الاستيرادات الخاضعة للتحويل العكسي', '10_zeroRated': 'المشتريات الخاضعة للنسبة الصفرية', '11_exempt': 'المشتريات المعفاة', '12_total': 'إجمالي المشتريات',
};
const SOURCE_AR: Record<string, string> = {
  SALES_INVOICE: 'فاتورة مبيعات', SALES_RETURN: 'مرتجع مبيعات', PURCHASE_INVOICE: 'فاتورة مشتريات', PURCHASE_RETURN: 'مرتجع مشتريات', EXPENSE: 'مصروف',
};
const SOURCE_LINK: Record<string, string> = {
  SALES_INVOICE: '/sales/invoices', SALES_RETURN: '/sales/returns', PURCHASE_INVOICE: '/purchases/invoices', PURCHASE_RETURN: '/purchases/returns', EXPENSE: '/expenses/list',
};

/** First and last day of the previous calendar quarter, in Riyadh time. */
function lastQuarter() {
  const [y, m] = new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh' }).format(new Date()).split('-').map(Number) as [number, number];
  const q = Math.floor((m - 1) / 3);
  const [year, startMonth] = q === 0 ? [y - 1, 10] : [y, q * 3 - 2];
  const end = new Date(Date.UTC(year, startMonth + 2, 0)).toISOString().slice(0, 10);
  return { from: `${year}-${String(startMonth).padStart(2, '0')}-01`, to: end };
}

export default function VatReturn() {
  const [range, setRange] = useState(lastQuarter);
  const [showTx, setShowTx] = useState(false);
  const report = useLoad(() => api<VatReturnData>('GET', `/api/reports/vat-return?dateFrom=${range.from}&dateTo=${range.to}`), [range.from, range.to]);
  const tx = useLoad(() => (showTx ? api<{ data: TaxTx[] }>('GET', `/api/reports/vat-transactions?dateFrom=${range.from}&dateTo=${range.to}`) : Promise.resolve(null)), [showTx, range.from, range.to]);
  const r = report.data;

  const section = (title: string, boxes: Record<string, Box>) => (
    <>
      <tr className="group-row"><th colSpan={5}>{title}</th></tr>
      {Object.entries(boxes).map(([key, b]) => {
        const [n] = key.split('_');
        const unsupported = r?.notSupported.includes(key);
        const total = key.endsWith('_total');
        return (
          <tr key={key} className={total ? 'total-row' : unsupported ? 'muted' : ''}>
            <td className="code">{n}</td>
            <td>{LABELS[key]}{unsupported && <span className="muted small"> — غير مدعوم بعد</span>}</td>
            <td className="num" dir="ltr">{formatAmount(b.amount)}</td>
            <td className="num" dir="ltr">{formatAmount(b.adjustment)}</td>
            <td className="num" dir="ltr">{formatAmount(b.vat)}</td>
          </tr>
        );
      })}
    </>
  );

  return (
    <ReportFrame title="إقرار ضريبة القيمة المضافة" subtitle={`من ${range.from} إلى ${range.to}`}>
      <div className="card">
        <div className="toolbar no-print">
          <label>من<input type="date" value={range.from} onChange={(e) => setRange((x) => ({ ...x, from: e.target.value }))} /></label>
          <label>إلى<input type="date" value={range.to} onChange={(e) => setRange((x) => ({ ...x, to: e.target.value }))} /></label>
        </div>
        <ErrorBox error={report.error} />
        {r && (
          <>
            <div className={`alert ${r.ledger.reconciled ? 'alert-ok' : 'alert-warn'}`}>
              {r.ledger.reconciled
                ? 'الإقرار مطابق لأرصدة حسابات ضريبة المخرجات والمدخلات في دفتر الأستاذ.'
                : `يوجد فرق بين الإقرار ودفتر الأستاذ: المخرجات ${formatAmount(r.ledger.outputDifference)}، المدخلات ${formatAmount(r.ledger.inputDifference)}. تحقق من القيود اليدوية على حسابات الضريبة.`}
            </div>
            <div className="table-wrap">
              <table className="table">
                <thead><tr><th>البند</th><th>الوصف</th><th className="num">المبلغ (ريال)</th><th className="num">التعديلات</th><th className="num">مبلغ الضريبة</th></tr></thead>
                <tbody>
                  {section('ضريبة القيمة المضافة على المبيعات', r.sales)}
                  {section('ضريبة القيمة المضافة على المشتريات', r.purchases)}
                  <tr className="group-row"><th colSpan={5}>صافي الضريبة</th></tr>
                  <tr><td className="code">13</td><td>إجمالي ضريبة القيمة المضافة المستحقة عن الفترة الحالية</td><td /><td /><td className="num" dir="ltr">{formatAmount(r['13_totalVatDue'])}</td></tr>
                  <tr className="muted"><td className="code">14</td><td>تصحيحات من الفترات السابقة <span className="small">— غير مدعوم بعد</span></td><td /><td /><td className="num" dir="ltr">{formatAmount(r['14_previousPeriodCorrections'])}</td></tr>
                  <tr className="muted"><td className="code">15</td><td>ضريبة القيمة المضافة المرحلة من الفترات السابقة <span className="small">— غير مدعوم بعد</span></td><td /><td /><td className="num" dir="ltr">{formatAmount(r['15_creditCarriedForward'])}</td></tr>
                  <tr className="total-row"><td className="code">16</td><td>صافي الضريبة المستحقة {Number(r['16_netVatDue']) < 0 ? '(رصيد لصالح المنشأة)' : ''}</td><td /><td /><td className="num" dir="ltr">{formatAmount(r['16_netVatDue'])}</td></tr>
                </tbody>
              </table>
            </div>
            <dl className="meta">
              <div><dt>رصيد ضريبة المخرجات في الدفتر</dt><dd dir="ltr">{formatAmount(r.ledger.vatOutput)}</dd></div>
              <div><dt>رصيد ضريبة المدخلات في الدفتر</dt><dd dir="ltr">{formatAmount(r.ledger.vatInput)}</dd></div>
            </dl>
            <p className="muted small">هذا تقرير مساعد لإعداد الإقرار وليس تقديمًا إلى هيئة الزكاة والضريبة والجمارك. البنود الموسومة بـ«غير مدعوم بعد» تظهر صفرًا لأن النظام لا يفصلها حاليًا.</p>
            <button className="btn" onClick={() => setShowTx((s) => !s)}>{showTx ? 'إخفاء الحركات الضريبية' : 'عرض الحركات الضريبية'}</button>
          </>
        )}
      </div>
      {showTx && (
        <div className="card">
          <ErrorBox error={tx.error} />
          <div className="table-wrap">
            <table className="table">
              <thead><tr><th>التاريخ</th><th>النوع</th><th>المستند</th><th>الطرف</th><th>الفئة</th><th className="num">المبلغ الخاضع</th><th className="num">الضريبة</th></tr></thead>
              <tbody>
                {tx.data?.data.map((t) => (
                  <tr key={t.id}>
                    <td dir="ltr">{t.date}</td>
                    <td>{t.direction === 'OUTPUT' ? 'مخرجات' : 'مدخلات'} — {SOURCE_AR[t.sourceType] ?? t.sourceType}{t.reversesId ? ' (إلغاء)' : ''}</td>
                    <td dir="ltr" className="code"><Link to={`${SOURCE_LINK[t.sourceType]}/${t.sourceId}`}>{t.sourceNumber ?? '—'}</Link></td>
                    <td>{t.partyName ?? ''}{t.partyVatNumber && <div className="muted small" dir="ltr">{t.partyVatNumber}</div>}</td>
                    <td>{t.vatCategory}</td>
                    <td className="num" dir="ltr">{formatAmount(t.taxableAmount)}</td>
                    <td className="num" dir="ltr">{formatAmount(t.taxAmount)}</td>
                  </tr>
                ))}
                {tx.data?.data.length === 0 && <tr><td colSpan={7} className="muted empty-row">لا توجد حركات في الفترة</td></tr>}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </ReportFrame>
  );
}
