---
description: Learn about usage and limits for Docker Hub.
keywords: Docker Hub, limit, usage
title: Docker Hub usage and limits
linkTitle: Usage and limits
weight: 30
aliases:
  /docker-hub/download-rate-limit/
---

{{% include "hub-limits.md" %}}

The following table provides an overview of the included usage and limits for each
user type, subject to fair use:


| User type                | Pull rate limit per hour               | Number of public repositories | Number of private repositories |
|--------------------------|----------------------------------------|---------------------|----------------------|
| Business (authenticated) | Unlimited                              | Unlimited           | Unlimited            |
| Team (authenticated)     | Unlimited                              | Unlimited           | Unlimited            |
| Pro (authenticated)      | Unlimited                              | Unlimited           | Unlimited            |
| Personal (authenticated) | 100                                    | Unlimited           | Up to 1              |
| Unauthenticated users    | 10 per IPv4 address or IPv6 /64 subnet | Not applicable      | Not applicable       |

For more details, see [Pull usage and limits](./pulls.md).

## Fair use

When utilizing the Docker Platform, users should be aware that excessive data
transfer, pull rates, or data storage can lead to throttling, or additional
charges. To ensure fair resource usage and maintain service quality, we reserve
the right to impose restrictions or apply additional charges to accounts
exhibiting excessive data and storage consumption.

### Abuse rate limit

Docker Hub has an abuse rate limit to protect the application and
infrastructure. This limit applies to all requests to Hub properties including
web pages, APIs, and image pulls. The limit is applied per IPv4 address or per
IPv6 /64 subnet, and while the limit changes over time depending on load and
other factors, it's in the order of thousands of requests per minute. The abuse
limit applies to all users equally regardless of account level.

You can differentiate between the pull rate limit and abuse rate limit by
looking at the error code. The abuse limit returns a simple `429 Too Many
Requests` response. The pull limit returns a longer error message that includes
a link to documentation.
