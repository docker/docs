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
}
