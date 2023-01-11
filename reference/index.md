---
title: Reference documentation
description: This section includes the reference documentation for the Docker platform’s various APIs, CLIs, and file formats.
notoc: true
---

This section includes the reference documentation for the Docker platform's
various APIs, CLIs, drivers and specifications, and file formats.

## File formats

<div class="component-container">
    <!--start row-->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
             <div class="component-icon">
                 <a href="/engine/reference/builder/"><img src="/assets/images/build-frontends.svg" alt="Download and install" width="70" height="70"></a>
                 </div>
                 <h2 id="dockerfile"><a href="/engine/reference/builder/">Dockerfile</a></h2>
                <p> Defines the contents and startup behavior of a single container.</p>
            </div>
        </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/compose/compose-file/"><img src="/assets/images/build-multi-platform.svg" alt="Release notes" width="70" height="70"></a>
            </div>
                <h2 id="compose-file"><a href="/compose/compose-file/">Compose file</a></h2>
                <p>Defines a multi-container application.</p>
            </div>
        </div>
    </div>
</div>

## Command-line interfaces (CLIs)

<div class="component-container">
<!--start row-->
    <div class="row">
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/engine/reference/commandline/cli/"><img src="/assets/images/terminal.svg" alt="Docker CLI" width="70" height="70"></a>
            </div>
                <h2 id="docker-for-linux"><a href="/engine/reference/commandline/cli/">Docker CLI</a></h2>
                <p>The main CLI for Docker, includes all <code>docker</code> commands.</p>
        </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
          <div class="component-icon">
                 <a href="/compose/reference/"><img src="/assets/images/compose-cli.svg" alt="Compose CLI" width="70" height="70"></a>
          </div>
                <h2 id="docker-for-windows/install/"><a href="/compose/reference/">Compose CLI</a></h2>
                <p>The CLI for Docker Compose, which allows you to build and run multi-container applications.</p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/engine/reference/commandline/dockerd/"><img src="/assets/images/manage.svg" alt="Give feedback" width="70" height="70"></a>
            </div>
                <h2 id="docker-for-windows/install/"><a href="/engine/reference/commandline/dockerd/">Daemon CLI (dockerd)</a></h2>
                <p>Persistent process that manages containers.</p>
        </div>
     </div>
    </div>
</div>

## Application programming interfaces (APIs)

<div class="component-container">
    <!--start row-->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
             <div class="component-icon">
                 <a href="/engine/api/"><img src="/assets/images/engine-api.svg" alt="Engine API" width="70" height="70"></a>
             </div>
                 <h2 id="dockerfile"><a href="/engine/api/">Engine API</a></h2>
                <p> The main API for Docker, provides programmatic access to a daemon.</p>
        </div>
      </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/registry/spec/api/"><img src="/assets/images/storage.svg" alt="Registry API" width="70" height="70"></a>
            </div>
                <h2 id="compose-file"><a href="/registry/spec/api/">Registry API</a></h2>
                <p>Facilitates distribution of images to the engine.</p>
            </div>
        </div>
  </div>  
    <!--start row-->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
             <div class="component-icon">
                 <a href="/docker-hub/api/latest/"><img src="/assets/images/sync.svg" alt="Docker Hub API" width="70" height="70"></a>
                 </div>
                 <h2 id="dockerfile"><a href="/docker-hub/api/latest/">Docker Hub API</a></h2>
                <p> API to interact with Docker Hub.</p>
            </div>
        </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/docker-hub/api/dvp/"><img src="/assets/images/data.svg" alt="DVP Data API" width="70" height="70"></a>
            </div>
                <h2 id="compose-file"><a href="/docker-hub/api/dvp/">DVP Data API</a></h2>
                <p> API for Docker Verified Publishers to fetch analytics data. </p>
            </div>
        </div>
    </div>
</div>

## Drivers and specifications

<div class="component-container">
<!--start row-->
    <div class="row">
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/registry/spec/manifest-v2-2/"><img src="/assets/images/image.svg" alt="Image specification" width="70" height="70"></a>
            </div>
                <h2 id="docker-for-linux"><a href="/registry/spec/manifest-v2-2/">Image specification</a></h2>
                <p>Describes the various components of a Docker image.</p>
        </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
          <div class="component-icon">
                 <a href="/registry/spec/auth/"><img src="/assets/images/authentication.svg" alt="Registry token authentication" width="70" height="70"></a>
          </div>
                <h2 id="docker-for-windows/install/"><a href="/registry/spec/auth/">Registry token authentication</a></h2>
                <p>Outlines the Docker Registry authentication schemes.</p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <a href="/registry/storage-drivers/"><img src="/assets/images/engine-storage.svg" alt="Registry storage drivers" width="70" height="70"></a>
            </div>
                <h2 id="docker-for-windows/install/"><a href="/registry/storage-drivers/">Registry storage drivers</a></h2>
                <p>Enables support for given cloud providers when storing images with Registry.</p>
        </div>
     </div>
    </div>
</div>

