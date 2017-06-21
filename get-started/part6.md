---
title: "Get Started, Part 6: Deploy your app"
keywords: deploy, production, datacenter, cloud, aws, azure, provider, admin, enterprise
description: Deploy your app to production using Docker CE or EE.
---
{% include_relative nav.html selected="6" %}

## Prerequisites

- [Install Docker version 1.13 or higher](/engine/installation/).
- Get [Docker Compose](/compose/overview.md) as described in [Part 3 prerequisites](/get-started/part3.md#prerequisites).
- Get [Docker Machine](/machine/overview.md) as described in [Part 4 prerequisites](/get-started/part4.md#prerequisites).
- Read the orientation in [Part 1](index.md).
- Learn how to create containers in [Part 2](part2.md).

- Make sure you have published the `friendlyhello` image you created by
[pushing it to a registry](/get-started/part2.md#share-your-image). We'll
be using that shared image here.

- Be sure your image works as a deployed container. Run this command,
slotting in your info for `username`, `repo`, and `tag`: `docker run -p 80:80
username/repo:tag`, then visit `http://localhost/`.

- Have [the final version of `docker-compose.yml` from Part 5](/get-started/part5.md#persisting-data) handy.

## Introduction

You've been editing the same Compose file for this entire tutorial. Well, we
have good news. That Compose file works just as well in production as it does
on your machine. Here, we'll go through some options for running your
Dockerized application.

## Choose an option

{% capture cloud %}
If you're okay with using Docker Community Edition in
production, you can use Docker Cloud to help manage your app on popular service providers such as Amazon Web Services, DigitalOcean, and Microsoft Azure.

To set up and deploy:

- Connect Docker Cloud with your preferred provider, granting Docker Cloud permission
  to automatically provision and "Dockerize" VMs for you.
- Use Docker Cloud to create your computing resources and create your swarm.
- Deploy your app.

> **Note**: We will be linking into the Docker Cloud documentation here; be sure
  to come back to this page after completing each step.

### Connect Docker Cloud

You can run Docker Cloud in [standard
mode](/docker-cloud/infrastructure/index.md) or in [Swarm
mode](/docker-cloud/cloud-swarm/index.md).

If you are running Docker Cloud in standard mode, follow instructions below to
link your service provider to Docker Cloud.

* [Amazon Web Services setup guide](/docker-cloud/cloud-swarm/link-aws-swarm/){: onclick="ga('send', 'event', 'Get Started Referral', 'Cloud', 'AWS');"}
* [DigitalOcean setup guide](/docker-cloud/infrastructure/link-do.md){: onclick="ga('send', 'event', 'Get Started Referral', 'Cloud', 'DigitalOcean');"}
* [Microsoft Azure setup guide](/docker-cloud/infrastructure/link-azure.md){: onclick="ga('send', 'event', 'Get Started Referral', 'Cloud', 'Azure');"}
* [Packet setup guide](/docker-cloud/infrastructure/link-packet.md){: onclick="ga('send', 'event', 'Get Started Referral', 'Cloud', 'Packet');"}
* [SoftLayer setup guide](/docker-cloud/infrastructure/link-softlayer.md){: onclick="ga('send', 'event', 'Get Started Referral', 'Cloud', 'SoftLayer');"}
* [Use the Docker Cloud Agent to Bring your Own Host](/docker-cloud/infrastructure/byoh.md){: onclick="ga('send', 'event', 'Get Started Referral', 'Cloud', 'BYOH');"}

If you are running in Swarm mode (recommended for Amazon Web Services or
Microsoft Azure), then skip to the next section on how to [create your
swarm](#create-your-swarm).

### Create your swarm

Ready to create a swarm?

* If you're on Amazon Web Services (AWS) you
  can [automatically create a
  swarm on AWS](/docker-cloud/cloud-swarm/create-cloud-swarm-aws/){: onclick="ga('send', 'event', 'Get Started Referral AWS', 'Cloud', 'Create AWS Swarm');"}.

* If you are on Microsoft Azure, you can [automatically create a
swarm on Azure](/docker-cloud/cloud-swarm/create-cloud-swarm-azure/){: onclick="ga('send', 'event', 'Get Started Referral Azure', 'Cloud', 'Create Azure Swarm');"}.

* Otherwise, [create your nodes](/docker-cloud/getting-started/your_first_node/){: onclick="ga('send', 'event', 'Get Started Referral', 'Cloud', 'Create Nodes');"}
  in the Docker Cloud UI, and run the `docker swarm init` and `docker swarm join`
  commands you learned in [part 4](part4.md) over [SSH via Docker
  Cloud](/docker-cloud/infrastructure/ssh-into-a-node/). Finally, [enable Swarm
  Mode](/docker-cloud/cloud-swarm/using-swarm-mode/) by clicking the toggle at
  the top of the screen, and [register the
  swarm](/docker-cloud/cloud-swarm/register-swarms/) you just created.

> **Note**: If you are [Using the Docker Cloud Agent to Bring your Own Host](/docker-cloud/infrastructure/byoh.md){: onclick="ga('send', 'event', 'Get
Started Referral', 'Cloud', 'BYOH');"}, this provider does not support swarm
mode. You can [register your own existing
swarms](/docker-cloud/cloud-swarm/register-swarms/) with Docker Cloud.

### Deploy your app

[Connect to your swarm via Docker
Cloud](/docker-cloud/cloud-swarm/connect-to-swarm.md). On Docker for
Mac or Docker for Windows (Edge releases), you can [connect to your swarms
directly through the desktop app
menus](/docker-cloud/cloud-swarm/connect-to-swarm.md#use-docker-for-mac-or-windows-edge-to-connect-to-swarms).

Either way, this opens a terminal whose context is your local machine, but whose
Docker commands are routed up to the swarm running on your cloud service
provider. This is a little different from the paradigm you've been following,
where you were sending commands via SSH. Now, you can directly access both your
local file system and your remote swarm, enabling some very tidy-looking
commands:

```shell
docker stack deploy -c docker-compose.yml getstartedlab
```

That's it! Your app is running in production and is managed by Docker Cloud.
{% endcapture %}
{% capture enterpriseboilerplate %}
Customers of Docker Enterprise Edition run a stable, commercially-supported
version of Docker Engine, and as an add-on they get our first-class management
software, Docker Datacenter. You can manage every aspect of your application
via UI using Universal Control Plane, run a private image registry with Docker
Trusted Registry, integrate with your LDAP provider, sign production images with
Docker Content Trust, and many other features.

[Take a tour of Docker Enterprise Edition](https://www.docker.com/enterprise-edition){: class="button outline-btn" onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'Take tour');" style="margin-bottom: 30px; margin-right:58%"}
{% endcapture %}
{% capture enterprisedeployapp %}
Once you're all set up and Datacenter is running, you can [deploy your Compose
file from directly within the UI](/datacenter/ucp/2.1/guides/user/services/){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'Deploy app in UI');"}.

![Deploy an app on DDC](/datacenter/ucp/2.1/guides/images/deploy-app-ui-1.png)

After that, you'll see it running, and can change any aspect of the application
you choose, or even edit the Compose file itself.

![Managing app on DDC](/datacenter/ucp/2.1/guides/images/deployed_visualizer.png)
{% endcapture %}
{% capture enterprisecloud %}
{{ enterpriseboilerplate }}

The bad news is: the only cloud providers with official Docker
Enterprise editions are Amazon Web Services and Microsoft Azure.

The good news is: there are one-click templates to quickly deploy Docker
Enterprise on each of these providers:

* [Docker Enterprise for AWS](https://store.docker.com/editions/enterprise/docker-ee-aws?tab=description){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'EE for AWS');"}
* [Docker Enterprise for Azure](https://store.docker.com/editions/enterprise/docker-ee-azure?tab=description){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'EE for Azure');"}

> **Note**: Having trouble with these? View [our setup guide for AWS](/datacenter/install/aws/){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'AWS setup guide');"}.
> You can also [view the WIP guide for Microsoft Azure](https://github.com/docker/docker.github.io/pull/2796){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'Azure setup guide');"}.

{{ enterprisedeployapp }}
{% endcapture %}
{% capture enterpriseonprem %}
{{ enterpriseboilerplate }}

Bringing your own server to Docker Enterprise and setting up Docker Datacenter
essentially involves two steps:

1. [Get Docker Enterprise Edition for your server's OS from Docker Store](https://store.docker.com/search?offering=enterprise&type=edition){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'Get Docker EE for your OS');"}.
2. Follow the [instructions to install Datacenter on your own host](/datacenter/install/linux/){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'BYOH setup guide');"}.

> **Note**: Running Windows containers? View our [Windows Server setup guide](/docker-ee-for-windows/install/){: onclick="ga('send', 'event', 'Get Started Referral', 'Enterprise', 'Windows Server setup guide');"}.

{{ enterprisedeployapp }}
{% endcapture %}

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" href="#cloud">Docker CE (Cloud provider)</a></li>
  <li><a data-toggle="tab" href="#enterprisecloud">Enterprise (Cloud provider)</a></li>
  <li><a data-toggle="tab" href="#enterpriseonprem">Enterprise (On-premise)</a></li>
</ul>
<div class="tab-content">
  <div id="cloud" class="tab-pane fade in active" markdown="1">{{ cloud }}</div>
  <div id="enterprisecloud" class="tab-pane fade" markdown="1">{{ enterprisecloud }}</div>
  <div id="enterpriseonprem" class="tab-pane fade" markdown="1">{{ enterpriseonprem }}</div>
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
