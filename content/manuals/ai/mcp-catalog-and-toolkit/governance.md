---
title: MCP governance and policy controls
linkTitle: Governance and policy
description: Control which MCP servers and tools users can access
keywords: mcp, governance, policy, security, access control, settings management
weight: 35
---

Organizations can define rules to control which MCP servers and tools
users can load and invoke. This lets administrators enforce security
policies and compliance requirements by allowing or denying specific
servers based on criteria like server type, source, or transport
mechanism.

Policies are configured through Settings Management using either the Admin
Console or the `admin-settings.json` file. Rules can apply globally to the
entire organization or target specific users. See
[Configuration](#configuration) for setup instructions.

## Policy structure

The policy is a JSON format evaluated by Open Policy Agent
([Rego](https://www.openpolicyagent.org/docs/policy-language)) under the hood.
A policy contains a set of rules that define conditions and actions. When an
MCP server attempts to load or a tool attempts to execute, the policy engine
evaluates each rule in order until it finds a match (first match wins). If a
rule matches the current context, the engine applies the rule's decision (allow
or deny). If no rules match, the default behavior applies.

Here's a basic policy that blocks a specific server:

```json
{
  "version": 1,
  "default": "allow",
  "rules": [
    {
      "action": "load",
      "server": "untrusted-server",
      "allow": false,
      "reason": "This server is not approved for organizational use"
    }
  ]
}
```

This policy allows all servers except `untrusted-server`. Rules override
the default behavior when they match.

### Policy schema

Policies consist of:

- `version` (optional): Policy schema version (currently `1`)
- `default`: Default action when no rules match ("allow" or "deny")
- `rules`: Array of rule objects that define specific conditions

Each rule can specify:

- `action`: Operation to control ("load", "invoke", "prompt")
- `server`: Specific server name to match
- `serverType`: Type of server ("registry", "image", "remote")
- `serverSource`: Source identifier for the server
  - For `registry` and `image` servers this is the image reference (for example `docker.io/mcp/github:1.2.3` or `ghcr.io/org/server:latest`)
  - For `remote` servers this is the endpoint URL (for example `https://mcp.example.com`)
- `transport`: Communication method ("stdio", "sse", "streamable")
- `tools`: Array of tool names to match
- `allow`: Whether to allow (`true`) or deny (`false`) the action. Can also use `deny: true` as an alternative to `allow: false`
- `reason`: Human-readable explanation for the rule

### Rule evaluation

The policy engine evaluates rules sequentially:

1. Check if a rule matches the current action and context
2. If a rule matches, apply its `allow` or `deny` decision
3. If no rules match, apply the `default` behavior

Order matters. Place more specific rules before general rules to ensure
they evaluate first. For example, a rule blocking a specific tool should
come before a rule allowing all tools on a server

## Policy examples

The following examples demonstrate common policy patterns. Each example
builds on the previous concepts, starting with basic server blocking and
progressing to comprehensive production policies.

### Block a specific server

Use default allow with explicit deny rules to block specific servers.
This approach requires minimal configuration while preventing access to
known problematic servers:

```json
{
  "version": 1,
  "default": "allow",
  "rules": [
    {
      "action": "load",
      "server": "untrusted-server",
      "allow": false,
      "reason": "This server is not approved for organizational use"
    }
  ]
}
```

### Default deny with allowlist

For security-sensitive environments, use default deny with an explicit
allowlist. This approach ensures users can only load pre-approved servers,
preventing unauthorized or untested servers from executing:

```json
{
  "version": 1,
  "default": "deny",
  "rules": [
    {
      "action": "load",
      "serverType": "registry",
      "serverSource": "docker.io/mcp/*",
      "allow": true,
      "reason": "Allow official Docker MCP servers from registry"
    },
    {
      "action": "load",
      "server": "github-official",
      "allow": true,
      "reason": "Allow GitHub official server"
    },
    {
      "action": "load",
      "server": "filesystem",
      "allow": true,
      "reason": "Allow filesystem access server"
    }
  ]
}
```

This policy:

- Blocks all servers by default
- Allows official Docker MCP servers from the registry
- Allows specific named servers that have been approved

### Block specific tools

Allow servers to load but restrict individual tools within them. This
provides granular control, permitting read-only operations while blocking
destructive actions:

```json
{
  "version": 1,
  "default": "allow",
  "rules": [
    {
      "action": "invoke",
      "server": "github-official",
      "tools": ["create_repository", "delete_repository"],
      "allow": false,
      "reason": "Block destructive GitHub operations"
    },
    {
      "action": "invoke",
      "server": "filesystem",
      "tools": ["write_file", "delete_file"],
      "allow": false,
      "reason": "Prevent file modifications"
    }
  ]
}
```

This policy allows servers to load but restricts specific tools that
could modify or delete data. Always specify a `server` field for tool rules
to avoid unintended matches when multiple servers have tools with the same name.

### Block servers from specific registries

Organizations using [custom catalogs](catalog.md#custom-catalogs) can
include servers from multiple OCI registries beyond Docker Hub (such as
GHCR, private registries, or community registries). Use policies to control
which registries users can load servers from:

- **`serverType: "registry"`** - Servers loaded from OCI registries (GHCR,
  private registries, etc.)
- **`serverType: "image"`** - Servers loaded from local Docker images
- **`serverSource`** - Matches the server image reference for `registry` and
  `image` servers, or the endpoint URL of `remote` servers. Supports wildcards.

This example blocks servers from GitHub Container Registry while allowing
other sources:

```json
{
  "version": 1,
  "default": "allow",
  "rules": [
    {
      "action": "load",
      "serverType": "registry",
      "serverSource": "ghcr.io/*",
      "allow": false,
      "reason": "Block servers from GitHub Container Registry"
    },
    {
      "action": "load",
      "serverType": "image",
      "allow": false,
      "reason": "Block all local image-based servers"
    }
  ]
}
```

This policy:

- Blocks servers from specific OCI registries (like GitHub Container Registry)
- Blocks all local image-based servers
- Allows servers from other sources (like Docker Hub)

### Filter by transport mechanism

Restrict communication protocols based on security requirements. Block
remote servers or specific transports to limit network exposure:

```json
{
  "version": 1,
  "default": "allow",
  "rules": [
    {
      "action": "load",
      "transport": "sse",
      "allow": false,
      "reason": "Block server-sent events transport for security"
    },
    {
      "action": "load",
      "serverType": "remote",
      "allow": false,
      "reason": "Block all remote servers"
    }
  ]
}
```

This policy restricts transport mechanisms and prevents remote server
connections.

### Combined policy for production

A comprehensive policy for production environments:

```json
{
  "version": 1,
  "default": "deny",
  "rules": [
    {
      "action": "load",
      "serverType": "registry",
      "serverSource": "docker.io/mcp/*",
      "allow": true,
      "reason": "Allow official Docker MCP servers"
    },
    {
      "action": "load",
      "serverType": "registry",
      "serverSource": "internal-registry.company.com/mcp/*",
      "allow": true,
      "reason": "Allow company-approved internal servers"
    },
    {
      "action": "invoke",
      "tools": ["delete_*", "remove_*", "destroy_*"],
      "allow": false,
      "reason": "Block destructive operations across all servers"
    },
    {
      "action": "load",
      "serverType": "remote",
      "allow": false,
      "reason": "Block remote servers for security"
    }
  ]
}
```

This production policy:

- Denies all servers by default
- Allows only official Docker and internal company servers
- Blocks tools with destructive names
- Prevents remote server connections

## Configuration

Configure MCP policies using either the Admin Console or the `admin-settings.json` file.

### Admin Console

To apply a policy through the Admin Console:

1. Open the [Admin Console](https://app.docker.com) and navigate to Settings
   Management.
2. Select the policy that you want to update.
3. Find the **AI > MCP Policy** option.
4. Enter your policy JSON into the configuration field.
5. Save and distribute the policy.

For detailed instructions, see [Configure with Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md).

### admin-settings.json file

For automated deployments or scripted installations, use the
[`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
placed on end-user machines.

To configure an MCP policy in `admin-settings.json`, add the policy as a
JSON-escaped string in the `mcpPolicy` field:

```json
{
  "configurationFileVersion": 2,
  "mcpPolicy": {
    "locked": true,
    "value": "{\"version\":1,\"default\":\"deny\",\"rules\":[{\"action\":\"load\",\"server\":\"approved-server\",\"allow\":true,\"reason\":\"Approved for use\"}]}"
  }
}
```

> [!NOTE]
> The policy JSON must be escaped as a string value. Admin Console policies
> take precedence over `admin-settings.json` policies.

### Policy propagation

Policy updates may take a moment to propagate to Docker Desktop. If a policy
change doesn't take effect immediately, restart Docker Desktop to apply the
updated policy.

## Policy strategies

Choose a policy approach based on your security requirements:

**Default deny**: Use `"default": "deny"` for hardened environments.
Explicitly allow only approved servers and tools. This prevents unauthorized
servers from loading.

With default deny, you must create rules for both `"action": "load"` and
`"action": "invoke"` to enable tools. For example:

```json
{
  "version": 1,
  "default": "deny",
  "rules": [
    {
      "action": "load",
      "server": "github-official",
      "allow": true,
      "reason": "Allow GitHub server to load"
    },
    {
      "action": "invoke",
      "server": "github-official",
      "allow": true,
      "reason": "Allow invoking GitHub tools"
    }
  ]
}
```

Without both rules, the server loads but its tools remain blocked.

**Default allow**: Use `"default": "allow"` for development environments.
Block specific problematic servers while permitting experimentation with new
tools.

Use wildcards (`*`) in `serverSource` to match multiple servers from the
same registry or organization. Combine action types (load and invoke) for
layered security - allow a server to load but block specific destructive
tools within it.
