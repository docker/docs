---
title: Interpolation
description: Learn about interpolation
keywords: compose, compose specification, interpolation, compose file reference
aliases: 
 - /compose/compose-file/12-interpolation/
weight: 90
---

{{% include "compose/interpolation.md" %}}

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

Interpolation can also be nested:

- `${VARIABLE:-${FOO}}`
- `${VARIABLE?$FOO}`
- `${VARIABLE:-${FOO:-default}}`

Other extended shell-style features, such as `${VARIABLE/foo/bar}`, are not
supported by Compose.

Compose processes any string following a `$` sign as long as it makes it
a valid variable definition - either an alphanumeric name (`[_a-zA-Z][_a-zA-Z0-9]*`)
or a braced string starting with `${`. In other circumstances, it will be preserved without attempting to interpolate a value.

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
user-defined strings, such as [labels](services.md#labels) or [environment](services.md#environment), an alternate equal sign syntax
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
