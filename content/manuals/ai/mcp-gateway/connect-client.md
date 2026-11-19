---
title: Connect and use MCP tools
linkTitle: Connect a client
description: Connect your service accounts and add an MCP Gateway endpoint to your AI client.
keywords: mcp gateway, mcp client, oauth, authentication, connect tools
weight: 10
---

Connect your MCP client to a gateway to use tools from its selected servers.
You need access to the gateway console and a gateway provided by your
administrator or [created by you](create-gateway.md).

These steps use the gateway console. If you use Docker Agentic Platform,
follow its [MCP connection instructions](/manuals/agentic-platform/mcp.md).

## Connect your service accounts

Signing in to the console doesn't sign you in to every service. Each server
has its own authentication requirements.

1. Open **Home**, select the **Gateways** tab, and find the gateway you want
   to use. Check which servers it includes.
2. In the **Servers** tab, open each server's details and check its
   authentication status. If a server is missing, ask its owner or a
   superadmin for **Read** on the server. Access to a shared gateway doesn't
   grant this permission automatically.
3. For a server showing **Needs auth**, select **Authenticate** and complete
   the service's sign-in and consent flow.
4. Return to the console and check that the server is ready. Repeat for each
   server that needs your account.

Some servers have credentials configured by an administrator or require no
upstream authentication. These don't need a separate sign-in from you. The
gateway's own authentication requirements still apply.

Servers labeled **Direct Auth** use a different path: add the server's own URL
to your client using the instructions in its details, and sign in from the
client. The console doesn't manage credentials for that direct connection.

## Add the gateway to your client

1. On **Home**, select **Gateways** and wait for your gateway to be ready.
2. Select **Add to your client** for that gateway.
3. Choose your MCP client and follow the instructions in the dialog. Use its
   configuration snippet or the endpoint shown there. Run any setup command
   on the computer where the client is installed.
4. Complete sign-in if prompted by your client.
5. Refresh or restart your client if needed, then check that the expected
   tools appear.

Use the endpoint from your own setup dialog. Different gateways can expose
different servers, even within the same deployment.

Once connected, use the tools from your client. Return to the console to
reconnect an account or change your server selection.

## If tools are missing

- Compare the configured endpoint with the gateway's **Add to your client**
  dialog.
- Check that the server is included in that gateway. Ask your administrator
  to change a centrally managed gateway.
- Check the server's authentication status and reconnect your account if
  needed. For Direct Auth, check sign-in in your client.
- Refresh the client so it can discover the tools again.

If these checks don't resolve the problem, give your administrator the gateway
name, server name, and client error. Don't include credentials.

Before using tools that change data, review the gateway's
[security boundaries](manage-access.md#security-boundaries).
