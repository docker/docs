{% capture tabChar %}	{% endcapture %}<!-- Make sure atom is using hard tabs -->
{% assign controller_data = site.data[include.datafolder][include.datafile] %}
{% assign parentPath = page.path | prepend: "/" | remove: page.name %}

<br />

{{ controller_data.short | replace_relative_links: page.path }}

{% if controller_data.min_api_version %}
{% comment %}
  To reduce unnecessary clutter on the page, we only show the minimum required
  API version if it requires a relatively recent version of the Engine, which
  is configured through the "min_api_threshold" option in _config.yml

  Below, we first convert the min_api_version from a string to a number, so that
  we can compare versions (see https://stackoverflow.com/a/27200972/1811501),
  then compare it, to decide whether to show the "minimum required API version".
{% endcomment %}
{% assign min_api_version = controller_data.min_api_version | plus: 0 %}
{% if min_api_version >= site.min_api_threshold %}

<a href="/engine/api/v{{ controller_data.min_api_version }}/" target="_blank" rel="noopener" class="_"><span class="badge badge-info" data-toggle="tooltip" data-placement="right" title="Open the {{ controller_data.min_api_version }} API reference (in a new window)">API {{ controller_data.min_api_version }}+</span></a>&nbsp;
The client and daemon API must both be at least
<a href="/engine/api/v{{ controller_data.min_api_version }}/" target="_blank" rel="noopener" class="_">{{ controller_data.min_api_version }}</a>
to use this command. Use the `docker version` command on the client to check
your client and daemon API versions.

{% endif %}
{% endif %}

{% if controller_data.deprecated %}

> This command is [deprecated](/engine/deprecated/){: target="_blank" rel="noopener" class="_"}.
>
> It may be removed in a future Docker version. For more information, see the [Docker Roadmap](https://github.com/docker/roadmap/issues/209){: target="_blank" rel="noopener" class="_"}.
{: .warning }

{% endif %}

{% if controller_data.experimental %}

> This command is experimental.
>
> This command is experimental on the Docker daemon. It should not be used in
> production environments.
> To enable experimental features on the Docker daemon, edit the
> [daemon.json](/engine/reference/commandline/dockerd/#daemon-configuration-file)
> and set `experimental` to `true`.
>
> {% include experimental.md %}

{% endif %}

{% if controller_data.experimentalcli %}

> This command is experimental.
>
> [Experimental features](https://docs.docker.com/engine/reference/commandline/cli/#experimental-features)
> are intended for testing and feedback as their functionality or UX may change
> between releases without warning or can be removed entirely in a future release.
{: .important }

{% endif %}

{% capture command-orchestrator %}
{% if controller_data.swarm %}

<span class="badge badge-info" data-toggle="tooltip" data-placement="right" title="This command works with the Swarm orchestrator.">Swarm</span> This command works with the Swarm orchestrator.

{% endif %}
{% if controller_data.kubernetes %}

<span class="badge badge-info" data-toggle="tooltip" data-placement="right" title="This command works with the Kubernetes orchestrator.">Kubernetes</span> This command works with the Kubernetes orchestrator.

{% endif %}
{% endcapture %}{{ command-orchestrator }}


{% if controller_data.usage %}

## Usage

```console
$ {{ controller_data.usage | replace: tabChar, "" | strip }}{% if controller_data.cname %} COMMAND{% endif %}
```

{% endif %}
{% unless controller_data.long == controller_data.short %}
{% if controller_data.options %}
Refer to the [options section](#options) for an overview of available [`OPTIONS`](#options) for this command.
{% endif %}

## Description

{: name="extended-description"}

{{ controller_data.long | replace_relative_links: page.path }}

{% endunless %}

{% if controller_data.examples %}
For example uses of this command, refer to the [examples section](#examples) below.
{% endif %}

{% if controller_data.options %}
  {% if controller_data.inherited_options %}
    {% assign alloptions = controller_data.options | concat:controller_data.inherited_options %}
  {% else %}
    {% assign alloptions = controller_data.options %}
  {% endif %}
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
{% for option in alloptions %}
  {% capture deprecated-badge %}{% if option.deprecated %}<a href="/engine/deprecated/" target="_blank" rel="noopener" class="_"><span class="badge badge-danger" data-toggle="tooltip" title="Read the deprecation reference (in a new window).">deprecated</span></a>{% endif %}{% endcapture %}
  {% capture experimental-daemon-badge %}{% if option.experimental %}<a href="/engine/reference/commandline/dockerd/#daemon-configuration-file" target="_blank" rel="noopener" class="_"><span class="badge badge-warning" data-toggle="tooltip" title="Read about experimental daemon options (in a new window).">experimental (daemon)</span></a>{% endif %}{% endcapture %}
  {% capture experimental-cli-badge %}{% if option.experimentalcli %}<a href="/engine/reference/commandline/cli/#configuration-files" target="_blank" rel="noopener" class="_"><span class="badge badge-warning"  data-toggle="tooltip" title="Read about experimental CLI options (in a new window).">experimental (CLI)</span></a>{% endif %}{% endcapture %}
  {%- if option.min_api_version -%}
    {% assign min_api_version = option.min_api_version | plus: 0 %}
    {% if min_api_version >= site.min_api_threshold %}
      {% capture min-api %}<a href="/engine/api/v{{ option.min_api_version }}/" target="_blank" rel="noopener" class="_"><span class="badge badge-info" data-toggle="tooltip" title="Open the {{ controller_data.min_api_version }} API reference (in a new window)">API {{ option.min_api_version }}+</span></a>{%endcapture%}
    {%- endif -%}
  {%- else -%}
    {% capture min-api %}{%endcapture%}
  {%- endif -%}
  {% capture flag-orchestrator %}{% if option.swarm %}<span class="badge badge-info" data-toggle="tooltip" title="This option works for the Swarm orchestrator.">Swarm</span>{% endif %}{% if option.kubernetes %}<span class="badge badge-info" data-toggle="tooltip" title="This option works for the Kubernetes orchestrator.">Kubernetes</span>{% endif %}{% endcapture %}
  {% capture all-badges %}{{ deprecated-badge }}{{ experimental-daemon-badge }}{{ experimental-cli-badge }}{{ min-api }}{{ flag-orchestrator }}{% endcapture %}
  {% assign defaults-to-skip = "[],map[],false,0,0s,default,'',\"\"" | split: ',' %}
  {% capture option-default %}{% if option.default_value %}{% unless defaults-to-skip contains option.default_value or defaults-to-skip == blank %}`{{ option.default_value }}`{% endunless %}{% endif %}{% endcapture %}
  <tr>
    {% if option.details_url and option.details_url != '' -%}
    <td markdown="span">[`--{{ option.option }}`]({{ option.details_url }}){% if option.shorthand %} , [`-{{ option.shorthand }}`]({{ option.details_url }}){% endif %}</td>
    {%- else -%}
    <td markdown="span">`--{{ option.option }}`{% if option.shorthand %} , `-{{ option.shorthand }}`{% endif %}</td>
    {%- endif %}
    <td markdown="span">{{ option-default }}</td>
    <td markdown="span">{% if all-badges != '' %}{{ all-badges | strip }}<br />{% endif %}{{ option.description | strip | escape }}</td>
  </tr>
{% endfor %} <!-- end for option -->
</tbody>
</table>
{% endif %} <!-- end if options -->

{% if controller_data.examples %}

## Examples

{{ controller_data.examples | replace_relative_links: page.path }}

{% endif %}

{% if controller_data.pname %}
{% unless controller_data.pname == include.datafile or controller_data.pname == "docker" %}

## Parent command

{% capture parentfile %}{{ controller_data.plink | remove_first: ".yaml" | remove_first: "docker_" }}{% endcapture %}
{% capture parentdatafile %}{{ controller_data.plink | remove_first: ".yaml" }}{% endcapture %}
{% capture parentDesc %}{{ site.data[include.datafolder][parentdatafile].short }}{% endcapture %}

| Command                                                        | Description      |
|:---------------------------------------------------------------|:-----------------|
| [{{ controller_data.pname }}]({{parentPath}}{{ parentfile }}/) | {{ parentDesc }} |

{% endunless %}
{% endif %}

{% if controller_data.cname %}

## Child commands

<table>
<thead>
  <tr>
    <td>Command</td>
    <td>Description</td>
  </tr>
</thead>
<tbody>
{% for command in controller_data.cname %}
  {% capture dataFileName %}{{ command | strip | replace: " ", "_" }}{% endcapture %}
  <tr>
    <td markdown="span">[{{ command }}]({{ parentPath }}{{ dataFileName | remove_first: "docker_" }}/)</td>
    <td markdown="span">{{ site.data[include.datafolder][dataFileName].short }}</td>
  </tr>
{% endfor %}
</tbody>
</table>
{% endif %}

{% unless controller_data.pname == "docker" or controller_data.pname == "dockerd" or include.datafile=="docker" %}

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
  {% capture dataFileName %}{{ command | strip | replace: " ", "_" }}{% endcapture %}
  <tr>
    <td markdown="span">[{{ command }}]({{ parentPath }}{{ dataFileName | remove_first: "docker_" }}/)</td>
    <td markdown="span">{{ site.data[include.datafolder][dataFileName].short }}</td>
  </tr>
{% endfor %}
</tbody>
</table>

{% endunless %}
