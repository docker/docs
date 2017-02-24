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
<div class="col-md-6">
{% capture basics %}
### Learn the basics of Docker

The basic tutorial introduces Docker concepts, tools, and commands. The examples show you how to build, push,
and pull Docker images, and run them as containers. This
tutorial stops short of teaching you how to deploy applications.
{% endcapture %}{{ basics | markdownify }}
{% capture basics %}[Start the basic tutorial](/engine/getstarted/){: class="button outline-btn"}{% endcapture %}{{ basics | markdownify }}
</div>

<div class="col-md-6">
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

<div class="block">
  <div class="component-container">
      <div class="row">
          <div class="col-md-4">
              <div class="component">
                  <div class="component-icon">
                      <a href="docker-for-mac/"> <img src="../images/icon-apple@2X.png" alt="Docker for Mac"> </a>
                  </div>
                  <h3 id="docker-for-mac"><a href="docker-for-mac/">Docker for Mac</a></h3>
                  <p>A native application using the macOS sandbox security model which delivers all Docker tools to your Mac.</p>
              </div>
          </div>
          <div class="col-md-4">
              <div class="component">
                  <div class="component-icon">
                      <a href="docker-for-mac/"> <img src="../images/icon-windows@2X.png" alt="Docker for Mac"> </a>
                  </div>
                  <h3 id="docker-for-windows"><a href="/#docker-for-windows">Docker for Windows</a></h3>
                  <p>A native Windows application which delivers all Docker tools to your Windows computer.</p>
              </div>
          </div>
          <div class="col-md-4">
              <div class="component">
                  <div class="component-icon">
                      <a href="docker-for-mac/"> <img src="../images/icon-linux@2X.png" alt="Docker for Mac"> </a>
                  </div>
                  <h3 id="docker-for-linux"><a href="/#docker-for-linux">Docker for Linux</a></h3>
                  <p>Install Docker on a computer which already has a Linux distribution installed.</p>
              </div>
          </div>
          <!--components-full-width-->
          <div class="col-md-12">
              <!--engine-->
              <div class="component-full">
                  <div class="component-full-icon">
                      <img src="../images/icon-engine@2X.png" alt="Docker Engine">
                  </div>
                  <div class="component-full-copy">
                      <h3 id="docker-engine"><a href="engine/installation/">Docker Engine</a></h3>
                      <p>Create Docker images and run Docker containers. As of v1.12.0, Engine includes <a href="/#">swarm mode</a> container orchestration features.</p>
                  </div>
              </div>
              <!--cloud-->
              <div class="component-full">
                  <div class="component-full-icon">
                      <img src="../images/icon-cloud@2X.png" alt="Docker Engine">
                  </div>
                  <div class="component-full-copy">
                      <h3 id="docker-cloud"><a href="/#">Docker Cloud</a></h3>
                      <p>A hosted service for building, testing, and deploying Docker images to your hosts.</p>
                  </div>
              </div>
              <!--UCP-->
              <div class="component-full">
                  <div class="component-full-icon">
                      <img src="../images/icon-ucp@2X.png" alt="Docker Universal Control Plane">
                  </div>
                  <div class="component-full-copy">
                      <h3 id="docker-cloud"><a href="datacenter/ucp/1.1/overview/">Docker Universal Control Plane</a></h3>
                      <p>(UCP) Manage a cluster of on-premises Docker hosts as if they were a single machine.</p>
                  </div>
              </div>
              <!--compose-->
              <div class="component-full">
                  <div class="component-full-icon">
                      <img src="../images/icon-compose@2X.png" alt="Docker Compose">
                  </div>
                  <div class="component-full-copy">
                      <h3 id="docker-cloud"><a href="compose/overview/">Docker Compose</a></h3>
                      <p>Define applications built using multiple containers.</p>
                  </div>
              </div>
              <!--hub-->
              <div class="component-full">
                  <div class="component-full-icon">
                      <img src="../images/icon-hub@2X.png" alt="Docker Hub">
                  </div>
                  <div class="component-full-copy">
                      <h3 id="docker-cloud"><a href="docker-hub/overview/">Docker Hub</a></h3>
                      <p>A hosted registry service for managing and building images.</p>
                  </div>
              </div>
              <!--dtr-->
              <div class="component-full">
                  <div class="component-full-icon">
                      <img src="../images/icon-registry@2X.png" alt="Docker Trusted Registry">
                  </div>
                  <div class="component-full-copy">
                      <h3 id="docker-cloud"><a href="docker-trusted-registry/">Docker Trusted Registry</a></h3>
                      <p>(DTR) stores and signs your images.</p>
                  </div>
              </div>
              <!--machine-->
              <div class="component-full">
                  <div class="component-full-icon">
                      <img src="../images/icon-machine@2X.png" alt="Docker Machine">
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
</div>
