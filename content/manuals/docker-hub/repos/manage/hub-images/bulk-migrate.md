---
title: Bulk migrate images
description: Learn how to migrate multiple Docker images and tags between organizations using scripts and automation.
keywords: docker hub, migration, bulk, images, tags, multi-arch
---

This guide shows you how to migrate Docker images in bulk between Docker Hub
organizations or namespaces. Whether you're consolidating repositories, changing
organization structure, or moving images to a new account, these techniques help
you migrate efficiently while preserving image integrity.

The topic is structured to build up in scale:

1. [Migrate a single image tag](#migrate-a-single-image-tag)
2. [Migrate all tags for a repository](#migrate-all-tags-for-a-repository)
3. [Migrate multiple repositories](#migrate-multiple-repositories)

The recommended tool for this workflow is `crane`. An equivalent alternative
using `regctl` is also shown. Both tools perform registry-to-registry copies
without pulling images locally and preserve multi-architecture images.

`crane` is recommended for its simplicity and focused image-copying workflow.
`regctl` is also a good choice, particularly if you already use it for broader
registry management tasks beyond image copying.

> [!NOTE]
>
> The main workflows in this topic operate on tagged images only. Untagged
> manifests or content no longer reachable from tags are not migrated. In
> practice, these are usually unused artifacts, but be aware of this limitation
> before migration. While you can migrate specific untagged manifests using
> [digest references](#migrate-by-digest), there is no API to enumerate untagged
> manifests in a repository.

## Prerequisites

Before you begin, ensure you have:

- One of the following installed and available in your `$PATH`:
  - [`crane`](https://github.com/google/go-containerregistry)
  - [`regctl`](https://regclient.org/usage/regctl/)
- Push access to both the source and destination organizations
- Registry authentication configured for your chosen tool

## Authenticate to registries

Both tools authenticate directly against registries:

- `crane` uses Docker credential helpers and `~/.docker/config.json`. See the
  [crane documentation](https://github.com/google/go-containerregistry/tree/main/cmd/crane/doc).
- `regctl` uses its own configuration file and can import Docker credentials.
  See the [regctl documentation](https://github.com/regclient/regclient/tree/main/docs).

Follow the authentication instructions for your registry and tool of choice.

## Migrate a single image tag

This is the simplest and most common migration scenario.

The following example script copies the image manifest directly between
registries and preserves multi-architecture images when present. Repeat this
process for each tag you want to migrate. Replace the environment variable
values with your source and destination organization names, repository name, and
tag.

```bash
SRC_ORG="oldorg"
DEST_ORG="neworg"
REPO="myapp"
TAG="1.2.3"

SRC_IMAGE="${SRC_ORG}/${REPO}:${TAG}"
DEST_IMAGE="${DEST_ORG}/${REPO}:${TAG}"

# Using crane (recommended)
crane cp "${SRC_IMAGE}" "${DEST_IMAGE}"

# Using regctl (alternative)
# regctl image copy "${SRC_IMAGE}" "${DEST_IMAGE}"
```

### Migrate by digest

To migrate a specific image by digest instead of tag, use the digest in the
source reference. This is useful when you need to migrate an exact image
version, even if the tag has been updated. Replace the environment variable
values with your source and destination organization names, repository name,
digest, and tag. You can choose between `crane` and `regctl` for the copy
operation.

```bash
SRC_ORG="oldorg"
DEST_ORG="neworg"
REPO="myapp"
DIGEST="sha256:abcd1234..."
TAG="stable"

SRC_IMAGE="${SRC_ORG}/${REPO}@${DIGEST}"
DEST_IMAGE="${DEST_ORG}/${REPO}:${TAG}"

# Using crane
crane cp "${SRC_IMAGE}" "${DEST_IMAGE}"

# Using regctl
# regctl image copy "${SRC_IMAGE}" "${DEST_IMAGE}"
```

## Migrate all tags for a repository

To migrate every tagged image in a repository, use the Docker Hub API to
enumerate tags and copy each one. The following example script retrieves all
tags for a given repository and migrates them in a loop. This approach scales to
repositories with many tags without overwhelming local resources. Note that
there is a rate limit on Docker Hub requests, so you may need to add delays or
pagination handling for very large repositories.

Replace the environment variable values with your source and destination
organization names and repository name. If your source repository is private,
also set `HUB_USER` and `HUB_TOKEN` with credentials that have pull access. You
can also choose between `crane` and `regctl` for the copy operation.

```bash
#!/usr/bin/env bash
set -euo pipefail

# Use environment variables if set, otherwise use defaults
SRC_ORG="${SRC_ORG:-oldorg}"
DEST_ORG="${DEST_ORG:-neworg}"
REPO="${REPO:-myapp}"

# Optional: for private repositories
# HUB_USER="your-username"
# HUB_TOKEN="your-access-token"
# AUTH="-u ${HUB_USER}:${HUB_TOKEN}"
AUTH=""

TOOL="crane"   # or: TOOL="regctl"

TAGS_URL="https://hub.docker.com/v2/repositories/${SRC_ORG}/${REPO}/tags?page_size=100"

while [[ -n "${TAGS_URL}" && "${TAGS_URL}" != "null" ]]; do
  RESP=$(curl -fsSL ${AUTH} "${TAGS_URL}")

  echo "${RESP}" | jq -r '.results[].name' | while read -r TAG; do
    [[ -z "${TAG}" ]] && continue

    SRC_IMAGE="${SRC_ORG}/${REPO}:${TAG}"
    DEST_IMAGE="${DEST_ORG}/${REPO}:${TAG}"

    echo "Migrating ${SRC_IMAGE} → ${DEST_IMAGE}"

    case "${TOOL}" in
      crane)
        crane cp "${SRC_IMAGE}" "${DEST_IMAGE}"
        ;;
      regctl)
        regctl image copy "${SRC_IMAGE}" "${DEST_IMAGE}"
        ;;
    esac
  done

  TAGS_URL=$(echo "${RESP}" | jq -r '.next')
done
```

> [!NOTE]
>
> Docker Hub automatically creates the destination repository on first push if
> your account has permission.

## Migrate multiple repositories

To migrate several repositories, create a list and run the single-repository
script for each one.

For example, create a `repos.txt` file with repository names:

```text
api
web
worker
```

Save the script from the previous section as `migrate-single-repo.sh`. Then, run
the following example script that processes each repository in the file. Replace
the environment variable values with your source and destination organization
names.

```bash
#!/usr/bin/env bash
set -euo pipefail

SRC_ORG="oldorg"
DEST_ORG="neworg"

while read -r REPO; do
  [[ -z "${REPO}" ]] && continue
  echo "==== Migrating repo: ${REPO}"
  SRC_ORG="${SRC_ORG}" DEST_ORG="${DEST_ORG}" REPO="${REPO}" ./migrate-single-repo.sh
done < repos.txt
```

## Verify migration integrity

After copying, verify that source and destination match by comparing digests.

### Basic digest verification

The following example script retrieves the image digest for a specific tag from
both source and destination and compares them. If the digests match, the
migration is successful. Replace the environment variable values with your
source and destination organization names, repository name, and tag. You can
choose between `crane` and `regctl` for retrieving digests.

```bash
SRC_ORG="oldorg"
DEST_ORG="neworg"
REPO="myapp"
TAG="1.2.3"

SRC_IMAGE="${SRC_ORG}/${REPO}:${TAG}"
DEST_IMAGE="${DEST_ORG}/${REPO}:${TAG}"

# Using crane
SRC_DIGEST=$(crane digest "${SRC_IMAGE}")
DEST_DIGEST=$(crane digest "${DEST_IMAGE}")

# Using regctl (alternative)
# SRC_DIGEST=$(regctl image digest "${SRC_IMAGE}")
# DEST_DIGEST=$(regctl image digest "${DEST_IMAGE}")

echo "Source:      ${SRC_DIGEST}"
echo "Destination: ${DEST_DIGEST}"

if [[ "${SRC_DIGEST}" == "${DEST_DIGEST}" ]]; then
  echo "✓ Migration verified: digests match"
else
  echo "✗ Migration failed: digests do not match"
  exit 1
fi
```

### Multi-arch verification

For multi-architecture images, also verify the manifest list to ensure all
platforms were copied correctly. Replace the environment variable values with
your source and destination organization names, repository name, and tag. You
can choose between `crane` and `regctl` for retrieving manifests.

```bash
SRC_ORG="oldorg"
DEST_ORG="neworg"
REPO="myapp"
TAG="1.2.3"

SRC_IMAGE="${SRC_ORG}/${REPO}:${TAG}"
DEST_IMAGE="${DEST_ORG}/${REPO}:${TAG}"

# Using crane
SRC_MANIFEST=$(crane manifest "${SRC_IMAGE}")
DEST_MANIFEST=$(crane manifest "${DEST_IMAGE}")

# Using regctl (alternative)
# SRC_MANIFEST=$(regctl image manifest --format raw-body "${SRC_IMAGE}")
# DEST_MANIFEST=$(regctl image manifest --format raw-body "${DEST_IMAGE}")

# Check if it's a manifest list (multi-arch)
if echo "${SRC_MANIFEST}" | jq -e '.manifests' > /dev/null 2>&1; then
  echo "Multi-arch image detected"
  
  # Compare platform list
  SRC_PLATFORMS=$(echo "${SRC_MANIFEST}" | jq -r '.manifests[] | "\(.platform.os)/\(.platform.architecture)"' | sort)
  DEST_PLATFORMS=$(echo "${DEST_MANIFEST}" | jq -r '.manifests[] | "\(.platform.os)/\(.platform.architecture)"' | sort)
  
  if [[ "${SRC_PLATFORMS}" == "${DEST_PLATFORMS}" ]]; then
    echo "✓ Platform list matches:"
    echo "${SRC_PLATFORMS}"
  else
    echo "✗ Platform lists do not match"
    echo "Source platforms:"
    echo "${SRC_PLATFORMS}"
    echo "Destination platforms:"
    echo "${DEST_PLATFORMS}"
    exit 1
  fi
else
  echo "Single-arch image"
fi
```

## Complete the migration

After migrating your images, complete these additional steps:

1. Copy repository metadata in the Docker Hub UI or via API:

   - README content
   - Repository description
   - Topics and tags

2. Configure repository settings to match the source:

   - Visibility (public or private)
   - Team permissions and access controls

3. Reconfigure integrations in the destination organization:

   - Webhooks
   - Automated builds
   - Security scanners

4. Update image references in your projects:

   - Change `FROM oldorg/repo:tag` to `FROM neworg/repo:tag` in Dockerfiles
   - Update deployment configurations
   - Update documentation

5. Deprecate the old location:
   - Update the source repository description to point to the new location
   - Consider adding a grace period before making the old repository private or
     read-only