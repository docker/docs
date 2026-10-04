import { Decimal, toMoney } from '../../lib/money.js';
import type { Db } from '../../db/tx.js';
import { badRequest } from '../../lib/errors.js';
import { taxGroups } from '../tax/ledger.js';
import type { Address, EInvoice, Line, TypeCode } from './xml.js';

/**
 * Turns an issued sales invoice or sales return into the data of a ZATCA
 * e-invoice, and checks everything ZATCA requires before anything is signed.
 */
export type ZatcaDocType = 'SALES_INVOICE' | 'SALES_RETURN';
interface Problem { path: string; message: string }

interface Seller { name: string; vatNumber: string | null; crn: string | null; street: string | null; buildingNumber: string | null; additionalNumber: string | null; district: string | null; city: string | null; postalCode: string | null; country: string }

export async function loadSeller(db: Db, companyId: string): Promise<Seller> {
  const { rows: [c] } = await db.query<Seller>(
    `SELECT COALESCE(legal_name, name) AS name, vat_number AS "vatNumber", commercial_registration AS crn, street,
            building_number AS "buildingNumber", additional_number AS "additionalNumber", district, city, postal_code AS "postalCode", country
       FROM companies WHERE id = $1`, [companyId]);
  return c!;
}

/** Seller requirements: VAT number, CRN and a complete national address. */
export function sellerProblems(s: Seller): Problem[] {
  const p: Problem[] = [];
  if (!s.vatNumber) p.push({ path: 'company.vatNumber', message: 'الرقم الضريبي للشركة مطلوب' });
  if (!s.crn) p.push({ path: 'company.commercialRegistration', message: 'رقم السجل التجاري للشركة مطلوب' });
  for (const [k, label] of [['street', 'الشارع'], ['buildingNumber', 'رقم المبنى'], ['district', 'الحي'], ['city', 'المدينة'], ['postalCode', 'الرمز البريدي']] as const) {
    if (!s[k]) p.push({ path: `company.${k}`, message: `${label} في العنوان الوطني للشركة مطلوب` });
  }
  return p;
}

export const sellerAddress = (s: Seller): Address => ({
  street: s.street!, buildingNumber: s.buildingNumber!, additionalNumber: s.additionalNumber, district: s.district!, city: s.city!, postalCode: s.postalCode!, country: s.country,
});

const OUT_OF_SCOPE = { code: 'VATEX-SA-OOS', reason: 'خارج نطاق ضريبة القيمة المضافة' };

export interface BuildContext { uuid: string; icv: number; previousHash: string; issueTime: string }

/** Builds the e-invoice data for an issued document, or throws 400 ZATCA_INVALID listing every problem. */
export async function buildEInvoice(db: Db, type: ZatcaDocType, documentId: string, c: BuildContext): Promise<{ invoice: EInvoice; kind: 'STANDARD' | 'SIMPLIFIED'; companyId: string; number: string }> {
  const table = type === 'SALES_INVOICE' ? 'sales_invoices' : 'sales_returns';
  const items = type === 'SALES_INVOICE' ? 'sales_invoice_items' : 'sales_return_items';
  const { rows: [d] } = await db.query<{
    company_id: string; doc_number: string; doc_date: string; due_date: string | null; invoice_kind: string | null; party_snapshot: any;
    taxable_amount: string; tax_amount: string; total: string; reason: string | null; original_number: string | null; original_date: string | null;
  }>(
    `SELECT d.company_id, d.doc_number, d.doc_date, ${type === 'SALES_INVOICE' ? 'd.due_date' : 'NULL::date AS due_date'},
            d.invoice_kind, d.party_snapshot, d.taxable_amount::text, d.tax_amount::text, d.total::text,
            ${type === 'SALES_RETURN'
              ? `d.reason, o.doc_number AS original_number, o.doc_date AS original_date FROM ${table} d JOIN sales_invoices o ON o.id = d.original_invoice_id`
              : `NULL AS reason, NULL AS original_number, NULL AS original_date FROM ${table} d`}
      WHERE d.id = $1`, [documentId]);
  if (!d) throw badRequest('DOCUMENT_NOT_FOUND', 'Document not found');
  const kind = (d.invoice_kind ?? 'SIMPLIFIED') as 'STANDARD' | 'SIMPLIFIED';
  const seller = await loadSeller(db, d.company_id);
  const problems = sellerProblems(seller);

  const snap = d.party_snapshot ?? {};
  const addr = snap.address ?? null;
  if (kind === 'STANDARD') {
    if (!snap.vatNumber) problems.push({ path: 'customer.vatNumber', message: 'الرقم الضريبي للعميل مطلوب في الفاتورة الضريبية' });
    for (const [k, label] of [['street', 'الشارع'], ['buildingNumber', 'رقم المبنى'], ['district', 'الحي'], ['city', 'المدينة'], ['postalCode', 'الرمز البريدي']] as const) {
      if (!addr?.[k]) problems.push({ path: `customer.address.${k}`, message: `${label} في عنوان العميل مطلوب في الفاتورة الضريبية` });
    }
  }

  const { rows: lines } = await db.query<{ line_no: number; name: string; quantity: string; unit_code: string | null; net_amount: string; vat_category: Line['vatCategory']; vat_rate: string; vat_amount: string; exemption_code: string | null; exemption_reason: string | null }>(
    `SELECT i.line_no, COALESCE(NULLIF(btrim(i.description), ''), p.name_ar, 'بند') AS name, i.quantity::text, u.code AS unit_code,
            i.net_amount::text, i.vat_category, i.vat_rate::text, i.vat_amount::text,
            p.vat_exemption_code AS exemption_code, p.vat_exemption_reason AS exemption_reason
       FROM ${items} i LEFT JOIN products p ON p.id = i.product_id LEFT JOIN units u ON u.id = i.unit_id
      WHERE i.document_id = $1 ORDER BY i.line_no`, [documentId]);

  const eLines: Line[] = lines.map((l) => {
    let exemptionCode = l.exemption_code;
    let exemptionReason = l.exemption_reason;
    if (l.vat_category === 'O' && !exemptionCode) ({ code: exemptionCode, reason: exemptionReason } = OUT_OF_SCOPE);
    if ((l.vat_category === 'Z' || l.vat_category === 'E') && !exemptionCode) {
      problems.push({ path: `lines.${l.line_no}`, message: `السطر ${l.line_no} (${l.name}): رمز سبب الإعفاء أو النسبة الصفرية مطلوب في بطاقة الصنف` });
    }
    return {
      id: l.line_no, name: l.name, quantity: l.quantity, unitCode: l.unit_code ?? 'PCE', netAmount: l.net_amount,
      vatCategory: l.vat_category, vatRate: l.vat_rate, vatAmount: l.vat_amount, exemptionCode, exemptionReason,
    };
  });

  // Tax breakdown equals what was posted; one exemption reason per category.
  const groups = taxGroups(eLines, d.tax_amount).map((g) => {
    const codes = [...new Set(eLines.filter((l) => l.vatCategory === g.vatCategory && new Decimal(l.vatRate).equals(g.vatRate)).map((l) => l.exemptionCode ?? null))];
    if (g.vatCategory !== 'S' && codes.length > 1) problems.push({ path: 'lines', message: `أسباب إعفاء مختلفة للفئة ${g.vatCategory} في مستند واحد غير مدعومة؛ افصلها في مستندات مستقلة` });
    const line = eLines.find((l) => l.vatCategory === g.vatCategory && l.exemptionCode === codes[0]);
    return { ...g, exemptionCode: line?.exemptionCode ?? null, exemptionReason: line?.exemptionReason ?? null };
  });

  if (problems.length) throw badRequest('ZATCA_INVALID', 'لا يمكن إنشاء الفاتورة الإلكترونية: بيانات ناقصة', problems);

  const typeCode: TypeCode = type === 'SALES_INVOICE' ? '388' : '381';
  const invoice: EInvoice = {
    number: d.doc_number, uuid: c.uuid, issueDate: d.doc_date, issueTime: c.issueTime,
    typeCode, subtype: kind === 'STANDARD' ? '0100000' : '0200000', icv: c.icv, previousHash: c.previousHash, currency: 'SAR',
    seller: { name: seller.name, vatNumber: seller.vatNumber!, crn: seller.crn, address: sellerAddress(seller) },
    buyer: snap.nameAr ? {
      name: snap.nameAr, vatNumber: snap.vatNumber ?? null, crn: null,
      address: addr?.street && addr?.buildingNumber ? { street: addr.street, buildingNumber: addr.buildingNumber, additionalNumber: addr.additionalNumber, district: addr.district ?? '', city: addr.city ?? '', postalCode: addr.postalCode ?? '', country: addr.country ?? 'SA' } : null,
    } : null,
    supplyDate: d.doc_date,
    billingReference: type === 'SALES_RETURN' ? { number: d.original_number!, date: d.original_date! } : null,
    reason: d.reason,
    // 30 = credit transfer (sold on credit), 10 = cash.
    paymentMeansCode: d.due_date && d.due_date > d.doc_date ? '30' : '10',
    lines: eLines,
    taxSubtotals: groups.map((g) => ({ vatCategory: g.vatCategory, vatRate: g.vatRate, taxableAmount: g.taxableAmount, taxAmount: g.taxAmount, exemptionCode: g.exemptionCode, exemptionReason: g.exemptionReason })),
    totals: { lineExtension: toMoney(d.taxable_amount), taxExclusive: toMoney(d.taxable_amount), taxAmount: toMoney(d.tax_amount), taxInclusive: toMoney(d.total), payable: toMoney(d.total) },
  };
  return { invoice, kind, companyId: d.company_id, number: d.doc_number };
}

/** Sample documents for the onboarding compliance checks (not part of any real chain). */
export function sampleInvoices(seller: Seller, types: string, at: { date: string; time: string }): EInvoice[] {
  const base = (n: number, typeCode: TypeCode, standard: boolean): Omit<EInvoice, 'uuid' | 'icv' | 'previousHash'> => ({
    number: `COMPLIANCE-${n}`, issueDate: at.date, issueTime: at.time, typeCode, subtype: standard ? '0100000' : '0200000', currency: 'SAR',
    seller: { name: seller.name, vatNumber: seller.vatNumber!, crn: seller.crn, address: sellerAddress(seller) },
    buyer: standard
      ? { name: 'عميل فحص الامتثال', vatNumber: '399999999800003', address: { street: 'شارع الملك فهد', buildingNumber: '1234', district: 'العليا', city: 'الرياض', postalCode: '12211', country: 'SA' } }
      : null,
    supplyDate: at.date,
    billingReference: typeCode === '388' ? null : { number: 'COMPLIANCE-1', date: at.date },
    reason: typeCode === '381' ? 'إرجاع بضاعة' : typeCode === '383' ? 'تعديل السعر' : null,
    paymentMeansCode: '10',
    lines: [{ id: 1, name: 'بند فحص الامتثال', quantity: '1', unitCode: 'PCE', netAmount: '100.00', vatCategory: 'S', vatRate: '0.15', vatAmount: '15.00' }],
    taxSubtotals: [{ vatCategory: 'S', vatRate: '0.15', taxableAmount: '100.00', taxAmount: '15.00' }],
    totals: { lineExtension: '100.00', taxExclusive: '100.00', taxAmount: '15.00', taxInclusive: '115.00', payable: '115.00' },
  });
  const out: Omit<EInvoice, 'uuid' | 'icv' | 'previousHash'>[] = [];
  if (types[0] === '1') out.push(base(1, '388', true), base(2, '381', true), base(3, '383', true));
  if (types[1] === '1') out.push(base(4, '388', false), base(5, '381', false), base(6, '383', false));
  return out as EInvoice[];
}
