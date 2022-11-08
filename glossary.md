---
title: "Glossary"
description: "Glossary of terms used around Docker"
keywords: "glossary, docker, terms, definitions"
notoc: true
skip_feedback: true
skip_read_time: true
redirect_from:
- /engine/reference/glossary/
- /reference/glossary/
---
<!--
To edit/add/remove glossary entries, visit the YAML file at:
https://github.com/docker/docs/blob/main/_data/glossary.yaml

To get a specific entry while writing a page in the docs, enter Liquid text
like so:
{{ site.data.glossary["aufs"] }}
-->
<table>
  <thead>
    <tr><th>Term</th><th>Definition</th></tr>
  </thead>
  <tbody>
  {%- for entry in site.data.glossary -%}
    {%- assign id = entry[0] | slugify -%}
    <tr>
      <td><a class="glossary" id="{{ id }}" href="#{{ id }}">{{ entry[0] }}</a></td>
      <td>{{ entry[1] | markdownify }}</td>
    </tr>
  {%- endfor -%}
  </tbody>
</table>
