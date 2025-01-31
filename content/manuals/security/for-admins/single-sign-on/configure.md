---
description: Learn how to configure single sign-on for your organization or company.
keywords: configure, sso, docker hub, hub, docker admin, admin, security
title: Configure single sign-on
linkTitle: Configure
aliases:
- /docker-hub/domains/
- /docker-hub/sso-connection/
- /docker-hub/enforcing-sso/
- /single-sign-on/configure/
- /admin/company/settings/sso-configuration/
- /admin/organization/security-settings/sso-configuration/
---

{{< summary-bar feature_name="SSO" >}}

Get started creating a single sign-on (SSO) connection for your organization or company. This guide walks through the steps to add and verify the domains your members use to sign in to Docker.

## Step one: Add your domain

{{< summary-bar feature_name="Admin console early access" >}}

{{< tabs >}}
{{< tab name="Admin Console" >}}

1. Sign in to the [Admin Console](https://admin.docker.com/).
2. Select your organization or company from the **Choose profile** page. Note that when an organization is part of a company, you must select the company and configure the domain for the organization at the company level.
3. Under **Security and access**, select **Domain management**.
4. Select **Add a domain**.
5. Enter your domain in the text box and select **Add domain**.
6. The pop-up modal will prompt you with steps to verify your domain. Copy the **TXT Record Value**.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

1. Sign in to [Docker Hub](https://hub.docker.com/).
2. Select **Organizations** and then your organization from the list.
3. On your organization page, select **Settings** and then **Security**.
4. Select **Add a domain**.
5. Enter your domain in the text box and select **Add domain**.
6. The pop-up modal will prompt you with steps to verify your domain. Copy the **TXT Record Value**.

{{< /tab >}}
{{< /tabs >}}

## Step two: Verify your domain

{{< summary-bar feature_name="Admin console early access" >}}

Verifying your domain ensures Docker knows you own it. Domain verification is done by adding your Docker TXT Record Value to your domain host. The TXT Record Value proves ownership, which signals the Domain Name System (DNS) to add this record. It can take up to 72 hours for DNS to recognize the change. When the change is reflected in DNS, Docker will automatically check the record to confirm your ownership.

{{< tabs >}}
{{< tab name="Admin Console" >}}

1. Navigate to your domain host, create a new TXT record, and paste the **TXT Record Value** from Docker.
2. TXT record verification can take 72 hours. Once you have waited for TXT record verification, return to the **Domain management** page of the Admin Console and select **Verify** next to your domain name.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

1. Navigate to your domain host, create a new TXT record, and paste the **TXT Record Value** from Docker.
2. TXT Record Verification can take 72 hours. Once you have waited for TXT record verification, return to the **Security** page of Docker Hub and select **Verify** next to your domain name.

{{< /tab >}}
{{< /tabs >}}

Once you have added and verified your domain, you are ready to create an SSO connection between Docker and your identity provider (IdP).

## More resources

The following videos walk through verifying your domain to create your SSO connection in Docker.

- [Video: Verify your domain for SSO with Okta](https://youtu.be/c56YECO4YP4?feature=shared&t=529)
- [Video: Verify your domain for SSO with Azure AD (OIDC)](https://youtu.be/bGquA8qR9jU?feature=shared&t=496)

## What's next?

[Connect Docker and your IdP](../single-sign-on/connect.md).

