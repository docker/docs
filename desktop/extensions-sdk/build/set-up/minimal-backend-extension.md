---
title: Set up a minimal backend extension
description: Minimal backend extension tutorial
keywords: Docker, extensions, sdk, build
redirect_from:
  - /desktop/extensions-sdk/tutorials/minimal-backend-extension/
---

To start creating your extension, you first need a directory with files which range from the extension’s source code to the required extension-specific files. This page provides information on how to set up a simple Docker Extension that runs CLI commands in the backend.

For extensions with a backend service running REST services over sockets or named pipes, see the `vm-ui extension` [sample](https://github.com/docker/extensions-sdk/tree/main/samples).

> Note
>
> Before you start, make sure you have installed the latest version of [Docker Desktop](../../../release-notes.md).

> Note
>
> If you want to start a codebase for your new extension, our [Quickstart guide](../../quickstart.md) and `docker extension init <my-extension>` will provide a better base for your extension, more up-to-date and related to your install of Docker Desktop.

## Extension folder structure

In the `minimal-backend` [sample folder](https://github.com/docker/extensions-sdk/tree/main/samples), you can find a ready-to-go example that represents a UI extension built on HTML that runs a backend service. We will go through this code example in this tutorial.

Although you can start from an empty directory, it is highly recommended that you start from the template below and change it accordingly to suit your needs.

```bash
.
├── Dockerfile # (1)
├── Makefile
├── client # (2)
│   └── src
│       ├── App.tsx
│       └── ... React aplication
├── hello.sh # (3)
└── metadata.json # (4)
└── ui # (5)
    └── index.html
```

1. Contains everything required to build the extension and run it in Docker Desktop.
2. The source folder that contains the UI application. In this example we use a React frontend, the main part of th extension is an App.tsx.
3. The script that runs inside the container.
4. A file that provides information about the extension such as the name, description, and version.
5. The source folder that contains all your HTML, CSS and JS files. There can also be other static assets such as logos and icons. For more information and guidelines on building the UI, see the [Design and UI styling section](../../design/design-guidelines.md).

If you want to set up user authentication for the extension, see [Authentication](../../dev/oauth2-flow.md).

## Invoke the extension backend from your javascript code

Using the [React extension example](./react-extension.md), we can invoke our extension backend from the App.tsx file.

Use the Docker Desktop Client object and then invoke a binary provided in our backend container (that lives inside the Docker Desktop VM) with `ddClient.docker.extension.vm.cli.exec()`.
In our example, our hello.sh script returns a string as result, we obtain it with `result?.stdout`.

```typescript
const ddClient = createDockerDesktopClient();
const [backendInfo, setBackendInfo] = useState<string | undefined>();

async function runExtensionBackend(inputText: string) {
  const result = await ddClient.extension.vm?.cli.exec("./hello.sh", [
    inputText,
  ]);
  setBackendInfo(result?.stdout);
}
```

## Create a Dockerfile

At minimum, your Dockerfile needs:

- Labels which provide extra information about the extension.
- The source code which in this case is an `index.html` that sits within the `ui` folder. `index.html` refers to javascript code in `script.js`.
- The `metadata.json` file.

```Dockerfile
FROM node:17.7-alpine3.14 AS client-builder
# ... build React application

FROM alpine:3.15

LABEL org.opencontainers.image.title="HelloBackend" \
    org.opencontainers.image.description="A sample extension that runs a shell script inside a container's Desktop VM." \
    org.opencontainers.image.vendor="Docker Inc." \
    com.docker.desktop.extension.api.version="1.0.0-beta.1" \
    com.docker.desktop.extension.icon="https://www.docker.com/wp-content/uploads/2022/03/Moby-logo.png"

COPY hello.sh .
COPY metadata.json .
COPY --from=client-builder /app/client/dist ui

CMD [ "sleep", "infinity" ]
```

## Configure the metadata file

A `metadata.json` file is required at the root of the image filesystem.

```json
{
  "vm": {
    "image": "${DESKTOP_PLUGIN_IMAGE}"
  },
  "ui": {
    "dashboard-tab": {
      "title": "Hello Backend Extension",
      "root": "/ui",
      "src": "index.html"
    }
  }
}
```

For more information on the `metadata.json`, see [Metadata](../../extensions/METADATA.md).

> **Warning**
>
> Do not replace the `${DESKTOP_PLUGIN_IMAGE}` placeholder in the `metadata.json` file. The placeholder is replaced automatically with the correct image name when the extension is installed.
> {: .warning}

## What's next?

Learn how to [build and install your extension](../build-install.md).
