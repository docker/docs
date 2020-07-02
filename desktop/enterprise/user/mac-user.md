---
title: Use Docker Desktop Enterprise on Mac
keywords: Docker EE, Windows, Mac, Docker Desktop, Enterprise
redirect_from:
- /ee/desktop/user/mac-user/
---

This page contains information about testing the installation and configuring Docker Desktop Enterprise (DDE) runtime options on Mac.

## Test your installation

1. Open a command-line terminal and test that your installation works by
    running the simple Docker image, [hello-world](https://hub.docker.com/_/hello-world/).

    ```shell
    $ docker run hello-world

    Unable to find image 'hello-world:latest' locally
    latest: Pulling from library/hello-world
    ca4f61b1923c: Pull complete
    Digest: sha256:ca0eeb6fb05351dfc8759c20733c91def84cb8007aa89a5bf606bc8b315b9fc7
    Status: Downloaded newer image for hello-world:latest

    Hello from Docker!
    This message shows that your installation appears to be working correctly.
    ...
    ```

2. Start a Dockerized web server. Like the `hello-world` image above, if the
    image is not found locally, Docker pulls it from Docker Hub.

   ```bash
    $ docker run --detach --publish=80:80 --name=webserver nginx
    ```

3. In a web browser, go to `http://localhost/` to view the nginx homepage.
    Because we specified the default HTTP port, it isn't necessary to append
    `:80` at the end of the URL.

    ![nginx home page](../images/hello-world-nginx.png)

4. View the details on the container while your web server is running (with
    `docker container ls` or `docker ps`):

    ```none
    $ docker container ls
    CONTAINER ID   IMAGE   COMMAND                  CREATED              STATUS              PORTS                         NAMES
    56f433965490   nginx   "nginx -g 'daemon off"   About a minute ago   Up About a minute   0.0.0.0:80->80/tcp, 443/tcp   webserver
    ```

5. Stop and remove containers and images with the following commands. Use the
    "all" flag (`--all` or `-a`) to view stopped containers.

 ```shell
$ docker container ls
$ docker container stop webserver
$ docker container ls -a
$ docker container rm webserver
$ docker image ls
$ docker image rm nginx
```

## Docker Desktop user interface

The Docker Desktop Enterprise user interface provides options to configure Docker Desktop preferences such as installation, version packs, Docker Hub login, and more. Right-click the Docker icon from the menu bar to open the Docker Desktop user interface.

### Version Selection

The **Version Selection** option lists the version packs installed on your Docker Desktop environment and allows you to switch between Docker Engine and Kubernetes versions using a single click. When you select a different version pack, Docker Desktop restarts and the selected versions of Docker Engine and Kubernetes will be used.

To switch to a different version pack, simply click on the version pack you would like to use.

### Preferences

Click on the Docker icon from the menu bar and then **Preferences** to configure the runtime options described below.

> **Note:** Administrators have the ability to lock some configuration options. Locked options cannot be selected, and are displayed with a lock icon.

![Docker context menu](../images/prefs.png)

#### General

![Preferences](../images/./prefs-general.png)

General settings include:

- **Start Docker Desktop when you log in:** Starts Docker Desktop when you open your session. (Enabled by default)

- **Include VM in Time Machine backups:** Backs up the Docker Desktop Enterprise virtual machine. (Disabled by default)

   **Securely store Docker logins in macOS keychain:** Stores your Docker login
  credentials. (Enabled by default)

- **Send usage statistics:** Sends diagnostics, crash reports, and usage
  data to Docker. This information helps Docker improve the application and get
  more context for troubleshooting problems. (Enabled by default)

#### File Sharing

Choose which local directories to share with your containers. File sharing is
required for volume mounting if the project lives outside of the `/Users`
directory. In that case, share the drive where the Dockerfile and volume are
located. Otherwise, you get `file not found` or `cannot start service` errors at
runtime.

![File Sharing](../images/prefs-fileshare.png)

File sharing settings include the following options:

- **Add a Directory**: Click `+` and navigate to the directory you want to add.

- **Apply & Restart** makes the directory available to containers using Docker's
  bind mount (`-v`) feature.

  There are some limitations on the directories that can be shared:

  - They cannot be a subdirectory of a directory that has been shared already.
  - They cannot already exist inside of Docker.

For more information, see:

- [Namespaces](https://docs.docker.com/docker-for-mac/osxfs/#namespaces) in [osxfs file system sharing](https://docs.docker.com/docker-for-mac/osxfs/).

- [Volume mounting requires file sharing for any project directories outside of `/Users`](https://docs.docker.com/docker-for-mac/troubleshoot/#volume-mounting-requires-file-sharing-for-any-project-directories-outside-of-users).

#### Disk

Specify the **Disk image location** of the Linux volume, where containers and
images are stored.

You can also move the disk image location. If you attempt to move the disk image
to a location that already has one, you get a prompt asking if you want to use
the existing image or replace it.

![Disk settings](../images/prefs-disk.png)

#### Advanced

On the Advanced tab, you can limit resources available to Docker.

![Advanced Preference settings](../images/prefs-advanced.png)

Advanced settings include the following options:

- **CPUs**: By default, Docker Desktop Enterprise is set to use half the number of processors
available on the host machine. To increase processing power, set this to a
higher number; to decrease, lower the number.

- **Memory**: By default, Docker Desktop Enterprise is set to use 2 GB runtime memory,
allocated from the total available memory on your Mac. To increase RAM, set this
to a higher number; to decrease it, lower the number.

- **Swap**: Configure swap file size as needed. The default is 1 GB.

#### Proxies

Docker Desktop Enterprise detects HTTP/HTTPS Proxy Settings from macOS and automatically
propagates these to Docker and to your containers. For example, if you set your
proxy settings to `http://proxy.example.com`, Docker uses this proxy when
pulling containers.

![Proxies settings](../images/prefs-proxies.png)

![macOS Proxy Settings](../images/proxy-settings.png)

When you start a container, your proxy settings propagate into the containers.
For example:

```
$ docker run -it alpine env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=b7edf988b2b5
TERM=xterm
HOME=/root
HTTP_PROXY=http://proxy.example.com:3128
http_proxy=http://proxy.example.com:3128
no_proxy=*.local, 169.254/16
```

You can see from the above output that the `HTTP_PROXY`, `http_proxy`, and
`no_proxy` environment variables are set. When your proxy configuration changes,
Docker restarts automatically to pick up the new settings. If you have
containers that you wish to keep running across restarts, you should consider
using [restart policies](https://docs.docker.com/engine/reference/run/#restart-policies---restart).

#### Daemon

You can configure options on the Docker daemon that determine how your
containers run.

Select **Basic** to configure the daemon with interactive settings, or select
**Advanced** to edit the JSON file directly.

![Daemon](../images/prefs-daemon-basic.png)

##### Experimental features

Docker Desktop Enterprise has experimental features enabled
on Docker Engine, as described in [Docker Experimental Features](https://github.com/docker/cli/blob/master/experimental/README.md) Readme. If you don't select **Experimental Features**, Docker Desktop Enterprise uses the current generally available release of Docker Engine.

> {% include experimental.md %}

You can see whether you are running experimental mode at the command line. If
`Experimental` is `true`, then Docker is running in experimental mode, as shown
here. (If `false`, Experimental mode is off.)

```bash
{% raw %}$ docker version -f {{.Server.Experimental}}{% endraw %}
true
```

##### Insecure registries

You can set up a custom and insecure [registry](https://docs.docker.com/registry/introduction) to store your public or private images (instead of
using [Docker Hub](https://hub.docker.com/). Add URLs for
your insecure registries and registry mirrors on which to host your images.

For more information, see:

- [How do I add custom CA certificates?](https://docs.docker.com/docker-for-mac/faqs/#how-do-i-add-custom-ca-certificates)

- [How do I add client certificates?](https://docs.docker.com/docker-for-mac/faqs/#how-do-i-add-client-certificates)

##### Daemon configuration file

Click the **Advanced** tab to configure the daemon from the JSON file. For a
full list of options, see the Docker Engine [dockerd command line reference](https://docs.docker.com/engine/reference/commandline/dockerd).

Click **Apply & Restart** to save your settings and reboot Docker. To cancel
changes, click another preference tab, then choose to discard or not apply
changes when prompted.

![Docker Daemon](../images/prefs-daemon-adv.png)

#### Kubernetes

Docker Desktop Enterprise includes a standalone Kubernetes server that runs on your Mac, so
that you can test deploying your Docker workloads on Kubernetes.

The Kubernetes client command, `kubectl`, is included and configured to connect
to the local Kubernetes server. If you have `kubectl` already installed and
pointing to some other environment, such as `minikube` or a GKE cluster, be sure
to change context so that `kubectl` is pointing to `docker-desktop`:

```bash
$ kubectl config get-contexts
$ kubectl config use-context docker-desktop
```

If you installed `kubectl` with Homebrew, or by some other method, and
experience conflicts, remove `/usr/local/bin/kubectl`.

To enable Kubernetes support and install a standalone instance of Kubernetes running as a Docker container, select **Enable Kubernetes**, choose the [default orchestrator](https://docs.docker.com/docker-for-mac/kubernetes/#override-the-default-orchestrator) and click the **Apply** button.

![Enable Kubernetes](../images/prefs-kubernetes.png)

Images required to run the Kubernetes server are downloaded and instantiated as containers, and the `/usr/local/bin/kubectl` command is installed on your Mac.

When Kubernetes is enabled and running, an additional status bar item displays at the bottom right of the Docker Desktop Enterprise **Preferences** dialog.

![Installation complete](../images/kubernetes-install-complete.png)

The status of Kubernetes shows in the Docker menu and the context points to `docker-desktop`.

![Docker Menu with Kubernetes](../images/kube-context.png)

By default, Kubernetes containers are hidden from commands like `docker
  service ls`, because managing them manually is not supported. To view these containers, select **Show system containers (advanced)** and click **Apply and restart**. Most users do not have to use this option.

To disable Kubernetes support at any time, clear the **Enable Kubernetes** check box. The
Kubernetes containers are stopped and removed, and the
`/usr/local/bin/kubectl` command is removed.

For more information about using the Kubernetes integration with Docker Desktop Enterprise, see [Deploy on Kubernetes](https://docs.docker.com/docker-for-mac/kubernetes).

#### Reset

Click on the Docker icon from the menu bar and then **Preferences**. Click **Reset** to reset to factory defaults, restart the Docker daemon, reset Kubernetes cluster, or to reset the disk image.

![Uninstall or reset Docker](../images/prefs-reset-mac.png)

### Troubleshoot

The **Diagnose and Feedback** option allows you troubleshoot any issues you may be experiencing with Docker Desktop Enterprise. For more information, see [Troubleshoot DDE issues on Mac](/desktop/enterprise/troubleshoot/mac-issues).

### Design new application

Select this option to open the Application Designer user interface. Application Designer provides a library of application and service templates to help Docker developers quickly create new Docker applications. For more information, see [Application Designer](/desktop/enterprise/app-designer).

### Docker Hub

Select **Sign in /Create Docker ID** from the Docker Desktop Enterprise menu to access your [Docker Hub](https://hub.docker.com/) account. Once logged in, select **Repositories**  on the Docker Desktop Enterprise menu to access your Docker Hub repositories directly.

See the following [Docker Hub topics](https://docs.docker.com/docker-hub/) to learn more:

- [Organizations and Teams in Docker Hub](https://docs.docker.com/docker-hub/orgs/)

- [Builds and Images](https://docs.docker.com/docker-hub/official_images/)

### Add TLS certificates

You can add trusted Certificate Authorities (CAs) (used to verify registry
server certificates) and client certificates (used to authenticate to
registries) to your Docker daemon.

#### Add custom CA certificates (server side)

All trusted CAs (root or intermediate) are supported. Docker Desktop Enterprise creates a
certificate bundle of all user-trusted CAs based on the Mac Keychain, and
appends it to Moby trusted certificates. So if an enterprise SSL certificate is
trusted by the user on the host, it is trusted by Docker Desktop Enterprise.

To manually add a custom, self-signed certificate, start by adding the
certificate to the macOS keychain, which is picked up by Docker Desktop Enterprise. Here is
an example.

```bash
$ sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ca.crt
```

If you prefer to add the certificate to your own local keychain only (rather
than for all users), run this command instead:

```
$ security add-trusted-cert -d -r trustRoot -k ~/Library/Keychains/login.keychain ca.crt
```

See also, [Directory structures for
certificates](#directory-structures-for-certificates).

> **Note:** You need to restart Docker Desktop Enterprise after making any changes to the
keychain or to the `~/.docker/certs.d` directory in order for the changes to
take effect.

For a complete explanation of how to do this, see the blog post [Adding
Self-signed Registry Certs to Docker & Docker Desktop Enterprise](http://container-solutions.com/adding-self-signed-registry-certs-docker-mac/).

#### Add client certificates

You can put your client certificates in
`~/.docker/certs.d/<MyRegistry>:<Port>/client.cert` and
`~/.docker/certs.d/<MyRegistry>:<Port>/client.key`.

When the Docker Desktop Enterprise application starts up, it copies the `~/.docker/certs.d`
folder on your Mac to the `/etc/docker/certs.d` directory on Moby (the Docker
Desktop Enterprise `xhyve` virtual machine).

> * You need to restart Docker Desktop Enterprise after making any changes to the keychain
>   or to the `~/.docker/certs.d` directory in order for the changes to take
>   effect.
>
> * The registry cannot be listed as an _insecure registry_ (see [Docker Daemon](#daemon)). Docker Desktop Enterprise ignores certificates listed
>   under insecure registries, and does not send client certificates. Commands
>   like `docker run` that attempt to pull from the registry produce error
>   messages on the command line, as well as on the registry.

#### Directory structures for certificates

If you have this directory structure, you do not need to manually add the CA
certificate to your macOS system login:

```
/Users/<user>/.docker/certs.d/
└── <MyRegistry>:<Port>
   ├── ca.crt
   ├── client.cert
   └── client.key
```

The following further illustrates and explains a configuration with custom
certificates:

```
/etc/docker/certs.d/        <-- Certificate directory
└── localhost:5000          <-- Hostname:port
   ├── client.cert          <-- Client certificate
   ├── client.key           <-- Client key
   └── ca.crt               <-- Certificate authority that signed
                                the registry certificate
```

You can also have this directory structure, as long as the CA certificate is
also in your keychain.

```
/Users/<user>/.docker/certs.d/
└── <MyRegistry>:<Port>
    ├── client.cert
    └── client.key
```

To learn more about how to install a CA root certificate for the registry and
how to set the client TLS certificate for verification, see [Verify repository client with certificates](https://docs.docker.com/engine/security/certificates) in the Docker Engine
topics.

### Install shell completion

Docker Desktop Enterprise comes with scripts to enable completion for `docker` and `docker-compose` commands. The completion scripts may be
found inside `Docker.app`, in the `Contents/Resources/etc/` directory and can be
installed both in Bash and Zsh.

#### Bash

Bash has [built-in support for completion](https://www.debian-administration.org/article/316/An_introduction_to_bash_completion_part_1). To activate completion for Docker commands, these files need to be
copied or symlinked to your `bash_completion.d/` directory. For example, if you have
installed bash through [Homebrew](http://brew.sh/).

```bash
etc=/Applications/Docker.app/Contents/Resources/etc
ln -s $etc/docker.bash-completion $(brew --prefix)/etc/bash_completion.d/docker
ln -s $etc/docker-compose.bash-completion $(brew --prefix)/etc/bash_completion.d/docker-compose
```

#### Zsh

In Zsh, the [completion
system](http://zsh.sourceforge.net/Doc/Release/Completion-System.html) takes care of things. To activate completion for Docker commands,
these files need to be copied or symlinked to your Zsh `site-functions/`
directory. For example, if you installed Zsh through [Homebrew](http://brew.sh/):

```bash
etc=/Applications/Docker.app/Contents/Resources/etc
ln -s $etc/docker.zsh-completion /usr/local/share/zsh/site-functions/_docker
ln -s $etc/docker-compose.zsh-completion /usr/local/share/zsh/site-functions/_docker-compose
```
