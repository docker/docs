---
description: Docker Hub API changelog
title: Docker Hub API changelog
linkTitle: Changelog
keywords: docker hub, hub, whats new, release notes, api, changelog
weight: 2
toc_min: 1
toc_max: 2
aliases:
  - /reference/api/hub/latest-changelog/
---

Here you can learn about the latest changes, new features, bug fixes, and known
issues for Docker Service APIs.

---

## 2025-09-25

### Updates

- Fix [Assign repository group](/reference/api/hub/latest/#tag/repositories/operation/CreateRepositoryGroup) endpoints request/response

---

## 2025-09-19

### New

- Add [Create repository](/reference/api/hub/latest/#tag/repositories/operation/CreateRepository) endpoints for a given `namespace`.
- Add [Get repository](/reference/api/hub/latest/#tag/repositories/operation/GetRepository) endpoints for a given `namespace`.
- Add [Check repository](/reference/api/hub/latest/#tag/repositories/operation/CheckRepository) endpoints for a given `namespace`.

### Deprecations

- [Deprecate POST /v2/repositories](/reference/api/hub/deprecated/#deprecate-legacy-createrepository)
- [Deprecate POST /v2/repositories/{namespace}](/reference/api/hub/deprecated/#deprecate-legacy-createrepository)
- [Deprecate GET /v2/repositories/{namespace}/{repository}](/reference/api/hub/deprecated/#deprecate-legacy-getrepository)
- [Deprecate HEAD /v2/repositories/{namespace}/{repository}](/reference/api/hub/deprecated/#deprecate-legacy-getrepository)

---

## 2025-07-29

### New

- Add [Update repository immutable tags settings](/reference/api/hub/latest/#tag/repositories/operation/UpdateRepositoryImmutableTags) endpoints for a given `namespace` and `repository`.
- Add [Verify repository immutable tags](/reference/api/hub/latest/#tag/repositories/operation/VerifyRepositoryImmutableTags) endpoints for a given `namespace` and `repository`.

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
