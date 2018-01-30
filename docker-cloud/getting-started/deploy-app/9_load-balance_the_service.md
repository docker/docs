---
description: Load-balance the service
keywords: load, balance, Python
redirect_from:
- /docker-cloud/getting-started/python/9_load-balance_the_service/
- /docker-cloud/getting-started/golang/9_load-balance_the_service/
title: Load-balance the service
notoc: true
---

To load-balance your application, you need to deploy a load-balancing service.
This service distributes incoming requests to all of the available containers in
the application.

In this example, you need a load balancer that forwards incoming requests to
both container #1 (web-1) and container #2 (web-2). For this tutorial, you'll
use [Docker Cloud's HAProxy image](https://github.com/docker/dockercloud-haproxy){: target="_blank" class="_"} to load balance, but you could also use other custom load balancers.

You can configure and run the `haproxy` load balancer service from the command line using a command like the example below. (If you are using the Go quickstart, edit the `link-service` value before running the command.)

```none
$ docker-cloud service run \
-p 80:80/tcp \
--role global \
--autorestart ALWAYS \
--link-service web:web \
--name lb \
dockercloud/haproxy
```

**-p 80:80/tcp** publishes port 80 of the container, and maps it to port 80 of the node.

**--role global** grants [API access](../../apps/api-roles.md) to this service. You can use this to query the Docker Cloud API from within the service.

**--autorestart ALWAYS** tells Docker Cloud to always [restart the containers](../../apps/autorestart.md) if they stop.

**--link-service web:web** links your load balancer service *haproxy* with the *web* service, and names the link *web*. (Learn more about Service Linking [here](../../apps/service-links.md).)

**--name lb** names the service *lb* (short for *load balancer*).

**dockercloud/haproxy** specifies the public image that we're using to make this service.

Run the `service ps` command to check if your service is already running.

```none
$ docker-cloud service ps
NAME                 UUID      STATUS     IMAGE                                          DEPLOYED
web                  68a6fb2c  ▶ Running  my-username/quickstart-python:latest           2 hours ago
lb                   e81f3815  ▶ Running  dockercloud/haproxy:latest                     11 minutes ago
```

Now let's check the container for this service. Run `docker-cloud container ps`.

```none
$ docker-cloud container ps
NAME                   UUID      STATUS     IMAGE                                          RUN COMMAND          EXIT CODE  DEPLOYED        PORTS
web-1                  6c89f20e  ▶ Running  my-username/quickstart-python:latest           python app.py                   2 hours ago     web-1.my-username.cont.dockerapp.io:49162->80/tcp
web-2                  ab045c42  ▶ Running  my-username/quickstart-python:latest           python app.py                   33 minutes ago  web-2.my-username.cont.dockerapp.io:49156->80/tcp
lb-1                   9793e58b  ▶ Running  dockercloud/haproxy:latest                           /run.sh                   14 minutes ago  443/tcp, lb-1.my-username.cont.dockerapp.io:80->80/tcp
```

You should notice an URL endpoint in the *PORT* column for haproxy-1. In the
example above, this is `lb-1.my-username.cont.dockerapp.io:80`. Open the `lb-1`
URL in your browser or curl from the CLI.

If you refresh or run curl multiple times, you should see requests distributed
between the two containers of the `web` service. You can see which container
responds to your request in the `Hostname` section of the response.

```none
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io
Hello Friendly Users!</br>Hostname: web-1</br>Counter: Redis Cache not found, counter disabled.%
$ curl lb-1.$DOCKER_ID_USER.cont.dockerapp.io
Hello Friendly Users!</br>Hostname: web-2</br>Counter: Redis Cache not found, counter disabled.%
```

You can learn more about *dockercloud/haproxy*, our free open source HAProxy image <a href="https://github.com/docker/dockercloud-haproxy" target="_blank">here</a>.

## What's Next?

[Provision a data backend for your service](10_provision_a_data_backend_for_your_service.md)
