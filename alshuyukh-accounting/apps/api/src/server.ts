import { buildApp } from './app.js';
import { loadEnv } from './config/env.js';
import { createPool } from './db/pool.js';
import { startZatcaWorker } from './modules/zatca/worker.js';

const env = loadEnv();
const pool = createPool(env.DATABASE_URL);
const app = await buildApp({ env, pool });

const stopWorker = env.ZATCA_WORKER === 'on' ? startZatcaWorker(pool, app.log, env.ZATCA_WORKER_INTERVAL_SECONDS) : () => undefined;

const shutdown = async () => {
  stopWorker();
  await app.close();
  await pool.end();
  process.exit(0);
};
process.on('SIGINT', shutdown);
process.on('SIGTERM', shutdown);

await app.listen({ port: env.PORT, host: env.HOST });
