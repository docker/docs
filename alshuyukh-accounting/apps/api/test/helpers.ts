import { randomUUID } from 'node:crypto';
import type { FastifyInstance, LightMyRequestResponse } from 'fastify';
import pg from 'pg';
import { buildApp } from '../src/app.js';
import { loadEnv } from '../src/config/env.js';
import { createPool } from '../src/db/pool.js';

export interface TestContext {
  app: FastifyInstance;
  pool: pg.Pool;        // restricted app role (RLS applies)
  ownerPool: pg.Pool;   // schema owner (bypasses RLS) — only for assertions
  close(): Promise<void>;
}

export async function setupApp(): Promise<TestContext> {
  const env = loadEnv();
  const pool = createPool(env.DATABASE_URL);
  const ownerPool = createPool(process.env.DATABASE_URL_MIGRATE!);
  const app = await buildApp({ env, pool, logger: false });
  await app.ready();
  return { app, pool, ownerPool, close: async () => { await app.close(); await pool.end(); await ownerPool.end(); } };
}

export const PASSWORD = 'Str0ng-Passw0rd!';

export interface Session {
  token: string;
  refreshCookie: string;
  tenantId: string;
  userId: string;
  email: string;
}

export function refreshCookieOf(res: LightMyRequestResponse): string {
  const c = res.cookies.find((x) => x.name === 'ash_rt');
  if (!c) throw new Error('refresh cookie not set');
  return c.value;
}

export async function register(app: FastifyInstance, overrides: Record<string, unknown> = {}): Promise<Session> {
  const email = `owner-${randomUUID()}@example.test`;
  const res = await app.inject({
    method: 'POST', url: '/api/auth/register',
    payload: { fullName: 'Test Owner', email, password: PASSWORD, tenantName: 'Test Org', companyName: `Company ${randomUUID().slice(0, 8)}`, ...overrides },
  });
  if (res.statusCode !== 201) throw new Error(`register failed: ${res.statusCode} ${res.body}`);
  const body = res.json();
  return { token: body.accessToken, refreshCookie: refreshCookieOf(res), tenantId: body.tenantId, userId: body.userId, email };
}

export async function login(app: FastifyInstance, email: string, password = PASSWORD, tenantId?: string): Promise<Session> {
  const res = await app.inject({ method: 'POST', url: '/api/auth/login', payload: { email, password, tenantId } });
  if (res.statusCode !== 200) throw new Error(`login failed: ${res.statusCode} ${res.body}`);
  const body = res.json();
  return { token: body.accessToken, refreshCookie: refreshCookieOf(res), tenantId: body.tenantId, userId: body.userId, email };
}

/** Small request helper with a bearer token. */
export function client(app: FastifyInstance, token: string) {
  const call = (method: 'GET' | 'POST' | 'PATCH' | 'PUT' | 'DELETE') => (url: string, payload?: unknown) =>
    app.inject({ method, url, payload: payload as never, headers: { authorization: `Bearer ${token}` } });
  return { get: call('GET'), post: call('POST'), patch: call('PATCH'), put: call('PUT'), del: call('DELETE') };
}

export async function roleId(app: FastifyInstance, token: string, code: string): Promise<string> {
  const res = await client(app, token).get('/api/roles');
  const role = res.json().data.find((r: { code: string }) => r.code === code);
  if (!role) throw new Error(`role ${code} not found`);
  return role.id;
}

/** Adds a member with the given system role and logs them in. */
export async function addMember(app: FastifyInstance, owner: Session, roleCode: string): Promise<Session> {
  const email = `${roleCode.toLowerCase()}-${randomUUID()}@example.test`;
  const res = await client(app, owner.token).post('/api/users', {
    email, fullName: `Member ${roleCode}`, initialPassword: PASSWORD, roleIds: [await roleId(app, owner.token, roleCode)],
  });
  if (res.statusCode !== 201) throw new Error(`addMember failed: ${res.statusCode} ${res.body}`);
  return login(app, email, PASSWORD, owner.tenantId);
}

// Accounting helpers ------------------------------------------------------------

export interface Chart { byCode: Map<string, string>; byKey: Map<string, string> }

export async function chart(app: FastifyInstance, token: string): Promise<Chart> {
  const res = await client(app, token).get('/api/accounts?includeInactive=true');
  if (res.statusCode !== 200) throw new Error(`accounts failed: ${res.body}`);
  const rows = res.json().data as { id: string; code: string; systemKey: string | null }[];
  return {
    byCode: new Map(rows.map((r) => [r.code, r.id])),
    byKey: new Map(rows.filter((r) => r.systemKey).map((r) => [r.systemKey!, r.id])),
  };
}

export interface FiscalYear {
  id: string; name: string; startDate: string; endDate: string; status: string;
  periods: { id: string; periodNumber: number; startDate: string; endDate: string; status: string }[];
}

export async function fiscalYears(app: FastifyInstance, token: string): Promise<FiscalYear[]> {
  return (await client(app, token).get('/api/fiscal-years')).json().data;
}

/** A date inside the first fiscal year, `month` 1-12 relative to its start. */
export async function dateInYear(app: FastifyInstance, token: string, month = 3, day = 15): Promise<string> {
  const [year] = await fiscalYears(app, token);
  return year!.periods[month - 1]!.startDate.slice(0, 8) + String(day).padStart(2, '0');
}

export const dr = (accountId: string, amount: string) => ({ accountId, debit: amount });
export const cr = (accountId: string, amount: string) => ({ accountId, credit: amount });
