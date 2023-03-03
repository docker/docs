---
title: FAQs
description: Common FAQs related to billing
keywords: billing, renewal, payments
---

## Where can I view my billing date?

Navigate to the **Plan** tab in your billing settings. The billing date is located near the bottom-right.

## What credit and debit cards are supported?

- Visa
- MasterCard
- American Express
- Discover
- JCB
- Diners
- UnionPay

## What happens if my subscription payment fails?

If your subscription payment fails, there is a grace period of 15 days, including the due date. Docker will retry to collect payment 8 times using the following schedule:

- Next day after the due date (aka day 1 after the due date)
- Next day after the first retry (aka day 2 after the due date)
- Then every other day (aka day 4, 6, 8, 10, 12, 14 after the due date)

Docker will also send an email notification `Action Required - Credit Card Payment Failed` with an attached unpaid invoice (PDF) after each failed payment attempt. 

Once the grace period is over and the invoice has still not been paid, the user account or organization is downgraded to a free plan and all paid features are disabled. 