{% for article in site.kb %}{% for tag in article.tags %}{% if tag.tag == include.tag %}- [{{ article.title }}]({{ article.url }})
{% endif %}{% endfor %}{% endfor %}
