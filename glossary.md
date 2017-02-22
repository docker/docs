---
title: "Docker Glossary"
description: "Glossary of terms used around Docker"
keywords: "glossary, docker, terms, definitions"
skip-right-nav: true
---
{% for entry in site.data.glossary.glossary %}
## {{ entry.term }}

{{ entry.def }}

<span id="related-{{ forloop.index }}" style="display:none" class="relatedGlossary">{{ entry.term }}</span>
{% endfor %}
