---
title: Manage webhooks
description: Learn how to create, configure, and test webhooks in Docker Trusted Registry.
keywords: registry, webhooks
redirect_from:
  - /datacenter/dtr/2.5/guides/user/create-and-manage-webhooks/
  - /ee/dtr/user/create-and-manage-webhooks/
---

>{% include enterprise_label_shortform.md %}

You can configure DTR to automatically post event notifications to a webhook URL of your choosing. This lets you build complex CI and CD pipelines with your Docker images. The following is a complete list of event types you can trigger webhook notifications for via the [web interface](use-the-web-ui) or the [API](use-the-API).

## Webhook types

| Event Type                              | Scope                   | Access Level     | Availability |
| --------------------------------------- | ----------------------- | ---------------- | ------------ |
| Tag pushed to repository (`TAG_PUSH`)               | Individual repositories | Repository admin | Web UI & API       |
| Tag pulled from repository (`TAG_PULL`)           | Individual repositories | Repository admin | Web UI & API         |
| Tag deleted from repository (`TAG_DELETE`)            | Individual repositories | Repository admin | Web UI & API         |
| Manifest pushed to repository (`MANIFEST_PUSH`)        | Individual repositories | Repository admin | Web UI & API          |
| Manifest pulled from repository (`MANIFEST_PULL`)     | Individual repositories | Repository admin | Web UI & API          |
| Manifest deleted from repository (`MANIFEST_DELETE`)      | Individual repositories | Repository admin | Web UI & API         |
| Security scan completed (`SCAN_COMPLETED`)           | Individual repositories | Repository admin | Web UI & API         |
| Security scan failed (`SCAN_FAILED`)                    | Individual repositories | Repository admin | Web UI & API         |
| Image promoted from repository (`PROMOTION`)          | Individual repositories | Repository admin | Web UI & API         |
| Image mirrored from repository (`PUSH_MIRRORING`)         | Individual repositories | Repository admin | Web UI & API          |
| Image mirrored from remote repository (`POLL_MIRRORING`)   | Individual repositories | Repository admin | Web UI & API         |
| Repository created, updated, or deleted (`REPO_CREATED`, `REPO_UPDATED`, and `REPO_DELETED`) | Namespaces / Organizations     | Namespace / Org owners  | API Only     |
| Security scanner update completed (`SCANNER_UPDATE_COMPLETED`)                                    |        Global               | DTR admin                 |         API only     |

You must have admin privileges to a repository or namespace in order to
subscribe to its webhook events. For example, a user must be an admin of repository "foo/bar" to subscribe to its tag push events. A DTR admin can subscribe to any event.

## Where to go next

- [Manage webhooks via the web interface](use-the-web-ui)
- [Manage webhooks via the the API](use-the-api)
