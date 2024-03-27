---
description: Overview of Docker Compose Bridge
keywords: compose, orchestration, kubernetes, bridge
title: Overview of Docker Compose Bridge
---

This page provides usage information for the `compose-bridge` command.

## Introduction

Docker Compose makes it easy to define a multi-container application
to be ran on a single-node Docker engine, relying on compose.yaml to
describe resources with a simple abstraction.

Compose bridge allows to reuse this exact same compose.yaml model but
translate it into another platform definition format, with a primary
focus on Kubernetes. This transformation can be customized to match
specific needs and requirements.

## Usage

Compose Bridge is a command line that consumes a compose.yaml model 
and run transformation to produce resource definitions for another platform.
[By default](transformation.md) it produces Kubernetes manifests
```
$ compose-bridge -f compose.yaml convert
Kubernetes resource api-deployment.yaml created
Kubernetes resource web-deployment.yaml created
Kubernetes resource api-expose.yaml created
Kubernetes resource web-expose.yaml created
Kubernetes resource 0-namespace.yaml created
Kubernetes resource default-network-policy.yaml created
Kubernetes resource api-service.yaml created
Kubernetes resource web-service.yaml created
```

Such manifests can then be used to run the application on Kubernetes using
the standard deployment command `kubectl apply`.

## Customization

The Kubernetes manifests produced by Compose Bridge based on a compose.yaml
model are designed to allow deployment on Docker Desktop with Kubernetes enabled. 

Kubernetes is such a veratile platform that there are many ways
to map compose concepts into a Kubernetes resource definitions. Compose
Bridge let you customize the transformation to match your own infrastructure
decisions and preferences, with various level of flexibility / investment.


### Use Kustomize

TODO

### Tweak the default templates

You can extract templates used by default transformation `docker/compose-bridge` by 
running `compose-bridge template duplicate` and adjust those to match your needs.

Once extracted, you can edit those templates to match your needs or add missing
declaration, then package those into a new Docker image you can use with
Compose Bridge as a custom transformation:

```Dockerfile
FROM docker/compose-bridge
COPY templates-customized /templates
```

### Add your own templates

For resources that are not managed by Compose Bridge default transformation, 
you can build your own templates. The compose.yaml model maybe does not offer
the configuration attributes required to populate the target manifest, you can
then rely on Compose custom extensions to let developers better describe the
application, and offer an agnostic transformation.

For illustration, let's consider developers can add `x-virtual-host` metadata
to service definitions in compose.yaml. You can use this custom attribute
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


### Build your own transformation

While Compose Bridge template make it easy to customize with minimal changes,
you maybe want to make significant changes, or rely on an existing conversion tool.

A Compose Bridge transformation is a docker image that is designed to get compose model
from `/in/compose.yaml` and produce platform manifests under `/out`. This simple 
contract make it easy to bundle an alternate transformation, as illustrated here using 
[kompose](https://kompose.io/):

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

This Dockerfile bundles kompose and defines command to run this tool according
to Compose Bridge transformation contract.

## Useful resources

TBD