import { useState, type FormEvent } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { STATUS_AR, statusClass, type Doc } from './config';

interface Payment {
  id: string; number: string; direction: string; date: string; partyName: string; methodName: string; amount: string;
  allocatedAmount: string; unallocatedAmount: string; status: string; reference: string | null; journalEntryId: string;
  allocations?: { id: string; amount: string; documentNumber: string; reversedAt: string | null }[];
}
interface Party { id: string; code: string; nameAr: string }
interface Method { id: string; nameAr: string; isActive: boolean }

const today = () => new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh' }).format(new Date());

/** Receipts from customers (sales) or payments to suppliers (purchases). */
export default function Payments({ party }: { party: 'customer' | 'supplier' }) {
  const { can } = useAuth();
  const direction = party === 'customer' ? 'RECEIPT' : 'DISBURSEMENT';
  const list = useLoad(() => api<{ data: Payment[]; total: number }>('GET', `/api/payments?direction=${direction}&limit=100`), [direction]);
  const [creating, setCreating] = useState(false);
  const [voiding, setVoiding] = useState<string | null>(null);
  const [error, setError] = useState<unknown>(null);

  async function voidPayment(e: FormEvent<HTMLFormElement>, id: string) {
    e.preventDefault();
    setError(null);
    try {
      await api('POST', `/api/payments/${id}/void`, { reason: String(new FormData(e.currentTarget).get('reason')) });
      setVoiding(null);
      list.reload();
    } catch (err) { setError(err); }
  }

  return (
    <>
      <ErrorBox error={error ?? list.error} />
      {creating && <NewPayment party={party} onDone={() => { setCreating(false); list.reload(); }} />}
      <div className="card">
        <div className="toolbar">
          <span className="spacer" />
          {can('payment.create') && !creating && <button className="btn btn-primary" onClick={() => setCreating(true)}>{party === 'customer' ? 'سند قبض جديد' : 'سند صرف جديد'}</button>}
        </div>
        <div className="table-wrap">
          <table className="table">
            <thead><tr><th>الرقم</th><th>التاريخ</th><th>{party === 'customer' ? 'العميل' : 'المورد'}</th><th>الطريقة</th><th className="num">المبلغ</th><th className="num">غير مخصص</th><th>الحالة</th><th /></tr></thead>
            <tbody>
              {!list.data && !list.error && <tr><td colSpan={8} className="muted empty-row">جارٍ التحميل…</td></tr>}
              {list.data?.data.map((p) => (
                <tr key={p.id}>
                  <td dir="ltr" className="code"><Link to={`/accounting/journal/${p.journalEntryId}`}>{p.number}</Link></td>
                  <td dir="ltr">{p.date}</td><td>{p.partyName}</td><td>{p.methodName}</td>
                  <td className="num" dir="ltr">{formatAmount(p.amount)}</td>
                  <td className="num" dir="ltr">{p.status === 'VOIDED' ? '' : formatAmount(p.unallocatedAmount)}</td>
                  <td><span className={`tag ${statusClass(p.status === 'POSTED' ? 'PAID' : p.status)}`}>{p.status === 'POSTED' ? 'مرحّل' : STATUS_AR[p.status]}</span></td>
                  <td>
                    {p.status === 'POSTED' && can('payment.void') && (voiding === p.id
                      ? <form className="inline-form" onSubmit={(e) => voidPayment(e, p.id)}><input name="reason" required minLength={3} placeholder="سبب الإلغاء" /><button className="btn btn-danger btn-small">تأكيد</button></form>
                      : <button className="btn btn-small" onClick={() => setVoiding(p.id)}>إلغاء</button>)}
                  </td>
                </tr>
              ))}
              {list.data?.data.length === 0 && <tr><td colSpan={8} className="muted empty-row">لا توجد سندات</td></tr>}
            </tbody>
          </table>
        </div>
      </div>
    </>
  );
}

function NewPayment({ party, onDone }: { party: 'customer' | 'supplier'; onDone: () => void }) {
  const { can } = useAuth();
  const parties = useLoad(() => api<{ data: Party[] }>('GET', `/api/${party}s?limit=200`));
  const methods = useLoad(() => api<{ data: Method[] }>('GET', '/api/payment-methods'));
  const [partyId, setPartyId] = useState('');
  const docPath = party === 'customer' ? 'invoices' : 'purchase-invoices';
  // Supplier payments can also settle credit expenses.
  const openDocs = useLoad(async () => {
    if (!partyId) return { data: [] as (Doc & { documentType: string })[] };
    const docs = (await api<{ data: Doc[] }>('GET', `/api/${docPath}?partyId=${partyId}&open=true&limit=200`)).data
      .map((d) => ({ ...d, documentType: party === 'customer' ? 'SALES_INVOICE' : 'PURCHASE_INVOICE' }));
    if (party === 'supplier' && can('expense.view')) {
      const exp = (await api<{ data: Doc[] }>('GET', `/api/expenses?supplierId=${partyId}&open=true&limit=200`)).data;
      docs.push(...exp.map((d) => ({ ...d, documentType: 'EXPENSE' })));
    }
    return { data: docs };
  }, [partyId]);
  const [alloc, setAlloc] = useState<Record<string, string>>({});
  const [error, setError] = useState<unknown>(null);

  async function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    setError(null);
    try {
      await api('POST', '/api/payments', {
        direction: party === 'customer' ? 'RECEIPT' : 'DISBURSEMENT',
        [party === 'customer' ? 'customerId' : 'supplierId']: partyId,
        paymentDate: String(f.get('date')), methodId: String(f.get('methodId')), amount: String(f.get('amount')),
        reference: f.get('reference') || null,
        allocations: Object.entries(alloc).filter(([, v]) => v.trim() && Number(v) > 0).map(([documentId, amount]) => ({
          documentType: openDocs.data?.data.find((d) => d.id === documentId)?.documentType, documentId, amount: amount.trim(),
        })),
      });
      onDone();
    } catch (err) { setError(err); }
  }

  return (
    <form className="card form" onSubmit={submit}>
      <h2>{party === 'customer' ? 'سند قبض جديد' : 'سند صرف جديد'}</h2>
      <ErrorBox error={error} />
      <div className="grid-3">
        <label>{party === 'customer' ? 'العميل' : 'المورد'}
          <select required value={partyId} onChange={(e) => { setPartyId(e.target.value); setAlloc({}); }}>
            <option value="">— اختر —</option>
            {parties.data?.data.map((p) => <option key={p.id} value={p.id}>{p.code} — {p.nameAr}</option>)}
          </select>
        </label>
        <label>المبلغ<input name="amount" dir="ltr" required pattern="[0-9]{1,16}(\.[0-9]{1,2})?" /></label>
        <label>طريقة الدفع
          <select name="methodId" required>{methods.data?.data.filter((m) => m.isActive).map((m) => <option key={m.id} value={m.id}>{m.nameAr}</option>)}</select>
        </label>
        <label>التاريخ<input name="date" type="date" required defaultValue={today()} /></label>
        <label>المرجع (اختياري)<input name="reference" /></label>
      </div>
      {partyId && (
        <>
          <h3>تخصيص على الفواتير المفتوحة (اختياري — الباقي يبقى رصيدًا مقدمًا)</h3>
          <table className="table">
            <thead><tr><th>الفاتورة</th><th>التاريخ</th><th className="num">المتبقي</th><th>المبلغ المخصص</th></tr></thead>
            <tbody>
              {openDocs.data?.data.map((d) => (
                <tr key={d.id}>
                  <td dir="ltr" className="code">{d.number}{d.documentType === 'EXPENSE' && <span className="muted small"> (مصروف)</span>}</td><td dir="ltr">{d.date}</td>
                  <td className="num" dir="ltr">{formatAmount(d.remainingAmount)}</td>
                  <td><input dir="ltr" inputMode="decimal" value={alloc[d.id] ?? ''} aria-label={`تخصيص ${d.number}`}
                    onChange={(e) => setAlloc((a) => ({ ...a, [d.id]: e.target.value }))} /></td>
                </tr>
              ))}
              {openDocs.data?.data.length === 0 && <tr><td colSpan={4} className="muted empty-row">لا توجد فواتير مفتوحة</td></tr>}
            </tbody>
          </table>
        </>
      )}
      <div className="form-actions"><button className="btn btn-primary">حفظ وترحيل</button><button type="button" className="btn btn-ghost" onClick={onDone}>إلغاء</button></div>
    </form>
  );
}
