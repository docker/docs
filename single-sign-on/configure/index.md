---
description: SSO configuration
keywords: configure, sso, docker hub, hub
title: Configure
redirect_from:
- /docker-hub/domains/
- /docker-hub/sso-connection/
- /docker-hub/enforcing-sso/
---

Follow the steps on this page to configure SSO for your organization or company. 

## Step one: Add and verify your domain

1. Sign in to Docker Hub, navigate to the **Organizations** page and select your organization or company.
2. Select **Settings**. If you are setting up SSO for an organization you then need to select **Security**. 
3. Select **Add Domain** and continue with the on-screen instructions to add the TXT Record Value to your domain name system (DNS).

    >**Note**
    >
    > Format your domains without protocol or www information, for example, yourcompany.com. This should include all email domains and subdomains users will use to access Docker. Public domains such as gmail.com, outlook.com, etc. aren’t permitted. Also, the email domain should be set as the primary email.

4. Once you have waited 72 hours for the TXT Record verification, you can then select **Verify** next to the domain you've added, and follow the on-screen instructions. 

![verify-domain](../images/verify-domain.png){: width="700px" }

## Step two: Create an SSO connection

> **Important**
>
> If your IdP setup requires an Entity ID and the ACS URL, you must select the
> **SAML** tab in the **Authentication Method** section. For example, if your
> Azure AD Open ID Connect (OIDC) setup uses SAML configuration within Azure
> AD, you must select **SAML**. If you are [configuring Open ID Connect with Azure AD](https://docs.microsoft.com/en-us/powerapps/maker/portals/configure/configure-openid-settings){: target="_blank" rel="noopener" class="_"} select
> **Azure AD** as the authentication method. Also, IdP initiated connections
> aren't supported at this time.
{: .important}


1. Once your domain is verified, in the **Single Sign-on Connection** table select **Create Connections**, and create a name for the connection. 

    > **Note**
    >
    > You have to verify at least one domain before creating the connections.

2. Select an authentication method, **SAML** or **Azure AD (OIDC)**.
3. Copy the following fields and add them to your IdP:

   - SAML: **Entity ID**, **ACS URL**
   - Azure AD (OIDC): **Redirect URL**

4. From your IdP, copy and paste the following values into the Docker **Settings** fields:

    - SAML: **SAML Sign-on URL**, **x509 Certificate**
    - Azure AD (OIDC): **Client ID**, **Client Secret**, **Azure AD Domain**

5. Select the verified domains you want to apply the connection to.

6. To provision your users, select the organization(s) and/or team(s).

    > **Note**
    >
    > If you are a company owner and have more than one organization, you need to select a default organization.

7. Review your summary and select **Create Connection**.

The SSO connection is now created. You can continue to set up [SSO Group Mapping and SCIM](../../docker-hub/scim.md) without enforcing SSO log-in.

## Optional step three: Test your SSO configuration

After you’ve completed the SSO configuration process in Docker Hub, you can test the configuration when you sign in to Docker Hub using an incognito browser. Log in to Docker Hub using your domain email address. You are then redirected to your IdP's login page to authenticate.

1. Authenticate through email instead of using your Docker ID, and test the login process.
2. To authenticate through CLI, your users must have a PAT before you enforce SSO for CLI users.

## Optional step four: Enforce SSO log-in in Docker Hub

1. In the **Single Sign-On Connections** table, select the **Action** icon and then **Enforce Single Sign-on**.
    When SSO is enforced, your users are unable to modify their email address and password, convert a user account to an organization, or set up 2FA through Docker Hub. You must enable 2FA through your IdP.
2. Continue with the on-screen instructions and verify that you’ve completed the tasks. 
3. Select **Turn on enforcement** to complete. 

To enforce SSO log-in for Docker Desktop, see [Enforce sign-in](../../docker-hub/configure-sign-in.md).
