import { existsSync } from 'node:fs';
import { resolve } from 'node:path';

const envFile = resolve(import.meta.dirname, '../../../.env');
if (existsSync(envFile)) process.loadEnvFile(envFile);

if (!process.env.TEST_DATABASE_URL || !process.env.TEST_DATABASE_URL_MIGRATE) {
  throw new Error('TEST_DATABASE_URL and TEST_DATABASE_URL_MIGRATE must be set (see .env.example)');
}
process.env.NODE_ENV = 'test';
process.env.DATABASE_URL = process.env.TEST_DATABASE_URL;
process.env.DATABASE_URL_MIGRATE = process.env.TEST_DATABASE_URL_MIGRATE;
process.env.JWT_SECRET ??= 'test-secret-that-is-long-enough-for-hs256-signing';
