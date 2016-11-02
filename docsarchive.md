---
title: Documentation Archive
---

This page lists the various ways you can view the docs as they were when a
prior version of Docker was shipped.

## View the docs archives locally

The docs archive is published as a [Docker repository at docs/archive](https://hub.docker.com/r/docs/archive/tags/).
To see any of these versions, run the following command, changing
the tag from `v1.4` to any tag you see in [the repo](https://hub.docker.com/r/docs/archive/tags/):

```shell
docker run -p 4000:4000 docs/archive:v1.4
```

The docs for `v1.4` will then be viewable at `http://localhost:4000`.

## Viewing the docs archives online

[This Docker Compose file](https://github.com/docker/docker.github.io/blob/master/_data/docsarchive/docker-compose.yml)
is used to stand up the images in docs/archive on an AWS instance, using the
following locations.

{% for item in site.data.docsarchive.docker-compose %}

### {{ item[0] }}

Docs for {{ item[0] }} are at [/{{ item[0] }}](http://54.71.194.30:{{ item[1].ports[0] | replace:':4000','' }})

{% endfor %}
