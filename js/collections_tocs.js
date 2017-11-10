---
layout: null
---
var collectionsTOC = new Array()
collectionsTOC["library"] = [
  {% for page in site.samples %}
  {
  "path":{{ page.url | jsonify }},
  "title":{{ page.title | jsonify }}
  }{% unless forloop.last%},{% endunless %}
  {% endfor %}
]
