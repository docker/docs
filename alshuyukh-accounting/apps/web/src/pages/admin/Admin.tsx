import { NavLink, Navigate, Route, Routes } from 'react-router-dom';
import { useAuth } from '../../auth';
import { PageHeader } from '../../ui';
import { Logs } from './Logs';
import { Overview } from './Overview';
import { Plans } from './Plans';
import { Revenue, Subscriptions, Usage } from './Billing';
import { TenantDetail, Tenants } from './Tenants';
import { Flags, Users } from './Users';

const TABS: [string, string][] = [
  ['overview', 'نظرة عامة'], ['tenants', 'المنشآت'], ['users', 'المستخدمون'], ['subscriptions', 'الاشتراكات'], ['plans', 'الباقات'],
  ['revenue', 'الإيرادات'], ['usage', 'الاستخدام'], ['logs', 'السجلات والأخطاء'], ['flags', 'خصائص التشغيل'],
];

/** Platform administration (the SaaS operator). Separate from tenant settings. */
export default function Admin() {
  const { me } = useAuth();
  if (!me?.user.isPlatformAdmin) return <Navigate to="/" replace />;
  return (
    <>
      <PageHeader title="إدارة المنصة" />
      <nav className="tabs tabs-scroll">
        {TABS.map(([to, label]) => <NavLink key={to} to={`/admin/${to}`} className={({ isActive }) => `tab ${isActive ? 'active' : ''}`}>{label}</NavLink>)}
      </nav>
      <Routes>
        <Route index element={<Navigate to="overview" replace />} />
        <Route path="overview" element={<Overview />} />
        <Route path="tenants" element={<Tenants />} />
        <Route path="tenants/:id" element={<TenantDetail />} />
        <Route path="users" element={<Users />} />
        <Route path="subscriptions" element={<Subscriptions />} />
        <Route path="plans" element={<Plans />} />
        <Route path="revenue" element={<Revenue />} />
        <Route path="usage" element={<Usage />} />
        <Route path="logs" element={<Logs />} />
        <Route path="flags" element={<Flags />} />
      </Routes>
    </>
  );
}
