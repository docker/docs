---
title: Getting started with Testcontainers for Go
linkTitle: Testcontainers for Go
description: Learn how to use Testcontainers for Go to test database interactions with a real PostgreSQL instance.
keywords: testcontainers, go, golang, testing, postgresql, integration testing
summary: |
  Learn how to create a Go application and test database interactions
  using Testcontainers for Go with a real PostgreSQL instance.
aliases:
  - /guides/testcontainers-go-getting-started/create-project/
  - /guides/testcontainers-go-getting-started/run-tests/
  - /guides/testcontainers-go-getting-started/test-suites/
  - /guides/testcontainers-go-getting-started/write-tests/
params:
  tags: [testing]
  time: 20 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-getting-started-with-testcontainers-for-go -->

In this guide, you will learn how to:

- Create a Go application with modules support
- Implement a Repository to manage customer data in a PostgreSQL database using the pgx driver
- Write integration tests using testcontainers-go
- Reuse containers across multiple tests using test suites

## Prerequisites

- Go 1.25+
- Your preferred IDE (VS Code, GoLand)
- A Docker environment supported by Testcontainers. For details, see
  the [testcontainers-go system requirements](https://golang.testcontainers.org/system_requirements/).

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.

## Create the Go project

### Initialize the project

Start by creating a Go project.

```console
$ mkdir testcontainers-go-demo
$ cd testcontainers-go-demo
$ go mod init github.com/testcontainers/testcontainers-go-demo
```

This guide uses the [jackc/pgx](https://github.com/jackc/pgx) PostgreSQL
driver to interact with the Postgres database and the testcontainers-go
[Postgres module](https://golang.testcontainers.org/modules/postgres/) to
spin up a Postgres Docker instance for testing. It also uses
[testify](https://github.com/stretchr/testify) for running multiple tests
as a suite and for writing assertions.

Install these dependencies:

```console
$ go get github.com/jackc/pgx/v5
$ go get github.com/testcontainers/testcontainers-go
$ go get github.com/testcontainers/testcontainers-go/modules/postgres
$ go get github.com/stretchr/testify
```

### Create Customer struct

Create a `types.go` file in the `customer` package and define the `Customer`
struct to model the customer details:

```go
package customer

type Customer struct {
	Id    int
	Name  string
	Email string
}
```

### Create Repository

Next, create `customer/repo.go`, define the `Repository` struct, and add
methods to create a customer and get a customer by email:

```go
package customer

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	return &Repository{
		conn: conn,
	}, nil
}

func (r Repository) CreateCustomer(ctx context.Context, customer Customer) (Customer, error) {
	err := r.conn.QueryRow(ctx,
		"INSERT INTO customers (name, email) VALUES ($1, $2) RETURNING id",
		customer.Name, customer.Email).Scan(&customer.Id)
	return customer, err
}

func (r Repository) GetCustomerByEmail(ctx context.Context, email string) (Customer, error) {
	var customer Customer
	query := "SELECT id, name, email FROM customers WHERE email = $1"
	err := r.conn.QueryRow(ctx, query, email).
		Scan(&customer.Id, &customer.Name, &customer.Email)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}
```

Here's what the code does:

- `Repository` holds a `*pgx.Conn` for performing database operations.
- `NewRepository(connStr)` takes a database connection string and initializes a `Repository`.
- `CreateCustomer()` and `GetCustomerByEmail()` are methods on the `Repository` receiver that insert and query customer records.

## Write tests with Testcontainers

You have the `Repository` implementation ready, but for testing you need a
PostgreSQL database. You can use testcontainers-go to spin up a Postgres
database in a Docker container and run your tests against that database.

### Set up the test database

In real applications you might use a database migration tool, but for this
guide, use a script to initialize the database.

Create a `testdata/init-db.sql` file to create the `CUSTOMERS` table and
insert sample data:

```sql
CREATE TABLE IF NOT EXISTS customers (id serial, name varchar(255), email varchar(255));

INSERT INTO customers(name, email) VALUES ('John', 'john@gmail.com');
```

### Understand the testcontainers-go API

The testcontainers-go library provides the generic `Container` abstraction
that can run any containerized service. To further simplify, testcontainers-go
provides technology-specific modules that reduce boilerplate and provide a
functional options pattern to construct the container instance.

For example, `PostgresContainer` provides `WithDatabase()`,
`WithUsername()`, `WithPassword()`, and other functions to set various
properties of Postgres containers.

### Write the test

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

## Reuse containers with test suites

In the previous section, you saw how to spin up a Postgres Docker container
for a single test. But often you have multiple tests in a single file, and you
may want to reuse the same Postgres Docker container for all of them.

You can use the [testify suite](https://pkg.go.dev/github.com/stretchr/testify/suite)
package to implement common test setup and teardown actions.

### Extract container setup

First, extract the `PostgresContainer` creation logic into a separate file
called `testhelpers/containers.go`:

```go
package testhelpers

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(t *testing.T, ctx context.Context) *PostgresContainer {
	t.Helper()

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

	return &PostgresContainer{
		PostgresContainer: ctr,
		ConnectionString:  connStr,
	}
}
```

In `containers.go`, `PostgresContainer` extends the testcontainers-go
`PostgresContainer` to provide easy access to `ConnectionString`. The
`CreatePostgresContainer()` function accepts `*testing.T` as its first
parameter, calls `t.Helper()` so that test failures point to the caller,
and uses `testcontainers.CleanupContainer()` to register automatic cleanup.

### Write the test suite

Create `customer/repo_suite_test.go` and implement tests for creating
a customer and getting a customer by email using the testify suite package:

```go
package customer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go-demo/testhelpers"
)

type CustomerRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *Repository
	ctx         context.Context
}

func (suite *CustomerRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.pgContainer = testhelpers.CreatePostgresContainer(suite.T(), suite.ctx)

	repository, err := NewRepository(suite.ctx, suite.pgContainer.ConnectionString)
	require.NoError(suite.T(), err)
	suite.repository = repository
}

func (suite *CustomerRepoTestSuite) TestCreateCustomer() {
	t := suite.T()

	customer, err := suite.repository.CreateCustomer(suite.ctx, Customer{
		Name:  "Henry",
		Email: "henry@gmail.com",
	})
	require.NoError(t, err)
	assert.NotNil(t, customer.Id)
}

func (suite *CustomerRepoTestSuite) TestGetCustomerByEmail() {
	t := suite.T()

	customer, err := suite.repository.GetCustomerByEmail(suite.ctx, "john@gmail.com")
	require.NoError(t, err)
	assert.Equal(t, "John", customer.Name)
	assert.Equal(t, "john@gmail.com", customer.Email)
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepoTestSuite))
}
```

Here's what the code does:

- `CustomerRepoTestSuite` extends `suite.Suite` and includes fields shared
  across multiple tests.
- `SetupSuite()` runs once before all tests. It calls
  `CreatePostgresContainer(suite.T(), ...)` which handles cleanup registration
  automatically via `CleanupContainer`, so no `TearDownSuite()` is needed.
- `TestCreateCustomer()` uses `require.NoError()` for the create operation
  (fail immediately if it errors) and `assert.NotNil()` for the ID check.
- `TestGetCustomerByEmail()` uses `require.NoError()` then asserts on the
  returned values.
- `TestCustomerRepoTestSuite(t *testing.T)` runs the test suite when you
  execute `go test`.

> [!TIP]
> For the purpose of this guide, the tests don't reset data in the database.
> In practice, it's a good idea to reset the database to a known state before
> running each test.

## Run tests and next steps

### Run the tests

Run all the tests using `go test ./...`. Optionally add the `-v` flag for
verbose output:

```console
$ go test -v ./...
```

You should see two Postgres Docker containers start automatically: one for the
suite and its two tests, and another for the initial standalone test. All tests
should pass. After the tests finish, the containers are stopped and removed
automatically.

### Summary

The Testcontainers for Go library helps you write integration tests by using
the same type of database (Postgres) that you use in production, instead of
mocks. Because you aren't using mocks and instead talk to real services, you're
free to refactor code and still verify that the application works as expected.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [Testcontainers for Go documentation](https://golang.testcontainers.org/)
- [Testcontainers for Go quickstart](https://golang.testcontainers.org/quickstart/)
- [Testcontainers Postgres module for Go](https://golang.testcontainers.org/modules/postgres/)
