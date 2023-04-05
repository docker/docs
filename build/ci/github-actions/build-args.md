---
title: Pass build arguments with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, args
---

To pass build arguments to the Dockerfile, you can use the build-args option.

Note that the list of build arguments should be passed as a multi-line string, e.g.

```yaml
      -
        name: Build and push
        uses: docker/build-push-action@v4
        with:
          push: true
          tags: ghcr.io/username/app:latest
          build-args: |
            UBUNTU_VERSION=${{ matrix.ubuntu_version }}
            RUBY_VERSION=${{ matrix.ruby_version }}
```
