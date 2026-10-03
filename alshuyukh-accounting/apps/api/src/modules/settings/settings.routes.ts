import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { parse, timezone } from '../../lib/validation.js';
import { requireAuth, requirePermission } from '../../plugins/auth.js';
import { writeAudit } from '../audit/audit.service.js';

const SELECT = `SELECT t.name AS "tenantName", s.default_currency AS "defaultCurrency", s.timezone, s.locale,
  s.fiscal_year_start_month AS "fiscalYearStartMonth", s.date_format AS "dateFormat", s.updated_at AS "updatedAt"
  FROM tenant_settings s JOIN tenants t ON t.id = s.tenant_id WHERE s.tenant_id = $1`;

const updateBody = z.object({
  tenantName: z.string().trim().min(2).max(200),
  defaultCurrency: z.string().regex(/^[A-Z]{3}$/),
  timezone,
  locale: z.enum(['ar', 'en']),
  fiscalYearStartMonth: z.number().int().min(1).max(12),
  dateFormat: z.enum(['YYYY-MM-DD', 'DD/MM/YYYY', 'MM/DD/YYYY']),
}).partial();

export default async function settingsRoutes(app: FastifyInstance) {
  app.get('/settings/tenant', { preHandler: requireAuth(app) }, async (req) =>
    req.tenantTx(async (db) => (await db.query(SELECT, [req.auth!.tenantId])).rows[0]));

  app.patch('/settings/tenant', { preHandler: requirePermission(app, 'settings.manage') }, async (req) => {
    const body = parse(updateBody, req.body);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      const before = (await db.query(SELECT, [a.tenantId])).rows[0];
      if (body.tenantName) await db.query(`UPDATE tenants SET name = $2 WHERE id = $1`, [a.tenantId, body.tenantName]);
      // TODO(Phase 2): block fiscalYearStartMonth changes once a fiscal year exists.
      await db.query(
        `UPDATE tenant_settings SET
           default_currency = COALESCE($2, default_currency), timezone = COALESCE($3, timezone),
           locale = COALESCE($4, locale), fiscal_year_start_month = COALESCE($5, fiscal_year_start_month),
           date_format = COALESCE($6, date_format)
         WHERE tenant_id = $1`,
        [a.tenantId, body.defaultCurrency ?? null, body.timezone ?? null, body.locale ?? null,
         body.fiscalYearStartMonth ?? null, body.dateFormat ?? null]);
      const after = (await db.query(SELECT, [a.tenantId])).rows[0];
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'SETTINGS_CHANGE', entityType: 'tenant_settings', entityId: a.tenantId, oldValues: before, newValues: after }, req.auditMeta());
      return after;
    });
  });
}
