---
title: Set up single sign-on
linkTitle: Setup
weight: 20
description: Connect Docker and your identity provider, test the setup, and enable enforcement
keywords: configure sso, set up sso, docker sso setup, docker identity provider, sso enforcement, docker hub, security
aliases:
  - /security/for-admins/single-sign-on/connect/
  - /docker-hub/domains/
  - /docker-hub/sso-connection/
  - /docker-hub/enforcing-sso/
  - /single-sign-on/configure/
  - /admin/company/settings/sso-configuration/
  - /admin/organization/security-settings/sso-configuration/
  - /security/for-admins/single-sign-on/configure/
  - /enterprise/security/single-sign-on/configure
---

{{< summary-bar feature_name="SSO" >}}

To set up a single sign-on (SSO), you need to establish a connection between Docker
and your identity provider (IdP). While this guide
uses Okta and Microsoft Entra ID as a working example, the general process remains the same for other IdPs.

If you're unfamiliar with the SSO process, first review [SSO overview](/manuals/enterprise/security/single-sign-on/_index.md) to learn about how SSO works.

## Prerequisites

Docker supports any SAML 2.0 or OIDC-compatible identity provider. Before you begin, make sure the following conditions are met:

- Notify your company about the upcoming SSO sign-in process.
- Confirm that each Docker user has a valid IdP account using the same
  email address as their Unique Primary Identifier (UPN).
- Ensure CI/CD pipelines use PATs or OATs instead of passwords.

## Set up an SSO connection

> [!TIP]
> These procedures have you copy and paste values between Docker and your IdP. Complete this guide in one session with separate browser windows open for Docker and your IdP.

### Step 1: Add a domain

To add a domain:

1. Sign in to [app.docker.com](https://app.docker.com), then choose your
   organization. If your organization is part of a company, then select the company to manage
   the domain at the company level.
1. Select **Admin Console**, then **Domain management**.
1. Select **Add a domain**.
1. Enter your domain in the text box and select **Add domain**.
1. In the modal, copy the **TXT Record Value** provided for domain verification.

### Step 2: Verify your domain

To confirm domain ownership, add a TXT record to your Domain Name System (DNS)
host using the TXT Record Value from Docker. DNS propagation can take up to
72 hours. Docker automatically checks for the record during this time.

> [!TIP]
>
> When adding a record name, use `@` or leave it empty for root domains like `example.com`. Avoid common values like `docker`, `docker-verification`, `www`, or your domain name itself. Always check your DNS provider's documentation to verify their specific record name requirements.

{{< tabs >}}
{{< tab name="AWS Route 53" >}}

1. To add your TXT record to AWS, see [Creating records by using the Amazon Route 53 console](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/resource-record-sets-creating.html).
1. Wait up to 72 hours for TXT record verification.
1. After the record is live, go to **Domain management** in the [Admin Console](https://app.docker.com/admin) and select **Verify**.

{{< /tab >}}
{{< tab name="Google Cloud DNS" >}}

1. To add your TXT record to Google Cloud DNS, see [Verifying your domain with a TXT record](https://cloud.google.com/identity/docs/verify-domain-txt).
1. Wait up to 72 hours for TXT record verification.
1. After the record is live, go to **Domain management** in the [Admin Console](https://app.docker.com/admin) and select **Verify**.

{{< /tab >}}
{{< tab name="GoDaddy" >}}

1. To add your TXT record to GoDaddy, see [Add a TXT record](https://www.godaddy.com/help/add-a-txt-record-19232).
1. Wait up to 72 hours for TXT record verification.
1. After the record is live, go to **Domain management** in the [Admin Console](https://app.docker.com/admin) and select **Verify**.

{{< /tab >}}
{{< tab name="Other providers" >}}

1. Sign in to your domain host.
1. Add a TXT record to your DNS settings and save the record.
1. Wait up to 72 hours for TXT record verification.
1. After the record is live, go to **Domain management** in the [Admin Console](https://app.docker.com/admin) and select **Verify**.

{{< /tab >}}
{{< /tabs >}}

### Step 3: Create an SSO connection in Docker

1. From [app.docker.com](https://app.docker.com), choose your
   organization. Select **Admin Console**.
1. Choose **SSO and SCIM** from the **Security** section.
1. Select **Create Connection**. Name the connection, then choose **SAML 2.0**.
1. Keep this window open. You need to copy and paste these values in your Okta window:
   - **Entity ID**
   - **ACS URL**
   - **Connection ID**

You will return here to finish connecting after creating an SSO connection in your IdP.

### Step 4: Create an SSO connection in your IdP

Use the following tabs based on your IdP provider.

{{< tabs >}}
{{< tab name="Okta SAML" >}}

You need [super admin permissions](https://help.okta.com/en-us/content/topics/security/administrators-super-admin.htm) for the Okta org.

1. Sign in to your Okta Admin account. From the top nav, select the **Admin** button to go to Okta's Admin Console.
1. From the **Applications** section in the left nav, choose **Applications**. Select **Create App Integration**.
1. Choose SAML 2.0 to match your selection from Docker's Admin Console.
1. For **General Settings**, name your app "Docker." Uploading a logo is optional.
1. For **Configure SAML**, input the following values:
   - For the **Single Sign On URL** value, paste the Docker ACS URL.
   - For the **Audience URI (SP Entity ID)** value, paste the Docker Entity ID.
   - For **Name ID format**, choose `EmailAddress`
   - For **Application username**, choose `Email`
   - For **Update application username on**, choose `Create and update`
   - Optional. Add [SAML attributes](/manuals/enterprise/security/provisioning/_index.md#sso-attributes), if required by your org.
1. For **Feedback**, choose **This is an internal app that we have created** checkbox before finishing.

Keep your Okta window open for the next step.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

To enable SSO with Microsoft Entra, you need [Cloud Application Administrator](https://learn.microsoft.com/en-us/entra/identity/role-based-access-control/permissions-reference#cloud-application-administrator) permissions.

1. From Microsoft Entra admin center, select **Entra ID**, then go to **Enterprise apps**. Select **All applications**.
1. Choose **Create your own application** and name your app "Docker". Select **Non-gallery**.
1. After creating your app, go to **Single Sign-On** and select **SAML**.
1. Select **Edit** on the **Basic SAML configuration** section. From **Basic SAML configuration**, choose **Edit** and paste the values you copied from creating an SSO connection in Docker:
   - For the **Identifier** value, paste the Docker Entity ID.
   - For the **Reply URL** value, paste Docker ACS URL.
1. Optional. Add [SAML attributes](/manuals/enterprise/security/provisioning/_index.md#sso-attributes), if required by your org.
1. From the **SAML Signing Certificate** section, download your **Certificate (Base64)**.

{{< /tab >}}
{{< tab name="Azure OpenID Connect (OIDC)" >}}

#### Register the app

1. Sign in to [Microsoft Entra admin center](https://entra.microsoft.com/).
1. Go to **App Registration** and select **New Registration**.
1. Name the application "Docker".
1. Set account types and paste the **Redirect URI** from Docker.
1. Select **Register**.
1. Copy the **Client ID**.

#### Create client secrets

1. In your app, go to **Certificates & secrets**.
1. Select **New client secret**, describe and configure duration, then **Add**.
1. Copy the **value** of the new secret.

#### Set API permissions

1. In your app, go to **API permissions**.
1. Select **Grant admin consent** and confirm.
1. Select **Add a permissions** > **Delegated permissions**.
1. Search and select `User.Read`.
1. Confirm that admin consent is granted.

{{< /tab >}}
{{< /tabs >}}

### Step 5: Connect Docker to your IdP

Complete the integration by pasting your IdP values into Docker.

> [!IMPORTANT]
> When prompted to copy a certificate, copy the entire certificate starting with
> `----BEGIN CERTIFICATE----` and including the `----END CERTIFICATE----` lines.

{{< tabs >}}
{{< tab name="Okta SAML" >}}

1. Go to **Applications** and choose **Applications**. Choose your app from the **ACTIVE** table.
1. From **Sign on**, go to **View SAML setup instructions**. This page contains the **SAML Sign-in URL** and **x509 Certificate**. Keep this page open.
1. Return to your opened Docker window for the **Create single sign-on connection** step. Paste in the **SAML Sign-in URL** and **x509 Certificate** values.
1. Optional. Select a default team, if required by your org.
1. Review and select **Create connection**.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

1. Open your downloaded **Certificate (Base64)** in a text editor.
1. Copy the following values:
   - From Azure AD: **Login URL**
   - **Certificate (Base64)** contents
1. Return to the Docker Admin Console, then paste the **Login URL** and **Certificate (Base64)** values.
1. Choose your domain from the drop-down.
1. Optional. Select a default team, if required by your org.
1. Review and select **Create connection**.

{{< /tab >}}
{{< tab name="Azure OpenID Connect (OIDC)" >}}

1. Return to the Docker Admin Console.
1. Paste the following values:
   - **Client ID**
   - **Client Secret**
   - **Azure AD Domain**
1. Optional. Select a default team, if required by your org.
1. Review and select **Create connection**.

{{< /tab >}}
{{< /tabs >}}

### Step 6: Test the connection

IdPs like Microsoft Entra and Okta may require that you assign a user to an application before testing SSO. You can review [Microsoft Entra](https://learn.microsoft.com/en-us/entra/identity/enterprise-apps/add-application-portal-setup-sso#test-single-sign-on)'s documentation and [Okta](https://help.okta.com/wf/en-us/content/topics/workflows/connector-reference/okta/actions/assignusertoapplicationforsso.htm)'s documentation to learn how to assign yourself or other users to an app.

After assigning yourself to an app:

1. Open an incognito browser window and sign in to the Admin Console using your domain email address.
1. When redirected to your IdP's sign in page, authenticate with your domain email instead of using your Docker ID.

If you have multiple IdPs, choose the sign-in option **Continue with SSO**. If you're using the CLI, you must authenticate using a personal access token.

## Configure multiple IdPs

Docker supports multiple identity provider (IdP) configurations by letting you associate one domain with more than one IdP. Each connection must use the same domain, which lets users choose their IdP when they select **Continue with SSO** at sign in.

To add multiple IdPs:

1. Use the same domain for each connection.
1. Repeat steps 3-6 from the [Set up an SSO connection](/manuals/enterprise/security/single-sign-on/connect.md#set-up-an-sso-connection) procedures on this page. Repeat these steps for each IdP your organization intends to use.

Because you must use the same domain for each IdP, you won't need to repeat steps to add and verify your domains.

## Enforce SSO

If SSO is not enforced, users can still sign in using Docker usernames and passwords. Enforcing SSO requires users to use SSO when signing into Docker, which centralizes authentication and enforces policies set by the IdP.

Before enforcing SSO, users accessing Docker through the CLI must [create a personal access token (PAT)](/manuals/security/access-tokens.md). The PAT replaces their username and password for authentication.

1. Sign in to [Docker Home](https://app.docker.com/) and select
   your organization or company.
1. Select **Admin Console**, then **SSO and SCIM**.
1. In the SSO connections table, select the **Action** menu, then **Enable enforcement**.
1. Follow the on-screen instructions.
1. Select **Turn on enforcement**.

When you enforce SSO, your users cannot modify their email address and
password, convert a user account to an organization, or set up 2FA through
Docker Hub. If you want to use 2FA, you must enable 2FA through your IdP.

## Next steps

- [Provision users](/manuals/enterprise/security/provisioning/_index.md).
- [Enforce sign-in](../enforce-sign-in/_index.md).
- [Create personal access tokens](/manuals/security/access-tokens.md).
- [Troubleshoot SSO](/manuals/enterprise/security/single-sign-on/troubleshoot-sso.md) issues.
