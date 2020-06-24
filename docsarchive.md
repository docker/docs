---
title: View the docs archives
---

This page lists the various ways you can view the docs as they were when a
prior version of Docker was shipped.

To view the docs offline on your local machine, run:

```
docker run -ti -p 4000:4000 {{ archive.image }}
```

## Accessing unsupported archived documentation

If you are using a version of the documentation that is no longer supported,
you can still access that documentation in the following ways:

- By entering your version number and selecting it from the branch selection list for this repo
- By directly accessing the Github URL for your version. For example, https://github.com/docker/docker.github.io/tree/v1.9 for `v1.9`
- By running a container of the specific [tag for your documentation version](https://hub.docker.com/r/docs/docker.github.io)
  in Docker Hub. For example, run the following to access `v1.9`:

  ```bash
  docker run  -it -p 4000:4000 docs/docker.github.io:v1.9
  ```
