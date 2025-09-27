---
title: Annotations
description: Annotations specify additional metadata about OCI images
keywords: build, buildkit, annotations, metadata
aliases:
- /build/building/annotations/
---

<!-- vale Docker.Spacing = NO -->

Annotations provide descriptive metadata for images. Use annotations to record
arbitrary information and attach it to your image, which helps consumers and
tools understand the origin, contents, and how to use the image.

Annotations are similar to, and in some sense overlap with, [labels]. Both
serve the same purpose: to attach metadata to a resource. As a general principle,
you can think of the difference between annotations and labels as follows:

- Annotations describe OCI image components, such as [manifests], [indexes],
  and [descriptors].
- Labels describe Docker resources, such as images, containers, networks, and
  volumes.

The OCI image [specification] defines the format of annotations, as well as a set
of pre-defined annotation keys. Adhering to the specified standards ensures
that metadata about images can be surfaced automatically and consistently, by
tools like Docker Scout.

Annotations are not to be confused with [attestations]:

- Attestations contain information about how an image was built and what it contains.
  An attestation is attached as a separate manifest on the image index.
  Attestations are not standardized by the Open Container Initiative.
- Annotations contain arbitrary metadata about an image.
  Annotations attach to the image [config] as labels,
  or on the image index or manifest as properties.

## Add annotations

You can add annotations to an image at build-time, or when creating the image
manifest or index.

> [!NOTE]
>
> The Docker Engine image store doesn't support loading images with
> annotations. To build with annotations, make sure to push the image directly
> to a registry, using the `--push` CLI flag or the
> [registry exporter](/manuals/build/exporters/image-registry.md).

To specify annotations on the command line, use the `--annotation` flag for the
`docker build` command:

```console
$ docker build --push --annotation "foo=bar" .
```

If you're using [Bake](/manuals/build/bake/_index.md), you can use the `annotations`
attribute to specify annotations for a given target:

```hcl
target "default" {
  output = ["type=registry"]
  annotations = ["foo=bar"]
}
```

For examples on how to add annotations to images built with GitHub Actions, see
[Add image annotations with GitHub Actions](/manuals/build/ci/github-actions/annotations.md)

You can also add annotations to an image created using `docker buildx
imagetools create`. This command only supports adding annotations to an index
or manifest descriptors, see
[CLI reference](/reference/cli/docker/buildx/imagetools/create.md#annotation).

## Inspect annotations

To view annotations on an **image index**, use the `docker buildx imagetools
inspect` command. This shows you any annotations for the index and descriptors
(references to manifests) that the index contains. The following example shows
an `org.opencontainers.image.documentation` annotation on a descriptor, and an
`org.opencontainers.image.authors` annotation on the index.

```console {hl_lines=["10-12","19-21"]}
$ docker buildx imagetools inspect <IMAGE> --raw
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "digest": "sha256:d20246ef744b1d05a1dd69d0b3fa907db007c07f79fe3e68c17223439be9fefb",
      "size": 911,
      "annotations": {
        "org.opencontainers.image.documentation": "https://foo.example/docs",
      },
      "platform": {
        "architecture": "amd64",
        "os": "linux"
      }
    },
  ],
  "annotations": {
    "org.opencontainers.image.authors": "dvdksn"
  }
}
```

To inspect annotations on a manifest, use the `docker buildx imagetools
inspect` command and specify `<IMAGE>@<DIGEST>`, where `<DIGEST>` is the digest
of the manifest:

```console {hl_lines="22-25"}
$ docker buildx imagetools inspect <IMAGE>@sha256:d20246ef744b1d05a1dd69d0b3fa907db007c07f79fe3e68c17223439be9fefb --raw
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "digest": "sha256:4368b6959a78b412efa083c5506c4887e251f1484ccc9f0af5c406d8f76ece1d",
    "size": 850
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:2c03dbb20264f09924f9eab176da44e5421e74a78b09531d3c63448a7baa7c59",
      "size": 3333033
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:4923ad480d60a548e9b334ca492fa547a3ce8879676685b6718b085de5aaf142",
      "size": 61887305
    }
  ],
  "annotations": {
    "index,manifest:org.opencontainers.image.vendor": "foocorp",
    "org.opencontainers.image.source": "https://git.example/foo.git",
  }
}
```

## Specify annotation level

By default, annotations are added to the image manifest. You can specify which
level (OCI image component) to attach the annotation to by prefixing the
annotation string with a special type declaration:

```console
$ docker build --annotation "<TYPE>:<KEY>=<VALUE>" .
```

The following types are supported:

- `manifest`: annotates manifests.
- `index`: annotates the root index.
- `manifest-descriptor`: annotates manifest descriptors in the index.
- `index-descriptor`:  annotates the index descriptor in the image layout.

For example, to build an image with the annotation `foo=bar` attached to the
image index:

```console
$ docker build --tag <IMAGE> --push --annotation "index:foo=bar" .
```

Note that the build must produce the component that you specify, or else the
build will fail. For example, the following does not work, because the `docker`
exporter does not produce an index:

```console
$ docker build --output type=docker --annotation "index:foo=bar" .
```

Likewise, the following example also does not work, because buildx creates a
`docker` output by default under some circumstances, such as when provenance
attestations are explicitly disabled:

```console
$ docker build --provenance=false --annotation "index:foo=bar" .
```

It is possible to specify types, separated by a comma, to add the annotation to
more than one level. The following example creates an image with the annotation
`foo=bar` on both the image index and the image manifest:

```console
$ docker build --tag <IMAGE> --push --annotation "index,manifest:foo=bar" .
```

You can also specify a platform qualifier within square brackets in the type
prefix, to annotate only components matching specific OS and architectures. The
following example adds the `foo=bar` annotation only to the `linux/amd64`
manifest:

```console
$ docker build --tag <IMAGE> --push --annotation "manifest[linux/amd64]:foo=bar" .
```

## Related information

Related articles:

- [Add image annotations with GitHub Actions](/manuals/build/ci/github-actions/annotations.md)
- [Annotations OCI specification][specification]

Reference information:

- [`docker buildx build --annotation`](/reference/cli/docker/buildx/build.md#annotation)
- [Bake file reference: `annotations`](/manuals/build/bake/reference.md#targetannotations)
- [`docker buildx imagetools create --annotation`](/reference/cli/docker/buildx/imagetools/create.md#annotation)

<!-- links -->

[specification]: https://github.com/opencontainers/image-spec/blob/main/annotations.md
[attestations]: /manuals/build/metadata/attestations/_index.md
[config]: https://github.com/opencontainers/image-spec/blob/main/config.md
[descriptors]: https://github.com/opencontainers/image-spec/blob/main/descriptor.md
[indexes]: https://github.com/opencontainers/image-spec/blob/main/image-index.md
[labels]: /manuals/engine/manage-resources/labels.md
[manifests]: https://github.com/opencontainers/image-spec/blob/main/manifest.md
