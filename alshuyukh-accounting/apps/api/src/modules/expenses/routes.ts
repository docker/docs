import type { FastifyInstance, FastifyRequest } from 'fastify';
import { z } from 'zod';
import { isoDate } from '../../lib/dates.js';
import { badRequest, forbidden, notFound } from '../../lib/errors.js';
import { amountString } from '../../lib/money.js';
import { parse, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from '../accounting/company-context.js';
import { writeAudit } from '../audit/audit.service.js';
import type { Ctx } from '../documents/drafts.js';
import { EXPENSE_SELECT, buildExpenseLines, cancelExpense, createExpense, deleteExpense, loadExpense, postExpense, updateExpense } from './service.js';

const line = z.object({
  categoryId: z.uuid(),
  amount: amountString,
  description: z.string().trim().max(500).nullish(),
  vatCategory: z.enum(['S', 'Z', 'E', 'O']).optional(),
  costCenterId: z.uuid().nullish(),
});
const fields = {
  companyId: z.uuid().optional(),
  expenseDate: isoDate,
  paymentType: z.enum(['CASH', 'BANK', 'CREDIT']),
  methodId: z.uuid().nullish(),
  supplierId: z.uuid().nullish(),
  payeeName: z.string().trim().max(200).nullish(),
  reference: z.string().trim().max(60).nullish(),
  vendorVatNumber: z.string().regex(/^3\d{13}3$/, 'Saudi VAT number must be 15 digits starting and ending with 3').nullish(),
  branchId: z.uuid().nullish(),
  pricesIncludeVat: z.boolean().optional(),
  notes: z.string().trim().max(2000).nullish(),
  lines: z.array(line).min(1).max(500),
};
const createBody = z.object({ ...fields, post: z.boolean().default(false) });
const updateBody = z.object(fields).omit({ companyId: true }).partial();

const ctxOf = (req: FastifyRequest): Ctx => ({ tenantId: req.auth!.tenantId, userId: req.auth!.userId, meta: req.auditMeta() });

export default async function expenseRoutes(app: FastifyInstance) {
  const canView = requirePermission(app, 'expense.view');

  // Categories ------------------------------------------------------------------
  const CAT_SELECT = `SELECT c.id, c.code, c.name_ar AS "nameAr", c.account_id AS "accountId", a.code AS "accountCode",
      a.name_ar AS "accountName", c.vat_category AS "vatCategory", c.is_active AS "isActive"
    FROM expense_categories c JOIN accounts a ON a.id = c.account_id`;

  const checkAccount = async (db: Parameters<Parameters<FastifyRequest['tenantTx']>[0]>[0], companyId: string, accountId: string) => {
    const { rows: [a] } = await db.query<{ account_type: string; is_postable: boolean; is_active: boolean }>(
      `SELECT account_type, is_postable, is_active FROM accounts WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, accountId]);
    if (!a || !a.is_postable || !a.is_active || !['EXPENSE', 'COST_OF_GOODS_SOLD', 'ASSET'].includes(a.account_type)) {
      throw badRequest('INVALID_ACCOUNT', 'An expense category needs an active postable expense, cost or asset account');
    }
  };

  app.get('/expense-categories', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      return { data: (await db.query(`${CAT_SELECT} WHERE c.tenant_id = $1 AND c.company_id = $2 ORDER BY c.name_ar`, [req.auth!.tenantId, companyId])).rows };
    });
  });

  app.post('/expense-categories', { preHandler: requirePermission(app, 'account.manage') }, async (req, reply) => {
    const body = parse(z.object({
      companyId: z.uuid().optional(), code: z.string().trim().toUpperCase().regex(/^[A-Z][A-Z0-9_]{1,30}$/),
      nameAr: z.string().trim().min(2).max(100), accountId: z.uuid(), vatCategory: z.enum(['S', 'Z', 'E', 'O']).default('S'),
    }), req.body);
    const a = req.auth!;
    const cat = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      await checkAccount(db, companyId, body.accountId);
      const { rows: [c] } = await db.query<{ id: string }>(
        `INSERT INTO expense_categories (tenant_id, company_id, code, name_ar, account_id, vat_category) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
        [a.tenantId, companyId, body.code, body.nameAr, body.accountId, body.vatCategory]);
      const created = (await db.query(`${CAT_SELECT} WHERE c.id = $1`, [c!.id])).rows[0];
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'expense_category', entityId: c!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return cat;
  });

  app.patch('/expense-categories/:id', { preHandler: requirePermission(app, 'account.manage') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ nameAr: z.string().trim().min(2).max(100), accountId: z.uuid(), vatCategory: z.enum(['S', 'Z', 'E', 'O']), isActive: z.boolean() }).partial(), req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const { rows: [before] } = await db.query<{ companyId: string }>(`${CAT_SELECT.replace('SELECT c.id', 'SELECT c.company_id AS "companyId", c.id')} WHERE c.tenant_id = $1 AND c.id = $2`, [a.tenantId, id]);
      if (!before) throw notFound('Expense category');
      if (body.accountId) await checkAccount(db, before.companyId, body.accountId);
      await db.query(`UPDATE expense_categories SET name_ar = COALESCE($2, name_ar), account_id = COALESCE($3, account_id),
          vat_category = COALESCE($4, vat_category), is_active = COALESCE($5, is_active) WHERE id = $1`,
        [id, body.nameAr ?? null, body.accountId ?? null, body.vatCategory ?? null, body.isActive ?? null]);
      const after = (await db.query(`${CAT_SELECT} WHERE c.id = $1`, [id])).rows[0];
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'expense_category', entityId: id, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });

  // Expenses --------------------------------------------------------------------
  app.get('/expenses', { preHandler: canView }, async (req) => {
    const q = parse(z.object({
      companyId: z.uuid().optional(), status: z.string().regex(/^[A-Z_]+(,[A-Z_]+)*$/).optional(), supplierId: z.uuid().optional(),
      open: z.enum(['true', 'false']).optional(), dateFrom: isoDate.optional(), dateTo: isoDate.optional(),
      limit: z.coerce.number().int().min(1).max(200).default(50), offset: z.coerce.number().int().min(0).default(0),
    }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query<{ total_count: string }>(
        `${EXPENSE_SELECT.replace('SELECT e.id', 'SELECT count(*) OVER () AS total_count, e.id')}
          WHERE e.tenant_id = $1 AND e.company_id = $2 AND e.deleted_at IS NULL
            AND ($3::text[] IS NULL OR e.status = ANY($3)) AND ($4::uuid IS NULL OR e.supplier_id = $4)
            AND (NOT $5::boolean OR (e.status IN ('POSTED', 'PARTIALLY_PAID') AND e.remaining_amount > 0))
            AND ($6::date IS NULL OR e.expense_date >= $6) AND ($7::date IS NULL OR e.expense_date <= $7)
          ORDER BY e.expense_date DESC, e.expense_number DESC NULLS FIRST, e.created_at DESC LIMIT $8 OFFSET $9`,
        [req.auth!.tenantId, companyId, q.status?.split(',') ?? null, q.supplierId ?? null, q.open === 'true',
         q.dateFrom ?? null, q.dateTo ?? null, q.limit, q.offset]);
      return { data: rows.map(({ total_count: _t, ...r }) => r), total: Number(rows[0]?.total_count ?? 0) };
    });
  });

  app.get('/expenses/:id', { preHandler: canView }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx((db) => loadExpense(db, req.auth!.tenantId, id));
  });

  app.post('/expenses/calculate', { preHandler: canView }, async (req) => {
    const body = parse(z.object({ companyId: z.uuid().optional(), expenseDate: isoDate, pricesIncludeVat: z.boolean().optional(), lines: fields.lines }), req.body);
    return req.tenantTx(async (db) => buildExpenseLines(db, await resolveCompanyId(db, req.auth!.tenantId, body.companyId), body));
  });

  app.post('/expenses', { preHandler: requirePermission(app, 'expense.create') }, async (req, reply) => {
    const body = parse(createBody, req.body);
    if (body.post && !req.auth!.permissions.has('expense.post')) throw forbidden('Missing permission: expense.post');
    const ctx = ctxOf(req);
    const expense = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, ctx.tenantId, body.companyId);
      const id = await createExpense(db, ctx, { ...body, companyId });
      return body.post ? postExpense(db, ctx, id) : loadExpense(db, ctx.tenantId, id);
    });
    reply.code(201);
    return expense;
  });

  app.patch('/expenses/:id', { preHandler: requirePermission(app, 'expense.create') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(updateBody, req.body);
    return req.tenantTx((db) => updateExpense(db, ctxOf(req), id, body));
  });

  app.delete('/expenses/:id', { preHandler: requirePermission(app, 'expense.create') }, async (req, reply) => {
    const { id } = parse(uuidParam, req.params);
    await req.tenantTx((db) => deleteExpense(db, ctxOf(req), id));
    return reply.code(204).send();
  });

  app.post('/expenses/:id/post', { preHandler: requirePermission(app, 'expense.post') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx((db) => postExpense(db, ctxOf(req), id));
  });

  app.post('/expenses/:id/cancel', { preHandler: requirePermission(app, 'expense.cancel') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const { reason } = parse(z.object({ reason: z.string().trim().min(3).max(500) }), req.body);
    return req.tenantTx((db) => cancelExpense(db, ctxOf(req), id, reason));
  });
}
