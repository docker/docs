To start Docker Desktop for Linux, search **Docker Desktop** on the
**Applications** menu and open it. This launches the Docker menu icon and opens
the Docker Dashboard, reporting the status of Docker Desktop.

Alternatively, open a terminal and run:

```console
$ systemctl --user start docker-desktop
```

When Docker Desktop starts, it creates a dedicated [context](/engine/context/working-with-contexts) that the Docker CLI
can use as a target and sets it as the current context in use. This is to avoid
a clash with a local Docker Engine that may be running on the Linux host and
using the default context. On shutdown, Docker Desktop resets the current
context to the previous one.

The Docker Desktop installer updates Docker Compose and the Docker CLI binaries
on the host. It installs Docker Compose V2 and gives users the choice to
link it as docker-compose from the Settings panel. Docker Desktop installs
the new Docker CLI binary that includes cloud-integration capabilities in `/usr/local/bin/com.docker.cli`
and creates a symlink to the classic Docker CLI at `/usr/local/bin`.

After you’ve successfully installed Docker Desktop, you can check the versions
of these binaries by running the following commands:

```console
$ docker compose version
Docker Compose version v2.17.3

$ docker --version
Docker version 23.0.5, build bc4487a

$ docker version
Client: Docker Engine - Community
 Cloud integration: v1.0.31
 Version:           23.0.5
 API version:       1.42
<...>
```

To enable Docker Desktop to start on sign in, from the Docker menu, select
**Settings** > **General** > **Start Docker Desktop when you sign in to your computer**.

Alternatively, open a terminal and run:

```console
$ systemctl --user enable docker-desktop
```

To stop Docker Desktop, select the Docker menu icon to open the Docker menu and select **Quit Docker Desktop**.

Alternatively, open a terminal and run:

```console
$ systemctl --user stop docker-desktop
```