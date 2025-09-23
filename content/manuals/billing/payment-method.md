---
title: Add or update a payment method
weight: 20
description: Learn how to add or update a payment method in Docker Hub
keywords: payments, billing, subscription, supported payment methods, failed payments, add credit card, bank transfer, Stripe Link, payment failure
aliases:
    - /billing/core-billing/payment-method/
---

This page describes how to add or update a payment method for your personal
account or for an organization.

You can add a payment method or update your account's existing payment method
at any time.

> [!IMPORTANT]
>
> If you want to remove all payment methods, you must first downgrade your
subscription to a free subscription. See [Downgrade](../subscription/change.md).

The following payment methods are supported:

- Cards
  - Visa
  - MasterCard
  - American Express
  - Discover
  - JCB
  - Diners
  - UnionPay
- Wallets
  - Stripe Link
- Bank accounts
  - Automated Clearing House (ACH) transfer with a
  [verified](manuals/billing/payment-method.md#verify-a-bank-account) US
  bank account
- [Pay by invoice](/manuals/billing/history.md)

All charges are in United States dollars (USD).

{{% include "tax-compliance.md" %}}

## Manage payment method

### Personal account

{{< tabs >}}
{{< tab name="Docker subscription" >}}

To add a payment method:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Billing**.
1. Select **Payment methods** from the left-hand menu.
1. Select **Add payment method**.
1. Enter your new payment information:
    - Add a card:
        - Select **Card** and fill out the card information form.
    - Add a Link payment:
        - Select **Secure, 1-click checkout with Link** and enter your
        Link **email address** and **phone number**.
        - If you don't already use Link, you must fill out the card information
        form to store a card for Link payments.
    - Add a bank account:
        - Select **US bank account**.
        - Verify your **Email** and **Full name**.
        - If your bank is listed, select your bank's name.
        - If your bank is not listed, select **Search for your bank**.
        - To verify your bank account, see
        [Verify a bank account](manuals/billing/payment-method.md#verify-a-bank-account).
1. Select **Add payment method**.
1. Optional. You can set a new default payment method by selecting
the **Set as default** action.
1. Optional. You can remove non-default payment methods by selecting
the **Delete** action.

> [!NOTE]
>
> If you want to set a US bank account as your default payment method, you must
> [verify the account](#verify-a-bank-account) first.

{{< /tab >}}
{{< tab name="Legacy Docker subscription" >}}

To add a payment method:

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **Billing**.
1. Select the **Payment methods** link.
1. Select **Add payment method**.
1. Enter your new payment information:
    - Add a card:
        - Select **Card** and fill out the card information form.
    - Add a Link payment:
        - Select **Secure, 1-click checkout with Link** and enter your
        Link **email address** and **phone number**.
        - If you are not an existing Link customer, you must fill out the
        card information form to store a card for Link payments.
1. Select **Add**.
1. Select the **Actions** icon, then select **Make default** to ensure that
your new payment method applies to all purchases and subscriptions.
1. Optional. You can remove non-default payment methods by selecting
the **Actions** icon. Then, select **Delete**.

{{< /tab >}}
{{< /tabs >}}

### Organization

> [!NOTE]
>
> You must be an organization owner to make changes to the payment information.

{{< tabs >}}
{{< tab name="Docker subscription" >}}

To add a payment method:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Billing**.
1. Select **Payment methods** from the left-hand menu.
1. Select **Add payment method**.
1. Enter your new payment information:
    - Add a card:
        - Select **Card** and fill out the card information form.
    - Add a Link payment:
        - Select **Secure, 1-click checkout with Link** and enter your
        Link **email address** and **phone number**.
        - If you are not an existing Link customer, you must fill out the
        card information form to store a card for Link payments.
    - Add a bank account:
        - Select **US bank account**.
        - Verify your **Email** and **Full name**.
        - If your bank is listed, select your bank's name.
        - If your bank is not listed, select **Search for your bank**.
        - To verify your bank account, see [Verify a bank account](manuals/billing/payment-method.md#verify-a-bank-account).
1. Select **Add payment method**.
1. Optional. You can set a new default payment method by selecting
the **Set as default** action.
1. Optional. You can remove non-default payment methods by selecting
the **Delete** action.

> [!NOTE]
>
> If you want to set a US bank account as your default payment method, you must
> verify the account first.

{{< /tab >}}
{{< tab name="Legacy Docker subscription" >}}

To add a payment method:

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select your organization, then select **Billing**.
1. Select the **Payment methods** link.
1. Select **Add payment method**.
1. Enter your new payment information:
    - Add a card:
        - Select **Card** and fill out the card information form.
    - Add a Link payment:
        - Select **Secure, 1-click checkout with Link** and enter your
        Link **email address** and **phone number**.
        - If you are not an existing Link customer, you must fill out the
        card information form to store a card for Link payments.
1. Select **Add payment method**.
1. Select the **Actions** icon, then select **Make default** to ensure that
your new payment method applies to all purchases and subscriptions.
1. Optional. You can remove non-default payment methods by selecting
the **Actions** icon. Then, select **Delete**.

{{< /tab >}}
{{< /tabs >}}

## Enable pay by invoice

Pay by invoice is available for Teams and Business customers with 25+ seats
on annual subscriptions, starting with your first renewal. When you select this
payment method, you'll pay upfront for your first subscription period using a
payment card or ACH bank transfer.

At renewal time, instead of automatic payment, you'll receive an invoice via
email that you must pay manually. Pay by invoice is not available for
subscription upgrades or changes.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Billing**.
1. Select **Payment methods**, then **Pay by invoice**.
1. To enable pay by invoice, select the toggle.
1. Confirm your billing contact details. If you need to change them, select
**Change** and enter your new details.

## Verify a bank account

There are two ways to verify a bank account as a payment method:

- Instant verification: Docker supports several major banks for instant
verification.
- Manual verification: All other banks must be verified manually.

{{< tabs >}}
{{< tab name="Instant verification" >}}

### Instant verification

To verify your bank account instantly, you must sign in to your bank account
from the Docker billing flow:

1. Choose **US bank account** as your payment method.
1. Verify your **Email** and **Full name**.
1. If your bank is listed, select your bank's name or
select **Search for your bank**.
1. Sign in to your bank and review the terms and conditions. This agreement
allows Docker to debit payments from your connected bank account.
1. Select **Agree and continue**.
1. Select an account to link and verify, and select **Connect account**.

When the account is verified, you will see a success message in the pop-up
modal.

{{< /tab >}}
{{< tab name="Manual verification" >}}

### Manual verification

To verify your bank account manually, you must enter the micro-deposit amount
from your bank statement:

1. Choose **US bank account** as your payment method.
1. Verify your **Email** and **First and last name**.
1. Select **Enter bank details manually instead**.
1. Enter your bank details: **Routing number** and **Account number**.
1. Select **Submit**.
1. You will receive an email with instructions on how to manually verify.

Manual verification uses micro-deposits. You’ll see a small deposit
(such as $0.01) in your bank account within 1–2 business days. Open your manual
verification email and enter the amount of this deposit to verify your account.

{{< /tab >}}
{{< /tabs >}}

## Failed payments

> [!NOTE]
>
> You can't manually retry a failed payment. Docker will retry failed payments
based on the retry schedule.

If your subscription payment fails, there is a grace period of 15 days,
including the due date. Docker retries to collect the payment 3 times using the
following schedule:

- 3 days after the due date
- 5 days after the previous attempt
- 7 days after the previous attempt

Docker also sends an email notification
`Action Required - Credit Card Payment Failed` with an attached unpaid invoice
after each failed payment attempt.

Once the grace period is over and the invoice is still not paid, the
subscription downgrades to a free subscription and all paid features are
disabled.
