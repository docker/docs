---
title: Compose Bridge templates
description: Learn about the Compose Bridge templates syntax
keywords: compose, bridge, templates
---

Compose Bridge's default transformation uses templates to produce Kubernetes manifests.
This page describes the templating mechanism.

## Syntax

Templates are plain text files, using [go-template](https://pkg.go.dev/text/template)
to allow logic and data injection based on the `compose.yaml` model.

When a template is executed, it must produce a YAML file. Multiple files can be generated
as long as those are separated by `---`

The first line, when creating the YAML file, defines the file being generated using a custom notation:
```yaml
#! manifest.yaml
```
With this header comment, `manifest.yaml` will be created by Compose Bridge with yalml document
content.

Combining these together, you can write a template to iterate over some of Compose resources,
then for each resource you can produce a dedicated manifest:

```yaml
{{ range $name, $service := .services }}
---
#! {{ $name }}-manifest.yaml
# Generated code, do not edit
key: value
## ...
{{ end }}
```

This example produces a manifest file for each and every Compose service in you Compose model.


## Input

The input compose model is the canonical yaml model you can get by running
 `docker compose config`. Within a template you can access model nodes using 
 dot notation:

 ```yaml
# iterate over a yaml sequence
{{ range $name, $service := .services }}
  # access a nested attribute using dot notation
  {{ if eq $service.deploy.mode "global" }}
kind: DaemonSet
  {{ end }}
{{ end }}
```

You can check the [Compose Specification json-spec file](https://github.com/compose-spec/compose-go/blob/main/schema/compose-spec.json) to have a full overview of the Compose model.

## Helpers

As part of the Go templating syntax, Compose Bridge offers a set of helper functions:

- `seconds`: convert a [duration](https://github.com/compose-spec/compose-spec/blob/master/11-extension.md#specifying-durations) into an integer
- `uppercase` convert a string into upper case characters
- `title`: convert a string by capitalizing first letter of each word
- `safe`: convert a string into a safe identifier, replacing all characters but \[a-z\] with `-`
- `truncate`: removes the N first elements from a list
- `join`: group elements from a list into a single string, using a separator
- `base64`: encode string as base64
- `map`: transform value according to mappings expressed as `"value -> newValue"` strings 
- `indent`: writes string content indented by N spaces
- `helmValue`: write the string content as a template value in final file
