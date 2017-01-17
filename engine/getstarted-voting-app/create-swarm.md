---
description: Setup for voting app example
keywords: multi-container, services, swarm mode, cluster, voting app
title: Create a swarm
---

Now, we'll transform our Docker machines into a swarm.

1. Log on to the manager.

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

        Boot2Docker version 1.13.0-rc6, build HEAD : 5ab2289 - Wed Jan 11 23:37:52 UTC 2017
        Docker version 1.13.0-rc6, build 2f2d055

2. Initialize a swarm.

        docker@manager:~$ docker swarm init --advertise-addr 192.168.99.100
        Swarm initialized: current node (ro5ak9ybe5qa62h7r81q29z0k) is now a manager.

  To add a worker to this swarm, run the following command:

        docker swarm join \
        --token SWMTKN-1-4bk4romozl9ap5xpt0it6gzdwabtezs399go3fyaw1hy8t1kam-0xnovvmchc4wfe7xmh85faiwe \
        192.168.99.100:2377

  To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.

### Add a worker node to the swarm

1. Log into the worker machine.

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

        Boot2Docker version 1.13.0-rc6, build HEAD : 5ab2289 - Wed Jan 11 23:37:52 UTC 2017
        Docker version 1.13.0-rc6, build 2f2d055

2. On the worker, run the join command given as the output of the `swarm init` command you ran on the manager.

        docker@worker:~$ docker swarm join \
        >     --token SWMTKN-1-4bk4romozl9ap5xpt0it6gzdwabtezs399go3fyaw1hy8t1kam-0xnovvmchc4wfe7xmh85faiwe \
        >     192.168.99.100:2377
        This node joined a swarm as a worker.

If you don't have the command, run `docker swarm join-token worker` on a manager node to retrieve the join command for a worker for this swarm.

### List the nodes in the swarm

Log into the manager and run `docker node ls`.

      docker@manager:~$ docker node ls
      ID                           HOSTNAME  STATUS  AVAILABILITY  MANAGER STATUS
      dlosx9b74cgu4lh2qy54ebfp8    worker    Ready   Active
      ro5ak9ybe5qa62h7r81q29z0k *  manager   Ready   Active        Leader

## Deploy the app

In these steps, you'll use the `docker-stack.yml` file to deploy the application
to the swarm you just created. To deploy an app  deployed to a swarm You'll run
the deploy command from the manager.

The `docker-stack.yml` file must be located on a manager for the swarm where you
want to deploy the application stack.

1. Get `docker-stack.yml` either from the lab or by copying it from the
example given here.

2. Copy `docker-stack.yml` from your host machine onto the manager.

        $ docker-machine scp ~/sandbox/voting-app/docker-stack.yml manager:/home/docker/.
        docker-stack.yml                                                                      100% 1558     1.5KB/s   00:00

3. Log onto the manager.

        $ docker-machine ssh manager

4. Check to make sure the `.yml` file is there.

        docker@manager:~$ ls
        docker-stack.yml  log.log

    You can use `vi` or other text editor to inspect it.

5. Deploy the application stack based on the `.yml`.

        docker@manager:~$ docker stack deploy --compose-file docker-stack.yml vote
        Creating network vote_default
        Creating network vote_backend
        Creating network vote_frontend
        Creating service vote_result
        Creating service vote_worker
        Creating service vote_visualizer
        Creating service vote_redis
        Creating service vote_db
        Creating service vote_vote

6. Verify that the stack deployed as expected with `docker stack services <appName>`.

        docker@manager:~$ docker stack services vote
        ID            NAME             MODE        REPLICAS  IMAGE
        0y3q6lgc0drn  vote_result      replicated  2/2       dockersamples/examplevotingapp_result:before
        fvsaqvuec4yw  vote_redis       replicated  2/2       redis:alpine
        igev2xk5s3zo  vote_worker      replicated  1/1       dockersamples/examplevotingapp_worker:latest
        vpfjr9b0qc01  vote_visualizer  replicated  1/1       dockersamples/visualizer:stable
        wctxjnwl22k4  vote_vote        replicated  2/2       dockersamples/examplevotingapp_vote:before
        zp0zyvgaguox  vote_db          replicated  1/1       postgres:9.4

## What's next?

In the next step, we'll [deploy the voting app](deploy-app.md).
