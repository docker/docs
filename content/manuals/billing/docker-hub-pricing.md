---
title: Docker Hub storage pricing
description: Learn how Docker Hub storage pricing is calculated
keywords: docker hub, storage, payments, billing, subscription
weight: 55
---

This guide explains how Docker Hub storage is measured, calculated, and billed
to help you understand your storage consumption and costs.

## How storage is measured

Docker Hub measures storage using:
- Hourly measurement: Storage usage is recorded every hour and expressed in **GB-hours**. This value represents the total storage your repositories consume during each hour.
- Monthly aggregation: At the end of each month, all hourly storage usage values for the month are summed up:
    - The total is divided by the number of hours in that month. For example, 730 hours for 30 days or 744 hours for 31 days.
    - The result is the average monthly usage, expressed in **GB-month**.

## How storage is calculated

Docker subscription plans include a specific amount of allocated
private repository storage:

- Personal plan: Includes up to 2GB of storage.
- Pro plan: Includes up to 5GB of storage.
- Team plan: Includes up to 50GB of storage.
- Business plan: Includes up to 500GB of storage.

Docker Hub determines additional charges based on your average monthly usage of private repository storage.

If you go over your allocated private repository storage, you will incur overage
costs. To calculate overage costs the included storage for your plan is subtracted from your average monthly
usage.

## Docker Hub consumption pricing

At the end of the month, Docker calculates your total storage usage
and compares it to your plan's included amount. If applicable, the overage cost
is billed to your account as an overage invoice.

There are two billing models for paying for additional Docker Hub storage:

- Pre-pay: Pay in advance for a specified amount of storage.

    > [!NOTE]
    >
    > Pre-purchased storage expires at the end of your subscription period.

- Post-pay: Receive an overage invoice for storage usage that exceeds your subscription plan's included amount
at the end of your billing cycle.

## Examples

### Business plan with pre-pay

In the following example, a customer with a Business plan has 500GB included in their subscription plan. They pre-pay
for 1700 GB.
- In January, they use 100 GB-month, meaning they did not use any of their pre-pay storage. Their pre-pay storage rolls over to the next month.
- In February, they use 650 GB-month, exceed their base allocation, and use 150GB from their pre-pay storage.
- In March, they use 1800 GB-month, exceed their base allocation, and use 1300GBs from their pre-pay storage.
- In April, they use 950 GB-month, exceed their base allocation, and going over their pre-pay storage. This results in an invoice of $14.00 for the storage overage.

|                          | January | February | March | April  |
|--------------------------|---------|----------|-------|--------|
| Included GB-month        | 500     | 500      | 500   | 500    |
| Pre-purchased GB         | 1700    | 1700     | 1700  | 1700   |
| Used storage in month    | 100     | 650      | 1800  | 950    |
| Remaining pre-purchased  | 1700    | 1550     | 250   | -200   |
| Overage invoice          | $0.00   | $0.00    | $0.00 | $14.00 |

For information on storage pricing, see the [Docker Pricing](https://www.docker.com/pricing/) page.

### Business plan with post-pay

In the following example, a customer with a Business plan has 500GB included in their subscription plan. They do
not pre-pay for additional storage consumption.
- In January, they use 100 GB-month and do not exceed their base allocation.
- In February, they use 650 GB-month, going over their base allocation by 150 GB-month. They are sent
an overage invoice for $10.50.
- In March, they use 1800 GB-month, going over their base allocation by 1300 GB-month. They are sent
an overage invoice for $91.00.
- In April, they use 950 GB-month, going over their base allocation by 450 GB-month. They are sent an
overage invoice for $31.50.

| Metric                            | January | February | March  | April  |
|-----------------------------------|---------|----------|--------|--------|
| Included GB-month                 | 500     | 500      | 500    | 500    |
| Used storage in month             | 100     | 650      | 1800   | 950    |
| Overage in GB-month               | 0       | 150      | 1300   | 450    |
| Overage invoice                   | $0.00   | $10.50   | $91.00 | $31.50 |

For information on storage pricing, see the [Docker Pricing](https://www.docker.com/pricing/) page.