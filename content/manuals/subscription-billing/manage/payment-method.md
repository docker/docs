---

title: Add or update a payment method
linkTitle: Payment methods
weight: 20
description: Learn how to manage cards, US bank accounts, Stripe Link, and pay by invoice for Docker subscriptions.
keywords: payments, billing, subscription, payment methods, credit card, ACH, US bank account, Stripe Link, pay by invoice, failed payments
aliases:

* /billing/payment-method/
* /billing/core-billing/payment-method/

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

* Verify a
  [bank account](#verify-a-bank-account) before choosing it as a payment method.
* Use an existing Stripe Link account, or enter your card information to set up
  Stripe Link.

## Manage payment method

Paid personal accounts and organizations follow the same procedures to add,
update, or remove payment methods.

### Add payment method

You can add multiple payment methods in the billing portal. When you add a new payment method, you set it as a new default payment method.

1. Sign in to [Docker Home](https://app.docker.com/).
2. Select your username for a personal account, or select your organization
   name for an organization.
3. Select **Billing** to go to the billing portal, then **Change** in the **Payment method** tile.
4. From the **Change payment method** modal, choose to add a card, a US bank account, or a Stripe Link payment.

   * To pay with a card, enter your card information.
   * To pay with a US bank account, verify your **Email** and **Full name**.

     * If your bank is listed, select your bank's name.
     * If your bank is not listed, select **Search for your bank**.
   * To pay through Stripe Link, select an existing payment method, then select
     **Use this card**.
5. For first-time setup, enter your billing information.
6. Finish adding the payment method by selecting **Save as default**.

### Change default payment method

After adding one or more payment methods, you can set one as a default method.

1. From **Billing**, go to the **Payment method** tile.
2. Select **Change** to open the **Change payment method** modal, then select **Change** next to your current default method.
3. Choose the payment method you want to set as default.
4. Verify your information, then select **Save as default**.

### Remove payment method

You can only remove secondary payment methods. To remove a secondary payment method:

1. From **Billing**, go to the **Payment method** tile.
2. Select **Change** to open the **Change payment method** modal.
3. Select the **Actions** menu next to the payment method you want to remove, then select **Remove**.
4. Verify your billing details, then select **Save as default**.

To remove your default payment method, first set a different payment method as default, or [downgrade to a free subscription](/manuals/subscription-billing/plans/docker.md#cancel-a-docker-plan).

## Enable and disable pay by invoice

> [!TIP]
> Do you need to pay by invoice? [Upgrade to a Docker Business or Docker Team plan](https://www.docker.com/pricing?ref=Docs&refAction=DocsBillingPaymentMethod) and choose the annual subscription.

Pay by invoice is available for annual Docker Team and Docker Business subscriptions. Docker Sales may need to enable pay by invoice on your account before it appears in Docker Home. For the first subscription period, you pay using a payment card or ACH bank transfer. At renewal time, instead of automatic payment, you'll receive an invoice by email that you must pay manually.

Follow these steps to enable or disable pay by invoice:

1. Sign in to [Docker Home](https://app.docker.com/) and select your organization.
2. Select **Billing**, then **Payment methods**.
3. Select **Pay by invoice**, then use the pay by invoice toggle to enable or disable.
4. Confirm your billing contact details. If you need to change them, select
   **Change** and enter your new details.

If the pay by invoice option isn't available in Docker Home, contact your Docker sales representative.

Pay by invoice isn't available for subscription upgrades or changes.

## Verify a bank account

You can verify your bank account with instant verification for supported banks. You must sign in to your US bank account when adding your bank as a payment method:

1. From **Billing**, go to the **Payment method** tile and select **Change**.
2. Choose **US bank account** as your payment method.
3. Verify your **Email** and **Full name**.
4. Search for your bank, then select it to initiate the sign-in flow.
5. Review the terms and conditions. This agreement
   allows Docker to debit payments from your connected bank account.
6. Select **Agree and continue**.
7. Select an account to link and verify, then select **Connect account**.

When the account is verified, a success message appears.

## Failed payments

If your payment fails, select **Pay now** to retry the payment through Stripe.

You have a grace period of 15 days, including the due date, when your payment
fails. Docker attempts to collect the payment three times using the following
schedule:

* 3 days after the due date
* 5 days after the previous attempt
* 7 days after the previous attempt

Docker also sends an email notification
`Action Required - Credit Card Payment Failed` with an attached unpaid invoice
after each failed payment attempt.

If the invoice remains unpaid after the grace period, the
subscription downgrades to a free subscription and all paid features are
disabled.

