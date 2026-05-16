---
title: Run tests and next steps
linkTitle: Run tests
description: Run your Testcontainers-based integration tests and explore next steps.
weight: 40
---

## Run the tests

Run all the tests using `go test ./...`. Optionally add the `-v` flag for
verbose output:

```console
$ go test -v ./...
```

You should see two Postgres Docker containers start automatically: one for the
suite and its two tests, and another for the initial standalone test. All tests
should pass. After the tests finish, the containers are stopped and removed
automatically.

## Summary

The Testcontainers for Go library helps you write integration tests by using
the same type of database (Postgres) that you use in production, instead of
mocks. Because you aren't using mocks and instead talk to real services, you're
free to refactor code and still verify that the application works as expected.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

## Further reading

- [Testcontainers for Go documentation](https://golang.testcontainers.org/)
- [Testcontainers for Go quickstart](https://golang.testcontainers.org/quickstart/)
- [Testcontainers Postgres module for Go](https://golang.testcontainers.org/modules/postgres/)
