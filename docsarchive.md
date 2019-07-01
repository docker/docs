---
title: View the docs archives
---

This page lists the various ways you can view the docs as they were when a
prior version of Docker was shipped.

**Note**: To access documentation for an unsupported version, refer to 
[Accessing unsupported archived documentation](#accessing-unsupported-archived-documentation). 

{% for archive in site.data.docsarchive.archives %}

{% if archive.current %}

## {{ archive.name }} (current)

Docs for {{ archive.name }} _(current)_ are accessible at [**https://docs.docker.com/**](/), or
to view the docs offline on your local machine, run:

```
docker run -ti -p 4000:4000 {{ archive.image }}
```

{% else %}

{% if archive.name != 'edge' %}

## {{ archive.name }}

Docs for {{ archive.name }} are accessible at [**https://docs.docker.com/{{ archive.name }}/**](/{{ archive.name }}/), or to view the docs offline on your local machine, 
run:

```
docker run -ti -p 4000:4000 {{ archive.image }}
```

{% endif %} <!-- edge check -->
{% endif %}
{% endfor %}

## Accessing unsupported archived documentation

Supported documentation includes the current version plus the previous five versions. 

If you are using a version of the documentation that is no longer supported, which means that the version number is not listed in the site dropdown list, you can still access that documentation in the following ways:

- By entering your version number and selecting it from the branch selection list for this repo 
- By directly accessing the Github URL for your version. For example, https://github.com/docker/docker.github.io/tree/v1.9 for `v1.9` 
- By running a container of the specific [tag for your documentation version](https://cloud.docker.com/u/docs/repository/docker/docs/docker.github.io/general#read-these-docs-offline) 
in Docker Hub. For example, run the following to access `v1.9`:

 ```bash
  docker run  -it -p 4000:4000 docs/docker.github.io:v1.9
  ```
