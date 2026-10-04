import { useEffect, useRef, useState, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import type { Account } from '../accounting/ChartOfAccounts';
import type { DocConfig, Doc } from './config';
import Totals, { type TotalsData } from './Totals';
import { loadWarehouses } from '../inventory/shared';

interface Product { id: string; sku: string; nameAr: string; salePrice: string; purchasePrice: string; productType: string; vatCategory: string }
interface Party { id: string; code: string; nameAr: string }
interface Line { productId: string; accountId: string; description: string; quantity: string; unitPrice: string; discountPercent: string; vatCategory: string }

const empty = (): Line => ({ productId: '', accountId: '', description: '', quantity: '1', unitPrice: '', discountPercent: '', vatCategory: 'S' });
const today = () => new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh' }).format(new Date());

/** Builds the API payload from the form lines; blank optional fields are left out so the server fills defaults. */
function payloadLines(lines: Line[]) {
  return lines.filter((l) => l.productId || l.accountId).map((l) => ({
    ...(l.productId ? { productId: l.productId } : { accountId: l.accountId, vatCategory: l.vatCategory }),
    quantity: l.quantity.trim(),
    ...(l.unitPrice.trim() ? { unitPrice: l.unitPrice.trim() } : {}),
    ...(l.discountPercent.trim() ? { discountPercent: l.discountPercent.trim() } : {}),
    ...(l.description.trim() ? { description: l.description.trim() } : {}),
  }));
}

export default function DocumentForm({ config, doc, onSaved }: { config: DocConfig; doc?: Doc; onSaved?: () => void }) {
  const navigate = useNavigate();
  const isPurchase = config.party === 'supplier';
  const parties = useLoad(() => api<{ data: Party[] }>('GET', `/api/${config.party}s?limit=200`));
  const products = useLoad(() => api<{ data: Product[] }>('GET', '/api/products?limit=200'));
  const accounts = useLoad(() => (isPurchase ? api<{ data: Account[] }>('GET', '/api/accounts?postableOnly=true') : Promise.resolve({ data: [] as Account[] })));
  const warehouses = useLoad(loadWarehouses);
  const [partyId, setPartyId] = useState(doc?.partyId ?? '');
  const [docDate, setDocDate] = useState(doc?.date ?? today());
  const [includeVat, setIncludeVat] = useState(doc?.pricesIncludeVat ?? false);
  const [lines, setLines] = useState<Line[]>(() => doc?.lines.length
    ? doc.lines.map((l) => ({
        productId: l.productId ?? '', accountId: l.productId ? '' : l.accountId ?? '', description: l.description ?? '',
        quantity: l.quantity, unitPrice: l.unitPrice, discountPercent: '', vatCategory: l.vatCategory,
      }))
    : [empty()]);
  // Existing discounts are kept as amounts when editing.
  const [discountAmounts] = useState<string[]>(() => doc?.lines.map((l) => l.discountBasis) ?? []);
  const [preview, setPreview] = useState<{ totals: TotalsData; lines: { netAmount: string; vatAmount: string; totalAmount: string }[] } | null>(null);
  const [previewError, setPreviewError] = useState<unknown>(null);
  const [error, setError] = useState<unknown>(null);
  const [busy, setBusy] = useState(false);
  const timer = useRef<ReturnType<typeof setTimeout>>(undefined);

  const update = (i: number, patch: Partial<Line>) => setLines((ls) => ls.map((l, j) => (j === i ? { ...l, ...patch } : l)));
  const withDiscounts = (ls: ReturnType<typeof payloadLines>) =>
    ls.map((l, i) => (!('discountPercent' in l) && discountAmounts[i] && discountAmounts[i] !== '0.00' ? { ...l, discountAmount: discountAmounts[i] } : l));

  // Totals always come from the server (POST …/calculate), debounced while typing.
  useEffect(() => {
    clearTimeout(timer.current);
    const body = withDiscounts(payloadLines(lines));
    if (!body.length) { setPreview(null); return; }
    timer.current = setTimeout(() => {
      api<{ totals: TotalsData; lines: { netAmount: string; vatAmount: string; totalAmount: string }[] }>(
        'POST', `/api/${config.api}/calculate`, { docDate, pricesIncludeVat: includeVat, lines: body })
        .then((r) => { setPreview(r); setPreviewError(null); })
        .catch((e) => { setPreview(null); setPreviewError(e); });
    }, 350);
    return () => clearTimeout(timer.current);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lines, docDate, includeVat, config.api]);

  async function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const opt = (k: string) => (f.get(k) ? String(f.get(k)) : null);
    const body: Record<string, unknown> = {
      partyId, docDate, pricesIncludeVat: includeVat, notes: opt('notes'), lines: withDiscounts(payloadLines(lines)),
    };
    if (config.key === 'SALES_INVOICE' || config.key === 'PURCHASE_INVOICE') body.dueDate = opt('dueDate');
    if (config.key === 'SALES_QUOTE') body.validUntil = opt('validUntil');
    if (config.key === 'PURCHASE_ORDER') body.expectedDate = opt('expectedDate');
    if (config.key === 'PURCHASE_INVOICE') body.supplierInvoiceNumber = opt('supplierInvoiceNumber');
    if (opt('warehouseId')) body.warehouseId = opt('warehouseId');
    setBusy(true);
    setError(null);
    try {
      if (doc) {
        await api('PATCH', `/api/${config.api}/${doc.id}`, body);
        onSaved?.();
      } else {
        const created = await api<{ id: string }>('POST', `/api/${config.api}`, body);
        navigate(`../${created.id}`, { relative: 'path' });
      }
    } catch (err) { setError(err); } finally { setBusy(false); }
  }

  const productById = new Map(products.data?.data.map((p) => [p.id, p]));
  return (
    <form className="card form" onSubmit={submit}>
      <h2>{doc ? `تعديل ${config.singular}` : config.newLabel}</h2>
      <ErrorBox error={error ?? parties.error ?? products.error} />
      <div className="grid-3">
        <label>{isPurchase ? 'المورد' : 'العميل'}
          <select required value={partyId} onChange={(e) => setPartyId(e.target.value)}>
            <option value="">— اختر —</option>
            {parties.data?.data.map((p) => <option key={p.id} value={p.id}>{p.code} — {p.nameAr}</option>)}
          </select>
        </label>
        <label>التاريخ<input type="date" required value={docDate} onChange={(e) => setDocDate(e.target.value)} /></label>
        {(config.key === 'SALES_INVOICE' || config.key === 'PURCHASE_INVOICE') && (
          <label>تاريخ الاستحقاق (افتراضيًا حسب مدة السداد)<input name="dueDate" type="date" defaultValue={doc?.dueDate ?? ''} /></label>
        )}
        {config.key === 'SALES_QUOTE' && <label>صالح حتى<input name="validUntil" type="date" defaultValue={doc?.validUntil ?? ''} /></label>}
        {config.key === 'PURCHASE_ORDER' && <label>تاريخ التوريد المتوقع<input name="expectedDate" type="date" /></label>}
        {(warehouses.data?.length ?? 0) > 1 && (
          <label>المستودع
            <select name="warehouseId" defaultValue={(doc as { warehouseId?: string } | undefined)?.warehouseId ?? ''}>
              <option value="">المستودع الرئيسي</option>
              {warehouses.data!.filter((w) => w.isActive).map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
            </select>
          </label>
        )}
        {config.key === 'PURCHASE_INVOICE' && <label>رقم فاتورة المورد<input name="supplierInvoiceNumber" dir="ltr" defaultValue={doc?.supplierInvoiceNumber ?? ''} /></label>}
      </div>
      <label className="check"><input type="checkbox" checked={includeVat} onChange={(e) => setIncludeVat(e.target.checked)} />الأسعار شاملة ضريبة القيمة المضافة</label>

      <div className="table-wrap">
        <table className="table lines-table doc-lines">
          <thead>
            <tr><th>{isPurchase ? 'الصنف أو الحساب' : 'الصنف'}</th><th>الكمية</th><th>سعر الوحدة</th><th>خصم %</th><th className="num">الصافي</th><th className="num">الضريبة</th><th /></tr>
          </thead>
          <tbody>
            {lines.map((l, i) => {
              const p = productById.get(l.productId);
              const pl = preview?.lines[payloadLines(lines.slice(0, i + 1)).length - 1];
              const hasTarget = !!(l.productId || l.accountId);
              return (
                <tr key={i}>
                  <td>
                    <select value={l.productId} aria-label={`صنف السطر ${i + 1}`} onChange={(e) => update(i, { productId: e.target.value, accountId: '', unitPrice: '' })}>
                      <option value="">{isPurchase ? '— صنف —' : '— اختر الصنف —'}</option>
                      {products.data?.data.map((pr) => <option key={pr.id} value={pr.id}>{pr.sku} — {pr.nameAr}</option>)}
                    </select>
                    {isPurchase && !l.productId && (
                      <div className="sub-row">
                        <select value={l.accountId} aria-label={`حساب السطر ${i + 1}`} onChange={(e) => update(i, { accountId: e.target.value })}>
                          <option value="">— أو حساب مصروف / أصل —</option>
                          {accounts.data?.data.filter((a) => ['EXPENSE', 'ASSET', 'COST_OF_GOODS_SOLD'].includes(a.type)).map((a) =>
                            <option key={a.id} value={a.id}>{a.code} — {a.nameAr}</option>)}
                        </select>
                        {l.accountId && (
                          <select value={l.vatCategory} aria-label="فئة الضريبة" onChange={(e) => update(i, { vatCategory: e.target.value })}>
                            <option value="S">خاضع</option><option value="Z">صفري</option><option value="E">معفى</option><option value="O">خارج النطاق</option>
                          </select>
                        )}
                      </div>
                    )}
                    {(l.accountId || p) && <input className="sub-row" placeholder="الوصف (اختياري)" value={l.description} onChange={(e) => update(i, { description: e.target.value })} />}
                  </td>
                  <td><input dir="ltr" inputMode="decimal" value={l.quantity} aria-label={`كمية السطر ${i + 1}`} onChange={(e) => update(i, { quantity: e.target.value })} /></td>
                  <td><input dir="ltr" inputMode="decimal" value={l.unitPrice} aria-label={`سعر السطر ${i + 1}`}
                    placeholder={p ? (isPurchase ? p.purchasePrice : p.salePrice).replace(/\.?0+$/, '') : ''} onChange={(e) => update(i, { unitPrice: e.target.value })} /></td>
                  <td><input dir="ltr" inputMode="decimal" value={l.discountPercent} aria-label={`خصم السطر ${i + 1}`} onChange={(e) => update(i, { discountPercent: e.target.value })} /></td>
                  <td className="num" dir="ltr">{hasTarget && pl ? formatAmount(pl.netAmount) : ''}</td>
                  <td className="num" dir="ltr">{hasTarget && pl ? formatAmount(pl.vatAmount) : ''}</td>
                  <td>{lines.length > 1 && <button type="button" className="icon-btn" aria-label="حذف السطر" onClick={() => setLines((ls) => ls.filter((_, j) => j !== i))}>×</button>}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
      <button type="button" className="btn btn-small" onClick={() => setLines((ls) => [...ls, empty()])}>+ سطر</button>

      <div className="doc-footer">
        <label className="grow">ملاحظات<input name="notes" defaultValue={doc?.notes ?? ''} /></label>
        <div>
          {preview && <Totals t={preview.totals} />}
          {previewError ? <ErrorBox error={previewError} /> : null}
          <p className="muted small">تُحسب المبالغ والضريبة في الخادم.</p>
        </div>
      </div>
      <div className="form-actions">
        <button className="btn btn-primary" disabled={busy || !partyId}>{doc ? 'حفظ' : 'حفظ كمسودة'}</button>
        {doc && <button type="button" className="btn btn-ghost" onClick={onSaved}>إلغاء</button>}
      </div>
    </form>
  );
}
