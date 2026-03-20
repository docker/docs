---
title: Create the Python project
linkTitle: Create the project
description: Set up a Python project with a PostgreSQL-backed customer service.
weight: 10
---

## Initialize the project

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
$ pip install psycopg pytest testcontainers[postgres]
$ pip freeze > requirements.txt
```

The `pip freeze` command generates a `requirements.txt` file so that others
can install the same package versions using `pip install -r requirements.txt`.

## Create the database helper

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

## Create the business logic

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
