---
datafolder: engine-cli
datafile: docker_image_ls
title: docker image ls
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

### Listing the images

To list the images in a local repository (not the registry) run:

    $ docker image ls

The list will contain the image repository name, a tag for the image, and an
image ID, when it was created and its virtual size. Columns: REPOSITORY, TAG,
IMAGE ID, CREATED, and SIZE.

The `docker image ls` command takes an optional `[REPOSITORY[:TAG]]` argument
that restricts the list to images that match the argument. If you specify
`REPOSITORY`but no `TAG`, the `docker image ls` command lists all images in the
given repository.

    $ docker image ls java

The `[REPOSITORY[:TAG]]` value must be an "exact match". This means that, for example,
`docker image ls jav` does not match the image `java`.

If both `REPOSITORY` and `TAG` are provided, only images matching that
repository and tag are listed.  To find all local images in the "java"
repository with tag "8" you can use:

    $ docker image ls java:8

To get a verbose list of images which contains all the intermediate images
used in builds use **-a**:

    $ docker image ls -a

Previously, the docker image ls command supported the --tree and --dot arguments,
which displayed different visualizations of the image data. Docker core removed
this functionality in the 1.7 version. If you liked this functionality, you can
still find it in the third-party dockviz tool: https://github.com/justone/dockviz.

### Listing images in a desired format

When using the --format option, the image command will either output the data
exactly as the template declares or, when using the `table` directive, will
include column headers as well. You can use special characters like `\t` for
inserting tab spacing between columns.

The following example uses a template without headers and outputs the ID and
Repository entries separated by a colon for all images:

```bash
{% raw %}
$ docker images --format "{{.ID}}: {{.Repository}}"

77af4d6b9913: <none>
b6fa739cedf5: committ
78a85c484bad: ipbabble
30557a29d5ab: docker
5ed6274db6ce: <none>
746b819f315e: postgres
746b819f315e: postgres
746b819f315e: postgres
746b819f315e: postgres
{% endraw %}
```

To list all images with their repository and tag in a table format you can use:

```bash
{% raw %}
$ docker images --format "table {{.ID}}\t{{.Repository}}\t{{.Tag}}"

IMAGE ID            REPOSITORY                TAG
77af4d6b9913        <none>                    <none>
b6fa739cedf5        committ                   latest
78a85c484bad        ipbabble                  <none>
30557a29d5ab        docker                    latest
5ed6274db6ce        <none>                    <none>
746b819f315e        postgres                  9
746b819f315e        postgres                  9.3
746b819f315e        postgres                  9.3.5
746b819f315e        postgres                  latest
{% endraw %}
```

Valid template placeholders are listed above.

### Listing only the shortened image IDs

Listing just the shortened image IDs. This can be useful for some automated
tools.

    $ docker image ls -q
