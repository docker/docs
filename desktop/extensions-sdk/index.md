---
title: Overview
description: Overall index for Docker Extensions SDK documentation
keywords: Docker, Extensions, sdk
redirect_from:
 - /desktop/extensions-sdk/dev/overview/
---

Use the resources in this section to create your own Docker Extension.

> Beta
>
> The Docker Extensions SDK is currently in [Beta](../../release-lifecycle.md#beta).
> Features and APIs detailed below are subject to change.

Extensions are packaged as specially formatted Docker images, which our CLI tool helps to build. At the root of the image filesystem is a `metadata.json` file which describes the content of the extension. It is a fundamental element of a Docker extension.

An extension can contain a UI part and backend parts that run either on the host or in the Desktop virtual machine. For further details, see [Extension metadata](extensions/METADATA.md).

Extensions are distributed through the Docker Hub.
Development of extensions can be done locally without the need to push the extension to Docker Hub. See [Extensions distribution](extensions/DISTRIBUTION.md) for further details.

## How to use the resources in this section

The documentation is structured so that it matches the steps you need to take when creating your extension. There are two main parts to creating a Docker Extension; build the foundations and then publish your extension.

### Part one: Build the foundations

The build process consists of:

- Installing the latest version of Docker Desktop.
- Setting up the directory with files which range from the extension’s source code to the required extension-specific files.
- Creating the Dockerfile to build, publish, and run your extension in Docker Desktop.
- Configuring the metadata file which is required at the root of the image filesystem.
- Building and installing the extension.
- Testing and debugging the extension.

There are also instructions on [how to set authentication](dev/oauth2-flow.md) for your extension.

This build section provides sample folders with ready-to-go examples that walk you through building:

- A frontend extension based on plain HTML
- A simple Docker extension that contains only a UI part and is based on ReactJS. This is useful if you want to develop an extension which consists exclusively of a visual part with no services running in the VM.
- An extension that invokes Docker CLI commands
- A simple backend extension

>Note
>
>Whilst you're building your extension, make sure you follow our [design](design/design-guidelines.md) and [UI styling](design/overview.md) guidelines to ensure visual consistency and [level AA accessibility standards](https://www.w3.org/WAI/WCAG2AA-Conformance).

If your extension requires additional services running in the Docker Desktop VM, see the [VM UI](https://github.com/docker/extensions-sdk/tree/main/samples/vm-service) example.

For further inspiration, see the other examples in the [samples folder](https://github.com/docker/extensions-sdk/tree/main/samples)

### Part two: Publish and distribute your extension

Docker Desktop displays published extensions in the Extensions Marketplace. If you want your extension to be published in the Marketplace, you can submit your extension [here](https://www.docker.com/products/extensions/submissions/). We’ll review your submission and provide feedback if changes are needed before we can validate and publish it to make it available to all Docker Desktop users.

## What’s next?
If you want to get up and running quickly with a Docker Extension, see the [Quickstart guide](quickstart.md).

Alternatively, get started with Part one: Build for more in-depth information about each step of the extension creation process.

For an in-depth tutorial of the entire build process, we recommend the following video walkthrough from DockerCon 2022:

<iframe width="560" height="315" src="https://www.youtube.com/embed/Yv7OG-EGJsg" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
