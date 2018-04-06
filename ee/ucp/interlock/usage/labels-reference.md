---
title: Layer 7 routing labels reference
description: Learn about the labels you can use in your swarm services to route
  layer 7 traffic to them.
keywords: routing, proxy
---

Once the layer 7 routing solution is enabled, you can
[start using it in your swarm services](index.md).

The following labels are available for you to use in swarm services:


| Label                                  | Description                                                                                                                                    | Example                |
|:---------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------|:-----------------------|
| `com.docker.lb.hosts`                  | Comma separated list of the hosts that the service should serve.                                                                               | `example.com,test.com` |
| `com.docker.lb.port`                   | Port to use for internal upstream communication.                                                                                               | `8080`                 |
| `com.docker.lb.network`                | Name of network the proxy service should attach to for upstream connectivity.                                                                  | `app-network-a`        |
| `com.docker.lb.context_root`           | Context or path to use for the application.                                                                                                    | `/app`                 |
| `com.docker.lb.context_root_rewrite`   | Boolean to enable rewrite for the context root.                                                                                                | `true`                 |
| `com.docker.lb.ssl_only`               | Boolean to force SSL for application.                                                                                                          | `true`                 |
| `com.docker.lb.ssl_cert`               | Docker secret to use for the SSL certificate.                                                                                                  | `example.com.cert`     |
| `com.docker.lb.ssl_key`                | Docker secret to use for the SSL key.                                                                                                          | `example.com.key`      |
| `com.docker.lb.websocket_endpoints`    | Comma separated list of endpoints to configure to be upgraded for websockets.                                                                  | `/ws,/foo`             |
| `com.docker.lb.service_cluster`        | Name of the service cluster to use for the application.                                                                                        | `us-east`              |
| `com.docker.lb.ssl_backend`            | Enable SSL communication to the upstreams.                                                                                                     | `true`                 |
| `com.docker.lb.ssl_backend_tls_verify` | Verification mode for the upstream TLS.                                                                                                        | `none`                 |
| `com.docker.lb.sticky_session_cookie`  | Cookie to use for sticky sessions.                                                                                                             | `none`                 |
| `com.docker.lb.redirects`              | Semi-colon separated list of redirects to add in the format of `<source>,<target>`.  Example: `http://old.example.com,http://new.example.com;` | `none`                 |
| `com.docker.lb.ssl_passthrough`        | Enable SSL passthrough.                                                                                                                        | `false`                |

