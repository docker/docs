---
description: On board with Docker Business
keywords: business, organizations, get started, onboarding
title: On board with Docker Business
toc_min: 1
toc_max: 2
---

The following section contains steps on how to get started onboarding your organization with Docker Business.

## Prerequisite

Before you start to on board your organization, you must have a Docker Business subscription. [Get in touch with us](https://www.docker.com/pricing/contact-sales/) if you haven't subscribed to Docker Business yet.

## Step 1: Identify your Docker users

To begin, you should identify which users you will need to add to your Docker Business organization.

Later on, in [Step 7: Enforce sign in for Docker Desktop](#step-7-enforce-sign-in-for-docker-desktop), you can enforce sign in and prevent unidentified users, who were not added to your Docker organization, from signing in to Docker Desktop on your organization's computers.

> **Note**
>
> If you will use Docker Single Sign-on (SSO), users from your IdP are automatically provisioned in your organization when they sign in. You can identify which users have signed in to Docker Hub by viewing your organization's members in Docker Hub. You can also perform the steps below to identify users before you configure SSO.

If your organization uses device management software, like MDM or JAMF, then you may use the device management software to help identify users that already have Docker Desktop installed. See your device management software's documentation for details.

To identify users without using device management software:

1. Instruct all your Docker users in your organization to update their existing Docker account's email address to an address that's in your organization's domain, or to create a new account using an email address in your organization's domain.
   * To update an account's email address, instruct your users to sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, go to [Account Settings](https://hub.docker.com/settings/general){: target="_blank" rel="noopener" class="_"}, and update the email address to their email address in your organization's domain.
   * To create a new account, instruct your users to go to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} and register account using their email address in your organization's domain.
2. Ask your Docker sales representative to provide a list of Docker accounts that use an email address in your organization's domain.

## Step 2: Create your organization

Next, you'll need an organization in Docker Hub. An organization is a collection of teams and repositories that can be managed together. [Learn more](/docker-hub/orgs/){: target="_blank" rel="noopener" class="_"}.

On the day that your Docker Business subscription begins, your organization's primary contact will receive a welcome email from Docker to guide you through creating a new organization or to let you choose an existing organization for your Docker Business subscription.

> **Note**
>
> If your organization's primary contact does not receive a welcome email from Docker on the day that your subscription begins:
>   * Check your email spam folder.
>   * Use the steps below to verify that your Docker Business organization does not already exist.
>   * Contact your Docker sales representative to verify your primary contact's email address.

After completing the steps from the welcome email, verify that your organization exists and that your Docker Business subscription is active:

1. Go to [Billing Details](https://hub.docker.com/billing){: target="_blank" rel="noopener" class="_"} and then click on your organization's name. 
2. Under **Plan**, view your subscription. You subscription should be **Docker Business**.

## Step 3: Create teams

All members in your organization need to be in at least one team. Teams are used to apply access control permissions to image repositories and organization settings.

Your organization will have at least one default team, the **owners** team, with at least a single member (you). Members of the **owners** team can help manage users, teams, and repositories in the organization. [Learn more](/docker-hub/orgs.md#the-owners-team){: target="_blank" rel="noopener" class="_"}.

Besides the default **owners** team, you should create at least one additional team. This additional team will contain members that can not modify organization settings.

> **Note**
>
> If you will use Docker Single Sign-on (SSO), a default team named **company** is automatically created. All members added through SSO are initially added to this team and are granted read-only access to your organization's repositories. If you are going to use SSO and you don't want to create teams yet, you can skip to [Step 4: Create Docker Hub repositories and configure access control](#step-4-create-docker-hub-repositories-and-configure-access-control). You can revisit this step at any time if you want to create teams.

To create a team:

1. Click **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} and then select your organization.
2. Click **Teams** and then click **Create Team**.
3. Fill out your team's information and then click **Create**.

## Step 4: Create Docker Hub repositories and configure repository access

If member's of your organization will share container images on Docker Hub, then you should create at least one image repository. [Learn more](/docker-hub/repos/){: target="_blank" rel="noopener" class="_"}.

To create a repository:

1. Click **Repositories** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and then click **Create Repository**.
2. Click the drop-down and select your organization to put the repository in your organization's namespace.
3. Select **Private** for the **Visibility** setting to make the repository only accessible to your organization or select **Public** to make the repository accessible by the public.
4. Link a GitHub or Bitbucket account now, or choose to do it later in the repository settings.
5. Click **Create**.

If you have created **Private** repositories, then only the **owners** team will have initial access.

> **Note**
>
> If you will use Docker Single Sign-on (SSO), a default team named **company** is automatically created. All members added through SSO are initially added to this team and are granted read-only access to your organization's repositories.

To configure access to private repositories for teams other than **owners** and **company**:

1. Click **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and then select your organization.
2. Click **Repositories** and then select the repository that you’d like to configure access to.
3. Click **Permissions**.
4. Select a team and a permission from the drop-down lists and click **Add**.

## Step 5: Configure Image Access Management

You may want to protect your organization from malicious content by restricting what images your members can pull from Docker Hub. By default, members can pull all images from Docker Hub.

To configure Image Access Management permissions:

1. Click **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and then select your organization.
2. Click **Settings** and then click **Org Permissions**.
3. Toggle Image Access Management to **Enabled**.
4. Toggle the permissions for each image type to **Allowed** or **Restricted**.

## Step 6: Add members

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

The following steps will help you quickly set up SSO. For more details, see [Configure Single Sign-on](/single-sign-on/index.md){: target="_blank" rel="noopener" class="_"} and [Single Sign-on FAQs](/single-sign-on/faqs.md){: target="_blank" rel="noopener" class="_"}.

1. Ensure that all members have at least Docker Desktop 4.4.2 installed on their machines.
2. If you have existing Docker CI/CD pipelines in your organization, replace their passwords with Personal Access Tokens.
See [Create a Personal Access (PAT)](/single-sign-on/index.md#create-a-personal-access-token-pat){: target="_blank" rel="noopener" class="_"} for more details.
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
   2. Navigate to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}.
   3. Authenticate through email instead of using your Docker ID. If you are able to authenticate, then SSO has been configured successfully.
7. To access Docker Hub through the CLI, each member of your organization must create a Personal Access Token. See [Create an access token](../docker-hub/access-tokens.md#create-an-access-token){: target="_blank" rel="noopener" class="_"} for details.
8. Perform the following to force users to sign in to Docker Hub using SSO.
   1. In [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, click **Organizations** select your organization, click **Settings**, and then select the **Security** tab.
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

## Step 7: Enforce sign in for Docker Desktop

At this point, your users can sign in to Docker Desktop on their computers using any Docker account, including accounts that are not a member of your Docker organization.

When users sign in to an account that is not a member of your Docker organization, they can circumvent your Docker organization's Image Access Management settings. To ensure Image Access Management is applied on your organization's computers, you can force users to sign in to Docker Desktop using an account that is a member of your organization.

Enforcing sign in is not required. You should enforce sign in when you want to enforce Image Access Management settings or you want to ensure that only your identified organization members sign in to Docker Desktop on your organization's computers.

To enforce sign in, first inform your users that they must sign in to Docker Desktop using only their Docker account that is a member of your organization, and then you need to create a `registry.json` file on each user's computer with the following contents, where `myorg` is replaced with your organization's name.

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

The Docker Desktop installer can create this file as part of the installation process or you can use other methods to deploy this file. For more details and examples of different ways to create the registry.json file, see [Create a registry.json file](/docker-hub/configure-sign-in.md#create-a-registryjson-file){: target="_blank" rel="noopener" class="_"}

## What's next

  As your organization grows, you can continue managing [teams](/docker-hub/orgs/#create-a-team){: target="_blank" rel="noopener" class="_"}, [members](/docker-hub/orgs/#invite-members){: target="_blank" rel="noopener" class="_"}, [repositories](/docker-hub/repos/){: target="_blank" rel="noopener" class="_"}, and [access control](/docker-hub/orgs/#configure-repository-permissions){: target="_blank" rel="noopener" class="_"}.

 Don't stop here. Get the most out of your Docker Business subscription by leveraging these popular features:

- Set up [Automated Builds](/docker-hub/builds/index.md){: target="_blank" rel="noopener" class="_"} triggered by code pushes to help streamline your CI process.
- Add [service accounts](/docker-hub/service-accounts.md){: target="_blank" rel="noopener" class="_"} to automate management of your container images or containerized applications.
- Configure [Vulnerability Scanning](/docker-hub/vulnerability-scanning.md){: target="_blank" rel="noopener" class="_"} to improve security awareness.
- Integrate your image pushes into your workflow with [Webhooks](/docker-hub/webhooks.md){: target="_blank" rel="noopener" class="_"}.
