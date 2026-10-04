import type { FastifyInstance, FastifyRequest } from 'fastify';
import { z } from 'zod';
import { isoDate, todayIn } from '../../lib/dates.js';
import { badRequest, notFound } from '../../lib/errors.js';
import { amountString } from '../../lib/money.js';
import { parse, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from '../accounting/company-context.js';
import { writeAudit } from '../audit/audit.service.js';
import { CalcError, totalsOf } from './calc.js';
import { buildLines, buildReturn, createDraft, createReturnDraft, deleteDraft, headerSelect, loadDocument, updateDraft, type Ctx } from './drafts.js';
import { KINDS, partyColumn, type DocKind } from './kinds.js';
import { PAYMENT_SELECT, allocate, createPayment, loadPayment, voidPayment } from './payments.js';
import { cancelDocument, convertDocument, issueDocument, setStatus } from './posting.js';

const quantity = z.string().trim().regex(/^\d{1,14}(\.\d{1,4})?$/, 'Quantity: up to 4 decimals, sent as a string');
const price = z.string().trim().regex(/^\d{1,14}(\.\d{1,4})?$/, 'Price: up to 4 decimals, sent as a string');
const percent = z.string().trim().regex(/^\d{1,3}(\.\d{1,2})?$/);

const line = z.object({
  productId: z.uuid().nullish(),
  accountId: z.uuid().nullish(),
  description: z.string().trim().max(500).nullish(),
  quantity,
  unitPrice: price.optional(),
  discountAmount: amountString.optional(),
  discountPercent: percent.optional(),
  unitId: z.uuid().nullish(),
  vatCategory: z.enum(['S', 'Z', 'E', 'O']).optional(),
});

const draftFields = {
  companyId: z.uuid().optional(),
  partyId: z.uuid(),
  docDate: isoDate,
  dueDate: isoDate.nullish(),
  validUntil: isoDate.nullish(),
  expectedDate: isoDate.nullish(),
  supplierInvoiceNumber: z.string().trim().min(1).max(60).nullish(),
  branchId: z.uuid().nullish(),
  warehouseId: z.uuid().nullish(),
  pricesIncludeVat: z.boolean().optional(),
  notes: z.string().trim().max(2000).nullish(),
  lines: z.array(line).min(1).max(500),
};
const draftBody = z.object(draftFields);
const draftUpdate = z.object(draftFields).omit({ companyId: true }).partial();

const returnBody = z.object({
  originalInvoiceId: z.uuid(),
  docDate: isoDate,
  reason: z.string().trim().min(3).max(500),
  notes: z.string().trim().max(2000).nullish(),
  lines: z.array(z.object({ sourceItemId: z.uuid(), quantity })).min(1).max(500),
});
const returnUpdate = returnBody.omit({ originalInvoiceId: true }).partial();

const listQuery = z.object({
  companyId: z.uuid().optional(),
  status: z.string().regex(/^[A-Z_]+(,[A-Z_]+)*$/).optional(),
  partyId: z.uuid().optional(),
  dateFrom: isoDate.optional(),
  dateTo: isoDate.optional(),
  open: z.enum(['true', 'false']).optional(),
  search: z.string().trim().max(60).optional(),
  limit: z.coerce.number().int().min(1).max(200).default(50),
  offset: z.coerce.number().int().min(0).default(0),
});

const ctxOf = (req: FastifyRequest): Ctx => ({ tenantId: req.auth!.tenantId, userId: req.auth!.userId, meta: req.auditMeta() });

const asCalcError = (e: unknown) => (e instanceof CalcError ? badRequest('INVALID_LINE', e.message) : e);

function documentRoutes(kind: DocKind) {
  return async (app: FastifyInstance) => {
    const base = `/${kind.path}`;
    const can = (p: keyof DocKind['perm']) => requirePermission(app, kind.perm[p]);

    app.get(base, { preHandler: can('view') }, async (req) => {
      const q = parse(listQuery, req.query);
      return req.tenantTx(async (db) => {
        const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
        const statuses = q.status?.split(',') ?? null;
        const search = q.search ? `%${q.search.replace(/[%_\\]/g, (m) => `\\${m}`)}%` : null;
        const hasRemaining = kind.legal;
        const { rows } = await db.query<{ total_count: string }>(
          `${headerSelect(kind).replace('SELECT d.id', 'SELECT count(*) OVER () AS total_count, d.id')}
            WHERE d.tenant_id = $1 AND d.company_id = $2 AND d.deleted_at IS NULL
              AND ($3::text[] IS NULL OR d.status = ANY($3))
              AND ($4::uuid IS NULL OR d.${partyColumn(kind)} = $4)
              AND ($5::date IS NULL OR d.doc_date >= $5) AND ($6::date IS NULL OR d.doc_date <= $6)
              AND (NOT $7::boolean ${hasRemaining ? `OR (d.status NOT IN ('DRAFT', 'CANCELLED') AND d.remaining_amount > 0)` : ''})
              AND ($8::text IS NULL OR d.doc_number ILIKE $8 OR p.name_ar ILIKE $8)
            ORDER BY d.doc_date DESC, d.doc_number DESC NULLS FIRST, d.created_at DESC
            LIMIT $9 OFFSET $10`,
          [req.auth!.tenantId, companyId, statuses, q.partyId ?? null, q.dateFrom ?? null, q.dateTo ?? null,
           q.open === 'true', search, q.limit, q.offset]);
        return { data: rows.map(({ total_count: _t, ...r }) => r), total: Number(rows[0]?.total_count ?? 0) };
      });
    });

    app.get(`${base}/:id`, { preHandler: can('view') }, async (req) => {
      const { id } = parse(uuidParam, req.params);
      return req.tenantTx((db) => loadDocument(db, kind, req.auth!.tenantId, id));
    });

    /** Server-side totals for the form, without saving anything. */
    app.post(`${base}/calculate`, { preHandler: can('view') }, async (req) => {
      return req.tenantTx(async (db) => {
        try {
          if (kind.isReturn) {
            const body = parse(returnBody.extend({ reason: z.string().optional() }), req.body);
            const r = await buildReturn(db, kind, req.auth!.tenantId, { ...body, reason: body.reason ?? '' }, null);
            return { lines: r.lines, totals: r.totals };
          }
          const body = parse(draftBody.extend({ partyId: z.uuid().optional() }), req.body);
          const companyId = await resolveCompanyId(db, req.auth!.tenantId, body.companyId);
          const lines = await buildLines(db, kind, companyId, body);
          return { lines, totals: totalsOf(lines) };
        } catch (e) { throw asCalcError(e); }
      });
    });

    app.post(base, { preHandler: can('create') }, async (req, reply) => {
      const ctx = ctxOf(req);
      const doc = await req.tenantTx(async (db) => {
        let id: string;
        if (kind.isReturn) {
          id = await createReturnDraft(db, kind, ctx, parse(returnBody, req.body));
        } else {
          const body = parse(draftBody, req.body);
          const companyId = await resolveCompanyId(db, ctx.tenantId, body.companyId);
          id = await createDraft(db, kind, ctx, { ...body, companyId });
        }
        return loadDocument(db, kind, ctx.tenantId, id);
      });
      reply.code(201);
      return doc;
    });

    app.patch(`${base}/:id`, { preHandler: can('edit') }, async (req) => {
      const { id } = parse(uuidParam, req.params);
      const body = kind.isReturn ? parse(returnUpdate, req.body) : parse(draftUpdate, req.body);
      return req.tenantTx((db) => updateDraft(db, kind, ctxOf(req), id, body));
    });

    app.delete(`${base}/:id`, { preHandler: can('delete') }, async (req, reply) => {
      const { id } = parse(uuidParam, req.params);
      await req.tenantTx((db) => deleteDraft(db, kind, ctxOf(req), id));
      return reply.code(204).send();
    });

    app.post(`${base}/:id/cancel`, { preHandler: can('cancel') }, async (req) => {
      const { id } = parse(uuidParam, req.params);
      const { reason } = parse(z.object({ reason: z.string().trim().min(3).max(500) }), req.body);
      return req.tenantTx((db) => cancelDocument(db, kind, ctxOf(req), id, reason));
    });

    if (kind.legal) {
      app.post(`${base}/:id/post`, { preHandler: can('issue') }, async (req) => {
        const { id } = parse(uuidParam, req.params);
        return req.tenantTx(async (db) => {
          try { return await issueDocument(db, kind, ctxOf(req), id); } catch (e) { throw asCalcError(e); }
        });
      });
    } else {
      app.post(`${base}/:id/status`, { preHandler: can(kind.key === 'PURCHASE_ORDER' ? 'issue' : 'edit') }, async (req) => {
        const { id } = parse(uuidParam, req.params);
        const { status } = parse(z.object({ status: z.enum(['SENT', 'ACCEPTED', 'REJECTED', 'APPROVED']) }), req.body);
        return req.tenantTx((db) => setStatus(db, kind, ctxOf(req), id, status));
      });
      app.post(`${base}/:id/convert`, { preHandler: can('create') }, async (req, reply) => {
        const { id } = parse(uuidParam, req.params);
        const body = parse(z.object({ docDate: isoDate.optional() }), req.body ?? {});
        const invoice = await req.tenantTx((db) => convertDocument(db, kind, ctxOf(req), id, body.docDate ?? todayIn('Asia/Riyadh')));
        reply.code(201);
        return invoice;
      });
    }
  };
}

const paymentBody = z.object({
  companyId: z.uuid().optional(),
  direction: z.enum(['RECEIPT', 'DISBURSEMENT']),
  customerId: z.uuid().nullish(),
  supplierId: z.uuid().nullish(),
  paymentDate: isoDate,
  methodId: z.uuid(),
  amount: amountString.refine((v) => Number(v) > 0, 'Amount must be greater than zero'),
  reference: z.string().trim().max(100).nullish(),
  notes: z.string().trim().max(1000).nullish(),
  branchId: z.uuid().nullish(),
  allocations: z.array(z.object({
    documentType: z.enum(['SALES_INVOICE', 'SALES_RETURN', 'PURCHASE_INVOICE', 'PURCHASE_RETURN', 'EXPENSE']),
    documentId: z.uuid(),
    amount: amountString,
  })).max(200).default([]),
});

async function paymentRoutes(app: FastifyInstance) {
  app.get('/payments', { preHandler: requirePermission(app, 'payment.view') }, async (req) => {
    const q = parse(z.object({
      companyId: z.uuid().optional(),
      direction: z.enum(['RECEIPT', 'DISBURSEMENT']).optional(),
      customerId: z.uuid().optional(),
      supplierId: z.uuid().optional(),
      status: z.enum(['POSTED', 'VOIDED']).optional(),
      limit: z.coerce.number().int().min(1).max(200).default(50),
      offset: z.coerce.number().int().min(0).default(0),
    }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query<{ total_count: string }>(
        `${PAYMENT_SELECT.replace('SELECT p.id', 'SELECT count(*) OVER () AS total_count, p.id')}
          WHERE p.tenant_id = $1 AND p.company_id = $2 AND ($3::text IS NULL OR p.direction = $3)
            AND ($4::uuid IS NULL OR p.customer_id = $4) AND ($5::uuid IS NULL OR p.supplier_id = $5)
            AND ($6::text IS NULL OR p.status = $6)
          ORDER BY p.payment_date DESC, p.payment_number DESC LIMIT $7 OFFSET $8`,
        [req.auth!.tenantId, companyId, q.direction ?? null, q.customerId ?? null, q.supplierId ?? null, q.status ?? null, q.limit, q.offset]);
      return { data: rows.map(({ total_count: _t, ...r }) => r), total: Number(rows[0]?.total_count ?? 0) };
    });
  });

  app.get('/payments/:id', { preHandler: requirePermission(app, 'payment.view') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx((db) => loadPayment(db, req.auth!.tenantId, id));
  });

  app.post('/payments', { preHandler: requirePermission(app, 'payment.create') }, async (req, reply) => {
    const body = parse(paymentBody, req.body);
    const payment = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, body.companyId);
      return createPayment(db, ctxOf(req), { ...body, companyId });
    });
    reply.code(201);
    return payment;
  });

  app.post('/payments/:id/allocations', { preHandler: requirePermission(app, 'payment.create') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const { allocations } = parse(paymentBody.pick({ allocations: true }), req.body);
    if (!allocations.length) throw badRequest('NO_ALLOCATIONS', 'Provide at least one allocation');
    return req.tenantTx((db) => allocate(db, ctxOf(req), id, allocations));
  });

  app.post('/payments/:id/void', { preHandler: requirePermission(app, 'payment.void') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const { reason } = parse(z.object({ reason: z.string().trim().min(3).max(500) }), req.body);
    return req.tenantTx((db) => voidPayment(db, ctxOf(req), id, reason));
  });

  // Payment methods ---------------------------------------------------------------
  const METHOD_SELECT = `SELECT m.id, m.code, m.name_ar AS "nameAr", m.method_type AS "methodType", m.account_id AS "accountId",
      a.code AS "accountCode", a.name_ar AS "accountName", m.is_active AS "isActive", m.sort_order AS "sortOrder"
    FROM payment_methods m JOIN accounts a ON a.id = m.account_id`;

  app.get('/payment-methods', { preHandler: requirePermission(app, 'payment.view') }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      return { data: (await db.query(`${METHOD_SELECT} WHERE m.tenant_id = $1 AND m.company_id = $2 ORDER BY m.sort_order, m.code`, [req.auth!.tenantId, companyId])).rows };
    });
  });

  const methodAccountCheck = async (db: Parameters<Parameters<FastifyRequest['tenantTx']>[0]>[0], companyId: string, accountId: string) => {
    const { rows: [a] } = await db.query<{ account_type: string; is_postable: boolean; is_active: boolean }>(
      `SELECT account_type, is_postable, is_active FROM accounts WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, accountId]);
    if (!a || a.account_type !== 'ASSET' || !a.is_postable || !a.is_active) {
      throw badRequest('INVALID_ACCOUNT', 'A payment method needs an active postable asset account (cash, bank or clearing)');
    }
  };

  app.post('/payment-methods', { preHandler: requirePermission(app, 'settings.manage') }, async (req, reply) => {
    const body = parse(z.object({
      companyId: z.uuid().optional(),
      code: z.string().trim().toUpperCase().regex(/^[A-Z][A-Z0-9_]{1,30}$/),
      nameAr: z.string().trim().min(2).max(100),
      methodType: z.enum(['CASH', 'BANK', 'CARD', 'STC_PAY', 'TAMARA', 'OTHER']),
      accountId: z.uuid(),
    }), req.body);
    const a = req.auth!;
    const method = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      await methodAccountCheck(db, companyId, body.accountId);
      const { rows: [m] } = await db.query<{ id: string }>(
        `INSERT INTO payment_methods (tenant_id, company_id, code, name_ar, method_type, account_id, sort_order)
         VALUES ($1, $2, $3, $4, $5, $6, 100) RETURNING id`, [a.tenantId, companyId, body.code, body.nameAr, body.methodType, body.accountId]);
      const created = (await db.query(`${METHOD_SELECT} WHERE m.id = $1`, [m!.id])).rows[0];
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'SETTINGS_CHANGE', entityType: 'payment_method', entityId: m!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return method;
  });

  app.patch('/payment-methods/:id', { preHandler: requirePermission(app, 'settings.manage') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ nameAr: z.string().trim().min(2).max(100), accountId: z.uuid(), isActive: z.boolean() }).partial(), req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const { rows: [before] } = await db.query<{ companyId: string }>(`${METHOD_SELECT.replace('SELECT m.id', 'SELECT m.company_id AS "companyId", m.id')} WHERE m.tenant_id = $1 AND m.id = $2`, [a.tenantId, id]);
      if (!before) throw notFound('Payment method');
      if (body.accountId) await methodAccountCheck(db, before.companyId, body.accountId);
      await db.query(`UPDATE payment_methods SET name_ar = COALESCE($2, name_ar), account_id = COALESCE($3, account_id), is_active = COALESCE($4, is_active) WHERE id = $1`,
        [id, body.nameAr ?? null, body.accountId ?? null, body.isActive ?? null]);
      const after = (await db.query(`${METHOD_SELECT} WHERE m.id = $1`, [id])).rows[0];
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'SETTINGS_CHANGE', entityType: 'payment_method', entityId: id, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });

  // Tax rates ---------------------------------------------------------------------
  // The VAT rate is data, not code: a new rate is added with an effective date
  // and the previous one is closed the day before. Phase 6 adds the VAT report.
  const RATE_SELECT = `SELECT id, vat_category AS "vatCategory", name_ar AS "nameAr", rate::text, effective_from AS "effectiveFrom",
      effective_to AS "effectiveTo", is_active AS "isActive" FROM tax_rates`;

  app.get('/tax-rates', { preHandler: requirePermission(app, 'invoice.view') }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      return { data: (await db.query(`${RATE_SELECT} WHERE tenant_id = $1 AND company_id = $2 ORDER BY effective_from DESC`, [req.auth!.tenantId, companyId])).rows };
    });
  });

  app.post('/tax-rates', { preHandler: requirePermission(app, 'tax.manage') }, async (req, reply) => {
    const body = parse(z.object({
      companyId: z.uuid().optional(),
      nameAr: z.string().trim().min(2).max(100),
      rate: z.string().regex(/^0\.\d{1,4}$/, 'Rate is a fraction, e.g. "0.15"'),
      effectiveFrom: isoDate,
    }), req.body);
    const a = req.auth!;
    const rate = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      // Close the rate that is open on the new start date.
      await db.query(
        `UPDATE tax_rates SET effective_to = $2::date - 1
          WHERE company_id = $1 AND vat_category = 'S' AND is_active AND effective_from < $2
            AND (effective_to IS NULL OR effective_to >= $2)`, [companyId, body.effectiveFrom]);
      const { rows: [r] } = await db.query<{ id: string }>(
        `INSERT INTO tax_rates (tenant_id, company_id, name_ar, rate, effective_from, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
        [a.tenantId, companyId, body.nameAr, body.rate, body.effectiveFrom, a.userId]);
      const created = (await db.query(`${RATE_SELECT} WHERE id = $1`, [r!.id])).rows[0];
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'SETTINGS_CHANGE', entityType: 'tax_rate', entityId: r!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return rate;
  });
}

export default async function documentsModule(app: FastifyInstance) {
  for (const kind of Object.values(KINDS)) await app.register(documentRoutes(kind));
  await app.register(paymentRoutes);
}
