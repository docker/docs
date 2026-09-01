---
title: Integrate Docker Scout with a container registry
linkTitle: Container registries
description: Integrate Docker Scout with any container registry using the docker scout watch CLI command
keywords: docker scout, registry, integration, image analysis, security, cves, watch, ecr, acr, artifactory, harbor, nexus
aliases:
  - /scout/integrations/registry/artifactory/
  - /scout/artifactory/
  - /scout/integrations/registry/ecr/
  - /scout/integrations/registry/acr/
---

[`docker scout watch`](/reference/cli/docker/scout/watch/) is a long-running
CLI process that indexes images from a container registry and pushes the
results to Docker Scout. It works with any Docker/OCI-compliant registry,
including Amazon ECR, Azure Container Registry, JFrog Artifactory, Harbor, and
Sonatype Nexus.

## How it works

You run `docker scout watch` on a host you control. The process can:

- Watch specific repositories or an entire registry
- Optionally ingest all existing images once, using `--all-images`
- Periodically refresh repository lists, using `--refresh-registry`
- Receive webhook callbacks from registries that support them, for
  near-real-time analysis instead of polling

After the integration, Docker Scout automatically pulls and analyzes images
that you push to the registry. Metadata about your images are stored on the
Docker Scout platform, but Docker Scout doesn't store the container images
themselves. For more information about how Docker Scout handles image data,
see [Data handling](/manuals/scout/deep-dive/data-handling.md).

## Set up `docker scout watch`

1. Pick a host on which to run `docker scout watch`.

   The host must have network access to your registry and be able to access
   the Scout API (`https://api.scout.docker.com`) over the internet. If
   you're using webhook callbacks, the registry must also be able to reach the
   `docker scout watch` host on the configured port.

2. Ensure you are running the latest version of Scout.

   ```console
   $ docker scout version
   ```

   If necessary, [install the latest version of Scout](/manuals/scout/install.md).

3. Authenticate Docker to your registry.

   ```console
   $ docker login <registry-hostname> --username <user> --password <password-or-access-token>
   ```

   For Amazon ECR, authenticate using the AWS CLI instead:

   ```console
   $ aws ecr get-login-password --region <region> | \
     docker login --username AWS --password-stdin \
     <aws_account_id>.dkr.ecr.<region>.amazonaws.com
   ```

   The AWS identity used must have at least `ecr:GetAuthorizationToken` and
   `ecr:BatchGetImage` permissions on the target registry.

   For Azure Container Registry:

   ```console
   $ docker login <registry-name>.azurecr.io \
     --username <username> \
     --password <password-or-access-token>
   ```

   > [!TIP]
   >
   > As a best practice, use a dedicated user or token with read-only access
   > to the registry.

4. Set up your Scout credentials.

   1. Generate an organization access token. For more details, see
      [Create an organization access token](/enterprise/security/access-tokens/#create-an-organization-access-token).
   2. Sign in to Docker using the organization access token.

      ```console
      $ docker login --username <your_organization_name>
      ```

      When prompted for a password, paste the organization access token.

   3. Connect your local Docker environment to your organization's Docker Scout service.

      ```console
      $ docker scout enroll <your_organization_name>
      ```

5. Index existing images. You only need to do this once.

   ```console
   $ docker scout watch --registry <registry-hostname> --all-images
   ```

6. Confirm the images have been indexed by viewing them on the
   [Scout Dashboard](https://scout.docker.com/).

7. Continuously watch for new images.

   ```console
   $ docker scout watch --registry <registry-hostname> --refresh-registry
   ```

   `docker scout watch` is a long-running process. Run it as a system
   service, for example using `systemd` or `nohup`, to ensure it continues
   running in the background. Use `--interval` (default 60 seconds) to
   control polling frequency, and `--repository` and `--tag` to narrow scope.

Reference: [`docker scout watch`](/reference/cli/docker/scout/watch/)

## Registry-specific options

Some registries need extra configuration beyond a hostname, passed through
the `--registry` flag as a `key=value` string, for example a REST API
endpoint for webhook callbacks, or a non-standard repository layout. Built-in
adapters exist for `type=artifactory`, `type=harbor`, and `type=nexus`, and a
`type=oci` adapter covers any OCI-compliant registry that implements the
`_catalog` endpoint. For the full option reference for each type, see
[`docker scout watch`](/reference/cli/docker/scout/watch/).

The following example walks through the `type=artifactory` adapter in detail.
See the CLI reference for equivalent Harbor, Nexus, and generic OCI examples.

### Example: JFrog Artifactory

These `type=artifactory` options override the generic registry handling for
the `--registry` option:

| Key              | Required | Description                                                                            |
|------------------|:--------:|----------------------------------------------------------------------------------------|
| `type`           |   Yes    | Must be `artifactory`.                                                                 |
| `registry`       |   Yes    | Docker/OCI registry hostname (e.g., `example.jfrog.io`).                               |
| `api`            |   Yes    | Artifactory REST API base URL (e.g., `https://example.jfrog.io/artifactory`).          |
| `repository`     |   Yes    | Repository to watch (replaces `--repository`).                                         |
| `includes`       |    No    | Globs to include (e.g., `*/frontend*`).                                                |
| `excludes`       |    No    | Globs to exclude (e.g., `*/legacy/*`).                                                 |
| `port`           |    No    | Local port to listen on for webhook callbacks.                                         |
| `subdomain-mode` |    No    | `true` or `false`; matches Artifactory's Docker layout (subdomain versus repository-path). |

Set up credentials for the Scout client to authenticate with Artifactory, and
a secret for Artifactory to authenticate its webhook callbacks:

```console
$ export DOCKER_SCOUT_ARTIFACTORY_API_USER=<user>
$ export DOCKER_SCOUT_ARTIFACTORY_API_PASSWORD=<password-or-access-token>
$ export DOCKER_SCOUT_ARTIFACTORY_WEBHOOK_SECRET=<random-64-128-character-secret>
```

> [!TIP]
>
> As a best practice, create a dedicated user with read-only access and use an
> access token instead of a password. Generate the webhook secret as a
> high-entropy random string of 64-128 characters.

Index existing images with the Artifactory-specific registry string:

```console
$ docker scout watch --registry \
  "type=artifactory,registry=example.jfrog.io,api=https://example.jfrog.io/artifactory,include=*/frontend*,exclude=*/dta/*,repository=docker-local,port=9000,subdomain-mode=true" \
  --all-images
```

Then configure Artifactory to call the webhook: in your Artifactory UI or via
REST API, set up a webhook for image push/update events, pointing to your
`docker scout watch` host and port, and include the
`DOCKER_SCOUT_ARTIFACTORY_WEBHOOK_SECRET` for authentication. For more
information, see the [JFrog Artifactory Webhooks
documentation](https://jfrog.com/help/r/jfrog-platform-administration-documentation/webhooks)
or the [JFrog Artifactory REST API Webhooks
documentation](https://jfrog.com/help/r/jfrog-rest-apis/webhooks).

Finally, run the same command with `--refresh-registry` instead of
`--all-images` as your long-running watch process, so new images are picked up
going forward:

```console
$ docker scout watch --registry \
  "type=artifactory,registry=example.jfrog.io,api=https://example.jfrog.io/artifactory,include=*/frontend*,exclude=*/dta/*,repository=docker-local,port=9000,subdomain-mode=true" \
  --refresh-registry
```
