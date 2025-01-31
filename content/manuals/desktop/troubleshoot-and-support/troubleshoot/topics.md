---
description: Explore common troubleshooting topics for Docker Desktop
keywords: Linux, Mac, Windows, troubleshooting, topics, Docker Desktop
title: Troubleshoot topics for Docker Desktop
linkTitle: Common topics
toc_max: 4
tags: [ Troubleshooting ]
weight: 10 
aliases:
 - /desktop/troubleshoot/topics/
---

> [!TIP]
>
> If you do not find a solution in troubleshooting, browse the GitHub repositories or create a new issue:
>
> - [docker/for-mac](https://github.com/docker/for-mac/issues)
> - [docker/for-win](https://github.com/docker/for-win/issues)
> - [docker/for-linux](https://github.com/docker/for-linux/issues)

## Topics for all platforms

### Make sure certificates are set up correctly

Docker Desktop ignores certificates listed under insecure registries, and
does not send client certificates to them. Commands like `docker run` that
attempt to pull from the registry produces error messages on the command line,
like this:

```console
Error response from daemon: Get http://192.168.203.139:5858/v2/: malformed HTTP response "\x15\x03\x01\x00\x02\x02"
```

As well as on the registry. For example:

```console
2017/06/20 18:15:30 http: TLS handshake error from 192.168.203.139:52882: tls: client didn't provide a certificate
2017/06/20 18:15:30 http: TLS handshake error from 192.168.203.139:52883: tls: first record does not look like a TLS handshake
```

### Docker Desktop's UI appears green, distorted, or has visual artifacts

Docker Desktop uses hardware-accelerated graphics by default, which may cause problems for some GPUs. In such cases,
Docker Desktop will launch successfully, but some screens may appear green, distorted,
or have some visual artifacts.

To work around this issue, disable hardware acceleration by creating a `"disableHardwareAcceleration": true` entry in Docker Desktop's `settings-store.json` file (or `settings.json` for Docker Desktop versions 4.34 and earlier). You can find this file at:

- Mac: `~/Library/Group Containers/group.com.docker/settings-store.json`
- Windows: `C:\Users\[USERNAME]\AppData\Roaming\Docker\settings-store.json`
- Linux: `~/.docker/desktop/settings-store.json.`

After updating the `settings-store.json` file, close and restart Docker Desktop to apply the changes.

## Topics for Linux and Mac

### Volume mounting requires file sharing for any project directories outside of `$HOME`

If you are using mounted volumes and get runtime errors indicating an
application file is not found, access to a volume mount is denied, or a service
cannot start, such as when using [Docker Compose](/manuals/compose/gettingstarted.md),
you might need to turn on [file sharing](/manuals/desktop/settings-and-maintenance/settings.md#file-sharing).

Volume mounting requires shared drives for projects that live outside of the
`/home/<user>` directory. From **Settings**, select **Resources** and then **File sharing**. Share the drive that contains the Dockerfile and volume.

### Docker Desktop fails to start on MacOS or Linux platforms

On MacOS and Linux, Docker Desktop creates Unix domain sockets used for inter-process communication.

Docker fails to start if the absolute path length of any of these sockets exceeds the OS limitation which is 104 characters on MacOS and 108 characters on Linux. These sockets are created under the user's home directory. If the user ID length is such that the absolute path of the socket exceeds the OS path length limitation, then Docker Desktop is unable to create the socket and fails to start. The workaround for this is to shorten the user ID which we recommend has a maximum length of 33 characters on MacOS and 55 characters on Linux. 

Following are the examples of errors on MacOS which indicate that the startup failure was due to exceeding the above mentioned OS limitation:

```console
[vpnkit-bridge][F] listen unix <HOME>/Library/Containers/com.docker.docker/Data/http-proxy-control.sock: bind: invalid argument
```

```console
[com.docker.backend][E] listen(vsock:4099) failed: listen unix <HOME>/Library/Containers/com.docker.docker/Data/vms/0/00000002.00001003: bind: invalid argument
```

## Topics for Mac

### Incompatible CPU detected

> [!TIP]
>
> If you are seeing this error, check you've installed the correct Docker Desktop for your architecture. 

Docker Desktop requires a processor (CPU) that supports virtualization and, more
specifically, the [Apple Hypervisor
framework](https://developer.apple.com/library/mac/documentation/DriversKernelHardware/Reference/Hypervisor/).
Docker Desktop is only compatible with Mac systems that have a CPU that supports the Hypervisor framework. Most Macs built in 2010 and later support it,as described in the Apple Hypervisor Framework documentation about supported hardware:

*Generally, machines with an Intel VT-x feature set that includes Extended Page
Tables (EPT) and Unrestricted Mode are supported.*

To check if your Mac supports the Hypervisor framework, run the following command in a terminal window.

```console
$ sysctl kern.hv_support
```

If your Mac supports the Hypervisor Framework, the command prints
`kern.hv_support: 1`.

If not, the command prints `kern.hv_support: 0`.

See also, [Hypervisor Framework
Reference](https://developer.apple.com/library/mac/documentation/DriversKernelHardware/Reference/Hypervisor/)
in the Apple documentation, and Docker Desktop [Mac system requirements](/manuals/desktop/setup/install/mac-install.md#system-requirements).

### VPNKit keeps breaking

In Docker Desktop version 4.19, gVisor replaced VPNKit to enhance the performance of VM networking when using the Virtualization framework on macOS 13 and above.

To continue using VPNKit, add `"networkType":"vpnkit"` to your `settings-store.json` file located at `~/Library/Group Containers/group.com.docker/settings-store.json`.

## Topics for Windows

### Volumes

#### Permissions errors on data directories for shared volumes

When sharing files from Windows, Docker Desktop sets permissions on [shared volumes](/manuals/desktop/settings-and-maintenance/settings.md#file-sharing)
to a default value of [0777](https://chmodcommand.com/chmod-0777/)
(`read`, `write`, `execute` permissions for `user` and for `group`).

The default permissions on shared volumes are not configurable. If you are
working with applications that require permissions different from the shared
volume defaults at container runtime, you need to either use non-host-mounted
volumes or find a way to make the applications work with the default file
permissions.

See also,
[Can I change permissions on shared volumes for container-specific deployment requirements?](/manuals/desktop/troubleshoot-and-support/faqs/windowsfaqs.md#can-i-change-permissions-on-shared-volumes-for-container-specific-deployment-requirements)
in the FAQs.

#### Volume mounting requires shared folders for Linux containers

If you are using mounted volumes and get runtime errors indicating an
application file is not found, access is denied to a volume mount, or a service
cannot start, such as when using [Docker Compose](/manuals/compose/gettingstarted.md),
you might need to turn on [shared folders](/manuals/desktop/settings-and-maintenance/settings.md#file-sharing).

With the Hyper-V backend, mounting files from Windows requires shared folders for Linux containers. From **Settings**, select **Shared Folders** and share the folder that contains the
Dockerfile and volume.

#### Support for symlinks

Symlinks work within and across containers. To learn more, see [How do symlinks work on Windows?](/manuals/desktop/troubleshoot-and-support/faqs/windowsfaqs.md#how-do-symlinks-work-on-windows).

#### Avoid unexpected syntax errors, use Unix style line endings for files in containers

Any file destined to run inside a container must use Unix style `\n` line
endings. This includes files referenced at the command line for builds and in
RUN commands in Docker files.

Docker containers and `docker build` run in a Unix environment, so files in
containers must use Unix style line endings: `\n`, _not_ Windows style: `\r\n`.
Keep this in mind when authoring files such as shell scripts using Windows
tools, where the default is likely to be Windows style line endings. These
commands ultimately get passed to Unix commands inside a Unix based container
(for example, a shell script passed to `/bin/sh`). If Windows style line endings
are used, `docker run` fails with syntax errors.

For an example of this issue and the resolution, see this issue on GitHub:
[Docker RUN fails to execute shell
script](https://github.com/moby/moby/issues/24388).

#### Path conversion on Windows

On Linux, the system takes care of mounting a path to another path. For example, when you run the following command on Linux:

```console
$ docker run --rm -ti -v /home/user/work:/work alpine
```

It adds a `/work` directory to the target container to mirror the specified path.

However, on Windows, you must update the source path. For example, if you are using 
the legacy Windows shell (`cmd.exe`), you can use the following command:

```console
$ docker run --rm -ti -v C:\Users\user\work:/work alpine
```

This starts the container and ensures the volume becomes usable. This is possible because Docker Desktop detects
the Windows-style path and provides the appropriate conversion to mount the directory.

Docker Desktop also allows you to use Unix-style path to the appropriate format. For example:

```console
$ docker run --rm -ti -v /c/Users/user/work:/work alpine ls /work
```

#### Working with Git Bash

Git Bash (or MSYS) provides a Unix-like environment on Windows. These tools apply their own
preprocessing on the command line. For example, if you run the following command in Git Bash, it gives an error:

```console
$ docker run --rm -ti -v C:\Users\user\work:/work alpine
docker: Error response from daemon: mkdir C:UsersUserwork: Access is denied.
```

This is because the `\` character has a special meaning in Git Bash. If you are using Git Bash, you must neutralize it using `\\`:

```console
$ docker run --rm -ti -v C:\\Users\\user\\work:/work alpine
```

Also, in scripts, the `pwd` command is used to avoid hard-coding file system locations. Its output is a Unix-style path.

```console
$ pwd
/c/Users/user/work
```

Combined with the `$()` syntax, the command below works on Linux, however, it fails on Git Bash.

```console
$ docker run --rm -ti -v $(pwd):/work alpine
docker: Error response from daemon: OCI runtime create failed: invalid mount {Destination:\Program Files\Git\work Type:bind Source:/run/desktop/mnt/host/c/Users/user/work;C Options:[rbind rprivate]}: mount destination \Program Files\Git\work not absolute: unknown.
```

You can work around this issue by using an extra `/`

```console
$ docker run --rm -ti -v /$(pwd):/work alpine
```

Portability of the scripts is not affected as Linux treats multiple `/` as a single entry.
Each occurrence of paths on a single line must be neutralized.

```console
$ docker run --rm -ti -v /$(pwd):/work alpine ls /work
ls: C:/Program Files/Git/work: No such file or directory
```

In this example, The `$(pwd)` is not converted because of the preceding '/'. However, the second '/work' is transformed by the
POSIX layer before passing it to Docker Desktop. You can also work around this issue by using an extra `/`.

```console
$ docker run --rm -ti -v /$(pwd):/work alpine ls //work
```

To verify whether the errors are generated from your script, or from another source, you can use an environment variable. For example:

```console
$ MSYS_NO_PATHCONV=1 docker run --rm -ti -v $(pwd):/work alpine ls /work
```

It only expects the environment variable here. The value doesn't matter.

In some cases, MSYS also transforms colons to semicolon. Similar conversions can also occur
when using `~` because the POSIX layer translates it to a DOS path. `MSYS_NO_PATHCONV` also works in this case.

### Virtualization

Your machine must have the following features for Docker Desktop to function correctly:

#### WSL 2 and Windows Home

1. Virtual Machine Platform
2. [Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/install-win10)
3. [Virtualization enabled in the BIOS](https://support.microsoft.com/en-gb/windows/enable-virtualization-on-windows-c5578302-6e43-4b4b-a449-8ced115f58e1)
   Note that many Windows devices already have virtualization enabled, so this may not apply.
4. Hypervisor enabled at Windows startup

![WSL 2 enabled](../../images/wsl2-enabled.png)

#### Hyper-V

On Windows 10 Pro or Enterprise, you can also use Hyper-V with the following features enabled:

1. [Hyper-V](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/hyper-v-technology-overview)
   installed and working
2. [Virtualization enabled in the BIOS](https://support.microsoft.com/en-gb/windows/enable-virtualization-on-windows-c5578302-6e43-4b4b-a449-8ced115f58e1)
   Note that many Windows devices already have virtualization enabled, so this may not apply.
3. Hypervisor enabled at Windows startup

![Hyper-V on Windows features](../../images/hyperv-enabled.png)

Docker Desktop requires Hyper-V as well as the Hyper-V Module for Windows
PowerShell to be installed and enabled. The Docker Desktop installer enables
it for you.

Docker Desktop also needs two CPU hardware features to use Hyper-V: Virtualization and Second Level Address Translation (SLAT), which is also called Rapid Virtualization Indexing (RVI). On some systems, Virtualization must be enabled in the BIOS. The steps required are vendor-specific, but typically the BIOS option is called `Virtualization Technology (VTx)` or something similar. Run the command `systeminfo` to check all required Hyper-V features. See [Pre-requisites for Hyper-V on Windows 10](https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/reference/hyper-v-requirements) for more details.

To install Hyper-V manually, see [Install Hyper-V on Windows 10](https://msdn.microsoft.com/en-us/virtualization/hyperv_on_windows/quick_start/walkthrough_install). A reboot is *required* after installation. If you install Hyper-V without rebooting, Docker Desktop does not work correctly.

From the start menu, type **Turn Windows features on or off** and press enter.
In the subsequent screen, verify that Hyper-V is enabled.

#### Virtualization must be turned on

In addition to [Hyper-V](#hyper-v) or [WSL 2](/manuals/desktop/features/wsl/_index.md), virtualization must be turned on. Check the
Performance tab on the Task Manager. Alternatively, you can type 'systeminfo' into your terminal. If you see 'Hyper-V Requirements:   A hypervisor has been detected. Features required for Hyper-V will not be displayed', then virtualization is enabled.

![Task Manager](../../images/virtualization-enabled.png)

If you manually uninstall Hyper-V, WSL 2 or turn off virtualization,
Docker Desktop cannot start. 

To turn on nested virtualization, see [Run Docker Desktop for Windows in a VM or VDI environment](/manuals/desktop/setup/vm-vdi.md#turn-on-nested-virtualization).

#### Hypervisor enabled at Windows startup

If you have completed the steps described above and are still experiencing
Docker Desktop startup issues, this could be because the Hypervisor is installed,
but not launched during Windows startup. Some tools (such as older versions of 
Virtual Box) and video game installers turn off hypervisor on boot. To turn it back on:

1. Open an administrative console prompt.
2. Run `bcdedit /set hypervisorlaunchtype auto`.
3. Restart Windows.

You can also refer to the [Microsoft TechNet article](https://social.technet.microsoft.com/Forums/en-US/ee5b1d6b-09e2-49f3-a52c-820aafc316f9/hyperv-doesnt-work-after-upgrade-to-windows-10-1809?forum=win10itprovirt) on Code flow guard (CFG) settings.

#### Turn on nested virtualization

If you are using Hyper-V and you get the following error message when running Docker Desktop in a VDI environment:

```console
The Virtual Machine Management Service failed to start the virtual machine 'DockerDesktopVM' because one of the Hyper-V components is not running
```

Try [enabling nested virtualization](/manuals/desktop/setup/vm-vdi.md#turn-on-nested-virtualization).


### Windows containers and Windows Server

Docker Desktop is not supported on Windows Server. If you have questions about how to run Windows containers on Windows 10, see
[Switch between Windows and Linux containers](/manuals/desktop/troubleshoot-and-support/faqs/windowsfaqs.md#how-do-i-switch-between-windows-and-linux-containers).

A full tutorial is available in [docker/labs](https://github.com/docker/labs) on
[Getting Started with Windows Containers](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md).

You can install a native Windows binary which allows you to develop and run
Windows containers without Docker Desktop. However, if you install Docker this way, you cannot develop or run Linux containers. If you try to run a Linux container on the native Docker daemon, an error occurs:

```none
C:\Program Files\Docker\docker.exe:
 image operating system "linux" cannot be used on this platform.
 See 'C:\Program Files\Docker\docker.exe run --help'.
```

### `Docker Desktop Access Denied` error message when starting Docker Desktop

Docker Desktop displays the **Docker Desktop - Access Denied** error if a Windows user is not part of the **docker-users** group.

If your admin account is different to your user account, add the **docker-users** group. Run **Computer Management** as an administrator and navigate to **Local Users and Groups** > **Groups** > **docker-users**.

Right-click to add the user to the group. Sign out and sign back in for the changes to take effect.

