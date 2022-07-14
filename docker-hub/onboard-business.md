---
description: Get started onboarding with Docker Business
keywords: business, organizations, get started, onboarding
title: Get started onboarding with Docker Business
toc_min: 1
toc_max: 2
---

The following section contains step-by-step instructions on how to get started onboarding your organization after you obtain a Docker Business subscription.

## Prerequisites

Before you start to on board your organization, ensure that you've completed the following:
- You have a Docker Business subscription. [Get in touch with us](https://www.docker.com/pricing/contact-sales/) if you haven't subscribed to Docker Business yet.
- Your Docker Business subscription is new. If you upgraded your Docker Team subscription or renewed your Docker Business subscription, see [what's next](#whats-next).
- Your Docker Business subscription has started. You cannot complete all the steps until after your subscription start date.
-  You are familiar with Docker terminology. If you discover any unfamiliar terms, see the [glossary](/glossary/#docker) or [FAQs](../docker-hub/onboarding-faqs.md).

## Step 1: Identify your Docker users and their Docker accounts

To begin, you should identify which users you will need to add to your Docker Business organization. Identifying your users will help you efficiently allocate your subscription's seats and manage access.

> **Note**
>
> If you will use Docker Single Sign-on (SSO), users from your identity provider (IdP) are automatically provisioned in your organization when they sign in. You can identify which users have signed in to Docker Hub by viewing your organization's members in Docker Hub. You can also perform the steps below to identify users before you configure SSO.

1. Identify the Docker users in your organization.
   - If your organization uses device management software, like MDM or JAMF, you may use the device management software to help identify Docker users. See your device management software's documentation for details. You can identify Docker users by checking if Docker Desktop is installed at the following location on each user's machine:
      - Mac: `/Applications/Docker.app`
      - Windows: `C:\Program Files\Docker\Docker`
   - If your organization does not use device management software, you may survey your users.
2. Instruct all your Docker users in your organization to update their existing Docker account's email address to an address that's in your organization's domain, or to create a new account using an email address in your organization's domain.
   - To update an account's email address, instruct your users to sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, go to [Account Settings](https://hub.docker.com/settings/general){: target="_blank" rel="noopener" class="_"}, and update the email address to their email address in your organization's domain.
   - To create a new account, instruct your users to go [sign up](https://hub.docker.com/signup){: target="_blank" rel="noopener" class="_"} using their email address in your organization's domain.
3. Ask your Docker sales representative to provide a list of Docker accounts that use an email address in your organization's domain.

## Step 2: Add your Docker Business subscription to an organization

On the day that your Docker Business subscription starts, your organization's primary contact will receive a welcome email from Docker to guide you through creating a new organization or to let you choose an existing organization for your Docker Business subscription.

> **Note**
>
> If your organization's primary contact does not receive a welcome email from Docker on the day that your subscription starts:
>   - Check your email spam folder.
>   - Use the steps below to verify that your Docker Business organization does not already exist.
>   - Contact your Docker sales representative to verify your primary contact's email address.

After completing the steps from the welcome email, verify that your organization exists and your organization has a Docker Business subscription:

1. Go to [Billing Details](https://hub.docker.com/billing){: target="_blank" rel="noopener" class="_"} and then select on your organization's name.
2. Under **Plan**, view your subscription. If you organization has a Docker Business subscription, you will see **Docker Business**.

## Step 3: Add members

Now that you have a Docker Business organization, it's time to start adding members. You can automatically add members to your organization by configuring Docker Single Sign-on (SSO), or invite members based their email address or Docker ID.

> **Note**
>
> If you are not ready to configure SSO, you can invite members using their email address or Docker ID and then configure SSO at a later time. Any members you invite by email address or Docker ID can continue to have access after configuring SSO.
>
> In addition, when SSO is configured, you can still invite members not in your identity provider (IdP) by using their email address or Docker ID.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#sso-configure">Configure Single Sign-on</a></li>
<li><a data-toggle="tab" data-target="#invite-member">Invite members</a></li>
</ul>
<div class="tab-content">
<div id="sso-configure" class="tab-pane fade in active" markdown="1">

### Configure Single Sign-on

The following steps will help you quickly set up SSO. For more details, see [Configure Single Sign-on](../single-sign-on/index.md){: target="_blank" rel="noopener" class="_"} and [Single Sign-on FAQs](../single-sign-on/faqs.md){: target="_blank" rel="noopener" class="_"}.

1. Ensure that all members have at least [Docker Desktop](../desktop/index.md/#download-and-install){: target="_blank" rel="noopener" class="_"} 4.4.2 installed on their machines.
2. If you have existing Docker CI/CD pipelines in your organization, replace their passwords with Personal Access Tokens.
See [Create a Personal Access (PAT)](../single-sign-on/index.md/#create-a-personal-access-token-pat){: target="_blank" rel="noopener" class="_"} for more details.
3. Configure either your SAML 2.0 identity provider (IdP) or your Azure AD IdP with Open ID Connect.

   <ul class="nav nav-tabs">
   <li class="active"><a data-toggle="tab" data-target="#saml">SAML 2.0 IdP configuration</a></li>
   <li><a data-toggle="tab" data-target="#azure-ad">Azure AD IdP configuration with Open ID Connect</a></li>
   </ul>
   <div class="tab-content">
   <div id="saml" class="tab-pane fade in active" markdown="1">

   1. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} as an administrator and navigate to **Organizations** and select the organization that you want to enable SSO on.
   2. Select **Settings** and select the **Security** tab.
   3. Select an authentication method for SAML 2.0.
   4. In the Identity Provider Set Up, copy the **Entity ID**, **ACS URL** and **Certificate Download URL**.
   5. Log in to your IdP to complete the IdP server configuration process. Refer to your IdP documentation for detailed instructions.

      > **Note**
      >
      > The NameID is your email address and is set as the default. For example, yourname@mycompany.com. We also support the optional name attribute. This attribute name must be lower-cased. The following is an example of this attribute in Okta.

   6. Complete the fields in the **Configuration Settings** section and select **Save**. If you want to change your IdP, you must delete your existing provider and configure SSO with your new IdP.

   </div>
   <div id="azure-ad" class="tab-pane fade" markdown="1">

   > **Note**
   >
   > This section is for users who only want to configure Open ID Connect with Azure AD. This connection is a basic OIDC connection, and there are no special customizations available when using it.

   1. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} as an administrator and navigate to **Organizations** and select the organization that you want to enable SSO on.
   2. Select **Settings** and select the **Security** tab.
   3. Select an authentication method for Azure AD.
   4. In the Identity Provider Set Up, copy the **Redirect URL / Reply URL**.
   5. Log in to your IdP to complete the IdP server configuration process. Refer to your IdP documentation for detailed instructions.

      > **Note**
      >
      > The NameID is your email address and is set as the default. For example: yourname@mycompany.com.

   6. Complete the fields in the Configuration Settings section and select **Save**. If you want to change your IdP, you must delete your existing provider and configure SSO with your new IdP.

   </div>
   </div>

4. Select **Add Domain** and specify the corporate domain you'd like to manage with SSO. Domains should be formatted without protocol or www information, for example, yourcompany.com. Docker currently supports multiple domains that are part of your IdP. Make sure that your domain is reachable through email.

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
   3. After you have updated the fields, select **Save**.

      > **Note**
      >
      > It can take up to 72 hours for DNS changes to take effect, depending on your DNS host. The Domains table will have an Unverified status during this time.
   4. In the Security section of your Docker organization, select **Verify** next to the domain you want to verify after 72 hours.
6. Perform the following to verify that SSO has been configured successfully.
   1. Open an incognito browser.
   2. Navigate to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}.
   3. Authenticate through email instead of using your Docker ID. If you are able to authenticate, then SSO has been configured successfully.
7. To access Docker Hub through the CLI, each member of your organization must create a Personal Access Token. See [Create an access token](../docker-hub/access-tokens.md/#create-an-access-token){: target="_blank" rel="noopener" class="_"} for details.
8. Perform the following to force users to sign in to Docker Hub using SSO.
   1. In [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, select **Organizations**, select your organization, select **Settings**, and then select the **Security** tab.
   2. Select **Turn ON Enforcement**.

</div>
<div id="invite-member" class="tab-pane fade" markdown="1">

### Invite members

All members in your organization need to be in at least one team. Teams are used to apply access control permissions to image repositories and organization settings.

Your organization will have at least one default team, the **owners** team, with at least a single member (you). Members of the **owners** team can help manage users, teams, and repositories in the organization. [Learn more](../docker-hub/orgs.md/#the-owners-team){: target="_blank" rel="noopener" class="_"}.

In the steps below, you will create a **members** team. Members that you invite to the **members** team will not be able to modify your organization settings.

To create the **members** team:

1. Select **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} and then select your organization.
2. Click **Teams** and then click **Create Team**.
3. Specify `members` for **Team name** and then click **Create**.

To invite a member to the **members** team in your organization:

1. Navigate to **Organizations** in Docker Hub, and select your organization.
2. In the **Members** tab, click **Invite Member**.
3. Enter the invitee's Docker ID or email, and select the **members** team from the drop-down list.
4. Click **Invite** to confirm.

</div>
</div>

## Step 4: Enforce sign in for Docker Desktop

By default, members of your organization can use Docker Desktop on their machines without signing in to any Docker account. To ensure that a user signs in to a Docker account that is a member of your organization and that the
organization’s settings apply to the user’s session, you can use a `registry.json` file.

The `registry.json` file is a configuration file that allows administrators to specify the Docker organization the user must belong to and ensure that the organization’s settings apply to the user’s session. The Docker Desktop installer can create this file on the users’ machines as part of the installation process.

After a `registry.json` file is configured on a user’s machine, Docker Desktop prompts the user to sign in. If a user doesn’t sign in, or tries to sign in using a different organization, other than the organization listed in the `registry.json` file, they will be denied access to Docker Desktop.

Deploying a `registry.json` file and forcing users to authenticate is not required, but offers the following benefits:

 - Allows administrators to configure features such as [Image Access Management](image-access-management.md) which allows team members to:
    - Only have access to Trusted Content on Docker Hub
    - Pull only from the specified categories of images
- Authenticated users get a higher pull rate limit compared to anonymous users. For example, if you are authenticated, you get 200 pulls per 6 hour period, compared to 100 pulls per 6 hour period per IP address for anonymous users. For more information, see [Download rate limit](download-rate-limit.md).
- Blocks users from accessing Docker Desktop until they are added to a specific organization.

### Create a registry.json file

Before creating a `registry.json` file, ensure that the user is a member of
your organization in Docker Hub.

Based on the user's operating system, you must create a `registry.json` file at the following location and ensure that the file can't be edited by the user:
   - Windows: `/ProgramData/DockerDesktop/registry.json`
   - Mac: `/Library/Application Support/com.docker.docker/registry.json`

The `registry.json` file must contain the following contents, where `myorg` is replaced with your organization's name.

```json
{
   "allowedOrgs":["myorg"]
}
```

You can use the following methods to create a `registry.json` file based on the user's operating system.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#windows">Windows</a></li>
<li><a data-toggle="tab" data-target="#mac">Mac</a></li>
</ul>
<div class="tab-content">
<div id="windows" class="tab-pane fade in active" markdown="1">


#### Windows

On Windows, you can use the following methods to create a `registry.json` file.


##### Create registry.json when installing Docker Desktop on Windows

To automatically create a `registry.json` file when installing Docker Desktop, download `Docker Desktop Installer.exe` and run one of the following commands from the directory containing `Docker Desktop Installer.exe`. Replace `myorg` with your organization's name.

If you're using PowerShell:

```powershell
PS> Start-Process '.\Docker Desktop Installer.exe' -Wait install --allowed-org=myorg
```

If you're using the Windows Command Prompt:

```console
C:\Users\Admin> "Docker Desktop Installer.exe" install --allowed-org=myorg
```

##### Create registry.json manually on Windows

To manually create a `registry.json` file, run the following PowerShell command as an Admin and replace `myorg` with your organization's name:

```powershell
PS>  Set-Content /ProgramData/DockerDesktop/registry.json '{"allowedOrgs":["myorg"]}'
```

This creates the `registry.json` file at `C:\ProgramData\DockerDesktop\registry.json` and includes the organization information the user belongs to. Make sure this file can't be edited by the user, only by the administrator.

</div>
<div id="mac" class="tab-pane fade" markdown="1">

#### Mac

On Mac, you can use the following methods to create a `registry.json` file.


#####  Create registry.json when installing Docker Desktop on Mac

To automatically create a registry.json file when installing Docker Desktop, download `Docker.dmg` and run the following commands in a terminal from the directory containing `Docker.dmg`. Replace `myorg` with your organization's name.

```bash
$ sudo hdiutil attach Docker.dmg 
$ sudo /Volumes/Docker/Docker.app/Contents/MacOS/install --allowed-org=myorg
$ sudo hdiutil detach /Volumes/Docker
```

#####  Create registry.json manually on Mac

To manually create a `registry.json` file, run the following commands in a terminal and replace `myorg` with your organization's name.

```bash
$ sudo touch /Library/Application Support/com.docker.docker/registry.json
$ sudo echo '{"allowedOrgs":["myorg"]}' >> /Library/Application Support/com.docker.docker/registry.json
```

This creates the `registry.json` file at `/Library/Application Support/com.docker.docker/registry.json` and includes the organization information the user belongs to. Make sure this file can't be edited by the user, only by the administrator.

</div></div>

## What's next

Get the most out of your Docker Business subscription by leveraging these popular features:

- If you haven't configured [Docker Single Sign-on](../single-sign-on/index.md) yet, configure it now for centralized account management.
- Create [repositories](../docker-hub/repos.md) to share container images.
- Create [teams](../docker-hub/orgs.md/#create-a-team) and configure [repository permissions](../docker-hub/orgs.md/#configure-repository-permissions).
- Control which images your members can access with [Image Access Management](../docker-hub/image-access-management.md/).
- Control which registries your members can access with [Registry Access Management](../docker-hub/registry-access-management.md/).

Your Docker Business subscription provides many more additional features. [Learn more](../subscription/index.md).
