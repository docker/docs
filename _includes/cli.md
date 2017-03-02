{% capture tabChar %}	{% endcapture %}<!-- Make sure atom is using hard tabs -->
{% capture dockerBaseDesc %}The base command for the Docker CLI.{% endcapture %}
{% if page.datafolder and page.datafile %}

## Description

{% if page.datafile=="docker" %}<!-- docker.yaml is textless, so override -->
{{ dockerBaseDesc }}
{% else %}
{{ site.data[page.datafolder][page.datafile].short }}
{% endif %}

{% if site.data[page.datafolder][page.datafile].usage %}

## Usage

```shell
{{ site.data[page.datafolder][page.datafile].usage | replace: tabChar,"" | strip }}{% if site.data[page.datafolder][page.datafile].cname %} COMMAND{% endif %}
```

{% endif %}
{% if site.data[page.datafolder][page.datafile].options %}

## Options

| Name, shorthand | Default | Description |
| ---- | ------- | ----------- |{% for option in  site.data[page.datafolder][page.datafile].options %}
| `--{{ option.option }}{% if option.shorthand %}, -{{ option.shorthand }}{% endif %}` | {% if option.default_value and option.default_value != "[]" %}`{{ option.default_value }}`{% endif %} | {{ option.description | replace: "|","&#124;" | strip }} | {% endfor %}

{% endif %}

{% if site.data[page.datafolder][page.datafile].cname %}

## Child commands

| Command | Description |
| ------- | ----------- |{% for command in site.data[page.datafolder][page.datafile].cname %}{% capture dataFileName %}{{ command | strip | replace: " ","_" }}{% endcapture %}
| [{{ command }}]({{ dataFileName | replace: "docker_","" }}/) | {{ site.data[page.datafolder][dataFileName].short }} |{% endfor %}

{% endif %}

{% if site.data[page.datafolder][page.datafile].pname and site.data[page.datafolder][page.datafile].pname != page.datafile %}

## Parent command

{% capture parentfile %}{{ site.data[page.datafolder][page.datafile].plink | replace: ".yaml", "" | replace: "docker_","" }}{% endcapture %}
{% capture parentdatafile %}{{ site.data[page.datafolder][page.datafile].plink | replace: ".yaml", "" }}{% endcapture %}

{% if site.data[page.datafolder][page.datafile].pname == "docker" %}
{% capture parentDesc %}{{ dockerBaseDesc }}{% endcapture %}
{% else %}
{% capture parentDesc %}{{ site.data[page.datafolder][parentdatafile].short }}{% endcapture %}
{% endif %}

| Command | Description |
| ------- | ----------- |
| [{{ site.data[page.datafolder][page.datafile].pname }}]({{ parentfile }}) | {{ parentDesc }}|

{% endif %}

{% if site.data[page.datafolder][page.datafile].pname != "docker" %}

## Related commands

| Command | Description |
| ------- | ----------- |{% for command in site.data[page.datafolder][parentdatafile].cname %}{% capture dataFileName %}{{ command | strip | replace: " ","_" }}{% endcapture %}
| [{{ command }}]({{ dataFileName | replace: "docker_","" }}/) | {{ site.data[page.datafolder][dataFileName].short }} |{% endfor %}

{% endif %}

{% if site.data[page.datafolder][page.datafile].long != site.data[page.datafolder][page.datafile].short %}

## Extended description

{{ site.data[page.datafolder][page.datafile].long }}

{% endif %}

{% if site.data[page.datafolder][page.datafile].examples %}

## Examples

{{ site.data[page.datafolder][page.datafile].examples }}

{% endif %}

{% endif %}
