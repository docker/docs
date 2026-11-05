---
title: Use a Docker Hardened Image
linktitle: Use an image
description: Learn how to pull, run, and reference Docker Hardened Images in Dockerfiles, CI pipelines, and standard development workflows.
keywords: use hardened image, docker pull secure image, non-root containers, multi-stage dockerfile, dev image variant
weight: 30
aliases:
  - /dhi/how-to/els/
  - /dhi/how-to/k8s/
---

You can use a Docker Hardened Image (DHI) just like any other image on Docker
Hub. DHIs follow the same familiar usage patterns. Pull them with `docker pull`,
reference them in your Dockerfile, and run containers with `docker run`.

The key difference is that DHIs are security-focused and intentionally minimal
to reduce the attack surface. This means some variants don't include a shell or
package manager, and may run as a non-root user by default.

> [!IMPORTANT]
>
> You must authenticate to the Docker Hardened Images registry (`dhi.io`) to
> pull DHI Community images. You can authenticate using either of the following:
>
> - **Docker ID and password:** Use your Docker Hub username and password. If
>   you don't have a Docker account, [create one](../../accounts/individual/create-account.md)
>   for free.
> - **Access token:** Use a [personal access token
>   (PAT)](../../security/access-tokens/personal-access-tokens.md) for personal accounts, or an
>   [organization access token
>   (OAT)](../../security/access-tokens/organization-access-tokens.md) with your organization
>   name as the username.
> - OIDC: In GitHub Actions, use a Docker-issued access token from an
>   [OIDC connection](/manuals/security/authentication/oidc-connections/create-manage.md)
>   with your organization name as the username.
>
> Run `docker login dhi.io` to authenticate.

## Considerations when adopting DHIs

Docker Hardened Images are intentionally minimal to improve security. If you're
updating existing Dockerfiles or frameworks to use DHIs, keep in mind that
runtime images don't include shells or package managers, run as non-root users
by default, and may have different configurations than images you're familiar
with.

For a comprehensive checklist of migration considerations and detailed guidance,
see [Migrate to Docker Hardened Images](../migration/_index.md).

## Pull, run, and reference DHIs

Docker Hardened Images use different image references depending on your
subscription:

| Subscription        | Image reference            | Authentication        |
|---------------------|----------------------------|-----------------------|
| Community           | `dhi.io/<image>:<tag>`     | `docker login dhi.io` |
| Select & Enterprise | `<your-org>/<image>:<tag>` | `docker login`        |

Select and Enterprise users should [mirror](./mirror.md) repositories to their
Docker Hub organization to access compliance variants and customization
features.

After authenticating, use the image reference in standard Docker commands and
Dockerfiles. For example:

```console
$ docker pull dhi.io/python:3.13
$ docker run --rm dhi.io/python:3.13 python -c "print('Hello from DHI')"
```

```dockerfile
FROM dhi.io/python:3.13
COPY . /app
CMD ["python", "/app/main.py"]
```

For multi-stage builds:
- Use a `-dev` tag for build stages that need a shell or package manager. See
  [Use dev variants for framework-based
  applications](#use-dev-variants-for-framework-based-applications).
- Use the `static` image for compiled executables with minimal runtime
  dependencies. See [Use a static image for compiled
  executables](#use-a-static-image-for-compiled-executables).

To learn how to search for available variants, see [Search and evaluate
images](./search-evaluate.md).

## Use a DHI in CI/CD pipelines

Docker Hardened Images work just like any other image in your CI/CD pipelines.
You can reference them in Dockerfiles, pull them as part of a pipeline step, or
run containers based on them during builds and tests.

In GitHub Actions, you can use an OIDC connection to authenticate with
short-lived Docker access tokens instead of storing a PAT or OAT. For a
mirrored image on Docker Hub, grant pull access to
`<your-organization>/<repository>`. Your organization's DHI subscription still
determines access to paid image variants.

See [Authenticate with OIDC in GitHub Actions](#authenticate-with-oidc-in-github-actions)
for an example that reads an image manifest and its private attestations.

Unlike typical container images, DHIs also include signed
[attestations](../explore/security-concepts/attestations.md) such as SBOMs and provenance
metadata. You can incorporate these into your pipeline to support supply chain
security, policy checks, or audit requirements if your tooling supports it.

To strengthen your software supply chain, consider adding your own attestations
when building images from DHIs. This lets you document how the image was built,
verify its integrity, and enable downstream validation and policy enforcement
using tools like Docker Scout.

To learn how to attach attestations during the build process, see [Docker Build
Attestations](/manuals/build/metadata/attestations.md).

### Authenticate with OIDC in GitHub Actions

Use an [OIDC connection](/manuals/security/authentication/oidc-connections/create-manage.md)
to read DHI images and attestations with short-lived Docker access tokens.
Use your Docker organization name as the username and the Docker-issued access
token as the password. The GitHub identity token is exchanged for this access
token. Don't use the GitHub identity token as a registry password.

Configure the connection's [rulesets and subject claims](/manuals/security/authentication/oidc-connections/rulesets-claims.md)
to match your workflow. Repository access depends on the registry you use:

- For free images and their attestations served through `dhi.io`, authenticate
  with your organization name and Docker-issued access token.
- For a mirrored or customized image on Docker Hub, grant pull access to
  `<your-organization>/<repository>`.
- Direct reads from `registry.scout.docker.com` require access to the referenced
  repository and use your organization as the Scout registry namespace.

For permission boundaries, see
[DHI and Scout compatibility exceptions](/manuals/security/authentication/oidc-connections/rulesets-claims.md#dhi-and-scout-compatibility-exceptions).

Repository grants can cover an exact repository, all repositories in a namespace
(`namespace/*`), or public Docker Hub repositories (`*/*/public`, with pull access
only). Public access applies across namespaces. It requires the repository to be
public when each read is authorized. It doesn't grant access to private
repositories or remove DHI subscription, FIPS, or Extended Lifecycle Support
requirements.

Grants are additive. Push access includes pull access. An exact pull grant cannot
restrict a matching namespace-wide push grant. To restrict a repository to
pull-only access, remove every matching push grant. The `dhi.io` endpoint accepts
read operations with OIDC credentials. Push grants also authorize supported
uploads to `registry.scout.docker.com`, but don't authorize deletion.

The repository grant also controls direct reads of the image's private
attestations from `registry.scout.docker.com`. Your organization's DHI
subscription also determines access to paid image variants.
Managing customization definitions requires the destination-repository permissions
described in [Customize a DHI](./customize.md).

The following workflow uses `regctl` to read a mirrored image manifest and list
its attestations. Set these GitHub Actions repository variables:

- `DOCKER_ORG`: Your Docker organization name
- `DOCKER_OIDC_CONNECTION_ID`: Your OIDC connection ID
- `DHI_REPOSITORY`: The mirrored repository name, such as `dhi-node`
- `DHI_TAG`: The image tag, such as `22`

Grant the connection pull access to `DOCKER_ORG/DHI_REPOSITORY`, and ensure its
subject claims match this workflow's `workflow_dispatch` run on your selected
branch. The exchange requests a token with a five-minute lifetime. Run the
registry operations before the token expires. Longer jobs must obtain a fresh
token before further registry operations.

```yaml
name: Read DHI attestations
on:
  workflow_dispatch:

permissions:
  contents: read
  id-token: write

jobs:
  verify:
    runs-on: ubuntu-latest
    env:
      DOCKER_ORG: ${{ vars.DOCKER_ORG }}
      DOCKER_OIDC_CONNECTION_ID: ${{ vars.DOCKER_OIDC_CONNECTION_ID }}
      DHI_REPOSITORY: ${{ vars.DHI_REPOSITORY }}
      DHI_TAG: ${{ vars.DHI_TAG }}
    steps:
      - name: Install regctl
        uses: regclient/actions/regctl-installer@v0.1.0
        with:
          release: v0.8.2
      - name: Exchange the GitHub identity token and sign in
        uses: actions/github-script@v8
        with:
          script: |
            const identityToken = await core.getIDToken('https://identity.docker.com');
            const response = await fetch('https://identity.docker.com/oauth/token', {
              method: 'POST',
              headers: {'Content-Type': 'application/x-www-form-urlencoded'},
              body: new URLSearchParams({
                grant_type: 'urn:ietf:params:oauth:grant-type:token-exchange',
                subject_token_type: 'urn:ietf:params:oauth:token-type:id_token',
                subject_token: identityToken,
                connection_id: process.env.DOCKER_OIDC_CONNECTION_ID,
                expires_in: '300'
              })
            });
            if (!response.ok) {
              throw new Error(`Docker token exchange failed: HTTP ${response.status}`);
            }
            const {access_token: token} = await response.json();
            if (typeof token !== 'string' || !token) {
              throw new Error('Docker token exchange returned no access token');
            }
            core.setSecret(token);
            for (const registry of ['docker.io', 'registry.scout.docker.com']) {
              await exec.getExecOutput('regctl', [
                'registry', 'login', registry,
                '--user', process.env.DOCKER_ORG, '--pass-stdin'
              ], {input: token, silent: true});
            }
      - name: Read the manifest and list attestations
        run: |
          regctl manifest get "docker.io/${DOCKER_ORG}/${DHI_REPOSITORY}:${DHI_TAG}"
          regctl artifact list "docker.io/${DOCKER_ORG}/${DHI_REPOSITORY}:${DHI_TAG}" \
            --external "registry.scout.docker.com/${DOCKER_ORG}/${DHI_REPOSITORY}" \
            --platform linux/amd64
      - name: Sign out
        if: always()
        run: |
          regctl registry logout docker.io
          regctl registry logout registry.scout.docker.com
```

For a free image served through `dhi.io`, sign in to `dhi.io` and use
`dhi.io/<image>:<tag>` for both commands. Omit `--external` so that image and
attestation reads both use `dhi.io`. Sign out of `dhi.io` when finished.

### Discover attestations with ORAS

You can use [ORAS](https://oras.land/) to discover and inspect the attestations
attached to Docker Hardened Images. This is particularly useful in CI/CD
pipelines for supply chain security validation and compliance checks.

For automated workflows, authenticate using an [organization access token
(OAT)](../../security/access-tokens/organization-access-tokens.md). OATs are owned by the
organization rather than an individual user, making them better suited for CI/CD
pipelines.

To discover attestations with ORAS:

1. [Generate an organization access
   token](../../security/access-tokens/organization-access-tokens.md) with **Read public
   repositories** scope.

   The following example shows how to discover attestations on DHI community
   images from `dhi.io`. If you're discovering attestations on images mirrored to
   your organization, generate an OAT scoped to read from your mirrored repository
   instead of **Read public repositories**.

2. Sign in to `dhi.io` using your organization name as the username and the OAT
   as the password.

   > [!WARNING]
   >
   > The following examples export credentials directly on the command line for
   > demonstration purposes. This exposes sensitive tokens in your shell history
   > and process list. In production environments, use secure methods such as
   > reading from files with restricted permissions, environment files loaded
   > at runtime, or secret management tools.

    ```console
    $ oras login dhi.io -u <YOUR_ORGANIZATION_NAME>
    ```

   Or non-interactively in a CI/CD pipeline, set your organization name and token:

   ```console
   $ export DOCKER_ORG="YOUR_ORGANIZATION_NAME"
   $ export OAT="YOUR_ORGANIZATION_ACCESS_TOKEN"
   $ echo $OAT | oras login dhi.io -u "$DOCKER_ORG" --password-stdin
   ```

3. Discover attestations on a DHI image:

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

## Use a static image for compiled executables

Docker Hardened Images include a `static` image repository designed specifically
for running compiled executables in an extremely minimal and secure runtime.
Unlike a non-hardened `FROM scratch` image, the DHI `static` image includes
attestations and essential packages like `ca-certificates`.

Use a `-dev` or other builder image to compile your binary, then copy the output
into a `static` image:

```dockerfile
FROM dhi.io/golang:1.22-dev AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o myapp

FROM dhi.io/static:20230311
COPY --from=build /app/myapp /myapp
ENTRYPOINT ["/myapp"]
```

For more multi-stage build patterns, see the [Go migration
example](../migration/examples/go.md).

## Use dev variants for framework-based applications

If you're building applications with frameworks that require package managers or
build tools (such as Python, Node.js, or Go), use a `-dev` variant during the
development or build stage. These variants include essential utilities like
shells, compilers, and package managers to support local iteration and CI
workflows.

Use `-dev` images in your inner development loop or in isolated CI stages to
maximize productivity. Once you're ready to produce artifacts for production,
switch to a smaller runtime variant to reduce the attack surface and image size.

For detailed multi-stage Dockerfile examples using dev variants, see the
migration examples:
- [Go](../migration/examples/go.md)
- [Python](../migration/examples/python.md)
- [Node.js](../migration/examples/node.md)

## Use Socket Firewall variants to monitor package installations

If you want supply chain protection during dependency installation, use a Socket
Firewall variant in place of the standard `-dev` variant in your build stage.
These variants come with [Socket](https://socket.dev/) preinstalled to monitor
package manager activity and block malicious packages before they reach your
image.

Two tiers are available. Use `-sfw-dev` for Socket Firewall Free, or
`-sfw-ent-dev` for Socket Firewall Enterprise (requires an API key from Socket).
The runtime stage stays the same regardless of which build stage variant you
use.

```dockerfile
FROM dhi.io/python:3.13-alpine3.23-sfw-dev AS build
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

FROM dhi.io/python:3.13-alpine3.23
COPY --from=build /app /app
CMD ["python", "app.py"]
```

For more information on Socket Firewall variants, see [Available image
types](../explore/available.md).

## Use compliance and ELS variants

{{< summary-bar feature_name="Docker Hardened Images" >}}

With a DHI Select or DHI Enterprise subscription, you can access additional
image variants:

- Compliance variants: FIPS-enabled and STIG-ready images for regulatory
  requirements
- ELS (Extended Lifecycle Support) variants (requires add-on): Security patches
  for end-of-life image versions

To access these variants, [mirror](./mirror.md) the repository to your Docker
Hub organization. For ELS, enable **Mirror end-of-life images** when setting up
mirroring. Once mirrored, use the compliance or EOL tags like any other image
tag.

## Use with Kubernetes

When deploying Docker Hardened Images to Kubernetes, the process is similar to
using any other container image with one key difference: you must configure
image pull secrets to authenticate to the DHI registry. This applies whether
you're pulling directly from `dhi.io`, from a mirror on Docker Hub, or from
your own third-party registry.

### Create an image pull secret

You can create an image pull secret using either an access token or Docker Desktop credentials.

For the `--docker-server` value:
- Use `dhi.io` for community images pulled directly from Docker Hardened Images
- Use `docker.io` for mirrored repositories on Docker Hub
- Use your registry's hostname for third-party registries

#### Using an access token

Create a secret using a [Personal Access Token
(PAT)](../../security/access-tokens/personal-access-tokens.md) or [Organization Access Token
(OAT)](../../security/access-tokens/organization-access-tokens.md). Ensure the token has at
least read-only access to the repositories.

```console
$ kubectl create -n <kubernetes namespace> secret docker-registry <secret name> --docker-server=<registry server> \
        --docker-username=<registry user> --docker-password=<access token> \
        --docker-email=<registry email>
```

#### Using Docker Desktop credentials

If you're already authenticated with Docker Desktop, you can create a secret
using your stored credentials. This method works for registries you've
authenticated to via Docker Desktop (using `docker login <registry>`).

```console
$ NS=<namespace>
$ kubectl create -n ${NS} secret docker-registry dhi-pull-secret \
    --docker-server=<registry server> \
    --docker-username=<registry user> \
    --docker-password="$(echo https://<registry server> | docker-credential-desktop get | jq -r .Secret)" \
    --docker-email=<registry email>
```

This method extracts credentials from Docker Desktop's credential store, avoiding the need to create a separate access token for local development.

### Test the image pull secret

After creating the secret, verify it works by deploying a test pod that
references the secret in its `imagePullSecrets` configuration.

Create a test pod:

```console
kubectl apply --wait -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: dhi-test
  namespace: <kubernetes namespace>
spec:
  containers:
  - name: test
    image: bash:5
    command: [ "sh", "-c", "echo 'Hello from DHI in Kubernetes!'" ]
  imagePullSecrets:
  - name: <secret name>
EOF
```

Check the pod status to ensure it completed successfully:

```console
$ kubectl get -n <kubernetes namespace> pods/dhi-test
```

A successful test shows `Completed` status:

```console
NAME       READY   STATUS      RESTARTS     AGE
dhi-test   0/1     Completed   ...          ...
```

If you see `ErrImagePull` status instead, there's an issue with your secret
configuration:

```console
NAME       READY   STATUS         RESTARTS   AGE
dhi-test   0/1     ErrImagePull   0          ...
```

Verify the pod output matches the expected message:

```console
$ kubectl logs -n <kubernetes namespace> pods/dhi-test
Hello from DHI in Kubernetes!
```

Clean up the test pod:

```console
$ kubectl delete -n <kubernetes namespace> pods/dhi-test
```
