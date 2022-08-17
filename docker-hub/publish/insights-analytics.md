---
title: Insights & analytics
description: Provides an API analytics tool for your products.
keywords: docker hub, hub, insights, analytics, api
---

The Insights & Analytics tool includes an API that provides analytics for your products. With this API, you now have self-serve access to metrics where you can download both raw data and summary data. You can view how many times your images have been pulled by tag or by digest, and get breakouts by geolocation, cloud provider, and tool.The API measures the number of 'docker pulls' that were triggered by users, and automated machines such as CICD systems. You can import this weekly or monthly data into your own systems or you can analyze it manually using spreadsheet reports. The data can then be exported as a downloadable GZipped CSV file. Review the Data Definitions section for more information.

There are automated systems that frequently check for new versions of your images, and we have tried to distinguish between those 'version checks' and real user downloads by inspecting the order and timing of API calls from the same IP address. This will provide more insight into customer behavior. You can inspect the rules we have applied in the Action Classification Rule Set section of this document. We welcome [feedback](https://forms.gle/nb7beTUQz9wzXy1b6)(: target="_blank" rel="noopener" class="_") and questions on these rules.

## Raw Data

In your raw data CSV you will have access to the data points listed below. You can request raw data for a complete week (Monday to Monday) or for a complete month (available on the 1st day of the following month). Please note that each Action is represented as a single row.

  * Type (industry)
  * Host (cloud provider)
  * Country (geolocation)
  * Timestamp
  * Namespace
  * Repo
  * Reference (digest will always be included, we provide tag when available)
  * HTTP Method
  * Action - can filter on the below
  * Pull by tag
  * Pull by digest
  * Version check
  * User-Agent Tool

## Summary Data

In your summary data CSV, you will have access to the data points listed below. You can request summary data for a complete week (Monday to Monday) or for a complete month (available on the 1st day of the following month).

For For every namespace, repo, and reference (tag or digest):

  * Unique IP address count
  * Pulls by tag count
  * Pulls by digest count
  * Version Check count

## Data Definitions

| Data Point | Definition |
|:-----|:--------|
| Action | An Action represents the multiple request events associated with a docker pull. We have applied rules to these events so that the data is more meaningful in analyzing user behavior and intent. An Action can be filtered into three distinct categories: version check, pull by tag, and pull by digest. Each Action is represented as a single row in the raw data. |
| Version Check  |  This is a filter on the Action data point. It is a speculation of user intent.                                                    * Includes: HEAD by tag not followed by a GET (by the same IP address within 5-second window) * Excludes: HEAD by digest|
| Pull by Tag | This is a filter on the Action data point. It is a speculation of user intent. Includes: GET (by digest or by tag)If the GET is immediately preceded by a HEAD by tag (by the same IP address within a 5-second window), then the GET and HEAD together are counted as 1 pull by tag. If the GET by tag is immediately followed by another GET (by the same IP address within a 5-second window, but a different digest), then both GETs together are counted as 1 pull by tag.|
|Pull by Digest  |This is a filter on the Action data point. It is a speculation of user intent. Includes: GET by digest. If the GET is immediately preceded by a HEAD by digest (by the same IP address within a 5-second window), then the GET and HEAD together are counted as 1 pull by digest. If the GET is immediately followed by another GET (same IP address within a 5-second window, but a different digest), then both GETs together are counted as 1 pull by digest. Includes: HEAD by digest, not followed by a GET|
|Type | The industry from which the event originates. Industry types include “business”, “ISP” (internet service provider), “hosting”, “education”, and “null” in cases where the industry could not be identified.|
|Host |The cloud service provider used in an event.|
|Country |The country from which the request originated.  |
|Timestamp |Date & time of an event in the following schema: YYYY-MM-DD 00:00:00|
|Namespace |The Docker organization that a repository is a part of. |
|Repo |The repository that an image belongs to.|
|Reference |The tag or digest of any given image.|
|HTTP Method | The HTTP method(s) used in a request by the client. More information on Docker Registry HTTP API protocols can be found here.|
|User-Agent Tool|The application a user used to pull an image.  Extracted from the UA string.|
|Unique IP Address|As part of our privacy-preserving policy, Docker only shares the count of distinct unique IP addresses that request an image.|



## Action Classification Rule Set

|Starting Event | Reference | Followed By | Resulting Action | Use Case(s) |Notes|
|:-----|:--------|:-----|:--------|:-----|:--------|
| HEAD | tag | N/A | Version Check | User already has all layers existing on local machine|This is similar to the use case of a pull by tag when the user already has all the image layers existing locally, however, we are able to differentiate the user intent and classify accordingly.|
| Get| tag | N/A |Pull by Tag | User already has all layers existing on local machine and / or Image is single arch |
| Get | tag | Get by different digest | Pull by tag | Image is multiarch |Following GET by digests must be different from the first digest|
| HEAD | tag | GET by same digest | Pull by Tag | Image is multi-arch but some part of the image already exists on the local machine and/or User already has all layers existing on local machine|The HEAD by tag will send the most current digest, the following GET must be by that same digest. There may be an additional GET, if the image is multi arch (use case on the next line). If the user doesn't want the most current digest, then the user would perform HEAD by digest.|
| HEAD| tag |GET by same digest, THEN a second GET by different digest| Pull by Tag| Image is multi-arch |The HEAD by tag will send the most current digest, the following GET must be by that same digest. Since the image is multi arch, there is a second GET by a different digest. If the user doesn't want the most current digest, then the user would perform HEAD by digest.|
| HEAD | tag | GET by same digest, THEN a second GET by different digest | Pull by Tag| Image is multi arch | The HEAD by tag will send the most current digest, the following GET must be by that same digest. Since the image is multi-arch, there is a second GET by a different digest. If the user doesn't want the most current digest, then the user would perform HEAD by digest.|
| GET | digest | N/A | Pull by Digest |  User already has all layers existing on local machine and/or Image is single arch |
| HEAD | digest | N/A | Pull by Digest | User already has all layers existing on local machine |
| GET | digest | GET by different digest | Pull by Digest | Image is multi arch |The following GET by digest must be different from the first digest|
| HEAD | digest | GET by same digest | Pull by digest |  Image is single arch and/or image is multi-arch but some part of the image already exists on the local machine|
| Head | digest | GET by same digest, tHEN a second GET by different digest | Pull by Digest | Image is multi arch |

