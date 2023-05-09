---
title: Docker CLI
description: Overview of the Docker CLI, how it works, and how to use it
keywords: docker, cli, command line, client
---

The Docker CLI is the terminal-based user interface for interacting with the
Docker daemon, and for invoking Docker commands.

The Docker CLI and the Docker daemon are separate components. Using the Docker
CLI doesn't necessarily mean interacting with the Docker daemon. Some `docker`
commands don't interface with the daemon at all, or do so optionally.

For more details about the Docker architecture, see
[Docker overview](../get-started/overview.md#docker-architecture).

This section describes how to use the Docker CLI, and how you can configure it.
[CLI reference](../engine/reference/commandline/index.md).

<div class="component-container">
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/cli/config-file/">
           <img src="/assets/images/build-configure-buildkit.svg" alt="Hammer and screwdriver" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/cli/config-file/">Configure</a></h2>
        <p>
          Read about the Docker CLI (client) configuration file.
        </p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/cli/io/">
           <img src="/assets/images/input_component.svg" alt="Cables" width="70px" height="70px">
          </a>
        </div>
        <h2><a href="/cli/io/">I/O streams</a></h2>
        <p>
          Learn how the Docker CLI uses stdin, stdout, and stderr.
        </p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href="/cli/prune/"><img src="/assets/images/engine-pruning.svg" alt="A pair of scissors" width="70px" height="70px"></a>
        </div>
        <h2><a href="/cli/prune/">Prune</a></h2>
        <p>Tidy up unused resources.</p>
      </div>
    </div>
  </div>
</div>
