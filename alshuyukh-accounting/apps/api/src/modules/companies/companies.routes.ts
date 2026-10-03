import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import type { Db } from '../../db/tx.js';
import { conflict, notFound } from '../../lib/errors.js';
import { optionalText, parse, timezone, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { setupCompanyAccounting } from '../accounting/setup.js';
import { writeAudit } from '../audit/audit.service.js';

const companyFields = {
  name: z.string().trim().min(2).max(200),
  legalName: optionalText(200),
  commercialRegistration: z.string().regex(/^\d{10}$/, 'Commercial registration must be 10 digits').nullish(),
  vatNumber: z.string().regex(/^3\d{13}3$/, 'Saudi VAT number must be 15 digits starting and ending with 3').nullish(),
  address: optionalText(500),
  city: optionalText(100),
  country: z.string().regex(/^[A-Z]{2}$/).optional(),
  currency: z.string().regex(/^[A-Z]{3}$/).optional(),
  timezone: timezone.optional(),
  // TODO(Storage): logo upload endpoint. For now only an https URL is accepted.
  logoUrl: z.url({ protocol: /^https$/ }).max(1000).nullish(),
  email: z.string().trim().toLowerCase().pipe(z.email().max(254)).nullish(),
  phone: z.string().regex(/^\+?[0-9 ()-]{6,20}$/).nullish(),
};
const companyCreate = z.object(companyFields);
const companyUpdate = z.object({ ...companyFields, isActive: z.boolean() }).partial();

const COMPANY_COLUMNS: Record<string, string> = {
  name: 'name', legalName: 'legal_name', commercialRegistration: 'commercial_registration',
  vatNumber: 'vat_number', address: 'address', city: 'city', country: 'country', currency: 'currency',
  timezone: 'timezone', logoUrl: 'logo_url', email: 'email', phone: 'phone', isActive: 'is_active',
};

const COMPANY_SELECT = `SELECT id, name, legal_name AS "legalName", commercial_registration AS "commercialRegistration",
  vat_number AS "vatNumber", address, city, country, currency, timezone, logo_url AS "logoUrl", email, phone,
  is_active AS "isActive", created_at AS "createdAt", updated_at AS "updatedAt" FROM companies`;

const branchFields = {
  code: z.string().regex(/^[A-Za-z0-9_-]{1,20}$/),
  name: z.string().trim().min(2).max(200),
  address: optionalText(500),
  city: optionalText(100),
  phone: z.string().regex(/^\+?[0-9 ()-]{6,20}$/).nullish(),
};
const BRANCH_COLUMNS: Record<string, string> = { code: 'code', name: 'name', address: 'address', city: 'city', phone: 'phone', isActive: 'is_active' };
const BRANCH_SELECT = `SELECT id, company_id AS "companyId", code, name, address, city, phone, is_main AS "isMain",
  is_active AS "isActive", created_at AS "createdAt", updated_at AS "updatedAt" FROM branches`;

const warehouseFields = {
  code: z.string().regex(/^[A-Za-z0-9_-]{1,20}$/),
  name: z.string().trim().min(2).max(200),
  branchId: z.uuid().nullish(),
  address: optionalText(500),
};
const WAREHOUSE_COLUMNS: Record<string, string> = { code: 'code', name: 'name', branchId: 'branch_id', address: 'address', isActive: 'is_active' };
const WAREHOUSE_SELECT = `SELECT id, company_id AS "companyId", branch_id AS "branchId", code, name, address,
  is_active AS "isActive", created_at AS "createdAt", updated_at AS "updatedAt" FROM warehouses`;

/** Builds "col = $n" pairs only from whitelisted columns, so no input reaches SQL text. */
function buildSet(body: Record<string, unknown>, columns: Record<string, string>, startAt: number) {
  const sets: string[] = [];
  const values: unknown[] = [];
  for (const [key, col] of Object.entries(columns)) {
    if (body[key] !== undefined) {
      values.push(body[key]);
      sets.push(`${col} = $${startAt + values.length - 1}`);
    }
  }
  return { sets, values };
}

async function findOne<T>(db: Db, select: string, tenantId: string, id: string): Promise<T | undefined> {
  const { rows } = await db.query(`${select} WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, [tenantId, id]);
  return rows[0] as T | undefined;
}

async function assertCompany(db: Db, tenantId: string, companyId: string) {
  if (!(await findOne(db, COMPANY_SELECT, tenantId, companyId))) throw notFound('Company');
}

export default async function companiesRoutes(app: FastifyInstance) {
  const canView = requirePermission(app, 'company.view');
  const canManage = requirePermission(app, 'company.manage');

  // Companies ---------------------------------------------------------------
  app.get('/companies', { preHandler: canView }, async (req) =>
    req.tenantTx(async (db) => {
      const { rows } = await db.query(`${COMPANY_SELECT} WHERE tenant_id = $1 AND deleted_at IS NULL ORDER BY created_at`, [req.auth!.tenantId]);
      return { data: rows };
    }));

  app.get('/companies/:id', { preHandler: canView }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const company = await req.tenantTx((db) => findOne(db, COMPANY_SELECT, req.auth!.tenantId, id));
    if (!company) throw notFound('Company');
    return company;
  });

  app.post('/companies', { preHandler: canManage }, async (req, reply) => {
    const body = parse(companyCreate, req.body);
    const a = req.auth!;
    // TODO(Phase 9): enforce the plan's company limit.
    const company = await req.tenantTx(async (db) => {
      const { sets, values } = buildSet(body, COMPANY_COLUMNS, 4);
      const cols = sets.map((s) => s.split(' = ')[0]);
      const { rows: [c] } = await db.query<{ id: string }>(
        `INSERT INTO companies (tenant_id, created_by, updated_by${cols.map((c) => `, ${c}`).join('')})
         VALUES ($1, $2, $3${values.map((_, i) => `, $${i + 4}`).join('')}) RETURNING id`,
        [a.tenantId, a.userId, a.userId, ...values]);
      const created = await findOne<{ timezone: string }>(db, COMPANY_SELECT, a.tenantId, c!.id);
      await setupCompanyAccounting(db, { tenantId: a.tenantId, companyId: c!.id, userId: a.userId }, created!.timezone);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'company', entityId: c!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return company;
  });

  app.patch('/companies/:id', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(companyUpdate, req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const before = await findOne(db, COMPANY_SELECT, a.tenantId, id);
      if (!before) throw notFound('Company');
      const { sets, values } = buildSet(body, COMPANY_COLUMNS, 4);
      if (sets.length) {
        await db.query(`UPDATE companies SET ${sets.join(', ')}, updated_by = $3 WHERE tenant_id = $1 AND id = $2`,
          [a.tenantId, id, a.userId, ...values]);
      }
      const after = await findOne(db, COMPANY_SELECT, a.tenantId, id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'company', entityId: id, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });

  app.delete('/companies/:id', { preHandler: canManage }, async (req, reply) => {
    const { id } = parse(uuidParam, req.params);
    const a = req.auth!;
    await req.tenantTx(async (db) => {
      const before = await findOne(db, COMPANY_SELECT, a.tenantId, id);
      if (!before) throw notFound('Company');
      const { rows: [count] } = await db.query<{ n: string }>(
        `SELECT count(*) AS n FROM companies WHERE tenant_id = $1 AND deleted_at IS NULL`, [a.tenantId]);
      if (Number(count!.n) <= 1) throw conflict('LAST_COMPANY', 'An organization must keep at least one company');
      const posted = await db.query(`SELECT 1 FROM journal_entries WHERE company_id = $1 AND status <> 'DRAFT' LIMIT 1`, [id]);
      if (posted.rowCount) throw conflict('COMPANY_HAS_LEDGER', 'A company with posted journal entries cannot be deleted. Deactivate it instead.');
      await db.query(`UPDATE companies SET deleted_at = now(), is_active = false, updated_by = $3 WHERE tenant_id = $1 AND id = $2`, [a.tenantId, id, a.userId]);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'DELETE', entityType: 'company', entityId: id, oldValues: before }, req.auditMeta());
    });
    return reply.code(204).send();
  });

  // Branches ----------------------------------------------------------------
  app.get('/companies/:id/branches', { preHandler: canView }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx(async (db) => {
      await assertCompany(db, req.auth!.tenantId, id);
      const { rows } = await db.query(`${BRANCH_SELECT} WHERE tenant_id = $1 AND company_id = $2 AND deleted_at IS NULL ORDER BY is_main DESC, code`, [req.auth!.tenantId, id]);
      return { data: rows };
    });
  });

  app.post('/companies/:id/branches', { preHandler: canManage }, async (req, reply) => {
    const { id: companyId } = parse(uuidParam, req.params);
    const body = parse(z.object(branchFields), req.body);
    const a = req.auth!;
    // TODO(Phase 9): enforce the plan's branch limit.
    const branch = await req.tenantTx(async (db) => {
      await assertCompany(db, a.tenantId, companyId);
      const { rows: [b] } = await db.query<{ id: string }>(
        `INSERT INTO branches (tenant_id, company_id, code, name, address, city, phone, created_by, updated_by)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8) RETURNING id`,
        [a.tenantId, companyId, body.code, body.name, body.address ?? null, body.city ?? null, body.phone ?? null, a.userId]);
      const created = await findOne(db, BRANCH_SELECT, a.tenantId, b!.id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'branch', entityId: b!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return branch;
  });

  app.patch('/branches/:id', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ ...branchFields, isActive: z.boolean() }).partial(), req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const before = await findOne(db, BRANCH_SELECT, a.tenantId, id);
      if (!before) throw notFound('Branch');
      const { sets, values } = buildSet(body, BRANCH_COLUMNS, 4);
      if (sets.length) {
        await db.query(`UPDATE branches SET ${sets.join(', ')}, updated_by = $3 WHERE tenant_id = $1 AND id = $2`, [a.tenantId, id, a.userId, ...values]);
      }
      const after = await findOne(db, BRANCH_SELECT, a.tenantId, id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'branch', entityId: id, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });

  // Warehouses --------------------------------------------------------------
  // Stock itself is handled by the Inventory Engine (Phase 5).
  app.get('/companies/:id/warehouses', { preHandler: canView }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx(async (db) => {
      await assertCompany(db, req.auth!.tenantId, id);
      const { rows } = await db.query(`${WAREHOUSE_SELECT} WHERE tenant_id = $1 AND company_id = $2 AND deleted_at IS NULL ORDER BY code`, [req.auth!.tenantId, id]);
      return { data: rows };
    });
  });

  app.post('/companies/:id/warehouses', { preHandler: canManage }, async (req, reply) => {
    const { id: companyId } = parse(uuidParam, req.params);
    const body = parse(z.object(warehouseFields), req.body);
    const a = req.auth!;
    // TODO(Phase 9): enforce the plan's warehouse limit.
    const warehouse = await req.tenantTx(async (db) => {
      await assertCompany(db, a.tenantId, companyId);
      if (body.branchId) {
        const b = await db.query(`SELECT 1 FROM branches WHERE tenant_id = $1 AND id = $2 AND company_id = $3 AND deleted_at IS NULL`, [a.tenantId, body.branchId, companyId]);
        if (!b.rowCount) throw notFound('Branch');
      }
      const { rows: [w] } = await db.query<{ id: string }>(
        `INSERT INTO warehouses (tenant_id, company_id, branch_id, code, name, address, created_by, updated_by)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $7) RETURNING id`,
        [a.tenantId, companyId, body.branchId ?? null, body.code, body.name, body.address ?? null, a.userId]);
      const created = await findOne(db, WAREHOUSE_SELECT, a.tenantId, w!.id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'warehouse', entityId: w!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return warehouse;
  });

  app.patch('/warehouses/:id', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ ...warehouseFields, isActive: z.boolean() }).partial(), req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const before = await findOne<{ companyId: string }>(db, WAREHOUSE_SELECT, a.tenantId, id);
      if (!before) throw notFound('Warehouse');
      if (body.branchId) {
        const b = await db.query(`SELECT 1 FROM branches WHERE tenant_id = $1 AND id = $2 AND company_id = $3 AND deleted_at IS NULL`, [a.tenantId, body.branchId, before.companyId]);
        if (!b.rowCount) throw notFound('Branch');
      }
      const { sets, values } = buildSet(body, WAREHOUSE_COLUMNS, 4);
      if (sets.length) {
        await db.query(`UPDATE warehouses SET ${sets.join(', ')}, updated_by = $3 WHERE tenant_id = $1 AND id = $2`, [a.tenantId, id, a.userId, ...values]);
      }
      const after = await findOne(db, WAREHOUSE_SELECT, a.tenantId, id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'warehouse', entityId: id, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });
}
