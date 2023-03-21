---
title: FAQs
description: Common FAQs related to billing
keywords: billing, renewal, payments
---

### Where can I view my billing date?

Navigate to the **Plan** tab in your billing settings. The billing date is located near the bottom-right.

### What credit and debit cards are supported?

- Visa
- MasterCard
- American Express
- Discover
- JCB
- Diners
- UnionPay

### What happens if my subscription payment fails?

If your subscription payment fails, there is a grace period of 15 days, including the due date. Docker retries to collect the payment 8 times using the following schedule:

- Next day after the due date (aka day 1 after the due date)
- Next day after the first retry (aka day 2 after the due date)
- Then every other day (aka day 4, 6, 8, 10, 12, 14 after the due date)

Docker also sends an email notification `Action Required - Credit Card Payment Failed` with an attached unpaid invoice after each failed payment attempt. 

Once the grace period is over and the invoice is still not paid, the subscription is downgraded to a free plan and all paid features are disabled. 

## What billing-related emails will I receive from Docker Hub?

Docker Hub sends the following billing-related emails:

- Confirmation of a new subscription.
- Confirmation of paid invoices. 
- Notifications of credit or debit card payment failures. 
- Notifications of credit or debit card expiration. 
- Confirmation of a cancelled subscription 
- Reminders of subscription renewals for annual subscribers. This is sent 14 days before the renewal date. 

### Does Docker offer academic pricing?

Contact the [Docker Sales Team](https://www.docker.com/company/contact){:target="blank" rel="noopener" class=""}.

### Do I need to do anything at the end of my subscription term?

No. All monthly and annual subscriptions are automatically renewed at the end of the term using the original form of payment.
