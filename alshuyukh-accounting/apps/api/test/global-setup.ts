import pg from 'pg';
import { migrate } from '../src/db/migrate.js';
import './env.js';

/** Rebuilds the test database schema from scratch, then applies all migrations. */
export default async function setup() {
  const url = process.env.TEST_DATABASE_URL_MIGRATE!;
  if (!/_test\b/.test(new URL(url).pathname)) throw new Error('Refusing to reset a database whose name does not end in _test');
  const client = new pg.Client({ connectionString: url });
  await client.connect();
  await client.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  await client.end();
  await migrate(url);
  // Module tests should not depend on the trial plan's limits; subscription tests set their own.
  const c2 = new pg.Client({ connectionString: url });
  await c2.connect();
  await c2.query(`UPDATE plans SET max_users = NULL, max_companies = NULL, max_branches = NULL, max_warehouses = NULL, max_products = NULL,
    max_invoices_per_month = NULL, max_storage_mb = NULL, max_api_calls_per_month = NULL WHERE code = 'TRIAL'`);
  await c2.end();
}
