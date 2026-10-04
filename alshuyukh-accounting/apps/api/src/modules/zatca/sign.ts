import { createHash, createPrivateKey, sign as ecSign } from 'node:crypto';
import { DOMParser, type Element, type Node as XmlNode } from '@xmldom/xmldom';
import { C14nCanonicalization } from 'xml-crypto';
import { certificateHash, type CertInfo } from './crypto.js';
import { encodeQr } from './qr.js';
import { buildInvoiceXml, esc, type EInvoice } from './xml.js';

/** PIH of the first document in a chain: base64 of the hex SHA-256 of "0", as defined by ZATCA. */
export const INITIAL_PIH = 'NWZlY2ViNjZmZmM4NmYzOGQ5NTI3ODZjNmQ2OTZjNzljMmRiYzIzOWRkNGU5MWI0NjcyOWQ3M2EyN2ZiNTdlOQ==';

const UBL_NS = {
  ext: 'urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2',
  cac: 'urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2',
  cbc: 'urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2',
};

/**
 * Invoice hash as ZATCA computes it: drop ext:UBLExtensions, cac:Signature and
 * the QR AdditionalDocumentReference, canonicalise (C14N 1.1; identical to 1.0
 * for this document), SHA-256, base64.
 */
export function invoiceHash(xml: string): string {
  const doc = new DOMParser().parseFromString(xml, 'text/xml');
  const root = doc.documentElement!;
  for (const child of Array.from(root.childNodes) as XmlNode[]) {
    if (child.nodeType !== 1) continue;
    const e = child as Element;
    const isQr = e.localName === 'AdditionalDocumentReference' && e.namespaceURI === UBL_NS.cac
      && Array.from(e.getElementsByTagNameNS(UBL_NS.cbc, 'ID')).some((id) => id.parentNode === e && id.textContent === 'QR');
    if ((e.localName === 'UBLExtensions' && e.namespaceURI === UBL_NS.ext) || (e.localName === 'Signature' && e.namespaceURI === UBL_NS.cac) || isQr) {
      root.removeChild(e);
    }
  }
  const canonical = new C14nCanonicalization().process(root as unknown as globalThis.Node, {}) as string;
  return createHash('sha256').update(canonical, 'utf8').digest('base64');
}

/** ECDSA (secp256k1, SHA-256) signature over the invoice hash bytes, base64. */
export function signHash(hashBase64: string, privateKeyPem: string): string {
  return ecSign('sha256', Buffer.from(hashBase64, 'base64'), createPrivateKey(privateKeyPem)).toString('base64');
}

/**
 * XAdES signed properties. ZATCA's SDK digests this element as text
 * (SHA-256 → hex → base64); the same text is embedded in the document.
 */
export function signedProperties(signingTime: string, cert: CertInfo): string {
  return `<xades:SignedProperties xmlns:xades="http://uri.etsi.org/01903/v1.3.2#" Id="xadesSignedProperties">
    <xades:SignedSignatureProperties>
        <xades:SigningTime>${signingTime}</xades:SigningTime>
        <xades:SigningCertificate>
            <xades:Cert>
                <xades:CertDigest>
                    <ds:DigestMethod xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
                    <ds:DigestValue xmlns:ds="http://www.w3.org/2000/09/xmldsig#">${certificateHash(cert.base64)}</ds:DigestValue>
                </xades:CertDigest>
                <xades:IssuerSerial>
                    <ds:X509IssuerName xmlns:ds="http://www.w3.org/2000/09/xmldsig#">${esc(cert.issuer)}</ds:X509IssuerName>
                    <ds:X509SerialNumber xmlns:ds="http://www.w3.org/2000/09/xmldsig#">${cert.serialNumber}</ds:X509SerialNumber>
                </xades:IssuerSerial>
            </xades:Cert>
        </xades:SigningCertificate>
    </xades:SignedSignatureProperties>
</xades:SignedProperties>`;
}

function ublExtensions(hash: string, signature: string, cert: CertInfo, props: string): string {
  const propsHash = Buffer.from(createHash('sha256').update(props, 'utf8').digest('hex')).toString('base64');
  return `<ext:UBLExtensions><ext:UBLExtension><ext:ExtensionURI>urn:oasis:names:specification:ubl:dsig:enveloped:xades</ext:ExtensionURI><ext:ExtensionContent>`
    + `<sig:UBLDocumentSignatures xmlns:sig="urn:oasis:names:specification:ubl:schema:xsd:CommonSignatureComponents-2" xmlns:sac="urn:oasis:names:specification:ubl:schema:xsd:SignatureAggregateComponents-2" xmlns:sbc="urn:oasis:names:specification:ubl:schema:xsd:SignatureBasicComponents-2">`
    + `<sac:SignatureInformation><cbc:ID>urn:oasis:names:specification:ubl:signature:1</cbc:ID><sbc:ReferencedSignatureID>urn:oasis:names:specification:ubl:signature:Invoice</sbc:ReferencedSignatureID>`
    + `<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="signature"><ds:SignedInfo>`
    + `<ds:CanonicalizationMethod Algorithm="http://www.w3.org/2006/12/xml-c14n11"/><ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256"/>`
    + `<ds:Reference Id="invoiceSignedData" URI=""><ds:Transforms>`
    + `<ds:Transform Algorithm="http://www.w3.org/TR/1999/REC-xpath-19991116"><ds:XPath>not(//ancestor-or-self::ext:UBLExtensions)</ds:XPath></ds:Transform>`
    + `<ds:Transform Algorithm="http://www.w3.org/TR/1999/REC-xpath-19991116"><ds:XPath>not(//ancestor-or-self::cac:Signature)</ds:XPath></ds:Transform>`
    + `<ds:Transform Algorithm="http://www.w3.org/TR/1999/REC-xpath-19991116"><ds:XPath>not(//ancestor-or-self::cac:AdditionalDocumentReference[cbc:ID='QR'])</ds:XPath></ds:Transform>`
    + `<ds:Transform Algorithm="http://www.w3.org/2006/12/xml-c14n11"/></ds:Transforms>`
    + `<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>${hash}</ds:DigestValue></ds:Reference>`
    + `<ds:Reference Type="http://www.w3.org/2000/09/xmldsig#SignatureProperties" URI="#xadesSignedProperties">`
    + `<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>${propsHash}</ds:DigestValue></ds:Reference>`
    + `</ds:SignedInfo><ds:SignatureValue>${signature}</ds:SignatureValue>`
    + `<ds:KeyInfo><ds:X509Data><ds:X509Certificate>${cert.base64}</ds:X509Certificate></ds:X509Data></ds:KeyInfo>`
    + `<ds:Object><xades:QualifyingProperties xmlns:xades="http://uri.etsi.org/01903/v1.3.2#" Target="signature">${props}</xades:QualifyingProperties></ds:Object>`
    + `</ds:Signature></sac:SignatureInformation></sig:UBLDocumentSignatures></ext:ExtensionContent></ext:UBLExtension></ext:UBLExtensions>`;
}

export interface SignedInvoice { xml: string; hash: string; signature: string; qr: string }

/**
 * Produces the final e-invoice: hash of the unsigned document, ECDSA signature,
 * XAdES block, and the phase-2 QR code (tags 1–9).
 */
export function signInvoice(inv: EInvoice, key: { privateKeyPem: string; cert: CertInfo }, signingTime: string): SignedInvoice {
  const hash = invoiceHash(buildInvoiceXml(inv));
  const signature = signHash(hash, key.privateKeyPem);
  const qr = encodeQr({
    sellerName: inv.seller.name, vatNumber: inv.seller.vatNumber, timestamp: `${inv.issueDate}T${inv.issueTime}`,
    total: inv.totals.taxInclusive, vatTotal: inv.totals.taxAmount, invoiceHash: hash, signature,
    publicKey: key.cert.publicKeyDer,
    // Tag 9 is for simplified invoices only (standard ones get their QR from ZATCA at clearance).
    certificateSignature: inv.subtype.startsWith('02') ? key.cert.signature : undefined,
  });
  const xml = buildInvoiceXml(inv, { ublExtensions: ublExtensions(hash, signature, key.cert, signedProperties(signingTime, key.cert)), qr });
  return { xml, hash, signature, qr };
}
