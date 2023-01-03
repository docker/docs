---
title: "Test and debug"
description: Test and debug your extension.
keywords: Docker, Extensions, sdk, preview, update, Chrome DevTools
redirect_from:
- /desktop/extensions-sdk/build/test-debug/
---

In order to improve the developer experience, Docker Desktop provides a set of tools to help you test and debug your extension.

### Open Chrome DevTools

In order to open the Chrome DevTools for your extension when you click on the extension tab, run:

```console
$ docker extension dev debug <name-of-your-extensions>
```

Each subsequent click on the extension tab will also open Chrome Dev Tools. To stop this behaviour, run:

```console
$ docker extension dev reset <name-of-your-extensions>
```

After an extension is deployed, it is also possible to open Chrome DevTools from the UI extension part using a variation of the [Konami Code](https://en.wikipedia.org/wiki/Konami_Code){:target="_blank" rel="noopener" class="_"}. Click on the extension tab, and then hit the key sequence `up, up, down, down, left, right, left, right, p, d, t`.

### Hot reloading whilst developing the UI

During UI development, it’s helpful to use hot reloading to test your changes without rebuilding your entire extension. To do this, you can configure Docker Desktop to load your UI from a development server, such as the one Create React App starts when invoked with yarn start.

Assuming your app runs on the default port, start your UI app and then run:

```console
$ cd ui
$ npm run dev
```

This starts a development server that listens on port 5173.

You can now tell Docker Desktop to use this as the frontend source. In another terminal run:

```console
$ docker extension dev ui-source <name-of-your-extensions> http://localhost:5173
```

Close and reopen the Docker Desktop dashboard and go to your extension. All the changes to the frontend code are immediately visible.

Once finished, you can reset the extension configuration to the original settings. This will also reset opening Chrome dev tools if you used `docker extension dev debug <name-of-your-extensions>`:

```console
$ docker extension dev reset <name-of-your-extensions>
```

## Show the extension containers

If your extension is composed of one or more services running as containers in the Docker Desktop VM, you can access them easily from the dashboard in Docker Desktop.

1. In Docker Desktop, navigate to **Settings** or **Preferences** if you’re a Mac user.
2. Under the **Extensions** tab, select the **Show Docker Desktop Extensions system containers** option. You can now view your extension containers and their logs.


## Clean up

To remove the extension, run:

```console
$ docker extension rm <name-of-your-extension>
```

## What's next

- Build an [advanced frontend](../build/frontend-extension-tutorial.md) extension.
- Learn more about extensions [architecture](../architecture/index.md).
- Explore our [design principles](../design/design-principles.md).
- Take a look at our [UI styling guidelines](../design/index.md).
- Learn how to [publish your extension](../extensions/index.md).
