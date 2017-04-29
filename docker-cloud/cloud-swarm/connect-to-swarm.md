---
previewflag: cloud-swarm
description: how to register and unregister swarms in Docker Cloud
keywords: swarm mode, swarms, orchestration Cloud, fleet management
title: Connect to a swarm through Docker Cloud
---

Docker Cloud allows you to connect your local Docker Engine to any swarm you
have access to in Docker Cloud. To do this, you run a proxy container in your
local Docker instance, which connects to a manager node on the target swarm.

## Connect to a swarm

1.  Log in to Docker Cloud in your web browser.
2.  Click **Swarms** in the top navigation, and click the name of the swarm you want to connect to.
3.  Copy the command provided in the dialog that appears.

    ![](images/swarm-connect.png)

4.  In a terminal window connected to your local Docker Engine, paste the command, and press **Enter**.

    The local Docker Engine downloads a containerized Docker Cloud client tool, and connects to the swarm.

5.  To complete the connection process, run the `export DOCKER_HOST` command found in the previous command's output, to connect your local shell to the client proxy.

    Be sure to include the client connection port in the URL. For example,  `export DOCKER_HOST=tcp://127.0.0.1:32768`.


6.  You can run `docker node ls` to verify that the swarm is running.

    ```
    $ docker node ls
    ID                            HOSTNAME              STATUS              AVAILABILITY        MANAGER STATUS
    6uotpiv8vyxsjzdtux13nkvj4     swarm-worker000001    Ready               Active              
    8j3e54j3vv16scrrvcyg5de2w     swarm-worker000002    Ready               Active              
    qmvk4swo9rdv1viu9t88dw0t3     swarm-worker000000    Ready               Active              
    w7kgzzdkka0k2svssz1dk1fzw *   swarm-manager000000   Ready               Active              Leader
    ```

    From this point on, you can use the
    [CLI commands](/engine/swarm/index.md#swarm-mode-cli-commands)
    to manage your cloud-hosted [swarm mode](/engine/swarm/) just as you
    would a local swarm.

> **Note**: To switch back to Docker hosts you can either run the `export` command again to overwrite it, or use `unset DOCKER_HOST`. If you are using Docker Machine, be sure to unset `DOCKER_TLS_VERIFY` as described in the [known issues](https://github.com/moby/mobycloud-federation#known-issues).

## Reconnect a swarm

If you accidentally unregister a swarm from Docker Cloud, or decide that you
want to re-register the swarm after it has been removed, you can
[re-register it](register-swarms.md#register-a-swarm) using the same
process as a normal registration. If the swarm is registered to
an organization, its access permissions were deleted when it was
unregistered, and must be recreated.

> **Note**: You cannot register a new or different swarm under the name of a
swarm that was unregistered. To re-register a swarm, it must have the same swarm
ID as it did when previously registered.

## Where to go next

Learn how to [create a new swarm in Docker Cloud](create-cloud-swarm.md).
