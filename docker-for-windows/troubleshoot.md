---
description: Troubleshooting, logs, and known issues
keywords: windows, troubleshooting, logs, issues
redirect_from:
- /windows/troubleshoot/
title: Logs and troubleshooting
---

This page explains how to diagnose and troubleshoot problems you may be having with Docker for Windows.

## Getting help

There are several ways to get the support you need with Docker for Windows. If you encounter problems not addressed here in the documentation:

- Refer to the knowledge base articles at the [Docker Success Center](https://success.docker.com/q/).
- Browse the logs (in `User\AppData\Local\Docker`) by clicking **log file** in the Diagnose & Feedback window.
- Ask questions on the [Docker for Windows forum](https://forums.docker.com/c/docker-for-windows).
- Upload diagnostics in the Diagnose & Feedback window. You'll get a unique ID in return. You can then use this ID to submit issues at the [Docker for Windows GitHub repo](https://github.com/docker/for-win/issues).

![Diagnose & Feedback with ID](images/diagnostic-id.png){:width="500px"}

## Troubleshooting topics

### Make sure certificates are set up correctly

Docker for Windows ignores certificates listed under insecure registries, and
does not send client certificates to them. Commands like `docker run` that
attempt to pull from the registry produces error messages on the command line,
like this:

```
Error response from daemon: Get http://192.168.203.139:5858/v2/: malformed HTTP response "\x15\x03\x01\x00\x02\x02"
```

As well as on the registry. For example:

```
2017/06/20 18:15:30 http: TLS handshake error from 192.168.203.139:52882: tls: client didn't provide a certificate
2017/06/20 18:15:30 http: TLS handshake error from 192.168.203.139:52883: tls: first record does not look like a TLS handshake
```

For more about using client and server side certificates, see [How do I add
custom CA certificates?](index.md#how-do-i-add-custom-ca
certificates) and [How do I add client
certificates?](index.md#how-do-i-add-client-certificates) in
the Getting Started topic.

### Volumes

#### Permissions errors on data directories for shared volumes

Docker for Windows sets permissions on [shared
volumes](index.md#shared-drives) to a default value of
[0755](http://permissions-calculator.org/decode/0755/) (`read`, `write`,
`execute` permissions for `user`, `read` and `execute` for `group`). If you are
working with applications that require permissions different than this default,
you may get errors similar to the following.

```none
Data directory (/var/www/html/data) is readable by other users. Please change the permissions to 0755 so that the directory cannot be listed by other users.
```

The default permissions on shared volumes are not configurable. If you are
working with applications that require permissions different from the shared
volume defaults at container runtime, you need to either use non-host-mounted
volumes or find a way to make the applications work with the default file
permissions.

Docker for Windows currrently implements host-mounted volumes based on the [Microsoft SMB protocol](https://msdn.microsoft.com/en-us/library/windows/desktop/aa365233(v=vs.85).aspx), which does not support fine-grained, `chmod` control over these permissions.

See also, [Can I change permissions on shared volumes for container-specific deployment requirements?](faqs.md#can-i-change-permissions-on-shared-volumes-for-container-specific-deployment-requirements) in the FAQs, and for more of an explanation, the GitHub issue, [Controlling Unix-style perms on directories passed through from shared Windows drives](https://github.com/docker/docker.github.io/issues/3298).

#### inotify on shared drives does not work

Currently, `inotify` does not work on Docker for Windows. This becomes evident,
for example, when an application needs to read/write to a container across a
mounted drive. Instead of relying on filesystem inotify, we recommend using
polling features for your framework or programming language.

* **Workaround for nodemon and Node.js** - If you are using [nodemon](https://github.com/remy/nodemon) with `Node.js`, try the fallback polling mode described here: [nodemon isn't restarting node applications](https://github.com/remy/nodemon#application-isnt-restarting)

* **Docker for Windows issue on GitHub** - See the issue [Inotify on shared drives does not work](https://github.com/docker/for-win/issues/56#issuecomment-242135705)

#### Volume mounting requires shared drives for Linux containers

If you are using mounted volumes and get runtime errors indicating an application file is not found, access is denied to a volume mount, or a service cannot start, such as when using [Docker Compose](/compose/gettingstarted.md), you might need to enable [shared drives](index.md#shared-drives).

Volume mounting requires shared drives for Linux containers (not for Windows
containers). Go to ![whale menu](/docker-for-mac/images/whale-x.png){: .inline}
--> **Settings** --> **Shared Drives** and share the drive that contains the
Dockerfile and volume.

#### Verify domain user has permissions for shared drives (volumes)

> **Tip**: Shared drives are only required for volume mounting [Linux
> containers](index.md#switch-between-windows-and-linux-containers),
> not Windows containers.

Permissions to access shared drives are tied to the username and password you
use to set up [shared drives](index.md#shared-drives). If
you run `docker` commands and tasks under a different username than the one used
to set up shared drives, your containers don't have permissions to access the
mounted volumes.  The volumes show as empty.

The solution to this is to switch to the domain user account and reset
credentials on shared drives.

Here is an example of how to de-bug this problem, given a scenario where you
shared the `C` drive as a local user instead of as the domain user. Assume the
local user is `samstevens` and the domain user is `merlin`.

1. Make sure you are logged in as the Windows domain user (for our example, `merlin`).

2. Run `net share c` to view user permissions for `<host>\<username>, FULL`.

   ```
   > net share c

   Share name        C
   Path              C:\
   Remark
   Maximum users     No limit
   Users             SAMSTEVENS
   Caching           Caching disabled
   Permission        windowsbox\samstevens, FULL
   ```

3. Run the following command to remove the share.

   ```
   > net share c /delete
   ```

4. Re-share the drive via the [Shared Drives
   dialog](index.md#shared-drives), and provide the Windows
   domain user account credentials.

5. Re-run `net share c`.

   ```
   > net share c

   Share name        C
   Path              C:\
   Remark
   Maximum users     No limit
   Users             MERLIN
   Caching           Caching disabled
   Permission        windowsbox\merlin, FULL
   ```

See also, the related issue on GitHub, [Mounted volumes are empty in the
container](https://github.com/docker/for-win/issues/25).

#### Volume mounts from host paths use a `nobrl` option to override database locking

You may encounter problems using volume mounts on the host, depending on the
database software and which options are enabled. Docker for Windows uses
[SMB/CIFS
protocols](https://msdn.microsoft.com/en-us/library/windows/desktop/aa365233(v=vs.85).aspx)
to mount host paths, and mounts them with the `nobrl` option, which prevents
lock requests from being sent to the database server
([docker/for-win#11](https://github.com/docker/for-win/issues/11),
[docker/for-win#694](https://github.com/docker/for-win/issues/694)). This is
done to ensure container access to database files shared from the host. Although
it solves the over-the-network database access problem, this "unlocked" strategy
can interfere with other aspects of database functionality (for example,
write-ahead logging (WAL) with SQLite, as described in
[docker/for-win#1886](https://github.com/Sonarr/Sonarr/issues/1886)).

If possible, avoid using shared drives for volume mounts on the host with
network paths, and instead mount on the MobyVM, or create a [data
volume](/engine/tutorials/dockervolumes.md#data-volumes) (named volume) or [data
container](/engine/tutorials/dockervolumes.md#creating-and-mounting-a-data-volume-container).
See also, the [volumes key under service
configuration](/compose/compose-file/index.md#volumes) and the [volume
configuration
reference](/compose/compose-file/index.md#volume-configuration-reference) in the
Compose file documentation.

#### Local security policies can block shared drives and cause login errors

You need permissions to mount shared drives to use the Docker for Windows
[shared drives](index.md#shared-drives) feature.

If local policy prevents this, you get errors when you attempt to enable shared
drives on Docker. This is not something Docker can resolve, since you do need
these permissions to use the feature.

Here are snip-its from example error messages:

```none
Logon failure: the user has not been granted the requested logon type at
this computer.

[19:53:26.900][SambaShare     ][Error  ] Unable to mount C drive: mount
error(5): I/O error Refer to the mount.cifs(8) manual page (e.g. man mount.cifs)
mount: mounting //10.0.75.1/C on /c failed: Invalid argument
```

See also, <a href="https://github.com/docker/for-win/issues/98">Docker for Windows issue #98</a>.

#### Understand symlinks limitations

Symlinks work within and across containers. However, symlinks created outside of
containers (for example, on the host) do not work. To learn more, see [Are
symlinks supported?](faqs.md#are-symlinks-supported) in the FAQs.

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

### Virtualization

In order for Docker for Windows to function properly your machine needs:

1. [Hyper-V](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/hyper-v-technology-overview) installed and working

2. Virtualization enabled


#### Hyper-V


Docker for Windows requires a Hyper-V as well as the Hyper-V Module for Windows
Powershell to be installed and enabled. The Docker for Windows installer enables
it for you.

See [these instructions](https://msdn.microsoft.com/en-us/virtualization/hyperv_on_windows/quick_start/walkthrough_install) to install Hyper-V manually. A reboot is *required*. If you install Hyper-V without the reboot, Docker for Windows does not work correctly. On some systems, Virtualization needs to be enabled in the BIOS. The steps to do so are Vendor specific, but typically the BIOS option is called `Virtualization Technology (VTx)` or similar.


From the start menu, type  in "Turn Windows features on or off" and hit enter. In the subequent screen, verify Hyper-V is enabled and has a checkmark:

![Hyper-V on Windows features](images/hyperv-enabled.png){:width="600px"}

#### Hyper-V driver for Docker Machine

Docker for Windows comes with the legacy tool Docker Machine which uses the old
[`boot2docker.iso`](https://github.com/boot2docker/boot2docker){:
target="_blank" class="_"}, and the [Microsoft Hyper-V
driver](/machine/drivers/hyper-v.md) to create local virtual machines. _This is
tangential to using Docker for Windows_, but if you want to use Docker Machine
to create multiple local VMs, or to provision remote machines, see the [Docker
Machine](/machine/index.md) topics. We mention this here only in case someone is
looking for information about Docker Machine on Windows, which requires that
Hyper-V is enabled, an external network switch is active, and referenced in the
flags for the `docker-machine create` command [as described in the Docker
Machine driver example](/machine/drivers/hyper-v.md#example).

#### Virtualization must be enabled

In addition to [Hyper-V](#hyper-v), virtualization must be enabled. Check the Performance tab on the
Task Manager:

![Task Manager](images/virtualization-enabled.png){:width="700px"}

If, at some point, if you manually uninstall Hyper-V or disable virtualization,
Docker for Windows cannot start. See: [Unable to run Docker for Windows on Windows 10 Enterprise](https://github.com/docker/for-win/issues/74).

### Networking and WiFi problems upon Docker for Windows install

Some users have encountered networking issues during install and startup of
Docker for Windows. For example, upon install or auto-reboot, network adapters
and/or WiFi gets disabled. In some scenarios, problems are due to having
VirtualBox or its network adapters still installed, but in other scenarios this
is not the case. (See also, Docker for Windows issue on GitHub:
[Enabling Hyper-V feature turns my wi-fi off
](https://github.com/docker/for-win/issues/139).)

Here are some steps to take if you encounter similar problems:

1.  Ensure **virtualization** is enabled, as described above in [Virtualization must be enabled](#virtualization-must-be-enabled).

2.  Ensure **Hyper-V** is installed and enabled, as described above in [Hyper-V must be enabled](#hyper-v-must-be-enabled).

3.  Ensure **DockerNAT** is enabled by checking the **Virtual Switch Manager** on the Actions tab on the right side of the **Hyper-V Manager**.

    ![Hyper-V manager](images/hyperv-manager.png)

4.  Set up an external network switch. If you plan at any point to use
    [Docker Machine](/machine/overview.md) to set up multiple local VMs, you
    need this anyway, as described in the topic on the
    [Hyper-V driver for Docker Machine](/machine/drivers/hyper-v.md#example).
    You can replace `DockerNAT` with this switch.

5.  If previous steps fail to solve the problems, follow steps on the
    [Cleanup README](https://github.com/Microsoft/Virtualization-Documentation/blob/master/windows-server-container-tools/CleanupContainerHostNetworking/README.md).

    > Read full description before you run Windows cleanup script
    >
    > The cleanup command has two flags, `-Cleanup` and `-ForceDeleteAllSwitches`.
    > Read the whole page before running any scripts, especially warnings about `-ForceDeleteAllSwitches`.
    > {: .warning}

### Windows containers and Windows Server 2016

If you have questions about how to set up and run Windows containers on Windows
Server 2016 or Windows 10, see [About Windows containers and Windows Server 2016](index.md#about-windows-containers-and-windows-server-2016).

A full tutorial is available in [docker/labs](https://github.com/docker/labs) at
[Getting Started with Windows Containers](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md).

You can install a native Windows binary which allows you to develop and run
Windows containers without Docker for Windows. However, if you install Docker
this way, you cannot develop or run Linux containers. If you try to run a Linux
container on the native Docker daemon, an error occurs:

```none
C:\Program Files\Docker\docker.exe:
 image operating system "linux" cannot be used on this platform.
 See 'C:\Program Files\Docker\docker.exe run --help'.
```

### Limitations of Windows containers for `localhost` and published ports

Docker for Windows provides the option to switch Windows and Linux containers.
If you are using Windows containers, keep in mind that there are some
limitations with regard to networking due to the current implementation of
Windows NAT (WinNAT). These limitations may potentially resolve as the Windows
containers project evolves.

One thing you may encounter rather immediately is that published ports on
Windows containers do not do loopback to the local host. Instead, container
endpoints are only reachable from the host using the container's IP and port.

So, in a scenario where you use Docker to pull an image and run a webserver with
a command like this:

```shell
> docker run -d -p 80:80 --name webserver nginx
```

Using `curl http://localhost`, or pointing your web browser at
`http://localhost` does not display the `nginx` web page (as it would do with
Linux containers).

To reach a Windows container from the local host, you need to specify
the IP address and port for the container that is running the service.

You can get the container IP address by using [`docker inspect`](/engine/reference/commandline/inspect.md) with some `--format` options and the ID or name of the container. For the example above, the command would look like this, using the name we gave to the container (`webserver`) instead of the container ID:

{% raw %}
```bash
$ docker inspect \
  --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' \
  webserver
```
{% endraw %}

This gives you the IP address of the container, for example:

{% raw %}
```bash
$ docker inspect \
  --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' \
  webserver

172.17.0.2
```
{% endraw %}

Now you can connect to the webserver by using `http://172.17.0.2:80` (or simply
`http://172.17.0.2`, since port `80` is the default HTTP port.)

For more information, see:

* Docker for Windows issue on GitHub: [Port binding does not work for
  locahost](https://github.com/docker/for-win/issues/458)

* [Published Ports on Windows Containers Don't Do
  Loopback](https://blog.sixeyed.com/published-ports-on-windows-containers-dont-do-loopback/)

* [Windows NAT capabilities and
  limitations](https://blogs.technet.microsoft.com/virtualization/2016/05/25/windows-nat-winnat-capabilities-and-limitations/)


### Running Docker for Windows in nested virtualization scenarios

Docker for Windows can run inside a Windows 10 virtual machine (VM) running on apps like Parallels or VMware Fusion on a Mac provided that the VM is properly configured. However, problems and intermittent failures may still occur due to the way these apps virtualize the hardware. For these reasons, _**Docker for Windows is not supported for nested virtualization scenarios**_. It might work in some cases, and not in others.

The better solution is to run Docker for Windows natively on a Windows system
(to work with Windows or Linux containers), or Docker for Mac on Mac to work
with Linux containers.

#### If you still want to use nested virtualization

* Make sure nested virtualization support is enabled in VMWare or Parallels.
Check the settings in **Hardware -> CPU & Memory -> Advanced Options -> Enable
nested virtualization** (the exact menu sequence might vary slightly).

* Configure your VM with at least 2 CPUs and sufficient memory to run your
  workloads.

* Make sure your system is more or less idle.

* Make sure your Windows OS is up-to-date. There have been several issues with
  some insider builds.

* The processor you have may also be relevant. For example, Westmere based Mac
  Pros have some additional hardware virtualization features over Nehalem based
  Mac Pros and so do newer generations of Intel processors.

#### Typical failures we see with nested virtualization

* Slow boot time of the Linux VM. If you look in the logs and find some entries
  prefixed with `Moby`. On real hardware, it takes 5-10 seconds to boot the
  Linux VM; roughly the time between the `Connected` log entry and the `*
  Starting Docker ... [ ok ]` log entry. If you boot the Linux VM inside a
  Windows VM, this may take considerably longer. We have a timeout of 60s or
  so. If the VM hasn't started by that time, we retry. If the retry fails we
  print an error. You can sometimes work around this by providing more resources
  to the Windows VM.

* Sometimes the VM fails to boot when Linux tries to calibrate the time stamp
  counter (TSC). This process is quite timing sensitive and may fail when
  executed inside a VM which itself runs inside a VM. CPU utilization is also
  likely to be higher.

#### Related issues

Discussion thread on GitHub at [Docker for Windows issue 267](https://github.com/docker/for-win/issues/267)

### Networking issues

Some users have reported problems connecting to Docker Hub on the Docker for Windows stable version. (See GitHub issue
[22567](https://github.com/moby/moby/issues/22567).)

Here is an example command and error message:

```shell
> docker run hello-world

Unable to find image 'hello-world:latest' locally
Pulling repository docker.io/library/hello-world
C:\Program Files\Docker\Docker\Resources\bin\docker.exe: Error while pulling image: Get https://index.docker.io/v1/repositories/library/hello-world/images: dial tcp: lookup index.docker.io on 10.0.75.1:53: no such host.
See 'C:\Program Files\Docker\Docker\Resources\bin\docker.exe run --help'.
```

As an immediate workaround to this problem, reset the DNS server to use the
Google DNS fixed address: `8.8.8.8`. You can configure this via the **Settings**
-> **Network** dialog, as described in the topic
[Network](index.md#network). Docker automatically restarts
when you apply this setting, which could take some time.

We are currently investigating this issue.

### NAT/IP configuration

By default, Docker for Windows uses an internal network prefix of
`10.0.75.0/24`. Should this clash with your normal network setup, you can change
the prefix from the **Settings** menu. See the [Network](index.md#network) topic
under [Settings](index.md#docker-settings).

## Workarounds

### `inotify` currently does not work on Docker for Windows

If you are using `Node.js` with `nodemon`, a temporary workaround is to try the
fallback polling mode described here:
[nodemon isn't restarting node applications](https://github.com/remy/nodemon#application-isnt-restarting). See
also this issue on GitHub
[Inotify on shared drives does not work](https://github.com/docker/for-win/issues/56#issuecomment-242135705).

### Reboot

Restart your PC to stop / discard any vestige of the daemon running from the
previously installed version.

### Unset `DOCKER_HOST`

The `DOCKER_HOST` environmental variable does not need to be set.  If you use
bash, use the command `unset ${!DOCKER_*}` to unset it.  For other shells,
consult the shell's documentation.

### Make sure Docker is running for webserver examples

For the `hello-world-nginx` example and others, Docker for Windows must be
running to get to the webserver on `http://localhost/`. Make sure that the
Docker whale is showing in the menu bar, and that you run the Docker commands in
a shell that is connected to the Docker for Windows Engine (not Engine from
Toolbox). Otherwise, you might start the webserver container but get a "web page
not available" error when you go to `docker`.

### How to solve `port already allocated` errors

If you see errors like `Bind for 0.0.0.0:8080 failed: port is already allocated`
or `listen tcp:0.0.0.0:8080: bind: address is already in use` ...

These errors are often caused by some other software on Windows using those
ports. To discover the identity of this software, either use the `resmon.exe`
GUI and click "Network" and then "Listening Ports" or in a powershell use
`netstat -aon | find /i "listening "` to discover the PID of the process
currently using the port (the PID is the number in the rightmost column). Decide
whether to shut the other process down, or to use a different port in your
docker app.

### Docker fails to start when firewall or anti-virus software is installed

**Some firewalls and anti-virus software might be incompatible with Microsoft
**Windows 10 builds**, such as Windows 10 Anniversary Update. The conflict
typically occurs after a Windows update or new install of the firewall, and
manifests as an error response from the Docker daemon and a **Docker for Windows
start failure**. The Comodo Firewall was one example of this problem, but users
report that software has since been updated to work with these Windows 10
builds.

See the Comodo forums topics
[Comodo Firewall conflict with Hyper-V](https://forums.comodo.com/bug-reports-cis/comodo-firewall-began-conflict-with-hyperv-t116351.0.html)
and
[Windows 10 Anniversary build doesn't allow Comodo drivers to be installed](https://forums.comodo.com/install-setup-configuration-help-cis/windows-10-aniversary-build-doesnt-allow-comodo-drivers-to-be-installed-t116322.0.html).
A Docker for Windows user-created issue describes the problem specifically as it
relates to Docker:
[Docker fails to start on Windows 10](https://github.com/docker/for-win/issues/27).

For a temporary workaround, uninstall the firewall or anti-virus software, or
explore other workarounds suggested on the forum.
