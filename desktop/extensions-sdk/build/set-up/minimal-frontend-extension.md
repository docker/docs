---
title: Set up a minimal frontend extension 
description: Minimal frontend extension tutorial
keywords: Docker, extensions, sdk, build
redirect_from:
- /desktop/extensions-sdk/tutorials/minimal-frontend-extension/
---

To start creating your extension, you first need a directory with files which range from the extension’s source code to the required extension-specific files. This page provides information on how to set up a minimal frontend extension based on plain HTML.

>Note
>
>Before you start, make sure you have installed the latest version of [Docker Desktop](../../release-notes.md).

## Extension folder structure

In the `minimal-frontend` [sample folder](https://github.com/docker/extensions-sdk/tree/main/samples), you can find a ready-to-go example that represents a UI Extension built on HTML. We will go through this code example in this tutorial.

Although you can start from an empty directory, it is highly recommended that you start from the template below and change it accordingly to suit your needs.

```bash
.
├── Dockerfile # (1)
├── metadata.json # (2)
└── ui # (3)
    └── index.html
```

1. Contains everything required to build the extension and run it in Docker Desktop.
2. A file that provides information about the extension such as the name, description, and version.
3. The source folder that contains all your HTML, CSS and JS files. There can also be other static assets such as logos and icons. For more information and guidelines on building the UI, see the [Design and UI styling section](../../design/design-guidelines.md).

If you want to set up user authentication for the extension, see [Authentication](../../dev/oauth2-flow.md).

## Create a Dockerfile

At a minimum, your Dockerfile needs:

- Labels which provide extra information about the extension.
- The source code which in this case is an `index.html` that sits within the `ui` folder.
- The `metadata.json` file.

```Dockerfile
FROM scratch

LABEL org.opencontainers.image.title="MinimalFrontEnd" \
    org.opencontainers.image.description="A sample extension to show how easy it's to get started with Desktop Extensions." \
    org.opencontainers.image.vendor="Docker Inc." \
    com.docker.desktop.extension.api.version="1.0.0-beta.1" \
    com.docker.desktop.extension.icon="https://www.docker.com/wp-content/uploads/2022/03/Moby-logo.png"

COPY ui ./ui
COPY metadata.json .
```

## Configure the metadata file

A `metadata.json` file is required at the root of the image filesystem.

```json
{
  "ui": {
    "dashboard-tab": {
      "title": "Min FrontEnd Extension",
      "root": "/ui",
      "src": "index.html"
    }
  }
}
```

For more information on the `metadata.json`, see [Metadata](../../extensions/METADATA.md).

## What's next?

Learn how to [build and install your extension](../build-install.md).
