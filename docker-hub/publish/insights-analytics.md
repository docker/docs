---
title: Insights & analytics
description: Provides usage statistics of your images on Docker Hub.
keywords: docker hub, hub, insights, analytics, api, verified publisher
---

Insights and analytics provides usage analytics for your Docker Verified
Publisher (DVP) images on Docker Hub. With this tool, you have self-serve access
to metrics as both raw data and summary data for a desired time span. You can
view number of image pulls by tag or by digest, and get breakdowns by
geolocation, cloud provider, and client, and more.

## Exporting analytics data

You can access the data either from the Docker Hub website, at
`hub.docker.com/orgs/{namespace}/insights`, or using the
[DVP Data API](/docker-hub/api/dvp/){: target="_blank" rel="noopener"
class="_"}. All members of an organization have access to the analytics data.

The data is available as a downloadable CSV file, in a weekly (Monday through
Sunday) or monthly format (available on the first day of the following calendar
month). You can import this data into your own systems, or you can analyze it
manually as a spreadsheet.

### Export data using the website

Here's how to export usage data for your organization's images using the Docker
Hub website.

1.  Sign in to [Docker Hub](https://hub.docker.com/){: target="_blank"
    rel="noopener" class="_"} and select **Organizations**.

2.  Choose your organization and select **Insights and analytics**.

    ![Organization overview page, with the Insights and Analytics tab](./images/organization-tabs.png)

3.  Set the time span for which you want to export analytics data.

    The downloadable CSV files for summary and raw data appear on the right-hand
    side.

    ![Filtering options and download links for analytics data](./images/download-analytics-data.png)

### Export data using the API

The HTTP API endpoints are available at:
`https://hub.docker.com/api/publisher/analytics/v1`. Learn how to export data
using the API in the [DVP Data API documentation](/docker-hub/api/dvp/){:
target="_blank" rel="noopener" class="_"}.

## Data points

Export data in either raw or summary format. Each format contains different data
points and with different structure.

The following sections describe the available data points for each format. The
**Available from** column shows when the field was first added.

### Raw data

The raw data format contains the following data points. Each row in the CSV file
represents an image pull.

| Data point                    | Value                                             | Description                                                                           | Available from  |
| ----------------------------- | ------------------------------------------------- | ------------------------------------------------------------------------------------- | --------------- |
| Timestamp                     | `YYYY-MM-DD 00:00:00`                             | Date and time of the request.                                                         | January 1, 2022 |
| Namespace                     | `String`                                          | Docker organization (image namespace).                                                | January 1, 2022 |
| Repository                    | `String`                                          | Image name.                                                                           | January 1, 2022 |
| Reference                     | `String`                                          | Image digest or tag, as requested.                                                    | January 1, 2022 |
| Digest                        | `String`                                          | Image digest.                                                                         | January 1, 2022 |
| Tag (included when available) | `String`                                          | Tag name. Only available if the request referred to a tag.                            | January 1, 2022 |
| Action day                    | `YYYY-MM-DD`                                      | The date part of the timestamp.                                                       | January 1, 2022 |
| HTTP method                   | `String`                                          | HTTP method used in the request, see [registry API documentation][1] for details.     | January 1, 2022 |
| Action                        | `pull_by_tag`, `pull_by_digest`, `version_check`  | Request type, see [Action classification rules][2].                                   | January 1, 2022 |
| Type                          | `business`, `isp`, `hosting`, `education`, `null` | The industry from which the event originates.                                         | January 1, 2022 |
| Host                          | `String`                                          | The cloud service provider used in an event.                                          | January 1, 2022 |
| Country                       | `String`                                          | Request origin country.                                                               | January 1, 2022 |
| Domain                        | `String`                                          | Request origin domain, see [Privacy][3].                                              | October n, 2022 |
| User agent tool               | `String`                                          | The application a user used to pull an image (for example, `docker` or `containerd`). | January 1, 2022 |
| User agent version            | `String`                                          | The version of the application used to pull an image.                                 | January 1, 2022 |

[1]: /registry/spec/api/
[2]: #action-classification-rules

### Summary data

The summary data format contains the following data points for each namespace,
repository, and reference (tag or digest), for the selected time span.

| Data point        | Value     | Description                                       | Available from  |
| ----------------- | --------- | ------------------------------------------------- | --------------- |
| Unique IP address | `String`  | Number of unique IP addresses, see [Privacy][3].  | January 1, 2022 |
| Pull by tag       | `Integer` | GET request, by digest or by tag.                 | January 1, 2022 |
| Pull by digest    | `Integer` | GET or HEAD request by digest, or HEAD by digest. | January 1, 2022 |
| Version check     | `Integer` | HEAD by tag, not followed by a GET                | January 1, 2022 |

[3]: #privacy

### Action classification rules

An action represents the multiple request events associated with a
`docker pull`. Pulls are grouped by category to make the data more meaningful
for understanding user behavior and intent. The categories are:

- Version check
- Pull by tag
- Pull by digest

Automated systems frequently check for new versions of your images. Being able
to distinguish between "version checks" in CI versus actual image pulls by a
user grants you more insight into your users' behavior.

The following table describe the rules applied for determining intent behind
pulls. To provide feedback or ask questions about these rules,
[fill out the Google Form](https://forms.gle/nb7beTUQz9wzXy1b6){:
target="_blank" rel="noopener" class="_"}.

| Starting event | Reference | Followed by                                                     | Resulting action | Use case(s)                                                                                                    | Notes                                                                                                                                                                                                                                                                                          |
| :------------- | :-------- | :-------------------------------------------------------------- | :--------------- | :------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| HEAD           | tag       | N/A                                                             | Version check    | User already has all layers existing on local machine                                                          | This is similar to the use case of a pull by tag when the user already has all the image layers existing locally, however, we are able to differentiate the user intent and classify accordingly.                                                                                              |
| GET            | tag       | N/A                                                             | Pull by tag      | User already has all layers existing on local machine and/or the image is single-arch                          |
| GET            | tag       | Get by different digest                                         | Pull by tag      | Image is multi-arch                                                                                            | Second GET by digests must be different from the first                                                                                                                                                                                                                                         |
| HEAD           | tag       | GET by same digest                                              | Pull by tag      | Image is multi-arch but some or all image layers already exist on the local machine.                           | The HEAD by tag will send the most current digest, the following GET must be by that same digest. There may occur an additional GET, if the image is multi-arch (see the next row in this table). If the user doesn't want the most recent digest, then the user would perform HEAD by digest. |
| HEAD           | tag       | GET by the same digest, then a second GET by a different digest | Pull by tag      | Image is multi-arch                                                                                            | The HEAD by tag will send the most recent digest, the following GET must be by that same digest. Since the image is multi-arch, there is a second GET by a different digest. If the user doesn't want the most recent digest, then the user would perform HEAD by digest.                      |
| HEAD           | tag       | GET by same digest, then a second GET by different digest       | Pull by tag      | Image is multi-arch                                                                                            | The HEAD by tag will send the most current digest, the following GET must be by that same digest. Since the image is multi-arch, there is a second GET by a different digest. If the user doesn't want the most recent digest, then the user would perform HEAD by digest.                     |
| GET            | digest    | N/A                                                             | Pull by digest   | User already has all layers existing on local machine and/or the image is single-arch                          |
| HEAD           | digest    | N/A                                                             | Pull by digest   | User already has all layers existing on their local machine.                                                   |
| GET            | digest    | GET by different digest                                         | Pull by digest   | Image is multi-arch                                                                                            | The second GET by digest must be different from the first                                                                                                                                                                                                                                      |
| HEAD           | digest    | GET by same digest                                              | Pull by digest   | Image is single arch and/or image is multi-arch but some part of the image already exists on the local machine |
| HEAD           | digest    | GET by same digest, then a second GET by different digest       | Pull by Digest   | Image is multi-arch                                                                                            |

## Changes in data over time

The insights and analytics service is continuously improved to increase the
value it brings to publishers. Some changes might include adding new data
points, or improving existing data to make it more useful.

When there is a change in the dataset provided by the service, such a change
doesn't get retroactively applied. As new data points get added, they're
available from the point of introduction and going forward.

Refer to the tables in the [Data points](#data-points) section to see from which
date a given data point is available.

## Privacy

This section contains information about privacy-protecting measures that ensures
consumers of content on Docker Hub remain completely anonymous.

> **Important**
>
> Docker never shares any Personally Identifiable Information (PII) as part of
> analytics data. {: .important }

The summary dataset includes Unique IP address count. This data point only
includes the number of distinct unique IP addresses that request an image.
Individual IP addresses are never shared.

The raw dataset includes user IP domains as a data point. That's the domain name
of the company associated with an IP address that pulled an image (for example,
`docker.com`). This data point is only included where the IP type is `business`,
which only includes a small subset of all pull data. On average, only about 30%
of all pulls classify as the `business` IP type (this varies between publishers
and images).
