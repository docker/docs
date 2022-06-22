---
title: Build, test, and install an extension
description: Docker extension CLI
keywords: Docker, extensions, sdk, cli
redirect_from:
- /desktop/extensions-sdk/extensions/validation/
- /desktop/extensions-sdk/dev/cli/build-test-install-extension/
---

The [sample folder](https://github.com/docker/extensions-sdk/tree/main/samples) contains multiple extensions.
These are Docker developed samples that are not meant to be final products.

To use one of them, navigate to the directory of the extension then build and install it on Docker Desktop.
The `docker extension` commands are carried out by the Extension CLI which is a developer tool. It is not included in the standard Docker Desktop package.

To build the extension, run:

```console
$ make build-extension
# or docker build -t my-extension .
```

To install the extension, run:

```console
$ docker extension install my-extension
```

> Using the CLI to install unpublished extensions
>
> Extensions can install binaries, invoke commands and access files on your machine. Make sure you trust extensions before installing them on your machine.
> {: .warning}

To list all your installed extensions, run:

```console
$ docker extension ls

ID                              PROVIDER            VERSION             UI                   VM                  HOST
docker/hub-explorer-extension   Docker Inc.         0.0.2               1 tab(Explore Hub)   Running(1)          1 binarie(s)
tailscale/docker-extension      Tailscale Inc.      0.0.2               1 tab(Tailscale)     Running(1)          1 binarie(s)
```

To remove the extension, run:

```console
$ docker extension rm my-extension
```

To update an extension with a newer version (local or remote image), run:

```console
$ docker extension update docker/hub-explorer-extension:0.0.3
```
