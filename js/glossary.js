---
layout: null
---
var glossary = [
{% for entry in site.data.glossary %}
{
"term": {{ entry[0] | jsonify }},
"def": {{ entry[1] | markdownify | replace:'href="#','href="/glossary/?term=' | jsonify }}
}{% unless forloop.last %},{% endunless %}
{% endfor %}
]
