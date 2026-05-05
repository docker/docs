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

For how VEX affects vulnerability counts and scanner selection, see [Scanner
integrations](/manuals/dhi/explore/scanner-integrations.md). To scan a DHI with
VEX support, see [Scan Docker Hardened Images](/manuals/dhi/how-to/scan.md).

## VEX status reference

Each VEX statement includes a `status` field that records Docker's
exploitability assessment for a given CVE and image. DHI uses three of the four
OpenVEX status values.

| Status | Meaning |
|---|---|
| `not_affected` | The CVE was reported against a package in the image, but Docker has assessed it is not exploitable as shipped |
| `under_investigation` | Docker is aware of the CVE and is actively evaluating whether it affects the image |
| `affected` | Docker has confirmed the CVE is exploitable in the image and a fix is not yet available |

You can view the VEX statements for any DHI using Docker Scout. See [Scan Docker
Hardened Images](/manuals/dhi/how-to/scan.md).

### `not_affected` justification codes

`not_affected` statements include a machine-readable `justification` field
explaining why the vulnerability does not apply:

| Justification | Meaning |
|---|---|
| `component_not_present` | The vulnerable component is not present in this image; the CVE matched by name against a different package |
| `vulnerable_code_not_present` | The vulnerable code path was not compiled into this build |
| `vulnerable_code_not_in_execute_path` | The vulnerable code exists in the package but is not called in this image's runtime configuration |
| `vulnerable_code_cannot_be_controlled_by_adversary` | The vulnerable code exists but an attacker cannot trigger it in this configuration |
| `inline_mitigations_already_exist` | Docker has applied a backport or patch that addresses the CVE |

### Why DHI does not use `fixed`

DHI does not use `fixed`. VEX-enabled scanners may not handle `fixed`
consistently, so when Docker backports an upstream patch where the version
number alone would not reflect the fix, it uses `not_affected` with
`inline_mitigations_already_exist` justification instead.
