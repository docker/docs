---
title: Subscription overview
linkTitle: Subscription
description: Learn how Docker subscriptions work, including the types of plans available and how they apply to personal and organization accounts.
keywords:
  docker subscription, pricing, billing, subscription types, subscription
  plans, docker hardened images, gordon, cloud sandboxes, subscription
  management
weight: 30
params:
  sidebar:
    group: Platform
grid_subscriptions:
  - title: Compare Docker subscriptions
    description: Visit the pricing page to see what's included in different Docker subscriptions.
    link: "https://www.docker.com/pricing?ref=Docs&refAction=DocsSubscription"
    icon: magnifying-glass
  - title: Set up your subscription
    description: Get started setting up a personal or organization subscription.
    link: /subscription/setup/
    icon: shopping-cart
  - title: Scale your subscription
    description: Scale your subscription to fit your needs.
    link: /subscription/scale/
    icon: chart-bar
  - title: Change your subscription
    description: Learn how to upgrade or downgrade your subscription.
    link: /subscription/change/
    icon: arrow-up-circle
  - title: Docker Desktop license agreement
    description: Review the terms of the Docker Subscription Service Agreement.
    link: /subscription/desktop-license/
    icon: document-text
  - title: Subscription FAQs
    description: Find the answers you need and explore common questions.
    link: /subscription/faq/
    icon: question-mark-circle
aliases:
  - /subscription/scale/
  - /subscription/change/
  - /subscription/details/
  - /docker-hub/billing/
  - /docker-hub/billing/faq/
---

After creating a Docker account, you can subscribe to a number of Docker plans. Docker plans may include different tier levels, ranging from basic tiers to a selection of paid tiers. Each tier within a plan upgrades your usage entitlements and feature sets.

This page breaks down the kinds of Docker plans you can subscribe to, with explanations about usage within each tier.

## Docker plan types

There are two types of Docker plans. Docker Core upgrades your basic personal or organization accounts to higher tiers with additional feature sets. Product-based plans are subscription types tied to discrete products in Docker's product catalog.

| Subscription type | Billing model                                    | Examples                                                                        |
| ----------------- | ------------------------------------------------ | ------------------------------------------------------------------------------- |
| Docker Core       | Flat-rate tiers tied to account type             | Docker Pro (Personal)<br>Docker Team and Docker Business (Organization)         |
| Product-based     | Prepaid entitlements added on top of a base plan | Gordon (Extends monthly usage limits)<br>DHI (Security and compliance features) |

Subscription types can be combined. A Docker Core plan provides the foundation for
most accounts, with prepaid or per-unit subscriptions added on top as needed. Depending on the product, you may not need a Docker Core plan at all.

## Top up your plan

Most subscriptions include a fixed amount of usage. You can purchase additional
units to extend usage without changing your plan tier.

| Unit         | Description                                                                             | Examples           |
| ------------ | --------------------------------------------------------------------------------------- | ------------------ |
| Seats        | Each seat extends the subscription entitlements to one more member.                     | Docker Core        |
| Licenses     | Access to specific products or product tiers, purchased separately from your core plan. | AI Governance      |
| Minutes      | Cloud build capacity, sold in blocks and consumed within the billing period.            | Docker Build Cloud |
| Repositories | Additional container repositories covered by security scanning and analysis features.   | DHI                |

## Support add-ons

Standard support is included with paid plans and scales with plan tier.

- Docker also offers premium support as an optional add-on for Docker Business and DHI customers.
- Premium support provides 24/7 response, priority SLAs, and a dedicated technical advisory manager.

## What's next

{{< grid items="grid_subscriptions" >}}
