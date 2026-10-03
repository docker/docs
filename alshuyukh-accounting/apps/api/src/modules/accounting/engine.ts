import type { Db } from '../../db/tx.js';
import { AppError, badRequest, conflict, notFound } from '../../lib/errors.js';
import { Decimal, sum, toMoney } from '../../lib/money.js';
import { writeAudit, type AuditMeta } from '../audit/audit.service.js';
import { PROFIT_AND_LOSS_TYPES } from './chart-template.js';

/**
 * Accounting Engine.
 *
 * Every financial event in the system becomes a journal entry through this
 * module. All functions take the caller's transaction (`db`) so a business
 * document and its journal entry commit or roll back together.
 *
 * Life cycle: DRAFT → POSTED → (REVERSED). Posted entries are immutable; the
 * database refuses changes even if this code had a bug. Corrections are made
 * with a reversal entry followed by a new correcting entry.
 */

export interface EngineContext {
  tenantId: string;
  userId: string;
  meta?: AuditMeta;
}

export interface LineInput {
  accountId: string;
  debit?: string;
  credit?: string;
  description?: string | null;
  costCenterId?: string | null;
  branchId?: string | null;
}

export interface EntryInput {
  companyId: string;
  entryDate: string;
  description: string;
  lines: LineInput[];
  referenceType?: string;
  referenceId?: string | null;
  source?: 'MANUAL' | 'SYSTEM';
  correctionOfId?: string | null;
}

interface NormalizedLine extends Required<Omit<LineInput, 'debit' | 'credit'>> { debit: string; credit: string }

const unbalanced = (debit: Decimal, credit: Decimal) =>
  new AppError(400, 'UNBALANCED_ENTRY', `Journal entry is not balanced: debit ${toMoney(debit)}, credit ${toMoney(credit)}`,
    { totalDebit: toMoney(debit), totalCredit: toMoney(credit), difference: toMoney(debit.minus(credit).abs()) });

/** Checks line shape and references. Balance is checked at posting. */
async function normalizeLines(db: Db, companyId: string, lines: LineInput[]): Promise<NormalizedLine[]> {
  if (lines.length < 2) throw badRequest('TOO_FEW_LINES', 'A journal entry needs at least two lines');
  if (lines.length > 1000) throw badRequest('TOO_MANY_LINES', 'A journal entry can have at most 1000 lines');

  const normalized = lines.map((l, i) => {
    const debit = new Decimal(l.debit || '0');
    const credit = new Decimal(l.credit || '0');
    if (debit.isNegative() || credit.isNegative()) throw badRequest('INVALID_AMOUNT', `Line ${i + 1}: amounts cannot be negative`);
    if (debit.decimalPlaces() > 2 || credit.decimalPlaces() > 2) throw badRequest('INVALID_AMOUNT', `Line ${i + 1}: at most 2 decimal places`);
    if (debit.isZero() === credit.isZero()) {
      throw badRequest('INVALID_LINE', `Line ${i + 1}: enter either a debit or a credit amount, not both or neither`);
    }
    return {
      accountId: l.accountId, debit: toMoney(debit), credit: toMoney(credit),
      description: l.description ?? null, costCenterId: l.costCenterId ?? null, branchId: l.branchId ?? null,
    };
  });

  const accountIds = [...new Set(normalized.map((l) => l.accountId))];
  const { rows: accounts } = await db.query<{ id: string; is_postable: boolean; is_active: boolean; code: string }>(
    `SELECT id, is_postable, is_active, code FROM accounts
      WHERE company_id = $1 AND id = ANY($2::uuid[]) AND deleted_at IS NULL`,
    [companyId, accountIds]);
  const byId = new Map(accounts.map((a) => [a.id, a]));
  normalized.forEach((l, i) => {
    const a = byId.get(l.accountId);
    if (!a) throw badRequest('INVALID_ACCOUNT', `Line ${i + 1}: account does not exist in this company`);
    if (!a.is_postable) throw badRequest('ACCOUNT_NOT_POSTABLE', `Line ${i + 1}: account ${a.code} is a header account`);
    if (!a.is_active) throw badRequest('ACCOUNT_INACTIVE', `Line ${i + 1}: account ${a.code} is inactive`);
  });

  const costCenters = [...new Set(normalized.map((l) => l.costCenterId).filter(Boolean))] as string[];
  if (costCenters.length) {
    const { rowCount } = await db.query(
      `SELECT 1 FROM cost_centers WHERE company_id = $1 AND id = ANY($2::uuid[]) AND deleted_at IS NULL AND is_active`,
      [companyId, costCenters]);
    if (rowCount !== costCenters.length) throw badRequest('INVALID_COST_CENTER', 'A cost center does not exist or is inactive');
  }
  const branches = [...new Set(normalized.map((l) => l.branchId).filter(Boolean))] as string[];
  if (branches.length) {
    const { rowCount } = await db.query(
      `SELECT 1 FROM branches WHERE company_id = $1 AND id = ANY($2::uuid[]) AND deleted_at IS NULL`,
      [companyId, branches]);
    if (rowCount !== branches.length) throw badRequest('INVALID_BRANCH', 'A branch does not exist in this company');
  }
  return normalized;
}

async function insertLines(db: Db, tenantId: string, companyId: string, entryId: string, lines: NormalizedLine[]) {
  // One statement for all lines keeps large entries fast.
  await db.query(
    `INSERT INTO journal_entry_lines
       (tenant_id, company_id, journal_entry_id, line_no, account_id, debit, credit, description, cost_center_id, branch_id)
     SELECT $1, $2, $3, l.ord::smallint, l.account_id, l.debit, l.credit, l.description, l.cost_center_id, l.branch_id
       FROM unnest($4::uuid[], $5::numeric[], $6::numeric[], $7::text[], $8::uuid[], $9::uuid[])
            WITH ORDINALITY AS l(account_id, debit, credit, description, cost_center_id, branch_id, ord)`,
    [tenantId, companyId, entryId,
     lines.map((l) => l.accountId), lines.map((l) => l.debit), lines.map((l) => l.credit),
     lines.map((l) => l.description), lines.map((l) => l.costCenterId), lines.map((l) => l.branchId)]);
}

async function assertCompany(db: Db, tenantId: string, companyId: string): Promise<{ currency: string }> {
  const { rows: [c] } = await db.query<{ currency: string }>(
    `SELECT currency FROM companies WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, [tenantId, companyId]);
  if (!c) throw notFound('Company');
  return c;
}

export async function createDraft(db: Db, ctx: EngineContext, input: EntryInput): Promise<string> {
  const company = await assertCompany(db, ctx.tenantId, input.companyId);
  const lines = await normalizeLines(db, input.companyId, input.lines);
  const { rows: [entry] } = await db.query<{ id: string }>(
    `INSERT INTO journal_entries (tenant_id, company_id, entry_date, description, reference_type, reference_id,
                                  source, currency, correction_of_id, created_by)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
    [ctx.tenantId, input.companyId, input.entryDate, input.description, input.referenceType ?? 'MANUAL',
     input.referenceId ?? null, input.source ?? 'MANUAL', company.currency, input.correctionOfId ?? null, ctx.userId]);
  await insertLines(db, ctx.tenantId, input.companyId, entry!.id, lines);
  await writeAudit(db, {
    tenantId: ctx.tenantId, userId: ctx.userId, action: 'CREATE', entityType: 'journal_entry', entityId: entry!.id,
    newValues: { ...input, lines },
  }, ctx.meta);
  return entry!.id;
}

async function lockEntry(db: Db, tenantId: string, id: string) {
  const { rows: [e] } = await db.query<{
    id: string; company_id: string; status: string; entry_date: string; description: string;
    source: string; reference_type: string; deleted_at: Date | null; fiscal_year_id: string | null;
  }>(
    `SELECT id, company_id, status, to_char(entry_date, 'YYYY-MM-DD') AS entry_date, description, source,
            reference_type, deleted_at, fiscal_year_id
       FROM journal_entries WHERE tenant_id = $1 AND id = $2 FOR UPDATE`,
    [tenantId, id]);
  if (!e || e.deleted_at) throw notFound('Journal entry');
  return e;
}

export async function updateDraft(db: Db, ctx: EngineContext, id: string, input: Partial<Omit<EntryInput, 'companyId'>>) {
  const e = await lockEntry(db, ctx.tenantId, id);
  if (e.status !== 'DRAFT') throw conflict('ENTRY_NOT_DRAFT', 'Only draft entries can be edited. Reverse a posted entry instead.');
  if (e.source !== 'MANUAL') throw conflict('SYSTEM_ENTRY', 'System-generated entries cannot be edited here');
  const before = await getEntry(db, ctx.tenantId, id);
  await db.query(
    `UPDATE journal_entries SET entry_date = COALESCE($2, entry_date), description = COALESCE($3, description) WHERE id = $1`,
    [id, input.entryDate ?? null, input.description ?? null]);
  if (input.lines) {
    const lines = await normalizeLines(db, e.company_id, input.lines);
    await db.query(`DELETE FROM journal_entry_lines WHERE journal_entry_id = $1`, [id]);
    await insertLines(db, ctx.tenantId, e.company_id, id, lines);
  }
  const after = await getEntry(db, ctx.tenantId, id);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'UPDATE', entityType: 'journal_entry', entityId: id, oldValues: before, newValues: after }, ctx.meta);
  return after;
}

export async function deleteDraft(db: Db, ctx: EngineContext, id: string) {
  const e = await lockEntry(db, ctx.tenantId, id);
  if (e.status !== 'DRAFT') throw conflict('ENTRY_NOT_DRAFT', 'Posted entries cannot be deleted. Reverse them instead.');
  if (e.source !== 'MANUAL') throw conflict('SYSTEM_ENTRY', 'System-generated entries cannot be deleted here');
  await db.query(`UPDATE journal_entries SET deleted_at = now() WHERE id = $1`, [id]);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'DELETE', entityType: 'journal_entry', entityId: id }, ctx.meta);
}

/**
 * Posts a draft: validates balance, accounts and the fiscal period, assigns
 * the next gap-free number, and makes the entry immutable.
 */
export async function post(db: Db, ctx: EngineContext, id: string) {
  const e = await lockEntry(db, ctx.tenantId, id);
  if (e.status !== 'DRAFT') throw conflict('ENTRY_NOT_DRAFT', 'This entry is already posted');

  // FOR SHARE blocks a concurrent period close until this transaction ends.
  const { rows: [period] } = await db.query<{ id: string; status: string; year_id: string; year_status: string; year_start: string }>(
    `SELECT p.id, p.status, y.id AS year_id, y.status AS year_status, to_char(y.start_date, 'YYYY') AS year_start
       FROM fiscal_periods p JOIN fiscal_years y ON y.id = p.fiscal_year_id
      WHERE p.company_id = $1 AND $2::date BETWEEN p.start_date AND p.end_date
      FOR SHARE OF p, y`,
    [e.company_id, e.entry_date]);
  if (!period) throw badRequest('NO_FISCAL_PERIOD', `No fiscal period covers ${e.entry_date}. Create the fiscal year first.`);
  if (period.status !== 'OPEN' || period.year_status !== 'OPEN') {
    throw conflict('PERIOD_CLOSED', `The fiscal period for ${e.entry_date} is closed`);
  }

  const { rows: lines } = await db.query<{ debit: string; credit: string; is_postable: boolean; is_active: boolean; deleted_at: Date | null; code: string }>(
    `SELECT l.debit, l.credit, a.is_postable, a.is_active, a.deleted_at, a.code
       FROM journal_entry_lines l JOIN accounts a ON a.id = l.account_id
      WHERE l.journal_entry_id = $1`,
    [id]);
  if (lines.length < 2) throw badRequest('TOO_FEW_LINES', 'A journal entry needs at least two lines');
  const bad = lines.find((l) => !l.is_postable || !l.is_active || l.deleted_at);
  if (bad) throw badRequest('ACCOUNT_INACTIVE', `Account ${bad.code} can no longer receive postings`);
  const debit = sum(lines.map((l) => l.debit));
  const credit = sum(lines.map((l) => l.credit));
  if (!debit.equals(credit) || debit.isZero()) throw unbalanced(debit, credit);

  const { rows: [seq] } = await db.query<{ last_number: number }>(
    `UPDATE journal_sequences SET last_number = last_number + 1 WHERE fiscal_year_id = $1 RETURNING last_number`,
    [period.year_id]);
  if (!seq) throw new AppError(500, 'SEQUENCE_MISSING', 'Journal sequence for the fiscal year is missing');
  const number = `JV-${period.year_start}-${String(seq.last_number).padStart(6, '0')}`;

  await db.query(
    `UPDATE journal_entries
        SET status = 'POSTED', entry_number = $2, fiscal_year_id = $3, fiscal_period_id = $4,
            posted_by = $5, posted_at = now()
      WHERE id = $1`,
    [id, number, period.year_id, period.id, ctx.userId]);
  await writeAudit(db, {
    tenantId: ctx.tenantId, userId: ctx.userId, action: 'POST', entityType: 'journal_entry', entityId: id,
    newValues: { entryNumber: number, totalDebit: toMoney(debit), totalCredit: toMoney(credit) },
  }, ctx.meta);
  return getEntry(db, ctx.tenantId, id);
}

/** Creates and posts an entry in one step. Used by sales, purchases, payments, etc. */
export async function postEntry(db: Db, ctx: EngineContext, input: EntryInput) {
  const id = await createDraft(db, ctx, input);
  return post(db, ctx, id);
}

export interface ReverseOptions {
  date?: string;
  reason: string;
  /** System documents (invoices…) reverse their own entries; manual reversal is limited to manual entries. */
  allowSystem?: boolean;
}

/**
 * Reverses a posted entry with a mirror entry (debits and credits swapped)
 * and marks the original REVERSED. The original stays in the ledger.
 */
export async function reverse(db: Db, ctx: EngineContext, id: string, opts: ReverseOptions) {
  const e = await lockEntry(db, ctx.tenantId, id);
  if (e.status === 'REVERSED') throw conflict('ALREADY_REVERSED', 'This entry has already been reversed');
  if (e.status !== 'POSTED') throw conflict('ENTRY_NOT_POSTED', 'Only posted entries can be reversed. Delete or edit the draft instead.');
  if (e.reference_type === 'REVERSAL') {
    throw conflict('CANNOT_REVERSE_REVERSAL', 'A reversal entry cannot be reversed. Post a new entry instead.');
  }
  if (e.source === 'SYSTEM' && !opts.allowSystem) {
    throw conflict('SYSTEM_ENTRY', 'This entry was generated by a document. Cancel or return the document instead.');
  }

  const { rows: lines } = await db.query<{ account_id: string; debit: string; credit: string; description: string | null; cost_center_id: string | null; branch_id: string | null }>(
    `SELECT account_id, debit, credit, description, cost_center_id, branch_id
       FROM journal_entry_lines WHERE journal_entry_id = $1 ORDER BY line_no`, [id]);

  const reversalId = await createDraft(db, ctx, {
    companyId: e.company_id,
    entryDate: opts.date ?? e.entry_date,
    description: `عكس القيد: ${e.description} — ${opts.reason}`.slice(0, 1000),
    referenceType: 'REVERSAL',
    referenceId: id,
    source: e.source as 'MANUAL' | 'SYSTEM',
    lines: lines.map((l) => ({
      accountId: l.account_id, debit: l.credit, credit: l.debit,
      description: l.description, costCenterId: l.cost_center_id, branchId: l.branch_id,
    })),
  });
  await db.query(`UPDATE journal_entries SET reversal_of_id = $2 WHERE id = $1`, [reversalId, id]);
  const reversal = await post(db, ctx, reversalId);

  await db.query(
    `UPDATE journal_entries SET status = 'REVERSED', reversed_by_entry_id = $2, reversed_at = now(), reversed_by = $3 WHERE id = $1`,
    [id, reversalId, ctx.userId]);
  await writeAudit(db, {
    tenantId: ctx.tenantId, userId: ctx.userId, action: 'REVERSE', entityType: 'journal_entry', entityId: id,
    newValues: { reversalEntryId: reversalId, reversalNumber: reversal.entryNumber, reason: opts.reason },
  }, ctx.meta);
  return reversal;
}

/**
 * Closes a fiscal year: transfers the net of all revenue, cost of sales and
 * expense accounts to retained earnings with a posted closing entry, then
 * closes every period and the year.
 */
export async function closeFiscalYear(db: Db, ctx: EngineContext, yearId: string) {
  const { rows: [year] } = await db.query<{ id: string; company_id: string; status: string; start_date: string; end_date: string }>(
    `SELECT id, company_id, status, to_char(start_date, 'YYYY-MM-DD') AS start_date, to_char(end_date, 'YYYY-MM-DD') AS end_date
       FROM fiscal_years WHERE tenant_id = $1 AND id = $2 FOR UPDATE`, [ctx.tenantId, yearId]);
  if (!year) throw notFound('Fiscal year');
  if (year.status === 'CLOSED') throw conflict('YEAR_CLOSED', 'This fiscal year is already closed');

  const earlierOpen = await db.query(
    `SELECT 1 FROM fiscal_years WHERE company_id = $1 AND end_date < $2 AND status = 'OPEN' LIMIT 1`, [year.company_id, year.start_date]);
  if (earlierOpen.rowCount) throw conflict('EARLIER_YEAR_OPEN', 'Close earlier fiscal years first');

  const drafts = await db.query(
    `SELECT 1 FROM journal_entries WHERE company_id = $1 AND status = 'DRAFT' AND deleted_at IS NULL
        AND entry_date BETWEEN $2 AND $3 LIMIT 1`, [year.company_id, year.start_date, year.end_date]);
  if (drafts.rowCount) throw conflict('DRAFTS_EXIST', 'Post or delete the draft entries of this year before closing it');

  const { rows: [re] } = await db.query<{ id: string }>(
    `SELECT id FROM accounts WHERE company_id = $1 AND system_key = 'RETAINED_EARNINGS' AND deleted_at IS NULL AND is_active AND is_postable`,
    [year.company_id]);
  if (!re) throw conflict('RETAINED_EARNINGS_MISSING', 'An active retained earnings account (system key RETAINED_EARNINGS) is required');

  const { rows: balances } = await db.query<{ account_id: string; net: string }>(
    `SELECT l.account_id, sum(l.debit - l.credit) AS net
       FROM journal_entry_lines l
       JOIN journal_entries e ON e.id = l.journal_entry_id
       JOIN accounts a ON a.id = l.account_id
      WHERE e.company_id = $1 AND e.status IN ('POSTED', 'REVERSED')
        AND e.entry_date BETWEEN $2 AND $3
        AND a.account_type = ANY($4::text[])
      GROUP BY l.account_id
     HAVING sum(l.debit - l.credit) <> 0`,
    [year.company_id, year.start_date, year.end_date, PROFIT_AND_LOSS_TYPES]);

  let closingEntryId: string | null = null;
  if (balances.length) {
    const lines: LineInput[] = balances.map((b) => {
      const net = new Decimal(b.net);
      return net.isPositive() ? { accountId: b.account_id, credit: toMoney(net) } : { accountId: b.account_id, debit: toMoney(net.abs()) };
    });
    // Positive total = expenses exceed revenue (a loss) → debit retained earnings.
    const result = sum(balances.map((b) => b.net));
    if (!result.isZero()) {
      lines.push(result.isPositive() ? { accountId: re.id, debit: toMoney(result) } : { accountId: re.id, credit: toMoney(result.abs()) });
    }
    const closing = await postEntry(db, ctx, {
      companyId: year.company_id, entryDate: year.end_date, description: 'قيد إقفال السنة المالية',
      referenceType: 'YEAR_CLOSING', referenceId: year.id, source: 'SYSTEM', lines,
    });
    closingEntryId = closing.id;
  }

  await db.query(
    `UPDATE fiscal_periods SET status = 'CLOSED', closed_at = COALESCE(closed_at, now()), closed_by = COALESCE(closed_by, $2)
      WHERE fiscal_year_id = $1 AND status = 'OPEN'`, [year.id, ctx.userId]);
  await db.query(
    `UPDATE fiscal_years SET status = 'CLOSED', closing_entry_id = $2, closed_at = now(), closed_by = $3 WHERE id = $1`,
    [year.id, closingEntryId, ctx.userId]);
  await writeAudit(db, {
    tenantId: ctx.tenantId, userId: ctx.userId, action: 'CLOSE', entityType: 'fiscal_year', entityId: year.id,
    newValues: { closingEntryId },
  }, ctx.meta);
  return { closingEntryId };
}

export interface EntryView {
  id: string;
  companyId: string;
  entryNumber: string | null;
  entryDate: string;
  description: string;
  referenceType: string;
  referenceId: string | null;
  source: string;
  status: string;
  currency: string;
  totalDebit: string;
  totalCredit: string;
  fiscalPeriodId: string | null;
  reversalOfId: string | null;
  reversedByEntryId: string | null;
  correctionOfId: string | null;
  createdBy: string | null;
  postedAt: Date | null;
  postedBy: string | null;
  reversedAt: Date | null;
  createdAt: Date;
  lines: {
    id: string; lineNo: number; accountId: string; accountCode: string; accountName: string;
    debit: string; credit: string; description: string | null; costCenterId: string | null; branchId: string | null;
  }[];
}

export async function getEntry(db: Db, tenantId: string, id: string): Promise<EntryView> {
  const { rows: [e] } = await db.query<Omit<EntryView, 'lines'>>(
    `SELECT e.id, e.company_id AS "companyId", e.entry_number AS "entryNumber",
            to_char(e.entry_date, 'YYYY-MM-DD') AS "entryDate", e.description,
            e.reference_type AS "referenceType", e.reference_id AS "referenceId", e.source, e.status, e.currency,
            CASE WHEN e.status = 'DRAFT' THEN COALESCE((SELECT sum(debit) FROM journal_entry_lines WHERE journal_entry_id = e.id), 0) ELSE e.total_debit END::numeric(18,2)::text AS "totalDebit",
            CASE WHEN e.status = 'DRAFT' THEN COALESCE((SELECT sum(credit) FROM journal_entry_lines WHERE journal_entry_id = e.id), 0) ELSE e.total_credit END::numeric(18,2)::text AS "totalCredit",
            e.fiscal_period_id AS "fiscalPeriodId", e.reversal_of_id AS "reversalOfId",
            e.reversed_by_entry_id AS "reversedByEntryId", e.correction_of_id AS "correctionOfId",
            e.created_by AS "createdBy", e.posted_at AS "postedAt", e.posted_by AS "postedBy",
            e.reversed_at AS "reversedAt", e.created_at AS "createdAt"
       FROM journal_entries e WHERE e.tenant_id = $1 AND e.id = $2 AND e.deleted_at IS NULL`,
    [tenantId, id]);
  if (!e) throw notFound('Journal entry');
  const { rows: lines } = await db.query<EntryView['lines'][number]>(
    `SELECT l.id, l.line_no AS "lineNo", l.account_id AS "accountId", a.code AS "accountCode", a.name_ar AS "accountName",
            l.debit::text AS debit, l.credit::text AS credit, l.description,
            l.cost_center_id AS "costCenterId", l.branch_id AS "branchId"
       FROM journal_entry_lines l JOIN accounts a ON a.id = l.account_id
      WHERE l.journal_entry_id = $1 ORDER BY l.line_no`, [id]);
  return { ...e, lines };
}
