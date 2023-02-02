---
description: manage companies
keywords: company, multiple organizations, manage companies
title: Overview
---

> **Note**
>
> The following features are only available to [Early Access](../release-lifecycle.md/#early-access-ea) participants.


To simplify the management of Docker organizations and settings, Docker has introduced a new view that provides a single point of visibility across multiple organizations called a Company. A company can become a parent to nested child organizations. A company lets Docker Business subscribers manage their organizations and configure settings centrally. With the new company owner role, you can control access to the company and company settings. These settings can affect all the organizations nested under the company. You can assign up to ten unique users to a company owner role without occupying a purchased seat.

Docker will work with your current Docker organization owners to create the company, associate your Docker Business organizations, and identify your company owner(s). Once created, users with a company owner role can navigate to a new page that displays the company name, organizations associated with the company, a list of company owners, and settings that include your Domain verification, Single Sign-on (SSO) connection to your identity provider, System for Cross-domain Identity Management (SCIM) setup.


 ![company-process](images/company-process-diagram.png){: width="700px" }

When a company owner makes adjustments to user management settings at the company level, this will affect all organizations associated with the company.

The company owner can:

- View all nested organizations.
- Configure SSO and SCIM for all nested organizations, including SCIM Group mapping.
- Enforce SSO log-in for all users in the company.
- Verify a domain separately from the organization namespace.
- Add and remove up to 10 company owners.

A company owner role is only available if your organization has a Docker Business subscription. If you don't have a Docker Business subscription, you must first [upgrade your subscription](../subscription/upgrade.md).

## Get started

You’ll need to send the following information to your CSM Docker team member to set up your company:

- The name of your company. For example, Docker uses the company name **dockerinc**.
- The organizations that you want to associate with the new company.
- The verified domains you want to move to the company level.
- Confirm if you want to migrate one of your organization’s SSO and SCIM settings to the company. Migrating SSO settings will also migrate verified domains from the organization to the parent company.

## Company overview and settings

To navigate to the company page:

1. Sign in to [Docker Hub](https://hub.docker.com/){: target="_blank" rel="noopener" class="_"} to view your company and organizations.
2. On the **Organizations** page, select your company to access the **Overview** tab. For example, the company listed below is **dockerinc** and the organization is **docker**.

    ![org-page](images/org-page.png){: width="700px" }
