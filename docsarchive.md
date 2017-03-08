---
title: View the docs archives
---

This page lists the various ways you can view the docs as they were when a
prior version of Docker was shipped.

{% for archive in site.data.docsarchive.archives %}

{% if archive.current %}

## {{ archive.name }} (current)

Docs for {{ archive.name }} _(current)_ are accessible at [**https://docs.docker.com/**](/), or
run:

```
docker run -ti -p 4000:4000 {{ archive.image }}
```
{% else %}

## {{ archive.name }}

Docs for {{ archive.name }} are accessible at [**https://docs.docker.com/{{ archive.name }}/**](/{{ archive.name }}/), or
run:

```
docker run -ti -p 4000:4000 {{ archive.image }}
```
{% endif %}
{% endfor %}
