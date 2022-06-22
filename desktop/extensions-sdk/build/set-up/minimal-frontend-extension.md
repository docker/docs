---
title: Minimal frontend extension tutorial
description: Minimal frontend extension tutorial
keywords: Docker, extensions, sdk, tutorial
redirect_from:
- /desktop/extensions-sdk/tutorials/minimal-frontend-extension/
---

Learn how to create a minimal frontend extension based on plain HTML.

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Docker Extensions CLI](https://github.com/docker/extensions-sdk/releases/)

## Extension folder structure

A Docker extension is made of several files which range from the extension's source code to the required extension-specific files.

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
3. The source folder that contains all your HTML, CSS and JS files. There can also be other static assets such as logos and icons.

## Create a Dockerfile

An extension requires a `Dockerfile` to build, publish, and run in Docker Desktop.

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

## Build the extension

To build the extension, run:

```bash
docker build -t desktop-frontend-minimal-extension:0.0.1 .
```

### Build the extension for multiple platforms

To build the extension for multiple platforms, run:

```bash
docker buildx build --platform=linux/amd64,linux/arm64 -t desktop-frontend-minimal-extension:0.0.1 .
```

## Validate the extension

Verify the extension image is compliant.

The validation checks if the extension's `Dockerfile` specifies all the required labels and if the metadata file is valid against the JSON schema file.

```bash
docker extension validate desktop-frontend-minimal-extension:0.0.1
```

If your extension is valid, the message below displays:

`The extension image "desktop-frontend-minimal-extension:0.0.1" is valid`.

## Install the extension

To install the extension in Docker Desktop, run:

```bash
docker extension install desktop-frontend-minimal-extension:0.0.1
```

If the installation is successful, the output below displays:

```bash
Installing new extension "MinimalFrontEnd" with desktop-frontend-minimal-extension:0.0.1 ...
Installing Desktop extension UI for tab "Min FrontEnd Extension"...
Extension UI tab "Min FrontEnd Extension" added.
Extension "MinimalFrontEnd" installed successfully
```

## Preview the extension

You can also enter the command below to verify the extension installed successfully:

```bash
docker extension ls
```

It outputs all the extensions installed:

```bash
PLUGIN              PROVIDER            IMAGE                                     UI                  VM  HOST
MyExtension         Docker Inc.         desktop-frontend-minimal-extension:0.0.1  1 tab(Min FrontEnd Extension) -   -
```

To preview the extension in Docker Desktop, close and open the Docker Desktop dashboard once the installation is complete.

The left-hand menu displays a new tab with the name `Min FrontEnd Extension`. When you select the new tab, `Hello, World!` displays in the top-left corner.

![UI Extension](images/ui-minimal-extension.png)

## Update the extension

To update the extension, you must first rebuild and revalidate your extension. You can then use the update command.

As an example, let's update the html file to include some inline CSS to support a dark mode.

```html
<head>
  <style>
    @media (prefers-color-scheme: dark) {
      body {
        background-color: #333;
      }

      h1 {
        color: white;
      }
    }
  </style>
  ...
</head>
```

Alternatively remove the `index.html` file and rename `updatedindex.html` to index.html in the ui directory.

Rebuild and revalidate the extension:

```bash
docker build -t desktop-frontend-minimal-extension:0.0.1 .
docker extension validate desktop-frontend-minimal-extension:0.0.1
```

Now update the extension:

```bash
docker extension update desktop-frontend-minimal-extension:0.0.1
```

If the update is successful, the following output displays:

```bash
Removing extension desktop-frontend-minimal-extension:0.0.1...
Extension UI tab Min FrontEnd Extension removed
Extension "MinimalFrontEnd" removed
Installing new extension "desktop-frontend-minimal-extension:0.0.1"
Installing Desktop extension UI for tab "Min FrontEnd Extension"...
Extension UI tab "Min FrontEnd Extension" added.
Extension "MinimalFrontEnd" installed successfully
```

When you run Docker Desktop in dark mode and click the `Min FrontEnd Extension` tab, it renders with dark mode colors.

![UI Extension](images/ui-minimal-extension-dark.png)

## Publish the extension

To publish the extension, upload the Docker image to [DockerHub](https://hub.docker.com).

Tag the previous image to prepend the account owner at the beginning of the image name:

```bash
docker tag desktop-frontend-minimal-extension:0.0.1 owner/desktop-frontend-minimal-extension:0.0.1
```

Push the image to DockerHub:

```bash
docker push owner/desktop-frontend-minimal-extension:0.0.1
```

> Publishing extensions in the marketplace
>
> For Docker Extensions images to be listed in Docker Desktop, they must be approved by Docker and the tags must follow semantic versioning, e.g: `0.0.1`.
>
> See [distribution and new releases](../extensions/DISTRIBUTION.md#distribution-and-new-releases) for more information.
>
> See <a href="https://semver.org/" target="__blank">semver.org</a> to learn more about semantic versioning.

> Having trouble to push the image?
>
> Ensure you are logged into DockerHub. Otherwise, run `docker login` to authenticate.

## Clean up

To remove the extension, run:

```bash
docker extension rm desktop-frontend-minimal-extension
```

The following output displays:

```bash
Removing extension desktop-frontend-minimal-extension...
Extension UI tab Min FrontEnd Extension removed
Extension "MinimalFrontEnd" removed
```

## What's next?

Learn how to [create an extension using React.](./react-extension.md)
