---
title: Pre-seeding database with schema and data at startup for development environment
linktitle: Pre-seeding database  
description: &desc Pre-seeding database with schema and data at startup for development environment
keywords: Pre-seeding, database, postgres, container-supported development
summary: *desc
tags: [app-dev, databases]
params:
  time: 20 minutes
---

Pre-seeding databases with essential data and schema during local development is a common practice to enhance the development and testing workflow. By simulating real-world scenarios, this practice helps catch frontend issues early, ensures alignment between Database Administrators and Software Engineers, and facilitates smoother collaboration. Pre-seeding offers benefits like confident deployments, consistency across environments, and early issue detection, ultimately improving the overall development process.

In this guide, you will learn how to:

- Use Docker to launch up a Postgres container
- Pre-seed Postgres using a SQL script
- Pre-seed Postgres by copying SQL files into Docker image
- Pre-seed Postgres using JavaScript code

## Using Postgres with Docker

The [official Docker image for Postgres](https://hub.docker.com/_/postgres) provides a convenient way to run Postgres database on your development machine. A Postgres Docker image is a pre-configured environment that encapsulates the PostgreSQL database system. It's a self-contained unit, ready to run in a Docker container. By using this image, you can quickly and easily set up a Postgres instance without the need for manual configuration.

## Prerequisites

The following prerequisites are required to follow along with this how-to guide:

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) 

## Launching Postgres

Launch a quick demo of Postgres by using the following steps:

1. Open the terminal and run the following command to start a Postgres container.

   This example will launch a Postgres container, expose port `5432` onto the host to let a native-running application to connect to it with the password `mysecretpassword`.

   ```console
   $ docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword postgres
   ```

2. Verify that Postgres is up and running by selecting the container and checking the logs on Docker Dashboard.

   ```plaintext
   PostgreSQL Database directory appears to contain a database; Skipping initialization
 
   2024-09-08 09:09:47.136 UTC [1] LOG:  starting PostgreSQL 16.4 (Debian 16.4-1.pgdg120+1) on aarch64-unknown-linux-gnu, compiled by gcc (Debian 12.2.0-14) 12.2.0, 64-bit
   2024-09-08 09:09:47.137 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
   2024-09-08 09:09:47.137 UTC [1] LOG:  listening on IPv6 address "::", port 5432
   2024-09-08 09:09:47.139 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
   2024-09-08 09:09:47.142 UTC [29] LOG:  database system was shut down at 2024-09-08 09:07:09 UTC
   2024-09-08 09:09:47.148 UTC [1] LOG:  database system is ready to accept connections
   ```

3. Connect to Postgres from the local system.

   The `psql` is the PostgreSQL interactive shell that is used to connect to a Postgres database and let you start executing SQL commands. Assuming that you already have `psql` utility installed on your local system, it's time to connect to the Postgres database. Run the following command on your local terminal:

   ```console
   $ docker exec -it postgres psql -h localhost -U postgres
   ```

   You can now execute any SQL queries or commands you need within the `psql` prompt.

   Use `\q` or `\quit` to exit from the Postgres interactive shell.

## Pre-seed the Postgres database using a SQL script

Now that you've familiarized yourself with Postgres, it's time to see how to pre-seed it with sample data. In this demonstration, you'll first create a script that holds SQL commands. The script defines the database, and table structure and inserts sample data. Then you will connect the database to verify the data.

Assuming that you have an existing Postgres database instance up and running, follow these steps to seed the database.

1. Create an empty file named `seed.sql` and add the following content.

   ```sql
   CREATE DATABASE sampledb;

   \c sampledb

   CREATE TABLE users (
     id SERIAL PRIMARY KEY,
     name VARCHAR(50),
     email VARCHAR(100) UNIQUE
   );

   INSERT INTO users (name, email) VALUES
     ('Alpha', 'alpha@example.com'),
     ('Beta', 'beta@example.com'),
     ('Gamma', 'gamma@example.com');  
   ```

   The SQL script creates a new database called `sampledb`, connects to it, and creates a `users` table. The table includes an auto-incrementing `id` as the primary key, a `name` field with a maximum length of 50 characters, and a unique `email` field with up to 100 characters.

   After creating the table, the `INSERT` command inserts three users into the `users` table with their respective names and emails. This setup forms a basic database structure to store user information with unique email addresses.

2. Seed the database.

   It’s time to feed the content of the `seed.sql` directly into the database by using the `<` operator. The command is used to execute a SQL script named `seed.sql` against a Postgres database named `sampledb`. 

   ```console
   $ cat seed.sql | docker exec -i postgres psql -h localhost -U postgres -f-
   ```

   Once the query is executed, you will see the following results:

   ```plaintext
   CREATE DATABASE
   You are now connected to database "sampledb" as user "postgres".
   CREATE TABLE
   INSERT 0 3
   ```

3. Run the following `psql` command to verify if the table named users is populated in the database `sampledb` or not. 

   ```console
   $ docker exec -it postgres psql -h localhost -U postgres sampledb
   ```

   You can now run `\l` in the `psql` shell to list all the databases on the Postgres server.

   ```console
   sampledb=# \l
                                                List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    | ICU Locale | Locale Provider |   Access privileges
   -----------+----------+----------+------------+------------+------------+-----------------+-----------------------
   postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 |            | libc            |
   sampledb  | postgres | UTF8     | en_US.utf8 | en_US.utf8 |            | libc            |
   template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 |            | libc            | =c/postgres          +
             |          |          |            |            |            |                 | postgres=CTc/postgres
   template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 |            | libc            | =c/postgres          +
             |          |          |            |            |            |                 | postgres=CTc/postgres
   (4 rows)
   ```

   To retrieve all the data from the users table, enter the following query:

   ```console
   sampledb=# SELECT * FROM users;
   id | name  |       email
   ----+-------+-------------------
    1 | Alpha | alpha@example.com
    2 | Beta  | beta@example.com
    3 | Gamma | gamma@example.com
   (3 rows)
   ```
  
   Use `\q` or `\quit` to exit from the Postgres interactive shell.

## Pre-seed the database by bind-mounting a SQL script

In Docker, mounting refers to making files or directories from the host system accessible within a container. This let you to share data or configuration files between the host and the container, enabling greater flexibility and persistence.

Now that you have learned how to launch Postgres and pre-seed the database using an SQL script, it’s time to learn how to mount an SQL file directly into the Postgres containers’ initialization directory (`/docker-entrypoint-initdb.d`). The `/docker-entrypoint-initdb.d` is a special directory in PostgreSQL Docker containers that is used for initializing the database when the container is first started

Make sure you stop any running Postgres containers (along with volumes) to prevent port conflicts before you follow the steps:

```console
$ docker container stop postgres
```

1. Modify the `seed.sql` with the following entries:

   ```sql
   CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    email VARCHAR(100) UNIQUE
   );

   INSERT INTO users (name, email) VALUES
    ('Alpha', 'alpha@example.com'),
    ('Beta', 'beta@example.com'),
    ('Gamma', 'gamma@example.com')
   ON CONFLICT (email) DO NOTHING;
   ```
   
2. Create a text file named `Dockerfile` and copy the following content.

   ```plaintext
   # syntax=docker/dockerfile:1
   FROM postgres:latest
   COPY seed.sql /docker-entrypoint-initdb.d/
   ```

   This Dockerfile copies the `seed.sql` script directly into the PostgreSQL container's initialization directory.
   

3. Use Docker Compose.
   
   Using Docker Compose makes it even easier to manage and deploy the PostgreSQL container with the seeded database. This compose.yml file defines a Postgres service named `db` using the latest Postgres image, which sets up a database with the name `sampledb`, along with a user `postgres` and a password `mysecretpassword`. 

   ```yaml
   services:
     db:
       build:
         context: .
         dockerfile: Dockerfile
       container_name: my_postgres_db
       environment:
         POSTGRES_USER: postgres
         POSTGRES_PASSWORD: mysecretpassword
         POSTGRES_DB: sampledb
       ports:
         - "5432:5432"
       volumes:
         - data_sql:/var/lib/postgresql/data   # Persistent data storage

   volumes:
     data_sql:
    ```
  
    It maps port `5432` on the host to the container's `5432`, let you access to the Postgres database from outside the container. It also define `data_sql` for persisting the database data, ensuring that data is not lost when the container is stopped.

    It is important to note that the port mapping to the host is only necessary if you want to connect to the database from non-containerized programs. If you containerize the service that connects to the DB, you should connect to the database over a custom bridge network.

4.  Bring up the Compose service.

    Assuming that you've placed the `seed.sql` file in the same directory as the Dockerfile, execute the following command:

    ```console
    $ docker compose up -d --build
    ```

5.  It’s time to verify if the table `users` get populated with the data. 

    ```console
    $ docker exec -it my_postgres_db psql -h localhost -U postgres sampledb
    ```

    ```sql 
    sampledb=# SELECT * FROM users;
      id | name  |       email
    ----+-------+-------------------
       1 | Alpha | alpha@example.com
       2 | Beta  | beta@example.com
       3 | Gamma | gamma@example.com
     (3 rows)

    sampledb=#
    ```


## Pre-seed the database using JavaScript code


Now that you have learned how to seed the database using various methods like SQL script, mounting volumes etc., it's time to try to achieve it using JavaScript code. 

1. Create a .env file with the following:

   ```plaintext
   POSTGRES_USER=postgres
   POSTGRES_DB_HOST=localhost
   POSTGRES_DB=sampledb
   POSTGRES_PASSWORD=mysecretpassword
   POSTGRES_PORT=5432
   ```

2. Create a new JavaScript file called seed.js with the following content:

   The following JavaScript code imports the `dotenv` package which is used to load environment variables from an `.env` file. The `.config()` method reads the `.env` file and sets the environment variables as properties of the `process.env` object. This let you to securely store sensitive information like database credentials outside of your code.

   Then, it creates a new Pool instance from the pg library, which provides a connection pool for efficient database interactions. The `seedData` function is defined to perform the database seeding operations.
It is called at the end of the script to initiate the seeding process. The try...catch...finally block is used for error handling. 

   ```plaintext
   require('dotenv').config();  // Load environment variables from .env file
   const { Pool } = require('pg');

   // Create a new pool using environment variables
   const pool = new Pool({
     user: process.env.POSTGRES_USER,
     host: process.env.POSTGRES_DB_HOST,
     database: process.env.POSTGRES_DB,
     port: process.env.POSTGRES_PORT,
     password: process.env.POSTGRES_PASSWORD,
   });

   const seedData = async () => {
     try {
        // Drop the table if it already exists (optional)
        await pool.query(`DROP TABLE IF EXISTS todos;`);

        // Create the table with the correct structure
        await pool.query(`
          CREATE TABLE todos (
            id SERIAL PRIMARY KEY,
            task VARCHAR(255) NOT NULL,
            completed BOOLEAN DEFAULT false
              );
        `   );

        // Insert seed data
        await pool.query(`
          INSERT INTO todos (task, completed) VALUES
          ('Watch netflix', false),
          ('Finish podcast', false),
          ('Pick up kid', false);
          `);
          console.log('Database seeded successfully!');
        } catch (err) {
          console.error('Error seeding the database', err);
        } finally {
          pool.end();
       }
     };

     // Call the seedData function to run the script
     seedData();
     ```

3.  Kick off the seeding process

    ```console
    $ node seed.js
    ```

    You should see the following command:

    ```plaintext
    Database seeded successfully!
    ```

4.  Verify if the database is seeded correctly:

    ```console
    $ docker exec -it postgres psql -h localhost -U postgres sampledb
    ```

    ```console
    sampledb=# SELECT * FROM todos;
    id |      task      | completed
    ----+----------------+-----------
    1 | Watch netflix  | f
    2 | Finish podcast | f
    3 | Pick up kid    | f
    (3 rows)  
    ```

## Recap

Pre-seeding a database with schema and data at startup is essential for creating a consistent and realistic testing environment, which helps in identifying issues early in development and aligning frontend and backend work. This guide has equipped you with the knowledge and practical steps to achieve pre-seeding using various methods, including SQL script, Docker integration, and JavaScript code. 

