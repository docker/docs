---
datafolder: engine-cli-edge
datafile: docker_config
title: docker config
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/cli
-->

{% if page.datafolder contains '-edge' %}
  {% include edge_only.md section="cliref" %}
{% endif %}
{% include cli.md datafolder=page.datafolder datafile=page.datafile %}

## More info

[Store configuration data using Docker Configs](/engine/swarm/configs.md)
