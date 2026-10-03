import type { FastifyInstance } from 'fastify';
import { z } from 'zod';
import type { Db } from '../../db/tx.js';
import { badRequest, conflict, forbidden, notFound } from '../../lib/errors.js';
import { parse, uuidParam } from '../../lib/validation.js';
import { requirePermission } from '../../plugins/auth.js';
import type { AuthContext } from '../../types.js';
import { writeAudit } from '../audit/audit.service.js';

/**
 * Resolves role IDs visible to the current tenant and the union of their
 * permissions. Unknown, deleted, or other-tenant roles are rejected.
 */
export async function resolveRoles(db: Db, roleIds: string[]) {
  const unique = [...new Set(roleIds)];
  if (unique.length === 0) return { roles: [], permissions: new Set<string>() };
  const { rows } = await db.query<{ id: string; code: string; perms: string[] }>(
    `SELECT r.id, r.code, COALESCE(array_agg(p.code) FILTER (WHERE p.code IS NOT NULL), '{}') AS perms
       FROM roles r
       LEFT JOIN role_permissions rp ON rp.role_id = r.id
       LEFT JOIN permissions p ON p.id = rp.permission_id
      WHERE r.id = ANY($1::uuid[]) AND r.deleted_at IS NULL
      GROUP BY r.id, r.code`,
    [unique],
  );
  if (rows.length !== unique.length) throw badRequest('INVALID_ROLE', 'One or more roles do not exist');
  return { roles: rows, permissions: new Set(rows.flatMap((r) => r.perms)) };
}

/** Prevents privilege escalation: nobody can grant a permission they do not hold. */
export function assertCanGrant(actor: AuthContext, permissions: Iterable<string>) {
  const missing = [...permissions].filter((p) => !actor.permissions.has(p));
  if (missing.length) throw forbidden(`You cannot grant permissions you do not have: ${missing.join(', ')}`);
}

async function validPermissionIds(db: Db, codes: string[]) {
  const { rows } = await db.query<{ id: string; code: string }>(
    `SELECT id, code FROM permissions WHERE code = ANY($1::text[])`, [codes]);
  const known = new Set(rows.map((r) => r.code));
  const unknown = codes.filter((c) => !known.has(c));
  if (unknown.length) throw badRequest('INVALID_PERMISSION', `Unknown permissions: ${unknown.join(', ')}`);
  return rows.map((r) => r.id);
}

const roleBody = z.object({
  code: z.string().trim().toUpperCase().regex(/^[A-Z][A-Z0-9_]{1,62}$/),
  nameAr: z.string().trim().min(2).max(100),
  nameEn: z.string().trim().min(2).max(100),
  description: z.string().trim().max(500).nullish(),
  permissions: z.array(z.string()).max(500),
});

type RoleRow = { id: string; code: string; name_ar: string; name_en: string; description: string | null; is_system: boolean; permissions: string[] };

const toRole = (r: RoleRow) => ({
  id: r.id, code: r.code, nameAr: r.name_ar, nameEn: r.name_en,
  description: r.description, isSystem: r.is_system, permissions: r.permissions,
});

async function loadRole(db: Db, id: string): Promise<RoleRow | undefined> {
  const { rows } = await db.query<RoleRow>(
    `SELECT r.id, r.code, r.name_ar, r.name_en, r.description, r.is_system,
            COALESCE(array_agg(p.code ORDER BY p.code) FILTER (WHERE p.code IS NOT NULL), '{}') AS permissions
       FROM roles r
       LEFT JOIN role_permissions rp ON rp.role_id = r.id
       LEFT JOIN permissions p ON p.id = rp.permission_id
      WHERE r.id = $1 AND r.deleted_at IS NULL
      GROUP BY r.id`,
    [id],
  );
  return rows[0];
}

export default async function rbacRoutes(app: FastifyInstance) {
  app.get('/permissions', { preHandler: requirePermission(app, 'role.view') }, async (req) =>
    req.tenantTx(async (db) => {
      const { rows } = await db.query(
        `SELECT code, module, description_ar AS "descriptionAr", description_en AS "descriptionEn"
           FROM permissions ORDER BY module, code`);
      return { data: rows };
    }));

  app.get('/roles', { preHandler: requirePermission(app, 'role.view') }, async (req) =>
    req.tenantTx(async (db) => {
      const { rows } = await db.query<RoleRow>(
        `SELECT r.id, r.code, r.name_ar, r.name_en, r.description, r.is_system,
                COALESCE(array_agg(p.code ORDER BY p.code) FILTER (WHERE p.code IS NOT NULL), '{}') AS permissions
           FROM roles r
           LEFT JOIN role_permissions rp ON rp.role_id = r.id
           LEFT JOIN permissions p ON p.id = rp.permission_id
          WHERE r.deleted_at IS NULL AND (r.tenant_id IS NULL OR r.tenant_id = $1)
          GROUP BY r.id
          ORDER BY r.is_system DESC, r.code`,
        [req.auth!.tenantId]);
      return { data: rows.map(toRole) };
    }));

  app.post('/roles', { preHandler: requirePermission(app, 'role.manage') }, async (req, reply) => {
    const body = parse(roleBody, req.body);
    const a = req.auth!;
    const role = await req.tenantTx(async (db) => {
      const sys = await db.query(`SELECT 1 FROM roles WHERE tenant_id IS NULL AND code = $1`, [body.code]);
      if (sys.rowCount) throw conflict('ROLE_CODE_RESERVED', 'This code is reserved for a system role');
      const permIds = await validPermissionIds(db, body.permissions);
      assertCanGrant(a, body.permissions);
      const { rows: [r] } = await db.query<{ id: string }>(
        `INSERT INTO roles (tenant_id, code, name_ar, name_en, description, is_system)
         VALUES ($1, $2, $3, $4, $5, false) RETURNING id`,
        [a.tenantId, body.code, body.nameAr, body.nameEn, body.description ?? null]);
      await db.query(
        `INSERT INTO role_permissions (role_id, permission_id) SELECT $1, unnest($2::uuid[])`, [r!.id, permIds]);
      const created = (await loadRole(db, r!.id))!;
      await writeAudit(db, {
        tenantId: a.tenantId, userId: a.userId, action: 'PERMISSION_CHANGE', entityType: 'role', entityId: r!.id,
        newValues: toRole(created),
      }, req.auditMeta());
      return created;
    });
    reply.code(201);
    return toRole(role);
  });

  app.patch('/roles/:id', { preHandler: requirePermission(app, 'role.manage') }, async (req) => {
    const { id } = parse(uuidParam, req.params);
    const body = parse(roleBody.omit({ code: true }).partial(), req.body);
    const a = req.auth!;
    const role = await req.tenantTx(async (db) => {
      const before = await loadRole(db, id);
      if (!before) throw notFound('Role');
      if (before.is_system) throw forbidden('System roles cannot be modified');
      // Editing a role you hold could strip your own access mid-session; allowed,
      // but you also cannot remove permissions you are not allowed to grant.
      assertCanGrant(a, before.permissions);
      await db.query(
        `UPDATE roles SET name_ar = COALESCE($2, name_ar), name_en = COALESCE($3, name_en),
                          description = CASE WHEN $4::boolean THEN $5 ELSE description END
          WHERE id = $1`,
        [id, body.nameAr ?? null, body.nameEn ?? null, body.description !== undefined, body.description ?? null]);
      if (body.permissions) {
        const permIds = await validPermissionIds(db, body.permissions);
        assertCanGrant(a, body.permissions);
        await db.query(`DELETE FROM role_permissions WHERE role_id = $1`, [id]);
        await db.query(`INSERT INTO role_permissions (role_id, permission_id) SELECT $1, unnest($2::uuid[])`, [id, permIds]);
      }
      const after = (await loadRole(db, id))!;
      await writeAudit(db, {
        tenantId: a.tenantId, userId: a.userId, action: 'PERMISSION_CHANGE', entityType: 'role', entityId: id,
        oldValues: toRole(before), newValues: toRole(after),
      }, req.auditMeta());
      return after;
    });
    return toRole(role);
  });

  app.delete('/roles/:id', { preHandler: requirePermission(app, 'role.manage') }, async (req, reply) => {
    const { id } = parse(uuidParam, req.params);
    const a = req.auth!;
    await req.tenantTx(async (db) => {
      const before = await loadRole(db, id);
      if (!before) throw notFound('Role');
      if (before.is_system) throw forbidden('System roles cannot be deleted');
      const used = await db.query(`SELECT 1 FROM user_roles WHERE role_id = $1 LIMIT 1`, [id]);
      if (used.rowCount) throw conflict('ROLE_IN_USE', 'Remove this role from all users before deleting it');
      await db.query(`UPDATE roles SET deleted_at = now() WHERE id = $1`, [id]);
      await writeAudit(db, {
        tenantId: a.tenantId, userId: a.userId, action: 'DELETE', entityType: 'role', entityId: id, oldValues: toRole(before),
      }, req.auditMeta());
    });
    return reply.code(204).send();
  });
}
