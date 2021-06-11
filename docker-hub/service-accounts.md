---
description: Docker Service accounts
keywords: Docker, service, accounts, Docker Hub
title: Service accounts
---

A service account is a Docker ID used for automated management of container images or containerized applications. Service accounts are typically used in automated workflows, and do not share Docker IDs with the members in the Team plan. Common use cases for service accounts include mirroring content on Docker Hub, or tying in image pulls from your CI/CD process.

> **Note**
>
> Service accounts included with the Team plan are limited to 15,000 pulls per day. If you require a higher number of pulls, you can purchase an Enhanced Service Account add-on.

## Enhanced Service Account add-on pricing

Refer to the following table for details on the Enhanced Service Account add-on pricing:

| Tier | Pull Rates Per Day* | Annual Fee |
| ------ | ------ | ------ |
| 1 | 5-10k | $9,950/yr |
| 2 | 10-25k | $17,950/yr |
| 3 | 25k-50k | $32,950/yr |
| 4 | 50-100k | $58,950/yr |
| 5 | 100k | [Contact Sales](https://www.docker.com/pricing/questions){:target="_blank" rel="noopener" class="_"} |

<sub>*Once the initial Tier is established, that is the minimum fee for the year.  Annual commitment required.  The service account may exceed Pulls by up to 25% for up to 20 days during the year without incurring additional fees.  Reports on consumption will be provided upon request.  At the end of the initial 1-year term, the appropriate Tier will be established for the following year.<sub>

## How a pull is defined

- A pull request is defined as up to two `GET` requests on registry manifest URLs (`/v2/*/manifests/*`).
- A normal image pull makes a single manifest request.
- A pull request for a multi-arch image makes two manifest requests.
- `HEAD` requests are not counted.
- Limits are applied based on the user doing the pull, and not based on the image being pulled or its owner.

## Creating a new service account

To create a new service account for your Team account:

1. Create a new Docker ID.
2. Create a [team](orgs.md#create-a-team) in your organization and grant it read-only access to your private repositories.
3. Add the new Docker ID to your [organization](orgs.md#working-with-organizations).
4. Add the new Docker ID  to the [team](orgs.md#add-a-member-to-a-team) you created earlier.
5. Create a new [personal access token (PAT)](/access-tokens.md) from the user account and use it for CI.

> **Note**
>
> If you want a read-only PAT just for your open-source repos, or to access
official images and other public images, you do not have to grant any access permissions to the new Docker ID.

## Additional information

Refer to the following topics for additional information:

- [Mirroring Docker Hub](../registry/recipes/mirror.md)
- [Resource Consumption Updates FAQ](https://www.docker.com/pricing/resource-consumption-updates){:target="_blank" rel="noopener" class="_"}
