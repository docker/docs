---
title: View the docs archives
---

This page lists the various ways you can view the docs as they were when a
prior version of Docker was shipped.

## View the docs archives locally

The docs archive is published as a [Docker Hub repository at docs/docker.github.io](https://hub.docker.com/r/docs/docker.github.io/tags/).
To see any of these versions, run the following command, changing
the tag from `v1.4` to any tag you see in [the repo](https://hub.docker.com/r/docs/docker.github.io/tags/):

```shell
docker run -p 4000:4000 docs/docker.github.io:v1.4
```

The docs for `v1.4` will then be viewable at `http://localhost:4000`.

## View the docs archives online

{% for item in site.data.docsarchive.docker-compose %}

### {{ item[0] }}

Docs for {{ item[0] }} are accessible at [/{{ item[0] }}/](/{{ item[0] }}/).

{% endfor %}
