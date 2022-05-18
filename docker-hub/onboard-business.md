---
description: On board with Docker Business
keywords: business, organizations, get started, onboarding
title: On board with Docker Business
toc_min: 1
toc_max: 2
---

The following section contains instructions on how to get started onboarding  your organization with Docker Business.

## Prerequisites

Before you continue to on board your organization, ensure that you've completed the following:

 - Subscribed to Docker Business

    If you haven't subscribed to Docker Business yet, [get in touch with us](https://www.docker.com/pricing/contact-sales/).

 - Created your organization and applied your Docker Business subscription to your organization

   On the day that your subscription becomes active, your organization's primary contact will receive a welcome email from Docker to guide you through creating an organization or using an existing organization for your Docker Business subscription.

## Step 1: Create teams

All members in your organization need to be in at least one team. Teams are used to apply access control permissions to image repositories and organization settings.

Your organization will have at least one default team, the **owners** team, with at least a single member (you). Members of the **owners** team can help manage users, teams, and repositories in the organization. [Learn more](../docker-hub/orgs.md#the-owners-team).

Besides the default **owners** team, you should create at least one additional team. This additional team will contain members that can not modify organization settings.

> **Note**
>
> When Docker Single Sign-on (SSO) is enabled, a default team named **company** is automatically created. All members added through SSO are initially added to this team and are granted read-only access to the organization's repositories. If you are going to use SSO and you don't want to create teams yet, you can skip to [Step 2: Create Docker Hub repositories and configure access control](#step-2-create-docker-hub-repositories-and-configure-access-control). You can revisit this step at any time if you want to create teams.

To create a team:

1. Click **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} and then select your organization.
2. Click **Teams** and then click **Create Team**.
3. Fill out your team's information and then click **Create**.

## Step 2: Create Docker Hub repositories and configure access control

If member's of your organization will access images on Docker Hub, then you should create at least one image repository.

To create a repository:

1. Click **Repositories** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and then click **Create Repository**.
2. Click the drop-down and select your organization to put the repository in your organization's namespace.
3. Select **Private** for the **Visibility** setting to make the repository only accessible to your organization or select **Public** to make the repository accessible by the public.
4. Link a GitHub or Bitbucket account now, or choose to do it later in the repository settings.
5. Click **Create**.

If you have created teams, then you must configure access to your private repositories:

1. Click **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and then select your organization.
2. Click **Repositories** and then select the repository that you’d like to configure access to.
3. Click **Permissions**.
4. Select a team and a permission from the drop-down lists and click **Add**.

## Step 3: Configure Image Access Management

You may want to protect your organization from malicious content by restricting what images your members can pull from Docker Hub. By default, members can pull all images from Docker Hub.

To configure Image Access Management permissions:

1. Click **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and then select your organization.
2. Click **Settings** and then click **Org Permissions**.
3. Toggle Image Access Management to **Enabled**.
4. Toggle the permissions for each image type to **Allowed** or **Restricted**.

## Step 4: Add members

Now that you have all your access control settings configured, it's time to start adding members. You can automatically add members to your organization by configuring Docker Single Sign-on (SSO), or invite members based their email address or Docker ID.

   > **Note**
   >
   > If you are not ready to configure SSO, you can invite members using their email address or Docker ID and then configure SSO at a later time. Any members you invite by email address or Docker ID will continue to have access after configuring SSO.
   >
   > In addition, when SSO is configured, you can still invite members not in your IdP by using their email address or Docker ID.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#sso-configure">Configure Single Sign-on</a></li>
<li><a data-toggle="tab" data-target="#invite-member">Invite members</a></li>
</ul>
<div class="tab-content">
<div id="sso-configure" class="tab-pane fade in active" markdown="1">

### Configure Single Sign-on

The following steps will help you quickly set up SSO. For more details, see [Configure Single Sign-on](../single-sign-on/index.md) and [Single Sign-on FAQs](../single-sign-on/faqs.md).

1. Ensure that all members have at least Docker Desktop 4.4.2 installed on their machines.
2. If you have existing Docker CI/CD pipelines in your organization, replace their passwords with Personal Access Tokens.
See [Create a Personal Access (PAT)](../single-sign-on/index.md#create-a-personal-access-token-pat) for more details.
3. Configure either your SAML 2.0 IdP or your Azure AD IdP with Open ID Connect.

   <ul class="nav nav-tabs">
   <li class="active"><a data-toggle="tab" data-target="#saml">SAML 2.0 IdP configuration</a></li>
   <li><a data-toggle="tab" data-target="#azure-ad">Azure AD IdP configuration with Open ID Connect</a></li>
   </ul>
   <div class="tab-content">
   <div id="saml" class="tab-pane fade in active" markdown="1">

   1. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} as an administrator and navigate to **Organizations** and select the organization that you want to enable SSO on.
   2. Click **Settings** and select the **Security** tab.
   3. Select an authentication method for SAML 2.0.
   4. In the Identity Provider Set Up, copy the **Entity ID**, **ACS URL** and **Certificate Download URL**.
   5. Log in to your IdP to complete the IdP server configuration process. Refer to your IdP documentation for detailed instructions.

      > **Note**
      >
      > The NameID is your email address and is set as the default. For example, yourname@mycompany.com. We also support the optional name attribute. This attribute name must be lower-cased. The following is an example of this attribute in Okta.

   6. Complete the fields in the **Configuration Settings** section and click **Save**. If you want to change your IdP, you must delete your existing provider and configure SSO with your new IdP.

   </div>
   <div id="azure-ad" class="tab-pane fade" markdown="1">

   > **Note**
   >
   > This section is for users who only want to configure Open ID Connect with Azure AD. This connection is a basic OIDC connection, and there are no special customizations available when using it.

   1. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} as an administrator and navigate to **Organizations** and select the organization that you want to enable SSO on.
   2. Click **Settings** and select the **Security** tab.
   3. Select an authentication method for Azure AD.
   4. In the Identity Provider Set Up, copy the **Redirect URL / Reply URL**.
   5. Log in to your IdP to complete the IdP server configuration process. Refer to your IdP documentation for detailed instructions.

      > **Note**
      >
      > The NameID is your email address and is set as the default. For example: yourname@mycompany.com.

   6. Complete the fields in the Configuration Settings section and click **Save**. If you want to change your IdP, you must delete your existing provider and configure SSO with your new IdP.

   </div>
   </div>

4. Click **Add Domain** and specify the corporate domain you'd like to manage with SSO. Domains should be formatted without protocol or www information, for example, yourcompany.com. Docker currently supports multiple domains that are part of your IdP. Make sure that your domain is reachable through email.

   > **Note**
   >
   > This should include all email domains and sub-domains users will use to access Docker. Public domains such as gmail.com and outlook.com are not permitted. Also, the email domain should be set as the primary email.

5. Perform the following steps to verify ownership of your domain by adding a TXT record to your Domain Name System (DNS) setting.
   1. Copy the provided **TXT record value** and navigate to your DNS host and locate the **Settings** page to add a new record.
   2. Select the option to add a new record and paste the TXT record value into the applicable field. For example, the **Value**, **Answer** or **Description** field.

      Your DNS record may have the following fields:

      * Record type: enter your 'TXT' record value
      * Name/Host/Alias: leave the default (@ or blank)
      * Time to live (TTL): enter **86400**
   3. After you have updated the fields, click **Save**.

      > **Note**
      >
      > It can take up to 72 hours for DNS changes to take effect, depending on your DNS host. The Domains table will have an Unverified status during this time.
   4. In the Security section of your Docker organization, click **Verify** next to the domain you want to verify after 72 hours.
6. Perform the following to verify that SSO has been configured successfully.
   1. Open an incognito browser.
   2. Navigate to [Docker Hub](https://hub.docker.com).
   3. Authenticate through email instead of using your Docker ID. If you are able to authenticate, then SSO has been configured successfully.
7. To access Docker Hub through the CLI, each member of your organization must create a Personal Access Token. See [Create an access token](../docker-hub/access-tokens.md#create-an-access-token) for details.
8. Perform the following to enforce users to sign in to Docker Hub using SSO.
   1. In [Docker Hub](https://hub.docker.com), click **Organizations** select your organization, click **Settings**, and then select the **Security** tab.
   2. Click **Turn ON Enforcement**.

</div>
<div id="invite-member" class="tab-pane fade" markdown="1">

### Invite members

To invite a member to your organization:

1. Navigate to **Organizations** in Docker Hub, and select your organization.
2. In the **Members** tab, click **Invite Member**.
3. Enter the invitee's Docker ID or email, and select a team from the drop-down list.

   > **Note**
   >
   > Remember, you shouldn't add all your members to the **owners** team. Members of the **owners** team have access to change all of your organization's settings.

4. Click **Invite** to confirm.

</div>
</div>

## Step 5: Enforce Sign in for Docker Desktop

In order to enforce access controls settings from Docker Business, you need to enforce members to sign in to Docker Desktop using a Docker account that is a member of your organization.

To enforce sign in, you need to create a `registry.json` file on each user's computer with the following contents, where `myorg` is replaced with your organization's name.

   ```console
   {
      "allowedOrgs":["myorg"]
   }
   ```

Based on your users' operating systems, you must create the registry.json file at:
- Mac: `/Library/Application Support/com.docker.docker/registry.json`
- Windows: `/ProgramData/DockerDesktop/registry.json`

> **Note**
>
> Ensure that only administrators have permission to modify the registry.json file. Users should not be able to edit the file.

The Docker Desktop installer can create this file and deploy it to the users' machines as part of the installation process. For more details and examples of different ways to create the registry.json file, see [Create a registry.json file](../docker-hub/configure-sign-in.md#create-a-registryjson-file)

## What's next

  As your organization grows, you can continue managing [teams](./docker-hub/orgs/#create-a-team), [members](./docker-hub/orgs/#invite-members), [repositories](./docker-hub/repos/), and [access control](./docker-hub/orgs/#configure-repository-permissions).

 Don't stop here. Get the most out of your Docker Business subscription by leveraging these popular features:

- Set up [Automated Builds](../docker-hub/builds/index.md) triggered by code pushes to help streamline your CI process.
- Add [service accounts](../docker-hub/service-accounts.md) to automate management of your container images or containerized applications.
- Configure [Vulnerability Scanning](../docker-hub/vulnerability-scanning.md) to improve security awareness.
- Integrate your image pushes into your workflow with [Webhooks](../docker-hub/webhooks.md.).
