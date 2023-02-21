---
description: Find the latest recommended version of the Docker Compose file format for defining multi-container applications.
keywords: docker compose file, docker compose yml, docker compose reference, docker compose cmd, docker compose user, docker compose image, yaml spec, docker compose syntax, yaml specification, docker compose specification
redirect_from:
- /compose/yaml/
- /compose/compose-file/compose-file-v1/
title: Compose file specification
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

<div class="row">
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/status/">Status of the specification</a></div>
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/model/">The Compose application model</a></div>
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/compose-file/">The Compose file</a></div>
</div>

<div class="row">
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/version-and-name/">Version and name top-level element</a></div>
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/services/">Services top-level element</a></div>
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/networks/">Networks top-level element</a></div>
</div>

<div class="row">
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/volumes/">Volumes top-level element</a></div>
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/configs/">Configs top-level element</a></div>
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/secrets/">Secrets top-level element</a></div>
</div>

<div class="row">
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/fragments/">Fragments</a></div>
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/extension/">Extension</a></div>
  <div class="panel col-xs-12 col-md-4"><a href= "/compose/compose-file/interpolation/">Interpolation</a></div>
</div>