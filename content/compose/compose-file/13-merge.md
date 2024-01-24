---
title: Merge
description: Learn about merging rules
keywords: compose, compose specification, merge, compose file reference
---

Compose lets you define a Compose application model through [multiple Compose files](https://docs.docker.com/compose/multiple-compose-files/). 
When doing so, Compose follows the rules declared in this section to merge Compose files.

## Mapping

A YAML `mapping` gets merged by adding missing entries and merging the conflicting ones.

Merging the following example YAML trees:

```yaml
services:
  foo:
    key1: value1
    key2: value2
```

```yaml
services:
  foo:
    key2: VALUE
    key3: value3
```

Results in a Compose application model equivalent to the YAML tree:

```yaml
services:
  foo:
    key1: value1
    key2: VALUE
    key3: value3
```

## Sequence

A YAML `sequence` is merged by appending values from the overriding Compose file to the previous one.

Merging the following example YAML trees:

```yaml
services:
  foo:
    DNS:
      - 1.1.1.1
```

```yaml
services:
  foo:
    DNS: 
      - 8.8.8.8
```

Results in a Compose application model equivalent to the YAML tree:

```yaml
services:
  foo:
    DNS:
      - 1.1.1.1
      - 8.8.8.8
```

## Exceptions

### Shell commands

When merging Compose files that use the services attributes [command](05-services.md#command), [entrypoint](05-services.md#entrypoint) and [healthcheck: `test`](05-services.md#healthcheck), the value is overridden by the latest Compose file, and not appended.

Merging the following example YAML trees:

```yaml
services:
  foo:
    command: ["echo", "foo"]
```

```yaml
services:
  foo:
    command: ["echo", "bar"]
```

Results in a Compose application model equivalent to the YAML tree:

```yaml
services:
  foo:
    command: ["echo", "bar"]
```

### Unique resources

Applies to the [ports](05-services.md#ports), [volumes](05-services.md#volumes), [secrets](05-services.md#secrets) and [configs](05-services.md#configs) services attributes.
While these types are modeled in a Compose file as a sequence, they have special uniqueness requirements:

| Attribute   | Unique key               |
|-------------|--------------------------|
| volumes     |  target                  |
| secrets     |  source                  |
| configs     |  source                  |
| ports       |  {ip, target, published, protocol}   |

When merging Compose files, Compose appends new entries that do not violate a uniqueness constraint and merge entries that share a unique key.

Merging the following example YAML trees:

```yaml
services:
  foo:
    volumes:
      - foo:/work
```

```yaml
services:
  foo:
    volumes:
      - bar:/work
```

Results in a Compose application model equivalent to the YAML tree:

```yaml
services:
  foo:
    volumes:
      - bar:/work
```

### Reset value

In addition to the previously described mechanism, an override Compose file can also be used to remove elements from your application model.
For this purpose, the custom [YAML tag](https://yaml.org/spec/1.2.2/#24-tags) `!reset` can be set to
override a value set by the overriden Compose file. A valid value for attribute must be provided,
but will be ignored and target attribute will be set with type's default value or `null`. 

For readability, it is recommended to explicitly set the attribute value to the null (`null`) or empty
array `[]` (with `!reset null` or `!reset []`) so that it is clear that resulting attribute will be
cleared.

Merging the following example YAML trees:

```yaml
services:
  foo:
    build: 
      dockerfile: foo.Dockerfile
    read_only: true
    environment:
      FOO: BAR
    ports:
      - "8080:80"            
```

```yaml
services:
  foo:
    image: foo
    build: !reset null
    read_only: !reset false
    environment:
      FOO: !reset null
    ports: !reset []
```

Result in a Compose application model equivalent to the YAML tree:

```yaml
services:
  foo:
    image: foo
    build: null
    read_only: false
    environment: {}
    ports: []
```

## Additional resources

For more information on how merge can be used to create a composite compose file, see [Working with multiple Compose files](../multiple-compose-files/_index.md)