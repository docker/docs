/**
 * Grants or revokes platform administrator access (the /admin panel).
 * Runs with the migration (owner) connection, outside the API:
 *   npm run admin -w apps/api -- grant someone@example.com
 *   npm run admin -w apps/api -- revoke someone@example.com
 */
import pg from 'pg';

const [command, email] = process.argv.slice(2);
if (!['grant', 'revoke'].includes(command ?? '') || !email) {
  console.error('Usage: admin <grant|revoke> <email>');
  process.exit(1);
}
const url = process.env.DATABASE_URL_MIGRATE ?? process.env.DATABASE_URL;
if (!url) { console.error('DATABASE_URL_MIGRATE is not set'); process.exit(1); }
const client = new pg.Client({ connectionString: url });
await client.connect();
const { rowCount } = await client.query(`UPDATE users SET is_platform_admin = $2 WHERE email = $1 AND deleted_at IS NULL`, [email, command === 'grant']);
await client.end();
if (!rowCount) { console.error(`No user with e-mail ${email}`); process.exit(1); }
console.log(`${command === 'grant' ? 'Granted' : 'Revoked'} platform administrator access for ${email}`);
