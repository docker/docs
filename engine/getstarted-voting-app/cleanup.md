---
description: Graceful shutdown, reboot, clean-up
keywords: voting app, docker-machine
title: Graceful shutdown, reboot, and clean-up
---


The voting app will continue to run on the swarm while the `manager` and `worker` machines are running, unless you explicitly stop it.

## Stopping the voting app

To shut down the voting app, simply stop the machines on which it is running. If you are using local hosts, follow the steps below. If you are using cloud hosts, stop them per your cloud setup.

1.  Open a terminal window and run `docker-machine ls` to list the current machines.

    ```
    $ docker-machine ls
    NAME      ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER    ERRORS
    manager   -        virtualbox   Running   tcp://192.168.99.100:2376           v1.13.1   
    worker    -        virtualbox   Running   tcp://192.168.99.101:2376           v1.13.1   
    ```
2.  Use `docker-machine stop` to shut down each machine, beginning with the worker.

    ```
    $ docker-machine stop worker
    Stopping "worker"...
    Machine "worker" was stopped.

    $ docker-machine stop manager
    Stopping "manager"...
    Machine "manager" was stopped.
    ```

## Restarting the voting app

If you want to come back to your `manager` and `worker` machines later, you can
keep them around. One advantage of this is that you can simply restart the
machines to launch the sample voting app again.

To restart local machines, follow the steps below. To restart cloud instances,
start them per your cloud setup.

1.  Open a terminal window and list the machines.

    ```
    $ docker-machine ls
    NAME      ACTIVE   DRIVER       STATE     URL   SWARM   DOCKER    ERRORS
    manager   -        virtualbox   Stopped                 Unknown   
    worker    -        virtualbox   Stopped                 Unknown   
    ```

3.  Run `docker-machine start` to start each machine, beginning with the manager.

    ```
    $ docker-machine start manager
    Starting "manager"...
    (manager) Check network to re-create if needed...
    (manager) Waiting for an IP...
    Machine "manager" was started.
    Waiting for SSH to be available...
    Detecting the provisioner...
    Started machines may have new IP addresses. You may need to re-run the `docker-machine env` command.

    $ docker-machine start worker
    Starting "worker"...
    (worker) Check network to re-create if needed...
    (worker) Waiting for an IP...
    Machine "worker" was started.
    Waiting for SSH to be available...
    Detecting the provisioner...
    Started machines may have new IP addresses. You may need to re-run the `docker-machine env` command.
    ```

3.  Run the following commands to log into the manager and see if the swarm is up.

    ```
    docker-machine ssh manager

    docker@manager:~$ docker stack services vote
    ID            NAME             MODE        REPLICAS  IMAGE
    74csdxb99tg9  vote_visualizer  replicated  1/1       dockersamples/visualizer:stable
    jm0g1vahcid9  vote_redis       replicated  2/2       redis:alpine
    mkk6lee494t4  vote_db          replicated  1/1       postgres:9.4
    o3sl1wr35yd6  vote_worker      replicated  1/1       dockersamples/examplevotingapp_worker:latest
    qcc8dw2zafc1  vote_vote        replicated  2/2       dockersamples/examplevotingapp_vote:after
    x5wcvknlnnh7  vote_result      replicated  1/1       dockersamples/examplevotingapp_result:after
    ```

At this point, the app is back up. The web pages you looked at in the [test drive](test-drive.md) should be available, and you could experiment, modify the app, and [redeploy](customize-app.md).

## Removing the machines

If you prefer to remove your local machines altogether, use `docker-machine rm`
to do so. (Or, `docker-machine rm -f` will force-remove running machines.)

```
$ docker-machine rm worker
About to remove worker
WARNING: This action will delete both local reference and remote instance.
Are you sure? (y/n): y
Successfully removed worker

$ docker-machine rm manager
About to remove manager
WARNING: This action will delete both local reference and remote instance.
Are you sure? (y/n): y
Successfully removed manager
```

The Docker images you pulled were all running on the virtual machines you
created (either local or cloud), so no other cleanup of images or processes is
needed once you stop and/or remove the virtual machines.

## What's next?

See the [Docker Machine topics](/machine/overview/) for more on working
with `docker-machine`.

Check out the [list of resources](customize-app.md#resources) for more on Docker
labs, sample apps, and swarm mode.
