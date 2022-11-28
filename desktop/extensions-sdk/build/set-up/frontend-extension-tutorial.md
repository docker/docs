---
title: Set up an advanced frontend extension
description: Advanced frontend extension tutorial
keywords: Docker, extensions, sdk, build
redirect_from:
  - /desktop/extensions-sdk/tutorials/react-extension/
  - /desktop/extensions-sdk/build/set-up/react-extension/
  - /desktop/extensions-sdk/build/set-up/minimal-frontend-using-docker-cli/
---

To start creating your extension, you first need a directory with files which range from the extension’s source code to the required extension-specific files. This page provides information on how to set up a simple Docker extension that contains only a UI part.

> Note
>
> Before you start, make sure you have installed the latest version of [Docker Desktop](https://www.docker.com/products/docker-desktop/).

## Extension folder structure

The quickest way to create a new extension is to run `docker extension init my-extension` as in the
[Quickstart](../../quickstart.md). This will create a new directory `my-extension` that contains a fully functional extension.

> **Tip**
>
> The `docker extension init` generates a React based extension. But you can still use it as a starting point for
> your own extension and use any other frontend framework, like Vue, Angular, Svelte, etc. or event stay with
> vanilla Javascript.
{: .tip }

Although you can start from an empty directory or from the `react-extension` [sample folder](https://github.com/docker/extensions-sdk/tree/main/samples){:target="_blank" rel="noopener" class="_"},
it's highly recommended that you start from the `docker extension init` command and change it to suit your needs.

```bash
.
├── Dockerfile # (1)
├── ui # (2)
│   ├── public # (3)
│   │   └── index.html
│   ├── src # (4)
│   │   ├── App.tsx
│   │   ├── index.tsx
│   ├── package.json
│   └── package-lock.lock
│   ├── tsconfig.json
├── docker.svg # (5)
└── metadata.json # (6)
```

1. Contains everything required to build the extension and run it in Docker Desktop.
2. High-level folder containing your front-end app source code.
3. Assets that aren’t compiled or dynamically generated are stored here. These can be static assets like logos or the robots.txt file.
4. The src, or source folder contains all the React components, external CSS files, and dynamic assets that are brought into the component files.
5. The icon that is displayed in the left-menu of the Docker Desktop Dashboard.
6. A file that provides information about the extension such as the name, description, and version.

## Adapting the Dockerfile

> **Note**
>
> When using the `docker extension init`, it creates a `Dockerfile` that already contains what is needed for a React
> extension.

Once the extension is created, you need to configure the `Dockerfile` to build the extension and configure the labels
that are used to populate the extension's card in the Marketplace. Here is an example of a `Dockerfile` for a React
extension:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#react-dockerfile" data-group="react">For React</a></li>
  <li><a data-toggle="tab" data-target="#vue-dockerfile" data-group="vue">For Vue</a></li>
  <li><a data-toggle="tab" data-target="#angular-dockerfile" data-group="angular">For Angular</a></li>
  <li><a data-toggle="tab" data-target="#svelte-dockerfile" data-group="svelte">For Svelte</a></li>
</ul>

<div class="tab-content">
  <div id="react-dockerfile" class="tab-pane fade in active" markdown="1">

```Dockerfile
FROM --platform=$BUILDPLATFORM node:18.9-alpine3.15 AS client-builder
WORKDIR /ui
# cache packages in layer
COPY ui/package.json /ui/package.json
COPY ui/package-lock.json /ui/package-lock.json
RUN --mount=type=cache,target=/usr/src/app/.npm \
    npm set cache /usr/src/app/.npm && \
    npm ci
# install
COPY ui /ui
RUN npm run build

FROM alpine
LABEL org.opencontainers.image.title="My extension" \
    org.opencontainers.image.description="Your Desktop Extension Description" \
    org.opencontainers.image.vendor="Awesome Inc." \
    com.docker.desktop.extension.api.version="0.3.0" \
    com.docker.desktop.extension.icon="https://www.docker.com/wp-content/uploads/2022/03/Moby-logo.png"
    com.docker.extension.screenshots="" \
    com.docker.extension.detailed-description="" \
    com.docker.extension.publisher-url="" \
    com.docker.extension.additional-urls="" \
    com.docker.extension.changelog=""

COPY metadata.json .
COPY docker.svg .
COPY --from=client-builder /ui/build ui

```

  </div>
  <div id="vue-dockerfile" class="tab-pane fade" markdown="1">

<br/>

> **Important**
>
> We don't have a working Dockerfile for Vue yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Vue)
> and let us know you'd like a Dockerfile for Vue.
{: .important }

  </div>
  <div id="angular-dockerfile" class="tab-pane fade" markdown="1">

<br/>

> **Important**
>
> We don't have a working Dockerfile for Angular yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Angular)
> and let us know you'd like a Dockerfile for Angular.
{: .important }

  </div>
  <div id="svelte-dockerfile" class="tab-pane fade" markdown="1">

<br/>

> **Important**
>
> We don't have a working Dockerfile for Svelte yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Svelte)
> and let us know you'd like a Dockerfile for Svelte.
{: .important }

  </div>
</div>

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

## Use the Extension APIs client

To use the Extension APIs and perform actions with Docker Desktop, the extension must first import the 
`@docker/extension-api-client` library. To install it, run the command below:

```bash
npm install @docker/extension-api-client
```

Then call the `createDockerDesktopClient` function to create a client object to call the extension APIs.

```js
import { createDockerDesktopClient } from '@docker/extension-api-client';

const ddClient = createDockerDesktopClient();
```

When using Typescript, you can also install `@docker/extension-api-client-types` as a dev dependency. This will 
provide you with type definitions for the extension APIs and auto-completion in your IDE.

```bash
npm install @docker/extension-api-client-types --save-dev
```

![types auto complete](images/types-autocomplete.png)

For example, you can use the `docker.cli.exec` function to get the list of all the containers via the `docker ps --all` 
command and display the result in a table.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#react-app" data-group="react">React</a></li>
  <li><a data-toggle="tab" data-target="#vue-app" data-group="vue">Vue</a></li>
  <li><a data-toggle="tab" data-target="#angular-app" data-group="angular">Angular</a></li>
  <li><a data-toggle="tab" data-target="#svelte-app" data-group="svelte">Svelte</a></li>
</ul>

<div class="tab-content">
  <div id="react-app" class="tab-pane fade in active" markdown="1">

Replace the `ui/src/App.tsx` file with the following code:

```tsx
{% raw %}
// ui/src/App.tsx
import React, { useEffect } from 'react';
import {
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography
} from "@mui/material";
import { createDockerDesktopClient } from "@docker/extension-api-client";

//obtain docker destkop extension client
const ddClient = createDockerDesktopClient();

export function App() {
  const [containers, setContainers] = React.useState([]);

  useEffect(() => {
    // List all containers
    ddClient.docker.cli.exec('ps', ['--all', '--format', '"{{json .}}"']).then((result) => {
      // result.parseJsonLines() parses the output of the command into an array of objects
      setContainers(result.parseJsonLines());
    });
  }, []);

  return (
    <Stack>
      <Typography data-testid="heading" variant="h3" role="title">
        Container list
      </Typography>
      <Typography
      data-testid="subheading"
      variant="body1"
      color="text.secondary"
      sx={{ mt: 2 }}
    >
      Simple list of containers using Docker Extensions SDK.
      </Typography>
      <TableContainer sx={{mt:2}}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Container id</TableCell>
              <TableCell>Image</TableCell>
              <TableCell>Command</TableCell>
              <TableCell>Created</TableCell>
              <TableCell>Status</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {containers.map((container) => (
              <TableRow
                key={container.ID}
                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
              >
                <TableCell>{container.ID}</TableCell>
                <TableCell>{container.Image}</TableCell>
                <TableCell>{container.Command}</TableCell>
                <TableCell>{container.CreatedAt}</TableCell>
                <TableCell>{container.Status}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Stack>
  );
}
{% endraw %}
```

  </div>
  <div id="vue-app" class="tab-pane fade" markdown="1">
    
<br/>

> **Important**
>
> We don't have an example for Vue yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Vue)
> and let us know you'd like a sample with Vue.
{: .important }
  
  </div>
  <div id="angular-app" class="tab-pane fade" markdown="1">

<br/>

> **Important**
>
> We don't have an example for Angular yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Angular)
> and let us know you'd like a sample with Angular.
{: .important }

  </div>
  <div id="svelte-app" class="tab-pane fade" markdown="1">

<br/>

> **Important**
>
> We don't have an example for Svelte yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Svelte)
> and let us know you'd like a sample with Svelte.
{: .important }

  </div>
</div>

![Screenshot of the container list.](images/react-extension.png)

## What's next?

- Learn how to [build and install your extension](../build-install.md).
- For more information and guidelines on building the UI, see the [Design and UI styling section](../../design/design-guidelines.md).
- If you want to set up user authentication for the extension, see [Authentication](../../guides/oauth2-flow.md).

