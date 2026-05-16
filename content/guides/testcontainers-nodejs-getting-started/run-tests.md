---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers-based integration tests and explore next steps.
weight: 30
---

## Run the tests

Add the test script to `package.json` if it isn't there already:

```json
{
  "scripts": {
    "test": "jest"
  }
}
```

Then run the tests:

```console
$ npm test
```

You should see output like:

```text
 PASS  src/customer-repository.test.js
  Customer Repository
    ✓ should create and return multiple customers (5 ms)

Test Suites: 1 passed, 1 total
Tests:       1 passed, 1 total
```

To see what Testcontainers is doing under the hood — which containers it
starts, what versions it uses — set the `DEBUG` environment variable:

```console
$ DEBUG=testcontainers* npm test
```

## Summary

The Testcontainers for Node.js library helps you write integration tests using
the same type of database (Postgres) that you use in production, instead of
mocks. Because you aren't using mocks and instead talk to real services, you're
free to refactor code and still verify that the application works as expected.

In addition to PostgreSQL, Testcontainers provides dedicated
[modules](https://github.com/testcontainers/testcontainers-node/tree/main/packages/modules)
for many SQL databases, NoSQL databases, messaging queues, and more.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Testcontainers for Node.js documentation](https://node.testcontainers.org)
