import { assertWithinLimit } from '../subscriptions/service.js';
import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import type { Db } from '../../db/tx.js';
import { badRequest, conflict, forbidden, notFound } from '../../lib/errors.js';
import { hashPassword } from '../../lib/password.js';
import { parse, password, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import { writeAudit } from '../audit/audit.service.js';
import { assertCanGrant, resolveRoles } from '../rbac/rbac.routes.js';

type MemberRow = {
  id: string; email: string; full_name: string; user_status: string; member_status: string;
  is_owner: boolean; last_login_at: Date | null; joined_at: Date;
  roles: { id: string; code: string; nameAr: string; nameEn: string }[];
};

const MEMBER_SELECT = `
  SELECT u.id, u.email, u.full_name, u.status AS user_status, ut.status AS member_status,
         ut.is_owner, u.last_login_at, ut.joined_at,
         COALESCE(json_agg(json_build_object('id', r.id, 'code', r.code, 'nameAr', r.name_ar, 'nameEn', r.name_en)
                  ORDER BY r.code) FILTER (WHERE r.id IS NOT NULL), '[]') AS roles
    FROM user_tenants ut
    JOIN users u ON u.id = ut.user_id
    LEFT JOIN user_roles ur ON ur.tenant_id = ut.tenant_id AND ur.user_id = ut.user_id
    LEFT JOIN roles r ON r.id = ur.role_id AND r.deleted_at IS NULL
   WHERE ut.tenant_id = $1 AND u.deleted_at IS NULL`;

/** An invited person's name and activity stay private until they accept. */
const toMember = (m: MemberRow) => {
  const invited = m.member_status === 'INVITED';
  return {
    id: m.id, email: m.email, fullName: invited ? null : m.full_name,
    status: invited ? 'INVITED' : m.member_status === 'ACTIVE' && m.user_status === 'ACTIVE' ? 'ACTIVE' : 'DISABLED',
    isOwner: m.is_owner, lastLoginAt: invited ? null : m.last_login_at, joinedAt: m.joined_at, roles: m.roles,
  };
};

async function loadMember(db: Db, tenantId: string, userId: string) {
  const { rows } = await db.query<MemberRow>(`${MEMBER_SELECT} AND u.id = $2 GROUP BY u.id, ut.id`, [tenantId, userId]);
  return rows[0];
}

const createBody = z.object({
  email: z.string().trim().toLowerCase().pipe(z.email().max(254)),
  fullName: z.string().trim().min(2).max(200),
  // Initial password for a new account. The user must change it at first login.
  // TODO(Notification System): replace with an e-mailed invitation link.
  initialPassword: password.optional(),
  roleIds: z.array(z.uuid()).min(1).max(20),
});

export default async function usersRoutes(app: FastifyInstance) {
  app.get('/users', { preHandler: requirePermission(app, 'user.view') }, async (req) =>
    req.tenantTx(async (db) => {
      const { rows } = await db.query<MemberRow>(`${MEMBER_SELECT} GROUP BY u.id, ut.id ORDER BY ut.joined_at`, [req.auth!.tenantId]);
      return { data: rows.map(toMember) };
    }));

  app.get('/users/:id', { preHandler: requirePermission(app, 'user.view') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const member = await req.tenantTx((db) => loadMember(db, req.auth!.tenantId, id));
    if (!member) throw notFound('User');
    return toMember(member);
  });

  /**
   * Adds a user to the tenant. A new e-mail gets an account (with the initial
   * password); an existing account gets an invitation it must accept, so no one
   * can be added to an organization without consent.
   */
  app.post('/users', { preHandler: requirePermission(app, 'user.invite', 'user.manage') }, async (req, reply) => {
    const body = parse(createBody, req.body);
    const a = req.auth!;
    const passwordHash = body.initialPassword ? await hashPassword(body.initialPassword) : null;
    const member = await req.tenantTx(async (db) => {
      await assertWithinLimit(db, a.tenantId, 'max_users');
      const { permissions } = await resolveRoles(db, body.roleIds);
      assertCanGrant(a, permissions);

      let { rows: [user] } = await db.query<{ id: string }>(`SELECT id FROM users WHERE email = $1 AND deleted_at IS NULL`, [body.email]);
      let created = false;
      if (!user) {
        if (!passwordHash) throw badRequest('INITIAL_PASSWORD_REQUIRED', 'initialPassword is required for a new account');
        ({ rows: [user] } = await db.query<{ id: string }>(
          `INSERT INTO users (email, password_hash, full_name, must_change_password) VALUES ($1, $2, $3, true) RETURNING id`,
          [body.email, passwordHash, body.fullName]));
        created = true;
      }
      const exists = await db.query(`SELECT 1 FROM user_tenants WHERE tenant_id = $1 AND user_id = $2`, [a.tenantId, user!.id]);
      if (exists.rowCount) throw conflict('ALREADY_MEMBER', 'This user is already a member of the organization');

      await db.query(`INSERT INTO user_tenants (tenant_id, user_id, status) VALUES ($1, $2, $3)`, [a.tenantId, user!.id, created ? 'ACTIVE' : 'INVITED']);
      await db.query(
        `INSERT INTO user_roles (tenant_id, user_id, role_id, created_by) SELECT $1, $2, unnest($3::uuid[]), $4`,
        [a.tenantId, user!.id, [...new Set(body.roleIds)], a.userId]);
      const m = (await loadMember(db, a.tenantId, user!.id))!;
      await writeAudit(db, {
        tenantId: a.tenantId, userId: a.userId, action: 'CREATE', entityType: 'user', entityId: user!.id,
        newValues: { email: body.email, accountCreated: created, roles: m.roles.map((r) => r.code) },
      }, req.auditMeta());
      return m;
    });
    reply.code(201);
    return toMember(member);
  });

  app.patch('/users/:id', { preHandler: requirePermission(app, 'user.manage') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ status: z.enum(['ACTIVE', 'DISABLED']) }), req.body);
    const a = req.auth!;
    if (id === a.userId) throw forbidden('You cannot change your own status');
    const member = await req.tenantTx(async (db) => {
      const before = await loadMember(db, a.tenantId, id);
      if (!before) throw notFound('User');
      if (before.is_owner) throw forbidden('The organization owner cannot be disabled');
      // Only the invited person can accept an invitation; it can be withdrawn by disabling it.
      if (before.member_status === 'INVITED' && body.status === 'ACTIVE') throw conflict('INVITATION_PENDING', 'The invitation has not been accepted yet');
      // You cannot act on someone who holds permissions you do not have.
      if (before.roles.length) assertCanGrant(a, (await resolveRoles(db, before.roles.map((r) => r.id))).permissions);
      if (body.status === 'ACTIVE' && before.member_status !== 'ACTIVE') await assertWithinLimit(db, a.tenantId, 'max_users');
      await db.query(`UPDATE user_tenants SET status = $3 WHERE tenant_id = $1 AND user_id = $2`, [a.tenantId, id, body.status]);
      if (body.status === 'DISABLED') {
        await db.query(`UPDATE user_sessions SET revoked_at = now() WHERE user_id = $1 AND tenant_id = $2 AND revoked_at IS NULL`, [id, a.tenantId]);
      }
      await writeAudit(db, {
        tenantId: a.tenantId, userId: a.userId, action: 'UPDATE', entityType: 'user', entityId: id,
        oldValues: { status: before.member_status }, newValues: { status: body.status },
      }, req.auditMeta());
      return (await loadMember(db, a.tenantId, id))!;
    });
    return toMember(member);
  });

  app.put('/users/:id/roles', { preHandler: requirePermission(app, 'user.manage') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(z.object({ roleIds: z.array(z.uuid()).min(1).max(20) }), req.body);
    const a = req.auth!;
    if (id === a.userId) throw forbidden('You cannot change your own roles');
    const member = await req.tenantTx(async (db) => {
      const before = await loadMember(db, a.tenantId, id);
      if (!before) throw notFound('User');
      const { roles, permissions } = await resolveRoles(db, body.roleIds);
      assertCanGrant(a, permissions);
      // Removing a role is also a privilege change: you must hold everything it grants.
      const removed = before.roles.filter((r) => !body.roleIds.includes(r.id)).map((r) => r.id);
      if (removed.length) assertCanGrant(a, (await resolveRoles(db, removed)).permissions);
      if (before.is_owner && !roles.some((r) => r.code === 'TENANT_OWNER')) {
        throw forbidden('The organization owner must keep the TENANT_OWNER role');
      }
      await db.query(`DELETE FROM user_roles WHERE tenant_id = $1 AND user_id = $2`, [a.tenantId, id]);
      await db.query(
        `INSERT INTO user_roles (tenant_id, user_id, role_id, created_by) SELECT $1, $2, unnest($3::uuid[]), $4`,
        [a.tenantId, id, [...new Set(body.roleIds)], a.userId]);
      const after = (await loadMember(db, a.tenantId, id))!;
      await writeAudit(db, {
        tenantId: a.tenantId, userId: a.userId, action: 'PERMISSION_CHANGE', entityType: 'user', entityId: id,
        oldValues: { roles: before.roles.map((r) => r.code) }, newValues: { roles: after.roles.map((r) => r.code) },
      }, req.auditMeta());
      return after;
    });
    return toMember(member);
  });
}
