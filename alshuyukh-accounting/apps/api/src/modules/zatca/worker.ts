import type pg from 'pg';
import type { FastifyBaseLogger } from 'fastify';
import { withTx } from '../../db/tx.js';
import { submitInvoice } from './service.js';

/**
 * Background submitter: reports simplified invoices (deadline: 24 hours) and
 * clears standard ones that could not be cleared when issued. Failures stay
 * PENDING with exponential backoff; rejections are final and shown in the UI.
 */
export function startZatcaWorker(pool: pg.Pool, log: FastifyBaseLogger, intervalSeconds: number) {
  let running = false;
  const tick = async () => {
    if (running) return;
    running = true;
    try {
      const { rows } = await pool.query<{ tenant_id: string; id: string }>(`SELECT tenant_id, id FROM zatca_due_invoices(20)`);
      for (const r of rows) {
        const ctx = { tenantId: r.tenant_id, userId: null };
        try {
          await submitInvoice((fn) => withTx(pool, ctx, fn), ctx, r.id, { due: true });
        } catch (e) {
          log.warn({ err: e, zatcaInvoiceId: r.id }, 'zatca submission failed');
        }
      }
    } catch (e) {
      log.error({ err: e }, 'zatca worker tick failed');
    } finally {
      running = false;
    }
  };
  const timer = setInterval(tick, intervalSeconds * 1000);
  timer.unref();
  return () => clearInterval(timer);
}
