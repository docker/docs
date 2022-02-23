---
description: Registry Access Management
keywords: registry, access, managment
title: Registry Access Management
---

Registry Access Management is a new feature available to organizations with a Docker Business subscription. This feature lets organization owners manage the registries that their users can access while using Docker Desktop.  When using this feature, organization owners can ensure that their users can only access their trusted registries, such as a secure private registry on Artifactory, thereby reducing the security risks that can occur when users interact with public registries.

## Configure Registry Access Management permissions

To configure Registry Access Management permissions, perform the following steps:

1. Sign into your [Docker Hub](https://hub.docker.com) account as an organization owner.
2. Select an organization, navigate to the **Settings** tab on the **Organizations** page and click **Org Permissions**.
3. Enable Registry Access Management to set the permissions for your registry.

> **Note**
>
> When enabled, the Docker Hub registry is set by default, however you can also restrict this registry for your developers.

4. Click **Add** and enter your registry details in the applicable fields, and click **Create** to add the registry to your list.
5. Verify that the registry appears in your list and click **Save & Apply**.  You can verify that your changes are saved in the Activity tab.

  > **Note**
  >
  > Once you add a registry, it can take a maximum of 24 hours for the changes to be enforced on your developers’ machines. If you want to apply the changes sooner, you must force a Docker logout on your developers’ machine and have the developer re-authenticate for Docker Desktop.  Also, there is no limit on the number of registries you can add. See the Caveats sections to learn more about limitations  when using this feature.

### Enforce authentication

To ensure that each org member uses Registry Access Management on their local machine, you can perform the following steps below to enforce sign-in under your organization. To do this:

1. Download the latest version of Docker Desktop, and then
2. Create a `registry.json` file.

Download Docker Desktop 4.5 or a later release.

- [Download and install for Windows](/desktop/windows/install/)
- [Download and install for Mac](/desktop/mac/install/)

#### Create a registry json file

Before you create a `registry.json` file, ensure that the developer is a member of at least one organization in Docker Hub. If the registry.json file matches at least one organization the developer is a member of, they can sign into Docker Desktop and access all of their organizations.

**On Windows**

On Windows, you must create a file `C:\ProgramData\DockerDesktop\registry.json` with file permissions that ensure that the developer using Docker Desktop cannot remove or edit the file (that is, only the system administrator can write to the file). The file must be `JSON` and contain one or more organization names in the `allowedOrgs` key.

To create your `registry.json` file on Windows:

1. Open Windows Powershell and select Run as Administrator.
2. Type the following command: `cd /ProgramData/DockerDesktop/`
3. In Notepad, type `registry.json` and enter the Docker Hub organization that the user belongs to in `allowedOrgs` key and click Save.

    For example:

    ```json
    {
    "allowedOrgs": ["myorg"]
    }
    ```

4. Navigate to Powershell and type ```start .```

Congratulations! You have just created the registry.json file.

**On macOS**:

On macOS, you must create a file at `/Library/Application Support/com.docker.docker/registry.json` with file permissions that ensure that the developer using Docker Desktop cannot remove or edit the file (that is, only the system administrator can write to the file). The file must be of type JSON and contain the name of the Docker Hub organization names in the allowedOrgs key.

To create your `registry.json` file on macOS:

1. Navigate to VS Code or any text editor of your choice.
2. Enter one or more organization names in the `allowedOrgs` key and save it in your Documents.

    For example:

    ```json
    {
     "allowedOrgs": ["myorg"]
    }
    ```

 3. Open a new terminal and type the following command:

    `sudo mkdir -p /Library/Application\ Support/com.docker.docker`

    Note: if prompted, type your password associated with your local computer.

4. Type the following command:

    `sudo cp Documents/registry.json /Library/Application\ Support/com.docker.docker/registry.json`

Congratulations! You have just created the `registry.json` file.

### Verify the restrictions

   After you’ve created the registry.json file and deployed it onto the users’ machines, you can verify whether the changes have taken effect by asking users to start Docker Desktop.

   If the configuration is successful, Docker Desktop prompts the user to authenticate using the organization credentials on start. If the user fails to authenticate, they will see an error message, and they will be denied access to Docker Desktop.

### Caveats

  There are certain limitations when using Registry Access Management; they are as follows:

  * Windows image pulls, and image builds are not restricted
  * Builds such as `docker buildx` using a Kubernetes driver are not restricted
  * Builds such as `docker buildx` using a custom docker-container driver are not restricted
  * Blocking is DNS-based; you must use a registry's access control mechanisms to distinguish between “push” and “pull”
  * You must disable HTTP proxy or use a corporate proxy which also blocks the registries
  * WSL 2 requires at least a 5.4 series Linux kernel  (this does not apply to earlier Linux kernel series)
 * Under the WSL 2 network, traffic from all Linux distributions is restricted (this will be resolved in the updated 5.15 series Linux kernel)

  Also, Registry Access Management operates on the level of hosts, not IP addresses, that can be bypassed if a user manipulates their domain resolution, such as by running Docker against a local proxy or modifying their operating system's `sts` file. Blocking these forms of manipulation is outside the remit of Docker Desktop.

