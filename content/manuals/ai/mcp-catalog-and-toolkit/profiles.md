---
title: MCP Profiles
linkTitle: Profiles
description: Organize MCP servers into profiles for different projects and environments
keywords: Docker MCP, profiles, MCP servers, configuration, sharing
weight: 25
---

{{< summary-bar feature_name="MCP Profiles" >}}

> [!NOTE]
> This page describes the MCP Toolkit interface in Docker Desktop 4.62 and
> later. Earlier versions have a different UI. Upgrade to follow these
> instructions exactly.

Profiles organize your MCP servers into named collections. Without profiles,
you'd configure servers separately for every AI application you use. Each time
you want to change which servers are available, you'd update Claude Desktop, VS
Code, Cursor, and other tools individually. Profiles solve this by centralizing
your server configurations.

## What profiles do

A profile is a named collection of MCP servers with their configurations and
settings. You select servers from the [MCP
Catalog](/manuals/ai/mcp-catalog-and-toolkit/catalog.md) (the source of
available servers) and add them to your profiles (your configured server
collections for specific work). Think of the catalog as a library of tools, and
profiles as your toolboxes organized for different jobs.

Your "web-dev" profile might include GitHub, Playwright, and database servers.
Your "data-analysis" profile might include spreadsheet, API, and visualization
servers. Connect different AI clients to different profiles, or switch between
profiles as you change tasks.

When you run the MCP Gateway or connect a client without specifying a profile,
Docker MCP uses your default profile. If you're upgrading from a previous
version of MCP Toolkit, your existing server configurations are already in the
default profile.

## Profile capabilities

Each profile maintains its own isolated collection of servers and
configurations. Your "web-dev" profile might include GitHub, Playwright, and
database servers, while your "data-analysis" profile includes spreadsheet, API,
and visualization servers. Create as many profiles as you need, each containing
only the servers relevant to that context.

You can connect different AI applications to different profiles. When you
connect a client, you specify which profile it should use. This means Claude
Desktop and VS Code can have access to different server collections if needed.

Profiles can be shared with your team. Push a profile to your registry, and
team members can pull it to get the exact same server collection and
configuration you use.

## Creating and managing profiles

### Create a profile

1. In Docker Desktop, select **MCP Toolkit** and select the **Profiles** tab.
2. Select **Create profile**.
3. Enter a name for your profile (e.g., "web-dev").
4. Optionally, search and add servers to your profile now, or add them later.
5. Optionally, search and add clients to connect to your profile.
6. Select **Create**.

Your new profile appears in the profiles list.

### View profile details

Select a profile in the **Profiles** tab to view its details. The profile view
has two tabs:

- **Overview**: Shows the servers in your profile, secrets configuration, and
  connected clients. Use the **+** buttons to add more servers or clients.
- **Tools**: Lists all available tools from your profile's servers. You can
  enable or disable individual tools.

### Remove a profile

1. In the **Profiles** tab, find the profile you want to remove.
2. Select ⋮ next to the profile name, and then **Delete**.
3. Confirm the removal.

> [!CAUTION]
> Removing a profile deletes all its server configurations and settings, and
> updates the client configuration (removes MCP Toolkit). This action can't be
> undone.

### Default profile

When you run the MCP Gateway or use MCP Toolkit without specifying a profile,
Docker MCP uses a profile named `default`, or an empty configuration if a
`default` profile does not exist.

If you're upgrading from a previous version of MCP Toolkit, your existing
server configurations automatically migrate to the `default` profile. You don't
need to manually recreate your setup - everything continues to work as before.

You can always specify a different profile using the `--profile` flag with the
gateway command:

```console
$ docker mcp gateway run --profile web-dev
```

## Adding servers to profiles

Profiles contain the MCP servers you select from the catalog. Add servers to
organize your tools for specific workflows.

### Add a server

You can add servers to a profile in two ways.

From the Catalog tab:

1. Select the **Catalog** tab.
2. Select the checkbox next to servers you want to add to see which profile to
   add them to.
3. Choose your profile from the drop-down.

From within a profile:

1. Select the **Profiles** tab and select your profile.
2. In the **Servers** section, select the **+** button.
3. Search for and select servers to add.

If a server requires OAuth authentication, you're prompted to authorize it. See
[OAuth authentication](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md#oauth-authentication)
for details.

### List servers in a profile

Select a profile in the **Profiles** tab to see all servers it contains.

### Remove a server

1. Select the **Profiles** tab and select your profile.
2. In the **Servers** section, find the server you want to remove.
3. Select the delete icon next to the server.

## Configuring profiles

### Server configuration

Some servers require configuration beyond authentication. Configure server
settings within your profile.

1. Select the **Profiles** tab and select your profile.
2. In the **Servers** section, select the configure icon next to the server.
3. Adjust the server's configuration settings as needed.

### OAuth credentials

OAuth credentials are shared across all profiles. When you authorize access to
a service like GitHub or Notion, that authorization is available to any server
in any profile that needs it.

This means all profiles use the same OAuth credentials for a given service. If
you need to use different accounts for different projects, you'll need to
revoke and re-authorize between switching profiles.

See [OAuth authentication](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md#oauth-authentication)
for details on authorizing servers.

### Configuration persistence

Profile configurations persist in your Docker installation. When you restart
Docker Desktop or your system, your profiles, servers, and configurations
remain intact.

## Sharing profiles

Profiles can be shared with your team by pushing them to OCI-compliant
registries as artifacts. This is useful for distributing standardized MCP
setups across your organization. Credentials are not included in shared
profiles for security reasons. Team members configure OAuth separately after
pulling.

### Push a profile

1. Select the profile you want to share in the **Profiles** tab.
2. Select **Push to Registry**.
3. Enter the registry destination (e.g., `registry.example.com/profiles/web-dev:v1`).
4. Complete authentication if required.

### Pull a profile

1. Select **Pull from Registry** in the **Profiles** tab.
2. Enter the registry reference (e.g., `registry.example.com/profiles/team-standard:latest`).
3. Complete authentication if required.

The profile is downloaded and added to your profiles list. Configure any
required OAuth credentials separately.

### Team collaboration workflow

A typical workflow for sharing profiles across a team:

1. Create and configure a profile with the servers your team needs.
2. Test the profile to ensure it works as expected.
3. Push the profile to your team's registry with a version tag (e.g.,
   `registry.example.com/profiles/team-dev:v1`).
4. Share the registry reference with your team.
5. Team members pull the profile and configure any required OAuth credentials.

This ensures everyone uses the same server collection and configuration,
reducing setup time and inconsistencies.

## Using profiles with clients

When you connect an AI client to the MCP Gateway, you specify which profile's
servers the client can access.

### Run the gateway with a profile

Connect clients to your profile through the **Clients** section in the MCP
Toolkit. You can add clients when creating a profile or add them to existing
profiles later.

### Configure clients for specific profiles

When setting up a client manually, you can specify which profile the client
uses. This lets different clients connect to different profiles.

For example, your Claude Desktop configuration might use:

```json
{
  "mcpServers": {
    "MCP_DOCKER": {
      "command": "docker",
      "args": ["mcp", "gateway", "run", "--profile", "claude-work"]
    }
  }
}
```

While your VS Code configuration uses a different profile:

```json
{
  "mcp": {
    "servers": {
      "MCP_DOCKER": {
        "command": "docker",
        "args": ["mcp", "gateway", "run", "--profile", "vscode-dev"],
        "type": "stdio"
      }
    }
  }
}
```

### Switching between profiles

To switch the profile your clients use, update the client configuration to
specify a different `--profile` value in the gateway command arguments.

## Further reading

- [Get started with MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/get-started.md)
- [Use MCP Toolkit from the CLI](/manuals/ai/mcp-catalog-and-toolkit/cli.md)
- [MCP Catalog](/manuals/ai/mcp-catalog-and-toolkit/catalog.md)
- [MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md)
