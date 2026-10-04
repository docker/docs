import { execFileSync } from 'node:child_process';
import { createPublicKey, verify } from 'node:crypto';
import { mkdtempSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { DOMParser } from '@xmldom/xmldom';
import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { withTx } from '../src/db/tx.js';
import { setTransport } from '../src/modules/zatca/client.js';
import { buildCsr, decrypt, encrypt, generateKeyPair, parseCertificate, certificateFromCsid } from '../src/modules/zatca/crypto.js';
import { decodeQr, encodeQr } from '../src/modules/zatca/qr.js';
import { INITIAL_PIH, invoiceHash, signInvoice } from '../src/modules/zatca/sign.js';
import { buildInvoiceXml, type EInvoice } from '../src/modules/zatca/xml.js';
import { commerce, type Commerce } from './commerce-helpers.js';
import { addMember, client, register, setupApp, type Session, type TestContext } from './helpers.js';
import { FakeZatca } from './zatca-fake.js';
import { runDueSubmissions } from '../src/modules/zatca/worker.js';

const UBL_CBC = 'urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2';
const SELLER_ADDRESS = { street: 'طريق الملك عبدالعزيز', buildingNumber: '2322', district: 'الملقا', city: 'الرياض', postalCode: '13521', country: 'SA' };

const sample = (over: Partial<EInvoice> = {}): EInvoice => ({
  number: 'INV-000001', uuid: '3cf5ee18-ee25-44ea-a444-2c37ba7f28be', issueDate: '2026-03-10', issueTime: '14:30:00',
  typeCode: '388', subtype: '0200000', icv: 1, previousHash: INITIAL_PIH, currency: 'SAR',
  seller: { name: 'مؤسسة الشيوخ & أبناؤه', vatNumber: '399999999900003', crn: '1010010000', address: SELLER_ADDRESS },
  buyer: null, supplyDate: '2026-03-10', paymentMeansCode: '10',
  lines: [{ id: 1, name: 'قلم <أزرق>', quantity: '3', unitCode: 'PCE', netAmount: '100.00', vatCategory: 'S', vatRate: '0.15', vatAmount: '15.00' }],
  taxSubtotals: [{ vatCategory: 'S', vatRate: '0.15', taxableAmount: '100.00', taxAmount: '15.00' }],
  totals: { lineExtension: '100.00', taxExclusive: '100.00', taxAmount: '15.00', taxInclusive: '115.00', payable: '115.00' },
  ...over,
});

describe('e-invoice building blocks', () => {
  it('encodes and decodes the TLV QR code, Arabic included', () => {
    const qr = encodeQr({ sellerName: 'مؤسسة الشيوخ', vatNumber: '399999999900003', timestamp: '2026-03-10T14:30:00', total: '115.00', vatTotal: '15.00' });
    const tags = decodeQr(qr);
    expect(tags.get(1)!.toString('utf8')).toBe('مؤسسة الشيوخ');
    expect([2, 3, 4, 5].map((t) => tags.get(t)!.toString())).toEqual(['399999999900003', '2026-03-10T14:30:00', '115.00', '15.00']);
  });

  it('builds a CSR that openssl accepts, with ZATCA template and SAN fields per environment', () => {
    const dir = mkdtempSync(join(tmpdir(), 'csr-'));
    for (const [environment, template] of [['DEVELOPER', 'TSTZATCA-Code-Signing'], ['SIMULATION', 'PREZATCA-Code-Signing'], ['PRODUCTION', 'ZATCA-Code-Signing']] as const) {
      const csr = buildCsr({ privateKeyPem: generateKeyPair().privateKeyPem, environment, commonName: 'EGS-1', organizationUnit: 'الفرع الرئيسي', organization: 'مؤسسة الشيوخ',
        country: 'SA', egsSerial: '1-ALSHUYUKH|2-1.0|3-x', vatNumber: '399999999900003', invoiceTypes: '1100', registeredAddress: 'Riyadh', businessCategory: 'Trading' });
      writeFileSync(join(dir, 'r.csr'), csr);
      const text = execFileSync('openssl', ['req', '-in', join(dir, 'r.csr'), '-noout', '-verify', '-text'], { stdio: ['ignore', 'pipe', 'pipe'] }).toString();
      expect(text).toContain('ASN1 OID: secp256k1');
      expect(text).toContain(template);
      expect(text).toContain('UID=399999999900003/title=1100');
    }
  });

  it('hashes the invoice without signature or QR, and the hash follows the content', () => {
    const k = generateKeyPair();
    const fake = new FakeZatca();
    const token = fake.issue(buildCsr({ privateKeyPem: k.privateKeyPem, environment: 'DEVELOPER', commonName: 'c', organizationUnit: 'u', organization: 'o', country: 'SA',
      egsSerial: '1-a|2-b|3-c', vatNumber: '399999999900003', invoiceTypes: '1100', registeredAddress: 'r', businessCategory: 'b' }));
    const cert = parseCertificate(certificateFromCsid(token));
    const signed = signInvoice(sample(), { privateKeyPem: k.privateKeyPem, cert }, '2026-03-10T14:30:00');
    expect(signed.hash).toBe(invoiceHash(buildInvoiceXml(sample())));
    expect(invoiceHash(signed.xml)).toBe(signed.hash);
    expect(invoiceHash(buildInvoiceXml(sample({ number: 'INV-000002' })))).not.toBe(signed.hash);
    expect(fake.verifyInvoice({ invoiceHash: signed.hash, uuid: sample().uuid, invoice: Buffer.from(signed.xml).toString('base64') })).toEqual([]);
    // Tampering is detected: a changed amount breaks the hash, a foreign signature fails verification.
    const tampered = signed.xml.replace('<cbc:PayableAmount currencyID="SAR">115.00', '<cbc:PayableAmount currencyID="SAR">15.00');
    expect(fake.verifyInvoice({ invoiceHash: signed.hash, uuid: sample().uuid, invoice: Buffer.from(tampered).toString('base64') })).toContain('hash mismatch');
    const other = signInvoice(sample({ number: 'X' }), { privateKeyPem: generateKeyPair().privateKeyPem, cert }, '2026-03-10T14:30:00');
    expect(fake.verifyInvoice({ invoiceHash: other.hash, uuid: sample().uuid, invoice: Buffer.from(other.xml).toString('base64') })).toContain('bad signature');
    // Simplified QR: tags 1–9; tag 8 is the public key, tag 9 the CA's signature.
    const tags = decodeQr(signed.qr);
    expect([...tags.keys()]).toEqual([1, 2, 3, 4, 5, 6, 7, 8, 9]);
    expect(tags.get(6)!.toString()).toBe(signed.hash);
    expect(verify('sha256', Buffer.from(signed.hash, 'base64'), createPublicKey({ key: tags.get(8)!, format: 'der', type: 'spki' }), Buffer.from(tags.get(7)!.toString(), 'base64'))).toBe(true);
    // Escaping keeps the XML well-formed.
    const doc = new DOMParser().parseFromString(signed.xml, 'text/xml');
    expect(doc.getElementsByTagNameNS(UBL_CBC, 'Name')[0]!.textContent).toBe('قلم <أزرق>');
    // Standard invoices carry no tag 9.
    const std = signInvoice(sample({ subtype: '0100000' }), { privateKeyPem: k.privateKeyPem, cert }, '2026-03-10T14:30:00');
    expect(decodeQr(std.qr).has(9)).toBe(false);
  });

  it('expresses inexact unit prices through BaseQuantity so quantity × price = net', () => {
    const xml = buildInvoiceXml(sample());
    expect(xml).toContain('<cbc:PriceAmount currencyID="SAR">100.00</cbc:PriceAmount><cbc:BaseQuantity unitCode="PCE">3</cbc:BaseQuantity>');
    const exact = buildInvoiceXml(sample({ lines: [{ ...sample().lines[0]!, quantity: '4', netAmount: '100.00' }] }));
    expect(exact).toContain('<cbc:PriceAmount currencyID="SAR">25.00</cbc:PriceAmount></cac:Price>');
  });

  it('encrypts secrets with authentication', () => {
    const key = Buffer.alloc(32, 7);
    const sealed = encrypt(key, 'private key');
    expect(sealed).not.toContain('private');
    expect(decrypt(key, sealed)).toBe('private key');
    const tampered = sealed.slice(0, -4) + (sealed.endsWith('AAAA') ? 'BBBB' : 'AAAA');
    expect(() => decrypt(key, tampered)).toThrow();
  });
});

let t: TestContext;
let owner: Session;
let c: Commerce;
let fake: FakeZatca;
let deviceId: string;
let companyId: string;

const issue = async (customerId: string, quantity = '1', extra: Record<string, unknown> = {}) => {
  const inv = await c.invoice(customerId, [{ productId: (await c.product({ salePrice: '100' })).id, quantity }], extra);
  return c.api.post(`/api/invoices/${inv.id}/post`);
};
const einvoiceOf = async (type: string, id: string) => (await c.api.get(`/api/zatca/document?type=${type}&id=${id}`)).json();
const db = <T>(fn: (d: Parameters<Parameters<typeof withTx>[2]>[0]) => Promise<T>) => withTx(t.pool, { tenantId: owner.tenantId, userId: owner.userId }, fn);

beforeAll(async () => {
  t = await setupApp();
  fake = new FakeZatca();
  setTransport(fake.transport);
  owner = await register(t.app);
  c = await commerce(t.app, owner);
  companyId = (await c.api.get('/api/companies')).json().data[0].id;
});
afterAll(async () => { setTransport(null); await t.close(); });

describe('onboarding', () => {
  it('requires the seller data ZATCA needs before creating a unit', async () => {
    const res = await c.api.post('/api/zatca/devices', { name: 'الوحدة الرئيسية', environment: 'DEVELOPER', businessCategory: 'تجارة' });
    expect(res.statusCode).toBe(400);
    expect(res.json().error.code).toBe('ZATCA_INVALID');
    expect(res.json().error.details.map((d: { path: string }) => d.path)).toEqual(expect.arrayContaining(['company.vatNumber', 'company.commercialRegistration', 'company.buildingNumber', 'company.postalCode']));
  });

  it('creates the unit, gets the compliance CSID, passes the checks and activates', async () => {
    expect((await c.api.patch(`/api/companies/${companyId}`, {
      vatNumber: '399999999900003', commercialRegistration: '1010010000', ...SELLER_ADDRESS, additionalNumber: '8888',
    })).statusCode).toBe(200);
    const created = await c.api.post('/api/zatca/devices', { name: 'الوحدة الرئيسية', environment: 'DEVELOPER', businessCategory: 'تجارة' });
    expect(created.statusCode).toBe(201);
    const device = created.json();
    deviceId = device.id;
    expect(device).toMatchObject({ status: 'NEW', invoiceTypes: '1100', icv: '0', lastHash: INITIAL_PIH });
    expect(device.csr).toContain('BEGIN CERTIFICATE REQUEST');
    expect(JSON.stringify(device)).not.toMatch(/PRIVATE KEY|private_key/);
    const stored = await db((d) => d.query(`SELECT private_key_enc FROM zatca_devices WHERE id = $1`, [deviceId]));
    expect(stored.rows[0].private_key_enc).toMatch(/^v1\./);

    expect((await c.api.post(`/api/zatca/devices/${deviceId}/compliance-csid`, { otp: '999999' })).json().error.code).toBe('ZATCA_REJECTED');
    expect((await c.api.post(`/api/zatca/devices/${deviceId}/activate`)).json().error.code).toBe('DEVICE_STATE');
    const compliance = await c.api.post(`/api/zatca/devices/${deviceId}/compliance-csid`, { otp: '123456' });
    expect(compliance.json().status).toBe('COMPLIANCE');
    expect((await c.api.post(`/api/zatca/devices/${deviceId}/activate`)).json().error.code).toBe('COMPLIANCE_INCOMPLETE');

    const checks = (await c.api.post(`/api/zatca/devices/${deviceId}/compliance-checks`)).json();
    expect(checks.complianceChecks).toEqual({
      STANDARD_INVOICE: 'PASSED', STANDARD_CREDIT_NOTE: 'PASSED', STANDARD_DEBIT_NOTE: 'PASSED',
      SIMPLIFIED_INVOICE: 'PASSED', SIMPLIFIED_CREDIT_NOTE: 'PASSED', SIMPLIFIED_DEBIT_NOTE: 'PASSED',
    });
    // Each sample was signed with the compliance certificate and verified by the fake gateway.
    expect(fake.calls.filter((x) => x.path === '/compliance/invoices')).toHaveLength(6);

    const active = (await c.api.post(`/api/zatca/devices/${deviceId}/activate`)).json();
    expect(active).toMatchObject({ status: 'ACTIVE', certificateIssuer: 'CN=TSZEINVOICE-SubCA-1, DC=extgazt, DC=gov, DC=local' });
    // The submission log never holds certificates or secrets.
    const log = await db((d) => d.query(`SELECT response::text FROM zatca_submissions WHERE device_id = $1`, [deviceId]));
    expect(log.rows.length).toBe(9);
    expect(log.rows.every((r) => !/secret|binarySecurityToken/.test(r.response))).toBe(true);
    // One live unit per company.
    expect((await c.api.post('/api/zatca/devices', { name: 'ثانية', environment: 'DEVELOPER', businessCategory: 'تجارة' })).statusCode).toBe(409);
  });
});

describe('issuing e-invoices', () => {
  let first: { id: string; hash: string };

  it('signs a simplified invoice at issue, chained from the initial hash', async () => {
    const cust = (await c.customer()).id; // no VAT number → simplified
    const res = await issue(cust);
    expect(res.statusCode).toBe(200);
    const doc = await einvoiceOf('SALES_INVOICE', res.json().id);
    expect(doc.einvoice).toMatchObject({ status: 'PENDING', invoiceKind: 'SIMPLIFIED', typeCode: '388', subtype: '0200000', icv: '1', previousHash: INITIAL_PIH });
    expect(doc.qrSvg).toContain('<svg');
    expect([...decodeQr(doc.qr).keys()]).toEqual([1, 2, 3, 4, 5, 6, 7, 8, 9]);
    expect(decodeQr(doc.qr).get(4)!.toString()).toBe('115.00');
    first = { id: doc.einvoice.id, hash: doc.einvoice.invoiceHash };
    const xml = (await c.api.get(`/api/zatca/invoices/${first.id}/xml`));
    expect(xml.headers['content-type']).toContain('application/xml');
    expect(fake.verifyInvoice({ invoiceHash: first.hash, uuid: doc.einvoice.uuid, invoice: Buffer.from(xml.body).toString('base64') })).toEqual([]);
  });

  it('chains the next invoice to the previous hash', async () => {
    const res = await issue((await c.customer()).id);
    const doc = await einvoiceOf('SALES_INVOICE', res.json().id);
    expect(doc.einvoice).toMatchObject({ icv: '2', previousHash: first.hash });
  });

  it('reports a simplified invoice and then refuses to cancel it', async () => {
    const res = await issue((await c.customer()).id);
    const z = (await einvoiceOf('SALES_INVOICE', res.json().id)).einvoice;
    const sent = (await c.api.post(`/api/zatca/invoices/${z.id}/submit`)).json();
    expect(sent).toMatchObject({ status: 'REPORTED', hasWarnings: false, attempts: 1 });
    expect(sent.submissions).toEqual([expect.objectContaining({ operation: 'REPORTING', httpStatus: 200, outcome: 'SUCCESS' })]);
    expect(fake.calls.at(-1)!.headers['Clearance-Status']).toBe('0');
    expect((await c.api.post(`/api/zatca/invoices/${z.id}/submit`)).json().error.code).toBe('NOT_PENDING');
    const cancel = await c.api.post(`/api/invoices/${res.json().id}/cancel`, { reason: 'خطأ' });
    expect(cancel.json().error.code).toBe('ZATCA_ISSUED');

    // The correction is a credit note, which joins the chain as type 381.
    const inv = res.json();
    const ret = (await c.api.post('/api/sales-returns', { originalInvoiceId: inv.id, docDate: c.date, reason: 'إرجاع', lines: [{ sourceItemId: inv.lines[0].id, quantity: '1' }] })).json();
    const posted = await c.api.post(`/api/sales-returns/${ret.id}/post`);
    expect(posted.json().invoiceKind).toBe('SIMPLIFIED');
    const note = (await einvoiceOf('SALES_RETURN', ret.id)).einvoice;
    expect(note).toMatchObject({ typeCode: '381', subtype: '0200000', status: 'PENDING' });
    const xml = (await c.api.get(`/api/zatca/invoices/${note.id}/xml`)).body;
    expect(xml).toContain(`<cac:BillingReference><cac:InvoiceDocumentReference><cbc:ID>${inv.number}</cbc:ID>`);
    expect(xml).toContain('<cbc:InstructionNote>إرجاع</cbc:InstructionNote>');
    expect((await c.api.post(`/api/zatca/invoices/${note.id}/submit`)).json().status).toBe('REPORTED');
  });

  it('clears a standard invoice and keeps ZATCA\'s cleared XML', async () => {
    const vatCustomer = (await c.customer({ vatNumber: '311111111111113', addresses: [{ ...SELLER_ADDRESS, buildingNumber: '1111' }] })).id;
    const res = await issue(vatCustomer);
    const z = (await einvoiceOf('SALES_INVOICE', res.json().id)).einvoice;
    expect(z).toMatchObject({ invoiceKind: 'STANDARD', subtype: '0100000' });
    const sent = (await c.api.post(`/api/zatca/invoices/${z.id}/submit`)).json();
    expect(sent).toMatchObject({ status: 'CLEARED' });
    expect(sent.documents.map((d: { kind: string }) => d.kind)).toEqual(['SIGNED', 'CLEARED']);
    expect(fake.calls.at(-1)!.path).toBe('/invoices/clearance/single');
    expect((await c.api.get(`/api/zatca/invoices/${z.id}/xml?kind=CLEARED`)).body).toContain('cleared by fake ZATCA');
  });

  it('refuses to issue a standard invoice without the buyer\'s address — nothing is posted', async () => {
    const vatCustomer = (await c.customer({ vatNumber: '322222222222223' })).id;
    const inv = await c.invoice(vatCustomer, [{ productId: (await c.product()).id, quantity: '1' }]);
    const res = await c.api.post(`/api/invoices/${inv.id}/post`);
    expect(res.statusCode).toBe(400);
    expect(res.json().error.code).toBe('ZATCA_INVALID');
    expect(res.json().error.details.map((d: { path: string }) => d.path)).toContain('customer.address.street');
    const after = (await c.api.get(`/api/invoices/${inv.id}`)).json();
    expect(after).toMatchObject({ status: 'DRAFT', number: null, journalEntryId: null });
  });

  it('requires an exemption reason for zero-rated and exempt items', async () => {
    const cust = (await c.customer()).id;
    const zero = (await c.product({ vatCategory: 'Z', salePrice: '50' })).id;
    const inv = await c.invoice(cust, [{ productId: zero, quantity: '1' }]);
    expect((await c.api.post(`/api/invoices/${inv.id}/post`)).json().error.code).toBe('ZATCA_INVALID');
    await c.api.patch(`/api/products/${zero}`, { vatExemptionCode: 'VATEX-SA-35', vatExemptionReason: 'أدوية' });
    const ok = await c.api.post(`/api/invoices/${inv.id}/post`);
    expect(ok.statusCode).toBe(200);
    const z = (await einvoiceOf('SALES_INVOICE', inv.id)).einvoice;
    const xml = (await c.api.get(`/api/zatca/invoices/${z.id}/xml`)).body;
    expect(xml).toContain('<cbc:ID>Z</cbc:ID><cbc:Percent>0.00</cbc:Percent><cbc:TaxExemptionReasonCode>VATEX-SA-35</cbc:TaxExemptionReasonCode>');
  });

  it('records rejections with their messages; a rejected invoice may be cancelled', async () => {
    const res = await issue((await c.customer()).id);
    const z = (await einvoiceOf('SALES_INVOICE', res.json().id)).einvoice;
    fake.behaviour = 'reject';
    const sent = (await c.api.post(`/api/zatca/invoices/${z.id}/submit`)).json();
    fake.behaviour = 'ok';
    expect(sent).toMatchObject({ status: 'REJECTED' });
    expect(sent.lastError).toContain('BR-KSA-37');
    expect(sent.submissions[0].messages).toEqual([expect.objectContaining({ level: 'ERROR', code: 'BR-KSA-37' })]);
    expect((await c.api.post(`/api/invoices/${res.json().id}/cancel`, { reason: 'رفضتها الهيئة' })).statusCode).toBe(200);
  });

  it('keeps a failed submission pending with backoff, and marks warnings', async () => {
    const res = await issue((await c.customer()).id);
    const z = (await einvoiceOf('SALES_INVOICE', res.json().id)).einvoice;
    for (const b of ['network-error', 'server-error'] as const) {
      fake.behaviour = b;
      const sent = (await c.api.post(`/api/zatca/invoices/${z.id}/submit`)).json();
      expect(sent.status).toBe('PENDING');
      expect(new Date(sent.nextAttemptAt).getTime()).toBeGreaterThan(Date.now());
    }
    fake.behaviour = 'warning';
    const sent = (await c.api.post(`/api/zatca/invoices/${z.id}/submit`)).json();
    fake.behaviour = 'ok';
    expect(sent).toMatchObject({ status: 'REPORTED', hasWarnings: true, attempts: 3, lastError: null });
    expect(sent.submissions.map((s: { outcome: string }) => s.outcome)).toEqual(['TRANSPORT_ERROR', 'TRANSPORT_ERROR', 'WARNING']);
  });

  it('serialises concurrent issues into one unbroken chain', async () => {
    const custs = await Promise.all(Array.from({ length: 6 }, () => c.customer()));
    const results = await Promise.all(custs.map((x) => issue(x.id)));
    expect(results.every((r) => r.statusCode === 200)).toBe(true);
    const chain = await db((d) => d.query<{ icv: string; invoice_hash: string; previous_hash: string }>(
      `SELECT icv, invoice_hash, previous_hash FROM invoice_hashes WHERE device_id = $1 ORDER BY icv`, [deviceId]));
    chain.rows.forEach((r, i) => {
      expect(Number(r.icv)).toBe(i + 1);
      expect(r.previous_hash).toBe(i === 0 ? INITIAL_PIH : chain.rows[i - 1]!.invoice_hash);
    });
    const dev = (await c.api.get(`/api/zatca/devices/${deviceId}`)).json();
    expect(dev).toMatchObject({ icv: String(chain.rows.length), lastHash: chain.rows.at(-1)!.invoice_hash });
  });

  it('protects generated e-invoices and the hash chain in the database', async () => {
    await expect(db((d) => d.query(`UPDATE zatca_invoices SET invoice_hash = 'x' WHERE device_id = $1`, [deviceId]))).rejects.toThrow(/cannot be modified/);
    await expect(db((d) => d.query(`UPDATE invoice_hashes SET invoice_hash = 'x' WHERE device_id = $1`, [deviceId]))).rejects.toThrow();
    await expect(db((d) => d.query(`DELETE FROM zatca_documents`))).rejects.toThrow();
    await expect(db((d) => d.query(`UPDATE zatca_invoices SET status = 'PENDING' WHERE status = 'REPORTED'`))).rejects.toThrow(/cannot change status/);
  });

  it('lists e-invoices with a status summary', async () => {
    const list = (await c.api.get('/api/zatca/invoices?limit=200')).json();
    expect(list.total).toBeGreaterThan(5);
    expect(list.summary).toMatchObject({ rejected: 1, overdue: 0 });
    expect(list.data[0].qr).toBeUndefined();
    expect((await c.api.get('/api/zatca/invoices?status=CLEARED')).json().total).toBe(1);
  });
});

describe('background submitter', () => {
  it('reports every due e-invoice in its tenant context', async () => {
    const res = await issue((await c.customer()).id);
    const z = (await einvoiceOf('SALES_INVOICE', res.json().id)).einvoice;
    const warnings: unknown[] = [];
    const attempted = await runDueSubmissions(t.pool, { warn: (...a: unknown[]) => warnings.push(a) } as never);
    expect(attempted).toBeGreaterThanOrEqual(1);
    expect(warnings).toEqual([]);
    expect((await einvoiceOf('SALES_INVOICE', res.json().id)).einvoice.status).toBe('REPORTED');
    expect(z.status).toBe('PENDING');
    // Nothing is due any more: a second pass attempts nothing for this tenant.
    expect((await c.api.get('/api/zatca/invoices?status=PENDING')).json().total).toBe(0);
  });
});

describe('without an active unit', () => {
  it('still prints a phase-1 QR and issues normally', async () => {
    const s = await register(t.app);
    const o = await commerce(t.app, s);
    const cid = (await o.api.get('/api/companies')).json().data[0].id;
    await o.api.patch(`/api/companies/${cid}`, { vatNumber: '399999999900003' });
    const inv = await o.postInvoice((await o.customer()).id, [{ productId: (await o.product()).id, quantity: '2' }]);
    const doc = (await o.api.get(`/api/zatca/document?type=SALES_INVOICE&id=${inv.id}`)).json();
    expect(doc.einvoice).toBeNull();
    expect([...decodeQr(doc.qr).keys()]).toEqual([1, 2, 3, 4, 5]);
    expect(decodeQr(doc.qr).get(4)!.toString()).toBe('230.00');
    expect((await o.api.post(`/api/invoices/${inv.id}/cancel`, { reason: 'تجربة' })).statusCode).toBe(200);
  });
});

describe('permissions and isolation', () => {
  it('separates viewing, managing and submitting', async () => {
    const accountant = client(t.app, (await addMember(t.app, owner, 'ACCOUNTANT')).token);
    expect((await accountant.get('/api/zatca/invoices')).statusCode).toBe(200);
    expect((await accountant.post('/api/zatca/devices', { name: 'x', environment: 'DEVELOPER', businessCategory: 'تجارة' })).statusCode).toBe(403);
    expect((await accountant.post(`/api/zatca/devices/${deviceId}/revoke`, { reason: 'test' })).statusCode).toBe(403);
    const clerk = client(t.app, (await addMember(t.app, owner, 'SALES_EMPLOYEE')).token);
    expect((await clerk.get('/api/zatca/invoices')).statusCode).toBe(403);
  });

  it('hides units and e-invoices from other tenants', async () => {
    const z = (await c.api.get('/api/zatca/invoices?limit=1')).json().data[0];
    const other = client(t.app, (await register(t.app)).token);
    expect((await other.get(`/api/zatca/invoices/${z.id}`)).statusCode).toBe(404);
    expect((await other.get(`/api/zatca/invoices/${z.id}/xml`)).statusCode).toBe(404);
    expect((await other.post(`/api/zatca/invoices/${z.id}/submit`)).statusCode).toBe(404);
    expect((await other.get(`/api/zatca/devices/${deviceId}`)).statusCode).toBe(404);
    expect((await other.get('/api/zatca/invoices')).json().total).toBe(0);
  });

  it('revokes a unit; new invoices are then issued without e-invoice', async () => {
    expect((await c.api.post(`/api/zatca/devices/${deviceId}/revoke`, { reason: 'استبدال الجهاز' })).json().status).toBe('REVOKED');
    const res = await issue((await c.customer()).id);
    expect((await einvoiceOf('SALES_INVOICE', res.json().id)).einvoice).toBeNull();
  });
});
