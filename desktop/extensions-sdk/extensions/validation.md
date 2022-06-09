---
title: Validate your extension
description: Docker extension validaiton
keywords: Docker, extensions, sdk, validation
---

The Extensions CLI lets you validate your extension before installing and running it locally:

```console
$ docker extension validate my-extension
```

It checks the content of the image and ensures the image has the right labels needed for extensions.

Before the image is built, it is also possible to validate only the metadata.json file:

```console
$ docker extension validate /path/to/metadata.json
```

## JSON schema

The JSON schema used to validate the `metadata.json` file against can be found under the [releases page](https://github.com/docker/extensions-sdk/releases/latest).

## Labels

See [labels](labels.md).
