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

## Step on: Add and verify your domain

1. Sign in to Docker Hub, navigate to the Organization page and select your organization or company.
2. Select **Settings** and then **Security**. 
3. Select **Add Domain** and continue with the on-screen instructions to add the TXT Record Value to your domain name system (DNS).

    >**Note**
    >
    > Format your domains without protocol or www information, for example, yourcompany.com. This should include all email domains and subdomains users will use to access Docker. Public domains such as gmail.com, outlook.com, etc aren’t permitted. Also, the email domain should be set as the primary email.

4. Once you have waited 72 hours for the TXT Record verification, you can then select **Verify** next to the domain you've added, and follow the instructions. 

![verify-domain](images/verify-domain.png){: width="700px" }

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


1. Once your domain is verified, continue to **Single Sign-on Connections** and select **Create Connections**, and create a name for the connection. 

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

**SSO connection is now created**. You can continue to set up SSO Group Mapping and SCIM without enforcing SSO log-in.

## Step three: Test your SSO configuration

After you’ve completed the SSO configuration process in Docker Hub, you can test the configuration when you sign in to Docker Hub using an incognito browser. Login using your domain email address. You will then get redirected to your identity provider’s login page to authenticate.

1. Authenticate through email instead of using your Docker ID, and test the login process.
2. To authenticate through CLI, your users must have a PAT before you enforce SSO for CLI users.

## Step four: Enforce SSO in Docker Hub

Without SSO enforcement, users can continue to sign in using Docker username and password. If users login with your Domain email, they will authenticate through your identity provider instead.


You must test your SSO connection first if you’d like to enforce SSO log-in. All users must authenticate with an email address instead of their Docker ID if SSO is enforced.

Before you enforce SSO in Docker Hub, you must complete the following:
Test SSO by logging in and out successfully, confirm that all members in your org have upgraded to Docker Desktop version 4.4.2, PATs are created for each member, CI/CD passwords are converted to PAT. Also, when using Docker partner products (for example, VS Code), you must use a PAT when you enforce SSO. For your service accounts add your additional domains in **Add Domains** or enable the accounts in your IdP.

Admins can force users to authenticate with Docker Desktop by provisioning a registry.json configuration file. The registry.json file will force users to authenticate as a user that's configured in the allowedOrgs list in the registry.json file. For info on how to configure a registry.json file see [Configure registry.json](../../docker-hub/image-access-management.md#enforce-authentication)

1. In the **Single Sign-On Connections** table, select the Action icon and **Enforce Single Sign-on**.
    When SSO is enforced, your users are unable to modify their email address and password, convert a user account to an organization, or set up 2FA through Docker Hub. You must enable 2FA through your IdP.

    > **Note**
    >
    > If you want to turn off SSO and revert back to Docker’s built-in
    > authentication, select **Turn OFF Enforcement**. Your users aren’t
    > forced to authenticate through your IdP and can sign in to Docker using
    > their personal credentials.

2. Continue with the on-screen instructions and verify that you’ve completed the tasks. 
3. Select **Turn on enforcement** to complete. 
