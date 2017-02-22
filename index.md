---
description: Home page for Docker's documentation
keywords: Docker, documentation, manual, guide, reference, api
layout: docs
title: Docker Documentation
notoc: true
---

Docker packages your app with its dependencies, freeing you from worrying about
your system configuration, and making your app more portable.

<table>
<tr valign="top">
<td width="50%">
{% capture basics %}
### Learn the basics of Docker

The basic tutorial introduces Docker concepts, tools, and commands. The examples
show you how to build, push, and pull Docker images, and run them as containers.
This tutorial stops short of teaching you how to deploy applications.
{% endcapture %}{{ basics | markdownify }}
</td>
<td width="50%">

{% capture apps %}
### Define and deploy applications

The define-and-deploy tutorial shows how to relate containers to each other and
define them as services in an application that is ready to deploy at scale in a
production environment. Highlights Compose Version 3 new features and swarm
mode.
{% endcapture %}{{ apps | markdownify }}

</td></tr>

<tr valign="top">
<td width="50%">{% capture basics %}[Start the basic tutorial](/engine/getstarted/){: class="button secondary-btn"}{% endcapture %}{{ basics | markdownify }}
</td>
<td width="50%">{% capture apps %}[Start the application tutorial](/engine/getstarted-voting-app/){: class="button secondary-btn"}{% endcapture %}{{ apps | markdownify }}
</td>
</tr>
</table>

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
		<p>A native application using the macOS sandbox security model which delivers all Docker tools to your Mac.</p>
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
	<a href="/docker-hub/index.md"><img src="/images/icon-hub@2X.png" alt="Docker Hub"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/docker-hub/index.md">Docker Hub</a></h3>
		<p>
    A hosted registry service for managing and building images.</p>
	</div>
	</div>
</li>
<li>
<div class="media_image">
	<a href="/docker-cloud/index.md"><img src="/images/icon-cloud@2X.png" alt="Docker Cloud"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/docker-cloud/index.md">Docker Cloud</a></h3>
		<p>
    A hosted service for building, testing, and deploying Docker images to your hosts.</p>
	</div>
	</div>
</li>
<li>
<div class="media_image">
	<a href="/datacenter/dtr/2.2/guides/"><img src="/images/icon-registry@2X.png" alt="Docker Trusted Registry"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/datacenter/dtr/2.2/guides/">Docker Trusted Registry</a></h3>
		<p>
    (DTR) stores and signs your images.</p>
	</div>
	</div>
</li>
<li>
<div class="media_image">
	<a href="/datacenter/ucp/2.1/guides/index.md"><img src="/images/icon-ucp@2X.png" alt="Docker Universal Control Plane"></a>
</div>
	<div class="media_content">
	<div data-mh="mh_docker_projects">
	<h3><a href="/datacenter/ucp/2.1/guides/index.md">Docker Universal Control Plane</a></h3>
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
    the cloud. Available for Windows, macOS, or Linux.</p>
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
</section>
