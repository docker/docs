import { useState, type FormEvent } from 'react';
import { api } from '../../api';
import { useAuth } from '../../auth';
import { formatAmount } from '../../money';
import { ErrorBox, PageHeader, useLoad } from '../../ui';

interface Address {
  addressType: 'BILLING' | 'SHIPPING'; isDefault?: boolean; buildingNumber: string | null; street: string | null;
  district: string | null; city: string | null; postalCode: string | null; additionalNumber: string | null; country?: string;
}
export interface Party {
  id: string; code: string; partyType: string; nameAr: string; nameEn: string | null; vatNumber: string | null;
  commercialRegistration: string | null; nationalId: string | null; email: string | null; phone: string | null;
  creditLimit: string | null; paymentTermsDays: number; notes: string | null; isActive: boolean; balance: string;
  addresses?: Address[];
}

export interface PartyConfig {
  path: 'customers' | 'suppliers';
  title: string;
  singular: string;
  view: string;
  manage: string;
  /** Customers owe us (debit balance); suppliers are owed (credit balance). */
  balanceLabel: (balance: string) => string;
}

export const CUSTOMERS: PartyConfig = {
  path: 'customers', title: 'العملاء', singular: 'عميل', view: 'customer.view', manage: 'customer.manage',
  balanceLabel: (b) => (b.startsWith('-') ? `دائن ${formatAmount(b.slice(1))}` : formatAmount(b)),
};
export const SUPPLIERS: PartyConfig = {
  path: 'suppliers', title: 'الموردون', singular: 'مورد', view: 'supplier.view', manage: 'supplier.manage',
  balanceLabel: (b) => (b.startsWith('-') ? formatAmount(b.slice(1)) : b === '0.00' ? '0.00' : `مدين ${formatAmount(b)}`),
};

const PAGE = 50;
const blank = (v: FormDataEntryValue | null) => (v === null || String(v).trim() === '' ? null : String(v).trim());

export default function Parties({ config }: { config: PartyConfig }) {
  const { can } = useAuth();
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('active');
  const [offset, setOffset] = useState(0);
  const [editing, setEditing] = useState<Party | 'new' | null>(null);
  const [error, setError] = useState<unknown>(null);
  const list = useLoad(() => {
    const q = new URLSearchParams({ limit: String(PAGE), offset: String(offset), status });
    if (search) q.set('search', search);
    return api<{ data: Party[]; total: number }>('GET', `/api/${config.path}?${q}`);
  }, [config.path, search, status, offset]);
  const manage = can(config.manage);

  async function open(p: Party) {
    setError(null);
    try { setEditing(await api<Party>('GET', `/api/${config.path}/${p.id}`)); } catch (e) { setError(e); }
  }

  async function save(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const body = {
      nameAr: String(f.get('nameAr')), nameEn: blank(f.get('nameEn')), partyType: String(f.get('partyType')),
      vatNumber: blank(f.get('vatNumber')), commercialRegistration: blank(f.get('commercialRegistration')),
      nationalId: blank(f.get('nationalId')), email: blank(f.get('email')), phone: blank(f.get('phone')),
      creditLimit: blank(f.get('creditLimit')), paymentTermsDays: Number(f.get('paymentTermsDays') || 0), notes: blank(f.get('notes')),
    };
    const addr: Address = {
      addressType: 'BILLING', buildingNumber: blank(f.get('buildingNumber')), street: blank(f.get('street')),
      district: blank(f.get('district')), city: blank(f.get('city')), postalCode: blank(f.get('postalCode')),
      additionalNumber: blank(f.get('additionalNumber')),
    };
    const hasAddress = Object.entries(addr).some(([k, v]) => k !== 'addressType' && v);
    setError(null);
    try {
      if (editing === 'new') {
        await api('POST', `/api/${config.path}`, { ...body, code: blank(f.get('code')) ?? undefined, addresses: hasAddress ? [addr] : [] });
      } else if (editing) {
        await api('PATCH', `/api/${config.path}/${editing.id}`, body);
        const others = (editing.addresses ?? []).filter((a) => a.addressType !== 'BILLING' || !a.isDefault);
        await api('PUT', `/api/${config.path}/${editing.id}/addresses`, { addresses: [...(hasAddress ? [addr] : []), ...others] });
      }
      setEditing(null);
      list.reload();
    } catch (err) { setError(err); }
  }

  async function toggle(p: Party) {
    setError(null);
    try { await api('PATCH', `/api/${config.path}/${p.id}`, { isActive: !p.isActive }); list.reload(); } catch (e) { setError(e); }
  }

  const current = editing && editing !== 'new' ? editing : null;
  const billing = current?.addresses?.find((a) => a.addressType === 'BILLING' && a.isDefault);

  return (
    <>
      <PageHeader title={config.title}>
        {manage && <button className="btn btn-primary" onClick={() => { setError(null); setEditing('new'); }}>{config.singular} جديد</button>}
      </PageHeader>
      <ErrorBox error={error ?? list.error} />

      {editing && (
        <form key={current?.id ?? 'new'} className="card form" onSubmit={save}>
          <h2>{current ? `تعديل: ${current.nameAr}` : `${config.singular} جديد`}</h2>
          <fieldset disabled={!manage}>
            <div className="grid-3">
              <label>الاسم بالعربية<input name="nameAr" required minLength={2} defaultValue={current?.nameAr} /></label>
              <label>الاسم بالإنجليزية<input name="nameEn" dir="ltr" defaultValue={current?.nameEn ?? ''} /></label>
              <label>النوع
                <select name="partyType" defaultValue={current?.partyType ?? 'BUSINESS'}>
                  <option value="BUSINESS">منشأة</option><option value="INDIVIDUAL">فرد</option>
                </select>
              </label>
              {!current && <label>الرمز (يُولَّد تلقائيًا إن تُرك فارغًا)<input name="code" dir="ltr" pattern="[0-9A-Za-z_\-]{1,30}" /></label>}
              <label>الرقم الضريبي<input name="vatNumber" dir="ltr" pattern="3[0-9]{13}3" placeholder="3XXXXXXXXXXXXX3" defaultValue={current?.vatNumber ?? ''} /></label>
              <label>السجل التجاري<input name="commercialRegistration" dir="ltr" pattern="[0-9]{10}" defaultValue={current?.commercialRegistration ?? ''} /></label>
              <label>الهوية / الإقامة<input name="nationalId" dir="ltr" pattern="[12][0-9]{9}" defaultValue={current?.nationalId ?? ''} /></label>
              <label>الجوال<input name="phone" dir="ltr" defaultValue={current?.phone ?? ''} /></label>
              <label>البريد الإلكتروني<input name="email" type="email" dir="ltr" defaultValue={current?.email ?? ''} /></label>
              <label>حد الائتمان<input name="creditLimit" dir="ltr" inputMode="decimal" pattern="[0-9]{1,16}(\.[0-9]{1,2})?" defaultValue={current?.creditLimit ?? ''} /></label>
              <label>مدة السداد (أيام)<input name="paymentTermsDays" type="number" min={0} max={365} defaultValue={current?.paymentTermsDays ?? 0} /></label>
            </div>
            <h3>العنوان الوطني</h3>
            <div className="grid-3">
              <label>رقم المبنى<input name="buildingNumber" dir="ltr" pattern="[0-9]{4}" defaultValue={billing?.buildingNumber ?? ''} /></label>
              <label>الشارع<input name="street" defaultValue={billing?.street ?? ''} /></label>
              <label>الحي<input name="district" defaultValue={billing?.district ?? ''} /></label>
              <label>المدينة<input name="city" defaultValue={billing?.city ?? ''} /></label>
              <label>الرمز البريدي<input name="postalCode" dir="ltr" pattern="[0-9]{5}" defaultValue={billing?.postalCode ?? ''} /></label>
              <label>الرقم الإضافي<input name="additionalNumber" dir="ltr" pattern="[0-9]{4}" defaultValue={billing?.additionalNumber ?? ''} /></label>
            </div>
            <label>ملاحظات<input name="notes" defaultValue={current?.notes ?? ''} /></label>
            <div className="form-actions">
              {manage && <button className="btn btn-primary">حفظ</button>}
              <button type="button" className="btn btn-ghost" onClick={() => setEditing(null)}>إغلاق</button>
            </div>
          </fieldset>
        </form>
      )}

      <div className="card">
        <div className="toolbar">
          <input type="search" placeholder="بحث بالاسم أو الرمز أو الرقم الضريبي أو الجوال" defaultValue={search} className="grow"
            onKeyDown={(e) => { if (e.key === 'Enter') { setSearch(e.currentTarget.value); setOffset(0); } }} />
          <select value={status} onChange={(e) => { setStatus(e.target.value); setOffset(0); }} aria-label="الحالة">
            <option value="active">النشطون</option><option value="inactive">غير النشطين</option><option value="all">الكل</option>
          </select>
        </div>
        <div className="table-wrap">
          <table className="table">
            <thead><tr><th>الرمز</th><th>الاسم</th><th>الرقم الضريبي</th><th>الجوال</th><th className="num">الرصيد</th><th /></tr></thead>
            <tbody>
              {list.data?.data.map((p) => (
                <tr key={p.id} className={p.isActive ? '' : 'row-muted'}>
                  <td dir="ltr" className="code">{p.code}</td>
                  <td><button className="link-btn" onClick={() => open(p)}>{p.nameAr}</button>{!p.isActive && <span className="tag tag-off">غير نشط</span>}</td>
                  <td dir="ltr">{p.vatNumber ?? '—'}</td>
                  <td dir="ltr">{p.phone ?? '—'}</td>
                  <td className="num" dir="ltr">{config.balanceLabel(p.balance)}</td>
                  <td className="row-actions">{manage && <button className="btn btn-small" onClick={() => toggle(p)}>{p.isActive ? 'إيقاف' : 'تفعيل'}</button>}</td>
                </tr>
              ))}
              {!list.data && !list.error && <tr><td colSpan={6} className="muted empty-row">جارٍ التحميل…</td></tr>}
              {list.data?.data.length === 0 && <tr><td colSpan={6} className="muted empty-row">لا توجد نتائج</td></tr>}
            </tbody>
          </table>
        </div>
        {list.data && list.data.total > PAGE && (
          <div className="pager">
            <span className="muted small">{offset + 1}–{Math.min(offset + PAGE, list.data.total)} من {list.data.total}</span>
            <button className="btn" disabled={offset === 0} onClick={() => setOffset(offset - PAGE)}>السابق</button>
            <button className="btn" disabled={offset + PAGE >= list.data.total} onClick={() => setOffset(offset + PAGE)}>التالي</button>
          </div>
        )}
        <p className="muted small">الرصيد محسوب من القيود المرحّلة المرتبطة بال{config.singular}.</p>
      </div>
    </>
  );
}
