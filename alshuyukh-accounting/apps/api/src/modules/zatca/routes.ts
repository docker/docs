import type { FastifyInstance, FastifyReply, FastifyRequest } from 'fastify';
import QRCode from 'qrcode';
import { z } from 'zod';
import type { Db } from '../../db/tx.js';
import { forbidden, notFound } from '../../lib/errors.js';
import { featuresOf } from '../subscriptions/service.js';
import { parse, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from '../accounting/company-context.js';
import { loadSeller } from './einvoice.js';
import { encodeQr } from './qr.js';
import {
  activateDevice, createDevice, DEVICE_SELECT, loadDevice, loadZatcaInvoice, requestComplianceCsid, revokeDevice, riyadhNow,
  runComplianceChecks, submitInvoice, ZATCA_INVOICE_SELECT, type Ctx, type Run,
} from './service.js';

const ctxOf = (req: FastifyRequest): Ctx => ({ tenantId: req.auth!.tenantId, userId: req.auth!.userId, meta: req.auditMeta() });
const runOf = (req: FastifyRequest): Run => (fn) => req.tenantTx(fn);
const qrSvg = (qr: string) => QRCode.toString(qr, { type: 'svg', errorCorrectionLevel: 'M', margin: 1 });

export default async function zatcaRoutes(app: FastifyInstance) {
  const canView = requirePermission(app, 'zatca.view');
  const canManage = requirePermission(app, 'zatca.manage');

  // Units (EGS) and onboarding ------------------------------------------------------------
  app.get('/zatca/devices', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      return { data: (await db.query(`${DEVICE_SELECT} WHERE d.tenant_id = $1 AND d.company_id = $2 ORDER BY d.created_at DESC`, [req.auth!.tenantId, companyId])).rows };
    });
  });

  app.get('/zatca/devices/:id', { preHandler: canView }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx((db) => loadDevice(db, req.auth!.tenantId, id));
  });

  app.post('/zatca/devices', { preHandler: canManage }, async (req, reply) => {
    const body = parse(z.object({
      companyId: z.uuid().optional(), branchId: z.uuid().nullish(),
      name: z.string().trim().min(2).max(100),
      environment: z.enum(['DEVELOPER', 'SIMULATION', 'PRODUCTION']),
      businessCategory: z.string().trim().min(2).max(100),
      registeredAddress: z.string().trim().min(3).max(200).nullish(),
      invoiceTypes: z.enum(['1100', '1000', '0100']).default('1100'),
    }), req.body);
    const device = await req.tenantTx(async (db) => {
      if (!(await featuresOf(db, req.auth!.tenantId)).zatca_einvoicing) throw forbidden('الفوترة الإلكترونية غير مفعلة لهذه المنشأة');
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, body.companyId);
      return createDevice(db, ctxOf(req), { ...body, companyId });
    });
    reply.code(201);
    return device;
  });

  app.post('/zatca/devices/:id/compliance-csid', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const { otp } = parse(z.object({ otp: z.string().regex(/^\d{6}$/, 'The OTP from the Fatoora portal is 6 digits') }), req.body);
    return requestComplianceCsid(runOf(req), ctxOf(req), id, otp);
  });

  app.post('/zatca/devices/:id/compliance-checks', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return runComplianceChecks(runOf(req), ctxOf(req), id);
  });

  app.post('/zatca/devices/:id/activate', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return activateDevice(runOf(req), ctxOf(req), id);
  });

  app.post('/zatca/devices/:id/revoke', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const { reason } = parse(z.object({ reason: z.string().trim().min(3).max(500) }), req.body);
    return req.tenantTx((db) => revokeDevice(db, ctxOf(req), id, reason));
  });

  // E-invoices ----------------------------------------------------------------------------
  app.get('/zatca/invoices', { preHandler: canView }, async (req) => {
    const q = parse(z.object({
      companyId: z.uuid().optional(), status: z.enum(['PENDING', 'REPORTED', 'CLEARED', 'REJECTED']).optional(),
      documentType: z.enum(['SALES_INVOICE', 'SALES_RETURN']).optional(),
      limit: z.coerce.number().int().min(1).max(200).default(50), offset: z.coerce.number().int().min(0).default(0),
    }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const where = `z.tenant_id = $1 AND z.company_id = $2 AND ($3::text IS NULL OR z.status = $3) AND ($4::text IS NULL OR z.document_type = $4)`;
      const params = [req.auth!.tenantId, companyId, q.status ?? null, q.documentType ?? null];
      const { rows } = await db.query(`${ZATCA_INVOICE_SELECT} WHERE ${where} ORDER BY z.created_at DESC LIMIT $5 OFFSET $6`, [...params, q.limit, q.offset]);
      const { rows: [c] } = await db.query<Record<string, string>>(
        `SELECT count(*) AS total, count(*) FILTER (WHERE status = 'PENDING') AS pending, count(*) FILTER (WHERE status = 'REJECTED') AS rejected,
                count(*) FILTER (WHERE status = 'PENDING' AND invoice_kind = 'SIMPLIFIED' AND created_at < now() - interval '24 hours') AS overdue
           FROM zatca_invoices z WHERE ${where}`, params);
      return { data: rows.map(({ qr: _q, ...r }) => r), total: Number(c!.total), summary: { pending: Number(c!.pending), rejected: Number(c!.rejected), overdue: Number(c!.overdue) } };
    });
  });

  app.get('/zatca/invoices/:id', { preHandler: canView }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx((db) => loadZatcaInvoice(db, req.auth!.tenantId, id));
  });

  app.get('/zatca/invoices/:id/xml', { preHandler: canView }, async (req, reply: FastifyReply) => {
    const { id } = parse(uuidParam, req.params);
    const { kind } = parse(z.object({ kind: z.enum(['SIGNED', 'CLEARED']).default('SIGNED') }), req.query);
    const doc = await req.tenantTx(async (db) => {
      const { rows: [d] } = await db.query<{ content: string; number: string }>(
        `SELECT d.content, z.document_number AS number FROM zatca_documents d JOIN zatca_invoices z ON z.id = d.zatca_invoice_id
          WHERE z.tenant_id = $1 AND z.id = $2 AND d.kind = $3`, [req.auth!.tenantId, id, kind]);
      if (!d) throw notFound('XML document');
      return d;
    });
    reply.header('Content-Type', 'application/xml; charset=utf-8');
    reply.header('Content-Disposition', `attachment; filename="${doc.number}-${kind.toLowerCase()}.xml"`);
    return doc.content;
  });

  // Whoever may issue invoices may push them to ZATCA (the obligation comes with issuing).
  app.post('/zatca/invoices/:id/submit', { preHandler: requirePermission(app, 'invoice.post') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return submitInvoice(runOf(req), ctxOf(req), id);
  });

  /**
   * E-invoicing data for printing a sales document: the signed e-invoice's QR
   * when it exists, otherwise a phase-1 QR (seller, VAT number, time, totals).
   */
  app.get('/zatca/document', { preHandler: requirePermission(app, 'invoice.view') }, async (req) => {
    const q = parse(z.object({ type: z.enum(['SALES_INVOICE', 'SALES_RETURN']), id: z.uuid() }), req.query);
    return req.tenantTx(async (db: Db) => {
      const { rows: [z_] } = await db.query<Record<string, any>>(`${ZATCA_INVOICE_SELECT} WHERE z.tenant_id = $1 AND z.document_type = $2 AND z.document_id = $3`,
        [req.auth!.tenantId, q.type, q.id]);
      if (z_) return { einvoice: z_, qr: z_.qr, qrSvg: await qrSvg(z_.qr) };
      const table = q.type === 'SALES_INVOICE' ? 'sales_invoices' : 'sales_returns';
      const { rows: [d] } = await db.query<{ company_id: string; doc_date: string; issued_at: Date | null; total: string; tax_amount: string; status: string }>(
        `SELECT company_id, doc_date, issued_at, total::text, tax_amount::text, status FROM ${table} WHERE tenant_id = $1 AND id = $2`, [req.auth!.tenantId, q.id]);
      if (!d) throw notFound('Document');
      if (d.status === 'DRAFT' || !d.issued_at) return { einvoice: null, qr: null, qrSvg: null };
      const seller = await loadSeller(db, d.company_id);
      if (!seller.vatNumber) return { einvoice: null, qr: null, qrSvg: null, missing: 'company.vatNumber' };
      const qr = encodeQr({ sellerName: seller.name, vatNumber: seller.vatNumber, timestamp: `${d.doc_date}T${riyadhNow(d.issued_at).time}`, total: d.total, vatTotal: d.tax_amount });
      return { einvoice: null, qr, qrSvg: await qrSvg(qr) };
    });
  });
}
