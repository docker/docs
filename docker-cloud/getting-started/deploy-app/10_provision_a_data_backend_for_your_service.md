---
description: Provision a data backend for the service
keywords: provision, Python, service
redirect_from:
- /docker-cloud/getting-started/python/10_provision_a_data_backend_for_your_service/
- /docker-cloud/getting-started/golang/10_provision_a_data_backend_for_your_service/
title: Provision a data backend for your service
---

Docker Cloud offers a large number of data stores in the *Jumpstart* library,
including Redis, MongoDB, PostgreSQL, and MySQL.

You may have noticed that your app has a visit counter that's been disabled up
until now. In this step you add a data backend for your service to use. In
this specific tutorial we use a Redis cache, but most concepts apply to any
data backend.

## Provision the service

The first step is to provision the data service itself. Run this command to
create and run the Redis service using the [redis](https://github.com/docker-library/redis/){: target="_blank" class="_"}
image:

```none
$ docker-cloud service run \
--env REDIS_PASS="password" \
--name redis \
redis
```
**--env REDIS_PASS="password"** defines an environment variable that sets the password for Redis. Because we are not publishing any ports for this service, only services **linked** to your *Redis service* can connect to it.

Use `docker-cloud service ps` to check if your new redis service is *running*. This might take a minute or two.

```none
$ docker-cloud service ps
NAME                 UUID      STATUS            IMAGE                                          DEPLOYED
redis                89806f93  ▶ Running         redis:latest                                   29 minutes ago
web                  bf644f91  ▶ Running         my-username/python-quickstart:latest           26 minutes ago
lb                   2f0d4b38  ▶ Running         dockercloud/haproxy:latest                     25 minutes ago
```

## Link the web service to the redis service

Next, we set up the link between the `redis` service and the `web` service.

```bash
$ docker-cloud service set --link redis:redis --redeploy web
```

In this command, we're creating a link from the `web` service (specified at the end of the command) to the `redis` service, and naming the link `redis`.

Next, visit or `curl` the load balanced web endpoint again. The web service now counts of the number of visits to the web service. This uses the Redis data backend, and is synchronized between all of the service's containers.

If you're using curl, you should see the counter incrementing like this:

```none
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io
Hello World</br>Hostname: web-1</br>Counter: 1%
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io
Hello World</br>Hostname: web-3</br>Counter: 2%
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io
Hello World</br>Hostname: web-2</br>Counter: 3%
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io
Hello World</br>Hostname: web-5</br>Counter: 4%
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io
Hello World</br>Hostname: web-3</br>Counter: 5%
```

## What's Next?

Next, we look at [Stackfiles for your service](11_service_stacks.md).
