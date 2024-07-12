---
description: Overview of Docker Compose Bridge
keywords: compose, orchestration, kubernetes, bridge
title: Overview of Docker Compose Bridge
---

{{< include "compose-bridge-early-access.md" >}}

## Introduction

Docker Compose makes it easy to define a multi-container application
to be run on a single-node Docker Engine, relying on a `compose.yaml` file to
describe resources with a simple abstraction.

Compose Bridge lets you reuse this exact same `compose.yaml` file but
translate it into another platform's definition format, with a primary
focus on Kubernetes. This transformation can be customized to match
specific needs and requirements.

## Usage

Compose Bridge is a command line tool that consumes a `compose.yaml` file  
and runs a transformation to produce resource definitions for another platform.
[By default](transformation.md), it produces Kubernetes manifests and a Kustomize overlay for Docker Desktop. For example: 
```console
$ compose-bridge -f compose.yaml convert
Kubernetes resource api-deployment.yaml created
Kubernetes resource db-deployment.yaml created
Kubernetes resource web-deployment.yaml created
Kubernetes resource api-expose.yaml created
Kubernetes resource db-expose.yaml created
Kubernetes resource web-expose.yaml created
Kubernetes resource 0-avatars-namespace.yaml created
Kubernetes resource default-network-policy.yaml created
Kubernetes resource private-network-policy.yaml created
Kubernetes resource public-network-policy.yaml created
Kubernetes resource db-db_data-persistentVolumeClaim.yaml created
Kubernetes resource api-service.yaml created
Kubernetes resource web-service.yaml created
Kubernetes resource kustomization.yaml created
Kubernetes resource db-db_data-persistentVolumeClaim.yaml created
Kubernetes resource api-service.yaml created
Kubernetes resource web-service.yaml created
Kubernetes resource kustomization.yaml created
```

Such manifests can then be used to run the application on Kubernetes using
the standard deployment command `kubectl apply -k out/overlays/desktop/`.

## Customization

The Kubernetes manifests produced by Compose Bridge 
are designed to allow deployment on Docker Desktop with Kubernetes enabled. 

Kubernetes is such a versatile platform that there are many ways
to map Compose concepts into a Kubernetes resource definitions. Compose
Bridge lets you customize the transformation to match your own infrastructure
decisions and preferences, with various level of flexibility / investment.


### Modify the default templates

You can extract templates used by default transformation `docker/compose-bridge-kubernetes`
by running `compose-bridge transformations create my-template --from docker/compose-bridge-kubernetes` 
and adjusting those to match your needs.

The templates will be extracted into a directory named after your template name (ie `my-template`).  
Inside, you will find a Dockerfile that allows you to create your own image to distribute your template, as well as a directory containing the templating files.  
You are free to edit the existing files, delete them, or [add new ones](#add-your-own-templates) to subsequently generate Kubernetes manifests that meet your needs.  
You can then use the generated Dockerfile to package your changes into a new Transformer image, which you can then use with Compose Bridge:

```console
$ docker build --tag mycompany/transform --push .
```

You can then use your transformation as a replacement:
```console
$ compose-bridge -f compose.yaml convert --transformation mycompany/transform 
```

For more information, see [Templates](./templates.md).

### Add your own templates

For resources that are not managed by Compose Bridge's default transformation, 
you can build your own templates. The `compose.yaml` model may not offer all 
the configuration attributes required to populate the target manifest. If this is the case, you can
then rely on Compose custom extensions to let developers better describe the
application, and offer an agnostic transformation.

As an illustration, if developers add `x-virtual-host` metadata
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
$ compose-bridge -f compose.yaml convert \
    --transformation docker/compose-bridge-kubernetes \
    --transformation mycompany/transform 
```

### Build your own transformation

While Compose Bridge templates make it easy to customize with minimal changes,
you may want to make significant changes, or rely on an existing conversion tool.

A Compose Bridge transformation is a Docker image that is designed to get a Compose model
from `/in/compose.yaml` and produce platform manifests under `/out`. This simple 
contract makes it easy to bundle an alternate transformation, as illustrated below using 
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

## Use `compose-bridge` as a `kubectl` plugin

To use the `compose-bridge` binary as a `kubectl` plugin, you need to make sure that the binary is available in your PATH and the name of the binary is prefixed with `kubectl-`. 

1. Rename or copy the `compose-bridge` binary to `kubectl-compose_bridge`:

    ```console
    $ mv /path/to/compose-bridge /usr/local/bin/kubectl-compose_bridge
    ```

2. Ensure that the binary is executable:
    
    ```console
    $ chmod +x /usr/local/bin/kubectl-compose_bridge
    ```

3. Verify that the plugin is recognized by `kubectl`:

    ```console
    $ kubectl plugin list
    ```

    In the output, you should see `kubectl-compose_bridge`.

4. Now you can use `compose-bridge` as a `kubectl` plugin:

    ```console
   $ kubectl compose-bridge [command]
    ```

Replace `[command]` with any `compose-bridge` command you want to use.
