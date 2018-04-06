---
title: Layer 7 routing configuration reference
description: Learn the configuration options for the UCP layer 7 routing solution
keywords: routing, proxy
---

Once you enable the layer 7 routing service, UCP creates the
`com.docker.ucp.interlock.conf-1` configuration and uses it to configure all
the internal components of this service.

The configuration is managed as a TOML file.

## Example configuration

Here's an example of the default configuration used by UCP:

```toml
ListenAddr = ":8080"
DockerURL = "unix:///var/run/docker.sock"
AllowInsecure = false
PollInterval = "3s"

[Extensions]
  [Extensions.default]
    Image = "dockereng/ucp-interlock-extension:3.0.0-2e21c4d"
    ServiceName = "ucp-interlock-extension"
    Args = [""]
    Constraints = ["node.labels.com.docker.ucp.orchestrator.swarm==true"]
    ProxyImage = "dockereng/ucp-interlock-proxy:3.0.0-2e21c4d"
    ProxyServiceName = "ucp-interlock-proxy"
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    ProxyReplicas = 2
    ProxyStopSignal = "SIGQUIT"
    ProxyStopGracePeriod = "5s"
    ProxyConstraints = ["node.labels.com.docker.ucp.orchestrator.swarm==true"]
    PublishMode = "ingress"
    PublishedPort = 80
    TargetPort = 80
    PublishedSSLPort = 8443
    TargetSSLPort = 443
    [Extensions.default.Labels]
      "com.docker.ucp.InstanceID" = "xybt14ctipeasd7lfggl4hw94"
    [Extensions.default.ContainerLabels]
      "com.docker.ucp.InstanceID" = "xybt14ctipeasd7lfggl4hw94"
    [Extensions.default.ProxyLabels]
      "com.docker.ucp.InstanceID" = "xybt14ctipeasd7lfggl4hw94"
    [Extensions.default.ProxyContainerLabels]
      "com.docker.ucp.InstanceID" = "xybt14ctipeasd7lfggl4hw94"
    [Extensions.default.Config]
      Version = ""
      User = "nginx"
      PidPath = "/var/run/proxy.pid"
      MaxConnections = 1024
      ConnectTimeout = 600
      SendTimeout = 600
      ReadTimeout = 600
      IPHash = false
      AdminUser = ""
      AdminPass = ""
      SSLOpts = ""
      SSLDefaultDHParam = 1024
      SSLDefaultDHParamPath = ""
      SSLVerify = "required"
      WorkerProcesses = 1
      RLimitNoFile = 65535
      SSLCiphers = "HIGH:!aNULL:!MD5"
      SSLProtocols = "TLSv1.2"
      AccessLogPath = "/dev/stdout"
      ErrorLogPath = "/dev/stdout"
      MainLogFormat = "'$remote_addr - $remote_user [$time_local] \"$request\" '\n\t\t    '$status $body_bytes_sent \"$http_referer\" '\n\t\t    '\"$http_user_agent\" \"$http_x_forwarded_for\"';"
      TraceLogFormat = "'$remote_addr - $remote_user [$time_local] \"$request\" $status '\n\t\t    '$body_bytes_sent \"$http_referer\" \"$http_user_agent\" '\n\t\t    '\"$http_x_forwarded_for\" $request_id $msec $request_time '\n\t\t    '$upstream_connect_time $upstream_header_time $upstream_response_time';"
      KeepaliveTimeout = "75s"
      ClientMaxBodySize = "32m"
      ClientBodyBufferSize = "8k"
      ClientHeaderBufferSize = "1k"
      LargeClientHeaderBuffers = "4 8k"
      ClientBodyTimeout = "60s"
      UnderscoresInHeaders = false
```

## Core configurations

These are the configurations used for the `ucp-interlock` service. The following
options are available:

| Option             | Type        | Description                                                                                    |
|:-------------------|:------------|:-----------------------------------------------------------------------------------------------|
| `ListenAddr`       | string      | Address to serve the Interlock GRPC API. Defaults to `8080`.                                   |
| `DockerURL`        | string      | Path to the socket or TCP address to the Docker API. Defaults to `unix:///var/run/docker.sock` |
| `TLSCACert`        | string      | Path to the CA certificate for connecting securely to the Docker API.                          |
| `TLSCert`          | string      | Path to the certificate for connecting securely to the Docker API.                             |
| `TLSKey`           | string      | Path to the key for connecting securely to the Docker API.                                     |
| `AllowInsecure`    | bool        | Skip TLS verification when connecting to the Docker API via TLS.                               |
| `PollInterval`     | string      | Interval to poll the Docker API for changes. Defaults to `3s`.                                 |
| `EndpointOverride` | string      | Override the default GRPC API endpoint for extensions. The default is  detected via Swarm.     |
| `Extensions`       | []Extension | Array of extensions as listed below.                                                           |

## Extension configuration

Interlock must contain at least one extension to service traffic.
The following options are available to configure the extensions:

| Option             | Type              | Description                                                                   |
|:-------------------|:------------------|:------------------------------------------------------------------------------|
| `Image`            | string            | Name of the Docker image to use for the extension service.                    |
| `Args`             | []string          | Arguments to be passed to the Docker extension service upon creation.         |
| `Labels`           | map[string]string | Labels to add to the extension service.                                       |
| `ServiceName`      | string            | Name of the extension service.                                                |
| `ProxyImage`       | string            | Name of the Docker image to use for the proxy service.                        |
| `ProxyArgs`        | []string          | Arguments to be passed to the proxy service upon creation.                    |
| `ProxyLabels`      | map[string]string | Labels to add to the proxy service.                                           |
| `ProxyServiceName` | string            | Name of the proxy service.                                                    |
| `ProxyConfigPath`  | string            | Path in the service for the generated proxy configuration.                    |
| `ServiceCluster`   | string            | Name of the cluster this extension services.                                  |
| `PublishMode`      | string            | Publish mode for the proxy service. Supported values are `ingress` or `host`. |
| `PublishedPort`    | int               | Port where the proxy service serves non-TLS traffic.                          |
| `PublishedSSLPort` | int               | Port where the proxy service serves TLS traffic.                              |
| `Template`         | string            | Docker configuration object that is used as the extension template.           |
| `Config`           | Config            | Proxy configuration used by the extensions as listed below.                   |

## Proxy configuration

By default NGINX is used as a proxy, so the following NGINX options are
available for the proxy service:

| Option                  | Type   | Description                                                                                          |
|:------------------------|:-------|:-----------------------------------------------------------------------------------------------------|
| `User`                  | string | User to be used in the proxy.                                                                        |
| `PidPath`               | string | Path to the pid file for the proxy service.                                                          |
| `MaxConnections`        | int    | Maximum number of connections for proxy service.                                                     |
| `ConnectTimeout`        | int    | Timeout in seconds for clients to connect.                                                           |
| `SendTimeout`           | int    | Timeout in seconds for the service to send a request to the proxied upstream.                        |
| `ReadTimeout`           | int    | Timeout in seconds for the service to read a response from the proxied upstream.                     |
| `IPHash`                | bool   | Specifies that requests are distributed between servers based on client IP addresses.                |
| `SSLOpts`               | string | Options to be passed when configuring SSL.                                                           |
| `SSLDefaultDHParam`     | int    | Size of DH parameters.                                                                               |
| `SSLDefaultDHParamPath` | string | Path to DH parameters file.                                                                          |
| `SSLVerify`             | string | SSL client verification.                                                                             |
| `WorkerProcesses`       | string | Number of worker processes for the proxy service.                                                    |
| `RLimitNoFile`          | int    | Number of maxiumum open files for the proxy service.                                                 |
| `SSLCiphers`            | string | SSL ciphers to use for the proxy service.                                                            |
| `SSLProtocols`          | string | Enable the specified TLS protocols.                                                                  |
| `AccessLogPath`         | string | Path to use for access logs (default: `/dev/stdout`).                                                |
| `ErrorLogPath`          | string | Path to use for error logs (default: `/dev/stdout`).                                                 |
| `MainLogFormat`         | string | [Format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format) to use for main logger.  |
| `TraceLogFormat`        | string | [Format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format) to use for trace logger. |


