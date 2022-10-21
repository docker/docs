---
description: Docker security announcements
keywords: Docker, CVEs, security, notice, Log4J 2, Log4Shell, Text4Shell, announcements
title: Docker security announcements
toc_min: 1
toc_max: 2
---

## Text4Shell CVE-2022-42889

[CVE-2022-42889](https://nvd.nist.gov/vuln/detail/CVE-2022-42889){:target="_blank" rel="noopener" class="_"} has been discovered in the popular Apache Commons Text library. Versions of this library up to but not including 1.10.0 are affected by this vulnerability.

We strongly encourage you to update to the latest version of [Apache Commons Text](https://commons.apache.org/proper/commons-text/download_text.cgi){:target="_blank" rel="noopener" class="_"}.

### Scan images using the `docker scan` command

`docker scan` as shipped with latest versions of Docker Desktop detects the Text4Shell CVE-2022-42889 vulnerability.

If an image is vulnerable to CVE-2022-42889, the output of `docker scan` will contain the following text:

```
  Upgrade org.apache.commons:commons-text@1.9 to org.apache.commons:commons-text@1.10.0 to fix
  ✗ Arbitrary Code Execution (new) [High Severity][https://snyk.io/vuln/SNYK-JAVA-ORGAPACHECOMMONS-3043138] in org.apache.commons:commons-text@1.9
    introduced by org.apache.commons:commons-text@1.9
```

### Scan images on Docker Hub

Docker Hub security scans triggered **after 1200 UTC 21 October 2021** are now
correctly identifying the Text4Shell CVE. Scans before this date **do not**
currently reflect the status of this vulnerability. Therefore, we recommend that
you trigger scans by pushing new images to Docker Hub to view the status of
the Text4Shell CVE in the vulnerability report. For detailed instructions, see [Scan images on Docker Hub](../docker-hub/vulnerability-scanning.md).

### Docker Official Images impacted by CVE-2022-42889

> **Important**
>
> We will be updating this section with the latest information. We recommend
> that you revisit this section to view the list of affected images and update
> images to the patched version as soon as possible to remediate the issue.
{: .important}

A number of [Docker Official Images](../docker-hub/official_images.md) contain the vulnerable versions of
Apache Commons Text. The following table lists Docker Official Images that
may contain the vulnerable versions of Apache Commons Text. We are working on updating
Apache Commons Text in these images to the latest version. Some of these images may not be
vulnerable for other reasons. We recommend that you also review the guidelines published on the upstream websites.

| Repository                | Patched version         | Additional documentation       |
|:------------------------|:-----------------------|:-----------------------|
| [bonita](https://hub.docker.com/_/bonita) |  | In progress |
| [Couchbase](https://hub.docker.com/_/couchbase) |  | In progress |
| [Geonetwork](https://hub.docker.com/_/geonetwork) |  | In progress |
| [neo4j](https://hub.docker.com/_/neo4j) |  | In progress |
| [sliverpeas](https://hub.docker.com/_/sliverpeas) |  | In progress |
| [solr](https://hub.docker.com/_/solr) |  | In progress |
| [xwiki](https://hub.docker.com/_/xwiki) |  | In progress |


## CVE-2021-45449

Docker Desktop versions 4.3.0 and 4.3.1 have a bug that may log sensitive information (access token or password) on the user's machine during login. This only affects users if they are on Docker Desktop 4.3.0, 4.3.1 and the user has logged in while on 4.3.0, 4.3.1. Gaining access to this data would require having access to the user’s local files. This vulnerability has been fixed in version 4.3.2 or higher. Users should update to version 4.3.2 and may want to update their password. Users should not send local log files to anyone. Users can manually delete their log files, they can be located in the following folder: `~/Library/Containers/com.docker.docker/Data/log/host/` on Mac, and in `C:\Users\<username>\AppData\Roaming\Docker\log\host\` on Windows. When a user installs 4.3.2 or higher, we will delete their local log files, so there is no risk of leakage after an update.

Additionally, these logs may be included when users upload diagnostics, meaning access tokens and passwords might have been shared with Docker. This only affects users if they are on Docker Desktop 4.3.0, 4.3.1, and the user has logged in while on 4.3.0, 4.3.1 and have gone through the process of submitting diagnostics to Docker. Only Docker support Engineers working on an active support case could have access to the diagnostic files, minimizing leakage risk from these files. We have deleted all potentially sensitive diagnostic files from our data storage and will continue to delete diagnostics reported from the affected versions on an ongoing basis.
For detailed information, see [CVE-2021-45449](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-45449){: target="_blank" rel="noopener" class="_"}.


### References

* [Release Notes (Windows)](../desktop/release-notes.md)
* [Release Notes (Mac)](../desktop/release-notes.md)

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
> incomplete. Additional issues have been identified and are tracked with
> [CVE-2021-45046](https://nvd.nist.gov/vuln/detail/CVE-2021-45046){: target="_blank" rel="noopener" class="_"} and
> [CVE-2021-45105](https://nvd.nist.gov/vuln/detail/CVE-2021-45105){: target="_blank" rel="noopener" class="_"}.
> For a more complete fix to this vulnerability, we recommended that you update to 2.17.0 where possible.
{: .important}

### Scan images using the `docker scan` command

The configurations for the `docker scan` command previously shipped in Docker
Desktop versions 4.3.0 and earlier unfortunately do not detect this
vulnerability on scans. You must update your Docker Desktop installation to
4.3.1 or higher to fix this issue. For detailed instructions, see [Scan images for Log4j 2 CVE](../engine/scan/index.md#scan-images-for-log4j-2-cve).

### Scan images on Docker Hub

Docker Hub security scans triggered **after 1700 UTC 13 December 2021** are now
correctly identifying the Log4j 2 CVEs. Scans before this date **do not**
currently reflect the status of this vulnerability. Therefore, we recommend that
you trigger scans by pushing new images to Docker Hub to view the status of
Log4j 2 CVE in the vulnerability report. For detailed instructions, see [Scan images on Docker Hub](../docker-hub/vulnerability-scanning.md).

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
> Although [xwiki](https://hub.docker.com/_/xwiki){:
target="_blank" rel="noopener" class="_"} images may be detected as vulnerable
by some scanners, the authors believe the images are not vulnerable by Log4j 2
CVE as the API jars do not contain the vulnerability.
> The [Nuxeo](https://hub.docker.com/_/nuxeo){: target="_blank" rel="noopener" class="_"}
> image is deprecated and will not be updated.
