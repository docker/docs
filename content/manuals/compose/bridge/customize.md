---
title: Customize Compose Bridge 
linkTitle: Customize
weight: 20
description: Learn about the Compose Bridge templates syntax
keywords: compose, bridge, templates
---

{{< summary-bar feature_name="Compose bridge" >}}

This page explains how Compose Bridge utilizes templating to efficiently translate Docker Compose files into Kubernetes manifests. It also explain how you can customize these templates for your specific requirements and needs, or how you can build your own transformation. 

## How it works 

Compose bridge uses transformations to let you convert a Compose model into another form. 

A transformation is packaged as a Docker image that receives the fully-resolved Compose model as `/in/compose.yaml` and can produce any target format file under `/out`.

Compose Bridge provides its transformation for Kubernetes using Go templates, so that it is easy to extend for customization by just replacing or appending your own templates.

### Syntax

Compose Bridge make use of templates to transform a Compose configuration file into Kubernetes manifests. Templates are plain text files that use the [Go templating syntax](https://pkg.go.dev/text/template). This enables the insertion of logic and data, making the templates dynamic and adaptable according to the Compose model.

When a template is executed, it must produce a YAML file which is the standard format for Kubernetes manifests. Multiple files can be generated as long as they are separated by `---`

Each YAML output file begins with custom header notation, for example:

```yaml
#! manifest.yaml
```

In the following example, a template iterates over services defined in a `compose.yaml` file. For each service, a dedicated Kubernetes manifest file is generated, named according to the service and containing specified configurations.

```yaml
{{ range $name, $service := .services }}
---
#! {{ $name }}-manifest.yaml
# Generated code, do not edit
key: value
## ...
{{ end }}
```

### Input

The input Compose model is the canonical YAML model you can get by running  `docker compose config`. Within the templates, data from the `compose.yaml` is accessed using dot notation, allowing you to navigate through nested data structures. For example, to access the deployment mode of a service, you would use `service.deploy.mode`:

 ```yaml
# iterate over a yaml sequence
{{ range $name, $service := .services }}
  # access a nested attribute using dot notation
  {{ if eq $service.deploy.mode "global" }}
kind: DaemonSet
  {{ end }}
{{ end }}
```

You can check the [Compose Specification JSON schema](https://github.com/compose-spec/compose-go/blob/main/schema/compose-spec.json) to have a full overview of the Compose model. This schema outlines all possible configurations and their data types in the Compose model. 

### Helpers

As part of the Go templating syntax, Compose Bridge offers a set of YAML helper functions designed to manipulate data within the templates efficiently:

- `seconds`: Converts a [duration](/reference/compose-file/extension.md#specifying-durations) into an integer
- `uppercase`: Converts a string into upper case characters
- `title`: Converts a string by capitalizing the first letter of each word
- `safe`: Converts a string into a safe identifier, replacing all characters (except lowercase a-z) with `-`
- `truncate`: Removes the N first elements from a list
- `join`: Groups elements from a list into a single string, using a separator
- `base64`: Encodes a string as base64 used in Kubernetes for encoding secrets
- `map`: Transforms a value according to mappings expressed as `"value -> newValue"` strings 
- `indent`: Writes string content indented by N spaces
- `helmValue`: Writes the string content as a template value in the final file

In the following example, the template checks if a healthcheck interval is specified for a service, applies the `seconds` function to convert this interval into seconds and assigns the value to the `periodSeconds` attribute.

```yaml
{{ if $service.healthcheck.interval }}
            periodSeconds: {{ $service.healthcheck.interval | seconds }}{{ end }}
{{ end }}
```

## Customization

As Kubernetes is a versatile platform, there are many ways
to map Compose concepts into Kubernetes resource definitions. Compose
Bridge lets you customize the transformation to match your own infrastructure
decisions and preferences, with various level of flexibility and effort.

### Modify the default templates

You can extract templates used by the default transformation `docker/compose-bridge-kubernetes`,
by running `compose-bridge transformations create --from docker/compose-bridge-kubernetes my-template` 
and adjusting the templates to match your needs.

The templates are extracted into a directory named after your template name, in this case `my-template`.  
It includes a Dockerfile that lets you create your own image to distribute your template, as well as a directory containing the templating files.  
You are free to edit the existing files, delete them, or [add new ones](#add-your-own-templates) to subsequently generate Kubernetes manifests that meet your needs.  
You can then use the generated Dockerfile to package your changes into a new transformation image, which you can then use with Compose Bridge:

```console
$ docker build --tag mycompany/transform --push .
```

You can then use your transformation as a replacement:

```console
$ compose-bridge convert --transformations mycompany/transform 
```

### Add your own templates

For resources that are not managed by Compose Bridge's default transformation, 
you can build your own templates. The `compose.yaml` model may not offer all 
the configuration attributes required to populate the target manifest. If this is the case, you can
then rely on Compose custom extensions to better describe the
application, and offer an agnostic transformation.

For example, if you add `x-virtual-host` metadata
to service definitions in the `compose.yaml` file, you can use the following custom attribute
to produce Ingress rules:

```yaml
{{ $project := .name }}
#! {{ $name }}-ingress.yaml
# Generated code, do not edit
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: virtual-host-ingress
  namespace: {{ $project }}
spec:
  rules:  
{{ range $name, $service := .services }}
{{ if $service.x-virtual-host }}
  - host: ${{ $service.x-virtual-host }}
    http:
      paths:
      - path: "/"
        backend:
          service:
            name: ${{ name }}
            port:
              number: 80  
{{ end }}
{{ end }}
```

Once packaged into a Docker image, you can use this custom template
when transforming Compose models into Kubernetes in addition to other
transformations:

```console
$ compose-bridge convert \
    --transformation docker/compose-bridge-kubernetes \
    --transformation mycompany/transform 
```

### Build your own transformation

While Compose Bridge templates make it easy to customize with minimal changes,
you may want to make significant changes, or rely on an existing conversion tool.

A Compose Bridge transformation is a Docker image that is designed to get a Compose model
from `/in/compose.yaml` and produce platform manifests under `/out`. This simple 
contract makes it easy to bundle an alternate transformation using 
[Kompose](https://kompose.io/):

```Dockerfile
FROM alpine

# Get kompose from github release page
RUN apk add --no-cache curl
ARG VERSION=1.32.0
RUN ARCH=$(uname -m | sed 's/armv7l/arm/g' | sed 's/aarch64/arm64/g' | sed 's/x86_64/amd64/g') && \
    curl -fsL \
    "https://github.com/kubernetes/kompose/releases/download/v${VERSION}/kompose-linux-${ARCH}" \
    -o /usr/bin/kompose
RUN chmod +x /usr/bin/kompose

CMD ["/usr/bin/kompose", "convert", "-f", "/in/compose.yaml", "--out", "/out"]
```

This Dockerfile bundles Kompose and defines the command to run this tool according
to the Compose Bridge transformation contract.

## What's next?

- [Explore the advanced integration](advanced-integration.md)
