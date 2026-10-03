import type { FastifyInstance, FastifyRequest } from 'fastify';
import { z } from 'zod';
import { isoDate } from '../../lib/dates.js';
import { conflict, notFound } from '../../lib/errors.js';
import { parse, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { writeAudit } from '../audit/audit.service.js';
import { resolveCompanyId } from './company-context.js';
import { closeFiscalYear } from './engine.js';
import { createFiscalYear } from './setup.js';

export default async function fiscalRoutes(app: FastifyInstance) {
  const canView = requirePermission(app, 'account.view');
  const canManage = requirePermission(app, 'fiscal.manage');

  app.get('/fiscal-years', { preHandler: canView }, async (req) => {
    const q = parse(z.object({ companyId: z.uuid().optional() }), req.query);
    return req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, req.auth!.tenantId, q.companyId);
      const { rows } = await db.query(
        `SELECT y.id, y.name, to_char(y.start_date, 'YYYY-MM-DD') AS "startDate", to_char(y.end_date, 'YYYY-MM-DD') AS "endDate",
                y.status, y.closing_entry_id AS "closingEntryId", y.closed_at AS "closedAt",
                COALESCE(json_agg(json_build_object(
                  'id', p.id, 'periodNumber', p.period_number, 'name', p.name,
                  'startDate', to_char(p.start_date, 'YYYY-MM-DD'), 'endDate', to_char(p.end_date, 'YYYY-MM-DD'),
                  'status', p.status) ORDER BY p.period_number), '[]') AS periods
           FROM fiscal_years y LEFT JOIN fiscal_periods p ON p.fiscal_year_id = y.id
          WHERE y.tenant_id = $1 AND y.company_id = $2
          GROUP BY y.id ORDER BY y.start_date DESC`,
        [req.auth!.tenantId, companyId]);
      return { data: rows };
    });
  });

  app.post('/fiscal-years', { preHandler: canManage }, async (req, reply) => {
    const body = parse(z.object({
      companyId: z.uuid().optional(),
      startDate: isoDate,
      endDate: isoDate.optional(),
      name: z.string().trim().min(1).max(100).optional(),
    }), req.body);
    const a = req.auth!;
    const id = await req.tenantTx(async (db) => {
      const companyId = await resolveCompanyId(db, a.tenantId, body.companyId);
      const yearId = await createFiscalYear(db, { tenantId: a.tenantId, companyId, userId: a.userId }, body.startDate, body.endDate, body.name);
      await writeAudit(db, { tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'fiscal_year', entityId: yearId, newValues: body }, req.auditMeta());
      return yearId;
    });
    reply.code(201);
    return { id };
  });

  app.post('/fiscal-years/:id/close', { preHandler: canManage }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const a = req.auth!;
    return req.tenantTx((db) => closeFiscalYear(db, { tenantId: a.tenantId, userId: a.userId, meta: req.auditMeta() }, id));
  });

  const setPeriodStatus = (status: 'OPEN' | 'CLOSED') => async (req: FastifyRequest) => {
    const { id } = parse(uuidParam, req.params);
    const a = req.auth!;
    return req.tenantTx(async (db) => {
      // FOR UPDATE waits for postings that hold the period FOR SHARE.
      const { rows: [p] } = await db.query<{ status: string; year_status: string }>(
        `SELECT p.status, y.status AS year_status FROM fiscal_periods p JOIN fiscal_years y ON y.id = p.fiscal_year_id
          WHERE p.tenant_id = $1 AND p.id = $2 FOR UPDATE OF p`, [a.tenantId, id]);
      if (!p) throw notFound('Fiscal period');
      if (p.year_status === 'CLOSED') throw conflict('YEAR_CLOSED', 'The fiscal year is closed');
      if (p.status === status) throw conflict('NO_CHANGE', `The period is already ${status.toLowerCase()}`);
      await db.query(
        `UPDATE fiscal_periods SET status = $2,
           closed_at = CASE WHEN $2 = 'CLOSED' THEN now() END, closed_by = CASE WHEN $2 = 'CLOSED' THEN $3::uuid END
         WHERE id = $1`, [id, status, a.userId]);
      await writeAudit(db, {
        tenantId: a.tenantId, userId: a.userId, action: status === 'CLOSED' ? 'CLOSE' : 'REOPEN',
        entityType: 'fiscal_period', entityId: id, oldValues: { status: p.status }, newValues: { status },
      }, req.auditMeta());
      return { id, status };
    });
  };
  app.post('/fiscal-periods/:id/close', { preHandler: canManage }, setPeriodStatus('CLOSED'));
  app.post('/fiscal-periods/:id/reopen', { preHandler: canManage }, setPeriodStatus('OPEN'));
}
