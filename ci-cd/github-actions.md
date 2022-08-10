---
description: Configure GitHub Actions
keywords: CI/CD, GitHub Actions,
title: Configure GitHub Actions
---

This page guides you through the process of setting up a GitHub Action CI/CD
pipeline with Docker. Before setting up a new pipeline, we recommend that you take a look at [Ben's blog](https://www.docker.com/blog/best-practices-for-using-docker-hub-for-ci-cd/){:target="_blank" rel="noopener" class="_"}
on CI/CD best practices.

This guide contains instructions on how to:

1. Use a sample Docker project as an example to configure GitHub Actions.
2. Set up the GitHub Actions workflow.
3. Optimize your workflow to reduce build time.
4. Push only specific versions to Docker Hub.

## Set up a Docker project

Let's get started. This guide uses a simple Docker project as an example. The
[SimpleWhaleDemo](https://github.com/usha-mandya/SimpleWhaleDemo){:target="_blank" rel="noopener" class="_"}
repository contains a Nginx alpine image. You can either clone this repository,
or use your own Docker project.

![SimpleWhaleDemo](images/simplewhaledemo.png){:width="500px"}

Before we start, ensure you can access [Docker Hub](https://hub.docker.com/)
from any workflows you create. To do this:

1. Add your Docker ID as a secret to GitHub. Navigate to your GitHub repository
and click **Settings** > **Secrets** > **New secret**.

2. Create a new secret with the name `DOCKER_HUB_USERNAME` and your Docker ID
as value.

3. Create a new Personal Access Token (PAT). To create a new token, go to
[Docker Hub Settings](https://hub.docker.com/settings/security) and then click
**New Access Token**.

4. Let's call this token **simplewhaleci**.

    ![New access token](images/github-access-token.png){:width="500px"}

5. Now, add this Personal Access Token (PAT) as a second secret into the GitHub
secrets UI with the name `DOCKER_HUB_ACCESS_TOKEN`.

    ![GitHub Secrets](images/github-secrets.png){:width="500px"}

## Set up the GitHub Actions workflow

In the previous section, we created a PAT and added it to GitHub to ensure we
can access Docker Hub from any workflow. Now, let's set up our GitHub Actions
workflow to build and store our images in Hub.

In this example, let us set the push flag to `true` as we also want to push.
We'll then add a tag to specify to always go to the latest version. Lastly,
we'll echo the image digest to see what was pushed.

To set up the workflow:

1. Go to your repository in GitHub and then click **Actions** > **New workflow**.
2. Click **set up a workflow yourself** and add the following content:

First, we will name this workflow:

{% raw %}
```yaml
name: ci
```
{% endraw %}

Then, we will choose when we run this workflow. In our example, we are going to
do it for every push against the main branch of our project:

{% raw %}
```yaml
on:
  push:
    branches:
      - 'main'
```
{% endraw %}

Now, we need to specify what we actually want to happen within our workflow
(what jobs), we are going to add our build one and select that it runs on the
latest Ubuntu instances available:

{% raw %}
```yaml
jobs:
  build:
    runs-on: ubuntu-latest
```
{% endraw %}

Now, we can add the steps required:
* The first one checks-out our repository under `$GITHUB_WORKSPACE`, so our workflow
can access it.
* The second one will use our PAT and username to log into Docker Hub.
* The third will setup Docker Buildx to create the builder instance using a
BuildKit container under the hood.

{% raw %}
```yaml
    steps:
      -
        name: Checkout 
        uses: actions/checkout@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v1
        with:
          username: ${{ secrets.DOCKER_HUB_USERNAME }}
          password: ${{ secrets.DOCKER_HUB_ACCESS_TOKEN }}
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v1
      -
        name: Build and push
        uses: docker/build-push-action@v2
        with:
          context: .
          file: ./Dockerfile
          push: true
          tags: ${{ secrets.DOCKER_HUB_USERNAME }}/simplewhale:latest
```
{% endraw %}

Now, let the workflow run for the first time and then tweak the Dockerfile to
make sure the CI is running and pushing the new image changes:

![CI to Docker Hub](images/ci-to-hub.png){:width="500px"}

## Optimizing the workflow

Next, let's look at how we can optimize the GitHub Actions workflow through
build cache using the registry. This allows to reduce the build time as it
will not have to run instructions that have not been impacted by changes in
your Dockerfile or source code and also reduce number of pulls we complete
against Docker Hub.

In this example, we need to add some extra attributes to the build and push
step:

{% raw %}
```yaml
      -
        name: Login to Docker Hub
        uses: docker/login-action@v1
        with:
          username: ${{ secrets.DOCKER_HUB_USERNAME }}
          password: ${{ secrets.DOCKER_HUB_ACCESS_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v2
        with:
          context: ./
          file: ./Dockerfile
          builder: ${{ steps.buildx.outputs.name }}
          push: true
          tags: ${{ secrets.DOCKER_HUB_USERNAME }}/simplewhale:latest
          cache-from: type=registry,ref=${{ secrets.DOCKER_HUB_USERNAME }}/simplewhale:buildcache
          cache-to: type=registry,ref=${{ secrets.DOCKER_HUB_USERNAME }}/simplewhale:buildcache,mode=max
```
{% endraw %}

As you can see, we are using the `type=registry` cache exporter to import/export
cache from a cache manifest or (special) image configuration. Here it will be
pushed as a specific tag named `buildcache` for our image build.

Now, run the workflow again and verify that it uses the build cache.

## Push tagged versions and handle pull requests

Earlier, we learnt how to set up a GitHub Actions workflow to a Docker project,
how to optimize the workflow by setting up cache. Let's now look at how we can
improve it further. We can do this by adding the ability to have tagged versions
behave differently to all commits to master. This means, only specific versions
are pushed, instead of every commit updating the latest version on Docker Hub.

You can consider this approach to have your commits pushed as an edge tag to
then use it in nightly tests. By doing this, you can always test the last
changes of your active branch while reserving your tagged versions for release
to Docker Hub.

First, let us modify our existing GitHub workflow to take into account pushed
tags and pull requests:

{% raw %}
```yaml
on:
  push:
    branches:
      - 'main'
    tags:
      - 'v*'
```
{% endraw %}

This ensures that the CI will trigger your workflow on push events
(branch and tags). If we tag our commit with something like `v1.0.2`:

{% raw %}
```console
$ git tag -a v1.0.2
$ git push origin v1.0.2
```
{% endraw %}

Now, go to GitHub and check your Actions

![Push tagged version](images/push-tagged-version.png){:width="500px"}

Let's reuse our current workflow to also handle pull requests for testing
purpose but also push our image in the GitHub Container Registry.

First we have to handle pull request events:

{% raw %}
```yaml
on:
  push:
    branches:
      - 'main'
    tags:
      - 'v*'
  pull_request:
    branches:
      - 'main'
```
{% endraw %}

To authenticate against the [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry){:target="_blank" rel="noopener" class="_"},
use the [`GITHUB_TOKEN`](https://docs.github.com/en/actions/reference/authentication-in-a-workflow){:target="_blank" rel="noopener" class="_"}
for the best security and experience.

Now let's change the Docker Hub login with the GitHub Container Registry one:

{% raw %}
```yaml
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v1
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
```
{% endraw %}

Remember to change how the image is tagged. The following example keeps ‘latest'
as the only tag. However, you can add any logic to this if you prefer:

{% raw %}
```yaml
  tags: ghcr.io/<username>/simplewhale:latest
```
{% endraw %}

> **Note**: Replace `<username>` with the repository owner. We could use
> {% raw %}`${{ github.repository_owner }}`{% endraw %} but this value can be mixed-case, so it could
> fail as [repository name must be lowercase](https://github.com/docker/build-push-action/blob/master/TROUBLESHOOTING.md#repository-name-must-be-lowercase){:target="_blank" rel="noopener" class="_"}.

![Update tagged images](images/ghcr-logic.png){:width="500px"}

Now, we will have two different flows: one for our changes to master, and one
for our pushed tags. Next, we need to modify what we had before to ensure we are
pushing our PRs to the GitHub registry rather than to Docker Hub.

## Conclusion

In this guide, you have learnt how to set up GitHub Actions workflow to an
existing Docker project, optimize your workflow to improve build times and
reduce the number of pull requests, and finally, we learnt how to push only
specific versions to Docker Hub.

## Next steps

You can now consider setting up nightly builds, test your image before pushing
it, setting up secrets, share images between jobs or automatically handle
tags and OCI Image Format Specification labels generation.

To look at how you can do one of these, or to get a full example on how to set
up what we have accomplished today, check out [our advanced examples](https://github.com/docker/build-push-action/tree/master/docs/advanced){:target="_blank" rel="noopener" class="_"}
which runs you through this and more details.
