import type { Db } from '../db/tx.js';

export interface SequenceDefaults { prefix: string; padding?: number }

/**
 * Returns the next number of a per-company document sequence, e.g.
 * "CUS-00001". The counter row is locked until the caller's transaction ends,
 * so numbers are gap-free and never reused (a rollback also rolls back the
 * increment).
 */
export async function nextDocumentNumber(
  db: Db, tenantId: string, companyId: string, docType: string, defaults: SequenceDefaults,
): Promise<string> {
  await db.query(
    `INSERT INTO document_sequences (tenant_id, company_id, doc_type, prefix, padding)
     VALUES ($1, $2, $3, $4, $5) ON CONFLICT (company_id, doc_type) DO NOTHING`,
    [tenantId, companyId, docType, defaults.prefix, defaults.padding ?? 5]);
  const { rows: [row] } = await db.query<{ prefix: string; padding: number; last_number: number }>(
    `UPDATE document_sequences SET last_number = last_number + 1
      WHERE company_id = $1 AND doc_type = $2 RETURNING prefix, padding, last_number`,
    [companyId, docType]);
  const n = String(row!.last_number).padStart(row!.padding, '0');
  return row!.prefix ? `${row!.prefix}-${n}` : n;
}

/**
 * Generates a code that is not already taken in `table` (a user may have
 * typed a code that matches a future generated one). Master-data codes only;
 * legal document numbers must not skip, so they use nextDocumentNumber alone.
 */
export async function nextFreeCode(
  db: Db, tenantId: string, companyId: string, docType: string, defaults: SequenceDefaults,
  table: 'customers' | 'suppliers' | 'products', column: 'code' | 'sku',
): Promise<string> {
  for (let i = 0; i < 1000; i++) {
    const code = await nextDocumentNumber(db, tenantId, companyId, docType, defaults);
    const taken = await db.query(`SELECT 1 FROM ${table} WHERE company_id = $1 AND ${column} = $2 AND deleted_at IS NULL`, [companyId, code]);
    if (!taken.rowCount) return code;
  }
  throw new Error(`Could not generate a free ${docType} code`);
}
