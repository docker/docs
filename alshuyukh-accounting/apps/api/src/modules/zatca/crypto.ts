import { createCipheriv, createDecipheriv, createHash, createPrivateKey, createPublicKey, generateKeyPairSync, hkdfSync, randomBytes, sign, X509Certificate } from 'node:crypto';
import { bits, children, ctx, int, octets, oid, printable, read, seq, set, utf8 } from './der.js';

/**
 * Key material for the cryptographic stamp. ZATCA requires ECDSA on the
 * secp256k1 curve with SHA-256.
 */
export function generateKeyPair() {
  const { privateKey, publicKey } = generateKeyPairSync('ec', { namedCurve: 'secp256k1' });
  return {
    privateKeyPem: privateKey.export({ type: 'sec1', format: 'pem' }) as string,
    publicKeyPem: publicKey.export({ type: 'spki', format: 'pem' }) as string,
  };
}

// --- Secrets at rest -------------------------------------------------------------------

/**
 * Private keys and CSID secrets are stored encrypted with AES-256-GCM. The key
 * comes from ZATCA_ENCRYPTION_KEY (base64, 32 bytes) or, in development, is
 * derived from JWT_SECRET with HKDF under a dedicated label.
 */
export function encryptionKey(env: { ZATCA_ENCRYPTION_KEY?: string; JWT_SECRET: string }): Buffer {
  if (env.ZATCA_ENCRYPTION_KEY) {
    const k = Buffer.from(env.ZATCA_ENCRYPTION_KEY, 'base64');
    if (k.length !== 32) throw new Error('ZATCA_ENCRYPTION_KEY must be 32 bytes, base64 encoded');
    return k;
  }
  return Buffer.from(hkdfSync('sha256', env.JWT_SECRET, 'alshuyukh', 'zatca-secrets-v1', 32));
}

export function encrypt(key: Buffer, plain: string): string {
  const iv = randomBytes(12);
  const c = createCipheriv('aes-256-gcm', key, iv);
  const body = Buffer.concat([c.update(plain, 'utf8'), c.final()]);
  return ['v1', iv.toString('base64'), c.getAuthTag().toString('base64'), body.toString('base64')].join('.');
}

export function decrypt(key: Buffer, sealed: string): string {
  const [v, iv, tag, body] = sealed.split('.');
  if (v !== 'v1' || !iv || !tag || !body) throw new Error('Unknown secret format');
  const d = createDecipheriv('aes-256-gcm', key, Buffer.from(iv, 'base64'));
  d.setAuthTag(Buffer.from(tag, 'base64'));
  return Buffer.concat([d.update(Buffer.from(body, 'base64')), d.final()]).toString('utf8');
}

// --- CSR -------------------------------------------------------------------------------

export interface CsrInput {
  privateKeyPem: string;
  environment: 'DEVELOPER' | 'SIMULATION' | 'PRODUCTION';
  commonName: string;        // CN
  organizationUnit: string;  // OU: branch name, or the 10-digit TIN for VAT groups
  organization: string;      // O: taxpayer name
  country: string;           // C
  egsSerial: string;         // SAN SN: 1-<solution>|2-<model>|3-<serial>
  vatNumber: string;         // SAN UID
  invoiceTypes: string;      // SAN title: 4 digits, e.g. 1100 = standard + simplified
  registeredAddress: string; // SAN registeredAddress
  businessCategory: string;  // SAN businessCategory
}

/** Certificate template name ZATCA expects in the CSR, per environment. */
export const TEMPLATE_NAME = { DEVELOPER: 'TSTZATCA-Code-Signing', SIMULATION: 'PREZATCA-Code-Signing', PRODUCTION: 'ZATCA-Code-Signing' } as const;

const rdn = (type: string, value: string) => set(seq(oid(type), utf8(value)));

/** Builds a PKCS#10 CSR (PEM) with the ZATCA certificate template and SAN fields. */
export function buildCsr(i: CsrInput): string {
  const key = createPrivateKey(i.privateKeyPem);
  const spki = createPublicKey(key).export({ type: 'spki', format: 'der' });
  const subject = seq(
    rdn('2.5.4.6', i.country), rdn('2.5.4.11', i.organizationUnit), rdn('2.5.4.10', i.organization), rdn('2.5.4.3', i.commonName),
  );
  const directoryName = seq(
    rdn('2.5.4.4', i.egsSerial),                      // SN (surname)
    rdn('0.9.2342.19200300.100.1.1', i.vatNumber),   // UID
    rdn('2.5.4.12', i.invoiceTypes),                  // title
    rdn('2.5.4.26', i.registeredAddress),             // registeredAddress
    rdn('2.5.4.15', i.businessCategory),              // businessCategory
  );
  const extensions = seq(
    seq(oid('1.3.6.1.4.1.311.20.2'), octets(printable(TEMPLATE_NAME[i.environment]))),
    seq(oid('2.5.29.17'), octets(seq(ctx(4, directoryName)))),
  );
  const attributes = ctx(0, seq(oid('1.2.840.113549.1.9.14'), set(extensions)));
  const info = seq(int(0), subject, Buffer.from(spki), attributes);
  const signature = sign('sha256', info, key); // DER-encoded ECDSA signature
  const csr = seq(info, seq(oid('1.2.840.10045.4.3.2')), bits(signature));
  return `-----BEGIN CERTIFICATE REQUEST-----\n${csr.toString('base64').replace(/(.{64})/g, '$1\n').trim()}\n-----END CERTIFICATE REQUEST-----\n`;
}

// --- Certificates ----------------------------------------------------------------------

/**
 * ZATCA returns the CSID as a base64 "binarySecurityToken" whose decoded text
 * is the base64 DER certificate (without PEM armour).
 */
export function certificateFromCsid(binarySecurityToken: string): string {
  return Buffer.from(binarySecurityToken, 'base64').toString('utf8').trim();
}

export interface CertInfo {
  /** Base64 DER, as embedded in the signature's X509Certificate. */
  base64: string;
  issuer: string;          // RFC 4514-style, most specific first (as ZATCA's SDK prints it)
  serialNumber: string;    // decimal
  validTo: Date;
  publicKeyDer: Buffer;    // SubjectPublicKeyInfo
  signature: Buffer;       // the CA's signature over the certificate (QR tag 9)
}

export function parseCertificate(base64Der: string): CertInfo {
  const der = Buffer.from(base64Der, 'base64');
  const cert = new X509Certificate(der);
  // Certificate ::= SEQUENCE { tbsCertificate, signatureAlgorithm, signatureValue BIT STRING }
  const [, , sigBits] = children(read(der));
  return {
    base64: base64Der,
    issuer: cert.issuer.split('\n').reverse().join(', '),
    serialNumber: BigInt(`0x${cert.serialNumber}`).toString(10),
    validTo: new Date(cert.validTo),
    publicKeyDer: cert.publicKey.export({ type: 'spki', format: 'der' }) as Buffer,
    signature: sigBits!.value.subarray(1),
  };
}

/** ZATCA's certificate digest: SHA-256 as lowercase hex, then base64 of that text. */
export const certificateHash = (base64Der: string) => Buffer.from(createHash('sha256').update(base64Der).digest('hex')).toString('base64');
