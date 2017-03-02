---
description: Setup for voting app example
keywords: multi-container, services, swarm mode, cluster, voting app, docker-stack.yml, docker stack deploy
title: Deploy the application
---

In these steps, you'll use the `docker-stack.yml` file to
deploy the voting application to the swarm you just created.

## Copy docker-stack.yml to the manager

The `docker-stack.yml` file must be located on a manager for the swarm where you want to deploy the application stack.

1.  Get `docker-stack.yml` either from the [source code in the lab](https://github.com/docker/example-voting-app/blob/master/docker-stack.yml) or by copying it from the example given [here](index.md#docker-stackyml-deployment-configuration-file).

    If you prefer to download the file directly from our GitHub
    repository rather than copy it from the documentation, you can use a tool like `curl`. This command downloads the raw file to the current directory on your local host. You can copy-paste it into your shell if you have `curl`:

    ```
     curl -o docker-stack.yml https://raw.githubusercontent.com/docker/example-voting-app/master/docker-stack.yml
    ```

    >**Tips:**
    >
    *  To get the URL for the raw file on GitHub, either use the link in the example command above, or go to the file on GitHub [here](https://github.com/docker/example-voting-app/blob/master/docker-stack.yml), then click **Raw** in the upper right.
    *  You might already have `curl` installed. If not, you can [get curl here](https://curl.haxx.se/).

2.  Copy `docker-stack.yml` from your host machine onto the manager.

    ```
    $ docker-machine scp ~/sandbox/voting-app/docker-stack.yml manager:/home/docker/.
    docker-stack.yml                                                                      100% 1558     1.5KB/s   00:00
    ```

3.  Log into the manager node.

    ```
    $ docker-machine ssh manager
    ```

    The `ssh` login should put you in `/home/docker/` by default.

4.  Check to make sure that the `.yml` file is there, using `ls`.

    ```
    docker@manager:~$ ls /home/docker/
    docker-stack.yml
    ```

    You can use `vi` or `cat` to inspect the file.

## Deploy the app

We'll deploy the application from the manager.

1.  Deploy the application stack based on the `.yml` using the command
[`docker stack deploy`](/engine/reference/commandline/stack_deploy.md) as follows.

    ```
    docker stack deploy --compose-file docker-stack.yml vote
    ```

    * The `--compose-file` option specifies the path to our stack file. In this case, we assume it's in the current directory so we simply name the stack file: `docker-stack.yml`.

    * For the example, we name this app `vote`, but we could name it anything we want.

      Here is an example of the command and the output.

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

2.  Verify that the stack deployed as expected with `docker stack services <APP-NAME>`.

    ```
    docker@manager:~$ docker stack services vote
    ID            NAME             MODE        REPLICAS  IMAGE
    1zkatkq7sf8n  vote_result      replicated  1/1       dockersamples/examplevotingapp_result:after
    hphnxyt93h42  vote_redis       replicated  2/2       redis:alpine
    jd0wafumrcil  vote_vote        replicated  2/2       dockersamples/examplevotingapp_vote:after
    msief4cqme29  vote_visualizer  replicated  1/1       dockersamples/visualizer:stable
    qa6y8sfmtjoz  vote_db          replicated  1/1       postgres:9.4
    w04bh1vumnep  vote_worker      replicated  1/1       dockersamples/examplevotingapp_worker:latest
    ```

## What's next?

In the next steps, we'll view components of the running app
on web pages, and [take the app for a test drive](test-drive.md).
