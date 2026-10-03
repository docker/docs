import { useState, type FormEvent } from 'react';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { ErrorBox, PageHeader, useLoad } from '../../ui';

interface Product {
  id: string; sku: string; barcode: string | null; nameAr: string; nameEn: string | null; productType: string;
  categoryId: string | null; categoryName: string | null; unitId: string; unitName: string; salePrice: string;
  salePriceIncludesVat: boolean; purchasePrice: string; vatCategory: string; trackInventory: boolean; isActive: boolean;
}
interface Unit { id: string; code: string; nameAr: string; isActive: boolean }
interface Category { id: string; parentId: string | null; nameAr: string; isActive: boolean }

const VAT_AR: Record<string, string> = { S: 'خاضع للنسبة الأساسية', Z: 'نسبة صفرية', E: 'معفى', O: 'خارج النطاق' };
const blank = (v: FormDataEntryValue | null) => (v === null || String(v).trim() === '' ? null : String(v).trim());
/** Shows a 4-decimal unit price without trailing zeros beyond 2 places. */
const price = (v: string) => v.replace(/(\.\d\d)00$/, '$1').replace(/(\.\d\d\d)0$/, '$1');

export default function Products() {
  const { can } = useAuth();
  const [tab, setTab] = useState<'products' | 'setup'>('products');
  return (
    <>
      <PageHeader title="المخزون" />
      <nav className="tabs">
        <button className={`tab ${tab === 'products' ? 'active' : ''}`} onClick={() => setTab('products')}>المنتجات والخدمات</button>
        {can('product.manage') && <button className={`tab ${tab === 'setup' ? 'active' : ''}`} onClick={() => setTab('setup')}>الوحدات والتصنيفات</button>}
      </nav>
      {tab === 'products' ? <ProductList /> : <Setup />}
      <p className="muted small">حركات المخزون والأرصدة والتكلفة تُضاف في المرحلة 5.</p>
    </>
  );
}

function ProductList() {
  const { can } = useAuth();
  const [search, setSearch] = useState('');
  const [editing, setEditing] = useState<Product | 'new' | null>(null);
  const [error, setError] = useState<unknown>(null);
  const products = useLoad(() => api<{ data: Product[]; total: number }>('GET', `/api/products?limit=200${search ? `&search=${encodeURIComponent(search)}` : ''}`), [search]);
  const units = useLoad(() => api<{ data: Unit[] }>('GET', '/api/units'));
  const categories = useLoad(() => api<{ data: Category[] }>('GET', '/api/product-categories'));
  const manage = can('product.manage');

  async function save(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const type = String(f.get('productType'));
    const body = {
      nameAr: String(f.get('nameAr')), nameEn: blank(f.get('nameEn')), barcode: blank(f.get('barcode')),
      productType: type, unitId: String(f.get('unitId')), categoryId: blank(f.get('categoryId')),
      salePrice: String(f.get('salePrice') || '0'), purchasePrice: String(f.get('purchasePrice') || '0'),
      salePriceIncludesVat: f.get('salePriceIncludesVat') === 'on', vatCategory: String(f.get('vatCategory')),
      trackInventory: type === 'GOODS' && f.get('trackInventory') === 'on',
    };
    setError(null);
    try {
      if (editing === 'new') await api('POST', '/api/products', { ...body, sku: blank(f.get('sku')) ?? undefined });
      else if (editing) await api('PATCH', `/api/products/${editing.id}`, body);
      setEditing(null);
      products.reload();
    } catch (err) { setError(err); }
  }

  const current = editing && editing !== 'new' ? editing : null;
  return (
    <>
      <ErrorBox error={error ?? products.error} />
      {editing && (
        <form key={current?.id ?? 'new'} className="card form" onSubmit={save}>
          <h2>{current ? `تعديل: ${current.nameAr}` : 'منتج أو خدمة جديدة'}</h2>
          <div className="grid-3">
            <label>الاسم بالعربية<input name="nameAr" required minLength={2} defaultValue={current?.nameAr} /></label>
            <label>الاسم بالإنجليزية<input name="nameEn" dir="ltr" defaultValue={current?.nameEn ?? ''} /></label>
            <label>النوع
              <select name="productType" defaultValue={current?.productType ?? 'GOODS'}>
                <option value="GOODS">سلعة</option><option value="SERVICE">خدمة</option>
              </select>
            </label>
            {!current && <label>رمز الصنف SKU (تلقائي إن تُرك فارغًا)<input name="sku" dir="ltr" /></label>}
            <label>الباركود<input name="barcode" dir="ltr" defaultValue={current?.barcode ?? ''} /></label>
            <label>الوحدة
              <select name="unitId" required defaultValue={current?.unitId ?? units.data?.data.find((u) => u.code === 'PCE')?.id}>
                {units.data?.data.filter((u) => u.isActive).map((u) => <option key={u.id} value={u.id}>{u.nameAr} ({u.code})</option>)}
              </select>
            </label>
            <label>التصنيف
              <select name="categoryId" defaultValue={current?.categoryId ?? ''}>
                <option value="">بدون تصنيف</option>
                {categories.data?.data.map((c) => <option key={c.id} value={c.id}>{c.nameAr}</option>)}
              </select>
            </label>
            <label>سعر البيع<input name="salePrice" dir="ltr" inputMode="decimal" pattern="[0-9]{1,14}(\.[0-9]{1,4})?" defaultValue={current?.salePrice ?? ''} /></label>
            <label>سعر الشراء<input name="purchasePrice" dir="ltr" inputMode="decimal" pattern="[0-9]{1,14}(\.[0-9]{1,4})?" defaultValue={current?.purchasePrice ?? ''} /></label>
            <label>فئة ضريبة القيمة المضافة
              <select name="vatCategory" defaultValue={current?.vatCategory ?? 'S'}>
                {Object.entries(VAT_AR).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
              </select>
            </label>
          </div>
          <label className="check"><input type="checkbox" name="salePriceIncludesVat" defaultChecked={current?.salePriceIncludesVat} />سعر البيع شامل الضريبة</label>
          <label className="check"><input type="checkbox" name="trackInventory" defaultChecked={current?.trackInventory ?? true} />تتبع المخزون (للسلع فقط)</label>
          <div className="form-actions">
            <button className="btn btn-primary">حفظ</button>
            <button type="button" className="btn btn-ghost" onClick={() => setEditing(null)}>إغلاق</button>
          </div>
        </form>
      )}
      <div className="card">
        <div className="toolbar">
          <input type="search" className="grow" placeholder="بحث بالاسم أو الرمز أو الباركود"
            onKeyDown={(e) => { if (e.key === 'Enter') setSearch(e.currentTarget.value); }} />
          {manage && <button className="btn btn-primary" onClick={() => { setError(null); setEditing('new'); }}>منتج جديد</button>}
        </div>
        <div className="table-wrap">
          <table className="table">
            <thead><tr><th>الرمز</th><th>الاسم</th><th>النوع</th><th>التصنيف</th><th>الوحدة</th><th className="num">سعر البيع</th><th className="num">سعر الشراء</th></tr></thead>
            <tbody>
              {products.data?.data.map((p) => (
                <tr key={p.id}>
                  <td dir="ltr" className="code">{p.sku}</td>
                  <td>{manage ? <button className="link-btn" onClick={() => setEditing(p)}>{p.nameAr}</button> : p.nameAr}</td>
                  <td>{p.productType === 'GOODS' ? 'سلعة' : 'خدمة'}</td>
                  <td className="muted">{p.categoryName ?? '—'}</td>
                  <td>{p.unitName}</td>
                  <td className="num" dir="ltr">{price(p.salePrice)}{p.salePriceIncludesVat && <span className="muted small"> شامل</span>}</td>
                  <td className="num" dir="ltr">{price(p.purchasePrice)}</td>
                </tr>
              ))}
              {!products.data && !products.error && <tr><td colSpan={7} className="muted empty-row">جارٍ التحميل…</td></tr>}
              {products.data?.data.length === 0 && <tr><td colSpan={7} className="muted empty-row">لا توجد منتجات</td></tr>}
            </tbody>
          </table>
        </div>
      </div>
    </>
  );
}

function Setup() {
  const units = useLoad(() => api<{ data: Unit[] }>('GET', '/api/units'));
  const categories = useLoad(() => api<{ data: Category[] }>('GET', '/api/product-categories'));
  const [error, setError] = useState<unknown>(null);

  async function submit(e: FormEvent<HTMLFormElement>, fn: (f: FormData) => Promise<unknown>, reload: () => void) {
    e.preventDefault();
    const form = e.currentTarget;
    setError(null);
    try { await fn(new FormData(form)); form.reset(); reload(); } catch (err) { setError(err); }
  }

  return (
    <>
      <ErrorBox error={error} />
      <div className="grid-cards">
        <div className="card">
          <h2>وحدات القياس</h2>
          <table className="table"><tbody>
            {units.data?.data.map((u) => <tr key={u.id}><td dir="ltr" className="code">{u.code}</td><td>{u.nameAr}</td></tr>)}
          </tbody></table>
          <form className="inline-form" onSubmit={(e) => submit(e, (f) => api('POST', '/api/units', { code: f.get('code'), nameAr: f.get('nameAr') }), units.reload)}>
            <input name="code" placeholder="الرمز" required dir="ltr" pattern="[A-Za-z0-9]{1,10}" />
            <input name="nameAr" placeholder="الاسم" required />
            <button className="btn">إضافة</button>
          </form>
        </div>
        <div className="card">
          <h2>التصنيفات</h2>
          <table className="table"><tbody>
            {categories.data?.data.map((c) => (
              <tr key={c.id}><td>{c.parentId ? `${categories.data!.data.find((p) => p.id === c.parentId)?.nameAr ?? ''} ← ` : ''}{c.nameAr}</td></tr>
            ))}
          </tbody></table>
          <form className="inline-form" onSubmit={(e) => submit(e, (f) => api('POST', '/api/product-categories', { nameAr: f.get('nameAr'), parentId: f.get('parentId') || null }), categories.reload)}>
            <input name="nameAr" placeholder="اسم التصنيف" required minLength={2} />
            <select name="parentId" defaultValue="">
              <option value="">تصنيف رئيسي</option>
              {categories.data?.data.filter((c) => !c.parentId).map((c) => <option key={c.id} value={c.id}>{c.nameAr}</option>)}
            </select>
            <button className="btn">إضافة</button>
          </form>
        </div>
      </div>
    </>
  );
}
