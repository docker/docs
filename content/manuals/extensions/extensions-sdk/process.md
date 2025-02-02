---
description: Understand the process of creating an extension.
title: The build and publish process
keyword: Docker Extensions, sdk, build, create, publish
aliases:
 - /desktop/extensions-sdk/process/
weight: 10
---

This documentation is structured so that it matches the steps you need to take when creating your extension. 

There are two main parts to creating a Docker extension:

1. Build the foundations
2. Publish the extension

> [!NOTE]
>
> You do not need to pay to create a Docker extension. The [Docker Extension SDK](https://www.npmjs.com/package/@docker/extension-api-client) is licensed under the Apache 2.0 License and is free to use. Anyone can create new extensions and share them without constraints.
> 
> There is also no constraint on how each extension should be licensed, this is up to you to decide when creating a new extension.

## Part one: Build the foundations

The build process consists of:

- Installing the latest version of Docker Desktop.
- Setting up the directory with files, including the extension’s source code and the required extension-specific files.
- Creating the `Dockerfile` to build, publish, and run your extension in Docker Desktop.
- Configuring the metadata file which is required at the root of the image filesystem.
- Building and installing the extension.

For further inspiration, see the other examples in the [samples folder](https://github.com/docker/extensions-sdk/tree/main/samples).

> [!TIP]
>
> Whilst creating your extension, make sure you follow the [design](design/design-guidelines.md) and [UI styling](design/_index.md) guidelines to ensure visual consistency and [level AA accessibility standards](https://www.w3.org/WAI/WCAG2AA-Conformance).

## Part two: Publish and distribute your extension

Docker Desktop displays published extensions in the Extensions Marketplace. The Extensions Marketplace is a curated space where developers can discover extensions to improve their developer experience and upload their own extension to share with the world.

If you want your extension published in the Marketplace, read the [publish documentation](extensions/publish.md).

{{% include "extensions-form.md" %}}

## What’s next?

If you want to get up and running with creating a Docker Extension, see the [Quickstart guide](quickstart.md).

Alternatively, get started with reading the "Part one: Build" section for more in-depth information about each step of the extension creation process.

For an in-depth tutorial of the entire build process, we recommend the following video walkthrough from DockerCon 2022.

<iframe width="560" height="315" src="https://www.youtube.com/embed/Yv7OG-EGJsg" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
