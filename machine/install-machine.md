---
description: How to install Docker Machine
keywords: machine, orchestration, install, installation, docker, documentation, uninstall Docker Machine, uninstall
title: Install Docker Machine
---

{% assign machineversion = '0.12.2' %}

On macOS and Windows, Machine is installed along with other Docker products when
you install the Docker Toolbox. For details on installing Docker Toolbox, see
the <a href="https://docs.docker.com/installation/mac/" target="_blank">macOS
installation</a> instructions or <a
href="https://docs.docker.com/docker-for-windows/" target="_blank">Windows
installation</a> instructions.

If you want only Docker Machine, you can install the Machine binaries directly
by following the instructions in the next section. You can find the latest
versions of the binaries on the [docker/machine release
page](https://github.com/docker/machine/releases/){: target="_blank" class="_" }
on GitHub.

## Installing Machine directly

<<<<<<< HEAD
1.  Install [Docker](/engine/installation/index.md){: target="_blank" class="_" }.
=======
1.  Install <a href="https://docs.docker.com/installation/"
    target="_blank">the Docker binary</a>.
>>>>>>> FreeBSD: add build support

2.  Download the Docker Machine binary and extract it to your PATH.

    If you are running on **macOS**:

    ```console
    $ curl -L https://github.com/docker/machine/releases/download/v{{machineversion}}/docker-machine-`uname -s`-`uname -m` >/usr/local/bin/docker-machine && \
  chmod +x /usr/local/bin/docker-machine
    ```

    If you are running on **Linux**:

    ```console
    $ curl -L https://github.com/docker/machine/releases/download/v{{machineversion}}/docker-machine-`uname -s`-`uname -m` >/tmp/docker-machine &&
    chmod +x /tmp/docker-machine &&
    sudo cp /tmp/docker-machine /usr/local/bin/docker-machine
    ```

<<<<<<< HEAD
    If you are running with **Windows** with [Git BASH](https://git-for-windows.github.io/){: target="_blank" class="_"}:
=======
    If you are running Windows with git bash
>>>>>>> FreeBSD: add build support

    ```console
    $ if [[ ! -d "$HOME/bin" ]]; then mkdir -p "$HOME/bin"; fi && \
curl -L https://github.com/docker/machine/releases/download/v{{machineversion}}/docker-machine-Windows-x86_64.exe > "$HOME/bin/docker-machine.exe" && \
chmod +x "$HOME/bin/docker-machine.exe"
    ```

    > The above command will work on Windows only if you use a
    terminal emulater such as [Git BASH](https://git-for-windows.github.io/){: target="_blank" class="_"}, which supports Linux commands like `chmod`.
    {: .important}

    Otherwise, download one of the releases from the [docker/machine release
    page](https://github.com/docker/machine/releases/){: target="_blank" class="_" } directly.

3.  Check the installation by displaying the Machine version:

        $ docker-machine version
        docker-machine version {{machineversion}}, build 9371605

## Installing bash completion scripts

The Machine repository supplies several `bash` scripts that add features such
as:

-   command completion
-   a function that displays the active machine in your shell prompt
-   a function wrapper that adds a `docker-machine use` subcommand to switch the
    active machine

Confirm the version and save scripts to `/etc/bash_completion.d` or
`/usr/local/etc/bash_completion.d`:

```bash
scripts=( docker-machine-prompt.bash docker-machine-wrapper.bash docker-machine.bash ); for i in "${scripts[@]}"; do sudo wget https://raw.githubusercontent.com/docker/machine/v{{machineversion}}/contrib/completion/bash/${i} -P /etc/bash_completion.d; done
```

To enable the `docker-machine` shell
prompt, add `$(__docker_machine_ps1)` to your `PS1` setting in `~/.bashrc`.

```
PS1='[\u@\h \W$(__docker_machine_ps1)]\$ '
```

<<<<<<< HEAD
You can find additional documentation in the comments at the [top of each script](https://github.com/docker/machine/tree/master/contrib/completion/bash){: target="_blank" class="_"}.

## How to uninstall Docker Machine

To uninstall Docker Machine:

*  Remove the executable: `rm $(which docker-machine)`

*  Optionally, remove the machines you created.

    To remove each machine individually: `docker-machine rm <machine-name>`

    To remove all machines: `docker-machine rm -f $(docker-machine ls -q)` (you might need to use `-force` on Windows)

  Removing machines is an optional step because there are cases where you might
  want to save and migrate existing machines to a [Docker for
  Mac](/docker-for-mac/index.md) or [Docker for
  Windows](/docker-for-windows/index.md) environment, for example.
=======
You can find additional documentation in the comments at the <a href="https://github.com/docker/machine/tree/master/contrib/completion/bash" target="_blank">top of each script</a>.
>>>>>>> FreeBSD: add build support

>**Note**: As a point of information, the `config.json`, certificates,
and other data related to each virtual machine created by `docker-machine`
is stored in `~/.docker/machine/machines/` on Mac and Linux and in
`~\.docker\machine\machines\` on Windows. We recommend that you do not edit or
remove those files directly as this will only affect information for the Docker
CLI, not the actual VMs, regardless of whether they are local or on remote
servers.

## Where to go next

-   [Docker Machine overview](overview.md)
-   Create and run a Docker host on your [local system using virtualization](get-started.md)
-   Provision multiple Docker hosts [on your cloud provider](get-started-cloud.md)
<<<<<<< HEAD
-  [Docker Machine driver reference](/machine/drivers/index.md)
-  [Docker Machine subcommand reference](/machine/reference/index.md)
=======
-   <a href="https://docs.docker.com/machine/drivers/" target="_blank">Docker Machine driver reference</a>
-   <a href="https://docs.docker.com/machine/reference/" target="_blank">Docker Machine subcommand reference</a>
>>>>>>> FreeBSD: add build support
