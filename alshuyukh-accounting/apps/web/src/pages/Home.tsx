import { useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api';
import { useAuth } from '../auth';
import { MonthlyBars, MonthlyLine } from '../components/Charts';
import { formatAmount } from '../money';
import { ErrorBox, PageHeader, useLoad } from '../ui';

interface Company { id: string; name: string; vatNumber: string | null; commercialRegistration: string | null; city: string | null; currency: string }

export default function Home() {
  const { me, can } = useAuth();
  const companies = useLoad(() => (can('company.view') ? api<{ data: Company[] }>('GET', '/api/companies') : Promise.resolve({ data: [] })));
  const company = companies.data?.data[0];

  const steps = [
    { done: !!company?.vatNumber, label: 'أضف الرقم الضريبي للشركة', to: '/settings/companies' },
    { done: !!company?.commercialRegistration, label: 'أضف رقم السجل التجاري', to: '/settings/companies' },
    { done: !!company?.city, label: 'أكمل عنوان الشركة', to: '/settings/companies' },
  ];

  return (
    <>
      <PageHeader title={`مرحبًا، ${me?.user.fullName ?? ''}`} />
      {can('financial_report.view') && <Dashboard />}
      <div className="grid-cards">
        {company && can('company.manage') && (
          <div className="card">
            <h2>إعداد الشركة</h2>
            <ul className="checklist">
              {steps.map((s) => (
                <li key={s.label} className={s.done ? 'done' : ''}>
                  <span aria-hidden>{s.done ? '✓' : '○'}</span>
                  {s.done ? s.label : <Link to={s.to}>{s.label}</Link>}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </>
  );
}

interface DashboardData {
  dateFrom: string; dateTo: string;
  kpis: Record<'totalSales' | 'totalPurchases' | 'costOfSales' | 'expenses' | 'netProfit' | 'receivables' | 'payables' | 'cash' | 'bank' | 'inventoryValue' | 'vatPayable', string>;
  monthly: { month: string; revenue: string; costs: string; netProfit: string; cashBalance: string }[];
}

const KPIS: { key: keyof DashboardData['kpis']; label: string; to: string; period?: boolean }[] = [
  { key: 'totalSales', label: 'المبيعات', to: '/reports/sales', period: true },
  { key: 'totalPurchases', label: 'المشتريات', to: '/reports/purchases', period: true },
  { key: 'expenses', label: 'المصروفات', to: '/reports/expenses', period: true },
  { key: 'netProfit', label: 'صافي الربح', to: '/reports/profit-loss', period: true },
  { key: 'receivables', label: 'الذمم المدينة', to: '/reports/receivables' },
  { key: 'payables', label: 'الذمم الدائنة', to: '/reports/payables' },
  { key: 'cash', label: 'النقدية', to: '/reports/cash-flow' },
  { key: 'bank', label: 'البنك', to: '/reports/cash-flow' },
  { key: 'inventoryValue', label: 'قيمة المخزون', to: '/inventory?tab=valuation' },
  { key: 'vatPayable', label: 'ضريبة القيمة المضافة المستحقة', to: '/reports/vat' },
];

function Dashboard() {
  const d = useLoad(() => api<DashboardData>('GET', '/api/dashboard'));
  const [table, setTable] = useState(false);
  if (d.error) return <ErrorBox error={d.error} />;
  if (!d.data) return <div className="card muted">جارٍ تحميل المؤشرات…</div>;
  const { kpis, monthly, dateFrom, dateTo } = d.data;
  return (
    <>
      <p className="muted small">المؤشرات من القيود المرحّلة فقط — الفترة <span dir="ltr">{dateFrom}</span> إلى <span dir="ltr">{dateTo}</span>؛ الأرصدة في نهاية الفترة.</p>
      <div className="kpi-grid">
        {KPIS.map((k) => (
          <Link key={k.key} to={k.to} className={`card kpi ${k.key === 'netProfit' ? 'kpi-main' : ''}`}>
            <span className="muted small">{k.label}{k.period ? '' : ' (رصيد)'}</span>
            <strong dir="ltr" className={kpis[k.key].startsWith('-') ? 'negative' : ''}>{formatAmount(kpis[k.key])}</strong>
          </Link>
        ))}
      </div>
      <div className="grid-cards charts">
        <div className="card">
          <div className="card-head"><h2>الإيرادات والتكاليف شهريًا</h2><button className="btn btn-small" onClick={() => setTable((v) => !v)}>{table ? 'عرض الرسم' : 'عرض كجدول'}</button></div>
          {table ? (
            <div className="table-wrap">
              <table className="table">
                <thead><tr><th>الشهر</th><th className="num">الإيرادات</th><th className="num">التكاليف والمصروفات</th><th className="num">صافي الربح</th><th className="num">النقدية آخر الشهر</th></tr></thead>
                <tbody>{monthly.map((m) => (
                  <tr key={m.month}><td dir="ltr">{m.month}</td><td className="num" dir="ltr">{formatAmount(m.revenue)}</td><td className="num" dir="ltr">{formatAmount(m.costs)}</td>
                    <td className="num" dir="ltr">{formatAmount(m.netProfit)}</td><td className="num" dir="ltr">{formatAmount(m.cashBalance)}</td></tr>
                ))}</tbody>
              </table>
            </div>
          ) : (
            <MonthlyBars data={monthly} series={[
              { key: 'revenue', label: 'الإيرادات', color: 'var(--series-1)' },
              { key: 'costs', label: 'التكاليف والمصروفات', color: 'var(--series-2)' },
            ]} />
          )}
        </div>
        <div className="card">
          <h2>رصيد النقدية والبنوك</h2>
          <MonthlyLine data={monthly} field="cashBalance" label="النقدية والبنوك آخر الشهر" color="var(--series-1)" />
        </div>
      </div>
    </>
  );
}
