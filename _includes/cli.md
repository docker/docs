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

**The client and daemon API must both be at least
{{ site.data[include.datafolder][include.datafile].min_api_version }}
to use this command.** Use the `docker version` command on the client to check
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
    <td>Stability</td>
    <td>Description</td>
  </tr>
</thead>
<tbody>
{% for option in site.data[include.datafolder][include.datafile].options %}
  <tr>
    <td markdown="span">`--{{ option.option }}{% if option.shorthand %} , -{{ option.shorthand }}{% endif %}`</td>
    <td markdown="span">{% if option.default_value and option.default_value != "[]" %}`{{ option.default_value }}`{% endif %}</td>
    <td markdown="span">
    {% if option.deprecated and option.experimental %}
      <span style="color: #ce4844">deprecated</span><span>,&nbsp;</span><span style="color: #aa6708">experimental</span>
    {% elsif option.deprecated %}
      <span style="color: #ce4844">deprecated</span>
    {% elsif option.experimental %}
      <span style="color: #aa6708">experimental</span>
    {% else %}
      <span>stable</span>
    {% endif %}
    </td>
    <td markdown="span">{{ option.description | replace: "|","&#124;" | strip }}</td>
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
