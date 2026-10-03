import { Link } from 'react-router-dom';
import { api } from '../api';
import { useAuth } from '../auth';
import { PageHeader, useLoad } from '../ui';

interface Company { id: string; name: string; vatNumber: string | null; commercialRegistration: string | null; city: string | null; currency: string }

export default function Home() {
  const { me, can } = useAuth();
  const companies = useLoad(() => (can('company.view') ? api<{ data: Company[] }>('GET', '/api/companies') : Promise.resolve({ data: [] })));
  const company = companies.data?.data[0];

  const steps = [
    { done: !!company?.vatNumber, label: 'أضف الرقم الضريبي للشركة', to: '/settings/companies' },
    { done: !!company?.commercialRegistration, label: 'أضف رقم السجل التجاري', to: '/settings/companies' },
    { done: !!company?.city, label: 'أكمل عنوان الشركة', to: '/settings/companies' },
  ];

  return (
    <>
      <PageHeader title={`مرحبًا، ${me?.user.fullName ?? ''}`} />
      <div className="grid-cards">
        <div className="card">
          <h2>المؤشرات المالية</h2>
          <p className="muted">
            تُحسب المبيعات والأرباح والذمم والنقدية من القيود المحاسبية المرحّلة فقط.
            ستظهر هنا بعد تفعيل محرك المحاسبة في المرحلة 2.
          </p>
        </div>
        {company && can('company.manage') && (
          <div className="card">
            <h2>إعداد الشركة</h2>
            <ul className="checklist">
              {steps.map((s) => (
                <li key={s.label} className={s.done ? 'done' : ''}>
                  <span aria-hidden>{s.done ? '✓' : '○'}</span>
                  {s.done ? s.label : <Link to={s.to}>{s.label}</Link>}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </>
  );
}
