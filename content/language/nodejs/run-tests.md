---
title: Run Node.js tests in a container
keywords: node.js, node, test
description: Learn how to run your Node.js tests in a container.
---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Node.js application](containerize.md).

## Overview

Testing is an essential part of modern software development. Testing can mean a
lot of things to different development teams. There are unit tests, integration
tests and end-to-end testing. In this guide you take a look at running your unit
tests in Docker when developing and when building.

## Run tests when developing locally

The sample application already has the Jest package for running tests and has tests inside the `spec` directory. When developing locally, you can use Compose to run your tests.

Run the following command to run the test script from the `package.json` file inside a container.

```console
$ docker compose run server npm run test
```

To learn more about the command, see [docker compose run](/engine/reference/commandline/compose_run/).

You should see output like the following.

```console
> docker-nodejs@1.0.0 test
> jest

 PASS  spec/routes/deleteItem.spec.js
 PASS  spec/routes/getItems.spec.js
 PASS  spec/routes/addItem.spec.js
 PASS  spec/routes/updateItem.spec.js
 PASS  spec/persistence/sqlite.spec.js
  ● Console

    console.log
      Using sqlite database at /tmp/todo.db

      at Database.log (src/persistence/sqlite.js:18:25)

    console.log
      Using sqlite database at /tmp/todo.db

      at Database.log (src/persistence/sqlite.js:18:25)

    console.log
      Using sqlite database at /tmp/todo.db

      at Database.log (src/persistence/sqlite.js:18:25)

    console.log
      Using sqlite database at /tmp/todo.db

      at Database.log (src/persistence/sqlite.js:18:25)

    console.log
      Using sqlite database at /tmp/todo.db

      at Database.log (src/persistence/sqlite.js:18:25)


Test Suites: 5 passed, 5 total
Tests:       9 passed, 9 total
Snapshots:   0 total
Time:        2.008 s
Ran all test suites.
```

## Run tests when building

To run your tests when building, you need to update your Dockerfile to add a new test stage.

The following is the updated Dockerfile.

```dockerfile {hl_lines="27-35"}
# syntax=docker/dockerfile:1

ARG NODE_VERSION=18.0.0

FROM node:${NODE_VERSION}-alpine as base
WORKDIR /usr/src/app
EXPOSE 3000

FROM base as dev
RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=package-lock.json,target=package-lock.json \
    --mount=type=cache,target=/root/.npm \
    npm ci --include=dev
USER node
COPY . .
CMD npm run dev

FROM base as prod
RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=package-lock.json,target=package-lock.json \
    --mount=type=cache,target=/root/.npm \
    npm ci --omit=dev
USER node
COPY . .
CMD node src/index.js

FROM base as test
ENV NODE_ENV test
RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=package-lock.json,target=package-lock.json \
    --mount=type=cache,target=/root/.npm \
    npm ci --include=dev
USER node
COPY . .
RUN npm run test
```

Instead of using `CMD` in the test stage, use `RUN` to run the tests. The reason is that the `CMD` instruction runs when the container runs, and the `RUN` instruction runs when the image is being built and the build will fail if the tests fail.

Run the following command to build a new image using the test stage as the target and view the test results. Include `--progress=plain` to view the build output, `--no-cache` to ensure the tests always run, and `--target test` to target the test stage.

```console
$ docker build -t node-docker-image-test --progress=plain --no-cache --target test .
```

To learn more about building and running tests, see the [Build with Docker guide](../../build/guide/_index.md).

You should see output containing the following.

```console
...

#11 [test 3/3] RUN npm run test
#11 1.058
#11 1.058 > docker-nodejs@1.0.0 test
#11 1.058 > jest
#11 1.058
#11 3.765 PASS spec/routes/getItems.spec.js
#11 3.767 PASS spec/routes/deleteItem.spec.js
#11 3.783 PASS spec/routes/updateItem.spec.js
#11 3.806 PASS spec/routes/addItem.spec.js
#11 4.179 PASS spec/persistence/sqlite.spec.js
#11 4.207
#11 4.208 Test Suites: 5 passed, 5 total
#11 4.208 Tests:       9 passed, 9 total
#11 4.208 Snapshots:   0 total
#11 4.208 Time:        2.168 s
#11 4.208 Ran all test suites.
#11 4.265 npm notice
#11 4.265 npm notice New major version of npm available! 8.6.0 -> 9.8.1
#11 4.265 npm notice Changelog: <https://github.com/npm/cli/releases/tag/v9.8.1>
#11 4.265 npm notice Run `npm install -g npm@9.8.1` to update!
#11 4.266 npm notice
#11 DONE 4.3s

...
```

## Summary

In this section, you learned how to run tests when developing locally using Compose and how to run tests when building your image.

Related information:
 - [docker compose run](/engine/reference/commandline/compose_run/)
 - [Build with Docker guide](../../build/guide/index.md)

## Next steps

Next, you’ll learn how to set up a CI/CD pipeline using GitHub Actions.

{{< button text="Configure CI/CD" url="configure-ci-cd.md" >}}