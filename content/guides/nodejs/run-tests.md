---
title: Run Node.js tests in a container
linkTitle: Run your tests
weight: 30
keywords: node.js, node, test, vitest
description: Learn how to run your Node.js tests in a container.
aliases:
  - /language/nodejs/run-tests/
  - /guides/language/nodejs/run-tests/
---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Node.js application](containerize.md).

## Overview

Testing is a core part of building reliable software. Docker makes it easy to
run your tests in the same environment used in CI and production, so failures
are caught before they reach your users.

In this section, you'll add [Vitest](https://vitest.dev/) to the project and
run tests both locally and inside a container.

## Update the application

You'll refactor `src/index.ts` to export the Express `app` instance so tests
can import it without starting a server. Add a test file and update
`package.json` to add Vitest and a test runner for HTTP requests. The file browser shows only the files that change in this step.

{{< files name="nodejs-docker-example" >}}

{{< file path="src/index.ts" status="modified" hl_lines="10,31,70-75" >}}
```typescript
// Express application backed by a PostgreSQL database.
// Creates a heroes table at startup.
// Endpoints: GET / (greeting), GET /health (health check), POST /heroes/ (create), GET /heroes/ (list).
// See https://expressjs.com/ and https://node-postgres.com/

import express, { type Request, type Response } from 'express';
import { Pool } from 'pg';
import { readFileSync } from 'fs';

export const app = express();
const port = parseInt(process.env.PORT ?? '3000', 10);

app.use(express.json());

function getPassword(): string {
  const passwordFile = process.env.POSTGRES_PASSWORD_FILE;
  if (passwordFile) {
    return readFileSync(passwordFile, 'utf8').trim();
  }
  return process.env.POSTGRES_PASSWORD ?? '';
}

const pool = new Pool({
  host: process.env.POSTGRES_SERVER,
  port: 5432,
  database: process.env.POSTGRES_DB,
  user: process.env.POSTGRES_USER,
  password: getPassword(),
});

if (process.env.POSTGRES_SERVER) {
  pool
    .query(
      `CREATE TABLE IF NOT EXISTS heroes (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        secret_name TEXT NOT NULL,
        age INTEGER
      )`,
    )
    .catch(console.error);
}

app.get('/', (_req: Request, res: Response) => {
  res.json({ message: 'Hello World' });
});

app.get('/health', (_req: Request, res: Response) => {
  res.json({ status: 'ok' });
});

app.post('/heroes/', async (req: Request, res: Response) => {
  const { name, secret_name, age } = req.body as {
    name: string;
    secret_name: string;
    age?: number;
  };
  const result = await pool.query(
    'INSERT INTO heroes (name, secret_name, age) VALUES ($1, $2, $3) RETURNING *',
    [name, secret_name, age],
  );
  res.json(result.rows[0]);
});

app.get('/heroes/', async (_req: Request, res: Response) => {
  const result = await pool.query('SELECT * FROM heroes');
  res.json(result.rows);
});

// Only start the server when this file is run directly.
if (require.main === module) {
  app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
  });
}
```
{{< /file >}}

{{< file path="src/index.test.ts" status="new" >}}
```typescript
// Unit tests for the Express application.
// Tests the root endpoint without starting a server.
// See https://vitest.dev/ for the test framework reference.

import { describe, it, expect } from 'vitest';
import request from 'supertest';
import { app } from './index';

describe('GET /', () => {
  it('returns a JSON greeting', async () => {
    const response = await request(app).get('/');
    expect(response.status).toBe(200);
    expect(response.body).toEqual({ message: 'Hello World' });
  });
});
```
{{< /file >}}

{{< file path="package.json" status="modified" hl_lines="10,20-22" >}}
```json
{
  "name": "nodejs-docker-example",
  "version": "1.0.0",
  "description": "A minimal Node.js TypeScript application.",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js",
    "dev": "tsx watch src/index.ts",
    "test": "vitest run"
  },
  "dependencies": {
    "express": "^4.21.2",
    "pg": "^8.16.0"
  },
  "devDependencies": {
    "@types/express": "^4.17.21",
    "@types/node": "^22.0.0",
    "@types/pg": "^8.11.0",
    "supertest": "^7.0.0",
    "@types/supertest": "^6.0.0",
    "tsx": "^4.19.3",
    "typescript": "^5.8.3",
    "vitest": "^3.0.0"
  }
}
```
{{< /file >}}

{{< /files >}}

## Run tests locally

Run the following command to run the tests locally:

```console
$ npm install
$ npm test
```

You should see output like the following:

```console
 RUN  v3.0.0 /app

 ✓ src/index.test.ts (1)
   ✓ GET / (1)
     ✓ returns a JSON greeting

 Test Files  1 passed (1)
      Tests  1 passed (1)
   Start at  12:00:00
   Duration  500ms
```

## Run tests in a container

Run the tests using the dev stage of your Dockerfile:

```console
$ docker compose run --build --rm --no-deps server npm test
```

The `--no-deps` flag skips starting the database, since the unit tests don't require it. The `--rm` flag removes the container when the tests finish.

You should see the same test output as when running locally.

## Run tests when building

To run tests during the Docker build process, add a `test` stage to your Dockerfile that runs after the dev stage.

```dockerfile {hl_lines="32-36"}
FROM dhi.io/node:24-alpine3.23-dev AS dev

WORKDIR /app

RUN --mount=type=cache,target=/root/.npm \
    --mount=type=bind,source=package.json,target=package.json \
    npm install

COPY . .
RUN npm run build

EXPOSE 3000
CMD ["npm", "run", "dev"]


FROM dhi.io/node:24-alpine3.23-dev AS deps
WORKDIR /app
RUN --mount=type=cache,target=/root/.npm \
    --mount=type=bind,source=package.json,target=package.json \
    npm install --omit=dev

FROM dhi.io/node:24-alpine3.23 AS runner
ENV PATH=/app/node_modules/.bin:$PATH
WORKDIR /app
COPY --from=deps --chown=node:node /app/node_modules ./node_modules
COPY --from=dev --chown=node:node /app/dist ./dist

EXPOSE 3000
CMD ["node", "dist/index.js"]


FROM dev AS test

ENV CI=true

CMD ["npm", "test"]
```

Then build and run the test stage:

```console
$ docker build --target test -t nodejs-app-test .
$ docker run --rm nodejs-app-test
```

## Summary

In this section, you learned how to run tests when developing locally and inside a container.

Related information:

- [Dockerfile reference](/reference/dockerfile/)
- [Compose file reference](/compose/compose-file/)
- [`docker compose run` CLI reference](/reference/cli/docker/compose/run/)

## Next steps

In the next section, you'll learn how to set up a CI/CD pipeline using GitHub Actions.
