---
title: Manage vulnerability exceptions
description: |
  Exceptions let you provide additional context and documentation for how
  vulnerabilities affect your artifacts, and provides the ability to
  suppress non-applicable vulnerabilities
keywords: scout, cves, suppress, vex, exceptions
---

Vulnerabilities found in container images sometimes need additional context.
Just because an image contains a vulnerable package, it doesn't mean that the
vulnerability is exploitable. **Exceptions** in Docker Scout lets you address
false positives in image analysis using Vulnerability Exploitability
eXchange (VEX) documents.

By negating non-applicable vulnerabilities, you can make it easier for yourself
and downstream consumers of your images to understand the security implications
of a vulnerability in the context of an image.

In Docker Scout, exceptions are automatically factored into the results.
If an image contains an exception that flags a CVE as non-applicable,
then that CVE is excluded from analysis results.

## Create an exception

To add an exception to an image, you need a VEX document. VEX is a standard
format for documenting vulnerabilities in the context of a software package or
product.

There are multiple implementations and formats of VEX. Docker Scout supports
the [OpenVex](https://github.com/openvex/spec) implementation. To create an
OpenVEX document, use the [`vexctl`](https://github.com/openvex/vexctl) command
line tool.

The following example command creates a VEX document stating that:

- The software product described by this VEX document is the Docker image
  `example/app:v1`
- The image contains the npm package `express@4.17.1`
- The npm package is affected by a known vulnerability: `CVE-2022-24999`
- The image is unaffected by the CVE, because the vulnerable code is never
  executed in containers that run this image

```console
$ vexctl create \
  --author="author@example.com" \
  --product="pkg:docker/example/app@v1" \
  --subcomponents="pkg:npm/express@4.17.1" \
  --vuln="CVE-2022-24999" \
  --status="not_affected" \
  --justification="vulnerable_code_not_in_execute_path" \
  --file="CVE-2022-24999.vex.json"
```

Here's a description of the options in this example:

`--author`
: The email of the author of the VEX document.

`--product`
: Package URL (PURL) of the Docker image. A PURL is an identifier
  for the image in a standardized format, defined in the PURL
  [specification](https://github.com/package-url/purl-spec/blob/master/PURL-TYPES.rst#docker).

  Docker image PURL strings begin with a `pkg:docker` type prefix, followed by
  the image repository and version (the image tag or SHA256 digest). Unlike
  image tags, where the version is specified like `example/app:v1`, in PURL the
  image repository and version are separated by an `@`.

`--subcomponents`
: PURL of the vulnerable package in the image. In this example, the
  vulnerability exists in an npm package, so the `--subcomponents` PURL is the
  identifier for the npm package name and version (`pkg:npm/express@4.17.1`).
  
  If the same vulnerability exists in multiple packages, `vexctl` lets you
  specify the `--subcomponents` flag multiple times for a single `create`
  command.

  You can also omit `--subcomponents`, in which case the VEX statement applies
  to the entire image.

`--vuln`
: ID of the CVE that the VEX statement addresses.

`--status`
: This is the status label of the vulnerability. This describes the
  relationship between the software (`--product`) and the CVE (`--vuln`).
  The possible values for the status label in OpenVEX are:

  - `not_affected`
  - `affected`
  - `fixed`
  - `under_investigation`

  In this example, the VEX statement asserts that the Docker image is
  `not_affected` by the vulnerability. The `not_affected` status is the only
  status that results in CVE suppression, where the CVE is filtered out of the
  analysis results. The other statuses are useful for documentation purposes,
  but they do not work for creating exceptions. For more information about all
  the possible status labels, see [Status Labels](https://github.com/openvex/spec/blob/main/OPENVEX-SPEC.md#status-labels)
  in the OpenVEX specification.

`--justification`
: Justifies the `not_affected` status label, informing why the product is not
  affected by the vulnerability. In this case, the justification given is
  `vulnerable_code_not_in_execute_path`, signalling that the vulnerability
  can't be executed as used by the product.

  In OpenVEX, status justifications can have one of the five possible values:

  - `component_not_present`
  - `vulnerable_code_not_present`
  - `vulnerable_code_not_in_execute_path`
  - `vulnerable_code_cannot_be_controlled_by_adversary`
  - `inline_mitigations_already_exist`

  For more information about these values and their definitions, see
  [Status Justifications](https://github.com/openvex/spec/blob/main/OPENVEX-SPEC.md#status-justifications)
  in the OpenVEX specification.

`--file`
: Filename of the VEX document output

Here's the OpenVEX JSON generated by this command:

```json
{
  "@context": "https://openvex.dev/ns/v0.2.0",
  "@id": "https://openvex.dev/docs/public/vex-749f79b50f5f2f0f07747c2de9f1239b37c2bda663579f87a35e5f0fdfc13de5",
  "author": "author@example.com",
  "timestamp": "2024-05-27T13:20:22.395824+02:00",
  "version": 1,
  "statements": [
    {
      "vulnerability": {
        "name": "CVE-2022-24999"
      },
      "timestamp": "2024-05-27T13:20:22.395829+02:00",
      "products": [
        {
          "@id": "pkg:docker/example/app@v1",
          "subcomponents": [
            {
              "@id": "pkg:npm/express@4.17.1"
            }
          ]
        }
      ],
      "status": "not_affected",
      "justification": "vulnerable_code_not_in_execute_path"
    }
  ]
}
```

Understanding how VEX documents are supposed to be structured can be a bit of a
mouthful. The [OpenVEX specification](https://github.com/openvex/spec)
describes the format and all the possible properties of documents and
statements. For the full details, refer to the specification to learn more
about the available fields and how to create a well-formed OpenVEX document.

To learn more about the available flags and syntax of the `vexctl` CLI tool and
how to install it, refer to the [`vexctl` GitHub repository](https://github.com/openvex/vexctl).

For an introduction to VEX, you may also want to check out this use-case guide:
[Suppress image vulnerabilities with VEX](/scout/guides/vex.md).

## Verifying VEX documents

To test whether the VEX documents you create are well-formed and produce the
expected results, use the `docker scout cves` command with the `--vex-location`
flag to apply a VEX document to a local image analysis using the CLI.

The following command invokes a local image analysis that incorporates all VEX
documents in the specified location, using the `--vex-location` flag. In this
example, the CLI is instructed to look for VEX documents in the current working
directory.

```console
$ docker scout cves <IMAGE> --vex-location .
```

The output of the `docker scout cves` command displays the results with any VEX
statements found in under the `--vex-location` location factored into the
results. For example, CVEs assigned a status of `not_affected` are filtered out
from the results. If the output doesn't seem to take the VEX statements into
account, that's an indication that the VEX documents might be invalid in some
way.

Things to look out for include:

- The PURL of a Docker image must begin with `pkg:docker/` followed by the image name.
- In a Docker image PURL, the image name and version is separated by `@`.
  An image named `example/myapp:1.0` has the following PURL: `pkg:docker/example/myapp@1.0`.
- Remember to specify an `author` (it's a mandatory field in OpenVEX)
- The [OpenVEX specification](https://github.com/openvex/spec) describes how
  and when to use `justification`, `impact_statement`, and other fields in the
  VEX documents. Specifying these in an incorrect way results in an invalid
  document. Make sure your VEX documents comply with the OpenVEX specification.

## Attach exceptions to images

When you've created an exception,
you can attach it to your image in the following ways:

- Attach the document as an [attestation](#attestation)
- Embed the document in the [image filesystem](#image-filesystem)

You can't remove a VEX document from an image once it's been added. For
documents attached as attestations, you can create a new VEX document and
attach it to the image again. Doing so will overwrite the previous VEX document
(but it won't remove the attestation). For images where the VEX document has
been embedded in the image's filesystem, you need to rebuild the image to
change the VEX document.

### Attestation

To attach exceptions as an attestation, you can use the `docker scout
attestation add` CLI command. Using attestations is the recommended option for attaching exceptions to
images.

You can attach attestations to images that have already been pushed to a
registry. You don't need to build or push the image again. Additionally, having
the exceptions attached to the image as attestations means consumers can
inspect the exceptions for an image, directly from the registry.

To attach an attestation to an image:

1. Build the image and push it to a registry.

   ```console
   $ docker build --provenance=true --sbom=true --tag <IMAGE> --push .
   ```

2. Attach the exception to the image as an attestation.

   ```console
   $ docker scout attestation add \
     --file <cve-id>.vex.json \
     --predicate-type https://openvex.dev/ns/v0.2.0 \
     <IMAGE>
   ```

   The options for this command are:

   - `--file`: the location and filename of the VEX document
   - `--predicate-type`: the in-toto `predicateType` for OpenVEX

### Image filesystem

Embedding exceptions directly on the image filesystem is a good option if you
know the exceptions ahead of time, before you build the image. And it's easy;
just `COPY` the VEX document to the image in your Dockerfile.

The downside with this approach is that you can't change or update the
exception later. Image layers are immutable, so anything you put in the image's
filesystem is there forever. Attaching the document as an
[attestation](#attestation) provides better flexibility.

> **Note**
>
> VEX documents embedded in the image filesystem are not considered for images
> that have attestations. If your image has **any** attestations, Docker Scout
> will only look for exceptions in the attestations, and not in the image
> filesystem.
>
> If you want to use the VEX document embedded in the image filesystem, you
> must remove the attestation from the image. Note that provenance attestations
> may be added automatically for images. To ensure that no attestations are
> added to the image, you can explicitly disable both SBOM and provenance
> attestations using the `--provenance=false` and `--sbom=false` flags when
> building the image.

To embed a VEX document on the image filesystem, `COPY` the file into the image
as part of the image build. The following example shows how to copy all VEX
documents under `.vex/` in the build context, to `/var/lib/db` in the image.

```dockerfile
# syntax=docker/dockerfile:1

FROM alpine
COPY .vex/* /var/lib/db/
```

The filename of the VEX document must match the `*.vex.json` glob pattern.
It doesn't matter where on the image's filesystem you store the file.

Note that the copied files must be part of the filesystem of the final image,
For multi-stage builds, the documents must persist in the final stage.

## View exceptions

The **Exceptions** page on the [Docker Scout Dashboard](https://scout.docker.com/)
lists the exceptions for all images in your organization. Selecting a row in
the list opens the exception side panel, which displays more information about
the exception and where it comes from.

To view all exceptions for a specific image tag:

{{< tabs >}}
{{< tab name="Docker Scout Dashboard" >}}

1. Go to the [Images page](https://scout.docker.com/reports/images).
2. Select the tag that you want to inspect.
3. Open the **Image attestations** tab.

{{< /tab >}}
{{< tab name="Docker Desktop" >}}

1. Open the **Images** view in Docker Desktop.
2. Open the **Hub** tab.
3. Select the tag you want to inspect.
4. Open the **Image attestations** tab.

{{< /tab >}}
{{< /tabs >}}
