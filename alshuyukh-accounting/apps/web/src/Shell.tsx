import { Suspense, useState } from 'react';
import { api } from './api';
import { ErrorBox, useLoad } from './ui';
import { Link, NavLink, Outlet } from 'react-router-dom';
import { useAuth } from './auth';

const NAV: { to: string; label: string; icon: string; phase?: number; feature?: string }[] = [
  { to: '/', label: 'الرئيسية', icon: '⌂' },
  { to: '/sales', label: 'المبيعات', icon: '↗' },
  { to: '/purchases', label: 'المشتريات', icon: '↙' },
  { to: '/inventory', label: 'المخزون', icon: '▦' },
  { to: '/customers', label: 'العملاء', icon: '☺' },
  { to: '/suppliers', label: 'الموردون', icon: '⚑' },
  { to: '/expenses', label: 'المصروفات', icon: '−' },
  { to: '/accounting', label: 'المحاسبة', icon: '⚖' },
  { to: '/reports', label: 'التقارير', icon: '▤' },
  { to: '/e-invoicing', label: 'الفوترة الإلكترونية', icon: '⎙', feature: 'zatca_einvoicing' },
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
          {NAV.filter((n) => !n.feature || me.features[n.feature] !== false).map((n) => (
            <NavLink key={n.to} to={n.to} end={n.to === '/'} onClick={() => setOpen(false)}
              className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
              <span className="nav-icon" aria-hidden>{n.icon}</span>
              <span>{n.label}</span>
              {n.phase && <span className="badge-soon">قريبًا</span>}
            </NavLink>
          ))}
          {me.user.isPlatformAdmin && (
            <NavLink to="/admin" onClick={() => setOpen(false)} className={({ isActive }) => `nav-item nav-admin ${isActive ? 'active' : ''}`}>
              <span className="nav-icon" aria-hidden>★</span><span>إدارة المنصة</span>
            </NavLink>
          )}
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
          <SubscriptionBanner />
          <Invitations />
          {me.user.mustChangePassword && (
            <div className="alert alert-warn">يجب تغيير كلمة المرور المؤقتة من صفحة الإعدادات.</div>
          )}
          <Suspense fallback={<div className="muted">جارٍ التحميل…</div>}><Outlet /></Suspense>
        </main>
      </div>
    </div>
  );
}

/** Trial ending soon, grace period, or expired (read-only) subscription. */
function SubscriptionBanner() {
  const { me, can } = useAuth();
  const s = me?.subscription;
  if (!s) return null;
  const days = s.periodEnd ? Math.ceil((new Date(s.periodEnd).getTime() - Date.now()) / 86_400_000) : 0;
  const link = can('subscription.manage') ? <Link to="/settings/subscription">إدارة الاشتراك</Link> : <span>تواصل مع مالك المنشأة.</span>;
  if (s.state === 'EXPIRED' || s.state === 'NONE' || s.state === 'CANCELLED') {
    return <div className="alert alert-error">انتهى الاشتراك؛ البيانات متاحة للاطلاع فقط ولا يمكن إضافة أو تعديل شيء حتى التجديد. {link}</div>;
  }
  if (s.state === 'GRACE') {
    const left = s.graceEnd ? Math.max(0, Math.ceil((new Date(s.graceEnd).getTime() - Date.now()) / 86_400_000)) : 0;
    return <div className="alert alert-warn">انتهت فترة الاشتراك. يبقى الوصول الكامل {left} يومًا ثم يتحول النظام للاطلاع فقط. {link}</div>;
  }
  if (s.state === 'TRIALING' && days <= 3) {
    return <div className="alert alert-warn">تنتهي الفترة التجريبية خلال {Math.max(days, 0)} يوم. {link}</div>;
  }
  return null;
}

/** Organizations that invited this user; joining needs the user's consent. */
function Invitations() {
  const { reload } = useAuth();
  const list = useLoad(() => api<{ data: { tenantId: string; tenantName: string }[] }>('GET', '/api/auth/invitations'));
  const [error, setError] = useState<unknown>(null);
  async function answer(tenantId: string, action: 'accept' | 'decline') {
    setError(null);
    try { await api('POST', `/api/auth/invitations/${tenantId}/${action}`); list.reload(); if (action === 'accept') await reload(); } catch (e) { setError(e); }
  }
  if (!list.data?.data.length) return null;
  return (
    <div className="alert alert-warn">
      <ErrorBox error={error} />
      {list.data.data.map((i) => (
        <div key={i.tenantId} className="invite-row">
          دعتك منشأة «{i.tenantName}» للانضمام إليها.
          <button className="btn btn-small btn-primary" onClick={() => answer(i.tenantId, 'accept')}>قبول</button>
          <button className="btn btn-small" onClick={() => answer(i.tenantId, 'decline')}>رفض</button>
        </div>
      ))}
    </div>
  );
}
