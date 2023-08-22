---
description: 'Docker Scout can integrate with runtime environments to give you realtime

  insights about your software supply chain.

  '
keywords: supply chain, security, streams, environments, workloads, deployments
title: Integrating Docker Scout with environments
---

{{< include "scout-early-access.md" >}}

You can integrate Docker Scout with your runtime environments, and get insights
for your running workloads. This gives you a realtime view of your security
status for your deployed artifacts.

Docker Scout lets you define multiple environments, and assign images to
different environments. This gives you a complete overview of your software
supply chain, and lets you view and compare deltas between environments, for
example staging and production.

How you define and name your environments is up to you. You can use patterns
that are meaningful to you and that matches how you ship your applications.

## Assign to environments

Each environment contains references to a number of images. These references
represent containers currently running in that particular environment.

For example, say you're running `myorg/webapp:3.1` in production, you can
assign that tag to your `production` environment. You might be running a
different version of the same image in staging, in which case you can assign
that version of the image to the `staging` environment.

## Comparing between environments

Assigning images to environments lets you make comparisons with and between
environments. This is useful for things like GitHub pull requests, for
comparing the image built from the code in the PR to the corresponding image in
staging or production.

You can also compare with streams using the `--to-stream` flag on the
[`docker scout compare`](../../../engine/reference/commandline/scout_compare.md)
CLI command:

```console
$ docker scout compare --to-stream production myorg/webapp:latest
```

## Assign images to environments

To add environments to Docker Scout, you can:

- Use the `docker scout stream` command in the Docker CLI:

  ```console
  $ docker scout stream <environment> <image>
  ```

- Use the [Docker Scout GitHub Action](https://github.com/marketplace/actions/docker-scout#record-an-image-deployed-to-a-stream-environment)

## View images for an environment

To view the images for an environment:

1. Go to the [Docker Scout Dashboard](https://scout.docker.com/).
2. Select the **Images** tab.
3. Open the **Environments** drop-down menu.
4. Select the environment that you want to view.

The list displays all images that have been assigned to the selected
environment. If you've deployed multiple versions of the same image in an
environment, all versions of the image appear in the list.

### Mismatching image tags

When you've selected an environment on the **Images** tab, tags in the list
represent the tag that was used to deploy the image. Tags are mutable, meaning
that you can change the image digest that a tag refers to. If Docker Scout
detects that a tag refers to an outdated digest, a warning icon displays next
to the image name.