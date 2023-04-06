---
description: Docker Hub overview
keywords: Docker, docker, docker hub, hub, overview
title: Overview
---

Docker Hub is a service provided by Docker for finding and sharing container images.

It's the world’s largest repository of container images with an array of content sources including container community developers, open source projects and independent software vendors (ISV) building and distributing their code in containers.

Docker Hub is also where you can go to [change your Docker account settings and carry out administrative tasks](admin-overview.md).

<style>
.tab-content > .tab-pane {s
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
<li class="active"><a data-toggle="tab" data-target="#includes">What key features are included in Docker Hub?</a></li>
<li><a data-toggle="tab" data-target="#features">What administrative tasks can I perform in Docker Hub?</a></li>
</ul>
<div class="tab-content">
<div id="includes" class="tab-pane fade in active" markdown="1">

* [Repositories](../docker-hub/repos/index.md): Push and pull container images.
* [Docker Official Images](official_images.md): Pull and use high-quality
container images provided by Docker.
* [Docker Verified Publisher Images](publish/index.md): Pull and use high-
quality container images provided by external vendors.
* [Docker-Sponsored Open Source Images](dsos-program.md): Pull and use high-
quality container images from non-commercial open source projects.
* [Builds](builds/index.md): Automatically build container images from
GitHub and Bitbucket and push them to Docker Hub.
* [Webhooks](webhooks.md): Trigger actions after a successful push
  to a repository to integrate Docker Hub with other services.
*[Docker Hub CLI](https://github.com/docker/hub-tool#readme){: target="_blank" rel="noopener" class="_"} tool (currently experimental) and an API that allows you to interact with Docker Hub. Browse through the [Docker Hub API](/docker-hub/api/latest/){: target="_blank" rel="noopener" class="_"} documentation to explore the supported endpoints.

</div>
<div id="features" class="tab-pane fade" markdown="1">

* [Create and manage teams & organizations](orgs.md)
* [Create a company](creating-companies.md)
* [Enforce sign in](configure-sign-in.md)
* Set up [SSO](../single-sign-on/index.md) and [SCIM](scim.md)
* Use [Group mapping](group-mapping.md)
* [Carry out domain audits](domain-audit.md)
* [Use Image Access Management](image-access-management.md) to control developers' access to certain types of images
* [Enable Registry Access Management](../desktop/hardened-desktop/registry-access-management.md)

</div>
</div>

<div class="component-container">
    <!--start row-->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
             <div class="component-icon">
                 <a href="/docker-id/"><img src="/assets/images/laptop.svg" alt="Docker ID" width="70" height="70"></a>
             </div>
                 <h2 id="docker-id"><a href="/docker-id/">Create an account</a></h2>
                <p>Sign up and create a new Docker ID</p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/docker-hub/repos/"><img src="/assets/images/explore.svg" alt="Docker Compose" width="70" height="70"></a>
            </div>
                <h2 id="docker-repos"><a href="/docker-hub/repos/">Create a repository</a></h2>
                <p>Create a repository to share your images with your team, customers, or Docker community. </p>
         </div>
     </div>
       <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/docker-hub/"><img src="/assets/images/checklist.svg" alt="quickstart" width="70" height="70"></a>
            </div>
                <h2 id="docker-hub"><a href="/docker-hub/">Quickstart</a></h2>
                <p>Step-by-step instructions on getting started on Docker Hub.</p>
    </div>
    </div>
        <!--start row-->
    <div class="row">
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
          <a href="/docker-hub/builds/">
           <img src="/assets/images/build-configure-buildkit.svg" alt="secure" width="70px" height="70px">
          </a>
            </div>
                <h2 id="docker-hub"><a href="/docker-hub/builds/">Use Automated builds</a></h2>
                <p>Create and manage automated builds and autotesting.</p>
        </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
          <div class="component-icon">
         <a href="/docker-hub/official_images"><img src="/assets/images/build-multi-platform.svg" alt="Stacked windows" alt="Staircase" width="70px" height="70px"></a>
          </div>
                <h2 id="docker-hub"><a href="/docker-hub/official_images">Official images</a></h2>
                <p>A curated set of Docker repositories hosted on Docker Hub.</p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/docker-hub/"><img src="/assets/images/note-add.svg" alt="Release notes" width="70" height="70"></a>
            </div>
                <h2 id="docker-release-notes"><a href="/docker-hub/release-notes/">Release notes</a></h2>
                <p>Find out about new features, improvements, and bug fixes.</p>
        </div>
     </div>
    </div>
</div>



