import { Navigate, Route, Routes } from 'react-router-dom';
import { useAuth } from './auth';
import Shell from './Shell';
import AuditLog from './pages/AuditLog';
import Companies from './pages/Companies';
import ComingSoon from './pages/ComingSoon';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Roles from './pages/Roles';
import Settings from './pages/Settings';
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
        <Route path="sales" element={<ComingSoon title="المبيعات" phase={4} />} />
        <Route path="purchases" element={<ComingSoon title="المشتريات" phase={4} />} />
        <Route path="inventory" element={<ComingSoon title="المخزون" phase={5} />} />
        <Route path="customers" element={<ComingSoon title="العملاء" phase={3} />} />
        <Route path="suppliers" element={<ComingSoon title="الموردون" phase={3} />} />
        <Route path="expenses" element={<ComingSoon title="المصروفات" phase={6} />} />
        <Route path="accounting" element={<ComingSoon title="المحاسبة" phase={2} />} />
        <Route path="reports" element={<ComingSoon title="التقارير" phase={7} />} />
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
