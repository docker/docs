import { Decimal } from '../../lib/money.js';

/**
 * UBL 2.1 e-invoice in the ZATCA (KSA) profile. The XML is generated without
 * whitespace between elements, so the canonical form used for the invoice
 * hash is unambiguous: removing the signature parts leaves no stray text.
 */

export type TypeCode = '388' | '381' | '383'; // invoice, credit note, debit note
export interface Address { street: string; buildingNumber: string; additionalNumber?: string | null; district: string; city: string; postalCode: string; country: string }
export interface Party { name: string; vatNumber?: string | null; crn?: string | null; address?: Address | null }
export interface Line {
  id: number; name: string; quantity: string; unitCode: string; netAmount: string;
  vatCategory: 'S' | 'Z' | 'E' | 'O'; vatRate: string; vatAmount: string; exemptionCode?: string | null; exemptionReason?: string | null;
}
export interface TaxSubtotal { vatCategory: Line['vatCategory']; vatRate: string; taxableAmount: string; taxAmount: string; exemptionCode?: string | null; exemptionReason?: string | null }
export interface EInvoice {
  number: string; uuid: string; issueDate: string; issueTime: string;
  typeCode: TypeCode; subtype: string; icv: number; previousHash: string; currency: string;
  seller: Party & { address: Address; vatNumber: string };
  buyer: Party | null;
  supplyDate?: string | null;
  /** Credit/debit notes: the original invoice and the reason. */
  billingReference?: { number: string; date: string } | null;
  reason?: string | null;
  paymentMeansCode: string;
  lines: Line[];
  taxSubtotals: TaxSubtotal[];
  totals: { lineExtension: string; taxExclusive: string; taxAmount: string; taxInclusive: string; payable: string };
}

const NS = [
  'xmlns="urn:oasis:names:specification:ubl:schema:xsd:Invoice-2"',
  'xmlns:cac="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2"',
  'xmlns:cbc="urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2"',
  'xmlns:ext="urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2"',
].join(' ');

export const esc = (s: string) => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
const el = (name: string, value: string, attrs = '') => `<${name}${attrs ? ` ${attrs}` : ''}>${esc(value)}</${name}>`;
const amt = (name: string, value: string, currency: string) => el(name, value, `currencyID="${currency}"`);
const percent = (rate: string) => new Decimal(rate).times(100).toFixed(2);

function address(a: Address) {
  return '<cac:PostalAddress>'
    + el('cbc:StreetName', a.street) + el('cbc:BuildingNumber', a.buildingNumber)
    + (a.additionalNumber ? el('cbc:PlotIdentification', a.additionalNumber) : '')
    + el('cbc:CitySubdivisionName', a.district) + el('cbc:CityName', a.city) + el('cbc:PostalZone', a.postalCode)
    + `<cac:Country>${el('cbc:IdentificationCode', a.country)}</cac:Country></cac:PostalAddress>`;
}

function party(tag: string, p: Party) {
  return `<cac:${tag}><cac:Party>`
    + (p.crn ? `<cac:PartyIdentification>${el('cbc:ID', p.crn, 'schemeID="CRN"')}</cac:PartyIdentification>` : '')
    + (p.address ? address(p.address) : '')
    + (p.vatNumber ? `<cac:PartyTaxScheme>${el('cbc:CompanyID', p.vatNumber)}<cac:TaxScheme>${el('cbc:ID', 'VAT')}</cac:TaxScheme></cac:PartyTaxScheme>` : '')
    + (p.name ? `<cac:PartyLegalEntity>${el('cbc:RegistrationName', p.name)}</cac:PartyLegalEntity>` : '')
    + `</cac:Party></cac:${tag}>`;
}

function taxCategory(tag: string, c: { vatCategory: string; vatRate: string; exemptionCode?: string | null; exemptionReason?: string | null }) {
  return `<cac:${tag}>${el('cbc:ID', c.vatCategory)}${el('cbc:Percent', percent(c.vatRate))}`
    + (c.vatCategory !== 'S' && c.exemptionCode ? el('cbc:TaxExemptionReasonCode', c.exemptionCode) + el('cbc:TaxExemptionReason', c.exemptionReason ?? c.exemptionCode) : '')
    + `<cac:TaxScheme>${el('cbc:ID', 'VAT')}</cac:TaxScheme></cac:${tag}>`;
}

/**
 * Unit price such that quantity × price = net exactly. When the division is not
 * exact the net is expressed per BaseQuantity = quantity, which keeps the
 * ZATCA rule "line net = quantity × price / base quantity" free of rounding.
 */
function price(l: Line, currency: string) {
  const q = new Decimal(l.quantity);
  const unit = new Decimal(l.netAmount).dividedBy(q);
  if (unit.decimalPlaces() <= 4 && unit.times(q).equals(l.netAmount)) return `<cac:Price>${amt('cbc:PriceAmount', unit.toFixed(Math.max(2, unit.decimalPlaces())), currency)}</cac:Price>`;
  return `<cac:Price>${amt('cbc:PriceAmount', l.netAmount, currency)}${el('cbc:BaseQuantity', q.toFixed(), `unitCode="${l.unitCode}"`)}</cac:Price>`;
}

/**
 * The invoice without signature or QR. `signatureBlock` and `qr` are inserted
 * later; they are excluded from the hash, so the hash can be computed first.
 */
export function buildInvoiceXml(inv: EInvoice, parts: { ublExtensions?: string; qr?: string } = {}): string {
  const c = inv.currency;
  const body = [
    parts.ublExtensions ?? '',
    el('cbc:ProfileID', 'reporting:1.0'),
    el('cbc:ID', inv.number),
    el('cbc:UUID', inv.uuid),
    el('cbc:IssueDate', inv.issueDate),
    el('cbc:IssueTime', inv.issueTime),
    el('cbc:InvoiceTypeCode', inv.typeCode, `name="${inv.subtype}"`),
    el('cbc:DocumentCurrencyCode', c),
    el('cbc:TaxCurrencyCode', c),
    inv.billingReference ? `<cac:BillingReference><cac:InvoiceDocumentReference>${el('cbc:ID', inv.billingReference.number)}${el('cbc:IssueDate', inv.billingReference.date)}</cac:InvoiceDocumentReference></cac:BillingReference>` : '',
    `<cac:AdditionalDocumentReference>${el('cbc:ID', 'ICV')}${el('cbc:UUID', String(inv.icv))}</cac:AdditionalDocumentReference>`,
    `<cac:AdditionalDocumentReference>${el('cbc:ID', 'PIH')}<cac:Attachment>${el('cbc:EmbeddedDocumentBinaryObject', inv.previousHash, 'mimeCode="text/plain"')}</cac:Attachment></cac:AdditionalDocumentReference>`,
    parts.qr ? `<cac:AdditionalDocumentReference>${el('cbc:ID', 'QR')}<cac:Attachment>${el('cbc:EmbeddedDocumentBinaryObject', parts.qr, 'mimeCode="text/plain"')}</cac:Attachment></cac:AdditionalDocumentReference>` : '',
    parts.ublExtensions ? `<cac:Signature>${el('cbc:ID', 'urn:oasis:names:specification:ubl:signature:Invoice')}${el('cbc:SignatureMethod', 'urn:oasis:names:specification:ubl:dsig:enveloped:xades')}</cac:Signature>` : '',
    party('AccountingSupplierParty', inv.seller),
    // A simplified invoice may have no identified buyer; the element itself is mandatory.
    party('AccountingCustomerParty', inv.buyer ?? { name: '' }),
    inv.supplyDate ? `<cac:Delivery>${el('cbc:ActualDeliveryDate', inv.supplyDate)}</cac:Delivery>` : '',
    `<cac:PaymentMeans>${el('cbc:PaymentMeansCode', inv.paymentMeansCode)}${inv.reason ? el('cbc:InstructionNote', inv.reason) : ''}</cac:PaymentMeans>`,
    // Two TaxTotals: the first in the tax currency (SAR), the second with the breakdown.
    `<cac:TaxTotal>${amt('cbc:TaxAmount', inv.totals.taxAmount, c)}</cac:TaxTotal>`,
    `<cac:TaxTotal>${amt('cbc:TaxAmount', inv.totals.taxAmount, c)}${inv.taxSubtotals.map((s) =>
      `<cac:TaxSubtotal>${amt('cbc:TaxableAmount', s.taxableAmount, c)}${amt('cbc:TaxAmount', s.taxAmount, c)}${taxCategory('TaxCategory', s)}</cac:TaxSubtotal>`).join('')}</cac:TaxTotal>`,
    `<cac:LegalMonetaryTotal>${amt('cbc:LineExtensionAmount', inv.totals.lineExtension, c)}${amt('cbc:TaxExclusiveAmount', inv.totals.taxExclusive, c)}`
      + `${amt('cbc:TaxInclusiveAmount', inv.totals.taxInclusive, c)}${amt('cbc:AllowanceTotalAmount', '0.00', c)}${amt('cbc:PrepaidAmount', '0.00', c)}`
      + `${amt('cbc:PayableAmount', inv.totals.payable, c)}</cac:LegalMonetaryTotal>`,
    ...inv.lines.map((l) => `<cac:InvoiceLine>${el('cbc:ID', String(l.id))}${el('cbc:InvoicedQuantity', new Decimal(l.quantity).toFixed(), `unitCode="${l.unitCode}"`)}`
      + `${amt('cbc:LineExtensionAmount', l.netAmount, c)}`
      + `<cac:TaxTotal>${amt('cbc:TaxAmount', l.vatAmount, c)}${amt('cbc:RoundingAmount', new Decimal(l.netAmount).plus(l.vatAmount).toFixed(2), c)}</cac:TaxTotal>`
      + `<cac:Item>${el('cbc:Name', l.name)}${taxCategory('ClassifiedTaxCategory', l)}</cac:Item>${price(l, c)}</cac:InvoiceLine>`),
  ].join('');
  return `<?xml version="1.0" encoding="UTF-8"?>\n<Invoice ${NS}>${body}</Invoice>`;
}
