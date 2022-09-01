---
title: Try Atomist
description:
keywords:
---

{% include atomist/disclaimer.md %}

The quickest way to try Atomist is to run it locally as a CLI tool. You can run
it as a Docker image that scans one of your local images. This eliminates the
need of having to integrate with and connect to a remote container registry. The
CLI uses your local Docker daemon directly to upload the Software Bill of
Materials (SBOM) to the Atomist control plane for analysis.

Before you can begin the setup, you’ll need a Docker ID. If you do not already
have one, you can [sign up here](https://hub.docker.com/signup){: target="blank"
rel="noopener" class=""}.

> Note
>
> We only recommend using this CLI-based method of scanning images for testing
> purposes. For continuous scanning, please integrate Atomist with your
> container registry. See [get started](./get-started.md).

1. Go to the [Atomist website](https://dso.docker.com) and sign in using your
   Docker ID.
2. Open the **Integrations** tab.
3. Under **API Keys**, create a new API key.
4. In your terminal of choice, invoke the Atomist CLI tool using `docker run`.
   Update the following values:

   - `--workspace`: the workspace ID found on the **Integrations** page on the
     Atomist website.
   - `--api-key`: the API key you just created.
   - `--image`: the Docker image that you want to scan.

   ```bash
   docker run \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -ti atomist/docker-registry-broker:0.0.1 \
      index-image local \
      --workspace AQ1K5FIKA \
      --api-key team::6016307E4DF885EAE0579AACC71D3507BB38E1855903850CF5D0D91C5C8C6DC0 \
      --image docker.io/david/myimage:latest
   ```

   > Note
   >
   > The image must have a tag (for example, `myimage:latest`) so that you are
   > able to identify the scan in the
   > [Atomist web UI](https://dso.docker.com/r/auth/overview/images).

   The output should look something like this:

   ```bash
   [info] Starting session with correlation-id c12e08d3-3bcc-4475-ab21-7114da599eaf
   [info] Starting atomist/docker-vulnerability-scanner-skill 'index_image' (1f99caa) atomist/skill:0.12.0-main.44 (fe90e3c) nodejs:16.15.0
   [info] Indexing image python:latest
   [info] Downloading image
   [info] Download completed
   [info] Indexing completed
   [info] Mapped packages to layers
   [info] Indexing completed successfully
   [info] Transacting image manifest for docker.io/david/myimage:latest with digest sha256:a8077d2b2ff4feb1588d941f00dd26560fe3a919c16a96305ce05f7b90f388f6
   [info] Successfully transacted entities in team AQ1K5FIKA
   [info] Image URL is https://dso.atomist.com/AR5C23OPM/overview/images/python/digests/sha256:a8077d2b2ff4feb1588d941f00dd26560fe3a919c16a96305ce05f7b90f388f6
   [info] Transacting SBOM...
   [info] Successfully transacted entities in team AQ1K5FIKA
   [info] Transacting SBOM...
   ```

5. When the scan is finished, open the
   [Atomist web UI](https://dso.docker.com/r/auth/overview/images), where you
   should see the scanned image in the list.

   ![scanned image in the image overview list](./images/images-overview.png){:
   width="700px"}

6. Click the image name. This gets you to the list of image tags scanned by
   Atomist.

   ![list of image tags](./images/tags-list.png){: width="700px"}

   Since we just ran our first scan, our list only contains one tag, for now.
   When you integrate Atomist with your container registry, any image tag you
   push to your registry will be scanned automatically, and show up in this
   list.

7. Click the tag name. This shows you the scan results for this tag. In this
   view, you can see how many vulnerabilities this image contains, their
   severity levels, whether there is a fix version available, and more.

   ![vulnerability breakdown view](./images/vulnerabilities-overview.png){:
   width="700px"}

This is the end of the tutorial. Take some time to explore the different data
views that Atomist presents about your image. When you're ready, head to the
[get started guide](./get-started.md) to learn how to start integrating Atomist
in your software supply chain.
