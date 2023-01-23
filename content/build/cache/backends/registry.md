---
title: "Registry cache"
keywords: build, buildx, cache, backend, registry
aliases:
  - /build/building/cache/backends/registry/
---

The `registry` cache storage can be thought of as an extension to the `inline`
cache. Unlike the `inline` cache, the `registry` cache is entirely separate from
the image, which allows for more flexible usage - `registry`-backed cache can do
everything that the inline cache can do, and more:

- Allows for separating the cache and resulting image artifacts so that you can
  distribute your final image without the cache inside.
- It can efficiently cache multi-stage builds in `max` mode, instead of only the
  final stage.
- It works with other exporters for more flexibility, instead of only the
  `image` exporter.

> **Note**
>
> This cache storage backend requires using a different driver than the default
> `docker` driver - see more information on selecting a driver
> [here](../../drivers/index.md). To create a new driver (which can act as a
> simple drop-in replacement):
>
> ```console
> $ docker buildx create --use --driver=docker-container
> ```

## Synopsis

Unlike the simpler `inline` cache, the `registry` cache supports several
configuration parameters:

```console
$ docker buildx build --push -t <registry>/<image> \
  --cache-to type=registry,ref=<registry>/<cache-image>[,parameters...] \
  --cache-from type=registry,ref=<registry>/<cache-image> .
```

The following table describes the available CSV parameters that you can pass to
`--cache-to` and `--cache-from`.

| Name                | Option                  | Type                    | Default | Description                                                          |
|---------------------|-------------------------|-------------------------|---------|----------------------------------------------------------------------|
| `ref`               | `cache-to`,`cache-from` | String                  |         | Full name of the cache image to import.                              |
| `dest`              | `cache-to`              | String                  |         | Path of the local directory where cache gets exported to.            |
| `mode`              | `cache-to`              | `min`,`max`             | `min`   | Cache layers to export, see [cache mode][1].                         |
| `oci-mediatypes`    | `cache-to`              | `true`,`false`          | `true`  | Use OCI media types in exported manifests, see [OCI media types][2]. |
| `compression`       | `cache-to`              | `gzip`,`estargz`,`zstd` | `gzip`  | Compression type, see [cache compression][3].                        |
| `compression-level` | `cache-to`              | `0..22`                 |         | Compression level, see [cache compression][3].                       |
| `force-compression` | `cache-to`              | `true`,`false`          | `false` | Forcibly apply compression, see [cache compression][3].              |

[1]: index.md#cache-mode
[2]: index.md#oci-media-types
[3]: index.md#cache-compression

You can choose any valid value for `ref`, as long as it's not the same as the
target location that you push your image to. You might choose different tags
(e.g. `foo/bar:latest` and `foo/bar:build-cache`), separate image names (e.g.
`foo/bar` and `foo/bar-cache`), or even different repositories (e.g.
`docker.io/foo/bar` and `ghcr.io/foo/bar`). It's up to you to decide the
strategy that you want to use for separating your image from your cache images.

If the `--cache-from` target doesn't exist, then the cache import step will
fail, but the build will continue.

## Further reading

For an introduction to caching see [Optimizing builds with cache](../index.md).

For more information on the `registry` cache backend, see the
[BuildKit README](https://github.com/moby/buildkit#registry-push-image-and-cache-separately){:target="blank" rel="noopener" class=""}.
