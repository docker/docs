---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Write your first integration test using testcontainers-go and PostgreSQL.
weight: 20
---

You have the `Repository` implementation ready, but for testing you need a
PostgreSQL database. You can use testcontainers-go to spin up a Postgres
database in a Docker container and run your tests against that database.

## Set up the test database

In real applications you might use a database migration tool, but for this
guide, use a script to initialize the database.

Create a `testdata/init-db.sql` file to create the `CUSTOMERS` table and
insert sample data:

```sql
CREATE TABLE IF NOT EXISTS customers (id serial, name varchar(255), email varchar(255));

INSERT INTO customers(name, email) VALUES ('John', 'john@gmail.com');
```

## Understand the testcontainers-go API

The testcontainers-go library provides the generic `Container` abstraction
that can run any containerized service. To further simplify, testcontainers-go
provides technology-specific modules that reduce boilerplate and provide a
functional options pattern to construct the container instance.

For example, `PostgresContainer` provides `WithDatabase()`,
`WithUsername()`, `WithPassword()`, and other functions to set various
properties of Postgres containers.

## Write the test

Create the `customer/repo_test.go` file and implement the test:

```go
package customer

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestCustomerRepository(t *testing.T) {
	ctx := context.Background()

	ctr, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("..", "testdata", "init-db.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	customerRepo, err := NewRepository(ctx, connStr)
	require.NoError(t, err)

	c, err := customerRepo.CreateCustomer(ctx, Customer{
		Name:  "Henry",
		Email: "henry@gmail.com",
	})
	assert.NoError(t, err)
	assert.NotNil(t, c)

	customer, err := customerRepo.GetCustomerByEmail(ctx, "henry@gmail.com")
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "Henry", customer.Name)
	assert.Equal(t, "henry@gmail.com", customer.Email)
}
```

Here's what the test does:

- Calls `postgres.Run()` with the `postgres:16-alpine` Docker image as the
  first argument. This is the v0.41.0 API — the image is a required positional
  parameter instead of an option.
- Configures initialization scripts using `WithInitScripts(...)` so that the
  `CUSTOMERS` table is created and sample data is inserted after the database
  starts.
- Uses `postgres.BasicWaitStrategies()` which combines waiting for the Postgres
  log message and for the port to be ready. This replaces manual wait strategy
  configuration.
- Calls `testcontainers.CleanupContainer(t, ctr)` right after `postgres.Run()`.
  This registers automatic cleanup with the test framework, replacing the manual
  `t.Cleanup` and `Terminate` pattern.
- Obtains the database `ConnectionString` from the container and initializes a
  `Repository`.
- Creates a customer with the email `henry@gmail.com` and verifies that the
  customer exists in the database.
