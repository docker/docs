import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import type { Db } from '../../db/tx.js';
import { badRequest, conflict, notFound } from '../../lib/errors.js';
import { Decimal } from '../../lib/money.js';
import { optionalText, parse, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { writeAudit } from '../audit/audit.service.js';
import { ACCOUNT_TYPES } from './chart-template.js';
import { resolveCompanyId } from './company-context.js';
import { setupCompanyAccounting } from './setup.js';

const ACCOUNT_SELECT = `
  SELECT a.id, a.company_id AS "companyId", a.code, a.name_ar AS "nameAr", a.name_en AS "nameEn",
         a.account_type AS "type", a.parent_id AS "parentId", a.level, a.group_id AS "groupId",
         a.is_postable AS "isPostable", a.is_active AS "isActive", a.is_system AS "isSystem",
         a.system_key AS "systemKey", a.description,
         COALESCE(b.debit, 0)::numeric(18,2)::text AS "totalDebit",
         COALESCE(b.credit, 0)::numeric(18,2)::text AS "totalCredit",
         COALESCE(b.debit - b.credit, 0)::numeric(18,2)::text AS "balance"
    FROM accounts a
    LEFT JOIN LATERAL (
      SELECT sum(l.debit) AS debit, sum(l.credit) AS credit
        FROM journal_entry_lines l JOIN journal_entries e ON e.id = l.journal_entry_id
       WHERE l.account_id = a.id AND e.status IN ('POSTED', 'REVERSED')
    ) b ON true`;

type AccountRow = { id: string; companyId: string; isSystem: boolean; isActive: boolean; isPostable: boolean; balance: string; type: string };

async function findAccount(db: Db, tenantId: string, id: string): Promise<AccountRow | undefined> {
  const { rows } = await db.query<AccountRow>(`${ACCOUNT_SELECT} WHERE a.tenant_id = $1 AND a.id = $2 AND a.deleted_at IS NULL`, [tenantId, id]);
  return rows[0];
}

const accountCode = z.string().trim().regex(/^[0-9A-Za-z.-]{1,20}$/, 'Code: up to 20 letters, digits, dots or dashes');

const createBody = z.object({
  companyId: z.uuid().optional(),
  code: accountCode,
  nameAr: z.string().trim().min(1).max(200),
  nameEn: optionalText(200),
  type: z.enum(ACCOUNT_TYPES as [string, ...string[]]).optional(),
  parentId: z.uuid().nullish(),
  groupId: z.uuid().nullish(),
  isPostable: z.boolean().default(true),
  description: optionalText(1000),
});

const updateBody = z.object({
  code: accountCode,
  nameAr: z.string().trim().min(1).max(200),
  nameEn: optionalText(200),
  parentId: z.uuid().nullable(),
  groupId: z.uuid().nullable(),
  isPostable: z.boolean(),
  isActive: z.boolean(),
  description: optionalText(1000),
}).partial();

export default async function accountsRoutes(app: FastifyInstance) {
  const canView = requirePermission(app, 'account.view');
  const canManage = requirePermission(app, 'account.manage');

  /** Seeds the default chart and current fiscal year for a company that has none (idempotent). */
  app.post('/accounting/setup', { preHandler: requirePermission(app, 'account.manage', 'fiscal.manage') }, async (req) => {
    const body = parse(z.object({ companyId: z.uuid().optional() }), req.body ?? {});
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      const { rows: [c] } = await db.query<{ timezone: string }>(`SELECT timezone FROM companies WHERE id = $1`, [companyId]);
      const result = await setupCompanyAccounting(db, { tenantId: a.tenantId, companyId, userId: a.userId }, c!.timezone);
      if (result.chartCreated || result.fiscalYearCreated) {
        await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'accounting_setup', entityId: companyId, newValues: result }, req.auditMeta());
      }
      return { companyId, ...result };
    });
  });

  app.get('/account-groups', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query(
        `SELECT id, code, name_ar AS "nameAr", name_en AS "nameEn", account_type AS "type", sort_order AS "sortOrder"
           FROM account_groups WHERE tenant_id = $1 AND company_id = $2 ORDER BY sort_order`, [req.auth!.tenantId, companyId]);
      return { data: rows };
    });
  });

  app.get('/accounts', { preHandler: canView }, async (req) => {
    const q = parse(z.object({
      companyId: z.uuid().optional(),
      includeInactive: z.enum(['true', 'false']).optional(),
      postableOnly: z.enum(['true', 'false']).optional(),
    }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query(
        `${ACCOUNT_SELECT}
          WHERE a.tenant_id = $1 AND a.company_id = $2 AND a.deleted_at IS NULL
            AND ($3::boolean OR a.is_active) AND (NOT $4::boolean OR a.is_postable)
          ORDER BY a.code`,
        [req.auth!.tenantId, companyId, q.includeInactive === 'true', q.postableOnly === 'true']);
      return { data: rows };
    });
  });

  app.get('/accounts/:id', { preHandler: canView }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const account = await req.tenantTx((db) => findAccount(db, req.auth!.tenantId, id));
    if (!account) throw notFound('Account');
    return account;
  });

  app.post('/accounts', { preHandler: canManage }, async (req, reply) => {
    const body = parse(createBody, req.body);
    const a = req.auth!;
    const account = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      let type = body.type;
      if (body.parentId) {
        const { rows: [parent] } = await db.query<{ account_type: string }>(
          `SELECT account_type FROM accounts WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, body.parentId]);
        if (!parent) throw badRequest('INVALID_PARENT', 'Parent account does not exist in this company');
        if (type && type !== parent.account_type) throw badRequest('TYPE_MISMATCH', 'Account type must match its parent');
        type = parent.account_type;
      }
      if (!type) throw badRequest('TYPE_REQUIRED', 'type is required for a top-level account');
      if (body.parentId) {
        // A postable parent with no postings can become a header automatically.
        const { rows: [p] } = await db.query<{ is_postable: boolean; has_lines: boolean }>(
          `SELECT is_postable, EXISTS (SELECT 1 FROM journal_entry_lines WHERE account_id = $1) AS has_lines FROM accounts WHERE id = $1`, [body.parentId]);
        if (p!.is_postable) {
          if (p!.has_lines) throw conflict('PARENT_HAS_POSTINGS', 'The parent account has journal lines and cannot become a header');
          await db.query(`UPDATE accounts SET is_postable = false, updated_by = $2 WHERE id = $1`, [body.parentId, a.userId]);
        }
      }
      const { rows: [row] } = await db.query<{ id: string }>(
        `INSERT INTO accounts (tenant_id, company_id, code, name_ar, name_en, account_type, parent_id, group_id,
                               is_postable, description, created_by, updated_by)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11) RETURNING id`,
        [a.tenantId, companyId, body.code, body.nameAr, body.nameEn ?? null, type, body.parentId ?? null,
         body.groupId ?? null, body.isPostable, body.description ?? null, a.userId]);
      const created = await findAccount(db, a.tenantId, row!.id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'account', entityId: row!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return account;
  });

  app.patch('/accounts/:id', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(updateBody, req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const before = await findAccount(db, a.tenantId, id);
      if (!before) throw notFound('Account');
      if (body.isActive === false && before.isActive) {
        if (before.isSystem) throw conflict('SYSTEM_ACCOUNT', 'System accounts are used by the accounting engine and cannot be deactivated');
        if (!new Decimal(before.balance).isZero()) throw conflict('ACCOUNT_HAS_BALANCE', 'Only accounts with a zero balance can be deactivated');
      }
      await db.query(
        `UPDATE accounts SET
           code = COALESCE($2, code), name_ar = COALESCE($3, name_ar),
           name_en = CASE WHEN $4::boolean THEN $5 ELSE name_en END,
           parent_id = CASE WHEN $6::boolean THEN $7::uuid ELSE parent_id END,
           group_id = CASE WHEN $8::boolean THEN $9::uuid ELSE group_id END,
           is_postable = COALESCE($10, is_postable), is_active = COALESCE($11, is_active),
           description = CASE WHEN $12::boolean THEN $13 ELSE description END,
           updated_by = $14
         WHERE tenant_id = $1 AND id = $15`,
        [a.tenantId, body.code ?? null, body.nameAr ?? null, body.nameEn !== undefined, body.nameEn ?? null,
         body.parentId !== undefined, body.parentId ?? null, body.groupId !== undefined, body.groupId ?? null,
         body.isPostable ?? null, body.isActive ?? null, body.description !== undefined, body.description ?? null,
         a.userId, id]);
      const after = await findAccount(db, a.tenantId, id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'account', entityId: id, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });

  app.delete('/accounts/:id', { preHandler: canManage }, async (req, reply) => {
    const { id } = parse(uuidParam, req.params);
    const a = req.auth!;
    await req.tenantTx(async (db) => {
      const before = await findAccount(db, a.tenantId, id);
      if (!before) throw notFound('Account');
      if (before.isSystem) throw conflict('SYSTEM_ACCOUNT', 'System accounts cannot be deleted');
      const used = await db.query(`SELECT 1 FROM journal_entry_lines WHERE account_id = $1 LIMIT 1`, [id]);
      if (used.rowCount) throw conflict('ACCOUNT_IN_USE', 'Accounts with journal lines cannot be deleted. Deactivate it instead.');
      const children = await db.query(`SELECT 1 FROM accounts WHERE parent_id = $1 AND deleted_at IS NULL LIMIT 1`, [id]);
      if (children.rowCount) throw conflict('ACCOUNT_HAS_CHILDREN', 'Delete or move the sub-accounts first');
      await db.query(`UPDATE accounts SET deleted_at = now(), is_active = false, updated_by = $2 WHERE id = $1`, [id, a.userId]);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'DELETE', entityType: 'account', entityId: id, oldValues: before }, req.auditMeta());
    });
    return reply.code(204).send();
  });

  // Cost centers ----------------------------------------------------------------
  const CC_SELECT = `SELECT id, company_id AS "companyId", code, name, is_active AS "isActive" FROM cost_centers`;

  app.get('/cost-centers', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query(`${CC_SELECT} WHERE tenant_id = $1 AND company_id = $2 AND deleted_at IS NULL ORDER BY code`, [req.auth!.tenantId, companyId]);
      return { data: rows };
    });
  });

  app.post('/cost-centers', { preHandler: canManage }, async (req, reply) => {
    const body = parse(z.object({ companyId: z.uuid().optional(), code: z.string().regex(/^[0-9A-Za-z_-]{1,20}$/), name: z.string().trim().min(2).max(200) }), req.body);
    const a = req.auth!;
    const cc = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      const { rows: [row] } = await db.query<{ id: string }>(
        `INSERT INTO cost_centers (tenant_id, company_id, code, name, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $5) RETURNING id`,
        [a.tenantId, companyId, body.code, body.name, a.userId]);
      const created = (await db.query(`${CC_SELECT} WHERE id = $1`, [row!.id])).rows[0];
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'cost_center', entityId: row!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return cc;
  });

  app.patch('/cost-centers/:id', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ name: z.string().trim().min(2).max(200), isActive: z.boolean() }).partial(), req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const before = (await db.query(`${CC_SELECT} WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, [a.tenantId, id])).rows[0];
      if (!before) throw notFound('Cost center');
      await db.query(`UPDATE cost_centers SET name = COALESCE($2, name), is_active = COALESCE($3, is_active), updated_by = $4 WHERE id = $1`,
        [id, body.name ?? null, body.isActive ?? null, a.userId]);
      const after = (await db.query(`${CC_SELECT} WHERE id = $1`, [id])).rows[0];
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'cost_center', entityId: id, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });
}
