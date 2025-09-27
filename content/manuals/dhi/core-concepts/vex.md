---
title: Vulnerability Exploitability eXchange (VEX)
linktitle: VEX
description: Learn how VEX helps you prioritize real risks by identifying which vulnerabilities in Docker Hardened Images are actually exploitable.
keywords: vex container security, vulnerability exploitability, filter false positives, docker scout vex, cve prioritization
---

## What is VEX?

Vulnerability Exploitability eXchange (VEX) is a standardized framework
developed by the U.S. Cybersecurity and Infrastructure Security Agency (CISA) to
document the exploitability of vulnerabilities within software components.
Unlike traditional CVE (Common Vulnerabilities and Exposures) databases, VEX
provides contextual assessments, indicating whether a vulnerability is
exploitable in a given environment. This approach helps organizations prioritize
remediation efforts by distinguishing between vulnerabilities that are
exploitable and those that are not relevant to their specific use cases.

## Why is VEX important?

VEX enhances traditional vulnerability management by:

- Reducing false positives: By providing context-specific assessments, VEX helps
  in filtering out vulnerabilities that do not pose a threat in a particular
  environment.

- Prioritizing remediation: Organizations can focus resources on addressing
  vulnerabilities that are exploitable in their specific context, improving
  efficiency in vulnerability management.

- Enhancing compliance: VEX reports provide detailed information that can assist
  in meeting regulatory requirements and internal security standards.

This approach is particularly beneficial in complex environments where numerous
components and configurations exist, and traditional CVE-based assessments may
lead to unnecessary remediation efforts.

## How Docker Hardened Images integrate VEX

To enhance vulnerability management, Docker Hardened Images (DHI) incorporate
VEX reports, providing context-specific assessments of known vulnerabilities.

This integration allows you to:

- Assess exploitability: Determine whether known vulnerabilities in the image's
components are exploitable in their specific environment.

- Prioritize actions: Focus remediation efforts on vulnerabilities that pose
  actual risks, optimizing resource allocation.

- Streamline audits: Utilize the detailed information provided by VEX reports to
  simplify compliance audits and reporting.

By combining the security features of DHI with the contextual insights of VEX,
organizations can achieve a more effective and efficient approach to
vulnerability management.

## Use VEX to filter known non-exploitable CVEs

When using Docker Scout, VEX statements are automatically applied and no
manual configuration is needed.

To manually retrieve the VEX attestation for tools that support it:

```console
$ docker scout vex get <your-namespace>/dhi-<image>:<tag> --output vex.json
```

> [!NOTE]
>
> The `docker scout vex get` command requires [Docker Scout
> CLI](https://github.com/docker/scout-cli/) version 1.18.3 or later.

For example:

```console
$ docker scout vex get docs/dhi-python:3.13 --output vex.json
```

This creates a `vex.json` file containing the VEX statements for the specified
image. You can then use this file with tools that support VEX to filter out
known non-exploitable CVEs.

For example, with Grype and Trivy, you can use the `--vex` flag to apply the VEX
statements during the scan:

```console
$ grype <your-namespace>/dhi-<image>:<tag> --vex vex.json
```