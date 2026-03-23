---
title: Reuse containers with test suites
linkTitle: Test suites
description: Share a single Postgres container across multiple tests using testify suites.
weight: 30
---

In the previous section, you saw how to spin up a Postgres Docker container
for a single test. But often you have multiple tests in a single file, and you
may want to reuse the same Postgres Docker container for all of them.

You can use the [testify suite](https://pkg.go.dev/github.com/stretchr/testify/suite)
package to implement common test setup and teardown actions.

## Extract container setup

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

## Write the test suite

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
