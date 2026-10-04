import { Link, useParams } from 'react-router-dom';
import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import type { Doc, DocConfig } from './config';
import { loadZatcaDocument } from './EInvoicePanel';

interface Company { id: string; name: string; legalName: string | null; vatNumber: string | null; commercialRegistration: string | null; buildingNumber: string | null; street: string | null; district: string | null; city: string | null; postalCode: string | null; additionalNumber: string | null; phone: string | null }

const joinAddress = (a: { buildingNumber?: string | null; street?: string | null; district?: string | null; city?: string | null; postalCode?: string | null; additionalNumber?: string | null } | null | undefined) =>
  a ? [a.buildingNumber && a.street ? `${a.buildingNumber} ${a.street}` : a.street, a.district, [a.city, a.postalCode].filter(Boolean).join(' '), a.additionalNumber ? `الرقم الإضافي ${a.additionalNumber}` : null].filter(Boolean).join('، ') : '';

/** Printable tax invoice / simplified tax invoice / credit note with its QR code. */
export default function PrintDocument({ config }: { config: DocConfig }) {
  const { id } = useParams();
  const type = config.key === 'SALES_RETURN' ? 'SALES_RETURN' : 'SALES_INVOICE';
  const data = useLoad(async () => {
    const [doc, companies, z] = await Promise.all([
      api<Doc>('GET', `/api/${config.api}/${id}`),
      api<{ data: Company[] }>('GET', '/api/companies'),
      loadZatcaDocument(type, id!),
    ]);
    return { doc, company: companies.data.find((c) => c.id === doc.companyId) ?? companies.data[0]!, z };
  }, [id, config.api]);
  if (!data.data) return <ErrorBox error={data.error} />;
  const { doc: d, company: co, z } = data.data;
  const standard = d.invoiceKind === 'STANDARD';
  const title = config.key === 'SALES_RETURN' ? (standard ? 'إشعار دائن ضريبي' : 'إشعار دائن ضريبي مبسط') : (standard ? 'فاتورة ضريبية' : 'فاتورة ضريبية مبسطة');
  const notCleared = standard && z.einvoice && z.einvoice.status !== 'CLEARED';
  const buyer = d.partySnapshot;
  return (
    <div className="print-page">
      <div className="toolbar no-print">
        <Link className="btn btn-ghost" to={`../${id}`} relative="path">رجوع</Link>
        <button className="btn btn-primary" onClick={() => window.print()}>طباعة</button>
      </div>
      {d.status === 'DRAFT' && <div className="alert alert-warn">مسودة — ليست مستندًا ضريبيًا.</div>}
      {notCleared && <div className="alert alert-error">لم تُعتمد من هيئة الزكاة والضريبة والجمارك بعد — لا تُسلَّم للعميل.</div>}
      <article className="invoice-sheet">
        <header className="invoice-head">
          <div>
            <h1>{title}</h1>
            <dl className="meta">
              <div><dt>الرقم</dt><dd dir="ltr">{d.number ?? '—'}</dd></div>
              <div><dt>التاريخ</dt><dd dir="ltr">{d.date}</dd></div>
              {d.originalInvoiceNumber && <div><dt>الفاتورة الأصلية</dt><dd dir="ltr">{d.originalInvoiceNumber}</dd></div>}
              {d.reason && <div><dt>سبب الإشعار</dt><dd>{d.reason}</dd></div>}
            </dl>
          </div>
          {z.qrSvg && !notCleared && <img className="invoice-qr" alt="رمز QR للفاتورة" src={`data:image/svg+xml;base64,${btoa(z.qrSvg)}`} />}
        </header>
        <section className="invoice-parties">
          <div>
            <h2>البائع</h2>
            <strong>{co.legalName ?? co.name}</strong>
            {co.vatNumber && <div>الرقم الضريبي: <span dir="ltr">{co.vatNumber}</span></div>}
            {co.commercialRegistration && <div>السجل التجاري: <span dir="ltr">{co.commercialRegistration}</span></div>}
            <div className="small">{joinAddress(co)}</div>
          </div>
          <div>
            <h2>المشتري</h2>
            <strong>{buyer?.nameAr ?? d.partyName}</strong>
            {buyer?.vatNumber && <div>الرقم الضريبي: <span dir="ltr">{buyer.vatNumber}</span></div>}
            <div className="small">{joinAddress(buyer?.address)}</div>
          </div>
        </section>
        <table className="table invoice-lines">
          <thead><tr><th>#</th><th>الصنف</th><th className="num">الكمية</th><th className="num">سعر الوحدة</th><th className="num">المبلغ الخاضع</th><th className="num">نسبة الضريبة</th><th className="num">الضريبة</th><th className="num">الإجمالي</th></tr></thead>
          <tbody>
            {d.lines.map((l) => (
              <tr key={l.id}>
                <td>{l.lineNo}</td><td>{l.description || l.productName}</td>
                <td className="num" dir="ltr">{l.quantity.replace(/\.?0+$/, '')}</td>
                <td className="num" dir="ltr">{formatAmount(l.unitPrice)}</td>
                <td className="num" dir="ltr">{formatAmount(l.netAmount)}</td>
                <td className="num" dir="ltr">{(Number(l.vatRate) * 100).toFixed(0)}%</td>
                <td className="num" dir="ltr">{formatAmount(l.vatAmount)}</td>
                <td className="num" dir="ltr">{formatAmount(l.totalAmount)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        <dl className="totals">
          <div><dt>الإجمالي غير شامل الضريبة</dt><dd dir="ltr">{formatAmount(d.taxableAmount)}</dd></div>
          {d.discountTotal !== '0.00' && <div><dt>الخصم</dt><dd dir="ltr">{formatAmount(d.discountTotal)}</dd></div>}
          <div><dt>ضريبة القيمة المضافة</dt><dd dir="ltr">{formatAmount(d.taxAmount)}</dd></div>
          <div className="strong"><dt>الإجمالي شامل الضريبة</dt><dd dir="ltr">{formatAmount(d.total)} ر.س</dd></div>
        </dl>
      </article>
    </div>
  );
}
