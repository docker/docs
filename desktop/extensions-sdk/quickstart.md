---
title: Quickstart
description: Guide on how to build an extension quickly
keywords: quickstart, extensions
redirect_from:
  - desktop/extensions-sdk/tutorials/initialize/
---

Follow the guide below to build a basic Docker Extension quickly. The Quickstart guide automatically generates boilerplate files for you.

> Note
>
> NodeJS and Go are only required when you follow the quickstart guide to build an extension. It uses the `docker extension init` command to automatically generate boilerplate files. This command uses a template based on a ReactJS and Go application.

## Prerequisites

- [Docker Desktop](../release-notes.md)
- [NodeJS](https://nodejs.org/)
- [Go](https://go.dev/dl/)

## Step one: Set up your directory

To set up your directory, use the `init` subcommand and provide a name for your extension.

`docker extension init my-extension`

You’ll be asked a series of questions about your extension, such as its name, a description, and the name of your Hub repository. This helps the CLI generate a set of boilerplate files for you to get started. The boilerplate files are stored in the directory `my-extension`.

The automatically generated extension contains:

- A backend service that listens on a socket. It has one endpoint `/hello` that returns a JSON payload.
- A React frontend that can call the backend and output the backend’s response.

For more information and guidelines on building the UI, see the [Design and UI styling section](design/design-guidelines.md).

## Step two: Build the extension

To build your extension, run:

`docker build -t <name-of-your-extension> .`

`docker build` builds your extension and also generates an image which is named after your chosen hub repository. For instance, if you typed `john/my-extension` as the answer to the following question:

`? Hub repository (eg. namespace/repository on hub): john/my-extension`

The `docker build` generates an image with name `john/my-extension`.

## Step three: Install and preview the extension

To install the extension in Docker Desktop, run:

`docker extension install <name-of-your-extension>`

To preview the extension in Docker Desktop, close and open Docker Dashboard once the installation is complete.

During UI development, it’s helpful to use hot reloading to test your changes without rebuilding your entire extension. See [Preview whilst developing the UI](build/test-debug.md#hot-reloading-whilst-developing-the-ui) for more information.

## Step four: Submit and publish your extension to the Marketplace

If you want to make your extension available to all Docker Desktop users, you can submit it for publication in the Marketplace. For more information, see [Publish](extensions/index.md).

## Clean up

To remove the extension, run:

`docker extension rm <name-of-your-extension>`

## What's next

Learn more about:

- [Building and installing an extension](build/build-install.md)
- [Testing and debugging](build/test-debug.md)
- [Setting up authentication](dev/oauth2-flow.md)
- [Designing the UI](design/design-guidelines.md)
