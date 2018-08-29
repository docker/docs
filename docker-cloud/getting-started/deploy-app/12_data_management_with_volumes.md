---
description: Data management with Volumes
keywords: Python, data, management
redirect_from:
- /docker-cloud/getting-started/python/12_data_management_with_volumes/
title: Data management with volumes
---

In the previous step, we set up Redis but didn't provide it a way to store the
data it's caching. This means that if you redeployed the redis service, or if
the container crashed, the data would be lost. To save the data so it persists
beyond the life of a container, or share data from one container to another,
you need to define a volume.

## Data persistence

To persist, data in Docker Cloud must be stored in a volume. The volume
can be defined on the image (for example in the Dockerfile), or specified when
you create a new service in the Docker Cloud web UI. Learn more about volumes in
Docker Cloud [here](/docker-cloud/apps/volumes.md).

### Test for lack of persistence

If you `redeploy` the Redis service you created earlier, the counter resets.

Let's try that. First, redeploy the redis service to reset the counter.

```bash
$ docker-cloud service redeploy redis --not-reuse-volumes
```

Check the container status using the `container ps` command, and wait until the new container is running again. In the example below you can see the original container in the "Terminated" state, and the new container that is "Starting".

```none
$ docker-cloud container ps --service redis
NAME     UUID      STATUS        IMAGE                RUN COMMAND      EXIT CODE  DEPLOYED        PORTS
redis-1  5ddc0d66  ✘ Terminated  redis:staging        /run.sh                  0  15 minutes ago  6379/tcp
redis-1  3eff67a9  ⚙ Starting    redis:staging        /run.sh
```

Once the container is running, get the web endpoint using `container ps`, then try curling or visiting the web endpoint again

```none
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io:80
<h3>Hello Friendly Users!</h3><b>Hostname:</b> web-1<br/><b>Visits:</b> 1%
```

The Redis cache service redeployment caused the counter to reset.

### Enabling persistence

The specific Redis image (*redis*) in this tutorial supports data persistence.
This is not a common requirement for a Redis cache and it's not enabled by
default in most images. However to activate this in *our* image, you only need
to set two environment variables.

Run the following command to create and set these two environment variables.

```none
$ docker-cloud service set \
-e REDIS_APPENDONLY=yes \
-e REDIS_APPENDFSYNC=always \
redis --redeploy
```

This command defines two new environment variables in the **redis** service and
then redeploys the service so they take effect. You can learn more about our
open source `redis` image [here](https://github.com/docker-library/redis/){: target="_blank" class="_"}.

With these settings, Redis can create and store its data in a volume. The volume is in `/data`.

Visit the web endpoint a few more times to make sure that the cache is working
as expected. Then redeploy the Redis service to see if the counter resets, or if
it persists even after the container is terminated and re-created.

Curl the service to increment the counter:

```none
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io:80
<h3>Hello Python users!!</h3><b>Hostname:</b> web-1<br/><b>Visits:</b> 1%
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io:80
<h3>Hello Python users!!</h3><b>Hostname:</b> web-2<br/><b>Visits:</b> 2%
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io:80
<h3>Hello Python users!!</h3><b>Hostname:</b> web-3<br/><b>Visits:</b> 3%
```

Next, redeploy the service using the `service redeploy` command:

```none
$ docker-cloud service redeploy redis
```

Check the service status:

```none
$ docker-cloud container ps --service redis
NAME     UUID      STATUS        IMAGE                RUN COMMAND      EXIT CODE  DEPLOYED        PORTS
cache-1  8193cc1b  ✘ Terminated  redis:staging        /run.sh                  0  10 minutes ago  6379/tcp
cache-1  61f63d97  ▶ Running     redis:staging        /run.sh                     37 seconds ago  6379/tcp
```

Once the service is running again, curl the web page again to see what the counter value is.

```none
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io:80
<h3>Hello Python users!!</h3><b>Hostname:</b> web-3<br/><b>Visits:</b> 4%
```

Congratulations! You've set up data persistence in Docker Cloud!

## Sharing/reusing data volumes between services

A service's volume can be accessed by another service. To do this you use the `--volumes-from` flag when creating the new service.

You might use this functionality to share data between two services, or to back
up, restore, or migrate a volume to a local host or a cloud storage provider.

## Download volume data for backup

In this next step, you download the `/data` volume from Redis to your local host using SCP (secure copy).

First, run an SSH service that mounts the volumes of the redis you want to back up:

```bash
$ docker-cloud service run -n download -p 2222:22 -e AUTHORIZED_KEYS="$(cat ~/.ssh/id_rsa.pub)" --volumes-from redis tutum/ubuntu
```

Then run **scp** to download the data volume files in Redis:

```bash
$ scp -r -P 2222 root@downloader-1.$DOCKER_ID_USER.svc.dockerapp.io:/data .
```

You now have a backup copy of the Redis data on your local host machine!

## What's Next?

Congratulations! You've completed the tutorials! You can now push an image to
Docker Cloud, deploy an app to your Cloud nodes, set environment variables,
scale the service, view logs, set up a load balancer and a data back end, and
set up a volume to save the data.

There's lots more to learn about Docker Cloud, so check out [the rest of our documentation](/docker-cloud/), the [API and CLI Documentation](../../../apidocs/docker-cloud.md), and our [Knowledge Hub](https://success.docker.com/Cloud) and [Docker Cloud Forums](https://forums.docker.com/c/docker-cloud).

You might also want to delete or remove all of your hello world Stacks, Services, and Nodes running in Docker Cloud. To clean up when you're finished with the tutorial:

- Click **Stacks** in the left navigation, hover over the stack you created and click the selection box that appears, then click **Terminate**.
- Once the Stack has terminated, click **Services** in the left navigation, hover over each service you created, click the selection box that appears, then click **Terminate**.
- Click **Node Clusters** in the left navigation, hover over the node cluster you created, click the selection box that appears, then click **Terminate**.

Objects (Stacks, Services, Node Clusters, and Containers and nodes) still appear
in the list in Docker Cloud for about five minutes after they are terminated.

Happy Docking!
