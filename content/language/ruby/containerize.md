---
title: Containerize a Ruby on Rails application
keywords: ruby, flask, containerize, initialize
description: Learn how to containerize a Ruby on Rails application.
aliases:
  - /language/ruby/build-images/
  - /language/ruby/run-containers/
---

## Prerequisites

* You have installed the latest version of [Docker Desktop](../../get-docker.md).
* You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

## Overview

This section walks you through containerizing and running a Ruby on Rails application.

## Get the sample application

The sample application uses the popular [Ruby on Rails](https://rubyonrails.org/) framework.

Clone the sample application to use with this guide. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/falconcr/docker-ruby-on-rails.git
```

## Initialize Docker assets

Now that you have an application, you can create the necessary Docker assets to
containerize your application. You can use Docker Desktop's built-in Docker Init
feature to help streamline the process, or you can manually create the assets.

Docker's docker init command provides predefined configurations tailored for specific programming languages. This feature simplifies the setup process by automatically generating Dockerfiles and other necessary configuration files based on the chosen language. For example, docker init has predefined configurations for languages like Python, Java, and Node.js, making it easier to get started with Docker for these environments.

However, it's important to note that as of now, docker init does not offer a predefined configuration for the Ruby programming language. This means that if you are working with Ruby, you'll need to create Dockerfiles and other related configurations manually.

Inside the `docker-ruby-on-rails` directory, you should create the following files:

Create a file named `Dockerfile` with the following contents.

```dockerfile {collapse=true,title=Dockerfile}
# syntax=docker/dockerfile:1

# Use the official Ruby image with version 3.2.0
FROM ruby:3.2.0

# Install dependencies
RUN apt-get update -qq && apt-get install -y \
  nodejs \
  postgresql-client \
  libssl-dev \
  libreadline-dev \
  zlib1g-dev \
  build-essential \
  curl

# Install rbenv
RUN git clone https://github.com/rbenv/rbenv.git ~/.rbenv && \
  echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.bashrc && \
  echo 'eval "$(rbenv init -)"' >> ~/.bashrc && \
  git clone https://github.com/rbenv/ruby-build.git ~/.rbenv/plugins/ruby-build && \
  echo 'export PATH="$HOME/.rbenv/plugins/ruby-build/bin:$PATH"' >> ~/.bashrc

# Install the specified Ruby version using rbenv
ENV PATH="/root/.rbenv/bin:/root/.rbenv/shims:$PATH"
RUN rbenv install 3.2.0 && rbenv global 3.2.0

# Set the working directory
WORKDIR /myapp

# Copy the Gemfile and Gemfile.lock
COPY Gemfile /myapp/Gemfile
COPY Gemfile.lock /myapp/Gemfile.lock

# Install Gems dependencies
RUN gem install bundler && bundle install

# Copy the application code
COPY . /myapp

# Precompile assets (optional, if using Rails with assets)
RUN bundle exec rake assets:precompile

# Expose the port the app runs on
EXPOSE 3000

# Command to run the server
CMD ["rails", "server", "-b", "0.0.0.0"]
```

Create a file named `compose.yaml` with the following contents.

```yaml {collapse=true,title=compose.yaml}
version: '3'
services:
  web:
    build: .
    command: bundle exec rails s -b '0.0.0.0'
    volumes:
      - .:/myapp
    ports:
      - "3000:3000"
```

Create a file named `.dockerignore` with the following contents.

```text {collapse=true,title=".dockerignore"}
# Include any files or directories that you don't want to be copied to your
# container here (e.g., local build artifacts, temporary files, etc.).
#
# For more help, visit the .dockerignore file reference guide at
# https://docs.docker.com/go/build-context-dockerignore/

# Ignore bundler config
/.bundle

# Ignore all log files and tempfiles
/log/*
/tmp/*
!/log/.keep
!/tmp/.keep

# Ignore the development and test databases
/db/*.sqlite3
/db/*.sqlite3-journal

# Ignore the production secrets file
/config/secrets.yml

# Ignore all files in the test, spec, and features folders
/test/*
/spec/*
/features/*

# Ignore system-specific files
*.swp
*.swo
*~
*.DS_Store

# Ignore coverage reports
/coverage/*

# Ignore node modules (if using a JavaScript front-end with Ruby on Rails)
/node_modules

# Ignore yarn lock file
/yarn.lock

# Ignore the .git directory and other VCS files
.git
.gitignore

**/docker-compose*
**/compose.y*ml
**/Dockerfile*
LICENSE
README.md
```

You should now have the following three files in your `docker-ruby-on-rails`
directory.


- .dockerignore
- compose.yaml
- Dockerfile


To learn more about the files, see the following:
 - [Dockerfile](../../reference/dockerfile.md)
 - [.dockerignore](../../reference/dockerfile.md#dockerignore-file)
 - [compose.yaml](../../compose/compose-file/_index.md)

## Run the application

Inside the `docker-ruby-on-rails` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:8000](http://localhost:8000). You should see a simple Ruby on Rails application.

In the terminal, press `ctrl`+`c` to stop the application.

### Run the application in the background

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
reference](../../compose/reference/_index.md).

## Summary

In this section, you learned how you can containerize and run your Ruby
application using Docker.

Related information:
 - [Build with Docker guide](../../build/guide/index.md)
 - [Docker Compose overview](../../compose/_index.md)

## Next steps

In the next section, you'll learn how you can develop your application using
containers.

{{< button text="Develop your application" url="develop.md" >}}
