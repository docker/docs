---
title: Local cache
description: Manage build cache with a local directory
keywords: build, buildx, cache, backend, local
aliases:
  - /build/building/cache/backends/local/
---

The `local` cache store is a simple cache option that stores your cache as files
in a directory on your filesystem, using an
[OCI image layout](https://github.com/opencontainers/image-spec/blob/main/image-layout.md)
for the underlying directory structure. Local cache is a good choice if you're
just testing, or if you want the flexibility to self-manage a shared storage
solution.

## Synopsis

```console
$ docker buildx build --push -t <registry>/<image> \
  --cache-to type=local,dest=path/to/local/dir[,parameters...] \
  --cache-from type=local,src=path/to/local/dir .
```

The following table describes the available CSV parameters that you can pass to
`--cache-to` and `--cache-from`.

| Name                | Option       | Type                    | Default | Description                                                                                                                     |
|---------------------|--------------|-------------------------|---------|---------------------------------------------------------------------------------------------------------------------------------|
| `src`               | `cache-from` | String                  |         | Path of the local directory where cache gets imported from.                                                                     |
| `digest`            | `cache-from` | String                  |         | Digest of manifest to import, see [cache versioning][4].                                                                        |
| `tag`               | `cache-to`,`cache-from` | String                  | `latest` | Tag of the cache manifest, see [cache versioning][4].                                                                          |
| `dest`              | `cache-to`   | String                  |         | Path of the local directory where cache gets exported to.                                                                       |
| `mode`              | `cache-to`   | `min`,`max`             | `min`   | Cache layers to export, see [cache mode][1].                                                                                    |
| `oci-mediatypes`    | `cache-to`   | `true`,`false`          | `true`  | Use OCI media types in exported manifests, see [OCI media types][2].                                                            |
| `image-manifest`    | `cache-to`   | `true`,`false`          | `true`  | When using OCI media types, generate an image manifest instead of an image index for the cache image, see [OCI media types][2]. |
| `compression`       | `cache-to`   | `gzip`,`estargz`,`zstd` | `gzip`  | Compression type, see [cache compression][3].                                                                                   |
| `compression-level` | `cache-to`   | `0..22`                 |         | Compression level, see [cache compression][3].                                                                                  |
| `force-compression` | `cache-to`   | `true`,`false`          | `false` | Forcibly apply compression, see [cache compression][3].                                                                         |
| `ignore-error`      | `cache-to`   | Boolean                 | `false` | Ignore errors caused by failed cache exports.                                                                                   |
| `reset`             | `cache-to`   | `true`,`false`          | `false` | Delete blobs that no tag references, see [cache versioning][4].                                                                 |

[1]: _index.md#cache-mode
[2]: _index.md#oci-media-types
[3]: _index.md#cache-compression
[4]: #cache-versioning

If the `src` cache doesn't exist, then the cache import step will fail, but the
build continues.

## Cache versioning

A local cache directory uses an OCI image layout. Its `index.json` file
associates tags with cache manifests, while the `blobs` directory stores the
manifest and cache data.

By default, BuildKit exports and imports the cache tagged `latest`. Use
different tags to keep multiple caches in the same directory:

```console
$ docker buildx build --cache-to type=local,dest=path/to/local/dir,tag=v1 .
$ docker buildx build --cache-to type=local,dest=path/to/local/dir,tag=v2 .
```

Exporting another cache with the same tag updates that tag to reference the new
manifest. Manifests referenced by other tags remain unchanged.

Import a cache by specifying its tag:

```console
$ docker buildx build --cache-from type=local,src=path/to/local/dir,tag=v1 .
```

A digest identifies an exact cache manifest. BuildKit reports the digest of
each exported manifest in the build output. Use `digest` instead of `tag` when
you need a specific manifest:

```console
$ docker buildx build \
  --cache-from type=local,src=path/to/local/dir,digest=sha256:<DIGEST> .
```

If you specify both `digest` and `tag`, BuildKit uses `digest`.

By default, updating a tag doesn't delete the blobs used by its previous
manifest. The previous manifest remains available by digest, so the local cache
directory grows over time.

Buildx version 0.35.0 and later supports `reset=true` on export, which deletes
blobs that no tag references:

```console
$ docker buildx build --cache-to type=local,dest=path/to/local/dir,reset=true .
```

Blobs that other tags reference are kept. Manifests that no tag references are
deleted, so you can no longer import them by digest.

## Further reading

For an introduction to caching see [Docker build cache](../_index.md).

For more information on the `local` cache backend, see the
[BuildKit README](https://github.com/moby/buildkit#local-directory-1).
