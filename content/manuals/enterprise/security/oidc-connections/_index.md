---
title: OIDC connections overview
linkTitle: OIDC connections
description: Authenticate GitHub Actions to Docker with short-lived OpenID Connect tokens
keywords: oidc connections, openid connect, github actions, jwt, subject claims, rulesets, enterprise security, workload authentication
tags: [admin]
weight: 35
---

{{< summary-bar feature_name="OIDC connections" >}}

OIDC connections create a trust relationship between Docker and a trusted
third party so you don't have to maintain long-lived credentials. When you
create an OIDC connection, Docker exchanges short-lived tokens with another
vendor that can grant fine-grained access to your Docker resources.

## How OIDC connections work

OIDC connections follow the OpenID Connect (OIDC) standard. Establishing a
trust relationship involves creating the connection, configuring a
workflow, and testing. For example, a trust relationship between Docker and
GitHub follows these steps:

- GitHub issues a JWT ID token for the workflow run.
- During authentication, Docker:
  - Verifies the token against GitHub's public key registry
  - Matches subject claims against rulesets created in
    [Docker Home](https://app.docker.com/)
- Docker returns an access token so the GitHub Action can sign in to Docker
  and access resources.

All tokens created and exchanged during an OIDC workflow are short-lived
and issued on a per-workflow basis.

## OIDC connections and OATs

[Organization access tokens (OATs)](/manuals/enterprise/security/access-tokens.md)
provide programmatic access to your Docker resources at the organization
level. Unlike personal access tokens, OATs aren't tied to individual
members, so access continues when membership changes.

OIDC connections don't replace OATs. OIDC connections authenticate a
workflow as if it were a user, then authorize access after authentication.

While OATs govern access to your Docker resources through organization
membership, OIDC connections authenticate GitHub Actions workflows when
they request a change to your Docker resources.

## What's next

- [Create an OIDC connection](/manuals/enterprise/security/oidc-connections/create-manage.md)
- [OIDC rulesets and subject claims](/manuals/enterprise/security/oidc-connections/rulesets-claims.md)
