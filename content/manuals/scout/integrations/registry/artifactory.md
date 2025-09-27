---
description: Integrate Artifactory Container Registry with Docker Scout
keywords: docker scout, artifactory, integration, image analysis, security, cves
title: Integrate Docker Scout with Artifactory Container Registry
linkTitle: Artifactory Container Registry
---

{{% experimental %}}

The `docker scout watch` command is experimental.

Experimental features are intended for testing and feedback as their
functionality or design may change between releases without warning or
can be removed entirely in a future release.

{{% /experimental %}}

Integrating Docker Scout with JFrog Artifactory lets you index and analyze
images from Artifactory. This integration is powered by a long-running
`docker scout watch` process. It pulls images from your selected repositories
(optionally filtered), can receive webhook callbacks from Artifactory, and
pushes image data to Docker Scout. View results in the Docker Scout Dashboard or
with `docker scout` CLI.

## How it works

You run [`docker scout watch`](/reference/cli/docker/scout/watch/) on a host you
control and configure the Artifactory-specific registry string via `--registry
"key=value,..."`. The watch process can:

- Watch specific repositories or an entire registry
- Optionally ingest all existing images once
- Periodically refresh repository lists
- Receive webhook callbacks from Artifactory on a local port you choose

After the integration, Docker Scout automatically pulls and analyzes images
that you push to the Artifactory registry. Metadata about your images are stored on the
Docker Scout platform, but Docker Scout doesn't store the container images
themselves. For more information about how Docker Scout handles image data, see
[Data handling](/manuals/scout/deep-dive/data-handling.md).

### Artifactory-specific registry string options

These `type=artifactory` options override the generic registry handling for the `--registry` option:

| Key              | Required | Description                                                                            |
|------------------|:--------:|----------------------------------------------------------------------------------------|
| `type`           |   Yes    | Must be `artifactory`.                                                                 |
| `registry`       |   Yes    | Docker/OCI registry hostname (e.g., `example.jfrog.io`).                               |
| `api`            |   Yes    | Artifactory REST API base URL (e.g., `https://example.jfrog.io/artifactory`).          |
| `repository`     |   Yes    | Repository to watch (replaces `--repository`).                                         |
| `includes`       |    No    | Globs to include (e.g., `*/frontend*`).                                                |
| `excludes`       |    No    | Globs to exclude (e.g., `*/legacy/*`).                                                 |
| `port`           |    No    | Local port to listen on for webhook callbacks.                                         |
| `subdomain-mode` |    No    | `true` or `false`; matches Artifactory’s Docker layout (subdomain versus repository-path). |

## Integrate an Artifactory registry

Use the following steps to integrate your Artifactory registry with Docker
Scout.

1. Pick the host on which to run `docker scout watch`.

   The host must have local or network access to your private registry and be able
   to access the Scout API (`https://api.scout.docker.com`) over the internet. If
   you're using webhook callbacks, Artifactory must also be able to reach the Scout
   client host on the configured port.
   Override the `--workers` option (default: `3`) for optimal performance based on
   the size of the host and the expected workload.

2. Ensure you are running the latest version of Scout.

   Check your current version:

   ```console
   $ docker scout version
   ```

   If necessary, [install the latest version of Scout](https://docs.docker.com/scout/install/).

3. Set up your Artifactory credentials.

   Store the credentials that the Scout client will use to authenticate with
   Artifactory. The following is an example using environment variables. Replace
   `<user>` and `<password-or-access-token>` with your actual values.

   ```console
   $ export DOCKER_SCOUT_ARTIFACTORY_API_USER=<user>
   $ export DOCKER_SCOUT_ARTIFACTORY_API_PASSWORD=<password-or-access-token>
   ```

   > [!TIP]
   >
   > As a best practice, create a dedicated user with read-only access and use
   > an access token instead of a password.

   Store the credential that Artifactory will use to authenticate webhook
   callbacks. The following is an example using an environment variable. Replace
   `<random-64-128-character-secret>` with an actual secret.

   ```console
   $ export DOCKER_SCOUT_ARTIFACTORY_WEBHOOK_SECRET=<random-64-128-character-secret>
   ````

   > [!TIP]
   >
   > As a best practice, generate a high-entropy random string of 64-128 characters.

4. Set up your Scout credentials.

   1. Generate an organization access token for accessing Scout. For more
      details, see [Create an organization access
      token](/enterprise/security/access-tokens/#create-an-organization-access-token).
   2. Sign in to Docker using the organization access token.

       ```console
       $ docker login --username <your_organization_name>
       ```

       When prompted for a password, paste the organization access token you
       generated.

   3. Connect your local Docker environment to your organization's Docker Scout service.

       ```console
       $ docker scout enroll <your_organization_name>
       ```

5. Index existing images. You only need to do this once.

    Run `docker scout watch` with the `--all-images` option to index all images in the specified Artifactory repository. The following is an example command:

   ```console
   $ docker scout watch --registry \
   "type=artifactory,registry=example.jfrog.io,api=https://example.jfrog.io/artifactory,include=*/frontend*,exclude=*/dta/*,repository=docker-local,port=9000,subdomain-mode=true" \
   --all-images
   ```

6. Confirm the images have been indexed by viewing them on the [Scout
   Dashboard](https://scout.docker.com/).

7. Configure Artifactory callbacks.

   In your Artifactory UI or via REST API, configure webhooks for image
   push/update events. Set the endpoint to your `docker scout watch` host and
   port, and include the `DOCKER_SCOUT_ARTIFACTORY_WEBHOOK_SECRET` for
   authentication.

   For more information, see the [JFrog Artifactory Webhooks
   documentation](https://jfrog.com/help/r/jfrog-platform-administration-documentation/webhooks)
   or the [JFrog Artifactory REST API Webhooks
   documentation](https://jfrog.com/help/r/jfrog-rest-apis/webhooks).

8. Continuously watch for new or updated images.

   Run `docker scout watch` with the `--refresh-registry` option to watch for
   new images to index. The following is an example command:

   ```console
   $ docker scout watch --registry \
   "type=artifactory,registry=example.jfrog.io,api=https://example.jfrog.io/artifactory,include=*/frontend*,exclude=*/dta/*,repository=docker-local,port=9000,subdomain-mode=true" \
   --refresh-registry
   ```

9. Optional. Set up Scout integration for real-time notifications from popular
   collaboration platforms. For details, see [Integrate Docker Scout with
   Slack](../team-collaboration/slack.md).