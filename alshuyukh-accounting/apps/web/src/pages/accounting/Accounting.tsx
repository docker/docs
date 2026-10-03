import { NavLink, Outlet } from 'react-router-dom';
import { useAuth } from '../../auth';
import { PageHeader } from '../../ui';

export default function Accounting() {
  const { can } = useAuth();
  const tabs = [
    { to: 'journal', label: 'القيود اليومية', show: can('journal.view') },
    { to: 'accounts', label: 'دليل الحسابات', show: can('account.view') },
    { to: 'fiscal', label: 'السنوات المالية', show: can('account.view') },
    { to: 'trial-balance', label: 'ميزان المراجعة', show: can('report.view') },
  ].filter((t) => t.show);
  return (
    <>
      <PageHeader title="المحاسبة" />
      <nav className="tabs" aria-label="أقسام المحاسبة">
        {tabs.map((t) => <NavLink key={t.to} to={t.to} className={({ isActive }) => `tab ${isActive ? 'active' : ''}`}>{t.label}</NavLink>)}
      </nav>
      <Outlet />
    </>
  );
}
