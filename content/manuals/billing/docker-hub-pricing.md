---
title: Docker Hub storage pricing
description: Learn how Docker Hub storage pricing is calculated
keywords: docker hub, storage, payments, billing, subscription
weight: 55
---

This guide explains how Docker Hub storage is measured, calculated, and billed
to help you understand your storage consumption and costs.

> [!NOTE]
>
> Docker Hub plan limits will take effect on March 1, 2025. No charges on Docker Hub pulls
or storage will be incurred from December 10, 2024 and February 28, 2025.

## How storage is measured

Docker Hub measures storage using:
- Hourly measurement: Storage usage is recorded every hour and expressed in **GB-hours**. This value represents the total storage your repositories consume during each hour.
- Monthly aggregation: At the end of each month, hourly storage usage is divided by the number of hours in that month.

  For example, 900 GB-hours of storage in a 30-day month (720 hours) equates to 900 / 720 = 1.25 gigabytes of monthly storage (GB-month).

## How storage is calculated

[Docker subscription plans](/manuals/subscription/details.md) include a specific amount of allocated
private repository storage:

- Personal plan: Includes up to 2GB of storage.
- Pro plan: Includes up to 5GB of storage.
- Team plan: Includes up to 50GB of storage.
- Business plan: Includes up to 500GB of storage.

Docker Hub determines additional charges based on your average monthly usage of private repository storage. If a repository is private at any point within an hour, it is counted as private for the full hour.

If you go over your allocated private repository storage, you will incur overage
costs. To calculate overage costs, the included storage for your plan is subtracted from your average monthly
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

### Storage carryover

If you pre-pay for storage, your purchased storage is valid for your entire subscription period. You can use it any time during that period, and any unused portion will roll over to the next month until the subscription period ends.

In the following example, a customer with an annual Business plan pre-pays for 500GB of storage for the year. Their plan includes a base allocation of 500GB of storage per month.
- In January, they use 510 GB-month, exceed their base allocation, and use 10GB from their pre-paid storage. Their remaining pre-paid
storage is 490GB.
- In February, they use 450 GB-month, and do not exceed their base allocation. They do not use any of their pre-paid storage, so it remains at 490GB.
- In March, they use 600 GB-month, exceed their base allocation, and use 100GB from their pre-paid storage. Their remaining pre-paid storage is 390GB.

|                                 | January  | February | March   |
|---------------------------------|----------|----------|---------|
| Included GB-month	              | 500	     | 500      | 500     |
| Used storage in month	          | 510	     | 450	    | 600     |
| Overage in GB-month	            | 10	     | 0	      | 100     |
| Remaining pre-purchased storage |	490	     | 490	    | 390     |

At the end of March, the customer has 390GB of pre-purchased storage left to use for the rest of the year.

## Examples

### Business plan with pre-pay

In the following example, a customer with a Business plan has 500GB included in their subscription plan. They pre-pay
for 1700 GB.
- In January, they use 100 GB-month, meaning they did not use any of their pre-pay storage. Their pre-pay storage rolls over to the next month.
- In February, they use 650 GB-month, exceed their base allocation, and use 150GB from their pre-pay storage.
- In March, they use 1800 GB-month, exceed their base allocation, and use 1300GB from their pre-pay storage.
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

|                                   | January | February | March   | April  |
|-----------------------------------|---------|----------|---------|--------|
| Included GB-month                 | 500     | 500      | 500     | 500    |
| Used storage in month             | 100     | 650      | 1800    | 950    |
| Overage in GB-month               | 0       | 150      | 1300    | 450    |
| Overage invoice                   | $0.00   | $10.50   | $130.00 | $45.00 |

For information on storage pricing, see the [Docker Pricing](https://www.docker.com/pricing/) page.