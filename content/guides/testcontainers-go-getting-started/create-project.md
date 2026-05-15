---
title: Create the Go project
linkTitle: Create the project
description: Set up a Go project with a PostgreSQL-backed repository.
weight: 10
---

## Initialize the project

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

## Create Customer struct

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

## Create Repository

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
