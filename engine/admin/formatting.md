---
description: CLI and log output formatting reference
keywords: format, formatting, output, templates, log
title: Format command and log output
---

Docker uses [Go templates](https://golang.org/pkg/text/template/) to allow users manipulate the output format
of certain commands and log drivers. Each command a driver provides a detailed
list of elements they support in their templates:

- [Docker Images formatting](../reference/commandline/images.md#formatting)
- [Docker Inspect formatting](../reference/commandline/inspect.md#examples)
- [Docker Log Tag formatting](logging/log_tags.md)
- [Docker Network Inspect formatting](../reference/commandline/network_inspect.md)
- [Docker PS formatting](../reference/commandline/ps.md#formatting)
- [Docker Volume Inspect formatting](../reference/commandline/volume_inspect.md)
- [Docker Version formatting](../reference/commandline/version.md#examples)

## Template functions

Docker provides a set of basic functions to manipulate template elements.
This is the complete list of the available functions with examples:

### Join

Join concatenates a list of strings to create a single string.
It puts a separator between each element in the list.

	{% raw %}
	$ docker ps --format '{{join .Names " or "}}'
	{% endraw %}

### Json

Json encodes an element as a json string.

	{% raw %}
	$ docker inspect --format '{{json .Mounts}}' container
	{% endraw %}

### Lower

Lower turns a string into its lower case representation.

	{% raw %}
	$ docker inspect --format "{{lower .Name}}" container
	{% endraw %}

### Split

Split slices a string into a list of strings separated by a separator.

	{% raw %}
	$ docker inspect --format '{{split (join .Names "/") "/"}}' container
	{% endraw %}

### Title

Title capitalizes a string.

	{% raw %}
	$ docker inspect --format "{{title .Name}}" container
	{% endraw %}

### Upper

Upper turns a string into its upper case representation.

	{% raw %}
	$ docker inspect --format "{{upper .Name}}" container
	{% endraw %}