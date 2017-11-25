# docs-base

This is the source for [docs/docker.github.io:docs-base](https://hub.docker.com/r/docs/docker.github.io/tags/).

The docs-base Dockerfile includes:

- Static HTML from each doc archive except master (using their Dockerfiles)
- A temporary index.html to make it browseable
- A script that runs across those builds to fix relative links

The Nginx config and Github Pages versions are no longer here! They are in the
`docs-builder` branch and the `docs/docker.github.io:docs-builder` image!

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

Adding a new archive version
