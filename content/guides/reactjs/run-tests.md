---
title: Run React.js tests in a container
linkTitle: Run your tests
weight: 40
keywords: react.js, react, test, vitest
description: Learn how to run your React.js tests in a container.

---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize React.js application](containerize.md).

## Overview

Testing is a critical part of the development process. In this section, you'll learn how to:

- Run unit tests using Vitest inside a Docker container.
- Use Docker Compose to run tests in an isolated, reproducible environment.

You’ll use [Vitest](https://vitest.dev) — a blazing fast test runner designed for Vite — along with [Testing Library](https://testing-library.com/) for assertions.

---

## Run tests during development

`docker-reactjs-sample` application includes a sample test file at location:

```console
$ src/App.test.tsx
```

This file uses Vitest and React Testing Library to verify the behavior of `App` component.

### Step 1: Install Vitest and React Testing Library

If you haven’t already added the necessary testing tools, install them by running:

```console
$ npm install --save-dev vitest @testing-library/react @testing-library/jest-dom jsdom
```

Then, update the scripts section of your `package.json` file to include the following:

```json
"scripts": {
  "test": "vitest run"
}
```

---

### Step 2: Configure Vitest

Update `vitest.config.ts` file in your project root with the following configuration:

```ts {hl_lines="14-18",linenos=true}
/// <reference types="vitest" />

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  base: "/",
  plugins: [react()],
  server: {
    host: true,
    port: 5173,
    strictPort: true,
  },
  test: {
    environment: "jsdom",
    setupFiles: "./src/setupTests.ts",
    globals: true,
  },
});
```

> [!NOTE]
> The `test` options in `vitest.config.ts` are essential for reliable testing inside Docker:
> - `environment: "jsdom"` simulates a browser-like environment for rendering and DOM interactions.  
> - `setupFiles: "./src/setupTests.ts"` loads global configuration or mocks before each test file (optional but recommended).  
> - `globals: true` enables global test functions like `describe`, `it`, and `expect` without importing them.
>
> For more details, see the official [Vitest configuration docs](https://vitest.dev/config/).

### Step 3: Update compose.yaml

Add a new service named `react-test` to your `compose.yaml` file. This service allows you to run your test suite in an isolated containerized environment.

```yaml {hl_lines="22-26",linenos=true}
services:
  react-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "5173:5173"
    develop:
      watch:
        - action: sync
          path: .
          target: /app

  react-prod:
    build:
      context: .
      dockerfile: Dockerfile
    image: docker-reactjs-sample
    ports:
      - "8080:8080"

  react-test:
    build:
      context: .
      dockerfile: Dockerfile.dev
    command: ["npm", "run", "test"]

```

The react-test service reuses the same `Dockerfile.dev` used for [development](develop.md) and overrides the default command to run tests with `npm run test`. This setup ensures a consistent test environment that matches your local development configuration.


After completing the previous steps, your project directory should contain the following files:

```text
├── docker-reactjs-sample/
│ ├── Dockerfile
│ ├── Dockerfile.dev
│ ├── .dockerignore
│ ├── compose.yaml
│ ├── nginx.conf
│ └── README.Docker.md
```

### Step 4: Run the tests

To execute your test suite inside the container, run the following command from your project root:

```console
$ docker compose run --rm react-test
```

This command will:
- Start the `react-test` service defined in your `compose.yaml` file.
- Execute the `npm run test` script using the same environment as development.
- Automatically remove the container after the tests complete [`docker compose run --rm`](/engine/reference/commandline/compose_run) command.

> [!NOTE]
> For more information about Compose commands, see the [Compose CLI
> reference](/reference/cli/docker/compose/_index.md).

---

## Summary

In this section, you learned how to run unit tests for your React.js application inside a Docker container using Vitest and Docker Compose.

What you accomplished:
- Installed and configured Vitest and React Testing Library for testing React components.
- Created a `react-test` service in `compose.yaml` to isolate test execution.
- Reused the development `Dockerfile.dev` to ensure consistency between dev and test environments.
- Ran tests inside the container using `docker compose run --rm react-test`.
- Ensured reliable, repeatable testing across environments without relying on local machine setup.

---

## Related resources

Explore official references and best practices to sharpen your Docker testing workflow:

- [Dockerfile reference](/reference/dockerfile/) – Understand all Dockerfile instructions and syntax.
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) – Write efficient, maintainable, and secure Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.  
- [`docker compose run` CLI reference](/reference/cli/docker/compose/run/) – Run one-off commands in a service container.
---

## Next steps

Next, you’ll learn how to set up a CI/CD pipeline using GitHub Actions to automatically build and test your React.js application in a containerized environment. This ensures your code is validated on every push or pull request, maintaining consistency and reliability across your development workflow.
