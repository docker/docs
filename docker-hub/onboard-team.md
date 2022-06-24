---
description: Docker Team onboarding
keywords: team, organizations, get started, onboarding
title: Docker Team onboarding
toc_min: 1
toc_max: 2
---

The following section contains step-by-step instructions on how to get started onboarding your organization after you obtain a Docker Team subscription.

## Prerequisites

Before you start to on board your organization, ensure that you've completed the following:
- You have a Docker Team subscription. [Buy now](https://www.docker.com/pricing/) if you haven't subscribed to Docker Team yet.
-  You are familiar with Docker terminology. If you discover any unfamiliar terms, see the [glossary](/glossary/#docker) or [FAQs](../docker-hub/onboarding-faqs.md).


## Step 1: Identify your Docker users and their Docker accounts

To begin, you should identify which users you will need to add to your Docker Team organization. Identifying your users will help you efficiently allocate your subscription's seats and manage access.

1. Identify the Docker users in your organization.
   - If your organization uses device management software, like MDM or JAMF, you may use the device management software to help identify Docker users. See your device management software's documentation for details. You can identify Docker users by checking if Docker Desktop is installed at the following location on each user's machine:
      - Mac: `/Applications/Docker.app`
      - Windows: `C:\Program Files\Docker\Docker`
   - If your organization does not use device management software, you may survey your users.
2. Instruct all your Docker users in your organization to update their existing Docker account's email address to an address that's in your organization's domain, or to create a new account using an email address in your organization's domain.
   - To update an account's email address, instruct your users to sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, go to [Account Settings](https://hub.docker.com/settings/general){: target="_blank" rel="noopener" class="_"}, and update the email address to their email address in your organization's domain.
   - To create a new account, instruct your users to go [sign up](https://hub.docker.com/signup){: target="_blank" rel="noopener" class="_"} using their email address in your organization's domain.
3. Ask your Docker sales representative to provide a list of Docker accounts that use an email address in your organization's domain.

## Step 2: Add members

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



## Step 3: Enforce sign in for Docker Desktop

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

Get the most out of your Docker Team subscription by leveraging these popular features:

- Create [repositories](../docker-hub/repos.md) to share container images.
- Create [teams](../docker-hub/orgs.md/#create-a-team) and configure [repository permissions](../docker-hub/orgs.md/#configure-repository-permissions).

Your Docker Team subscription provides many more additional features. [Learn more](../subscription/index.md).
