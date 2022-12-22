---
title: Security
description: Aspects of the security model of extensions
keywords: Docker, extensions, sdk, security
---

## Extension capabilities

An extension can have the following optional parts; a user interface in HTML or JavaScript, a backend part that runs as a container, and executables deployed on the host machine.

Extensions are executed with the same permissions as the Docker Desktop user. Extension capabilities include running any docker commands (including running containers and mounting folders), running extension binaries, and accessing files on your machine.

The Extensions SDK provides a set of JavaScript APIs to invoke commands or execute these binaries from the extension UI code. Extensions can also provide a backend part that starts a long-lived running container in the background.

Learn more in the [architecture section](https://docs.docker.com/desktop/extensions-sdk/architecture/)

An extension's code is loaded into the Docker Desktop UI electron app. The extension UI code is subject to some constraints for security reasons, described below.

## Sandboxing

Extension UI code is rendered in a sandboxed environment and does not have a node.js environment initialized, nor direct access to the electron APIs.
The extension UI code cannot perform privileged tasks, such as making changes to the system, or spawning subprocesses, except by using the SDK APIs provided with the extension framework.
It can also perform interactions with Docker Desktop, such as navigating to various places in the Dashboard, only through the extension SDK APIs.

## Extension isolation

Extensions are isolated from each other and extension UI code is running in its own session for each extension. Extensions cannot access other extensions’ session data.

`localStorage` is one of the mechanisms of a browser’s web storage. It allows users to save data as key-value pairs in the browser for later use. `localStorage` does not clear data when the browser (the extension pane) closes. This makes it ideal for persisting data when navigating out of the extension to other parts of Docker Desktop.

If your extension uses `localStorage` to store data, other extensions running in Docker Desktop cannot access the local storage of your extension. The extension’s local storage is persisted even after Docker Desktop is stopped or restarted. When an extension is upgraded, its local storage is persisted, whereas when it is uninstalled, its local storage is completely removed.

## Cross site-scripting

CORS rules are enforced with same-origin policy enabled, so the extension UI code cannot load external scripts.
