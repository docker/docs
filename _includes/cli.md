{% capture tabChar %}	{% endcapture %}<!-- Make sure atom is using hard tabs -->
{% capture dockerBaseDesc %}The base command for the Docker CLI.{% endcapture %}
{% if include.datafolder and include.datafile %}

## Description

{% if include.datafile=="docker" %}<!-- docker.yaml is textless, so override -->
{{ dockerBaseDesc }}
{% else %}
{{ site.data[include.datafolder][include.datafile].short }}
{% endif %}

{% if site.data[include.datafolder][include.datafile].min_api_version %}

<span class="badge badge-info">API {{ site.data[include.datafolder][include.datafile].min_api_version }}+</span>&nbsp;
The client and daemon API must both be at least
{{ site.data[include.datafolder][include.datafile].min_api_version }}
to use this command. Use the `docker version` command on the client to check
your client and daemon API versions.

{% endif %}

{% if site.data[include.datafolder][include.datafile].deprecated %}

> This command is deprecated.
>
> It may be removed in a future Docker version.
{: .warning }

{% endif %}

{% if site.data[include.datafolder][include.datafile].experimental %}

> This command is experimental.
>
> It should not be used in production environments.
{: .important }

{% endif %}

{% if site.data[include.datafolder][include.datafile].usage %}

## Usage

```none
{{ site.data[include.datafolder][include.datafile].usage | replace: tabChar,"" | strip }}{% if site.data[include.datafolder][include.datafile].cname %} COMMAND{% endif %}
```

{% endif %}
{% if site.data[include.datafolder][include.datafile].options %}

## Options

<table>
<thead>
  <tr>
    <td>Name, shorthand</td>
    <td>Default</td>
    <td>Description</td>
  </tr>
</thead>
<tbody>
{% for option in site.data[include.datafolder][include.datafile].options %}

  {% capture min-api %}{% if option.min_api_version %}<span class="badge badge-info">API {{ option.min_api_version }}+</span>&nbsp;{% endif %}{%endcapture%}
  {% capture stability-string %}{% if option.deprecated and option.experimental %}<span class="badge badge-danger">deprecated</span>&nbsp;<span class="badge badge-warning">experimental</span>&nbsp;{% elsif option.deprecated %}<span class="badge badge-danger">deprecated</span>&nbsp;{% elsif option.experimental %}<span class="badge badge-warning">experimental</span>&nbsp;{% endif %}{% endcapture %}
  {% capture all-badges %}{% unless min-api == '' and stability-string == '' %}{{ min-api }}{{ stability-string }}<br />{% endunless %}{% endcapture %}
  {% assign defaults-to-skip = "[],map[],false,0,0s,default,'',\"\"" | split: ',' %}
  {% capture option-default %}{% if option.default_value %}{% unless defaults-to-skip contains option.default_value or defaults-to-skip == blank %}`{{ option.default_value }}`{% endunless %}{% endif %}{% endcapture %}
  <tr>
    <td markdown="span">`--{{ option.option }}{% if option.shorthand %} , -{{ option.shorthand }}{% endif %}`</td>
    <td markdown="span">{{ option-default }}</td>
    <td markdown="span">{{ all-badges | strip }}{{ option.description | strip }}</td>
  </tr>

{% endfor %} <!-- end for option -->
</tbody>
</table>

{% endif %} <!-- end if options -->

{% if site.data[include.datafolder][include.datafile].cname %}

## Child commands

<table>
<thead>
  <tr>
    <td>Command</td>
    <td>Description</td>
  </tr>
</thead>
<tbody>
{% for command in site.data[include.datafolder][include.datafile].cname %}
  {% capture dataFileName %}{{ command | strip | replace: " ","_" }}{% endcapture %}
  <tr>
    <td markdown="span">[{{ command }}]({{ dataFileName | replace: "docker_","" }}/)</td>
    <td markdown="span">{{ site.data[include.datafolder][dataFileName].short }}</td>
  </tr>
{% endfor %}
</tbody>
</table>
{% endif %}

{% if site.data[include.datafolder][include.datafile].pname %}
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
{% endif %}

{% unless site.data[include.datafolder][include.datafile].pname == "docker" or site.data[include.datafolder][include.datafile].pname == "dockerd" or include.datafile=="docker" %}

## Related commands

<table>
<thead>
  <tr>
    <td>Command</td>
    <td>Description</td>
  </tr>
</thead>
<tbody>
{% for command in site.data[include.datafolder][parentdatafile].cname %}
  {% capture dataFileName %}{{ command | strip | replace: " ","_" }}{% endcapture %}
  <tr>
    <td markdown="span">[{{ command }}]({{ dataFileName | replace: "docker_","" }}/)</td>
    <td markdown="span">{{ site.data[include.datafolder][dataFileName].short }}</td>
  </tr>
{% endfor %}
</tbody>
</table>

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
