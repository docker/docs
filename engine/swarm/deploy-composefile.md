---
description: How to deploy a stack to a swarm
keywords:
- guide, swarm mode, composefile, stack, compose, deploy
title: Deploy a stack to a swarm
---

When your are running Docker Engine in swarm mode, you run `docker
stack deploy` to deploy a complete application stack in the swarm. The
`stack deploy` commands currently accepts compose file in version 3.0
and experimental DAB files as description to your stack. The built-in
swarm orchestrator and scheduler deploy your application stack to
nodes in your swarm to achieve and maintain the desired state.

This guide assumes you are working with the Docker Engine running in
swarm mode. You must run all `docker stack deploy` and `docker
service` commands from a manager node. As DAB files are still
experimental, this guide will only use compose file to deploy a stack.

The `docker deploy` command supports any compose file of version "3.0"
or above. Any compose file with a lower version will be refused as
*Unsupported*.

If you haven't already, read through [Swarm mode key concepts](key-concepts.md)
and [How services work](how-swarm-mode-works/services.md).

## Create our stack

The stack used in this guide is inspired by the [example voting
app](https://github.com/docker/example-voting-app). This stack consists of :

* a Python webapp which lets you vote between two options
* a Redis queue which collects new votes
* a .NET worker which consumes votes and stores them in…
* a Postgres database backed by a Docker volume
* a Node.js webapp which shows the results of the voting in real time

There is two networks : front-tier and back-tier. And only the vote
and result application are published.

The `docker-compose.yml` file looks is the following.

```yaml
version: "3.0"

services:
  vote:
    image: vote-app
    build: ./vote
    command: python app.py
    volumes:
     - ./vote:/app
    ports:
      - "5000:80"
    networks:
      - front-tier
      - back-tier

  result:
    image: vdemeester/result-app
    build: ./result
    command: nodemon --debug server.js
    volumes:
      - ./result:/app
    ports:
      - "5001:80"
      - "5858:5858"
    networks:
      - front-tier
      - back-tier

  worker:
    image: vdemeester/worker
    build: ./worker
    networks:
      - back-tier

  redis:
    image: redis:alpine
    ports: ["6379"]
    networks:
      - back-tier

  db:
    image: postgres:9.4
    container_name: db
    volumes:
      - "db-data:/var/lib/postgresql/data"
    networks:
      - back-tier

volumes:
  db-data:

networks:
  front-tier:
  back-tier:
```

### Validate your stack with `docker-compose`

Run the stack with `docker-compose` and make sure it works as
expected, i.e. you can vote and the result application is updated.

```bash
$ docker-compose up -d
docker-compose up -d                                                                                                     ~/src/github/docker/example-voting-app
Creating network "examplevotingapp_front-tier" with the default driver
Creating network "examplevotingapp_back-tier" with the default driver
Creating volume "examplevotingapp_db-data" with default driver
Building vote
Step 1 : FROM python:2.7-alpine
2.7-alpine: Pulling from library/python
3690ec4760f9: Pull complete
# […]
Successfully built 961a8478fb8e
WARNING: Image for service vote was built because it did not already exist. To rebuild this image you must use `docker-compose build` or `docker-compose up --build`.
Building worker
Step 1 : FROM microsoft/dotnet:1.0.0-preview2-sdk
1.0.0-preview2-sdk: Pulling from microsoft/dotnet
386a066cd84a: Already exists
# […]
Successfully built 8e3d6b2d9d7b
WARNING: Image for service worker was built because it did not already exist. To rebuild this image you must use `docker-compose build` or `docker-compose up --build`.
Pulling redis (redis:alpine)...
alpine: Pulling from library/redis
3690ec4760f9: Already exists
# […]
Digest: sha256:f1ed3708f538b537eb9c2a7dd50dc90a706f7debd7e1196c9264edeea521a86d
Status: Downloaded newer image for redis:alpine
Pulling db (postgres:9.4)...
9.4: Pulling from library/postgres
386a066cd84a: Already exists
# […]
Digest: sha256:9db811348585075eddb7b6938fa65c0236fe8e4f7feaf3c2890a3c4b7f4c9bfc
Status: Downloaded newer image for postgres:9.4
Building result
Step 1 : FROM node:5.11.0-slim
5.11.0-slim: Pulling from library/node
8b87079b7a06: Pull complete
# […]
Successfully built 7ee28b2ba725
WARNING: Image for service result was built because it did not already exist. To rebuild this image you must use `docker-compose build` or `docker-compose up --build`.
Creating examplevotingapp_vote_1
Creating examplevotingapp_worker_1
Creating redis
Creating examplevotingapp_result_1
Creating db

$ docker ps
CONTAINER ID        IMAGE                   COMMAND                  CREATED             STATUS              PORTS                                          NAMES
4332daa2d077        postgres:9.4            "/docker-entrypoint.s"   5 minutes ago       Up 5 minutes        5432/tcp                                       db
66d216122008        vdemeester/result-app   "nodemon --debug serv"   5 minutes ago       Up 5 minutes        0.0.0.0:5858->5858/tcp, 0.0.0.0:5001->80/tcp   examplevotingapp_result_1
c6401d390656        vdemeester/worker       "/bin/sh -c 'dotnet W"   5 minutes ago       Up 4 minutes                                                       examplevotingapp_worker_1
1fc0858853e7        redis:alpine            "docker-entrypoint.sh"   5 minutes ago       Up 5 minutes        0.0.0.0:32784->6379/tcp                        redis
e9a897645542        vote-app                "python app.py"          5 minutes ago       Up 5 minutes        0.0.0.0:5000->80/tcp                           examplevotingapp_vote_1
4a5083351397        docker-dev:master       "hack/dind bash"         4 hours ago         Up 4 hours                                                         nostalgic_meninsky
```

Our stack is running locally on our development enviroment.

### Push generated images

To be able to deploy it into our cluster, we need to make sure we push
our images to a docker registry so that it's available for all the
node of the swarm.

```bash
$ docker-compose push
Pushing vote (vdemeester/vote-app:latest)...
The push refers to a repository [docker.io/vdemeester/vote-app]
9d5d2ef2fa43: Pushed
5b1685972742: Pushed
1c0ffb0b324a: Pushed
3d414ea916c4: Pushed
bd7330a79bcf: Pushed
c9fc143a069a: Pushed
011b303988d2: Pushed
latest: digest: sha256:a34a439778ef33516572864d4f64bda18e3d4c04a4836e1ac70eb7e957724acf size: 1782
Pushing worker (vdemeester/worker:latest)...
The push refers to a repository [docker.io/vdemeester/worker]
c9e5cefbbd5b: Pushed
b84c856bd3c2: Pushed
67067d973315: Pushed
487eb54dfb73: Mounted from microsoft/dotnet
afad0eb67159: Mounted from microsoft/dotnet
8e3baf9e95c0: Mounted from microsoft/dotnet
9f17712cba0b: Mounted from microsoft/dotnet
223c0d04a137: Mounted from library/logstash
fe4c16cbf7a4: Mounted from library/debian
latest: digest: sha256:5c65afd7a77737a0542ac66f6d7a269ba30162959c7d4721af3a8ec9821880d0 size: 2213
Pushing result (vdemeester/result-app:latest)...
The push refers to a repository [docker.io/vdemeester/result-app]
ae6eb49768dc: Pushed
c0f894cc87b4: Pushed
ddede8a0bb7d: Pushed
660addd7320b: Pushed
8892cb213866: Pushed
517442a963f1: Pushed
35229e387818: Pushed
5f70bf18a086: Mounted from library/node
d42503f40afa: Mounted from library/node
ffe0b65c4aff: Mounted from library/node
3d5a262d6929: Mounted from library/node
6eb35183d3b8: Mounted from library/node
latest: digest: sha256:9936d18524b020f4a8a8f4696243240e54fa7c46ca015c9e0a5ddcfe32b5442a size: 3448
```

Our application stack is now ready to be deployed.

## Deploy a stack

We can now deploy the application stack on the swarm.

```bash
$ docker deploy --compose-file docker-compose.yml myapp
Ignoring unsupported options: build

Creating network myapp_front-tier
Creating network myapp_back-tier
Creating network myapp_default
Creating service myapp_db
Creating service myapp_vote
Creating service myapp_result
Creating service myapp_worker
Creating service myapp_redis
```

Let's list the service to make sure our stack is deployed and running.

```bash
$ docker service ls
ID            NAME           MODE        REPLICAS  IMAGE
6b6i54p2zuuf  myapp_result   replicated  1/1       vdemeester/result-app@sha256:9936d18524b020f4a8a8f4696243240e54fa7c46ca015c9e0a5ddcfe32b5442a
j8smsivm6pwe  myapp_worker   replicated  1/1       vdemeester/worker@sha256:5c65afd7a77737a0542ac66f6d7a269ba30162959c7d4721af3a8ec9821880d0
wf5wdh50egti  myapp_vote     replicated  1/1       vdemeester/vote-app@sha256:a34a439778ef33516572864d4f64bda18e3d4c04a4836e1ac70eb7e957724acf
wk6u16k4836u  myapp_db       replicated  1/1       postgres:9.4@sha256:9db811348585075eddb7b6938fa65c0236fe8e4f7feaf3c2890a3c4b7f4c9bfc
z8ocla9noxei  myapp_redis    replicated  1/1       redis:alpine@sha256:f1ed3708f538b537eb9c2a7dd50dc90a706f7debd7e1196c9264edeea521a86d
```

And we can now connect to any ip of our cluster, we should see the
voting app on port 5000 and the result app on port 5001.