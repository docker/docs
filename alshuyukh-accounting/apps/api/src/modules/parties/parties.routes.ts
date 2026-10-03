import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import type { Db } from '../../db/tx.js';
import { badRequest, conflict, notFound } from '../../lib/errors.js';
import { amountString } from '../../lib/money.js';
import { nextFreeCode } from '../../lib/sequences.js';
import { optionalText, parse, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from '../accounting/company-context.js';
import { writeAudit } from '../audit/audit.service.js';

/**
 * Customers and suppliers share one implementation. Each kind has its own
 * tables, permissions, numbering, and control account (AR / AP). Balances
 * come from journal lines tagged with the party, never from a stored field.
 */
interface PartyKind {
  path: string;                    // URL segment
  table: 'customers' | 'suppliers';
  addressTable: 'customer_addresses' | 'supplier_addresses';
  fk: 'customer_id' | 'supplier_id';
  accountColumn: 'receivable_account_id' | 'payable_account_id';
  // Control account must be of this type (AR is an asset, AP a liability).
  accountType: 'ASSET' | 'LIABILITY';
  defaultAccountKey: 'ACCOUNTS_RECEIVABLE' | 'ACCOUNTS_PAYABLE';
  docType: 'CUSTOMER' | 'SUPPLIER';
  prefix: string;
  viewPermission: string;
  managePermission: string;
  entity: 'customer' | 'supplier';
  label: string;
}

export const CUSTOMER: PartyKind = {
  path: 'customers', table: 'customers', addressTable: 'customer_addresses', fk: 'customer_id',
  accountColumn: 'receivable_account_id', accountType: 'ASSET',
  defaultAccountKey: 'ACCOUNTS_RECEIVABLE', docType: 'CUSTOMER', prefix: 'CUS',
  viewPermission: 'customer.view', managePermission: 'customer.manage', entity: 'customer', label: 'Customer',
};

export const SUPPLIER: PartyKind = {
  path: 'suppliers', table: 'suppliers', addressTable: 'supplier_addresses', fk: 'supplier_id',
  accountColumn: 'payable_account_id', accountType: 'LIABILITY',
  defaultAccountKey: 'ACCOUNTS_PAYABLE', docType: 'SUPPLIER', prefix: 'SUP',
  viewPermission: 'supplier.view', managePermission: 'supplier.manage', entity: 'supplier', label: 'Supplier',
};

const address = z.object({
  addressType: z.enum(['BILLING', 'SHIPPING']).default('BILLING'),
  isDefault: z.boolean().default(false),
  buildingNumber: z.string().regex(/^\d{4}$/, 'Building number must be 4 digits').nullish(),
  street: optionalText(200),
  district: optionalText(100),
  city: optionalText(100),
  postalCode: z.string().regex(/^\d{5}$/, 'Postal code must be 5 digits').nullish(),
  additionalNumber: z.string().regex(/^\d{4}$/, 'Additional number must be 4 digits').nullish(),
  country: z.string().regex(/^[A-Z]{2}$/).default('SA'),
});
type AddressInput = z.infer<typeof address>;

const fields = {
  code: z.string().trim().regex(/^[0-9A-Za-z_-]{1,30}$/, 'Code: letters, digits, _ or -').optional(),
  partyType: z.enum(['BUSINESS', 'INDIVIDUAL']),
  nameAr: z.string().trim().min(2).max(200),
  nameEn: optionalText(200),
  vatNumber: z.string().regex(/^3\d{13}3$/, 'Saudi VAT number must be 15 digits starting and ending with 3').nullish(),
  commercialRegistration: z.string().regex(/^\d{10}$/, 'Commercial registration must be 10 digits').nullish(),
  nationalId: z.string().regex(/^[12]\d{9}$/, 'National ID / Iqama must be 10 digits starting with 1 or 2').nullish(),
  email: z.string().trim().toLowerCase().pipe(z.email().max(254)).nullish(),
  phone: z.string().regex(/^\+?[0-9 ()-]{6,20}$/).nullish(),
  creditLimit: amountString.nullish(),
  paymentTermsDays: z.number().int().min(0).max(365),
  controlAccountId: z.uuid().nullish(),
  notes: optionalText(2000),
};

const createBody = z.object({
  companyId: z.uuid().optional(),
  ...fields,
  partyType: fields.partyType.default('BUSINESS'),
  paymentTermsDays: fields.paymentTermsDays.default(0),
  addresses: z.array(address).max(10).default([]),
});
const updateBody = z.object({ ...fields, isActive: z.boolean() }).partial();

function checkAddresses(addresses: AddressInput[]) {
  for (const type of ['BILLING', 'SHIPPING'] as const) {
    const ofType = addresses.filter((a) => a.addressType === type);
    const defaults = ofType.filter((a) => a.isDefault).length;
    if (defaults > 1) throw badRequest('MULTIPLE_DEFAULT_ADDRESSES', `Only one default ${type.toLowerCase()} address is allowed`);
    // The first address of each type becomes the default when none is marked.
    if (ofType.length && defaults === 0) ofType[0]!.isDefault = true;
  }
}

export function partyRoutes(kind: PartyKind) {
  const SELECT = `
    SELECT p.id, p.company_id AS "companyId", p.code, p.party_type AS "partyType", p.name_ar AS "nameAr", p.name_en AS "nameEn",
           p.vat_number AS "vatNumber", p.commercial_registration AS "commercialRegistration", p.national_id AS "nationalId",
           p.email, p.phone, p.credit_limit::text AS "creditLimit", p.payment_terms_days AS "paymentTermsDays",
           p.${kind.accountColumn} AS "controlAccountId", p.notes, p.is_active AS "isActive",
           p.created_at AS "createdAt", p.updated_at AS "updatedAt",
           COALESCE(b.balance, 0)::numeric(18,2)::text AS balance
      FROM ${kind.table} p
      LEFT JOIN LATERAL (
        SELECT sum(l.debit - l.credit) AS balance
          FROM journal_entry_lines l JOIN journal_entries e ON e.id = l.journal_entry_id
         WHERE l.${kind.fk} = p.id AND e.status IN ('POSTED', 'REVERSED')
      ) b ON true`;

  const ADDRESS_SELECT = `
    SELECT id, address_type AS "addressType", is_default AS "isDefault", building_number AS "buildingNumber",
           street, district, city, postal_code AS "postalCode", additional_number AS "additionalNumber", country
      FROM ${kind.addressTable}`;

  async function load(db: Db, tenantId: string, id: string) {
    const { rows: [party] } = await db.query<{ companyId: string; [key: string]: unknown }>(
      `${SELECT} WHERE p.tenant_id = $1 AND p.id = $2 AND p.deleted_at IS NULL`, [tenantId, id]);
    if (!party) return undefined;
    const { rows: addresses } = await db.query(
      `${ADDRESS_SELECT} WHERE ${kind.fk} = $1 AND deleted_at IS NULL ORDER BY address_type, is_default DESC, created_at`, [id]);
    return { ...party, addresses };
  }

  async function checkControlAccount(db: Db, companyId: string, accountId: string | null | undefined) {
    if (!accountId) return;
    const { rows: [a] } = await db.query<{ account_type: string; is_postable: boolean; is_active: boolean }>(
      `SELECT account_type, is_postable, is_active FROM accounts WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, accountId]);
    if (!a) throw badRequest('INVALID_ACCOUNT', 'Account does not exist in this company');
    if (a.account_type !== kind.accountType || !a.is_postable || !a.is_active) {
      throw badRequest('INVALID_CONTROL_ACCOUNT', `The ${kind.entity} account must be an active postable ${kind.accountType.toLowerCase()} account`);
    }
  }

  async function insertAddresses(db: Db, tenantId: string, companyId: string, partyId: string, addresses: AddressInput[]) {
    for (const a of addresses) {
      await db.query(
        `INSERT INTO ${kind.addressTable} (tenant_id, company_id, ${kind.fk}, address_type, is_default, building_number,
                                         street, district, city, postal_code, additional_number, country)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
        [tenantId, companyId, partyId, a.addressType, a.isDefault, a.buildingNumber ?? null, a.street ?? null,
         a.district ?? null, a.city ?? null, a.postalCode ?? null, a.additionalNumber ?? null, a.country]);
    }
  }

  return async function routes(app: FastifyInstance) {
    const canView = requirePermission(app, kind.viewPermission);
    const canManage = requirePermission(app, kind.managePermission);

    app.get(`/${kind.path}`, { preHandler: canView }, async (req) => {
      const q = parse(z.object({
        companyId: z.uuid().optional(),
        search: z.string().trim().max(100).optional(),
        status: z.enum(['active', 'inactive', 'all']).default('active'),
        limit: z.coerce.number().int().min(1).max(200).default(50),
        offset: z.coerce.number().int().min(0).default(0),
      }), req.query);
      return req.tenantTx(async (db) => {
        const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
        const search = q.search ? `%${q.search.replace(/[%_\\]/g, (m) => `\\${m}`)}%` : null;
        const { rows } = await db.query<{ total: string }>(
          `${SELECT.replace('SELECT p.id', 'SELECT count(*) OVER () AS total, p.id')}
            WHERE p.tenant_id = $1 AND p.company_id = $2 AND p.deleted_at IS NULL
              AND ($3 = 'all' OR p.is_active = ($3 = 'active'))
              AND ($4::text IS NULL OR p.name_ar ILIKE $4 OR p.name_en ILIKE $4 OR p.code ILIKE $4
                   OR p.vat_number ILIKE $4 OR p.phone ILIKE $4 OR p.email::text ILIKE $4)
            ORDER BY p.name_ar LIMIT $5 OFFSET $6`,
          [req.auth!.tenantId, companyId, q.status, search, q.limit, q.offset]);
        return { data: rows.map(({ total: _t, ...r }) => r), total: Number(rows[0]?.total ?? 0) };
      });
    });

    app.get(`/${kind.path}/:id`, { preHandler: canView }, async (req) => {
      const { id } = parse(uuidParam, req.params);
      const party = await req.tenantTx((db) => load(db, req.auth!.tenantId, id));
      if (!party) throw notFound(kind.label);
      return party;
    });

    app.post(`/${kind.path}`, { preHandler: canManage }, async (req, reply) => {
      const body = parse(createBody, req.body);
      checkAddresses(body.addresses);
      const a = req.auth!;
      const party = await req.tenantTx(async (db) => {
        const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
        await checkControlAccount(db, companyId, body.controlAccountId);
        const code = body.code ?? await nextFreeCode(db, a.tenantId, companyId, kind.docType, { prefix: kind.prefix }, kind.table, 'code');
        const { rows: [row] } = await db.query<{ id: string }>(
          `INSERT INTO ${kind.table} (tenant_id, company_id, code, party_type, name_ar, name_en, vat_number,
             commercial_registration, national_id, email, phone, credit_limit, payment_terms_days, ${kind.accountColumn},
             notes, created_by, updated_by)
           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $16) RETURNING id`,
          [a.tenantId, companyId, code, body.partyType, body.nameAr, body.nameEn ?? null, body.vatNumber ?? null,
           body.commercialRegistration ?? null, body.nationalId ?? null, body.email ?? null, body.phone ?? null,
           body.creditLimit ?? null, body.paymentTermsDays, body.controlAccountId ?? null, body.notes ?? null, a.userId]);
        await insertAddresses(db, a.tenantId, companyId, row!.id, body.addresses);
        const created = await load(db, a.tenantId, row!.id);
        await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: kind.entity, entityId: row!.id, newValues: created }, req.auditMeta());
        return created;
      });
      reply.code(201);
      return party;
    });

    app.patch(`/${kind.path}/:id`, { preHandler: canManage }, async (req) => {
      const { id } = parse(uuidParam, req.params);
      const body = parse(updateBody, req.body);
      const a = req.auth!;
      return req.tenantTx(async (db) => {
        const before = await load(db, a.tenantId, id);
        if (!before) throw notFound(kind.label);
        if (body.controlAccountId) await checkControlAccount(db, before.companyId, body.controlAccountId);
        // Whitelisted column names only; values are parameters.
        const map: [keyof typeof body, string][] = [
          ['code', 'code'], ['partyType', 'party_type'], ['nameAr', 'name_ar'], ['nameEn', 'name_en'],
          ['vatNumber', 'vat_number'], ['commercialRegistration', 'commercial_registration'], ['nationalId', 'national_id'],
          ['email', 'email'], ['phone', 'phone'], ['creditLimit', 'credit_limit'], ['paymentTermsDays', 'payment_terms_days'],
          ['controlAccountId', kind.accountColumn], ['notes', 'notes'], ['isActive', 'is_active'],
        ];
        const sets: string[] = [];
        const values: unknown[] = [];
        for (const [key, col] of map) {
          if (body[key] !== undefined) { values.push(body[key]); sets.push(`${col} = $${values.length + 3}`); }
        }
        if (sets.length) {
          await db.query(`UPDATE ${kind.table} SET ${sets.join(', ')}, updated_by = $3 WHERE tenant_id = $1 AND id = $2`,
            [a.tenantId, id, a.userId, ...values]);
        }
        const after = await load(db, a.tenantId, id);
        await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: kind.entity, entityId: id, oldValues: before, newValues: after }, req.auditMeta());
        return after;
      });
    });

    /** Replaces the address list. Old addresses are soft-deleted (documents keep their own copy). */
    app.put(`/${kind.path}/:id/addresses`, { preHandler: canManage }, async (req) => {
      const { id } = parse(uuidParam, req.params);
      const body = parse(z.object({ addresses: z.array(address).max(10) }), req.body);
      checkAddresses(body.addresses);
      const a = req.auth!;
      return req.tenantTx(async (db) => {
        const before = await load(db, a.tenantId, id);
        if (!before) throw notFound(kind.label);
        await db.query(`UPDATE ${kind.addressTable} SET deleted_at = now() WHERE ${kind.fk} = $1 AND deleted_at IS NULL`, [id]);
        await insertAddresses(db, a.tenantId, before.companyId, id, body.addresses);
        const after = await load(db, a.tenantId, id);
        await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: kind.entity, entityId: id, oldValues: { addresses: before.addresses }, newValues: { addresses: after!.addresses } }, req.auditMeta());
        return after;
      });
    });

    app.delete(`/${kind.path}/:id`, { preHandler: canManage }, async (req, reply) => {
      const { id } = parse(uuidParam, req.params);
      const a = req.auth!;
      await req.tenantTx(async (db) => {
        const before = await load(db, a.tenantId, id);
        if (!before) throw notFound(kind.label);
        const docTables = kind.table === 'customers' ? ['sales_quotes', 'sales_invoices', 'sales_returns'] : ['purchase_orders', 'purchase_invoices', 'purchase_returns'];
        const used = await db.query(
          `SELECT 1 FROM journal_entry_lines WHERE ${kind.fk} = $1
           UNION ALL SELECT 1 FROM payments WHERE ${kind.fk} = $1
           ${docTables.map((t) => `UNION ALL SELECT 1 FROM ${t} WHERE ${kind.fk} = $1`).join(' ')} LIMIT 1`, [id]);
        if (used.rowCount) throw conflict('PARTY_IN_USE', `This ${kind.entity} has accounting transactions. Deactivate it instead.`);
        await db.query(`UPDATE ${kind.table} SET deleted_at = now(), is_active = false, updated_by = $2 WHERE id = $1`, [id, a.userId]);
        await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'DELETE', entityType: kind.entity, entityId: id, oldValues: before }, req.auditMeta());
      });
      return reply.code(204).send();
    });
  };
}
