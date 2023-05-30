---
description: More details on the Advisory Database and CVE-to-package matching service behind Docker Scout analysis.
keywords: scanning, analysis, vulnerabilities, Hub, supply chain, security
title: Advisory Database sources and matching service
---

{% include scout-early-access.md %}

Docker Scout is a service that helps developers and security teams build and
maintain a secure software supply chain. A key component of this is the ability
to assess your software artifacts against a reliable source of vulnerability
information. Different tools collect vulnerability information from different
sources, and use different methods to identify matches against software
artifacts. This can lead to differing results between tools.

To help you understand why different tools can provide different results when
assessing software for vulnerabilities, this page explains how the Docker Scout
Advisory Database vulnerability database and CVE-to-package matching service
works.

## Docker Scout’s Advisory Database sources

Docker Scout creates and maintains its vulnerability database by ingesting and
collating vulnerability data from multiple sources on an hourly basis. These
sources include many recognizable package repositories and trusted security
trackers, such as:

- [Alpine secdb](https://secdb.alpinelinux.org/)
- [Amazon Linux Security Center](https://alas.aws.amazon.com/)
- [CISA Known Exploited Vulnerability
  Catalog](https://www.cisa.gov/known-exploited-vulnerabilities-catalog)
- [Debian Security Bug Tracker](https://security-tracker.debian.org/tracker/)
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
- [RustSec Advisory Database](https://github.com/rustsec/advisory-db)
- [SUSE Security CVRF](http://ftp.suse.com/pub/projects/security/cvrf/)
- [Ubuntu CVE Tracker](https://people.canonical.com/~ubuntu-security/cve/)
- [Wolfi Security Feed](https://packages.wolfi.dev/os/security.json)

Docker Scout correlates this data by making a full inventory of a container
image and storing that inventory in a [software bill of materials
(SBOM)](https://ntia.gov/sites/default/files/publications/sbom_at_a_glance_apr2021_0.pdf).

The SBOM summarizes the contents of the image and how the contents got there
meaning that when there is information about a new vulnerability, Docker Scout
correlates it with the SBOM. If Docker Scout finds a match for a vulnerability,
it can identify the artifact that’s now vulnerable, why, and where it’s in use.

When a customer enrolls with Docker Scout, the organization receives their own
instance of the database. This database tracks timestamped metadata about your
images that Docker Scout can then match to CVEs. Find more details on how this
works in the [Advanced image analysis document](./advanced-image-analysis.md).

Docker Scout is ideal for analyzing images in Docker Desktop and Docker Hub, but
the flexibility of the approach also means it can integrate with other image
sources, for example, [JFrog
Artifactory](https://docs.docker.com/scout/artifactory/).

## How Docker Scout makes more precise matches

Many other tools use fuzzy [Common Product Enumeration
(CPE)](https://en.wikipedia.org/wiki/Common_Platform_Enumeration) matching with
wild cards to known vulnerabilities with the versions of software packages they affect.
This can return a lot of false positives which you need to triage.

The typical structure of a CPE match looks like this:

```
cpe:<cpe_version>:<part>:<vendor>:<product>:<version>:<update>:<edition>:<language>:<sw_edition>:<target_sw>:<target_hw>:<other>
```

For example `cpe:*:*:*:calendar:*:*:*:*:*:*:*` returns a match on anything with
the product name “calendar”. If there is a vulnerability present in an NPM
package, this CPE match would also return packages and modules for all other
languages too.

Instead, Docker Scout matches CVEs to SBOMs using [package URL (PURL)
links](https://github.com/package-url/purl-spec) that are a more precise,
universal schema for matching software packages. A PURL link can help you only
identify the relevant packages with far less false positives.

Continuing this example, a PURL can match the specific package name to a
language and version.

```
pkg:npm/calendar@12.0.2
```

This only matches a node package with the name “calendar” and the version
“12.0.2”. For relevant packages, you can specify architectures and operating
system versions to make more precise matches.

In summary, Docker Scout’s technique improves matching accuracy and reduces the
number of results that turn out to be false-positives.

## Package ecosystems supported by the Docker Scout Advisory Database

By sourcing vulnerability data from the providers above, Docker Scout is able to support analyzing the following package ecosystems:

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
