---
description: Scale the service
keywords: scale, Python, service
redirect_from:
- /docker-cloud/getting-started/python/7_scale_the_service/
- /docker-cloud/getting-started/golang/7_scale_the_service/
title: Scale the service
notoc: true
---

Right now, your service is running on a single container. That's great for now.

You can check how many containers are running using the `docker-cloud container ps` command.

```none
$ docker-cloud container ps
NAME                   UUID      STATUS     IMAGE                                          RUN COMMAND          EXIT CODE  DEPLOYED     PORTS
web-1                  6c89f20e  ▶ Running  my-username/python-quickstart:latest           python app.py                   1 hour ago   web-1.my-username.cont.dockerapp.io:49162->80/tcp
```

A single container works just fine for now, but it could be a problem if that container becomes unresponsive. To avoid this, you can scale to more than one container. You do this with the `service scale` command:

```bash
$ docker-cloud service scale web 2
```

In this example, you can see we're scaling the service called `web` to `2` containers.

Run `service ps` again, and you should now see your service scaling:

```none
$ docker-cloud service ps
NAME                 UUID      STATUS     IMAGE                                          DEPLOYED
web                  68a6fb2c  ⚙ Scaling  my-username/python-quickstart:latest           1 hour ago
```

If you run `container ps` you should see multiple containers:

```none
$ docker-cloud container ps
NAME                   UUID      STATUS      IMAGE                                          RUN COMMAND          EXIT CODE  DEPLOYED     PORTS
web-1                  6c89f20e  ▶ Running   my-username/python-quickstart:latest           python app.py                   1 hour ago   web-1.my-username.cont.dockerapp.io:49162->80/tcp
web-2                  ab045c42  ⚙ Starting  my-username/python-quickstart:latest                                                        80/tcp
```

Containers aren't assigned a *PORT* until they are *running*, so you need to wait until the Service status goes from *Scaling* to *Running* to see what port is assigned to them.

```none
$ docker-cloud container ps
NAME                   UUID      STATUS     IMAGE                                          RUN COMMAND          EXIT CODE  DEPLOYED      PORTS
web-1                  6c89f20e  ▶ Running  my-username/python-quickstart:latest           python app.py                   1 hour ago    web-1.my-username.cont.dockerapp.io:49162->80/tcp
web-2                  ab045c42  ▶ Running  my-username/python-quickstart:latest           python app.py                   1 minute ago  web-2.my-username.cont.dockerapp.io:49156->80/tcp
```

Use either of the URLs from the `container ps` command to visit one of your service's containers, either using your browser or curl.

In the example output above, the URL `web-1.my-username.cont.dockerapp.io:49162` reaches the web app on the first container, and `web-2.my-username.cont.dockerapp.io:49156` reaches the web app on the second container.

If you use curl to visit the pages, you should see something like this:

```none
$ curl web-1.$DOCKER_ID_USER.cont.dockerapp.io:49166
Hello Python Users!</br>Hostname: web-1</br>Counter: Redis Cache not found, counter disabled.%
$ curl web-2.$DOCKER_ID_USER.cont.dockerapp.io:49156
Hello Python Users!</br>Hostname: web-2</br>Counter: Redis Cache not found, counter disabled.%
```

Congratulations! You now have *two* containers running in your **web** service.

## What's Next?

[View service logs](8_view_logs.md)
