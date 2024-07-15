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

## Get the sample application

You'll need to clone a new repository to get a sample application that includes logic to connect to the database.

1. Change to a directory where you want to clone the repository and run the following command.

   ```console
   $ git clone https://github.com/falconcr/docker-ruby-on-rails.git
   ```

2. In the cloned repository's directory, manually create the following files in your project directory.

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



Create a file named `compose.yaml` with the following contents.

```yaml {collapse=true,title=compose.yaml}
# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Docker Compose reference guide at
# https://docs.docker.com/go/compose-spec-reference/
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

## Add a local database and persist data

You can use containers to set up local services, like a database. In this section, you'll update the `compose.yaml` file to define a database service and a volume to persist data.

In the cloned repository's directory, open the `compose.yaml` file in an IDE or text editor. Yyou need to add the database password file as an environment variable to the server service and specify the secret file to use .

The following is the updated `compose.yaml` file.

```yaml {hl_lines="10-30"}
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
├── docker-ruby-on-rails/
├── app
├── bin
├── config
│── db/
│ │ └─ password.txt
├── lib
├── log
├── public
├── storage
├── test
├── tmp
├── vendor
│ ├── .dockerignore
│ ├── .gitignore
│ ├── confi.ru
│ ├── Gemfile
│ ├── Gemfile.lock
│ ├── compose.yaml
│ ├── Rakefile
│ ├── Dockerfile
│ └── README.md
```

Now, run the following `docker compose up` command to start your application.

```console
$ docker compose up --build
```

Refresh http://localhost:3000 in your browser and verify that the Whale items persisted, even after the containers were removed and ran again.

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
