---
description: Learn how to complete your single-sign on connection and next steps for enabling SSO.
keywords: configure, sso, docker hub, hub, docker admin, admin, security
title: Create an SSO connection
linkTitle: Connect
---

{{< summary-bar feature_name="SSO" >}}

Creating a single sign-on (SSO) connection requires setting up the connection in Docker first, followed by setting up the connection in your identity provider (IdP). This guide provides steps for setting up your SSO connection in Docker and your IdP.

> [!TIP]
>
> This guide requires copying and pasting values in both Docker and your IdP. To ensure a seamless connection process, complete all the steps in this guide in one session and keep separate browsers open for both Docker and your IdP.

## Prerequisites

Make sure you have completed the following before you begin:

- Your domain is verified
- You have an account set up with an IdP
- You have completed the steps in the [Configure single sign-on](../single-sign-on/configure.md) guide

## Step one: Create an SSO connection in Docker

{{< summary-bar feature_name="Admin console early access" >}}

>[!NOTE]
>
> Before creating an SSO connection in Docker, you must verify at least one domain.

{{< tabs >}}
{{< tab name="Admin Console" >}}

1. Sign in to the [Admin Console](https://admin.docker.com/).
2. Select your organization or company from the **Choose profile** page. Note that when an organization is part of a company, you must select the company and configure the domain for the organization at the company level.
3. Under Security and access, select **SSO and SCIM**.
4. Select **Create Connection** and provide a name for the connection.
5. Select an authentication method, **SAML** or **Azure AD (OIDC)**.
6. Copy the following fields to add to your IdP:
    - Okta SAML: **Entity ID**, **ACS URL**
    - Azure OIDC: **Redirect URL**
7. Keep this window open so you can paste the connection information from your IdP here at the end of this guide.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

1. Sign in to Docker Hub.
2. Select **Organizations** and then your organization from the list.
3. On your organization page, select **Settings** and then **Security**.
4. In the SSO connection table, select **Create Connection** and provide a name for the connection.
5. Select an authentication method, **SAML** or **Azure AD (OIDC)**.
6. Copy the following fields to add to your IdP:
    - Okta SAML: **Entity ID**, **ACS URL**
    - Azure OIDC: **Redirect URL**
7. Keep this window open so you can paste the connection information from your IdP here at the end of this guide.

{{< /tab >}}
{{< /tabs >}}

## Step two: Create an SSO connection in your IdP

The user interface for your IdP may differ slightly from the following steps. Refer to the documentation for your IdP to verify.

{{< tabs >}}
{{< tab name="Okta SAML" >}}

1. Sign in to your Okta account.
2. Select **Admin** to open the Okta Admin portal.
3. From the left-hand navigation, select **Administration**.
4. Select **Administration** and then **Create App Integration**.
5. Select **SAML 2.0** and then **Next**.
6. Enter "Docker Hub" as your **App Name**.
7. Optional. Upload a logo.
8. Select **Next**.
9. Enter the following values from Docker into their corresponding Okta fields:
    - Docker ACS URL: **Single Sign On URL**
    - Docker Entity ID: **Audience URI (SP Entity ID)**
10. Configure the following settings in Okta:
    - Name ID format: `EmailAddress`
    - Application username: `Email`
    - Update application on: `Create and update`
11. Select **Next**.
12. Select the **This is an internal app that we have created** checkbox.
13. Select **Finish**.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

1. Sign in to your Azure AD admin portal.
2. Select **Default Directory** and then **Add**.
3. Choose **Enterprise Application** and select **Create your own application**.
4. Enter "Docker" for application name and select the **non-gallery** option.
5. After the application is created, go to **Single Sign-On** and select **SAML**.
6. Select **Edit** on the **Basic SAML configuration** section.
7. Enter the following values from Docker into their corresponding Azure fields:
    - Docker Entity ID: **Identifier**
    - Docker ACS URL: **Reply URL**
8. Save configuration.
9. From the **SAML Signing Certificate** section, download your **Certificate (Base64)**.

{{< /tab >}}
{{< tab name="Azure Connect (OIDC)" >}}

To create an Azure Connect (OIDC) connection, you must create an app registration, client secrets, and configure API permissions for Docker:

### Create app registration

1. Sign in to your Azure AD admin portal.
2. Select **App Registration** and then **New Registration**.
3. Enter "Docker Hub SSO" or similar for application name.
4. Under **Supported account types**, specify who can use this application or access the app.
5. In the **Redirect URI** section, select **Web** from the drop-down menu and paste the **Redirect URI** value from the Docker console into this field.
6. Select **Register** to register the app.
7. Copy the **Client ID** from the app's overview page. You need this information to continue configuring SSO in Docker.

### Create client secrets

1. Open your app in Azure AD and select **Certificates & secrets**.
2. Select **+ New client secret**.
3. Specify the description of the secret and set how long keys can be used.
4. Select **Add** to continue.
5. Copy the secret **Value** field. You need this to continue configuring SSO in Docker.

### Configure API permissions

1. Open your app in Azure AD and navigate to your app settings.
2. Select **API permission** and then **Grant admin consent for [your tenant name]**.
3. Select **Yes** to confirm.
4. After confirming, select **Add a permission** and then **Delegated permissions**.
5. Search for `User.Read` and select this option.
6. Select **Add permissions** to confirm.
7. Verify admin consent was granted for each permission by checking the **Status** column.

{{< /tab >}}
{{< /tabs >}}

## Step three: Connect Docker and your IdP

After creating your connection in Docker and your IdP, you can cross-connect them to complete your SSO connection:

{{< tabs >}}
{{< tab name="Okta SAML" >}}

1. Open your app you created in Okta and select **View SAML setup instructions**.
2. Copy the following values from the Okta SAML setup instruction page:
    - **SAML Sign-in URL**
    - **x509 Certificate**
3. Open Docker Hub or the Admin Console. Your SSO configuration page should still be open from Step one of this guide.
4. Select **Next** to open the **Update single-sign on connection** page.
5. Paste your Okta **SAML Sign-in URL** and **x509 Certificate** values in Docker.
6. Select **Next**.
7. Optional. Select a default team to provision users to and select **Next**.
8. Verify your SSO connection details and select **Create Connection**.

{{< /tab >}}
{{< tab name="Entra ID SAML 2.0" >}}

1. Open your app in Azure AD.
2. Open your downloaded **Certificate (Base64)** in a text editor.
3. Copy the following values:
    - From Azure AD: **Login URL**
    - Copy your the contents of your **Certificate (Base64)** file from your text editor
4. Open Docker Hub or the Admin Console. Your SSO configuration page should still be open from Step one of this guide.
5. Paste your **Login URL** and **Certificate (Base64)** values in Docker.
6. Select **Next**.
7. Optional. Select a default team to provision users to and select **Next**.
8. Verify your SSO connection details and select **Create Connection**.

{{< /tab >}}
{{< tab name="Azure Connect (OIDC)" >}}

1. Open Docker Hub or the Admin Console. Your SSO configuration page should still be open from Step one of this guide.
2. Paste the following values from Azure AD in to Docker:
    - **Client ID**
    - **Client Secret**
    - **Azure AD Domain**
3. Select **Next**.
4. Optional. Select a default team to provision users to and select **Next**.
5. Verify your SSO connection details and select **Create Connection**.

{{< /tab >}}
{{< /tabs >}}

## Step four: Test your connection

After you've completed the SSO connection process in Docker, we recommend testing it:

1. Open an incognito browser.
2. Sign in to the Admin Console using your **domain email address**.
3. The browser will redirect to your IdP's login page to authenticate.
4. Authenticate through your domain email instead of using your Docker ID.

You can also test your SSO connection through the command-line interface (CLI). If you want to test through the CLI, your users must have a personal access token (PAT).


## Optional: Enforce SSO

>[!IMPORTANT]
>
> If SSO isn't enforced, users can choose to sign in with either their Docker username and password or SSO.

Enforcing SSO requires users to use SSO when signing into Docker. This centralizes authentication and enforces policies set by the IdP.

1. Sign in to the [Admin Console](https://admin.docker.com/).
2. Select your organization or company from the **Choose profile** page. Note that when an organization is part of a company, you must select the company and configure the domain for the organization at the company level.
3. Under Security and access, select **SSO and SCIM**.
4. In the SSO connections table, select the **Action** icon and then **Enable enforcement**. When SSO is enforced, your users are unable to modify their email address and password, convert a user account to an organization, or set up 2FA through Docker Hub. If you want to use 2FA, you must enable 2FA through your IdP.
5. Continue with the on-screen instructions and verify you've completed all tasks.
6. Select **Turn on enforcement** to complete.

Your users must now sign in to Docker with SSO.

> [!NOTE]
>
> When SSO is enforced, [users can't use passwords to access the Docker CLI](/security/security-announcements/#deprecation-of-password-logins-on-cli-when-sso-enforced). Users must use a [personal access token](/manuals/security/for-admins/access-tokens.md) (PAT) for authentication to access the Docker CLI.

## More resources

The following videos demonstrate how to enforce SSO.

- [Video: Enforce SSO with Okta SAML](https://youtu.be/c56YECO4YP4?feature=shared&t=1072)
- [Video: Enforce SSO with Azure AD (OIDC)](https://youtu.be/bGquA8qR9jU?feature=shared&t=1087)


## What's next

- [Provision users](/manuals/security/for-admins/provisioning/_index.md)
- [Enforce sign-in](../enforce-sign-in/_index.md)
- [Create access tokens](/manuals/security/for-admins/access-tokens.md)
