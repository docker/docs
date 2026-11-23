---
title: Connect agents to services with API v2
description: Understand cloud credential bindings for services declared by sandbox kits.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 40
linkTitle: Connect agents to services
---

An agent's tools can need credentials for model providers, source repositories,
or other services. Kits declare those requirements so an application can supply
credentials when it creates the sandbox.

For cloud sandboxes, a credential binding connects a stored secret to a service
declared by a kit. The secret must match the declared service, and each binding
must answer a credential requirement in the kit composition.

Managed credential bindings are not exposed by the local API subset described
in [Local and cloud capabilities](../capabilities.md). Keep local credential
configuration separate from the cloud example.

See the [kit capability documentation](/manuals/ai/sandboxes/customize/author/capabilities.md)
for credential declarations, and
[Connect and authenticate](../authentication.md) for the distinction between
application and agent authentication.

<!-- Add a cloud secret-create, kit-preview, credential-bind, create, and cleanup
example. Verify the release's service matching and injection rules. Do not copy
v1 storage.secrets examples. -->
