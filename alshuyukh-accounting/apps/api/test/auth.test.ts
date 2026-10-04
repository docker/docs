import { randomUUID } from 'node:crypto';
import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { PASSWORD, client, login, refreshCookieOf, register, setupApp, type TestContext } from './helpers.js';

let t: TestContext;
beforeAll(async () => { t = await setupApp(); });
afterAll(async () => { await t.close(); });

const refresh = (cookie: string, csrf = true) =>
  t.app.inject({
    method: 'POST', url: '/api/auth/refresh',
    cookies: { ash_rt: cookie },
    headers: csrf ? { 'x-csrf-protection': '1' } : {},
  });

describe('registration', () => {
  it('creates the user, tenant, settings, default company, main branch and warehouse', async () => {
    const s = await register(t.app, { companyName: 'شركة الشيوخ التجارية', vatNumber: '300000000000003' });
    const api = client(t.app, s.token);

    const me = (await api.get('/api/auth/me')).json();
    expect(me.tenant.id).toBe(s.tenantId);
    expect(me.tenant.isOwner).toBe(true);
    expect(me.roles.map((r: { code: string }) => r.code)).toEqual(['TENANT_OWNER']);
    expect(me.permissions).toContain('subscription.manage');

    const companies = (await api.get('/api/companies')).json().data;
    expect(companies).toHaveLength(1);
    expect(companies[0]).toMatchObject({ name: 'شركة الشيوخ التجارية', currency: 'SAR', timezone: 'Asia/Riyadh', country: 'SA', vatNumber: '300000000000003' });

    const branches = (await api.get(`/api/companies/${companies[0].id}/branches`)).json().data;
    expect(branches).toEqual([expect.objectContaining({ code: 'MAIN', isMain: true })]);
    const warehouses = (await api.get(`/api/companies/${companies[0].id}/warehouses`)).json().data;
    expect(warehouses).toEqual([expect.objectContaining({ code: 'MAIN', branchId: branches[0].id })]);

    const settings = (await api.get('/api/settings/tenant')).json();
    expect(settings).toMatchObject({ defaultCurrency: 'SAR', timezone: 'Asia/Riyadh', locale: 'ar' });
  });

  it('stores an argon2id hash, never the plain password', async () => {
    const s = await register(t.app);
    const { rows: [u] } = await t.ownerPool.query('SELECT password_hash FROM users WHERE id = $1', [s.userId]);
    expect(u.password_hash).toMatch(/^\$argon2id\$/);
    expect(u.password_hash).not.toContain(PASSWORD);
  });

  it('rejects a duplicate e-mail, weak password and malformed VAT number', async () => {
    const s = await register(t.app);
    const base = { fullName: 'X Y', tenantName: 'Org', companyName: 'Co' };
    const dup = await t.app.inject({ method: 'POST', url: '/api/auth/register', payload: { ...base, email: s.email, password: PASSWORD } });
    expect(dup.statusCode).toBe(409);
    const weak = await t.app.inject({ method: 'POST', url: '/api/auth/register', payload: { ...base, email: `w-${randomUUID()}@example.test`, password: 'short' } });
    expect(weak.statusCode).toBe(400);
    expect(weak.json().error.code).toBe('VALIDATION_ERROR');
    const vat = await t.app.inject({ method: 'POST', url: '/api/auth/register', payload: { ...base, email: `v-${randomUUID()}@example.test`, password: PASSWORD, vatNumber: '123' } });
    expect(vat.statusCode).toBe(400);
  });
});

describe('login', () => {
  it('logs in with valid credentials (e-mail is case-insensitive)', async () => {
    const s = await register(t.app);
    const l = await login(t.app, s.email.toUpperCase());
    expect(l.tenantId).toBe(s.tenantId);
    expect((await client(t.app, l.token).get('/api/auth/me')).statusCode).toBe(200);
  });

  it('returns the same 401 for a wrong password and an unknown e-mail', async () => {
    const s = await register(t.app);
    const wrong = await t.app.inject({ method: 'POST', url: '/api/auth/login', payload: { email: s.email, password: 'Wrong-passw0rd' } });
    const unknown = await t.app.inject({ method: 'POST', url: '/api/auth/login', payload: { email: `nobody-${randomUUID()}@example.test`, password: 'Wrong-passw0rd' } });
    expect(wrong.statusCode).toBe(401);
    expect(unknown.statusCode).toBe(401);
    expect(wrong.json().error.message).toBe(unknown.json().error.message);
  });

  it('locks sign-in from an address after repeated failures, without revealing it', async () => {
    const s = await register(t.app);
    const login = (password: string, remoteAddress: string) =>
      t.app.inject({ method: 'POST', url: '/api/auth/login', remoteAddress, payload: { email: s.email, password } });
    for (let i = 0; i < 5; i++) await login('Wrong-passw0rd', '203.0.113.7');
    // Same answer as a wrong password: an attacker cannot tell the account exists or is locked.
    const locked = await login(PASSWORD, '203.0.113.7');
    expect(locked.statusCode).toBe(401);
    expect(locked.json().error.code).toBe('UNAUTHORIZED');
    // Failures from one address do not lock the owner out everywhere.
    expect((await login(PASSWORD, '198.51.100.20')).statusCode).toBe(200);
  });

  it('ignores X-Forwarded-For unless a proxy is trusted', async () => {
    const s = await register(t.app);
    for (let i = 0; i < 5; i++) {
      await t.app.inject({ method: 'POST', url: '/api/auth/login', remoteAddress: '203.0.113.8',
        headers: { 'x-forwarded-for': `10.0.0.${i}` }, payload: { email: s.email, password: 'Wrong-passw0rd' } });
    }
    const res = await t.app.inject({ method: 'POST', url: '/api/auth/login', remoteAddress: '203.0.113.8',
      headers: { 'x-forwarded-for': '10.0.0.99' }, payload: { email: s.email, password: PASSWORD } });
    expect(res.statusCode).toBe(401);
  });

  it('rejects requests without a token or with a tampered token', async () => {
    expect((await t.app.inject({ method: 'GET', url: '/api/companies' })).statusCode).toBe(401);
    const s = await register(t.app);
    const tampered = s.token.slice(0, -2) + (s.token.endsWith('AA') ? 'BB' : 'AA');
    expect((await client(t.app, tampered).get('/api/companies')).statusCode).toBe(401);
  });
});

describe('sessions', () => {
  it('rotates the refresh token', async () => {
    const s = await register(t.app);
    const r1 = await refresh(s.refreshCookie);
    expect(r1.statusCode).toBe(200);
    const next = refreshCookieOf(r1);
    expect(next).not.toBe(s.refreshCookie);
    expect((await client(t.app, r1.json().accessToken).get('/api/auth/me')).statusCode).toBe(200);
  });

  it('requires the CSRF header on refresh', async () => {
    const s = await register(t.app);
    expect((await refresh(s.refreshCookie, false)).statusCode).toBe(400);
  });

  it('revokes all sessions when a rotated refresh token is reused', async () => {
    const s = await register(t.app);
    const r1 = await refresh(s.refreshCookie);
    const newAccess = r1.json().accessToken;
    const replay = await refresh(s.refreshCookie);
    expect(replay.statusCode).toBe(401);
    expect((await refresh(refreshCookieOf(r1))).statusCode).toBe(401);
    expect((await client(t.app, newAccess).get('/api/auth/me')).statusCode).toBe(401);
    expect((await client(t.app, s.token).get('/api/auth/me')).statusCode).toBe(401);
  });

  it('logout invalidates the access token immediately', async () => {
    const s = await register(t.app);
    const out = await t.app.inject({ method: 'POST', url: '/api/auth/logout', headers: { authorization: `Bearer ${s.token}`, 'x-csrf-protection': '1' } });
    expect(out.statusCode).toBe(204);
    expect((await client(t.app, s.token).get('/api/auth/me')).statusCode).toBe(401);
    expect((await refresh(s.refreshCookie)).statusCode).toBe(401);
  });

  it('changing the password signs out other sessions', async () => {
    const s = await register(t.app);
    const other = await login(t.app, s.email);
    const res = await client(t.app, s.token).post('/api/auth/change-password', { currentPassword: PASSWORD, newPassword: 'An0ther-Strong-Pass' });
    expect(res.statusCode).toBe(204);
    expect((await client(t.app, other.token).get('/api/auth/me')).statusCode).toBe(401);
    expect((await client(t.app, s.token).get('/api/auth/me')).statusCode).toBe(200);
    await login(t.app, s.email, 'An0ther-Strong-Pass');
  });
});

describe('security headers', () => {
  it('sets secure headers', async () => {
    const res = await t.app.inject({ method: 'GET', url: '/api/health' });
    expect(res.statusCode).toBe(200);
    expect(res.headers['x-content-type-options']).toBe('nosniff');
    expect(res.headers['content-security-policy']).toContain("default-src 'none'");
    expect(res.headers['strict-transport-security']).toBeDefined();
  });
});
