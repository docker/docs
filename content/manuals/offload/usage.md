---
title: Docker Offload usage and billing
linktitle: Usage & billing
weight: 30
description: Learn about Docker Offload usage and how to monitor your cloud resources.
keywords: cloud, usage, cloud minutes, shared cache, top repositories, cloud builder, Docker Offload
---

{{< summary-bar feature_name="Docker Offload" >}}

> [!NOTE]
>
> All free trial usage granted for the Docker Offload Beta expire after 90 days
> from the time they are granted. To continue using Docker Offload Beta after
> your usage expires, you can enable on-demand usage at [Docker Home
> Billing](https://app.docker.com/billing).

## Understand usage and billing models

Docker Offload offers two usage models to fit different team needs and usage
patterns:

- Committed usage: This provides a committed amount of cloud compute time for
  your organization.
- On-demand usage: This provides pay-as-you-go flexibility. You can enable or
  disable on-demand usage in [Billing](#manage-billing).

## Manage billing

For Docker Offload, you can view and configure billing on the **Docker Offload**
page in [Docker Home Billing](https://app.docker.com/billing). On this page, you
can:

- View your committed usage
- View rates for cloud resources
- Manage on-demand billing, including setting a monthly limit
- Track your organization's Docker Offload usage
- Add or change payment methods

You must be an organization owner to manage billing. For more general
information about billing, see [Billing](../billing/_index.md).

## Monitor your usage

The **Offload activity** page in Docker Home provides visibility into how you
are using cloud resources to run containers.

To monitor your usage:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Select the account for which you want to monitor usage.
3. Select **Offload** > **Offload activity**.

### Overview metrics

Key metrics at the top of the page summarize your Docker Offload usage:

- **Total duration**: The total time spent in Offload sessions
- **Average duration**: The average time per Offload session
- **Total sessions**: The total number of Offload sessions
- **Unique images used**: The number of distinct container images used across
  sessions
- **Unique users**: The number of different users in Docker Offload sessions

### Filter and export your data

You can filter the Offload activity data by:

- **Period**: Select a preset time period or choose a custom date range
- **Users**: Organization owners and members with analytics permissions can
  filter by specific users

Export your session data by selecting the **Download CSV** button. The exported
file includes:

- Session ID
- Username
- Image
- Started time
- Ended time
- Duration (in seconds)
- Status
- Container count

The CSV export includes data for your selected date range and user filters,
letting you download exactly what you're viewing.

### Activity cards

The following cards provide insights into your Docker Offload usage:

- **Offload usage**: Shows your usage trends over time and cloud resource
  consumption patterns. If you have billing permissions, you can select the card
  to view detailed billing information.
- **Popular images**: Shows the top 4 most frequently used container images in
  your Docker Offload sessions. Select the card to see more images.
- **Top users**: Shows the top 4 users by session count and duration. Select
  the card to see more users.

### Offload sessions

A detailed list of Offload sessions appears below the activity cards. The list:

- Starts with any currently active sessions
- Shows session details including start time, duration, images used, and user
  information
- Can be filtered using the date and user filters described above
- Displays **Offload sessions** if you have organization-wide analytics
  permissions, or **My Offload sessions** if viewing only your own data

Select any session to view more details in a side panel.