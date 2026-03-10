---
title: Set up single sign-on
linkTitle: Connect
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
and your identity provider (IdP). This guide walks you through adding and verifying your domain, setting up Docker and your IdP for an SSO connection, and finalizing and testing the connection. 

## Overview 

Docker supports any SAML 2.0 or OIDC-compatible identity provider. While this guide
focuses on detailed set-up instructions for Okta and Microsoft Entra ID, the general process remains the same for other IdPs. 

You will:

- Add and verify a domain.
- Configure the connection in Docker.
- Set up the application in your IdP with the values from Docker.
- Complete the connection by entering your IdP's values back into Docker.
- Test the connection.

## Prerequisites

Before you begin, make sure the following conditions are met:

- Notify your company about the upcoming SSO sign-in process.
- Ensure all users have Docker Desktop version 4.42 or later installed.
- Confirm that each Docker user has a valid IdP account using the same
email address as their Unique Primary Identifier (UPN).
- If you plan to [enforce SSO](/manuals/enterprise/security/single-sign-on/connect.md#optional-enforce-sso),
users accessing Docker through the CLI must [create a personal access token (PAT)](/docker-hub/access-tokens/). The PAT replaces their username and password for authentication.
- Ensure CI/CD pipelines use PATs or OATs instead of passwords.

## Set up an SSO connection 

> [!TIP]
> These procedures have you copy and paste values between Docker and your IdP. Complete this guide in one session with separate browser windows open for Docker and your IdP. 

### Step 1: Add a domain

To add a domain:

1. Sign in to [Docker Home](https://app.docker.com), then choose your
organization. If your organization is part of a company, then select the company to manage
the domain at the company level.
1. Select **Admin Console**, then **Domain management**.
2. Select **Add a domain**.
3. Enter your domain in the text box and select **Add domain**.
4. In the modal, copy the **TXT Record Value** provided for domain verification.

### Step 2: Verify your domain

To confirm domain ownership, add a TXT record to your Domain Name System (DNS)
host using the TXT Record Value from Docker. DNS propagation can take up to
72 hours. Docker automatically checks for the record during this time.

> [!TIP]
>
> When adding a record name, **use `@` or leave it empty** for root domains like `example.com`. **Avoid common values** like `docker`, `docker-verification`, `www`, or your domain name itself. Always **check your DNS provider's documentation** to verify their specific record name requirements.

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

### Step 3. Create an SSO connection in Docker

1. From [Docker Home](https://app.docker.com), choose your
organization and toggle the **Admin Console** dropdown. Select **SSO and SCIM** from the **Security** section. 
1. Select **Create Connection** and name the connection. Choose either **SAML** or **Azure AD (OIDC)** for your authentication method.
1. Copy the required values for your IdP and store these values in a text editor:
    - Okta SAML: **Entity ID**, **ACS URL**
    - Azure OIDC: **Redirect URL**

Keep this window open to paste values from your IdP later.

### Step 4. Create an SSO connection in your IdP

Use the following tabs based on your IdP provider.

{{< tabs >}}
{{< tab name="Okta SAML" >}}

To enable SSO with Okta, you need [super admin](https://help.okta.com/en-us/content/topics/security/administrators-super-admin.htm) permissions for the Okta org. 

1. Open the Admin portal from your Okta account and select **Administration**. 
1. Choose **Create App Integration** and select **SAML 2.0**. 
   - When prompted, name your app "Docker." 
   - You may upload a logo, but it's not required. 
1. Paste the values you copied from creating an SSO connection in Docker: 
    - For the **Single Sign On URL** value, paste the Docker ACS URL. 
    - For the **Audience URI (SP Entity ID)** value, paste the Docker Entity ID. 
1. Configure the following settings. These settings determine the primary identification method your IdP sends to Docker for verification:
    - Name ID format: `EmailAddress`
    - Application username: `Email`
    - Update application on: `Create and update`
1. Optional. Add [SAML attributes](/manuals/enterprise/security/provisioning/_index.md#sso-attributes), if required by your org. 
1. Select the **This is an internal app that we have created** checkbox before finishing. 

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

The following procedures reproduce instructions from Microsoft Learn documentation for [configuring an app service with OIDC](https://learn.microsoft.com/en-us/azure/app-service/configure-authentication-provider-openid-connect#-register-your-app-with-the-oidc-identity-provider). If you're uncertain, review the official Microsoft documentation and return here for the rest of the procedures.  

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

### Step 5. Connect Docker to your IdP

Complete the integration by pasting your IdP values into Docker.

> [!IMPORTANT]
> When prompted to copy a certificate, copy the entire certificate starting with 
> `----BEGIN CERTIFICATE----` and including the `----END CERTIFICATE----` lines.

{{< tabs >}}
{{< tab name="Okta SAML" >}}

1. In Okta, select your app and go to **View SAML setup instructions**.
1. Copy the **SAML Sign-in URL** and **x509 Certificate**, then return to the Docker Admin Console.
1. Paste the **SAML Sign-in URL** and **x509 Certificate** values.
1. Optional. Select a default team, if required by your org.
1. Review and select **Create connection**.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

1. Open your downloaded **Certificate (Base64)** in a text editor.
1. Copy the following values:
    - From Azure AD: **Login URL**
    - **Certificate (Base64)** contents
1. Return to the Docker Admin Console, then paste the **Login URL** and **Certificate (Base64)** values.
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

### Step 6. Test the connection

IdPs like Microsoft Entra and Okta may require that you assign a user to an application before testing SSO. You can review [Microsoft Entra](https://learn.microsoft.com/en-us/entra/identity/enterprise-apps/add-application-portal-setup-sso#test-single-sign-on)'s documentation and [Okta](https://help.okta.com/wf/en-us/content/topics/workflows/connector-reference/okta/actions/assignusertoapplicationforsso.htm)'s documentation to learn how to assign yourself or other users to an app.

After assigning yourself to an app: 

1. Open an incognito browser window and sign in to the Admin Console using your domain email address.
1. When redirected to your IdP's sign in page, authenticate with your domain email instead of using your Docker ID. 

If you have multiple IdPs, choose the sign-in option **Continue with SSO**. If you're using the CLI, you must authenticate using a personal access token.

## Configure multiple IdPs

Docker supports multiple identity provider (IdP) configurations by letting you associate one domain with more than one IdP. Each connection must use the same domain, which lets users choose their IdP when they select **Continue with SSO** at login. 

To add multiple IdPs:

1. Use the same domain for each connection. 
1. Repeat steps 3-6 from the [Set up an SSO connection] procedures on this page. Repeat these steps for each IdP your organization intends to use.

Because you must use the same domain for each IdP, you don't need to repeat steps 1 and 2. 

## Enforce SSO

If SSO is not enforced, users can still sign in using Docker usernames and passwords. Enforcing SSO requires users to use SSO when signing into Docker, which centralizes authentication and enforces policies set by the IdP.

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
- [Create personal access tokens](/manuals/enterprise/security/access-tokens.md).
- [Troubleshoot SSO](/manuals/enterprise/troubleshoot/troubleshoot-sso.md) issues.
