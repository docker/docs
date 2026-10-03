import type { Db } from '../../db/tx.js';
import { badRequest, notFound } from '../../lib/errors.js';

/**
 * Resolves the company a request works on. When the tenant has a single
 * company, companyId may be omitted.
 */
export async function resolveCompanyId(db: Db, tenantId: string, companyId?: string | null): Promise<string> {
  if (companyId) {
    const { rowCount } = await db.query(`SELECT 1 FROM companies WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, [tenantId, companyId]);
    if (!rowCount) throw notFound('Company');
    return companyId;
  }
  const { rows } = await db.query<{ id: string }>(`SELECT id FROM companies WHERE tenant_id = $1 AND deleted_at IS NULL LIMIT 2`, [tenantId]);
  if (rows.length === 1) return rows[0]!.id;
  throw badRequest('COMPANY_REQUIRED', 'companyId is required because this organization has more than one company');
}
