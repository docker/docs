# docs-base

This is the source for [docs/docker.github.io:docs-base](https://hub.docker.com/r/docs/docker.github.io/tags/).

The docs-base Dockerfile includes:

- The GitHub Pages environment (w/Jekyll)
- nginx
- Builds of all previous versions of Docker's documentation

Having this large amount of stuff that stays relatively static in a base image
helps keep build times for the docs low as we can use Docker Cloud's caching
when running auto-builds out of GitHub.

While you would only see the docs archives by doing so, you can run docs-base
locally and peruse by running:

```
docker run -ti -p 4000:4000 docs/docker.github.io:docs-base
```

The contents of docs-base will then be viewable in your browser at
`localhost:4000`.

## Reasons to update this branch

- Changing the nginx configuration
- Publishing a new docs archive, or anything else that stays static between
  doc builds.
- Updating the GitHub Pages version used in production to a newer release
