---
title: 'Mirror a Docker Hardened Image repository <span class="not-prose bg-blue-500 dark:bg-blue-400 rounded-sm px-1 text-xs text-white whitespace-nowrap">DHI Enterprise</span>'
linktitle: Mirror a repository
description: Learn how to mirror an image into your organization's namespace and optionally push it to another private registry.
weight: 20
keywords: mirror docker image, private container registry, docker hub automation, webhook image sync, secure image distribution, internal registry, jfrog artifactory, harbor registry, amazon ecr, google artifact registry, github container registry
---

{{< summary-bar feature_name="Docker Hardened Images" >}}

Mirroring requires a DHI Select or Enterprise subscription. Without a
subscription, you can pull Docker Hardened Images directly from `dhi.io` without
mirroring. With a DHI Select or Enterprise subscription, you must mirror to get:

- Compliance variants (FIPS-enabled or STIG-ready images)
- Extended Lifecycle Support (ELS) variants (requires add-on)
- Image or Helm chart customization
- Air-gapped or restricted network environments
- SLA-backed security updates

## How to mirror

This topic covers two types of mirroring for Docker Hardened Image (DHI)
repositories:

- [Mirror to Docker Hub](#mirror-a-dhi-repository-to-docker-hub): Mirror a DHI
  repository to your organization's namespace on Docker Hub. This requires a DHI
  Enterprise subscription and is used to [customize an image or
  chart](./customize.md) and access compliance variants and ELS variants
  (requires add-on). This must be done through the Docker Hub web interface.

- [Mirror to a third-party
  registry](#mirror-a-dhi-repository-to-a-third-party-registry): Mirror a
  repository to another container registry, such as Amazon ECR, Google Artifact
  Registry, or a private Harbor instance.

## Mirror a DHI repository to Docker Hub

Mirroring a repository to Docker Hub requires a DHI Enterprise subscription and
enables access to compliance variants, Extended Lifecycle Support (ELS) variants
(requires add-on), and customization capabilities:

- Image repositories: Mirroring lets you customize images by adding packages,
  OCI artifacts (such as custom certificates or additional tools), environment
  variables, labels, and other configuration settings. For more details, see
  [Customize a Docker Hardened Image](./customize.md#customize-a-docker-hardened-image).

- Chart repositories: Mirroring lets you customize image references within
  the chart. This is particularly useful when using customized images or when
  you've mirrored images to a third-party registry and need the chart to
  reference those custom locations. For more details, see [Customize a Docker
  Hardened Helm chart](./customize.md#customize-a-docker-hardened-helm-chart).

Only organization owners can perform mirroring. Once mirrored, the repository
becomes available in your organization's namespace, and you can customize it as
needed.

You can mirror repositories using either the Docker Hub web interface or the DHI CLI.

### Mirror using the DHI CLI

The DHI CLI provides a command-line interface for managing Docker Hardened
Images, including mirroring operations. For installation instructions and usage
details, see [Use the DHI CLI](./cli.md#mirror-dhi-images).

### Stop mirroring with the CLI

```console
$ docker dhi mirror stop --org my-org dhi-golang
```

After stopping mirroring, the repository remains but will no longer receive updates.

### Mirror using the Docker Hub web interface

To mirror a Docker Hardened Image repository using the web interface:

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub**.
3. In the namespace drop-down, select your organization.
4. Select **Hardened Images** > **Catalog**.
5. Select a DHI repository to view its details.
6. Mirror the repository:
    - To mirror an image repository, select **Use this image** > **Mirror
      repository**, and then follow the on-screen instructions. If you have the ELS add-on, you can also
      select **Enable support for end-of-life versions**.
    - To mirror a Helm chart repository, select **Get Helm chart**, and then follow the on-screen instructions.

It may take a few minutes for all the tags to finish mirroring.

After mirroring a repository, the repository appears in your organization's
repository list, prefixed by `dhi-`. It will continue to receive updated images.

Once mirrored, the repository is a standard Docker Hub repository in your
organization's namespace. It behaves exactly like any other Hub repository,
which means you can manage access and permissions, configure webhooks, and use
other standard Hub features. See [Docker Hub
repositories](/manuals/docker-hub/repos/_index.md) for details.

Additionally, mirrored DHI repositories let you customize images and charts. To
learn more, see [Customize a Docker Hardened Image or chart](./customize.md).

### Webhook integration for syncing and alerts

To keep external registries or systems in sync with your mirrored Docker
Hardened Images, and to receive notifications when updates occur, you can
configure a [webhook](/docker-hub/repos/manage/webhooks/) on the mirrored
repository in Docker Hub. A webhook sends a `POST` request to a URL you define
whenever a new image tag is pushed or updated.

For example, you might configure a webhook to call a CI/CD system at
`https://ci.example.com/hooks/dhi-sync` whenever a new tag is mirrored. The
automation triggered by this webhook can pull the updated image from Docker Hub
and push it to an internal registry such as Amazon ECR, Google Artifact
Registry, or GitHub Container Registry.

Other common webhook use cases include:

- Triggering validation or vulnerability scanning workflows
- Signing or promoting images
- Sending notifications to downstream systems

#### Example webhook payload

When a webhook is triggered, Docker Hub sends a JSON payload like the following:

```json
{
  "callback_url": "https://registry.hub.docker.com/u/exampleorg/dhi-python/hook/abc123/",
  "push_data": {
    "pushed_at": 1712345678,
    "pusher": "trustedbuilder",
    "tag": "3.13-alpine3.21"
  },
  "repository": {
    "name": "dhi-python",
    "namespace": "exampleorg",
    "repo_name": "exampleorg/dhi-python",
    "repo_url": "https://hub.docker.com/r/exampleorg/dhi-python",
    "is_private": true,
    "status": "Active",
    ...
  }
}
```

### Stop mirroring a repository

Only organization owners can stop mirroring a repository. After you stop
mirroring, the repository remains, but it will
no longer receive updates. You can still use the last images or charts that were mirrored,
but the repository will not receive new tags or updates from the original
repository.

> [!NOTE]
>
> If you only want to stop mirroring ELS versions, you can clear the ELS
> option in the mirrored repository's **Settings** tab. For more details, see
> [Disable ELS for a repository](./els.md#disable-els-for-a-repository).

 To stop mirroring a repository:

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub**.
3. In the namespace drop-down, select your organization that has access to DHI.
4. Select **Hardened Images** > **Manage**.
5. Select the **Mirrored Images** or **Mirrored Helm charts** tab.
6. In the far right column of the repository you want to stop mirroring, select the menu icon.
7. Select **Stop mirroring**.

## Mirror a DHI repository to a third-party registry

You can optionally mirror a DHI repository to another container registry, such as Amazon
ECR, Google Artifact Registry, GitHub Container Registry, or a private Harbor
instance.

You can use any standard workflow to mirror the image, such as the
[Docker CLI](/reference/cli/docker/), [Docker Hub Registry
API](/reference/api/registry/latest/), third-party registry tools, or CI/CD
automation.

However, to preserve the full security context, including attestations, you must
also mirror its associated OCI artifacts. DHI repositories store the image
layers on `dhi.io` (or `docker.io` for customized images) and the signed
attestations in a separate registry (`registry.scout.docker.com`).

To copy both, you can use [`regctl`](https://regclient.org/cli/regctl/), an
OCI-aware CLI that supports mirroring images along with attached artifacts such
as SBOMs, vulnerability reports, and SLSA provenance. For ongoing synchronization,
you can use [`regsync`](https://regclient.org/cli/regsync/).

### Authenticate to `dhi.io` with an organization access token

You can authenticate to `dhi.io` using an [organization access token
(OAT)](../../enterprise/security/access-tokens.md) instead of a personal access
token (PAT). OATs are owned by the organization rather than an individual user,
which makes them better suited for CI/CD pipelines and automated workflows.

> [!NOTE]
>
> When using an OAT, use your organization name as the username, not your
> personal Docker ID. OATs are org-scoped and will return a `401 Unauthorized`
> error if presented under an individual user's username.

To authenticate using an OAT:

1. Sign in to [Docker Home](https://app.docker.com) and select your organization.
2. Select **Admin Console**, then **Access tokens**.
3. Select **Generate access token**.
4. Give the token a descriptive name, for example `dhi-pull-automation`.
5. Expand the **Repository** drop-down and select **Read public repositories**.
6. Select **Generate token**, then copy and save the token. You won't be able
   to retrieve it after closing the screen.
7. Sign in to `dhi.io` using your organization name as the username and the OAT
   as the password:

    ```console
    $ oras login dhi.io -u <YOUR_ORGANIZATION_NAME>
    ```

   Or non-interactively in a CI/CD pipeline:

   ```console
   $ echo $OAT | oras login dhi.io -u "$DOCKER_ORG" --password-stdin
   ```

8. Verify access by discovering attestations on a DHI image:

   ```console
   $ oras discover dhi.io/node:24-dev --platform linux/amd64
   ```

   > [!NOTE]
   >
   > The `--platform` flag is required. Without it, `oras discover` resolves to
   > the multi-arch image index, which returns only an index-level signature
   > rather than the full set of per-platform attestations.

   A successful response lists the attestations attached to the image,
   including SBOMs, provenance, vulnerability reports, and changelog metadata.

### Example mirroring with `regctl`

The following example shows how to mirror a specific tag of a Docker Hardened
Image from Docker Hub to another registry, along with its associated
attestations using `regctl`. You must [install
`regctl`](https://github.com/regclient/regclient) first.

The example assumes you have mirrored the DHI repository to your organization's
namespace on Docker Hub as described in the previous section. You can apply the
same steps to a non-mirrored image by updating the `SRC_ATT_REPO` and
`SRC_REPO` variables accordingly.

1. Set environment variables for your specific environment. Replace the
   placeholders with your actual values.

   In this example, you use a Docker username to represent a member of the Docker
   Hub organization that the DHI repositories are mirrored in. Prepare a
   [personal access token (PAT)](../../security/access-tokens.md) for the user
   with `read only` access. Alternatively, you can use your organization name and
   an [organization access token
   (OAT)](../../enterprise/security/access-tokens.md) to authenticate with `docker.io`.
   Note that OATs are not supported for `registry.scout.docker.com`. If your
   workflow requires authenticating to the Scout registry, use a personal access
   token (PAT) for that step.

   ```console
   $ export DOCKER_USERNAME="YOUR_DOCKER_USERNAME"
   $ export DOCKER_PAT="YOUR_DOCKER_PAT"
   $ export DOCKER_ORG="YOUR_DOCKER_ORG"
   $ export DEST_REG="registry.example.com"
   $ export DEST_REPO="mirror/dhi-python"
   $ export DEST_REG_USERNAME="YOUR_DESTINATION_REGISTRY_USERNAME"
   $ export DEST_REG_TOKEN="YOUR_DESTINATION_REGISTRY_TOKEN"
   $ export SRC_REPO="docker.io/${DOCKER_ORG}/dhi-python"
   $ export SRC_ATT_REPO="registry.scout.docker.com/${DOCKER_ORG}/dhi-python"
   $ export TAG="3.13-alpine3.21"
   ```

2. Sign in via `regctl` to Docker Hub, the Scout registry that contains
   the attestations, and your destination registry.

   ```console
   $ echo $DOCKER_PAT | regctl registry login -u "$DOCKER_USERNAME" --pass-stdin docker.io
   $ echo $DOCKER_PAT | regctl registry login -u "$DOCKER_USERNAME" --pass-stdin registry.scout.docker.com
   $ echo $DEST_REG_TOKEN | regctl registry login -u "$DEST_REG_USERNAME" --pass-stdin "$DEST_REG"
   ```

3. Mirror the image and attestations using `--referrers` and referrer endpoints:

   ```console
   $ regctl image copy \
        "${SRC_REPO}:${TAG}" \
        "${DEST_REG}/${DEST_REPO}:${TAG}" \
        --referrers \
        --referrers-src "${SRC_ATT_REPO}" \
        --referrers-tgt "${DEST_REG}/${DEST_REPO}" \
        --force-recursive
   ```

4. Verify that artifacts were preserved.

   First, get a digest for a specific tag and platform. For example, `linux/amd64`.

   ```console
   DIGEST="$(regctl manifest head "${DEST_REG}/${DEST_REPO}:${TAG}" --platform linux/amd64)"
   ```

   List attached artifacts (SBOM, provenance, VEX, vulnerability reports).

   ```console
   $ regctl artifact list "${DEST_REG}/${DEST_REPO}@${DIGEST}"
   ```

   Or, list attached artifacts with `docker scout`.

   ```console
   $ docker scout attest list "registry://${DEST_REG}/${DEST_REPO}@${DIGEST}"
   ```

### Example ongoing mirroring with `regsync`

`regsync` automates pulling from your organizations mirrored DHI repositories on
Docker Hub and pushing to your external registry including attestations. It
reads a YAML configuration file and can filter tags.

The following example uses a `regsync.yaml` file that syncs Node 24 and Python
3.12 Debian 13 variants, excluding Alpine and Debian 12.

```yaml{title="regsync.yaml"}
version: 1
# Optional: inline creds if not relying on prior CLI logins
# creds:
#   - registry: docker.io
#     user: <your-docker-username>
#     pass: "{{file \"/run/secrets/docker_token\"}}"
#   - registry: registry.scout.docker.com
#     user: <your-docker-username>
#     pass: "{{file \"/run/secrets/docker_token\"}}"
#   - registry: registry.example.com
#     user: <service-user>
#     pass: "{{file \"/run/secrets/dest_token\"}}"

sync:
  - source: docker.io/<your-org>/dhi-node
    target: registry.example.com/mirror/dhi-node
    type: repository
    fastCopy: true
    referrers: true
    referrerSource: registry.scout.docker.com/<your-org>/dhi-node
    referrerTarget: registry.example.com/mirror/dhi-node
    tags:
      allow: [ "24.*" ]
      deny: [ ".*alpine.*", ".*debian12.*" ]

  - source: docker.io/<your-org>/dhi-python
    target: registry.example.com/mirror/dhi-python
    type: repository
    fastCopy: true
    referrers: true
    referrerSource: registry.scout.docker.com/<your-org>/dhi-python
    referrerTarget: registry.example.com/mirror/dhi-python
    tags:
      allow: [ "3.12.*" ]
      deny: [ ".*alpine.*", ".*debian12.*" ]
```

To do a dry run with the configuration file, you can run the following command.
You must [install `regsync`](https://github.com/regclient/regclient) first.

```console
$ regsync check -c regsync.yaml
```

To run the sync with the configuration file:

```console
$ regsync once -c regsync.yaml
```

## What next

After mirroring, see [Pull a DHI](./use.md#pull-a-dhi) to learn how to pull and use mirrored images.
