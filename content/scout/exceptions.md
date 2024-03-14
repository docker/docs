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
false positives in image analysis using VEX documents.

By negating non-applicable vulnerabilities, you can make it easier for yourself
and downstream consumers of your images to understand the security implications
of a vulnerability in the context an image.

In Docker Scout, exceptions are automatically factored into the results,
provided that you trust the authority that issued the exception. If an image
contains an exception that flags a CVE as non-applicable, then that CVE is
excluded from analysis results.

## Create an exception

To add an exception to an image, you need a Vulnerability Exploitability
eXchange (VEX) document. VEX is a standard format for documenting
vulnerabilites in the context of a software package or product.

There are multiple implementations and formats of VEX. Docker Scout supports
the [OpenVex](https://github.com/openvex/spec) implementation.

To create an OpenVEX document, you can use the `vexctl` command line tool.
The inputs you provide to `vexctl` when creating a VEX document for an image
typically includes:

- Author (an email address)
- Product (the PURL of the container image)
- CVE ID
- Subcomponents (one or more PURLs of packages affected by the CVE)
- The status of the vulnerability
- A status justification, impact statement, and/or an action statement,
  depending on the assigned vulnerability status
- The filename that you want to output the VEX document to

The following example command creates a VEX document `CVE-2022-24999.vex.json`.
The VEX document states that the image `example/app:v1` contains an npm package
`express@4.17.1` that is unaffected by CVE with ID `CVE-2022-24999`. The
justification for the `not_affected` status is that the vulnerable code is
never executed in containers that run this image.

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

The [OpenVEX specification](https://github.com/openvex/spec) dictates the
format and structure of OpenVEX documents and statements. Refer to the
specification to learn more about the available fields and how to create a
well-formed OpenVEX document.

To learn more about the available flags and syntax of the `vexctl` CLI tool and how to install it,
refer to the [`vexctl` GitHub repository](https://github.com/openvex/vexctl).

For an introduction to VEX, you may also want to check out this use-case guide:
[Suppress image vulnerabilities with VEX](./guides/vex.md).

## Attach exceptions to images

When you've created an exception,
you can attach it to your image in the following ways:

- Embed the document in the [image filesystem](#image-filesystem)
- Attach the document as an [attestation](#attestation)

### Image filesystem

Embedding exceptions directly on the image filesystem ensures that they follow
the image wherever it goes. This is a good option if you know the exceptions
ahead of time, and you want to ensure that the exceptions are embedded into the
image directly at build-time.

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

### Attestation

To embed exceptions as an attestation, you can use the `docker scout
attestation add` CLI command.

Just like when you embed the exception on the image's filesystem, exceptions in
attestations also follow the image. This is a good option for image publishers
that want to avoid bundling the exceptions to the image's filesystem. It also
lets you attach the exception later, after the image was already built and
pushed to a registry.

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

## View exceptions

The **Exceptions** page on the [Docker Scout Dashboard](https://scout.docker.com/)
lists the exceptions for all images in your organization. Selecting a row in
the list opens the exception side panel, which displays more information about
the exception and where it comes from.

To view all exceptions for a specific image tag:

{{< tabs >}}
{{< tab name="Dashboard" >}}

1. Open the [Docker Scout Dashboard](https://scout.docker.com/).
2. Go to the **Images** page.
3. Select the tag that you want to inspect.
4. Open to the **Image attachments** tab.

{{< /tab >}}
{{< tab name="Docker Desktop" >}}

1. Open the **Images** view in Docker Desktop.
2. Open the **Hub** tab.
3. Select the tag you want to inspect.
4. Open the **Image attachments** tab.

{{< /tab >}}
{{< /tabs >}}
