---
title: Write tests with Testcontainers
linkTitle: Write tests
description: Write integration tests using testcontainers-python and pytest with a real PostgreSQL database.
weight: 20
---

You'll create a PostgreSQL container using Testcontainers and use it for all
the tests. Before each test, you'll delete all customer records so that tests
run with a clean database.

## Set up pytest fixtures

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

## Create the test file

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

## Write the tests

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
