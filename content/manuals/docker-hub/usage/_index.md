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

When using Docker Hub, unauthenticated and Docker Personal users are subject to
strict limits. In contrast, Docker Pro, Team, and Business users benefit from a
consumption-based model with a base amount of included usage. This included
usage is not a hard limit; users can scale or upgrade their subscriptions to
receive additional usage or use on-demand usage.

The following table provides an overview of the included usage and limits for each
user type, subject to fair use:


| User type                | Pulls per month | Pull rate limit per hour               | Public repositories | Public repository storage | Private repositories | Private repository storage |
|--------------------------|-----------------|----------------------------------------|---------------------|---------------------------|----------------------|----------------------------|
| Business (authenticated) | 1M              | Unlimited                              | Unlimited           | Unlimited                 | Unlimited            | Up to 500 GB               |
| Team (authenticated)     | 100K            | Unlimited                              | Unlimited           | Unlimited                 | Unlimited            | Up to 50 GB                |
| Pro (authenticated)      | 25K             | Unlimited                              | Unlimited           | Unlimited                 | Unlimited            | Up to 5 GB                 |
| Personal (authenticated) | Not applicable  | 40                                     | Unlimited           | Unlimited                 | Up to 1              | Up to 2 GB                 |
| Unauthenticated users    | Not applicable  | 10 per IPv4 address or IPv6 /64 subnet | Not applicable      | Not applicable            | Not applicable       | Not applicable             |

For more details, see the following:

- [Pull usage and limits](./pulls.md)
- [Storage usage and limits](./storage.md)

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
