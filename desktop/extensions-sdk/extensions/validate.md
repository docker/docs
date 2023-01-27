---
title: Validate your extension
description: Step three in the extension creation process
keywords: Docker, Extensions, sdk, validate, install
redirect_from:
- /desktop/extensions-sdk/extensions/validation/
- /desktop/extensions-sdk/build/build-install/
- /desktop/extensions-sdk/dev/cli/build-test-install-extension/
---

Validate your extension before you share or publish it. Validating the extension ensures:

- That the extension is built with the [image labels](labels.md) it requires to display correctly in the marketplace
- That the extension installs and runs without problems

The Extensions CLI lets you validate your extension before installing and running it locally.

The validation checks if the extension’s `Dockerfile` specifies all the required labels and if the metadata file is valid against the JSON schema file.

To validate, run:

```console
$ docker extension validate <name-of-your-extension>
```

If your extension is valid, the message below displays:

```console
The extension image "name-of-your-extension" is valid
```

Before the image is built, it is also possible to validate only the metadata.json file:

```console
$ docker extension validate /path/to/metadata.json
```

The JSON schema used to validate the `metadata.json` file against can be found under the [releases page](https://github.com/docker/extensions-sdk/releases/latest){:target="_blank" rel="noopener" class="_"}.

