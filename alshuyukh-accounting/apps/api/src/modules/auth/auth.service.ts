import { startSubscription } from '../subscriptions/service.js';
import { randomUUID } from 'node:crypto';
import { withTx, type Db } from '../../db/tx.js';
import { AppError, conflict, forbidden, unauthorized } from '../../lib/errors.js';
import { getDummyHash, hashPassword, verifyPassword } from '../../lib/password.js';
import { hashToken, newRefreshToken } from '../../lib/tokens.js';
import type { Deps } from '../../types.js';
import { setupCompanyAccounting } from '../accounting/setup.js';
import { writeAudit, type AuditMeta } from '../audit/audit.service.js';

export interface Membership {
  tenantId: string;
  tenantName: string;
  tenantStatus: string;
  isOwner: boolean;
}

export interface IssuedTokens {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  tenantId: string;
  userId: string;
}

export interface RegisterInput {
  fullName: string;
  email: string;
  password: string;
  tenantName: string;
  companyName: string;
  vatNumber?: string | null;
  commercialRegistration?: string | null;
}

export class AuthService {
  constructor(private readonly deps: Deps) {}

  /** Self-service sign-up: creates the user, tenant, default company, main branch and warehouse. */
  async register(input: RegisterInput, meta: AuditMeta): Promise<IssuedTokens> {
    const passwordHash = await hashPassword(input.password);
    const tenantId = randomUUID();
    const slug = `org-${tenantId.slice(0, 8)}`;

    const userId = await withTx(this.deps.pool, { tenantId, userId: null }, (db) =>
      provisionTenant(db, { ...input, tenantId, passwordHash, slug }, meta));

    return this.issueSession(userId, tenantId, meta);
  }

  async login(email: string, password: string, tenantId: string | undefined, meta: AuditMeta): Promise<IssuedTokens> {
    const { env, pool } = this.deps;

    const ip = meta.ip ?? 'unknown';
    const user = await withTx(pool, {}, async (db) => {
      const { rows } = await db.query<{
        id: string; password_hash: string; status: string; deleted_at: Date | null; locked_until: Date | null;
      }>(
        `SELECT u.id, u.password_hash, u.status, u.deleted_at,
                (SELECT f.locked_until FROM login_failures f WHERE f.user_id = u.id AND f.ip = $2) AS locked_until
           FROM users u WHERE u.email = $1`,
        [email, ip],
      );
      return rows[0];
    });

    if (!user) {
      await verifyPassword(await getDummyHash(), password);
      await this.auditFailedLogin(null, email, 'unknown_email', meta);
      throw unauthorized('Invalid e-mail or password');
    }
    // The password is always checked (same timing) and every failure looks the
    // same, so a lock reveals neither that the account exists nor that it is locked.
    const ok = await verifyPassword(user.password_hash, password);
    if (user.locked_until && user.locked_until > new Date()) {
      await this.auditFailedLogin(user.id, email, 'locked', meta);
      throw unauthorized('Invalid e-mail or password');
    }
    if (!ok) {
      // Locks only this account from this address; the owner can still sign in from elsewhere.
      await withTx(pool, { userId: user.id }, async (db) => {
        await db.query(
          `INSERT INTO login_failures (user_id, ip, failures, locked_until) VALUES ($1, $2, 1, NULL)
           ON CONFLICT (user_id, ip) DO UPDATE SET
             failures = CASE WHEN login_failures.locked_until < now() THEN 1 ELSE login_failures.failures + 1 END,
             locked_until = CASE WHEN (CASE WHEN login_failures.locked_until < now() THEN 1 ELSE login_failures.failures + 1 END) >= $3
                                 THEN now() + make_interval(mins => $4) ELSE NULL END,
             updated_at = now()`,
          [user.id, ip, env.LOGIN_MAX_ATTEMPTS, env.LOGIN_LOCK_MINUTES],
        );
      });
      await this.auditFailedLogin(user.id, email, 'bad_password', meta);
      throw unauthorized('Invalid e-mail or password');
    }
    if (user.status !== 'ACTIVE' || user.deleted_at) {
      await this.auditFailedLogin(user.id, email, 'disabled', meta);
      throw unauthorized('Account is disabled');
    }

    const memberships = await this.listMemberships(user.id);
    const active = memberships.filter((m) => m.tenantStatus === 'ACTIVE');
    const target = tenantId ? active.find((m) => m.tenantId === tenantId) : active[0];
    if (!target) {
      await this.auditFailedLogin(user.id, email, 'no_active_membership', meta);
      throw forbidden('No active organization is available for this account');
    }

    await withTx(pool, { userId: user.id }, (db) =>
      db.query(`WITH c AS (DELETE FROM login_failures WHERE user_id = $1 AND ip = $2) UPDATE users SET last_login_at = now() WHERE id = $1`, [user.id, ip]),
    );
    return this.issueSession(user.id, target.tenantId, meta, 'LOGIN');
  }

  /**
   * Rotates the refresh token. Presenting an already-rotated token means it
   * was stolen or replayed, so every session of that user is revoked.
   */
  async refresh(refreshToken: string, meta: AuditMeta): Promise<IssuedTokens> {
    const tokenHash = hashToken(refreshToken);
    const { pool, env, tokens } = this.deps;
    const newToken = newRefreshToken();

    const result = await withTx(pool, {}, async (db) => {
      const { rows: [s] } = await db.query<{
        id: string; user_id: string; tenant_id: string; expires_at: Date; revoked_at: Date | null; replaced_by: string | null;
      }>(
        `SELECT id, user_id, tenant_id, expires_at, revoked_at, replaced_by
           FROM user_sessions WHERE token_hash = $1 FOR UPDATE`,
        [tokenHash],
      );
      if (!s) throw unauthorized('Invalid refresh token');
      if (s.revoked_at) {
        // Handled after this transaction, because throwing here would roll back the revocation.
        if (s.replaced_by) return { kind: 'reused', session: { id: s.id, user_id: s.user_id, tenant_id: s.tenant_id } } as const;
        throw unauthorized('Invalid refresh token');
      }
      if (s.expires_at < new Date()) throw unauthorized('Refresh token expired');

      await db.query(`SELECT set_config('app.tenant_id', $1, true), set_config('app.user_id', $2, true)`, [s.tenant_id, s.user_id]);
      await this.assertCanUseTenant(db, s.user_id, s.tenant_id);

      const { rows: [next] } = await db.query<{ id: string }>(
        `INSERT INTO user_sessions (user_id, tenant_id, token_hash, expires_at, ip_address, user_agent)
         VALUES ($1, $2, $3, now() + make_interval(days => $4), $5, $6) RETURNING id`,
        [s.user_id, s.tenant_id, hashToken(newToken), env.REFRESH_TOKEN_TTL_DAYS, meta.ip ?? null, meta.userAgent ?? null],
      );
      await db.query(
        `UPDATE user_sessions SET revoked_at = now(), replaced_by = $2, last_used_at = now() WHERE id = $1`,
        [s.id, next!.id],
      );
      return {
        kind: 'ok',
        tokens: {
          accessToken: await tokens.signAccess({ sub: s.user_id, tid: s.tenant_id, sid: next!.id }),
          refreshToken: newToken,
          expiresIn: env.ACCESS_TOKEN_TTL_SECONDS,
          tenantId: s.tenant_id,
          userId: s.user_id,
        },
      } as const;
    });

    if (result.kind === 'reused') {
      const s = result.session;
      await withTx(pool, { tenantId: s.tenant_id, userId: s.user_id }, async (db) => {
        await db.query(`UPDATE user_sessions SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL`, [s.user_id]);
        await writeAudit(db, { tenantId: s.tenant_id, userId: s.user_id, action: 'TOKEN_REUSE_DETECTED', entityType: 'session', entityId: s.id }, meta);
      });
      throw unauthorized('Invalid refresh token');
    }
    return result.tokens;
  }

  async logout(sessionId: string, userId: string, tenantId: string, meta: AuditMeta): Promise<void> {
    await withTx(this.deps.pool, { tenantId, userId }, async (db) => {
      await db.query(`UPDATE user_sessions SET revoked_at = now() WHERE id = $1 AND user_id = $2 AND revoked_at IS NULL`, [sessionId, userId]);
      await writeAudit(db, { tenantId, userId, action: 'LOGOUT', entityType: 'session', entityId: sessionId }, meta);
    });
  }

  /** Moves the current session to another tenant the user belongs to. */
  async switchTenant(sessionId: string, userId: string, fromTenantId: string, toTenantId: string, meta: AuditMeta) {
    const { pool, tokens, env } = this.deps;
    await withTx(pool, { tenantId: toTenantId, userId }, async (db) => {
      await this.assertCanUseTenant(db, userId, toTenantId);
      await db.query(`UPDATE user_sessions SET tenant_id = $3 WHERE id = $1 AND user_id = $2`, [sessionId, userId, toTenantId]);
      await writeAudit(db, {
        tenantId: toTenantId, userId, action: 'TENANT_SWITCH', entityType: 'tenant', entityId: toTenantId,
        oldValues: { tenantId: fromTenantId }, newValues: { tenantId: toTenantId },
      }, meta);
    });
    return {
      accessToken: await tokens.signAccess({ sub: userId, tid: toTenantId, sid: sessionId }),
      expiresIn: env.ACCESS_TOKEN_TTL_SECONDS,
      tenantId: toTenantId,
    };
  }

  async changePassword(userId: string, tenantId: string, sessionId: string, current: string, next: string, meta: AuditMeta) {
    const { pool } = this.deps;
    const hash = await withTx(pool, { userId }, async (db) =>
      (await db.query<{ password_hash: string }>('SELECT password_hash FROM users WHERE id = $1', [userId])).rows[0]?.password_hash);
    if (!hash || !(await verifyPassword(hash, current))) throw badCurrentPassword();
    if (current === next) throw new AppError(400, 'PASSWORD_REUSED', 'New password must differ from the current password');
    const newHash = await hashPassword(next);
    await withTx(pool, { tenantId, userId }, async (db) => {
      await db.query(
        `UPDATE users SET password_hash = $2, password_changed_at = now(), must_change_password = false WHERE id = $1`,
        [userId, newHash],
      );
      // Sign out every other device.
      await db.query(`UPDATE user_sessions SET revoked_at = now() WHERE user_id = $1 AND id <> $2 AND revoked_at IS NULL`, [userId, sessionId]);
      await writeAudit(db, { tenantId, userId, action: 'PASSWORD_CHANGE', entityType: 'user', entityId: userId }, meta);
    });
  }

  async listMemberships(userId: string): Promise<Membership[]> {
    return withTx(this.deps.pool, { userId }, async (db) => {
      const { rows } = await db.query<{ tenant_id: string; name: string; status: string; is_owner: boolean }>(
        `SELECT t.id AS tenant_id, t.name, t.status, ut.is_owner
           FROM user_tenants ut JOIN tenants t ON t.id = ut.tenant_id
          WHERE ut.user_id = $1 AND ut.status = 'ACTIVE' AND t.deleted_at IS NULL
          ORDER BY ut.joined_at`,
        [userId],
      );
      return rows.map((r) => ({ tenantId: r.tenant_id, tenantName: r.name, tenantStatus: r.status, isOwner: r.is_owner }));
    });
  }

  private async assertCanUseTenant(db: Db, userId: string, tenantId: string): Promise<void> {
    const { rows: [r] } = await db.query<{ member_status: string; tenant_status: string; user_status: string }>(
      `SELECT ut.status AS member_status, t.status AS tenant_status, u.status AS user_status
         FROM user_tenants ut JOIN tenants t ON t.id = ut.tenant_id JOIN users u ON u.id = ut.user_id
        WHERE ut.user_id = $1 AND ut.tenant_id = $2 AND u.deleted_at IS NULL`,
      [userId, tenantId],
    );
    // Not-a-member and no-such-tenant look identical to the caller.
    if (!r || r.member_status !== 'ACTIVE' || r.user_status !== 'ACTIVE') throw forbidden('No access to this organization');
    if (r.tenant_status !== 'ACTIVE') throw new AppError(403, 'TENANT_SUSPENDED', 'This organization is suspended');
  }

  private async issueSession(userId: string, tenantId: string, meta: AuditMeta, auditAction?: 'LOGIN'): Promise<IssuedTokens> {
    const { pool, env, tokens } = this.deps;
    const refreshToken = newRefreshToken();
    const sessionId = await withTx(pool, { tenantId, userId }, async (db) => {
      const { rows: [s] } = await db.query<{ id: string }>(
        `INSERT INTO user_sessions (user_id, tenant_id, token_hash, expires_at, ip_address, user_agent)
         VALUES ($1, $2, $3, now() + make_interval(days => $4), $5, $6) RETURNING id`,
        [userId, tenantId, hashToken(refreshToken), env.REFRESH_TOKEN_TTL_DAYS, meta.ip ?? null, meta.userAgent ?? null],
      );
      if (auditAction) {
        await writeAudit(db, { tenantId, userId, action: auditAction, entityType: 'session', entityId: s!.id }, meta);
      }
      return s!.id;
    });
    return {
      accessToken: await tokens.signAccess({ sub: userId, tid: tenantId, sid: sessionId }),
      refreshToken,
      expiresIn: env.ACCESS_TOKEN_TTL_SECONDS,
      tenantId,
      userId,
    };
  }

  private async auditFailedLogin(userId: string | null, email: string, reason: string, meta: AuditMeta) {
    await withTx(this.deps.pool, { userId }, (db) =>
      writeAudit(db, { tenantId: null, userId, action: 'LOGIN_FAILED', newValues: { email, reason } }, meta));
  }
}

const badCurrentPassword = () => new AppError(400, 'INVALID_CURRENT_PASSWORD', 'Current password is incorrect');

export interface ProvisionInput extends RegisterInput {
  tenantId: string; passwordHash: string; slug: string;
  /** Admin-created owners must change the temporary password at first sign-in. */
  mustChangePassword?: boolean;
  planId?: string | null;
  createdBy?: string | null;
}

/**
 * Creates the owner user, the tenant, its default company, main branch,
 * warehouse, chart of accounts and subscription. Used by self-service sign-up
 * and by the platform administrator. Runs in the caller's transaction.
 */
export async function provisionTenant(db: Db, input: ProvisionInput, meta: AuditMeta): Promise<string> {
  const { tenantId, passwordHash, slug } = input;
  // Lets RLS accept the organization's first subscription rows (see migration 0014).
  await db.query(`SELECT set_config('app.provisioning', 'on', true)`);
  const existing = await db.query('SELECT 1 FROM users WHERE email = $1', [input.email]);
  if (existing.rowCount) throw conflict('EMAIL_TAKEN', 'An account with this e-mail already exists');

  const { rows: [user] } = await db.query<{ id: string }>(
    `INSERT INTO users (email, password_hash, full_name, must_change_password) VALUES ($1, $2, $3, $4) RETURNING id`,
    [input.email, passwordHash, input.fullName, input.mustChangePassword ?? false],
  );
  const uid = user!.id;
  await db.query(`SELECT set_config('app.user_id', $1, true)`, [uid]);

  await db.query(`INSERT INTO tenants (id, name, slug) VALUES ($1, $2, $3)`, [tenantId, input.tenantName, slug]);
  await db.query(`INSERT INTO tenant_settings (tenant_id) VALUES ($1)`, [tenantId]);
  await db.query(
    `INSERT INTO user_tenants (tenant_id, user_id, is_owner) VALUES ($1, $2, true)`,
    [tenantId, uid],
  );
  await db.query(
    `INSERT INTO user_roles (tenant_id, user_id, role_id, created_by)
     SELECT $1, $2, id, $2 FROM roles WHERE code = 'TENANT_OWNER' AND tenant_id IS NULL`,
    [tenantId, uid],
  );
  const { rows: [company] } = await db.query<{ id: string }>(
    `INSERT INTO companies (tenant_id, name, legal_name, vat_number, commercial_registration, created_by)
     VALUES ($1, $2, $2, $3, $4, $5) RETURNING id`,
    [tenantId, input.companyName, input.vatNumber ?? null, input.commercialRegistration ?? null, uid],
  );
  const { rows: [branch] } = await db.query<{ id: string }>(
    `INSERT INTO branches (tenant_id, company_id, code, name, is_main, created_by)
     VALUES ($1, $2, 'MAIN', 'الفرع الرئيسي', true, $3) RETURNING id`,
    [tenantId, company!.id, uid],
  );
  await db.query(
    `INSERT INTO warehouses (tenant_id, company_id, branch_id, code, name, created_by)
     VALUES ($1, $2, $3, 'MAIN', 'المستودع الرئيسي', $4)`,
    [tenantId, company!.id, branch!.id, uid],
  );
  await setupCompanyAccounting(db, { tenantId, companyId: company!.id, userId: uid });
  // Every new organization starts on the default plan (a trial when the plan has trial days).
  await startSubscription(db, tenantId, input.createdBy ?? uid, input.planId ?? null);
  await writeAudit(db, {
    tenantId, userId: uid, action: 'REGISTER', entityType: 'tenant', entityId: tenantId,
    newValues: { tenantName: input.tenantName, companyName: input.companyName, email: input.email },
  }, meta);
  return uid;
}
