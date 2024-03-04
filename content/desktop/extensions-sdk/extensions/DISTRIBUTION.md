---
title: Package and release your extension
description: Docker extension disctribution
keywords: Docker, extensions, sdk, distribution
---

This page contains additional information on how to package and distribute extensions.

## Package your extension

Docker extensions are packaged as Docker images. The entire extension runtime including the UI, backend services (host or VM), and any necessary binary must be included in the extension image.
Every extension image must contain a `metadata.json` file at the root of its filesystem that defines the [contents of the extension](../architecture/metadata.md).

The Docker image must have several [image labels](labels.md), providing information about the extension. See how to use [extension labels](labels.md) to provide extension overview information.

To package and release an extension, you must build a Docker image (`docker build`), and push the image to [Docker Hub](https://hub.docker.com/) (`docker push`) with a specific tag that lets you manage versions of the extension.

## Release your extension

Docker image tags must follow semver conventions in order to allow fetching the latest version of the extension, and to know if there are updates available. See [semver.org](https://semver.org/) to learn more about semantic versioning.

Extension images must be multi-arch images so that users can install extensions on ARM/AMD hardware. These multi-arch images can include ARM/AMD specific binaries. Mac users will automatically use the right image based on their architecture.
Extensions that install binaries on the host must also provide Windows binaries in the same extension image. See how to [build a multi-arch image](multi-arch.md) for your extension.

You can implement extensions without any constraints on the code repository. Docker doesn't need access to the code repository in order to use the extension. Also, you can manage new releases of your extension, without any dependency on Docker Desktop releases.

## New releases and updates

You can release a new version of your Docker extension by pushing a new image with a new tag to Docker Hub.

Any new image pushed to an image repository corresponding to an extension defines a new version of that extension. Image tags are used to identify version numbers. Extension versions must follow semver to make it easy to understand and compare versions.

Docker Desktop scans the list of extensions published in the marketplace for new versions, and provides notifications to users when they can upgrade a specific extension. Extensions that aren't part of the Marketplace don't have automatic update notifications at the moment.

Users can download and install the newer version of any extension without updating Docker Desktop itself.

## Extension API dependencies

Extensions must specify the Extension API version they rely on. Docker Desktop checks the extensions required version, and only proposes to install extensions that are compatible with the current Docker Desktop version installed. Users might need to update Docker Desktop in order to install the latest extensions available.

Extension image labels must specify the API version that the extension relies upon. This allows Docker Desktop to inspect newer versions of extension images without downloading the full extension image upfront.

## License on extensions and the extension SDK

The [Docker Extension SDK](https://www.npmjs.com/package/@docker/extension-api-client) is licensed under the Apache 2.0 License and is free to use.

There is no constraint on how each extension should be licensed, this is up to you to decide when creating a new extension.