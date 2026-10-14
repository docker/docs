---
title: "API Tool"
description: "Create custom tools that call HTTP APIs."
keywords: docker agent, ai agents, tools, toolsets, api tool
linkTitle: "API"
weight: 240
canonical: https://docs.docker.com/ai/docker-agent/tools/api/
---

_Create custom tools that call HTTP APIs._

## Overview

The API tool type lets you define custom tools that make HTTP requests to external APIs. This is useful for integrating agents with REST APIs, webhooks, or any HTTP-based service without writing code.

> [!NOTE]
> **When to Use**
>
> - Integrating with REST APIs that don't have an MCP server
> - Simple HTTP operations (GET, POST)
> - Quick prototyping before building a full MCP server

## Configuration

```yaml
agents:
  assistant:
    model: openai/gpt-4o
    description: Assistant with API access
    instruction: You can look up weather information.
    toolsets:
      - type: api
        api_config:
          name: get_weather
          method: GET
          endpoint: "https://api.weather.example/v1/current?city=${city}"
          instruction: Get current weather for a city
          args:
            city:
              type: string
              description: City name to get weather for
          required: ["city"]
          headers:
            Authorization: "Bearer ${env.WEATHER_API_KEY}"
```

## Properties

The `api` toolset accepts the following toolset-level fields in addition to the `api_config` block:

| Property            | Type    | Required | Description                                                                                                                                                                                                                                                       |
| ------------------- | ------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `api_config`        | object  | ✓        | The HTTP tool definition. See the table below.                                                                                                                                                                                                                    |
| `timeout`           | int     | ✗        | HTTP client timeout in seconds (default: `30`). Applies to every call the generated tool makes.                                                                                                                                                                   |
| `allow_private_ips` | boolean | ✗        | Opt in to dialling **non-public** IP addresses (loopback, RFC1918, link-local — including the cloud-metadata endpoint at `169.254.169.254` — multicast and the unspecified address). Set to `true` only when the configured endpoint legitimately targets internal services. See [Reaching internal services](#reaching-internal-services). |

### `api_config`

| Property        | Type   | Required | Description                                      |
| --------------- | ------ | -------- | ------------------------------------------------ |
| `name`          | string | ✓        | Tool name (how the agent references it)          |
| `method`        | string | ✓        | HTTP method: `GET` or `POST`                     |
| `endpoint`      | string | ✓        | URL endpoint (supports `${param}` interpolation) |
| `instruction`   | string | ✗        | Description shown to the agent                   |
| `args`          | object | ✗        | Parameter definitions (JSON Schema properties)   |
| `required`      | array  | ✗        | List of required parameter names                 |
| `headers`       | object | ✗        | HTTP headers to include. Values support `${env.VAR}` and `${headers.NAME}` placeholders (the latter forwards a header from the caller's incoming request, useful when docker agent is itself exposed as an HTTP server). |
| `output_schema` | object | ✗        | JSON Schema for the response. Used by MCP / Code Mode consumers; tool responses are still returned to the model as raw strings.                                                                                          |

## HTTP Methods

### GET Requests

For GET requests, parameters are interpolated into the URL:

```yaml
toolsets:
  - type: api
    api_config:
      name: search_users
      method: GET
      endpoint: "https://api.example.com/users?q=${query}&limit=${limit}"
      instruction: Search for users by name
      args:
        query:
          type: string
          description: Search query
        limit:
          type: integer
          description: Maximum results (default 10)
      required: ["query"]
```

### POST Requests

For POST requests, parameters are sent as JSON in the request body:

```yaml
toolsets:
  - type: api
    api_config:
      name: create_task
      method: POST
      endpoint: "https://api.example.com/tasks"
      instruction: Create a new task
      args:
        title:
          type: string
          description: Task title
        description:
          type: string
          description: Task description
        priority:
          type: string
          enum: ["low", "medium", "high"]
          description: Task priority
      required: ["title"]
      headers:
        Content-Type: "application/json"
        Authorization: "Bearer ${env.API_TOKEN}"
```

## URL Interpolation

Use `${param}` syntax to insert parameter values into URLs:

```yaml
endpoint: "https://api.example.com/users/${user_id}/posts/${post_id}"
```

Parameter values are inserted as strings by the template expansion. Add URL encoding in the template when needed (for example, `${encodeURIComponent(city)}`).

## Headers

Headers can include environment variables:

```yaml
headers:
  Authorization: "Bearer ${env.API_KEY}"
  X-Custom-Header: "static-value"
  Content-Type: "application/json"
```

## Output Schema

Optionally document the expected response format:

```yaml
toolsets:
  - type: api
    api_config:
      name: get_user
      method: GET
      endpoint: "https://api.example.com/users/${id}"
      instruction: Get user details by ID
      args:
        id:
          type: string
          description: User ID
      required: ["id"]
      output_schema:
        type: object
        properties:
          id:
            type: string
          name:
            type: string
          email:
            type: string
          created_at:
            type: string
```

## Example: GitHub API

```yaml
agents:
  github_assistant:
    model: openai/gpt-4o
    description: Assistant that can query GitHub
    instruction: You can look up GitHub repositories and users.
    toolsets:
      - type: api
        api_config:
          name: get_repo
          method: GET
          endpoint: "https://api.github.com/repos/${owner}/${repo}"
          instruction: Get information about a GitHub repository
          args:
            owner:
              type: string
              description: Repository owner (user or org)
            repo:
              type: string
              description: Repository name
          required: ["owner", "repo"]
          headers:
            Accept: "application/vnd.github.v3+json"
            Authorization: "Bearer ${env.GITHUB_TOKEN}"

      - type: api
        api_config:
          name: get_user
          method: GET
          endpoint: "https://api.github.com/users/${username}"
          instruction: Get information about a GitHub user
          args:
            username:
              type: string
              description: GitHub username
          required: ["username"]
          headers:
            Accept: "application/vnd.github.v3+json"
```

## Limitations

- Only supports GET and POST methods
- Response body is limited to 1MB
- Default 30-second timeout per request (override with the `timeout` field)
- Only HTTP and HTTPS URLs are supported
- No support for file uploads or multipart forms
- By default, requests to non-public IP ranges (loopback, RFC1918, link-local, the cloud-metadata endpoint, multicast, the unspecified address) are refused at dial time — even when DNS for an otherwise-public host resolves there. Set `allow_private_ips: true` to disable that check.

## Reaching internal services

```yaml
toolsets:
  - type: api
    timeout: 60
    allow_private_ips: true
    api_config:
      name: get_local_status
      method: GET
      endpoint: "http://localhost:8080/health"
      instruction: Check the local service health
```

> [!WARNING]
> **SSRF**
>
> Setting `allow_private_ips: true` re-exposes the SSRF surface for this tool. Only enable it when the configured `endpoint` is a trusted internal service — a prompt-injected agent cannot redirect the call elsewhere because the endpoint is fixed in config, but redirects from the configured host can still reach unexpected places.

> [!TIP]
> **For Complex APIs**
>
> For APIs that need authentication flows, pagination, or complex request/response handling, consider using an MCP server instead. The API tool is best for simple, stateless HTTP operations.

> [!WARNING]
> **Security**
>
> API keys and tokens in headers are visible in debug logs. Use environment variables (`${env.VAR}`) rather than hardcoding secrets in configuration files.
