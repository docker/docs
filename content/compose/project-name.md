---
title: Specify a project name
description: Understand the different ways you can set a project name in Compose and what the precedence is.
keywords: name, compose, project, -p flag, name top-level element
---

In Compose, the default project name is derived from the base name of the project directory. However, you have the flexibility to set a custom project name. 

This page offers examples of scenarios where custom project names can be helpful, outlines the various methods to set a project name, and provides the order of precedence for each approach.

> **Note**
>
> The default project directory is the base directory of the Compose file. A custom value can also be set
> for it using the `--project-directory` command line option.

## Example use cases

Compose uses a project name to isolate environments from each other. You can make use of this project name in several different contexts:

- On a development host: Create multiple copies of a single environment, useful for running stable copies for each feature branch of a project.
- On a CI server: Prevent interference between builds by setting the project name to a unique build number.
- On a shared or development host: Avoid interference between different projects that might share the same service names.

## Set a project name

Project names must contain only lowercase letters, decimal digits, dashes, and
underscores, and must begin with a lowercase letter or decimal digit. If the
base name of the project directory or current directory violates this
constraint, alternative mechanisms are available.

The precedence order for each method, from highest to lowest, is as follows:

1. The `-p` command line flag. 
2. The [COMPOSE_PROJECT_NAME environment variable](environment-variables/envvars.md).
3. The [top-level `name:` attribute](compose-file/04-version-and-name.md) in your Compose file. Or the last `name:` if you [specify multiple Compose files](multiple-compose-files/merge.md) in the command line with the `-f` flag.
4. The base name of the project directory containing your Compose file. Or the base name of the first Compose file if you if you [specify multiple Compose files](multiple-compose-files/merge.md) in the command line with the `-f` flag. 
5. The base name of the current directory if no Compose file is specified.

## What's next?

- Read up on [working with multiple Compose files](multiple-compose-files/_index.md).
- Explore some [sample apps](samples-for-compose.md).