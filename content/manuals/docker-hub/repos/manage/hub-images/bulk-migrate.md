---
title: Bulk migrate Docker images
description: Learn how to migrate multiple Docker images and tags between organizations using scripts and automation.
keywords: docker hub, migration, bulk, images, tags, multi-arch, buildx
---

This guide shows you how to migrate Docker images in bulk between Docker Hub
organizations or namespaces. Whether you're consolidating repositories,
changing organization structure, or moving images to a new account, these
techniques help you migrate efficiently while preserving image integrity.

## Prerequisites

Before you begin, ensure you have:

- Docker CLI version 20.10 or later installed
- Docker Buildx (optional but recommended for multi-architecture images)
- Push access to both source and destination organizations
- `jq` installed for JSON parsing in scripts
- `curl` for API calls

## Authenticate to Docker Hub

Sign in to Docker Hub to authenticate your session:

```console
$ docker login
```

Enter your credentials when prompted. This authentication persists for your
session and prevents rate limiting issues.

## Migrate a single image tag

The basic workflow for migrating a single image tag involves three steps: pull,
tag, and push.

1. Set your source and destination variables:

   ```bash
   SRC_ORG=oldorg
   DEST_ORG=neworg
   REPO=myapp
   TAG=1.2.3
   ```

2. Pull the image from the source organization:

   ```console
   $ docker pull ${SRC_ORG}/${REPO}:${TAG}
   ```

3. Tag the image for the destination organization:

   ```console
   $ docker tag ${SRC_ORG}/${REPO}:${TAG} ${DEST_ORG}/${REPO}:${TAG}
   ```

4. Push the image to the destination organization:

   ```console
   $ docker push ${DEST_ORG}/${REPO}:${TAG}
   ```

Repeat these steps for any additional tags you need to migrate, including
`latest` if applicable.

## Migrate all tags for a repository

To migrate all tags from a single repository, use this script that queries the
Docker Hub API and processes each tag:

```bash
#!/usr/bin/env bash
set -euo pipefail

SRC_ORG="oldorg"
DEST_ORG="neworg"
REPO="myapp"

# Paginate through tags
TAGS_URL="https://hub.docker.com/v2/repositories/${SRC_ORG}/${REPO}/tags?page_size=100"
while [[ -n "${TAGS_URL}" && "${TAGS_URL}" != "null" ]]; do
  RESP=$(curl -fsSL "${TAGS_URL}")
  echo "${RESP}" | jq -r '.results[].name' | while read -r TAG; do
    echo "==> Migrating ${SRC_ORG}/${REPO}:${TAG} → ${DEST_ORG}/${REPO}:${TAG}"
    docker pull "${SRC_ORG}/${REPO}:${TAG}"
    docker tag  "${SRC_ORG}/${REPO}:${TAG}" "${DEST_ORG}/${REPO}:${TAG}"
    docker push "${DEST_ORG}/${REPO}:${TAG}"
  done
  TAGS_URL=$(echo "${RESP}" | jq -r '.next')
done
```

This script automatically handles pagination when a repository has more than
100 tags.

> [!NOTE]
>
> Docker Hub automatically creates the destination repository on first push if
> your account has the necessary permissions.

### Migrate private repository tags

For private repositories, authenticate your API calls with a Docker Hub access
token:

1. Create a personal access token in your
   [Docker Hub account settings](https://hub.docker.com/settings/security).

2. Set your credentials as variables:

   ```bash
   HUB_USER="your-username"
   HUB_TOKEN="your-access-token"
   ```

3. Modify the `curl` command in the script to include authentication:

   ```bash
   RESP=$(curl -fsSL -u "${HUB_USER}:${HUB_TOKEN}" "${TAGS_URL}")
   ```

> [!IMPORTANT]
>
> If you encounter pull rate or throughput limits, keep `docker login` active
> to avoid anonymous pulls. Consider adding throttling or careful parallelization
> if migrating large numbers of images.

## Migrate multiple repositories

To migrate multiple repositories at once, create a list of repository names
and process them in a loop.

1. Create a file named `repos.txt` with one repository name per line:

   ```text
   api
   web
   worker
   database
   ```

2. Save the single-repository script from the previous section as
   `migrate-single-repo.sh` and make it executable.

3. Use this wrapper script to process all repositories:

   ```bash
   #!/usr/bin/env bash
   set -euo pipefail

   SRC_ORG="oldorg"
   DEST_ORG="neworg"

   while read -r REPO; do
     [[ -z "${REPO}" ]] && continue
     echo "==== Migrating repo: ${REPO}"
     export REPO
     ./migrate-single-repo.sh
   done < repos.txt
   ```

## Preserve multi-architecture images

Standard `docker pull` only retrieves the image for your current platform.
For multi-architecture images, this approach loses other platform variants.

### Use Buildx imagetools (recommended)

The recommended approach uses Buildx to copy the complete manifest without
pulling images locally:

```console
$ docker buildx imagetools create \
  -t ${DEST_ORG}/${REPO}:${TAG} \
     ${SRC_ORG}/${REPO}:${TAG}
```

This command copies the source manifest with all platforms directly to the
destination tag.

Verify the migration by inspecting both manifests:

```console
$ docker buildx imagetools inspect ${SRC_ORG}/${REPO}:${TAG}
$ docker buildx imagetools inspect ${DEST_ORG}/${REPO}:${TAG}
```

Compare the platforms and digests in the output to confirm they match.

### Manual manifest creation

If you need to use the pull/tag/push workflow for multi-architecture images,
you must pull each platform variant and recreate the manifest using
`docker manifest create` and `docker manifest push`. This approach is slower
and more error-prone than using Buildx imagetools.

## Verify migration integrity

After migrating images, verify that they transferred correctly.

### Single-architecture images

Compare image digests between source and destination:

```console
$ docker pull ${SRC_ORG}/${REPO}:${TAG}
$ docker inspect --format='{{index .RepoDigests 0}}' ${SRC_ORG}/${REPO}:${TAG}

$ docker pull ${DEST_ORG}/${REPO}:${TAG}
$ docker inspect --format='{{index .RepoDigests 0}}' ${DEST_ORG}/${REPO}:${TAG}
```

The SHA256 digests should match if the migration succeeded.

### Multi-architecture images

For multi-arch images, compare the output from Buildx imagetools:

```console
$ docker buildx imagetools inspect ${SRC_ORG}/${REPO}:${TAG}
$ docker buildx imagetools inspect ${DEST_ORG}/${REPO}:${TAG}
```

Verify that the platforms and manifest digest match between source and
destination.

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
