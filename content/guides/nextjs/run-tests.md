---
title: Run Next.js tests in a container
linkTitle: Run your tests and lint
weight: 40
keywords: next.js, test, jest, vitest, lint, eslint
description: Learn how to run your Next.js tests and lint in a container.

---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize Next.js application](containerize.md).

## Overview

Testing is a critical part of the development process. In this section, you'll learn how to:

- Run unit tests using Vitest (or Jest) inside a Docker container.
- Run lint (e.g. ESLint) inside a Docker container.
- Use Docker Compose to run tests and lint in an isolated, reproducible environment.

The [sample project](https://github.com/kristiyan-velkov/docker-nextjs-sample) uses [Vitest](https://vitest.dev/) with [Testing Library](https://testing-library.com/) for component testing. You can use the same setup or follow the alternative Jest configuration later.

---

## Run tests during development

The [sample project](https://github.com/kristiyan-velkov/docker-nextjs-sample) already includes lint (ESLint) and sample tests (Vitest, `app/page.test.tsx`) in place. If you're using the sample app, you can skip to **Step 3: Update compose.yaml** and run tests or lint with the commands below. If you're using your own project, follow the install and configuration steps to add the packages and scripts.

The sample includes a test file at:

```text
app/page.test.tsx
```

This file uses Vitest and React Testing Library to verify the behavior of page components.

### Step 1: Install Vitest and React Testing Library (custom projects)

If you're using a custom project and haven't already added the necessary testing tools, install them by running:

```console
$ npm install --save-dev vitest @vitejs/plugin-react @testing-library/react @testing-library/dom jsdom
```

Then, update the scripts section of your `package.json` file to include:

```json
"scripts": {
  "test": "vitest",
  "test:run": "vitest run"
}
```

For lint, add a `lint` script (and optionally `lint:fix`). For example, with [ESLint](https://eslint.org/):

```json
"scripts": {
  "test": "vitest",
  "test:run": "vitest run",
  "lint": "eslint .",
  "lint:fix": "eslint . --fix"
}
```

The sample project uses `eslint` and `eslint-config-next` for Next.js. Install them in a custom project with:

```console
$ npm install --save-dev eslint eslint-config-next @eslint/eslintrc
```

Create an ESLint config file (e.g. `eslint.config.cjs`) in your project root with Next.js rules and global ignores:

```js
const { defineConfig, globalIgnores } = require("eslint/config");
const { FlatCompat } = require("@eslint/eslintrc");

const compat = new FlatCompat({ baseDirectory: __dirname });

module.exports = defineConfig([
  ...compat.extends(
    "eslint-config-next/core-web-vitals",
    "eslint-config-next/typescript"
  ),
  globalIgnores([
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
    "node_modules/**",
    "eslint.config.cjs",
  ]),
]);
```

---

### Step 2: Configure Vitest (custom projects)

If you're using a custom project, create a `vitest.config.ts` file in your project root (matching the [sample project](https://github.com/kristiyan-velkov/docker-nextjs-sample)):

```ts
import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    setupFiles: "./vitest.setup.ts",
    globals: true,
  },
});
```

Create a `vitest.setup.ts` file in your project root:

```ts
import "@testing-library/jest-dom/vitest";
```

> [!NOTE]
> Vitest works well with Next.js and provides fast execution and ESM support. For more details, see the [Next.js testing documentation](https://nextjs.org/docs/app/building-your-application/testing) and [Vitest docs](https://vitest.dev/).

### Step 3: Update compose.yaml

Add `nextjs-test` and `nextjs-lint` services to your `compose.yaml` file. In the sample project these services use the `tools` profile so they don't start with a normal `docker compose up`. Both reuse `Dockerfile.dev` and run the test or lint command:

```yaml
services:
  nextjs-prod-standalone:
    build:
      context: .
      dockerfile: Dockerfile
    image: nextjs-sample:prod
    container_name: nextjs-sample-prod
    ports:
      - "3000:3000"

  nextjs-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    image: nextjs-sample:dev
    container_name: nextjs-sample-dev
    ports:
      - "3000:3000"
    environment:
      - WATCHPACK_POLLING=true
    develop:
      watch:
        - action: sync
          path: .
          target: /app
          ignore:
            - node_modules/
            - .next/
        - action: rebuild
          path: package.json

  nextjs-test:
    build:
      context: .
      dockerfile: Dockerfile.dev
    image: nextjs-sample:dev
    container_name: nextjs-sample-test
    command:
      [
        "sh",
        "-c",
        "if [ -f package-lock.json ]; then npm run test:run 2>/dev/null || npm run test -- --run; elif [ -f yarn.lock ]; then yarn test:run 2>/dev/null || yarn test --run; elif [ -f pnpm-lock.yaml ]; then pnpm run test:run; else npm run test -- --run; fi",
      ]
    profiles:
      - tools

  nextjs-lint:
    build:
      context: .
      dockerfile: Dockerfile.dev
    image: nextjs-sample:dev
    container_name: nextjs-sample-lint
    command:
      [
        "sh",
        "-c",
        "if [ -f package-lock.json ]; then npm run lint; elif [ -f yarn.lock ]; then yarn lint; elif [ -f pnpm-lock.yaml ]; then pnpm lint; else npm run lint; fi",
      ]
    profiles:
      - tools
```

The `nextjs-test` and `nextjs-lint` services reuse the same `Dockerfile.dev` used for [development](develop.md) and override the default command to run tests or lint. The `profiles: [tools]` means these services only run when you use the `--profile tools` option.

After completing the previous steps, your project directory should contain:

```text
├── docker-nextjs-sample/
│ ├── Dockerfile
│ ├── Dockerfile.dev
│ ├── .dockerignore
│ ├── compose.yaml
│ ├── vitest.config.ts
│ ├── vitest.setup.ts
│ ├── next.config.ts
│ └── README.Docker.md
```

### Step 4: Run the tests

To execute your test suite inside the container, run from your project root:

```console
$ docker compose --profile tools run --rm nextjs-test
```

This command will:
- Start the `nextjs-test` service (because of `--profile tools`).
- Run your test script (`test:run` or `test -- --run`) in the same environment as development.
- Remove the container after the tests complete ([`docker compose run --rm`](/reference/cli/docker/compose/run/)).

> [!NOTE]
> For more information about Compose commands and profiles, see the [Compose CLI reference](/reference/cli/docker/compose/).

### Step 5: Run lint in the container

To run your linter (e.g. ESLint) inside the container, use the `nextjs-lint` service with the same `tools` profile:

```console
$ docker compose --profile tools run --rm nextjs-lint
```

This command will:
- Start the `nextjs-lint` service (because of `--profile tools`).
- Run your lint script (`npm run lint`, `yarn lint`, or `pnpm lint` depending on your lockfile) in the same environment as development.
- Remove the container after lint completes.

Ensure your `package.json` includes a `lint` script. The sample project already has `"lint": "eslint ."` and `"lint:fix": "eslint . --fix"`; for a custom project, add the same and install `eslint` and `eslint-config-next` if needed.

---

## Summary

In this section, you learned how to run unit tests for your Next.js application inside a Docker container using Vitest and Docker Compose.

What you accomplished:
- Installed and configured Vitest and React Testing Library for testing Next.js components.
- Created `nextjs-test` and `nextjs-lint` services in `compose.yaml` (with `tools` profile) to isolate test and lint execution.
- Reused the development `Dockerfile.dev` to ensure consistency between dev, test, and lint environments.
- Ran tests inside the container using `docker compose --profile tools run --rm nextjs-test`.
- Ran lint inside the container using `docker compose --profile tools run --rm nextjs-lint`.
- Ensured reliable, repeatable testing and linting across environments without relying on local machine setup.

---

## Related resources

Explore official references and best practices to sharpen your Docker testing workflow:

- [Dockerfile reference](/reference/dockerfile/) – Understand all Dockerfile instructions and syntax.
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) – Write efficient, maintainable, and secure Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.  
- [`docker compose run` CLI reference](/reference/cli/docker/compose/run/) – Run one-off commands in a service container.
- [Next.js Testing Documentation](https://nextjs.org/docs/app/building-your-application/testing) – Official Next.js testing guide.
---

## Next steps

Next, you'll learn how to set up a CI/CD pipeline using GitHub Actions to automatically build and test your Next.js application in a containerized environment. This ensures your code is validated on every push or pull request, maintaining consistency and reliability across your development workflow.
