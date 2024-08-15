---
title: Suppress image vulnerabilities with VEX
description: An introduction to using VEX with Docker Scout
keywords: scout, vex, Vulnerability Exploitability eXchange, openvex, vulnerabilities, cves
---

Vulnerability Exploitability eXchange (VEX) is a way to add context
about how a software product is affected by a vulnerable component.
In this guide, you will learn about:

- How VEX can help you suppress non-applicable or fixed vulnerabilities found in your images
- How to create VEX documents and statements
- How to apply and consume VEX data with Docker Scout

> **Experimental**
>
> Support for VEX is an experimental feature in Docker Scout.
> We recommend that you do not use this feature in production environments
> as this feature may change or be removed from future releases.
{ .experimental }

## Prerequisites

If you want to follow along the steps of this guide, you'll need the following:

- The latest version of Docker Desktop
- The [containerd image store](../../desktop/containerd.md) must be enabled
- Git
- A [Docker account](/accounts/create-account/)
- A GitHub account

## Introduction to VEX

Just because a software product contains a vulnerable component,
it doesn't mean that the vulnerability is exploitable.
For example, the vulnerable component might never be
loaded into memory when the application runs.
Or the vulnerability might have been fixed,
but the system that detects the vulnerability is unaware.

The concept of VEX is defined by a working group by
the United States Cybersecurity and Infrastructure Security Agency (CISA).
At the core of VEX are exploitability assessments.
These assessments describe the status of a given CVE for a product.
The possible vulnerability statuses in VEX are:

- Not affected: No remediation is required regarding this vulnerability.
- Affected: Actions are recommended to remediate or address this vulnerability.
- Fixed: These product versions contain a fix for the vulnerability.
- Under investigation: It is not yet known whether these product versions are affected by the vulnerability. An update will be provided in a later release.

There are multiple implementations of VEX,
and the exact format and properties of a VEX document differ
depending on the implementation of VEX that you use.
Examples of VEX implementations include:

- OpenVEX
- Open Security Advisory Framework (OSAF)
- CycloneDX
- SPDX

This guide uses the OpenVEX implementation in its examples.
For more information about OpenVEX,
refer to the [specification](https://github.com/openvex/spec).

In all implementations, the core idea is the same:
to provide a framework for describing the impact of vulnerabilities.
Key components of VEX regardless of implementation includes:

VEX document
: A type of security advisory for storing VEX statements.
  The format of the document depends on the specific implementation.

VEX statement
: Describes the status of a vulnerability in a product,
  whether it's exploitable, and whether there are ways to remediate the issue.

Justification and impact
: Depending on the vulnerability status, statements include a justification
  or impact statement describing why a product is or isn't affected.

Action statements
: Describe how to remediate or mitigate the vulnerability.

This added context that VEX provides around vulnerabilities
helps organizations perform risk assessment of software products.

## Create and enable a repository

To get started, create a sample project to work with.
Use the [Docker Scout demo service](https://github.com/docker/scout-demo-service)
template repository to bootstrap a new repository in your own GitHub organization.

1. Create the repository from the template.
2. Clone the Git repository to your machine.
3. Build an image from the repository and push it to a new Docker Hub repository.

   ```console
   $ cd scout-demo-service
   $ docker build --provenance=true --sbom=true --tag <ORG>/scout-demo-service:v1 --push .
   ```

4. Enable Docker Scout on the repository.

   ```console
   $ docker scout repo enable <ORG>/scout-demo-service
   ```

## Inspect vulnerability

Use the `docker scout cves` command to view the vulnerabilities for the image.
For the purpose of this guide, we'll concentrate on a particular CVE: CVE-2022-24999.

```console
$ docker scout cves --only-cve-id CVE-2022-24999
```

The output from this command shows that this CVE affects two JavaScript packages
in the image: `express@4.17.1` and `qs@6.7.0`.

Now, let's imagine that you've reviewed this vulnerability
and determined that the exploit doesn't affect your application
because the vulnerable code is never executed.
Your image would still show as containing a `HIGH` severity vulnerability,
merely by having the component in the SBOM.

This is a situation where VEX can help suppress the CVE for your product.

## Create a VEX document

OpenVEX, the implementation of VEX used in this guide,
provides a CLI tool for generating VEX documents called `vexctl`.
If you have Go installed on your system,
you can install `vexctl` using the `go install` command:

```console
$ go install github.com/openvex/vexctl@latest
```

Otherwise, you can grab a prebuilt binary from the `vexctl`
[releases page](https://github.com/openvex/vexctl/releases) on GitHub.

Use the `vexctl create` command to create an OpenVEX document.
The following command creates a VEX document `CVE-2022-24999.vex.json`
which states that the Docker image and its sub-components
are unaffected by the CVE because the vulnerable code is never executed.

```console
$ vexctl create \
  --author="author@example.com" \
  --product="pkg:docker/<ORG>/scout-demo-service@v1" \
  --subcomponents="pkg:npm/express@4.17.1" \
  --subcomponents="pkg:npm/qs@6.7.0" \
  --vuln="CVE-2022-24999" \
  --status="not_affected" \
  --justification="vulnerable_code_not_in_execute_path" \
  --file="CVE-2022-24999.vex.json"
```

This creates a VEX document that looks like this:

```json
{
  "@context": "https://openvex.dev/ns/v0.2.0",
  "@id": "https://openvex.dev/docs/public/vex-a7e7c69be57bc49dec9cc2f3a3e3329b12674ca53b53a53ab42134dcc0510779",
  "author": "author@example.com",
  "timestamp": "2024-02-03T09:33:28.913572+01:00",
  "version": 1,
  "statements": [
    {
      "vulnerability": {
        "name": "CVE-2022-24999"
      },
      "timestamp": "2024-02-03T09:33:28.913574+01:00",
      "products": [
        {
          "@id": "pkg:docker/<ORG>/scout-demo-service@v1",
          "subcomponents": [
            {
              "@id": "pkg:npm/express@4.17.1"
            },
            {
              "@id": "pkg:npm/qs@6.7.0"
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

## Verify CVE suppression

To test whether the VEX statement provides the expected security advisory,
use the `docker scout cves` command to review the CVE status.
The `--vex-location` flag lets you specify a directory containing VEX documents.

```console
$ docker scout cves --only-cve-id CVE-2022-24999 --vex-location .
```

You should now see `VEX` fields appear in the output of the command.

```text {hl_lines=[7,8]}
✗ HIGH CVE-2022-24999 [OWASP Top Ten 2017 Category A9 - Using Components with Known Vulnerabilities]
  https://scout.docker.com/v/CVE-2022-24999
  Affected range : <4.17.3                                             
  Fixed version  : 4.17.3                                              
  CVSS Score     : 7.5                                                 
  CVSS Vector    : CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H        
  VEX            : not affected [vulnerable code not in execute path]  
                 : author@example.com     
```

## Attach VEX documents to images

To distribute VEX statements together with your images,
you can embed the VEX document in an in-toto attestation.
This is also how BuildKit attaches SBOM and provenance attestations to images.

When you encapsulate a VEX document in an in-toto attestation,
the OpenVEX specification recommends a different format for the VEX document.
In-toto attestations already contain the image reference in the attestation predicate,
making the image reference in the `products` key of the VEX document redundant.
Instead, `products` should refer to the packages that contain the vulnerabilities
(`subcomponents` in the VEX document shown earlier).

1. Create a new `in-toto.vex.json` document with the following contents:

   ```json
   {
     "@context": "https://openvex.dev/ns/v0.2.0",
     "@id": "https://openvex.dev/docs/public/vex-a7e7c69be57bc49dec9cc2f3a3e3329b12674ca53b53a53ab42134dcc0510779",
     "author": "author@example.com",
     "timestamp": "2024-02-03T09:33:28.913572+01:00",
     "version": 1,
     "statements": [
       {
         "vulnerability": {
           "name": "CVE-2022-24999"
         },
         "timestamp": "2024-02-03T09:33:28.913574+01:00",
         "products": [
           {
             "@id": "pkg:npm/express@4.17.1"
           },
           {
             "@id": "pkg:npm/qs@6.7.0"
           }
         ],
         "status": "not_affected",
         "justification": "vulnerable_code_not_in_execute_path"
       }
     ]
   }
   ```

2. Attach the `in-toto.vex.json` VEX document as an attestation:

   ```console
   $ docker scout attestation add \
     --file in-toto.vex.json \
     --predicate-type https://openvex.dev/ns/v0.2.0 \
     <ORG>/scout-demo-service:v1
   ```

   This adds an in-toto attestation to the image,
   with a predicate type of `https://openvex.dev/ns/v0.2.0`.

3. Analyze the image with `docker scout cves`.

   > [!NOTE]
   >
   > This only works when analyzing remote images in a registry.
   > To force Docker Scout to analyze a registry image instead of a local one,
   > specify the image reference with a `registry://` prefix.

   ```console
   $ docker scout cves \
     --only-cve-id CVE-2022-24999 \
     registry://<ORG>/scout-demo-service:v1
   ```

## Summary

In this guide, you've learned about:

- How VEX is a concept that can help you optimize vulnerability triaging
- How to create VEX documents and use them in image analysis
- How you can distribute VEX with your images as in-toto attestations

For more information about VEX:

- [VEX Status Justifications (CISA)](https://www.cisa.gov/sites/default/files/publications/VEX_Status_Justification_Jun22.pdf)
- [VEX Use Cases (CISA)](https://www.cisa.gov/sites/default/files/publications/VEX_Use_Cases_Document_508c.pdf)
- [OpenVEX specification](https://github.com/openvex/spec/blob/main/OPENVEX-SPEC.md)
- [Manage vulnerability exceptions](/scout/explore/exceptions.md)
