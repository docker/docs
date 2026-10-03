import type { Db } from '../../db/tx.js';

export type AuditAction =
  | 'REGISTER' | 'LOGIN' | 'LOGIN_FAILED' | 'LOGOUT' | 'TOKEN_REUSE_DETECTED'
  | 'TENANT_SWITCH' | 'PASSWORD_CHANGE'
  | 'CREATE' | 'UPDATE' | 'DELETE'
  | 'SETTINGS_CHANGE' | 'PERMISSION_CHANGE'
  | 'POST' | 'REVERSE' | 'CLOSE' | 'REOPEN';

export interface AuditMeta {
  ip?: string | null;
  userAgent?: string | null;
  requestId?: string | null;
}

export interface AuditEntry {
  tenantId: string | null;
  userId: string | null;
  action: AuditAction;
  entityType?: string | null;
  entityId?: string | null;
  oldValues?: unknown;
  newValues?: unknown;
}

const SECRET_KEYS = new Set(['password', 'password_hash', 'passwordHash', 'token', 'token_hash']);

function redact(value: unknown): unknown {
  if (value === null || value === undefined) return null;
  if (Array.isArray(value)) return value.map(redact);
  if (typeof value === 'object') {
    return Object.fromEntries(
      Object.entries(value as Record<string, unknown>)
        .filter(([k]) => !SECRET_KEYS.has(k))
        .map(([k, v]) => [k, redact(v)]),
    );
  }
  return value;
}

/**
 * Writes an audit record inside the caller's transaction, so the audit row
 * commits or rolls back together with the change it describes.
 */
export async function writeAudit(db: Db, entry: AuditEntry, meta: AuditMeta = {}): Promise<void> {
  await db.query(
    `INSERT INTO audit_logs (tenant_id, user_id, action, entity_type, entity_id,
                             old_values, new_values, ip_address, user_agent, request_id)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
    [
      entry.tenantId, entry.userId, entry.action, entry.entityType ?? null, entry.entityId ?? null,
      entry.oldValues === undefined ? null : JSON.stringify(redact(entry.oldValues)),
      entry.newValues === undefined ? null : JSON.stringify(redact(entry.newValues)),
      meta.ip ?? null, meta.userAgent?.slice(0, 500) ?? null, meta.requestId ?? null,
    ],
  );
}
