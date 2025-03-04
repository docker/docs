---
description: Learn how to troubleshoot common Docker Hub issues.
keywords: hub, troubleshoot
title: Troubleshoot Docker Hub
linkTitle: Troubleshoot
weight: 60
tags: [Troubleshooting]
---

If you experience issues with Docker Hub, refer to the following solutions.

## You have reached your pull rate limit (429 response code)

### Error message

When this issue occurs, you receive following error message in the Docker CLI or
in the Docker Engine logs:

```text
You have reached your pull rate limit. You may increase the limit by authenticating and upgrading: https://www.docker.com/increase-rate-limits
```

### Possible causes

- You have reached your pull rate limit as an authenticated Docker Personal
  user.
- You have reached your pull rate limit as an unauthenticated user based on your
  IPv4 address or IPv6 /64 subnet.

### Solution

You can use one of the following solutions:

- [Authenticate](./usage/pulls.md#authentication) or
  [upgrade](../subscription/change.md#upgrade-your-subscription) your Docker
  account.
- [View your pull rate limit](./usage/pulls.md#view-hourly-pull-rate-and-limit),
  wait until your pull rate limit decreases, and then try again.

## Too many requests (429 response code)

### Error message

When this issue occurs, you receive following error message in the Docker CLI or
in the Docker Engine logs:

```text
Too Many Requests
```

### Possible causes

- You have reached the [Abuse rate limit](./usage/_index.md#abuse-rate-limit)

### Solution

1. Check for broken CI/CD pipelines accessing Docker Hub and fix them.
2. Implement a retry with back-off solution in your automated scripts to ensure
   that you're not resending thousands of requests per minute.

## 500 response code

### Error message

When this issue occurs, the following error message is common in the Docker CLI
or in the Docker Engine logs:

```text
Unexpected status code 500
```

### Possible causes

- There is a temporary Docker Hub service issue.

### Solution

1. View the [Docker System Status Page](https://www.dockerstatus.com/) and
   verify that all services are operational.
2. Try accessing Docker Hub again. It may be a temporary issue.
3. [Contact Docker Support](https://www.docker.com/support/) to report the issue.