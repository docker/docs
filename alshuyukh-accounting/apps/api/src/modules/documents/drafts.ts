import type { Db } from '../../db/tx.js';
import { AppError, badRequest, conflict, notFound } from '../../lib/errors.js';
import { Decimal, toMoney } from '../../lib/money.js';
import { nextDocumentNumber } from '../../lib/sequences.js';
import { addDays } from '../../lib/dates.js';
import { writeAudit, type AuditMeta } from '../audit/audit.service.js';
import { CalcError, calculateLine, totalsOf, type CalcLine, type CalcTotals, type VatCategory } from './calc.js';
import { KINDS, partyColumn, partyTable, type DocKind } from './kinds.js';

export interface Ctx { tenantId: string; userId: string; meta?: AuditMeta }

export interface LineInput {
  productId?: string | null;
  accountId?: string | null;
  description?: string | null;
  quantity: string;
  unitPrice?: string;
  discountAmount?: string;
  discountPercent?: string;
  unitId?: string | null;
  vatCategory?: VatCategory;
}

export interface DraftInput {
  companyId: string;
  partyId: string;
  docDate: string;
  dueDate?: string | null;
  validUntil?: string | null;
  expectedDate?: string | null;
  supplierInvoiceNumber?: string | null;
  branchId?: string | null;
  warehouseId?: string | null;
  pricesIncludeVat?: boolean;
  notes?: string | null;
  sourceId?: string | null;
  lines: LineInput[];
}

export interface ReturnInput {
  originalInvoiceId: string;
  docDate: string;
  reason: string;
  notes?: string | null;
  lines: { sourceItemId: string; quantity: string }[];
}

export interface BuiltLine extends CalcLine {
  productId: string | null;
  accountId: string | null;
  sourceItemId: string | null;
  description: string | null;
  unitId: string | null;
  discountBasis: string;
}

/** The standard VAT rate effective on `date`. Z, E and O are always 0. */
export async function vatRateFor(db: Db, companyId: string, category: VatCategory, date: string): Promise<string> {
  if (category !== 'S') return '0';
  const { rows: [r] } = await db.query<{ rate: string }>(
    `SELECT rate::text FROM tax_rates
      WHERE company_id = $1 AND vat_category = 'S' AND is_active
        AND effective_from <= $2 AND (effective_to IS NULL OR effective_to >= $2)`,
    [companyId, date]);
  if (!r) throw conflict('TAX_RATE_MISSING', `No standard VAT rate is configured for ${date}`);
  return r.rate;
}

const calcError = (e: unknown) => {
  if (e instanceof CalcError) return badRequest('INVALID_LINE', e.message);
  return e;
};

/** Turns line input into calculated lines, filling defaults from the product. */
export async function buildLines(db: Db, kind: DocKind, companyId: string, input: Pick<DraftInput, 'lines' | 'docDate' | 'pricesIncludeVat'>): Promise<BuiltLine[]> {
  if (!input.lines.length) throw badRequest('NO_LINES', 'A document needs at least one line');
  const includeVat = input.pricesIncludeVat ?? false;
  const productIds = [...new Set(input.lines.map((l) => l.productId).filter(Boolean))] as string[];
  const { rows: products } = await db.query<{
    id: string; name_ar: string; unit_id: string; sale_price: string; sale_price_includes_vat: boolean;
    purchase_price: string; vat_category: VatCategory; is_active: boolean;
  }>(
    `SELECT id, name_ar, unit_id, sale_price::text, sale_price_includes_vat, purchase_price::text, vat_category, is_active
       FROM products WHERE company_id = $1 AND id = ANY($2::uuid[]) AND deleted_at IS NULL`,
    [companyId, productIds]);
  const byId = new Map(products.map((p) => [p.id, p]));

  const accountIds = [...new Set(input.lines.map((l) => l.accountId).filter(Boolean))] as string[];
  if (accountIds.length) {
    if (kind.party === 'customer') throw badRequest('INVALID_LINE', 'Sales lines must reference a product');
    const { rows } = await db.query<{ id: string }>(
      `SELECT id FROM accounts WHERE company_id = $1 AND id = ANY($2::uuid[]) AND deleted_at IS NULL AND is_active AND is_postable
          AND account_type IN ('ASSET', 'EXPENSE', 'COST_OF_GOODS_SOLD')`, [companyId, accountIds]);
    if (rows.length !== accountIds.length) throw badRequest('INVALID_ACCOUNT', 'Purchase lines need an active postable asset, expense or cost account');
  }
  const unitIds = [...new Set(input.lines.map((l) => l.unitId).filter(Boolean))] as string[];
  if (unitIds.length) {
    const { rowCount } = await db.query(`SELECT 1 FROM units WHERE company_id = $1 AND id = ANY($2::uuid[])`, [companyId, unitIds]);
    if (rowCount !== unitIds.length) throw badRequest('INVALID_UNIT', 'Unit does not exist in this company');
  }

  const built: BuiltLine[] = [];
  for (const [i, l] of input.lines.entries()) {
    const n = i + 1;
    const product = l.productId ? byId.get(l.productId) : undefined;
    if (l.productId && !product) throw badRequest('INVALID_PRODUCT', `Line ${n}: product does not exist in this company`);
    if (product && !product.is_active) throw badRequest('PRODUCT_INACTIVE', `Line ${n}: product is inactive`);
    if (!product && !l.accountId) throw badRequest('INVALID_LINE', `Line ${n}: choose a product${kind.party === 'supplier' ? ' or an account' : ''}`);

    const category: VatCategory = product?.vat_category ?? l.vatCategory ?? 'S';
    const rate = await vatRateFor(db, companyId, category, input.docDate);
    let unitPrice = l.unitPrice;
    if (unitPrice === undefined) {
      if (!product) throw badRequest('INVALID_LINE', `Line ${n}: unitPrice is required`);
      if (kind.party === 'customer') {
        // Convert the product's list price to the document's price basis.
        let p = new Decimal(product.sale_price);
        if (product.sale_price_includes_vat && !includeVat) p = p.dividedBy(new Decimal(rate).plus(1));
        if (!product.sale_price_includes_vat && includeVat) p = p.times(new Decimal(rate).plus(1));
        unitPrice = p.toDecimalPlaces(4).toString();
      } else {
        unitPrice = new Decimal(product.purchase_price).toString();
      }
    }
    try {
      const calc = calculateLine({
        quantity: l.quantity, unitPrice, discountAmount: l.discountAmount, discountPercent: l.discountPercent,
        vatCategory: category, vatRate: rate,
      }, includeVat, n);
      const discountBasis = l.discountPercent
        ? toMoney(new Decimal(l.quantity).times(unitPrice).toDecimalPlaces(2).times(l.discountPercent).dividedBy(100))
        : toMoney(l.discountAmount ?? '0');
      built.push({
        ...calc,
        productId: product?.id ?? null,
        accountId: l.accountId ?? null,
        sourceItemId: null,
        description: l.description ?? product?.name_ar ?? null,
        unitId: l.unitId ?? product?.unit_id ?? null,
        discountBasis,
      });
    } catch (e) {
      throw calcError(e);
    }
  }
  return built;
}

interface OriginalItem {
  id: string; line_no: number; product_id: string | null; account_id: string | null; description: string | null;
  quantity: string; unit_id: string | null; unit_price: string; discount_basis: string; gross_amount: string;
  discount_amount: string; net_amount: string; vat_category: VatCategory; vat_rate: string;
  returned_qty: string; returned_net: string; returned_discount: string; returned_basis: string;
}

/**
 * Computes a return from the original invoice lines. Amounts are prorated by
 * quantity; returning the last remaining quantity of a line takes exactly
 * what is left, and the last return of an invoice takes exactly the VAT that
 * is left, so a full return always nets the invoice to zero.
 */
export async function buildReturn(db: Db, kind: DocKind, tenantId: string, input: ReturnInput, excludeReturnId: string | null) {
  const orig = KINDS[kind.originalKind!];
  const { rows: [invoice] } = await db.query<{
    id: string; company_id: string; status: string; party_id: string; tax_amount: string; prices_include_vat: boolean;
    branch_id: string | null; warehouse_id: string | null; doc_date: string;
  }>(
    `SELECT id, company_id, status, ${partyColumn(kind)} AS party_id, tax_amount::text, prices_include_vat, branch_id, warehouse_id,
            to_char(doc_date, 'YYYY-MM-DD') AS doc_date
       FROM ${orig.table} WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`,
    [tenantId, input.originalInvoiceId]);
  if (!invoice) throw notFound(orig.label);
  if (invoice.status === 'DRAFT' || invoice.status === 'CANCELLED') throw conflict('INVOICE_NOT_ISSUED', 'Only issued invoices can be returned');
  if (input.docDate < invoice.doc_date) throw badRequest('INVALID_DATE', 'A return cannot be dated before its invoice');

  const { rows: items } = await db.query<OriginalItem>(
    `SELECT i.id, i.line_no, i.product_id, i.account_id, i.description, i.quantity::text, i.unit_id, i.unit_price::text,
            i.discount_basis::text, i.gross_amount::text, i.discount_amount::text, i.net_amount::text, i.vat_category, i.vat_rate::text,
            COALESCE(r.qty, 0)::text AS returned_qty, COALESCE(r.net, 0)::text AS returned_net,
            COALESCE(r.discount, 0)::text AS returned_discount, COALESCE(r.basis, 0)::text AS returned_basis
       FROM ${orig.itemsTable} i
       LEFT JOIN LATERAL (
         SELECT sum(ri.quantity) AS qty, sum(ri.net_amount) AS net, sum(ri.discount_amount) AS discount, sum(ri.discount_basis) AS basis
           FROM ${kind.itemsTable} ri JOIN ${kind.table} rd ON rd.id = ri.document_id
          WHERE ri.source_item_id = i.id AND rd.status = 'ISSUED' AND rd.id IS DISTINCT FROM $2
       ) r ON true
      WHERE i.document_id = $1 ORDER BY i.line_no`,
    [invoice.id, excludeReturnId]);
  const byId = new Map(items.map((i) => [i.id, i]));
  const { rows: [prior] } = await db.query<{ tax: string }>(
    `SELECT COALESCE(sum(tax_amount), 0)::text AS tax FROM ${kind.table}
      WHERE original_invoice_id = $1 AND status = 'ISSUED' AND id IS DISTINCT FROM $2`, [invoice.id, excludeReturnId]);

  if (!input.lines.length) throw badRequest('NO_LINES', 'Choose at least one line to return');
  const seen = new Set<string>();
  const lines: BuiltLine[] = input.lines.map((l, idx) => {
    const n = idx + 1;
    const o = byId.get(l.sourceItemId);
    if (!o) throw badRequest('INVALID_LINE', `Line ${n}: the line does not belong to the original invoice`);
    if (seen.has(o.id)) throw badRequest('INVALID_LINE', `Line ${n}: the same invoice line appears twice`);
    seen.add(o.id);
    const q = new Decimal(l.quantity);
    const remaining = new Decimal(o.quantity).minus(o.returned_qty);
    if (!q.greaterThan(0)) throw badRequest('INVALID_LINE', `Line ${n}: quantity must be greater than zero`);
    if (q.greaterThan(remaining)) {
      throw badRequest('RETURN_EXCEEDS_INVOICE', `Line ${n}: only ${remaining.toString()} can still be returned`);
    }
    let net: Decimal; let discount: Decimal; let basis: Decimal;
    if (q.equals(remaining)) {
      net = new Decimal(o.net_amount).minus(o.returned_net);
      discount = new Decimal(o.discount_amount).minus(o.returned_discount);
      basis = new Decimal(o.discount_basis).minus(o.returned_basis);
    } else {
      const share = q.dividedBy(o.quantity);
      net = new Decimal(o.net_amount).times(share).toDecimalPlaces(2);
      discount = new Decimal(o.discount_amount).times(share).toDecimalPlaces(2);
      basis = new Decimal(o.discount_basis).times(share).toDecimalPlaces(2);
    }
    const vat = net.times(o.vat_rate).toDecimalPlaces(2);
    return {
      quantity: q.toString(), unitPrice: o.unit_price, grossAmount: toMoney(net.plus(discount)), discountAmount: toMoney(discount),
      netAmount: toMoney(net), vatCategory: o.vat_category, vatRate: o.vat_rate, vatAmount: toMoney(vat), totalAmount: toMoney(net.plus(vat)),
      productId: o.product_id, accountId: o.account_id, sourceItemId: o.id, description: o.description, unitId: o.unit_id,
      discountBasis: toMoney(basis),
    };
  });

  const totals: CalcTotals = totalsOf(lines);
  const taxLeft = new Decimal(invoice.tax_amount).minus(prior!.tax);
  const fullyReturned = items.every((o) => {
    const line = lines.find((l) => l.sourceItemId === o.id);
    return new Decimal(o.returned_qty).plus(line?.quantity ?? 0).equals(o.quantity);
  });
  let tax = new Decimal(totals.taxAmount);
  if (fullyReturned || tax.greaterThan(taxLeft)) tax = taxLeft;
  totals.taxAmount = toMoney(tax);
  totals.total = toMoney(new Decimal(totals.taxableAmount).plus(tax));
  return { invoice, lines, totals, fullyReturned };
}

async function checkParty(db: Db, kind: DocKind, companyId: string, partyId: string) {
  const { rows: [p] } = await db.query<{ is_active: boolean }>(
    `SELECT is_active FROM ${partyTable(kind)} WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, partyId]);
  if (!p) throw badRequest('INVALID_PARTY', `${kind.party === 'customer' ? 'Customer' : 'Supplier'} does not exist in this company`);
  if (!p.is_active) throw badRequest('PARTY_INACTIVE', `${kind.party === 'customer' ? 'Customer' : 'Supplier'} is inactive`);
}

async function checkBranchWarehouse(db: Db, companyId: string, branchId?: string | null, warehouseId?: string | null) {
  if (branchId) {
    const b = await db.query(`SELECT 1 FROM branches WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, branchId]);
    if (!b.rowCount) throw badRequest('INVALID_BRANCH', 'Branch does not exist in this company');
  }
  if (warehouseId) {
    const w = await db.query(`SELECT 1 FROM warehouses WHERE company_id = $1 AND id = $2 AND deleted_at IS NULL`, [companyId, warehouseId]);
    if (!w.rowCount) throw badRequest('INVALID_WAREHOUSE', 'Warehouse does not exist in this company');
  }
}

export async function insertItems(db: Db, kind: DocKind, tenantId: string, companyId: string, documentId: string, lines: BuiltLine[]): Promise<string[]> {
  const sourceCol = kind.isReturn ? ', source_item_id' : '';
  const ids: string[] = [];
  for (const [i, l] of lines.entries()) {
    const { rows: [row] } = await db.query<{ id: string }>(
      `INSERT INTO ${kind.itemsTable} (tenant_id, company_id, document_id, line_no, product_id, account_id, description,
         quantity, unit_id, unit_price, gross_amount, discount_basis, discount_amount, net_amount, vat_category, vat_rate,
         vat_amount, total_amount${sourceCol})
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18${kind.isReturn ? ', $19' : ''})
       RETURNING id`,
      [tenantId, companyId, documentId, i + 1, l.productId, l.accountId, l.description, l.quantity, l.unitId, l.unitPrice,
       l.grossAmount, l.discountBasis, l.discountAmount, l.netAmount, l.vatCategory, l.vatRate, l.vatAmount, l.totalAmount,
       ...(kind.isReturn ? [l.sourceItemId] : [])]);
    ids.push(row!.id);
  }
  return ids;
}

const totalsParams = (t: CalcTotals) => [t.subtotal, t.discountTotal, t.taxableAmount, t.taxAmount, t.total];

/** Creates a draft (or, for quotes and orders, a numbered open document). */
export async function createDraft(db: Db, kind: DocKind, ctx: Ctx, input: DraftInput): Promise<string> {
  if (kind.isReturn) throw new AppError(500, 'INTERNAL', 'use createReturnDraft');
  await checkParty(db, kind, input.companyId, input.partyId);
  await checkBranchWarehouse(db, input.companyId, input.branchId, input.warehouseId);
  const lines = await buildLines(db, kind, input.companyId, input);
  const totals = totalsOf(lines);
  const { rows: [company] } = await db.query<{ currency: string }>(`SELECT currency FROM companies WHERE id = $1`, [input.companyId]);
  const number = kind.legal ? null : await nextDocumentNumber(db, ctx.tenantId, input.companyId, kind.key, { prefix: kind.prefix, padding: 6 });

  const extra: Record<string, unknown> = {};
  if (kind.key === 'SALES_QUOTE') extra.valid_until = input.validUntil ?? null;
  if (kind.key === 'PURCHASE_ORDER') extra.expected_date = input.expectedDate ?? null;
  if (kind.key === 'SALES_INVOICE' || kind.key === 'PURCHASE_INVOICE') extra.due_date = input.dueDate ?? null;
  if (kind.key === 'PURCHASE_INVOICE') extra.supplier_invoice_number = input.supplierInvoiceNumber ?? null;
  if (kind.key === 'SALES_INVOICE') extra.source_quote_id = input.sourceId ?? null;
  if (kind.key === 'PURCHASE_INVOICE') extra.source_order_id = input.sourceId ?? null;
  const extraCols = Object.keys(extra);

  const { rows: [doc] } = await db.query<{ id: string }>(
    `INSERT INTO ${kind.table} (tenant_id, company_id, ${partyColumn(kind)}, doc_number, doc_date, branch_id, warehouse_id,
       currency, prices_include_vat, notes, subtotal, discount_total, taxable_amount, tax_amount, total, created_by
       ${extraCols.map((c) => `, ${c}`).join('')})
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
       ${extraCols.map((_, i) => `, $${17 + i}`).join('')})
     RETURNING id`,
    [ctx.tenantId, input.companyId, input.partyId, number, input.docDate, input.branchId ?? null, input.warehouseId ?? null,
     company!.currency, input.pricesIncludeVat ?? false, input.notes ?? null, ...totalsParams(totals), ctx.userId,
     ...Object.values(extra)]);
  await insertItems(db, kind, ctx.tenantId, input.companyId, doc!.id, lines);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'CREATE', entityType: kind.entity, entityId: doc!.id, newValues: { ...input, totals } }, ctx.meta);
  return doc!.id;
}

export async function createReturnDraft(db: Db, kind: DocKind, ctx: Ctx, input: ReturnInput): Promise<string> {
  const { invoice, lines, totals } = await buildReturn(db, kind, ctx.tenantId, input, null);
  const { rows: [company] } = await db.query<{ currency: string }>(`SELECT currency FROM companies WHERE id = $1`, [invoice.company_id]);
  const { rows: [doc] } = await db.query<{ id: string }>(
    `INSERT INTO ${kind.table} (tenant_id, company_id, ${partyColumn(kind)}, original_invoice_id, reason, doc_date, branch_id,
       warehouse_id, currency, prices_include_vat, notes, subtotal, discount_total, taxable_amount, tax_amount, total, created_by)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id`,
    [ctx.tenantId, invoice.company_id, invoice.party_id, invoice.id, input.reason, input.docDate, invoice.branch_id,
     invoice.warehouse_id, company!.currency, invoice.prices_include_vat, input.notes ?? null, ...totalsParams(totals), ctx.userId]);
  await insertItems(db, kind, ctx.tenantId, invoice.company_id, doc!.id, lines);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'CREATE', entityType: kind.entity, entityId: doc!.id, newValues: { ...input, totals } }, ctx.meta);
  return doc!.id;
}

export async function lockDocument(db: Db, kind: DocKind, tenantId: string, id: string) {
  const { rows: [d] } = await db.query<{
    id: string; company_id: string; status: string; doc_number: string | null; doc_date: string; party_id: string;
    total: string; tax_amount: string; taxable_amount: string; branch_id: string | null; journal_entry_id: string | null;
    deleted_at: Date | null; [k: string]: unknown;
  }>(
    `SELECT *, to_char(doc_date, 'YYYY-MM-DD') AS doc_date, ${partyColumn(kind)} AS party_id,
            total::text, tax_amount::text, taxable_amount::text
       FROM ${kind.table} WHERE tenant_id = $1 AND id = $2 FOR UPDATE`, [tenantId, id]);
  if (!d || d.deleted_at) throw notFound(kind.label);
  return d;
}

/** Replaces a draft's header fields and lines. Quotes and orders stay editable until sent/approved. */
export type UpdateInput = Partial<Omit<DraftInput, 'lines'>> & Partial<Omit<ReturnInput, 'lines'>> & { lines?: LineInput[] | ReturnInput['lines'] };

export async function updateDraft(db: Db, kind: DocKind, ctx: Ctx, id: string, input: UpdateInput) {
  const d = await lockDocument(db, kind, ctx.tenantId, id);
  if (d.status !== 'DRAFT') throw conflict('DOCUMENT_NOT_DRAFT', 'Only drafts can be edited');
  const before = await loadDocument(db, kind, ctx.tenantId, id);

  let lines: BuiltLine[];
  let totals: CalcTotals;
  if (kind.isReturn) {
    const r = await buildReturn(db, kind, ctx.tenantId, {
      originalInvoiceId: d.original_invoice_id as string,
      docDate: input.docDate ?? d.doc_date,
      reason: input.reason ?? (d.reason as string),
      lines: input.lines as ReturnInput['lines'] ?? before.lines.map((l) => ({ sourceItemId: l.sourceItemId!, quantity: l.quantity })),
    }, id);
    ({ lines, totals } = r);
    await db.query(`UPDATE ${kind.table} SET doc_date = $2, reason = COALESCE($3, reason), notes = COALESCE($4, notes) WHERE id = $1`,
      [id, input.docDate ?? d.doc_date, input.reason ?? null, input.notes ?? null]);
  } else {
    if (input.partyId) await checkParty(db, kind, d.company_id, input.partyId);
    await checkBranchWarehouse(db, d.company_id, input.branchId, input.warehouseId);
    const pricesIncludeVat = input.pricesIncludeVat ?? (d.prices_include_vat as boolean);
    const docDate = input.docDate ?? d.doc_date;
    const lineInput: LineInput[] = (input.lines as LineInput[] | undefined) ?? before.lines.map((l) => ({
      productId: l.productId, accountId: l.accountId, description: l.description, quantity: l.quantity,
      unitPrice: l.unitPrice, discountAmount: l.discountBasis, unitId: l.unitId, vatCategory: l.vatCategory,
    }));
    lines = await buildLines(db, kind, d.company_id, { lines: lineInput, docDate, pricesIncludeVat });
    totals = totalsOf(lines);
    const sets: [string, unknown][] = [
      [partyColumn(kind), input.partyId ?? d.party_id], ['doc_date', docDate], ['prices_include_vat', pricesIncludeVat],
      ['branch_id', input.branchId !== undefined ? input.branchId : d.branch_id],
      ['warehouse_id', input.warehouseId !== undefined ? input.warehouseId : d.warehouse_id],
      ['notes', input.notes !== undefined ? input.notes : d.notes],
    ];
    if (kind.key === 'SALES_QUOTE' && input.validUntil !== undefined) sets.push(['valid_until', input.validUntil]);
    if (kind.key === 'PURCHASE_ORDER' && input.expectedDate !== undefined) sets.push(['expected_date', input.expectedDate]);
    if ((kind.key === 'SALES_INVOICE' || kind.key === 'PURCHASE_INVOICE') && input.dueDate !== undefined) sets.push(['due_date', input.dueDate]);
    if (kind.key === 'PURCHASE_INVOICE' && input.supplierInvoiceNumber !== undefined) sets.push(['supplier_invoice_number', input.supplierInvoiceNumber]);
    await db.query(`UPDATE ${kind.table} SET ${sets.map(([c], i) => `${c} = $${i + 2}`).join(', ')} WHERE id = $1`,
      [id, ...sets.map(([, v]) => v)]);
  }
  await db.query(`DELETE FROM ${kind.itemsTable} WHERE document_id = $1`, [id]);
  await insertItems(db, kind, ctx.tenantId, d.company_id, id, lines);
  await db.query(
    `UPDATE ${kind.table} SET subtotal = $2, discount_total = $3, taxable_amount = $4, tax_amount = $5, total = $6 WHERE id = $1`,
    [id, ...totalsParams(totals)]);
  const after = await loadDocument(db, kind, ctx.tenantId, id);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'UPDATE', entityType: kind.entity, entityId: id, oldValues: before, newValues: after }, ctx.meta);
  return after;
}

export async function deleteDraft(db: Db, kind: DocKind, ctx: Ctx, id: string) {
  const d = await lockDocument(db, kind, ctx.tenantId, id);
  if (d.status !== 'DRAFT') throw conflict('DOCUMENT_NOT_DRAFT', kind.legal ? 'Issued documents cannot be deleted. Cancel them or issue a return.' : 'Only drafts can be deleted. Cancel it instead.');
  await db.query(`UPDATE ${kind.table} SET deleted_at = now() WHERE id = $1`, [id]);
  await writeAudit(db, { tenantId: ctx.tenantId, userId: ctx.userId, action: 'DELETE', entityType: kind.entity, entityId: id }, ctx.meta);
}

export interface DocumentView {
  id: string; companyId: string; number: string | null; date: string; partyId: string; partyName: string; status: string;
  total: string; [k: string]: unknown;
  lines: {
    id: string; lineNo: number; productId: string | null; productName: string | null; sku: string | null; accountId: string | null;
    sourceItemId: string | null; description: string | null; quantity: string; unitId: string | null; unitCode: string | null;
    unitPrice: string; grossAmount: string; discountBasis: string; discountAmount: string; netAmount: string;
    vatCategory: VatCategory; vatRate: string; vatAmount: string; totalAmount: string;
  }[];
}

export function headerSelect(kind: DocKind) {
  const extra: string[] = [];
  if (kind.key === 'SALES_QUOTE') extra.push(`to_char(d.valid_until, 'YYYY-MM-DD') AS "validUntil"`, `d.converted_invoice_id AS "convertedInvoiceId"`);
  if (kind.key === 'PURCHASE_ORDER') extra.push(`to_char(d.expected_date, 'YYYY-MM-DD') AS "expectedDate"`, `d.converted_invoice_id AS "convertedInvoiceId"`);
  if (kind.key === 'SALES_INVOICE' || kind.key === 'PURCHASE_INVOICE') {
    extra.push(`to_char(d.due_date, 'YYYY-MM-DD') AS "dueDate"`, `d.paid_amount::text AS "paidAmount"`,
      `d.returned_amount::text AS "returnedAmount"`, `d.remaining_amount::text AS "remainingAmount"`);
  }
  if (kind.key === 'SALES_INVOICE') extra.push(`d.invoice_kind AS "invoiceKind"`, `d.source_quote_id AS "sourceQuoteId"`);
  if (kind.key === 'PURCHASE_INVOICE') extra.push(`d.supplier_invoice_number AS "supplierInvoiceNumber"`, `d.source_order_id AS "sourceOrderId"`);
  if (kind.isReturn) {
    extra.push(`d.original_invoice_id AS "originalInvoiceId"`, `d.reason`, `d.applied_amount::text AS "appliedAmount"`,
      `d.refunded_amount::text AS "refundedAmount"`, `d.remaining_amount::text AS "remainingAmount"`,
      `(SELECT doc_number FROM ${KINDS[kind.originalKind!].table} WHERE id = d.original_invoice_id) AS "originalInvoiceNumber"`);
  }
  return `
    SELECT d.id, d.company_id AS "companyId", d.doc_number AS number, to_char(d.doc_date, 'YYYY-MM-DD') AS date,
           d.${partyColumn(kind)} AS "partyId", p.name_ar AS "partyName", p.code AS "partyCode", d.status,
           d.branch_id AS "branchId", d.warehouse_id AS "warehouseId", d.currency, d.prices_include_vat AS "pricesIncludeVat",
           d.subtotal::text, d.discount_total::text AS "discountTotal", d.taxable_amount::text AS "taxableAmount",
           d.tax_amount::text AS "taxAmount", d.total::text, d.notes, d.party_snapshot AS "partySnapshot",
           d.journal_entry_id AS "journalEntryId", d.issued_at AS "issuedAt", d.cancelled_at AS "cancelledAt",
           d.cancel_reason AS "cancelReason", d.created_at AS "createdAt"
           ${extra.map((e) => `, ${e}`).join('')}
      FROM ${kind.table} d JOIN ${partyTable(kind)} p ON p.id = d.${partyColumn(kind)}`;
}

export async function loadDocument(db: Db, kind: DocKind, tenantId: string, id: string): Promise<DocumentView> {
  const { rows: [doc] } = await db.query<Omit<DocumentView, 'lines'>>(
    `${headerSelect(kind)} WHERE d.tenant_id = $1 AND d.id = $2 AND d.deleted_at IS NULL`, [tenantId, id]);
  if (!doc) throw notFound(kind.label);
  const { rows: lines } = await db.query<DocumentView['lines'][number]>(
    `SELECT i.id, i.line_no AS "lineNo", i.product_id AS "productId", pr.name_ar AS "productName", pr.sku,
            i.account_id AS "accountId", ${kind.isReturn ? 'i.source_item_id' : 'NULL::uuid'} AS "sourceItemId",
            i.description, i.quantity::text, i.unit_id AS "unitId", u.code AS "unitCode", i.unit_price::text AS "unitPrice",
            i.gross_amount::text AS "grossAmount", i.discount_basis::text AS "discountBasis", i.discount_amount::text AS "discountAmount",
            i.net_amount::text AS "netAmount", i.vat_category AS "vatCategory", i.vat_rate::text AS "vatRate",
            i.vat_amount::text AS "vatAmount", i.total_amount::text AS "totalAmount"
       FROM ${kind.itemsTable} i
       LEFT JOIN products pr ON pr.id = i.product_id
       LEFT JOIN units u ON u.id = i.unit_id
      WHERE i.document_id = $1 ORDER BY i.line_no`, [id]);
  const view = { ...doc, lines } as DocumentView;
  if (kind.key === 'SALES_INVOICE' || kind.key === 'PURCHASE_INVOICE' || kind.isReturn) {
    const col = { SALES_INVOICE: 'sales_invoice_id', PURCHASE_INVOICE: 'purchase_invoice_id', SALES_RETURN: 'sales_return_id', PURCHASE_RETURN: 'purchase_return_id' }[kind.key as string]!;
    const { rows: allocations } = await db.query(
      `SELECT a.id, a.amount::text, a.created_at AS "createdAt", p.id AS "paymentId", p.payment_number AS "paymentNumber",
              to_char(p.payment_date, 'YYYY-MM-DD') AS "paymentDate"
         FROM payment_allocations a JOIN payments p ON p.id = a.payment_id
        WHERE a.${col} = $1 AND a.reversed_at IS NULL ORDER BY p.payment_date`, [id]);
    view.allocations = allocations;
  }
  if (kind.key === 'SALES_INVOICE' || kind.key === 'PURCHASE_INVOICE') {
    const ret = kind.key === 'SALES_INVOICE' ? KINDS.SALES_RETURN : KINDS.PURCHASE_RETURN;
    const { rows: returns } = await db.query(
      `SELECT id, doc_number AS number, to_char(doc_date, 'YYYY-MM-DD') AS date, status, total::text
         FROM ${ret.table} WHERE original_invoice_id = $1 AND deleted_at IS NULL ORDER BY created_at`, [id]);
    view.returns = returns;
  }
  return view;
}

/** Default due date: document date + the party's payment terms. */
export async function defaultDueDate(db: Db, kind: DocKind, partyId: string, docDate: string): Promise<string> {
  const { rows: [p] } = await db.query<{ payment_terms_days: number }>(`SELECT payment_terms_days FROM ${partyTable(kind)} WHERE id = $1`, [partyId]);
  return addDays(docDate, p?.payment_terms_days ?? 0);
}
