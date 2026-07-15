---
title: Billing FAQs
linkTitle: FAQs
description: Find answers to common questions about Docker billing, failed payments, taxes, and pay by invoice.
keywords: billing, renewal, failed payments, sales tax, VAT, academic pricing, pay by invoice
tags: [FAQ]
weight: 80
---

## What happens if my subscription payment fails?

If your subscription payment fails, there is a grace period of 15 days,
including the due date. Docker attempts to collect the payment three times using
the following schedule:

- 3 days after the due date
- 5 days after the previous attempt
- 7 days after the previous attempt

Docker also sends an email notification
`Action Required - Credit Card Payment Failed` with an attached unpaid invoice
after each failed payment attempt.

If the invoice remains unpaid after the grace period, the
subscription downgrades to a free subscription and all paid features are
disabled.

## Can I manually retry a failed payment?

Yes. If your payment fails, select **Pay now** to retry the payment through
Stripe.

Before retrying, verify that your default payment method is up to date. For
instructions, see
[Manage a payment method](/manuals/billing/payment-method.md#manage-payment-method).

## Does Docker collect sales tax and VAT?

Docker collects sales tax or VAT from the following customers:

- For United States customers, Docker began collecting sales tax on
  July 1, 2024.
- For European customers, Docker began collecting VAT on March 1, 2025.
- For United Kingdom customers, Docker began collecting VAT on May 1, 2025.

To help ensure correct tax assessments, keep your
[billing information](/manuals/billing/details.md) up to date. You can add a
tax ID or VAT ID when you
[purchase a Docker plan](/manuals/subscription/manage.md#set-up-a-new-plan).

If you're exempt from sales tax, see
[Submit a tax certificate](/manuals/billing/tax-certificate.md).

## Does Docker offer academic pricing?

For academic pricing, contact the
[Docker Sales Team](https://www.docker.com/company/contact).

## Can I use pay by invoice for upgrades or additional seats?

No. Pay by invoice is only available for renewing annual subscriptions, not for
purchasing upgrades or additional seats. You must use card payment or US bank
accounts for these changes.

For a list of supported payment methods, see
[Add or update a payment method](/manuals/billing/payment-method.md).
