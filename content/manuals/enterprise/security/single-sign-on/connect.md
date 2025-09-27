---
title: Connect single sign-on
linkTitle: Connect
description: Connect Docker and your identity provider, test the setup, and enable enforcement
keywords: configure sso, set up sso, docker sso setup, docker identity provider, sso enforcement, docker hub, security
aliases:
 - /security/for-admins/single-sign-on/connect/
---

{{< summary-bar feature_name="SSO" >}}

Setting up a single sign-on (SSO) connection involves configuring both Docker
and your identity provider (IdP). This guide walks you through setup
in Docker, setup in your IdP, and final connection.

> [!TIP]
>
> You’ll copy and paste values between Docker and your IdP. Complete this guide
in one session with separate browser windows open for Docker and your IdP.

## Supported identity providers

Docker supports any SAML 2.0 or OIDC-compatible identity provider. This guide
provides detailed setup instructions for the most commonly
used providers: Okta and Microsoft Entra ID.

If you're using a
different IdP, the general process remains the same:

1. Configure the connection in Docker.
1. Set up the application in your IdP using the values from Docker.
1. Complete the connection by entering your IdP's values back into Docker.
1. Test the connection.

## Prerequisites

Before you begin:

- Verify your domain
- Set up an account with your identity provider (IdP)
- Complete the steps in the [Configure single sign-on](configure.md) guide

## Step one: Create an SSO connection in Docker

> [!NOTE]
>
> You must [verify at least one domain](/manuals/enterprise/security/single-sign-on/configure.md) before creating an SSO connection.

1. Sign in to [Docker Home](https://app.docker.com) and choose your
organization.
1. Select **Admin Console**, then **SSO and SCIM**.
1. Select **Create Connection** and provide a name for the connection.
1. Select an authentication method: **SAML** or **Azure AD (OIDC)**.
1. Copy the required values for your IdP:
    - Okta SAML: **Entity ID**, **ACS URL**
    - Azure OIDC: **Redirect URL**

Keep this window open to paste values from your IdP later.

## Step two: Create an SSO connection in your IdP

Use the following tabs based on your IdP provider.

{{< tabs >}}
{{< tab name="Okta SAML" >}}

1. Sign in to your Okta account and open the Admin portal.
1. Select **Administration** and then **Create App Integration**.
1. Select **SAML 2.0**, then **Next**.
1. Name your app "Docker".
1. Optional. Upload a logo.
1. Paste values from Docker:
    - Docker ACS URL -> **Single Sign On URL**
    - Docker Entity ID -> **Audience URI (SP Entity ID)**
1. Configure the following settings:
    - Name ID format: `EmailAddress`
    - Application username: `Email`
    - Update application on: `Create and update`
1. Optional. Add SAML attributes. See [SSO attributes](/manuals/enterprise/security/provisioning/_index.md#sso-attributes).
1. Select **Next**.
1. Select the **This is an internal app that we have created** checkbox.
1. Select **Finish**.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

1. Sign in to Microsoft Entra (formerly Azure AD).
1. Select **Default Directory** > **Add** > **Enterprise Application**.
1. Choose **Create your own application**, name it "Docker", and choose **Non-gallery**.
1. After creating your app, go to **Single Sign-On** and select **SAML**.
1. Select **Edit** on the **Basic SAML configuration** section.
1. Edit **Basic SAML configuration** and paste values from Docker:
    - Docker Entity ID -> **Identifier**
    - Docker ACS URL -> **Reply URL**
1. Optional. Add SAML attributes. See [SSO attributes](/manuals/enterprise/security/provisioning/_index.md#sso-attributes).
1. Save the configuration.
1. From the **SAML Signing Certificate** section, download your **Certificate (Base64)**.

{{< /tab >}}
{{< tab name="Azure Connect (OIDC)" >}}

### Register the app

1. Sign in to Microsoft Entra (formerly Azure AD).
1. Select **App Registration** > **New Registration**.
1. Name the application "Docker".
1. Set account types and paste the **Redirect URI** from Docker.
1. Select **Register**.
1. Copy the **Client ID**.

### Create client secrets

1. In your app, go to **Certificates & secrets**.
1. Select **New client secret**, describe and configure duration, then **Add**.
1. Copy the **value** of the new secret.

### Set API permissions

1. In your app, go to **API permissions**.
1. Select **Grant admin consent** and confirm.
1. Select **Add a permissions** > **Delegated permissions**.
1. Search and select `User.Read`.
1. Confirm that admin consent is granted.

{{< /tab >}}
{{< /tabs >}}

## Step three: Connect Docker to your IdP

Complete the integration by pasting your IdP values into Docker.

{{< tabs >}}
{{< tab name="Okta SAML" >}}

1. In Okta, select your app and go to **View SAML setup instructions**.
1. Copy the **SAML Sign-in URL** and **x509 Certificate**.

    > [!IMPORTANT]
    >
    > Copy the entire certificate, including `----BEGIN CERTIFICATE----` and `----END CERTIFICATE----` lines.
1. Return to the Docker Admin Console.
1. Paste the **SAML Sign-in URL** and **x509 Certificate** values.
1. Optional. Select a default team.
1. Review and select **Create connection**.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

1. Open your downloaded **Certificate (Base64)** in a text editor.
1. Copy the following values:
    - From Azure AD: **Login URL**
    - **Certificate (Base64)** contents

    > [!IMPORTANT]
    >
    > Copy the entire certificate, including `----BEGIN CERTIFICATE----` and `----END CERTIFICATE----` lines.
1. Return to the Docker Admin Console.
1. Paste the **Login URL** and **Certificate (Base64)** values.
1. Optional. Select a default team.
1. Review and select **Create connection**.

{{< /tab >}}
{{< tab name="Azure Connect (OIDC)" >}}

1. Return to the Docker Admin Console.
1. Paste the following values:
    - **Client ID**
    - **Client Secret**
    - **Azure AD Domain**
1. Optional. Select a default team.
1. Review and select **Create connection**.

{{< /tab >}}
{{< /tabs >}}

## Step four: Test the connection

1. Open an incognito browser window.
1. Sign in to the Admin Console using your **domain email address**.
1. The browser will redirect to your identity provider's sign in page to authenticate. If you have [multiple IdPs](#optional-configure-multiple-idps), choose the sign sign-in option **Continue with SSO**.
1. Authenticate through your domain email instead of using your Docker ID.

If you're using the CLI, you must authenticate using a personal access token.

## Optional: Configure multiple IdPs

Docker supports multiple IdP configurations. To use multiple IdPs with one domain:

- Repeat Steps 1-4 on this page for each IdP.
- Each connection must use the same domain.
- Users will select **Continue with SSO** to choose their IdP at sign in.

## Optional: Enforce SSO

> [!IMPORTANT]
>
> If SSO is not enforced, users can still sign in using Docker usernames and passwords.

Enforcing SSO requires users to use SSO when signing into Docker. This centralizes authentication and enforces policies set by the IdP.

1. Sign in to [Docker Home](https://app.docker.com/) and select
your organization or company.
1. Select **Admin Console**, then **SSO and SCIM**.
1. In the SSO connections table, select the **Action** menu, then **Enable enforcement**.
1. Follow the on-screen instructions.
1. Select **Turn on enforcement**.

When SSO is enforced, your users are unable to modify their email address and
password, convert a user account to an organization, or set up 2FA through
Docker Hub. If you want to use 2FA, you must enable 2FA through your IdP.

## Next steps

- [Provision users](/manuals/enterprise/security/provisioning/_index.md).
- [Enforce sign-in](../enforce-sign-in/_index.md).
- [Create personal access tokens](/manuals/enterprise/security/access-tokens.md).
- [Troubleshoot SSO](/manuals/enterprise/troubleshoot/troubleshoot-sso.md) issues.
