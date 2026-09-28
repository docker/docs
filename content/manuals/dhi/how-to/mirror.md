---
title: Mirror a Docker Hardened Image repository
linktitle: Mirror a repository
description: Learn how to mirror an image into your organization's namespace and optionally push it to another private registry.
weight: 20
keywords: mirror docker image, private container registry, docker hub automation, webhook image sync, secure image distribution, internal registry, jfrog artifactory, harbor registry, amazon ecr, google artifact registry, github container registry, terraform, infrastructure as code
---

{{< summary-bar feature_name="Docker Hardened Images" >}}

Mirroring requires a DHI Select or Enterprise subscription. Without a
subscription, you can pull Docker Hardened Images directly from `dhi.io` without
mirroring. With a DHI Select or Enterprise subscription, you must mirror to your
organization to get:

- Compliance variants (FIPS-enabled or STIG-ready images)
- Extended Lifecycle Support (ELS) variants (requires add-on)
- Image or Helm chart customization
- Air-gapped or restricted network environments
- [SLA-backed security updates](https://docs.docker.com/go/dhi-sla/)

## How to mirror

This topic covers two types of mirroring for Docker Hardened Image (DHI)
repositories:

- [Mirror to your organization](#mirror-a-dhi-repository-to-your-organization):
  Mirror a DHI repository to your organization's namespace on Docker Hub.

- [Mirror to a third-party
  registry](#mirror-a-dhi-repository-to-a-third-party-registry): Mirror a
  repository to another container registry, such as Amazon ECR, Google Artifact
  Registry, or a private Harbor instance.

## Mirror a DHI repository to your organization

Organization owners, editors, and members with a [custom role](../../security/roles-and-permissions/custom-roles/_index.md)
that includes the DHI mirroring permission can create, view, and manage mirrors.
When using the CLI or Terraform, you can also mirror using an [organization
access token (OAT)](../../security/access-tokens/organization-access-tokens.md) with the
appropriate permission scopes, without requiring role-based access.

When a member with a custom role that includes the DHI mirroring permission
creates a mirror, Docker automatically creates and manages the
`dhi-mirroring-admins` team in your organization, adds that member to it, and
grants the team access to the new mirror. This lets the member manage mirrors
they create without organization owner or editor access. Mirrors created by
organization owners or editors don't use this team. Removing members from this
team may affect their ability to view and manage mirrors.

You can mirror image and chart repositories to your organization's namespace on
Docker Hub. Mirroring makes the repositories available within your organization
and lets you customize them for your environment:

- Image repositories: Mirroring lets you customize images by adding packages,
  OCI artifacts (such as custom certificates or additional tools), environment
  variables, labels, and other configuration settings. For more details, see
  [Customize a Docker Hardened Image](./customize.md#customize-a-docker-hardened-image).

- Chart repositories: Mirroring lets you customize image references within
  the chart. This is particularly useful when using customized images or when
  you've mirrored images to a third-party registry and need the chart to
  reference those custom locations. For more details, see [Customize a Docker
  Hardened Helm chart](./customize.md#customize-a-docker-hardened-helm-chart).

{{< tabs >}}
{{< tab name="Docker Hub" >}}

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

{{< /tab >}}
{{< tab name="CLI" >}}

Authenticate with `docker login` using your Docker credentials, a [personal
access token (PAT)](../../security/access-tokens/personal-access-tokens.md) with **Read & Write**
permissions, or an [organization access token
(OAT)](../../security/access-tokens/organization-access-tokens.md). When using an OAT, the
available operations depend on the token's permission scope:

- To list mirrored repositories, the OAT must have read (pull) access to the
  relevant repositories. Results are scoped to repositories the OAT can access.
- To create a mirror to an existing destination repository, the OAT must have
  push access to that repository. To create a mirror to a new destination
  repository that doesn't yet exist, the OAT must have org-wide repository
  access (for example, `<org>/*` with pull or push). Repository-scoped access to
  the future repository name is not sufficient.
- To stop mirroring, the OAT must have push access to the relevant repository.
- OATs with public repository read-only access cannot list or manage mirrored
  repositories.

Use the [`docker dhi mirror`](/reference/cli/docker/dhi/mirror/) command:

```console
$ docker dhi mirror start --org my-org \
  dhi/golang,my-org/dhi-golang \
  dhi/nginx,my-org/dhi-nginx \
  dhi/prometheus-chart,my-org/dhi-prometheus-chart
```

Mirror with dependencies:

```console
$ docker dhi mirror start --org my-org dhi/golang,my-org/dhi-golang --dependencies
```

List mirrored images in your organization:

```console
$ docker dhi mirror list --org my-org
```

Filter mirrored images by name or type:

```console
$ docker dhi mirror list --org my-org --filter python
$ docker dhi mirror list --org my-org --type image
$ docker dhi mirror list --org my-org --type helm-chart
```

{{< /tab >}}
{{< tab name="Terraform" >}}

You can manage DHI mirrors as infrastructure-as-code using the [DHI Terraform
provider](/dhi/tools/terraform/).

Define a `dhi_mirror` resource for each repository you want to mirror:

```hcl
resource "dhi_mirror" "golang" {
  source_namespace = "dhi"
  source_name      = "golang"
  destination_name = "dhi-golang"
}

resource "dhi_mirror" "nginx" {
  source_namespace = "dhi"
  source_name      = "nginx"
  destination_name = "dhi-nginx"
}
```

To enable Extended Lifecycle Support (ELS) variants, set the `els` attribute:

```hcl
resource "dhi_mirror" "golang" {
  source_namespace = "dhi"
  source_name      = "golang"
  destination_name = "dhi-golang"
  els              = true
}
```

Run `terraform apply` to create the mirrors.

For the full list of resource attributes, see the [Terraform Registry
documentation](https://registry.terraform.io/providers/docker-hardened-images/dhi/latest/docs/resources/mirror).

{{< /tab >}}
{{< /tabs >}}

After mirroring, the repository appears in your organization's repository list,
prefixed by `dhi-`, and continues to receive updated images. It behaves like any
other Docker Hub repository, so you can manage access and permissions, configure
webhooks, and use other standard Hub features. See [Docker Hub
repositories](/manuals/docker-hub/repos/_index.md) for details.

### Stop mirroring a repository

After you stop mirroring, the repository remains, but it no longer receives
updates. You can still use the last images or charts that were mirrored.

> [!NOTE]
>
> If you only want to stop mirroring ELS versions, you can clear the ELS
> option in the mirrored repository's **Settings** tab.

{{< tabs >}}
{{< tab name="Docker Hub" >}}

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub**.
3. In the namespace drop-down, select your organization that has access to DHI.
4. Select **Hardened Images** > **Manage**.
5. Select the **Mirrored Images** or **Mirrored Helm charts** tab.
6. In the far right column of the repository you want to stop mirroring, select the menu icon.
7. Select **Stop mirroring**.

{{< /tab >}}
{{< tab name="CLI" >}}

Authenticate with `docker login` using your Docker credentials, a [personal
access token (PAT)](../../security/access-tokens/personal-access-tokens.md) with **Read & Write**
permissions, or an [organization access token
(OAT)](../../security/access-tokens/organization-access-tokens.md) with push access to the
relevant repository.

Use the [`docker dhi mirror`](/reference/cli/docker/dhi/mirror/) command:

```console
$ docker dhi mirror stop --org my-org dhi-golang
```

{{< /tab >}}
{{< tab name="Terraform" >}}

To stop mirroring, remove the `dhi_mirror` resource from your Terraform
configuration and run `terraform apply`. The repository remains in your
organization but no longer receives updates.

{{< /tab >}}
{{< /tabs >}}

## Mirror a DHI repository to a third-party registry

After mirroring a DHI repository to your organization on Docker Hub, you can
optionally mirror it to another container registry, such as Amazon ECR, Google
Artifact Registry, GitHub Container Registry, or a private Harbor instance.

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

### Automate syncing with webhooks

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

When a webhook fires, Docker Hub sends the standard [webhook
payload](/docker-hub/repos/manage/webhooks/#example-webhook-payload). For a
mirrored DHI repository, the payload also includes an additional
`dhi_metadata` object. This object describes what changed between the newly
pushed build and the previous build of the same tag, including vulnerability
fixes, package changes, and configuration changes.

> [!NOTE]
>
> Docker Hub adds `dhi_metadata` only to pushes on mirrored DHI repositories.
> Webhooks on other repositories deliver the standard payload.

Each DHI build produces a signed changelog attestation. At webhook delivery
time, Docker Hub retrieves the changelog for the pushed image and embeds it in
the payload as `dhi_metadata`.

DHI changelogs are generated per architecture, so `dhi_metadata` is a map
keyed by the architecture-specific manifest digest. A multi-platform image
push contains an entry for each platform that has a changelog. Match the
digest key against the platform you care about instead of assuming a single
entry.

#### `dhi_metadata` fields

Each platform entry contains the following fields.

| Field | Type | Description |
| :---- | :---- | :---- |
| `schema_version` | integer | Version of the `dhi_metadata` schema. |
| `change_categories` | array of strings | High-level summary of what changed in this build. See [Change categories](#change-categories). |
| `previous_version` | object | The prior build this one is compared against. Contains `tag` and `digest`. |
| `changes` | object | Detailed diff versus the previous version. See the following table. |

The `changes` object contains:

| Field | Type | Description |
| :---- | :---- | :---- |
| `vulnerabilities_fixed` | array | CVEs resolved in this build. Each entry has `cve_id`, `severity`, `package`, and `fixed_in_version`. |
| `packages_updated` | array | Packages whose version changed. Each entry has `name`, `type`, `old_version`, and `new_version`. |
| `packages_added` | array | Packages added in this build. Each entry has `name`, `type`, and `version`. |
| `packages_removed` | array | Packages removed in this build. Each entry has `name`, `type`, and `version`. |
| `environment_variables_changed` | array | Changes to environment variables. Each entry has `change`, `key`, and `from_value` or `to_value` as applicable. |
| `labels_changed` | array | Changes to image labels, in the same shape as environment variable changes. |
| `configuration_changed` | array | Changes to other image configuration. For example, the entrypoint. |

When a change type has no entries, its array is present but empty, shown as `[]`.

#### Change categories

`change_categories` gives a quick, machine-readable summary of the build.

| Value | Meaning |
| :---- | :---- |
| `vulnerability_fix` | The build resolves one or more CVEs. See `changes.vulnerabilities_fixed`. |
| `version_upgrade` | One or more packages changed version. See `changes.packages_updated`. |
| `other` | The build has package, environment variable, label, or configuration changes that don't fall into either category above. |

A build can have more than one category. For example, a build that fixes a CVE
and also bumps a package version returns both `vulnerability_fix` and
`version_upgrade`. A build with no changes at all returns an empty array.

#### Example: vulnerability fix and version upgrade

The following excerpt shows the `dhi_metadata` object from a webhook payload
for a push to a mirrored DHI repository. The example is trimmed to a single
platform and a subset of changes for readability. A real payload contains one
`dhi_metadata` entry per architecture.

```json {collapse=true}
{
  ...
  "dhi_metadata": {
    "sha256:04639747b6d72bcf1d0322f2a5b122ee76d963e31bb4a070891b25f15a5001c5": {
      "schema_version": 1,
      "change_categories": ["vulnerability_fix", "version_upgrade"],
      "previous_version": {
        "tag": "2-compat-fips-dev",
        "digest": "sha256:1738aa35838f520431c898b85d7cd60da71d8f997965287db4f3be27c1df32a1"
      },
      "changes": {
        "vulnerabilities_fixed": [
          {
            "cve_id": "CVE-2019-9192",
            "severity": "low",
            "package": "glibc",
            "fixed_in_version": "2.41-12+deb13u4+dhi0"
          },
          {
            "cve_id": "CVE-2018-20796",
            "severity": "low",
            "package": "glibc",
            "fixed_in_version": "2.41-12+deb13u4+dhi0"
          }
        ],
        "packages_updated": [
          {
            "name": "glibc",
            "type": "deb",
            "old_version": "2.41-12+deb13u4",
            "new_version": "2.41-12+deb13u4+dhi0"
          },
          {
            "name": "libc6",
            "type": "deb",
            "old_version": "2.41-12+deb13u4",
            "new_version": "2.41-12+deb13u4+dhi0"
          }
        ],
        "packages_added": [],
        "packages_removed": [],
        "environment_variables_changed": [],
        "labels_changed": [
          {
            "change": "changed",
            "key": "com.docker.dhi.chain-id",
            "from_value": "sha256:4567092c648d813b8c4c60c7d100fc34df817dd5cb4c7968e9a5c43bafb9e7a5",
            "to_value": "sha256:62d4e2090951e812a87fb599db362677f72dee095f85889ea56df63c0999b02a"
          }
        ],
        "configuration_changed": []
      }
    }
  }
}
```

#### Example: version bump with no CVEs

When a build only bumps package versions, `change_categories` contains
`version_upgrade` and `vulnerabilities_fixed` is empty.

```json {collapse=true}
{
  ...
  "dhi_metadata": {
    "sha256:2982980b6bb3cdedafa9377bcc37405c20ed48702deef11faf13ec99d596057d": {
      "schema_version": 1,
      "change_categories": ["version_upgrade"],
      "previous_version": {
        "tag": "5-fips-dev",
        "digest": "sha256:81355a1301ecc5f78dd87b68a284642d7b6bfbd86f3a37f3932fad7ecf1141e6"
      },
      "changes": {
        "vulnerabilities_fixed": [],
        "packages_updated": [
          {
            "name": "sqlite3",
            "type": "deb",
            "old_version": "3.46.1-7+deb13u2+dhi0",
            "new_version": "3.46.1-7+deb13u2+dhi1"
          }
        ],
        "packages_added": [],
        "packages_removed": [],
        "environment_variables_changed": [],
        "labels_changed": [],
        "configuration_changed": []
      }
    }
  }
}
```

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

   In this example, you authenticate as your Docker organization using an
   [organization access token
   (OAT)](../../security/access-tokens/organization-access-tokens.md). The OAT must have at
   least pull access to every DHI repository you want to mirror. Only
   repositories in the token's scope are accessible. Alternatively, you can
   authenticate as a Docker Hub user with a [personal access token
   (PAT)](../../security/access-tokens/personal-access-tokens.md) that has `read only` access.

   > [!WARNING]
   >
   > The following examples export credentials directly on the command line for
   > demonstration purposes. This exposes sensitive tokens in your shell history
   > and process list. In production environments, use secure methods such as
   > reading from files with restricted permissions, environment files loaded
   > at runtime, or secret management tools.

   ```console
   $ export DOCKER_ORG="YOUR_DOCKER_ORG"
   $ export DOCKER_OAT="YOUR_DOCKER_OAT"
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
   $ echo $DOCKER_OAT | regctl registry login -u "$DOCKER_ORG" --pass-stdin docker.io
   $ echo $DOCKER_OAT | regctl registry login -u "$DOCKER_ORG" --pass-stdin registry.scout.docker.com
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

```yaml{title="regsync.yaml",collapse=true}
version: 1
# Optional: inline creds if not relying on prior CLI logins
# creds:
#   - registry: docker.io
#     user: <your-docker-org>
#     pass: "{{file \"/run/secrets/docker_oat\"}}"
#   - registry: registry.scout.docker.com
#     user: <your-docker-org>
#     pass: "{{file \"/run/secrets/docker_oat\"}}"
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
