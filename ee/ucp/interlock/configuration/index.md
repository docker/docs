---
title: Configure Interlock
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

Interlock configuration is managed via file as [TOML](https://github.com/toml-lang/toml).
The following will describe how to configure the various components of Interlock.

## Core
The core configuration handles the Interlock service itself.  The following options
are available:

| Option             | Type        | Description                                                                                   |
|:-------------------|:------------|:----------------------------------------------------------------------------------------------|
| `ListenAddr`       | string      | address to serve the Interlock GRPC API (default: `:8080`)                                    |
| `DockerURL`        | string      | path to the socket or TCP address to the Docker API (default: `unix:///var/run/docker.sock`)  |
| `TLSCACert`        | string      | path to the CA certificate for connecting securely to the Docker API                          |
| `TLSCert`          | string      | path to the certificate for connecting securely to the Docker API                             |
| `TLSKey`           | string      | path to the key for connecting securely to the Docker API                                     |
| `AllowInsecure`    | bool        | skip TLS verification when connecting to the Docker API via TLS                               |
| `PollInterval`     | string      | interval to poll the Docker API for changes (default: `3s`)                                   |
| `EndpointOverride` | string      | override the default GRPC API endpoint for extensions (by default this is detected via Swarm) |
| `Extensions`       | []Extension | array of extensions as listed below                                                           |

## Extension
Interlock must contain at least one extension to service traffic.  The following options are
available to configure the extensions.

| Option             | Type                         | Description                                                          |
|:-------------------|:-----------------------------|:---------------------------------------------------------------------|
| `Image`            | string                       | name of the Docker Image to use for the extension service            |
| `Args`             | []string                     | arguments to be passed to the Docker extension service upon creation |
| `Labels`           | map[string]string            | labels to be added to the extension service                          |
| `ServiceName`      | string                       | name of the extension service                                        |
| `ProxyImage`       | string                       | name of the Docker Image to use for the proxy service                |
| `ProxyArgs`        | []string                     | arguments to be passed to the Docker proxy service upon creation     |
| `ProxyLabels`      | map[string]string            | labels to be added to the proxy service                              |
| `ProxyServiceName` | string                       | name of the proxy service                                            |
| `ProxyConfigPath`  | string                       | path in the service for the generated proxy config                   |
| `ServiceCluster`   | string                       | name of the cluster this extension services                          |
| `PublishMode`      | string (`ingress` or `host`) | publish mode that the proxy service uses                             |
| `PublishedPort`    | int                          | port that the proxy service serves non-SSL traffic                   |
| `PublishedSSLPort` | int                          | port that the proxy service serves SSL traffic                       |
| `Template`         | string                       | Docker config object that is used as the extension template          |
| `Config`           | Config                       | proxy configuration used by the extensions as listed below           |

## Proxy
The following options are made available to the extensions.  The extensions use whichever they need to configure
the proxy service.  This provides a way for the user to provide overrides to the extension configuration.

Interlock passes extension configuration through directly to the extension.  Therefore, each extension has
different configuration options available.  See the docs for each extension for the officially supported options.

## Example Configuration
The following is an example configuration to use with the Nginx extension.

```toml
ListenAddr = ":8080"
DockerURL = "unix:///var/run/docker.sock"
PollInterval = "3s"

[Extensions]
  [Extensions.default]
    Image = "docker/interlock-extension-nginx:latest"
    Args = ["-D"]
    ProxyImage = "nginx:alpine"
    ProxyArgs = []
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    ServiceCluster = ""
    PublishMode = "ingress"
    PublishedPort = 80
    TargetPort = 80
    PublishedSSLPort = 443
    TargetSSLPort = 443
    [Extensions.default.Config]
      User = "nginx"
      PidPath = "/var/run/proxy.pid"
      WorkerProcesses = 1
      RlimitNoFile = 65535
      MaxConnections = 2048
    [Extensions.default.Labels]
      extension_name = "defaultExtension"
    [Extensions.default.ProxyLabels]
      proxy_name = "defaultProxy"
```

{% endif %}
