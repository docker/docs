---
description: More details on the advisory database and CVE-to-package matching service
  behind Docker Scout analysis.
keywords: scout, scanning, analysis, vulnerabilities, Hub, supply chain, security, packages, repositories, ecosystem
title: Advisory database sources and matching service
aliases:
  /scout/advisory-db-sources/
---

Reliable information sources are key for Docker Scout's ability to
surface relevant and accurate assessments of your software artifacts.
Given the diversity of sources and methodologies in the industry,
discrepancies in vulnerability assessment results can and do happen.
This page describes how the Docker Scout advisory database
and its CVE-to-package matching approach works to deal with these discrepancies.

## Advisory database sources

Docker Scout aggregates vulnerability data from multiple sources.
The data is continuously updated to ensure that your security posture
is represented using the latest available information, in real-time.

Docker Scout uses the following package repositories and security trackers:

- [Alpine secdb](https://secdb.alpinelinux.org/)
- [AlmaLinux Security Advisory](https://errata.almalinux.org/)
- [Amazon Linux Security Center](https://alas.aws.amazon.com/)
- [Bitnami Vulnerability Database](https://github.com/bitnami/vulndb)
- [CISA Known Exploited Vulnerability
  Catalog](https://www.cisa.gov/known-exploited-vulnerabilities-catalog)
- [CISA Vulnrichment](https://github.com/cisagov/vulnrichment)
- [Debian Security Bug Tracker](https://security-tracker.debian.org/tracker/)
- [Exploit Prediction Scoring System (EPSS)](https://api.first.org/epss/)
- [GitHub Advisory Database](https://github.com/advisories/)
- [GitLab Advisory
  Database](https://gitlab.com/gitlab-org/advisories-community/)
- [Golang VulnDB](https://github.com/golang/vulndb)
- [inTheWild, a community-driven open database of vulnerability
  exploitation](https://github.com/gmatuz/inthewilddb)
- [National Vulnerability Database](https://nvd.nist.gov/)
- [Oracle Linux Security](https://linux.oracle.com/security/)
- [Python Packaging Advisory
  Database](https://github.com/pypa/advisory-database)
- [RedHat Security Data](https://www.redhat.com/security/data/metrics/)
- [Rocky Linux Security Advisory](https://errata.rockylinux.org/)
- [RustSec Advisory Database](https://github.com/rustsec/advisory-db)
- [SUSE Security CVRF](http://ftp.suse.com/pub/projects/security/cvrf/)
- [Ubuntu CVE Tracker](https://people.canonical.com/~ubuntu-security/cve/)
- [Wolfi Security Feed](https://packages.wolfi.dev/os/security.json)
- [Chainguard Security Feed](https://packages.cgr.dev/chainguard/security.json)

When you enable Docker Scout for your Docker organization,
a new database instance is provisioned on the Docker Scout platform.
The database stores the Software Bill of Materials (SBOM) and other metadata about your images.
When a security advisory has new information about a vulnerability,
your SBOM is cross-referenced with the CVE information to detect how it affects you.

For more details on how image analysis works, see the [image analysis page](/scout/explore/analysis.md).

## Vulnerability matching

Traditional tools often rely on broad [Common Product Enumeration (CPE)](https://en.wikipedia.org/wiki/Common_Platform_Enumeration) matching,
which can lead to many false-positive results.

Docker Scout uses [Package URLs (PURLs)](https://github.com/package-url/purl-spec)
to match packages against CVEs, which yields more precise identification of vulnerabilities.
PURLs significantly reduce the chances of false positives, focusing only on genuinely affected packages.

## Supported package ecosystems

Docker Scout supports the following package ecosystems:

- .NET
- GitHub packages
- Go
- Java
- JavaScript
- PHP
- Python
- RPM
- Ruby
- `alpm` (Arch Linux)
- `apk` (Alpine Linux)
- `deb` (Debian Linux and derivatives)
