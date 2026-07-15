---
title: Docker AI Cloud plans
linkTitle: AI Cloud
description:
  Learn about Docker AI Cloud plans for individual developers, including usage
  entitlements, pay-as-you-go billing behaviors, and cancellation
keywords:
  ai cloud, cloud sandboxes, mcp gateway, agentic platform, pay-as-you-go,
  paygo, metered billing, personal subscription, promotional credit
weight: 20
---

> [!TIP]
> AI Cloud Personal sign-ups receive a one-time
> promotional credit toward cloud compute usage.

Docker AI Cloud plans are a Docker Personal and Docker Pro exclusive. Docker AI
Cloud combines the Agentic Platform, MCP Gateway, and Cloud Sandboxes together
in a single product for individual developers.

- AI Cloud Basic is the free default for every new sign-up. It doesn't require
  a payment method.
- AI Cloud Personal is a pay-as-you-go plan for developers who need Cloud
  Sandboxes. It has no recurring subscription fee, so you pay only for the
  usage you accrue.

To add an AI Cloud Personal plan to your Docker Personal or Docker Pro
account, see [Manage plans](/manuals/subscription/manage.md).

## Usage

AI Cloud Basic includes MCP Gateway access with a monthly tool-call
allowance. When you upgrade from AI Cloud Basic to AI Cloud Personal, you
receive:

- Cloud Sandboxes
- Higher MCP Gateway limits
- Concurrent sandboxes and snapshot storage
- Private Docker Hub repositories

When you subscribe to an AI Cloud plan, your invoice reflects metered,
pay-as-you-go usage. Docker meters compute usage based on the vCPU and memory
that your sandboxes consume over time, along with outbound data transfer.
When you use AI Cloud Basic or Personal, you bring your own API keys for
inference. Those costs are handled by your inference provider.

## Billing behaviors

AI Cloud Basic requires no payment method and incurs no charges. When you
upgrade from AI Cloud Basic to AI Cloud Personal, you add a payment method to
enable Cloud Sandboxes with higher limits. AI Cloud Personal bills your
accrued usage monthly on the first of the month.

## Cancel a plan

When you cancel AI Cloud Personal, access ends immediately and Docker charges
your accrued pay-as-you-go usage right away rather than waiting for your next
billing date. Canceling also terminates any running background processes,
including active sandboxes. Your account returns to AI Cloud Basic.
