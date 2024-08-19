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
* You have a [Git client](https://git-scm.com/downloads). The examples in this section show the Git CLI, but you can use any client.

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

`docker init`, the command for bootstrapping the Docker-related assets for a project, does not yet support the Ruby programming language. This means that if you are working with Ruby, you'll need to create Dockerfiles and other related configurations manually.

Inside the `docker-ruby-on-rails` directory, create the following files:

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
git
.gitignore

# Created by https://www.gitignore.io/api/git,ruby,rails,jetbrains+all
# Edit at https://www.gitignore.io/?templates=git,ruby,rails,jetbrains+all

### Git ###
# Created by git for backups. To disable backups in Git:
# $ git config --global mergetool.keepBackup false
*.orig

# Created by git when using merge tools for conflicts
*.BACKUP.*
*.BASE.*
*.LOCAL.*
*.REMOTE.*
*_BACKUP_*.txt
*_BASE_*.txt
*_LOCAL_*.txt
*_REMOTE_*.txt

### JetBrains+all ###
# Covers JetBrains IDEs: IntelliJ, RubyMine, PhpStorm, AppCode, PyCharm, CLion, Android Studio and WebStorm
# Reference: https://intellij-support.jetbrains.com/hc/en-us/articles/206544839

# User-specific stuff
.idea/**/workspace.xml
.idea/**/tasks.xml
.idea/**/usage.statistics.xml
.idea/**/dictionaries
.idea/**/shelf

# Generated files
.idea/**/contentModel.xml

# Sensitive or high-churn files
.idea/**/dataSources/
.idea/**/dataSources.ids
.idea/**/dataSources.local.xml
.idea/**/sqlDataSources.xml
.idea/**/dynamic.xml
.idea/**/uiDesigner.xml
.idea/**/dbnavigator.xml

# Gradle
.idea/**/gradle.xml
.idea/**/libraries

# Gradle and Maven with auto-import
# When using Gradle or Maven with auto-import, you should exclude module files,
# since they will be recreated, and may cause churn.  Uncomment if using
# auto-import.
# .idea/modules.xml
# .idea/*.iml
# .idea/modules
# *.iml
# *.ipr

# CMake
cmake-build-*/

# Mongo Explorer plugin
.idea/**/mongoSettings.xml

# File-based project format
*.iws

# IntelliJ
out/

# mpeltonen/sbt-idea plugin
.idea_modules/

# JIRA plugin
atlassian-ide-plugin.xml

# Cursive Clojure plugin
.idea/replstate.xml

# Crashlytics plugin (for Android Studio and IntelliJ)
com_crashlytics_export_strings.xml
crashlytics.properties
crashlytics-build.properties
fabric.properties

# Editor-based Rest Client
.idea/httpRequests

# Android studio 3.1+ serialized cache file
.idea/caches/build_file_checksums.ser

### JetBrains+all Patch ###
# Ignores the whole .idea folder and all .iml files
# See https://github.com/joeblau/gitignore.io/issues/186 and https://github.com/joeblau/gitignore.io/issues/360

.idea/

# Reason: https://github.com/joeblau/gitignore.io/issues/186#issuecomment-249601023

*.iml
modules.xml
.idea/misc.xml
*.ipr

# Sonarlint plugin
.idea/sonarlint

### Rails ###
*.rbc
capybara-*.html
.rspec
/db/*.sqlite3
/db/*.sqlite3-journal
/public/system
/coverage/
/spec/tmp
rerun.txt
pickle-email-*.html

# Ignore all logfiles and tempfiles.
/log/*
/tmp/*
!/log/.keep
!/tmp/.keep

# TODO Comment out this rule if you are OK with secrets being uploaded to the repo
config/initializers/secret_token.rb
config/master.key

# Only include if you have production secrets in this file, which is no longer a Rails default
# config/secrets.yml

# dotenv
# TODO Comment out this rule if environment variables can be committed
.env

## Environment normalization:
/.bundle
/vendor/bundle

# these should all be checked in to normalize the environment:
# Gemfile.lock, .ruby-version, .ruby-gemset

# unless supporting rvm < 1.11.0 or doing something fancy, ignore this:
.rvmrc

# if using bower-rails ignore default bower_components path bower.json files
/vendor/assets/bower_components
*.bowerrc
bower.json

# Ignore pow environment settings
.powenv

# Ignore Byebug command history file.
.byebug_history

# Ignore node_modules
node_modules/

# Ignore precompiled javascript packs
/public/packs
/public/packs-test
/public/assets

# Ignore yarn files
/yarn-error.log
yarn-debug.log*
.yarn-integrity

# Ignore uploaded files in development
/storage/*
!/storage/.keep

### Ruby ###
*.gem
/.config
/InstalledFiles
/pkg/
/spec/reports/
/spec/examples.txt
/test/tmp/
/test/version_tmp/
/tmp/

# Used by dotenv library to load environment variables.
# .env

# Ignore Byebug command history file.

## Specific to RubyMotion:
.dat*
.repl_history
build/
*.bridgesupport
build-iPhoneOS/
build-iPhoneSimulator/

## Specific to RubyMotion (use of CocoaPods):
#
# We recommend against adding the Pods directory to your .gitignore. However
# you should judge for yourself, the pros and cons are mentioned at:
# https://guides.cocoapods.org/using/using-cocoapods.html#should-i-check-the-pods-directory-into-source-control
# vendor/Pods/

## Documentation cache and generated files:
/.yardoc/
/_yardoc/
/doc/
/rdoc/

/.bundle/
/lib/bundler/man/

# for a library or gem, you might want to ignore these files since the code is
# intended to run in multiple environments; otherwise, check them in:
# Gemfile.lock
# .ruby-version
# .ruby-gemset

# unless supporting rvm < 1.11.0 or doing something fancy, ignore this:

# End of https://www.gitignore.io/api/git,ruby,rails,jetbrains+all
```

You should now have the following three files in your `docker-ruby-on-rails`
directory.


- .dockerignore
- compose.yaml
- Dockerfile


To learn more about the files, see the following:
 - [Dockerfile](../../reference/dockerfile.md)
 - [.dockerignore](../../reference/dockerfile.md#dockerignore-file)
 - [compose.yaml](/reference/compose-file/_index.md)

## Run the application

Inside the `docker-ruby-on-rails` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000). You should see a simple Ruby on Rails application.

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
