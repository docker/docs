---
title: Invoices and billing history
weight: 40
description: Discover how to view invoices and your billing history
keywords: payments, billing, subscription, invoices, renewals, invoice management, billing administration, pay invoice
aliases:
    - /billing/core-billing/history/
---

Learn how to view and pay invoices, view your billing history, and verify
your billing renewal date. All monthly and annual subscriptions are
automatically renewed at the end of the subscription term using your default
payment method.

{{% include "tax-compliance.md" %}}

## View an invoice

Your invoice includes the following:

- Invoice number
- Date of issue
- Due date
- Your "Bill to" information
- Amount due (in USD)
- Pay online: Select this link to pay your invoice online
- Description of your order, quantity if applicable, unit price, and amount
(in USD)
- Subtotal, discount (if applicable), and total

The information listed in the "Bill to" section of your invoice is based on your billing information. Not all fields are required. The billing information
includes the following:

- Name (required): The name of the administrator or company
- Address (required)
- Email address (required): The email address that receives all billing-related
emails for the account
- Phone number
- Tax ID or VAT

You can’t make changes to a paid or unpaid billing invoice. When you update your billing information, this change won't update an existing invoice. If you need
to update your billing information, make sure you do so before your subscription renewal date when your invoice is finalized.
For more information, see [Update billing information](details.md).

## Pay an invoice

> [!NOTE]
>
> Pay by invoice is only available for subscribers on an annual billing cycle.
To change your billing cycle, see [Change your billing cycle](/manuals/billing/cycle.md).

You can pay an invoice from the Docker Billing Console:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Under Settings and administration, select **Billing**.
3. Select **Invoices** and locate the invoice you want to pay.
4. In the **Actions** column, select **Pay invoice**.
5. Fill out your payment details and select **Pay**.

When your payment has processed, the invoice's **Status** column will update to
**Paid** and you will receive a confirmation email.

If you choose to pay using a US bank account, you must verify the account. For
more information, see [Verify a bank account](manuals/billing/payment-method.md#verify-a-bank-account).

## View renewal date

{{< tabs >}}
{{< tab name="Docker plan" >}}

You receive your invoice when the subscription renews. To verify your renewal
date, sign in to the [Docker Home Billing](https://app.docker.com/billing). Your renewal date and amount are displayed on your subscription plan card.

{{< /tab >}}
{{< tab name="Legacy Docker plan" >}}

You receive your invoice when the subscription renews. To verify your renewal
date:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select your user avatar to open the drop-down menu.
3. Select **Billing**.
4. Select the user or organization account to view the billing details. Here you can find your renewal date and the renewal amount.

{{< /tab >}}
{{< /tabs >}}

## Include your VAT number on your invoice

> [!NOTE]
>
> If the VAT number field is not available, complete the
[Contact Support form](https://hub.docker.com/support/contact/). This field
may need to be manually added.

{{< tabs >}}
{{< tab name="Docker plan" >}}

To add or update your VAT number:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Under Settings and administration, select **Billing**.
3. Select **Billing information** from the left-hand menu.
4. Select **Change** on your billing information card.
5. Ensure the **I'm purchasing as a business** checkbox is checked.
6. Enter your VAT number in the Tax ID section.

    > [!IMPORTANT]
    >
    > Your VAT number must include your country prefix. For example, if you are
    entering a VAT number for Germany, you would enter `DE123456789`.

7. Select **Update**.

Your VAT number will be included on your next invoice.

{{< /tab >}}
{{< tab name="Legacy Docker plan" >}}

To add or update your VAT number:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. For user accounts, select your user avatar, then **Billing**. For
organizations, select the name of the organization.
3. Select the **Billing address** link.
4. In the **Billing Information** section, select **Update information**.
5. Enter your VAT number in the Tax ID section.

    > [!IMPORTANT]
    >
    > Your VAT number must include your country prefix. For example, if you are
    entering a VAT number for Germany, you would enter `DE123456789`.

6. Select **Save**.

Your VAT number will be included on your next invoice.

{{< /tab >}}
{{< /tabs >}}

## View billing history

You can view your billing history and download past invoices for a personal
account or organization.

### Personal account

{{< tabs >}}
{{< tab name="Docker plan" >}}

To view billing history:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Under Settings and administration, select **Billing**.
3. Select **Invoices** from the left-hand menu.
4. Optional. Select the **Invoice number** to open invoice details.
5. Optional. Select the **Download** button to download an invoice.

{{< /tab >}}
{{< tab name="Legacy Docker plan" >}}

To view billing history:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select your avatar in the top-right corner.
3. From the drop-down menu select **Billing**.
4. Select the **Payment methods and billing history** link.
    You can find your past invoices in the **Invoice History** section.

From here you can download an invoice.

{{< /tab >}}
{{< /tabs >}}

### Organization

You must be an owner of the organization to view the billing history.

{{< tabs >}}
{{< tab name="Docker plan" >}}

To view billing history:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Under Settings and administration, select **Billing**.
3. Select **Invoices** from the left-hand menu.
4. Optional. Select the **invoice number** to open invoice details.
5. Optional. Select the **download** button to download an invoice.

{{< /tab >}}
{{< tab name="Legacy Docker plan" >}}

To view billing history:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select your avatar in the top-right corner.
3. From the drop-down menu select **Billing**.
4. Select the **Payment methods and billing history** link. You can find your
past invoices in the **Invoice History** section.

You can also download an invoice from this page if needed.

{{< /tab >}}
{{< /tabs >}}
