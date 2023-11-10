---
title: Build attestations
keywords: build, attestations, sbom, provenance, metadata
description: |
  Introduction to SBOM and provenance attestations with Docker Build,
  what they are, and why they exist
---

Build attestations describe how an image was built, and what it contains. The
attestations are created at build-time by BuildKit, and become attached to the
final image as metadata.

The purpose of attestations is to make it possible to inspect an image and see
where it comes from, who created it and how, and what it contains. This enables
you to make informed decisions about how an image impacts the supply chain security
of your application. It also enables the use of policy engines for validating
images based on policy rules you've defined.

Two types of build annotations are available:

- Software Bill of Material (SBOM): list of software artifacts that an image
  contains, or that were used to build the image.
- Provenance: how an image was built.

## Purpose of attestations

The use of open source and third-party packages is more widespread than ever
before. Developers share and reuse code because it helps increase productivity,
allowing teams to create better products, faster.

Importing and using code created elsewhere without vetting it introduces a
severe security risk. Even if you do review the software that you consume, new
zero-day vulnerabilities are frequently discovered, requiring development teams
take action to remediate them.

Build attestations make it easier to see the contents of an image, and where it
comes from. Use attestations to analyze and decide whether to use an image, or
to see if images you are already using are exposed to vulnerabilities.

## Creating attestations

When you build an image with `docker buildx build`, you can add attestation
records to the resulting image using the `--provenance` and `--sbom` options.
You can opt in to add either the SBOM or provenance attestation type, or both.

```console
$ docker buildx build --sbom=true --provenance=true .
```

> **Note**
>
> The default image store doesn't support attestations. If you're using the
> default image store and you build an image using the default `docker` driver,
> or using a different driver with the `--load` flag, the attestations are
> lost.
>
> To make sure the attestations are preserved, you can:
>
> - Use a `docker-container` driver with the `--push` flag to push the image to
>   a registry directly.
> - Enable the [containerd image store](../../desktop/containerd/_index.md).

> **Note**
>
> Provenance attestations are enabled by default, with the `mode=min` option.
> You can disable provenance attestations using the `--provenance=false` flag,
> or by setting the [`BUILDX_NO_DEFAULT_ATTESTATIONS`](../building/env-vars.md#buildx_no_default_attestations) environment variable.
>
> Using the `--provenance=true` flag attaches provenance attestations with `mode=max`
> by default. See [Provenance attestation](./slsa-provenance.md) for more details.

BuildKit generates the attestations when building the image. The attestation
records are wrapped in the in-toto JSON format and attached to the image
index in a manifest for the final image.

## Storage

BuildKit produces attestations in the [in-toto format](https://github.com/in-toto/attestation),
as defined by the [in-toto framework](https://in-toto.io/),
a standard supported by the Linux Foundation.

Attestations attach to images as a manifest in the image index. The data records
of the attestations are stored as JSON blobs.

Because attestations attach to images as a manifest, it means that you can
inspect the attestations for any image in a registry without having to pull the
whole image.

All BuildKit exporters support attestations. The `local` and `tar` can't save
the attestations to an image manifest, since it's outputting a directory of
files or a tarball, not an image. Instead, these exporters write the
attestations to one or more JSON files in the root directory of the export.

The following example shows a truncated in-toto JSON representation of an SBOM
attestation.

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "predicateType": "https://spdx.dev/Document",
  "subject": [
    {
      "name": "pkg:docker/<registry>/<image>@<tag/digest>?platform=<platform>",
      "digest": {
        "sha256": "e8275b2b76280af67e26f068e5d585eb905f8dfd2f1918b3229db98133cb4862"
      }
    }
  ],
  "predicate": {
    "SPDXID": "SPDXRef-DOCUMENT",
    "creationInfo": {
      "created": "2022-12-15T11:47:54.546747383Z",
      "creators": ["Organization: Anchore, Inc", "Tool: syft-v0.60.3"],
      "licenseListVersion": "3.18"
    },
    "dataLicense": "CC0-1.0",
    "documentNamespace": "https://anchore.com/syft/dir/run/src/core-da0f600b-7f0a-4de0-8432-f83703e6bc4f",
    "name": "/run/src/core",
    // list of files that the image contains, e.g.:
    "files": [
      {
        "SPDXID": "SPDXRef-1ac501c94e2f9f81",
        "comment": "layerID: sha256:9b18e9b68314027565b90ff6189d65942c0f7986da80df008b8431276885218e",
        "fileName": "/bin/busybox",
        "licenseConcluded": "NOASSERTION"
      }
    ],
    // list of packages that were identified for this image:
    "packages": [
      {
        "name": "busybox",
        "originator": "Person: Sören Tempel <soeren+alpine@soeren-tempel.net>",
        "sourceInfo": "acquired package info from APK DB: lib/apk/db/installed",
        "versionInfo": "1.35.0-r17",
        "SPDXID": "SPDXRef-980737451f148c56",
        "description": "Size optimized toolbox of many common UNIX utilities",
        "downloadLocation": "https://busybox.net/",
        "licenseConcluded": "GPL-2.0-only",
        "licenseDeclared": "GPL-2.0-only"
        // ...
      }
    ],
    // files-packages relationship
    "relationships": [
      {
        "relatedSpdxElement": "SPDXRef-1ac501c94e2f9f81",
        "relationshipType": "CONTAINS",
        "spdxElementId": "SPDXRef-980737451f148c56"
      },
      ...
    ],
    "spdxVersion": "SPDX-2.2"
  }
}
```

To deep-dive into the specifics about how attestations are stored, see
[Image Attestation Storage (BuildKit)](attestation-storage.md).

## What's next

Learn more about the available attestation types and how to use them:

- [Provenance](slsa-provenance.md)
- [SBOM](sbom.md)
