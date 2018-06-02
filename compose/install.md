---
description: How to install Docker Compose
keywords: compose, orchestration, install, installation, docker, documentation
title: Install Docker Compose
toc_max: 2
---

You can run Compose on macOS, Windows, and 64-bit Linux.

## Prerequisites

Docker Compose relies on Docker Engine for any meaningful work, so make sure you
have Docker Engine installed either locally or remote, depending on your setup.

- On desktop systems like Docker for Mac and Windows, Docker Compose is
included as part of those desktop installs.

- On Linux systems, first install the
[Docker](/install/index.md#server){: target="_blank" class="_"}
for your OS as described on the Get Docker page, then come back here for
instructions on installing Compose on
Linux systems.

- To run Compose as a non-root user, see [Manage Docker as a non-root user](/install/linux/linux-postinstall.md).

## Install Compose

Follow the instructions below to install Compose on Mac, Windows, Windows Server
2016, or Linux systems, or find out about alternatives like using the `pip`
Python package manager or installing Compose as a container.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#macOS">Mac</a></li>
<li><a data-toggle="tab" data-target="#windows">Windows</a></li>
<li><a data-toggle="tab" data-target="#linux">Linux</a></li>
<li><a data-toggle="tab" data-target="#alternatives">Alternative Install Options</a></li>
</ul>
<div class="tab-content">
<div id="macOS" class="tab-pane fade in active" markdown="1">
### Install Compose on macOS

**Docker for Mac** and **Docker Toolbox** already include Compose along
with other Docker apps, so Mac users do not need to install Compose separately.
Docker install instructions for these are here:

  * [Get Docker for Mac](/docker-for-mac/install.md)
  * [Get Docker Toolbox](/toolbox/overview.md) (for older systems)
<hr>
</div>
<div id="windows" class="tab-pane fade" markdown="1">
### Install Compose on Windows systems

**Docker for Windows** and **Docker Toolbox** already include Compose
along with other Docker apps, so most Windows users do not need to
install Compose separately. Docker install instructions for these are here:

* [Get Docker for Windows](/docker-for-windows/install.md)
* [Get Docker Toolbox](/toolbox/overview.md) (for older systems)

**If you are running the Docker daemon and client directly on Microsoft
Windows Server 2016** (with [Docker EE for Windows Server 2016](/install/windows/docker-ee.md), you _do_ need to install
Docker Compose. To do so, follow these steps:

1.  Start an "elevated" PowerShell (run it as administrator).
    Search for PowerShell, right-click, and choose
    **Run as administrator**. When asked if you want to allow this app
    to make changes to your device, click **Yes**.
    
    In Powershell, since Github now requires TLS1.2, run the following:
    
    ```none  
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    ```

    Then run the following command to download
    Docker Compose, replacing `$dockerComposeVersion` with the specific
    version of Compose you want to use:

    ```none
    Invoke-WebRequest "https://github.com/docker/compose/releases/download/$dockerComposeVersion/docker-compose-Windows-x86_64.exe" -UseBasicParsing -OutFile $Env:ProgramFiles\docker\docker-compose.exe
    ```

    For example, to download Compose version {{site.compose_version}},
    the command is:

    ```none
    Invoke-WebRequest "https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-Windows-x86_64.exe" -UseBasicParsing -OutFile $Env:ProgramFiles\docker\docker-compose.exe
    ```
    > Use the latest Compose release number in the download command.
    >
    > As already mentioned, the above command is an _example_, and
    it may become out-of-date once in a while. Always follow the
    command pattern shown above it. If you cut-and-paste an example,
    check which release it specifies and, if needed,
    replace `$dockerComposeVersion` with the release number that
    you want. Compose releases are also available for direct download
    on the [Compose repository release page on GitHub](https://github.com/docker/compose/releases){:target="_blank" class="_"}.
    {: .important}

2.  Run the executable to install Compose.
<hr>
</div>
<div id="linux" class="tab-pane fade" markdown="1">
### Install Compose on Linux systems

On **Linux**, you can download the Docker Compose binary from the [Compose
repository release page on GitHub](https://github.com/docker/compose/releases){:
target="_blank" class="_"}. Follow the instructions from the link, which involve
running the `curl` command in your terminal to download the binaries. These step
by step instructions are also included below.

1.  Run this command to download the latest version of Docker Compose:

    ```bash
    sudo curl -L https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
    ```

    > Use the latest Compose release number in the download command.
    >
    The above command is an _example_, and it may become out-of-date. To ensure you have the latest version, check the [Compose repository release page on GitHub](https://github.com/docker/compose/releases){: target="_blank" class="_"}.
    {: .important}

    If you have problems installing with `curl`, see
    [Alternative Install Options](install.md#alternative-install-options) tab above.

2.  Apply executable permissions to the binary:

    ```bash
    sudo chmod +x /usr/local/bin/docker-compose
    ```

3.  Optionally, install [command completion](completion.md) for the
    `bash` and `zsh` shell.

4.  Test the installation.

    ```bash
    $ docker-compose --version
    docker-compose version {{site.compose_version}}, build 1719ceb
    ```
<hr>
</div>
<div id="alternatives" class="tab-pane fade" markdown="1">
### Alternative install options

- [Install using pip](#install-using-pip)
- [Install as a container](#install-as-a-container)

#### Install using pip

Compose can be installed from
[pypi](https://pypi.python.org/pypi/docker-compose) using `pip`. If you install
using `pip`, we recommend that you use a
[virtualenv](https://virtualenv.pypa.io/en/latest/) because many operating
systems have python system packages that conflict with docker-compose
dependencies. See the [virtualenv
tutorial](http://docs.python-guide.org/en/latest/dev/virtualenvs/) to get
started.

```bash
pip install docker-compose
```
If you are not using virtualenv,

```bash
sudo pip install docker-compose
```

> pip version 6.0 or greater is required.

#### Install as a container

Compose can also be run inside a container, from a small bash script wrapper. To
install compose as a container run this command. Be sure to replace the version
number with the one that you want, if this example is out-of-date:

```bash
$ sudo curl -L --fail https://github.com/docker/compose/releases/download/{{site.compose_version}}/run.sh -o /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose
```

>  Use the latest Compose release number in the download command.
>
The above command is an _example_, and it may become out-of-date once in a
while. Check which release it specifies and, if needed, replace the given
release number with the one that you want. Compose releases are also listed and
available for direct download on the [Compose repository release page on
GitHub](https://github.com/docker/compose/releases){: target="_blank"
class="_"}.
{: .important}
<hr>
</div>
</div>

## Master builds

If you're interested in trying out a pre-release build, you can download a binary
from
[https://dl.bintray.com/docker-compose/master/](https://dl.bintray.com/docker-compose/master/).
Pre-release builds allow you to try out new features before they are released,
but may be less stable.


## Upgrading

If you're upgrading from Compose 1.2 or earlier, remove or
migrate your existing containers after upgrading Compose. This is because, as of
version 1.3, Compose uses Docker labels to keep track of containers, and your
containers need to be recreated to add the labels.

If Compose detects containers that were created without labels, it refuses
to run so that you don't end up with two sets of them. If you want to keep using
your existing containers (for example, because they have data volumes you want
to preserve), you can use Compose 1.5.x to migrate them with the following
command:

```bash
docker-compose migrate-to-labels
```

Alternatively, if you're not worried about keeping them, you can remove them.
Compose just creates new ones.

```bash
docker container rm -f -v myapp_web_1 myapp_db_1 ...
```

## Uninstallation

To uninstall Docker Compose if you installed using `curl`:

```bash
sudo rm /usr/local/bin/docker-compose
```

To uninstall Docker Compose if you installed using `pip`:

```bash
pip uninstall docker-compose
```

> Got a "Permission denied" error?
>
> If you get a "Permission denied" error using either of the above
> methods, you probably do not have the proper permissions to remove
> `docker-compose`. To force the removal, prepend `sudo` to either of the above
> commands and run again.


## Where to go next

- [User guide](index.md)
- [Getting Started](gettingstarted.md)
- [Get started with Django](django.md)
- [Get started with Rails](rails.md)
- [Get started with WordPress](wordpress.md)
- [Command line reference](/compose/reference/index.md)
- [Compose file reference](compose-file.md)
