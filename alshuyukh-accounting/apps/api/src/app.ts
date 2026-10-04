import cookie from '@fastify/cookie';
import cors from '@fastify/cors';
import helmet from '@fastify/helmet';
import rateLimit from '@fastify/rate-limit';
import Fastify, { type FastifyInstance } from 'fastify';
import type pg from 'pg';
import type { Env } from './config/env.js';
import { AppError } from './lib/errors.js';
import { TokenService } from './lib/tokens.js';
import accountsRoutes from './modules/accounting/accounts.routes.js';
import fiscalRoutes from './modules/accounting/fiscal.routes.js';
import journalRoutes from './modules/accounting/journal.routes.js';
import trialBalanceRoutes from './modules/accounting/trial-balance.routes.js';
import auditRoutes from './modules/audit/audit.routes.js';
import { CUSTOMER, SUPPLIER, partyRoutes } from './modules/parties/parties.routes.js';
import productsRoutes from './modules/products/products.routes.js';
import documentsModule from './modules/documents/routes.js';
import inventoryRoutes from './modules/inventory/routes.js';
import expenseRoutes from './modules/expenses/routes.js';
import taxRoutes from './modules/tax/routes.js';
import commercialReportRoutes from './modules/reports/commercial.routes.js';
import dashboardRoutes from './modules/reports/dashboard.routes.js';
import financialReportRoutes from './modules/reports/financial.routes.js';
import partyReportRoutes from './modules/reports/parties.routes.js';
import authRoutes from './modules/auth/auth.routes.js';
import companiesRoutes from './modules/companies/companies.routes.js';
import rbacRoutes from './modules/rbac/rbac.routes.js';
import settingsRoutes from './modules/settings/settings.routes.js';
import usersRoutes from './modules/users/users.routes.js';
import authPlugin from './plugins/auth.js';

// PostgreSQL error codes mapped to client errors.
const PG_ERRORS: Record<string, [number, string, string]> = {
  '23505': [409, 'DUPLICATE', 'A record with the same unique value already exists'],
  '23503': [400, 'INVALID_REFERENCE', 'A referenced record does not exist'],
  '23514': [400, 'CONSTRAINT_VIOLATION', 'A value is outside the allowed range or format'],
  '42501': [403, 'FORBIDDEN', 'Operation not permitted'],
};

export interface BuildOptions {
  env: Env;
  pool: pg.Pool;
  logger?: boolean;
}

export async function buildApp({ env, pool, logger = true }: BuildOptions): Promise<FastifyInstance> {
  const app = Fastify({
    logger: logger ? { level: env.NODE_ENV === 'production' ? 'info' : 'debug', redact: ['req.headers.authorization', 'req.headers.cookie'] } : false,
    trustProxy: env.NODE_ENV === 'production',
    bodyLimit: 1_048_576,
    genReqId: () => crypto.randomUUID(),
  });

  app.decorate('deps', { env, pool, tokens: new TokenService(env.JWT_SECRET, env.ACCESS_TOKEN_TTL_SECONDS) });

  await app.register(helmet, {
    contentSecurityPolicy: { directives: { defaultSrc: ["'none'"], frameAncestors: ["'none'"] } },
  });
  await app.register(cors, {
    origin: env.CORS_ORIGINS.split(',').map((o) => o.trim()).filter(Boolean),
    credentials: true,
    allowedHeaders: ['Content-Type', 'Authorization', 'X-CSRF-Protection'],
  });
  await app.register(cookie);
  await app.register(rateLimit, { max: env.NODE_ENV === 'test' ? 10_000 : 300, timeWindow: '1 minute' });
  await app.register(authPlugin);

  app.setErrorHandler((err: Error, req, reply) => {
    if (err instanceof AppError) {
      return reply.code(err.statusCode).send({ error: { code: err.code, message: err.message, details: err.details } });
    }
    const pgCode = (err as { code?: string }).code;
    const mapped = pgCode ? PG_ERRORS[pgCode] : undefined;
    if (mapped) {
      req.log.warn({ err }, 'database constraint');
      // Messages raised by our own ledger triggers are safe and useful to show.
      const fromTrigger = /PL\/pgSQL function/.test((err as { where?: string }).where ?? '');
      return reply.code(mapped[0]).send({ error: { code: fromTrigger ? 'LEDGER_RULE' : mapped[1], message: fromTrigger ? err.message : mapped[2] } });
    }
    const status = (err as { statusCode?: number }).statusCode;
    if (status && status >= 400 && status < 500) {
      return reply.code(status).send({ error: { code: (err as { code?: string }).code ?? 'BAD_REQUEST', message: err.message } });
    }
    req.log.error({ err }, 'unhandled error');
    return reply.code(500).send({ error: { code: 'INTERNAL_ERROR', message: 'Internal server error' } });
  });

  app.get('/api/health', async () => {
    await pool.query('SELECT 1');
    return { status: 'ok' };
  });

  await app.register(authRoutes, { prefix: '/api/auth' });
  await app.register(async (api) => {
    await api.register(usersRoutes);
    await api.register(rbacRoutes);
    await api.register(companiesRoutes);
    await api.register(settingsRoutes);
    await api.register(auditRoutes);
    await api.register(accountsRoutes);
    await api.register(fiscalRoutes);
    await api.register(journalRoutes);
    await api.register(trialBalanceRoutes);
    await api.register(partyRoutes(CUSTOMER));
    await api.register(partyRoutes(SUPPLIER));
    await api.register(productsRoutes);
    await api.register(documentsModule);
    await api.register(inventoryRoutes);
    await api.register(expenseRoutes);
    await api.register(taxRoutes);
    await api.register(financialReportRoutes);
    await api.register(partyReportRoutes);
    await api.register(commercialReportRoutes);
    await api.register(dashboardRoutes);
  }, { prefix: '/api' });

  return app;
}
