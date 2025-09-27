---
title: Billing FAQs
linkTitle: FAQs
description: Frequently asked questions related to billing
keywords: billing, renewal, payments, faq
tags: [FAQ]
weight: 60
---

### What happens if my subscription payment fails?

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

### Can I manually retry a failed payment?

No. Docker retries failed payments on a [retry schedule](/manuals/billing/faqs.md#what-happens-if-my-subscription-payment-fails).

To ensure a retired payment is successful, verify your default payment is
updated. If you need to update your default payment method, see
[Manage payment method](/manuals/billing/payment-method.md#manage-payment-method).

### Does Docker collect sales tax and/or VAT?

Docker collects sales tax and/or VAT from the following:

- For United States customers, Docker began collecting sales tax on July 1, 2024.
- For European customers, Docker began collecting VAT on March 1, 2025.
- For United Kingdom customers, Docker began collecting VAT on May 1, 2025.

To ensure that tax assessments are correct, make sure that your billing
information and VAT/Tax ID, if applicable, are updated. See
[Update the billing information](/billing/details/).

If you're exempt from sales tax, see
[Register a tax certificate](/billing/tax-certificate/).

### Does Docker offer academic pricing?

For academic pricing, contact the
[Docker Sales Team](https://www.docker.com/company/contact).
