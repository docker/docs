---
title: Preview and update
description: Step five in the extension creation process
keywords: Docker, Extensions, sdk, preview, update, Chrome DevTools
---

Once your extension has been validated and installed, you can preview your extension in Docker Desktop. 

## Preview the extension

To preview the extension in Docker Desktop, close and open Docker Dashboard once the installation is complete.

The left-hand menu displays a new tab with the name of your extension. 
The example below shows the [`Min FrontEnd Extension`](set-up/minimal-frontend-extension.md). 

![minimal-frontend-extension](images/ui-minimal-extension.png)

### Open Chome DevTools**

In order to open the Chrome DevTools for your extension when you click on the extension tab, run:

`$ docker extension dev debug john/my-extension`

Each subsequent click on the extension tab will also open Chrome Dev Tools. To stop this behaviour, run:

`$ docker extension dev reset john/my-extension`

After an extension is deployed, it is also possible to open Chrome DevTools from the UI extension part using a variation of the [Konami Code](https://en.wikipedia.org/wiki/Konami_Code). Click on the extension tab, and then hit the key sequence `up, up, down, down, left, right, left, right, p, d, t`.

### Preview whilst developing the UI

If your extension has a UI, you can see it directly inside Docker Desktop whilst you develop it directly. For this you need to first install the extension. If you then run a development server locally, with `yarn start` for example, enter the following command:

`$ docker extension dev ui-source my-extension http://localhost:8080`

This changes the source of the extension UI to your local development server. Auto and hot-reload now work.

> Note
> 
> Make sure to reopen the Dashboard when you set a new source for the extension’s UI.

Once finished, you can reset the extension configuration to the original settings. This will also reset opening Chrome dev tools if you used `docker extension dev debug my-extension`:

`$ docker extension dev reset my-extension`

## Show the extension containers

If your extension is composed of one or more services running as containers in the Docker Desktop VM, you can access them easily from the dashboard in Docker Desktop.

1. In Docker Desktop, navigate to **Settings** or **Preferences** if you’re a Mac user.
2. Under the **Extensions** tab, select the **Show Docker Desktop Extensions system containers** option. You can now view your extension containers and their logs.


## Update the extension

To update the extension, you must first [rebuild](build.md) and [revalidate](validate-install.md) your extension. You can then use the update command.

`docker extension update <name-of-your-extensions>`

## What's next?

- Explore our [design principles](../design/design-principles.md).
- Take a look at our [UI styling guidelines](../design/overview.md).
- Set up [authentication for your extension](oauth2-flow.md)
- Learn how to [publish your extension](../extensions/Overview.md).
