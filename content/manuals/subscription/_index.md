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
  - title: Compare Docker subscriptions
    description: Visit the pricing page to see what's included in different Docker subscriptions.
    link: "https://www.docker.com/pricing?ref=Docs&refAction=DocsSubscription"
    icon: feature_search
  - title: Set up your subscription
    description: Get started setting up a personal or organization subscription.
    link: /subscription/setup/
    icon: shopping_cart
  - title: Scale your subscription
    description: Scale your subscription to fit your needs.
    link: /subscription/scale/
    icon: leaderboard
  - title: Change your subscription
    description: Learn how to upgrade or downgrade your subscription.
    link: /subscription/change/
    icon: upgrade
  - title: Manage seats
    description: Learn how to add or remove seats from your subscription.
    link: /subscription/manage-seats/
    icon: group_add
  - title: Docker Desktop license agreement
    description: Review the terms of the Docker Subscription Service Agreement.
    link: /subscription/desktop-license/
    icon: license
  - title: Subscription FAQs
    description: Find the answers you need and explore common questions.
    link: /subscription/faq/
    icon: help
aliases:
  - /docker-hub/billing/
  - /docker-hub/billing/faq/
  - /subscription/overview/
---

This page summarizes Docker subscription plans and optional add-ons for personal accounts, organization accounts, and sales-led enterprise programs. For complete tier comparisons and commercial terms, see [Docker pricing](https://www.docker.com/pricing?ref=Docs&refAction=DocsSubscription).

## Docker personal accounts

Personal accounts suit single-user development workflows. Billing and upgrades for these subscriptions are tied to an individual Docker ID.

- Docker Personal is the no-cost tier on a personal account.
- Docker Pro is a paid, per-user subscription on the same account type.

### Docker Personal and Docker Pro

| Feature                                                        | Personal                       | Pro                                                                                                             |
| -------------------------------------------------------------- | ------------------------------ | --------------------------------------------------------------------------------------------------------------- |
| Hub pull rate                                                  | 200 pulls per 6 hours per user | Unlimited                                                                                                       |
| Private repositories                                           | 1                              | Unlimited                                                                                                       |
| Scout-enabled repositories                                     | 1                              | 2                                                                                                               |
| Scout SDLC integrations                                        | None                           | Up to 5                                                                                                         |
| Docker Build Cloud included build minutes per month            | Trial (no paid pool)           | 200                                                                                                             |
| Docker Build Cloud included cache                              | None                           | 50 GB                                                                                                           |
| Docker Build Cloud maximum parallel builds                     | None                           | 4                                                                                                               |
| Docker Build Cloud extra build minutes                         | None                           | Additional Docker Build Cloud build minutes (prepaid, expire when the subscription ends)                        |
| Docker Testcontainers Cloud included runtime minutes per month | Trial                          | 100                                                                                                             |
| Docker Testcontainers Cloud extra runtime minutes              | None                           | Additional Testcontainers Cloud runtime minutes (prepaid through Docker sales or on-demand on the monthly bill) |

## Docker organization accounts

Organization accounts group members, repositories, and billing under a Docker Hub organization.

- Docker Team is an organization subscription that improves collaboration across small developer teams.
- Docker Business adds security and governance capabilities for larger regulated environments.

| Feature                                         | Team                 | Business             |
| ----------------------------------------------- | -------------------- | -------------------- |
| Members                                         | Up to 100            | Unlimited            |
| Docker Hub organizations                        | 1                    | Unlimited            |
| Organization access tokens                      | 10                   | 100                  |
| Build Cloud included minutes per month          | 500                  | 1500                 |
| Build Cloud included cache                      | 100 GB               | 200 GB               |
| Build Cloud capacity add-ons                    | Contact Docker sales | Contact Docker sales |
| Testcontainers Cloud included minutes per month | 500                  | 1500                 |

You can purchase Docker Team and Docker Business subscriptions on a per-organization basis.

## Scale your subscription

- Docker Build Cloud: Purchase prepaid build minutes
- Testcontainers Cloud: Use on-demand runtime minutes with billing at the end of each monthly cycle

[Contact sales](https://www.docker.com/pricing/contact-sales/) to add prepaid Testcontainers Cloud runtime minutes or Docker Build Cloud capacity add-ons to your base subscriptions.

## Docker Enterprise

Docker Enterprise covers sales-led agreements for large organizations that need tailored commercial terms, deployment programs, and coordinated support beyond standard self-serve subscriptions.

Engage Docker through your procurement process or reach out using [Contact sales](https://www.docker.com/pricing/contact-sales/).

## Manage your subscription

Use these guides to compare plans, set up billing, adjust seats, and find answers to common questions.

{{< grid items="grid_subscriptions" >}}
