---
description: Home page for Docker's documentation
keywords: Docker, documentation, manual, guide, reference, api, samples
landing: true
title: Docker Documentation
notoc: true
---

Docker packages your app with its dependencies, freeing you from worrying about your
system configuration, and making your app more portable.

<div class="row">
<div class="col-md-6 block">
{% capture basics %}
### Learn the basics of Docker

The basic tutorial introduces Docker concepts, tools, and commands. The examples show you how to build, push,
and pull Docker images, and run them as containers. This
tutorial stops short of teaching you how to deploy applications.
{% endcapture %}{{ basics | markdownify }}
{% capture basics %}[Start the basic tutorial](/engine/getstarted/){: class="button outline-btn"}{% endcapture %}{{ basics | markdownify }}
</div>

<div class="col-md-6 block">
{% capture apps %}
### Define and deploy applications

The define-and-deploy tutorial shows how to relate
containers to each other and define them as services in an application that is ready to deploy at scale in a
production environment. Highlights [Compose Version 3 new features](/engine/getstarted-voting-app/index.md#compose-version-3-features-and-compatibility) and swarm mode.
{% endcapture %}{{ apps | markdownify }}
{% capture apps %}[Start the application tutorial](/engine/getstarted-voting-app/){: class="button outline-btn"}{% endcapture %}{{ apps | markdownify }}
</div>
</div>


## Components

<div class="component-container">
    <div class="row">
        <div class="col-md-4">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-mac/"> <img src="../images/apple_48.svg" alt="Docker for Mac"> </a>
                </div>
                <h3 id="docker-for-mac"><a href="docker-for-mac/">Docker for Mac</a></h3>
                <p>A native application using the macOS sandbox security model which delivers all Docker tools to your Mac.</p>
            </div>
        </div>
        <div class="col-md-4">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-mac/"> <img src="../images/windows_48.svg" alt="Docker for Mac"> </a>
                </div>
                <h3 id="docker-for-windows"><a href="/#docker-for-windows">Docker for Windows</a></h3>
                <p>A native Windows application which delivers all Docker tools to your Windows computer.</p>
            </div>
        </div>
        <div class="col-md-4">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-mac/"> <img src="../images/linux_48.svg" alt="Docker for Mac"> </a>
                </div>
                <h3 id="docker-for-linux"><a href="/#docker-for-linux">Docker for Linux</a></h3>
                <p>Install Docker on a computer which already has a Linux distribution installed.</p>
            </div>
        </div>
        <!--components-full-width-->
        <div class="col-md-12">
            <!--editions-->
            <div class="component-full">
                <div class="component-full-icon">
                    <img src="../images/engine_48.svg" alt="Docker Editions">
                </div>
                <div class="component-full-copy">
                    <h3 id="docker-engine"><a href="engine/installation/">Docker Editions</a></h3>
                    <p>Get started with containers quickly with Docker Community edition (Docker CE)
                       or Docker Enterprise Edition (Docker EE).</p>
		            <p>Editions are available for desktops, servers, or cloud providers.</p>
                </div>
            </div>
            <!--cloud-->
            <div class="component-full">
                <div class="component-full-icon">
                    <img src="../images/cloud_48.svg" alt="Docker Cloud">
                </div>
                <div class="component-full-copy">
                    <h3 id="docker-cloud"><a href="/#">Docker Cloud</a></h3>
                    <p>A hosted service for building, testing, and deploying Docker images to your hosts.</p>
                </div>
            </div>
            <!--UCP-->
            <div class="component-full">
                <div class="component-full-icon">
                    <img src="../images/UCP_48.svg" alt="Docker Universal Control Plane">
                </div>
                <div class="component-full-copy">
                    <h3 id="docker-cloud"><a href="datacenter/ucp/1.1/overview/">Docker Universal Control Plane</a></h3>
                    <p>(UCP) Manage a cluster of on-premises Docker hosts as if they were a single machine.</p>
                </div>
            </div>
            <!--compose-->
            <div class="component-full">
                <div class="component-full-icon">
                    <img src="../images/compose_48.svg" alt="Docker Compose">
                </div>
                <div class="component-full-copy">
                    <h3 id="docker-cloud"><a href="compose/overview/">Docker Compose</a></h3>
                    <p>Define applications built using multiple containers.</p>
                </div>
            </div>
            <!--hub-->
            <!-- <div class="component-full">
                <div class="component-full-icon">
                    <img src="../images/hub_48.svg" alt="Docker Hub">
                </div>
                <div class="component-full-copy">
                    <h3 id="docker-cloud"><a href="docker-hub/overview/">Docker Hub</a></h3>
                    <p>A hosted registry service for managing and building images.</p>
                </div>
            </div> -->
            <!--dtr-->
            <div class="component-full">
                <div class="component-full-icon">
                    <img src="../images/dtr_48.svg" alt="Docker Trusted Registry">
                </div>
                <div class="component-full-copy">
                    <h3 id="docker-cloud"><a href="docker-trusted-registry/">Docker Trusted Registry</a></h3>
                    <p>(DTR) stores and signs your images.</p>
                </div>
            </div>
            <!--machine-->
            <div class="component-full">
                <div class="component-full-icon">
                    <img src="../images/machine_48.svg" alt="Docker Machine">
                </div>
                <div class="component-full-copy">
                    <h3 id="docker-cloud"><a href="machine/install-machine/">Docker Machine</a></h3>
                    <p>Automate container provisioning on your network or in the cloud. Available for Windows, macOS, or Linux.</p>
                </div>
            </div>
            <!-- end col-12-->
        </div>
        <!-- end component-container-->
    </div>
</div>
