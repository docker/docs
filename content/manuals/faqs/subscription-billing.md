---
title: Subscription and billing FAQs
linkTitle: Subscription and billing
description: Frequently asked questions about Docker subscriptions, billing, failed payments, taxes, and plans.
keywords: subscription faqs, billing, docker plans, renewal, failed payments, sales tax, VAT, academic pricing, pay by invoice, subscription transfer
tags: [FAQ]
toc_max: 2
weight: 35
aliases:
  - /subscription/faq/
  - /billing/faqs/
  - /subscription-billing/faqs/subscription/
  - /subscription-billing/faqs/billing/
---

For more information on Docker subscriptions, see
[Subscription and billing](/manuals/subscription-billing/_index.md).

## Subscriptions

### Can I transfer my subscription from one user or organization account to another?

Subscriptions are non-transferable between accounts or organizations.

### Can I pause or delay my Docker subscription?

You can't pause or delay a subscription, but you can downgrade your
subscription. If a subscription invoice isn't paid by the due date, there's a
15-day grace period starting from the due date.

### Does Docker offer academic pricing?

For academic pricing, contact the
[Docker Sales Team](https://www.docker.com/company/contact).

### How can I contribute to Docker content?

Docker offers two content contribution programs:

- [Docker-Sponsored Open Source Program (DSOS)](/manuals/docker-hub/repos/manage/trusted-content/dsos-program.md)
  for open source projects
- [Docker Verified Publisher (DVP)](/manuals/docker-hub/repos/manage/trusted-content/dvp-program.md)
  for commercial publishers

You can also join the
[Developer Preview Program](https://www.docker.com/community/get-involved/developer-preview/)
or sign up for early access programs to participate in research and try new
features.

## Payments

### What happens if my subscription payment fails?

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

### Can I manually retry a failed payment?

Yes. If your payment fails, select **Pay now** to retry the payment through
Stripe.

Before retrying, verify that your default payment method is up to date. For
instructions, see
[Manage a payment method](/manuals/subscription-billing/manage/payment-method.md#manage-payment-method).

### Can I use pay by invoice for upgrades or additional seats?

No. Pay by invoice is only available for renewing annual subscriptions, not for
purchasing upgrades or additional seats. You must use card payment or US bank
accounts for these changes.

For a list of supported payment methods, see
[Add or update a payment method](/manuals/subscription-billing/manage/payment-method.md).

> [!TIP]
>
> Need to upgrade? <a href="https://www.docker.com/pricing?ref=Docs&refAction=DocsSubscriptionFaq" id="pricing-link" class="link" rel="noopener">Compare Docker Team and Docker Business</a> to choose the plan that best fits your team's needs.

## Taxes

### Does Docker collect sales tax and VAT?

Docker collects sales tax or VAT from the following customers:

- For United States customers, Docker began collecting sales tax on
  July 1, 2024.
- For European customers, Docker began collecting VAT on March 1, 2025.
- For United Kingdom customers, Docker began collecting VAT on May 1, 2025.

To help ensure correct tax assessments, keep your
[billing information](/manuals/subscription-billing/manage/details.md) up to date. For details on
adding a VAT number or submitting a US tax exemption certificate, see
[Taxes](/manuals/subscription-billing/manage/tax-certificate.md).
