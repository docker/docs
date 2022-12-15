---
title: Build attestations
keywords: build, attestations, sbom, provenance
description: >
  Introduction to SBOM and provenance attestations with Docker Build; what they
  are and why they exist
---

Build attestations describe how an image was built, and what it contains. The
attestations are created at build-time by BuildKit, and attach to the final
image as metadata.

The purpose of attestations is to make it possible to inspect an image and see
where it comes from, who created it and how, and what it contains. This enables
policy engines to evaluate and validate an image.

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

BuildKit generates the attestations when building the image. The attestation
records are wrapped in the in-toto JSON format and attached it to the image
index in a manifest for the final image.

## Storage

<!-- prettier-ignore -->
BuildKit produces attestations in the
[in-toto format](https://github.com/in-toto/attestation){: target="blank" rel="noopener" class="\_" },
as defined by the 
[in-toto framework](https://in-toto.io/){: target="blank" rel="noopener" class="\_" },
a standard supported by the Linux Foundation.

Attestations attach to images as a manifest in the image index. The data records
of the attestations are stored as JSON blobs.

Because attestations attach to images as a manifest, it means that you can
inspect the attestations for any image in a registry without having to pull the
whole image.

The following example shows a truncated in-toto JSON representation of an SBOM
attestation. The `subject` key contains the index of software artifacts included
in the image.

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "predicateType": "https://spdx.dev/Document",
  "subject": [
    {
      "name": "bin/awk",
      "digest": {
        "sha256": "e99b0b53b1ede6f76e8a48451d29d1554c04c9d2c88da68519cfefd01d648681"
      }
    },
    {
      "name": "bin/base64",
      "digest": {
        "sha256": "e99b0b53b1ede6f76e8a48451d29d1554c04c9d2c88da68519cfefd01d648681"
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
    "packages": [],
    "spdxVersion": "SPDX-2.2"
  }
}
```

<!-- prettier-ignore -->
To deep-dive into the specifics about how attestations are stored, see
[Image Attestation Storage (BuildKit)](https://github.com/moby/buildkit/blob/master/docs/attestation-storage.md){: target="blank" rel="noopener" class="_"}.

## What's next

Learn more about the available attestation types and how to use them:

- [Provenance](provenance.md)
- [SBOM](sbom.md)
