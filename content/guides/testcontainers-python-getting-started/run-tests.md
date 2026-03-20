---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers-based integration tests and explore next steps.
weight: 30
---

## Run the tests

Run the tests using pytest:

```console
$ pytest -v
```

You should see output similar to:

```text
============================= test session starts ==============================
platform linux -- Python 3.13.x, pytest-9.x.x
collected 2 items

tests/test_customers.py::test_get_all_customers PASSED                   [ 50%]
tests/test_customers.py::test_get_customer_by_email PASSED               [100%]

============================== 2 passed in 1.90s ===============================
```

The tests run against a real PostgreSQL database instead of mocks, which gives
more confidence in the implementation.

## Summary

The Testcontainers for Python library helps you write integration tests using the
same type of database (Postgres) that you use in production, instead of mocks.
Because you aren't using mocks and instead talk to real services, you're free
to refactor code and still verify that the application works as expected.

In addition to PostgreSQL, Testcontainers for Python provides modules for many
SQL databases, NoSQL databases, messaging queues, and more. You can use
Testcontainers to run any containerized dependency for your tests.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [testcontainers-python documentation](https://testcontainers-python.readthedocs.io/)
- [Getting started with Testcontainers for Go](/guides/testcontainers-go-getting-started/)
- [Getting started with Testcontainers for Java](https://testcontainers.com/guides/getting-started-with-testcontainers-for-java/)
- [Getting started with Testcontainers for Node.js](https://testcontainers.com/guides/getting-started-with-testcontainers-for-nodejs/)
