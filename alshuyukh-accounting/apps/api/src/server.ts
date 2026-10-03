import { buildApp } from './app.js';
import { loadEnv } from './config/env.js';
import { createPool } from './db/pool.js';

const env = loadEnv();
const pool = createPool(env.DATABASE_URL);
const app = await buildApp({ env, pool });

const shutdown = async () => {
  await app.close();
  await pool.end();
  process.exit(0);
};
process.on('SIGINT', shutdown);
process.on('SIGTERM', shutdown);

await app.listen({ port: env.PORT, host: env.HOST });
