<!-- This section is included in topics that contain instructions on how to configure registry.json file to enforce users to sign into Docker Desktop-->

## Create a registry.json file

When creating a `registry.json` file, ensure that the developer is a member of
at least one organization in Docker Hub. If the `registry.json` file matches at
least one organization the developer is a member of, they can sign in to Docker
Desktop and access all their organizations.

### Windows

On Windows, you can run a command in a terminal to install Docker Desktop, or you can download Docker Desktop and manually create your `registry.json` file.

Run the following command in a terminal to install Docker Desktop:

```console
C:\Users\Admin> "Docker Desktop Installer.exe" install
```

If you’re using PowerShell, you should run it as:

```console
PS> Start-Process '.\win\build\Docker Desktop Installer.exe' -Wait install
```

If using the Windows Command Prompt:

```console
C:\Users\Admin> start /w "" "Docker Desktop Installer.exe" install
```

The `install` command accepts the following flag:

`--allowed-org=<org name>`

This requires the user to sign in and be part of the specified Docker Hub organization when running the application. For example:

```console
C:\Users\Admin> "Docker Desktop Installer.exe" install --allowed-org=acmeinc
```

To manually create a `registry.json` file:

1. Open Windows PowerShell and select **Run as Administrator**.
2. Type the following command `cd /ProgramData/DockerDesktop/`
3. Type `notepad registry.json` and enter the name of the Docker Hub
   organization that the developer belongs to in the `allowedOrgs` key and click
   **Save**. For example:

    ```json
    {
        "allowedOrgs": ["myorg"]
    }
    ```

This creates the `registry.json` file at `C:\ProgramData\DockerDesktop\registry.json` and includes the organization information the user belongs to. Make sure this file can't be edited by the individual developer, only by the administrator.

### Mac

On macOS, you can run a command in a terminal to install Docker Desktop, or you can download Docker Desktop and manually create your `registry.json` file.

Download `Docker.dmg` and run the following commands in a terminal to install Docker Desktop in the Applications folder:

```console
$ sudo hdiutil attach Docker.dmg
$ sudo /Volumes/Docker/Docker.app/Contents/MacOS/install
$ sudo hdiutil detach /Volumes/Docker
```

The `install` command accepts the following flags:

`--allowed-org=<org name>`

This requires the user to sign in and be part of the specified Docker Hub
organization when running the application. For example:

```console
$ sudo hdiutil attach Docker.dmg --allowed-org=acmeinc
```

To manually create a `registry.json` file on macOS, you must create a file at `/Library/Application Support/com.docker.docker/registry.json` with file permissions that ensure that the developer using Docker Desktop cannot remove or edit the file (that is, only the system administrator can write to the file). The file must be of type `JSON` and contain the name of the Docker Hub organization names in the `allowedOrgs` key.

To create your `registry.json` file:

1. Navigate to VS Code or any text editor of your choice.
2. Enter the name of the Docker Hub organization that the developer belongs to in the  `allowedOrgs` key and save it in your Documents. For example:

    ```json
    {
        "allowedOrgs": ["myorg"]
    }
    ```
3. Open a new terminal and type the following command:

    ```console
    $ sudo mkdir -p /Library/Application\ Support/com.docker.docker
    ```
    If prompted, type your password associated with your local computer.
4. Type the following command:
    ```console
    $ sudo cp Documents/registry.json /Library/Application\ Support/com.docker.docker/registry.json
    ```

This creates the `registry.json` file at `/Library/Application Support/com.docker.docker/registry.json`
and includes the organization information the user belongs to. Make sure this file
can't be edited by the individual developer, only by the administrator.