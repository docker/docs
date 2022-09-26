> **Note**
>
> Samples compatible with [Docker Dev Environments](/desktop/dev-environments/) require [Docker Desktop](../get-docker.md) version 4.10 or later.

| Name | Description | Docker Dev Environment (if compatible) |
| ---- | ----------- | -------------------------------------- |
{% for sample in site.data.samples.samples -%}
{% for service in sample.services -%}
{% if service == page.service -%}
| [{{sample.title}}]({{sample.url}}) | {{sample.description}} | {% if sample.dev_env -%} [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url={{sample.url}}) {% else -%}-{% endif -%}|
{% endif -%}
{%- endfor -%}
{%- endfor -%}

{% include_relative samples-footer.md %}