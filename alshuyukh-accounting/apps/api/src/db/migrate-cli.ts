import { migrate } from './migrate.js';

const url = process.env.DATABASE_URL_MIGRATE;
if (!url) {
  console.error('DATABASE_URL_MIGRATE is not set');
  process.exit(1);
}

migrate(url, (msg) => console.log(msg))
  .then(() => console.log('migrations complete'))
  .catch((err: Error) => {
    console.error(err.message);
    process.exit(1);
  });
