import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { client, login, register, setupApp, type Session, type TestContext } from './helpers.js';

let t: TestContext;
let owner: Session;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
});
afterAll(async () => { await t.close(); });

const logs = async (query = '') => (await client(t.app, owner.token).get(`/api/audit-logs?limit=200${query}`)).json().data as Array<Record<string, any>>;

describe('audit log', () => {
  it('records registration and login', async () => {
    await login(t.app, owner.email);
    const actions = (await logs()).map((l) => l.action);
    expect(actions).toContain('REGISTER');
    expect(actions).toContain('LOGIN');
  });

  it('records failed logins with the reason and request metadata', async () => {
    await t.app.inject({ method: 'POST', url: '/api/auth/login', payload: { email: owner.email, password: 'Wrong-passw0rd' }, headers: { 'user-agent': 'audit-test' } });
    const { rows } = await t.ownerPool.query(
      `SELECT new_values, user_agent, host(ip_address) AS ip FROM audit_logs WHERE action = 'LOGIN_FAILED' AND user_id = $1`, [owner.userId]);
    expect(rows[0].new_values).toMatchObject({ reason: 'bad_password' });
    expect(rows[0].user_agent).toBe('audit-test');
    expect(rows[0].ip).toBe('127.0.0.1');
  });

  it('records updates with old and new values', async () => {
    const api = client(t.app, owner.token);
    const company = (await api.get('/api/companies')).json().data[0];
    await api.patch(`/api/companies/${company.id}`, { city: 'الرياض' });
    const [entry] = await logs(`&entityId=${company.id}&action=UPDATE`);
    expect(entry!.oldValues.city).toBeNull();
    expect(entry!.newValues.city).toBe('الرياض');
    expect(entry!.userId).toBe(owner.userId);
  });

  it('records settings and permission changes', async () => {
    const api = client(t.app, owner.token);
    await api.patch('/api/settings/tenant', { fiscalYearStartMonth: 7 });
    await api.post('/api/roles', { code: 'AUDIT_TEST', nameAr: 'اختبار', nameEn: 'Test', permissions: ['company.view'] });
    const actions = (await logs()).map((l) => l.action);
    expect(actions).toContain('SETTINGS_CHANGE');
    expect(actions).toContain('PERMISSION_CHANGE');
  });

  it('never stores password hashes or secrets', async () => {
    await client(t.app, owner.token).post('/api/users', {
      email: `aud-${Date.now()}@example.test`, fullName: 'Audit Member', initialPassword: 'Str0ng-Passw0rd!',
      roleIds: [(await client(t.app, owner.token).get('/api/roles')).json().data.find((r: { code: string }) => r.code === 'VIEWER').id],
    });
    const { rows } = await t.ownerPool.query(`SELECT old_values::text || new_values::text AS v FROM audit_logs WHERE tenant_id = $1`, [owner.tenantId]);
    for (const r of rows) {
      expect(r.v ?? '').not.toContain('argon2');
      expect(r.v ?? '').not.toContain('Str0ng-Passw0rd!');
    }
  });

  it('is append-only, even for the schema owner', async () => {
    await expect(t.ownerPool.query(`UPDATE audit_logs SET action = 'TAMPERED' WHERE tenant_id = $1`, [owner.tenantId])).rejects.toThrow(/append-only/);
    await expect(t.ownerPool.query(`DELETE FROM audit_logs WHERE tenant_id = $1`, [owner.tenantId])).rejects.toThrow(/append-only/);
  });

  it('paginates with a cursor', async () => {
    const first = (await client(t.app, owner.token).get('/api/audit-logs?limit=2')).json();
    expect(first.data).toHaveLength(2);
    expect(first.nextCursor).toBeTruthy();
    const second = (await client(t.app, owner.token).get(`/api/audit-logs?limit=2&cursor=${encodeURIComponent(first.nextCursor)}`)).json();
    expect(second.data[0].id).not.toBe(first.data[1].id);
  });
});
