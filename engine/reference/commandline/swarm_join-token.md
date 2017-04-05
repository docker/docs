---
datafolder: engine-cli
datafile: docker_swarm_join-token
title: docker swarm join-token
redirect_from:
- /engine/reference/commandline/swarm_join_token/
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% if site.edge == true %}
  {% assign datafolder = "engine-cli-edge" %}
{% else %}
  {% assign datafolder = page.datafolder %}
{% endif %}

{% include cli.md datafolder=datafolder datafile=page.datafile %}
