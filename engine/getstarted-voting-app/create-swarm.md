---
description: Setup for voting app example
keywords: multi-container, services, swarm mode, cluster, voting app
title: Create a swarm
---

Now, we'll add our Docker machines to a [swarm](/engine/swarm/index.md).

## Initialize the swarm

1.  Log into the manager.

    ```
    $ docker-machine ssh manager
                            ##         .
                      ## ## ##        ==
                   ## ## ## ## ##    ===
               /"""""""""""""""""\___/ ===
          ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~ /  ===- ~~~
               \______ o           __/
                 \    \         __/
                  \____\_______/
     _                 _   ____     _            _
    | |__   ___   ___ | |_|___ \ __| | ___   ___| | _____ _ __
    | '_ \ / _ \ / _ \| __| __) / _` |/ _ \ / __| |/ / _ \ '__|
    | |_) | (_) | (_) | |_ / __/ (_| | (_) | (__|   <  __/ |
    |_.__/ \___/ \___/ \__|_____\__,_|\___/ \___|_|\_\___|_|

      WARNING: this is a build from test.docker.com, not a stable release.

    Boot2Docker version 1.13.0-rc7, build HEAD : b2cde29 - Sat Jan 14 00:29:39 UTC 2017
    Docker version 1.13.0-rc7, build 48a9e53
    docker@manager:~$ ls
    log.log
    docker@manager:~$ docker swarm init --advertise-addr 192.168.99.100
    Swarm initialized: current node (2tjrasfqfu945b7n4753374sw) is now a manager.
    ```

2.  Initialize a swarm.

    The command to initialize a swarm is:

    ```
    docker swarm init --advertise-addr <MANAGER-IP>
    ```

    Use the IP address of the manager. (See [Verify machines are running and get IP addresses](node-setup.md#verify-machines-are-running-and-get-ip-addresses)).


    ```
    docker@manager:~$ docker swarm init --advertise-addr 192.168.99.100
    Swarm initialized: current node (2tjrasfqfu945b7n4753374sw) is now a manager.

    To add a worker to this swarm, run the following command:

        docker swarm join \
        --token SWMTKN-1-144pfsupfo25h43zzr6b6bghjson8uedxjsndo5vuehqlyarsk-9k4q84axm008whv9zl4a8m8ct \
        192.168.99.100:2377

    To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
    ```

### Add a worker node to the swarm

1.  Log into the worker machine.

    ```
    $ docker-machine ssh worker
                            ##         .
                      ## ## ##        ==
                   ## ## ## ## ##    ===
               /"""""""""""""""""\___/ ===
          ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~ /  ===- ~~~
               \______ o           __/
                 \    \         __/
                  \____\_______/
     _                 _   ____     _            _
    | |__   ___   ___ | |_|___ \ __| | ___   ___| | _____ _ __
    | '_ \ / _ \ / _ \| __| __) / _` |/ _ \ / __| |/ / _ \ '__|
    | |_) | (_) | (_) | |_ / __/ (_| | (_) | (__|   <  __/ |
    |_.__/ \___/ \___/ \__|_____\__,_|\___/ \___|_|\_\___|_|

      WARNING: this is a build from test.docker.com, not a stable release.

    Boot2Docker version 1.13.0-rc7, build HEAD : b2cde29 - Sat Jan 14 00:29:39 UTC 2017
    Docker version 1.13.0-rc7, build 48a9e53
    ```

2.  On the worker, run the `join` command given as the output of the `swarm init` command you ran on the manager.

    ```
    docker@worker:~$ docker swarm join \
    >     --token SWMTKN-1-144pfsupfo25h43zzr6b6bghjson8uedxjsndo5vuehqlyarsk-9k4q84axm008whv9zl4a8m8ct \
    >     192.168.99.100:2377
    This node joined a swarm as a worker.
    ```

    If you don't have the command, run `docker swarm join-token worker` on a manager node to retrieve the `join` command for a worker for this swarm.

### List the nodes in the swarm

Log into the manager (e.g., `docker-machine ssh manager`) and run `docker node ls`.

```
  docker@manager:~$ docker node ls
ID                           HOSTNAME  STATUS  AVAILABILITY  MANAGER STATUS
2tjrasfqfu945b7n4753374sw *  manager   Ready   Active        Leader
syc46yimgtyz9ljcsfqiurvp0    worker    Ready   Active        
```

## What's next?

In the next step, we'll [deploy the voting app](deploy-app.md).
