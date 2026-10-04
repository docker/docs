import { useState } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { useAuth } from './auth';

const NAV: { to: string; label: string; icon: string; phase?: number }[] = [
  { to: '/', label: 'الرئيسية', icon: '⌂' },
  { to: '/sales', label: 'المبيعات', icon: '↗' },
  { to: '/purchases', label: 'المشتريات', icon: '↙' },
  { to: '/inventory', label: 'المخزون', icon: '▦' },
  { to: '/customers', label: 'العملاء', icon: '☺' },
  { to: '/suppliers', label: 'الموردون', icon: '⚑' },
  { to: '/expenses', label: 'المصروفات', icon: '−' },
  { to: '/accounting', label: 'المحاسبة', icon: '⚖' },
  { to: '/reports', label: 'التقارير', icon: '▤' },
  { to: '/e-invoicing', label: 'الفوترة الإلكترونية', icon: '⎙' },
  { to: '/settings', label: 'الإعدادات', icon: '⚙' },
];

export default function Shell() {
  const { me, logout, switchTenant } = useAuth();
  const [open, setOpen] = useState(false);
  if (!me) return null;

  return (
    <div className={`shell ${open ? 'nav-open' : ''}`}>
      <aside className="sidebar" aria-label="القائمة الرئيسية">
        <div className="brand">
          <span className="brand-mark" aria-hidden>ش</span>
          <span>الشيوخ للمحاسبة</span>
        </div>
        <nav>
          {NAV.map((n) => (
            <NavLink key={n.to} to={n.to} end={n.to === '/'} onClick={() => setOpen(false)}
              className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
              <span className="nav-icon" aria-hidden>{n.icon}</span>
              <span>{n.label}</span>
              {n.phase && <span className="badge-soon">قريبًا</span>}
            </NavLink>
          ))}
        </nav>
      </aside>
      {open && <div className="scrim" onClick={() => setOpen(false)} />}

      <div className="main">
        <header className="topbar">
          <button className="icon-btn menu-btn" aria-label="فتح القائمة" onClick={() => setOpen(true)}>☰</button>
          <div className="tenant">
            {me.memberships.length > 1 ? (
              <select aria-label="المنشأة" value={me.tenant.id} onChange={(e) => switchTenant(e.target.value)}>
                {me.memberships.map((m) => <option key={m.tenantId} value={m.tenantId}>{m.tenantName}</option>)}
              </select>
            ) : <strong>{me.tenant.name}</strong>}
          </div>
          <div className="user">
            <span className="user-name">{me.user.fullName}</span>
            <span className="muted small">{me.roles.map((r) => r.nameAr).join('، ')}</span>
          </div>
          <button className="btn btn-ghost" onClick={logout}>تسجيل الخروج</button>
        </header>
        <main className="content">
          {me.user.mustChangePassword && (
            <div className="alert alert-warn">يجب تغيير كلمة المرور المؤقتة من صفحة الإعدادات.</div>
          )}
          <Outlet />
        </main>
      </div>
    </div>
  );
}
