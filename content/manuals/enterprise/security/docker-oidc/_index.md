---
title: Docker OIDC overview
linkTitle: Docker OIDC
description: Learn how Docker OIDC connections let organizations authenticate workloads and control access with OpenID Connect
keywords: docker oidc, openid connect, oidc connections, subject claims, rulesets, enterprise security, authentication, admin
tags: [admin]
weight: 35
---

{{< summary-bar feature_name="Docker OIDC" >}}

Docker OIDC creates a trust relationship between GitHub and Docker so you don’t have to maintain long-lived credentials. When you create an OIDC connection, Docker and GitHub exchange short-lived tokens that still grant fine-grained access to your Docker resources.

Docker OIDC is available for organizations with Docker Core subscriptions. You receive 10 OIDC connections out of the box. To upgrade your subscription, see [Change your subscription](/manuals/subscription/details.md).

## Prerequisites

To create an OIDC connection, you need:

- A Docker Core subscription
- Organization ownership

## How Docker OIDC works

Docker OIDC mirrors implementations of the OIDC standard. Establishing the trust relationship between GitHub and Docker involves broad phases:

- GitHub issues a JWT ID token for the workflow run.
- During the authentication process, Docker then:
  - Verifies the token against GitHub’s public key registry
  - Matches subject claims against specified rulesets created in the Admin Console.
- Docker returns an access token, allowing the GitHub Action login to Docker to access resources.

All tokens created and exchanged during an OIDC workflow are short-lived and issued on a per-GitHub Action basis.

## Docker OIDC and OATs

Organization access tokens (OATs) programmatically extend organization-level access to your Docker resources for all members. When membership changes in your organization, OATs ensure that access is granted or revoked without manual administrative oversight.

Docker OIDC doesn’t replace OATs. Rather, Docker OIDC authenticates a workflow process as if it were a user, then extends authorization after authentication.

While OATs govern access to your Docker resources through organization membership, Docker OIDC authenticates GitHub Action workflows when they request a change to your Docker resources.

## What’s next

- [Create an OIDC connection](/manuals/enterprise/security/docker-oidc/create-manage.md)
- Refer to [Docker OIDC rulesets](/manuals/enterprise/security/docker-oidc/rulesets-claims.md).
