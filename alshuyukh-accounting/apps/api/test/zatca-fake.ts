import { execFileSync } from 'node:child_process';
import { createPublicKey, verify, X509Certificate } from 'node:crypto';
import { mkdtempSync, readFileSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { DOMParser } from '@xmldom/xmldom';
import { invoiceHash } from '../src/modules/zatca/sign.js';
import type { Transport } from '../src/modules/zatca/client.js';

/**
 * A stand-in for the ZATCA gateway. It issues real X.509 certificates from a
 * throw-away secp256k1 CA (via the openssl CLI), and on every invoice it
 * recomputes the hash and verifies the ECDSA signature with the certificate
 * embedded in the XML — so the tests check our cryptography, not just our flow.
 */
export interface Call { path: string; headers: Record<string, string>; body: any }
export type Behaviour = 'ok' | 'warning' | 'reject' | 'server-error' | 'network-error';

export class FakeZatca {
  readonly calls: Call[] = [];
  behaviour: Behaviour = 'ok';
  private dir = mkdtempSync(join(tmpdir(), 'zatca-ca-'));
  private serial = 1000;

  constructor() {
    execFileSync('openssl', ['ecparam', '-name', 'secp256k1', '-genkey', '-noout', '-out', join(this.dir, 'ca.key')]);
    execFileSync('openssl', ['req', '-x509', '-new', '-key', join(this.dir, 'ca.key'), '-subj', '/DC=local/DC=gov/DC=extgazt/CN=TSZEINVOICE-SubCA-1', '-days', '30', '-out', join(this.dir, 'ca.pem')]);
  }

  /** Signs a CSR; returns the token ZATCA would return (base64 of the base64 DER certificate). */
  issue(csrPem: string): string {
    const csr = join(this.dir, `req-${++this.serial}.csr`);
    const out = join(this.dir, `cert-${this.serial}.der`);
    writeFileSync(csr, csrPem);
    execFileSync('openssl', ['x509', '-req', '-in', csr, '-CA', join(this.dir, 'ca.pem'), '-CAkey', join(this.dir, 'ca.key'),
      '-set_serial', String(this.serial), '-days', '10', '-outform', 'DER', '-out', out]);
    return Buffer.from(readFileSync(out).toString('base64')).toString('base64');
  }

  /** Checks the hash and signature of a submitted invoice; returns problems found. */
  verifyInvoice(body: { invoiceHash: string; uuid: string; invoice: string }): string[] {
    const xml = Buffer.from(body.invoice, 'base64').toString('utf8');
    const problems: string[] = [];
    if (invoiceHash(xml) !== body.invoiceHash) problems.push('hash mismatch');
    const doc = new DOMParser().parseFromString(xml, 'text/xml');
    const text = (ns: string, name: string) => doc.getElementsByTagNameNS(ns, name)[0]?.textContent ?? '';
    const ds = 'http://www.w3.org/2000/09/xmldsig#';
    if (text('urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2', 'UUID') !== body.uuid) problems.push('uuid mismatch');
    if (text(ds, 'DigestValue') !== body.invoiceHash) problems.push('digest mismatch');
    const cert = new X509Certificate(Buffer.from(text(ds, 'X509Certificate'), 'base64'));
    const ok = verify('sha256', Buffer.from(body.invoiceHash, 'base64'), createPublicKey(cert.publicKey.export({ type: 'spki', format: 'pem' })), Buffer.from(text(ds, 'SignatureValue'), 'base64'));
    if (!ok) problems.push('bad signature');
    return problems;
  }

  transport: Transport = async (url, init) => {
    // Match on the endpoint so any base URL (real gateway or ZATCA_GATEWAY_URL) works.
    const path = url.match(/\/(compliance\/invoices|compliance|production\/csids|invoices\/reporting\/single|invoices\/clearance\/single)$/)?.[0] ?? url;
    const body = JSON.parse(init.body);
    this.calls.push({ path, headers: init.headers, body });
    const reply = (status: number, payload: unknown) => ({ status, text: async () => JSON.stringify(payload) });
    if (this.behaviour === 'network-error') throw new Error('connect ETIMEDOUT');
    if (this.behaviour === 'server-error') return reply(503, { message: 'Service unavailable' });

    if (path === '/compliance') {
      if (init.headers.OTP !== '123456') return reply(400, { errors: ['Invalid-OTP'] });
      const csr = Buffer.from(body.csr, 'base64').toString('utf8');
      return reply(200, { requestID: 1234567890123, dispositionMessage: 'ISSUED', binarySecurityToken: this.issue(csr), secret: 'compliance-secret' });
    }
    if (path === '/production/csids') {
      if (!init.headers.Authorization) return reply(401, { message: 'unauthorized' });
      const token = Buffer.from(init.headers.Authorization.slice(6), 'base64').toString('utf8').split(':')[0]!;
      // Re-issue for the same key: read the public key back from the compliance certificate.
      const cert = new X509Certificate(Buffer.from(Buffer.from(token, 'base64').toString('utf8'), 'base64'));
      return reply(200, { requestID: 99, binarySecurityToken: this.reissue(cert), secret: 'production-secret' });
    }
    const problems = this.verifyInvoice(body);
    if (problems.length || this.behaviour === 'reject') {
      return reply(400, {
        validationResults: { status: 'ERROR', infoMessages: [], warningMessages: [],
          errorMessages: (problems.length ? problems : ['BR-KSA-37 seller address']).map((p) => ({ type: 'ERROR', code: 'BR-KSA-37', category: 'KSA', message: p, status: 'ERROR' })) },
        reportingStatus: 'NOT_REPORTED', clearanceStatus: 'NOT_CLEARED',
      });
    }
    const warn = this.behaviour === 'warning';
    const validationResults = { status: warn ? 'WARNING' : 'PASS', infoMessages: [], errorMessages: [],
      warningMessages: warn ? [{ type: 'WARNING', code: 'BR-KSA-08', category: 'KSA', message: 'buyer identification', status: 'WARNING' }] : [] };
    if (path === '/compliance/invoices') return reply(warn ? 202 : 200, { validationResults, reportingStatus: 'REPORTED', clearanceStatus: 'CLEARED' });
    if (path === '/invoices/reporting/single') return reply(warn ? 202 : 200, { validationResults, reportingStatus: 'REPORTED' });
    if (path === '/invoices/clearance/single') {
      const cleared = Buffer.from(body.invoice, 'base64').toString('utf8').replace('</Invoice>', '<!-- cleared by fake ZATCA --></Invoice>');
      return reply(warn ? 202 : 200, { validationResults, clearanceStatus: 'CLEARED', clearedInvoice: Buffer.from(cleared).toString('base64') });
    }
    return reply(404, { message: 'not found' });
  };

  /** Issues a second certificate for the public key of an existing one (stands in for the production CSID). */
  private reissue(cert: X509Certificate): string {
    const key = join(this.dir, `pub-${++this.serial}.pem`);
    const out = join(this.dir, `prod-${this.serial}.der`);
    writeFileSync(key, cert.publicKey.export({ type: 'spki', format: 'pem' }) as string);
    execFileSync('openssl', ['x509', '-new', '-force_pubkey', key, '-subj', '/CN=EGS-production', '-CA', join(this.dir, 'ca.pem'),
      '-CAkey', join(this.dir, 'ca.key'), '-set_serial', String(this.serial), '-days', '10', '-outform', 'DER', '-out', out]);
    return Buffer.from(readFileSync(out).toString('base64')).toString('base64');
  }
}
