---
description: Work with swarms in Docker Cloud
keywords: swarm mode, swarms, orchestration Cloud, fleet management
title: Swarms in Docker Cloud
---

<br>

<b>Note</b>: All Swarm management features in Docker Cloud are free while in Beta.
{: style="text-align:center" }

--------
Docker Cloud now allows you to connect to clusters of Docker Engines running in v1.13 [Swarm Mode](/engine/swarm/).

## Enable Swarm Mode in Docker Cloud

By default, Docker Cloud allows you to manage [node clusters](/docker-cloud/infrastructure/), but you can opt in to use the Beta
Swarm Mode features. Swarm management replaces the node cluster management
features when the Swarm Mode interface is enabled.

Click the **Swarm Mode** toggle to enable the swarm mode interface.

![the Swarm Mode toggle](images/swarm-toggle.png)

You can switch between node cluster and swarm mode at any time, and enabling swarm mode does _not_ remove or disconnect existing node clusters.

## Swarm mode and organizations

If you use Docker Cloud in an [organization](/docker-cloud/orgs/), you can
use Swarm Mode to access any Docker swarms available to your organization.
Members of the `owners` team grant each team in an organization access to the
swarms they need. If necessary, you can create new teams to manage beta swarm
access.

If you use swarm mode as a member of a team other than the `owners` team,
you will only see the swarms that you have been granted access to.

Members of the `owners` team must switch to the Swarm Mode Docker Cloud
interface to grant teams access to an organization's swarms. Swarms only appear
in the [resource management](/docker-cloud/orgs/#/set-team-permissions) screens
for teams when in the swarm mode interface.

## Register an existing swarm

At this time, you cannot _create_ swarms from within Docker Cloud. However you can register existing swarms as part of the beta functionality.

Before you begin, you need the following:

- a Docker ID
- a Docker swarm composed of v1.13 (or later) Docker Engine nodes
- a terminal session connected to one of the swarm's manager nodes
- incoming port 2376 unblocked on that manager node

> **Note**: The IP to the manager node for your swarm must be open and publicly accessible so that Docker Cloud can connect and run commands.

To register an existing swarm in Docker Cloud:

1. Log in to Docker Cloud if necessary.
2. If necessary, click the **Swarm mode** toggle to activate the Swarm Mode interface.
3. Click **Swarms** in the top navigation.
4. Click **Bring your own swarm**
5. Select the whole command displayed in the dialog, and copy it to your clipboard.
6. In terminal or another shell, connect to the Docker Engine running in the swarm's manager node using SSH.
7. Paste the command you copied into the terminal session connected to the manager node.
8.  When prompted, log in using your Docker ID and password.

    The registration process uses your Docker ID to determine which namespaces you have access to<!--are allowed to register the swarm under TODO:CLOUD-4079 -->. Once you log in, the CLI lists these namespaces to help you with the next step.

9.  Enter a name, with a namespace before the name if needed, and press Enter.

    If you do not enter a name, the swarm is registered to your Docker ID account using the swarm ID, which the long string displayed before the shell prompt. For example, the prompt might look like this:

    ```none
    Enter a name for the new cluster [mydockerid/5rdshkgzn1sw016zimgckzx3j]:
    ```

    Enter a name at the prompt to prevent Docker Cloud from registering the swarm using the long swarm ID as the name.

    To register a swarm with an organization, prefix the new name with the organization name, for example `myorganization/myteamswarm`.

The manager node pulls the `dockercloud/registration` container which creates a
global service called `dockercloud-server-proxy`. This service runs on _all_ of
the swarm's manager nodes.

The swarm then appears in the **Swarms** screen in Docker Cloud.

### Swarm Registration example

```none
$ docker run -ti --rm -v /var/run/docker.sock:/var/run/docker.sock dockercloud/registration
Use your Docker ID credentials to authenticate:
Username: myusername
Password:

Available namespaces:
* myorganization
* pacificocean
* sealife
Enter name for the new cluster [myusername/1btbwtge4xwjj0mjpdpr7jutn]: myusername/myswarm
Registering this Docker cluster with Docker Cloud...
Successfully registered the node as myswarm
You can now access this cluster using the following command in any Docker Engine:
	docker run --rm -ti -v /var/run/docker.sock:/var/run/docker.sock -e DOCKER_HOST dockercloud/client myswarm
```

## Swarm statuses in Docker Cloud

Swarms that are registered in Docker Cloud appear in the Swarms list. Each line in the list shows the swarm's status. The statuses are:

<!-- TODO - **DEPLOYING**: Docker Cloud is provisioning the nodes of this swarm. -->
- **DEPLOYED**: the swarm is sending heartbeat pings to Docker Cloud, and Cloud can contact it to run a health check.
- **UNREACHABLE**: the swarm is sending heartbeart pings, but Docker Cloud cannot contact the swarm.
- **UNAVAILABLE**: Docker Cloud is not receiving any heartbeats from the swarm.
- **REMOVED**: the swarm has been unregistered from Docker Cloud and will be removed from the list soon.

> **Note**: [Removing a swarm](#unregister-a-swarm-from-Docker-cloud) only removes the swarm from the interface in Docker Cloud. It does not change the swarm itself or any processes running on the swarm.

## Connect to a swarm through Docker Cloud

Docker Cloud allows you to connect your local Docker Engine to any swarm you
have access to in Docker Cloud. To do this, you run a proxy container in your local Docker instance, which connects to a manager node on the target swarm.

1. Log in to Docker Cloud in your web browser.
2. Click **Swarms** in the top navigation, and click the name of the swarm you want to connect to.
3. Copy the command provided in the dialog that appears.
4. In a terminal window connected to your local Docker Engine, paste the command, and press **Enter**.

    The local Docker Engine downloads a containerized Docker Cloud client tool, and connects to the swarm.

5. To complete the connection process, run the `export DOCKER_HOST` command found in the previous command's output, to connect your local shell to the client proxy.

    Be sure to include the client connection port in the URL. For example `export DOCKER_HOST=tcp://127.0.0.1:32768`.


To switch Docker hosts you can either run the `export` command again to overwrite it, or use `unset DOCKER_HOST`.

> **Note**: If you are using Docker Machine, be sure to unset `DOCKER_TLS_VERIFY` as described in the [known issues](https://github.com/docker/dockercloud-federation#known-issues).

## Unregister a swarm from Docker Cloud

Unregistering a swarm from Docker Cloud only removes the swarm from Docker
Cloud, deletes any access rights granted to teams, and disables proxy
connections. Unregistering does not stop the services, containers, or processes on the swarm, and it does not disband the swarm or terminate the nodes.

To unregister a swarm from Docker Cloud:

1. Log in to Docker Cloud if necessary.
2. Click **Swarms** in the top navigation.
3. Put your mouse cursor on the swarm you want to unregister.
4. Click the trash can icon that appears.
5. In the confirmation dialog that appears, click **Remove**.

Docker Cloud marks the swarm as `REMOVED` and removes the swarm from the list in the next few minutes.  

## Reconnect a swarm

If you accidentally unregister a swarm from Docker Cloud, or decide that you
want to re-register the swarm after it has been removed, you can re-register it
using the same process as a normal registration. If the swarm is registered to
an organization, its access permissions were deleted when it was unregistered,
and must be recreated.

> **Note**: You cannot register a new or different swarm under the name of a
swarm that was unregistered. To re-register a swarm, it must have the same swarm
ID as it did when previously registered.
