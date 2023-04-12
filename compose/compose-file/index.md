---
description: Find the latest recommended version of the Docker Compose file format for defining multi-container applications.
keywords: docker compose file, docker compose yml, docker compose reference, docker compose cmd, docker compose user, docker compose image, yaml spec, docker compose syntax, yaml specification, docker compose specification
redirect_from:
- /compose/yaml/
- /compose/compose-file/compose-file-v1/
title: Overview
toc_max: 4
toc_min: 1
---
{% include compose-eol.md %}

The Compose file is a [YAML](https://yaml.org){: target="_blank" rel="noopener" class="_"} file defining services,
networks, and volumes for a Docker application. The latest and recommended
version of the Compose file format is defined by the [Compose
Specification](https://github.com/compose-spec/compose-spec/blob/master/spec.md){:
target="_blank" rel="noopener" class="_"}. The Compose spec merges the legacy
2.x and 3.x versions, aggregating properties across these formats and is implemented by **Compose 1.27.0+**.

Use the links below to easily navigate key sections of the Compose specification. 

<div class="component-container">
  <!--start row-->
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href= "/compose/compose-file/01-status/"><img src="/assets/images/engine-api.svg" alt="Arrow pointing downwards" width="70px" height="70px"></a>
        </div>
        <h2><a href= "/compose/compose-file/01-status/">Status of the Specification</a></h2>
        <p>Read about the status of the specification.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href= "/compose/compose-file/02-model/"><img src="/assets/images/storage.svg" alt="Data disks" width="70px" height="70px"></a>
        </div>
        <h2><a href= "/compose/compose-file/02-model/">The Compose application model</a></h2>
        <p>Learn about the Compose application model.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href= "/compose/compose-file/03-compose-file/"><img src="/assets/images/build-multi-platform.svg" alt="Computers on a local area network" width="70px" height="70px"></a>
        </div>
        <h2><a href= "/compose/compose-file/03-compose-file/">The Compose file</a></h2>
        <p>Understand the Compose file.</p>
      </div>
    </div>
  </div>
  <!--start row-->
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href= "/compose/compose-file/04-version-and-name/"><img src="/assets/images/engine-logging.svg" alt="Document with a text outline" width="70px" height="70px"></a>
        </div>
        <h2><a href= "/compose/compose-file/04-version-and-name/">Version and name top-level element</a></h2>
        <p>Understand version and name attributes for Compose.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href= "/compose/compose-file/05-services/"><img src="/assets/images/build-configure-buildkit.svg" alt="A pair of scissors" width="70px" height="70px"></a>
        </div>
        <h2><a href= "/compose/compose-file/05-services/">Services top-level element</a></h2>
        <p>Explore all services attributes for Compose.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href= "/compose/compose-file/06-networks/"><img src="/assets/images/engine-networking.svg" alt="Settings cogwheel with stars" width="70px" height="70px"></a>
        </div>
        <h2><a href= "/compose/compose-file/06-networks/">Networks top-level element</a></h2>
        <p>Find all networks attributes for Compose.</p>
      </div>
    </div>
  </div>
  <!--start row-->
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href= "/compose/compose-file/07-volumes/"><img src="/assets/images/engine-storage.svg" alt="Checkered shield" width="70px" height="70px"></a>
        </div>
        <h2><a href= "/compose/compose-file/07-volumes/">Volumes top-level element</a></h2>
        <p>Explore all volumes attributes for Compose.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href= "/compose/compose-file/08-configs/"><img src="/assets/images/engine-configure-daemon.svg" alt="Alarm bell with an exclamation mark" width="70px" height="70px"></a>
        </div>
        <h2><a href= "/compose/compose-file/08-configs/">Configs top-level element</a></h2>
        <p>Find out about configs in Compose.</p>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
      <div class="component">
        <div class="component-icon">
          <a href= "/compose/compose-file/09-secrets/"><img src="/assets/images/lock.svg" alt="Document with an overlaying plus sign" width="70px" height="70px"></a>
        </div>
        <h2><a href= "/compose/compose-file/09-secrets/">Secrets top-level element</a></h2>
        <p>Learn about secrets in Compose.</p>
      </div>
    </div>
  </div>
</div>