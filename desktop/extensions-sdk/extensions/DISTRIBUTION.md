---
title: Distribute your extension
description: Docker extension disctribution
keywords: Docker, extensions, sdk, distribution
---

## Packaging

Docker Extensions are packaged as Docker images. The entire extension runtime including the UI, backend services (host or VM) and any necessary binary must be included in the extension image.
Every extension image must contain a metadata.json file at the root of its filesystem that [defines the contents of the extension](METADATA.md).

The image must have [several labels](labels.md):

- `org.opencontainers.image.title` for the name of the extension.
- `org.opencontainers.image.vendor` for the provider for the extension.
- `com.docker.desktop.extension.api.version`for the Docker API version the extension is compatible with.

Packaging and releasing an extension is done by running `docker build` to create the image, and `docker push` to make the image available on Docker Hub with a specific tag that allows you to manage versions of the extension.

Take advantage of multi-arch images to build images that include ARM/AMD binaries. The right image will be used for Mac users depending on their architecture.
For extensions on Docker Desktop for Windows, Windows binaries that are to be installed on the host must be included in the same extension image. We will revisit this with some tag conventions to allow some images specific to Windows, and other images specific to Mac, based on a tag prefix. See [how to build extensions for multiple architectures](../build/build.md).

You can implement extensions without any constraints on the code repository. Docker does not need access to the code repository in order to use the extension. Release of new versions of the extension is managed you.

## Distribution and new releases

Releasing a Docker Desktop extension is done by running `docker push` to push the extension image to Docker Hub.

Docker Desktop includes an allow-list of extensions available to users. The extension image Hub repository (like `mycompany/my-desktop-extension`) must be part of the Docker Desktop allow-list to be recognized as an extension.

This allow-list specifies which Hub repositories are to be used by Docker Desktop to download and install extensions with a specific version at a given point in time.

Any new image pushed to a repository that is part of the allow-list corresponds to a new version of the extension. Image tags are used to identify version numbers. Extension versions must follow semver to make it easy to understand and compare versions.

With a given release of Docker Desktop (including some extensions), users should not need to upgrade Docker Desktop in order to obtain new versions of a specific extension. Newer versions of the extension can be released independently of Docker Desktop releases, provided there is no Extension API mismatchs.

Docker Desktop scans the list of published extensions for new versions regularly, and provides notifications to users when they can upgrade a specific extension.

Users can download and install the newer version of an extension without updating Docker Desktop itself.

## API dependencies

Extensions must specify the Extension API version they rely on. Currently there is no technical validation of this version, as the extension framework is still experimental.

Docker Desktop can use this Extension API version to detect if a newer version of an extension is valid given the user's current version of Docker Desktop. If it is, the user sees a notification to upgrade the corresponding extension.

The API version that the extension relies upon must be specified in the extension image labels. This allows Docker Desktop to inspect newer versions of extension images without downloading the full extension image upfront.
