{% for article in site.kb %}
{% for tag in article.tags %}
{% if tag==include.tag %}- [{{ article.title }}]({{ article.url }})
{% endfor %}
{% endfor %}
