import { Link, Navigate, Route, Routes } from 'react-router-dom';
import { useAuth } from '../../auth';
import { PageHeader } from '../../ui';
import { ExpenseReport, SalesPurchases } from './Commercial';
import { BalanceSheet, CashFlow, ProfitLoss } from './Financial';
import { GeneralLedger, JournalReport } from './Ledger';
import { Aging, Statement } from './Parties';
import VatReturn from './VatReturn';

interface Item { to: string; title: string; hint: string; perm: string }
const GROUPS: { title: string; items: Item[] }[] = [
  { title: 'القوائم المالية', items: [
    { to: 'profit-loss', title: 'قائمة الدخل', hint: 'الإيرادات والتكاليف وصافي الربح', perm: 'financial_report.view' },
    { to: 'balance-sheet', title: 'الميزانية العمومية', hint: 'الأصول والالتزامات وحقوق الملكية في تاريخ', perm: 'financial_report.view' },
    { to: 'cash-flow', title: 'التدفقات النقدية', hint: 'مصادر النقد واستخداماته', perm: 'financial_report.view' },
    { to: '/accounting/trial-balance', title: 'ميزان المراجعة', hint: 'أرصدة كل الحسابات', perm: 'financial_report.view' },
  ] },
  { title: 'الدفاتر', items: [
    { to: 'general-ledger', title: 'دفتر الأستاذ العام', hint: 'حركات حساب برصيد متحرك', perm: 'financial_report.view' },
    { to: 'journal', title: 'تقرير اليومية', hint: 'القيود المرحّلة بسطورها', perm: 'financial_report.view' },
    { to: 'expenses', title: 'تقرير المصروفات', hint: 'حسب الحساب أو مركز التكلفة أو الشهر', perm: 'financial_report.view' },
  ] },
  { title: 'العملاء والموردون', items: [
    { to: 'receivables', title: 'أعمار الذمم المدينة', hint: 'المستحق على العملاء حسب التأخير', perm: 'report.view' },
    { to: 'payables', title: 'أعمار الذمم الدائنة', hint: 'المستحق للموردين حسب التأخير', perm: 'report.view' },
    { to: 'customer-statement', title: 'كشف حساب عميل', hint: 'الحركات والرصيد المتحرك', perm: 'report.view' },
    { to: 'supplier-statement', title: 'كشف حساب مورد', hint: 'الحركات والرصيد المتحرك', perm: 'report.view' },
  ] },
  { title: 'المبيعات والمشتريات', items: [
    { to: 'sales', title: 'تقرير المبيعات', hint: 'حسب العميل أو الصنف أو الشهر، مع مجمل الربح', perm: 'report.view' },
    { to: 'purchases', title: 'تقرير المشتريات', hint: 'حسب المورد أو الصنف أو الشهر', perm: 'report.view' },
  ] },
  { title: 'الضريبة والمخزون', items: [
    { to: 'vat', title: 'إقرار ضريبة القيمة المضافة', hint: 'بنود نموذج الهيئة مع مطابقة الدفتر', perm: 'financial_report.view' },
    { to: '/inventory?tab=valuation', title: 'تقييم المخزون', hint: 'الكمية والتكلفة مع مطابقة الدفتر', perm: 'inventory.view' },
    { to: '/inventory?tab=card', title: 'حركة الصنف', hint: 'كل حركات صنف برصيد متحرك', perm: 'inventory.view' },
  ] },
];

function Index() {
  const { can } = useAuth();
  const groups = GROUPS.map((g) => ({ ...g, items: g.items.filter((i) => can(i.perm)) })).filter((g) => g.items.length);
  return (
    <>
      <PageHeader title="التقارير" />
      {groups.length === 0 && <div className="card muted">لا تملك صلاحية عرض أي تقرير.</div>}
      {groups.map((g) => (
        <section key={g.title} className="report-group">
          <h2>{g.title}</h2>
          <div className="report-grid">
            {g.items.map((i) => (
              <Link key={i.to} to={i.to.startsWith('/') ? i.to : `/reports/${i.to}`} className="card link-card report-card">
                <strong>{i.title}</strong><span className="muted small">{i.hint}</span>
              </Link>
            ))}
          </div>
        </section>
      ))}
    </>
  );
}

export default function Reports() {
  return (
    <Routes>
      <Route index element={<Index />} />
      <Route path="profit-loss" element={<ProfitLoss />} />
      <Route path="balance-sheet" element={<BalanceSheet />} />
      <Route path="cash-flow" element={<CashFlow />} />
      <Route path="general-ledger" element={<GeneralLedger />} />
      <Route path="journal" element={<JournalReport />} />
      <Route path="expenses" element={<ExpenseReport />} />
      <Route path="receivables" element={<Aging key="c" party="customer" />} />
      <Route path="payables" element={<Aging key="s" party="supplier" />} />
      <Route path="customer-statement" element={<Statement key="c" party="customer" />} />
      <Route path="supplier-statement" element={<Statement key="s" party="supplier" />} />
      <Route path="sales" element={<SalesPurchases key="sales" kind="sales" />} />
      <Route path="purchases" element={<SalesPurchases key="purchases" kind="purchases" />} />
      <Route path="vat" element={<VatReturn />} />
      <Route path="*" element={<Navigate to="/reports" replace />} />
    </Routes>
  );
}
