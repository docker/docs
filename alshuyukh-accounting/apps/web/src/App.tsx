import { lazy } from 'react';
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
import Parties, { CUSTOMERS, SUPPLIERS } from './pages/parties/Parties';
import Products from './pages/products/Products';
import Section from './pages/documents/Section';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Users from './pages/Users';

// Less frequently used sections load on first visit, keeping the initial bundle small.
const Reports = lazy(() => import('./pages/reports/Reports'));
const EInvoicing = lazy(() => import('./pages/einvoicing/EInvoicing'));
const Admin = lazy(() => import('./pages/admin/Admin'));
const Subscription = lazy(() => import('./pages/subscription/Subscription'));
const Expenses = lazy(() => import('./pages/expenses/Expenses'));
const AuditLog = lazy(() => import('./pages/AuditLog'));
const Roles = lazy(() => import('./pages/Roles'));
const Companies = lazy(() => import('./pages/Companies'));
const Settings = lazy(() => import('./pages/Settings'));

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
        <Route path="reports/*" element={<Reports />} />
        <Route path="e-invoicing" element={<EInvoicing />} />
        <Route path="settings" element={<Settings />} />
        <Route path="settings/companies" element={<Companies />} />
        <Route path="settings/users" element={<Users />} />
        <Route path="settings/roles" element={<Roles />} />
        <Route path="settings/audit" element={<AuditLog />} />
        <Route path="settings/subscription" element={<Subscription />} />
        <Route path="admin/*" element={<Admin />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  );
}
