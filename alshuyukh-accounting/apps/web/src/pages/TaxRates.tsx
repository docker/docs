import { useState, type FormEvent } from 'react';
import { api } from '../api';
import { useAuth } from '../auth';
import { ErrorBox, useLoad } from '../ui';

interface Rate { id: string; vatCategory: string; nameAr: string; rate: string; effectiveFrom: string; effectiveTo: string | null; isActive: boolean }

/** Standard VAT rate history. A new rate closes the previous one the day before it starts. */
export default function TaxRates() {
  const { can } = useAuth();
  const rates = useLoad(() => api<{ data: Rate[] }>('GET', '/api/tax-rates'));
  const [error, setError] = useState<unknown>(null);

  async function add(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const f = Object.fromEntries(new FormData(form)) as Record<string, string>;
    const pct = Number(f.percent);
    setError(null);
    try {
      // The API takes the rate as a fraction ("0.15"); build it from the percent text without floating-point math.
      const [int = '0', frac = ''] = f.percent!.trim().split('.');
      if (!(pct > 0 && pct < 100) || frac.length > 2) throw new Error('أدخل نسبة بين 0 و100 بحد أقصى منزلتين عشريتين');
      const rate = `0.${int.padStart(2, '0')}${frac}`.replace(/0+$/, '').replace(/\.$/, '.0');
      await api('POST', '/api/tax-rates', { nameAr: f.nameAr, rate, effectiveFrom: f.effectiveFrom });
      form.reset();
      rates.reload();
    } catch (err) { setError(err); }
  }

  return (
    <div className="card">
      <h2>نسب ضريبة القيمة المضافة</h2>
      <ErrorBox error={error ?? rates.error} />
      <div className="table-wrap">
        <table className="table">
          <thead><tr><th>الاسم</th><th className="num">النسبة</th><th>من</th><th>إلى</th></tr></thead>
          <tbody>
            {rates.data?.data.map((r) => (
              <tr key={r.id}>
                <td>{r.nameAr}</td>
                <td className="num" dir="ltr">{(Number(r.rate) * 100).toFixed(2).replace(/\.?0+$/, '')}%</td>
                <td dir="ltr">{r.effectiveFrom}</td><td dir="ltr">{r.effectiveTo ?? 'مستمرة'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <p className="muted small">تُطبَّق النسبة حسب تاريخ المستند، فلا تتأثر المستندات السابقة بتغيير النسبة.</p>
      {can('tax.manage') && (
        <form className="form" onSubmit={add}>
          <div className="grid-3">
            <label>الاسم<input name="nameAr" required minLength={2} placeholder="النسبة الأساسية" /></label>
            <label>النسبة %<input name="percent" required dir="ltr" inputMode="decimal" pattern="[0-9]{1,2}(\.[0-9]{1,2})?" /></label>
            <label>تبدأ من<input name="effectiveFrom" type="date" required /></label>
          </div>
          <button className="btn btn-primary">إضافة نسبة</button>
        </form>
      )}
    </div>
  );
}
