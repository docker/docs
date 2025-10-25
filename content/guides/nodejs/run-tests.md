---
title: Run Node.js tests in a container
linkTitle: Run your tests
weight: 30
keywords: node.js, node, test
description: Learn how to run your Node.js tests in a container.
aliases:
  - /language/nodejs/run-tests/
  - /guides/language/nodejs/run-tests/
---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Node.js application](containerize.md).

## Overview

Testing is an essential part of modern software development. Testing can mean a
lot of things to different development teams. There are unit tests, integration
tests and end-to-end testing. In this guide you take a look at running your unit
tests in Docker when developing and when building.

## Run tests when developing locally

The sample application uses Vitest for testing with comprehensive test coverage across client and server components. The test suite includes 101 passing tests covering React components, custom hooks, API routes, database operations, and utility functions.

### Run tests locally (without Docker)

```console
$ npm run test
```

### Add test service to Docker Compose

To run tests in a containerized environment, you need to add a dedicated test service to your `compose.yml` file. Add the following service configuration:

```yaml
services:
  # ... existing services ...

  # ========================================
  # Test Service
  # ========================================
  app-test:
    build:
      context: .
      dockerfile: Dockerfile
      target: test
    container_name: todoapp-test
    environment:
      NODE_ENV: test
      POSTGRES_HOST: db
      POSTGRES_PORT: 5432
      POSTGRES_DB: todoapp_test
      POSTGRES_USER: todoapp
      POSTGRES_PASSWORD: '${POSTGRES_PASSWORD:-todoapp_password}'
    depends_on:
      db:
        condition: service_healthy
    command: ['npm', 'run', 'test:coverage']
    networks:
      - todoapp-network
    profiles:
      - test
```

This test service configuration:

- **Builds from test stage**: Uses the `test` target from your multi-stage Dockerfile
- **Isolated test database**: Uses a separate `todoapp_test` database for testing
- **Profile-based**: Uses the `test` profile so it only runs when explicitly requested
- **Health dependency**: Waits for the database to be healthy before starting tests

### Run tests in a container

You can run tests using the dedicated test service:

```console
$ docker compose up app-test --build
```

Or run tests against the development service:

```console
$ docker compose run --rm app-dev npm run test
```

For a one-off test run with coverage:

```console
$ docker compose run --rm app-dev npm run test:coverage
```

### Run tests with coverage

To generate a coverage report:

```console
$ npm run test:coverage
```

You should see output like the following:

```console
> docker-nodejs-sample@1.0.0 test
> vitest --run

 ✓ src/server/__tests__/routes/todos.test.ts (5 tests) 16ms
 ✓ src/shared/utils/__tests__/validation.test.ts (15 tests) 6ms
 ✓ src/client/components/__tests__/LoadingSpinner.test.tsx (8 tests) 67ms
 ✓ src/server/database/__tests__/postgres.test.ts (13 tests) 136ms
 ✓ src/client/components/__tests__/ErrorMessage.test.tsx (8 tests) 127ms
 ✓ src/client/components/__tests__/TodoList.test.tsx (8 tests) 147ms
 ✓ src/client/components/__tests__/TodoItem.test.tsx (8 tests) 218ms
 ✓ src/client/__tests__/App.test.tsx (13 tests) 259ms
 ✓ src/client/components/__tests__/AddTodoForm.test.tsx (12 tests) 323ms
 ✓ src/client/hooks/__tests__/useTodos.test.ts (11 tests) 569ms

 Test Files  10 passed (10)
      Tests  101 passed (101)
   Start at  15:32:56
   Duration  1.98s (transform 456ms, setup 1.26s, collect 1.74s, tests 1.87s, environment 5.82s, prepare 916ms)
```

### Test Structure

The test suite covers:

- **Client Components** (`src/client/components/__tests__/`): React component testing with React Testing Library
- **Custom Hooks** (`src/client/hooks/__tests__/`): React hooks testing with proper mocking
- **Server Routes** (`src/server/__tests__/routes/`): API endpoint testing with Supertest
- **Database Layer** (`src/server/database/__tests__/`): PostgreSQL database operations testing
- **Utility Functions** (`src/shared/utils/__tests__/`): Validation and helper function testing
- **Integration Tests** (`src/client/__tests__/`): Full application integration testing

## Run tests when building

To run tests during the Docker build process, you need to add a dedicated test stage to your Dockerfile. If you haven't already added this stage, add the following to your multi-stage Dockerfile:

```dockerfile
# ========================================
# Test Stage
# ========================================
FROM build-deps AS test

# Set environment
ENV NODE_ENV=test \
    CI=true

# Copy source files
COPY --chown=nodejs:nodejs . .

# Switch to non-root user
USER nodejs

# Run tests with coverage
CMD ["npm", "run", "test:coverage"]
```

This test stage:

- **Test environment**: Sets `NODE_ENV=test` and `CI=true` for proper test execution
- **Non-root user**: Runs tests as the `nodejs` user for security
- **Flexible execution**: Uses `CMD` instead of `RUN` to allow running tests during build or as a separate container
- **Coverage support**: Configured to run tests with coverage reporting

### Build and run tests during image build

To build an image that runs tests during the build process, you can create a custom Dockerfile or modify the existing one temporarily:

```console
$ docker build --target test -t node-docker-image-test .
```

### Run tests in a dedicated test container

The recommended approach is to use the test service defined in `compose.yml`:

```console
$ docker compose --profile test up app-test --build
```

Or run it as a one-off container:

```console
$ docker compose run --rm app-test
```

### Run tests with coverage in CI/CD

For continuous integration, you can run tests with coverage:

```console
$ docker build --target test --progress=plain --no-cache -t test-image .
$ docker run --rm test-image npm run test:coverage
```

You should see output containing the following:

```console
 ✓ src/server/__tests__/routes/todos.test.ts (5 tests) 16ms
 ✓ src/shared/utils/__tests__/validation.test.ts (15 tests) 6ms
 ✓ src/client/components/__tests__/LoadingSpinner.test.tsx (8 tests) 67ms
 ✓ src/server/database/__tests__/postgres.test.ts (13 tests) 136ms
 ✓ src/client/components/__tests__/ErrorMessage.test.tsx (8 tests) 127ms
 ✓ src/client/components/__tests__/TodoList.test.tsx (8 tests) 147ms
 ✓ src/client/components/__tests__/TodoItem.test.tsx (8 tests) 218ms
 ✓ src/client/__tests__/App.test.tsx (13 tests) 259ms
 ✓ src/client/components/__tests__/AddTodoForm.test.tsx (12 tests) 323ms
 ✓ src/client/hooks/__tests__/useTodos.test.ts (11 tests) 569ms

 Test Files  10 passed (10)
      Tests  101 passed (101)
   Start at  07:33:25
   Duration  2.11s (transform 339ms, setup 619ms, collect 1.12s, tests 1.43s, environment 3.52s, prepare 901ms)
```

## Summary

In this section, you learned how to run tests when developing locally using Docker Compose and how to run tests when building your image.

Related information:

- [Dockerfile reference](/reference/dockerfile/) – Understand all Dockerfile instructions and syntax.
- [Best practices for writing Dockerfiles](/develop/develop-images/dockerfile_best-practices/) – Write efficient, maintainable, and secure Dockerfiles.
- [Compose file reference](/compose/compose-file/) – Learn the full syntax and options available for configuring services in `compose.yaml`.
- [`docker compose run` CLI reference](/reference/cli/docker/compose/run/) – Run one-off commands in a service container.

## Next steps

Next, you’ll learn how to set up a CI/CD pipeline using GitHub Actions.
