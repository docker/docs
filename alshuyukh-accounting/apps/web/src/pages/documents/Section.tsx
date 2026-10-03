import { NavLink, Navigate, Route, Routes } from 'react-router-dom';
import { useAuth } from '../../auth';
import { PageHeader } from '../../ui';
import { DOCS, type DocConfig } from './config';
import DocumentDetail from './DocumentDetail';
import DocumentForm from './DocumentForm';
import DocumentList from './DocumentList';
import Payments from './Payments';

/** Sales or purchases section: tabs for each document type and payments. */
export default function Section({ party }: { party: 'customer' | 'supplier' }) {
  const { can } = useAuth();
  const docs: DocConfig[] = party === 'customer'
    ? [DOCS.SALES_INVOICE, DOCS.SALES_QUOTE, DOCS.SALES_RETURN]
    : [DOCS.PURCHASE_INVOICE, DOCS.PURCHASE_ORDER, DOCS.PURCHASE_RETURN];
  const paymentsRoute = party === 'customer' ? 'receipts' : 'payments';
  const visible = docs.filter((d) => can(d.perm.view));
  // Absolute paths: inside a splat route, relative links resolve against the full URL.
  const base = party === 'customer' ? '/sales' : '/purchases';
  return (
    <>
      <PageHeader title={party === 'customer' ? 'المبيعات' : 'المشتريات'} />
      <nav className="tabs">
        {visible.map((d) => <NavLink key={d.key} to={`${base}/${d.route}`} className={({ isActive }) => `tab ${isActive ? 'active' : ''}`}>{d.title}</NavLink>)}
        {can('payment.view') && <NavLink to={`${base}/${paymentsRoute}`} className={({ isActive }) => `tab ${isActive ? 'active' : ''}`}>{party === 'customer' ? 'سندات القبض' : 'سندات الصرف'}</NavLink>}
      </nav>
      <Routes>
        <Route index element={<Navigate to={visible[0]?.route ?? paymentsRoute} replace />} />
        {visible.map((d) => (
          <Route key={d.key} path={d.route}>
            <Route index element={<DocumentList config={d} />} />
            {!d.isReturn && <Route path="new" element={<DocumentForm config={d} />} />}
            <Route path=":id" element={<DocumentDetail config={d} />} />
          </Route>
        ))}
        <Route path={paymentsRoute} element={<Payments party={party} />} />
      </Routes>
    </>
  );
}
