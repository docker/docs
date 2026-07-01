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

Docker AI Cloud brings the Agentic Platform, [MCP Gateway](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md), and [Cloud Sandboxes](/manuals/ai/sandboxes/_index.md) together in a single product for individual developers.

- AI Cloud Basic is the free default for every new signup. It doesn't require a payment method.
- AI Cloud Personal is a pay-as-you-go plan for developers who need Cloud Sandboxes. It has no recurring subscription fee, so you pay only for the usage you accrue.

To upgrade to AI Cloud Personal, see [Manage plans](/manuals/subscription/manage.md).

## Usage

AI Cloud Basic includes MCP Gateway access with a monthly tool-call allowance. When you upgrade to AI Cloud Personal, you receive Cloud Sandboxes, higher MCP Gateway limits, concurrent sandboxes, snapshot storage, and private Docker Hub repositories.

When you subscribe to an AI Cloud plan, your invoice reflects metered, pay-as-you-go usage. Docker meters compute usage based on the vCPU and memory that your sandboxes consume over time, along with outbound data transfer.

When you use AI Cloud Basic or Personal, you bring your own API keys for inference. Those costs are handled by your inference provider.

> [!NOTE]
> Docker AI Cloud plans are a Docker Personal 
> and Docker Pro exclusive.

## Billing behaviors

AI Cloud Personal is pay-as-you-go. Docker meters your cloud sandbox compute and outbound data transfer as you use them, then bills the accrued usage monthly on the first of the month. AI Cloud Basic requires no payment method and incurs no charges.

When you upgrade from AI Cloud Basic to AI Cloud Personal, you add a payment method to enable Cloud Sandboxes and higher limits.

AI Cloud Personal signups receive a one-time promotional credit toward cloud compute usage. To review your balance, see [View your credits](/manuals/subscription/manage.md#view-your-credits).

## Cancel a plan

When you cancel AI Cloud Personal, access ends immediately and Docker charges your accrued metered usage right away rather than waiting for the next billing date. Canceling also terminates any running background processes, including active sandboxes. Your account returns to AI Cloud Basic.
