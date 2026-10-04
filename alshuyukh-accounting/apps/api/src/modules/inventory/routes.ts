import type { FastifyInstance, FastifyRequest } from 'fastify';
import { z } from 'zod';
import type { Db } from '../../db/tx.js';
import { isoDate } from '../../lib/dates.js';
import { badRequest, notFound } from '../../lib/errors.js';
import { Decimal, toMoney } from '../../lib/money.js';
import { nextDocumentNumber } from '../../lib/sequences.js';
import { parse, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { resolveCompanyId } from '../accounting/company-context.js';
import { postEntry } from '../accounting/engine.js';
import { writeAudit } from '../audit/audit.service.js';
import { systemAccount } from '../documents/posting.js';
import { issue, receive, type Ref } from './engine.js';

const quantity = z.string().trim().regex(/^\d{1,14}(\.\d{1,4})?$/, 'Quantity: up to 4 decimals, sent as a string');
const cost = z.string().trim().regex(/^\d{1,14}(\.\d{1,6})?$/, 'Unit cost: up to 6 decimals, sent as a string');

async function checkWarehouse(db: Db, companyId: string, warehouseId: string) {
  const { rowCount } = await db.query(`SELECT 1 FROM warehouses WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL AND is_active`, [companyId, warehouseId]);
  if (!rowCount) throw badRequest('INVALID_WAREHOUSE', 'Warehouse does not exist in this company or is inactive');
}

async function checkStocked(db: Db, companyId: string, productIds: string[]) {
  const ids = [...new Set(productIds)];
  const { rows } = await db.query<{ id: string; track_inventory: boolean; name_ar: string }>(
    `SELECT id, track_inventory, name_ar FROM products WHERE company_id = $1 AND id = ANY($2::uuid[]) AND deleted_at IS NULL`, [companyId, ids]);
  if (rows.length !== ids.length) throw badRequest('INVALID_PRODUCT', 'A product does not exist in this company');
  const untracked = rows.find((p) => !p.track_inventory);
  if (untracked) throw badRequest('NOT_STOCKED', `"${untracked.name_ar}" does not track stock`);
  if (ids.length !== productIds.length) throw badRequest('INVALID_LINE', 'The same product appears twice');
}

const ctxRef = (req: FastifyRequest, companyId: string, date: string, referenceType: string, referenceId: string): Ref =>
  ({ tenantId: req.auth!.tenantId, companyId, userId: req.auth!.userId, date, referenceType, referenceId });

export default async function inventoryRoutes(app: FastifyInstance) {
  const canView = requirePermission(app, 'inventory.view');

  /** Quantity, value and average cost per product and warehouse. */
  app.get('/inventory/balances', { preHandler: canView }, async (req) => {
    const q = parse(z.object({
      companyId: z.uuid().optional(), warehouseId: z.uuid().optional(), productId: z.uuid().optional(),
      includeZero: z.enum(['true', 'false']).optional(),
    }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query(
        `SELECT b.product_id AS "productId", p.sku, p.name_ar AS "productName", u.code AS "unitCode",
                b.warehouse_id AS "warehouseId", w.name AS "warehouseName", b.quantity::text, b.value::text,
                CASE WHEN b.quantity > 0 THEN round(b.value / b.quantity, 4) ELSE 0 END::text AS "averageCost"
           FROM inventory_balances b
           JOIN products p ON p.id = b.product_id JOIN units u ON u.id = p.unit_id JOIN warehouses w ON w.id = b.warehouse_id
          WHERE b.tenant_id = $1 AND b.company_id = $2 AND ($3::uuid IS NULL OR b.warehouse_id = $3)
            AND ($4::uuid IS NULL OR b.product_id = $4) AND ($5::boolean OR b.quantity > 0)
          ORDER BY p.name_ar, w.name`,
        [req.auth!.tenantId, companyId, q.warehouseId ?? null, q.productId ?? null, q.includeZero === 'true']);
      return { data: rows };
    });
  });

  /** Stock card: movements of one product with running balances. */
  app.get('/inventory/movements', { preHandler: canView }, async (req) => {
    const q = parse(z.object({
      companyId: z.uuid().optional(), productId: z.uuid().optional(), warehouseId: z.uuid().optional(),
      dateFrom: isoDate.optional(), dateTo: isoDate.optional(),
      limit: z.coerce.number().int().min(1).max(500).default(100), offset: z.coerce.number().int().min(0).default(0),
    }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query<{ total_count: string }>(
        `SELECT count(*) OVER () AS total_count, m.id, m.movement_date AS date, m.movement_type AS type, m.direction,
                m.product_id AS "productId", p.name_ar AS "productName", p.sku, m.warehouse_id AS "warehouseId", w.name AS "warehouseName",
                m.quantity::text, m.unit_cost::text AS "unitCost", m.total_cost::text AS "totalCost",
                m.balance_quantity::text AS "balanceQuantity", m.balance_value::text AS "balanceValue",
                m.reference_type AS "referenceType", m.reference_id AS "referenceId", m.created_at AS "createdAt"
           FROM stock_movements m JOIN products p ON p.id = m.product_id JOIN warehouses w ON w.id = m.warehouse_id
          WHERE m.tenant_id = $1 AND m.company_id = $2 AND ($3::uuid IS NULL OR m.product_id = $3)
            AND ($4::uuid IS NULL OR m.warehouse_id = $4) AND ($5::date IS NULL OR m.movement_date >= $5)
            AND ($6::date IS NULL OR m.movement_date <= $6)
          ORDER BY m.created_at DESC, m.id DESC LIMIT $7 OFFSET $8`,
        [req.auth!.tenantId, companyId, q.productId ?? null, q.warehouseId ?? null, q.dateFrom ?? null, q.dateTo ?? null, q.limit, q.offset]);
      return { data: rows.map(({ total_count: _t, ...r }) => r), total: Number(rows[0]?.total_count ?? 0) };
    });
  });

  /**
   * Inventory valuation, reconciled with the general ledger: the stock value
   * must equal the balance of the Inventory account.
   */
  app.get('/inventory/valuation', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query<{ value: string }>(
        `SELECT p.id AS "productId", p.sku, p.name_ar AS "productName", sum(b.quantity)::text AS quantity, sum(b.value)::text AS value
           FROM inventory_balances b JOIN products p ON p.id = b.product_id
          WHERE b.tenant_id = $1 AND b.company_id = $2
          GROUP BY p.id HAVING sum(b.quantity) > 0 ORDER BY p.name_ar`, [req.auth!.tenantId, companyId]);
      const stockValue = rows.reduce((s, r) => s.plus(r.value), new Decimal(0));
      const { rows: [gl] } = await db.query<{ balance: string }>(
        `SELECT COALESCE(sum(l.debit - l.credit), 0)::text AS balance
           FROM journal_entry_lines l JOIN journal_entries e ON e.id = l.journal_entry_id
           JOIN accounts a ON a.id = l.account_id
          WHERE e.company_id = $1 AND e.status IN ('POSTED', 'REVERSED') AND a.system_key = 'INVENTORY'`, [companyId]);
      const difference = new Decimal(gl!.balance).minus(stockValue);
      return {
        data: rows, stockValue: toMoney(stockValue), ledgerBalance: toMoney(gl!.balance),
        difference: toMoney(difference), reconciled: difference.isZero(),
      };
    });
  });

  // Transfers --------------------------------------------------------------------
  const TRANSFER_SELECT = `
    SELECT t.id, t.transfer_number AS number, t.transfer_date AS date, t.from_warehouse_id AS "fromWarehouseId",
           fw.name AS "fromWarehouseName", t.to_warehouse_id AS "toWarehouseId", tw.name AS "toWarehouseName",
           t.notes, t.total_cost::text AS "totalCost", t.created_at AS "createdAt"
      FROM stock_transfers t JOIN warehouses fw ON fw.id = t.from_warehouse_id JOIN warehouses tw ON tw.id = t.to_warehouse_id`;

  async function loadTransfer(db: Db, tenantId: string, id: string) {
    const { rows: [t] } = await db.query(`${TRANSFER_SELECT} WHERE t.tenant_id = $1 AND t.id = $2`, [tenantId, id]);
    if (!t) throw notFound('Transfer');
    const { rows: lines } = await db.query(
      `SELECT i.line_no AS "lineNo", i.product_id AS "productId", p.sku, p.name_ar AS "productName", i.quantity::text, i.total_cost::text AS "totalCost"
         FROM stock_transfer_items i JOIN products p ON p.id = i.product_id WHERE i.transfer_id = $1 ORDER BY i.line_no`, [id]);
    return { ...t, lines };
  }

  app.get('/stock-transfers', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      return { data: (await db.query(`${TRANSFER_SELECT} WHERE t.tenant_id = $1 AND t.company_id = $2 ORDER BY t.transfer_date DESC, t.transfer_number DESC LIMIT 200`, [req.auth!.tenantId, companyId])).rows };
    });
  });

  app.get('/stock-transfers/:id', { preHandler: canView }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx((db) => loadTransfer(db, req.auth!.tenantId, id));
  });

  /** Moves stock between warehouses at average cost. No journal entry: the value stays in Inventory. */
  app.post('/stock-transfers', { preHandler: requirePermission(app, 'inventory.transfer') }, async (req, reply) => {
    const body = parse(z.object({
      companyId: z.uuid().optional(), fromWarehouseId: z.uuid(), toWarehouseId: z.uuid(), date: isoDate,
      notes: z.string().trim().max(1000).nullish(),
      lines: z.array(z.object({ productId: z.uuid(), quantity })).min(1).max(500),
    }), req.body);
    if (body.fromWarehouseId === body.toWarehouseId) throw badRequest('SAME_WAREHOUSE', 'Choose two different warehouses');
    const a = req.auth!;
    const transfer = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      await checkWarehouse(db, companyId, body.fromWarehouseId);
      await checkWarehouse(db, companyId, body.toWarehouseId);
      await checkStocked(db, companyId, body.lines.map((l) => l.productId));
      const number = await nextDocumentNumber(db, a.tenantId, companyId, 'STOCK_TRANSFER', { prefix: 'TRF', padding: 6 });
      const id = crypto.randomUUID();
      const ref = ctxRef(req, companyId, body.date, 'STOCK_TRANSFER', id);
      const costs: string[] = [];
      for (const [i, l] of body.lines.entries()) {
        if (!new Decimal(l.quantity).greaterThan(0)) throw badRequest('INVALID_LINE', `Line ${i + 1}: quantity must be greater than zero`);
        const out = await issue(db, ref, { productId: l.productId, warehouseId: body.fromWarehouseId, quantity: l.quantity, type: 'TRANSFER_OUT' });
        await receive(db, ref, { productId: l.productId, warehouseId: body.toWarehouseId, quantity: l.quantity, type: 'TRANSFER_IN' }, out.cost);
        costs.push(out.cost);
      }
      const total = toMoney(costs.reduce((s2, c) => s2.plus(c), new Decimal(0)));
      await db.query(
        `INSERT INTO stock_transfers (id, tenant_id, company_id, transfer_number, transfer_date, from_warehouse_id, to_warehouse_id, notes, total_cost, created_by)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
        [id, a.tenantId, companyId, number, body.date, body.fromWarehouseId, body.toWarehouseId, body.notes ?? null, total, a.userId]);
      for (const [i, l] of body.lines.entries()) {
        await db.query(
          `INSERT INTO stock_transfer_items (tenant_id, company_id, transfer_id, line_no, product_id, quantity, total_cost) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
          [a.tenantId, companyId, id, i + 1, l.productId, l.quantity, costs[i]]);
      }
      const t = { id };
      const created = await loadTransfer(db, a.tenantId, t!.id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'stock_transfer', entityId: t!.id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return transfer;
  });

  // Adjustments ------------------------------------------------------------------
  const ADJ_SELECT = `
    SELECT s.id, s.adjustment_number AS number, s.adjustment_date AS date, s.warehouse_id AS "warehouseId", w.name AS "warehouseName",
           s.reason, s.offset_account_id AS "offsetAccountId", a.code AS "offsetAccountCode", a.name_ar AS "offsetAccountName",
           s.total_increase::text AS "totalIncrease", s.total_decrease::text AS "totalDecrease",
           s.journal_entry_id AS "journalEntryId", s.created_at AS "createdAt"
      FROM stock_adjustments s JOIN warehouses w ON w.id = s.warehouse_id JOIN accounts a ON a.id = s.offset_account_id`;

  async function loadAdjustment(db: Db, tenantId: string, id: string) {
    const { rows: [s] } = await db.query(`${ADJ_SELECT} WHERE s.tenant_id = $1 AND s.id = $2`, [tenantId, id]);
    if (!s) throw notFound('Adjustment');
    const { rows: lines } = await db.query(
      `SELECT i.line_no AS "lineNo", i.product_id AS "productId", p.sku, p.name_ar AS "productName", i.direction,
              i.quantity::text, i.counted_quantity::text AS "countedQuantity", i.system_quantity::text AS "systemQuantity",
              i.unit_cost::text AS "unitCost", i.total_cost::text AS "totalCost"
         FROM stock_adjustment_items i JOIN products p ON p.id = i.product_id WHERE i.adjustment_id = $1 ORDER BY i.line_no`, [id]);
    return { ...s, lines };
  }

  app.get('/stock-adjustments', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      return { data: (await db.query(`${ADJ_SELECT} WHERE s.tenant_id = $1 AND s.company_id = $2 ORDER BY s.adjustment_date DESC, s.adjustment_number DESC LIMIT 200`, [req.auth!.tenantId, companyId])).rows };
    });
  });

  app.get('/stock-adjustments/:id', { preHandler: canView }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    return req.tenantTx((db) => loadAdjustment(db, req.auth!.tenantId, id));
  });

  /**
   * Adjusts stock in one warehouse. Each line is either a direction and
   * quantity, or a counted quantity (stock count) from which the difference
   * is derived. Increases are valued at the given unit cost or the current
   * average; decreases at average cost. The journal entry posts the
   * difference against the offset account (default: inventory adjustments;
   * use an equity account for opening stock).
   */
  app.post('/stock-adjustments', { preHandler: requirePermission(app, 'inventory.adjust') }, async (req, reply) => {
    const line = z.union([
      z.object({ productId: z.uuid(), countedQuantity: z.string().trim().regex(/^\d{1,14}(\.\d{1,4})?$/), unitCost: cost.optional() }),
      z.object({ productId: z.uuid(), direction: z.enum(['IN', 'OUT']), quantity, unitCost: cost.optional() }),
    ]);
    const body = parse(z.object({
      companyId: z.uuid().optional(), warehouseId: z.uuid(), date: isoDate, reason: z.string().trim().min(3).max(500),
      offsetAccountId: z.uuid().optional(), lines: z.array(line).min(1).max(500),
    }), req.body);
    const a = req.auth!;
    const adjustment = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      await checkWarehouse(db, companyId, body.warehouseId);
      await checkStocked(db, companyId, body.lines.map((l) => l.productId));
      const offset = body.offsetAccountId ?? await systemAccount(db, companyId, 'INVENTORY_ADJUSTMENT');
      const { rows: [acc] } = await db.query<{ account_type: string; is_postable: boolean; is_active: boolean; system_key: string | null }>(
        `SELECT account_type, is_postable, is_active, system_key FROM accounts WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, offset]);
      if (!acc || !acc.is_postable || !acc.is_active || acc.system_key === 'INVENTORY'
          || !['EXPENSE', 'COST_OF_GOODS_SOLD', 'EQUITY', 'REVENUE'].includes(acc.account_type)) {
        throw badRequest('INVALID_OFFSET_ACCOUNT', 'The offset account must be an active postable expense, cost, equity or revenue account');
      }

      const number = await nextDocumentNumber(db, a.tenantId, companyId, 'STOCK_ADJUSTMENT', { prefix: 'ADJ', padding: 6 });
      const id = crypto.randomUUID();
      const ref = ctxRef(req, companyId, body.date, 'STOCK_ADJUSTMENT', id);
      let increase = new Decimal(0);
      let decrease = new Decimal(0);
      const items: unknown[][] = [];
      for (const [i, l] of body.lines.entries()) {
        const n = i + 1;
        const { rows: [bal] } = await db.query<{ quantity: string; value: string }>(
          `SELECT quantity::text, value::text FROM inventory_balances WHERE product_id = $1 AND warehouse_id = $2`, [l.productId, body.warehouseId]);
        const systemQty = new Decimal(bal?.quantity ?? 0);
        let direction: 'IN' | 'OUT';
        let qty: Decimal;
        let counted: string | null = null;
        if ('countedQuantity' in l) {
          const diff = new Decimal(l.countedQuantity).minus(systemQty);
          if (diff.isZero()) continue; // counted matches the books
          direction = diff.greaterThan(0) ? 'IN' : 'OUT';
          qty = diff.abs();
          counted = l.countedQuantity;
        } else {
          direction = l.direction;
          qty = new Decimal(l.quantity);
          if (!qty.greaterThan(0)) throw badRequest('INVALID_LINE', `Line ${n}: quantity must be greater than zero`);
        }
        let lineCost: string;
        if (direction === 'IN') {
          let unit: Decimal;
          if (l.unitCost !== undefined) unit = new Decimal(l.unitCost);
          else if (systemQty.greaterThan(0)) unit = new Decimal(bal!.value).dividedBy(systemQty);
          else throw badRequest('UNIT_COST_REQUIRED', `Line ${n}: enter a unit cost; there is no current cost to use`);
          lineCost = toMoney(unit.times(qty).toDecimalPlaces(2, Decimal.ROUND_HALF_UP));
          await receive(db, ref, { productId: l.productId, warehouseId: body.warehouseId, quantity: qty.toString(), type: 'ADJUSTMENT_IN' }, lineCost);
          increase = increase.plus(lineCost);
        } else {
          lineCost = (await issue(db, ref, { productId: l.productId, warehouseId: body.warehouseId, quantity: qty.toString(), type: 'ADJUSTMENT_OUT' })).cost;
          decrease = decrease.plus(lineCost);
        }
        items.push([l.productId, direction, qty.toString(), counted, counted === null ? null : systemQty.toString(),
          new Decimal(lineCost).dividedBy(qty).toDecimalPlaces(6).toString(), lineCost]);
      }
      if (!items.length) throw badRequest('NOTHING_TO_ADJUST', 'Every counted quantity matches the system quantity');

      let journalEntryId: string | null = null;
      const journal = [];
      const inventory = await systemAccount(db, companyId, 'INVENTORY');
      if (increase.greaterThan(0)) journal.push({ accountId: inventory, debit: toMoney(increase) }, { accountId: offset, credit: toMoney(increase) });
      if (decrease.greaterThan(0)) journal.push({ accountId: offset, debit: toMoney(decrease) }, { accountId: inventory, credit: toMoney(decrease) });
      if (journal.length) {
        const entry = await postEntry(db, { tenantId: a.tenantId, userId: a.userId, meta: req.auditMeta() }, {
          companyId, entryDate: body.date, description: `تسوية مخزون ${number}: ${body.reason}`,
          referenceType: 'STOCK_ADJUSTMENT', referenceId: id, source: 'SYSTEM', lines: journal,
        });
        journalEntryId = entry.id;
      }
      await db.query(
        `INSERT INTO stock_adjustments (id, tenant_id, company_id, adjustment_number, adjustment_date, warehouse_id, reason,
           offset_account_id, total_increase, total_decrease, journal_entry_id, created_by)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
        [id, a.tenantId, companyId, number, body.date, body.warehouseId, body.reason, offset, toMoney(increase), toMoney(decrease), journalEntryId, a.userId]);
      for (const [i, item] of items.entries()) {
        await db.query(
          `INSERT INTO stock_adjustment_items (tenant_id, company_id, adjustment_id, line_no, product_id, direction, quantity,
             counted_quantity, system_quantity, unit_cost, total_cost)
           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, [a.tenantId, companyId, id, i + 1, ...item]);
      }
      const created = await loadAdjustment(db, a.tenantId, id);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'stock_adjustment', entityId: id, newValues: created }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return adjustment;
  });
}
