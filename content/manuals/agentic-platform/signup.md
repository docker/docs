---
title: Sign up for Docker Agentic Platform
linkTitle: Sign up
description: Activate Docker Agentic Platform access for the Console, cloud sandbox CLI, and API or SDKs.
keywords: docker agentic platform, cloud sandboxes, signup, billing, subscription, account access
weight: 5
---

> [!NOTE]
> Docker Agentic Platform is experimental. Features and behavior may change.

To use Docker Agentic Platform or run cloud sandboxes through the CLI, API, or
SDKs, activate a Docker Agentic Platform pay-as-you-go subscription. The same
subscription provides cloud access for all these interfaces.

## Before you begin

Subscribe using your personal Docker account. You can use this account even
if you belong to an organization. The Docker Agentic Platform subscription
is attached to your personal account and billed separately from your Docker
subscription. Review the [billing details](#billing) before subscribing.

## Activate cloud access

1. Open the [Docker Agentic Platform Console](https://agentic-platform.docker.com/)
   and sign in with the Docker account you want to use for cloud sandboxes.
2. If your account doesn't have access, follow the redirect to Docker Billing.
   Review the Docker Agentic Platform pay-as-you-go plan, provide the requested
   billing and payment details, and complete checkout.

   If you see **Docker Agentic Platform access required** instead of a redirect,
   select **Go to Docker Billing** to subscribe.
3. Return to the Console after checkout. If your account already has an active
   subscription, you can skip checkout.

To create a sandbox in the Console, follow [Get started](get-started.md).
Use the same Docker account for CLI sign-in or API authentication that you
used to subscribe.
Agent credentials are configured separately from subscription activation;
follow your interface's instructions to set them up.

## Check account access

If the Console, CLI, or API reports that your account doesn't have access,
confirm that checkout completed and that you're using the account with the
active subscription. You can review the plan in Docker Home under
**Billing** > **Active subscriptions** > **Docker Agentic Platform**.

## Billing

Cloud sandboxes use pay-as-you-go compute with no recurring subscription fee.
Docker meters compute in seconds, based on sandbox runtime and the CPU and
memory configuration. Your model provider bills inference separately.

Open [Usage & billing](https://agentic-platform.docker.com/usage) in the Console
to review pricing, usage, and available credits. Compute credits offset
eligible cloud compute usage; they don't cover model-provider charges.

Billing is monthly, on the day you subscribed. Your invoice reflects usage
accrued during the previous billing period. Stop or delete sandboxes you no
longer need.

For instructions on reviewing your plan, billing dates, and cancellation, see
[Docker Agentic Platform plans](/manuals/subscription-billing/plans/docker-agentic-platform.md).
