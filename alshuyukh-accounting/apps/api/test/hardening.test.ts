import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { buildApp } from '../src/app.js';
import { loadEnv } from '../src/config/env.js';
import { createPool } from '../src/db/pool.js';
import { migrate } from '../src/db/migrate.js';
import { addMember, client, register, setupApp, type TestContext } from './helpers.js';

let t: TestContext;
beforeAll(async () => { t = await setupApp(); });
afterAll(async () => { await t.close(); });

describe('configuration', () => {
  const base = { DATABASE_URL: 'postgres://x', JWT_SECRET: 'a'.repeat(40) };
  it('refuses unsafe production settings', () => {
    expect(() => loadEnv({ ...base, NODE_ENV: 'production', JWT_SECRET: 'change-me-' + 'a'.repeat(40), ZATCA_ENCRYPTION_KEY: 'x' })).toThrow(/JWT_SECRET/);
    expect(() => loadEnv({ ...base, NODE_ENV: 'production' })).toThrow(/ZATCA_ENCRYPTION_KEY/);
    expect(() => loadEnv({ ...base, JWT_SECRET: 'short' })).toThrow(/at least 32/);
    expect(() => loadEnv({ ...base, NODE_ENV: 'production', ZATCA_ENCRYPTION_KEY: Buffer.alloc(32).toString('base64') })).not.toThrow();
  });
});

describe('probes', () => {
  it('reports liveness and readiness', async () => {
    expect((await t.app.inject({ method: 'GET', url: '/api/health' })).json()).toEqual({ status: 'ok' });
    expect((await t.app.inject({ method: 'GET', url: '/api/ready' })).json()).toEqual({ status: 'ready' });
  });
});

describe('migrations', () => {
  it('are idempotent and refuse a migration edited after it ran', async () => {
    const url = process.env.TEST_DATABASE_URL_MIGRATE!;
    await expect(migrate(url)).resolves.toBeUndefined();
    const { rows: [m] } = await t.ownerPool.query<{ version: string; checksum: string }>(`SELECT version, checksum FROM schema_migrations ORDER BY version LIMIT 1`);
    await t.ownerPool.query(`UPDATE schema_migrations SET checksum = 'tampered' WHERE version = $1`, [m!.version]);
    try {
      await expect(migrate(url)).rejects.toThrow(/was modified after it was applied/);
    } finally {
      await t.ownerPool.query(`UPDATE schema_migrations SET checksum = $2 WHERE version = $1`, [m!.version, m!.checksum]);
    }
  });
});

describe('errors', () => {
  it('hides internal errors from clients and records them for the platform', async () => {
    const env = loadEnv();
    const pool = createPool(env.DATABASE_URL);
    const app = await buildApp({ env, pool, logger: false });
    app.get('/api/__boom', async () => { throw new Error('secret internal detail'); });
    await app.ready();
    try {
      const res = await app.inject({ method: 'GET', url: '/api/__boom?x=1' });
      expect(res.statusCode).toBe(500);
      expect(res.json()).toEqual({ error: { code: 'INTERNAL_ERROR', message: 'Internal server error' } });
      expect(res.body).not.toContain('secret');
      // Logged asynchronously; the response never waits for it.
      let rows: { path: string; message: string }[] = [];
      for (let i = 0; i < 20 && !rows.length; i++) {
        rows = (await t.ownerPool.query(`SELECT path, message FROM system_errors WHERE message = 'secret internal detail'`)).rows;
        if (!rows.length) await new Promise((r) => setTimeout(r, 50));
      }
      expect(rows).toEqual([{ path: '/api/__boom', message: 'secret internal detail' }]);
    } finally {
      await app.close();
      await pool.end();
    }
  });

  it('maps database constraint errors to safe client errors', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    const a = await api.post('/api/cost-centers', { code: 'DUP', name: 'مركز' });
    expect(a.statusCode).toBe(201);
    const dup = await api.post('/api/cost-centers', { code: 'DUP', name: 'مركز آخر' });
    expect(dup.statusCode).toBe(409);
    expect(dup.json().error.code).toBe('DUPLICATE');
    expect(dup.body).not.toMatch(/cost_centers|constraint|duplicate key/i);
  });

  it('rejects malformed and oversized input', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    expect((await api.get('/api/journal-entries?limit=100000')).statusCode).toBe(400);
    expect((await api.get('/api/reports/profit-loss?dateFrom=2026-02-30&dateTo=2026-03-01')).statusCode).toBe(400);
    expect((await api.get('/api/customers/not-a-uuid')).statusCode).toBe(400);
    const big = await t.app.inject({ method: 'POST', url: '/api/customers', headers: { authorization: `Bearer ${s.token}`, 'content-type': 'application/json' },
      payload: JSON.stringify({ nameAr: 'x'.repeat(2 * 1024 * 1024) }) });
    expect(big.statusCode).toBe(413);
  });
});

describe('expense categories', () => {
  it('creates, renames, deactivates and validates categories', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    const accounts = (await api.get('/api/accounts?postableOnly=true')).json().data as { id: string; code: string; type: string }[];
    const travel = accounts.find((a) => a.code === '6200')!;
    const cat = await api.post('/api/expense-categories', { code: 'travel', nameAr: 'سفر', accountId: travel.id, vatCategory: 'S' });
    expect(cat.statusCode).toBe(201);
    expect(cat.json()).toMatchObject({ code: 'TRAVEL', accountCode: '6200', isActive: true });
    expect((await api.post('/api/expense-categories', { code: 'TRAVEL', nameAr: 'سفر', accountId: travel.id })).statusCode).toBe(409);
    const cash = accounts.find((a) => a.code === '1100')!;
    expect((await api.post('/api/expense-categories', { code: 'BAD', nameAr: 'حساب أصل', accountId: cash.id })).json().error.code).toBe('INVALID_ACCOUNT');
    const calc = (await api.post('/api/expenses/calculate', { expenseDate: '2026-03-10', lines: [{ categoryId: cat.json().id, amount: '200' }] })).json();
    expect(calc.totals).toEqual({ subtotal: '200.00', taxAmount: '30.00', total: '230.00' });
    const upd = await api.patch(`/api/expense-categories/${cat.json().id}`, { nameAr: 'سفر وتنقلات', isActive: false });
    expect(upd.json()).toMatchObject({ nameAr: 'سفر وتنقلات', isActive: false });
    const other = client(t.app, (await register(t.app)).token);
    expect((await other.patch(`/api/expense-categories/${cat.json().id}`, { nameAr: 'اختراق' })).statusCode).toBe(404);
    // A deactivated category cannot be used on new expenses.
    expect((await api.post('/api/expenses/calculate', { expenseDate: '2026-03-10', lines: [{ categoryId: cat.json().id, amount: '200' }] })).statusCode).toBe(400);
  });
});

describe('member management', () => {
  it('does not let a user manager disable or edit someone with more access', async () => {
    const owner = await register(t.app);
    const api = client(t.app, owner.token);
    const role = await api.post('/api/roles', { code: 'USER_DESK', nameAr: 'مكتب المستخدمين', nameEn: 'User desk', permissions: ['user.view', 'user.manage', 'role.view'] });
    expect(role.statusCode).toBe(201);
    const desk = await addMember(t.app, owner, 'VIEWER');
    const admin = await addMember(t.app, owner, 'COMPANY_ADMIN');
    expect((await api.put(`/api/users/${desk.userId}/roles`, { roleIds: [role.json().id] })).statusCode).toBe(200);
    const deskApi = client(t.app, desk.token);
    expect((await deskApi.patch(`/api/users/${admin.userId}`, { status: 'DISABLED' })).statusCode).toBe(403);
    expect((await deskApi.put(`/api/users/${admin.userId}/roles`, { roleIds: [role.json().id] })).statusCode).toBe(403);
    // A member with no more access than the manager can still be managed.
    const peer = await addMember(t.app, owner, 'VIEWER');
    expect((await deskApi.patch(`/api/users/${peer.userId}`, { status: 'DISABLED' })).statusCode).toBe(403);
    expect((await api.put(`/api/users/${peer.userId}/roles`, { roleIds: [role.json().id] })).statusCode).toBe(200);
    expect((await deskApi.patch(`/api/users/${peer.userId}`, { status: 'DISABLED' })).statusCode).toBe(200);
  });
});

describe('reports', () => {
  it('refuses report ranges longer than five years', async () => {
    const s = await register(t.app);
    const res = await client(t.app, s.token).get('/api/reports/profit-loss?dateFrom=2015-01-01&dateTo=2026-12-31');
    expect(res.statusCode).toBe(400);
    expect(res.json().error.code).toBe('RANGE_TOO_LONG');
  });
});
