---
description: Declare default environment variables in a file
keywords: fig, composition, compose, docker, orchestration, environment, env file
title: Declare default environment variables in file
---

Compose supports declaring default environment variables in an environment file
named `.env` placed in the project directory. Docker Compose versions earlier than `1.28`,
load the `.env` file from the current working directory, where the command is executed, or from the
project directory if this is explicitly set with the `--project-directory` option. This
inconsistency has been addressed starting with `+v1.28` by limiting the default `.env` file path
to the project directory. You can use the `--env-file` commandline option to override the default
`.env` and specify the path to a custom environment file.

The project directory is specified by the order of precedence:

- `--project-directory` flag
- Folder of the first `--file` flag
- Current directory

## Syntax rules

The "dotEnv" file format is a _de facto_ convention, not a standard with no clear specification, 
and as such syntax might differ between tools using this file format.

The following syntax rules apply to the `.env` file as supported by Docker Compose:

### Compose 1.x

- Compose expects each line in an `env` file to be in `VAR=VAL` format.
- Lines beginning with `#` are processed as comments and ignored.
- Blank lines are ignored.
- There is no special handling of quotation marks. This means that
  **they are part of the VAL**.

### Compose 2.x

Compose V2 adopt the [Ruby-style dotEnv file syntax](https://github.com/bkeepers/dotenv#usage), 
which means:

- Compose expects each line in an `env` file to be in `VAR=VAL` format. Can be preceeded by `export ` which will be ignored
- Lines beginning with `#` are processed as comments and ignored.
- Blank lines are ignored.
- Multi-line values can be wrapped using double quotes. 
- Values can be wrapped by single or double quotes, which will be remived from actual value
- Variables defined by dolar sign `$` inside value are substitued, unless   value is wrapped by single quotes

```
# Some comment, ignored
USER=docker

# $word won't be considered a variable thanks to single quotes
PASSWORD='pas$word'   

# ${USER} will be substitued by value from USER variable
DATABASE_URL="postgres://${USER}@localhost/my_database"

# export is ignored by Docker Compose, but usefull to use the same file for shell scripts
export S3_BUCKET=YOURS3BUCKET

# A multi-line value, wrapped by double quotes
PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----
...
-----END DSA PRIVATE KEY-----"
```

## Compose file and CLI variables

The environment variables you define here are used for
[variable substitution](compose-file/compose-file-v3.md#variable-substitution)
in your Compose file, and can also be used to define the following
[CLI variables](reference/envvars.md):

- `COMPOSE_API_VERSION`
- `COMPOSE_CONVERT_WINDOWS_PATHS`
- `COMPOSE_FILE`
- `COMPOSE_HTTP_TIMEOUT`
- `COMPOSE_PROFILES`
- `COMPOSE_PROJECT_NAME`
- `COMPOSE_TLS_VERSION`
- `DOCKER_CERT_PATH`
- `DOCKER_HOST`
- `DOCKER_TLS_VERIFY`

> **Notes**
>
> * Values present in the environment at runtime always override those defined
>   inside the `.env` file. Similarly, values passed via command-line arguments
>   take precedence as well.
> * Environment variables defined in the `.env` file are not automatically
>   visible inside containers. To set container-applicable environment variables,
>   follow the guidelines in the topic
>   [Environment variables in Compose](environment-variables.md), which
>   describes how to pass shell environment variables through to containers,
>   define environment variables in Compose files, and more.

## More Compose documentation

- [User guide](index.md)
- [Installing Compose](install.md)
- [Getting Started](gettingstarted.md)
- [Command line reference](reference/index.md)
- [Compose file reference](compose-file/index.md)
- [Sample apps with Compose](samples-for-compose.md)
