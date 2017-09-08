{% assign section = include.section %}

{% if section == "option" %}
> **Edge only**: This option is only available in Docker CE Edge versions. See [Docker CE Edge](/edge/index.md).

{% elsif section == "options" %}
> **Edge only**: These options are only available in Docker CE Edge versions. See [Docker CE Edge](/edge/index.md).

{% elsif section == "page" %}
> **Edge only**: This topic is only applicable to Docker CE Edge versions. See [Docker CE Edge](/edge/index.md).

{% elsif section == "cliref" %}
> **Edge only**: This is the CLI reference for Docker CE Edge versions. Some of these options may not be available
> to Docker CE stable or Docker EE. You can
> [view the stable version of this CLI reference]({{ page.url | replace:"/edge/", "/"}})
> or [learn about Docker CE Edge](/edge/index.md).

{% elsif section == "dockerd" %}
> **Edge only**: This is the `dockerd` configuration reference for Docker CE Edge versions. Some of these options may not be available
> to Docker CE stable or Docker EE. You can
> [view the stable version of this `dockerd` configuration reference]({{ page.url | replace:"/edge/", "/"}})
> or [learn about Docker CE Edge](/edge/index.md).

{% endif %}

