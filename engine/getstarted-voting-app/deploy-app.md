---
description: Setup for voting app example
keywords: multi-container, services, swarm mode, cluster, voting app, docker-stack.yml, docker stack deploy
title: Deploy the application
---

In these steps, you'll use the `docker-stack.yml` file to deploy the application to the swarm you just created. You'll deploy from the `manager`.

## Copy `docker-stack.yml` to the manager

The `docker-stack.yml` file must be located on a manager for the swarm where you want to deploy the application stack.

1. Get `docker-stack.yml` either from the lab or by copying it from the example given here.

2. Copy `docker-stack.yml` from your host machine onto the manager.

        $ docker-machine scp ~/sandbox/voting-app/docker-stack.yml manager:/home/docker/.
        docker-stack.yml                                                                      100% 1558     1.5KB/s   00:00

3. Log onto the `manager` node.

        $ docker-machine ssh manager

4. Check to make sure the `.yml` file is there.

        docker@manager:~$ ls
        docker-stack.yml  log.log

  You can use `vi` or other text editor to inspect it.

## Deploy the app

1. Deploy the application stack based on the `.yml` using the command `docker stack deploy` as follows.

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

2. Verify that the stack deployed as expected with `docker stack services <appName>`.

        docker@manager:~$ docker stack services vote
        ID            NAME             MODE        REPLICAS  IMAGE
        0y3q6lgc0drn  vote_result      replicated  2/2       dockersamples/examplevotingapp_result:before
        fvsaqvuec4yw  vote_redis       replicated  2/2       redis:alpine
        igev2xk5s3zo  vote_worker      replicated  1/1       dockersamples/examplevotingapp_worker:latest
        vpfjr9b0qc01  vote_visualizer  replicated  1/1       dockersamples/visualizer:stable
        wctxjnwl22k4  vote_vote        replicated  2/2       dockersamples/examplevotingapp_vote:before
        zp0zyvgaguox  vote_db          replicated  1/1       postgres:9.4

## What's next?

In the next steps, we'll view components of the running app on web pages. We
will vote for cats and dogs, view the results, and monitor the manager and
worker nodes, containers and services on a visualizer.

TO BE CONTINUED - WORK IN PROGRESS
