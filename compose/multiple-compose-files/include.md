---
description: How to use Docker Compose's include top-level element
keywords: compose, docker, include, compose file
title: Include
---

With Docker Compose version 2.20 and later, you can include a whole Compose file directly in your local Compose file using the [`include` top-level element](../compose-file/14-include.md). This solves the relative path problem that `extends` and merging present. 

With `include` support, Compose makes it easier to modularize complex applications into sub-Compose files. This allows application configurations to be made simpler and more explicit. This also helps to reflect in the config file organization the engineering team responsible for the code.

Each path listed in the `include` section is loaded as an individual Compose application model, with it’s own project directory, in order to resolve relative paths.

Once the included Compose application is loaded, all resources definitions are copied into the current Compose application model.

> **Note**
>
> `include` applies recursively so an included Compose file which declares its own `include` section, triggers those other files to be included as well.

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

In the example above, `my-compose-include.yaml` manages `serviceB` which details some replicas, web UI to inspect data, isolated networks, volumes for data persistence, etc. The application relying on `serviceB` doesn’t need to know about the infrastructure details, and consumes the Compose file as a building block it can rely on. 

This means the team managing `serviceB` can refactor its own database component to introduce additional services without any dependent teams being impacted. It also means that the dependent teams don't need to include additional flags on each Compose command they run.

## Reference information

[`include` top-level element](../compose-file/14-include.md)
