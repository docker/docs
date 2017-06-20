---
description: How to install Docker Machine
keywords: machine, orchestration, install, installation, docker, documentation
title: Install Docker Machine
---

On macOS and Windows, Machine is installed along with other Docker products when
you install the [Docker for Mac](/docker-for-mac/index.md), [Docker for
Windows](/docker-for-windows/index.md), or [Docker
Toolbox](/toolbox/overview.md).

If you want only Docker Machine, you can install the Machine binaries directly
by following the instructions in the next section. You can find the latest
versions of the binaries on the <a
href="https://github.com/docker/machine/releases/" target="_blank">
docker/machine release page</a> on GitHub.

## Installing Machine Directly

1.  Install <a href="/engine/installation/"
    target="_blank">the Docker binary</a>.

2.  Download the Docker Machine binary and extract it to your PATH.

    If you are running on **macOS**:

    ```console
    $ curl -L https://github.com/docker/machine/releases/download/v0.12.0/docker-machine-`uname -s`-`uname -m` >/usr/local/bin/docker-machine && \
  chmod +x /usr/local/bin/docker-machine
    ```

    If you are running on **Linux**:

    ```console
    $ curl -L https://github.com/docker/machine/releases/download/v0.12.0/docker-machine-`uname -s`-`uname -m` >/tmp/docker-machine &&
    chmod +x /tmp/docker-machine &&
    sudo cp /tmp/docker-machine /usr/local/bin/docker-machine
    ```

    If you are running with **Windows** with git bash:

    ```console
    $ if [[ ! -d "$HOME/bin" ]]; then mkdir -p "$HOME/bin"; fi && \
curl -L https://github.com/docker/machine/releases/download/v0.12.0/docker-machine-Windows-x86_64.exe > "$HOME/bin/docker-machine.exe" && \
chmod +x "$HOME/bin/docker-machine.exe"
    ```

    Otherwise, download one of the releases from the <a href="https://github.com/docker/machine/releases/" target="_blank"> docker/machine release page</a> directly.

3.  Check the installation by displaying the Machine version:

        $ docker-machine version
        docker-machine version 0.12.0, build 45c69ad

## Installing bash completion scripts

The Machine repository supplies several `bash` scripts that add features such
as:

-   command completion
-   a function that displays the active machine in your shell prompt
-   a function wrapper that adds a `docker-machine use` subcommand to switch the
    active machine

To install the scripts, copy or link them into your `/etc/bash_completion.d` or
`/usr/local/etc/bash_completion.d` directory. To enable the `docker-machine` shell
prompt, add `$(__docker_machine_ps1)` to your `PS1` setting in `~/.bashrc`.

    PS1='[\u@\h \W$(__docker_machine_ps1)]\$ '

You can find additional documentation in the comments at the <a
href="https://github.com/docker/machine/tree/master/contrib/completion/bash"
target="_blank">top of each script</a>.

### How to uninstall

To uninstall Docker Machine:

*  Remove the executable: `rm $(which docker-machine)`

*  Optionally, remove the machines you created.

    To remove each machine individually: `docker-machine rm <machine-name>`

    To remove all machines: `docker-machine rm -f $(docker-machine ls -q)`

  Removing machines is an optional step because there are cases where you might
  want to save and migrate existing machines to a [Docker for
  Mac](/docker-for-mac/index.md) or [Docker for
  Windows](/docker-for-windows/index.md) environment, for example.

## Where to go next

-   [Docker Machine overview](overview.md)
-   Create and run a Docker host on your [local system using virtualization](get-started.md)
-   Provision multiple Docker hosts [on your cloud provider](get-started-cloud.md)
-   <a href="/machine/drivers/" target="_blank">Docker Machine driver reference</a>
-   <a href="/machine/reference/" target="_blank">Docker Machine subcommand reference</a>
