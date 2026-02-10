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
and your identity provider (IdP). This guide walks you through set-up
in Docker, set-up in your IdP, and final connection. 

## Prerequisites

Before you begin:

- Verify your domain. You must [verify at least one domain](/manuals/enterprise/security/single-sign-on/configure.md) before creating an SSO connection.
- Set up an account with your identity provider (IdP).
- Complete the steps in the [Configure single sign-on](configure.md) guide.

## Set up SSO for Docker 

Docker supports any SAML 2.0 or OIDC-compatible identity provider. This guide
provides detailed setup instructions for the most commonly
used providers: Okta and Microsoft Entra ID. If you're using a
different IdP, the general process remains the same:

- Configure the connection in Docker.
- Set up the application in your IdP using the values from Docker.
-  Complete the connection by entering your IdP's values back into Docker.
- Test the connection.

These procedures prompt you to navigate between Docker docs and IdP docs. You will also need to copy and paste values 
between Docker and your IdP. Complete this guide in one session with separate browser windows open for Docker and your IdP. 

### 1. Create an SSO connection in Docker

1. From [Docker Home](https://app.docker.com), choose your
organization and toggle the **Admin Console** dropdown. Select **SSO and SCIM** from the **Security** section. 
1. Select **Create Connection** and name the connection. Choose either **SAML** or **Azure AD (OIDC)** for your authentication method.
1. Copy the required values for your IdP and store these values in a text editor:
    - Okta SAML: **Entity ID**, **ACS URL**
    - Azure OIDC: **Redirect URL**

Keep this window open to paste values from your IdP later.

### 2. Create an SSO connection in your IdP

Use the following tabs based on your IdP provider.

{{< tabs >}}
{{< tab name="Okta SAML" >}}

To enable SSO with Okta, you need [super admin]() permissions for the Okta org. 

1. Open the Admin portal from your Okta account and select **Administration**. 
1. Choose **Create App Integration** and select **SAML 2.0**. 
   - When prompted, name your app "Docker." 
   - You may upload a logo, but it's required. 
1. Paste the values you copied from creating an SSO connection in Docker: 
    - For the **Single Sign On URL** value, paste the Docker ACS URL. 
    - For the **Audience URI (SP Entity ID)** value, paste the Docker Entity ID. 
1. Configure the following settings. These settings determine the primary identification method your IdP sends to Docker for verification:
    - Name ID format: `EmailAddress`
    - Application username: `Email`
    - Update application on: `Create and update`
1. Optional. Add [SAML attributes](/manuals/enterprise/security/provisioning/_index.md#sso-attributes), if required by your org. 
1.  Select the **This is an internal app that we have created** checkbox before finishing. 

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

To enable SSO with Microsoft Entra, you need [Cloud Application Administrator](https://learn.microsoft.com/en-us/entra/identity/role-based-access-control/permissions-reference#cloud-application-administrator) permissions. 

1. From Microsoft Entra admin center, select **Entra ID**, then go to **Enterprise apps**. Select **All applications**.  
2. Choose **Create your own application** and name your app "Docker". Select **Non-gallery**.
3. After creating your app, go to **Single Sign-On** and select **SAML**.
4. Select **Edit** on the **Basic SAML configuration** section. From **Basic SAML configuration**, choose **Edit** and paste the values you copied from creating an SSO connection in Docker:
    - For the **Identifier** value, paste the Docker Entity ID.
    - For the **Reply URL** value, paste Docker ACS URL. 
5. Optional. Add [SAML attributes](/manuals/enterprise/security/provisioning/_index.md#sso-attributes), if required by your org. 
6. From the **SAML Signing Certificate** section, download your **Certificate (Base64)**.

{{< /tab >}}
{{< tab name="Azure Connect (OIDC)" >}}

#### Register the app

1. Sign in to Microsoft Entra (formerly Azure AD).
1. Select **App Registration** > **New Registration**.
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

### 3. Connect Docker to your IdP

Complete the integration by pasting your IdP values into Docker.

    > [!IMPORTANT]
    >
    > When prompted to copy a certificate, copy the entire certificate 
    > starting with `----BEGIN CERTIFICATE----` and including the `----END CERTIFICATE----` lines.

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
{{< tab name="Azure Connect (OIDC)" >}}

1. Return to the Docker Admin Console.
1. Paste the following values:
    - **Client ID**
    - **Client Secret**
    - **Azure AD Domain**
1. Optional. Select a default team, if required by your org.
1. Review and select **Create connection**.

{{< /tab >}}
{{< /tabs >}}

### 4. Test the connection

IdPs like Microsoft Entra and Okta may require that you assign a user to an application before testing SSO. You can review [Microsoft Entra](https://learn.microsoft.com/en-us/entra/identity/enterprise-apps/add-application-portal-setup-sso#test-single-sign-on)'s documentation and [Okta](https://help.okta.com/wf/en-us/content/topics/workflows/connector-reference/okta/actions/assignusertoapplicationforsso.htm)'s documentation to learn how to assign yourself or other users to an app.

After assigning yourself to an app: 

1. Open an incognito browser window and sign in to the Admin Console using your domain email address.
2. When redirected to your IdP's sign in page, authenticate with your domain email instead of using your Docker ID. 

If you have [multiple IdPs](#optional-configure-multiple-idps), choose the sign-in option **Continue with SSO**. If you're using the CLI, you must authenticate using a personal access token.

## Configure multiple IdPs

Docker supports multiple IdP configurations. To use multiple IdPs with one domain:

- Repeat Steps 1-4 on this page for each IdP.
- Each connection must use the same domain.
- Users will select **Continue with SSO** to choose their IdP at sign in.

## Enforce SSO

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
