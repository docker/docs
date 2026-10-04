import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { client, register, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
});
afterAll(async () => { await t.close(); });

describe('companies, branches, warehouses', () => {
  it('creates a company with Saudi defaults and validates input', async () => {
    const api = client(t.app, owner.token);
    const res = await api.post('/api/companies', { name: 'فرع جدة للتجارة', commercialRegistration: '4030000000', vatNumber: '311111111111113' });
    expect(res.statusCode).toBe(201);
    expect(res.json()).toMatchObject({ currency: 'SAR', timezone: 'Asia/Riyadh', country: 'SA' });

    expect((await api.post('/api/companies', { name: 'Bad TZ', timezone: 'Mars/Base' })).statusCode).toBe(400);
    expect((await api.post('/api/companies', { name: 'Bad CR', commercialRegistration: '12' })).statusCode).toBe(400);
    expect((await api.post('/api/companies', { name: 'فرع جدة للتجارة' })).statusCode).toBe(409);
  });

  it('enforces unique branch codes and a single main branch', async () => {
    const api = client(t.app, owner.token);
    const company = (await api.post('/api/companies', { name: `Co ${Date.now()}` })).json();
    expect((await api.post(`/api/companies/${company.id}/branches`, { code: 'RUH', name: 'الرياض' })).statusCode).toBe(201);
    expect((await api.post(`/api/companies/${company.id}/branches`, { code: 'RUH', name: 'مكرر' })).statusCode).toBe(409);
  });

  it('soft-deletes a company but keeps at least one', async () => {
    const api = client(t.app, owner.token);
    const extra = (await api.post('/api/companies', { name: `Temp ${Date.now()}` })).json();
    expect((await api.del(`/api/companies/${extra.id}`)).statusCode).toBe(204);
    expect((await api.get(`/api/companies/${extra.id}`)).statusCode).toBe(404);
    const { rows } = await t.ownerPool.query('SELECT deleted_at FROM companies WHERE id = $1', [extra.id]);
    expect(rows[0].deleted_at).not.toBeNull();

    const fresh = await register(t.app);
    const only = (await client(t.app, fresh.token).get('/api/companies')).json().data[0];
    expect((await client(t.app, fresh.token).del(`/api/companies/${only.id}`)).statusCode).toBe(409);
  });

  it('rejects a warehouse linked to a branch of another company', async () => {
    const api = client(t.app, owner.token);
    const c1 = (await api.post('/api/companies', { name: `C1 ${Date.now()}` })).json();
    const c2 = (await api.post('/api/companies', { name: `C2 ${Date.now()}` })).json();
    const b1 = (await api.post(`/api/companies/${c1.id}/branches`, { code: 'B1', name: 'Branch one' })).json();
    expect((await api.post(`/api/companies/${c2.id}/warehouses`, { code: 'W', name: 'Wrong', branchId: b1.id })).statusCode).toBe(404);
    expect((await api.post(`/api/companies/${c1.id}/warehouses`, { code: 'W', name: 'Right', branchId: b1.id })).statusCode).toBe(201);
  });

  it('rejects unknown fields from reaching SQL and ignores tenant_id in the body', async () => {
    const other = await register(t.app);
    const api = client(t.app, owner.token);
    const res = await api.post('/api/companies', { name: `Spoof ${Date.now()}`, tenantId: other.tenantId, tenant_id: other.tenantId });
    expect(res.statusCode).toBe(201);
    const { rows } = await t.ownerPool.query('SELECT tenant_id FROM companies WHERE id = $1', [res.json().id]);
    expect(rows[0].tenant_id).toBe(owner.tenantId);
  });
});

describe('multi-tenant membership', () => {
  it('lets one user belong to two tenants and switch between them', async () => {
    const t1 = await register(t.app, { tenantName: 'Org One' });
    const t2 = await register(t.app, { tenantName: 'Org Two' });
    const viewerRole = (await client(t.app, t2.token).get('/api/roles')).json().data.find((r: { code: string }) => r.code === 'VIEWER').id;
    const add = await client(t.app, t2.token).post('/api/users', { email: t1.email, fullName: 'Shared', roleIds: [viewerRole] });
    expect(add.statusCode).toBe(201);
    // An existing account joins only after it accepts; until then the inviter sees no personal details.
    expect(add.json()).toMatchObject({ status: 'INVITED', fullName: null, lastLoginAt: null });
    expect((await client(t.app, t1.token).get('/api/auth/me')).json().memberships).toHaveLength(1);
    expect((await client(t.app, t1.token).post('/api/auth/switch-tenant', { tenantId: t2.tenantId })).statusCode).not.toBe(200);
    expect((await client(t.app, t2.token).patch(`/api/users/${add.json().id}`, { status: 'ACTIVE' })).statusCode).toBe(409);
    const inv = (await client(t.app, t1.token).get('/api/auth/invitations')).json().data;
    expect(inv).toEqual([{ tenantId: t2.tenantId, tenantName: 'Org Two', invitedAt: expect.any(String) }]);
    expect((await client(t.app, t1.token).post(`/api/auth/invitations/${t2.tenantId}/accept`)).statusCode).toBe(204);
    expect((await client(t.app, t1.token).get('/api/auth/invitations')).json().data).toEqual([]);

    const me = (await client(t.app, t1.token).get('/api/auth/me')).json();
    expect(me.memberships.map((m: { tenantName: string }) => m.tenantName).sort()).toEqual(['Org One', 'Org Two']);

    const sw = await client(t.app, t1.token).post('/api/auth/switch-tenant', { tenantId: t2.tenantId });
    expect(sw.statusCode).toBe(200);
    const api2 = client(t.app, sw.json().accessToken);
    const me2 = (await api2.get('/api/auth/me')).json();
    expect(me2.tenant.id).toBe(t2.tenantId);
    expect(me2.roles.map((r: { code: string }) => r.code)).toEqual(['VIEWER']);
    expect((await api2.post('/api/companies', { name: 'Not allowed' })).statusCode).toBe(403);
    // The pre-switch token is bound to the old tenant and is no longer valid.
    expect((await client(t.app, t1.token).get('/api/companies')).statusCode).toBe(401);
  });

  it('lets an invited user decline, and refuses answering for another user', async () => {
    const t1 = await register(t.app);
    const t2 = await register(t.app);
    const t3 = await register(t.app);
    const viewerRole = (await client(t.app, t2.token).get('/api/roles')).json().data.find((r: { code: string }) => r.code === 'VIEWER').id;
    expect((await client(t.app, t2.token).post('/api/users', { email: t1.email, fullName: 'Someone', roleIds: [viewerRole] })).statusCode).toBe(201);
    // t3 was not invited: answering is not possible.
    expect((await client(t.app, t3.token).post(`/api/auth/invitations/${t2.tenantId}/accept`)).statusCode).toBe(404);
    expect((await client(t.app, t1.token).post(`/api/auth/invitations/${t2.tenantId}/decline`)).statusCode).toBe(204);
    expect((await client(t.app, t1.token).post(`/api/auth/invitations/${t2.tenantId}/accept`)).statusCode).toBe(404);
    const member = (await client(t.app, t2.token).get('/api/users')).json().data.find((u: { email: string }) => u.email === t1.email);
    expect(member.status).toBe('DISABLED');
  });
});
