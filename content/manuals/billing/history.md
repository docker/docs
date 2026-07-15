---
title: View and manage invoices and billing history
linkTitle: Invoices
weight: 60
description: Learn how to view and pay Docker invoices, download billing history, and find your subscription renewal date.
keywords: payments, billing, subscription, invoices, renewals, billing history, pay invoice, VAT number
aliases:
  - /billing/core-billing/history/
---

Learn how to view and pay invoices, view your billing history, and verify
your billing renewal date. All monthly and annual subscriptions are
automatically renewed at the end of the subscription term using your default
payment method.

## View an invoice

Your invoice includes the following:

- Invoice number
- Date of issue
- Due date
- Your **Bill to** information
- Amount due (in USD)
- **Pay online** link
- Description of your order, quantity if applicable, unit price, and
  amount (in USD)
- Subtotal, discount (if applicable), and total

The information listed in the **Bill to** section of your invoice is based on
your billing information. Not all fields are required. The billing information
includes the following:

- Name (required): The name of the administrator or company
- Address (required)
- Email address (required): The email address that receives all billing-related
  emails for the account
- Phone number
- Tax ID or VAT

You can't change a paid or unpaid invoice. Updating your billing information
doesn't update an existing invoice.

Update your billing information before your subscription renewal date, when
Docker finalizes your invoice.

For more information, see [Update billing information](details.md).

## Pay an invoice

> [!NOTE]
>
> Pay by invoice is only available for subscribers on an annual billing cycle.
> To change your billing cycle, see
> [Change your billing cycle](/manuals/billing/cycle.md).

If you've selected pay by invoice for your subscription, you'll receive email
reminders to pay your invoice at 10 days before the due date, on the due date,
and 15 days after the due date.

You can pay an invoice from the Docker Billing Console:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Billing**.
1. Select **Invoices** and locate the invoice you want to pay.
1. In the **Actions** column, select **Pay invoice**.
1. Fill out your payment details and select **Pay**.

After Docker processes your payment, the invoice's **Status** changes to
**Paid**, and you receive a confirmation email.

If you pay using a US bank account, you must
[verify the bank account](/manuals/billing/payment-method.md#verify-a-bank-account).

## View renewal date

You receive your invoice when the subscription renews. To verify your renewal
date:

1. Sign in to [Docker Home Billing](https://app.docker.com/billing).
1. Find your renewal date and amount on your subscription plan card.

## Include your VAT number on your invoice

> [!NOTE]
>
> If the VAT number field is not available, complete the
> [Contact Support form](https://hub.docker.com/support/contact/). This field
> may need to be manually added.

To add or update your VAT number:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Billing**.
1. Select **Billing information** from the left-hand menu.
1. Select **Change** on your billing information card.
1. Select the **I'm purchasing as a business** checkbox.
1. Enter your VAT number in the Tax ID section.

   > [!IMPORTANT]
   >
   > Your VAT number must include your country prefix. For example, enter
   > `DE123456789` for a German VAT number.

1. Select **Update**.

Your VAT number will be included on your next invoice.

## View billing history

You can view your billing history and download past invoices for a personal
account or organization.

### Personal account

To view billing history:

1. Sign in to [Docker Home](https://app.docker.com/) and select your personal
   account.
1. Select **Billing**.
1. Select **Invoices** from the left-hand menu.
1. Optional. Select the **Invoice number** to open invoice details.
1. Optional. Select the **Download** button to download an invoice.

### Organization

You must be an owner of the organization to view the billing history.

To view billing history:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Billing**.
1. Select **Invoices** from the left-hand menu.
1. Optional. Select the **Invoice number** to open invoice details.
1. Optional. Select **Download** to download an invoice.
