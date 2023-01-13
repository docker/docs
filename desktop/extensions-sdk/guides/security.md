---
title: Security
description: Aspects of the security model of extensions
keywords: Docker, extensions, sdk, security
---

## Extension capabilities

An extension can have the following optional parts: 
* a user interface in HTML or JavaScript, displayed in Docker Desktop Dashboard
* a backend part that runs as a container
* executables deployed on the host machine.

Extensions are executed with the same permissions as the Docker Desktop user. Extension capabilities include running any Docker commands (including running containers and mounting folders), running extension binaries, and accessing files on your machine that are accessible by the user running Docker Desktop.

The Extensions SDK provides a set of JavaScript APIs to invoke commands or invoke these binaries from the extension UI code. Extensions can also provide a backend part that starts a long-lived running container in the background.

> Note
> Make sure you trust the publisher or author of the extension when you install it, as the extension has the same access rights as the user running Docker Desktop.
{: .important}

Learn more in the [architecture section](https://docs.docker.com/desktop/extensions-sdk/architecture/)
