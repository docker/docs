---
title: Gordon usage and limits
linkTitle: Usage and limits
description: Understand Gordon's usage limits and how they vary by Docker
  subscription tier
weight: 50
---

{{< summary-bar feature_name="Gordon" >}}

Gordon is available with all Docker subscriptions. Usage limits apply based on
your subscription tier to ensure fair resource allocation.

## Usage limits by subscription

Usage limits increase with higher subscription tiers:

| Subscription | Usage Multiplier | Description                          |
| ------------ | ---------------- | ------------------------------------ |
| Personal     | 1x (baseline)    | Standard usage for personal projects |
| Pro          | 3x               | Three times Personal tier usage      |
| Team         | 3x               | Three times Personal tier usage      |
| Business     | 6x               | Six times Personal tier usage        |

> [!NOTE]
> Limits are per user, not per organization for Team and Business
> subscriptions. Team provides the same multiplier as Pro (3x), but limits
> apply per individual user rather than being shared across the organization.

## How usage is measured

Gordon usage is activity-based. Each interaction consumes resources based on the
complexity of your request and the work Gordon performs. Simple queries use less
than complex multi-step tasks.

## What happens when you reach a limit

As you approach your usage limit, a message appears near the chat input
prompting you to upgrade your subscription.

When you reach your usage limit, Gordon becomes unavailable until the limit
resets. The interface displays when the limit will reset and suggests upgrading
your subscription for higher limits.

## Fair use policy

Usage limits are designed for typical development workflows. Gordon is intended
for:

- Active development and debugging
- Learning Docker concepts
- Optimizing Docker configurations
- Troubleshooting issues

Automated scripting or excessive programmatic access may be subject to
additional restrictions.
