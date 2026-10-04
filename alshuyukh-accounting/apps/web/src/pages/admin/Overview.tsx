import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { STATE_AR } from '../subscription/Subscription';

interface Overview { tenants: number; suspended: number; users: number; newTenants30d: number; errors24h: number; revenueThisMonth: string; mrr: string; subscriptions: Record<string, number> }
interface Health {
  status: string; version: string; node: string; uptimeSeconds: number; memoryMb: { rss: number; heapUsed: number };
  database: { latencyMs: number; postgres: string; sizeMb: number; migrations: number; lastMigration: string; pool: { total: number; idle: number; waiting: number } };
  zatca: { pending: number; overdue: number; worker: string }; errors24h: number;
}

export function Overview() {
  const o = useLoad(() => api<Overview>('GET', '/api/admin/overview'));
  const h = useLoad(() => api<Health>('GET', '/api/admin/health'));
  const tiles: [string, string | number, string?][] = o.data ? [
    ['المنشآت', o.data.tenants], ['موقوفة', o.data.suspended, o.data.suspended ? 'negative' : ''], ['المستخدمون', o.data.users],
    ['منشآت جديدة (30 يومًا)', o.data.newTenants30d], ['الإيراد الشهري المتكرر (MRR)', `${formatAmount(o.data.mrr)} ر.س`],
    ['المحصّل هذا الشهر', `${formatAmount(o.data.revenueThisMonth)} ر.س`], ['أخطاء آخر 24 ساعة', o.data.errors24h, o.data.errors24h ? 'negative' : ''],
  ] : [];
  return (
    <>
      <ErrorBox error={o.error ?? h.error} />
      <div className="kpi-grid">{tiles.map(([l, v, c]) => <div key={l} className="card kpi"><span className="muted small">{l}</span><strong className={c}>{v}</strong></div>)}</div>
      {o.data && (
        <div className="card">
          <h2>الاشتراكات حسب الحالة</h2>
          <div className="chips">{Object.entries(o.data.subscriptions).map(([k, n]) => <span key={k} className={`tag ${STATE_AR[k]?.[1] ?? ''}`}>{STATE_AR[k]?.[0] ?? k}: {n}</span>)}</div>
        </div>
      )}
      {h.data && (
        <div className="card">
          <div className="card-head"><h2>صحة النظام</h2><span className="tag status-posted">{h.data.status === 'ok' ? 'يعمل' : h.data.status}</span></div>
          <dl className="meta">
            <div><dt>الإصدار</dt><dd dir="ltr">{h.data.version} · Node {h.data.node}</dd></div>
            <div><dt>مدة التشغيل</dt><dd>{Math.floor(h.data.uptimeSeconds / 3600)} ساعة {Math.floor((h.data.uptimeSeconds % 3600) / 60)} دقيقة</dd></div>
            <div><dt>الذاكرة</dt><dd dir="ltr">{h.data.memoryMb.heapUsed} / {h.data.memoryMb.rss} MB</dd></div>
            <div><dt>قاعدة البيانات</dt><dd dir="ltr">{h.data.database.postgres} · {h.data.database.latencyMs} ms · {h.data.database.sizeMb} MB</dd></div>
            <div><dt>الاتصالات</dt><dd dir="ltr">{h.data.database.pool.total} / idle {h.data.database.pool.idle} / waiting {h.data.database.pool.waiting}</dd></div>
            <div><dt>الترحيلات</dt><dd dir="ltr">{h.data.database.migrations} · {h.data.database.lastMigration}</dd></div>
            <div><dt>الفوترة الإلكترونية</dt><dd>معلقة {h.data.zatca.pending} · متأخرة {h.data.zatca.overdue} · المرسل {h.data.zatca.worker === 'on' ? 'يعمل' : 'متوقف'}</dd></div>
          </dl>
        </div>
      )}
    </>
  );
}
