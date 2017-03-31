---
description: Home page for Docker's documentation
keywords: Docker, documentation, manual, guide, reference, api, samples
landing: true
title: Docker Documentation
notoc: true
---

Docker packages your app with its dependencies, freeing you from worrying about
your system configuration, and making your app more portable.

<div class="row">
<div class="col-sm-12 col-md-12 col-lg-6">
{% capture basics %}
### Learn Docker basics

Get started learning Docker concepts, tools, and commands. The examples show you
how to build, push, and pull Docker images, and run them as containers. This
tutorial stops short of teaching you how to deploy applications.
{% endcapture %}{{ basics | markdownify }}
{% capture basics %}[Start the basic tutorial](/engine/getstarted/){: class="button outline-btn"}{% endcapture %}{{ basics | markdownify }}
</div>

<div class="col-sm-12 col-md-12 col-lg-6 block">
{% capture apps %}
### Define and deploy apps in Swarm Mode

Learn how to relate containers to each other, define them as services, and
configure an application stack ready to deploy at scale in a production
environment. Highlights Compose Version 3 new features and swarm mode.
{% endcapture %}{{ apps | markdownify }}
{% capture apps %}[Start the application tutorial](/engine/getstarted-voting-app/){: class="button outline-btn"}{% endcapture %}{{ apps | markdownify }}
</div>
</div>

## Components

<div class="component-container">
    <!--start row-->
    <div class="row">
    <!--organic row 1-->
    <!--Docker for Mac-->
        <div class="col-sm-4 col-md-12 col-lg-3 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-mac/"> <img src="../images/apple_48.svg" alt="Docker for Mac"> </a>
                </div>
                <h3 id="docker-for-mac"><a href="docker-for-mac/">Docker for Mac</a></h3>
                <p>A native application using the macOS sandbox security model which delivers all Docker tools to your Mac.</p>
            </div>
        </div>
    <!--Docker for Windows-->
        <div class="col-sm-4 col-md-12 col-lg-3 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-windows/"> <img src="../images/windows_48.svg" alt="Docker for Windows"> </a>
                </div>
                <h3 id="docker-for-windows"><a href="docker-for-windows/">Docker for Windows</a></h3>
                <p>A native Windows application which delivers all Docker tools to your Windows computer.</p>
            </div>
        </div>
    <!--Docker for Linux-->
        <div class="col-sm-4 col-md-12 col-lg-3 block">
            <div class="component">
                <div class="component-icon">
                    <a href="engine/installation/linux/ubuntu/"> <img src="../images/linux_48.svg" alt="Docker for Linux"> </a>
                </div>
                <h3 id="docker-for-linux"><a href="engine/installation/linux/ubuntu/">Docker for Linux</a></h3>
                <p>Install Docker on a computer which already has a Linux distribution installed.</p>
            </div>
        </div>
    <!--organic row 2-->
    <!--editions-->
    <div class="col-sm-4 col-md-12 col-lg-3 block">
        <div class="component">
            <div class="component-icon">
                <a href="engine/installation/"> <img src="../images/apple_48.svg" alt="Docker Editions"> </a>
            </div>
            <h3 id="editions"><a href="engine/installation/">Docker Editions</a></h3>
            <p>Platform matrix and superset of installers for Docker for desktops, servers, or cloud providers.</p>
        </div>
    </div>
    <!--compose-->
        <div class="col-sm-4 col-md-12 col-lg-3 block">
            <div class="component">
                <div class="component-icon">
                    <a href="compose/overview/"> <img src="../images/compose_48.svg" alt="Docker Compose"> </a>
                </div>
                <h3 id="compose"><a href="compose/overview/">Docker Compose</a></h3>
                <p>Define application stacks built using multiple containers, services, and swarm configurations.</p>
        </div>
    </div>
    <!--machine-->
        <div class="col-sm-4 col-md-12 col-lg-3 block">
            <div class="component">
                <div class="component-icon">
                    <a href="machine/overview/"> <img src="../images/machine_48.svg" alt="Docker Trusted Registry"> </a>
                </div>
                <h3 id="machine"><a href="machine/overview/">Docker Machine</a></h3>
                <p>Automate container provisioning on your network or in the cloud. Available for Windows, macOS, or Linux.</p>
        </div>
    </div>
    <!--organic row 3-->
    <!--cloud-->
        <div class="col-sm-4 col-md-12 col-lg-3 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-cloud/"> <img src="../images/cloud_48.svg" alt="Docker Cloud"> </a>
                </div>
                <h3 id="docker-cloud"><a href="docker-cloud/">Docker Cloud</a></h3>
                <p>A hosted service for building, testing, and deploying Docker images to your hosts.</p>
            </div>
        </div>
    <!--UCP-->
        <div class="col-sm-4 col-md-12 col-lg-3 block">
            <div class="component">
                <div class="component-icon">
                    <a href="datacenter/ucp/2.1/guides/"> <img src="../images/UCP_48.svg" alt="Docker Universal Control Plane"> </a>
                </div>
                <h3 id="ucp"><a href="datacenter/ucp/2.1/guides/">Docker Universal Control Plane</a></h3>
                <p>(UCP) Manage a cluster of on-premise Docker hosts like a single machine with this enterprise product.</p>
            </div>
        </div>
    <!--DTR-->
        <div class="col-sm-4 col-md-12 col-lg-3 block">
            <div class="component">
                <div class="component-icon">
                    <a href="datacenter/dtr/2.2/guides/"> <img src="../images/dtr_48.svg" alt="Docker Trusted Registry"> </a>
                </div>
                <h3 id="dtr"><a href="datacenter/dtr/2.2/guides/">Docker Trusted Registry</a></h3>
                <p>(DTR) An enterprise image storage solution you can install behind a firewall to manage images and access.</p>
        </div>
    </div>
    <!-- end real row-->
    </div>
<!-- end component-container 2-->
</div>
