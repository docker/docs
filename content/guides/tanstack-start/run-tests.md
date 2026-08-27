---
title: Run TanStack Start tests in a container
linkTitle: Run your tests
weight: 40
keywords: tanstack start, test, vitest
description: Learn how to run your TanStack Start tests in a container.
---

## Prerequisites

Complete all the previous sections of this guide, starting with
[Containerize TanStack Start application](containerize.md).

## Overview

Testing is a critical part of the development process. In this section, you'll
learn how to:

- Run unit tests using Vitest inside a Docker container.
- Use Docker Compose to run tests in an isolated, reproducible environment.

The [sample project](https://github.com/kristiyan-velkov/docker-tanstack-start-sample) uses [Vitest](https://vitest.dev/) with
[Testing Library](https://testing-library.com/) for component testing.

---

## Run tests during development

The sample project includes a `test` script in `package.json`:

```json
"scripts": {
  "test": "vitest run"
}
```

The sample project includes a test file at:

```text
src/lib/utils.test.ts
```

If you're using your own project and haven't added Vitest yet, install the
testing tools:

```console
$ npm install --save-dev vitest @vitejs/plugin-react @testing-library/react @testing-library/dom jsdom
```

Then add the `test` script to `package.json` as shown above.

### Step 1: Configure Vitest

Add a `test` block to your `vite.config.ts`:

```ts {hl_lines="20-23",linenos=true}
import { defineConfig } from "vite";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import viteTsConfigPaths from "vite-tsconfig-paths";
import tailwindcss from "@tailwindcss/vite";
import { nitro } from "nitro/vite";

const config = defineConfig({
  plugins: [
    devtools(),
    nitro(),
    viteTsConfigPaths({ projects: ["./tsconfig.json"] }),
    tailwindcss(),
    tanstackStart(),
    viteReact(),
  ],
  server: {
    host: true,
    port: 3000,
    strictPort: true,
  },
  test: {
    environment: "jsdom",
    globals: true,
  },
});

export default config;
```

> [!NOTE]
> The `test` options configure Vitest for React component testing:
>
> - `environment: "jsdom"` simulates a browser-like environment.
> - `globals: true` exposes `describe`, `it`, and `expect` without imports.
>
> For more details, see the
> [Vitest configuration docs](https://vitest.dev/config/).

### Step 2: Update compose.yml

Add a `tanstack-start-test` service to your `compose.yml` file:

```yaml {hl_lines="35-39",linenos=true}
services:
  tanstack-start-prod:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        NODE_VERSION: 24.14.0-alpine
    image: tanstack-start:prod
    ports:
      - "3000:3000"

  tanstack-start-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
      args:
        NODE_VERSION: 24.14.0-alpine
    image: tanstack-start:dev
    ports:
      - "3000:3000"
    develop:
      watch:
        - action: sync
          path: .
          target: /app
          ignore:
            - node_modules/
            - .output/

  tanstack-start-test:
    build:
      context: .
      dockerfile: Dockerfile.dev
    command: ["npm", "run", "test"]
```

The `tanstack-start-test` service reuses `Dockerfile.dev` and overrides the
default command to run `npm run test`.

### Step 3: Run the tests

Execute your test suite inside the container:

```console
$ docker compose run --rm tanstack-start-test
```

This command starts the test service, runs Vitest, and removes the container
when tests finish.

> [!NOTE]
> For more information about Compose commands, see the
> [Compose CLI reference](/reference/cli/docker/compose/).

---

## Summary

In this section, you learned how to run unit tests for your TanStack Start
application inside a Docker container using Vitest and Docker Compose.

What you accomplished:

- Configured Vitest in `vite.config.ts` for component testing
- Created a `tanstack-start-test` service in `compose.yml`
- Ran tests with `docker compose run --rm tanstack-start-test`

---

## Related resources

- [Dockerfile reference](/reference/dockerfile/) – Dockerfile instructions and
  syntax
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) –
  Write maintainable and secure Dockerfiles
- [Compose file reference](/compose/compose-file/) – Configure services in
  `compose.yml`
- [`docker compose run` CLI reference](/reference/cli/docker/compose/run/) –
  Run one-off commands in a service container

## Next steps

Next, you'll set up a CI/CD pipeline using GitHub Actions to build, test, and
push your TanStack Start application image to Docker Hub.
