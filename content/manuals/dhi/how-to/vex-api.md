---
title: Query VEX for a Docker Hardened Image
linktitle: Query VEX suppressions via API
description: Use the DHI GraphQL API to fetch VEX statements and suppressed CVEs for a Docker Hardened Image by digest.
keywords: dhi vex api, graphql api, suppressed cves, base vex, image digest, docker hardened images
weight: 47
---

This is a guided example of using the [DHI API](/manuals/dhi/tools/api.md)
to fetch VEX statements and suppressed CVEs for a Docker Hardened Image (DHI)
by digest, using the `imagePackagesForImageCoords` query. The query takes an
image digest and returns every package, every CVE reported against it,
whether Docker suppresses that CVE, and why.

> [!NOTE]
> For one-off lookups or local scanning workflows, [`docker scout vex
> get`](/manuals/dhi/how-to/scan.md#export-vex-attestations) may be simpler.
> Use the API described here when you need to query many images
> programmatically.

For the API endpoint and authentication, see [Use the DHI
API](/manuals/dhi/tools/api.md). This page assumes you already have a valid
token.

## Query suppressed CVEs

This is the `imagePackagesForImageCoords` query used in this example. It
fetches every package in an image, every CVE reported against it, and
whether Docker suppresses that CVE:

```graphql
query BaseVex($ctx: Context!, $q: IpImagePackagesForImageCoordsQuery!) {
  imagePackagesForImageCoords(context: $ctx, query: $q) {
    imagePackages {
      packages {
        package {
          purl
          name
          version
          vulnerabilities {
            sourceId
            isExcepted
            fixedBy
            cvss {
              severity
            }
            vulnerabilityExceptions {
              id
              sourceType
              type
              justification
              additionalDetails
              isDhiStatement
            }
          }
        }
      }
    }
  }
}
```

With variables:

```json
{
  "ctx": { "organization": "your-org" },
  "q": {
    "digest": "sha256:<platform-manifest-digest>",
    "hostName": "hub.docker.com",
    "repoName": "your-org/your-repo",
    "includeExcepted": true,
    "includeNodsa": true
  }
}
```

For the full argument reference, see
[`imagePackagesForImageCoords`](/manuals/dhi/tools/api.md#imagepackagesforimagecoords).

### Get an access token

Exchange your organization access token (OAT) or personal access token (PAT)
for an access token, as described in [Use the DHI
API](/manuals/dhi/tools/api.md#authentication):

```console
$ DHI_API_TOKEN=$(curl -s -X POST https://hub.docker.com/v2/auth/token \
  -H "Content-Type: application/json" \
  -d "{\"identifier\": \"<identifier>\", \"secret\": \"<token>\"}" \
  | jq -r .access_token)
```

Use your Docker Hub username as `identifier` for a PAT, or your organization
name for an OAT.

### Send the request

Combine the query and variables into the request body described in [Use the
DHI API](/manuals/dhi/tools/api.md#request-format). The following example
uses `jq` to build the body safely, since the query spans multiple lines:

```console
$ QUERY='query BaseVex($ctx: Context!, $q: IpImagePackagesForImageCoordsQuery!) {
  imagePackagesForImageCoords(context: $ctx, query: $q) {
    imagePackages {
      packages {
        package {
          purl
          name
          version
          vulnerabilities {
            sourceId
            isExcepted
            fixedBy
            cvss {
              severity
            }
            vulnerabilityExceptions {
              id
              sourceType
              type
              justification
              additionalDetails
              isDhiStatement
            }
          }
        }
      }
    }
  }
}'

$ curl https://api.dso.docker.com/v1/graphql \
  -H "Authorization: Bearer $DHI_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d "$(jq -n --arg query "$QUERY" '{
    query: $query,
    variables: {
      ctx: { organization: "your-org" },
      q: {
        digest: "sha256:<platform-manifest-digest>",
        hostName: "hub.docker.com",
        repoName: "your-org/your-repo",
        includeExcepted: true,
        includeNodsa: true
      }
    }
  }')"
```

Substitute your organization, digest, host, and repository.

## Sample response

The following example is trimmed to a single package:

```json
{
  "package": {
    "purl": "pkg:deb/debian/tar@1.35%2Bdfsg-3.1%2Bdhi1?os_distro=trixie&os_name=debian&os_version=13",
    "name": "tar",
    "version": "1.35+dfsg-3.1+dhi1",
    "vulnerabilities": [
      {
        "sourceId": "CVE-2025-45582",
        "cvss": { "severity": "MEDIUM" },
        "fixedBy": null,
        "isExcepted": false,
        "vulnerabilityExceptions": []
      },
      {
        "sourceId": "CVE-2026-18477",
        "cvss": { "severity": "MEDIUM" },
        "fixedBy": null,
        "isExcepted": true,
        "vulnerabilityExceptions": [
          {
            "id": "debian-nodsa-CVE-2026-18477",
            "sourceType": "EXTERNAL",
            "type": "FALSE_POSITIVE",
            "justification": null,
            "additionalDetails": "Debian NODSA",
            "isDhiStatement": false
          }
        ]
      }
    ]
  }
}
```

Use `isExcepted` to tell whether a CVE is suppressed. For the full response
field reference and the OpenVEX mapping, see
[`imagePackagesForImageCoords`](/manuals/dhi/tools/api.md#imagepackagesforimagecoords).

## Query customization images

DHI VEX statements apply directly to images built with the [customization
feature](/manuals/dhi/how-to/customize.md). Query the customized image's own
digest to get its packages, CVEs, and suppressions in a single call. You
don't need to query the base image separately.

## Caveats

The suppression set for a given digest can change as new advisories and
assessments are published. Cache results per digest with a short TTL rather
than treating a response as permanent.
