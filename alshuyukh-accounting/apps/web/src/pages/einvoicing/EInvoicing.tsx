import { Fragment, useState, type FormEvent } from 'react';
import { Link } from 'react-router-dom';
import { api, downloadFile } from '../../api';
import { useAuth } from '../../auth';
import { ErrorBox, formatDateTime, PageHeader, useLoad } from '../../ui';

interface Device {
  id: string; name: string; branchName: string | null; environment: string; status: string; egsSerial: string; invoiceTypes: string;
  csr: string; complianceChecks: Record<string, string>; certificateSerial: string | null; certificateIssuer: string | null;
  certificateExpiresAt: string | null; icv: string; lastHash: string; activatedAt: string | null; createdAt: string;
}
interface EInvoice {
  id: string; documentType: string; documentId: string; documentNumber: string; invoiceKind: string; typeCode: string; uuid: string; icv: string;
  status: string; hasWarnings: boolean; attempts: number; nextAttemptAt: string; lastError: string | null; submittedAt: string | null;
  createdAt: string; reportingOverdue: boolean;
}
interface Detail extends EInvoice {
  invoiceHash: string; previousHash: string;
  submissions: { id: string; operation: string; httpStatus: number | null; outcome: string; createdAt: string; messages: { level: string; code: string | null; message: string }[] }[];
  documents: { kind: string }[];
}

const ENV_AR: Record<string, string> = { DEVELOPER: 'بوابة المطورين (تجريبي)', SIMULATION: 'بيئة المحاكاة', PRODUCTION: 'الإنتاج' };
const DEVICE_STATUS: Record<string, [string, string]> = {
  NEW: ['بانتظار رمز OTP', 'status-draft'], COMPLIANCE: ['فحوص الامتثال', 'status-partial'], ACTIVE: ['مفعّلة', 'status-posted'], REVOKED: ['ملغاة', 'status-reversed'],
};
export const EINVOICE_STATUS: Record<string, [string, string]> = {
  PENDING: ['بانتظار الإرسال', 'status-partial'], REPORTED: ['مُبلَّغ عنها', 'status-posted'], CLEARED: ['معتمدة', 'status-posted'], REJECTED: ['مرفوضة', 'status-reversed'],
};
const TYPE_AR: Record<string, string> = { '388': 'فاتورة', '381': 'إشعار دائن', '383': 'إشعار مدين' };
const CHECK_AR: Record<string, string> = {
  STANDARD_INVOICE: 'فاتورة ضريبية', STANDARD_CREDIT_NOTE: 'إشعار دائن ضريبي', STANDARD_DEBIT_NOTE: 'إشعار مدين ضريبي',
  SIMPLIFIED_INVOICE: 'فاتورة مبسطة', SIMPLIFIED_CREDIT_NOTE: 'إشعار دائن مبسط', SIMPLIFIED_DEBIT_NOTE: 'إشعار مدين مبسط',
};

export default function EInvoicing() {
  const { can } = useAuth();
  const devices = useLoad(() => api<{ data: Device[] }>('GET', '/api/zatca/devices'));
  const [status, setStatus] = useState('');
  const list = useLoad(() => api<{ data: EInvoice[]; total: number; summary: { pending: number; rejected: number; overdue: number } }>(
    'GET', `/api/zatca/invoices?limit=100${status ? `&status=${status}` : ''}`), [status]);
  const live = devices.data?.data.filter((d) => d.status !== 'REVOKED') ?? [];
  return (
    <>
      <PageHeader title="الفوترة الإلكترونية (فاتورة)" />
      <div className="alert alert-warn">
        التكامل مع منصة فاتورة (ZATCA) مبني وفق المواصفات ومختبر داخليًا فقط، ولم يُختبر بعد مع بوابة الهيئة الفعلية.
        ابدأ ببيئة «بوابة المطورين» ثم «المحاكاة» قبل «الإنتاج».
      </div>
      <ErrorBox error={devices.error ?? list.error} />
      {devices.data && live.length === 0 && can('zatca.manage') && <NewDevice onDone={devices.reload} />}
      {devices.data && live.length === 0 && !can('zatca.manage') && <div className="card muted">لم تُربط الشركة بمنصة فاتورة بعد. الفواتير تُطبع برمز QR للمرحلة الأولى.</div>}
      {devices.data?.data.map((d) => <DeviceCard key={d.id} d={d} onChange={() => { devices.reload(); list.reload(); }} />)}

      <div className="card">
        <div className="card-head">
          <h2>الفواتير الإلكترونية</h2>
          {list.data && (
            <span className="muted small">
              {list.data.summary.pending} بانتظار الإرسال · {list.data.summary.rejected} مرفوضة
              {list.data.summary.overdue > 0 && <strong className="negative"> · {list.data.summary.overdue} تجاوزت مهلة 24 ساعة للإبلاغ</strong>}
            </span>
          )}
        </div>
        <div className="toolbar">
          <select aria-label="الحالة" value={status} onChange={(e) => setStatus(e.target.value)}>
            <option value="">كل الحالات</option>
            {Object.entries(EINVOICE_STATUS).map(([k, [l]]) => <option key={k} value={k}>{l}</option>)}
          </select>
        </div>
        <EInvoiceTable rows={list.data?.data ?? []} onChange={list.reload} />
      </div>
    </>
  );
}

function NewDevice({ onDone }: { onDone: () => void }) {
  const [error, setError] = useState<unknown>(null);
  async function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = Object.fromEntries(new FormData(e.currentTarget)) as Record<string, string>;
    setError(null);
    try { await api('POST', '/api/zatca/devices', { name: f.name, environment: f.environment, businessCategory: f.businessCategory, invoiceTypes: f.invoiceTypes }); onDone(); }
    catch (err) { setError(err); }
  }
  return (
    <form className="card form" onSubmit={submit}>
      <h2>ربط الشركة بمنصة فاتورة</h2>
      <p className="muted small">تُنشأ وحدة توليد فواتير (EGS) بمفتاح تشفير خاص يبقى مشفرًا في الخادم، وطلب شهادة (CSR) يُرسل إلى الهيئة مع رمز OTP من بوابة فاتورة.</p>
      <ErrorBox error={error} />
      <div className="grid-2">
        <label>اسم الوحدة<input name="name" required minLength={2} defaultValue="الوحدة الرئيسية" /></label>
        <label>البيئة
          <select name="environment" defaultValue="DEVELOPER">{Object.entries(ENV_AR).map(([k, v]) => <option key={k} value={k}>{v}</option>)}</select>
        </label>
        <label>نشاط المنشأة (يظهر في الشهادة)<input name="businessCategory" required minLength={2} placeholder="تجارة التجزئة" /></label>
        <label>أنواع الفواتير
          <select name="invoiceTypes" defaultValue="1100">
            <option value="1100">ضريبية ومبسطة</option><option value="1000">ضريبية فقط (للمنشآت)</option><option value="0100">مبسطة فقط (للأفراد)</option>
          </select>
        </label>
      </div>
      <p className="muted small">يجب أن تكون بيانات الشركة مكتملة: الرقم الضريبي، السجل التجاري، والعنوان الوطني (<Link to="/settings/companies">إعدادات الشركة</Link>).</p>
      <button className="btn btn-primary">إنشاء الوحدة وطلب الشهادة</button>
    </form>
  );
}

function DeviceCard({ d, onChange }: { d: Device; onChange: () => void }) {
  const { can } = useAuth();
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState(false);
  const manage = can('zatca.manage') && d.status !== 'REVOKED';
  async function run(fn: () => Promise<unknown>) {
    setError(null); setBusy(true);
    try { await fn(); onChange(); } catch (e) { setError(e); } finally { setBusy(false); }
  }
  const [label, cls] = DEVICE_STATUS[d.status] ?? [d.status, ''];
  const steps = [
    { done: true, text: 'إنشاء الوحدة وطلب الشهادة (CSR)' },
    { done: d.status !== 'NEW', text: 'إدخال رمز OTP من بوابة فاتورة والحصول على شهادة الامتثال' },
    { done: d.status === 'ACTIVE' || (Object.keys(d.complianceChecks).length > 0 && !Object.values(d.complianceChecks).includes('FAILED')), text: 'فحوص الامتثال لكل أنواع المستندات' },
    { done: d.status === 'ACTIVE', text: 'الحصول على شهادة الإنتاج وتفعيل التوقيع' },
  ];
  return (
    <div className="card">
      <div className="entry-head">
        <div><h2>{d.name}{d.branchName ? ` — ${d.branchName}` : ''}</h2><p className="muted small">{ENV_AR[d.environment]} · <span dir="ltr">{d.egsSerial}</span></p></div>
        <span className={`tag ${cls}`}>{label}</span>
      </div>
      <ErrorBox error={error} />
      <ol className="steps">{steps.map((s) => <li key={s.text} className={s.done ? 'done' : ''}><span aria-hidden>{s.done ? '✓' : '○'}</span>{s.text}</li>)}</ol>
      {Object.keys(d.complianceChecks).length > 0 && (
        <div className="chips">{Object.entries(d.complianceChecks).map(([k, v]) => <span key={k} className={`tag ${v === 'FAILED' ? 'status-reversed' : 'status-posted'}`}>{CHECK_AR[k] ?? k}: {v === 'FAILED' ? 'فشل' : v === 'PASSED' ? 'نجح' : 'نجح مع تحذيرات'}</span>)}</div>
      )}
      {d.status === 'ACTIVE' && (
        <dl className="meta">
          <div><dt>آخر عداد (ICV)</dt><dd dir="ltr">{d.icv}</dd></div>
          <div><dt>الرقم التسلسلي للشهادة</dt><dd dir="ltr">{d.certificateSerial}</dd></div>
          <div><dt>تنتهي الشهادة</dt><dd>{formatDateTime(d.certificateExpiresAt)}</dd></div>
          <div><dt>الجهة المصدرة</dt><dd dir="ltr" className="small">{d.certificateIssuer}</dd></div>
        </dl>
      )}
      {manage && d.status === 'NEW' && (
        <form className="inline-form" onSubmit={(e) => { e.preventDefault(); const otp = String(new FormData(e.currentTarget).get('otp')); void run(() => api('POST', `/api/zatca/devices/${d.id}/compliance-csid`, { otp })); }}>
          <input name="otp" required pattern="[0-9]{6}" inputMode="numeric" dir="ltr" placeholder="OTP (6 أرقام)" aria-label="رمز OTP" />
          <button className="btn btn-primary" disabled={busy}>طلب شهادة الامتثال</button>
          <button type="button" className="btn" onClick={() => navigator.clipboard?.writeText(d.csr)}>نسخ CSR</button>
        </form>
      )}
      {manage && d.status === 'COMPLIANCE' && (
        <div className="form-actions">
          <button className="btn" disabled={busy} onClick={() => run(() => api('POST', `/api/zatca/devices/${d.id}/compliance-checks`))}>{busy ? 'جارٍ الفحص…' : 'تشغيل فحوص الامتثال'}</button>
          <button className="btn btn-primary" disabled={busy} onClick={() => run(() => api('POST', `/api/zatca/devices/${d.id}/activate`))}>طلب شهادة الإنتاج والتفعيل</button>
        </div>
      )}
      {manage && (
        <details className="no-print"><summary className="muted small">إلغاء الوحدة</summary>
          <form className="inline-form" onSubmit={(e) => { e.preventDefault(); const reason = String(new FormData(e.currentTarget).get('reason')); void run(() => api('POST', `/api/zatca/devices/${d.id}/revoke`, { reason })); }}>
            <input name="reason" required minLength={3} placeholder="سبب الإلغاء" />
            <button className="btn btn-danger" disabled={busy}>إلغاء الوحدة</button>
          </form>
        </details>
      )}
    </div>
  );
}

export function EInvoiceTable({ rows, onChange }: { rows: EInvoice[]; onChange: () => void }) {
  const { can } = useAuth();
  const [open, setOpen] = useState<string | null>(null);
  const [detail, setDetail] = useState<Detail | null>(null);
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState<string | null>(null);
  async function toggle(id: string) {
    if (open === id) { setOpen(null); return; }
    setOpen(id); setDetail(null);
    try { setDetail(await api<Detail>('GET', `/api/zatca/invoices/${id}`)); } catch (e) { setError(e); }
  }
  async function submit(id: string) {
    setError(null); setBusy(id);
    try { const d = await api<Detail>('POST', `/api/zatca/invoices/${id}/submit`); if (open === id) setDetail(d); onChange(); } catch (e) { setError(e); } finally { setBusy(null); }
  }
  const docLink = (r: EInvoice) => `${r.documentType === 'SALES_INVOICE' ? '/sales/invoices' : '/sales/returns'}/${r.documentId}`;
  return (
    <>
      <ErrorBox error={error} />
      <div className="table-wrap">
        <table className="table">
          <thead><tr><th>المستند</th><th>النوع</th><th>ICV</th><th>الإنشاء</th><th>الحالة</th><th /></tr></thead>
          <tbody>
            {rows.map((r) => {
              const [label, cls] = EINVOICE_STATUS[r.status] ?? [r.status, ''];
              return (
                <Fragment key={r.id}>
                  <tr>
                    <td><Link to={docLink(r)} dir="ltr" className="code">{r.documentNumber}</Link></td>
                    <td>{TYPE_AR[r.typeCode]} {r.invoiceKind === 'STANDARD' ? 'ضريبي' : 'مبسط'}</td>
                    <td dir="ltr">{r.icv}</td>
                    <td>{formatDateTime(r.createdAt)}</td>
                    <td>
                      <span className={`tag ${cls}`}>{label}{r.hasWarnings ? ' (تحذيرات)' : ''}</span>
                      {r.reportingOverdue && <span className="tag status-reversed">متأخرة عن 24 ساعة</span>}
                      {r.status === 'PENDING' && r.attempts > 0 && <div className="muted small">محاولات: {r.attempts} — التالية {formatDateTime(r.nextAttemptAt)}</div>}
                    </td>
                    <td className="actions-cell">
                      {r.status === 'PENDING' && can('invoice.post') && <button className="btn btn-small btn-primary" disabled={busy === r.id} onClick={() => submit(r.id)}>{busy === r.id ? '…' : r.invoiceKind === 'STANDARD' ? 'طلب الاعتماد' : 'إبلاغ'}</button>}
                      <button className="btn btn-small" onClick={() => toggle(r.id)}>{open === r.id ? 'إخفاء' : 'التفاصيل'}</button>
                    </td>
                  </tr>
                  {open === r.id && (
                    <tr className="sub-detail"><td colSpan={6}>
                      {!detail ? 'جارٍ التحميل…' : (
                        <div className="einvoice-detail">
                          {detail.lastError && <div className="alert alert-error">{detail.lastError}</div>}
                          <dl className="meta">
                            <div><dt>UUID</dt><dd dir="ltr" className="small">{detail.uuid}</dd></div>
                            <div><dt>بصمة الفاتورة</dt><dd dir="ltr" className="small code">{detail.invoiceHash}</dd></div>
                            <div><dt>بصمة السابقة (PIH)</dt><dd dir="ltr" className="small code">{detail.previousHash}</dd></div>
                          </dl>
                          <div className="form-actions">
                            {detail.documents.map((doc) => (
                              <button key={doc.kind} className="btn btn-small" onClick={() => downloadFile(`/api/zatca/invoices/${detail.id}/xml?kind=${doc.kind}`, `${detail.documentNumber}-${doc.kind.toLowerCase()}.xml`).catch(setError)}>
                                تنزيل XML {doc.kind === 'CLEARED' ? 'المعتمد من الهيئة' : 'الموقّع'}
                              </button>
                            ))}
                          </div>
                          {detail.submissions.length > 0 && (
                            <table className="table">
                              <thead><tr><th>العملية</th><th>الوقت</th><th>HTTP</th><th>النتيجة</th><th>الرسائل</th></tr></thead>
                              <tbody>{detail.submissions.map((s) => (
                                <tr key={s.id}>
                                  <td>{s.operation === 'CLEARANCE' ? 'اعتماد' : 'إبلاغ'}</td><td>{formatDateTime(s.createdAt)}</td><td dir="ltr">{s.httpStatus ?? '—'}</td>
                                  <td>{{ SUCCESS: 'نجاح', WARNING: 'نجاح مع تحذيرات', REJECTED: 'رفض', TRANSPORT_ERROR: 'تعذر الاتصال' }[s.outcome]}</td>
                                  <td>{s.messages.map((m, i) => <div key={i} className={m.level === 'ERROR' ? 'negative' : 'muted'}><span dir="ltr" className="code">{m.code}</span> {m.message}</div>)}</td>
                                </tr>
                              ))}</tbody>
                            </table>
                          )}
                        </div>
                      )}
                    </td></tr>
                  )}
                </Fragment>
              );
            })}
            {rows.length === 0 && <tr><td colSpan={6} className="muted empty-row">لا توجد فواتير إلكترونية</td></tr>}
          </tbody>
        </table>
      </div>
    </>
  );
}
