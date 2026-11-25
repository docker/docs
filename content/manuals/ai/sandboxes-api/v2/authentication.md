---
title: Connect and authenticate with API v2
description: Understand local daemon connections, cloud authentication, and workload credentials.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 20
linkTitle: Connect and authenticate
---

Where your sandbox runs determines how your application connects to its API.
Docker account credentials and credentials used by an agent serve different
purposes.

## Connect locally

Local requests go to the `sbx` daemon through local interprocess communication,
such as a Unix socket. Follow the daemon's local connection and sign-in
requirements when configuring your application.

<!-- Add verified daemon discovery, supported IPC transports, sign-in
requirements, and SDK connection configuration after release validation. -->

## Connect to cloud

Cloud requests go to a hosted service and authenticate your application with
Docker account credentials. The account also needs
[cloud access](/manuals/agentic-platform/signup.md#activate-cloud-access).

<!-- Add supported interactive and automation authentication, token renewal,
service URLs, and permissions. Do not copy v1 authentication helpers or URLs. -->

## Authenticate an agent

An agent can need credentials for a model provider or another service in
addition to your application's Docker authentication. A Docker sign-in does
not sign the agent in to its model provider.

Kits can declare the credentials their tools need. See
[Connect agents to services](recipes/connect-agents-to-services.md).
