---
title: "OpenAPI Tool"
description: "Automatically generate tools from an OpenAPI specification."
keywords: docker agent, ai agents, tools, toolsets, openapi tool
linkTitle: "OpenAPI"
weight: 230
canonical: https://docs.docker.com/ai/docker-agent/tools/openapi/
---

_Automatically generate tools from an OpenAPI specification._

## Overview

The OpenAPI tool fetches an OpenAPI 3.x specification from a URL and creates one tool per API operation. Each endpoint's parameters, request body, and description are translated into a callable tool that the agent can invoke directly.

## Configuration

```yaml
toolsets:
  - type: openapi
    url: "https://petstore3.swagger.io/api/v3/openapi.json"
```

### With custom headers

Pass custom headers to every HTTP request made by the generated tools (for example, for authentication):

```yaml
toolsets:
  - type: openapi
    url: "https://api.example.com/openapi.json"
    headers:
      Authorization: "Bearer ${env.API_TOKEN}"
      X-Custom-Header: "my-value"
```

### Custom timeout

Override the default 30-second HTTP timeout (applies both to fetching the spec and to the generated tool calls):

```yaml
toolsets:
  - type: openapi
    url: "https://api.example.com/openapi.json"
    timeout: 60
```

### Reaching internal services

By default the OpenAPI tool refuses connections to non-public IP addresses, blocking SSRF attempts even when DNS resolves an otherwise-public host to an internal range. Opt in with `allow_private_ips` when the spec or its `servers` entries legitimately target localhost or your internal network:

```yaml
toolsets:
  - type: openapi
    url: "http://localhost:8080/openapi.json"
    allow_private_ips: true
```

## Properties

| Property            | Type              | Required | Description                                                                                                                                                                                                                                                       |
| ------------------- | ----------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `url`               | string            | ✓        | URL of the OpenAPI specification (JSON format). Supports `${env.VAR}` interpolation.                                                                                                                                                                              |
| `headers`           | map[string]string | ✗        | Custom HTTP headers sent with every request — both the spec fetch and every generated tool call. Values support `${env.VAR}` and `${headers.NAME}` placeholders (the latter forwards a header from the caller's incoming request when docker agent is exposed as a server). |
| `timeout`           | int               | ✗        | HTTP client timeout in seconds (default: `30`). Applies to both the spec fetch and the generated tools' requests.                                                                                                                                                 |
| `allow_private_ips` | boolean           | ✗        | Opt in to dialling **non-public** IP addresses (loopback, RFC1918, link-local — including the cloud-metadata endpoint at `169.254.169.254` — multicast and the unspecified address). Set to `true` only when the spec or its servers legitimately target internal services. By default such addresses are refused at dial time, after DNS resolution, so DNS rebinding cannot bypass the check. |

## How it works

1. The spec is fetched from the configured `url` at startup.
2. Each operation (GET, POST, PUT, …) becomes a separate tool named after its `operationId` (or `method_path` when no `operationId` is set).
3. Path and query parameters are exposed as tool parameters. Request body properties are prefixed with `body_`.
4. Read-only operations (GET, HEAD, OPTIONS) are annotated accordingly.
5. Responses are returned as text; errors include the HTTP status code.

## Limits

- The OpenAPI spec must be **10 MB or less**.
- Individual API responses are truncated at **1 MB**.

## Example

See the full [Pet Store example](https://github.com/docker/docker-agent/blob/main/examples/openapi-petstore.yaml) for a working agent configuration.
