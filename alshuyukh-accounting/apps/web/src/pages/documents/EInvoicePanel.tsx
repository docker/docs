import { useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { ErrorBox, useLoad } from '../../ui';
import { EINVOICE_STATUS } from '../einvoicing/EInvoicing';

export interface ZatcaDocument {
  einvoice: { id: string; status: string; invoiceKind: string; icv: string; hasWarnings: boolean; lastError: string | null; reportingOverdue: boolean } | null;
  qr: string | null; qrSvg: string | null; missing?: string;
}

export const loadZatcaDocument = (type: string, id: string) => api<ZatcaDocument>('GET', `/api/zatca/document?type=${type}&id=${id}`);

/** E-invoice status on a sales invoice or credit note, with the submit action and the print link. */
export default function EInvoicePanel({ type, id, printPath }: { type: 'SALES_INVOICE' | 'SALES_RETURN'; id: string; printPath: string }) {
  const { can } = useAuth();
  const z = useLoad(() => loadZatcaDocument(type, id), [type, id]);
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState(false);
  const e = z.data?.einvoice;
  async function submit() {
    setError(null); setBusy(true);
    try { await api('POST', `/api/zatca/invoices/${e!.id}/submit`); z.reload(); } catch (err) { setError(err); } finally { setBusy(false); }
  }
  if (!z.data) return <ErrorBox error={z.error} />;
  const [label, cls] = e ? EINVOICE_STATUS[e.status] ?? [e.status, ''] : ['', ''];
  return (
    <div className="einvoice-panel">
      <ErrorBox error={error} />
      {e ? (
        <>
          <span>الفاتورة الإلكترونية:</span> <span className={`tag ${cls}`}>{label}{e.hasWarnings ? ' (تحذيرات)' : ''}</span>
          <span className="muted small"> ICV <span dir="ltr">{e.icv}</span></span>
          {e.reportingOverdue && <span className="tag status-reversed">تجاوزت مهلة الإبلاغ (24 ساعة)</span>}
          {e.status === 'PENDING' && can('invoice.post') && (
            <button className="btn btn-small btn-primary" disabled={busy} onClick={submit}>{busy ? '…' : e.invoiceKind === 'STANDARD' ? 'طلب اعتماد الهيئة' : 'إبلاغ الهيئة'}</button>
          )}
          {e.lastError && <div className="alert alert-error small">{e.lastError}</div>}
          {e.invoiceKind === 'STANDARD' && e.status !== 'CLEARED' && <p className="muted small">الفاتورة الضريبية لا تُسلَّم للعميل قبل اعتمادها من الهيئة.</p>}
          {can('zatca.view') && <Link className="small" to="/e-invoicing">سجل الفوترة الإلكترونية</Link>}
        </>
      ) : (
        <span className="muted small">{z.data.qr ? 'لا توجد وحدة فوترة مفعّلة؛ تُطبع الفاتورة برمز QR للمرحلة الأولى.' : z.data.missing ? 'أضف الرقم الضريبي للشركة لطباعة رمز QR.' : ''}</span>
      )}
      <Link className="btn btn-small" to={printPath}>طباعة</Link>
    </div>
  );
}
