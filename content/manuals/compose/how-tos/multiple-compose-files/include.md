---
description: How to use Docker Compose's include top-level element
keywords: compose, docker, include, compose file
title: Include
aliases:
- /compose/multiple-compose-files/include/
---

{{< summary-bar feature_name="Compose include" >}}

{{% include "compose/include.md" %}}

The [`include` top-level element](/reference/compose-file/include.md) helps to reflect the engineering team responsible for the code directly in the config file's organization. It also solves the relative path problem that [`extends`](extends.md) and [merge](merge.md) present. 

Each path listed in the `include` section loads as an individual Compose application model, with its own project directory, in order to resolve relative paths.

Once the included Compose application loads, all resources are copied into the current Compose application model.

> [!NOTE]
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

## Include and overrides

Compose reports an error if any resource from `include` conflicts with resources from the included Compose file. This rule prevents
unexpected conflicts with resources defined by the included compose file author. However, there may be some circumstances where you might want to tweak the
included model. This can be achieved by adding an override file to the include directive:
```yaml
include:
  - path : 
      - third-party/compose.yaml
      - override.yaml  # local override for third-party model
```

The main limitation with this approach is that you need to maintain a dedicated override file per include. For complex projects with multiple
includes this would result into many Compose files.

The other option is to use a `compose.override.yaml` file. While conflicts will be rejected from the file using `include` when same
resource is declared, a global Compose override file can override the resulting merged model, as demonstrated in following example:

Main `compose.yaml` file:
```yaml
include:
  - team-1/compose.yaml # declare service-1
  - team-2/compose.yaml # declare service-2
```

Override `compose.override.yaml` file:
```yaml
services:
  service-1:
    # override included service-1 to enable debugger port
    ports:
      - 2345:2345

  service-2:
    # override included service-2 to use local data folder containing test data
    volumes:
      - ./data:/data
```

Combined together, this allows you to benefit from third-party reusable components, and adjust the Compose model for your needs.

## Reference information

[`include` top-level element](/reference/compose-file/include.md)
