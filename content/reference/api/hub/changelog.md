---
description: Docker Hub API changelog
title: Docker Hub API changelog
linkTitle: Changelog
keywords: docker hub, hub, whats new, release notes, api, changelog
weight: 1
toc_min: 1
toc_max: 2
---

Here you can learn about the latest changes, new features, bug fixes, and known
issues for Docker Service APIs.

---

## 2025-06-27

### New

- Add [List repositories](/reference/api/hub/latest/#tag/repositories/operation/listNamespaceRepositories) endpoints for a given `namespace`.

### Deprecations

- [Deprecate /v2/repositories/{namespace}](/reference/api/hub/deprecated/#deprecate-legacy-listnamespacerepositories)

---

## 2025-03-25

### New

- Add [APIs](/reference/api/hub/latest/#tag/org-access-tokens) for organization access token (OATs) management.

---

## 2025-03-18

### New

- Add access to [audit logs](/reference/api/hub/latest/#tag/audit-logs) for org
  access tokens.
