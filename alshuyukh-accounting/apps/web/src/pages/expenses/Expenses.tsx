import { useEffect, useRef, useState, type FormEvent } from 'react';
import { Link, NavLink, Navigate, Route, Routes, useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { formatAmount } from '../../money';
import { ErrorBox, PageHeader, useLoad } from '../../ui';
import { STATUS_AR, statusClass } from '../documents/config';

export interface Category { id: string; code: string; nameAr: string; accountId: string; accountCode: string; accountName: string; vatCategory: string; isActive: boolean }
interface Method { id: string; code: string; nameAr: string; isActive: boolean; methodType?: string }
interface Party { id: string; code: string; nameAr: string }
interface CostCenter { id: string; code: string; name: string; isActive: boolean }
interface ExpenseLine { id: string; categoryId: string; categoryName: string; accountCode: string; description: string | null; costCenterId: string | null; amount: string; netAmount: string; vatCategory: string; vatAmount: string; totalAmount: string }
export interface Expense {
  id: string; number: string | null; date: string; paymentType: 'CASH' | 'BANK' | 'CREDIT'; methodId: string | null; methodName: string | null;
  supplierId: string | null; supplierName: string | null; payeeName: string | null; reference: string | null; vendorVatNumber: string | null;
  pricesIncludeVat: boolean; subtotal: string; taxAmount: string; total: string; paidAmount: string; remainingAmount: string;
  notes: string | null; status: string; journalEntryId: string | null; cancelReason: string | null; lines: ExpenseLine[];
}

const today = () => new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh' }).format(new Date());
export const VAT_AR: Record<string, string> = { S: 'خاضع', Z: 'صفري', E: 'معفى', O: 'خارج النطاق' };
const PAY_AR: Record<string, string> = { CASH: 'نقدًا', BANK: 'تحويل بنكي', CREDIT: 'آجل على مورد' };
const tab = ({ isActive }: { isActive: boolean }) => `tab ${isActive ? 'active' : ''}`;

export default function Expenses() {
  const { can } = useAuth();
  return (
    <>
      <PageHeader title="المصروفات" />
      <nav className="tabs">
        <NavLink to="/expenses/list" className={tab}>سندات المصروفات</NavLink>
        <NavLink to="/expenses/categories" className={tab}>فئات المصروفات</NavLink>
      </nav>
      <Routes>
        <Route index element={<Navigate to="list" replace />} />
        <Route path="list" element={<ExpenseList />} />
        {can('expense.create') && <Route path="list/new" element={<ExpenseForm />} />}
        <Route path="list/:id" element={<ExpenseDetail />} />
        <Route path="categories" element={<Categories />} />
      </Routes>
    </>
  );
}

function ExpenseList() {
  const { can } = useAuth();
  const [status, setStatus] = useState('');
  const list = useLoad(() => {
    const q = new URLSearchParams({ limit: '100' });
    if (status === 'OPEN') q.set('open', 'true'); else if (status) q.set('status', status);
    return api<{ data: Expense[]; total: number }>('GET', `/api/expenses?${q}`);
  }, [status]);
  return (
    <div className="card">
      <div className="toolbar">
        <select value={status} onChange={(e) => setStatus(e.target.value)} aria-label="الحالة">
          <option value="">كل الحالات</option><option value="OPEN">غير مسددة</option>
          {['DRAFT', 'POSTED', 'PARTIALLY_PAID', 'PAID', 'CANCELLED'].map((s) => <option key={s} value={s}>{STATUS_AR[s]}</option>)}
        </select>
        <span className="grow" />
        {can('expense.create') && <Link className="btn btn-primary" to="/expenses/list/new">مصروف جديد</Link>}
      </div>
      <ErrorBox error={list.error} />
      <div className="table-wrap">
        <table className="table">
          <thead><tr><th>الرقم</th><th>التاريخ</th><th>المستفيد</th><th>الدفع</th><th className="num">الضريبة</th><th className="num">الإجمالي</th><th className="num">المتبقي</th><th>الحالة</th></tr></thead>
          <tbody>
            {!list.data && !list.error && <tr><td colSpan={8} className="muted empty-row">جارٍ التحميل…</td></tr>}
            {list.data?.data.map((e) => (
              <tr key={e.id}>
                <td dir="ltr" className="code"><Link to={`/expenses/list/${e.id}`}>{e.number ?? 'مسودة'}</Link></td>
                <td dir="ltr">{e.date}</td>
                <td>{e.supplierName ?? e.payeeName ?? '—'}</td>
                <td>{PAY_AR[e.paymentType]}{e.methodName ? ` — ${e.methodName}` : ''}</td>
                <td className="num" dir="ltr">{formatAmount(e.taxAmount)}</td>
                <td className="num" dir="ltr">{formatAmount(e.total)}</td>
                <td className="num" dir="ltr">{['DRAFT', 'CANCELLED'].includes(e.status) ? '' : formatAmount(e.remainingAmount)}</td>
                <td><span className={`tag ${statusClass(e.status)}`}>{STATUS_AR[e.status] ?? e.status}</span></td>
              </tr>
            ))}
            {list.data?.data.length === 0 && <tr><td colSpan={8} className="muted empty-row">لا توجد مصروفات</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  );
}

interface FormLine { categoryId: string; amount: string; description: string; vatCategory: string; costCenterId: string }
const emptyLine = (): FormLine => ({ categoryId: '', amount: '', description: '', vatCategory: '', costCenterId: '' });
const payloadLines = (ls: FormLine[]) => ls.filter((l) => l.categoryId && l.amount.trim()).map((l) => ({
  categoryId: l.categoryId, amount: l.amount.trim(),
  ...(l.vatCategory ? { vatCategory: l.vatCategory } : {}),
  ...(l.description.trim() ? { description: l.description.trim() } : {}),
  ...(l.costCenterId ? { costCenterId: l.costCenterId } : {}),
}));

function ExpenseForm({ expense, onSaved }: { expense?: Expense; onSaved?: () => void }) {
  const navigate = useNavigate();
  const { can } = useAuth();
  const cats = useLoad(() => api<{ data: Category[] }>('GET', '/api/expense-categories'));
  const methods = useLoad(() => (can('payment.view') ? api<{ data: Method[] }>('GET', '/api/payment-methods') : Promise.resolve({ data: [] as Method[] })));
  const suppliers = useLoad(() => api<{ data: Party[] }>('GET', '/api/suppliers?limit=200').catch(() => ({ data: [] as Party[] })));
  const centers = useLoad(() => api<{ data: CostCenter[] }>('GET', '/api/cost-centers').catch(() => ({ data: [] as CostCenter[] })));
  const [date, setDate] = useState(expense?.date ?? today());
  const [paymentType, setPaymentType] = useState<string>(expense?.paymentType ?? 'CASH');
  const [methodId, setMethodId] = useState(expense?.methodId ?? '');
  const [supplierId, setSupplierId] = useState(expense?.supplierId ?? '');
  const [includeVat, setIncludeVat] = useState(expense?.pricesIncludeVat ?? false);
  const [lines, setLines] = useState<FormLine[]>(() => expense?.lines.length
    ? expense.lines.map((l) => ({ categoryId: l.categoryId, amount: l.amount, description: l.description ?? '', vatCategory: l.vatCategory, costCenterId: l.costCenterId ?? '' }))
    : [emptyLine()]);
  const [preview, setPreview] = useState<{ lines: { netAmount: string; vatAmount: string }[]; totals: { subtotal: string; taxAmount: string; total: string } } | null>(null);
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState(false);
  const timer = useRef<ReturnType<typeof setTimeout>>(undefined);
  const update = (i: number, patch: Partial<FormLine>) => setLines((ls) => ls.map((l, j) => (j === i ? { ...l, ...patch } : l)));

  useEffect(() => {
    clearTimeout(timer.current);
    const body = payloadLines(lines);
    if (!body.length) { setPreview(null); return; }
    timer.current = setTimeout(() => {
      api<typeof preview>('POST', '/api/expenses/calculate', { expenseDate: date, pricesIncludeVat: includeVat, lines: body })
        .then(setPreview).catch(() => setPreview(null));
    }, 350);
    return () => clearTimeout(timer.current);
  }, [lines, date, includeVat]);

  // Cash expenses default to a cash method and bank expenses to a bank method.
  const methodOptions = methods.data?.data.filter((m) => m.isActive) ?? [];
  useEffect(() => {
    if (paymentType === 'CREDIT' || methodId || !methodOptions.length) return;
    const match = methodOptions.find((m) => m.code === paymentType) ?? methodOptions[0];
    if (match) setMethodId(match.id);
  }, [paymentType, methodOptions, methodId]);

  async function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const opt = (k: string) => (String(f.get(k) ?? '').trim() || null);
    const post = (e.nativeEvent as SubmitEvent).submitter?.getAttribute('value') === 'post';
    const body = {
      expenseDate: date, paymentType, pricesIncludeVat: includeVat,
      methodId: paymentType === 'CREDIT' ? null : methodId || null,
      supplierId: supplierId || null,
      payeeName: opt('payeeName'), reference: opt('reference'), vendorVatNumber: opt('vendorVatNumber'), notes: opt('notes'),
      lines: payloadLines(lines),
    };
    setBusy(true); setError(null);
    try {
      if (expense) {
        await api('PATCH', `/api/expenses/${expense.id}`, body);
        if (post) await api('POST', `/api/expenses/${expense.id}/post`);
        onSaved?.();
      } else {
        const created = await api<{ id: string }>('POST', '/api/expenses', { ...body, post });
        navigate(`/expenses/list/${created.id}`);
      }
    } catch (err) { setError(err); } finally { setBusy(false); }
  }

  const catById = new Map(cats.data?.data.map((c) => [c.id, c]));
  const ready = cats.data && methods.data;
  return (
    <form className="card form" onSubmit={submit}>
      <h2>{expense ? 'تعديل مسودة المصروف' : 'مصروف جديد'}</h2>
      <ErrorBox error={error ?? cats.error} />
      <div className="grid-3">
        <label>التاريخ<input type="date" required value={date} onChange={(e) => setDate(e.target.value)} /></label>
        <label>طريقة السداد
          <select value={paymentType} onChange={(e) => { setPaymentType(e.target.value); setMethodId(''); }}>
            <option value="CASH">نقدًا</option><option value="BANK">تحويل بنكي</option><option value="CREDIT">آجل على مورد</option>
          </select>
        </label>
        {paymentType !== 'CREDIT' && (
          <label>الصندوق / البنك
            <select required value={methodId} onChange={(e) => setMethodId(e.target.value)}>
              <option value="">— اختر —</option>
              {methodOptions.map((m) => <option key={m.id} value={m.id}>{m.nameAr}</option>)}
            </select>
          </label>
        )}
        <label>المورد{paymentType === 'CREDIT' ? '' : ' (اختياري)'}
          <select required={paymentType === 'CREDIT'} value={supplierId} onChange={(e) => setSupplierId(e.target.value)}>
            <option value="">— بدون —</option>
            {suppliers.data?.data.map((s) => <option key={s.id} value={s.id}>{s.code} — {s.nameAr}</option>)}
          </select>
        </label>
        {!supplierId && <label>اسم المستفيد<input name="payeeName" defaultValue={expense?.payeeName ?? ''} /></label>}
        <label>رقم فاتورة المورد / المرجع<input name="reference" dir="ltr" defaultValue={expense?.reference ?? ''} /></label>
        <label>الرقم الضريبي للمورد<input name="vendorVatNumber" dir="ltr" pattern="3[0-9]{13}3" title="15 رقمًا يبدأ وينتهي بالرقم 3" defaultValue={expense?.vendorVatNumber ?? ''} /></label>
      </div>
      <label className="check"><input type="checkbox" checked={includeVat} onChange={(e) => setIncludeVat(e.target.checked)} />المبالغ شاملة ضريبة القيمة المضافة</label>

      <div className="table-wrap">
        <table className="table lines-table">
          <thead><tr><th>الفئة</th><th>المبلغ</th><th>الضريبة</th>{(centers.data?.data.length ?? 0) > 0 && <th>مركز التكلفة</th>}<th className="num">الصافي</th><th className="num">ضريبة المدخلات</th><th /></tr></thead>
          <tbody>
            {lines.map((l, i) => {
              const pl = preview?.lines[payloadLines(lines.slice(0, i + 1)).length - 1];
              const ok = l.categoryId && l.amount.trim();
              return (
                <tr key={i}>
                  <td>
                    <select value={l.categoryId} aria-label={`فئة السطر ${i + 1}`} onChange={(e) => update(i, { categoryId: e.target.value, vatCategory: '' })}>
                      <option value="">— اختر الفئة —</option>
                      {cats.data?.data.filter((c) => c.isActive).map((c) => <option key={c.id} value={c.id}>{c.nameAr}</option>)}
                    </select>
                    <input className="sub-row" placeholder="الوصف (اختياري)" value={l.description} onChange={(e) => update(i, { description: e.target.value })} />
                  </td>
                  <td><input dir="ltr" inputMode="decimal" value={l.amount} aria-label={`مبلغ السطر ${i + 1}`} onChange={(e) => update(i, { amount: e.target.value })} /></td>
                  <td>
                    <select value={l.vatCategory} aria-label={`ضريبة السطر ${i + 1}`} onChange={(e) => update(i, { vatCategory: e.target.value })}>
                      <option value="">{l.categoryId ? VAT_AR[catById.get(l.categoryId)?.vatCategory ?? 'S'] : 'افتراضي'}</option>
                      {Object.entries(VAT_AR).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
                    </select>
                  </td>
                  {(centers.data?.data.length ?? 0) > 0 && (
                    <td><select value={l.costCenterId} aria-label="مركز التكلفة" onChange={(e) => update(i, { costCenterId: e.target.value })}>
                      <option value="">—</option>
                      {centers.data!.data.filter((c) => c.isActive).map((c) => <option key={c.id} value={c.id}>{c.code} — {c.name}</option>)}
                    </select></td>
                  )}
                  <td className="num" dir="ltr">{ok && pl ? formatAmount(pl.netAmount) : ''}</td>
                  <td className="num" dir="ltr">{ok && pl ? formatAmount(pl.vatAmount) : ''}</td>
                  <td>{lines.length > 1 && <button type="button" className="icon-btn" aria-label="حذف السطر" onClick={() => setLines((ls) => ls.filter((_, j) => j !== i))}>×</button>}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
      <button type="button" className="btn btn-small" onClick={() => setLines((ls) => [...ls, emptyLine()])}>+ سطر</button>

      <div className="doc-footer">
        <label className="grow">ملاحظات<input name="notes" defaultValue={expense?.notes ?? ''} /></label>
        <div>
          {preview && (
            <dl className="totals">
              <div><dt>الإجمالي قبل الضريبة</dt><dd dir="ltr">{formatAmount(preview.totals.subtotal)}</dd></div>
              <div><dt>ضريبة المدخلات</dt><dd dir="ltr">{formatAmount(preview.totals.taxAmount)}</dd></div>
              <div className="strong"><dt>الإجمالي</dt><dd dir="ltr">{formatAmount(preview.totals.total)}</dd></div>
            </dl>
          )}
          <p className="muted small">تُحسب المبالغ والضريبة في الخادم.</p>
        </div>
      </div>
      <div className="form-actions">
        <button className="btn" value="draft" disabled={busy || !ready}>حفظ كمسودة</button>
        {can('expense.post') && <button className="btn btn-primary" value="post" disabled={busy || !ready}>حفظ وترحيل</button>}
        {expense && <button type="button" className="btn btn-ghost" onClick={onSaved}>إلغاء</button>}
      </div>
    </form>
  );
}

function ExpenseDetail() {
  const { id } = useParams();
  const { can } = useAuth();
  const navigate = useNavigate();
  const exp = useLoad(() => api<Expense>('GET', `/api/expenses/${id}`), [id]);
  const [mode, setMode] = useState<'view' | 'edit' | 'cancel'>('view');
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState(false);
  async function run(fn: () => Promise<unknown>) {
    setError(null); setBusy(true);
    try { await fn(); setMode('view'); exp.reload(); } catch (e) { setError(e); } finally { setBusy(false); }
  }
  const e = exp.data;
  if (!e) return <ErrorBox error={exp.error} />;
  if (mode === 'edit') return <ExpenseForm expense={e} onSaved={() => { setMode('view'); exp.reload(); }} />;
  const posted = !['DRAFT', 'CANCELLED'].includes(e.status);

  return (
    <div className="card">
      <ErrorBox error={error} />
      <div className="entry-head">
        <div>
          <h2 dir="ltr" className="code">{e.number ?? 'مسودة'}</h2>
          <p>سند مصروف — {e.supplierId ? <Link to="/suppliers">{e.supplierName}</Link> : e.payeeName ?? 'بدون مستفيد'}</p>
        </div>
        <span className={`tag ${statusClass(e.status)}`}>{STATUS_AR[e.status] ?? e.status}</span>
      </div>
      <dl className="meta">
        <div><dt>التاريخ</dt><dd dir="ltr">{e.date}</dd></div>
        <div><dt>السداد</dt><dd>{PAY_AR[e.paymentType]}{e.methodName ? ` — ${e.methodName}` : ''}</dd></div>
        {e.reference && <div><dt>المرجع</dt><dd dir="ltr">{e.reference}</dd></div>}
        {e.vendorVatNumber && <div><dt>الرقم الضريبي للمورد</dt><dd dir="ltr">{e.vendorVatNumber}</dd></div>}
        {e.journalEntryId && <div><dt>القيد</dt><dd><Link to={`/accounting/journal/${e.journalEntryId}`}>عرض القيد المحاسبي</Link></dd></div>}
        {e.cancelReason && <div><dt>سبب الإلغاء</dt><dd>{e.cancelReason}</dd></div>}
        <div><dt>المبالغ</dt><dd>{e.pricesIncludeVat ? 'شاملة الضريبة' : 'غير شاملة الضريبة'}</dd></div>
      </dl>
      <div className="table-wrap">
        <table className="table">
          <thead><tr><th>الفئة</th><th>الحساب</th><th>الوصف</th><th>الضريبة</th><th className="num">الصافي</th><th className="num">الضريبة</th><th className="num">الإجمالي</th></tr></thead>
          <tbody>
            {e.lines.map((l) => (
              <tr key={l.id}>
                <td>{l.categoryName}</td><td dir="ltr" className="code">{l.accountCode}</td><td>{l.description ?? ''}</td><td>{VAT_AR[l.vatCategory]}</td>
                <td className="num" dir="ltr">{formatAmount(l.netAmount)}</td><td className="num" dir="ltr">{formatAmount(l.vatAmount)}</td><td className="num" dir="ltr">{formatAmount(l.totalAmount)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <dl className="totals">
        <div><dt>الإجمالي قبل الضريبة</dt><dd dir="ltr">{formatAmount(e.subtotal)}</dd></div>
        <div><dt>ضريبة المدخلات</dt><dd dir="ltr">{formatAmount(e.taxAmount)}</dd></div>
        <div className="strong"><dt>الإجمالي</dt><dd dir="ltr">{formatAmount(e.total)}</dd></div>
        {posted && <div><dt>المدفوع</dt><dd dir="ltr">{formatAmount(e.paidAmount)}</dd></div>}
        {posted && <div><dt>المتبقي</dt><dd dir="ltr">{formatAmount(e.remainingAmount)}</dd></div>}
      </dl>
      {e.notes && <p className="muted">{e.notes}</p>}
      {e.paymentType === 'CREDIT' && posted && Number(e.remainingAmount) > 0 && (
        <p className="muted small">يُسدَّد المصروف الآجل من <Link to="/purchases/payments">سندات الصرف</Link> بتخصيص الدفعة عليه.</p>
      )}

      {mode === 'cancel' && (
        <form className="inline-form" onSubmit={(ev) => { ev.preventDefault(); const reason = String(new FormData(ev.currentTarget).get('reason')); void run(() => api('POST', `/api/expenses/${e.id}/cancel`, { reason })); }}>
          <input name="reason" required minLength={3} placeholder="سبب الإلغاء" />
          <button className="btn btn-danger" disabled={busy}>تأكيد الإلغاء بقيد عكسي</button>
          <button type="button" className="btn btn-ghost" onClick={() => setMode('view')}>تراجع</button>
        </form>
      )}
      {mode === 'view' && (
        <div className="form-actions">
          {e.status === 'DRAFT' && can('expense.create') && <button className="btn" onClick={() => setMode('edit')}>تعديل</button>}
          {e.status === 'DRAFT' && can('expense.post') && <button className="btn btn-primary" disabled={busy} onClick={() => run(() => api('POST', `/api/expenses/${e.id}/post`))}>ترحيل</button>}
          {e.status === 'DRAFT' && can('expense.create') && (
            <button className="btn btn-danger" disabled={busy} onClick={() => run(async () => { await api('DELETE', `/api/expenses/${e.id}`); navigate('/expenses/list'); })}>حذف المسودة</button>
          )}
          {posted && can('expense.cancel') && <button className="btn btn-danger" onClick={() => setMode('cancel')}>إلغاء المصروف</button>}
        </div>
      )}
    </div>
  );
}

function Categories() {
  const { can } = useAuth();
  const cats = useLoad(() => api<{ data: Category[] }>('GET', '/api/expense-categories'));
  const manage = can('account.manage');
  const accounts = useLoad(() => (manage ? api<{ data: { id: string; code: string; nameAr: string; type: string }[] }>('GET', '/api/accounts?postableOnly=true') : Promise.resolve({ data: [] })), [manage]);
  const [error, setError] = useState<unknown>(null);

  async function add(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const f = Object.fromEntries(new FormData(form)) as Record<string, string>;
    setError(null);
    try { await api('POST', '/api/expense-categories', f); form.reset(); cats.reload(); } catch (err) { setError(err); }
  }
  async function toggle(c: Category) {
    setError(null);
    try { await api('PATCH', `/api/expense-categories/${c.id}`, { isActive: !c.isActive }); cats.reload(); } catch (err) { setError(err); }
  }

  return (
    <div className="card">
      <ErrorBox error={error ?? cats.error} />
      <p className="muted small">كل فئة مرتبطة بحساب مصروف في دليل الحسابات وبمعاملة ضريبية افتراضية يمكن تغييرها في كل سطر.</p>
      <div className="table-wrap">
        <table className="table">
          <thead><tr><th>الرمز</th><th>الفئة</th><th>الحساب</th><th>الضريبة الافتراضية</th><th>الحالة</th>{manage && <th />}</tr></thead>
          <tbody>
            {cats.data?.data.map((c) => (
              <tr key={c.id}>
                <td dir="ltr" className="code">{c.code}</td><td>{c.nameAr}</td><td><span dir="ltr" className="code">{c.accountCode}</span> {c.accountName}</td>
                <td>{VAT_AR[c.vatCategory]}</td><td>{c.isActive ? 'نشطة' : 'موقوفة'}</td>
                {manage && <td><button className="btn btn-small" onClick={() => toggle(c)}>{c.isActive ? 'إيقاف' : 'تفعيل'}</button></td>}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {manage && accounts.data && (
        <form className="form" onSubmit={add}>
          <h3>فئة جديدة</h3>
          <div className="grid-3">
            <label>الرمز<input name="code" required dir="ltr" pattern="[A-Za-z][A-Za-z0-9_]{1,30}" placeholder="TRAVEL" /></label>
            <label>الاسم<input name="nameAr" required minLength={2} /></label>
            <label>الحساب
              <select name="accountId" required>
                <option value="">— اختر —</option>
                {accounts.data.data.filter((a) => a.type === 'EXPENSE' || a.type === 'COST_OF_GOODS_SOLD').map((a) => <option key={a.id} value={a.id}>{a.code} — {a.nameAr}</option>)}
              </select>
            </label>
            <label>الضريبة الافتراضية
              <select name="vatCategory" defaultValue="S">{Object.entries(VAT_AR).map(([k, v]) => <option key={k} value={k}>{v}</option>)}</select>
            </label>
          </div>
          <button className="btn btn-primary">إضافة</button>
        </form>
      )}
    </div>
  );
}
