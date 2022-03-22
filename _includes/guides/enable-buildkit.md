<!-- This text will be included in Build images topic in the Get started guides -->

### Enable BuildKit

Before we start building images, ensure you have enabled BuildKit on your machine.
BuildKit allows you to build Docker images efficiently. For more information,
see [Building images with BuildKit](/develop/develop-images/build_enhancements/).

BuildKit is enabled by default for all users on Docker Desktop. If you have
installed Docker Desktop, you don't have to manually enable BuildKit. If you are
running Docker on Linux, you can enable BuildKit either by using an environment
variable or by making BuildKit the default setting.

To set the BuildKit environment variable when running the `docker build` command,
run:

```console
$ DOCKER_BUILDKIT=1 docker build .
```

To enable docker BuildKit by default, set daemon configuration in `/etc/docker/daemon.json` feature to `true` and restart the daemon.
If the `daemon.json` file doesn't exist, create new file called `daemon.json` and then add the following to the file.

```json
{
  "features":{"buildkit" : true}
}
```

Restart the Docker daemon.
