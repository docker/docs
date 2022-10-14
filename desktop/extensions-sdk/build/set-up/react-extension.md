---
title: Set up a minimal react extension
description: Minimal react extension tutorial
keywords: Docker, extensions, sdk, build
redirect_from:
  - /desktop/extensions-sdk/tutorials/react-extension/
---

To start creating your extension, you first need a directory with files which range from the extension’s source code to the required extension-specific files. This page provides information on how to set up a simple Docker extension that contains only a UI part and is based on ReactJS.

> Note
>
> Before you start, make sure you have installed the latest version of [Docker Desktop](../../../release-notes.md).

> Note
>
> If you want to start a codebase for your new extension, our [Quickstart guide](../../quickstart.md) and `docker extension init <my-extension>` will provide a better base for your extension, more up-to-date and related to your install of Docker Desktop.

## Extension folder structure

In the `react-extension` [sample folder](https://github.com/docker/extensions-sdk/tree/main/samples), you can find a ready-to-go example that represents a UI Extension built on ReactJS. We will go through this code example in this tutorial.

Although you can start from an empty directory, it is highly recommended that you start from the template below and change it accordingly to suit your needs.

```bash
.
├── Dockerfile # (1)
├── client # (2)
│   ├── package.json
│   ├── public # (3)
│   │   └── index.html
│   ├── src # (4)
│   │   ├── App.tsx
│   │   ├── globals.d.ts
│   │   ├── index.tsx
│   │   └── react-app-env.d.ts
│   ├── tsconfig.json
│   └── yarn.lock
├── docker.svg # (5)
└── metadata.json # (6)
```

1. Contains everything required to build the extension and run it in Docker Desktop.
2. High-level folder containing your front-end app source code.
3. Assets that aren’t compiled or dynamically generated are stored here. These can be static assets like logos or the robots.txt file.
4. The src, or source folder contains all the React components, external CSS files, and dynamic assets that are brought into the component files.
5. The icon that is displayed in the left-menu of the Docker Desktop Dashboard.
6. A file that provides information about the extension such as the name, description, and version.

For more information and guidelines on building the UI, see the [Design and UI styling section](../../design/design-guidelines.md).

If you want to set up user authentication for the extension, see [Authentication](../../dev/oauth2-flow.md).

## Create a Dockerfile

Use the Dockerfile below as a template and change it accordingly to suit your needs.

```Dockerfile
FROM node:14.17-alpine3.13 AS client-builder
WORKDIR /app/client
# cache packages in layer
COPY client/package.json /app/client/package.json
COPY client/yarn.lock /app/client/yarn.lock
ARG TARGETARCH
RUN yarn config set cache-folder /usr/local/share/.cache/yarn-${TARGETARCH}
RUN --mount=type=cache,target=/usr/local/share/.cache/yarn-${TARGETARCH} yarn
# install
COPY client /app/client
RUN --mount=type=cache,target=/usr/local/share/.cache/yarn-${TARGETARCH} yarn build

FROM debian:bullseye-slim
LABEL org.opencontainers.image.title="ui-extension" \
    org.opencontainers.image.description="Your Desktop Extension Description" \
    org.opencontainers.image.vendor="Docker Inc." \
    com.docker.desktop.extension.api.version="1.0.0-beta.1" \
    com.docker.desktop.extension.icon="https://www.docker.com/wp-content/uploads/2022/03/Moby-logo.png"

COPY --from=client-builder /app/client/dist ui
COPY docker.svg .
COPY metadata.json .

```

## Configure the metadata file

A `metadata.json` file is required at the root of your extension directory.

```json
{
  "icon": "docker.svg",
  "ui": {
    "dashboard-tab": {
      "title": "UI Extension",
      "root": "/ui",
      "src": "index.html"
    }
  }
}
```

## Use extension APIs in the application code

The React application can import `@docker/extension-api-client` and use extension APIs to perform actions with Docker Desktop.

```ts
import { Box, Button } from '@mui/material';
import { createDockerDesktopClient } from '@docker/extension-api-client';

export function App() {
  //obtain docker destkop extension client
  const ddClient = createDockerDesktopClient();

  function sayHello() {
    ddClient.desktopUI.toast.success('Hello, World!');
  }

  ...
}
```

For more information on the `metadata.json`, see [Metadata](../../extensions/METADATA.md).

## What's next?

Learn how to [build and install your extension](../build-install.md).
