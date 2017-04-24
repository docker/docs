{% capture tabChar %}	{% endcapture %}<!-- Make sure atom is using hard tabs -->
{% capture dockerBaseDesc %}The base command for the Docker CLI.{% endcapture %}
{% if include.datafolder and include.datafile %}

## Description

{% if include.datafile=="docker" %}<!-- docker.yaml is textless, so override -->
{{ dockerBaseDesc }}
{% else %}
{{ site.data[include.datafolder][include.datafile].short }}
{% endif %}

{% if site.data[include.datafolder][include.datafile].usage %}

## Usage

```none
{{ site.data[include.datafolder][include.datafile].usage | replace: tabChar,"" | strip }}{% if site.data[include.datafolder][include.datafile].cname %} COMMAND{% endif %}
```

{% endif %}
{% if site.data[include.datafolder][include.datafile].options %}

## Options

| Name, shorthand | Default | Description |
| ---- | ------- | ----------- |{% for option in  site.data[include.datafolder][include.datafile].options %}
| `--{{ option.option }}{% if option.shorthand %}, -{{ option.shorthand }}{% endif %}` | {% if option.default_value and option.default_value != "[]" %}`{{ option.default_value }}`{% endif %} | {{ option.description | replace: "|","&#124;" | strip }} | {% endfor %}

{% endif %}

{% if site.data[include.datafolder][include.datafile].cname %}

## Child commands

| Command | Description |
| ------- | ----------- |{% for command in site.data[include.datafolder][include.datafile].cname %}{% capture dataFileName %}{{ command | strip | replace: " ","_" }}{% endcapture %}
| [{{ command }}]({{ dataFileName | replace: "docker_","" }}/) | {{ site.data[include.datafolder][dataFileName].short }} |{% endfor %}

{% endif %}

{% unless site.data[include.datafolder][include.datafile].pname == include.datafile %}

## Parent command

{% capture parentfile %}{{ site.data[include.datafolder][include.datafile].plink | replace: ".yaml", "" | replace: "docker_","" }}{% endcapture %}
{% capture parentdatafile %}{{ site.data[include.datafolder][include.datafile].plink | replace: ".yaml", "" }}{% endcapture %}

{% if site.data[include.datafolder][include.datafile].pname == "docker" %}
{% capture parentDesc %}{{ dockerBaseDesc }}{% endcapture %}
{% else %}
{% capture parentDesc %}{{ site.data[include.datafolder][parentdatafile].short }}{% endcapture %}
{% endif %}

| Command | Description |
| ------- | ----------- |
| [{{ site.data[include.datafolder][include.datafile].pname }}]({{ parentfile }}) | {{ parentDesc }}|

{% endunless %}

{% unless site.data[include.datafolder][include.datafile].pname == "docker" or site.data[include.datafolder][include.datafile].pname == "dockerd" %}

## Related commands

| Command | Description |
| ------- | ----------- |{% for command in site.data[include.datafolder][parentdatafile].cname %}{% capture dataFileName %}{{ command | strip | replace: " ","_" }}{% endcapture %}
| [{{ command }}]({{ dataFileName | replace: "docker_","" }}/) | {{ site.data[include.datafolder][dataFileName].short }} |{% endfor %}

{% endunless %}

{% unless site.data[include.datafolder][include.datafile].long == site.data[include.datafolder][include.datafile].short %}

## Extended description

{{ site.data[include.datafolder][include.datafile].long }}

{% endunless %}

{% if site.data[include.datafolder][include.datafile].examples %}

## Examples

{{ site.data[include.datafolder][include.datafile].examples }}

{% endif %}
{% else %}

The include.datafolder or include.datafile was not set.

{% endif %}
