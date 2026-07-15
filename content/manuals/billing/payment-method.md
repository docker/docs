---
title: Add or update a payment method
linkTitle: Payment methods
weight: 20
description: Learn how to manage cards, US bank accounts, Stripe Link, and pay by invoice for Docker subscriptions.
keywords: payments, billing, subscription, payment methods, credit card, ACH, US bank account, Stripe Link, pay by invoice, failed payments
aliases:
  - /billing/core-billing/payment-method/
---

Docker supports several payment methods for paid personal accounts and
organizations. This page describes supported payment types, how to manage
payments from [Docker Home](https://app.docker.com/), and how to set up pay by
invoice.

## Supported payment types

You can add or update a payment method at any time. All charges are in United
States dollars (USD). Docker supports the following payment methods:

| Category      | Payment type                                                            |
| ------------- | ----------------------------------------------------------------------- |
| Cards         | Visa, MasterCard, American Express, Discover, JCB, Diners, UnionPay     |
| Wallets       | Stripe Link                                                             |
| Bank accounts | Automated Clearing House (ACH) transfer with a verified US bank account |

## Prerequisites

Some payment methods require additional setup:

- Verify a
  [bank account](#verify-a-bank-account) before choosing it as a payment method.
- Use a Docker Business or Docker Team plan to
  [pay by invoice](#enable-and-disable-pay-by-invoice).
- Use an existing Stripe Link account, or enter your card information to set up
  Link.

## Manage payment method

Paid personal accounts and organizations follow the same procedures to add,
update, or remove payment methods.

### Add payment method

1. Sign in to [Docker Home](https://app.docker.com/).
1. Select your account name for a personal account, or select your organization
   name for an organization.
1. Select **Billing**, then **Payment methods**.
1. Select **Add payment method** and enter your new payment information:
   - For first-time setup, enter your billing information.
   - To purchase as a business, provide your tax ID.
1. Choose to add a card, a US bank account, or a Link payment.
   - To pay with a card, enter your card information.
   - To pay with a US bank account:
     - Verify your **Email** and **Full name**.
     - If your bank is listed, select your bank's name.
     - If your bank is not listed, select **Search for your bank**.
   - To pay through Link, select an existing payment method, then select
     **Use this card**.
1. Finish adding the payment method by selecting **Add payment method**.

### Set default payment method

After adding one or more payment methods, you can set one as a default method.

1. From **Billing**, go to **Payment methods**.
1. In the **Payment method** table, find the payment method you want to set as
   default.
1. Select the three dots, then choose **Set as default**.

### Remove payment method

To remove a single payment method:

1. From **Billing**, go to **Payment methods**.
1. In the **Payment method** table, find the payment method you want to remove.
1. Select the three dots, then choose **Remove**.

To remove your default payment method, first set a different default or
[downgrade to a free subscription](/manuals/subscription/plans/docker.md#cancel-a-docker-plan).

## Enable and disable pay by invoice

> [!TIP]
>
> To pay by invoice,
> [upgrade to a Docker Business or Docker Team plan](https://www.docker.com/pricing?ref=Docs&refAction=DocsBillingPaymentMethod)
> and choose an annual subscription.

Pay by invoice requires upfront payment for your first subscription period
using a payment card or ACH bank transfer. At renewal, Docker emails you an
invoice to pay manually instead of charging your default payment method.

Follow these steps to enable or disable pay by invoice:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Billing**, then **Payment methods**.
1. Select **Pay by invoice**, then turn the pay by invoice toggle on or off.
1. Confirm your billing contact details. To change them, select
   **Change** and enter your new details.

Pay by invoice isn't available for subscription upgrades or changes.

## Verify a bank account

There are two ways to verify a bank account as a payment method:

- Instant verification for supported banks
- Manual verification using micro-deposits

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

When the account is verified, a success message appears.

### Manual verification

To verify your bank account manually, you must enter the micro-deposit amount
from your bank statement:

1. Choose **US bank account** as your payment method.
1. Verify your **Email** and **First and last name**.
1. Select **Enter bank details manually instead**.
1. Enter your bank details: **Routing number** and **Account number**.
1. Select **Submit**.
1. Follow the instructions in the verification email.

Manual verification uses micro-deposits. You'll see a small deposit, such as
$0.01, in your bank account within one to two business days. Open your manual
verification email and enter the amount of this deposit to verify your account.

## Failed payments

If your payment fails, select **Pay now** to retry the payment through Stripe.

You have a grace period of 15 days, including the due date, when your payment
fails. Docker attempts to collect the payment three times using the following
schedule:

- 3 days after the due date
- 5 days after the previous attempt
- 7 days after the previous attempt

Docker also sends an email notification
`Action Required - Credit Card Payment Failed` with an attached unpaid invoice
after each failed payment attempt.

If the invoice remains unpaid after the grace period, the
subscription downgrades to a free subscription and all paid features are
disabled.
