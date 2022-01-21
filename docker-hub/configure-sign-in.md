---
description: Configure registry.json to require sign in
keywords: authentication, registry.json, configure,
title: Configure registry.json to require sign in
---

The `registry.json` file is a configuration file that allows administrators to
specify the Docker organization the user must belong to, and thereby ensure
that the organization's settings are applied to the user's session. Docker
Desktop installation requires admin access. In large enterprises where admin
access is restricted, administrators can create a `registry.json` file and
deploy it to the users’ machines using a device management software as part of
the Docker Desktop installation process.

After you deploy a `registry.json` file to a user’s machine, it prompts the user to sign into Docker Desktop. If a user doesn’t sign in, or tries to sign in using a different organization, other than the organization listed in the `registry.json` file, they will be denied access to Docker Desktop.
Deploying a `registry.json` file and forcing users to authenticate offers the following benefits:

1. Allows administrators to configure features such as [Image Access Management](image-access-management.md) which allows team members to:
    - Only have access to Trusted Content on Docker Hub
    - Pull only from the specified categories of images
2. Authenticated users get a higher pull rate limit compared to anonymous users. For example, if you are authenticated, you get 200 pulls per 6 hour period, compared to 100 pulls per 6 hour period per IP address for anonymous users. For more information, see [Download rate limit](download-rate-limit.md).
3. Blocks users from accessing Docker Desktop until they are added to a specific organization.

## Create a registry.json file

Before you create a `registry.json` file, ensure that the user is a member of at least one organization in Docker Hub. If the `registry.json` file matches at least one organization the user is a member of, they can sign into Docker Desktop, and then access all their organizations.

### Windows

On Windows, you must create a file at
`C:\ProgramData\DockerDesktop\registry.json` with file permissions that ensure
that the developer using Docker Desktop cannot remove or edit the file (that is,
only the system administrator can write to the file). The file must be of type
`JSON` and contain the name of the organization in the `allowedOrgs` key.

To create your `registry.json` file on Windows:

1. Open Windows Powershell and select Run as Administrator.
2. Type the following command `cd /ProgramData/DockerDesktop/`
3. Type `notepad registry.json` and enter the name of the Docker Hub
   organization that the user belongs to in the `allowedOrgs` key and click
   **Save**. For example:

    ```json
    {
        "allowedOrgs": ["myorg"]
    }
    ```

### Mac

On macOS, you must create a file at `/Library/Application Support/com.docker.docker/registry.json` with file permissions that ensure that
the developer using Docker Desktop cannot remove or edit the file (that is, only
the system administrator can write to the file). The file must be of type `JSON`
and contain the name of the Docker Hub organization names in the `allowedOrgs`
key.

To create your `registry.json` file on macOS:

1. Navigate to VS Code or any text editor of your choice.
2. Enter the name of the Docker Hub organization that the user belongs to in the
   `allowedOrgs` key and save it in your Documents. For example:

    ```json
    {
        "allowedOrgs": ["myorg"]
    }
    ```

3. Open a new terminal and type the following command:

    ```console
    sudo mkdir -p /Library/Application\ Support/com.docker.docker
    ```

    If prompted, type your password associated with your local computer.

4. Type the following command:

     ```console
    sudo cp Documents/registry.json /Library/Application\ Support/com.docker.docker/registry.json
    ```

## Verify the changes

After you’ve created the `registry.json` file and deployed it onto the users’ machines, you can verify whether the changes have taken effect by asking users to start Docker Desktop.

If the configuration is successful, Docker Desktop prompts the user to authenticate using the organization credentials on start. If the user fails to authenticate, they will see an error message, and they will be denied access to Docker Desktop.
