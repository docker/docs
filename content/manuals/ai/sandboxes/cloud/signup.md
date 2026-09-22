---
title: Sign up for cloud sandboxes
linkTitle: Signup and billing
description: Activate cloud sandbox access, sign in with the sbx CLI, and manage cloud compute usage and billing.
keywords: docker sandboxes, cloud sandboxes, signup, billing, subscription, docker agentic platform, sbx login
weight: 5
---

To run cloud sandboxes, activate a pay-as-you-go subscription for your Docker
account. Signup and billing use Docker Agentic Platform, so you'll see that
name in the web console, checkout, and your active plans.

## Before you begin

You need:

- A Docker Personal or Docker Pro account
- The [`sbx` CLI](../install.md), version 0.45.0 or later

Cloud compute is billed separately from your Docker Personal or Pro plan.
Review the [billing details](#billing) before subscribing.

## Activate cloud access

1. Open the [Docker Agentic Platform console](https://agentic-platform.docker.com/)
   and sign in with the Docker account you want to use for cloud sandboxes.
2. If your account doesn't have access, you're redirected to Docker Billing.
   Review the Docker Agentic Platform pay-as-you-go plan, provide the requested
   billing and payment details, and complete checkout.

   If you see **Docker Agentic Platform access required** instead of a redirect,
   select **Go to Docker Billing** to subscribe.
3. Return to the console after checkout. If your account already has an active
   subscription, you can skip checkout.

You can create your first sandbox from the CLI. You don't need to create one
in the web console first.

## Sign in from the CLI

Sign in with the same Docker account you used to subscribe:

```console
$ sbx login
```

Check cloud connectivity and account access:

```console
$ sbx --cloud diagnose
```

If the check reports that your account doesn't have access, confirm that
checkout completed and that the CLI is signed in to the account with the
subscription.

Then [create your first cloud sandbox](_index.md#get-started). Configure agent
credentials with [`sbx --cloud secret`](credentials.md), even if you previously
configured credentials for local sandboxes or through the web console.

## Billing

Cloud sandboxes use pay-as-you-go compute with no recurring subscription fee.
Docker meters compute in seconds, based on sandbox runtime and the CPU and
memory configuration. Inference charges are separate and billed by your model
provider.

Open [Usage & billing](https://agentic-platform.docker.com/usage) in the console
to review pricing, usage, and available credits. Compute credits offset eligible
cloud compute usage; they don't cover model-provider charges.

Billing is monthly, on the day you subscribed. Your invoice reflects usage
accrued during the previous billing period. Stop or remove sandboxes you no
longer need, and [check their expiration settings](usage.md#configure-expiration).

### Manage your subscription

1. Sign in to [Docker Home](https://app.docker.com/) and select your account.
2. Go to **Billing**, then **Active plans**.
3. Select **Manage** next to **Docker Agentic Platform** to review the plan.

To cancel, select **Cancel subscription**, review your usage, and confirm.
Cancellation disables renewal at the start of the next plan period. Docker
bills the usage you've accrued at the end of the current period.
