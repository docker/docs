---
title: Vulnerability scanning for Docker local images
description: Vulnerability scanning for Docker local images
keywords: Docker, scan, Snyk, images, local, CVE, vulnerability, security
toc_min: 1
toc_max: 2
---

{% include sign-up-cta.html
  body="Did you know that you can now get 10 free scans per month? Sign in to Docker to start scanning your images for vulnerabilities."
  header-text="Scan your images for free"
  target-url="https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade_scan"
%}

Looking to speed up your development cycles? Quickly detect and learn how to remediate CVEs in your images by running `docker scan IMAGE_NAME`. Check out [How to scan images](#how-to-scan-images) for details.

Vulnerability scanning for Docker local images  allows developers and development teams to review the security state of the container images and take actions to fix issues identified during the scan, resulting in more secure deployments. Docker Scan runs on Snyk engine, providing users with visibility into the security posture of their local Dockerfiles and local images.

Users trigger vulnerability scans through the CLI, and use the CLI to view the scan results. The scan results contain a list of Common Vulnerabilities and Exposures (CVEs), the sources, such as OS packages and libraries, versions in which they were introduced, and a recommended fixed version (if available) to remediate the CVEs discovered.

For information about the system requirements to run vulnerability scanning, see [Prerequisites](#prerequisites).

This page contains information about the `docker scan` CLI command. For information about automatically scanning Docker images through Docker Hub, see [Hub Vulnerability Scanning](/docker-hub/vulnerability-scanning/).

## How to scan images

The `docker scan` command allows you to scan existing Docker images using the image name or ID.  For example, run the following command to scan the hello-world image:

```console
$ docker scan hello-world

Testing hello-world...

Organization:      docker-desktop-test
Package manager:   linux
Project name:      docker-image|hello-world
Docker image:      hello-world
Licenses:          enabled

✓ Tested 0 dependencies for known issues, no vulnerable paths found.

Note that we do not currently have vulnerability data for your image.
```

### Get a detailed scan report

You can get a detailed scan report about a Docker image by providing the Dockerfile used to create the image. The syntax is `docker scan --file PATH_TO_DOCKERFILE DOCKER_IMAGE`.

For example, if you apply the option to the `docker-scan` test image, it displays the following result:

```console
$ docker scan --file Dockerfile docker-scan:e2e
Testing docker-scan:e2e
...
✗ High severity vulnerability found in perl
  Description: Integer Overflow or Wraparound
  Info: https://snyk.io/vuln/SNYK-DEBIAN10-PERL-570802
  Introduced through: git@1:2.20.1-2+deb10u3, meta-common-packages@meta
  From: git@1:2.20.1-2+deb10u3 > perl@5.28.1-6
  From: git@1:2.20.1-2+deb10u3 > liberror-perl@0.17027-2 > perl@5.28.1-6
  From: git@1:2.20.1-2+deb10u3 > perl@5.28.1-6 > perl/perl-modules-5.28@5.28.1-6
  and 3 more...
  Introduced by your base image (golang:1.14.6)

Organization:      docker-desktop-test
Package manager:   deb
Target file:       Dockerfile
Project name:      docker-image|99138c65ebc7
Docker image:      99138c65ebc7
Base image:        golang:1.14.6
Licenses:          enabled

Tested 200 dependencies for known issues, found 157 issues.

According to our scan, you are currently using the most secure version of the selected base image
```

### Excluding the base image

When using docker scan with the `--file` flag, you can also add the `--exclude-base` tag. This excludes the base image (specified in the Dockerfile using the `FROM` directive) vulnerabilities from your report. For example:

```console
$ docker scan --file Dockerfile --exclude-base docker-scan:e2e
Testing docker-scan:e2e
...
✗ Medium severity vulnerability found in libidn2/libidn2-0
  Description: Improper Input Validation
  Info: https://snyk.io/vuln/SNYK-DEBIAN10-LIBIDN2-474100
  Introduced through: iputils/iputils-ping@3:20180629-2+deb10u1, wget@1.20.1-1.1, curl@7.64.0-4+deb10u1, git@1:2.20.1-2+deb10u3
  From: iputils/iputils-ping@3:20180629-2+deb10u1 > libidn2/libidn2-0@2.0.5-1+deb10u1
  From: wget@1.20.1-1.1 > libidn2/libidn2-0@2.0.5-1+deb10u1
  From: curl@7.64.0-4+deb10u1 > curl/libcurl4@7.64.0-4+deb10u1 > libidn2/libidn2-0@2.0.5-1+deb10u1
  and 3 more...
  Introduced in your Dockerfile by 'RUN apk add -U --no-cache wget tar'



Organization:      docker-desktop-test
Package manager:   deb
Target file:       Dockerfile
Project name:      docker-image|99138c65ebc7
Docker image:      99138c65ebc7
Base image:        golang:1.14.6
Licenses:          enabled

Tested 200 dependencies for known issues, found 16 issues.
```

### Viewing the JSON output

You can also display the scan result as a JSON output by adding the `--json` flag to the command. For example:

```console
$ docker scan --json hello-world
{
  "vulnerabilities": [],
  "ok": true,
  "dependencyCount": 0,
  "org": "docker-desktop-test",
  "policy": "# Snyk (https://snyk.io) policy file, patches or ignores known vulnerabilities.\nversion: v1.19.0\nignore: {}\npatch: {}\n",
  "isPrivate": true,
  "licensesPolicy": {
    "severities": {},
    "orgLicenseRules": {
      "AGPL-1.0": {
        "licenseType": "AGPL-1.0",
        "severity": "high",
        "instructions": ""
      },
      ...
      "SimPL-2.0": {
        "licenseType": "SimPL-2.0",
        "severity": "high",
        "instructions": ""
      }
    }
  },
  "packageManager": "linux",
  "ignoreSettings": null,
  "docker": {
    "baseImageRemediation": {
      "code": "SCRATCH_BASE_IMAGE",
      "advice": [
        {
          "message": "Note that we do not currently have vulnerability data for your image.",
          "bold": true,
          "color": "yellow"
        }
      ]
    },
    "binariesVulns": {
      "issuesData": {},
      "affectedPkgs": {}
    }
  },
  "summary": "No known vulnerabilities",
  "filesystemPolicy": false,
  "uniqueCount": 0,
  "projectName": "docker-image|hello-world",
  "path": "hello-world"
}
```

In addition to the `--json` flag, you can also use the `--group-issues` flag to display a vulnerability only once in the scan report:

```console
$ docker scan --json --group-issues docker-scan:e2e
{
    {
      "title": "Improper Check for Dropped Privileges",
      ...
      "packageName": "bash",
      "language": "linux",
      "packageManager": "debian:10",
      "description": "## Overview\nAn issue was discovered in disable_priv_mode in shell.c in GNU Bash through 5.0 patch 11. By default, if Bash is run with its effective UID not equal to its real UID, it will drop privileges by setting its effective UID to its real UID. However, it does so incorrectly. On Linux and other systems that support \"saved UID\" functionality, the saved UID is not dropped. An attacker with command execution in the shell can use \"enable -f\" for runtime loading of a new builtin, which can be a shared object that calls setuid() and therefore regains privileges. However, binaries running with an effective UID of 0 are unaffected.\n\n## References\n- [CONFIRM](https://security.netapp.com/advisory/ntap-20200430-0003/)\n- [Debian Security Tracker](https://security-tracker.debian.org/tracker/CVE-2019-18276)\n- [GitHub Commit](https://github.com/bminor/bash/commit/951bdaad7a18cc0dc1036bba86b18b90874d39ff)\n- [MISC](http://packetstormsecurity.com/files/155498/Bash-5.0-Patch-11-Privilege-Escalation.html)\n- [MISC](https://www.youtube.com/watch?v=-wGtxJ8opa8)\n- [Ubuntu CVE Tracker](http://people.ubuntu.com/~ubuntu-security/cve/CVE-2019-18276)\n",
      "identifiers": {
        "ALTERNATIVE": [],
        "CVE": [
          "CVE-2019-18276"
        ],
        "CWE": [
          "CWE-273"
        ]
      },
      "severity": "low",
      "severityWithCritical": "low",
      "cvssScore": 7.8,
      "CVSSv3": "CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:H/I:H/A:H/E:F",
      ...
      "from": [
        "docker-image|docker-scan@e2e",
        "bash@5.0-4"
      ],
      "upgradePath": [],
      "isUpgradable": false,
      "isPatchable": false,
      "name": "bash",
      "version": "5.0-4"
    },
    ...
    "summary": "880 vulnerable dependency paths",
      "filesystemPolicy": false,
      "filtered": {
        "ignore": [],
        "patch": []
      },
      "uniqueCount": 158,
      "projectName": "docker-image|docker-scan",
      "platform": "linux/amd64",
      "path": "docker-scan:e2e"
}
```

You can find all the sources of the vulnerability in the `from` section.

### Checking the dependency tree

To view the dependency tree of your image, use the --dependency-tree flag. This displays all the dependencies before the scan result. For example:

```console
$ docker scan --dependency-tree debian:buster

$ docker-image|99138c65ebc7 @ latest
     ├─ ca-certificates @ 20200601~deb10u1
     │  └─ openssl @ 1.1.1d-0+deb10u3
     │     └─ openssl/libssl1.1 @ 1.1.1d-0+deb10u3
     ├─ curl @ 7.64.0-4+deb10u1
     │  └─ curl/libcurl4 @ 7.64.0-4+deb10u1
     │     ├─ e2fsprogs/libcom-err2 @ 1.44.5-1+deb10u3
     │     ├─ krb5/libgssapi-krb5-2 @ 1.17-3
     │     │  ├─ e2fsprogs/libcom-err2 @ 1.44.5-1+deb10u3
     │     │  ├─ krb5/libk5crypto3 @ 1.17-3
     │     │  │  └─ krb5/libkrb5support0 @ 1.17-3
     │     │  ├─ krb5/libkrb5-3 @ 1.17-3
     │     │  │  ├─ e2fsprogs/libcom-err2 @ 1.44.5-1+deb10u3
     │     │  │  ├─ krb5/libk5crypto3 @ 1.17-3
     │     │  │  ├─ krb5/libkrb5support0 @ 1.17-3
     │     │  │  └─ openssl/libssl1.1 @ 1.1.1d-0+deb10u3
     │     │  └─ krb5/libkrb5support0 @ 1.17-3
     │     ├─ libidn2/libidn2-0 @ 2.0.5-1+deb10u1
     │     │  └─ libunistring/libunistring2 @ 0.9.10-1
     │     ├─ krb5/libk5crypto3 @ 1.17-3
     │     ├─ krb5/libkrb5-3 @ 1.17-3
     │     ├─ openldap/libldap-2.4-2 @ 2.4.47+dfsg-3+deb10u2
     │     │  ├─ gnutls28/libgnutls30 @ 3.6.7-4+deb10u4
     │     │  │  ├─ nettle/libhogweed4 @ 3.4.1-1
     │     │  │  │  └─ nettle/libnettle6 @ 3.4.1-1
     │     │  │  ├─ libidn2/libidn2-0 @ 2.0.5-1+deb10u1
     │     │  │  ├─ nettle/libnettle6 @ 3.4.1-1
     │     │  │  ├─ p11-kit/libp11-kit0 @ 0.23.15-2
     │     │  │  │  └─ libffi/libffi6 @ 3.2.1-9
     │     │  │  ├─ libtasn1-6 @ 4.13-3
     │     │  │  └─ libunistring/libunistring2 @ 0.9.10-1
     │     │  ├─ cyrus-sasl2/libsasl2-2 @ 2.1.27+dfsg-1+deb10u1
     │     │  │  └─ cyrus-sasl2/libsasl2-modules-db @ 2.1.27+dfsg-1+deb10u1
     │     │  │     └─ db5.3/libdb5.3 @ 5.3.28+dfsg1-0.5
     │     │  └─ openldap/libldap-common @ 2.4.47+dfsg-3+deb10u2
     │     ├─ nghttp2/libnghttp2-14 @ 1.36.0-2+deb10u1
     │     ├─ libpsl/libpsl5 @ 0.20.2-2
     │     │  ├─ libidn2/libidn2-0 @ 2.0.5-1+deb10u1
     │     │  └─ libunistring/libunistring2 @ 0.9.10-1
     │     ├─ rtmpdump/librtmp1 @ 2.4+20151223.gitfa8646d.1-2
     │     │  ├─ gnutls28/libgnutls30 @ 3.6.7-4+deb10u4
     │     │  ├─ nettle/libhogweed4 @ 3.4.1-1
     │     │  └─ nettle/libnettle6 @ 3.4.1-1
     │     ├─ libssh2/libssh2-1 @ 1.8.0-2.1
     │     │  └─ libgcrypt20 @ 1.8.4-5
     │     └─ openssl/libssl1.1 @ 1.1.1d-0+deb10u3
     ├─ gnupg2/dirmngr @ 2.2.12-1+deb10u1
    ...

Organization:      docker-desktop-test
Package manager:   deb
Project name:      docker-image|99138c65ebc7
Docker image:      99138c65ebc7
Licenses:          enabled

Tested 200 dependencies for known issues, found 157 issues.

For more free scans that keep your images secure, sign up to Snyk at https://dockr.ly/3ePqVcp.
```

For more information about the vulnerability data, see [Docker Vulnerability Scanning CLI Cheat Sheet](https://goto.docker.com/rs/929-FJL-178/images/cheat-sheet-docker-desktop-vulnerability-scanning-CLI.pdf){: target="_blank" rel="noopener" class="_"}.

### Limiting the level of vulnerabilities displayed

Docker scan allows you to choose the level of vulnerabilities displayed in your scan report using the `--severity` flag.
You can set the severity flag to `low`, `medium`, or` high` depending on the level of vulnerabilities you’d like to see in your report.  
For example, if you set the severity level as `medium`, the scan report displays all vulnerabilities that are classified as medium and high.
 
 ```console
$ docker scan --severity=medium docker-scan:e2e 
./bin/docker-scan_darwin_amd64 scan --severity=medium docker-scan:e2e

Testing docker-scan:e2e...

✗ Medium severity vulnerability found in sqlite3/libsqlite3-0
  Description: Divide By Zero
  Info: https://snyk.io/vuln/SNYK-DEBIAN10-SQLITE3-466337
  Introduced through: gnupg2/gnupg@2.2.12-1+deb10u1, subversion@1.10.4-1+deb10u1, mercurial@4.8.2-1+deb10u1
  From: gnupg2/gnupg@2.2.12-1+deb10u1 > gnupg2/gpg@2.2.12-1+deb10u1 > sqlite3/libsqlite3-0@3.27.2-3
  From: subversion@1.10.4-1+deb10u1 > subversion/libsvn1@1.10.4-1+deb10u1 > sqlite3/libsqlite3-0@3.27.2-3
  From: mercurial@4.8.2-1+deb10u1 > python-defaults/python@2.7.16-1 > python2.7@2.7.16-2+deb10u1 > python2.7/libpython2.7-stdlib@2.7.16-2+deb10u1 > sqlite3/libsqlite3-0@3.27.2-3

✗ Medium severity vulnerability found in sqlite3/libsqlite3-0
  Description: Uncontrolled Recursion
...
✗ High severity vulnerability found in binutils/binutils-common
  Description: Missing Release of Resource after Effective Lifetime
  Info: https://snyk.io/vuln/SNYK-DEBIAN10-BINUTILS-403318
  Introduced through: gcc-defaults/g++@4:8.3.0-1
  From: gcc-defaults/g++@4:8.3.0-1 > gcc-defaults/gcc@4:8.3.0-1 > gcc-8@8.3.0-6 > binutils@2.31.1-16 > binutils/binutils-common@2.31.1-16
  From: gcc-defaults/g++@4:8.3.0-1 > gcc-defaults/gcc@4:8.3.0-1 > gcc-8@8.3.0-6 > binutils@2.31.1-16 > binutils/libbinutils@2.31.1-16 > binutils/binutils-common@2.31.1-16
  From: gcc-defaults/g++@4:8.3.0-1 > gcc-defaults/gcc@4:8.3.0-1 > gcc-8@8.3.0-6 > binutils@2.31.1-16 > binutils/binutils-x86-64-linux-gnu@2.31.1-16 > binutils/binutils-common@2.31.1-16
  and 4 more...

Organization:      docker-desktop-test
Package manager:   deb
Project name:      docker-image|docker-scan
Docker image:      docker-scan:e2e
Platform:          linux/amd64
Licenses:          enabled

Tested 200 dependencies for known issues, found 37 issues.
```

## Provider authentication

If you have an existing Snyk account, you can directly use your Snyk [API token](https://app.snyk.io/account){: target="_blank" rel="noopener" class="_"}:

```console
$ docker scan --login --token SNYK_AUTH_TOKEN

Your account has been authenticated. Snyk is now ready to be used.
```

If you use the `--login` flag without any token, you will be redirected to the Snyk website to login.

## Prerequisites

To run vulnerability scanning on your Docker images, you must meet the following requirements:

1. Download and install Docker Desktop.

    - [Download for Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-amd64)
    - [Download for Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-arm64)
    - [Download for Windows](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe)

2. Sign into [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}.

3. From the Docker Desktop menu, select **Sign in/ Create Docker ID**. Alternatively, open a terminal and run the command `docker login`.

4. (Optional) You can create a [Snyk account](https://dockr.ly/3ePqVcp){: target="_blank" rel="noopener" class="_"} for scans, or use the additional monthly free scans provided by Snyk with your Docker Hub account.

Check your installation by running `docker scan --version`, it should print the current version of docker scan and the Snyk engine version. For example:

```console
$ docker scan --version
Version:    v0.5.0
Git commit: 5a09266
Provider:   Snyk (1.432.0)
```

> **Note:**
>
> Docker Scan uses the Snyk binary installed in your environment by default. If 
this is not available, it uses the Snyk binary embedded in Docker Desktop.
> The minimum version required for Snyk is `1.385.0`.

## Supported options

The high-level `docker scan` command scans local images using the image name or the image ID. It supports the following options:

| Option                                                       | Description                                   |
|:------------------------------------------------------------------ :------------------------------------------------|
| `--accept-license` | Accept the license agreement of the third-party scanning provider    |
| `--dependency-tree` | Display the dependency tree of the image along with scan results |
| `--exclude-base` | Exclude the base image during scanning. This option requires the --file option to be set |
| `-f`, `--file string` | Specify the location of the Dockerfile associated with the image. This option displays a detailed scan result |
| `--json` | Display the result of the scan in JSON format|
| `--login` | Log into Snyk using an optional token (using the flag --token), or by using a web-based token |
| `--reject-license` | Reject the license agreement of the third-party scanning provider |
| `--severity string` | Only report vulnerabilities of provided level or higher (low, medium, high) |
| `--token string`  | Use the authentication token to log into the third-party scanning provider |
| `--version` | Display the Docker Scan plugin version |

## Known issues

**WSL 2**

- The Vulnerability scanning feature doesn’t work with Alpine distributions.
- If you are using Debian and OpenSUSE distributions, the login process only works with the `--token` flag, you won’t be redirected to the Snyk website for authentication.

## Feedback

Your feedback is very important to us. Let us know your feedback by creating an issue in the [scan-cli-plugin](https://github.com/docker/scan-cli-plugin/issues/new){: target="_blank" rel="noopener" class="_"} GitHub repository.
