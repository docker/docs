---
title: "Docker Glossary"
description: "Glossary of terms used around Docker"
keywords: "glossary, docker, terms, definitions"
notoc: true
---
<!--
To edit/add/remove glossary entries, visit the YAML file at:
https://github.com/docker/docker.github.io/blob/master/_data/glossary.yaml

To get a specific entry while writing a page in the docs, enter Liquid text
like so:
{{ site.data.glossary["aufs"] }}
-->

{% for entry in site.data.glossary %}{% assign newEntry="true" %}

## {{ entry[0] }}

{{ entry[1] }}

{% endfor %}
