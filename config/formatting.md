---
description: CLI and log output formatting reference
keywords: format, formatting, output, templates, log
title: Format command and log output
redirect_from:
- /engine/admin/formatting/
---

Docker uses [Go templates](https://golang.org/pkg/text/template/) which you can
use to manipulate the output format of certain commands and log drivers.

Docker provides a set of basic functions to manipulate template elements.
All of these examples use the `docker inspect` command, but many other CLI
commands have a `--format` flag, and many of the CLI command references
include examples of customizing the output format.

>**Note**
>
> When using the `--format` flag, you need observe your shell environment. 
> In a Posix shell, you can run the following with a single quote:
>
> {% raw %}
> ```console
> $ docker inspect --format '{{join .Args " , "}}'
> ```
> {% endraw %}
>
> Otherwise, in a Windows shell (for example, PowerShell), you need to use single quotes, but
> escape the double quotes inside the params as follows:
>
> {% raw %}
> ```console
> $ docker inspect --format '{{join .Args \" , \"}}'
> ```
> {% endraw %}
>
{:.important}

## join

`join` concatenates a list of strings to create a single string.
It puts a separator between each element in the list.

{% raw %}
```console
$ docker inspect --format '{{join .Args " , "}}' container
```
{% endraw %}

## table

`table` specifies which fields you want to see its output.

{% raw %}
```console
$ docker image list --format "table {{.ID}}\t{{.Repository}}\t{{.Tag}}\t{{.Size}}"
```
{% endraw %}

## json

`json` encodes an element as a json string.


{% raw %}
```console
$ docker inspect --format '{{json .Mounts}}' container
```
{% endraw %}

## lower

`lower` transforms a string into its lowercase representation.

{% raw %}
```console
$ docker inspect --format "{{lower .Name}}" container
```
{% endraw %}

## split

`split` slices a string into a list of strings separated by a separator.

{% raw %}
```console
$ docker inspect --format '{{split .Image ":"}}'
```
{% endraw %}

## title

`title` capitalizes the first character of a string.

{% raw %}
```console
$ docker inspect --format "{{title .Name}}" container
```
{% endraw %}

## upper

`upper` transforms a string into its uppercase representation.

{% raw %}
```console
$ docker inspect --format "{{upper .Name}}" container
```
{% endraw %}


## println

`println` prints each value on a new line.

{% raw %}
```console
$ docker inspect --format='{{range .NetworkSettings.Networks}}{{println .IPAddress}}{{end}}' container
```
{% endraw %}

# Hint

To find out what data can be printed, show all content as json:

{% raw %} 
```console
$ docker container ls --format='{{json .}}'
```
{% endraw %} 
