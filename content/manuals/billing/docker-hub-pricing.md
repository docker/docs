---
title: Docker Hub storage pricing
description: Learn how Docker Hub storage pricing is calculated
keywords: docker hub, storage, payments, billing, subscription
weight: 55
---

This guide explains how Docker Hub storage is measured, calculated, and billed
to help you understand your storage consumption and costs.

## How Docker Hub storage is measured

Docker Hub uses the following to track your storage usage:

- Hourly measurement: Storage usage is recorded every hour and measured in
GB-hours. This value reflects the total storage your private repositories
consume during that hour.
- Monthly aggregation: At the end of the month, all hourly storage usage values
are summed up. The total values are divided by the number of hours in the month
to determine the average monthly usage. This value is expressed in GB-month.

For example, if your storage usage fluctuates between 50GB and 100GB over
a 30-day month, the total usage is averaged across the month to provide a
single GB-month value.

This approach ensures your billing reflects your overall storage needs, even if
usage varies day to day.

## How Docker Hub storage is calculated

Docker Hub uses your average monthly usage to determine whether additional
charges apply.

Each Docker subscription plan includes a specific amount of allocated
private repository storage:

- Personal plan: Includes up to 2GB of storage.
- Pro plan: Includes up to 5GB of storage.
- Team plan: Includes up to 50GB of storage.
- Business plan: Includes up to 500GB of storage.

If you go over your allocated private repository storage, you will incur overage
costs. To calculate overage costs:

- The included storage for your plan is subtracted from your average monthly
usage.
- Any remaining amount represents your overage amount, billed as
additional storage consumption.

Hourly storage is measured in GB-hour, but billing calculations are based on
the monthly aggregated GB-month rate.