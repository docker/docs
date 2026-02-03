---
title: Vulnerability Exploitability eXchange (VEX)
linktitle: VEX
description: Learn how VEX helps you prioritize real risks by identifying which vulnerabilities in Docker Hardened Images are actually exploitable.
keywords: vex container security, vulnerability exploitability, filter false positives, docker scout vex, cve prioritization
---

## What is VEX?

Vulnerability Exploitability eXchange (VEX) is a specification for documenting
the exploitability status of vulnerabilities within software components. VEX is
primarily defined through industry standards such as CSAF (OASIS) and CycloneDX
VEX, with the U.S. Cybersecurity and Infrastructure Security Agency (CISA)
encouraging its adoption. VEX complements CVE (Common Vulnerabilities and
Exposures) identifiers by adding producer-asserted status information,
indicating whether a vulnerability is exploitable in the product as shipped.
This helps organizations prioritize remediation efforts by identifying
vulnerabilities that do not affect their specific product configurations.

## Why is VEX important?

VEX enhances traditional vulnerability management by:

- Suppressing non-applicable vulnerabilities: By providing product-level
  exploitability assertions from the supplier, VEX helps filter out
  vulnerabilities that do not affect the product as shipped.

- Prioritizing remediation: Organizations can focus resources on addressing
  vulnerabilities that the producer has confirmed are exploitable in the
  product, improving efficiency in vulnerability management.

- Supporting vulnerability documentation: VEX statements can support audit
  discussions and help document why certain vulnerabilities do not require
  remediation.

This approach is particularly beneficial when working with complex software
components where not all reported CVEs apply to the specific product
configuration.

## How Docker Hardened Images integrate VEX

To enhance vulnerability management, Docker Hardened Images (DHI) incorporate
VEX reports, providing context-specific assessments of known vulnerabilities.

This integration allows you to:

- Consume producer assertions: Review Docker's assertions about whether known
  vulnerabilities in the image's components are exploitable in the product as
  shipped.

- Prioritize actions: Focus remediation efforts on vulnerabilities that Docker
  has confirmed are exploitable in the image, optimizing resource allocation.

- Support audit documentation: Use VEX statements to document why certain
  reported vulnerabilities do not require immediate action.

By combining the security features of DHI with VEX's product-level
exploitability assertions, organizations can achieve a more effective and
efficient approach to vulnerability management.

> [!TIP]
>
> To understand which scanners support VEX and why it matters for your security
> workflow, see [Scanner integrations](/manuals/dhi/explore/scanner-integrations.md).

## Use VEX to suppress non-applicable CVEs

Docker Hardened Images include VEX attestations that can be consumed by
vulnerability scanners to suppress non-applicable CVEs. For detailed
instructions on scanning with VEX support across different tools including
Docker Scout, Trivy, and Grype, see [Scan Docker Hardened
Images](/manuals/dhi/how-to/scan.md).