{% for override in site.data.not_edited_here.overrides %}
{% if page.url contains override.path %}
{% assign overridden="true" %}
{% assign editsource=override.source %}
{% break %}
{% endif %}
{% endfor %}
{% if overridden != "true" %}{% capture editsource %}https://github.com/docker/docker.github.io/edit/master/{{ page.path }}{% endcapture %}{% endif %}
<span><a href="{{ editsource }}" class="button darkblue-btn nomunge" style="color:#FFFFFF; width:100%; margin: 0px;">Edit This Page</a></span>
