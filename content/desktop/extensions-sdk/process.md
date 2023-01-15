---
description: Understand the process of creating an extension.
title: The build and publish process
keyword: Docker Extensions, sdk, build, create, publish
---

The documentation is structured so that it matches the steps you need to take when creating your extension. There are two main parts to creating a Docker Extension; build the foundations and then publish your extension.

### Part one: Build the foundations

The build process consists of:

- Installing the latest version of Docker Desktop.
- Setting up the directory with files which range from the extension’s source code to the required extension-specific files.
- Creating the Dockerfile to build, publish, and run your extension in Docker Desktop.
- Configuring the metadata file which is required at the root of the image filesystem.
- Building and installing the extension.

> Note
>
> Whilst you're building your extension, make sure you follow our [design](design/design-guidelines.md) and [UI styling](design/index.md) guidelines to ensure visual consistency and [level AA accessibility standards](https://www.w3.org/WAI/WCAG2AA-Conformance){:target="_blank" rel="noopener" class="_"}.

For further inspiration, see the other examples in the [samples folder](https://github.com/docker/extensions-sdk/tree/main/samples){:target="_blank" rel="noopener" class="_"}.

### Part two: Publish and distribute your extension

Docker Desktop displays published extensions in the Extensions Marketplace. The Extensions Marketplace is a curated space where developers from all over can discover extensions to improve their developer experience and upload their own extension to share with the world.

All extensions submitted to the Extension Marketplace are reviewed and approved by our team before listing. This review process ensures a level of trust, security, and quality for developers using Extensions and allows for Extension developers to get feedback on what will improve their Extensions experience.

If you want your extension to be published in the Marketplace, you can submit your extension [here](https://www.docker.com/products/extensions/submissions/){:target="_blank" rel="noopener" class="_"}. We’ll review your submission and provide feedback if changes are needed before we can validate and publish it to make it available to all Docker Desktop users.

{% include extensions-form.md %}

## What’s next?

If you want to get up and running quickly with a Docker Extension, see the [Quickstart guide](quickstart.md).

Alternatively, get started with Part one: Build for more in-depth information about each step of the extension creation process.

For an in-depth tutorial of the entire build process, we recommend the following video walkthrough from DockerCon 2022:

<iframe width="560" height="315" src="https://www.youtube.com/embed/Yv7OG-EGJsg" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
