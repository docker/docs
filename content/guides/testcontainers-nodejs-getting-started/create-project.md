---
title: Create the Node.js project
linkTitle: Create the project
description: Set up a Node.js project with a PostgreSQL-backed customer repository.
weight: 10
---

## Initialize the project

Create a new Node.js project:

```console
$ npm init -y
```

Add `pg`, `jest`, and `@testcontainers/postgresql` as dependencies:

```console
$ npm install pg --save
$ npm install jest @testcontainers/postgresql --save-dev
```

## Implement the customer repository

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
