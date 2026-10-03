import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import type { Db } from '../../db/tx.js';
import { badRequest, conflict, notFound } from '../../lib/errors.js';
import { nextFreeCode } from '../../lib/sequences.js';
import { optionalText, parse, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from '../accounting/company-context.js';
import { writeAudit } from '../audit/audit.service.js';

/** Unit price: up to 4 decimals, as a string. */
const price = z.string().trim().regex(/^\d{1,14}(\.\d{1,4})?$/, 'Price must be a non-negative number with at most 4 decimals, sent as a string');

const PRODUCT_SELECT = `
  SELECT p.id, p.company_id AS "companyId", p.sku, p.barcode, p.name_ar AS "nameAr", p.name_en AS "nameEn",
         p.description, p.product_type AS "productType", p.category_id AS "categoryId", c.name_ar AS "categoryName",
         p.unit_id AS "unitId", u.code AS "unitCode", u.name_ar AS "unitName",
         p.sale_price::text AS "salePrice", p.sale_price_includes_vat AS "salePriceIncludesVat",
         p.purchase_price::text AS "purchasePrice", p.vat_category AS "vatCategory", p.track_inventory AS "trackInventory",
         p.sales_account_id AS "salesAccountId", p.purchase_account_id AS "purchaseAccountId",
         p.is_active AS "isActive", p.created_at AS "createdAt", p.updated_at AS "updatedAt"
    FROM products p
    JOIN units u ON u.id = p.unit_id
    LEFT JOIN product_categories c ON c.id = p.category_id`;

const productFields = {
  sku: z.string().trim().regex(/^[0-9A-Za-z._/-]{1,40}$/, 'SKU: letters, digits, . _ / -'),
  barcode: z.string().trim().regex(/^[0-9A-Za-z-]{4,40}$/).nullish(),
  nameAr: z.string().trim().min(2).max(200),
  nameEn: optionalText(200),
  description: optionalText(2000),
  productType: z.enum(['GOODS', 'SERVICE']),
  categoryId: z.uuid().nullish(),
  unitId: z.uuid(),
  salePrice: price,
  salePriceIncludesVat: z.boolean(),
  purchasePrice: price,
  vatCategory: z.enum(['S', 'Z', 'E', 'O']),
  trackInventory: z.boolean(),
  salesAccountId: z.uuid().nullish(),
  purchaseAccountId: z.uuid().nullish(),
};

const createBody = z.object({
  companyId: z.uuid().optional(),
  ...productFields,
  sku: productFields.sku.optional(),
  unitId: productFields.unitId.optional(),
  productType: productFields.productType.default('GOODS'),
  salePrice: price.default('0'),
  salePriceIncludesVat: z.boolean().default(false),
  purchasePrice: price.default('0'),
  vatCategory: productFields.vatCategory.default('S'),
  trackInventory: z.boolean().optional(),
});
const updateBody = z.object({ ...productFields, isActive: z.boolean() }).partial();

const COLUMNS: [keyof z.infer<typeof updateBody>, string][] = [
  ['sku', 'sku'], ['barcode', 'barcode'], ['nameAr', 'name_ar'], ['nameEn', 'name_en'], ['description', 'description'],
  ['productType', 'product_type'], ['categoryId', 'category_id'], ['unitId', 'unit_id'], ['salePrice', 'sale_price'],
  ['salePriceIncludesVat', 'sale_price_includes_vat'], ['purchasePrice', 'purchase_price'], ['vatCategory', 'vat_category'],
  ['trackInventory', 'track_inventory'], ['salesAccountId', 'sales_account_id'], ['purchaseAccountId', 'purchase_account_id'],
  ['isActive', 'is_active'],
];

async function loadProduct(db: Db, tenantId: string, id: string) {
  const { rows: [p] } = await db.query<{ companyId: string; productType: string; trackInventory: boolean; [k: string]: unknown }>(
    `${PRODUCT_SELECT} WHERE p.tenant_id = $1 AND p.id = $2 AND p.deleted_at IS NULL`, [tenantId, id]);
  return p;
}

/** Checks that referenced rows belong to the same company and fit their role. */
async function checkReferences(db: Db, companyId: string, body: Partial<z.infer<typeof updateBody>>) {
  if (body.unitId) {
    const u = await db.query(`SELECT 1 FROM units WHERE company_id = $1 AND id = $2 AND is_active`, [companyId, body.unitId]);
    if (!u.rowCount) throw badRequest('INVALID_UNIT', 'Unit does not exist in this company');
  }
  if (body.categoryId) {
    const c = await db.query(`SELECT 1 FROM product_categories WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, body.categoryId]);
    if (!c.rowCount) throw badRequest('INVALID_CATEGORY', 'Category does not exist in this company');
  }
  const checkAccount = async (id: string | null | undefined, allowed: string[], label: string) => {
    if (!id) return;
    const { rows: [a] } = await db.query<{ account_type: string; is_postable: boolean; is_active: boolean }>(
      `SELECT account_type, is_postable, is_active FROM accounts WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, id]);
    if (!a || !a.is_postable || !a.is_active || !allowed.includes(a.account_type)) {
      throw badRequest('INVALID_PRODUCT_ACCOUNT', `${label} must be an active postable account of type ${allowed.join(' or ')}`);
    }
  };
  await checkAccount(body.salesAccountId, ['REVENUE'], 'Sales account');
  await checkAccount(body.purchaseAccountId, ['ASSET', 'EXPENSE', 'COST_OF_GOODS_SOLD'], 'Purchase account');
}

export default async function productsRoutes(app: FastifyInstance) {
  const canView = requirePermission(app, 'product.view');
  const canManage = requirePermission(app, 'product.manage');

  // Units ------------------------------------------------------------------------
  app.get('/units', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query(
        `SELECT id, code, name_ar AS "nameAr", name_en AS "nameEn", is_active AS "isActive"
           FROM units WHERE tenant_id = $1 AND company_id = $2 ORDER BY code`, [req.auth!.tenantId, companyId]);
      return { data: rows };
    });
  });

  app.post('/units', { preHandler: canManage }, async (req, reply) => {
    const body = parse(z.object({
      companyId: z.uuid().optional(),
      code: z.string().trim().toUpperCase().regex(/^[A-Z0-9]{1,10}$/),
      nameAr: z.string().trim().min(1).max(50),
      nameEn: optionalText(50),
    }), req.body);
    const a = req.auth!;
    const unit = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      const { rows: [u] } = await db.query<{ id: string }>(
        `INSERT INTO units (tenant_id, company_id, code, name_ar, name_en) VALUES ($1, $2, $3, $4, $5)
         RETURNING id, code, name_ar AS "nameAr", name_en AS "nameEn", is_active AS "isActive"`,
        [a.tenantId, companyId, body.code, body.nameAr, body.nameEn ?? null]);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'unit', entityId: u!.id, newValues: u }, req.auditMeta());
      return u;
    });
    reply.code(201);
    return unit;
  });

  app.patch('/units/:id', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ nameAr: z.string().trim().min(1).max(50), nameEn: optionalText(50), isActive: z.boolean() }).partial(), req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const { rows: [before] } = await db.query(`SELECT id, code, name_ar AS "nameAr", name_en AS "nameEn", is_active AS "isActive" FROM units WHERE tenant_id = $1 AND id = $2`, [a.tenantId, id]);
      if (!before) throw notFound('Unit');
      if (body.isActive === false) {
        const used = await db.query(`SELECT 1 FROM products WHERE unit_id = $1 AND deleted_at IS NULL AND is_active LIMIT 1`, [id]);
        if (used.rowCount) throw conflict('UNIT_IN_USE', 'Active products use this unit');
      }
      const { rows: [after] } = await db.query(
        `UPDATE units SET name_ar = COALESCE($2, name_ar), name_en = CASE WHEN $3::boolean THEN $4 ELSE name_en END,
                          is_active = COALESCE($5, is_active)
          WHERE id = $1 RETURNING id, code, name_ar AS "nameAr", name_en AS "nameEn", is_active AS "isActive"`,
        [id, body.nameAr ?? null, body.nameEn !== undefined, body.nameEn ?? null, body.isActive ?? null]);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'unit', entityId: id, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });

  // Categories -------------------------------------------------------------------
  const CATEGORY_SELECT = `SELECT id, parent_id AS "parentId", name_ar AS "nameAr", name_en AS "nameEn", is_active AS "isActive" FROM product_categories`;

  app.get('/product-categories', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query(`${CATEGORY_SELECT} WHERE tenant_id = $1 AND company_id = $2 AND deleted_at IS NULL ORDER BY name_ar`, [req.auth!.tenantId, companyId]);
      return { data: rows };
    });
  });

  app.post('/product-categories', { preHandler: canManage }, async (req, reply) => {
    const body = parse(z.object({ companyId: z.uuid().optional(), nameAr: z.string().trim().min(2).max(100), nameEn: optionalText(100), parentId: z.uuid().nullish() }), req.body);
    const a = req.auth!;
    const cat = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      if (body.parentId) {
        const p = await db.query(`SELECT 1 FROM product_categories WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, body.parentId]);
        if (!p.rowCount) throw badRequest('INVALID_CATEGORY', 'Parent category does not exist in this company');
      }
      const { rows: [c] } = await db.query<{ id: string }>(
        `INSERT INTO product_categories (tenant_id, company_id, parent_id, name_ar, name_en) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
        [a.tenantId, companyId, body.parentId ?? null, body.nameAr, body.nameEn ?? null]);
      const created = (await db.query(`${CATEGORY_SELECT} WHERE id = $1`, [c!.id])).rows[0];
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'product_category', entityId: c!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return cat;
  });

  app.patch('/product-categories/:id', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ nameAr: z.string().trim().min(2).max(100), nameEn: optionalText(100), isActive: z.boolean() }).partial(), req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const before = (await db.query(`${CATEGORY_SELECT} WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, [a.tenantId, id])).rows[0];
      if (!before) throw notFound('Category');
      await db.query(
        `UPDATE product_categories SET name_ar = COALESCE($2, name_ar), name_en = CASE WHEN $3::boolean THEN $4 ELSE name_en END,
                is_active = COALESCE($5, is_active) WHERE id = $1`,
        [id, body.nameAr ?? null, body.nameEn !== undefined, body.nameEn ?? null, body.isActive ?? null]);
      const after = (await db.query(`${CATEGORY_SELECT} WHERE id = $1`, [id])).rows[0];
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'product_category', entityId: id, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });

  // Products ---------------------------------------------------------------------
  app.get('/products', { preHandler: canView }, async (req) => {
    const q = parse(z.object({
      companyId: z.uuid().optional(),
      search: z.string().trim().max(100).optional(),
      categoryId: z.uuid().optional(),
      productType: z.enum(['GOODS', 'SERVICE']).optional(),
      status: z.enum(['active', 'inactive', 'all']).default('active'),
      limit: z.coerce.number().int().min(1).max(200).default(50),
      offset: z.coerce.number().int().min(0).default(0),
    }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const search = q.search ? `%${q.search.replace(/[%_\\]/g, (m) => `\\${m}`)}%` : null;
      const { rows } = await db.query<{ total: string }>(
        `${PRODUCT_SELECT.replace('SELECT p.id', 'SELECT count(*) OVER () AS total, p.id')}
          WHERE p.tenant_id = $1 AND p.company_id = $2 AND p.deleted_at IS NULL
            AND ($3 = 'all' OR p.is_active = ($3 = 'active'))
            AND ($4::text IS NULL OR p.name_ar ILIKE $4 OR p.name_en ILIKE $4 OR p.sku ILIKE $4 OR p.barcode = $5)
            AND ($6::uuid IS NULL OR p.category_id = $6)
            AND ($7::text IS NULL OR p.product_type = $7)
          ORDER BY p.name_ar LIMIT $8 OFFSET $9`,
        [req.auth!.tenantId, companyId, q.status, search, q.search ?? null, q.categoryId ?? null, q.productType ?? null, q.limit, q.offset]);
      return { data: rows.map(({ total: _t, ...r }) => r), total: Number(rows[0]?.total ?? 0) };
    });
  });

  app.get('/products/:id', { preHandler: canView }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const product = await req.tenantTx((db) => loadProduct(db, req.auth!.tenantId, id));
    if (!product) throw notFound('Product');
    return product;
  });

  app.post('/products', { preHandler: canManage }, async (req, reply) => {
    const body = parse(createBody, req.body);
    const trackInventory = body.trackInventory ?? body.productType === 'GOODS';
    if (body.productType === 'SERVICE' && trackInventory) throw badRequest('SERVICE_NO_STOCK', 'Services cannot track inventory');
    const a = req.auth!;
    const product = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      let unitId = body.unitId;
      if (!unitId) {
        const { rows: [u] } = await db.query<{ id: string }>(`SELECT id FROM units WHERE company_id = $1 AND code = $2`, [companyId, body.productType === 'SERVICE' ? 'HUR' : 'PCE']);
        if (!u) throw badRequest('UNIT_REQUIRED', 'unitId is required');
        unitId = u.id;
      }
      await checkReferences(db, companyId, { ...body, unitId });
      const sku = body.sku ?? await nextFreeCode(db, a.tenantId, companyId, 'PRODUCT', { prefix: 'PRD' }, 'products', 'sku');
      const { rows: [row] } = await db.query<{ id: string }>(
        `INSERT INTO products (tenant_id, company_id, sku, barcode, name_ar, name_en, description, product_type, category_id,
           unit_id, sale_price, sale_price_includes_vat, purchase_price, vat_category, track_inventory,
           sales_account_id, purchase_account_id, created_by, updated_by)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $18) RETURNING id`,
        [a.tenantId, companyId, sku, body.barcode ?? null, body.nameAr, body.nameEn ?? null, body.description ?? null,
         body.productType, body.categoryId ?? null, unitId, body.salePrice, body.salePriceIncludesVat, body.purchasePrice,
         body.vatCategory, trackInventory, body.salesAccountId ?? null, body.purchaseAccountId ?? null, a.userId]);
      const created = await loadProduct(db, a.tenantId, row!.id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'product', entityId: row!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return product;
  });

  app.patch('/products/:id', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(updateBody, req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const before = await loadProduct(db, a.tenantId, id);
      if (!before) throw notFound('Product');
      const type = body.productType ?? before.productType;
      const track = body.trackInventory ?? (body.productType === 'SERVICE' ? false : before.trackInventory);
      if (type === 'SERVICE' && track) throw badRequest('SERVICE_NO_STOCK', 'Services cannot track inventory');
      // TODO(Phase 5): refuse changing type or tracking once stock movements exist.
      if (body.productType === 'SERVICE' && body.trackInventory === undefined) body.trackInventory = false;
      await checkReferences(db, before.companyId, body);
      const sets: string[] = [];
      const values: unknown[] = [];
      for (const [key, col] of COLUMNS) {
        if (body[key] !== undefined) { values.push(body[key]); sets.push(`${col} = $${values.length + 3}`); }
      }
      if (sets.length) {
        await db.query(`UPDATE products SET ${sets.join(', ')}, updated_by = $3 WHERE tenant_id = $1 AND id = $2`, [a.tenantId, id, a.userId, ...values]);
      }
      const after = await loadProduct(db, a.tenantId, id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'product', entityId: id, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });

  app.delete('/products/:id', { preHandler: canManage }, async (req, reply) => {
    const { id } = parse(uuidParam, req.params);
    const a = req.auth!;
    await req.tenantTx(async (db) => {
      const before = await loadProduct(db, a.tenantId, id);
      if (!before) throw notFound('Product');
      // TODO(Phase 4/5): refuse when invoices or stock movements reference the product.
      await db.query(`UPDATE products SET deleted_at = now(), is_active = false, updated_by = $2 WHERE id = $1`, [id, a.userId]);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'DELETE', entityType: 'product', entityId: id, oldValues: before }, req.auditMeta());
    });
    return reply.code(204).send();
  });
}
