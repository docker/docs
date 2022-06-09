---
title: Overview
description: Docker extension development
keywords: Docker, extensions, sdk, development
---

The section below describes how to get started developing your custom Docker extension.

Extensions can be composed of a visual part displayed in the Docker Desktop dashboard, and one or more optional services that run inside the Docker Desktop VM.

If you intend to develop an extension which consists exclusively of a visual part with no services running in the VM, see the [React extension](../tutorials/react-extension.md) tutorial.

If your extension needs to invoke docker commands, see the [Docker cli extension](../tutorials/minimal-frontend-using-docker-cli.md) tutorial.

If your extension requires additional services running in the Docker Desktop VM, see the [VM UI](https://github.com/docker/extensions-sdk/tree/main/samples/vm-service) example.

For further inspiration, see the other examples in the [samples folder](https://github.com/docker/extensions-sdk/tree/main/samples)

### Open Dev Tools

In order to open the Chrome Dev Tools for your extension when you click on the extension tab, run:

```console
$ docker extension dev debug john/my-extension
```

Each subsequent click on the extension tab will also open Chrome Dev Tools.
To stop this behaviour, run:

```console
$ docker extension dev reset john/my-extension
```

After an extension is deployed, it is also possible to open the Chrome Dev Tools from the UI extension part using a variation of the [Konami Code](https://en.wikipedia.org/wiki/Konami_Code).
Click on the extension tab, and then hit the key sequence `up, up, down, down, left, right, left, right, p, d, t`.

### Develop the Extension UI

If your extension has a UI, you can see it directly inside Docker Desktop whilst you develop it directly.
For this you need to first install the extension.
If you then run a development server locally, with `yarn start` for example, enter the following command:

```console
$ docker extension dev ui-source john/my-extension http://localhost:8080
```

This changes the source of the extension UI to your local development server. Auto and hot-reload now work.

> Make sure to reopen the Dashboard when you set a new source for the extension's UI.

Once finished, you can reset the extension configuration to the original settings. This will also reset opening Chrome dev tools if you used `docker extension dev debug my-extension`:

```console
$ docker extension dev reset john/my-extension
```

## Show the extension containers

If your extension is composed of one or more services running as containers in the Docker Desktop VM, you can access them easily from the dashboard in Docker Desktop.

1. In Docker Desktop, navigate to **Settings**, or **Preferences** if you're a Mac user.
2. Under the **Extensions** tab, select the **Show Docker Desktop Extensions system containers** option. You can now view your extension containers and their logs.
