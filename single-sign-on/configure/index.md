---
description: SSO configuration
keywords: configure, sso, docker hub, hub
title: Configure
---

To configure SSO, sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} to complete the IdP server configuration process. You can only configure SSO with a single IdP.  When this is complete, log back in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} and complete the SSO enablement process.

> **Important**
>
> If your IdP setup requires an Entity ID and the ACS URL, you must select the
> **SAML** tab in the **Authentication Method** section. For example, if your
> Azure AD Open ID Connect (OIDC) setup uses SAML configuration within Azure
> AD, you must select **SAML**. If you are [configuring Open ID Connect with Azure AD](https://docs.microsoft.com/en-us/powerapps/maker/portals/configure/configure-openid-settings){: target="_blank" rel="noopener" class="_"} select
> **Azure AD** as the authentication method. Also, IdP initiated connections
> aren't supported at this time.
{: .important}

The following video walks you through the process of configuring SSO.

<iframe width="560" height="315" src="https://www.youtube-nocookie.com/embed/QY0j02ggf64" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

### Configuring your IdP

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#SAML-2.0">SAML 2.0</a></li>
<li><a data-toggle="tab" data-target="#Azure-AD">Azure AD (OIDC)</a></li>
</ul>
<div class="tab-content">
<div id="SAML-2.0" class="tab-pane fade in active" markdown="1">

#### SAML 2.0

1. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} as an administrator and navigate to **Organizations** and select the organization that you want to enable SSO on.
2. Select **Settings** and select the **Security** tab.
3. Select an authentication method for **SAML 2.0**.

    ![SSO SAML1](/single-sign-on/images/sso-saml1.png){:width="500px"}

4. In the Identity Provider Set Up, copy the **Entity ID**, **ACS URL** and **Certificate Download URL**.

    ![SSO SAML2](/single-sign-on/images/sso-saml2.png){:width="500px"}

5. Sign in to your IdP to complete the IdP server configuration process. Refer to your IdP documentation for detailed instructions.

    > **Note**
    >
    > The NameID is your email address and is set as the default.
    > For example, yourname@mycompany.com. The optional `name` attribute is also supported. This attribute name must be lower-cased. _The following is an example of this attribute in Okta._

    ![SSO Attribute](/single-sign-on/images/sso-attribute.png){:width="500px"}

6. Complete the fields in the **Configuration Settings** section and select **Save**. If you want to change your IdP, you must delete your existing provider and configure SSO with your new IdP.

    ![SSO SAML3](/single-sign-on/images/sso-saml3.png){:width="500px"}

7. Proceed to **add your domain** before you test and enforce SSO.

<hr>
</div>
<div id="Azure-AD" class="tab-pane fade" markdown="1">

### Azure AD (OIDC)

>**Note**
>
> This section is for users who only want to configure Open ID Connect with
> Azure AD. This connection is a basic OIDC connection, and there are no
> special customizations available when using it.

1. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} as an administrator and navigate to **Organizations** and select the organization that you want to enable SSO on.
2. Select **Settings** and select the **Security** tab.
3. Select an authentication method for **Azure AD**.
4. In the Identity Provider Set Up, copy the **Redirect URL / Reply URL**.

    ![SSO Azure AD OIDC](/single-sign-on/images/sso-azure-oidc.png){:width="500px"}

5. Sign in to your IdP to complete the IdP server configuration process. Refer to your IdP documentation for detailed instructions.

    > **Note**
    >
    > The NameID is your email address and is set as the default.
    > For example: yourname@mycompany.com.

6. Complete the fields in the **Configuration Settings** section and click **Save**. If you want to change your IdP, you must delete your existing provider and configure SSO with your new IdP.

    ![SSO Azure3](/single-sign-on/images/sso-azure3.png){:width="500px"}

7. Proceed to **add your domain** before you test and enforce SSO.

<hr>
</div>
</div>

## Domain control

Select **Add Domain** and specify the corporate domain you’d like to manage with SSO. Format your domains without protocol or www information, for example, yourcompany.com. Docker supports multiple domains that are part of your IdP. Make sure that your domain is reachable through email.

> **Note**
>
> This should include all email domains and sub-domains users will use to access Docker.
> Public domains such as gmail.com, outlook.com, etc aren't permitted.
> Also, the email domain should be set as the primary email.

![SSO Domain](/single-sign-on/images/sso-domain.png){:width="500px"}

## Domain verification

To verify ownership of a domain, add a TXT record to your Domain Name System (DNS) settings.

1. Copy the provided TXT record value and navigate to your DNS host and locate the **Settings** page to add a new record.
2. Select the option to add a new record and paste the TXT record value into the applicable field. For example, the **Value**, **Answer** or **Description** field.

    Your DNS record may have the following fields:
    * Record type: enter your 'TXT' record value
    * Name/Host/Alias: leave the default (@ or blank)
    * Time to live (TTL): enter **86400**

3. After you have updated the fields, select **Save**.

    > **Note**
    >
    > It can take up to 72 hours for DNS changes to take effect, depending on
    > your DNS host. The Domains table will have an Unverified status during
    > this time.

4. In the Security section of your Docker organization, select **Verify** next to the domain you want to verify after 72 hours.

    > **Note**
    >
    > Once you've verified your domain, you can move forward to test your
    > configuration and enforce SSO, or you can configure your [System Cross-domain Identity Management (SCIM)](../../docker-hub/scim.md).

## Test your SSO configuration

After you’ve completed the SSO configuration process in Docker Hub, you can test the configuration when you sign in to Docker Hub using an incognito browser. Login using your domain email address and IdP password. You will then get redirected to your identity provider’s login page to authenticate.

1. Authenticate through email instead of using your Docker ID, and test the login process.
2. To authenticate through CLI, your users must have a PAT before you enforce SSO for CLI users.

## Enforce SSO in Docker Hub

Before you enforce SSO in Docker Hub, you must complete the following:
Test SSO by logging in and out successfully, confirm that all members in your org have upgraded to Docker Desktop version 4.4.2, PATs are created for each member, CI/CD passwords are converted to PAT. Also, when using Docker partner products (for example, VS Code), you must use a PAT when you enforce SSO. For your service accounts add your additional domains in **Add Domains** or enable the accounts in your IdP.

Admins can force users to authenticate with Docker Desktop by provisioning a registry.json configuration file. The registry.json file will force users to authenticate as a user that's configured in the allowedOrgs list in the registry.json file. For info on how to configure a registry.json file see [Configure registry.json](../../docker-hub/image-access-management.md#enforce-authentication)

1. On the Single Sign-On page in Docker Hub, select **Turn ON Enforcement** to enable your SSO.
2. When SSO is enforced, your users are unable to modify their email address and password, convert a user account to an organization, or set up 2FA through Docker Hub. You must enable 2FA through your IdP.

    > **Note**
    >
    > If you want to turn off SSO and revert back to Docker’s built-in
    > authentication, select **Turn OFF Enforcement**. Your users aren’t
    > forced to authenticate through your IdP and can sign in to Docker using
    > their personal credentials.

![SSO Enforced](/single-sign-on/images/sso-enforce.png){:width="500px"}
