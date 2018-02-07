---
description: Learn how to troubleshoot your Docker Universal Control Plane cluster.
keywords: docker, ucp, troubleshoot
title: Troubleshoot your cluster
---

If you detect problems in your UCP cluster, you can start your troubleshooting
session by checking the logs of the
[individual UCP components](../../architecture.md). Only administrator users can
see information about UCP system containers.

## Check the logs from the UI

To see the logs of the UCP system containers, navigate to the **Containers**
page of UCP. By default the UCP system containers are hidden. Click the
**Show all containers** option for the UCP system containers to be listed as
well.

![](../../images/troubleshoot-with-logs-1.png){: .with-border}

You can click on a container to see more details like its configurations and
logs.


## Check the logs from the CLI

You can also check the logs of UCP system containers from the CLI. This is
specially useful if the UCP web application is not working.

1. Get a client certificate bundle.

    When using the Docker CLI client you need to authenticate using client
    certificates.
    [Learn how to use client certificates](../../user/access-ucp/cli-based-access.md).

    If your client certificate bundle is for a non-admin user, you don't have
    permissions to see the UCP system containers.

2.  Check the logs of UCP system containers.

    ```bash
    # By default system containers are not displayed. Use the -a flag to display them
    $ docker ps -a

    CONTAINER ID    IMAGE                             COMMAND                  CREATED         STATUS           PORTS                            NAMES
    922503c2102a    docker/ucp-controller:1.1.0-rc2   "/bin/controller serv"   4 hours ago     Up 30 minutes    192.168.10.100:444->8080/tcp     ucp/ucp-controller
    1b6d429f1bd5    docker/ucp-swarm:1.1.0-rc2        "/swarm join --discov"   4 hours ago     Up 4 hours       2375/tcp                         ucp/ucp-swarm-join

    # See the logs of the ucp/ucp-controller container
    $ docker logs ucp/ucp-controller

    {"level":"info","license_key":"PUagrRqOXhMH02UgxWYiKtg0kErLY8oLZf1GO4Pw8M6B","msg":"/v1.22/containers/ucp/ucp-controller/json","remote_addr":"192.168.10.1:59546","tags":["api","v1.22","get"],"time":"2016-04-25T23:49:27Z","type":"api","username":"dave.lauper"}
    {"level":"info","license_key":"PUagrRqOXhMH02UgxWYiKtg0kErLY8oLZf1GO4Pw8M6B","msg":"/v1.22/containers/ucp/ucp-controller/logs","remote_addr":"192.168.10.1:59546","tags":["api","v1.22","get"],"time":"2016-04-25T23:49:27Z","type":"api","username":"dave.lauper"}
    ```

## Get a support dump

Before making any changes to UCP, download a [support dump](../../get-support.md).
This allows you to troubleshoot problems which were already happening before
changing UCP configurations.

Then you can increase the UCP log level to debug, making it easier to understand
the status of the UCP cluster. Changing the UCP log level restarts all UCP
system components and introduces a small downtime window to UCP. Your
applications won't be affected by this.

To increase the UCP log level, navigate to the **UCP web UI**, go to the
**Admin Settings** tab, and choose **Logs**.

![](../../images/troubleshoot-with-logs-2.png){: .with-border}

Once you change the log level to **Debug** the UCP containers are restarted.
Now that the UCP components are creating more descriptive logs, you can download
again a support dump and use it to troubleshoot the component causing the
problem.

Depending on the problem you are experiencing, it's more likely that you'll
find related messages in the logs of specific components on manager nodes:

* If the problem occurs after a node was added or removed, check the logs
of the `ucp-reconcile` container.
* If the problem occurs in the normal state of the system, check the logs
of the `ucp-controller` container.
* If you can browse to the UCP web UI but can't log in, check the
logs of the `ucp-auth-api` and `ucp-auth-store` containers.

It's normal for the `ucp-reconcile` container to be in a stopped state. This
container is only started when the `ucp-agent` detects that a node needs to
transition to a different state, and it is responsible for creating and removing
containers, issuing certificates, and pulling missing images.


## Where to go next

* [Troubleshoot configurations](troubleshoot-configurations.md)
