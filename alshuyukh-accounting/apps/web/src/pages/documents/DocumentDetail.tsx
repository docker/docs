import { useState, type FormEvent } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { DOCS, STATUS_AR, statusClass, type DocConfig, type Doc } from './config';
import DocumentForm from './DocumentForm';
import Totals from './Totals';

const today = () => new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh' }).format(new Date());

interface Method { id: string; code: string; nameAr: string; isActive: boolean }

export default function DocumentDetail({ config }: { config: DocConfig }) {
  const { id } = useParams();
  const { can } = useAuth();
  const navigate = useNavigate();
  const doc = useLoad(() => api<Doc>('GET', `/api/${config.api}/${id}`), [id, config.api]);
  const [mode, setMode] = useState<'view' | 'edit' | 'return' | 'pay' | 'cancel'>('view');
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState(false);

  async function run(fn: () => Promise<unknown>) {
    setError(null);
    setBusy(true);
    try { await fn(); setMode('view'); doc.reload(); } catch (e) { setError(e); } finally { setBusy(false); }
  }

  const d = doc.data;
  if (!d) return <ErrorBox error={doc.error} />;
  if (mode === 'edit') return <DocumentForm config={config} doc={d} onSaved={() => { setMode('view'); doc.reload(); }} />;

  const isInvoice = config.key === 'SALES_INVOICE' || config.key === 'PURCHASE_INVOICE';
  const open = !['DRAFT', 'CANCELLED'].includes(d.status) && Number(d.remainingAmount ?? 0) > 0;
  const returnConfig = config.key === 'SALES_INVOICE' ? DOCS.SALES_RETURN : DOCS.PURCHASE_RETURN;
  const invoiceConfig = config.party === 'customer' ? DOCS.SALES_INVOICE : DOCS.PURCHASE_INVOICE;
  const canPay = can('payment.create') && open && (isInvoice || config.isReturn);

  return (
    <div className="card">
      <ErrorBox error={error} />
      <div className="entry-head">
        <div>
          <h2 dir="ltr" className="code">{d.number ?? 'مسودة'}</h2>
          <p>{config.singular} — <Link to={`/${config.party}s`}>{d.partyName}</Link></p>
        </div>
        <span className={`tag ${statusClass(d.status)}`}>{STATUS_AR[d.status] ?? d.status}</span>
      </div>
      <dl className="meta">
        <div><dt>التاريخ</dt><dd dir="ltr">{d.date}</dd></div>
        {d.dueDate && <div><dt>الاستحقاق</dt><dd dir="ltr">{d.dueDate}</dd></div>}
        {d.invoiceKind && d.status !== 'DRAFT' && <div><dt>نوع الفاتورة الضريبية</dt><dd>{d.invoiceKind === 'STANDARD' ? 'فاتورة ضريبية' : 'فاتورة ضريبية مبسطة'}</dd></div>}
        {d.supplierInvoiceNumber && <div><dt>رقم فاتورة المورد</dt><dd dir="ltr">{d.supplierInvoiceNumber}</dd></div>}
        {d.originalInvoiceId && <div><dt>الفاتورة الأصلية</dt><dd><Link to={`../../invoices/${d.originalInvoiceId}`} relative="path" dir="ltr">{d.originalInvoiceNumber}</Link></dd></div>}
        {d.reason && <div><dt>سبب الإرجاع</dt><dd>{d.reason}</dd></div>}
        {d.journalEntryId && <div><dt>القيد</dt><dd><Link to={`/accounting/journal/${d.journalEntryId}`}>عرض القيد المحاسبي</Link></dd></div>}
        {d.convertedInvoiceId && <div><dt>الفاتورة</dt><dd><Link to={`../../invoices/${d.convertedInvoiceId}`} relative="path">عرض الفاتورة</Link></dd></div>}
        {d.cancelReason && <div><dt>سبب الإلغاء</dt><dd>{d.cancelReason}</dd></div>}
        <div><dt>الأسعار</dt><dd>{d.pricesIncludeVat ? 'شاملة الضريبة' : 'غير شاملة الضريبة'}</dd></div>
      </dl>

      <div className="table-wrap">
        <table className="table">
          <thead><tr><th>#</th><th>الصنف</th><th className="num">الكمية</th><th className="num">السعر</th><th className="num">الخصم</th><th className="num">الصافي</th><th className="num">الضريبة</th><th className="num">الإجمالي</th></tr></thead>
          <tbody>
            {d.lines.map((l) => (
              <tr key={l.id}>
                <td>{l.lineNo}</td>
                <td>{l.sku && <><span className="code" dir="ltr">{l.sku}</span>{' — '}</>}{l.description ?? l.productName}</td>
                <td className="num" dir="ltr">{l.quantity.replace(/\.?0+$/, '')} {l.unitCode}</td>
                <td className="num" dir="ltr">{l.unitPrice.replace(/(\.\d\d)00$/, '$1')}</td>
                <td className="num" dir="ltr">{l.discountAmount === '0.00' ? '' : formatAmount(l.discountAmount)}</td>
                <td className="num" dir="ltr">{formatAmount(l.netAmount)}</td>
                <td className="num" dir="ltr">{formatAmount(l.vatAmount)} <span className="muted small">{l.vatCategory === 'S' ? `${(Number(l.vatRate) * 100).toFixed(0)}%` : l.vatCategory}</span></td>
                <td className="num" dir="ltr">{formatAmount(l.totalAmount)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="doc-footer">
        <div className="grow">{d.notes && <p><strong>ملاحظات:</strong> {d.notes}</p>}</div>
        <Totals t={d} extra={config.legal && d.status !== 'DRAFT' ? [
          ...(isInvoice ? [['المدفوع', d.paidAmount], ['المرتجع', d.returnedAmount]] as [string, string | undefined][] : [['المطبّق على الفاتورة', d.appliedAmount], ['المسترد', d.refundedAmount]] as [string, string | undefined][]),
          ['المتبقي', d.remainingAmount],
        ] : undefined} />
      </div>

      {mode === 'view' && (
        <div className="form-actions">
          {d.status === 'DRAFT' && can(config.perm.create) && !config.isReturn && <button className="btn" onClick={() => setMode('edit')}>تعديل</button>}
          {d.status === 'DRAFT' && config.legal && can(config.perm.issue) && (
            <button className="btn btn-primary" disabled={busy} onClick={() => run(() => api('POST', `/api/${config.api}/${d.id}/post`))}>
              {config.party === 'customer' ? 'إصدار' : 'ترحيل'}
            </button>
          )}
          {config.key === 'SALES_QUOTE' && d.status === 'DRAFT' && <button className="btn" onClick={() => run(() => api('POST', `/api/${config.api}/${d.id}/status`, { status: 'SENT' }))}>تحديد كمرسل</button>}
          {config.key === 'SALES_QUOTE' && d.status === 'SENT' && (
            <>
              <button className="btn" onClick={() => run(() => api('POST', `/api/${config.api}/${d.id}/status`, { status: 'ACCEPTED' }))}>قبول</button>
              <button className="btn" onClick={() => run(() => api('POST', `/api/${config.api}/${d.id}/status`, { status: 'REJECTED' }))}>رفض</button>
            </>
          )}
          {config.key === 'PURCHASE_ORDER' && d.status === 'DRAFT' && can(config.perm.issue) && <button className="btn btn-primary" onClick={() => run(() => api('POST', `/api/${config.api}/${d.id}/status`, { status: 'APPROVED' }))}>اعتماد</button>}
          {((config.key === 'SALES_QUOTE' && ['DRAFT', 'SENT', 'ACCEPTED'].includes(d.status)) || (config.key === 'PURCHASE_ORDER' && d.status === 'APPROVED')) && can(config.perm.create) && (
            <button className="btn btn-primary" onClick={() => run(async () => {
              const inv = await api<{ id: string }>('POST', `/api/${config.api}/${d.id}/convert`, {});
              navigate(`../../invoices/${inv.id}`, { relative: 'path' });
            })}>تحويل إلى فاتورة</button>
          )}
          {canPay && <button className="btn btn-primary" onClick={() => setMode('pay')}>{config.isReturn ? 'تسجيل استرداد' : 'تسجيل دفعة'}</button>}
          {isInvoice && ['ISSUED', 'POSTED', 'PARTIALLY_PAID', 'PAID'].includes(d.status) && can(config.perm.create) && <button className="btn" onClick={() => setMode('return')}>إنشاء مرتجع</button>}
          {d.status === 'DRAFT' && can(config.perm.create) && (
            <button className="btn btn-danger" onClick={() => { if (confirm('حذف المسودة؟')) void run(async () => { await api('DELETE', `/api/${config.api}/${d.id}`); navigate('..', { relative: 'path' }); }); }}>حذف المسودة</button>
          )}
          {!['DRAFT', 'CANCELLED', 'CONVERTED', 'REJECTED'].includes(d.status) && can(config.perm.cancel) && <button className="btn btn-danger" onClick={() => setMode('cancel')}>إلغاء</button>}
        </div>
      )}

      {mode === 'cancel' && (
        <form className="reverse-box" onSubmit={(e) => { e.preventDefault(); const reason = String(new FormData(e.currentTarget).get('reason')); void run(() => api('POST', `/api/${config.api}/${d.id}/cancel`, { reason })); }}>
          <p>{config.legal ? 'يُعكس القيد المحاسبي بتاريخ المستند ويبقى الرقم محفوظًا. لا يمكن إلغاء مستند عليه دفعات أو مرتجعات.' : 'سيُلغى المستند.'}</p>
          <label>سبب الإلغاء<input name="reason" required minLength={3} /></label>
          <div className="form-actions"><button className="btn btn-danger" disabled={busy}>تأكيد الإلغاء</button><button type="button" className="btn btn-ghost" onClick={() => setMode('view')}>رجوع</button></div>
        </form>
      )}

      {mode === 'pay' && <PaymentBox config={config} doc={d} invoiceConfig={invoiceConfig} busy={busy} run={run} onCancel={() => setMode('view')} />}

      {mode === 'return' && (
        <form className="reverse-box neutral" onSubmit={(e: FormEvent<HTMLFormElement>) => {
          e.preventDefault();
          const f = new FormData(e.currentTarget);
          const lines = d.lines.map((l) => ({ sourceItemId: l.id, quantity: String(f.get(`q-${l.id}`) ?? '').trim() })).filter((l) => l.quantity && Number(l.quantity) > 0);
          void run(async () => {
            const r = await api<{ id: string }>('POST', `/api/${returnConfig.api}`, { originalInvoiceId: d.id, docDate: String(f.get('date')), reason: String(f.get('reason')), lines });
            navigate(`../../returns/${r.id}`, { relative: 'path' });
          });
        }}>
          <h3>{returnConfig.singular} على {d.number}</h3>
          <table className="table">
            <thead><tr><th>الصنف</th><th className="num">الكمية في الفاتورة</th><th>الكمية المرتجعة</th></tr></thead>
            <tbody>{d.lines.map((l) => (
              <tr key={l.id}><td>{l.description ?? l.productName}</td><td className="num" dir="ltr">{l.quantity.replace(/\.?0+$/, '')}</td>
                <td><input name={`q-${l.id}`} dir="ltr" inputMode="decimal" placeholder="0" aria-label={`كمية مرتجعة ${l.lineNo}`} /></td></tr>
            ))}</tbody>
          </table>
          <div className="grid-2">
            <label>التاريخ<input name="date" type="date" required defaultValue={today()} /></label>
            <label>سبب الإرجاع<input name="reason" required minLength={3} /></label>
          </div>
          <div className="form-actions"><button className="btn btn-primary" disabled={busy}>إنشاء مسودة المرتجع</button><button type="button" className="btn btn-ghost" onClick={() => setMode('view')}>رجوع</button></div>
        </form>
      )}

      {(d.allocations?.length || d.returns?.length) ? (
        <div className="grid-cards related">
          {d.allocations && d.allocations.length > 0 && (
            <div><h3>الدفعات</h3><table className="table"><tbody>{d.allocations.map((a) => (
              <tr key={a.id}><td dir="ltr" className="code">{a.paymentNumber}</td><td dir="ltr">{a.paymentDate}</td><td className="num" dir="ltr">{formatAmount(a.amount)}</td></tr>
            ))}</tbody></table></div>
          )}
          {d.returns && d.returns.length > 0 && (
            <div><h3>المرتجعات</h3><table className="table"><tbody>{d.returns.map((r) => (
              <tr key={r.id}><td dir="ltr" className="code"><Link to={`../../returns/${r.id}`} relative="path">{r.number ?? 'مسودة'}</Link></td><td dir="ltr">{r.date}</td>
                <td className="num" dir="ltr">{formatAmount(r.total)}</td><td><span className={`tag ${statusClass(r.status)}`}>{STATUS_AR[r.status]}</span></td></tr>
            ))}</tbody></table></div>
          )}
        </div>
      ) : null}
    </div>
  );
}

function PaymentBox({ config, doc, busy, run, onCancel }: {
  config: DocConfig; doc: Doc; invoiceConfig: DocConfig; busy: boolean; run: (fn: () => Promise<unknown>) => Promise<void>; onCancel: () => void;
}) {
  const methods = useLoad(() => api<{ data: Method[] }>('GET', '/api/payment-methods'));
  const direction = (config.party === 'customer') !== config.isReturn ? 'RECEIPT' : 'DISBURSEMENT';
  return (
    <form className="reverse-box neutral" onSubmit={(e) => {
      e.preventDefault();
      const f = new FormData(e.currentTarget);
      const amount = String(f.get('amount'));
      void run(() => api('POST', '/api/payments', {
        direction, [config.party === 'customer' ? 'customerId' : 'supplierId']: doc.partyId,
        paymentDate: String(f.get('date')), methodId: String(f.get('methodId')), amount, reference: f.get('reference') || null,
        allocations: [{ documentType: config.key, documentId: doc.id, amount }],
      }));
    }}>
      <h3>{direction === 'RECEIPT' ? 'سند قبض' : 'سند صرف'} لـ {doc.number}</h3>
      <div className="grid-3">
        <label>المبلغ<input name="amount" dir="ltr" required pattern="[0-9]{1,16}(\.[0-9]{1,2})?" defaultValue={doc.remainingAmount} /></label>
        <label>طريقة الدفع
          <select name="methodId" required>{methods.data?.data.filter((m) => m.isActive).map((m) => <option key={m.id} value={m.id}>{m.nameAr}</option>)}</select>
        </label>
        <label>التاريخ<input name="date" type="date" required defaultValue={today()} /></label>
        <label>المرجع (اختياري)<input name="reference" /></label>
      </div>
      <div className="form-actions"><button className="btn btn-primary" disabled={busy}>حفظ</button><button type="button" className="btn btn-ghost" onClick={onCancel}>رجوع</button></div>
    </form>
  );
}
