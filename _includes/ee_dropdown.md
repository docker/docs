{% if page.path contains "ucp" %}
  {% capture currentDoc %}UCP version {{ ucp_version }}{% endcapture %}
  {% capture dropdown %}
  <div class="dropdown">
    <button class="btn btn-primary dropdown-toggle" type="button" data-toggle="dropdown">{{ ucp_version }}
    <span class="caret"></span></button>
    <ul class="dropdown-menu">
      {% for version in ucp_versions %}
        {% if version != ucp_version %}<li><a href="/datacenter/ucp/{{ version }}">{{ version }}{% if version.latest %} (latest){% endif %}</a></li>
      {% endfor %}
    </ul>
  </div>
  {% endcapture %}
{% endif %}
{% if page.path contains "dtr" %}
  {% capture currentDoc %}DTR version {{ dtr_version }}{% endcapture %}
  {% capture dropdown %}
  <div class="dropdown">
    <button class="btn btn-primary dropdown-toggle" type="button" data-toggle="dropdown">{{ dtr_version }}
    <span class="caret"></span></button>
    <ul class="dropdown-menu">
      {% for version in dtr_versions %}
        {% if version != dtr_version %}<li><a href="/datacenter/dtr/{{ version }}">{{ version }}{% if version.latest %} (latest){% endif %}</a></li>
      {% endfor %}
    </ul>
  </div>
  {% endcapture %}
{% endif %}

> These are the docs for {{ currentDoc }}
>
> To select a different version, use the selector below.
>
> {{ dropdown }}
