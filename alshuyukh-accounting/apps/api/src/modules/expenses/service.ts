import type { Db } from '../../db/tx.js';
import { badRequest, conflict, notFound } from '../../lib/errors.js';
import { Decimal, toMoney } from '../../lib/money.js';
import { nextDocumentNumber } from '../../lib/sequences.js';
import { postEntry, reverse, type LineInput as JournalLine } from '../accounting/engine.js';
import { writeAudit } from '../audit/audit.service.js';
import { CalcError, calculateLine, totalsOf, type VatCategory } from '../documents/calc.js';
import { vatRateFor, type Ctx } from '../documents/drafts.js';
import { controlAccount, systemAccount } from '../documents/posting.js';
import { recordTax, reverseTax, taxGroups } from '../tax/ledger.js';

/**
 * Expenses.
 *
 *   Paid now (CASH / BANK):  Dr Expense accounts + Dr VAT input  /  Cr payment method account
 *   On credit (CREDIT):      Dr Expense accounts + Dr VAT input  /  Cr Payable (supplier)
 *
 * A credit expense is settled later with a supplier payment allocated to it.
 */

export interface ExpenseLineInput {
  categoryId: string;
  amount: string;
  description?: string | null;
  vatCategory?: VatCategory;
  costCenterId?: string | null;
}

export interface ExpenseInput {
  companyId: string;
  expenseDate: string;
  paymentType: 'CASH' | 'BANK' | 'CREDIT';
  methodId?: string | null;
  supplierId?: string | null;
  payeeName?: string | null;
  reference?: string | null;
  vendorVatNumber?: string | null;
  branchId?: string | null;
  pricesIncludeVat?: boolean;
  notes?: string | null;
  lines: ExpenseLineInput[];
}

interface BuiltLine {
  categoryId: string; accountId: string; description: string | null; costCenterId: string | null; amount: string;
  netAmount: string; vatCategory: VatCategory; vatRate: string; vatAmount: string; totalAmount: string;
}

export async function buildExpenseLines(db: Db, companyId: string, input: Pick<ExpenseInput, 'lines' | 'expenseDate' | 'pricesIncludeVat'>) {
  if (!input.lines.length) throw badRequest('NO_LINES', 'An expense needs at least one line');
  const ids = [...new Set(input.lines.map((l) => l.categoryId))];
  const { rows: cats } = await db.query<{ id: string; account_id: string; vat_category: VatCategory; is_active: boolean; name_ar: string }>(
    `SELECT id, account_id, vat_category, is_active, name_ar FROM expense_categories WHERE company_id = $1 AND id = ANY($2::uuid[])`, [companyId, ids]);
  const byId = new Map(cats.map((c) => [c.id, c]));
  const ccIds = [...new Set(input.lines.map((l) => l.costCenterId).filter(Boolean))] as string[];
  if (ccIds.length) {
    const { rowCount } = await db.query(`SELECT 1 FROM cost_centers WHERE company_id = $1 AND id = ANY($2::uuid[]) AND is_active AND deleted_at IS NULL`, [companyId, ccIds]);
    if (rowCount !== ccIds.length) throw badRequest('INVALID_COST_CENTER', 'A cost center does not exist or is inactive');
  }
  const lines: BuiltLine[] = [];
  for (const [i, l] of input.lines.entries()) {
    const cat = byId.get(l.categoryId);
    if (!cat) throw badRequest('INVALID_CATEGORY', `Line ${i + 1}: expense category does not exist in this company`);
    if (!cat.is_active) throw badRequest('CATEGORY_INACTIVE', `Line ${i + 1}: "${cat.name_ar}" is inactive`);
    const category = l.vatCategory ?? cat.vat_category;
    const rate = await vatRateFor(db, companyId, category, input.expenseDate);
    let calc;
    try {
      calc = calculateLine({ quantity: '1', unitPrice: l.amount, vatCategory: category, vatRate: rate }, input.pricesIncludeVat ?? false, i + 1);
    } catch (e) {
      if (e instanceof CalcError) throw badRequest('INVALID_LINE', e.message);
      throw e;
    }
    if (!new Decimal(l.amount).greaterThan(0)) throw badRequest('INVALID_LINE', `Line ${i + 1}: amount must be greater than zero`);
    lines.push({
      categoryId: cat.id, accountId: cat.account_id, description: l.description ?? null, costCenterId: l.costCenterId ?? null,
      amount: toMoney(l.amount), netAmount: calc.netAmount, vatCategory: category, vatRate: calc.vatRate, vatAmount: calc.vatAmount,
      totalAmount: calc.totalAmount,
    });
  }
  const t = totalsOf(lines.map((l) => ({ grossAmount: l.netAmount, discountAmount: '0', netAmount: l.netAmount, vatCategory: l.vatCategory, vatRate: l.vatRate })));
  return { lines, totals: { subtotal: t.taxableAmount, taxAmount: t.taxAmount, total: t.total } };
}

async function checkHeader(db: Db, companyId: string, input: Partial<ExpenseInput>) {
  if (input.paymentType === 'CREDIT' && !input.supplierId) throw badRequest('SUPPLIER_REQUIRED', 'An expense on credit needs a supplier');
  if (input.paymentType && input.paymentType !== 'CREDIT' && !input.methodId) throw badRequest('METHOD_REQUIRED', 'Choose how the expense was paid');
  if (input.supplierId) {
    const s = await db.query(`SELECT 1 FROM suppliers WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL AND is_active`, [companyId, input.supplierId]);
    if (!s.rowCount) throw badRequest('INVALID_PARTY', 'Supplier does not exist in this company or is inactive');
  }
  if (input.methodId) {
    const m = await db.query(`SELECT 1 FROM payment_methods WHERE company_id = $1 AND id = $2 AND is_active`, [companyId, input.methodId]);
    if (!m.rowCount) throw badRequest('INVALID_METHOD', 'Payment method does not exist in this company or is inactive');
  }
  if (input.branchId) {
    const b = await db.query(`SELECT 1 FROM branches WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, input.branchId]);
    if (!b.rowCount) throw badRequest('INVALID_BRANCH', 'Branch does not exist in this company');
  }
}

async function insertLines(db: Db, ctx: Ctx, companyId: string, expenseId: string, lines: BuiltLine[]) {
  for (const [i, l] of lines.entries()) {
    await db.query(
      `INSERT INTO expense_items (tenant_id, company_id, document_id, line_no, category_id, account_id, description, cost_center_id,
         amount, net_amount, vat_category, vat_rate, vat_amount, total_amount)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
      [ctx.tenantId, companyId, expenseId, i + 1, l.categoryId, l.accountId, l.description, l.costCenterId, l.amount,
       l.netAmount, l.vatCategory, l.vatRate, l.vatAmount, l.totalAmount]);
  }
}

export async function createExpense(db: Db, ctx: Ctx, input: ExpenseInput): Promise<string> {
  await checkHeader(db, input.companyId, input);
  const { lines, totals } = await buildExpenseLines(db, input.companyId, input);
  const { rows: [e] } = await db.query<{ id: string }>(
    `INSERT INTO expenses (tenant_id, company_id, branch_id, expense_date, payment_type, method_id, supplier_id, payee_name,
       reference, vendor_vat_number, prices_include_vat, subtotal, tax_amount, total, notes, created_by)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id`,
    [ctx.tenantId, input.companyId, input.branchId ?? null, input.expenseDate, input.paymentType,
     input.paymentType === 'CREDIT' ? null : input.methodId, input.supplierId ?? null, input.payeeName ?? null,
     input.reference ?? null, input.vendorVatNumber ?? null, input.pricesIncludeVat ?? false,
     totals.subtotal, totals.taxAmount, totals.total, input.notes ?? null, ctx.userId]);
  await insertLines(db, ctx, input.companyId, e!.id, lines);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'CREATE', entityType: 'expense', entityId: e!.id, newValues: { ...input, totals } }, ctx.meta);
  return e!.id;
}

async function lockExpense(db: Db, tenantId: string, id: string) {
  const { rows: [e] } = await db.query<{
    id: string; company_id: string; status: string; expense_date: string; payment_type: 'CASH' | 'BANK' | 'CREDIT';
    method_id: string | null; supplier_id: string | null; payee_name: string | null; vendor_vat_number: string | null;
    prices_include_vat: boolean; branch_id: string | null; journal_entry_id: string | null; expense_number: string | null;
    paid_amount: string; deleted_at: Date | null; [k: string]: unknown;
  }>(`SELECT *, paid_amount::text FROM expenses WHERE tenant_id = $1 AND id = $2 FOR UPDATE`, [tenantId, id]);
  if (!e || e.deleted_at) throw notFound('Expense');
  return e;
}

export async function updateExpense(db: Db, ctx: Ctx, id: string, input: Partial<Omit<ExpenseInput, 'companyId'>>) {
  const e = await lockExpense(db, ctx.tenantId, id);
  if (e.status !== 'DRAFT') throw conflict('DOCUMENT_NOT_DRAFT', 'Only draft expenses can be edited');
  const before = await loadExpense(db, ctx.tenantId, id);
  const merged = {
    expenseDate: input.expenseDate ?? e.expense_date, paymentType: input.paymentType ?? e.payment_type,
    methodId: input.methodId !== undefined ? input.methodId : e.method_id,
    supplierId: input.supplierId !== undefined ? input.supplierId : e.supplier_id,
    payeeName: input.payeeName !== undefined ? input.payeeName : e.payee_name,
    reference: input.reference !== undefined ? input.reference : (e.reference as string | null),
    vendorVatNumber: input.vendorVatNumber !== undefined ? input.vendorVatNumber : e.vendor_vat_number,
    branchId: input.branchId !== undefined ? input.branchId : e.branch_id,
    pricesIncludeVat: input.pricesIncludeVat ?? e.prices_include_vat,
    notes: input.notes !== undefined ? input.notes : (e.notes as string | null),
    lines: input.lines ?? before.lines.map((l) => ({ categoryId: l.categoryId, amount: l.amount, description: l.description, vatCategory: l.vatCategory, costCenterId: l.costCenterId })),
  };
  await checkHeader(db, e.company_id, merged);
  const { lines, totals } = await buildExpenseLines(db, e.company_id, merged);
  await db.query(
    `UPDATE expenses SET expense_date = $2, payment_type = $3, method_id = $4, supplier_id = $5, payee_name = $6, reference = $7,
       vendor_vat_number = $8, branch_id = $9, prices_include_vat = $10, notes = $11, subtotal = $12, tax_amount = $13, total = $14
     WHERE id = $1`,
    [id, merged.expenseDate, merged.paymentType, merged.paymentType === 'CREDIT' ? null : merged.methodId, merged.supplierId,
     merged.payeeName, merged.reference, merged.vendorVatNumber, merged.branchId, merged.pricesIncludeVat, merged.notes,
     totals.subtotal, totals.taxAmount, totals.total]);
  await db.query(`DELETE FROM expense_items WHERE document_id = $1`, [id]);
  await insertLines(db, ctx, e.company_id, id, lines);
  const after = await loadExpense(db, ctx.tenantId, id);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'UPDATE', entityType: 'expense', entityId: id, oldValues: before, newValues: after }, ctx.meta);
  return after;
}

export async function deleteExpense(db: Db, ctx: Ctx, id: string) {
  const e = await lockExpense(db, ctx.tenantId, id);
  if (e.status !== 'DRAFT') throw conflict('DOCUMENT_NOT_DRAFT', 'Posted expenses cannot be deleted. Cancel them instead.');
  await db.query(`UPDATE expenses SET deleted_at = now() WHERE id = $1`, [id]);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'DELETE', entityType: 'expense', entityId: id }, ctx.meta);
}

export async function postExpense(db: Db, ctx: Ctx, id: string) {
  const e = await lockExpense(db, ctx.tenantId, id);
  if (e.status !== 'DRAFT') throw conflict('DOCUMENT_NOT_DRAFT', 'This expense is already posted');
  const current = await loadExpense(db, ctx.tenantId, id);
  // Recalculate with the rates in force on the expense date.
  const { lines, totals } = await buildExpenseLines(db, e.company_id, {
    expenseDate: e.expense_date, pricesIncludeVat: e.prices_include_vat,
    lines: current.lines.map((l) => ({ categoryId: l.categoryId, amount: l.amount, description: l.description, vatCategory: l.vatCategory, costCenterId: l.costCenterId })),
  });
  await checkHeader(db, e.company_id, { paymentType: e.payment_type, methodId: e.method_id, supplierId: e.supplier_id });

  const number = await nextDocumentNumber(db, ctx.tenantId, e.company_id, 'EXPENSE', { prefix: 'EXP', padding: 6 });
  let creditAccount: string;
  if (e.payment_type === 'CREDIT') {
    creditAccount = await controlAccount(db, 'supplier', e.company_id, e.supplier_id!);
  } else {
    const { rows: [m] } = await db.query<{ account_id: string }>(`SELECT account_id FROM payment_methods WHERE id = $1`, [e.method_id]);
    creditAccount = m!.account_id;
  }
  // Debit each expense account and cost center once.
  const debits = new Map<string, { accountId: string; costCenterId: string | null; amount: Decimal }>();
  for (const l of lines) {
    const key = `${l.accountId}|${l.costCenterId ?? ''}`;
    const d = debits.get(key) ?? { accountId: l.accountId, costCenterId: l.costCenterId, amount: new Decimal(0) };
    d.amount = d.amount.plus(l.netAmount);
    debits.set(key, d);
  }
  const journal: JournalLine[] = [...debits.values()].filter((d) => !d.amount.isZero())
    .map((d) => ({ accountId: d.accountId, debit: toMoney(d.amount), costCenterId: d.costCenterId, branchId: e.branch_id }));
  if (new Decimal(totals.taxAmount).greaterThan(0)) {
    journal.push({ accountId: await systemAccount(db, e.company_id, 'VAT_INPUT'), debit: totals.taxAmount, branchId: e.branch_id });
  }
  journal.push({
    accountId: creditAccount, credit: totals.total, branchId: e.branch_id,
    ...(e.payment_type === 'CREDIT' ? { supplierId: e.supplier_id } : {}),
  });
  const entry = await postEntry(db, ctx, {
    companyId: e.company_id, entryDate: e.expense_date, description: `مصروف ${number}${e.payee_name ? ` — ${e.payee_name}` : ''}`,
    referenceType: 'EXPENSE', referenceId: id, source: 'SYSTEM', lines: journal,
  });

  const paidNow = e.payment_type !== 'CREDIT';
  await db.query(`DELETE FROM expense_items WHERE document_id = $1`, [id]);
  await insertLines(db, ctx, e.company_id, id, lines);
  await db.query(
    `UPDATE expenses SET status = $2, expense_number = $3, journal_entry_id = $4, posted_by = $5, posted_at = now(),
       subtotal = $6, tax_amount = $7, total = $8, paid_amount = $9 WHERE id = $1`,
    [id, paidNow ? 'PAID' : 'POSTED', number, entry.id, ctx.userId, totals.subtotal, totals.taxAmount, totals.total,
     paidNow ? totals.total : '0']);

  const { rows: [party] } = await db.query<{ name: string | null; vat: string | null }>(
    `SELECT COALESCE(s.name_ar, $2) AS name, COALESCE($3, s.vat_number) AS vat FROM (SELECT 1) x LEFT JOIN suppliers s ON s.id = $1`,
    [e.supplier_id, e.payee_name, e.vendor_vat_number]);
  await recordTax(db, {
    tenantId: ctx.tenantId, companyId: e.company_id, sourceType: 'EXPENSE', sourceId: id, sourceNumber: number,
    journalEntryId: entry.id, date: e.expense_date, groups: taxGroups(lines, totals.taxAmount),
    partyName: party?.name ?? null, partyVatNumber: party?.vat ?? null,
  });
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'POST', entityType: 'expense', entityId: id, newValues: { number, total: totals.total, journalEntryId: entry.id } }, ctx.meta);
  return loadExpense(db, ctx.tenantId, id);
}

/** Cancels a posted expense by reversing its entry and VAT. Credit expenses must have no payments. */
export async function cancelExpense(db: Db, ctx: Ctx, id: string, reason: string) {
  const e = await lockExpense(db, ctx.tenantId, id);
  if (e.status === 'DRAFT') throw conflict('DOCUMENT_NOT_ISSUED', 'Delete the draft instead');
  if (e.status === 'CANCELLED') throw conflict('ALREADY_CANCELLED', 'This expense is already cancelled');
  const allocated = await db.query(`SELECT 1 FROM payment_allocations WHERE expense_id = $1 AND reversed_at IS NULL LIMIT 1`, [id]);
  if (allocated.rowCount) throw conflict('HAS_PAYMENTS', 'Void the payments allocated to this expense first');
  await reverse(db, ctx, e.journal_entry_id!, { reason: `إلغاء مصروف ${e.expense_number}: ${reason}`, date: e.expense_date, allowSystem: true });
  await reverseTax(db, 'EXPENSE', id);
  await db.query(`UPDATE expenses SET status = 'CANCELLED', cancelled_by = $2, cancelled_at = now(), cancel_reason = $3 WHERE id = $1`, [id, ctx.userId, reason]);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'CANCEL', entityType: 'expense', entityId: id, newValues: { reason } }, ctx.meta);
  return loadExpense(db, ctx.tenantId, id);
}

export const EXPENSE_SELECT = `
  SELECT e.id, e.company_id AS "companyId", e.expense_number AS number, e.expense_date AS date, e.payment_type AS "paymentType",
         e.method_id AS "methodId", m.name_ar AS "methodName", e.supplier_id AS "supplierId", s.name_ar AS "supplierName",
         e.payee_name AS "payeeName", e.reference, e.vendor_vat_number AS "vendorVatNumber", e.branch_id AS "branchId",
         e.prices_include_vat AS "pricesIncludeVat", e.subtotal::text, e.tax_amount::text AS "taxAmount", e.total::text,
         e.paid_amount::text AS "paidAmount", e.remaining_amount::text AS "remainingAmount", e.notes, e.status,
         e.journal_entry_id AS "journalEntryId", e.cancel_reason AS "cancelReason", e.created_at AS "createdAt"
    FROM expenses e
    LEFT JOIN payment_methods m ON m.id = e.method_id
    LEFT JOIN suppliers s ON s.id = e.supplier_id`;

export async function loadExpense(db: Db, tenantId: string, id: string) {
  const { rows: [e] } = await db.query<Record<string, unknown>>(`${EXPENSE_SELECT} WHERE e.tenant_id = $1 AND e.id = $2 AND e.deleted_at IS NULL`, [tenantId, id]);
  if (!e) throw notFound('Expense');
  const { rows: lines } = await db.query<{
    id: string; categoryId: string; amount: string; description: string | null; vatCategory: VatCategory; costCenterId: string | null;
  }>(
    `SELECT i.id, i.line_no AS "lineNo", i.category_id AS "categoryId", c.name_ar AS "categoryName", i.account_id AS "accountId",
            a.code AS "accountCode", i.description, i.cost_center_id AS "costCenterId", i.amount::text, i.net_amount::text AS "netAmount",
            i.vat_category AS "vatCategory", i.vat_rate::text AS "vatRate", i.vat_amount::text AS "vatAmount", i.total_amount::text AS "totalAmount"
       FROM expense_items i JOIN expense_categories c ON c.id = i.category_id JOIN accounts a ON a.id = i.account_id
      WHERE i.document_id = $1 ORDER BY i.line_no`, [id]);
  return { ...e, lines };
}
