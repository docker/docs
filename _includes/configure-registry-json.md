<!-- This section is included in topics that contain instructions on how to configure registry.json file to enforce users to sign into Docker Desktop-->

## Create a registry.json file

Before creating a `registry.json` file, ensure that the user is a member of
at least one organization in Docker Hub. If the `registry.json` file matches at
least one organization the user is a member of, they can sign in to Docker
Desktop and access all their organizations.

Based on the user's operating system, you must create a `registry.json` file at the following location and make sure the file can't be edited by the user:
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


### Windows

On Windows, you can use the following methods to create a `registry.json` file.


#### Create registry.json when installing Docker Desktop on Windows

To automatically create a `registry.json` file when installing Docker Desktop, download `Docker Desktop Installer.exe` and run one of the following commands from the directory containing `Docker Desktop Installer.exe`. Replace `myorg` with your organization's name.

If you're using PowerShell:

```powershell
PS> Start-Process '.\Docker Desktop Installer.exe' -Wait install --allowed-org=myorg
```

If you're using the Windows Command Prompt:

```console
C:\Users\Admin> "Docker Desktop Installer.exe" install --allowed-org=myorg
```

#### Create registry.json manually on Windows

To manually create a `registry.json` file, run the following PowerShell command as an Admin and replace `myorg` with your organization's name:

```powershell
PS>  Set-Content /ProgramData/DockerDesktop/registry.json '{"allowedOrgs":["myorg"]}'
```

This creates the `registry.json` file at `C:\ProgramData\DockerDesktop\registry.json` and includes the organization information the user belongs to. Make sure this file can't be edited by the user, only by the administrator.

</div>
<div id="mac" class="tab-pane fade" markdown="1">

### Mac

On Mac, you can use the following methods to create a `registry.json` file.


####  Create registry.json when installing Docker Desktop on Mac

To automatically create a registry.json file when installing Docker Desktop, download `Docker.dmg` and run the following commands in a terminal from the directory containing `Docker.dmg`. Replace `myorg` with your organization's name.

```bash
$ sudo hdiutil attach Docker.dmg 
$ sudo /Volumes/Docker/Docker.app/Contents/MacOS/install --allowed-org=myorg
$ sudo hdiutil detach /Volumes/Docker
```

####  Create registry.json manually on Mac

To manually create a `registry.json` file, run the following commands in a terminal and replace `myorg` with your organization's name.

```bash
$ sudo touch /Library/Application Support/com.docker.docker/registry.json
$ sudo echo '{"allowedOrgs":["myorg"]}' >> /Library/Application Support/com.docker.docker/registry.json
```

This creates the `registry.json` file at `/Library/Application Support/com.docker.docker/registry.json` and includes the organization information the user belongs to. Make sure this file can't be edited by the user, only by the administrator.

</div></div>