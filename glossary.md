---
title: "Docker Glossary"
description: "Glossary of terms used around Docker"
keywords: "glossary, docker, terms, definitions"
notoc: true
noratings: true
---
<!--
To edit/add/remove glossary entries, visit the YAML file at:
https://github.com/moby/moby.github.io/blob/master/_data/glossary.yaml

To get a specific entry while writing a page in the docs, enter Liquid text
like so:
{{ site.data.glossary["aufs"] }}
-->
<span id="glossaryMatch" />
<span id="topicMatch" />

## Glossary terms

To see a definition for a term, and all topics in the documentation that have
been tagged with that term, click any entry below:

{% for entry in site.data.glossary %}
- [{{ entry[0] }}](/glossary/?term={{ entry[0]}})
{% endfor %}
