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

Some payment methods require additional setup before changing it to your default payment method. For example:

- Verify a
  [bank account](#verify-a-bank-account) before choosing it as a payment method.
- Use an existing Stripe Link account, or enter your card information to set up
  Link.

## Manage payment method

Paid personal accounts and organizations follow the same procedures to add,
update, or remove payment methods.

### Add payment method

You can add multiple payment methods in the billing portal. When you add a new payment method, you set it as a new default payment method. 

1. Sign in to [Docker Home](https://app.docker.com/).
1. Select your account name for a personal account, or select your organization
   name for an organization.
1. Select **Billing** to go to the billing portal, then **Change** in the **Payment methods** tile.
1. From the **Change payment method** modal, choose to add a card, a US bank account, or a Link payment.
   - To pay with a card, enter your card information.
   - To pay with a US bank account, verify your **Email** and **Full name**.
     - If your bank is listed, select your bank's name.
     - If your bank is not listed, select **Search for your bank**.
   - To pay through Link, select an existing payment method, then select
     **Use this card**.
1. For first-time setup, enter your billing information. 
1. Finish adding the payment method by selecting **Save as default**.

### Change default payment method

After adding one or more payment methods, you can set one as a default method.

1. From **Billing**, go to the **Payment methods** tile.
1. Select **Change** to open the **Change payment method** modal, then select **Change** next to your current default method. 
1. Choose the payment method you want to set as default.
1. Verify your information, then select **Save as default**.

### Remove payment method

You can only remove secondary payment methods. To remove a secondary payment method:

1. From **Billing**, go to the **Payment methods** tile.
1. Select **Change** to open the **Change payment method** modal, then select **Change** next to your current default method.
1. Select the **Actions** menu next to the payment method you want to remove, then select **Remove**.
1. Verify your billing details, then select **Save as default**. 

You can't remove a default payment method. If you want to remove your default payment method, you must change your default payment method then follow the remove payment method procedures.

## Pay by invoice

> [!TIP]
>
> To pay by invoice,
> [upgrade to a Docker Business or Docker Team plan](https://www.docker.com/pricing?ref=Docs&refAction=DocsBillingPaymentMethod)
> and choose an annual subscription.

Pay by invoice requires upfront payment for your first subscription period
using a payment card or ACH bank transfer. At renewal, Docker emails you an
invoice to pay manually instead of charging your default payment method.

- To add pay by invoice as a payment method, contact your Docker sales representative.
- You can only pay by invoice by choosing it as a payment method when subscribing to Docker Team or Docker Business plans.  
- Pay by invoice isn't available for subscription upgrades or changes.

## Verify a bank account

You can verify your bank account with instant verification for supported banks. You must sign in to your US bank account when adding your bank as a payment method:

1. From **Billing**, go to the **Payment method** tile and select **Change**.
1. Choose **US bank account** as your payment method.
1. Verify your **Email** and **Full name**.
1. Search for your bank, then select it to initiate the signin flow.
1. Review the terms and conditions. This agreement
   allows Docker to debit payments from your connected bank account.
1. Select **Agree and continue**.
1. Select an account to link and verify, then select **Connect account**.

When the account is verified, a success message appears.

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
