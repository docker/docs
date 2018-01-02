---
description: Home page for Docker's documentation
keywords: Docker, documentation, manual, guide, reference, api, samples
landing: true
title: Docker Documentation
notoc: true
notags: true
---
{% assign page.title = site.name %}

<div class="row">
<div markdown="1" class="col-xs-12 col-sm-12 col-md-12 col-lg-6 block">

## Get started with Docker

Try our new multi-part walkthrough that covers writing your first app,
data storage, networking, and swarms, and ends with your app running on
production servers in the cloud. Total reading time is less than an hour.

[Get started with Docker](/get-started/){: class="button outline-btn"}

</div>
<div markdown="1" class="col-xs-12 col-sm-12 col-md-12 col-lg-6 block">

## Try Docker Enterprise Edition

Run your solution in production with Docker Enterprise Edition to get a
management dashboard, security scanning, LDAP integration, content signing,
multi-cloud support, and more. Click below to test-drive a running instance of
Docker EE without installing anything.

[Try Docker Enterprise Edition](https://dockertrial.com){: class="button outline-btn" onclick="ga('send', 'event', 'EE Trial Referral', 'Front Page', 'Click');"}

</div>
</div>

## Docker Editions

<div class="row">
<div markdown="1" class="col-xs-12 col-sm-12 col-md-12 col-lg-6 block">

### Docker Community Edition

Get started with Docker and experimenting with container-based apps. Docker CE
is available on many platforms, from desktop to cloud to server. Build and share
containers and automate the development pipeline from a single environment.
Choose the Edge channel to get access to the latest features, or the Stable
channel for more predictability.

[Learn more about Docker CE](/install/index.md#platform-support-matrix){: class="button outline-btn"}

</div>
<div markdown="1" class="col-xs-12 col-sm-12 col-md-12 col-lg-6 block">

### Docker Enterprise Edition

Designed for enterprise development and IT teams who build, ship, and run
business critical applications in production at scale. Integrated, certified,
and supported to provide enterprises with the most secure container platform in
the industry to modernize all applications. Docker EE Advanced comes with enterprise
[add-ons](#docker-ee-add-ons) like UCP and DTR.

[Learn more about Docker EE](/install/#platform-support-matrix){: class="button outline-btn"}

</div>
</div><!-- end row -->

## Run Docker anywhere

<div class="component-container">
    <!--start row-->
    <div class="row">
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-mac/"> <img src="../images/apple_48.svg" alt="Docker for Mac"> </a>
                </div>
                <h3 id="docker-for-mac"><a href="docker-for-mac/">Docker for Mac</a></h3>
                <p>A native application using the macOS sandbox security model which delivers all Docker tools to your Mac.</p>
            </div>
        </div>
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-windows/"> <img src="../images/windows_48.svg" alt="Docker for Windows"> </a>
                </div>
                <h3 id="docker-for-windows"><a href="docker-for-windows/">Docker for Windows</a></h3>
                <p>A native Windows application which delivers all Docker tools to your Windows computer.</p>
            </div>
        </div>
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="install/linux/ubuntu/"> <img src="../images/linux_48.svg" alt="Docker for Linux"> </a>
                </div>
                <h3 id="docker-for-linux"><a href="install/linux/ubuntu/">Docker for Linux</a></h3>
                <p>Install Docker on a computer which already has a Linux distribution installed.</p>
            </div>
        </div>
    </div>
</div>

<div class="component-container">
    <!--start row-->
    <div class="row">
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-cloud/"> <img src="../images/cloud_48.svg" alt="Docker Cloud"> </a>
                </div>
                <h3 id="docker-cloud"><a href="docker-cloud/">Docker Cloud</a></h3>
                <p>A hosted service for building, testing, and deploying Docker images to your hosts.</p>
            </div>
        </div>
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-aws/"> <img src="../images/cloud_48.svg" alt="Docker for AWS"> </a>
                </div>
                <h3 id="docker-cloud-providers"><a href="docker-for-aws/">Docker for AWS</a></h3>
                <p>Deploy your Docker apps on AWS.</p>
            </div>
        </div>
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="docker-for-azure/"> <img src="../images/cloud_48.svg" alt="Docker for Azure"> </a>
                </div>
                <h3 id="docker-cloud-providers"><a href="docker-for-azure/">Docker for Azure</a></h3>
                <p>Deploy your Docker apps on Azure.</p>
            </div>
        </div>
    </div>
</div>

## Components

<h5>Docker EE Add-ons</h5>

<div class="component-container">
    <!--start row-->
    <div class="row">
    <!--UCP-->
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="datacenter/ucp/{{ site.ucp_version }}/guides/"> <img src="../images/UCP_48.svg" alt="Universal Control Plane"> </a>
                </div>
                <h3 id="ucp"><a href="datacenter/ucp/{{ site.ucp_version }}/guides/">Universal Control Plane</a></h3>
                <p>(UCP) Manage a cluster of on-premise Docker hosts like a single machine with this enterprise product.</p>
            </div>
        </div>
    <!--DTR-->
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="datacenter/dtr/{{ site.dtr_version }}/guides/"> <img src="../images/dtr_48.svg" alt="Docker Trusted Registry"> </a>
                </div>
                <h3 id="dtr"><a href="datacenter/dtr/{{ site.dtr_version }}/guides/">Docker Trusted Registry</a></h3>
                <p>(DTR) An enterprise image storage solution you can install behind a firewall to manage images and access.</p>
            </div>
        </div>
    </div>
    <!-- end real row-->
</div>

<h5>Docker Tools</h5>

<div class="component-container">
    <!--start row-->
    <div class="row">
    <!--compose-->
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="compose/overview/"> <img src="../images/compose_48.svg" alt="Docker Compose"> </a>
                </div>
                <h3 id="compose"><a href="compose/overview/">Docker Compose</a></h3>
                <p>Define application stacks built using multiple containers, services, and swarm configurations.</p>
            </div>
        </div>
    <!--machine-->
        <div class="col-sm-12 col-md-12 col-lg-4 block">
            <div class="component">
                <div class="component-icon">
                    <a href="machine/overview/"> <img src="../images/machine_48.svg" alt="Docker Trusted Registry"> </a>
                </div>
                <h3 id="machine"><a href="machine/overview/">Docker Machine</a></h3>
                <p>Automate container provisioning on your network or in the cloud. Available for Windows, macOS, or Linux.</p>
        </div>
    </div>
</div>


<!-- end component-container 2-->
</div>
