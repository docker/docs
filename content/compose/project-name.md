---
title: Specify a project name
description: Understand the different ways you can set a project name in Compose and what the precedence is.
keywords: name, compose, project, -p flag, name top-level element
---

## Use `-p` to specify a project name

Each configuration has a project name which Compose can set in different ways. The level of precedence (from highest to lowest) for each method is as follows: 

1. The `-p` command line flag 
2. The [COMPOSE_PROJECT_NAME environment variable][]
3. The top level `name:` variable from the config file (or the last `name:` from
  a series of config files specified using `-f`)
4. The `basename` of the project directory containing the config file (or
  containing the first config file specified using `-f`)
5. The `basename` of the current directory if no config file is specified



Project names must contain only lowercase letters, decimal digits, dashes, and
underscores, and must begin with a lowercase letter or decimal digit. If the
`basename` of the project directory or current directory violates this
constraint, you must use one of the other mechanisms.

### Have multiple isolated environments on a single host

Compose uses a project name to isolate environments from each other. You can make use of this project name in several different contexts:

* On a dev host, to create multiple copies of a single environment, such as when you want to run a stable copy for each feature branch of a project
* On a CI server, to keep builds from interfering with each other, you can set
  the project name to a unique build number
* On a shared host or dev host, to prevent different projects, which may use the
  same service names, from interfering with each other

The default project name is the base name of the project directory. You can set
a custom project name by using the
[`-p` command line option](reference/index.md) or the
[`COMPOSE_PROJECT_NAME` environment variable](environment-variables/envvars.md#compose_project_name).

The default project directory is the base directory of the Compose file. A custom value
for it can be defined with the `--project-directory` command line option.
