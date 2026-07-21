---
title: Getting started with Testcontainers for Node.js
linkTitle: Testcontainers for Node.js
description: Learn how to use Testcontainers for Node.js to test database interactions with a real PostgreSQL instance.
keywords: testcontainers, nodejs, javascript, testing, postgresql, integration testing, jest
summary: |
  Learn how to create a Node.js application and test database interactions
  using Testcontainers for Node.js with a real PostgreSQL instance.
aliases:
  - /guides/testcontainers-nodejs-getting-started/create-project/
  - /guides/testcontainers-nodejs-getting-started/run-tests/
  - /guides/testcontainers-nodejs-getting-started/write-tests/
params:
  tags: [testing]
  time: 15 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-getting-started-with-testcontainers-for-nodejs -->

In this guide, you will learn how to:

- Create a Node.js application that stores and retrieves customers from PostgreSQL
- Write integration tests using Testcontainers and Jest
- Run tests against a real PostgreSQL database in a Docker container

## Prerequisites

- Node.js 18+
- npm
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.

## Create the Node.js project

### Initialize the project

Create a new Node.js project:

```console
$ npm init -y
```

Add `pg`, `jest`, and `@testcontainers/postgresql` as dependencies:

```console
$ npm install pg --save
$ npm install jest @testcontainers/postgresql --save-dev
```

### Implement the customer repository

Create `src/customer-repository.js` with functions to manage customers in
PostgreSQL:

```javascript
async function createCustomerTable(client) {
  const sql =
    "CREATE TABLE IF NOT EXISTS customers (id INT NOT NULL, name VARCHAR NOT NULL, PRIMARY KEY (id))";
  await client.query(sql);
}

async function createCustomer(client, customer) {
  const sql = "INSERT INTO customers (id, name) VALUES($1, $2)";
  await client.query(sql, [customer.id, customer.name]);
}

async function getCustomers(client) {
  const sql = "SELECT * FROM customers";
  const result = await client.query(sql);
  return result.rows;
}

module.exports = { createCustomerTable, createCustomer, getCustomers };
```

The module provides three functions:

- `createCustomerTable()` creates the `customers` table if it doesn't exist.
- `createCustomer()` inserts a customer record.
- `getCustomers()` fetches all customer records.

## Write tests with Testcontainers

Create `src/customer-repository.test.js` with the test:

```javascript
const { Client } = require("pg");
const { PostgreSqlContainer } = require("@testcontainers/postgresql");
const {
  createCustomerTable,
  createCustomer,
  getCustomers,
} = require("./customer-repository");

describe("Customer Repository", () => {
  jest.setTimeout(60000);

  let postgresContainer;
  let postgresClient;

  beforeAll(async () => {
    postgresContainer = await new PostgreSqlContainer().start();
    postgresClient = new Client({
      connectionString: postgresContainer.getConnectionUri(),
    });
    await postgresClient.connect();
    await createCustomerTable(postgresClient);
  });

  afterAll(async () => {
    await postgresClient.end();
    await postgresContainer.stop();
  });

  it("should create and return multiple customers", async () => {
    const customer1 = { id: 1, name: "John Doe" };
    const customer2 = { id: 2, name: "Jane Doe" };

    await createCustomer(postgresClient, customer1);
    await createCustomer(postgresClient, customer2);

    const customers = await getCustomers(postgresClient);
    expect(customers).toEqual([customer1, customer2]);
  });
});
```

Here's what the test does:

- The `beforeAll` block starts a real PostgreSQL container using
  `PostgreSqlContainer`. It then creates a `pg` client connected to the
  container and sets up the `customers` table.
- The `afterAll` block closes the client connection and stops the container.
- The test inserts two customers, fetches all customers, and asserts the
  results match.

The test timeout is set to 60 seconds to allow time for the container to start
on the first run (when the Docker image needs to be pulled).

## Run tests and next steps

### Run the tests

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

### Summary

The Testcontainers for Node.js library helps you write integration tests using
the same type of database (Postgres) that you use in production, instead of
mocks. Because you aren't using mocks and instead talk to real services, you're
free to refactor code and still verify that the application works as expected.

In addition to PostgreSQL, Testcontainers provides dedicated
[modules](https://github.com/testcontainers/testcontainers-node/tree/main/packages/modules)
for many SQL databases, NoSQL databases, messaging queues, and more.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [Testcontainers for Node.js documentation](https://node.testcontainers.org)
