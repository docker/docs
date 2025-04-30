---
title: Add or update a payment method
weight: 20
description: Learn how to add or update a payment method in Docker Hub
keywords: payments, billing, subscription, supported payment methods, failed payments, coupons
alisases:
    - /billing/core-billing/payment-method/
---

This page describes how to add or update a payment method for your personal account or for an organization.

You can add a payment method or update your account's existing payment method at any time.

> [!IMPORTANT]
>
> If you want to remove all payment methods, you must first downgrade your subscription to a free plan. See [Downgrade](../subscription/change.md).

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
  - ACH transfer with a [verified](manuals/billing/payment-method.md#verify-a-bank-account) US bank account

All currency, for example the amount listed on your billing invoice, is in United States dollar (USD).

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
    - If you are adding a card:
        - Select **Card** and fill out the card information form.
    - If you are adding a Link payment:
        - Select **Secure, 1-click checkout with Link** and enter your Link **email address** and **phone number**.
        - If you are not an existing Link customer, you must fill out the card information form to store a card for Link payments.
    - If you are adding a bank account:
        - Select **US bank account**.
        - Verify your **Email** and **Full name**.
        - If your bank is listed, select your bank's name.
        - If your bank is not listed, select **Search for your bank**.
        - To verify your bank account, see [Verify a bank account](manuals/billing/payment-method.md#verify-a-bank-account).
6. Select **Add payment method**.
7. Optional. You can set a new default payment method by selecting the **Set as default** action.
8. Optional. You can remove non-default payment methods by selecting the **Delete** action.

> [!NOTE]
>
> If you want to set a US bank account as your default payment method, you must
> verify the account first.

{{< /tab >}}
{{< tab name="Legacy Docker plan" >}}

To add a payment method:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select your avatar in the top-right corner.
3. From the drop-down menu select **Billing**.
4. Select the **Payment methods and billing history** link.
5. In the **Payment method** section, select **Add payment method**.
6. Enter your new payment information:
    - If you are adding a card:
        - Select **Card** and fill out the card information form.
    - If you are adding a Link payment:
        - Select **Secure, 1-click checkout with Link** and enter your Link **email address** and **phone number**.
        - If you are not an existing Link customer, you must fill out the card information form to store a card for Link payments.
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
    - If you are adding a card:
        - Select **Card** and fill out the card information form.
    - If you are adding a Link payment:
        - Select **Secure, 1-click checkout with Link** and enter your Link **email address** and **phone number**.
        - If you are not an existing Link customer, you must fill out the card information form to store a card for Link payments.
    - If you are adding a bank account:
        - Select **US bank account**.
        - Verify your **Email** and **Full name**.
        - If your bank is listed, select your bank's name.
        - If your bank is not listed, select **Search for your bank**.
        - To verify your bank account, see [Verify a bank account](manuals/billing/payment-method.md#verify-a-bank-account).
7. Select **Add payment method**.
8. Select **Add payment method**.
9. Optional. You can set a new default payment method by selecting the **Set as default** action.
10. Optional. You can remove non-default payment methods by selecting the **Delete** action.

> [!NOTE]
>
> If you want to set a US bank account as your default payment method, you must
> verify the account first.

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
    - If you are adding a card:
        - Select **Card** and fill out the card information form.
    - If you are adding a Link payment:
        - Select **Secure, 1-click checkout with Link** and enter your Link **email address** and **phone number**.
        - If you are not an existing Link customer, you must fill out the card information form to store a card for Link payments.
8. Select **Add payment method**.
9. Select the **Actions** icon, then select **Make default** to ensure that your new payment method applies to all purchases and subscriptions.
10. Optional. You can remove non-default payment methods by selecting the **Actions** icon. Then, select **Delete**.

{{< /tab >}}
{{< /tabs >}}

## Verify a bank account

There are two ways to verify a bank account as a payment method:

- Instant verification: Docker supports several major banks for instant verification.
- Manual verification: All other banks must be verified manually.

### Instant verification

To verify your bank account instantly, you must sign in to your bank account
from the Docker billing flow:

1. Choose **US bank account** as your payment method.
2. Verify your **Email** and **Full name**.
3. If your bank is listed, select your bank's name or select **Search for your bank**.
4. Sign in to your bank and review the terms and conditions. This agreement
allows Docker to debit payments from your connected bank account.
5. Select **Agree and continue**.
6. Select an account to link and verify, and select **Connect account**.

When the account is verified, you will see a success message in the pop-up modal.

### Manual verification

To verify your bank account manually, you must enter the micro-deposit amount from your bank statement:

1. Choose **US bank account** as your payment method.
2. Verify your **Email** and **First and last name**.
3. Select **Enter bank details manually instead**.
4. Enter your bank details: **Routing number** and **Account number**.
5. Select **Submit**.
6. You will receive an email with instructions on how to manually verify.

Manual verification uses micro-deposits. You should see a small deposit
(e.g. $-0.01) in your bank account in 1-2 business days. Open your manual verification email and enter the amount of this deposit to verify your account.

## Failed payments

> [!NOTE]
>
> You can't manually retry a failed payment. Docker will retry failed payments
based on the retry schedule.

If your subscription payment fails, there is a grace period of 15 days, including the due date. Docker retries to collect the payment 3 times using the following schedule:

- 3 days after the due date
- 5 days after the previous attempt
- 7 days after the previous attempt

Docker also sends an email notification `Action Required - Credit Card Payment Failed` with an attached unpaid invoice after each failed payment attempt.

Once the grace period is over and the invoice is still not paid, the subscription downgrades to a free plan and all paid features are disabled.

## Redeem a coupon

You can redeem a coupon for any paid Docker subscription.

A coupon can be used when you:
- Sign up to a new paid subscription from a free subscription
- Upgrade an existing paid subscription

You are asked to enter your coupon code when you confirm or enter your payment method.

If you use a coupon to pay for a subscription, when the coupon expires, your payment method is charged the full cost of your subscription. If you don't have a saved payment method, your account downgrades to a free subscription.
