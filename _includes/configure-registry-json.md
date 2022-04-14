<!-- This section is included in topics that contain instructions on how to configure registry.json file to enforce users to sign into Docker Desktop-->

## Create a registry.json file

When creating a `registry.json` file, ensure that the developer is a member of at least one organization in Docker Hub. If the `registry.json` file matches at least one organization the developer is a member of, they can sign in to Docker Desktop and access all their organizations.

### Windows

On Windows, run the following command in a terminal to install Docker Desktop:

`"Docker Desktop Installer.exe" install`

If you’re using PowerShell, you should run it as:

`Start-Process '.\win\build\Docker Desktop Installer.exe' -Wait install`

If using the Windows Command Prompt:

`start /w "" "Docker Desktop Installer.exe" install`

The `install` command accepts the following flag:

`--allowed-org=<org name>`

This requires the user to sign in and be part of the specified Docker Hub organization when running the application. For example:

 `Docker Desktop Installer.exe" install --allowed-org=acmeinc`

 This creates the registry.json file at `C:\ProgramData\DockerDesktop\registry.json` and includes the organization information the user belongs to. Make sure this file can't be edited by the individual developer, only by the administrator.

### Mac

After downloading `Docker.dmg`, run the following commands in a terminal to install Docker Desktop in the Applications folder:



```
sudo hdiutil attach Docker.dmg
sudo /Volumes/Docker/Docker.app/Contents/MacOS/install
sudo hdiutil detach /Volumes/Docker
```

The `install` command accepts the following flags:

`--allowed-org=<org name>`

This requires the user to sign in and be part of the specified Docker Hub organization when running the application. For example:

 `sudo hdiutil attach Docker.dmg --allowed-org=acmeinc`

This creates the registry.json file at `C:\ProgramData\DockerDesktop\registry.json` and includes the organization information the user belongs to. Make sure this file can't be edited by the individual developer, only by the administrator.
