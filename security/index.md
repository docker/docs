---
description: Docker security announcements
keywords: Docker, CVEs, security, notice, Log4J 2, Log4Shell, announcements
title: Docker security announcements
toc_min: 1
toc_max: 2
---

## Log4j 2 CVE-2021-44228

The [Log4j 2 CVE-2021-44228](https://nvd.nist.gov/vuln/detail/CVE-2021-44228){:
target="_blank" rel="noopener" class="_"} vulnerability in Log4j 2, a very common Java logging library, allows remote code execution, often from a context that is easily available to an attacker. For example, it was found in Minecraft servers which allowed the commands to be typed into chat logs as these were then sent to the logger. This makes it a very serious vulnerability, as the logging library is used so widely and it may be simple to exploit. Many open source maintainers are working hard with fixes and updates to the software ecosystem.

The vulnerable versions of Log4j 2 are versions 2.0 to version 2.14.1 inclusive. The first fixed version is 2.15.0. We strongly encourage you to update to the [latest version](https://logging.apache.org/log4j/2.x/download.html) if you can. If you are using a version before 2.0, you are also not vulnerable.

You may not be vulnerable if you are using these versions, as your configuration
may already mitigate this (see the Mitigations section below), or the things you
log may not include any user input. This may be difficult to validate however
without understanding all the code paths that may log in detail, and where they
may get input from. So you probably will want to upgrade all code using
vulnerable versions.

> CVE-2021-45046
>
> As an update to
> [CVE-2021-44228](https://nvd.nist.gov/vuln/detail/CVE-2021-44228){:
target="_blank" rel="noopener" class="_"}, the fix made in version 2.15.0 was
> incomplete in certain non-default configurations. An additional issue has
> been identified and is tracked with
> [CVE-2021-45046](https://nvd.nist.gov/vuln/detail/CVE-2021-45046){: target="_blank" rel="noopener" class="_"}. For a more complete fix to this
> vulnerability, we recommended that you update to 2.16.0 where possible.
{: .important}

## Scan images using the `docker scan` command

The configuration for the `docker scan` command previously shipped in Docker
Desktop versions 4.3.0 and earlier unfortunately do not detect this
vulnerability on scans. You must update your Docker Desktop installation to
4.3.1 or higher to fix this issue. For detailed instructions, see [Scan images for Log4j2 CVE](../engine/scan/index.md#scan-images-for-log4j-2-cve).

## Scan images on Docker Hub

Docker Hub security scans triggered **after 1700 UTC 13 December 2021** are now
correctly identifying the Log4j2 CVE. Scans before this date **do not**
currently reflect the status of this vulnerability. Therefore, we recommend that
you trigger scans by pushing the image to Docker Hub to view the status of
Log4j2 CVE in the vulnerability report. For detailed instructions, see [Scan images on Docker Hub](../docker-hub/vulnerability-scanning.md).

## Docker Official Images impacted by Log4j 2 CVE

> **Important**
>
> We will be updating this section with the latest information. We recommend
> that you revisit this section to view the list of affected images and update
> images to the patched version as soon as possible to remediate the issue.
{: .important}

A number of [Docker Official Images](../docker-hub/official_images.md) contain the vulnerable versions of
Log4j 2 CVE-2021-44228. The following table lists Docker Official Images that
may contain the vulnerable versions of Log4j 2. We are working on updating
Log4j 2 in these images to the latest version. Some of these images may not be
vulnerable for other reasons. We recommend that you also review the guidelines published on the upstream websites.

| Repository                | Patched version         | Additional documentation       |
|:------------------------|:-----------------------|:-----------------------|
| [couchbase](https://hub.docker.com/_/couchbase)    | 7.0.3 | [Couchbase blog](https://blog.couchbase.com/what-to-know-about-the-log4j-vulnerability-cve-2021-44228/) |
| [Elasticsearch](https://hub.docker.com/_/elasticsearch)    | 7.16.1, 6.8.21 | [Elasticsearch security announcement](https://discuss.elastic.co/t/apache-log4j2-remote-code-execution-rce-vulnerability-cve-2021-44228-esa-2021-31/291476) |
| [Flink](https://hub.docker.com/_/flink)    | Awaiting info  | [Flink advice on Log4j CVE](https://flink.apache.org/2021/12/10/log4j-cve.html) |
| [Geonetwork](https://hub.docker.com/_/geonetwork)    | 3.10.10 | [Geonetwork GitHub discussion](https://github.com/geonetwork/core-geonetwork/issues/6076) |
| [lightstreamer](https://hub.docker.com/_/lightstreamer)     | Awaiting info | Awaiting info  |
| [logstash](https://hub.docker.com/_/logstash)    | 7.16.1 | Awaiting info |
| [neo4j](https://hub.docker.com/_/neo4j)     | 4.4.1 | [Neo4j announcement](https://community.neo4j.com/t/log4j-cve-mitigation-for-neo4j/48856) |
| [nuxeo](https://hub.docker.com/_/nuxeo)     | Awaiting info | [nuxeo blog](https://www.nuxeo.com/blog/apache-log4j/) |
| [solr](https://hub.docker.com/_/solr)    | Awaiting info | [Solr security news](https://solr.apache.org/security.html#apache-solr-affected-by-apache-log4j-cve-2021-44228) |
| [sonarqube](https://hub.docker.com/_/sonarqube)    | 9.2.2 | [SonarQube announcement](https://community.sonarsource.com/t/sonarqube-sonarcloud-and-the-log4j-vulnerability/54721) |
| [storm](https://hub.docker.com/_/storm)    | Awaiting info | Awaiting info |

> **Note**
>
> Although [xwiki](https://hub.docker.com/_/xwiki){:
target="_blank" rel="noopener" class="_"} images may be detected as vulnerable
by some scanners, the authors believe the images are not vulnerable by Log4j 2
CVE as the API jars do not contain the vulnerability.
