---
title: Local cache
description: Manage build cache with Amazon S3 buckets
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
| ------------------- | ------------ | ----------------------- | ------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `src`               | `cache-from` | String                  |         | Path of the local directory where cache gets imported from.                                                                     |
| `digest`            | `cache-from` | String                  |         | Digest of manifest to import, see [cache versioning][4].                                                                        |
| `dest`              | `cache-to`   | String                  |         | Path of the local directory where cache gets exported to.                                                                       |
| `mode`              | `cache-to`   | `min`,`max`             | `min`   | Cache layers to export, see [cache mode][1].                                                                                    |
| `oci-mediatypes`    | `cache-to`   | `true`,`false`          | `true`  | Use OCI media types in exported manifests, see [OCI media types][2].                                                            |
| `image-manifest`    | `cache-to`   | `true`,`false`          | `false` | When using OCI media types, generate an image manifest instead of an image index for the cache image, see [OCI media types][2]. |
| `compression`       | `cache-to`   | `gzip`,`estargz`,`zstd` | `gzip`  | Compression type, see [cache compression][3].                                                                                   |
| `compression-level` | `cache-to`   | `0..22`                 |         | Compression level, see [cache compression][3].                                                                                  |
| `force-compression` | `cache-to`   | `true`,`false`          | `false` | Forcibly apply compression, see [cache compression][3].                                                                         |
| `ignore-error`      | `cache-to`   | Boolean                 | `false` | Ignore errors caused by failed cache exports.                                                                                   |

[1]: _index.md#cache-mode
[2]: _index.md#oci-media-types
[3]: _index.md#cache-compression
[4]: #cache-versioning

If the `src` cache doesn't exist, then the cache import step will fail, but the
build continues.

## Cache versioning

<!-- FIXME: update once https://github.com/moby/buildkit/pull/3111 is released -->

This section describes how versioning works for caches on a local filesystem,
and how you can use the `digest` parameter to use older versions of cache.

If you inspect the cache directory manually, you can see the resulting OCI image
layout:

```console
$ ls cache
blobs  index.json  ingest
$ cat cache/index.json | jq
{
  "schemaVersion": 2,
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.index.v1+json",
      "digest": "sha256:6982c70595cb91769f61cd1e064cf5f41d5357387bab6b18c0164c5f98c1f707",
      "size": 1560,
      "annotations": {
        "org.opencontainers.image.ref.name": "latest"
      }
    }
  ]
}
```

Like other cache types, local cache gets replaced on export, by replacing the
contents of the `index.json` file. However, previous caches will still be
available in the `blobs` directory. These old caches are addressable by digest,
and kept indefinitely. Therefore, the size of the local cache will continue to
grow (see [`moby/buildkit#1896`](https://github.com/moby/buildkit/issues/1896)
for more information).

When importing cache using `--cache-to`, you can specify the `digest` parameter
to force loading an older version of the cache, for example:

```console
$ docker buildx build --push -t <registry>/<image> \
  --cache-to type=local,dest=path/to/local/dir \
  --cache-from type=local,ref=path/to/local/dir,digest=sha256:6982c70595cb91769f61cd1e064cf5f41d5357387bab6b18c0164c5f98c1f707 .
```

## Further reading

For an introduction to caching see [Docker build cache](../_index.md).

For more information on the `local` cache backend, see the
[BuildKit README](https://github.com/moby/buildkit#local-directory-1).
