---
description: Docker Desktop overview
keywords: Desktop, Docker, GUI, run, docker, local, machine, dashboard
title: Docker Desktop
redirect_from:
- /desktop/opensource/
- /docker-for-mac/dashboard/
- /docker-for-mac/opensource/
- /docker-for-windows/dashboard/
- /docker-for-windows/opensource/
---

> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a paid
> subscription.

Docker Desktop is an easy-to-install application for your Mac, Linux, or Windows environment
that enables you to build and share containerized applications and microservices. 

It provides a simple interface that enables you to manage your containers, applications, and images directly from your machine without having to use the CLI to perform core actions.

<style>
.tab-content > .tab-pane {
  background-color: #fafafb;
  border: 1px solid #ddd;
  border-top: 0;
  padding: 10px;
  border-bottom-left-radius: 10px;
  border-bottom-right-radius: 10px;
  margin-bottom: 10px;
}
.night .tab-content > .tab-pane {
  background-color: #0e1c25;
  border: 1px solid #4f6071;
}
</style>
<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#includes">What's included in Docker Desktop?</a></li>
<li><a data-toggle="tab" data-target="#features">What are the key features of Docker Desktop?</a></li>
</ul>
<div class="tab-content">
<div id="includes" class="tab-pane fade in active" markdown="1">

- [Docker Engine](../engine/index.md)
- Docker CLI client
- [Docker Buildx](../build/index.md)
- [Docker Compose](../compose/index.md)
- [Docker Content Trust](../engine/security/trust/index.md)
- [Kubernetes](https://github.com/kubernetes/kubernetes/)
- [Credential Helper](https://github.com/docker/docker-credential-helpers/)

</div>
<div id="features" class="tab-pane fade" markdown="1">

* Ability to containerize and share any application on any cloud platform, in multiple languages and frameworks.
* Easy installation and setup of a complete Docker development environment.
* Includes the latest version of Kubernetes.
* On Windows, the ability to toggle between Linux and Windows Server environments to build applications.
* Fast and reliable performance with native Windows Hyper-V virtualization.
* Ability to work natively on Linux through WSL 2 on Windows machines.
* Volume mounting for code and data, including file change notifications and easy access to running containers on the localhost network.

</div>
</div>

Docker Desktop works with your choice of development tools and languages and
gives you access to a vast library of certified images and templates in
[Docker Hub](https://hub.docker.com/). This enables development teams to extend
their environment to rapidly auto-build, continuously integrate, and collaborate
using a secure repository.

<div class="component-container">
    <!--start row-->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
             <div class="component-icon">
                 <img src="/assets/images/download.svg" alt="Download and install" width="70" height="70">
                 </div>
                 <h2 id="docker-for-mac">Install Docker Desktop</h2>
                <p> <a href="/desktop/install/mac-install/">On Mac </a>, <a href="/desktop/install/windows-install/">Windows</a> or <a href="/desktop/install/linux-install/">Linux</a></p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/desktop/use-desktop/"><img src="/assets/images/explore.svg" alt="Docker Desktop" width="70" height="70"></a>
            </div>
                <h2 id="docker-for-mac"><a href="/desktop/use-desktop/">Explore Docker Desktop</a></h2>
                <p>Navigate Docker Desktop and learn about its key features.</p>
         </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/desktop/release-notes/"><img src="/assets/images/note-add.svg" alt="Release notes" width="70" height="70"></a>
            </div>
                <h2 id="docker-for-linux"><a href="/desktop/release-notes/">View the release notes</a></h2>
                <p>Find out about new features, improvements, and bug fixes.</p>
        </div>
    </div>
    </div>
    <!--start row-->
    <div class="row">
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/desktop/faqs/general/"><img src="/assets/images/help.svg" alt="FAQs" width="70" height="70"></a>
            </div>
                <h2 id="docker-for-linux"><a href="/desktop/faqs/general/">Browse common FAQs</a></h2>
                <p>Explore general FAQs or FAQs for specific platforms.</p>
        </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
          <div class="component-icon">
                 <a href="/desktop/kubernetes/"><img src="/assets/images/all-inbox.svg" alt="Additional resources" width="70" height="70"></a>
          </div>
                <h2 id="docker-for-windows/install/"><a href="/desktop/kubernetes/">Find additional resources</a></h2>
                <p>Find information on networking features, deploying on Kuberntes and more.</p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/desktop/feedback/"><img src="/assets/images/sms.svg" alt="Give feedback" width="70" height="70"></a>
            </div>
                <h2 id="docker-for-windows/install/"><a href="/desktop/feedback/">Give feedback</a></h2>
                <p>Provide feedback on Docker Desktop or Docker Desktop features.</p>
        </div>
     </div>
    </div>
</div>
