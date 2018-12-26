---
title: "Get Started, Part 6: Deploy your app"
keywords: deploy, production, datacenter, cloud, aws, azure, provider, admin, enterprise
description: Deploy your app to production using Docker CE or EE.
---
{% include_relative nav.html selected="6" %}

## Prerequisites

- [Install Docker](/install/index.md).
- Get [Docker Compose](/compose/overview.md) as described in [Part 3 prerequisites](/get-started/part3.md#prerequisites).
- Get [Docker Machine](/machine/overview.md) as described in [Part 4 prerequisites](/get-started/part4.md#prerequisites).
- Read the orientation in [Part 1](index.md).
- Learn how to create containers in [Part 2](part2.md).

- Make sure you have published the `friendlyhello` image you created by
[pushing it to a registry](/get-started/part2.md#share-your-image). We use that
shared image here.

- Be sure your image works as a deployed container. Run this command,
slotting in your info for `username`, `repo`, and `tag`: `docker run -p 80:80
username/repo:tag`, then visit `http://localhost/`.

- Have [the final version of `docker-compose.yml` from Part 5](/get-started/part5.md#persist-the-data) handy.

## Introduction

You've been editing the same Compose file for this entire tutorial. Well, we
have good news. That Compose file works just as well in production as it does
on your machine. In this section, we will go through some options for running your
Dockerized application.

## Choose an option

{% capture community %}

### Install Docker Engine --- Community

Find the [install instructions](/install/#supported-platforms) for Docker Engine --- Community on the platform of your choice.

### Create your swarm

Run `docker swarm init` to create a swarm on the node.

### Deploy your app

Run `docker stack deploy -c docker-compose.yml getstartedlab` to deploy
the app on the cloud hosted swarm.

```shell
docker stack deploy -c docker-compose.yml getstartedlab

Creating network getstartedlab_webnet
Creating service getstartedlab_web
Creating service getstartedlab_visualizer
Creating service getstartedlab_redis
```

Your app is now running on your cloud provider.

#### Run some swarm commands to verify the deployment

You can use the swarm command line, as you've done already, to browse and manage
the swarm. Here are some examples that should look familiar by now:

* Use `docker node ls` to list the nodes in your swarm.

```shell
[getstartedlab] ~ $ docker node ls
ID                            HOSTNAME                                      STATUS              AVAILABILITY        MANAGER STATUS
n2bsny0r2b8fey6013kwnom3m *   ip-172-31-20-217.us-west-1.compute.internal   Ready               Active              Leader
```

* Use `docker service ls` to list services.

```shell
[getstartedlab] ~/sandbox/getstart $ docker service ls
ID                  NAME                       MODE                REPLICAS            IMAGE                             PORTS
ioipby1vcxzm        getstartedlab_redis        replicated          0/1                 redis:latest                      *:6379->6379/tcp
u5cxv7ppv5o0        getstartedlab_visualizer   replicated          0/1                 dockersamples/visualizer:stable   *:8080->8080/tcp
vy7n2piyqrtr        getstartedlab_web          replicated          5/5                 sam/getstarted:part6    *:80->80/tcp
```

* Use `docker service ps <service>` to view tasks for a service.

```shell
[getstartedlab] ~/sandbox/getstart $ docker service ps vy7n2piyqrtr
ID                  NAME                  IMAGE                            NODE                                          DESIRED STATE       CURRENT STATE            ERROR               PORTS
qrcd4a9lvjel        getstartedlab_web.1   sam/getstarted:part6   ip-172-31-20-217.us-west-1.compute.internal   Running             Running 20 seconds ago                       
sknya8t4m51u        getstartedlab_web.2   sam/getstarted:part6   ip-172-31-20-217.us-west-1.compute.internal   Running             Running 17 seconds ago                       
ia730lfnrslg        getstartedlab_web.3   sam/getstarted:part6   ip-172-31-20-217.us-west-1.compute.internal   Running             Running 21 seconds ago                       
1edaa97h9u4k        getstartedlab_web.4   sam/getstarted:part6   ip-172-31-20-217.us-west-1.compute.internal   Running             Running 21 seconds ago                       
uh64ez6ahuew        getstartedlab_web.5   sam/getstarted:part6   ip-172-31-20-217.us-west-1.compute.internal   Running             Running 22 seconds ago        
```

#### Open ports to services on cloud provider machines

At this point, your app is deployed as a swarm on your cloud provider servers,
as evidenced by the `docker` commands you just ran. But, you still need to
open ports on your cloud servers in order to:

* if using many nodes, allow communication between the `redis` service and `web` service

* allow inbound traffic to the `web` service on any worker nodes so that
Hello World and Visualizer are accessible from a web browser.

* allow inbound SSH traffic on the server that is running the `manager` (this may be already set on your cloud provider)

{: id="table-of-ports"}

These are the ports you need to expose for each service:

| Service        | Type    | Protocol |  Port   |
| :---           | :---    | :---     | :---    |
| `web`          | HTTP    | TCP      |  80     |
| `visualizer`   | HTTP    | TCP      |  8080   |
| `redis`        | TCP     | TCP      |  6379   |

Methods for doing this vary depending on your cloud provider.

We use Amazon Web Services (AWS) as an example.

> What about the redis service to persist data?
>
> To get the `redis` service working, you need to `ssh` into
the cloud server where the `manager` is running, and make a `data/`
directory in `/home/docker/` before you run `docker stack deploy`.
Another option is to change the data path in the `docker-stack.yml` to
a pre-existing path on the `manager` server. This example does not
include this step, so the `redis` service is not up in the example output.

### Iteration and cleanup

From here you can do everything you learned about in previous parts of the
tutorial.

* Scale the app by changing the `docker-compose.yml` file and redeploy
on-the-fly with the `docker stack deploy` command.

* Change the app behavior by editing code, then rebuild, and push the new image.
(To do this, follow the same steps you took earlier to [build the
app](part2.md#build-the-app) and [publish the
image](part2.md#publish-the-image)).

* You can tear down the stack with `docker stack rm`. For example:

  ```
  docker stack rm getstartedlab
  ```

Unlike the scenario where you were running the swarm on local Docker machine
VMs, your swarm and any apps deployed on it continue to run on cloud
servers regardless of whether you shut down your local host.

{% endcapture %}
{% capture enterpriseboilerplate %}
Customers of Docker Enterprise Edition run a stable, commercially-supported
version of Docker Engine, and as an add-on they get our first-class management
software, Docker Datacenter. You can manage every aspect of your application
through the interface using Universal Control Plane, run a private image registry with Docker
Trusted Registry, integrate with your LDAP provider, sign production images with
Docker Content Trust, and many other features.

{% endcapture %}
{% capture enterprisedeployapp %}
Once you're all set up and Docker Enterprise is running, you can [deploy your Compose
file from directly within the UI](/ee/ucp/swarm/deploy-multi-service-app/){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'Deploy app in UI');"}.

![Deploy an app on Docker Enterprise](/ee/ucp/images/deploy-multi-service-app-2.png)

After that, you can see it running, and can change any aspect of the application
you choose, or even edit the Compose file itself.

![Managing app on Docker Enterprise](/ee/ucp/images/deploy-multi-service-app-4.png)
{% endcapture %}
{% capture enterprise %}
{{ enterpriseboilerplate }}

Bringing your own server to Docker Enterprise and setting up Docker Datacenter
essentially involves two steps:

1. [Get Docker Enterprise for your server's OS from Docker Hub](https://hub.docker.com/search?offering=enterprise&type=edition){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'Get Docker EE for your OS');"}.
2. Follow the [instructions to install Docker Enterprise on your own host](/datacenter/install/linux/){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'BYOH setup guide');"}.

> **Note**: Running Windows containers? View our [Windows Server setup guide](/install/windows/docker-ee.md){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'Windows Server setup guide');"}.

{{ enterprisedeployapp }}
{% endcapture %}

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" href="#enterprise">Docker Enterprise</a></li>
  <li><a data-toggle="tab" href="#community">Docker Engine - Community</a></li>
</ul>
<div class="tab-content">
  <div id="enterprise" class="tab-pane fade in active" markdown="1">{{ enterprise }}</div>
  <div id="community" class="tab-pane fade" markdown="1">{{ community }}</div>
</div>

## Congratulations!

You've taken a full-stack, dev-to-deploy tour of the entire Docker platform.

There is much more to the Docker platform than what was covered here, but you
have a good idea of the basics of containers, images, services, swarms, stacks,
scaling, load-balancing, volumes, and placement constraints.

Want to go deeper? Here are some resources we recommend:

- [Samples](/samples/): Our samples include multiple examples of popular software
  running in containers, and some good labs that teach best practices.
- [User Guide](/engine/userguide/): The user guide has several examples that
  explain networking and storage in greater depth than was covered here.
- [Admin Guide](/engine/admin/): Covers how to manage a Dockerized production
  environment.
- [Training](https://training.docker.com/): Official Docker courses that offer
  in-person instruction and virtual classroom environments.
- [Blog](https://blog.docker.com): Covers what's going on with Docker lately.
