---
title: Use NGINX with Interlock
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

The following configuration options are available for the Nginx extension:

| Option                  | Type   | Description                                                                                         |
|:------------------------|:-------|:----------------------------------------------------------------------------------------------------|
| `User`                  | string | user to be used in the proxy                                                                        |
| `PidPath`               | string | path to the pid file for the proxy service                                                          |
| `MaxConnections`        | int    | maximum number of connections for proxy service                                                     |
| `ConnectTimeout`        | int    | timeout in seconds for clients to connect                                                           |
| `SendTimeout`           | int    | timeout in seconds for the service to send a request to the proxied upstream                        |
| `ReadTimeout`           | int    | timeout in seconds for the service to read a response from the proxied upstream                     |
| `IPHash`                | bool   | specifies that requests are distributed between servers based on client IP addresses                |
| `SSLOpts`               | string | options to be passed when configuring SSL                                                           |
| `SSLDefaultDHParam`     | int    | size of DH parameters                                                                               |
| `SSLDefaultDHParamPath` | string | path to DH parameters file                                                                          |
| `SSLVerify`             | string | SSL client verification                                                                             |
| `WorkerProcesses`       | string | number of worker processes for the proxy service                                                    |
| `RLimitNoFile`          | int    | number of maxiumum open files for the proxy service                                                 |
| `SSLCiphers`            | string | SSL ciphers to use for the proxy service                                                            |
| `SSLProtocols`          | string | enable the specified TLS protocols                                                                  |
| `AccessLogPath`         | string | Path to use for access logs (default: `/dev/stdout`)                                                |
| `ErrorLogPath`          | string | Path to use for error logs (default: `/dev/stdout`)                                                 |
| `MainLogFormat`         | string | [Format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format) to use for main logger  |
| `TraceLogFormat`        | string | [Format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format) to use for trace logger |

{% endif %}
