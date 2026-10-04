import { randomUUID } from 'node:crypto';
import type { Db } from '../../db/tx.js';
import { badRequest, conflict, notFound } from '../../lib/errors.js';
import { Decimal, toMoney } from '../../lib/money.js';
import { nextDocumentNumber } from '../../lib/sequences.js';
import { postEntry, reverse } from '../accounting/engine.js';
import { writeAudit } from '../audit/audit.service.js';
import type { Ctx } from './drafts.js';
import { KINDS, type DocKind } from './kinds.js';
import { controlAccount, refreshInvoiceStatus } from './posting.js';

/**
 * Payments move money between a cash/bank account (from the payment method)
 * and the party's control account:
 *
 *   Receipt from customer       Dr Cash/Bank  /  Cr Receivable (customer)
 *   Refund to customer          Dr Receivable (customer)  /  Cr Cash/Bank
 *   Payment to supplier         Dr Payable (supplier)  /  Cr Cash/Bank
 *   Refund from supplier        Dr Cash/Bank  /  Cr Payable (supplier)
 *
 * Allocation links a payment to the documents it settles. It does not post
 * anything: the ledger is already right per party. An unallocated remainder
 * stays on the party's account as an advance.
 */

export type Direction = 'RECEIPT' | 'DISBURSEMENT';
export type TargetType = 'SALES_INVOICE' | 'SALES_RETURN' | 'PURCHASE_INVOICE' | 'PURCHASE_RETURN' | 'EXPENSE';

export interface AllocationInput { documentType: TargetType; documentId: string; amount: string }

export interface PaymentInput {
  companyId: string;
  direction: Direction;
  customerId?: string | null;
  supplierId?: string | null;
  paymentDate: string;
  methodId: string;
  amount: string;
  reference?: string | null;
  notes?: string | null;
  branchId?: string | null;
  allocations?: AllocationInput[];
}

const TARGET_COLUMN: Record<TargetType, string> = {
  SALES_INVOICE: 'sales_invoice_id', SALES_RETURN: 'sales_return_id',
  PURCHASE_INVOICE: 'purchase_invoice_id', PURCHASE_RETURN: 'purchase_return_id', EXPENSE: 'expense_id',
};

/** Which documents a payment may settle, by party and direction. */
function allowedTargets(party: 'customer' | 'supplier', direction: Direction): TargetType[] {
  if (party === 'customer') return direction === 'RECEIPT' ? ['SALES_INVOICE'] : ['SALES_RETURN'];
  return direction === 'DISBURSEMENT' ? ['PURCHASE_INVOICE', 'EXPENSE'] : ['PURCHASE_RETURN'];
}

/** Where a settlement target lives and how its paid amount and status are kept. */
interface Target { table: string; partyColumn: string; settleColumn: string; numberColumn: string; refresh(db: Db, id: string): Promise<void> }

function target(type: TargetType): Target {
  if (type === 'EXPENSE') {
    return {
      table: 'expenses', partyColumn: 'supplier_id', settleColumn: 'paid_amount', numberColumn: 'expense_number',
      refresh: async (db, id) => {
        await db.query(
          `UPDATE expenses SET status = CASE WHEN remaining_amount = 0 THEN 'PAID' WHEN paid_amount > 0 THEN 'PARTIALLY_PAID' ELSE 'POSTED' END
            WHERE id = $1 AND status <> 'CANCELLED'`, [id]);
      },
    };
  }
  const kind: DocKind = KINDS[type];
  return {
    table: kind.table, partyColumn: `${kind.party}_id`, settleColumn: kind.isReturn ? 'refunded_amount' : 'paid_amount', numberColumn: 'doc_number',
    refresh: (db, id) => (kind.isReturn ? Promise.resolve() : refreshInvoiceStatus(db, kind, id)),
  };
}

export async function createPayment(db: Db, ctx: Ctx, input: PaymentInput) {
  const party = input.customerId ? 'customer' : 'supplier';
  const partyId = (input.customerId ?? input.supplierId)!;
  if (!!input.customerId === !!input.supplierId) throw badRequest('INVALID_PARTY', 'Choose either a customer or a supplier');
  const { rows: [p] } = await db.query<{ is_active: boolean }>(
    `SELECT is_active FROM ${party}s WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [input.companyId, partyId]);
  if (!p) throw badRequest('INVALID_PARTY', `The ${party} does not exist in this company`);

  const { rows: [method] } = await db.query<{ account_id: string; is_active: boolean; name_ar: string }>(
    `SELECT account_id, is_active, name_ar FROM payment_methods WHERE company_id = $1 AND id = $2`, [input.companyId, input.methodId]);
  if (!method) throw badRequest('INVALID_METHOD', 'Payment method does not exist in this company');
  if (!method.is_active) throw badRequest('METHOD_INACTIVE', 'Payment method is inactive');
  if (input.branchId) {
    const b = await db.query(`SELECT 1 FROM branches WHERE company_id = $1 AND id = $2`, [input.companyId, input.branchId]);
    if (!b.rowCount) throw badRequest('INVALID_BRANCH', 'Branch does not exist in this company');
  }

  const number = await nextDocumentNumber(db, ctx.tenantId, input.companyId, input.direction === 'RECEIPT' ? 'RECEIPT' : 'DISBURSEMENT',
    { prefix: input.direction === 'RECEIPT' ? 'RCPT' : 'PAY', padding: 6 });
  const control = await controlAccount(db, party, input.companyId, partyId);
  const tag = party === 'customer' ? { customerId: partyId } : { supplierId: partyId };
  const moneyIn = input.direction === 'RECEIPT';
  const paymentId = randomUUID();
  const entry = await postEntry(db, ctx, {
    companyId: input.companyId, entryDate: input.paymentDate,
    description: `${moneyIn ? 'سند قبض' : 'سند صرف'} ${number}${input.reference ? ` — ${input.reference}` : ''}`,
    referenceType: moneyIn ? 'PAYMENT_RECEIPT' : 'PAYMENT_DISBURSEMENT', referenceId: paymentId, source: 'SYSTEM',
    lines: [
      { accountId: method.account_id, ...(moneyIn ? { debit: input.amount } : { credit: input.amount }), branchId: input.branchId ?? null },
      { accountId: control, ...(moneyIn ? { credit: input.amount } : { debit: input.amount }), ...tag, branchId: input.branchId ?? null },
    ],
  });
  const { rows: [row] } = await db.query<{ id: string }>(
    `INSERT INTO payments (id, tenant_id, company_id, branch_id, payment_number, direction, customer_id, supplier_id, payment_date,
       method_id, amount, reference, notes, journal_entry_id, created_by)
     VALUES ($15, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`,
    [ctx.tenantId, input.companyId, input.branchId ?? null, number, input.direction, input.customerId ?? null, input.supplierId ?? null,
     input.paymentDate, input.methodId, input.amount, input.reference ?? null, input.notes ?? null, entry.id, ctx.userId, paymentId]);
  await writeAudit(db, {
    tenantId: ctx.tenantId, userId: ctx.userId, action: 'PAYMENT', entityType: 'payment', entityId: row!.id,
    newValues: { number, direction: input.direction, party, partyId, amount: input.amount, journalEntryId: entry.id },
  }, ctx.meta);
  if (input.allocations?.length) await allocate(db, ctx, row!.id, input.allocations);
  return loadPayment(db, ctx.tenantId, row!.id);
}

/** Settles documents with an existing payment's unallocated amount. */
export async function allocate(db: Db, ctx: Ctx, paymentId: string, allocations: AllocationInput[]) {
  const { rows: [pay] } = await db.query<{
    id: string; company_id: string; status: string; direction: Direction; customer_id: string | null; supplier_id: string | null;
    unallocated_amount: string; payment_date: string;
  }>(
    `SELECT id, company_id, status, direction, customer_id, supplier_id, unallocated_amount::text, payment_date
       FROM payments WHERE tenant_id = $1 AND id = $2 FOR UPDATE`, [ctx.tenantId, paymentId]);
  if (!pay) throw notFound('Payment');
  if (pay.status !== 'POSTED') throw conflict('PAYMENT_VOIDED', 'This payment is voided');
  const party = pay.customer_id ? 'customer' : 'supplier';
  const partyId = (pay.customer_id ?? pay.supplier_id)!;
  const expected = allowedTargets(party, pay.direction);

  const total = allocations.reduce((s, a) => s.plus(a.amount), new Decimal(0));
  if (total.greaterThan(pay.unallocated_amount)) {
    throw badRequest('OVER_ALLOCATED', `Allocations (${toMoney(total)}) exceed the unallocated amount (${pay.unallocated_amount})`);
  }
  const seen = new Set<string>();
  for (const a of allocations) {
    if (!expected.includes(a.documentType)) throw badRequest('INVALID_ALLOCATION', `This payment can only settle: ${expected.join(', ')}`);
    if (seen.has(a.documentId)) throw badRequest('INVALID_ALLOCATION', 'The same document appears twice');
    seen.add(a.documentId);
    const amount = new Decimal(a.amount);
    if (!amount.greaterThan(0)) throw badRequest('INVALID_ALLOCATION', 'Allocation amounts must be positive');
    const tg = target(a.documentType);
    const { rows: [doc] } = await db.query<{ status: string; party_id: string; remaining_amount: string; doc_number: string }>(
      `SELECT status, ${tg.partyColumn} AS party_id, remaining_amount::text, ${tg.numberColumn} AS doc_number
         FROM ${tg.table} WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL FOR UPDATE`, [pay.company_id, a.documentId]);
    if (!doc || doc.party_id !== partyId) throw badRequest('INVALID_ALLOCATION', `Document does not belong to this ${party}`);
    if (doc.status === 'DRAFT' || doc.status === 'CANCELLED') throw conflict('DOCUMENT_NOT_OPEN', `${doc.doc_number ?? 'The draft'} is not open for settlement`);
    if (amount.greaterThan(doc.remaining_amount)) {
      throw badRequest('OVER_ALLOCATED', `${doc.doc_number}: only ${doc.remaining_amount} is still open`);
    }
    await db.query(`UPDATE ${tg.table} SET ${tg.settleColumn} = ${tg.settleColumn} + $2 WHERE id = $1`, [a.documentId, toMoney(amount)]);
    await tg.refresh(db, a.documentId);
    await db.query(
      `INSERT INTO payment_allocations (tenant_id, company_id, payment_id, ${TARGET_COLUMN[a.documentType]}, amount, created_by)
       VALUES ($1, $2, $3, $4, $5, $6)`, [ctx.tenantId, pay.company_id, paymentId, a.documentId, toMoney(amount), ctx.userId]);
  }
  await db.query(`UPDATE payments SET allocated_amount = allocated_amount + $2 WHERE id = $1`, [paymentId, toMoney(total)]);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'UPDATE', entityType: 'payment', entityId: paymentId, newValues: { allocations } }, ctx.meta);
  return loadPayment(db, ctx.tenantId, paymentId);
}

/** Voids a payment: reverses its entry and releases every allocation. */
export async function voidPayment(db: Db, ctx: Ctx, paymentId: string, reason: string) {
  const { rows: [pay] } = await db.query<{ status: string; journal_entry_id: string; payment_date: string; payment_number: string }>(
    `SELECT status, journal_entry_id, payment_date, payment_number FROM payments WHERE tenant_id = $1 AND id = $2 FOR UPDATE`,
    [ctx.tenantId, paymentId]);
  if (!pay) throw notFound('Payment');
  if (pay.status === 'VOIDED') throw conflict('ALREADY_VOIDED', 'This payment is already voided');

  const { rows: allocations } = await db.query<{ id: string; amount: string; target: TargetType; document_id: string }>(
    `SELECT id, amount::text,
            CASE WHEN sales_invoice_id IS NOT NULL THEN 'SALES_INVOICE' WHEN sales_return_id IS NOT NULL THEN 'SALES_RETURN'
                 WHEN purchase_invoice_id IS NOT NULL THEN 'PURCHASE_INVOICE' WHEN expense_id IS NOT NULL THEN 'EXPENSE'
                 ELSE 'PURCHASE_RETURN' END AS target,
            COALESCE(sales_invoice_id, sales_return_id, purchase_invoice_id, purchase_return_id, expense_id) AS document_id
       FROM payment_allocations WHERE payment_id = $1 AND reversed_at IS NULL`, [paymentId]);
  for (const a of allocations) {
    const tg = target(a.target);
    await db.query(`SELECT 1 FROM ${tg.table} WHERE id = $1 FOR UPDATE`, [a.document_id]);
    await db.query(`UPDATE ${tg.table} SET ${tg.settleColumn} = ${tg.settleColumn} - $2 WHERE id = $1`, [a.document_id, a.amount]);
    await tg.refresh(db, a.document_id);
    await db.query(`UPDATE payment_allocations SET reversed_at = now() WHERE id = $1`, [a.id]);
  }
  await reverse(db, ctx, pay.journal_entry_id, { reason: `إلغاء ${pay.payment_number}: ${reason}`, date: pay.payment_date, allowSystem: true });
  await db.query(
    `UPDATE payments SET status = 'VOIDED', allocated_amount = 0, voided_by = $2, voided_at = now(), void_reason = $3 WHERE id = $1`,
    [paymentId, ctx.userId, reason]);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'VOID', entityType: 'payment', entityId: paymentId, newValues: { reason } }, ctx.meta);
  return loadPayment(db, ctx.tenantId, paymentId);
}

export const PAYMENT_SELECT = `
  SELECT p.id, p.company_id AS "companyId", p.payment_number AS number, p.direction, p.payment_date AS date,
         p.customer_id AS "customerId", p.supplier_id AS "supplierId",
         COALESCE(c.name_ar, s.name_ar) AS "partyName", p.method_id AS "methodId", m.name_ar AS "methodName",
         p.amount::text, p.allocated_amount::text AS "allocatedAmount", p.unallocated_amount::text AS "unallocatedAmount",
         p.reference, p.notes, p.status, p.journal_entry_id AS "journalEntryId", p.void_reason AS "voidReason",
         p.created_at AS "createdAt"
    FROM payments p
    JOIN payment_methods m ON m.id = p.method_id
    LEFT JOIN customers c ON c.id = p.customer_id
    LEFT JOIN suppliers s ON s.id = p.supplier_id`;

export async function loadPayment(db: Db, tenantId: string, id: string) {
  const { rows: [p] } = await db.query(`${PAYMENT_SELECT} WHERE p.tenant_id = $1 AND p.id = $2`, [tenantId, id]);
  if (!p) throw notFound('Payment');
  const { rows: allocations } = await db.query(
    `SELECT a.id, a.amount::text, a.reversed_at AS "reversedAt",
            COALESCE(si.doc_number, sr.doc_number, pi.doc_number, pr.doc_number, ex.expense_number) AS "documentNumber",
            CASE WHEN a.sales_invoice_id IS NOT NULL THEN 'SALES_INVOICE' WHEN a.sales_return_id IS NOT NULL THEN 'SALES_RETURN'
                 WHEN a.purchase_invoice_id IS NOT NULL THEN 'PURCHASE_INVOICE' WHEN a.expense_id IS NOT NULL THEN 'EXPENSE'
                 ELSE 'PURCHASE_RETURN' END AS "documentType",
            COALESCE(a.sales_invoice_id, a.sales_return_id, a.purchase_invoice_id, a.purchase_return_id, a.expense_id) AS "documentId"
       FROM payment_allocations a
       LEFT JOIN sales_invoices si ON si.id = a.sales_invoice_id
       LEFT JOIN sales_returns sr ON sr.id = a.sales_return_id
       LEFT JOIN purchase_invoices pi ON pi.id = a.purchase_invoice_id
       LEFT JOIN purchase_returns pr ON pr.id = a.purchase_return_id
       LEFT JOIN expenses ex ON ex.id = a.expense_id
      WHERE a.payment_id = $1 ORDER BY a.created_at`, [id]);
  return { ...p, allocations };
}
