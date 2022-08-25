---
title: Try Atomist
description:
keywords:
---

{% include atomist/disclaimer.md %}

The quickest way to try Atomist is to run it locally as a CLI tool. We've built
a Docker image that can scan your local images, without having to integrate with
and connect to a container registry. The CLI uses your local Docker daemon
directly to upload the Software Bill of Materials (SBOM) to the Atomist control
plane for analysis.

Before you can begin the setup, you’ll need a Docker ID. If you do not already
have one, you can [sign up here](https://hub.docker.com/signup){: target="blank"
rel="noopener" class=""}.

1. Go to the [Atomist website](https://dso.docker.com) and sign in using your
   Docker ID.
2. Open the **Integrations** tab.
3. Under **API Keys**, create a new API key.
4. Open your terminal of choice.
5. Invoke the CLI tool using `docker run`. Update the following values:

   - `--workspace`: the workspace ID found on the **Integrations** page on the
     Atomist website.
   - `--api-key`: the API key you just created.
   - `--image`: the Docker image that you want to scan.

This CLI is hosted in a container itself and can be invoked like this:

```shell
docker run \
   -v /var/run/docker.sock:/var/run/docker.sock \
   -ti atomist/docker-registry-broker:latest \
   index-image local \
   --workspace AQ1K5FIKA \
   --api-key team::6016307E4DF885EAE0579AACC71D3507BB38E1855903850CF5D0D91C5C8C6DC0 \
   --image docker.io/atomist/skill:alpine_3.15-node_16
```

which should output something like this:

```shell
 [info] Starting session with correlation-id 5c4f2a81-5370-4536-bc81-2af2ecbaf802
 [info] Starting atomist/docker-vulnerability-scanner-skill 'index_image' (bb5674d) atomist/skill:0.12.0-main.15 (522efce) nodejs:16.14.0
 [info] Indexing image docker.io/atomist/skill:alpine_3.15-node_16
 [info] Downloading image
 [info] Download completed
 [info] Indexing completed
 [info] Mapped packages to layers
 [info] Transacting 33 packages
 [info] Indexing completed successfully
 [info] Transacting image manifest for docker.io/atomist/skill:alpine_3.15-node_16 with digest sha256:9c3c9b88e031466f446471ee8a4233c60c326e785b849985776efbb890f0ec51
 [info] Successfully transacted entities in team AQ1K5FIKA
 [info] Transacting SBOM...
 [info] Successfully transacted entities in team AQ1K5FIKA
```

> Note
>
> The image must have a tag (e.g. `myimage:latest`) so that you can identify the
> scan in the [web GUI](https://dso.docker.com/r/auth/overview/images).
