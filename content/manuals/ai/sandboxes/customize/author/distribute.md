---
title: Build and distribute kits
description: Build and publish v3 kit images, share source through Git, and sign images for distribution.
keywords: sandboxes, sbx, kits, v3, workloads, mixins
weight: 40
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Once you've [authored a kit](/manuals/ai/sandboxes/customize/author/_index.md),
share it by publishing an image to a registry or its source to Git.
Publishing an image saves others from building the kit themselves. Sharing
the source lets them build it when they create a sandbox. Either way, `sbx`
reads the kit's descriptor to configure the sandbox.

This page shows how to publish and sign v3 kits. To run a kit someone else has
shared, see [Use kits](/manuals/ai/sandboxes/customize/use-kits.md).

## Publish an image

Use Docker Buildx to build and publish your kit as an OCI image. Pass the YAML
descriptor with `-f` and the source directory as the build context:

```console
$ docker login
$ docker buildx build ./my-kit -f ./my-kit/my-kit.yaml \
    -t docker.io/<NAMESPACE>/my-kit:1.0.0 --push
```

Replace `<NAMESPACE>` with a Docker Hub namespace you can push to. Buildx
uses your `docker login` credentials to push the image.

Push the image before running it by its registry reference. `sbx` pulls the
published image. It can't use images stored only in your host's Docker image
store. During development, you can pass a local source directory to `sbx`
instead. For private images, configure
[registry credentials](/manuals/ai/sandboxes/configuration/credentials.md#registry-credentials)
for your sandbox.

Use Buildx for v3 kits. The `sbx kit pack`, `push`, and `pull` commands are
for v1 and v2 kits.

### Support both Linux architectures

To share the kit with people who use different machines, build for both
supported Linux architectures:

```console
$ docker buildx build ./my-kit -f ./my-kit/my-kit.yaml \
    --platform linux/amd64,linux/arm64 \
    -t docker.io/<NAMESPACE>/my-kit:1.0.0 --push
```

## Publish a kit set

Publish a set with the same Buildx command, passing its YAML descriptor with
`-f`. The build fetches the published images listed in `kits:`, checks whether
their declarations are compatible, and combines their files and settings into
one image. See
[Compose a kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md) for a complete example.

The result is a workload or mixin image that people can use directly. It
records each component's manifest digest, which identifies the exact image
used in the build. Updating a component's tag doesn't change a published
set. To share the update, rebuild the set and publish another version.
If you want rebuilds to use the same component images, add a `digest` beside
each `ref` in `kits:`.

## Share source through Git

Commit the kit's source directory to a Git repository, then share a reference
in this format. Use `dir` to select the kit directory and `ref` to select the
commit:

```text
git+https://github.com/<ORG>/<REPOSITORY>.git#ref=<COMMIT>&dir=my-kit
```

When someone creates a sandbox from that reference, `sbx` builds the kit
from the selected commit.

## Sign and verify kits

Sign your published image so others can verify who signed it:

```console
$ sbx kit sign docker.io/<NAMESPACE>/my-kit:1.0.0
$ sbx kit verify docker.io/<NAMESPACE>/my-kit:1.0.0 \
    --certificate-identity <SIGNER_IDENTITY> \
    --certificate-oidc-issuer <ISSUER_URL>
```

These commands use Sigstore signatures, which are compatible with Cosign.
The example signs without a key, so verification needs the signer's certificate
identity and OpenID Connect issuer. To sign with a key instead, pass
`--key cosign.key` to `sign` and `--key cosign.pub` to `verify`.

To require a valid signature before using a kit, see
[Verify kit signatures](/manuals/ai/sandboxes/customize/use-kits.md#verify-kit-signatures).

V3 kits shared as source directories or Git references can't be signed.
If you need signatures, publish and sign an OCI image.

## Published format

A published kit image contains both the kit's content and its descriptor.
Docker image tools can inspect and distribute it, and Docker Sandboxes reads
the descriptor when creating the sandbox.

You can also inspect the descriptor inside the sandbox, under
`/usr/share/sandbox/kit/<stem>/`.
For the image annotations and file layout, see
[Published image format](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/SPEC-v3.md#10-the-oci-layout).

Use `sbx` to run a workload with the settings and hooks declared in its kit
descriptor. Running it with `docker run` uses only the image configuration,
so those settings and hooks aren't applied.
