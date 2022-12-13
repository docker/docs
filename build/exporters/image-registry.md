---
title: "Image and registry exporters"
keywords: >
  build, buildx, buildkit, exporter, image, registry
redirect_from:
  - /build/building/exporters/image-registry/
---

The `image` exporter outputs the build result into a container image format. The
`registry` exporter is identical, but it automatically pushes the result by
setting `push=true`.

## Synopsis

Build a container image using the `image` and `registry` exporters:

```console
$ docker buildx build --output type=image[,parameters] .
```

```console
$ docker buildx build --output type=registry[,parameters] .
```

The following table describes the available parameters that you can pass to
`--output` for `type=image`:

| Parameter              | Type                                   | Default | Description                                                                                                                                                                                                                         |
|------------------------|----------------------------------------|---------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`                 | String                                 |         | Specify image name(s)                                                                                                                                                                                                               |
| `push`                 | `true`,`false`                         | `false` | Push after creating the image.                                                                                                                                                                                                      |
| `push-by-digest`       | `true`,`false`                         | `false` | Push image without name.                                                                                                                                                                                                            |
| `registry.insecure`    | `true`,`false`                         | `false` | Allow pushing to insecure registry.                                                                                                                                                                                                 |
| `dangling-name-prefix` | `<value>`                              |         | Name image with `prefix@<digest>`, used for anonymous images                                                                                                                                                                        |
| `name-canonical`       | `true`,`false`                         |         | Add additional canonical name `name@<digest>`                                                                                                                                                                                       |
| `compression`          | `uncompressed`,`gzip`,`estargz`,`zstd` | `gzip`  | Compression type, see [compression][1]                                                                                                                                                                                              |
| `compression-level`    | `0..22`                                |         | Compression level, see [compression][1]                                                                                                                                                                                             |
| `force-compression`    | `true`,`false`                         | `false` | Forcefully apply compression, see [compression][1]                                                                                                                                                                                  |
| `oci-mediatypes`       | `true`,`false`                         | `false` | Use OCI media types in exporter manifests, see [OCI Media types][2]                                                                                                                                                                 |
| `buildinfo`            | `true`,`false`                         | `true`  | Attach inline [build info][3]                                                                                                                                                                                                       |
| `buildinfo-attrs`      | `true`,`false`                         | `false` | Attach inline [build info attributes][3]                                                                                                                                                                                            |
| `unpack`               | `true`,`false`                         | `false` | Unpack image after creation (for use with containerd)                                                                                                                                                                               |
| `store`                | `true`,`false`                         | `true`  | Store the result images to the worker's (for example, containerd) image store, and ensures that the image has all blobs in the content store. Ignored if the worker doesn't have image store (when using OCI workers, for example). |
| `annotation.<key>`     | String                                 |         | Attach an annotation with the respective `key` and `value` to the built image,see [annotations][4]                                                                                                                                  |

[1]: index.md#compression
[2]: index.md#oci-media-types
[3]: index.md#build-info
[4]: #annotations

## Annotations

These exporters support adding OCI annotation using `annotation.*` dot notation
parameter. The following example sets the `org.opencontainers.image.title`
annotation for a build:

```console
$ docker buildx build \
    --output "type=<type>,name=<registry>/<image>,annotation.org.opencontainers.image.title=<title>" .
```

For more information about annotations, see
[BuildKit documentation](https://github.com/moby/buildkit/blob/master/docs/annotations.md){:target="blank" rel="noopener" class=""}.

## Further reading

For more information on the `image` or `registry` exporters, see the
[BuildKit README](https://github.com/moby/buildkit/blob/master/README.md#imageregistry){:target="blank" rel="noopener" class=""}.
