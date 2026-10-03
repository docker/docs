import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { parse } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';

const query = z.object({
  action: z.string().regex(/^[A-Z_]+$/).optional(),
  entityType: z.string().max(50).optional(),
  entityId: z.uuid().optional(),
  userId: z.uuid().optional(),
  from: z.iso.datetime({ offset: true }).optional(),
  to: z.iso.datetime({ offset: true }).optional(),
  // Keyset pagination: pass the previous page's nextCursor.
  cursor: z.string().regex(/^[^|]+\|[0-9a-f-]{36}$/).optional(),
  limit: z.coerce.number().int().min(1).max(200).default(50),
});

export default async function auditRoutes(app: FastifyInstance) {
  app.get('/audit-logs', { preHandler: requirePermission(app, 'audit.view') }, async (req) => {
    const q = parse(query, req.query);
    const where = ['a.tenant_id = $1'];
    const params: unknown[] = [req.auth!.tenantId];
    const add = (sql: string, value: unknown) => { params.push(value); where.push(sql.replace('?', `$${params.length}`)); };
    if (q.action) add('a.action = ?', q.action);
    if (q.entityType) add('a.entity_type = ?', q.entityType);
    if (q.entityId) add('a.entity_id = ?', q.entityId);
    if (q.userId) add('a.user_id = ?', q.userId);
    if (q.from) add('a.created_at >= ?', q.from);
    if (q.to) add('a.created_at <= ?', q.to);
    if (q.cursor) {
      const [ts, id] = q.cursor.split('|');
      params.push(ts, id);
      where.push(`(a.created_at, a.id) < ($${params.length - 1}::timestamptz, $${params.length}::uuid)`);
    }
    params.push(q.limit + 1);
    return req.tenantTx(async (db) => {
      const { rows } = await db.query<{ id: string; createdAt: Date }>(
        `SELECT a.id, a.action, a.entity_type AS "entityType", a.entity_id AS "entityId",
                a.user_id AS "userId", u.full_name AS "userName", a.old_values AS "oldValues",
                a.new_values AS "newValues", host(a.ip_address) AS "ipAddress", a.user_agent AS "userAgent",
                a.created_at AS "createdAt"
           FROM audit_logs a LEFT JOIN users u ON u.id = a.user_id
          WHERE ${where.join(' AND ')}
          ORDER BY a.created_at DESC, a.id DESC
          LIMIT $${params.length}`,
        params);
      const page = rows.slice(0, q.limit);
      const last = page[page.length - 1];
      return {
        data: page,
        nextCursor: rows.length > q.limit && last ? `${last.createdAt.toISOString()}|${last.id}` : null,
      };
    });
  });
}
