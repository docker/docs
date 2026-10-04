import { Navigate, Route, Routes } from 'react-router-dom';
import { useAuth } from './auth';
import Shell from './Shell';
import Accounting from './pages/accounting/Accounting';
import ChartOfAccounts from './pages/accounting/ChartOfAccounts';
import FiscalYears from './pages/accounting/FiscalYears';
import JournalDetail from './pages/accounting/JournalDetail';
import JournalForm from './pages/accounting/JournalForm';
import JournalList from './pages/accounting/JournalList';
import TrialBalance from './pages/accounting/TrialBalance';
import AuditLog from './pages/AuditLog';
import Parties, { CUSTOMERS, SUPPLIERS } from './pages/parties/Parties';
import Products from './pages/products/Products';
import Section from './pages/documents/Section';
import Companies from './pages/Companies';
import ComingSoon from './pages/ComingSoon';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Roles from './pages/Roles';
import Settings from './pages/Settings';
import Expenses from './pages/expenses/Expenses';
import VatReturn from './pages/reports/VatReturn';
import Users from './pages/Users';

export default function App() {
  const { me, loading } = useAuth();
  if (loading) return <div className="center-screen">جارٍ التحميل…</div>;

  if (!me) {
    return (
      <Routes>
        <Route path="/register" element={<Register />} />
        <Route path="*" element={<Login />} />
      </Routes>
    );
  }

  return (
    <Routes>
      <Route element={<Shell />}>
        <Route index element={<Home />} />
        <Route path="sales/*" element={<Section key="sales" party="customer" />} />
        <Route path="purchases/*" element={<Section key="purchases" party="supplier" />} />
        <Route path="inventory" element={<Products />} />
        <Route path="customers" element={<Parties key="customers" config={CUSTOMERS} />} />
        <Route path="suppliers" element={<Parties key="suppliers" config={SUPPLIERS} />} />
        <Route path="expenses/*" element={<Expenses />} />
        <Route path="accounting" element={<Accounting />}>
          <Route index element={<Navigate to="journal" replace />} />
          <Route path="journal" element={<JournalList />} />
          <Route path="journal/new" element={<JournalForm />} />
          <Route path="journal/:id" element={<JournalDetail />} />
          <Route path="accounts" element={<ChartOfAccounts />} />
          <Route path="fiscal" element={<FiscalYears />} />
          <Route path="trial-balance" element={<TrialBalance />} />
        </Route>
        <Route path="reports" element={<Navigate to="vat" replace />} />
        <Route path="reports/vat" element={<VatReturn />} />
        <Route path="e-invoicing" element={<ComingSoon title="الفوترة الإلكترونية" phase={8} />} />
        <Route path="settings" element={<Settings />} />
        <Route path="settings/companies" element={<Companies />} />
        <Route path="settings/users" element={<Users />} />
        <Route path="settings/roles" element={<Roles />} />
        <Route path="settings/audit" element={<AuditLog />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  );
}
