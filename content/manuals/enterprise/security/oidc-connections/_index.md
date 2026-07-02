---
title: OIDC connections overview
linkTitle: OIDC connections
description: Learn how OIDC connections let organizations authenticate workloads and control access with OpenID Connect
keywords: oidc connections, openid connect, subject claims, rulesets, enterprise security, authentication, admin
tags: [admin]
weight: 35
params:
  sidebar:
    badge:
      color: blue
      text: Beta
---

{{< summary-bar feature_name="OIDC connections" >}}

OIDC connections create a trust relationship between Docker and a trusted third-party so you don't have to maintain long-lived credentials. When you create an OIDC connection, Docker exchanges short-lived tokens with another vendor that can grant fine-grained access to your Docker resources.

## How OIDC connections work

OIDC connections mirror implementations of the OIDC standard. Establishing a trust relationship involves creating the connection, configuring a workflow, and testing. For example, a trust relationship between Docker and GitHub follows these steps:

- GitHub issues a JWT ID token for the workflow run.
- During the authentication process, Docker then:
  - Verifies the token against GitHub's public key registry
  - Matches subject claims against specified rulesets created in the Admin Console.
- Docker returns an access token, allowing the GitHub Action login to Docker to access resources.

All tokens created and exchanged during an OIDC workflow are short-lived and issued on a per-workflow basis.

## OIDC connections and OATs

Organization access tokens (OATs) programmatically extend organization-level access to your Docker resources for all members. When membership changes in your organization, OATs ensure that access is granted or revoked without manual administrative oversight.

OIDC connections don't replace OATs. Rather, OIDC connections authenticate a workflow process as if it were a user, then extend authorization after authentication.

While OATs govern access to your Docker resources through organization membership, OIDC connections authenticate GitHub Actions workflows when they request a change to your Docker resources.

## What's next

- [Create an OIDC connection](/manuals/enterprise/security/oidc-connections/create-manage.md)
- [OIDC connections rulesets](/manuals/enterprise/security/oidc-connections/rulesets-claims.md)
