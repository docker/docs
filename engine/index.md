---
description: Engine
keywords: Engine
redirect_from:
  - /edge/
  - /engine/ce-ee-node-activate/
  - /engine/misc/
  - /linux/
  - /manuals/ # TODO remove this redirect after we've created a landing page for the product manuals section
title: Docker Engine overview
---

Docker Engine is an open source containerization technology for building and
containerizing your applications. Docker Engine acts as a client-server
application with:

- A server with a long-running daemon process
  [`dockerd`](/engine/reference/commandline/dockerd).
- APIs which specify interfaces that programs can use to talk to and instruct
  the Docker daemon.
- A command line interface (CLI) client
  [`docker`](/engine/reference/commandline/cli/).

The CLI uses [Docker APIs](api/index.md) to control or interact with the Docker
daemon through scripting or direct CLI commands. Many other Docker applications
use the underlying API and CLI. The daemon creates and manage Docker objects,
such as images, containers, networks, and volumes.

For more details, see
[Docker Architecture](../get-started/overview.md#docker-architecture).

<div class="component-container">
  <!--start row-->
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/engine/install/"><img src="/assets/images/download.svg" alt="Arrow pointing downwards" width="70px" height="70px"></a>
        </div>
        <h2><a href="/engine/install/">Install Docker Engine</a></h2>
        <p>Learn how to install the open source Docker Engine for your distribution.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/storage/"><img src="/assets/images/engine-storage.svg" alt="Data disks" width="70px" height="70px"></a>
        </div>
        <h2><a href="/storage/">Storage</a></h2>
        <p>Use persistent data with Docker containers.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/network/"><img src="/assets/images/engine-networking.svg" alt="Computers on a local area network" width="70px" height="70px"></a>
        </div>
        <h2><a href="/network/">Networking</a></h2>
        <p>Manage network connections between containers.</p>
      </div>
    </div>
  </div>
  <!--start row-->
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/config/containers/logging/"><img src="/assets/images/engine-logging.svg" alt="Document with a text outline" width="70px" height="70px"></a>
        </div>
        <h2><a href="/config/containers/logging/">Container logs</a></h2>
        <p>Learn how to view and read container logs.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/config/pruning/"><img src="/assets/images/engine-pruning.svg" alt="A pair of scissors" width="70px" height="70px"></a>
        </div>
        <h2><a href="/config/pruning/">Prune</a></h2>
        <p>Tidy up unused resources.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/config/daemon/"><img src="/assets/images/engine-configure-daemon.svg" alt="Settings cogwheel with stars" width="70px" height="70px"></a>
        </div>
        <h2><a href="/config/daemon/">Configure the daemon</a></h2>
        <p>Delve into the configuration options of the Docker daemon.</p>
      </div>
    </div>
  </div>
  <!--start row-->
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/engine/security/rootless/"><img src="/assets/images/engine-rootless.svg" alt="Checkered shield" width="70px" height="70px"></a>
        </div>
        <h2><a href="/engine/security/rootless/">Rootless mode</a></h2>
        <p>Run Docker without root privileges.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/engine/deprecated/"><img src="/assets/images/engine-deprecated.svg" alt="Alarm bell with an exclamation mark" width="70px" height="70px"></a>
        </div>
        <h2><a href="/engine/deprecated/">Deprecated features</a></h2>
        <p>Find out what features of Docker Engine you should stop using.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/engine/release-notes/"><img src="/assets/images/note-add.svg" alt="Document with an overlaying plus sign" width="70px" height="70px"></a>
        </div>
        <h2><a href="/engine/release-notes/">Release notes</a></h2>
        <p>Read the release notes for the latest version.</p>
      </div>
    </div>
  </div>
</div>

## Licensing

The Docker Engine is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/moby/moby/blob/master/LICENSE) for the full license
text.

