---
title: "Glossary"
description: "Glossary of terms used around Docker"
keywords: "glossary, docker, terms, definitions"
notoc: true
noratings: true
skip_read_time: true
redirect_from:
- /engine/reference/glossary/
- /reference/glossary/
---
<!--
To edit/add/remove glossary entries, visit the YAML file at:
https://github.com/docker/docker.github.io/blob/master/_data/glossary.yaml

To get a specific entry while writing a page in the docs, enter Liquid text
like so:
{{ site.data.glossary["aufs"] }}
-->
<span id="glossaryMatch" />
<span id="topicMatch" />

<table border="1">
  {% for entry in site.data.glossary %}
    <tr>
      <td>{{ entry[0] }}</td>
      <td>{{ entry[1] | markdownify }}</td>
    </tr>
  {% endfor %}
</table>
