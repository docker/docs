---
title: Getting started with Testcontainers for Python
linkTitle: Testcontainers for Python
description: Learn how to use Testcontainers for Python to test database interactions with a real PostgreSQL instance.
keywords: testcontainers, python, testing, postgresql, integration testing, pytest
summary: |
  Learn how to create a Python application and test database interactions
  using Testcontainers for Python with a real PostgreSQL instance.
aliases:
  - /guides/testcontainers-python-getting-started/create-project/
  - /guides/testcontainers-python-getting-started/run-tests/
  - /guides/testcontainers-python-getting-started/write-tests/
params:
  tags: [testing]
  time: 15 minutes
---


<!-- Source: https://github.com/testcontainers/tc-guide-getting-started-with-testcontainers-for-python -->

In this guide, you will learn how to:

- Create a Python application that uses PostgreSQL to store customer data
- Use `psycopg` to interact with the database
- Write integration tests using `testcontainers-python` and `pytest`
- Manage container lifecycle with pytest fixtures

## Prerequisites

- Python 3.10+
- pip
- A Docker environment supported by Testcontainers

> [!NOTE]
> If you're new to Testcontainers, visit the
> [Testcontainers overview](https://testcontainers.com/getting-started/) to learn more about
> Testcontainers and the benefits of using it.

## Create the Python project

### Initialize the project

Start by creating a Python project with a virtual environment:

```console
$ mkdir tc-python-demo
$ cd tc-python-demo
$ python3 -m venv venv
$ source venv/bin/activate
```

This guide uses [psycopg3](https://www.psycopg.org/psycopg3/) to interact
with the Postgres database, [pytest](https://pytest.org/) for testing, and
[testcontainers-python](https://testcontainers-python.readthedocs.io/) for
running a PostgreSQL database in a container.

Install the dependencies:

```console
$ pip install "psycopg[binary]" pytest testcontainers[postgres]
$ pip freeze > requirements.txt
```

The `pip freeze` command generates a `requirements.txt` file so that others
can install the same package versions using `pip install -r requirements.txt`.

### Create the database helper

Create a `db/connection.py` file with a function to get a database connection:

```python
import os

import psycopg


def get_connection():
    host = os.getenv("DB_HOST", "localhost")
    port = os.getenv("DB_PORT", "5432")
    username = os.getenv("DB_USERNAME", "postgres")
    password = os.getenv("DB_PASSWORD", "postgres")
    database = os.getenv("DB_NAME", "postgres")
    return psycopg.connect(f"host={host} dbname={database} user={username} password={password} port={port}")
```

Instead of hard-coding the database connection parameters, the function uses
environment variables. This makes it possible to run the application in
different environments without changing code.

### Create the business logic

Create a `customers/customers.py` file and define the `Customer` class:

```python
class Customer:
    def __init__(self, cust_id, name, email):
        self.id = cust_id
        self.name = name
        self.email = email

    def __str__(self):
        return f"Customer({self.id}, {self.name}, {self.email})"
```

Add a `create_table()` function to create the `customers` table:

```python
from db.connection import get_connection


def create_table():
    with get_connection() as conn:
        with conn.cursor() as cur:
            cur.execute("""
                CREATE TABLE customers (
                    id serial PRIMARY KEY,
                    name varchar not null,
                    email varchar not null unique)
                """)
            conn.commit()
```

The function obtains a database connection using `get_connection()` and creates
the `customers` table. The `with` statement automatically closes the connection
when done.

Add the remaining CRUD functions:

```python
def create_customer(name, email):
    with get_connection() as conn:
        with conn.cursor() as cur:
            cur.execute(
                "INSERT INTO customers (name, email) VALUES (%s, %s)", (name, email))
            conn.commit()


def get_all_customers() -> list[Customer]:
    with get_connection() as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT * FROM customers")
            return [Customer(cid, name, email) for cid, name, email in cur]


def get_customer_by_email(email) -> Customer:
    with get_connection() as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT id, name, email FROM customers WHERE email = %s", (email,))
            (cid, name, email) = cur.fetchone()
            return Customer(cid, name, email)


def delete_all_customers():
    with get_connection() as conn:
        with conn.cursor() as cur:
            cur.execute("DELETE FROM customers")
            conn.commit()
```

> [!NOTE]
> To keep it straightforward for this guide, each function creates a new
> connection. In a real-world application, use a connection pool to reuse
> connections.

## Write tests with Testcontainers

You'll create a PostgreSQL container using Testcontainers and use it for all
the tests. Before each test, you'll delete all customer records so that tests
run with a clean database.

### Set up pytest fixtures

This guide uses [pytest fixtures](https://pytest.org/en/stable/how-to/fixtures.html)
for setup and teardown logic. A recommended approach is to use
[finalizers](https://pytest.org/en/stable/how-to/fixtures.html#adding-finalizers-directly)
to guarantee cleanup runs even if setup fails:

```python
@pytest.fixture
def setup(request):
    # setup code

    def cleanup():
        # teardown code

    request.addfinalizer(cleanup)
    return some_value
```

### Create the test file

Create a `tests/__init__.py` file with empty content to enable pytest
[auto-discovery](https://pytest.org/explanation/goodpractices.html#test-discovery).

Then create `tests/test_customers.py` with the fixtures:

```python
import os
import pytest
from testcontainers.postgres import PostgresContainer

from customers import customers

postgres = PostgresContainer("postgres:16-alpine")


@pytest.fixture(scope="module", autouse=True)
def setup(request):
    postgres.start()

    def remove_container():
        postgres.stop()

    request.addfinalizer(remove_container)
    os.environ["DB_CONN"] = postgres.get_connection_url()
    os.environ["DB_HOST"] = postgres.get_container_host_ip()
    os.environ["DB_PORT"] = str(postgres.get_exposed_port(5432))
    os.environ["DB_USERNAME"] = postgres.username
    os.environ["DB_PASSWORD"] = postgres.password
    os.environ["DB_NAME"] = postgres.dbname
    customers.create_table()


@pytest.fixture(scope="function", autouse=True)
def setup_data():
    customers.delete_all_customers()
```

Here's what the fixtures do:

- The `setup` fixture has `scope="module"`, so it runs once for all tests in
  the file. It starts a PostgreSQL container, sets environment variables with
  the connection details, and creates the `customers` table. A cleanup
  function removes the container after all tests complete.
- The `setup_data` fixture has `scope="function"`, so it runs before every
  test. It deletes all records to give each test a clean database.

### Write the tests

Add the test functions to the same file:

```python
def test_get_all_customers():
    customers.create_customer("Siva", "siva@gmail.com")
    customers.create_customer("James", "james@gmail.com")
    customers_list = customers.get_all_customers()
    assert len(customers_list) == 2


def test_get_customer_by_email():
    customers.create_customer("John", "john@gmail.com")
    customer = customers.get_customer_by_email("john@gmail.com")
    assert customer.name == "John"
    assert customer.email == "john@gmail.com"
```

- `test_get_all_customers()` inserts two customer records, fetches all
  customers, and asserts the count.
- `test_get_customer_by_email()` inserts a customer, fetches it by email, and
  asserts the details.

Because `setup_data` deletes all records before each test, the tests can run in
any order.

## Run tests and next steps

### Run the tests

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

### Summary

The Testcontainers for Python library helps you write integration tests using the
same type of database (Postgres) that you use in production, instead of mocks.
Because you aren't using mocks and instead talk to real services, you're free
to refactor code and still verify that the application works as expected.

In addition to PostgreSQL, Testcontainers for Python provides modules for many
SQL databases, NoSQL databases, messaging queues, and more. You can use
Testcontainers to run any containerized dependency for your tests.

To learn more about Testcontainers, visit the
[Testcontainers overview](https://testcontainers.com/getting-started/).

### Further reading

- [testcontainers-python documentation](https://testcontainers-python.readthedocs.io/)
- [Getting started with Testcontainers for Go](/guides/testcontainers-go-getting-started/)
- [Getting started with Testcontainers for Java](https://testcontainers.com/guides/getting-started-with-testcontainers-for-java/)
- [Getting started with Testcontainers for Node.js](https://testcontainers.com/guides/getting-started-with-testcontainers-for-nodejs/)
