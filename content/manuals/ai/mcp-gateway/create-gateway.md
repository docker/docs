---
title: Create and share a gateway
linkTitle: Create a gateway
description: Choose MCP servers for a gateway and share it with people or teams through the gateway console.
keywords: mcp gateway, create gateway, share gateway, mcp servers
weight: 20
---

Create a gateway when you need a different selection of MCP servers or want
to provide a shared endpoint for a team. Creating a gateway selects from
existing servers. To register another server, see [Publish an MCP server](publish-server.md).

You need access to the gateway console. If permissions are enabled, a role
you belong to must permit gateway creation, or you must be a superadmin.
If **Create gateway** is missing, ask a superadmin to check your roles.

## Create a gateway

1. Select **Create gateway**, or the **+** next to **Gateways** in the sidebar.
2. If the form shows **Owner**, keep yourself or choose a role you belong to.
   Members of an owning role get owner rights.
3. Enter a name and an optional description. The name can't be changed later.
4. Select the servers whose tools you want. Servers labeled **Direct Auth**
   can't be included in a gateway you create.
5. Select **Create gateway** and wait for it to become ready.
6. [Connect your accounts and MCP client](connect-client.md).

You can select a server before connecting your account, but you must complete
any required sign-in before using its tools.

## Share with a team

If your deployment has permissions enabled, a gateway is visible only to its
owner and superadmins until shared. To share it, you must be its owner or a
superadmin.

1. Open the gateway and select **Permissions**.
2. Select **Add permission** and choose a role or enter a person.
3. Select **Read** so recipients can find the gateway. Add **Update** if they
   should maintain its server selection, or **Delete** if they should remove it.
4. Select **Save**.
5. For servers that recipients need to authenticate to through the console,
   also [grant them **Read** on those servers](publish-server.md#make-the-server-available).
   If you don't own a server, ask its owner or a superadmin to share it.

Gateway and server permissions are separate. Recipients need server **Read**
to find the server and connect their service account through the console.
They don't need **Update** for that sign-in. For how grants and ownership
work, see [Manage access](manage-access.md).

## Change the server selection

Open the gateway, select **Edit gateway**, change the selected servers or
description, and select **Save gateway**. Wait for the update to finish, then
refresh your client if it doesn't pick up the change.

The name and endpoint can't be changed. If the console says server changes
aren't available for that gateway, you can edit only its description.
For gateways labeled **Centrally Managed**, ask your deployment administrator
to change the configuration outside the console.
