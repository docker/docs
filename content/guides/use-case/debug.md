---
description: Learn how to debug containers.
keywords: debug, container access
title: Debug containers
toc_max: 2
---

Debugging is a critical aspect of software development, but accessing containers
to run debugging tools can present several challenges due to the following
factors:

- Containers are isolated. While this isolation simplifies dependency
  management, it can complicate debugging because traditional tools may not be
  available.
- Containers are slim. To minimize their size and attack surface, containers
  often include only the binaries necessary to run the application, which can
  limit the availability of debugging tools.
- Containers are ephemeral. Containers are often designed to be short-lived and
  stateless, making it difficult to access logs and other diagnostic data for
  post-mortem analysis.

This guide introduces you to a few commands that will help you access your
containers to debug. By the end of this guide, you’ll be able to access and
debug:

- Running containers
- Containers without shell
- Stopped containers

## Access a running container

With Docker, there are a few ways you can access a container that’s running your
application, including [`docker exec`](/reference/cli/docker/container/exec/)
and [`docker debug`](/reference/cli/docker/debug/). The `docker exec` command is
available to all Docker users, while `docker debug` requires a paid
subscription. The following instructions show you how to use each command to
access and debug a running container.

1. Run the following command to start an [Nginx](https://hub.docker.com/_/nginx)
   container, which you will then access to simulate a debugging scenario.

   ```console
   $ docker run --rm -d --name my-container nginx
   ```

2. Access the container using the following `docker exec` or `docker debug`
   command.

   > [!NOTE]
   >
   > In Docker Desktop, you can access the container using `docker exec` or
   > `docker debug` via the graphical user interface. In the **Containers**
   > view, under the **Actions** column, select **Show container actions**.
   > Then, select **Open in terminal** to connect via `docker exec`, or **Use
   > Docker Debug** to connect via `docker debug`.

   {{< tabs >}}
   {{< tab name="docker debug" >}}
   
   ```console
   $ docker debug my-container
   ```

   {{< /tab >}}
   {{< tab name="docker exec" >}}

   ```console
   $ docker exec -it my-container bash
   ```

   The `-it` flags are used to start an interactive terminal session, and `bash`
   specifies the shell to be used inside the container.

   {{< /tab >}}
   {{< /tabs >}}


   Best practices dictate that containers shouldn’t contain unnecessary
   packages, so your container may not have the tools you’re familiar with or
   need to debug it. For example, the Nginx container doesn’t contain the `ps`
   program to list processes. To use tools with `docker exec`, you need to
   install the tool in the container, which can increase the attack surface. For
   `docker debug`, it contains a toolbox of tools that you can use without
   altering the container.

3. Install if necessary, and then run `ps` to list the container's processes.
   Select how you accessed the container in the following tabs, and then following the instructions to run `ps`.

   {{< tabs >}}
   {{< tab name="docker debug" >}}
   
   ```console
   $ ps aux
   ```

   Docker debug comes with many tools already available in its toolbox, including `ps`. To install other tools into the toolbox, use the `install` command. For example to install `nmap`, run the following command:
  
   ```console
   $ install nmap
   ```

   {{< /tab >}}
   {{< tab name="docker exec" >}}

   Use apt-get to install `procps`.

   ```console
   $ apt-get update && apt-get install -y procps
   ```

   After installation, you can run `ps` in the container.

   ```console
   $ ps aux
   ```

   {{< /tab >}}
   {{< /tabs >}}

Unlike `docker exec`, the `docker debug` tools are available in the
container, but are not installed in the container. This means the tools
available via `docker debug` don’t increase size nor the attack surface of
the container. In addition, the tools in the toolbox are available in every
container you access without the need to reinstall them.

## Access a container with no shell

With Docker, there are a few ways you can access a container with no shell. The
quickest and easiest way is to use `docker debug`. Another way is to manually
start a second container in the same pid and network namespace. The following
instructions cover both methods.

1. Run the following command to start an [NATS](https://hub.docker.com/_/nats)
   container, which you will then access to simulate a debugging scenario. The
   NATS image has no shell and you won’t be able to connect to it using `docker
   exec`.

   ```console
   $ docker run --rm -d --name my-container-2 nats
   ```

2. Access the container using `docker debug` or a secondary container.

   {{< tabs >}}
   {{< tab name="docker debug" >}}
   
   ```console
   $ docker debug my-container-3
   ```

   {{< /tab >}}
   {{< tab name="Secondary container" >}}

   ```console
   $  docker run -it --rm --pid=container:my-container-2 --network=container:my-container-2 busybox sh
   ```

   This command starts a busybox container in the same pid and network namespace
   as the container named my-container-2, and then connects to the busybox
   container.

   {{< /tab >}}
   {{< /tabs >}}

3. Run your debugging tools. In this scenario, you’ll run `ps` to list the
   processes.

   ```console
   $ ps aux
   ```

   You should see at least the following process running, indicating that you can see the processes in the NATS container.

   ```text
   /nats-server --config nats-server.conf
   ```

4. Access the container’s filesystem.

   {{< tabs >}}
   {{< tab name="docker debug" >}}
   
   `docker debug` gives you quick access to the filesystem. Run `ls` to see the contents.

   ```console
   $ ls
   ```

   You should see the following contents, indicating that you can access the filesystem in the NATS container.

   ```text
   dev  etc  nats-server  nats-server.conf  nix  proc  sys
   ```

   {{< /tab >}}
   {{< tab name="Secondary container" >}}

    After running `ps` in the previous step, you see that NATS is running as pid 1 and the user is root. You can access the file system at `/proc/1/root/`. Use the following command to list the contents.

   ```console
   $ ls /proc/1/root/
   ```
    You should see the following contents, indicating that you can access the filesystem in the NATS container.

   ```text
   dev  etc  nats-server  nats-server.conf  proc  sys
   ```

   {{< /tab >}}
   {{< /tabs >}}

## Access a stopped container

With Docker, there are a few ways you can access the contents of a stopped
container. The quickest and easiest way is to use `docker debug`. Another way is
to copy the files out of the container. The following instructions cover both
methods.

1. Run the following commands to create and then stop an
   [Nginx](https://hub.docker.com/_/nginx) container, which you will then access
   to simulate a debugging scenario.

   ```console
   $ docker run -d --name my-container-3 nginx
   $ docker stop my-container-3
   ```

2. Inspect the `entrypoint.sh` file by accessing the stopped container using
   `docker debug` or by [copying](/reference/cli/docker/container/cp/) the contents out of the container.


   {{< tabs >}}
   {{< tab name="docker debug" >}}
   
   Access the stopped container with `docker debug`.

   ```console
   $ docker debug my-container-3
   ```

   Then, use `cat` to inspect the `docker-entrypoint.sh` file in the container.

   ```console
   $ cat docker-entrypoint.sh
   ```

   {{< /tab >}}
   {{< tab name="Copy files" >}}

   The following command copies the `/docker-entrypoint.sh` file from `my-container-3` to your current directory (`.`).

   ```console
   $ docker cp my-container-3:/docker-entrypoint.sh .
   ```

   You can then use a tool on your own system to open the
   `docker-entrypoint.sh` file and inspect it. For example `cat
   docker-entrypoint.sh` or `notepad docker-entrypoint.sh` on Windows.

   {{< /tab >}}
   {{< /tabs >}}

## Summary

Debugging containers can be challenging due to their isolated, slim, and
ephemeral nature. However, with the right tools and techniques, you can
efficiently diagnose and resolve issues. This guide has equipped you with
essential methods for:

- Debugging a running container using `docker exec` or `docker debug`.
- Debugging containers with no shell by using `docker debug` or a secondary
  container.
- Debugging a stopped container to inspect its contents with `docker debug` or
  by copying files.

The `docker debug` command stands out by providing a powerful toolbox without
altering the container, ensuring both efficiency and security in your debugging
process. With these techniques, you’re now prepared to effectively manage and
debug your Docker containers.

Related information:
- [`docker debug` reference](/reference/cli/docker/debug/)
- [`docker exec` reference](/reference/cli/docker/container/exec/)
- [`docker cp` reference](/reference/cli/docker/container/cp/)
- [`docker run` reference](/reference/cli/docker/container/run/)