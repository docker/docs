import { createHash, randomUUID } from 'node:crypto';
import type { Db } from '../../db/tx.js';
import { AppError, badRequest, conflict, notFound } from '../../lib/errors.js';
import { writeAudit, type AuditMeta } from '../audit/audit.service.js';
import { clientFor, validationMessages, type Credentials, type Environment, type ZatcaResponse } from './client.js';
import { buildCsr, certificateFromCsid, decrypt, encrypt, generateKeyPair, parseCertificate } from './crypto.js';
import { buildEInvoice, loadSeller, sampleInvoices, sellerProblems, type ZatcaDocType } from './einvoice.js';
import { INITIAL_PIH, signInvoice } from './sign.js';

export interface Ctx { tenantId: string; userId: string | null; meta?: AuditMeta }
/** Runs a function in its own short tenant transaction; network calls happen between them. */
export type Run = <T>(fn: (db: Db) => Promise<T>) => Promise<T>;

let secretKey: Buffer | null = null;
export const configureZatca = (key: Buffer) => { secretKey = key; };
const key = () => {
  if (!secretKey) throw new AppError(500, 'ZATCA_NOT_CONFIGURED', 'ZATCA encryption key is not configured');
  return secretKey;
};

/** Riyadh wall-clock date and time, as used on Saudi invoices. */
export function riyadhNow(d = new Date()) {
  const parts = new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh', year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit', hourCycle: 'h23' })
    .formatToParts(d).reduce<Record<string, string>>((a, p) => ({ ...a, [p.type]: p.value }), {});
  return { date: `${parts.year}-${parts.month}-${parts.day}`, time: `${parts.hour}:${parts.minute}:${parts.second}` };
}

export const DEVICE_SELECT = `SELECT d.id, d.company_id AS "companyId", d.branch_id AS "branchId", b.name AS "branchName", d.name, d.environment,
  d.egs_serial AS "egsSerial", d.common_name AS "commonName", d.organization_unit AS "organizationUnit", d.invoice_types AS "invoiceTypes",
  d.registered_address AS "registeredAddress", d.business_category AS "businessCategory", d.status, d.csr,
  d.compliance_checks AS "complianceChecks", d.certificate_serial AS "certificateSerial", d.certificate_issuer AS "certificateIssuer",
  d.certificate_expires_at AS "certificateExpiresAt", d.icv, d.last_hash AS "lastHash", d.activated_at AS "activatedAt",
  d.revoked_at AS "revokedAt", d.created_at AS "createdAt"
  FROM zatca_devices d LEFT JOIN branches b ON b.id = d.branch_id`;

interface DeviceRow {
  id: string; company_id: string; environment: Environment; status: string; csr: string; invoice_types: string;
  private_key_enc: string; compliance_csid: string | null; compliance_secret_enc: string | null; compliance_request_id: string | null;
  compliance_checks: Record<string, string>; production_csid: string | null; production_secret_enc: string | null; icv: string; last_hash: string;
}

async function deviceRow(db: Db, tenantId: string, id: string, lock = false): Promise<DeviceRow> {
  const { rows: [d] } = await db.query<DeviceRow>(`SELECT * FROM zatca_devices WHERE tenant_id = $1 AND id = $2${lock ? ' FOR UPDATE' : ''}`, [tenantId, id]);
  if (!d) throw notFound('Device');
  return d;
}

export async function loadDevice(db: Db, tenantId: string, id: string) {
  const { rows: [d] } = await db.query(`${DEVICE_SELECT} WHERE d.tenant_id = $1 AND d.id = $2`, [tenantId, id]);
  if (!d) throw notFound('Device');
  return d;
}

// --- Onboarding ------------------------------------------------------------------------

export interface DeviceInput {
  companyId: string; branchId?: string | null; name: string; environment: Environment;
  businessCategory: string; registeredAddress?: string | null; invoiceTypes?: string;
}

/** Creates an EGS unit: key pair (encrypted at rest) and the CSR to submit with the portal OTP. */
export async function createDevice(db: Db, ctx: Ctx, input: DeviceInput) {
  const seller = await loadSeller(db, input.companyId);
  const problems = sellerProblems(seller);
  if (problems.length) throw badRequest('ZATCA_INVALID', 'أكمل بيانات الشركة قبل ربطها بالفوترة الإلكترونية', problems);
  let unit = 'الفرع الرئيسي';
  if (input.branchId) {
    const { rows: [b] } = await db.query<{ name: string }>(`SELECT name FROM branches WHERE company_id = $1 AND id = $2`, [input.companyId, input.branchId]);
    if (!b) throw badRequest('INVALID_BRANCH', 'Branch not found in this company');
    unit = b.name;
  }
  const serial = randomUUID();
  const keys = generateKeyPair();
  const csrInput = {
    privateKeyPem: keys.privateKeyPem, environment: input.environment,
    commonName: `EGS-${serial.slice(0, 8)}`, organizationUnit: unit, organization: seller.name, country: seller.country,
    egsSerial: `1-ALSHUYUKH|2-1.0|3-${serial}`, vatNumber: seller.vatNumber!, invoiceTypes: input.invoiceTypes ?? '1100',
    registeredAddress: input.registeredAddress ?? `${seller.buildingNumber} ${seller.street}, ${seller.district}, ${seller.city} ${seller.postalCode}`,
    businessCategory: input.businessCategory,
  };
  const csr = buildCsr(csrInput);
  const { rows: [row] } = await db.query<{ id: string }>(
    `INSERT INTO zatca_devices (tenant_id, company_id, branch_id, name, environment, egs_serial, common_name, organization_unit, invoice_types,
       registered_address, business_category, private_key_enc, public_key, csr, created_by)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`,
    [ctx.tenantId, input.companyId, input.branchId ?? null, input.name, input.environment, csrInput.egsSerial, csrInput.commonName,
     unit, csrInput.invoiceTypes, csrInput.registeredAddress, input.businessCategory, encrypt(key(), keys.privateKeyPem), keys.publicKeyPem, csr, ctx.userId]);
  const created = await loadDevice(db, ctx.tenantId, row!.id);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'CREATE', entityType: 'zatca_device', entityId: row!.id,
    newValues: { name: input.name, environment: input.environment, egsSerial: csrInput.egsSerial } }, ctx.meta);
  return created;
}

type Operation = 'COMPLIANCE_CSID' | 'COMPLIANCE_CHECK' | 'PRODUCTION_CSID' | 'REPORTING' | 'CLEARANCE';

/** Never store certificates, secrets or the (large) cleared invoice in the submissions log. */
function sanitize(body: Record<string, any> | null) {
  if (!body) return null;
  const { binarySecurityToken: _t, secret: _s, clearedInvoice, ...rest } = body;
  return { ...rest, ...(clearedInvoice ? { clearedInvoice: '[stored as document]' } : {}) };
}

async function recordSubmission(db: Db, ctx: Ctx, deviceId: string, zatcaInvoiceId: string | null, operation: Operation, res: ZatcaResponse, outcome: 'SUCCESS' | 'WARNING' | 'REJECTED' | 'TRANSPORT_ERROR') {
  const { rows: [s] } = await db.query<{ id: string }>(
    `INSERT INTO zatca_submissions (tenant_id, device_id, zatca_invoice_id, operation, http_status, outcome, response, duration_ms, created_by)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
    [ctx.tenantId, deviceId, zatcaInvoiceId, operation, res.status || null, outcome,
     JSON.stringify(res.transportError ? { transportError: res.transportError } : sanitize(res.body)), res.durationMs, ctx.userId]);
  const messages = validationMessages(res.body);
  for (const m of messages) {
    await db.query(`INSERT INTO zatca_errors (tenant_id, submission_id, zatca_invoice_id, level, code, category, message) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
      [ctx.tenantId, s!.id, zatcaInvoiceId, m.level, m.code, m.category, m.message.slice(0, 2000)]);
  }
  return messages;
}

const errorSummary = (res: ZatcaResponse) => res.transportError
  ? `تعذر الاتصال بمنصة فاتورة: ${res.transportError}`
  : validationMessages(res.body).filter((m) => m.level === 'ERROR').map((m) => `${m.code ?? ''} ${m.message}`.trim()).join(' | ') || `HTTP ${res.status}`;

/** Step 1: exchange the CSR and the portal OTP for a compliance CSID. */
export async function requestComplianceCsid(run: Run, ctx: Ctx, deviceId: string, otp: string) {
  const d = await run((db) => deviceRow(db, ctx.tenantId, deviceId));
  if (d.status !== 'NEW' && d.status !== 'COMPLIANCE') throw conflict('DEVICE_STATE', 'This unit is already active or revoked');
  const res = await clientFor(d.environment).complianceCsid(d.csr, otp);
  const ok = res.status === 200 && typeof res.body?.binarySecurityToken === 'string' && typeof res.body?.secret === 'string';
  await run(async (db) => {
    await recordSubmission(db, ctx, deviceId, null, 'COMPLIANCE_CSID', res, ok ? 'SUCCESS' : res.transportError ? 'TRANSPORT_ERROR' : 'REJECTED');
    if (!ok) return;
    await db.query(
      `UPDATE zatca_devices SET status = 'COMPLIANCE', compliance_csid = $2, compliance_secret_enc = $3, compliance_request_id = $4, compliance_checks = '{}'
        WHERE id = $1`, [deviceId, res.body!.binarySecurityToken, encrypt(key(), res.body!.secret), String(res.body!.requestID ?? '')]);
    await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'UPDATE', entityType: 'zatca_device', entityId: deviceId, newValues: { status: 'COMPLIANCE' } }, ctx.meta);
  });
  if (!ok) throw new AppError(502, 'ZATCA_REJECTED', errorSummary(res));
  return run((db) => loadDevice(db, ctx.tenantId, deviceId));
}

const credentials = (csid: string, secretEnc: string): Credentials => ({ username: csid, password: decrypt(key(), secretEnc) });
const isAccepted = (res: ZatcaResponse) => res.status === 200 || res.status === 202;

/** Step 2: sign and submit one sample of every document type the unit will issue. */
export async function runComplianceChecks(run: Run, ctx: Ctx, deviceId: string) {
  const d = await run((db) => deviceRow(db, ctx.tenantId, deviceId));
  if (d.status !== 'COMPLIANCE' || !d.compliance_csid || !d.compliance_secret_enc) throw conflict('DEVICE_STATE', 'Request the compliance CSID first');
  const seller = await run((db) => loadSeller(db, d.company_id));
  const cert = parseCertificate(certificateFromCsid(d.compliance_csid));
  const privateKeyPem = decrypt(key(), d.private_key_enc);
  const auth = credentials(d.compliance_csid, d.compliance_secret_enc);
  const now = riyadhNow();
  const results: Record<string, string> = {};
  let previousHash = INITIAL_PIH;
  for (const [i, sample] of sampleInvoices(seller, d.invoice_types, now).entries()) {
    const label = `${sample.subtype === '0100000' ? 'STANDARD' : 'SIMPLIFIED'}_${{ 388: 'INVOICE', 381: 'CREDIT_NOTE', 383: 'DEBIT_NOTE' }[sample.typeCode]}`;
    const uuid = randomUUID();
    const signed = signInvoice({ ...sample, uuid, icv: i + 1, previousHash }, { privateKeyPem, cert }, `${now.date}T${now.time}`);
    previousHash = signed.hash;
    const res = await clientFor(d.environment).complianceCheck(auth, { invoiceHash: signed.hash, uuid, invoice: Buffer.from(signed.xml).toString('base64') });
    const passed = isAccepted(res);
    results[label] = passed ? (res.status === 202 ? 'PASSED_WITH_WARNINGS' : 'PASSED') : 'FAILED';
    await run((db) => recordSubmission(db, ctx, deviceId, null, 'COMPLIANCE_CHECK', res, passed ? (res.status === 202 ? 'WARNING' : 'SUCCESS') : res.transportError ? 'TRANSPORT_ERROR' : 'REJECTED'));
  }
  await run((db) => db.query(`UPDATE zatca_devices SET compliance_checks = $2 WHERE id = $1`, [deviceId, JSON.stringify(results)]));
  return run((db) => loadDevice(db, ctx.tenantId, deviceId));
}

/** Step 3: once every sample passed, obtain the production CSID; the unit starts signing invoices. */
export async function activateDevice(run: Run, ctx: Ctx, deviceId: string) {
  const d = await run((db) => deviceRow(db, ctx.tenantId, deviceId));
  if (d.status !== 'COMPLIANCE' || !d.compliance_csid || !d.compliance_secret_enc) throw conflict('DEVICE_STATE', 'Request the compliance CSID first');
  const checks = Object.values(d.compliance_checks ?? {});
  if (!checks.length || checks.some((c) => c === 'FAILED')) throw conflict('COMPLIANCE_INCOMPLETE', 'Run the compliance checks until every sample passes');
  const res = await clientFor(d.environment).productionCsid(credentials(d.compliance_csid, d.compliance_secret_enc), d.compliance_request_id ?? '');
  const ok = res.status === 200 && typeof res.body?.binarySecurityToken === 'string' && typeof res.body?.secret === 'string';
  await run(async (db) => {
    await recordSubmission(db, ctx, deviceId, null, 'PRODUCTION_CSID', res, ok ? 'SUCCESS' : res.transportError ? 'TRANSPORT_ERROR' : 'REJECTED');
    if (!ok) return;
    const cert = parseCertificate(certificateFromCsid(res.body!.binarySecurityToken));
    await db.query(
      `UPDATE zatca_devices SET status = 'ACTIVE', production_csid = $2, production_secret_enc = $3, certificate_serial = $4,
              certificate_issuer = $5, certificate_expires_at = $6, activated_at = now() WHERE id = $1`,
      [deviceId, res.body!.binarySecurityToken, encrypt(key(), res.body!.secret), cert.serialNumber, cert.issuer, cert.validTo]);
    await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'UPDATE', entityType: 'zatca_device', entityId: deviceId,
      newValues: { status: 'ACTIVE', certificateSerial: cert.serialNumber, expiresAt: cert.validTo } }, ctx.meta);
  });
  if (!ok) throw new AppError(502, 'ZATCA_REJECTED', errorSummary(res));
  return run((db) => loadDevice(db, ctx.tenantId, deviceId));
}

export async function revokeDevice(db: Db, ctx: Ctx, deviceId: string, reason: string) {
  const d = await deviceRow(db, ctx.tenantId, deviceId, true);
  if (d.status === 'REVOKED') throw conflict('DEVICE_STATE', 'Already revoked');
  await db.query(`UPDATE zatca_devices SET status = 'REVOKED', revoked_at = now() WHERE id = $1`, [deviceId]);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'UPDATE', entityType: 'zatca_device', entityId: deviceId, newValues: { status: 'REVOKED', reason } }, ctx.meta);
  return loadDevice(db, ctx.tenantId, deviceId);
}

// --- Generation ------------------------------------------------------------------------

/**
 * Called inside the issuing transaction of a sales invoice or return. When
 * the company (or branch) has an active unit, the document gets the next ICV,
 * the previous hash, a signed XML and a QR code. The unit row is locked, so
 * concurrent issues are serialised and the chain cannot fork. Any missing
 * data throws, which rolls back the issue itself.
 */
export async function generateForDocument(db: Db, ctx: Ctx, type: ZatcaDocType, documentId: string, scope: { companyId: string; branchId: string | null }) {
  const { rows: [device] } = await db.query<DeviceRow>(
    `SELECT * FROM zatca_devices WHERE tenant_id = $1 AND company_id = $2 AND status = 'ACTIVE' AND (branch_id = $3 OR branch_id IS NULL)
      ORDER BY branch_id NULLS LAST LIMIT 1 FOR UPDATE`, [ctx.tenantId, scope.companyId, scope.branchId]);
  if (!device) return null;
  const icv = Number(device.icv) + 1;
  const now = riyadhNow();
  const uuid = randomUUID();
  const { invoice, kind, number } = await buildEInvoice(db, type, documentId, { uuid, icv, previousHash: device.last_hash, issueTime: now.time });
  const signed = signInvoice(invoice, { privateKeyPem: decrypt(key(), device.private_key_enc), cert: parseCertificate(certificateFromCsid(device.production_csid!)) }, `${now.date}T${now.time}`);

  const { rows: [z] } = await db.query<{ id: string }>(
    `INSERT INTO zatca_invoices (tenant_id, company_id, device_id, document_type, document_id, document_number, invoice_kind, type_code, subtype,
       uuid, icv, invoice_hash, previous_hash, issue_date, issue_time, qr)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id`,
    [ctx.tenantId, scope.companyId, device.id, type, documentId, number, kind, invoice.typeCode, invoice.subtype,
     uuid, icv, signed.hash, device.last_hash, invoice.issueDate, invoice.issueTime, signed.qr]);
  await db.query(`INSERT INTO zatca_documents (tenant_id, zatca_invoice_id, kind, content, sha256) VALUES ($1, $2, 'SIGNED', $3, $4)`,
    [ctx.tenantId, z!.id, signed.xml, createHash('sha256').update(signed.xml).digest('hex')]);
  await db.query(`INSERT INTO invoice_hashes (tenant_id, company_id, device_id, icv, invoice_hash, previous_hash, zatca_invoice_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
    [ctx.tenantId, scope.companyId, device.id, icv, signed.hash, device.last_hash, z!.id]);
  await db.query(`UPDATE zatca_devices SET icv = $2, last_hash = $3 WHERE id = $1`, [device.id, icv, signed.hash]);
  return z!.id;
}

/** An e-invoice in the chain can only be corrected by a credit note, unless ZATCA rejected it. */
export async function assertCancellable(db: Db, type: string, documentId: string) {
  const { rows: [z] } = await db.query<{ status: string }>(`SELECT status FROM zatca_invoices WHERE document_type = $1 AND document_id = $2`, [type, documentId]);
  if (z && z.status !== 'REJECTED') {
    throw conflict('ZATCA_ISSUED', 'هذا المستند صدر فاتورة إلكترونية؛ لا يمكن إلغاؤه، أصدر إشعارًا دائنًا (مرتجعًا) بدلًا من ذلك');
  }
}

// --- Submission ------------------------------------------------------------------------

const backoffMinutes = (attempts: number) => Math.min(2 ** attempts, 360);

/**
 * Reports (simplified) or clears (standard) one pending e-invoice. A short
 * lease on next_attempt_at keeps two submitters from sending it twice.
 * `due` = true is the background worker: it only takes rows whose time has come.
 */
export async function submitInvoice(run: Run, ctx: Ctx, zatcaInvoiceId: string, opts: { due?: boolean } = {}) {
  const lease = await run(async (db) => {
    const { rows: [z] } = await db.query<{ id: string; device_id: string; invoice_kind: string; uuid: string; invoice_hash: string; attempts: number }>(
      `UPDATE zatca_invoices SET next_attempt_at = now() + interval '5 minutes'
        WHERE tenant_id = $1 AND id = $2 AND status = 'PENDING' ${opts.due ? 'AND next_attempt_at <= now()' : ''}
        RETURNING id, device_id, invoice_kind, uuid, invoice_hash, attempts`, [ctx.tenantId, zatcaInvoiceId]);
    if (!z) return null;
    const { rows: [doc] } = await db.query<{ content: string }>(`SELECT content FROM zatca_documents WHERE zatca_invoice_id = $1 AND kind = 'SIGNED'`, [z.id]);
    const device = await deviceRow(db, ctx.tenantId, z.device_id);
    return { z, xml: doc!.content, device };
  });
  if (!lease) {
    const exists = await run((db) => db.query(`SELECT status FROM zatca_invoices WHERE tenant_id = $1 AND id = $2`, [ctx.tenantId, zatcaInvoiceId]));
    if (!exists.rowCount) throw notFound('E-invoice');
    if (opts.due) return null;
    throw conflict('NOT_PENDING', 'This e-invoice is not waiting for submission (or is being submitted right now)');
  }
  const { z, xml, device } = lease;
  if (!device.production_csid || !device.production_secret_enc) throw conflict('DEVICE_STATE', 'The unit has no production CSID');
  const clearance = z.invoice_kind === 'STANDARD';
  const client = clientFor(device.environment);
  const payload = { invoiceHash: z.invoice_hash, uuid: z.uuid, invoice: Buffer.from(xml).toString('base64') };
  const res = clearance ? await client.clear(credentials(device.production_csid, device.production_secret_enc), payload)
    : await client.report(credentials(device.production_csid, device.production_secret_enc), payload);

  await run(async (db) => {
    const accepted = isAccepted(res);
    const rejected = res.status === 400 || res.status === 409 || res.status === 422;
    const outcome = accepted ? (res.status === 202 ? 'WARNING' : 'SUCCESS') : rejected ? 'REJECTED' : 'TRANSPORT_ERROR';
    const messages = await recordSubmission(db, ctx, device.id, z.id, clearance ? 'CLEARANCE' : 'REPORTING', res, outcome);
    if (accepted) {
      await db.query(
        `UPDATE zatca_invoices SET status = $2, has_warnings = $3, attempts = attempts + 1, submitted_at = now(), last_error = NULL WHERE id = $1`,
        [z.id, clearance ? 'CLEARED' : 'REPORTED', res.status === 202 || messages.some((m) => m.level === 'WARNING')]);
      if (clearance && typeof res.body?.clearedInvoice === 'string') {
        const cleared = Buffer.from(res.body.clearedInvoice, 'base64').toString('utf8');
        await db.query(`INSERT INTO zatca_documents (tenant_id, zatca_invoice_id, kind, content, sha256) VALUES ($1, $2, 'CLEARED', $3, $4) ON CONFLICT DO NOTHING`,
          [ctx.tenantId, z.id, cleared, createHash('sha256').update(cleared).digest('hex')]);
      }
    } else if (rejected) {
      await db.query(`UPDATE zatca_invoices SET status = 'REJECTED', attempts = attempts + 1, submitted_at = now(), last_error = $2 WHERE id = $1`, [z.id, errorSummary(res)]);
    } else {
      // Network failure, 401/403 or 5xx: try again later with exponential backoff.
      await db.query(
        `UPDATE zatca_invoices SET attempts = attempts + 1, last_error = $2, next_attempt_at = now() + make_interval(mins => $3) WHERE id = $1`,
        [z.id, errorSummary(res), backoffMinutes(z.attempts + 1)]);
    }
  });
  return run((db) => loadZatcaInvoice(db, ctx.tenantId, z.id));
}

export const ZATCA_INVOICE_SELECT = `SELECT z.id, z.company_id AS "companyId", z.device_id AS "deviceId", z.document_type AS "documentType",
  z.document_id AS "documentId", z.document_number AS "documentNumber", z.invoice_kind AS "invoiceKind", z.type_code AS "typeCode",
  z.subtype, z.uuid, z.icv, z.invoice_hash AS "invoiceHash", z.previous_hash AS "previousHash", z.issue_date AS "issueDate",
  z.issue_time AS "issueTime", z.qr, z.status, z.has_warnings AS "hasWarnings", z.attempts, z.next_attempt_at AS "nextAttemptAt",
  z.last_error AS "lastError", z.submitted_at AS "submittedAt", z.created_at AS "createdAt",
  (z.status = 'PENDING' AND z.invoice_kind = 'SIMPLIFIED' AND z.created_at < now() - interval '24 hours') AS "reportingOverdue"
  FROM zatca_invoices z`;

export async function loadZatcaInvoice(db: Db, tenantId: string, id: string) {
  const { rows: [z] } = await db.query(`${ZATCA_INVOICE_SELECT} WHERE z.tenant_id = $1 AND z.id = $2`, [tenantId, id]);
  if (!z) throw notFound('E-invoice');
  const { rows: submissions } = await db.query(
    `SELECT s.id, s.operation, s.http_status AS "httpStatus", s.outcome, s.duration_ms AS "durationMs", s.created_at AS "createdAt",
            COALESCE((SELECT json_agg(json_build_object('level', e.level, 'code', e.code, 'category', e.category, 'message', e.message) ORDER BY e.level)
                        FROM zatca_errors e WHERE e.submission_id = s.id), '[]') AS messages
       FROM zatca_submissions s WHERE s.zatca_invoice_id = $1 ORDER BY s.created_at`, [id]);
  const { rows: docs } = await db.query(`SELECT kind, sha256, created_at AS "createdAt" FROM zatca_documents WHERE zatca_invoice_id = $1 ORDER BY created_at`, [id]);
  return { ...z, submissions, documents: docs };
}
