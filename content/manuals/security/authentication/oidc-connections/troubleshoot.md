---
title: Troubleshoot OIDC connections
linkTitle: Troubleshoot
description: >-
  Use the Failures tab to diagnose OIDC connection token-exchange errors and
  fix common GitHub Actions authentication issues
weight: 30
keywords: >-
  oidc troubleshooting, oidc connection errors, Failures tab, github actions
  oidc, subject claims, rulesets, docker login-action,
  DOCKERHUB_OIDC_CONNECTIONID, access denied, connection_id is invalid,
  subject_token
tags: [Troubleshooting, admin]
toc_max: 2
---

{{< summary-bar feature_name="OIDC connections" >}}

When a GitHub Actions workflow fails to sign in to Docker with an OIDC
connection, open the connection's **Failures** tab. Compare the incoming
token claims with your rulesets, then use the sections below to fix common
errors.

## View connection failures

Docker records a failure when a workflow reaches the token-exchange step
for your connection and no ruleset matches the token's claims. For GitHub
Actions, Docker only logs the failure when the token's repository
namespace matches a namespace in a `sub` rule on the connection. An
unrelated repository can still receive `access denied` without a
**Failures** row.

1. Sign in to [Docker Home](https://app.docker.com/), select your
   organization, then go to **Identity & auth**.
1. Select **OIDC connections**.
1. Find the connection, open the action menu, and select **Edit**.
1. On the **Edit OIDC connection** page, select the **Failures** tab.

The **Failures** tab is view-only. It lists past failed exchanges so you
can compare claims to rulesets. It does not run a test or accept inputs.

If Docker has not recorded any failures for the connection, the tab shows
**No failures have been reported for this connection**.

> [!NOTE]
> Wrong connection IDs, deactivated connections, and token problems that
> fail before claim matching do not appear on **Failures**. Use the error
> sections later on this page for those cases.

### Read a failure entry

Each failure row includes:

- The date and time of the failed exchange
- The ruleset rule that did not match, shown as `claim = expected value`
  (for example, `sub = repo:my-org/my-repo:ref:refs/heads/main`)
- The incoming claim value from the GitHub ID token, marked with ❌ when
  it does not match

To inspect every claim in the token, expand **Show all claims**. The
expanded table lists **Incoming claim** and **Value**. Claims that failed
validation appear with ❌.

Use the unmatched expected value and the incoming value together to
decide whether to update the ruleset or the workflow that produced the
token.

## Subject claim does not match a ruleset

### Error message

In GitHub Actions, the login step fails during token exchange. Docker
returns:

```text
access denied
```

When Docker records the failure, the **Failures** tab shows the expected
subject claim next to the incoming `sub` value from GitHub. If there is
no failure row, confirm the workflow uses the correct connection ID and
that a ruleset `sub` pattern covers that repository's namespace.

### Causes

- The workflow's repository, branch, tag, pull request, or environment
  produces a different `sub` value than the ruleset expects.
- The ruleset uses a narrow pattern and the workflow runs from a
  different ref (for example, a feature branch instead of `main`).
- The ruleset still uses a mutable `repo:org/repo:...` subject while the
  GitHub repository issues an immutable subject claim. See
  [Subject claims](/manuals/security/authentication/oidc-connections/rulesets-claims.md#subject-claims).

### Solutions

1. Open the connection's **Failures** tab, if a row exists, and compare
   the expected `sub` value with the incoming value marked ❌.
1. Update the ruleset's **Subject claim** on the **Rulesets** tab so it
   matches the workflow, or use a wildcard where appropriate. See
   [Rulesets and subject claims](/manuals/security/authentication/oidc-connections/rulesets-claims.md).
1. Select **Save connection**, then re-run the GitHub Actions workflow.

## Connection ID is invalid or the connection is inactive

### Error message

In GitHub Actions, the login step fails during token exchange. Docker
returns:

```text
connection_id is invalid
```

If the workflow omits the connection ID entirely, Docker returns
`connection_id is required when subject_token_type is id_token`
instead.

`connection_id is invalid` also appears when the ID is malformed, does
not exist, or the connection is deactivated or deleted. Deactivated and
deleted connections are not available for token exchange, so Docker
treats the ID as invalid.

This failure does not add a row on the **Failures** tab.

### Causes

- `DOCKERHUB_OIDC_CONNECTIONID` is missing from the workflow, or its
  value does not match the connection ID shown in Docker Home.
- The connection was deactivated from the **OIDC connections** action
  menu. Inactive connections use the same invalid-ID error as unknown
  IDs.
- The connection was deleted and the workflow still references the old
  ID.

### Solutions

1. In Docker Home, open **OIDC connections** and confirm the connection
   status.
1. If the connection is inactive, open the action menu and select
   **Activate**.
1. If the connection was deleted, create a new connection and update
   `DOCKERHUB_OIDC_CONNECTIONID` in every affected workflow. See
   [Create and manage OIDC connections](/manuals/security/authentication/oidc-connections/create-manage.md).
1. If the connection is active, copy the connection ID from the **Edit
   OIDC connection** page and confirm the workflow uses that exact
   value.

## Subject token errors

### Error message

In GitHub Actions, the login step fails before or during token exchange.
Docker may return one of:

```text
subject_token is required
subject_token is invalid
subject_token has expired
subject_token audience mismatch
subject_token issuer mismatch
```

These failures do not appear on the **Failures** tab because the
exchange fails before Docker evaluates your rulesets.

### Causes

- The workflow does not supply a GitHub OIDC ID token to the login step.
- The token expired, is not yet valid, or does not match the expected
  issuer or audience.
- The token cannot be verified against the provider's keys.

### Solutions

1. Confirm the job grants GitHub permission to issue an OIDC token:

   ```yaml
   permissions:
     contents: read
     id-token: write
   ```

1. Use `docker/login-action` version 4.5.0 or later with your
   organization name and connection ID. See
   [Configure a GitHub Actions workflow](/manuals/security/authentication/oidc-connections/create-manage.md#configure-a-github-actions-workflow).
1. Re-run the workflow so GitHub issues a fresh ID token.

## Username is not an organization

### Causes

- The `username` input to `docker/login-action` is a personal Docker Hub
  account instead of an organization name.
- Only organization accounts can sign in with OIDC connections.

### Solutions

1. Set `username` to the Docker organization that owns the OIDC
   connection.
1. Confirm that organization matches the one where you created the
   connection in Docker Home.

## Related information

- [Create and manage OIDC connections](/manuals/security/authentication/oidc-connections/create-manage.md)
- [OIDC rulesets and subject claims](/manuals/security/authentication/oidc-connections/rulesets-claims.md)
- [OIDC connections overview](/manuals/security/authentication/oidc-connections/_index.md)
