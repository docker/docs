---
title: Use containers for Ruby on Rails development
keywords: ruby, local, development
description: Learn how to develop your Ruby on Rails application locally.
---

## Prerequisites

Complete [Containerize a Ruby on Rails application](containerize.md).

## Overview

In this section, you'll learn how to set up a development environment for your containerized application. This includes:

- Adding a local database and persisting data
- Configuring Compose to automatically update your running Compose services as you edit and save your code

## Add a local database and persist data

You can use containers to set up local services, like a database. In this section, you'll update the `compose.yaml` file to define a database service and a volume to persist data.

In the cloned repository's directory, open the `compose.yaml` file in an IDE or text editor. Yyou need to add the database password file as an environment variable to the server service and specify the secret file to use .

The following is the updated `compose.yaml` file.

```yaml {hl_lines="10-30"}
services:
  web:
    build: .
    command: bundle exec rails s -b '0.0.0.0'
    volumes:
      - .:/myapp
    ports:
      - "3000:3000"
    depends_on:
      - db
    secrets:
      - db-password
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
      - RAILS_ENV=test
  db:
    image: postgres:latest
    secrets:
      - db-password
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
secrets:
  db-password:
    file: db/password.txt
```

> **Note**
>
> To learn more about the instructions in the Compose file, see [Compose file
> reference](/compose/compose-file/).

Before you run the application using Compose, notice that this Compose file specifies a `password.txt` file to hold the database's password. You must create this file as it's not included in the source repository.

In the cloned repository's directory, create a new directory named `db` and inside that directory create a file named `password.txt` that contains the password for the database. Using your favorite IDE or text editor, add the following contents to the `password.txt` file.

```text
mysecretpassword
```

Save and close the `password.txt` file.

You should now have the following contents in your `docker-ruby-on-rails`
directory.

```text
.
├── Dockerfile
├── Gemfile
├── Gemfile.lock
├── README.md
├── Rakefile
├── app/
├── bin/
├── compose.yaml
├── config/
├── config.ru
├── db/
│   ├── development.sqlite3
│   ├── migrate
│   ├── password.txt
│   ├── schema.rb
│   └── seeds.rb
├── lib/
├── log/
├── public/
├── storage/
├── test/
├── tmp/
└── vendor
```

Now, run the following `docker compose up` command to start your application.

```console
$ docker compose up --build
```
Refresh <http://localhost:3000> in your browser and verify that the Whale items persisted, even after the containers were removed and ran again.u

Press `ctrl+c` in the terminal to stop your application.

## Automatically update services

Use Compose Watch to automatically update your running Compose services as you
edit and save your code. For more details about Compose Watch, see [Use Compose
Watch](../../compose/file-watch.md).

Open your `compose.yaml` file in an IDE or text editor and then add the Compose
Watch instructions. The following is the updated `compose.yaml` file.

```yaml {hl_lines="17-2 0"}
version: '3'
services:
  web:
    build: .
    command: bundle exec rails s -b '0.0.0.0'
    volumes:
      - .:/myapp
    ports:
      - "3000:3000"
    depends_on:
      - db
    secrets:
      - db-password
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
      - RAILS_ENV=test
    develop:
      watch:
        - action: rebuild
          path: .
  db:
    image: postgres:latest
    secrets:
      - db-password
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
secrets:
  db-password:
    file: db/password.txt
```

Run the following command to run your application with Compose Watch.

```console
$ docker compose watch
```

Any changes to the application's source files on your local machine will now be immediately reflected in the running container.

Open `docker-ruby-on-rails/app/views/whales/index.html.erb` in an IDE or text editor and update the `Whales` string by adding a exclamation marks.

```diff
-    <h1>Whales</h1>
+    <h1>Whales!</h1>
```

Save the changes to `index.html.erb` and then wait a few seconds for the application to rebuild. Go to the application again and verify that the updated text appears.

Press `ctrl+c` in the terminal to stop your application.

## Summary

In this section, you took a look at setting up your Compose file to add a local
database and persist data. You also learned how to use Compose Watch to automatically rebuild and run your container when you update your code.

Related information:
 - [Compose file reference](/compose/compose-file/)
 - [Compose file watch](../../compose/file-watch.md)
 - [Multi-stage builds](../../build/building/multi-stage.md)

## Next steps

In the next section, you'll take a look at how to set up a CI/CD pipeline using GitHub Actions.

{{< button text="Configure CI/CD" url="configure-ci-cd.md" >}}
