---
title: "Getting Started, Part 3: Stateful, Multi-container Applications"
---

# Getting Started, Part 3: Stateful, Multi-container Applications

In [Getting Started, Part 2: Creating and Building Your App](part2.md), we
wrote, built, ran, and shared our first Dockerized app, which all fit in a
single container.

In part 3, we will expand this application so that it is comprised of two
containers running simultaneously: one running the web app we have already
written, and another that stores data on the web app's behalf.

## Understanding services

In a world where every executable is running in a container, things are very
fluid and portable, which is exciting. There's just one problem: if you run
two containers at the same time, they don't know about each other. Each
container is isolated from the host environment, by design -- that's how Docker
enables environment-agnostic deployment.

We need something that defines some connective tissue between containers, so
that they run at the same time, and have the right ports open
so they can talk to each other. It's obvious why: having a front-end application
is all well and good, but it's going to need to store data at some point,
and that's going to happen via a different executable entirely.

In a distributed application, these different pieces of the app are called
"services." For example, if you imagine a video sharing site, there will
probably be a service for storing application data in a database, a service
for video transcoding in the background after a user uploads something, a
service for the front-end, and so on, and they all need to work in concert.

The easiest way to organize your containerized app into services using
is using Docker Compose. We're going to add a data storage service
to our simple Hello World app. Don't worry, it's shockingly easy.

## Your first `docker-compose.yml` File

A `docker-compose.yml` file is a YAML markup file that is hierarchical in
structure, and defines how multiple Docker images should work together when
they are running in containers.

We saw that the "Hello World" app we created looked for a running instance of
Redis, and if it failed, it produced an error message. All we need is a running
Redis instance, and that error message will be replaced with a visitor counter.

Well, just as we grabbed the base image of Python earlier, we can grab the
official image of Redis, and run that right alongside our app.

Save this `docker-compose.yml` file:

{% gist johndmulhausen/7b8e955ccc939d9cef83a015e06ed8e7 %}

Yes, that's all you need to specify, and Redis will be pulled and run. You could
make a `Dockerfile` that pulls in the base image of Redis and builds a custom
image that has all your preferences "baked in," but we're just going to point to
the base image here, and accept the default settings. (Redis documents these
defaults on [the page for the official Redis
image](https://store.docker.com/images/1f6ef28b-3e48-4da1-b838-5bd8710a2053)).

This `docker-compose.yml` file tells Docker to do the following:

- Pull and run [the image we uploaded to Docker Hub in step 2](/getting-started/part2/#/share-the-app) as a service called `web`
- Map port 4000 on the host to `web`'s port 80
- Link the `web` service to the service we named `redis`; this ensures that the
  dependency between `redis` and `web` is expressed, and these containers will
  run together in the same subnet.
- Our service named `redis` just runs the official Redis image, so go get it from Docker Hub.

## Run and scale up your first multi-container app

Run this command in the directory where you saved `docker-compose.yml`:

```shell
docker-compose up
```

This will pull all the necessary images and run them in concert. Now when you
visit `http://localhost:4000`, you'll see a number next to the visitor counter
instead of the error message. It really works -- just keep hitting refresh.

## Connecting to containers with port mapping

With a containerized instance of Redis running, you're probably wondering --
how do I break through the wall of isolation and manage my data? The answer is,
port mapping. [The page for the official Redis
image](https://store.docker.com/images/1f6ef28b-3e48-4da1-b838-5bd8710a2053)
states that the normal management ports are open in their image, so you would
be able to connect to it at `localhost:6379` if you add a `ports:` section to
`docker-compose.yml` under `redis` that maps `6379` to your host, just as port
`80` is mapped for `web`. Same with MySQL or any other data solution; once you
map your ports, you can use your fave UI tools like MySQL Workbench, Redis
Desktop Manager, etc, to connect to your Dockerized instance.

Redis port mapping isn't necessary in `docker-compose.yml` because the two
services (`web` and `redis`) are linked, ensuring they run on the same host (VM
or physical machine), in a private subnet that is automatically created by the
Docker runtime. Containers within
that subnet can already talk to each other; it's connecting from the outside
that necessitates port mapping.

## Cheat sheet and recap: Hosts, subnets, and Docker Compose

You learned that by creating a `docker-compose.yml` file, you can define the
entire stack for your application. This ensures that your services run
together in a private subnet that lets them connect to each
other, but only to the world as specifically dircted. This means that if you
want to connect your favorite data management software to your data storage
service, you'll have to ensure the container has the proper port exposed and
your host has that port mapped to the container in `docker-compose.yml`.

```shell
docker-compose up #Pull and run images specified in `docker-compose.yml` as services
docker-compose up -d #Same thing, but in background mode
docker-compose stop #Stop all running containers for this app
docker-compose rm -f #Remove all containers for this app
```

## Get ready to scale

Until now, I've been able to shield you from worrying too much about host
management. That's because installing Docker always sets up a default way
to run containers on that machine. Docker for Windows and Mac
comes with a virtual machine host running a lighweight operating system
we call Moby, which is just a very slimmed-down Linux. Docker for Linux
just works without a VM at all. And Docker for Windows can even run Microsoft
Windows containers using native Hyper-V support. When you've run `docker
run` and `docker-compose up` so far, Docker has used these solutions
to run your containers. That's because we want you to be able to install
Docker and get straight to the work of development and building images.

But when it comes to getting your app into production, we all know that
you're not going to run just one host machine that has Redis, Python, and
all your other sevices. That won't scale. You need to learn how to run not
just multiple containers on your local host, but multiple containers on
multiple hosts. And that's precisely what we're going to get into next.

[On to "Part 4: Running our App in Production" >>](part4.md){: class="button darkblue-btn"}
