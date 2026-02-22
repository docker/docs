---
title: Use MCP Toolkit from the CLI
linkTitle: Use with CLI
description: Manage MCP profiles, servers, and catalogs using the Docker MCP CLI.
keywords: docker mcp, cli, profiles, servers, catalog, gateway
weight: 35
---

{{< summary-bar feature_name="Docker MCP Toolkit" >}}

> [!NOTE]
> The `docker mcp` commands documented here are available in Docker Desktop
> 4.62 and later. Earlier versions may not support all commands shown.

The `docker mcp` commands let you manage MCP profiles, servers, OAuth
credentials, and catalogs from the terminal. Use the CLI for scripting,
automation, and headless environments.

## Profiles

### Create a profile

```console
$ docker mcp profile create --name <profile-id>
```

The profile ID is used to reference the profile in subsequent commands:

```console
$ docker mcp profile create --name web-dev
```

### List profiles

```console
$ docker mcp profile list
```

### View a profile

```console
$ docker mcp profile show <profile-id>
```

### Remove a profile

```console
$ docker mcp profile remove <profile-id>
```

> [!CAUTION]
> Removing a profile deletes all its server configurations and settings. This
> action can't be undone.

## Servers

### Browse the catalog

List available servers and their IDs:

```console
$ docker mcp catalog server ls mcp/docker-mcp-catalog
```

The output lists each server by name. The name (for example, `playwright` or
`github-official`) is the server ID to use in `catalog://` URIs.

To look up a server ID in Docker Desktop, open **MCP Toolkit** > **Catalog**,
select a server, and check the **Server ID** field.

### Add servers to a profile

Servers are referenced by URI. The URI format depends on where the server
comes from:

| Format | Source |
| --- | --- |
| `catalog://<catalog-ref>/<server-id>` | An OCI catalog |
| `docker://<image>:<tag>` | A Docker image |
| `https://<url>/v0/servers/<uuid>` | The MCP community registry |
| `file://<path>` | A local YAML or JSON file |

The most common format is `catalog://`, where `<catalog-ref>` matches the
**Catalog** field and `<server-id>` matches the **Server ID** field shown in
Docker Desktop or in the `catalog server ls` output:

```console
$ docker mcp profile server add <profile-id> \
  --server catalog://<catalog-ref>/<server-id>
```

Add multiple servers in one command:

```console
$ docker mcp profile server add web-dev \
  --server catalog://mcp/docker-mcp-catalog/github-official \
  --server catalog://mcp/docker-mcp-catalog/playwright
```

To add a server defined in a local YAML file:

```console
$ docker mcp profile server add my-profile \
  --server file://./my-server.yaml
```

The YAML file defines the server image and configuration:

```yaml
name: my-server
title: My Server
type: server
image: myimage:latest
description: Description of the server
```

If the server requires OAuth authentication, authorize it in Docker Desktop
after adding. See [OAuth authentication](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md#oauth-authentication).

### List servers

List all servers across all profiles:

```console
$ docker mcp profile server ls
```

Filter by profile:

```console
$ docker mcp profile server ls --filter profile=web-dev
```

### Remove a server

```console
$ docker mcp profile server remove <profile-id> --name <server-name>
```

Remove multiple servers at once:

```console
$ docker mcp profile server remove web-dev \
  --name github-official \
  --name playwright
```

### Configure server settings

Set and retrieve configuration values for servers in a profile:

```console
$ docker mcp profile config <profile-id> --set <server-id>.<key>=<value>
$ docker mcp profile config <profile-id> --get-all
$ docker mcp profile config <profile-id> --del <server-id>.<key>
```

Server configuration keys and their expected values are defined by each server.
Check the server's documentation or its entry in Docker Desktop under
**MCP Toolkit** > **Catalog** > **Configuration**.

## Gateway

Run the MCP Gateway with a specific profile:

```console
$ docker mcp gateway run --profile <profile-id>
```

Omit `--profile` to use the default profile.

### Connect a client manually

To connect any client that isn't listed in Docker Desktop, configure it to run
the gateway over `stdio`. For example, in a JSON-based client configuration:

```json
{
  "servers": {
    "MCP_DOCKER": {
      "command": "docker",
      "args": ["mcp", "gateway", "run", "--profile", "web-dev"],
      "type": "stdio"
    }
  }
}
```

For Claude Desktop, the format is:

```json
{
  "mcpServers": {
    "MCP_DOCKER": {
      "command": "docker",
      "args": ["mcp", "gateway", "run", "--profile", "web-dev"]
    }
  }
}
```

### Connect a named client

Connect a supported client to a profile:

```console
$ docker mcp client connect <client> --profile <profile-id>
```

For example, to connect VS Code to a project-specific profile:

```console
$ docker mcp client connect vscode --profile my-project
```

This creates a `.vscode/mcp.json` file in the current directory. Because this
is a user-specific file, add it to `.gitignore`:

```console
$ echo ".vscode/mcp.json" >> .gitignore
```

## Share profiles

Profiles are shared as OCI artifacts via any OCI-compatible registry.
Credentials are not included for security reasons. Team members configure
OAuth separately after pulling.

### Push a profile

```console
$ docker mcp profile push <profile-id> <registry-reference>
```

For example:

```console
$ docker mcp profile push web-dev registry.example.com/profiles/web-dev:v1
```

### Pull a profile

```console
$ docker mcp profile pull <registry-reference>
```

For example:

```console
$ docker mcp profile pull registry.example.com/profiles/team-standard:latest
```

## Custom catalogs

Custom catalogs let you curate a focused collection of servers for your team or
organization. For an overview of what custom catalogs are and when to use them,
see [Custom catalogs](/manuals/ai/mcp-catalog-and-toolkit/catalog.md#custom-catalogs).

Catalogs are referenced by OCI reference, for example
`registry.example.com/mcp/my-catalog:latest`. Servers within a catalog use
the same URI schemes as when
[adding servers to a profile](#add-servers-to-a-profile).

### Customize the Docker catalog

Use the Docker catalog as a base, then add or remove servers to fit your
organization's needs. Copy it first:

```console
$ docker mcp catalog tag mcp/docker-mcp-catalog registry.example.com/mcp/company-tools:latest
```

List the servers it contains:

```console
$ docker mcp catalog server ls registry.example.com/mcp/company-tools:latest
```

Remove servers your organization doesn't approve:

```console
$ docker mcp catalog server remove registry.example.com/mcp/company-tools:latest \
  --name <server-name>
```

Add your own private servers, packaged as Docker images:

```console
$ docker mcp catalog server add registry.example.com/mcp/company-tools:latest \
  --server docker://registry.example.com/mcp/internal-api:latest \
  --server docker://registry.example.com/mcp/data-pipeline:latest
```

Push when ready:

```console
$ docker mcp catalog push registry.example.com/mcp/company-tools:latest
```

### Build a catalog from scratch

To include exactly what you choose and nothing else, create a catalog from
scratch. You can include servers from the Docker catalog, your own private
images, or both.

Create a catalog and specify which servers to include:

```console
$ docker mcp catalog create registry.example.com/mcp/data-tools:latest \
  --title "Data Analysis Tools" \
  --server catalog://mcp/docker-mcp-catalog/postgres \
  --server catalog://mcp/docker-mcp-catalog/brave-search \
  --server docker://registry.example.com/mcp/analytics:latest
```

View the result:

```console
$ docker mcp catalog show registry.example.com/mcp/data-tools:latest
```

Push to distribute:

```console
$ docker mcp catalog push registry.example.com/mcp/data-tools:latest
```

### Distribute a catalog

Push your catalog so team members can import it:

```console
$ docker mcp catalog push <oci-reference>
```

Team members can pull it using the CLI:

```console
$ docker mcp catalog pull <oci-reference>
```

Or import it using Docker Desktop: select **MCP Toolkit** > **Catalog** >
**Import catalog** and enter the OCI reference.

### Use a custom catalog with the gateway

Run the gateway with your catalog instead of the default Docker catalog:

```console
$ docker mcp gateway run --catalog <oci-reference>
```

For [Dynamic MCP](/manuals/ai/mcp-catalog-and-toolkit/dynamic-mcp.md), where
agents discover and add servers during conversations, this limits what agents
can find to your curated set.

To enable specific servers from your catalog without using a profile:

```console
$ docker mcp gateway run --catalog <oci-reference> --servers <name1> --servers <name2>
```

## Further reading

- [Get started with MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/get-started.md)
- [MCP Profiles](/manuals/ai/mcp-catalog-and-toolkit/profiles.md)
- [MCP Catalog](/manuals/ai/mcp-catalog-and-toolkit/catalog.md)
- [MCP Gateway](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md)
