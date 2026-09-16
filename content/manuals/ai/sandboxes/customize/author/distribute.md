---
title: Build and distribute kits
description: Build and publish v3 kit images, share source through Git, and sign images for distribution.
keywords: sandboxes, sbx, kits, v3, workloads, mixins
weight: 40
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

V3 kits are experimental. This page covers building and distributing kits you've
[authored](/manuals/ai/sandboxes/customize/author/_index.md). To run a kit someone else has published, see
[Use kits](/manuals/ai/sandboxes/customize/use-kits.md).

Share kits as published images in a container registry or as source files in
Git. Consumers can [run a kit](/manuals/ai/sandboxes/customize/use-kits.md#run-a-kit) using either type of reference.

## Publish an image

A published v3 kit is an OCI image. Use Docker Buildx to build and push it,
passing the descriptor with `-f` and the source directory as the build context:

```console
$ docker login
$ docker buildx build ./my-kit -f ./my-kit/my-kit.yaml \
    -t docker.io/<NAMESPACE>/my-kit:1.0.0 --push
```

Replace `<NAMESPACE>` with a Docker Hub namespace you can push to. Buildx
uses your `docker login` credentials. Include `docker.io/` explicitly in
Docker Hub references. For pulling from private registries in a sandbox,
configure [Registry credentials](/manuals/ai/sandboxes/configuration/credentials.md#registry-credentials).

Build for both supported Linux architectures when distributing across machines:

```console
$ docker buildx build ./my-kit -f ./my-kit/my-kit.yaml \
    --platform linux/amd64,linux/arm64 \
    -t docker.io/<NAMESPACE>/my-kit:1.0.0 --push
```

An image present only in the host Docker image store isn't available to the
sandbox runtime by registry reference. Push it to a registry, or pass a local
source directory to `sbx` for development. The `sbx kit pack`, `push`, and
`pull` packaging commands are for v1 and v2 kits. V3 uses the Buildx workflow
shown here.

## Share source through Git

Commit the kit's source directory to a Git repository. Share a reference that
identifies the kit directory with `dir` and pins a revision with `ref`:

```text
git+https://github.com/<ORG>/<REPOSITORY>.git#ref=<COMMIT>&dir=my-kit
```

Docker Sandboxes builds the source when a consumer creates a sandbox from
that reference.

## Sign and verify kits

Sign the OCI image after pushing it:

```console
$ sbx kit sign docker.io/<NAMESPACE>/my-kit:1.0.0
$ sbx kit verify docker.io/<NAMESPACE>/my-kit:1.0.0 \
    --certificate-identity <SIGNER_IDENTITY> \
    --certificate-oidc-issuer <ISSUER_URL>
```

These commands use Cosign-compatible Sigstore signatures. For keyless signing,
verification must specify the signer's certificate identity and OpenID Connect
issuer. For key-based signing, pass `--key cosign.key` to `sign` and
`--key cosign.pub` to `verify`.

For signature enforcement when loading kits, see
[Verify kit signatures](/manuals/ai/sandboxes/customize/use-kits.md#verify-kit-signatures).

V3 source directories and Git sources don't support source signing. Publish
and sign an OCI image when signatures are required.

## Published format

A published kit image contains both the kit's content and its descriptor.
Docker image tools can inspect and distribute it, and Docker Sandboxes reads
the descriptor when creating the sandbox.

The built kit also includes its descriptor under
`/usr/share/sandbox/kit/<stem>/`, so you can inspect it inside the sandbox.
For the image annotations and file layout, see
[Published image format](/manuals/ai/sandboxes/customize/author/kit-reference.md#published-image-format).

Running a workload with `docker run` uses its image configuration, but doesn't
apply the kit's capability declarations or lifecycle hooks. Use `sbx` to run
it with those behaviors.

