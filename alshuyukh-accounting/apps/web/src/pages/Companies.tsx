import { useState, type FormEvent } from 'react';
import { api } from '../api';
import { useAuth } from '../auth';
import { ErrorBox, PageHeader, useLoad } from '../ui';

interface Company { id: string; name: string; legalName: string | null; vatNumber: string | null; commercialRegistration: string | null; city: string | null; address: string | null; phone: string | null; email: string | null; currency: string }
interface Branch { id: string; code: string; name: string; city: string | null; isMain: boolean; isActive: boolean }
interface Warehouse { id: string; code: string; name: string; branchId: string | null; isActive: boolean }

const nullIfEmpty = (v: FormDataEntryValue | null) => (v === null || String(v).trim() === '' ? null : String(v).trim());

export default function Companies() {
  const { can } = useAuth();
  const companies = useLoad(() => api<{ data: Company[] }>('GET', '/api/companies'));
  const [selected, setSelected] = useState<string | null>(null);
  const companyId = selected ?? companies.data?.data[0]?.id ?? null;
  const company = companies.data?.data.find((c) => c.id === companyId);
  const branches = useLoad(() => (companyId ? api<{ data: Branch[] }>('GET', `/api/companies/${companyId}/branches`) : Promise.resolve({ data: [] })), [companyId]);
  const warehouses = useLoad(() => (companyId ? api<{ data: Warehouse[] }>('GET', `/api/companies/${companyId}/warehouses`) : Promise.resolve({ data: [] })), [companyId]);
  const [error, setError] = useState<unknown>(null);
  const [saved, setSaved] = useState(false);
  const manage = can('company.manage');

  async function run(fn: () => Promise<unknown>, after: () => void) {
    setError(null); setSaved(false);
    try { await fn(); setSaved(true); after(); } catch (err) { setError(err); }
  }

  function saveCompany(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    void run(() => api('PATCH', `/api/companies/${companyId}`, {
      name: String(f.get('name')), legalName: nullIfEmpty(f.get('legalName')), vatNumber: nullIfEmpty(f.get('vatNumber')),
      commercialRegistration: nullIfEmpty(f.get('commercialRegistration')), city: nullIfEmpty(f.get('city')),
      address: nullIfEmpty(f.get('address')), phone: nullIfEmpty(f.get('phone')), email: nullIfEmpty(f.get('email')),
    }), companies.reload);
  }

  function addBranch(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const f = new FormData(form);
    void run(() => api('POST', `/api/companies/${companyId}/branches`, { code: String(f.get('code')), name: String(f.get('name')), city: nullIfEmpty(f.get('city')) }),
      () => { form.reset(); branches.reload(); });
  }

  function addWarehouse(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const f = new FormData(form);
    void run(() => api('POST', `/api/companies/${companyId}/warehouses`, { code: String(f.get('code')), name: String(f.get('name')), branchId: nullIfEmpty(f.get('branchId')) }),
      () => { form.reset(); warehouses.reload(); });
  }

  return (
    <>
      <PageHeader title="الشركات والفروع">
        {(companies.data?.data.length ?? 0) > 1 && (
          <select value={companyId ?? ''} onChange={(e) => setSelected(e.target.value)} aria-label="الشركة">
            {companies.data!.data.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select>
        )}
      </PageHeader>
      <ErrorBox error={error ?? companies.error} />
      {saved && <div className="alert alert-ok">تم الحفظ</div>}

      {company && (
        <form key={company.id} className="card form" onSubmit={saveCompany}>
          <h2>بيانات الشركة</h2>
          <fieldset disabled={!manage}>
            <div className="grid-2">
              <label>الاسم التجاري<input name="name" defaultValue={company.name} required /></label>
              <label>الاسم القانوني<input name="legalName" defaultValue={company.legalName ?? ''} /></label>
              <label>الرقم الضريبي<input name="vatNumber" defaultValue={company.vatNumber ?? ''} pattern="3[0-9]{13}3" dir="ltr" /></label>
              <label>السجل التجاري<input name="commercialRegistration" defaultValue={company.commercialRegistration ?? ''} pattern="[0-9]{10}" dir="ltr" /></label>
              <label>المدينة<input name="city" defaultValue={company.city ?? ''} /></label>
              <label>العنوان<input name="address" defaultValue={company.address ?? ''} /></label>
              <label>الهاتف<input name="phone" defaultValue={company.phone ?? ''} dir="ltr" /></label>
              <label>البريد الإلكتروني<input name="email" type="email" defaultValue={company.email ?? ''} dir="ltr" /></label>
            </div>
            <p className="muted small">العملة: {company.currency}</p>
            {manage && <button className="btn btn-primary">حفظ</button>}
          </fieldset>
        </form>
      )}

      <div className="grid-cards">
        <div className="card">
          <h2>الفروع</h2>
          <table className="table">
            <thead><tr><th>الرمز</th><th>الاسم</th><th>المدينة</th><th /></tr></thead>
            <tbody>
              {branches.data?.data.map((b) => (
                <tr key={b.id}><td dir="ltr">{b.code}</td><td>{b.name}</td><td>{b.city ?? '—'}</td><td>{b.isMain && <span className="tag">رئيسي</span>}</td></tr>
              ))}
            </tbody>
          </table>
          {manage && (
            <form className="inline-form" onSubmit={addBranch}>
              <input name="code" placeholder="الرمز" required pattern="[A-Za-z0-9_\-]{1,20}" dir="ltr" />
              <input name="name" placeholder="اسم الفرع" required minLength={2} />
              <input name="city" placeholder="المدينة" />
              <button className="btn">إضافة فرع</button>
            </form>
          )}
        </div>

        <div className="card">
          <h2>المستودعات</h2>
          <table className="table">
            <thead><tr><th>الرمز</th><th>الاسم</th><th>الفرع</th></tr></thead>
            <tbody>
              {warehouses.data?.data.map((w) => (
                <tr key={w.id}><td dir="ltr">{w.code}</td><td>{w.name}</td><td>{branches.data?.data.find((b) => b.id === w.branchId)?.name ?? '—'}</td></tr>
              ))}
            </tbody>
          </table>
          {manage && (
            <form className="inline-form" onSubmit={addWarehouse}>
              <input name="code" placeholder="الرمز" required pattern="[A-Za-z0-9_\-]{1,20}" dir="ltr" />
              <input name="name" placeholder="اسم المستودع" required minLength={2} />
              <select name="branchId" defaultValue="">
                <option value="">بدون فرع</option>
                {branches.data?.data.map((b) => <option key={b.id} value={b.id}>{b.name}</option>)}
              </select>
              <button className="btn">إضافة مستودع</button>
            </form>
          )}
        </div>
      </div>
    </>
  );
}
