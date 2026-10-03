import { createHash } from 'node:crypto';
import { readdir, readFile } from 'node:fs/promises';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import pg from 'pg';
import { PERMISSIONS, SYSTEM_ROLES, resolveRolePermissions } from '../modules/rbac/catalog.js';

const MIGRATIONS_DIR = join(dirname(fileURLToPath(import.meta.url)), 'migrations');

/**
 * Applies pending SQL migrations in file-name order, each in its own
 * transaction, then syncs the RBAC catalog. An already-applied migration
 * whose content changed aborts the run: write a new migration instead.
 */
export async function migrate(connectionString: string, log: (msg: string) => void = () => {}): Promise<void> {
  const client = new pg.Client({ connectionString, options: '-c timezone=UTC' });
  await client.connect();
  try {
    await client.query(`
      CREATE TABLE IF NOT EXISTS schema_migrations (
        version    text PRIMARY KEY,
        checksum   text NOT NULL,
        applied_at timestamptz NOT NULL DEFAULT now()
      )`);
    // Serialize concurrent migration runs.
    await client.query(`SELECT pg_advisory_lock(hashtext('alshuyukh_migrations'))`);

    const applied = new Map<string, string>(
      (await client.query<{ version: string; checksum: string }>('SELECT version, checksum FROM schema_migrations'))
        .rows.map((r) => [r.version, r.checksum]),
    );
    const files = (await readdir(MIGRATIONS_DIR)).filter((f) => f.endsWith('.sql')).sort();

    for (const file of files) {
      const sql = await readFile(join(MIGRATIONS_DIR, file), 'utf8');
      const checksum = createHash('sha256').update(sql).digest('hex');
      const existing = applied.get(file);
      if (existing) {
        if (existing !== checksum) {
          throw new Error(`Migration ${file} was modified after it was applied. Create a new migration instead.`);
        }
        continue;
      }
      await client.query('BEGIN');
      try {
        await client.query(sql);
        await client.query('INSERT INTO schema_migrations (version, checksum) VALUES ($1, $2)', [file, checksum]);
        await client.query('COMMIT');
        log(`applied ${file}`);
      } catch (err) {
        await client.query('ROLLBACK');
        throw new Error(`Migration ${file} failed: ${(err as Error).message}`);
      }
    }

    await syncRbacCatalog(client);
    log('rbac catalog synced');
  } finally {
    await client.query(`SELECT pg_advisory_unlock(hashtext('alshuyukh_migrations'))`).catch(() => undefined);
    await client.end();
  }
}

async function syncRbacCatalog(client: pg.Client): Promise<void> {
  await client.query('BEGIN');
  try {
    for (const p of PERMISSIONS) {
      await client.query(
        `INSERT INTO permissions (code, module, description_ar, description_en)
         VALUES ($1, $2, $3, $4)
         ON CONFLICT (code) DO UPDATE SET module = EXCLUDED.module,
           description_ar = EXCLUDED.description_ar, description_en = EXCLUDED.description_en`,
        [p.code, p.module, p.ar, p.en],
      );
    }
    for (const role of SYSTEM_ROLES) {
      const { rows } = await client.query<{ id: string }>(
        `INSERT INTO roles (code, name_ar, name_en, is_system, tenant_id)
         VALUES ($1, $2, $3, true, NULL)
         ON CONFLICT (code) WHERE tenant_id IS NULL
         DO UPDATE SET name_ar = EXCLUDED.name_ar, name_en = EXCLUDED.name_en
         RETURNING id`,
        [role.code, role.ar, role.en],
      );
      const roleId = rows[0]!.id;
      const codes = resolveRolePermissions(role);
      await client.query(
        `DELETE FROM role_permissions WHERE role_id = $1
           AND permission_id NOT IN (SELECT id FROM permissions WHERE code = ANY($2::text[]))`,
        [roleId, codes],
      );
      await client.query(
        `INSERT INTO role_permissions (role_id, permission_id)
         SELECT $1, id FROM permissions WHERE code = ANY($2::text[])
         ON CONFLICT DO NOTHING`,
        [roleId, codes],
      );
    }
    await client.query('COMMIT');
  } catch (err) {
    await client.query('ROLLBACK');
    throw err;
  }
}
