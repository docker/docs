import type pg from 'pg';
import type { Env } from './config/env.js';
import type { Db } from './db/tx.js';
import type { TokenService } from './lib/tokens.js';
import type { AuditMeta } from './modules/audit/audit.service.js';

export interface Deps {
  env: Env;
  pool: pg.Pool;
  tokens: TokenService;
}

export interface AuthContext {
  userId: string;
  tenantId: string;
  sessionId: string;
  isOwner: boolean;
  isPlatformAdmin: boolean;
  permissions: ReadonlySet<string>;
}

declare module 'fastify' {
  interface FastifyInstance {
    deps: Deps;
    authenticate: (req: FastifyRequest) => Promise<void>;
  }
  interface FastifyRequest {
    auth: AuthContext | null;
    /** Runs fn in a transaction scoped to the authenticated user's tenant. */
    tenantTx<T>(fn: (db: Db) => Promise<T>): Promise<T>;
    auditMeta(): AuditMeta;
  }
}
