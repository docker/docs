---
title: Interpolation
description: Learn about interpolation
keywords: compose, compose specification, interpolation, compose file reference
---

Values in a Compose file can be set by variables and interpolated at runtime. Compose files use a Bash-like
syntax `${VARIABLE}`.

Both `$VARIABLE` and `${VARIABLE}` syntax is supported. Default values can be defined inline using typical shell syntax:

- `${VARIABLE:-default}` evaluates to `default` if `VARIABLE` is unset or
  empty in the environment.
- `${VARIABLE-default}` evaluates to `default` only if `VARIABLE` is unset
  in the environment.

Similarly, the following syntax allows you to specify mandatory variables:

- `${VARIABLE:?err}` exits with an error message containing `err` if
  `VARIABLE` is unset or empty in the environment.
- `${VARIABLE?err}` exits with an error message containing `err` only if
  `VARIABLE` is unset in the environment.

Interpolation can also be nested:

- `${VARIABLE:-${FOO}}`
- `${VARIABLE?$FOO}`
- `${VARIABLE:-${FOO:-default}}`

Other extended shell-style features, such as `${VARIABLE/foo/bar}`, are not
supported by Compose.

You can use a `$$` (double-dollar sign) when your configuration needs a literal
dollar sign. This also prevents Compose from interpolating a value, so a `$$`
allows you to refer to environment variables that you don't want processed by
Compose.

```yml
web:
  build: .
  command: "$$VAR_NOT_INTERPOLATED_BY_COMPOSE"
```

If Compose can't resolve a substituted variable and no default value is defined, it displays a warning and substitutes the variable with an empty string.

As any values in a Compose file can be interpolated with variable substitution, including compact string notation
for complex elements, interpolation is applied before a merge on a per-file basis.

Interpolation applies only to YAML values, not to keys. For the few places where keys are actually arbitrary
user-defined strings, such as [labels](05-services.md#labels) or [environment](05-services.md#environment), an alternate equal sign syntax
must be used for interpolation to apply. For example:

```yml
services:
  foo:
    labels:
      "$VAR_NOT_INTERPOLATED_BY_COMPOSE": "BAR"
```

```yml
services:
  foo:
    labels:
      - "$VAR_INTERPOLATED_BY_COMPOSE=BAR"
```
