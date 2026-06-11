---
title: Ruby on Rails language-specific guide
linkTitle: Ruby
description: Containerize Ruby on Rails apps using Docker
keywords: Docker, getting started, ruby, language
summary: |
  This guide explains how to containerize Ruby on Rails applications using
  Docker.
aliases:
  - /language/ruby/
  - /guides/language/ruby/
  - /language/ruby/build-images/
  - /language/ruby/run-containers/
  - /language/ruby/containerize/
  - /language/ruby/configure-ci-cd/
  - /guides/language/ruby/configure-ci-cd/
  - /language/ruby/develop/
  - /language/ruby/deploy/
  - /guides/ruby/configure-github-actions/
  - /guides/ruby/containerize/
  - /guides/ruby/deploy/
  - /guides/ruby/develop/
params:
  tags: [cicd]
  time: 20 minutes
---


The Ruby language-specific guide teaches you how to containerize a Ruby on Rails application using Docker. In this guide, you’ll learn how to:

- Containerize and run a Ruby on Rails application
- Configure a GitHub Actions workflow to build and push a Docker image to Docker Hub
- Set up a local environment to develop a Ruby on Rails application using containers
- Deploy your containerized Ruby on Rails application locally to Kubernetes to test and debug your deployment

Start by containerizing an existing Ruby on Rails application.

## Containerize a Ruby on Rails application

### Prerequisites

- You have installed the latest version of [Docker Desktop](/get-started/get-docker.md).
- You have a [Git client](https://git-scm.com/downloads). The examples in this section show the Git CLI, but you can use any client.

### Overview

This section walks you through containerizing and running a [Ruby on Rails](https://rubyonrails.org/) application.

Starting from Rails 7.1 [Docker is supported out of the box](https://guides.rubyonrails.org/7_1_release_notes.html#generate-dockerfiles-for-new-rails-applications). This means that you will get a `Dockerfile`, `.dockerignore` and `bin/docker-entrypoint` files generated for you when you create a new Rails application.

If you have an existing Rails application, you will need to create the Docker assets manually from the examples below.

### 1. Create Docker assets

> [!TIP]
>
> [Gordon](/ai/gordon/), Docker's AI assistant, can generate Docker assets for your project. Ask Gordon to create a Dockerfile, Compose file, and `.dockerignore` tailored to your application.

Rails 7.1 and newer generates multistage Dockerfile out of the box. Following are two versions of such a file: one using Docker Hardened Images (DHIs) and another using the Docker Official Image (DOIs). Although the Dockerfile is generated automatically, understanding its purpose and functionality is important. Reviewing the following example is highly recommended.

[Docker Hardened Images (DHIs)](https://docs.docker.com/dhi/) are minimal, secure, and production-ready container base and application images maintained by Docker. DHIs are recommended whenever it is possible for better security. They are designed to reduce vulnerabilities and simplify compliance, freely available to everyone with no subscription required, no usage restrictions, and no vendor lock-in.

Multistage Dockerfiles help create smaller, more efficient images by separating build and runtime dependencies, ensuring only necessary components are included in the final image. Read more in the [Multi-stage builds guide](/get-started/docker-concepts/building-images/multi-stage-builds/).



{{< tabs >}}
{{< tab name="Using DHIs" >}}

You must authenticate to `dhi.io` before you can pull Docker Hardened Images. Run `docker login dhi.io` to authenticate.

```dockerfile {title=Dockerfile}
# syntax=docker/dockerfile:1
# check=error=true

# This Dockerfile is designed for production, not development.
# docker build -t app .
# docker run -d -p 80:80 -e RAILS_MASTER_KEY=<value from config/master.key> --name app app

# For a containerized dev environment, see Dev Containers: https://guides.rubyonrails.org/getting_started_with_devcontainer.html

# Make sure RUBY_VERSION matches the Ruby version in .ruby-version
ARG RUBY_VERSION=3.4.8
FROM dhi.io/ruby:$RUBY_VERSION-dev AS base

# Rails app lives here
WORKDIR /rails

# Install base packages
# Replace libpq-dev with sqlite3 if using SQLite, or libmysqlclient-dev if using MySQL
RUN apt-get update -qq && \
    apt-get install --no-install-recommends -y curl libjemalloc2 libvips libpq-dev && \
    rm -rf /var/lib/apt/lists /var/cache/apt/archives

# Set production environment
ENV RAILS_ENV="production" \
    BUNDLE_DEPLOYMENT="1" \
    BUNDLE_PATH="/usr/local/bundle" \
    BUNDLE_WITHOUT="development"

# Throw-away build stage to reduce size of final image
FROM base AS build

# Install packages needed to build gems
RUN apt-get update -qq && \
    apt-get install --no-install-recommends -y build-essential curl git pkg-config libyaml-dev && \
    rm -rf /var/lib/apt/lists /var/cache/apt/archives

# Install JavaScript dependencies and Node.js for asset compilation
#
# Uncomment the following lines if you are using NodeJS need to compile assets
#
# ARG NODE_VERSION=18.12.0
# ARG YARN_VERSION=1.22.19
# ENV PATH=/usr/local/node/bin:$PATH
# RUN curl -sL https://github.com/nodenv/node-build/archive/master.tar.gz | tar xz -C /tmp/ && \
#     /tmp/node-build-master/bin/node-build "${NODE_VERSION}" /usr/local/node && \
#     npm install -g yarn@$YARN_VERSION && \
#     npm install -g mjml && \
#     rm -rf /tmp/node-build-master

# Install application gems
COPY Gemfile Gemfile.lock ./
RUN bundle install && \
    rm -rf ~/.bundle/ "${BUNDLE_PATH}"/ruby/*/cache "${BUNDLE_PATH}"/ruby/*/bundler/gems/*/.git && \
    bundle exec bootsnap precompile --gemfile

# Install node modules
#
# Uncomment the following lines if you are using NodeJS need to compile assets
#
# COPY package.json yarn.lock ./
# RUN --mount=type=cache,id=yarn,target=/rails/.cache/yarn YARN_CACHE_FOLDER=/rails/.cache/yarn \
#     yarn install --frozen-lockfile

# Copy application code
COPY . .

# Precompile bootsnap code for faster boot times
RUN bundle exec bootsnap precompile app/ lib/

# Precompiling assets for production without requiring secret RAILS_MASTER_KEY
RUN SECRET_KEY_BASE_DUMMY=1 ./bin/rails assets:precompile

# Final stage for app image
FROM base

# Copy built artifacts: gems, application
COPY --from=build "${BUNDLE_PATH}" "${BUNDLE_PATH}"
COPY --from=build /rails /rails

# Run and own only the runtime files as a non-root user for security
RUN groupadd --system --gid 1000 rails && \
    useradd rails --uid 1000 --gid 1000 --create-home --shell /bin/bash && \
    chown -R rails:rails db log storage tmp
USER 1000:1000

# Entrypoint prepares the database.
ENTRYPOINT ["/rails/bin/docker-entrypoint"]

# Start server via Thruster by default, this can be overwritten at runtime
EXPOSE 80
CMD ["./bin/thrust", "./bin/rails", "server"]
```

{{< /tab >}}
{{< tab name="Using DOIs" >}}

```dockerfile {title=Dockerfile}
# syntax=docker/dockerfile:1
# check=error=true

# This Dockerfile is designed for production, not development.
# docker build -t app .
# docker run -d -p 80:80 -e RAILS_MASTER_KEY=<value from config/master.key> --name app app

# For a containerized dev environment, see Dev Containers: https://guides.rubyonrails.org/getting_started_with_devcontainer.html

# Make sure RUBY_VERSION matches the Ruby version in .ruby-version
ARG RUBY_VERSION=3.4.8
FROM docker.io/library/ruby:$RUBY_VERSION-slim AS base

# Rails app lives here
WORKDIR /rails

# Install base packages
# Replace libpq-dev with sqlite3 if using SQLite, or libmysqlclient-dev if using MySQL
RUN apt-get update -qq && \
    apt-get install --no-install-recommends -y curl libjemalloc2 libvips libpq-dev && \
    rm -rf /var/lib/apt/lists /var/cache/apt/archives

# Set production environment
ENV RAILS_ENV="production" \
    BUNDLE_DEPLOYMENT="1" \
    BUNDLE_PATH="/usr/local/bundle" \
    BUNDLE_WITHOUT="development"

# Throw-away build stage to reduce size of final image
FROM base AS build

# Install packages needed to build gems
RUN apt-get update -qq && \
    apt-get install --no-install-recommends -y build-essential curl git pkg-config libyaml-dev && \
    rm -rf /var/lib/apt/lists /var/cache/apt/archives

# Install JavaScript dependencies and Node.js for asset compilation
#
# Uncomment the following lines if you are using NodeJS need to compile assets
#
# ARG NODE_VERSION=18.12.0
# ARG YARN_VERSION=1.22.19
# ENV PATH=/usr/local/node/bin:$PATH
# RUN curl -sL https://github.com/nodenv/node-build/archive/master.tar.gz | tar xz -C /tmp/ && \
#     /tmp/node-build-master/bin/node-build "${NODE_VERSION}" /usr/local/node && \
#     npm install -g yarn@$YARN_VERSION && \
#     npm install -g mjml && \
#     rm -rf /tmp/node-build-master

# Install application gems
COPY Gemfile Gemfile.lock ./
RUN bundle install && \
    rm -rf ~/.bundle/ "${BUNDLE_PATH}"/ruby/*/cache "${BUNDLE_PATH}"/ruby/*/bundler/gems/*/.git && \
    bundle exec bootsnap precompile --gemfile

# Install node modules
#
# Uncomment the following lines if you are using NodeJS need to compile assets
#
# COPY package.json yarn.lock ./
# RUN --mount=type=cache,id=yarn,target=/rails/.cache/yarn YARN_CACHE_FOLDER=/rails/.cache/yarn \
#     yarn install --frozen-lockfile

# Copy application code
COPY . .

# Precompile bootsnap code for faster boot times
RUN bundle exec bootsnap precompile app/ lib/

# Precompiling assets for production without requiring secret RAILS_MASTER_KEY
RUN SECRET_KEY_BASE_DUMMY=1 ./bin/rails assets:precompile

# Final stage for app image
FROM base

# Copy built artifacts: gems, application
COPY --from=build "${BUNDLE_PATH}" "${BUNDLE_PATH}"
COPY --from=build /rails /rails

# Run and own only the runtime files as a non-root user for security
RUN groupadd --system --gid 1000 rails && \
    useradd rails --uid 1000 --gid 1000 --create-home --shell /bin/bash && \
    chown -R rails:rails db log storage tmp
USER 1000:1000

# Entrypoint prepares the database.
ENTRYPOINT ["/rails/bin/docker-entrypoint"]

# Start server via Thruster by default, this can be overwritten at runtime
EXPOSE 80
CMD ["./bin/thrust", "./bin/rails", "server"]
```

{{< /tab >}}
{{< /tabs >}}

The Dockerfile above assumes you are using Thruster together with Puma as an application server. In case you are using any other server, you can replace the last three lines with the following:

```dockerfile
# Start the application server
EXPOSE 3000
CMD ["./bin/rails", "server"]
```

This Dockerfile uses a script at `./bin/docker-entrypoint` as the container's entrypoint. This script prepares the database and runs the application server. Below is an example of such a script.

```bash {title=docker-entrypoint}
#!/bin/bash -e

# Enable jemalloc for reduced memory usage and latency.
if [ -z "${LD_PRELOAD+x}" ]; then
    LD_PRELOAD=$(find /usr/lib -name libjemalloc.so.2 -print -quit)
    export LD_PRELOAD
fi

# If running the rails server then create or migrate existing database
if [ "${@: -2:1}" == "./bin/rails" ] && [ "${@: -1:1}" == "server" ]; then
  ./bin/rails db:prepare
fi

exec "${@}"
```

Besides the two files above you will also need a `.dockerignore` file. This file is used to exclude files and directories from the context of the build. Below is an example of a `.dockerignore` file.

```text {collapse=true,title=".dockerignore"}
# See https://docs.docker.com/engine/reference/builder/#dockerignore-file for more about ignoring files.

# Ignore git directory.
/.git/
/.gitignore

# Ignore bundler config.
/.bundle

# Ignore all environment files.
/.env*

# Ignore all default key files.
/config/master.key
/config/credentials/*.key

# Ignore all logfiles and tempfiles.
/log/*
/tmp/*
!/log/.keep
!/tmp/.keep

# Ignore pidfiles, but keep the directory.
/tmp/pids/*
!/tmp/pids/.keep

# Ignore storage (uploaded files in development and any SQLite databases).
/storage/*
!/storage/.keep
/tmp/storage/*
!/tmp/storage/.keep

# Ignore assets.
/node_modules/
/app/assets/builds/*
!/app/assets/builds/.keep
/public/assets

# Ignore CI service files.
/.github

# Ignore development files
/.devcontainer

# Ignore Docker-related files
/.dockerignore
/Dockerfile*
```

The last optional file that you may want is the `compose.yaml` file, which is used by Docker Compose to define the services that make up the application. Since SQLite is being used as the database, there is no need to define a separate service for the database. The only service required is the Rails application itself.

```yaml {title=compose.yaml}
services:
  web:
    build: .
    environment:
      - RAILS_MASTER_KEY
    ports:
      - "3000:80"
```

You should now have the following files in your application folder:

- `.dockerignore`
- `compose.yaml`
- `Dockerfile`
- `bin/docker-entrypoint`

To learn more about the files, see the following:

- [Dockerfile](/reference/dockerfile)
- [.dockerignore](/reference/dockerfile#dockerignore-file)
- [compose.yaml](/reference/compose-file/_index.md)
- [docker-entrypoint](/reference/dockerfile/#entrypoint)

### 2. Run the application

To run the application, run the following command in a terminal inside the application's directory.

```console
$ RAILS_MASTER_KEY=<master_key_value> docker compose up --build
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000). You should see a simple Ruby on Rails application.

In the terminal, press `ctrl`+`c` to stop the application.

### 3. Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `docker-ruby-on-rails` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000).

You should see a simple Ruby on Rails application.

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the [Compose CLI
reference](/reference/cli/docker/compose/).

### Summary

In this section, you learned how you can containerize and run your Ruby
application using Docker.

Related information:

- [Docker Compose overview](/manuals/compose/_index.md)

### Next steps

In the next section, you'll take a look at how to set up a CI/CD pipeline using GitHub Actions.

## Automate your builds with GitHub Actions

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Ruby on Rails application](containerize.md). You must have a [GitHub](https://github.com/signup) account and a verified [Docker](https://hub.docker.com/signup) account to complete this section.

If you didn't create a [GitHub repository](https://github.com/new) for your project yet, it is time to do it. After creating the repository, don't forget to [add a remote](https://docs.github.com/en/get-started/getting-started-with-git/managing-remote-repositories) and ensure you can commit and [push your code](https://docs.github.com/en/get-started/using-git/pushing-commits-to-a-remote-repository#about-git-push) to GitHub.

1. In your project's GitHub repository, open **Settings**, and go to **Secrets and variables** > **Actions**.

2. Under the **Variables** tab, create a new **Repository variable** named `DOCKER_USERNAME` and your Docker ID as a value.

3. Create a new [Personal Access Token (PAT)](/manuals/security/access-tokens.md#create-an-access-token) for Docker Hub. You can name this token `docker-tutorial`. Make sure access permissions include Read and Write.

4. Add the PAT as a **Repository secret** in your GitHub repository, with the name
   `DOCKERHUB_TOKEN`.

### Overview

GitHub Actions is a CI/CD (Continuous Integration and Continuous Deployment) automation tool built into GitHub. It allows you to define custom workflows for building, testing, and deploying your code when specific events occur (e.g., pushing code, creating a pull request, etc.). A workflow is a YAML-based automation script that defines a sequence of steps to be executed when triggered. Workflows are stored in the `.github/workflows/` directory of a repository.

In this section, you'll learn how to set up and use GitHub Actions to build your Docker image as well as push it to Docker Hub. You will complete the following steps:

1. Define the GitHub Actions workflow.
2. Run the workflow.

### 1. Define the GitHub Actions workflow

You can create a GitHub Actions workflow by creating a YAML file in the `.github/workflows/` directory of your repository. To do this use your favorite text editor or the GitHub web interface. The following steps show you how to create a workflow file using the GitHub web interface.

If you prefer to use the GitHub web interface, follow these steps:

1. Go to your repository on GitHub and then select the **Actions** tab.

2. Select **set up a workflow yourself**.

   This takes you to a page for creating a new GitHub Actions workflow file in
   your repository. By default, the file is created under `.github/workflows/main.yml`, let's change it name to `build.yml`.

If you prefer to use your text editor, create a new file named `build.yml` in the `.github/workflows/` directory of your repository.

Add the following content to the file:

```yaml
name: Build and push Docker image

on:
  push:
    branches:
      - main

jobs:
  build_and_push:
    runs-on: ubuntu-latest
    steps:
      - name: Login to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Build and push
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          push: true
          tags: ${{ vars.DOCKER_USERNAME }}/${{ github.event.repository.name }}:latest
```

Each GitHub Actions workflow includes one or several jobs. Each job consists of steps. Each step can either run a set of commands or use already [existing actions](https://github.com/marketplace?type=actions). The action above has three steps:

1. [Login to Docker Hub](https://github.com/docker/login-action): Action logs in to Docker Hub using the Docker ID and Personal Access Token (PAT) you created earlier.

2. [Set up Docker Buildx](https://github.com/docker/setup-buildx-action): Action sets up Docker [Buildx](https://github.com/docker/buildx), a CLI plugin that extends the capabilities of the Docker CLI.

3. [Build and push](https://github.com/docker/build-push-action): Action builds and pushes the Docker image to Docker Hub. The `tags` parameter specifies the image name and tag. The `latest` tag is used in this example.

### 2. Run the workflow

Commit the changes and push them to the `main` branch. This workflow is runs every time you push changes to the `main` branch. You can find more information about workflow triggers [in the GitHub documentation](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows).

Go to the **Actions** tab of you GitHub repository. It displays the workflow. Selecting the workflow shows you the breakdown of all the steps.

When the workflow is complete, go to your [repositories on Docker Hub](https://hub.docker.com/repositories). If you see the new repository in that list, it means the GitHub Actions workflow successfully pushed the image to Docker Hub.

### Summary

In this section, you learned how to set up a GitHub Actions workflow for your Ruby on Rails application.

Related information:

- [Introduction to GitHub Actions](/guides/gha.md)
- [Docker Build GitHub Actions](/manuals/build/ci/github-actions/_index.md)
- [Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

### Next steps

In the next section, you'll learn how you can develop your application using containers.

## Use containers for Ruby on Rails development

### Prerequisites

Complete [Containerize a Ruby on Rails application](containerize.md).

### Overview

In this section, you'll learn how to set up a development environment for your containerized application. This includes:

- Adding a local database and persisting data
- Configuring Compose to automatically update your running Compose services as you edit and save your code

### Add a local database and persist data

You can use containers to set up local services, like a database. In this section, you'll update the `compose.yaml` file to define a database service and a volume to persist data.

In the cloned repository's directory, open the `compose.yaml` file in an IDE or text editor. You need to add the database password file as an environment variable to the server service and specify the secret file to use.

The following is the updated `compose.yaml` file.

```yaml {hl_lines="07-25"}
services:
  web:
    build: .
    command: bundle exec rails s -b '0.0.0.0'
    ports:
      - "3000:3000"
    depends_on:
      - db
    environment:
      - RAILS_ENV=test
    env_file: "webapp.env"
  db:
    image: postgres:18
    secrets:
      - db-password
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    volumes:
      - postgres_data:/var/lib/postgresql

volumes:
  postgres_data:
secrets:
  db-password:
    file: db/password.txt
```

> [!NOTE]
>
> To learn more about the instructions in the Compose file, see [Compose file
> reference](/reference/compose-file/).

Before you run the application using Compose, notice that this Compose file specifies a `password.txt` file to hold the database's password. You must create this file as it's not included in the source repository.

In the cloned repository's directory, create a new directory named `db` and inside that directory create a file named `password.txt` that contains the password for the database. Using your favorite IDE or text editor, add the following contents to the `password.txt` file.

```text
mysecretpassword
```

Save and close the `password.txt` file. In addition, in the file `webapp.env` you can change the password to connect to the database.

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

In Ruby on Rails, `db:migrate` is a Rake task that is used to run migrations on the database. Migrations are a way to alter the structure of your database schema over time in a consistent and easy way.

```console
$ docker exec -it docker-ruby-on-rails-web-1 rake db:migrate RAILS_ENV=test
```

You will see a similar message like this:

`console
== 20240710193146 CreateWhales: migrating =====================================
-- create_table(:whales)
   -> 0.0126s
== 20240710193146 CreateWhales: migrated (0.0127s) ============================
`

Refresh <http://localhost:3000> in your browser and add the whales.

Press `ctrl+c` in the terminal to stop your application and run `docker compose up` again, the whales are being persisted.

### Automatically update services

Use Compose Watch to automatically update your running Compose services as you
edit and save your code. For more details about Compose Watch, see [Use Compose
Watch](/manuals/compose/how-tos/file-watch.md).

Open your `compose.yaml` file in an IDE or text editor and then add the Compose
Watch instructions. The following is the updated `compose.yaml` file.

```yaml {hl_lines="13-16"}
services:
  web:
    build: .
    command: bundle exec rails s -b '0.0.0.0'
    ports:
      - "3000:3000"
    depends_on:
      - db
    environment:
      - RAILS_ENV=test
    env_file: "webapp.env"

    develop:
      watch:
        - action: rebuild
          path: .
  db:
    image: postgres:18
    secrets:
      - db-password
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    volumes:
      - postgres_data:/var/lib/postgresql

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

Open `docker-ruby-on-rails/app/views/whales/index.html.erb` in an IDE or text editor and update the `Whales` string by adding an exclamation mark.

```diff
-    <h1>Whales</h1>
+    <h1>Whales!</h1>
```

Save the changes to `index.html.erb` and then wait a few seconds for the application to rebuild. Go to the application again and verify that the updated text appears.

Press `ctrl+c` in the terminal to stop your application.

### Summary

In this section, you took a look at setting up your Compose file to add a local
database and persist data. You also learned how to use Compose Watch to automatically rebuild and run your container when you update your code.

Related information:

- [Compose file reference](/reference/compose-file/)
- [Compose file watch](/manuals/compose/how-tos/file-watch.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)

### Next steps

In the next section, you'll learn how you can locally test and debug your workloads on Kubernetes before deploying.

## Test your Ruby on Rails deployment

### Prerequisites

- Complete all the previous sections of this guide, starting with [Containerize a Ruby on Rails application](containerize.md).
- [Turn on Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

### Overview

In this section, you'll learn how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine. This lets you to test and debug your workloads on Kubernetes locally before deploying.

### Create a Kubernetes YAML file

In your `docker-ruby-on-rails` directory, create a file named
`docker-ruby-on-rails-kubernetes.yaml`. Open the file in an IDE or text editor and add
the following contents. Replace `DOCKER_USERNAME/REPO_NAME` with your Docker
username and the name of the repository that you created in [Configure CI/CD for
your Ruby on Rails application](configure-github-actions.md).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-ruby-on-rails-demo
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      service: ruby-on-rails
  template:
    metadata:
      labels:
        service: ruby-on-rails
    spec:
      containers:
        - name: ruby-on-rails-container
          image: DOCKER_USERNAME/REPO_NAME
          imagePullPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  name: docker-ruby-on-rails-demo
  namespace: default
spec:
  type: NodePort
  selector:
    service: ruby-on-rails
  ports:
    - port: 3000
      targetPort: 3000
      nodePort: 30001
```

In this Kubernetes YAML file, there are two objects, separated by the `---`:

- A Deployment, describing a scalable group of identical pods. In this case,
  you'll get just one replica, or copy of your pod. That pod, which is
  described under `template`, has just one container in it. The
  container is created from the image built by GitHub Actions in [Configure CI/CD for
  your Ruby on Rails application](configure-github-actions.md).
- A NodePort service, which will route traffic from port 30001 on your host to
  port 8001 inside the pods it routes to, allowing you to reach your app
  from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

### Deploy and check your application

1. In a terminal, navigate to `docker-ruby-on-rails` and deploy your application to
   Kubernetes.

   ```console
   $ kubectl apply -f docker-ruby-on-rails-kubernetes.yaml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```shell
   deployment.apps/docker-ruby-on-rails-demo created
   service/docker-ruby-on-rails-demo created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```shell
   NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
   docker-ruby-on-rails-demo  1/1     1            1           15s
   ```

   This indicates all one of the pods you asked for in your YAML are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```shell
   NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
   kubernetes                  ClusterIP   10.96.0.1       <none>        443/TCP          23h
   docker-ruby-on-rails-demo   NodePort    10.99.128.230   <none>        3000:30001/TCP   75s
   ```

   In addition to the default `kubernetes` service, you can see your `docker-ruby-on-rails-demo` service, accepting traffic on port 30001/TCP.

3. To create and migrate the database in a Ruby on Rails application running on Kubernetes, you need to follow these steps.

   **Get the Current Pods**:
   First, you need to identify the pods running in your Kubernetes cluster. Execute the following command to list the current pods in the `default` namespace:

   ```sh
   # Get the current pods in the cluster in the namespace default
   $ kubectl get pods
   ```

   This command will display a list of all pods in the `default` namespace. Look for the pod with the prefix `docker-ruby-on-rails-demo-`. Here is an example output:

   ```console
   NAME                                         READY   STATUS    RESTARTS      AGE
   docker-ruby-on-rails-demo-7cbddb5d6f-qh44l   1/1     Running   2 (22h ago)   9d
   ```

   **Execute the Migration Command**:
   Once you've identified the correct pod, use the `kubectl exec` command to run the database migration inside the pod.

   ```sh
   $ kubectl exec -it docker-ruby-on-rails-demo-7cbddb5d6f-qh44l -- rails db:migrate RAILS_ENV=development
   ```

   This command opens an interactive terminal session (`-it`) in the specified pod and runs the `rails db:migrate` command with the environment set to development (`RAILS_ENV=development`).

   By following these steps, you ensure that your database is properly migrated within the Ruby on Rails application running in your Kubernetes cluster. This process helps maintain the integrity and consistency of your application's data structure during deployment and updates.

4. Open the browser and go to [http://localhost:30001](http://localhost:30001), you should see the ruby on rails application working.

5. Run the following command to tear down your application.

   ```console
   $ kubectl delete -f docker-ruby-on-rails-kubernetes.yaml
   ```

### Summary

In this section, you learned how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine.

Related information:

- [Kubernetes documentation](https://kubernetes.io/docs/home/)
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md)
- [Swarm mode overview](/manuals/engine/swarm/_index.md)
