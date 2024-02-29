---
description: Docker security announcements
keywords: Docker, CVEs, security, notice, Log4J 2, Log4Shell, Text4Shell, announcements
title: Docker security announcements
toc_min: 1
toc_max: 2
---

## Docker Security Advisory: Multiple Vulnerabilities in runc, BuildKit, and Moby

_Last updated February 2, 2024_

We at Docker prioritize the security and integrity of our software and the trust of our users. Security researchers at Snyk Labs identified and reported four security vulnerabilities in the container ecosystem. One of the vulnerabilities, [CVE-2024-21626](https://scout.docker.com/v/CVE-2024-21626), concerns the runc container runtime, and the other three affect BuildKit ([CVE-2024-23651](https://scout.docker.com/v/CVE-2024-23651), [CVE-2024-23652](https://scout.docker.com/v/CVE-2024-23652), and [CVE-2024-23653](https://scout.docker.com/v/CVE-2024-23653)). We want to assure our community that our team, in collaboration with the reporters and open source maintainers, has been diligently working on coordinating and implementing necessary remediations.

We are committed to maintaining the highest security standards. We have published patched versions of runc, BuildKit, and Moby on January 31 and released an update for Docker Desktop on February 1 to address these vulnerabilities.  Additionally, our latest BuildKit and Moby releases included fixes for [CVE-2024-23650](https://scout.docker.com/v/CVE-2024-23650) and [CVE-2024-24557](https://scout.docker.com/v/CVE-2024-24557), discovered respectively by an independent researcher and through Docker’s internal research initiatives.

|                        | Versions Impacted         |
|:-----------------------|:--------------------------|
| `runc`                 | <= 1.1.11                 |
| `BuildKit`             | <= 0.12.4                 |
| `Moby (Docker Engine)` | <= 25.0.1 and <= 24.0.8   |
| `Docker Desktop`       | <= 4.27.0                 |

### What should I do if I’m on an affected version?

If you are using affected versions of runc, BuildKit, Moby, or Docker Desktop, make sure to update to the latest versions, linked in the following table:

|                        | Patched Versions          |
|:-----------------------|:--------------------------|
| `runc`                 | >= [1.1.12](https://github.com/opencontainers/runc/releases/tag/v1.1.12)                 |
| `BuildKit`             | >= [0.12.5](https://github.com/moby/buildkit/releases/tag/v0.12.5)                 |
| `Moby (Docker Engine)` | >= [25.0.2](https://github.com/moby/moby/releases/tag/v25.0.2) and >= [24.0.9](https://github.com/moby/moby/releases/tag/v24.0.9)   |
| `Docker Desktop`       | >= [4.27.1](../desktop/release-notes.md#4271)                 |


If you are unable to update to an unaffected version promptly, follow these best practices to mitigate risk: 

* Only use trusted Docker images (such as [Docker Official Images](../trusted-content/official-images.md)).
* Don’t build Docker images from untrusted sources or untrusted Dockerfiles.
* If you are a Docker Business customer using Docker Desktop and unable to update to v4.27.1, make sure to enable [Hardened Docker Desktop](../desktop/hardened-desktop/_index.md) features such as:
  * [Enhanced Container Isolation](../desktop/hardened-desktop/enhanced-container-isolation/_index.md), which mitigates the impact of CVE-2024-21626 in the case of running containers from malicious images.
  * [Image Access Management](./for-admins/image-access-management.md), and [Registry Access Management](./for-admins/registry-access-management.md), which give organizations control over which images and repositories their users can access.
* For CVE-2024-23650, CVE-2024-23651, CVE-2024-23652, and CVE-2024-23653, avoid using BuildKit frontend from an untrusted source. A frontend image is usually specified as the #syntax line on your Dockerfile, or with `--frontend` flag when using the `buildctl build` command.
* To mitigate CVE-2024-24557, make sure to either use BuildKit or disable caching when building images. From the CLI this can be done via the `DOCKER_BUILDKIT=1` environment variable (default for Moby >= v23.0 if the buildx plugin is installed) or the `--no-cache flag`. If you are using the HTTP API directly or through a client, the same can be done by setting `nocache` to `true` or `version` to `2` for the [/build API endpoint](https://docs.docker.com/engine/api/v1.44/#tag/Image/operation/ImageBuild).

### Technical details and impact

#### CVE-2024-21626 (High)
In runc v1.1.11 and earlier, due to certain leaked file descriptors, an attacker can gain access to the host filesystem by causing a newly-spawned container process (from `runc exec`) to have a working directory in the host filesystem namespace, or by tricking a user to run a malicious image and allow a container process to gain access to the host filesystem through `runc run`. The attacks can also be adapted to overwrite semi-arbitrary host binaries, allowing for complete container escapes. Note that when using higher-level runtimes (such as Docker or Kubernetes), this vulnerability can be exploited by running a malicious container image without additional configuration or by passing specific workdir options when starting a container. The vulnerability can also be exploited from within Dockerfiles in the case of Docker.

_The issue has been fixed in runc v1.1.12._

#### CVE-2024-23651 (High)
In BuildKit <= v0.12.4, two malicious build steps running in parallel sharing the same cache mounts with subpaths could cause a race condition, leading to files from the host system being accessible to the build container. This will only occur if a user is trying to build a Dockerfile of a malicious project.

_The issue has been fixed in BuildKit v0.12.5._

#### CVE-2024-23652 (High)
In BuildKit <= v0.12.4, a malicious BuildKit frontend or Dockerfile using `RUN --mount` could trick the feature that removes empty files created for the mountpoints into removing a file outside the container from the host system. This will only occur if a user is using a malicious Dockerfile.

_The issue has been fixed in BuildKit v0.12.5._

#### CVE-2024-23653 (High)
In addition to running containers as build steps, BuildKit also provides APIs for running interactive containers based on built images. In BuildKit <= v0.12.4, it is possible to use these APIs to ask BuildKit to run a container with elevated privileges. Normally, running such containers is only allowed if special `security.insecure` entitlement is enabled both by buildkitd configuration and allowed by the user initializing the build request.

_The issue has been fixed in BuildKit v0.12.5._

#### CVE-2024-23650 (Medium)
In BuildKit <= v0.12.4, a malicious BuildKit client or frontend could craft a request that could lead to BuildKit daemon crashing with a panic.

_The issue has been fixed in BuildKit v0.12.5._

#### CVE-2024-24557 (Medium)
In Moby <= v25.0.1 and <= v24.0.8, the classic builder cache system is prone to cache poisoning if the image is built FROM scratch. Also, changes to some instructions (most important being `HEALTHCHECK` and `ONBUILD`) would not cause a cache miss. An attacker with knowledge of the Dockerfile someone is using could poison their cache by making them pull a specially crafted image that would be considered a valid cache candidate for some build steps.

_The issue has been fixed in Moby >= v25.0.2 and >= v24.0.9._

### How are Docker products affected? 

#### Docker Desktop
Docker Desktop v4.27.0 and earlier are affected. Docker Desktop v4.27.1 was released on February 1 and includes runc, BuildKit, and dockerd binaries patches. In addition to updating to this new version, we encourage all Docker users to diligently use Docker images and Dockerfiles and ensure you only use trusted content in your builds.

As always, you should check Docker Desktop system requirements for your operating system ([Windows](../desktop/install/windows-install.md#system-requirements), [Linux](../desktop/install/linux-install.md#general-system-requirements), [Mac](../desktop/install/mac-install.md#system-requirements)) before updating to ensure full compatibility.

#### Docker Build Cloud
Any new Docker Build Cloud builder instances will be provisioned with the latest Docker Engine and BuildKit versions and will, therefore, be unaffected by these CVEs. Updates have also been rolled out to existing Docker Build Cloud builders.

_No other Docker products are affected by these vulnerabilities._

### Advisory links
* Runc
  * [CVE-2024-21626](https://github.com/opencontainers/runc/security/advisories/GHSA-xr7r-f8xq-vfvv)
* BuildKit
  * [CVE-2024-23650](https://github.com/moby/buildkit/security/advisories/GHSA-9p26-698r-w4hx)
  * [CVE-2024-23651](https://github.com/moby/buildkit/security/advisories/GHSA-m3r6-h7wv-7xxv)
  * [CVE-2024-23652](https://github.com/moby/buildkit/security/advisories/GHSA-4v98-7qmw-rqr8)
  * [CVE-2024-23653](https://github.com/moby/buildkit/security/advisories/GHSA-wr6v-9f75-vh2g)
* Moby
  * [CVE-2024-24557](https://github.com/moby/moby/security/advisories/GHSA-xw73-rw38-6vjc)

## Text4Shell CVE-2022-42889

_Last updated October 2022_

[CVE-2022-42889](https://nvd.nist.gov/vuln/detail/CVE-2022-42889) has been discovered in the popular Apache Commons Text library. Versions of this library up to but not including 1.10.0 are affected by this vulnerability.

We strongly encourage you to update to the latest version of [Apache Commons Text](https://commons.apache.org/proper/commons-text/download_text.cgi). 

### Scan images on Docker Hub

Docker Hub security scans triggered after 1200 UTC 21 October 2021 are now
correctly identifying the Text4Shell CVE. Scans before this date do not
currently reflect the status of this vulnerability. Therefore, we recommend that
you trigger scans by pushing new images to Docker Hub to view the status of
the Text4Shell CVE in the vulnerability report. For detailed instructions, see [Scan images on Docker Hub](../docker-hub/vulnerability-scanning.md).

### Docker Official Images impacted by CVE-2022-42889

A number of [Docker Official Images](../trusted-content/official-images.md) contain the vulnerable versions of
Apache Commons Text. The following lists Docker Official Images that
may contain the vulnerable versions of Apache Commons Text:

- [bonita](https://hub.docker.com/_/bonita) 
- [Couchbase](https://hub.docker.com/_/couchbase)
- [Geonetwork](https://hub.docker.com/_/geonetwork) 
- [neo4j](https://hub.docker.com/_/neo4j)
- [sliverpeas](https://hub.docker.com/_/sliverpeas)
- [solr](https://hub.docker.com/_/solr) 
- [xwiki](https://hub.docker.com/_/xwiki) 

We have updated
Apache Commons Text in these images to the latest version. Some of these images may not be
vulnerable for other reasons. We recommend that you also review the guidelines published on the upstream websites.

## Log4j 2 CVE-2021-44228

_Last updated December 2021_

The [Log4j 2 CVE-2021-44228](https://nvd.nist.gov/vuln/detail/CVE-2021-44228) vulnerability in Log4j 2, a very common Java logging library, allows remote code execution, often from a context that is easily available to an attacker. For example, it was found in Minecraft servers which allowed the commands to be typed into chat logs as these were then sent to the logger. This makes it a very serious vulnerability, as the logging library is used so widely and it may be simple to exploit. Many open source maintainers are working hard with fixes and updates to the software ecosystem.

The vulnerable versions of Log4j 2 are versions 2.0 to version 2.14.1 inclusive. The first fixed version is 2.15.0. We strongly encourage you to update to the [latest version](https://logging.apache.org/log4j/2.x/download.html) if you can. If you are using a version before 2.0, you are also not vulnerable.

You may not be vulnerable if you are using these versions, as your configuration
may already mitigate this, or the things you
log may not include any user input. This may be difficult to validate however
without understanding all the code paths that may log in detail, and where they
may get input from. So you probably will want to upgrade all code using
vulnerable versions.

> CVE-2021-45046
>
> As an update to
> [CVE-2021-44228](https://nvd.nist.gov/vuln/detail/CVE-2021-44228), the fix made in version 2.15.0 was
> incomplete. Additional issues have been identified and are tracked with
> [CVE-2021-45046](https://nvd.nist.gov/vuln/detail/CVE-2021-45046) and
> [CVE-2021-45105](https://nvd.nist.gov/vuln/detail/CVE-2021-45105).
> For a more complete fix to this vulnerability, we recommended that you update to 2.17.0 where possible.
{ .important }

### Scan images on Docker Hub

Docker Hub security scans triggered after 1700 UTC 13 December 2021 are now
correctly identifying the Log4j 2 CVEs. Scans before this date do not
currently reflect the status of this vulnerability. Therefore, we recommend that
you trigger scans by pushing new images to Docker Hub to view the status of
Log4j 2 CVE in the vulnerability report. For detailed instructions, see [Scan images on Docker Hub](../docker-hub/vulnerability-scanning.md).

## Docker Official Images impacted by Log4j 2 CVE

_Last updated December 2021_

A number of [Docker Official Images](../trusted-content/official-images.md) contain the vulnerable versions of
Log4j 2 CVE-2021-44228. The following table lists Docker Official Images that
may contained the vulnerable versions of Log4j 2. We updated Log4j 2 in these images to the latest version. Some of these images may not be
vulnerable for other reasons. We recommend that you also review the guidelines published on the upstream websites.

| Repository                | Patched version         | Additional documentation       |
|:------------------------|:-----------------------|:-----------------------|
| [couchbase](https://hub.docker.com/_/couchbase)    | 7.0.3 | [Couchbase blog](https://blog.couchbase.com/what-to-know-about-the-log4j-vulnerability-cve-2021-44228/) |
| [Elasticsearch](https://hub.docker.com/_/elasticsearch)    | 6.8.22, 7.16.2 | [Elasticsearch announcement](https://www.elastic.co/blog/new-elasticsearch-and-logstash-releases-upgrade-apache-log4j2) |
| [Flink](https://hub.docker.com/_/flink)    | 1.11.6, 1.12.7, 1.13.5, 1.14.2  | [Flink advice on Log4j CVE](https://flink.apache.org/2021/12/10/log4j-cve.html) |
| [Geonetwork](https://hub.docker.com/_/geonetwork)    | 3.10.10 | [Geonetwork GitHub discussion](https://github.com/geonetwork/core-geonetwork/issues/6076) |
| [lightstreamer](https://hub.docker.com/_/lightstreamer)     | Awaiting info | Awaiting info  |
| [logstash](https://hub.docker.com/_/logstash)    | 6.8.22, 7.16.2 | [Elasticsearch announcement](https://www.elastic.co/blog/new-elasticsearch-and-logstash-releases-upgrade-apache-log4j2) |
| [neo4j](https://hub.docker.com/_/neo4j)     | 4.4.2 | [Neo4j announcement](https://community.neo4j.com/t/log4j-cve-mitigation-for-neo4j/48856) |
| [solr](https://hub.docker.com/_/solr)    | 8.11.1 | [Solr security news](https://solr.apache.org/security.html#apache-solr-affected-by-apache-log4j-cve-2021-44228) |
| [sonarqube](https://hub.docker.com/_/sonarqube)    | 8.9.5, 9.2.2 | [SonarQube announcement](https://community.sonarsource.com/t/sonarqube-sonarcloud-and-the-log4j-vulnerability/54721) |
| [storm](https://hub.docker.com/_/storm)    | Awaiting info | Awaiting info |

> **Note**
>
> Although [xwiki](https://hub.docker.com/_/xwiki) images may be detected as vulnerable
by some scanners, the authors believe the images are not vulnerable by Log4j 2
CVE as the API jars do not contain the vulnerability.
> The [Nuxeo](https://hub.docker.com/_/nuxeo)
> image is deprecated and will not be updated.
