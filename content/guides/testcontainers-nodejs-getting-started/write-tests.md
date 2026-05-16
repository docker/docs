---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Write integration tests using Testcontainers for Node.js and Jest with a real PostgreSQL database.
weight: 20
---

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
