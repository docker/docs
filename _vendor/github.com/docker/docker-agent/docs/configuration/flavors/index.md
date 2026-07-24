---
title: "Flavors"
description: "Ship one agent file with named variants, enabled at run time as YAML patches."
keywords: docker agent, ai agents, configuration, yaml, flavors, patch, variants, overrides
weight: 95
canonical: https://docs.docker.com/ai/docker-agent/configuration/flavors/
---

_Ship one agent file with named variants, enabled at run time as YAML patches._

## Overview

A flavor is a named YAML patch declared in the agent file itself, under the
top-level `flavors` section. Enabling a flavor applies its patch on top of the
rest of the document before the config is parsed, so one file can carry
several variants — a cheaper model for local runs, extra tools for CI, a more
verbose instruction for debugging — without duplicating the whole config.

```yaml
agents:
  root:
    model: claude
    instruction: You are a helpful assistant.

models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-5

flavors:
  cheap:
    models:
      claude:
        model: claude-3-5-haiku-latest
```

Enable flavors with the repeatable `--flavor` flag:

```bash
$ docker agent run agent.yaml --flavor cheap
```

The flag works on every command that runs an agent — `run`, `chat`, `eval`,
`serve api`, `serve a2a`, `serve mcp` — and order matters: patches apply in
the order the flavors are requested, each on top of the previous result.

```bash
$ docker agent run agent.yaml --flavor cheap --flavor verbose
```

Flavors the file does not define are ignored (with a debug log), so you can
enable the same flavor set across a fleet of agents and each file only reacts
to the names it declares. External sub-agents loaded from OCI or URL
references receive the same enabled flavors.

## Merge Semantics

Patches follow [JSON Merge Patch](https://www.rfc-editor.org/rfc/rfc7386)
semantics, with two extensions for arrays:

| Patch value | Effect |
| --- | --- |
| Object | Merged recursively into the existing object. |
| Scalar or array | Replaces the existing value. |
| `null` | Deletes the key. |
| Key ending in `+` | Appends the items to the existing array. |
| Key ending in `-` | Removes matching entries from an array or object. |

### Merging and replacing

An object patch only touches the keys it names — siblings survive:

```yaml
flavors:
  verbose:
    agents:
      root:
        instruction: Explain your reasoning in detail.  # model, tools, ... unchanged
```

### Deleting a key

Set it to `null`:

```yaml
flavors:
  no-limit:
    models:
      claude:
        max_tokens: null
```

### Appending to an array

Plain arrays replace wholesale. To add entries instead, suffix the key
with `+`:

```yaml
agents:
  root:
    toolsets:
      - type: think

flavors:
  with-shell:
    agents:
      root:
        toolsets+:
          - type: shell
```

With `--flavor with-shell` the root agent gets both `think` and `shell`.

### Removing entries

Suffix the key with `-`. Each item in the patch value selects what to remove:

- From an **array**: a scalar removes equal elements; an object removes every
  element it partially matches (all of the matcher's keys must be present
  with matching values).
- From an **object**: items are key names to drop.

```yaml
flavors:
  slim:
    agents:
      root:
        toolsets-:
          - type: shell   # drop every shell toolset, however configured
        sub_agents-:
          - checker       # drop by value
    models-:
      - spare             # drop the named model definition
```

> [!NOTE]
> The `+` and `-` suffixes are reserved inside flavor patches: a patch cannot
> set a literal key ending in either character. Base documents are unaffected.

## Inspecting the Result

`docker agent debug config` prints the config exactly as the runtime sees it,
flavors applied:

```bash
$ docker agent debug config agent.yaml --flavor cheap --flavor with-shell
```

## HCL

Flavors work in [HCL configs](../hcl/index.md) too, as labeled blocks. The
append/remove operators need quoted attribute names inside object
expressions:

```hcl
flavors "with-shell" {
  agents = {
    root = {
      "toolsets+" = [{ type = "shell" }]
    }
  }
}
```

## Notes

- Flavors require config schema version 13 or later; older versions reject
  the `flavors` key with a hint to bump the top-level `version` field.
- Patches apply before validation, so a flavored config is validated exactly
  like a hand-written one.
- `docker agent push` publishes the raw document, `flavors` section included,
  so consumers of a pushed agent can enable its flavors too.
