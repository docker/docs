import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { isoDate } from '../../lib/dates.js';
import { forbidden } from '../../lib/errors.js';
import { amountString } from '../../lib/money.js';
import { parse, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from './company-context.js';
import { createDraft, deleteDraft, getEntry, post, reverse, updateDraft, type EngineContext } from './engine.js';

const line = z.object({
  accountId: z.uuid(),
  debit: amountString.optional(),
  credit: amountString.optional(),
  description: z.string().trim().max(500).nullish(),
  costCenterId: z.uuid().nullish(),
  branchId: z.uuid().nullish(),
  customerId: z.uuid().nullish(),
  supplierId: z.uuid().nullish(),
});
const lines = z.array(line).min(2).max(1000);

const createBody = z.object({
  companyId: z.uuid().optional(),
  entryDate: isoDate,
  description: z.string().trim().min(1).max(1000),
  lines,
  /** Post immediately (requires journal.post). */
  post: z.boolean().default(false),
});

const listQuery = z.object({
  companyId: z.uuid().optional(),
  status: z.enum(['DRAFT', 'POSTED', 'REVERSED']).optional(),
  dateFrom: isoDate.optional(),
  dateTo: isoDate.optional(),
  accountId: z.uuid().optional(),
  referenceType: z.string().regex(/^[A-Z_]+$/).optional(),
  search: z.string().trim().max(100).optional(),
  limit: z.coerce.number().int().min(1).max(200).default(50),
  offset: z.coerce.number().int().min(0).default(0),
});

export default async function journalRoutes(app: FastifyInstance) {
  const ctxOf = (req: { auth: { tenantId: string; userId: string } | null; auditMeta: () => EngineContext['meta'] }): EngineContext =>
    ({ tenantId: req.auth!.tenantId, userId: req.auth!.userId, meta: req.auditMeta() });

  app.get('/journal-entries', { preHandler: requirePermission(app, 'journal.view') }, async (req) => {
    const q = parse(listQuery, req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query<{ total: string }>(
        `SELECT e.id, e.entry_number AS "entryNumber", to_char(e.entry_date, 'YYYY-MM-DD') AS "entryDate",
                e.description, e.reference_type AS "referenceType", e.source, e.status,
                CASE WHEN e.status = 'DRAFT' THEN COALESCE((SELECT sum(debit) FROM journal_entry_lines WHERE journal_entry_id = e.id), 0)
                     ELSE e.total_debit END::numeric(18,2)::text AS "totalDebit",
                e.posted_at AS "postedAt", e.created_at AS "createdAt",
                count(*) OVER () AS total
           FROM journal_entries e
          WHERE e.tenant_id = $1 AND e.company_id = $2 AND e.deleted_at IS NULL
            AND ($3::text IS NULL OR e.status = $3)
            AND ($4::date IS NULL OR e.entry_date >= $4)
            AND ($5::date IS NULL OR e.entry_date <= $5)
            AND ($6::uuid IS NULL OR EXISTS (SELECT 1 FROM journal_entry_lines l WHERE l.journal_entry_id = e.id AND l.account_id = $6))
            AND ($7::text IS NULL OR e.reference_type = $7)
            AND ($8::text IS NULL OR e.description ILIKE '%' || $8 || '%' OR e.entry_number ILIKE '%' || $8 || '%')
          ORDER BY e.entry_date DESC, e.entry_number DESC NULLS FIRST, e.created_at DESC
          LIMIT $9 OFFSET $10`,
        [req.auth!.tenantId, companyId, q.status ?? null, q.dateFrom ?? null, q.dateTo ?? null, q.accountId ?? null,
         q.referenceType ?? null, q.search?.replace(/[%_\\]/g, (m) => `\\${m}`) ?? null, q.limit, q.offset]);
      return { data: rows.map(({ total: _t, ...r }) => r), total: Number(rows[0]?.total ?? 0) };
    });
  });

  app.get('/journal-entries/:id', { preHandler: requirePermission(app, 'journal.view') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx((db) => getEntry(db, req.auth!.tenantId, id));
  });

  app.post('/journal-entries', { preHandler: requirePermission(app, 'journal.create') }, async (req, reply) => {
    const body = parse(createBody, req.body);
    if (body.post && !req.auth!.permissions.has('journal.post')) throw forbidden('Missing permission: journal.post');
    const ctx = ctxOf(req);
    const entry = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, ctx.tenantId, body.companyId);
      const id = await createDraft(db, ctx, { companyId, entryDate: body.entryDate, description: body.description, lines: body.lines });
      return body.post ? post(db, ctx, id) : getEntry(db, ctx.tenantId, id);
    });
    reply.code(201);
    return entry;
  });

  app.patch('/journal-entries/:id', { preHandler: requirePermission(app, 'journal.create') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ entryDate: isoDate, description: z.string().trim().min(1).max(1000), lines }).partial(), req.body);
    return req.tenantTx((db) => updateDraft(db, ctxOf(req), id, body));
  });

  app.delete('/journal-entries/:id', { preHandler: requirePermission(app, 'journal.create') }, async (req, reply) => {
    const { id } = parse(uuidParam, req.params);
    await req.tenantTx((db) => deleteDraft(db, ctxOf(req), id));
    return reply.code(204).send();
  });

  app.post('/journal-entries/:id/post', { preHandler: requirePermission(app, 'journal.post') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx((db) => post(db, ctxOf(req), id));
  });

  /**
   * Reverses a posted manual entry. With `correction`, also creates the
   * correcting entry as a draft linked to the original, in the same transaction.
   */
  app.post('/journal-entries/:id/reverse', { preHandler: requirePermission(app, 'journal.reverse') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({
      reason: z.string().trim().min(3).max(500),
      date: isoDate.optional(),
      correction: z.object({ entryDate: isoDate, description: z.string().trim().min(1).max(1000), lines }).optional(),
    }), req.body);
    if (body.correction && !req.auth!.permissions.has('journal.create')) throw forbidden('Missing permission: journal.create');
    const ctx = ctxOf(req);
    return req.tenantTx(async (db) => {
      const original = await getEntry(db, ctx.tenantId, id);
      const reversal = await reverse(db, ctx, id, { reason: body.reason, date: body.date });
      let correction = null;
      if (body.correction) {
        const cid = await createDraft(db, ctx, { companyId: original.companyId, ...body.correction, correctionOfId: id });
        correction = await getEntry(db, ctx.tenantId, cid);
      }
      return { original: await getEntry(db, ctx.tenantId, id), reversal, correction };
    });
  });
}
