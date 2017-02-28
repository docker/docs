---
title: View the docs archives
---

This page lists the various ways you can view the docs as they were when a
prior version of Docker was shipped.

{% for item in site.data.docsarchive.docker-compose %}

## {{ item[0] }}

Docs for {{ item[0] }} are accessible at [**https://docs.docker.com/{{ item[0] }}/**](/{{ item[0] }}/), or
run:

```
docker run -ti -p 4000:4000 docs/docker.github.io:{{ item[0] }}
```

{% endfor %}
