---
description: Understand the accepted syntax for declaring environment variables with an environment file
keywords: fig, composition, compose, docker, orchestration, environment, env file
title: Syntax for environment files in Docker Compose
aliases:
- /compose/env-file/
---

This page provides information on the syntax rules and guidelines when using an `.env` file. It defines the rules for commenting, and explains how values are processed. Additionally, it introduces the concept of interpolation, which allows the use of variables within environment files.

> **Important**
>
> Environment variables from an environment file have lower precedence than those from any other [method](set-environment-variables.md).
For more information, see [Environment variables precedence](envvars-precedence.md).
{ .important }

## Syntax
The following syntax rules apply to environment files:

- Lines beginning with `#` are processed as comments and ignored.
- Blank lines are ignored.
- Unquoted and double-quoted (`"`) values have [interpolation](#interpolation) applied.
- Each line represents a key-value pair. Values can optionally be quoted.
  - `VAR=VAL` -> `VAL`
  - `VAR="VAL"` -> `VAL`
  - `VAR='VAL'` -> `VAL`
- Inline comments for unquoted values must be preceded with a space.
  - `VAR=VAL # comment` -> `VAL`
  - `VAR=VAL# not a comment` -> `VAL# not a comment`
- Inline comments for quoted values must follow the closing quote.
  - `VAR="VAL # not a comment"` -> `VAL # not a comment`
  - `VAR="VAL" # comment` -> `VAL`
- Single-quoted (`'`) values are used literally.
  - `VAR='$OTHER'` -> `$OTHER`
  - `VAR='${OTHER}'` -> `${OTHER}`
- Quotes can be escaped with `\`.
  - `VAR='Let\'s go!'` -> `Let's go!`
  - `VAR="{\"hello\": \"json\"}"` -> `{"hello": "json"}`
- Common shell escape sequences including `\n`, `\r`, `\t`, and `\\` are supported in double-quoted values.
  - `VAR="some\tvalue"` -> `some  value`
  - `VAR='some\tvalue'` -> `some\tvalue`
  - `VAR=some\tvalue` -> `some\tvalue`

### Interpolation
Compose supports interpolation in environment files.
Interpolation is applied for unquoted and double-quoted values.
Both braced (`${VAR}`) and unbraced (`$VAR`) expressions are supported.

For braced expressions, the following formats are supported:
- Direct substitution
  - `${VAR}` -> value of `VAR`
- Default value
  - `${VAR:-default}` -> value of `VAR` if set and non-empty, otherwise `default`
  - `${VAR-default}` -> value of `VAR` if set, otherwise `default`
- Required value
  - `${VAR:?error}` -> value of `VAR` if set and non-empty, otherwise exit with error
  - `${VAR?error}` -> value of `VAR` if set, otherwise exit with error
- Alternative value
  - `${VAR:+replacement}` -> `replacement` if `VAR` is set and non-empty, otherwise empty
  - `${VAR+replacement}` -> `replacement` if `VAR` is set, otherwise empty

For more information, see [Interpolation](../compose-file/12-interpolation.md) in the Compose specification. 
