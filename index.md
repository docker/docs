---
description: Home page for Docker's documentation
keywords:
- Docker, documentation, manual, guide, reference, api
title: Welcome to the Docker Documentation
layout: docs
---

# Welcome to the Docker Documentation

At its core, Docker provides a way to run almost any application securely
isolated in a container. The isolation and security allow you to run many
containers simultaneously on your host. The lightweight nature of containers,
which run without the extra load of a hypervisor, means you can get more out of
your hardware. Additionally, your application can always be packaged with its
dependencies and environment variables right in the build image, making testing
and deployment simpler than ever.

## Getting Started

Try [Getting Started withDocker Engine](engine/getstarted/index.md) to learn the
basics of building and running containers, [Getting Started with Docker Compose](compose/gettingstarted.md) to see
how easy it is to deploy a multi-container app (standing up a Python + Redis app in the process), or [our Swarm tutorial](engine/swarm/swarm-tutorial/index.md), which
shows how to run containers on many hosts (VMs or physical machines) at once,
as a cluster. Prefer a more linear, guided, visual tour? [Check out our self-paced training](https://training.docker.com/self-paced-training)!

### Typical Docker Platform Workflow

1. Get your code and its dependencies into Docker [containers](engine/getstarted/step_two.md):
   - [Write a Dockerfile](engine/getstarted/step_four.md) that specifies the execution
     environment and pulls in your code.
   - If your app depends on external applications (such as Redis, or
     MySQL), simply [find them on a registry such as Docker Hub](docker-hub/repos.md), and refer to them in
     [a Docker Compose file](compose/overview.md), along with a reference to your application, so they'll run
     simultaneously.
     - Software providers also distribute paid software via the [Docker Store](https://store.docker.com).
   - Build, then run your containers on a virtual host via [Docker Machine](machine/overview.md) as you develop.
2. Configure [networking](engine/tutorials/networkingcontainers.md) and
   [storage](engine/tutorials/dockervolumes.md) for your solution, if needed.
3. Upload builds to a registry ([ours](engine/tutorials/dockerrepos.md), [yours](docker-trusted-registry/index.md), or your cloud provider's), to collaborate with your team.
4. If you're gonna need to scale your solution across multiple hosts (VMs or physical machines), [plan
   for how you'll set up your Swarm cluster](engine/swarm/key-concepts.md) and [scale it to meet demand](engine/swarm/swarm-tutorial/index.md).
   - Note: Use [Universal Control Plane](ucp/overview.md) and you can manage your
     Swarm cluster using a friendly UI!
5. Finally, deploy to your preferred
   cloud provider (or, for redundancy, *multiple* cloud providers) with [Docker Cloud](docker-cloud/overview.md). Or, use [Docker Datacenter](https://www.docker.com/products/docker-datacenter), and deploy to your own on-premise hardware.

## Components

<section class="section projects_items_section GenericDev" style="margin-top:-150px; margin-bottom:-150px">
<ul class="items widthcol3 media">
<li>
	<div class="media_image">
		<a href="/docker-for-mac/"><img src="/images/icon-apple@2X.png" alt="Docker for Mac"></a>
	</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/docker-for-mac/">Docker for Mac</a></h3>
		<p>A native application using the OS X sandbox security model which delivers all Docker tools to your Mac.</p>
	</div>
	</div>
</li>
<li>
	<div class="media_image">
		<a href="/docker-for-windows/"><img src="/images/icon-windows@2X.png" alt="Docker for Windows"></a>
	</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/docker-for-windows/">Docker for Windows</a></h3>
		<p>A native Windows application which delivers all Docker tools to your Windows computer.</p>
	</div>
	</div>
</li>
<li>
	<div class="media_image">
		<a href="/engine/installation/linux/"><img src="/images/icon-linux@2X.png" alt="Docker for Linux"></a>
	</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/engine/installation/linux/">Docker for Linux</a></h3>
		<p>Install Docker on a computer which already has a Linux distribution installed.</p>
	</div>
	</div>
</li>
</ul>
<ul class="items widthcol media">
<li>
<div class="media_image">
	<a href="/engine/installation/"><img src="/images/icon-engine@2X.png" alt="Docker Engine"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/engine/installation/">Docker Engine</a></h3>
		<p>
    Create Docker images and run Docker containers.</p>
    <p>
		As of v1.12.0, Engine includes <a href="/engine/swarm/">swarm mode</a> container orchestration features.</p>
	</div>
	</div>
</li>
</ul>
<ul class="items widthcol2 media">
<li>
<div class="media_image">
	<a href="/docker-hub/overview/"><img src="/images/icon-hub@2X.png" alt="Docker Hub"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/docker-hub/overview/">Docker Hub</a></h3>
		<p>
    A hosted registry service for managing and building images.</p>
	</div>
	</div>
</li>
<li>
<div class="media_image">
	<a href="/docker-cloud/overview/"><img src="/images/icon-cloud@2X.png" alt="Docker Cloud"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/docker-cloud/overview/">Docker Cloud</a></h3>
		<p>
    A hosted service for building, testing, and deploying Docker images to your hosts.</p>
	</div>
	</div>
</li>
<li>
<div class="media_image">
	<a href="/docker-trusted-registry/"><img src="/images/icon-registry@2X.png" alt="Docker Trusted Registry"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/docker-trusted-registry/">Docker Trusted Registry</a></h3>
		<p>
    (DTR) stores and signs your images.</p>
	</div>
	</div>
</li>
<li>
<div class="media_image">
	<a href="/ucp/overview/"><img src="/images/icon-ucp@2X.png" alt="Docker Universal Control Plane"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/ucp/overview/">Docker Universal Control Plane</a></h3>
		<p>
    (UCP) Manage a cluster of on-premises Docker hosts as if they were a single machine.
    </p>
	</div>
	</div>
</li>
<li>
<div class="media_image">
	<a href="/machine/install-machine/"><img src="/images/icon-machine@2X.png" alt="Docker Machine"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/machine/install-machine/">Docker Machine</a></h3>
		<p>
    Automate container provisioning on your network or in
    the cloud. Available for Windows, Mac OS X, or Linux.</p>
	</div>
	</div>
</li>
<li>
<div class="media_image">
	<a href="/compose/overview/"><img src="/images/icon-compose@2X.png" alt="Docker Compose"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/compose/overview/">Docker Compose</a></h3>
		<p>
    Define applications built using multiple containers.</p>
	</div>
	</div>
</li>
</ul>

<ul class="media">
<li>
<div class="media_image">
	<a href="mailto:feedback@docker.com?subject=Docker%20Feedback"><img src="/images/chat.png" alt="chat icon"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="mailto:feedback@docker.com?subject=Docker%20Feedback">Feedback!</a></h3>
		<p>
    Questions? Suggestions? Spot a typo?! 😱<br/>
    Email us at <a href="mailto:feedback@docker.com?subject=Docker%20Feedback">feedback@docker.com</a>.
    </p>
	</div>
	</div>


</li>
</ul>
</section>
