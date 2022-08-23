---
title: "Step two: Build and install"
description: Step three in the extension creation process
keywords: Docker, Extensions, sdk, validate, install
redirect_from:
- /desktop/extensions-sdk/extensions/validation/
- /desktop/extensions-sdk/dev/cli/build-test-install-extension/
---

Once you have set up your directory correctly, you can build your extension for single or multiple architectures. 

## Build the extension for a single architecture

To build the extension, run:

`docker build -t <name-of-your-extension> .`

## Validate your extension

The Extensions CLI lets you validate your extension before installing and running it locally.

The validation checks if the extension’s `Dockerfile` specifies all the required labels and if the metadata file is valid against the JSON schema file.

To validate, run:

`docker extension validate <name-of-your-extension>`

If your extension is valid, the message below displays:

`The extension image "name-of-your-extension" is valid`.

Before the image is built, it is also possible to validate only the metadata.json file:

`$ docker extension validate /path/to/metadata.json`

The JSON schema used to validate the `metadata.json` file against can be found under the [releases page](https://github.com/docker/extensions-sdk/releases/latest).

## Install the extension

To install the extension in Docker Desktop, run:

`docker extension install <name-of-your-extension>`

> Note 
> 
> Extensions can install binaries, invoke commands and access files on your machine. Make sure you trust extensions before installing them on your machine.

To list all your installed extensions, run:

```typescript
$ docker extension ls

ID                              PROVIDER            VERSION             UI                   VM                  HOST
docker/hub-explorer-extension   Docker Inc.         0.0.2               1 tab(Explore Hub)   Running(1)          1 binarie(s)
```

## What's next?

Learn how to [test and debug](test-debug.md) your extension.
