---
description: How to use Docker Compose's include top-level element
keywords: compose, docker, include, compose file
title: Include
---

> **Note**
>
> `include` is available in Docker Compose version 2.20 and later, and Docker Desktop version 4.22 and later. 

With the [`include` top-level element](../compose-file/14-include.md), you can include a separate Compose file directly in your local Compose file. This solves the relative path problem that [`extends`](extends.md) and [merge](merge.md) present. 

`include` makes it easier to modularize complex applications into sub-Compose files. This allows application configurations to be made simpler and more explicit. This also helps to reflect in the config file organization the engineering team responsible for the code.

Each path listed in the `include` section loads as an individual Compose application model, with its own project directory, in order to resolve relative paths.

Once the included Compose application loads, all resources are copied into the current Compose application model.

> **Note**
>
> `include` applies recursively so an included Compose file which declares its own `include` section, results in those other files being included as well.

## Example

```yaml
include:
  - my-compose-include.yaml  #with serviceB declared
services:
  serviceA:
    build: .
    depends_on:
      - serviceB #use serviceB directly as if it was declared in this Compose file
```

`my-compose-include.yaml` manages `serviceB` which details some replicas, web UI to inspect data, isolated networks, volumes for data persistence, etc. The application relying on `serviceB` doesn’t need to know about the infrastructure details, and consumes the Compose file as a building block it can rely on. 

This means the team managing `serviceB` can refactor its own database component to introduce additional services without impacting any dependent teams. It also means that the dependent teams don't need to include additional flags on each Compose command they run.

## Reference information

[`include` top-level element](../compose-file/14-include.md)
