---
title: Subscription
description: High-level summary of Docker subscription tiers and optional add-ons for personal accounts, organizations, and enterprise.
keywords:
  - docker subscription
  - subscription plans
  - add-ons
  - docker personal
  - docker pro
  - docker team
  - docker business
  - pricing
  - billing
  - subscription management
weight: 50
params:
  sidebar:
    group: Platform
grid_subscriptions:
  - title: Create an account
    description: Create your first Docker account.
    link: /accounts/
    icon: group_add
  - title: Set up your subscription
    description: Get started setting up a personal or organization subscription.
    link: /subscription/setup/
    icon: shopping_cart
  - title: Scale your subscription
    description: Scale your subscription to fit your needs.
    link: /subscription/scale/
    icon: leaderboard
aliases:
  - /docker-hub/billing/
  - /docker-hub/billing/faq/
  - /subscription/overview/
---

This page summarizes Docker subscription plans and optional add-ons for personal accounts, organization accounts, and sales-led enterprise programs.

> [!TIP]
> For complete tier comparisons and commercial terms, see
> [Docker pricing](https://www.docker.com/pricing?ref=Docs&refAction=DocsSubscription).

## Docker subscription types

Personal accounts serve users developing with Docker for smaller-scale projects. When you upgrade to a Pro account or purchase add-ons, your billing details are tied to an individual Docker ID.

- Docker Personal is the no-cost tier on a personal account.
- Docker Pro is a paid, per-user subscription on the same account type.

Docker also offers organization accounts that group individual members, repositories, and billing under a Docker Hub organization.

- Docker Team is an organization subscription that improves collaboration across small developer teams.
- Docker Business adds security and governance capabilities for larger regulated environments.

Whether you're developing with personal or organization accounts, Docker offers a suite of products that enhances your development workflows.

### Docker Personal and Docker Pro

The table provides an abridged comparison of Docker Personal and Docker Pro accounts.

| Feature                                                        | Personal                    | Pro                                                                                                             |
| -------------------------------------------------------------- | --------------------------- | --------------------------------------------------------------------------------------------------------------- |
| Hub pull rate                                                  | 100 pulls per hour per user | Unlimited                                                                                                       |
| Private repositories                                           | 1                           | Unlimited                                                                                                       |
| Scout-enabled repositories                                     | 1                           | 2                                                                                                               |
| Scout SDLC integrations                                        | None                        | Up to 5                                                                                                         |
| Docker Build Cloud included build minutes per month            | Trial (no paid pool)        | 200                                                                                                             |
| Docker Build Cloud included cache                              | None                        | 50 GB                                                                                                           |
| Docker Build Cloud maximum parallel builds                     | None                        | 4                                                                                                               |
| Docker Build Cloud extra build minutes                         | None                        | Additional Docker Build Cloud build minutes (prepaid, expire when the subscription ends)                        |
| Docker Testcontainers Cloud included runtime minutes per month | Trial                       | 100                                                                                                             |
| Docker Testcontainers Cloud extra runtime minutes              | None                        | Additional Testcontainers Cloud runtime minutes (prepaid through Docker sales or on-demand on the monthly bill) |

To learn more about creating a Docker account, see [Accounts overview](/manuals/accounts/_index.md).

### Docker Team and Docker Business

Docker Team and Docker Business are subscription types that serve organizations. Administrators can oversee identity management, security, and subscriptions across scalable teams.

| Feature                                         | Team                 | Business             |
| ----------------------------------------------- | -------------------- | -------------------- |
| Members                                         | Up to 100            | Unlimited            |
| Docker Hub organizations                        | 1                    | Unlimited            |
| Organization access tokens                      | 10                   | 100                  |
| Build Cloud included minutes per month          | 500                  | 1500                 |
| Build Cloud included cache                      | 100 GB               | 200 GB               |
| Build Cloud capacity add-ons                    | Contact Docker sales | Contact Docker sales |
| Testcontainers Cloud included minutes per month | 500                  | 1500                 |

You can purchase Docker Team and Docker Business subscriptions on a per-organization basis. To learn more about managing organizations through Docker Team or Docker Business, see [Administration overview](/manuals/admin/_index.md).

### Docker Enterprise

Docker Enterprise covers sales-led agreements for large organizations that need tailored commercial terms, deployment programs, and coordinated support beyond standard self-serve subscriptions.

Engage Docker through your procurement process or reach out using [Contact sales](https://www.docker.com/pricing/contact-sales/).

## Scaling your subscription

You can scale your subscriptions with self-serve add-ons:

- For organization accounts, you can purchase hardened repos with DHI Select.
- For personal accounts, you can add Gordon Plus, Max, or Ultra as a subscription layer that multiplies your Gordon usage.

Products available for add-ons can be found by signing into [Docker Home](https://app.docker.com/) from your personal or organization account, then selecting **Billing > Browse Products**.

To add build minutes to your Docker Build Cloud and Testcontainers Cloud subscription, [reach out to sales](https://www.docker.com/pricing/contact-sales/).

## What's next

{{< grid items="grid_subscriptions" >}}
