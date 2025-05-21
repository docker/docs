---
title: Billing FAQs
linkTitle: FAQs
description: Frequently asked questions related to billing
keywords: billing, renewal, payments, faq
tags: [FAQ]
weight: 60
---

### What credit and debit cards are supported?

- Visa
- MasterCard
- American Express
- Discover
- JCB
- Diners
- UnionPay
- Link
- ACH transfer with a [verified](manuals/billing/payment-method.md#verify-a-bank-account) US bank account

### What currency is supported?

United States dollar (USD).

### What happens if my subscription payment fails?

If your subscription payment fails, there is a grace period of 15 days, including the due date. Docker retries to collect the payment 3 times using the following schedule:

- 3 days after the due date
- 5 days after the previous attempt
- 7 days after the previous attempt

Docker also sends an email notification `Action Required - Credit Card Payment Failed` with an attached unpaid invoice after each failed payment attempt.

Once the grace period is over and the invoice is still not paid, the subscription downgrades to a free plan and all paid features are disabled.

### Can I manually retry a failed payment?

No. Docker retries failed payments on a [retry schedule](/manuals/billing/faqs.md#what-happens-if-my-subscription-payment-fails).

To ensure a retired payment is successful, verify your default payment is
updated. If you need to update your default payment method, see
[Manage payment method](/manuals/billing/payment-method.md#manage-payment-method).

### Does Docker collect sales tax and/or VAT?

Docker collects sales tax and/or VAT from the following:

- Docker began collecting sales tax from United States
customers on July 1, 2024.
- For European customers, Docker began collecting VAT on March 1, 2025.
- For United Kingdom customers, Docker began collecting VAT on May 1, 2025.

To ensure that tax assessments are correct, make sure that your billing
information and VAT/Tax ID, if applicable, are updated. See
[Update the billing information](/billing/details/).

If you're exempt from sales tax, see
[Register a tax certificate](/billing/tax-certificate/).

### How do I certify my tax exempt status?

If you're exempt from sales tax, you can [register a valid tax exemption certificate](./tax-certificate.md) with Docker's Support team. [Contact Support](https://hub.docker.com/support/contact) to get started.

### Does Docker offer academic pricing?

Contact the [Docker Sales Team](https://www.docker.com/company/contact).

### Do I need to do anything at the end of my subscription term?

No. All monthly and annual subscriptions are automatically renewed at the end of the term using the original form of payment.