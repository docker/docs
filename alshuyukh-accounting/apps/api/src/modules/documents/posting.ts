import type { Db } from '../../db/tx.js';
import { badRequest, conflict } from '../../lib/errors.js';
import { Decimal, toMoney } from '../../lib/money.js';
import { nextDocumentNumber } from '../../lib/sequences.js';
import { postEntry, reverse, type LineInput as JournalLine } from '../accounting/engine.js';
import { writeAudit } from '../audit/audit.service.js';
import { issue, receive, resolveWarehouse, reverseDocumentMovements, trackedProducts, type Ref } from '../inventory/engine.js';
import { recordTax, reverseTax, taxGroups, type TaxSource } from '../tax/ledger.js';
import { totalsOf } from './calc.js';
import {
  buildLines, buildReturn, createDraft, defaultDueDate, insertItems, loadDocument, lockDocument,
  type BuiltLine, type Ctx, type LineInput,
} from './drafts.js';
import { KINDS, partyColumn, partyTable, type DocKind } from './kinds.js';

/**
 * Issuing documents and the journal entries they create.
 *
 *  Sales invoice     Dr Receivable (customer)  /  Cr Revenue per account, Cr VAT output
 *  Sales return      Dr Revenue, Dr VAT output  /  Cr Receivable (customer)
 *  Purchase invoice  Dr Inventory/Expense per account, Dr VAT input  /  Cr Payable (supplier)
 *  Purchase return   Dr Payable (supplier)  /  Cr Inventory/Expense, Cr VAT input
 *
 * Stocked goods also move inventory, in the same entry:
 *  Sales invoice     Dr Cost of goods sold / Cr Inventory (average cost)
 *  Sales return      Dr Inventory / Cr Cost of goods sold (the cost the goods left at)
 *  Purchase invoice  Inventory is the line account (purchase price)
 *  Purchase return   Cr Inventory at average cost; the difference from the
 *                    purchase price goes to cost of goods sold
 *
 * The document, its stock movements and its entry are written in the same transaction.
 */

export async function systemAccount(db: Db, companyId: string, key: string): Promise<string> {
  const { rows: [a] } = await db.query<{ id: string }>(
    `SELECT id FROM accounts WHERE company_id = $1 AND system_key = $2 AND deleted_at IS NULL AND is_active AND is_postable`,
    [companyId, key]);
  if (!a) throw conflict('SYSTEM_ACCOUNT_MISSING', `The ${key} account is missing or inactive`);
  return a.id;
}

/** The party's own control account, or the company default (receivable / payable). */
export async function controlAccount(db: Db, party: 'customer' | 'supplier', companyId: string, partyId: string): Promise<string> {
  const column = party === 'customer' ? 'receivable_account_id' : 'payable_account_id';
  const { rows: [p] } = await db.query<{ account_id: string | null }>(
    `SELECT ${column} AS account_id FROM ${party}s WHERE id = $1`, [partyId]);
  return p?.account_id ?? systemAccount(db, companyId, party === 'customer' ? 'ACCOUNTS_RECEIVABLE' : 'ACCOUNTS_PAYABLE');
}

/** Posting account of every line: revenue for sales, inventory/expense for purchases. */
async function lineAccounts(db: Db, kind: DocKind, companyId: string, lines: BuiltLine[]): Promise<string[]> {
  const productIds = [...new Set(lines.map((l) => l.productId).filter(Boolean))] as string[];
  const { rows } = await db.query<{ id: string; sales_account_id: string | null; purchase_account_id: string | null; track_inventory: boolean; name_ar: string }>(
    `SELECT id, sales_account_id, purchase_account_id, track_inventory, name_ar FROM products WHERE id = ANY($1::uuid[])`, [productIds]);
  const products = new Map(rows.map((p) => [p.id, p]));
  const result: string[] = [];
  for (const [i, l] of lines.entries()) {
    if (l.accountId) { result.push(l.accountId); continue; }
    const p = products.get(l.productId!)!;
    if (kind.party === 'customer') {
      result.push(p.sales_account_id ?? await systemAccount(db, companyId, 'SALES'));
    } else if (p.track_inventory) {
      // Stocked goods always go through the Inventory account so the stock
      // ledger and the general ledger stay equal.
      result.push(await systemAccount(db, companyId, 'INVENTORY'));
    } else if (p.purchase_account_id) {
      result.push(p.purchase_account_id);
    } else {
      throw badRequest('ACCOUNT_REQUIRED', `Line ${i + 1}: choose an expense account for "${p.name_ar}" or set a purchase account on the product`);
    }
  }
  return result;
}

async function partySnapshot(db: Db, kind: DocKind, partyId: string) {
  const { rows: [p] } = await db.query(
    `SELECT p.code, p.name_ar AS "nameAr", p.name_en AS "nameEn", p.vat_number AS "vatNumber",
            p.commercial_registration AS "commercialRegistration", p.party_type AS "partyType",
            (SELECT jsonb_build_object('buildingNumber', a.building_number, 'street', a.street, 'district', a.district,
                      'city', a.city, 'postalCode', a.postal_code, 'additionalNumber', a.additional_number, 'country', a.country)
               FROM ${kind.party}_addresses a
              WHERE a.${kind.party}_id = p.id AND a.address_type = 'BILLING' AND a.is_default AND a.deleted_at IS NULL) AS address
       FROM ${partyTable(kind)} p WHERE p.id = $1`, [partyId]);
  return p;
}

/** Re-derives an invoice's status from its payments and issued returns. */
export async function refreshInvoiceStatus(db: Db, kind: DocKind, invoiceId: string) {
  const ret = kind.key === 'SALES_INVOICE' ? KINDS.SALES_RETURN : KINDS.PURCHASE_RETURN;
  const { rows: [inv] } = await db.query<{ status: string; paid_amount: string; remaining_amount: string; fully_returned: boolean }>(
    `SELECT d.status, d.paid_amount::text, d.remaining_amount::text,
            NOT EXISTS (
              SELECT 1 FROM ${kind.itemsTable} i
               WHERE i.document_id = d.id AND i.quantity > COALESCE((
                 SELECT sum(ri.quantity) FROM ${ret.itemsTable} ri JOIN ${ret.table} r ON r.id = ri.document_id
                  WHERE ri.source_item_id = i.id AND r.status = 'ISSUED'), 0)
            ) AS fully_returned
       FROM ${kind.table} d WHERE d.id = $1`, [invoiceId]);
  if (!inv || inv.status === 'CANCELLED' || inv.status === 'DRAFT') return;
  const status = inv.fully_returned ? 'RETURNED'
    : new Decimal(inv.remaining_amount).isZero() ? 'PAID'
    : new Decimal(inv.paid_amount).greaterThan(0) ? 'PARTIALLY_PAID'
    : kind.issuedStatus;
  if (status !== inv.status) await db.query(`UPDATE ${kind.table} SET status = $2 WHERE id = $1`, [invoiceId, status]);
}

/** Recalculates a draft with current tax rates right before issue. */
async function recalculate(db: Db, kind: DocKind, ctx: Ctx, d: Awaited<ReturnType<typeof lockDocument>>) {
  const current = await loadDocument(db, kind, ctx.tenantId, d.id);
  let lines: BuiltLine[];
  let totals;
  let fullyReturned = false;
  if (kind.isReturn) {
    const r = await buildReturn(db, kind, ctx.tenantId, {
      originalInvoiceId: d.original_invoice_id as string, docDate: d.doc_date, reason: d.reason as string,
      lines: current.lines.map((l) => ({ sourceItemId: l.sourceItemId!, quantity: l.quantity })),
    }, d.id);
    ({ lines, totals, fullyReturned } = r);
  } else {
    const input: LineInput[] = current.lines.map((l) => ({
      productId: l.productId, accountId: l.accountId, description: l.description, quantity: l.quantity,
      unitPrice: l.unitPrice, discountAmount: l.discountBasis, unitId: l.unitId, vatCategory: l.vatCategory,
    }));
    lines = await buildLines(db, kind, d.company_id, { lines: input, docDate: d.doc_date, pricesIncludeVat: d.prices_include_vat as boolean });
    totals = totalsOf(lines);
  }
  return { lines, totals, fullyReturned };
}

export async function issueDocument(db: Db, kind: DocKind, ctx: Ctx, id: string) {
  if (!kind.legal) throw badRequest('NOT_POSTABLE', `${kind.label}s do not post to the ledger`);
  const d = await lockDocument(db, kind, ctx.tenantId, id);
  if (d.status !== 'DRAFT') throw conflict('DOCUMENT_NOT_DRAFT', 'This document has already been issued');

  const { lines, totals } = await recalculate(db, kind, ctx, d);
  if (!new Decimal(totals.total).greaterThan(0)) throw badRequest('ZERO_TOTAL', 'A document with a zero total cannot be issued');

  const { rows: [party] } = await db.query<{ is_active: boolean; credit_limit: string | null; vat_number: string | null }>(
    `SELECT is_active, credit_limit::text, vat_number FROM ${partyTable(kind)} WHERE id = $1`, [d.party_id]);
  if (!party!.is_active && !kind.isReturn) throw conflict('PARTY_INACTIVE', `The ${kind.party} is inactive`);

  if (kind.key === 'SALES_INVOICE' && party!.credit_limit !== null) {
    const { rows: [bal] } = await db.query<{ balance: string }>(
      `SELECT COALESCE(sum(l.debit - l.credit), 0)::text AS balance
         FROM journal_entry_lines l JOIN journal_entries e ON e.id = l.journal_entry_id
        WHERE l.customer_id = $1 AND e.status IN ('POSTED', 'REVERSED')`, [d.party_id]);
    const exposure = new Decimal(bal!.balance).plus(totals.total);
    if (exposure.greaterThan(party!.credit_limit)) {
      throw conflict('CREDIT_LIMIT_EXCEEDED', `This invoice takes the customer's balance to ${toMoney(exposure)}, above the credit limit of ${toMoney(party!.credit_limit)}`);
    }
  }

  // Lines are final: store them with their posting accounts.
  const accounts = await lineAccounts(db, kind, d.company_id, lines);
  lines.forEach((l, i) => { l.accountId = accounts[i]!; });
  await db.query(`DELETE FROM ${kind.itemsTable} WHERE document_id = $1`, [id]);
  const lineIds = await insertItems(db, kind, ctx.tenantId, d.company_id, id, lines);

  const number = await nextDocumentNumber(db, ctx.tenantId, d.company_id, kind.key, { prefix: kind.prefix, padding: 6 });
  const control = await controlAccount(db, kind.party, d.company_id, d.party_id);
  const vatAccount = await systemAccount(db, d.company_id, kind.party === 'customer' ? 'VAT_OUTPUT' : 'VAT_INPUT');
  const tag = kind.party === 'customer' ? { customerId: d.party_id } : { supplierId: d.party_id };
  const branchId = d.branch_id;

  // Revenue / cost side, grouped per account.
  const perAccount = new Map<string, Decimal>();
  const addTo = (accountId: string, amount: Decimal.Value) => perAccount.set(accountId, (perAccount.get(accountId) ?? new Decimal(0)).plus(amount));
  lines.forEach((l) => addTo(l.accountId!, l.netAmount));

  // Stock movements. `cogs` > 0 means Dr Cost of goods sold / Cr Inventory.
  const stock = await moveStock(db, kind, ctx, d, lines, lineIds);
  const cogs = stock.cogs;
  if (stock.warehouseId) d.warehouse_id = stock.warehouseId;
  if (kind.key === 'PURCHASE_RETURN' && stock.costs.size) {
    // Inventory leaves at average cost; the gap to the purchase price is a cost difference.
    // The line accounts are credited: Inventory at cost, the difference to cost of goods sold
    // (a credit when the purchase price was above average, a debit when below).
    const inventory = await systemAccount(db, d.company_id, 'INVENTORY');
    const cogsAccount = await systemAccount(db, d.company_id, 'COGS');
    for (const [i, l] of lines.entries()) {
      const cost = stock.costs.get(i);
      if (cost === undefined) continue;
      addTo(inventory, new Decimal(cost).minus(l.netAmount));
      addTo(cogsAccount, new Decimal(l.netAmount).minus(cost));
    }
  }

  // Sales invoice & purchase return credit the line accounts; the other two debit them.
  const linesSideCredit = kind.key === 'SALES_INVOICE' || kind.key === 'PURCHASE_RETURN';
  /** A signed amount on the given side; negative amounts flip to the other side. */
  const side = (credit: boolean, amount: Decimal.Value) => {
    const a = new Decimal(amount);
    return (credit !== a.isNegative()) ? { credit: toMoney(a.abs()) } : { debit: toMoney(a.abs()) };
  };
  const journal: JournalLine[] = [
    { accountId: control, ...side(!linesSideCredit, totals.total), ...tag, branchId },
    ...[...perAccount.entries()].filter(([, v]) => !v.isZero()).map(([accountId, v]) => ({ accountId, ...side(linesSideCredit, v), branchId })),
  ];
  if (new Decimal(totals.taxAmount).greaterThan(0)) journal.push({ accountId: vatAccount, ...side(linesSideCredit, totals.taxAmount), branchId });
  if (!cogs.isZero()) {
    journal.push({ accountId: await systemAccount(db, d.company_id, 'COGS'), ...side(false, cogs), branchId });
    journal.push({ accountId: await systemAccount(db, d.company_id, 'INVENTORY'), ...side(true, cogs), branchId });
  }

  const entry = await postEntry(db, ctx, {
    companyId: d.company_id, entryDate: d.doc_date, description: `${kind.labelAr} ${number}`,
    referenceType: kind.key, referenceId: id, source: 'SYSTEM', lines: journal,
  });

  const sets: [string, unknown][] = [
    ['status', kind.issuedStatus], ['doc_number', number], ['journal_entry_id', entry.id], ['issued_by', ctx.userId],
    ['issued_at', new Date()], ['party_snapshot', JSON.stringify(await partySnapshot(db, kind, d.party_id))],
    ['subtotal', totals.subtotal], ['discount_total', totals.discountTotal], ['taxable_amount', totals.taxableAmount],
    ['tax_amount', totals.taxAmount], ['total', totals.total], ['warehouse_id', d.warehouse_id],
  ];
  if (kind.key === 'SALES_INVOICE' || kind.key === 'PURCHASE_INVOICE') {
    sets.push(['due_date', d.due_date ?? await defaultDueDate(db, kind, d.party_id, d.doc_date)]);
  }
  if (kind.key === 'SALES_INVOICE') sets.push(['invoice_kind', party!.vat_number ? 'STANDARD' : 'SIMPLIFIED']);

  let originalKind: DocKind | null = null;
  if (kind.isReturn) {
    // Apply the credit to the original invoice up to what is still open on it.
    originalKind = KINDS[kind.originalKind!];
    const { rows: [orig] } = await db.query<{ remaining_amount: string; status: string }>(
      `SELECT remaining_amount::text, status FROM ${originalKind.table} WHERE id = $1 FOR UPDATE`, [d.original_invoice_id]);
    if (orig!.status === 'CANCELLED') throw conflict('INVOICE_CANCELLED', 'The original invoice was cancelled');
    const applied = Decimal.min(totals.total, orig!.remaining_amount);
    sets.push(['applied_amount', toMoney(applied)]);
    await db.query(`UPDATE ${originalKind.table} SET returned_amount = returned_amount + $2 WHERE id = $1`, [d.original_invoice_id, toMoney(applied)]);
  }
  const snapshot = JSON.parse(sets.find(([c]) => c === 'party_snapshot')![1] as string) as { nameAr: string; vatNumber: string | null };
  await db.query(`UPDATE ${kind.table} SET ${sets.map(([c], i) => `${c} = $${i + 2}`).join(', ')} WHERE id = $1`, [id, ...sets.map(([, v]) => v)]);
  if (originalKind) await refreshInvoiceStatus(db, originalKind, d.original_invoice_id as string);

  await recordTax(db, {
    tenantId: ctx.tenantId, companyId: d.company_id, sourceType: kind.key as TaxSource, sourceId: id, sourceNumber: number,
    journalEntryId: entry.id, date: d.doc_date, groups: taxGroups(lines, totals.taxAmount),
    partyName: snapshot.nameAr, partyVatNumber: snapshot.vatNumber,
  });

  await writeAudit(db, {
    tenantId: ctx.tenantId, userId: ctx.userId, action: 'POST', entityType: kind.entity, entityId: id,
    newValues: { number, total: totals.total, taxAmount: totals.taxAmount, journalEntryId: entry.id },
  }, ctx.meta);
  return loadDocument(db, kind, ctx.tenantId, id);
}

/**
 * Cancels an issued invoice or return that nothing has been settled
 * against: its journal entry is reversed on the original date.
 * TODO(Phase 8): once an invoice is reported to ZATCA it can only be corrected by a return (credit note).
 */
export async function cancelDocument(db: Db, kind: DocKind, ctx: Ctx, id: string, reason: string) {
  const d = await lockDocument(db, kind, ctx.tenantId, id);
  if (!kind.legal) {
    const cancellable = kind.key === 'SALES_QUOTE' ? ['DRAFT', 'SENT', 'ACCEPTED'] : ['DRAFT', 'APPROVED'];
    if (!cancellable.includes(d.status)) throw conflict('CANNOT_CANCEL', `A ${d.status.toLowerCase()} ${kind.label.toLowerCase()} cannot be cancelled`);
  } else {
    if (d.status === 'DRAFT') throw conflict('DOCUMENT_NOT_ISSUED', 'Delete the draft instead');
    if (d.status === 'CANCELLED') throw conflict('ALREADY_CANCELLED', 'This document is already cancelled');
    if (kind.isReturn) {
      if (new Decimal(d.refunded_amount as string).greaterThan(0)) throw conflict('HAS_PAYMENTS', 'Void the refund payments first');
    } else {
      if (new Decimal(d.paid_amount as string).greaterThan(0)) throw conflict('HAS_PAYMENTS', 'Void the payments allocated to this invoice first');
      const ret = kind.key === 'SALES_INVOICE' ? KINDS.SALES_RETURN : KINDS.PURCHASE_RETURN;
      const returns = await db.query(`SELECT 1 FROM ${ret.table} WHERE original_invoice_id = $1 AND status = 'ISSUED' LIMIT 1`, [id]);
      if (returns.rowCount) throw conflict('HAS_RETURNS', 'Cancel the returns issued against this invoice first');
    }
    await reverse(db, ctx, d.journal_entry_id!, { reason: `إلغاء ${kind.labelAr} ${d.doc_number}: ${reason}`, date: d.doc_date, allowSystem: true });
    await reverseTax(db, kind.key as TaxSource, id);
    const residual = await reverseDocumentMovements(db, stockRef(ctx, d.company_id, d.doc_date, kind.key, id), kind.key, id);
    if (new Decimal(residual).greaterThan(0)) {
      // The warehouse emptied with value left over: expense it so the ledgers stay equal.
      await postEntry(db, ctx, {
        companyId: d.company_id, entryDate: d.doc_date, description: `تسوية قيمة مخزون بعد إلغاء ${kind.labelAr} ${d.doc_number}`,
        referenceType: 'INVENTORY_RESIDUAL', referenceId: id, source: 'SYSTEM',
        lines: [
          { accountId: await systemAccount(db, d.company_id, 'COGS'), debit: residual },
          { accountId: await systemAccount(db, d.company_id, 'INVENTORY'), credit: residual },
        ],
      });
    }
  }
  await db.query(`UPDATE ${kind.table} SET status = 'CANCELLED', cancelled_by = $2, cancelled_at = now(), cancel_reason = $3 WHERE id = $1`,
    [id, ctx.userId, reason]);
  if (kind.isReturn) {
    const orig = KINDS[kind.originalKind!];
    await db.query(`SELECT 1 FROM ${orig.table} WHERE id = $1 FOR UPDATE`, [d.original_invoice_id]);
    await db.query(`UPDATE ${orig.table} SET returned_amount = returned_amount - $2 WHERE id = $1`, [d.original_invoice_id, d.applied_amount]);
    await refreshInvoiceStatus(db, orig, d.original_invoice_id as string);
  }
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'CANCEL', entityType: kind.entity, entityId: id, newValues: { reason } }, ctx.meta);
  return loadDocument(db, kind, ctx.tenantId, id);
}

const stockRef = (ctx: Ctx, companyId: string, date: string, referenceType: string, referenceId: string): Ref =>
  ({ tenantId: ctx.tenantId, companyId, userId: ctx.userId, date, referenceType, referenceId });

/**
 * Records the stock movements of an invoice or return being issued.
 * Returns the net cost moved to cost of goods sold and, for purchase
 * returns, the cost of each line.
 */
async function moveStock(db: Db, kind: DocKind, ctx: Ctx, d: Awaited<ReturnType<typeof lockDocument>>, lines: BuiltLine[], lineIds: string[]) {
  const tracked = await trackedProducts(db, lines.map((l) => l.productId));
  const result = { cogs: new Decimal(0), costs: new Map<number, string>(), warehouseId: null as string | null };
  if (!tracked.size) return result;
  const ref = stockRef(ctx, d.company_id, d.doc_date, kind.key, d.id);
  const docWarehouse = await resolveWarehouse(db, d.company_id, d.warehouse_id as string | null);
  result.warehouseId = docWarehouse;

  for (const [i, l] of lines.entries()) {
    if (!l.productId || !tracked.has(l.productId)) continue;
    const move = { productId: l.productId, quantity: l.quantity, lineId: lineIds[i]! };
    if (kind.key === 'SALES_INVOICE') {
      const { cost } = await issue(db, ref, { ...move, warehouseId: docWarehouse, type: 'SALE' });
      result.cogs = result.cogs.plus(cost);
    } else if (kind.key === 'PURCHASE_INVOICE') {
      await receive(db, ref, { ...move, warehouseId: docWarehouse, type: 'PURCHASE' }, l.netAmount);
    } else {
      // Returns go back to (or leave from) the warehouse of the original movement.
      const origType = kind.key === 'SALES_RETURN' ? 'SALE' : 'PURCHASE';
      const { rows: [orig] } = await db.query<{ warehouse_id: string; quantity: string; total_cost: string }>(
        `SELECT warehouse_id, quantity::text, total_cost::text FROM stock_movements
          WHERE reference_line_id = $1 AND movement_type = $2 LIMIT 1`, [l.sourceItemId, origType]);
      const warehouseId = orig?.warehouse_id ?? docWarehouse;
      if (kind.key === 'PURCHASE_RETURN') {
        const { cost } = await issue(db, ref, { ...move, warehouseId, type: 'PURCHASE_RETURN' });
        result.costs.set(i, cost);
        continue;
      }
      if (!orig) continue; // sold before stock tracking: nothing to bring back at cost
      // Sales return: the cost the goods left at, prorated; the last return takes the rest.
      const { rows: [prior] } = await db.query<{ qty: string; cost: string }>(
        `SELECT COALESCE(sum(m.quantity), 0)::text AS qty, COALESCE(sum(m.total_cost), 0)::text AS cost
           FROM stock_movements m
           JOIN ${kind.itemsTable} ri ON ri.id = m.reference_line_id
           JOIN ${kind.table} r ON r.id = ri.document_id
          WHERE ri.source_item_id = $1 AND m.movement_type = 'SALE_RETURN' AND r.status = 'ISSUED'`, [l.sourceItemId]);
      const remainingQty = new Decimal(orig.quantity).minus(prior!.qty);
      const cost = new Decimal(l.quantity).equals(remainingQty)
        ? new Decimal(orig.total_cost).minus(prior!.cost)
        : new Decimal(orig.total_cost).times(l.quantity).dividedBy(orig.quantity).toDecimalPlaces(2, Decimal.ROUND_HALF_UP);
      await receive(db, ref, { ...move, warehouseId, type: 'SALE_RETURN' }, toMoney(cost));
      result.cogs = result.cogs.minus(cost);
    }
  }
  return result;
}

const TRANSITIONS: Record<string, Record<string, string[]>> = {
  SALES_QUOTE: { DRAFT: ['SENT'], SENT: ['ACCEPTED', 'REJECTED'] },
  PURCHASE_ORDER: { DRAFT: ['APPROVED'] },
};

/** Moves a quote or purchase order along its workflow. */
export async function setStatus(db: Db, kind: DocKind, ctx: Ctx, id: string, status: string) {
  const d = await lockDocument(db, kind, ctx.tenantId, id);
  const allowed = TRANSITIONS[kind.key]?.[d.status] ?? [];
  if (!allowed.includes(status)) throw conflict('INVALID_TRANSITION', `Cannot move from ${d.status} to ${status}`);
  await db.query(`UPDATE ${kind.table} SET status = $2 WHERE id = $1`, [id, status]);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'UPDATE', entityType: kind.entity, entityId: id, oldValues: { status: d.status }, newValues: { status } }, ctx.meta);
  return loadDocument(db, kind, ctx.tenantId, id);
}

/** Copies a quote into a draft invoice, or an approved order into a draft purchase invoice. */
export async function convertDocument(db: Db, kind: DocKind, ctx: Ctx, id: string, docDate: string) {
  const target = kind.key === 'SALES_QUOTE' ? KINDS.SALES_INVOICE : kind.key === 'PURCHASE_ORDER' ? KINDS.PURCHASE_INVOICE : null;
  if (!target) throw badRequest('NOT_CONVERTIBLE', 'Only quotes and purchase orders can be converted');
  const d = await lockDocument(db, kind, ctx.tenantId, id);
  const ok = kind.key === 'SALES_QUOTE' ? ['DRAFT', 'SENT', 'ACCEPTED'] : ['APPROVED'];
  if (!ok.includes(d.status)) throw conflict('INVALID_TRANSITION', `A ${d.status.toLowerCase()} ${kind.label.toLowerCase()} cannot be converted`);
  const source = await loadDocument(db, kind, ctx.tenantId, id);
  const invoiceId = await createDraft(db, target, ctx, {
    companyId: d.company_id, partyId: d.party_id, docDate, branchId: d.branch_id, warehouseId: d.warehouse_id as string | null,
    pricesIncludeVat: d.prices_include_vat as boolean, notes: d.notes as string | null, sourceId: id,
    lines: source.lines.map((l) => ({
      productId: l.productId, accountId: l.accountId, description: l.description, quantity: l.quantity,
      unitPrice: l.unitPrice, discountAmount: l.discountBasis, unitId: l.unitId, vatCategory: l.vatCategory,
    })),
  });
  await db.query(`UPDATE ${kind.table} SET status = 'CONVERTED', converted_invoice_id = $2 WHERE id = $1`, [id, invoiceId]);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'UPDATE', entityType: kind.entity, entityId: id, newValues: { status: 'CONVERTED', invoiceId } }, ctx.meta);
  return loadDocument(db, target, ctx.tenantId, invoiceId);
}

