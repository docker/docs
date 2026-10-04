import type { Db } from '../../db/tx.js';
import { badRequest, conflict } from '../../lib/errors.js';
import { Decimal, toMoney } from '../../lib/money.js';
import { COSTING, type Balance, type CostingMethod } from './costing.js';

/**
 * Inventory engine. Every quantity change goes through receive() or issue(),
 * which lock the balance row, write an append-only stock movement, and
 * update the balance in the caller's transaction.
 *
 * Accounting stays with the caller: it posts the movement's cost to the
 * Inventory account in the same journal entry, so the stock ledger and the
 * general ledger always hold the same value.
 */

export type MovementType =
  | 'PURCHASE' | 'SALE' | 'SALE_RETURN' | 'PURCHASE_RETURN' | 'TRANSFER_IN' | 'TRANSFER_OUT'
  | 'ADJUSTMENT_IN' | 'ADJUSTMENT_OUT' | 'CANCELLATION_IN' | 'CANCELLATION_OUT';

export interface Ref {
  tenantId: string;
  companyId: string;
  userId: string;
  date: string;
  referenceType: string;
  referenceId: string;
}

export interface MovementInput {
  productId: string;
  warehouseId: string;
  quantity: string;
  type: MovementType;
  lineId?: string | null;
}

async function costingFor(db: Db, companyId: string): Promise<CostingMethod> {
  const { rows: [c] } = await db.query<{ costing_method: string }>(`SELECT costing_method FROM companies WHERE id = $1`, [companyId]);
  return COSTING[c!.costing_method]!;
}

async function lockBalance(db: Db, ref: Ref, productId: string, warehouseId: string): Promise<Balance> {
  await db.query(
    `INSERT INTO inventory_balances (tenant_id, company_id, product_id, warehouse_id) VALUES ($1, $2, $3, $4)
     ON CONFLICT (product_id, warehouse_id) DO NOTHING`,
    [ref.tenantId, ref.companyId, productId, warehouseId]);
  const { rows: [b] } = await db.query<Balance>(
    `SELECT quantity::text, value::text FROM inventory_balances WHERE product_id = $1 AND warehouse_id = $2 FOR UPDATE`,
    [productId, warehouseId]);
  return b!;
}

async function write(db: Db, ref: Ref, m: MovementInput, direction: 'IN' | 'OUT', cost: string, after: Balance, reversesId?: string) {
  await db.query(`UPDATE inventory_balances SET quantity = $3, value = $4 WHERE product_id = $1 AND warehouse_id = $2`,
    [m.productId, m.warehouseId, after.quantity, after.value]);
  const unitCost = new Decimal(cost).dividedBy(m.quantity).toDecimalPlaces(6).toString();
  const { rows: [row] } = await db.query<{ id: string }>(
    `INSERT INTO stock_movements (tenant_id, company_id, product_id, warehouse_id, movement_type, direction, quantity,
       unit_cost, total_cost, balance_quantity, balance_value, reference_type, reference_id, reference_line_id,
       reverses_id, movement_date, created_by)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id`,
    [ref.tenantId, ref.companyId, m.productId, m.warehouseId, m.type, direction, m.quantity, unitCost, cost,
     after.quantity, after.value, ref.referenceType, ref.referenceId, m.lineId ?? null, reversesId ?? null, ref.date, ref.userId]);
  return row!.id;
}

/** Adds stock at a known total cost (purchase price, or the cost a return or transfer carries). */
export async function receive(db: Db, ref: Ref, m: MovementInput, totalCost: string, reversesId?: string) {
  const method = await costingFor(db, ref.companyId);
  const before = await lockBalance(db, ref, m.productId, m.warehouseId);
  const after = method.receive(before, m.quantity, toMoney(totalCost));
  const id = await write(db, ref, m, 'IN', toMoney(totalCost), after, reversesId);
  return { id, cost: toMoney(totalCost) };
}

async function productName(db: Db, productId: string) {
  return (await db.query<{ name_ar: string }>(`SELECT name_ar FROM products WHERE id = $1`, [productId])).rows[0]?.name_ar ?? productId;
}

/** Removes stock at its current cost. Negative stock is not allowed. */
export async function issue(db: Db, ref: Ref, m: MovementInput) {
  const method = await costingFor(db, ref.companyId);
  const before = await lockBalance(db, ref, m.productId, m.warehouseId);
  if (new Decimal(m.quantity).greaterThan(before.quantity)) {
    throw conflict('INSUFFICIENT_STOCK', `Not enough stock of "${await productName(db, m.productId)}": ${new Decimal(before.quantity).toString()} available, ${new Decimal(m.quantity).toString()} needed`);
  }
  const { cost, balance } = method.issue(before, m.quantity);
  const id = await write(db, ref, m, 'OUT', cost, balance);
  return { id, cost };
}

/**
 * Removes stock at a fixed cost, to undo an earlier receipt (cancelling a
 * purchase or a sales return). If that empties the warehouse while value is
 * left over, the leftover is returned so the caller can expense it.
 */
export async function issueAtCost(db: Db, ref: Ref, m: MovementInput, totalCost: string, reversesId: string) {
  const before = await lockBalance(db, ref, m.productId, m.warehouseId);
  const name = await productName(db, m.productId);
  if (new Decimal(m.quantity).greaterThan(before.quantity) || new Decimal(totalCost).greaterThan(before.value)) {
    throw conflict('INSUFFICIENT_STOCK', `The stock of "${name}" received by this document has already been used. Issue a return instead.`);
  }
  const qty = new Decimal(before.quantity).minus(m.quantity);
  let value = new Decimal(before.value).minus(totalCost);
  let residual = new Decimal(0);
  if (qty.isZero() && !value.isZero()) { residual = value; value = new Decimal(0); }
  const id = await write(db, ref, m, 'OUT', toMoney(totalCost), { quantity: qty.toString(), value: toMoney(value) }, reversesId);
  return { id, residual: toMoney(residual) };
}

/**
 * Undoes every movement of a document with opposite movements at the same
 * cost. Returns any residual value to expense (see issueAtCost).
 */
export async function reverseDocumentMovements(db: Db, ref: Ref, originalType: string, originalId: string) {
  const { rows } = await db.query<{ id: string; product_id: string; warehouse_id: string; direction: 'IN' | 'OUT'; quantity: string; total_cost: string; reference_line_id: string | null }>(
    `SELECT m.id, m.product_id, m.warehouse_id, m.direction, m.quantity::text, m.total_cost::text, m.reference_line_id
       FROM stock_movements m
      WHERE m.reference_type = $1 AND m.reference_id = $2 AND m.reverses_id IS NULL
        AND NOT EXISTS (SELECT 1 FROM stock_movements r WHERE r.reverses_id = m.id)
      ORDER BY m.created_at`, [originalType, originalId]);
  let residual = new Decimal(0);
  for (const m of rows) {
    const input: MovementInput = {
      productId: m.product_id, warehouseId: m.warehouse_id, quantity: m.quantity, lineId: m.reference_line_id,
      type: m.direction === 'IN' ? 'CANCELLATION_OUT' : 'CANCELLATION_IN',
    };
    if (m.direction === 'OUT') await receive(db, ref, input, m.total_cost, m.id);
    else residual = residual.plus((await issueAtCost(db, ref, input, m.total_cost, m.id)).residual);
  }
  return toMoney(residual);
}

/** The document's warehouse, or the company's main warehouse. */
export async function resolveWarehouse(db: Db, companyId: string, warehouseId?: string | null): Promise<string> {
  if (warehouseId) return warehouseId;
  const { rows: [w] } = await db.query<{ id: string }>(
    `SELECT id FROM warehouses WHERE company_id = $1 AND deleted_at IS NULL AND is_active
      ORDER BY (code = 'MAIN') DESC, created_at LIMIT 1`, [companyId]);
  if (!w) throw badRequest('WAREHOUSE_REQUIRED', 'Choose a warehouse for stocked items');
  return w.id;
}

/** Products on these lines that hold stock. */
export async function trackedProducts(db: Db, productIds: (string | null)[]): Promise<Set<string>> {
  const ids = [...new Set(productIds.filter(Boolean))] as string[];
  if (!ids.length) return new Set();
  const { rows } = await db.query<{ id: string }>(`SELECT id FROM products WHERE id = ANY($1::uuid[]) AND track_inventory`, [ids]);
  return new Set(rows.map((r) => r.id));
}
