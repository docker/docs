---
title: Add or update a payment method
weight: 20
description: Learn how to add or update a payment method in Docker Hub
keywords: payments, billing, subscription, supported payment methods, failed payments, coupons
alisases:
    - /billing/core-billing/payment-method/
---

This page describes how to add or update a payment method for your personal
account or for an organization.

You can add a payment method or update your account's existing payment method
at any time.

> [!IMPORTANT]
>
> If you want to remove all payment methods, you must first downgrade your
subscription to a free plan. See [Downgrade](../subscription/change.md).

The following payment methods are supported:

- Visa
- MasterCard
- American Express
- Discover
- JCB
- Diners
- UnionPay
- Link
- Pay by invoice

All currency, for example the amount listed on your billing invoice, is in
United States dollar (USD).

{{% include "tax-compliance.md" %}}

## Manage payment method

### Personal account

{{< tabs >}}
{{< tab name="Docker plan" >}}

To add a payment method:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Under Settings and administration, select **Billing**.
3. Select **Payment methods** from the left-hand menu.
4. Select **Add payment method**.
5. Enter your new payment information:
    - If you are adding a card, fill out the card information form.
    - If you are adding a Link payment, select **Secure, 1-click checkout with Link**
    and enter your Link **email address** and **phone number**. If you are not
    an existing Link customer, you must fill out the card information form to
    store a card for Link payments.
6. Select **Add payment method**.
7. Optional. You can set a new default payment method by selecting the **Set as default** action.
8. Optional. You can remove non-default payment methods by selecting the **Delete** action.

{{< /tab >}}
{{< tab name="Legacy Docker plan" >}}

To add a payment method:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select your avatar in the top-right corner.
3. From the drop-down menu select **Billing**.
4. Select the **Payment methods and billing history** link.
5. In the **Payment method** section, select **Add payment method**.
6. Enter your new payment information:
    - If you are adding a card, fill out the card information form.
    - IIf you are adding a Link payment, select **Secure, 1-click checkout with Link**
    and enter your Link **email address** and **phone number**. If you are not
    an existing Link customer, you must fill out the card information form to
    store a card for Link payments.
7. Select **Add**.
8. Select the **Actions** icon, then select **Make default** to ensure that your new payment method applies to all purchases and subscriptions.
9. Optional. You can remove non-default payment methods by selecting the **Actions** icon. Then, select **Delete**.

{{< /tab >}}
{{< /tabs >}}

### Organization

> [!NOTE]
>
> You must be an organization owner to make changes to the payment information.

{{< tabs >}}
{{< tab name="Docker plan" >}}

To add a payment method:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Under Settings and administration, select **Billing**.
3. Choose your organization from the top-left drop-down.
4. Select **Payment methods** from the left-hand menu.
5. Select **Add payment method**.
6. Enter your new payment information:
    - If you are adding a card, fill out the card information form.
    - If you are adding a Link payment, select **Secure, 1-click checkout with Link**
    and enter your Link **email address** and **phone number**. If you are not
    an existing Link customer, you must fill out the card information form to
    store a card for Link payments.
7. Select **Add payment method**.
8. Optional. You can set a new default payment method by selecting the **Set as default** action.
9. Optional. You can remove non-default payment methods by selecting the **Delete** action.

{{< /tab >}}
{{< tab name="Legacy Docker plan" >}}

To add a payment method:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select your avatar in the top-right corner.
3. From the drop-down menu select **Billing**.
4. Select the organization account you want to update.
5. Select the **Payment methods and billing history** link.
6. In the **Payment Method** section, select **Add payment method**.
7. Enter your new payment information:
    - If you are adding a card, fill out the card information form.
    - If you are adding a Link payment, select **Secure, 1-click checkout with Link**
    and enter your Link **email address** and **phone number**. If you are not
    an existing Link customer, you must fill out the card information form to
    store a card for Link payments.
8. Select **Add payment method**.
9. Select the **Actions** icon, then select **Make default** to ensure that your new payment method applies to all purchases and subscriptions.
10. Optional. You can remove non-default payment methods by selecting the **Actions** icon. Then, select **Delete**.

{{< /tab >}}
{{< /tabs >}}

## Enable pay by invoice

Pay by invoice is available during first-time checkouts and subscription
renewals for subscribers with an annual subscription. When enabled, subscription
renewals will be sent to your billing contact's email address via invoice.
You can pay an invoice using card or US bank account.

> [!NOTE]
>
> Pay by invoice is only available for annual subscriptions. You can't use pay
by invoice for subscription upgrades or purchasing additional seats.

1. Sign in to [Docker Home](https://app.docker.com/).
2. Under **Settings and administration**, select **Billing**.
3. Select **Payment methods**, then **Pay by invoice**.
4. To enable pay by invoice, select the **Pay by invoice** toggle.
5. Confirm your billing contact details. If you need to change them, select
**Change** and enter your new billing contact details.

## Failed payments

> [!NOTE]
>
> You can't manually retry a failed payment. Docker will retry failed payments
based on the retry schedule.

If your subscription payment fails, there is a grace period of 15 days,
including the due date. Docker retries to collect the payment 3 times using the following schedule:

- 3 days after the due date
- 5 days after the previous attempt
- 7 days after the previous attempt

Docker also sends an email notification
`Action Required - Credit Card Payment Failed` with an attached unpaid
invoice after each failed payment attempt.

Once the grace period is over and the invoice is still not paid, the
subscription downgrades to a free plan and all paid features are disabled.

## Redeem a coupon

You can redeem a coupon for any paid Docker subscription.

A coupon can be used when you:

- Sign up to a new paid subscription from a free subscription
- Upgrade an existing paid subscription

You are asked to enter your coupon code when you confirm or enter your
payment method.

If you use a coupon to pay for a subscription, when the coupon expires, your
payment method is charged the full cost of your subscription. If you don't have
a saved payment method, your account downgrades to a free subscription.
