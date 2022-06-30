---
title: "Step one: Get started"
description: The first step to building an extension.
keywords: Docker, Extensions, sdk, prerequisites
redirect_from:
  - desktop/extensions-sdk/tutorials/initialize/
---

To start creating your extension, you first need a directory with files which range from the extension’s source code to the required extension-specific files.

## Prerequisites

Before you create your own extension, you need to install [Docker Desktop](../../release-notes.md).

You can list installed extensions (the list should be empty if you have not installed extensions already):

```console
$ docker extension ls
ID                  PROVIDER            VERSION             UI                  VM                  HOST
```

![Extensions enabled](images/extensions-enabled.png)

You can now continue to step two and set up your directory.

## What’s next?

Explore how to set up:

- [A frontend extension based on plain HTML](set-up/minimal-frontend-extension.md)
- [A simple Docker extension that contains only a UI part and is based on ReactJS](set-up/react-extension.md)
- [An extension that invokes Docker CLI commands](set-up/minimal-frontend-using-docker-cli.md)
- [A simple backend extension](set-up/minimal-backend-extension.md)
