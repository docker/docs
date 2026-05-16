---
title: Docker Offload usage
linktitle: Usage
weight: 30
description: Learn about Docker Offload usage and how to monitor your cloud resources.
keywords: cloud, usage, Offload
---

{{< summary-bar feature_name="Docker Offload" >}}


The **Offload activity** page in Docker Home provides visibility into user
activity and session metrics for Docker Offload.

To monitor your usage:

1. Sign in to [Docker Home](https://app.docker.com/).
2. If you have access to multiple organizations, select the organization
   associated with your Docker Offload subscription.
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
- **Additional Filters**: Filter by active sessions and session duration.

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
  consumption patterns.
- **Popular images**: Shows the top 4 most frequently used container images in
  your Docker Offload sessions. Select the card to see more images.
- **Top Offload users**: Shows the top 4 users by session count and duration. Select
  the card to see more users.

### Offload sessions

A detailed list of Offload sessions appears following the activity cards. The list:

- Starts with any currently active sessions
- Shows session details including start time, duration, images used, and user
  information
- Can be filtered using the date and user filters described previously
- Displays **Offload sessions** if you have organization-wide analytics
  permissions, or **My Offload sessions** if viewing only your own data

Select any session to view more details in a side panel.