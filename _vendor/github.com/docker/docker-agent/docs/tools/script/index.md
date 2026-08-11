---
title: "Script Tool"
description: "Define custom shell scripts as named tools with typed parameters."
keywords: docker agent, ai agents, tools, toolsets, script tool
linkTitle: "Script"
weight: 30
canonical: https://docs.docker.com/ai/docker-agent/tools/script/
---

_Define custom shell scripts as named tools with typed parameters._

## Overview

The script tool lets you define custom shell scripts as named tools. Unlike the generic [shell tool](../shell/index.md) where the agent writes the command, script tools execute predefined commands — ideal for exposing safe, well-scoped operations with descriptive names.

## Configuration

### Simple Scripts

```yaml
toolsets:
  - type: script
    shell:
      run_tests:
        cmd: task test
        description: Run the project test suite
      lint:
        cmd: task lint
        description: Run the linter
```

### Scripts with Parameters

Use `${param}` interpolation and JSON Schema to define typed arguments:

```yaml
toolsets:
  - type: script
    shell:
      deploy:
        cmd: ./scripts/deploy.sh ${env}
        description: Deploy to an environment
        args:
          env:
            type: string
            enum: [staging, production]
        required: [env]
```

## Properties

| Property                          | Type   | Description                                                |
| --------------------------------- | ------ | ---------------------------------------------------------- |
| `shell.<name>.cmd`                | string | Shell command to execute (supports `${arg}` interpolation) |
| `shell.<name>.description`        | string | Description shown to the model                             |
| `shell.<name>.args`               | object | Parameter definitions (JSON Schema properties)             |
| `shell.<name>.required`           | array  | Required parameter names                                   |
| `shell.<name>.env`                | object | Environment variables for this script                      |
| `shell.<name>.working_dir`        | string | Working directory for script execution                     |

> [!TIP]
> **Script vs. Shell**
>
> Use the [shell tool](../shell/index.md) when the agent needs to run arbitrary commands. Use the script tool when you want to expose specific, predefined operations with clear names and typed parameters — giving the agent less freedom but more safety.
