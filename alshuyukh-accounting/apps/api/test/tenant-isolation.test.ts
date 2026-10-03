import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { withTx } from '../src/db/tx.js';
import { client, register, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let a: Session;
let b: Session;
let companyA: { id: string };
let branchA: { id: string };

beforeAll(async () => {
  t = await setupApp();
  a = await register(t.app, { tenantName: 'Tenant A', companyName: 'Company A' });
  b = await register(t.app, { tenantName: 'Tenant B', companyName: 'Company B' });
  companyA = (await client(t.app, a.token).get('/api/companies')).json().data[0];
  branchA = (await client(t.app, a.token).get(`/api/companies/${companyA.id}/branches`)).json().data[0];
});
afterAll(async () => { await t.close(); });

describe('API-level isolation (tenant B attacking tenant A by ID)', () => {
  it('lists only its own companies', async () => {
    const list = (await client(t.app, b.token).get('/api/companies')).json().data;
    expect(list.map((c: { name: string }) => c.name)).toEqual(['Company B']);
  });

  it('cannot read, update or delete another tenant\'s company (404, no existence leak)', async () => {
    const api = client(t.app, b.token);
    expect((await api.get(`/api/companies/${companyA.id}`)).statusCode).toBe(404);
    expect((await api.patch(`/api/companies/${companyA.id}`, { name: 'Hacked' })).statusCode).toBe(404);
    expect((await api.del(`/api/companies/${companyA.id}`)).statusCode).toBe(404);
    const still = (await client(t.app, a.token).get(`/api/companies/${companyA.id}`)).json();
    expect(still.name).toBe('Company A');
  });

  it('cannot list or create branches/warehouses under another tenant\'s company', async () => {
    const api = client(t.app, b.token);
    expect((await api.get(`/api/companies/${companyA.id}/branches`)).statusCode).toBe(404);
    expect((await api.post(`/api/companies/${companyA.id}/branches`, { code: 'X1', name: 'Injected' })).statusCode).toBe(404);
    expect((await api.patch(`/api/branches/${branchA.id}`, { name: 'Hacked' })).statusCode).toBe(404);
    expect((await api.post(`/api/companies/${companyA.id}/warehouses`, { code: 'W1', name: 'Injected' })).statusCode).toBe(404);
  });

  it('cannot see or manage another tenant\'s users', async () => {
    const api = client(t.app, b.token);
    const users = (await api.get('/api/users')).json().data;
    expect(users.map((u: { id: string }) => u.id)).toEqual([b.userId]);
    expect((await api.get(`/api/users/${a.userId}`)).statusCode).toBe(404);
    expect((await api.patch(`/api/users/${a.userId}`, { status: 'DISABLED' })).statusCode).toBe(404);
  });

  it('cannot read another tenant\'s audit log', async () => {
    const logs = (await client(t.app, b.token).get('/api/audit-logs?limit=200')).json().data;
    expect(logs.every((l: { userId: string | null }) => l.userId === b.userId)).toBe(true);
    const byEntity = (await client(t.app, b.token).get(`/api/audit-logs?entityId=${a.tenantId}`)).json().data;
    expect(byEntity).toEqual([]);
  });

  it('cannot switch into a tenant it is not a member of', async () => {
    const res = await client(t.app, b.token).post('/api/auth/switch-tenant', { tenantId: a.tenantId });
    expect(res.statusCode).toBe(403);
  });

  it('cannot assign another tenant\'s custom role', async () => {
    const role = await client(t.app, a.token).post('/api/roles', { code: 'A_ONLY', nameAr: 'دور خاص', nameEn: 'A only', permissions: ['company.view'] });
    expect(role.statusCode).toBe(201);
    const res = await client(t.app, b.token).post('/api/users', {
      email: 'victim@example.test', fullName: 'Victim', initialPassword: 'Str0ng-Passw0rd!', roleIds: [role.json().id],
    });
    expect(res.statusCode).toBe(400);
    expect(res.json().error.code).toBe('INVALID_ROLE');
  });
});

describe('database-level isolation (Row-Level Security)', () => {
  it('runs the API with a role that cannot bypass RLS', async () => {
    const { rows: [r] } = await t.pool.query(`SELECT rolsuper, rolbypassrls FROM pg_roles WHERE rolname = current_user`);
    expect(r).toEqual({ rolsuper: false, rolbypassrls: false });
  });

  it('returns only the context tenant\'s rows even without a WHERE tenant_id filter', async () => {
    const rows = await withTx(t.pool, { tenantId: b.tenantId, userId: b.userId }, async (db) =>
      (await db.query<{ tenant_id: string }>('SELECT tenant_id FROM companies')).rows);
    expect(rows.length).toBeGreaterThan(0);
    expect(rows.every((r) => r.tenant_id === b.tenantId)).toBe(true);
  });

  it('returns nothing when no tenant context is set', async () => {
    const rows = await withTx(t.pool, {}, async (db) => (await db.query('SELECT id FROM companies')).rows);
    expect(rows).toEqual([]);
  });

  it('rejects writing a row into another tenant', async () => {
    await expect(withTx(t.pool, { tenantId: b.tenantId, userId: b.userId }, (db) =>
      db.query(`INSERT INTO companies (tenant_id, name) VALUES ($1, 'Injected')`, [a.tenantId]),
    )).rejects.toThrow(/row-level security/);
  });

  it('rejects linking a branch to another tenant\'s company (composite FK)', async () => {
    await expect(withTx(t.pool, { tenantId: b.tenantId, userId: b.userId }, (db) =>
      db.query(`INSERT INTO branches (tenant_id, company_id, code, name) VALUES ($1, $2, 'X', 'Injected')`, [b.tenantId, companyA.id]),
    )).rejects.toThrow(/foreign key/);
  });

  it('makes cross-tenant UPDATEs affect zero rows', async () => {
    const res = await withTx(t.pool, { tenantId: b.tenantId, userId: b.userId }, (db) =>
      db.query(`UPDATE companies SET name = 'Hacked' WHERE id = $1`, [companyA.id]));
    expect(res.rowCount).toBe(0);
  });

  it('does not allow the app role to modify system roles', async () => {
    const res = await withTx(t.pool, { tenantId: b.tenantId, userId: b.userId }, (db) =>
      db.query(`UPDATE roles SET name_en = 'pwned' WHERE code = 'TENANT_OWNER' AND tenant_id IS NULL`));
    expect(res.rowCount).toBe(0);
  });

  it('does not allow the app role to hard-delete business records', async () => {
    await expect(withTx(t.pool, { tenantId: a.tenantId, userId: a.userId }, (db) =>
      db.query('DELETE FROM companies WHERE id = $1', [companyA.id]),
    )).rejects.toThrow(/permission denied/);
  });
});
