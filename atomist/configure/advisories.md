---
title: Advisory sources
description:
keywords:
---

{% include atomist/disclaimer.md %}

Vulnerabilities and advisories in the Atomist vulnerability database are updated
continuously. Data from advisories is ingested, processed, enriched, and
transformed from the following sources, by default:

|           | Name                               | URL                                                   |
| :-------: | ---------------------------------- | ----------------------------------------------------- |
| `alpine`  | Alpine `secdb`                     | <https://secdb.alpinelinux.org/?>                     |
| `amazon`  | Amazon Linux Security Center       | <https://alas.aws.amazon.com/>                        |
| `debian`  | Debian Security Bug Tracker        | <https://security-tracker.debian.org/tracker/>        |
|  `nist`   | National Vulnerability Database    | <https://nvd.nist.gov/>                               |
| `redhat`  | RedHat Security Data               | <https://www.redhat.com/security/data/metrics/>       |
| `ubuntu`  | Ubuntu CVE Tracker                 | <https://people.canonical.com/~ubuntu-security/cve/>  |
|  `suse`   | SUSE Security CVRF                 | <http://ftp.suse.com/pub/projects/security/cvrf/>     |
| `github`  | GitHub Advisory Database           | <https://github.com/advisories/>                      |
| `gitlab`  | GitLab Advisory Database           | <https://gitlab.com/gitlab-org/advisories-community/> |
| `golang`  | Golang VulnDB                      | <https://github.com/golang/vulnbd>                    |
| `rustsec` | RustSec Advisory Database          | <https://github.com/rustsec/advisory-db>              |
|  `pypa`   | Python Packaging Advisory Database | <https://github.com/pypa/advisory-database>           |

The ingestion process for vulnerabilities and advisories operates around the
clock and on incremental diffs detected in the source systems making updates
immediately available to users of the database. Our ingestion architecture is
modular, allowing us to incorporate new advisory streams within a consistent
model. We see the ingestion process and our database as significant competitive
differentiators.

During ingestion, source advisory records are transformed into a common,
proprietary data format. This transformation includes URLs, references to CVEs,
scoring data and most importantly vulnerable version ranges and fix versions.

To support matching vulnerabilities and advisories to software packages -
commonly identified by purls - we create our own internal `advisory-url`. The
`advisory-url` can be created using package information from vendor advisories
like Debian, Alpine and open source vulnerability databases like GitHub and
GitLab. With the help of the `advisory-url`, we can identify potentially
affected packages with a very fast index lookup.

Once packages are identified via the `advisory-url` search, the database uses a
proprietary algorithm to check if a package version satisfies the vulnerable
version range of an advisory.

Additionally we normalize version ranges as provided in vendor advisories. This
is important as we deal with a variety of versioning systems, e.g.
[SemVer for NPM](https://github.com/npm/node-semver#ranges),
[Maven versioning](https://maven.apache.org/pom.html#Dependency_Version_Requirement_Specification),
or
[Debian](https://www.debian.org/doc/debian-policy/ch-controlfields.html#version)
package versioning for Debian and Ubuntu based images.

## Adding and updating advisories

Users can create their own vulnerabilities and advisories by following the steps
outlined below:

- Create a repository called `atomist-advisories` in an GitHub organzation/user
  that has the Atomist GitHub App installed making sure that the app has read
  access to the repository.
- In the default branch of `atomist-advisories` repository, add a new JSON file
  called `<source>/<source id>.json` where `source` should be the name of your
  company and `source-id` has to be a unique id for the advisory within
  `source`.
- The content of the JSON file must strictly follow the schema defined in
  [Open Source Vulnerability format](https://ossf.github.io/osv-schema/). Take a
  look at the
  [GitHub Advisory Database](https://github.com/github/advisory-database/tree/main/advisories/github-reviewed)
  for examples of advisories.

### Deleting advisories

Deleting an advisory from the database can be achieved by removing the
corresponding JSON advisory file from the `atomist-advisories` repository.

> Only additions, changes and removals of JSON advisory files in the
> repository's default branch are being processed and mirrored into the
> database.
