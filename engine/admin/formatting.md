---
description: CLI and log output formatting reference
keywords: format, formatting, output, templates, log
title: Format command and log output
---

Docker uses [Go templates](https://golang.org/pkg/text/template/) which allow users to manipulate the output format
of certain commands and log drivers. Each command a driver provides has a detailed
list of elements they support in their templates:

- [Docker Images formatting](../reference/commandline/images.md#formatting)
- [Docker Inspect formatting](../reference/commandline/inspect.md#examples)
- [Docker Log Tag formatting](logging/log_tags.md)
- [Docker Network Inspect formatting](../reference/commandline/network_inspect.md)
- [Docker PS formatting](../reference/commandline/ps.md#formatting)
- [Docker Stats formatting](../reference/commandline/stats.md#formatting)
- [Docker Volume Inspect formatting](../reference/commandline/volume_inspect.md)
- [Docker Version formatting](../reference/commandline/version.md#examples)

## Template functions

Docker provides a set of basic functions to manipulate template elements.
This is the complete list of the available functions with examples:

### `join`

`join` concatenates a list of strings to create a single string.
It puts a separator between each element in the list.

	{% raw %}
	$ docker inspect --format '{{join .Args " , "}}' container
	{% endraw %}

### `json`

`json` encodes an element as a json string.

	{% raw %}
	$ docker inspect --format '{{json .Mounts}}' container
	{% endraw %}

### `lower`

`lower` transforms a string into its lowercase representation.

	{% raw %}
	$ docker inspect --format "{{lower .Name}}" container
	{% endraw %}

### `split`

`split` slices a string into a list of strings separated by a separator.

	{% raw %}
	$ docker inspect --format '{{split (join .Names "/") "/"}}' container
  {% endraw %}

### `title`

`title` capitalizes the first character of a string.

	{% raw %}
	$ docker inspect --format "{{title .Name}}" container
	{% endraw %}

### `upper`

`upper` transforms a string into its uppercase representation.

	{% raw %}
	$ docker inspect --format "{{upper .Name}}" container
	{% endraw %}
