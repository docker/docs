import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { ALL_PERMISSION_CODES, SYSTEM_ROLES, resolveRolePermissions } from '../src/modules/rbac/catalog.js';
import { addMember, client, register, roleId, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;
let companyId: string;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
  companyId = (await client(t.app, owner.token).get('/api/companies')).json().data[0].id;
});
afterAll(async () => { await t.close(); });

describe('catalog', () => {
  it('only references permissions that exist', () => {
    for (const role of SYSTEM_ROLES) {
      for (const p of resolveRolePermissions(role)) expect(ALL_PERMISSION_CODES).toContain(p);
    }
  });

  it('gives VIEWER read-only permissions', () => {
    const viewer = SYSTEM_ROLES.find((r) => r.code === 'VIEWER')!;
    expect(resolveRolePermissions(viewer).every((p) => p.endsWith('.view'))).toBe(true);
  });

  it('is synced to the database', async () => {
    const { rows } = await t.ownerPool.query(`SELECT code FROM roles WHERE is_system ORDER BY code`);
    expect(rows.map((r) => r.code).sort()).toEqual(SYSTEM_ROLES.map((r) => r.code).sort());
  });
});

describe('permission checks', () => {
  it('lets an accountant view companies but not create them or manage users', async () => {
    const acc = await addMember(t.app, owner, 'ACCOUNTANT');
    const api = client(t.app, acc.token);
    expect((await api.get('/api/companies')).statusCode).toBe(200);
    expect((await api.post('/api/companies', { name: 'Nope' })).statusCode).toBe(403);
    expect((await api.get('/api/users')).statusCode).toBe(403);
    expect((await api.patch('/api/settings/tenant', { locale: 'en' })).statusCode).toBe(403);
    expect((await api.get('/api/audit-logs')).statusCode).toBe(403);
  });

  it('lets a viewer read but not write', async () => {
    const viewer = await addMember(t.app, owner, 'VIEWER');
    const api = client(t.app, viewer.token);
    expect((await api.get('/api/users')).statusCode).toBe(200);
    expect((await api.patch(`/api/companies/${companyId}`, { city: 'Riyadh' })).statusCode).toBe(403);
    expect((await api.post(`/api/companies/${companyId}/branches`, { code: 'B2', name: 'Jeddah' })).statusCode).toBe(403);
  });

  it('applies role changes to existing tokens immediately', async () => {
    const viewer = await addMember(t.app, owner, 'VIEWER');
    const api = client(t.app, viewer.token);
    expect((await api.post(`/api/companies/${companyId}/branches`, { code: 'DMM', name: 'Dammam' })).statusCode).toBe(403);
    const res = await client(t.app, owner.token).put(`/api/users/${viewer.userId}/roles`, { roleIds: [await roleId(t.app, owner.token, 'COMPANY_ADMIN')] });
    expect(res.statusCode).toBe(200);
    expect((await api.post(`/api/companies/${companyId}/branches`, { code: 'DMM', name: 'Dammam' })).statusCode).toBe(201);
  });

  it('blocks a disabled member immediately', async () => {
    const m = await addMember(t.app, owner, 'VIEWER');
    expect((await client(t.app, owner.token).patch(`/api/users/${m.userId}`, { status: 'DISABLED' })).statusCode).toBe(200);
    expect((await client(t.app, m.token).get('/api/companies')).statusCode).toBe(401);
    const relogin = await t.app.inject({ method: 'POST', url: '/api/auth/login', payload: { email: m.email, password: 'Str0ng-Passw0rd!' } });
    expect(relogin.statusCode).toBe(403);
  });
});

describe('privilege escalation', () => {
  it('prevents a company admin from granting TENANT_OWNER', async () => {
    const admin = await addMember(t.app, owner, 'COMPANY_ADMIN');
    const res = await client(t.app, admin.token).post('/api/users', {
      email: `esc-${Date.now()}@example.test`, fullName: 'Escalate', initialPassword: 'Str0ng-Passw0rd!',
      roleIds: [await roleId(t.app, owner.token, 'TENANT_OWNER')],
    });
    expect(res.statusCode).toBe(403);
  });

  it('prevents creating a custom role with permissions the actor lacks', async () => {
    const admin = await addMember(t.app, owner, 'COMPANY_ADMIN');
    const res = await client(t.app, admin.token).post('/api/roles', {
      code: 'SNEAKY', nameAr: 'خفي', nameEn: 'Sneaky', permissions: ['subscription.manage'],
    });
    expect(res.statusCode).toBe(403);
  });

  it('prevents users from changing their own roles or status', async () => {
    const admin = await addMember(t.app, owner, 'COMPANY_ADMIN');
    const api = client(t.app, admin.token);
    expect((await api.put(`/api/users/${admin.userId}/roles`, { roleIds: [await roleId(t.app, owner.token, 'VIEWER')] })).statusCode).toBe(403);
    expect((await api.patch(`/api/users/${admin.userId}`, { status: 'DISABLED' })).statusCode).toBe(403);
  });

  it('protects the owner from being disabled or demoted', async () => {
    const admin = await addMember(t.app, owner, 'COMPANY_ADMIN');
    const api = client(t.app, admin.token);
    expect((await api.patch(`/api/users/${owner.userId}`, { status: 'DISABLED' })).statusCode).toBe(403);
    expect((await api.put(`/api/users/${owner.userId}/roles`, { roleIds: [await roleId(t.app, owner.token, 'VIEWER')] })).statusCode).toBe(403);
  });

  it('refuses to modify or delete system roles', async () => {
    const id = await roleId(t.app, owner.token, 'ACCOUNTANT');
    const api = client(t.app, owner.token);
    expect((await api.patch(`/api/roles/${id}`, { nameEn: 'Changed' })).statusCode).toBe(403);
    expect((await api.del(`/api/roles/${id}`)).statusCode).toBe(403);
  });
});

describe('custom roles', () => {
  it('creates, assigns and soft-deletes a custom role', async () => {
    const api = client(t.app, owner.token);
    const created = await api.post('/api/roles', {
      code: 'branch_auditor', nameAr: 'مدقق فروع', nameEn: 'Branch auditor', permissions: ['company.view', 'audit.view'],
    });
    expect(created.statusCode).toBe(201);
    expect(created.json()).toMatchObject({ code: 'BRANCH_AUDITOR', isSystem: false, permissions: ['audit.view', 'company.view'] });

    const dup = await api.post('/api/roles', { code: 'BRANCH_AUDITOR', nameAr: 'مكرر', nameEn: 'Dup', permissions: [] });
    expect(dup.statusCode).toBe(409);
    const reserved = await api.post('/api/roles', { code: 'VIEWER', nameAr: 'مشاهد', nameEn: 'Viewer', permissions: [] });
    expect(reserved.statusCode).toBe(409);
    const unknown = await api.post('/api/roles', { code: 'BAD', nameAr: 'سيء', nameEn: 'Bad', permissions: ['does.not_exist'] });
    expect(unknown.statusCode).toBe(400);

    const m = await addMember(t.app, owner, 'VIEWER');
    await api.put(`/api/users/${m.userId}/roles`, { roleIds: [created.json().id] });
    expect((await api.del(`/api/roles/${created.json().id}`)).statusCode).toBe(409);
    await api.put(`/api/users/${m.userId}/roles`, { roleIds: [await roleId(t.app, owner.token, 'VIEWER')] });
    expect((await api.del(`/api/roles/${created.json().id}`)).statusCode).toBe(204);
    const roles = (await api.get('/api/roles')).json().data;
    expect(roles.find((r: { code: string }) => r.code === 'BRANCH_AUDITOR')).toBeUndefined();
  });
});
