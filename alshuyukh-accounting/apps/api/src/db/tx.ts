import type pg from 'pg';

export type Db = pg.PoolClient;

export interface DbContext {
  tenantId?: string | null;
  userId?: string | null;
}

/**
 * Runs `fn` inside a transaction with the request context applied.
 *
 * `app.tenant_id` and `app.user_id` are set with is_local = true, so they
 * exist only for this transaction and cannot leak to the next request that
 * reuses the pooled connection. Row-Level Security policies read them.
 */
export async function withTx<T>(pool: pg.Pool, ctx: DbContext, fn: (db: Db) => Promise<T>): Promise<T> {
  const client = await pool.connect();
  try {
    await client.query('BEGIN');
    await client.query(
      `SELECT set_config('app.tenant_id', $1, true), set_config('app.user_id', $2, true)`,
      [ctx.tenantId ?? '', ctx.userId ?? ''],
    );
    const result = await fn(client);
    await client.query('COMMIT');
    return result;
  } catch (err) {
    await client.query('ROLLBACK').catch(() => undefined);
    throw err;
  } finally {
    client.release();
  }
}

/** Switches the tenant context inside an open transaction (used during registration). */
export async function setTenantContext(db: Db, tenantId: string): Promise<void> {
  await db.query(`SELECT set_config('app.tenant_id', $1, true)`, [tenantId]);
}
