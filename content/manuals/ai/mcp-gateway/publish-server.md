---
title: Publish an MCP server
linkTitle: Publish a server
description: Register a hosted container or remote MCP endpoint, configure authentication, and share the server through MCP Gateway.
keywords: mcp gateway, publish mcp server, hosted container, remote endpoint, secrets
weight: 30
---

Publish an existing MCP server through the gateway console so people can use
its tools. You can run a container through the gateway service or connect a
server that's already running elsewhere.

## Before you start

You need permission to create servers, or superadmin membership. If
**Add server** is missing, ask a superadmin to enable server creation for your
role. Prepare the following:

| Server source | What you need |
| --- | --- |
| Hosted container | A versioned container image that serves MCP over HTTP, its port and MCP path, and any required environment variables |
| Remote endpoint | An MCP endpoint URL reachable from the gateway, and the transport it supports |

Check how the server authenticates. For gateway-managed OAuth with a
pre-registered client, obtain the client ID and, if the provider requires one,
the client secret. When registering that client with the provider, use the
gateway's upstream OAuth callback URL supplied by your deployment administrator.

If the server needs an API key or OAuth client secret, store it as a secret
before adding the server. In **Secrets**, select **New secret**, enter a name
and value, and select **Create secret**. You need
permission to create secrets, or someone with that permission can create and
share them with you.

Selecting a secret requires **Read** and **Read Secret Value** on it. Someone
who can use a secret in a server can also reveal its value. Use secret
references for credentials rather than plain-text environment variables or
header values.

## Add the server

1. Open **Home**, select **Servers**, and select **Add server**.
2. Enter a unique server name. It becomes part of endpoint URLs and can't be
   changed later. Add a display name and description to help people identify
   the server. If **Owner** is shown, choose yourself or a role you belong to.
3. Choose **Hosted container** or **Remote endpoint**, then select **Continue**.
4. Configure the runtime:
   - For a hosted container, enter the image, container port, and MCP path.
     Set resource requests and replicas as needed. Add environment variables,
     using **Secret reference** for credentials.
   - For a remote endpoint, enter the full MCP endpoint URL.
   - Choose **Streamable HTTP**, unless the server requires
     **Server-sent events (SSE)**.
5. Select **Continue** and configure authentication:
   - **Gateway managed** lets the gateway handle authentication. Users connect
     their service accounts through the console when required.
   - **Client managed** leaves authentication to the MCP client.
   - For gateway-managed authentication with a pre-registered OAuth client,
     enter **Client ID**. If the provider requires a client secret, use
     **Choose client secret** to select it. Leave both empty if the server uses
     dynamic client registration.
6. Leave **Expose per-server endpoint** selected if users should be able to
   connect to the server's own endpoint. Clear it to expose the server only
   through gateways that include it.
7. If the server needs an API key in a header, use **Add header secret reference**
   to supply it. Select **Continue**.
8. Review the configuration and select **Create and deploy**.

Passing configuration checks doesn't confirm that the image starts or the
remote endpoint is reachable. Wait for the deployment to become ready. If it
fails, use the server's **Status** tab to check readiness and tool discovery,
then correct the configuration in **Settings** and select **Update and deploy**.

## Make the server available

The server is visible only to its owner and superadmins until shared.
Members of an owning role also have access.

1. On the server's **Permissions** tab, select **Add permission**.
2. Choose a team role or person and grant **Read**. Add **Update** only if
   they should change the server's configuration.
3. Select **Save**.
4. [Add the server to a gateway](create-gateway.md#change-the-server-selection),
   or use its own endpoint under **Connection** if enabled.
5. [Connect a client](connect-client.md) and verify that the tools appear.

Sharing a server doesn't share its secrets or connect another person's
upstream account. Client-managed servers use their own direct connection
instructions rather than the combined gateway setup.
