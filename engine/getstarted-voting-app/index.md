---
description: Getting started with multi-container Docker apps
keywords: docker, container, multi-container, services, swarm mode, cluster, stack deploy, compose, voting app
title: Get started with multi-container apps and services in swarm mode
---

This tutorial is built around a web-based voting application that collects, tallies, and returns the results of votes (for cats and dogs, or other choices you specify). The voting app includes several services. We'll deploy the app as a _stack_ to introduce some new concepts surfaced in [Compose v.3](/compose/compose-file.md), and also use [swarm mode](/engine/swarm/index.md), which is built into Docker Engine.

## Got Docker?

If you haven't yet downloaded Docker or installed it, go to [Get Docker](https://www.docker.com/) and grab Docker for your platform. Once you have Docker installed, you can run `docker hello-world` or other commands described in the newcomer tutorial to [verify your installation](/engine/getstarted/step_one.md#step-3-verify-your-installation). If you are totally new to Docker, you might try the quick [newcomer tutorial](/engine/getstarted/index.md) first, then come back.

## What you'll learn and do

You'll learn how to:

* Use `docker machine` to create multiple virtual local hosts or dockerized cloud servers
* Use `docker` commands to set up and run a swarm with manager and worker nodes
* Deploy the `vote` app by feeding our example `docker-stack.yml` file to `docker stack deploy`
* Test the app by voting for cats and dogs, and viewing the results
* Use the `visualizer` to explore and understand the runtime app and services
* Update the `docker-stack.yml` and re-deploy the app using a different `vote` image to implement a vote on different choices


## Anatomy of the voting app

The voting app you are about to deploy is composed of several services:


| Service        | Description | Base Image  |
| ------------- |--------------| -----|
| `vote`      | Presents the voting interface via port `5000`. Viewable at `<manager-IP>:5000` | Based on a Python image, `dockersamples/examplevotingapp_vote` |
| `result`      | Displays the voting results via port 5001.  Viewable at `<manager-IP>:5001`     |  Based on a Node.js image, `dockersamples/examplevotingapp_result` |
| `visualizer` | A web app that shows a map of the deployment of the various services across the available nodes via port `8080`. Viewable at `<manager-IP>:8080`  |  Based on a .NET image, `dockersamples/examplevotingapp_worker` |
| `redis` | Collects raw voting data and stores it in a key/value queue     |  Based on a `redis` image, `redis:alpine` |
| `db` | A PostgreSQL service which provides permanent storage on a host volume     |  Based on a `postgres` image, `postgres:9.4` |
| `worker` | A background service that transfers votes from the queue to permanent storage     |  Based on a .NET image, `dockersamples/examplevotingapp_worker` |

Each service will run in its own container. Using swarm mode, we can also scale the application to deploy replicas of containerized services distributed across multiple nodes.

## docker-stack.yml
For this tutorial, you need only have Docker running and the copy of `docker-stack.yml` we provide here.

This file defines all the services, their base images, and provides . You will use this file to deploy the app.

```
version: "3"
services:

  redis:
    image: redis:alpine
    ports:
      - "6379"
    networks:
      - frontend
    deploy:
      replicas: 2
      update_config:
        parallelism: 2
        delay: 10s
      restart_policy:
        condition: on-failure
  db:
    image: postgres:9.4
    volumes:
      - db-data:/var/lib/postgresql/data
    networks:
      - backend
    deploy:
      placement:
        constraints: [node.role == manager]
  vote:
    image: dockersamples/examplevotingapp_vote:before
    ports:
      - 5000:80
    networks:
      - frontend
    depends_on:
      - redis
    deploy:
      replicas: 2
      update_config:
        parallelism: 2
      restart_policy:
        condition: on-failure
  result:
    image: dockersamples/examplevotingapp_result:before
    ports:
      - 5001:80
    networks:
      - backend
    depends_on:
      - db
    deploy:
      replicas: 2
      update_config:
        parallelism: 2
        delay: 10s
      restart_policy:
        condition: on-failure

  worker:
    image: dockersamples/examplevotingapp_worker
    networks:
      - frontend
      - backend
    deploy:
      mode: replicated
      replicas: 1
      labels: [APP=VOTING]
      restart_policy:
        condition: on-failure
        delay: 10s
        max_attempts: 3
        window: 120s

  visualizer:
    image: dockersamples/visualizer:stable
    ports:
      - "8080:8080"
    stop_grace_period: 1m30s
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock"

networks:
  frontend:
  backend:

volumes:
  db-data:
```

## Set up Dockerized hosts

1. Create `manager1` and `worker1` machines:

  ```
  $  docker-machine create --driver virtualbox manager1
  Running pre-create checks...
  Creating machine...
  (manager1) Copying /Users/victoriabialas/.docker/machine/cache/boot2docker.iso to /Users/victoriabialas/.docker/machine/machines/manager1/boot2docker.iso...
  (manager1) Creating VirtualBox VM...
  (manager1) Creating SSH key...
  (manager1) Starting the VM...
  (manager1) Check network to re-create if needed...
  (manager1) Waiting for an IP...
  Waiting for machine to be running, this may take a few minutes...
  Detecting operating system of created instance...
  Waiting for SSH to be available...
  Detecting the provisioner...
  Provisioning with boot2docker...
  Copying certs to the local machine directory...
  Copying certs to the remote machine...
  Setting Docker configuration on the remote daemon...
  Checking connection to Docker...
  Docker is up and running!
  To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env manager1
  ```

2. Get the IP addresses of the machines.

  ```
  $ docker-machine ls
  NAME       ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER        ERRORS
  manager1   *        virtualbox   Running   tcp://192.168.99.100:2376           v1.13.0-rc6   
  worker1    -        virtualbox   Running   tcp://192.168.99.101:2376           v1.13.0-rc6   
  ```

## Create a swarm

1. Log on to the manager.

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

  Boot2Docker version 1.13.0-rc6, build HEAD : 5ab2289 - Wed Jan 11 23:37:52 UTC 2017
  Docker version 1.13.0-rc6, build 2f2d055
  ```

2. Initialize a swarm.

  ```
  docker@manager:~$ docker swarm init --advertise-addr 192.168.99.100
  Swarm initialized: current node (ro5ak9ybe5qa62h7r81q29z0k) is now a manager.

  To add a worker to this swarm, run the following command:

      docker swarm join \
      --token SWMTKN-1-4bk4romozl9ap5xpt0it6gzdwabtezs399go3fyaw1hy8t1kam-0xnovvmchc4wfe7xmh85faiwe \
      192.168.99.100:2377

  To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
  ```

### Add a worker node to the swarm

1. Log into the worker machine.

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

  Boot2Docker version 1.13.0-rc6, build HEAD : 5ab2289 - Wed Jan 11 23:37:52 UTC 2017
  Docker version 1.13.0-rc6, build 2f2d055
  ```

2. On the worker, run the join command given as the output of the `swarm init` command you ran on the manager.

  ```
  docker@worker:~$ docker swarm join \
  >     --token SWMTKN-1-4bk4romozl9ap5xpt0it6gzdwabtezs399go3fyaw1hy8t1kam-0xnovvmchc4wfe7xmh85faiwe \
  >     192.168.99.100:2377
  This node joined a swarm as a worker.
  ```

If you don't have the command, run `docker swarm join-token worker` on a manager node to retrieve the join command for a worker for this swarm.


### List the nodes in the swarm

Log into the manager and run `docker node ls`.

```
docker@manager:~$ docker node ls
ID                           HOSTNAME  STATUS  AVAILABILITY  MANAGER STATUS
dlosx9b74cgu4lh2qy54ebfp8    worker    Ready   Active        
ro5ak9ybe5qa62h7r81q29z0k *  manager   Ready   Active        Leader
```

## Deploy the app

In these steps, you'll use the `docker-stack.yml` file to deploy the application to the swarm you just created. To deploy an app  deployed to a swarm You'll run the deploy command from the manager.

The `docker-stack.yml` file must be located on a manager for the swarm where you want to deploy the application stack.  

1. Get `docker-stack.yml` either from the lab or by copying it from the example given here.

2. Copy `docker-stack.yml` from your host machine onto the manager.

  ```
  $ docker-machine scp ~/sandbox/voting-app/docker-stack.yml manager:/home/docker/.
  docker-stack.yml                                                                      100% 1558     1.5KB/s   00:00    
  ```

3. Log onto the manager.

  ```
  $ docker-machine ssh manager1
  ```

4. Check to make sure the `.yml` file is there.

  ```
  docker@manager:~$ ls
  docker-stack.yml  log.log
  ```

  You can use `vi` or other text editor to inspect it.

5. Deploy the application stack based on the `.yml`.

  ```
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
  ```

6. Verify that the stack deployed as expected with `docker stack services <appName>`.

  ```
  docker@manager:~$ docker stack services vote
  ID            NAME             MODE        REPLICAS  IMAGE
  0y3q6lgc0drn  vote_result      replicated  2/2       dockersamples/examplevotingapp_result:before
  fvsaqvuec4yw  vote_redis       replicated  2/2       redis:alpine
  igev2xk5s3zo  vote_worker      replicated  1/1       dockersamples/examplevotingapp_worker:latest
  vpfjr9b0qc01  vote_visualizer  replicated  1/1       dockersamples/visualizer:stable
  wctxjnwl22k4  vote_vote        replicated  2/2       dockersamples/examplevotingapp_vote:before
  zp0zyvgaguox  vote_db          replicated  1/1       postgres:9.4
  ```

## Next topic

TBD

## Next topic

TBD

## Where to go next

The voting app is also available as a [lab on
GitHub](https://github.com/docker/labs/blob/master/beginner/chapters/votingapp.md)
along with the complete [source
code]((https://github.com/docker/example-voting-app).

The lab is a deeper dive, and includes a few more tasks, like cloning a GitHub
repository, manually changing source code, and rebuilding an image instead of
using the ready-baked images referenced here.


&nbsp;
