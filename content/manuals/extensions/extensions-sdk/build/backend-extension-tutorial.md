---
title: Add a backend to your extension
description: Learn how to add a backend to your extension.
keywords: Docker, extensions, sdk, build
aliases:
 - /desktop/extensions-sdk/tutorials/minimal-backend-extension/
 - /desktop/extensions-sdk/build/minimal-backend-extension/
 - /desktop/extensions-sdk/build/set-up/backend-extension-tutorial/
 - /desktop/extensions-sdk/build/backend-extension-tutorial/
---

Your extension can ship a backend part with which the frontend can interact with. This page provides information on why and how to add a backend.

Before you start, make sure you have installed the latest version of [Docker Desktop](https://www.docker.com/products/docker-desktop/).

> Tip
>
> Check the [Quickstart guide](../quickstart.md) and `docker extension init <my-extension>`. They provide a better base for your extension as it's more up-to-date and related to your install of Docker Desktop.

## Why add a backend?

Thanks to the Docker Extensions SDK, most of the time you should be able to do what you need from the Docker CLI
directly from [the frontend](frontend-extension-tutorial.md#use-the-extension-apis-client).

Nonetheless, there are some cases where you might need to add a backend to your extension. So far, extension
builders have used the backend to:
- Store data in a local database and serve them back with a REST API.
- Store the extension state, for example when a button starts a long-running process, so that if you navigate away from the extension user interface and comes back, the frontend can pick up where it left off.

For more information about extension backends, see [Architecture](../architecture/index.md#the-backend).

## Add a backend to the extension

If you created your extension using the `docker extension init` command, you already have a backend setup. Otherwise, you have to first create a `vm` directory that contains the code and updates the Dockerfile to
containerize it.

Here is the extension folder structure with a backend:

```bash
.
├── Dockerfile # (1)
├── Makefile
├── metadata.json
├── ui
    └── index.html
└── vm # (2)
    ├── go.mod
    └── main.go
```

1. Contains everything required to build the backend and copy it in the extension's container filesystem.
2. The source folder that contains the backend code of the extension.

Although you can start from an empty directory or from the `vm-ui extension` [sample](https://github.com/docker/extensions-sdk/tree/main/samples),
it is highly recommended that you start from the `docker extension init` command and change it to suit your needs.

> [!TIP]
>
> The `docker extension init` generates a Go backend. But you can still use it as a starting point for
> your own extension and use any other language like Node.js, Python, Java, .Net, or any other language and framework.

In this tutorial, the backend service simply exposes one route that returns a JSON payload that says "Hello".

```json
{ "Message": "Hello" }
```

> [!IMPORTANT]
>
> We recommend that, the frontend and the backend communicate through sockets, and named pipes on Windows, instead of
> HTTP. This prevents port collision with any other running application or container running
> on the host. Also, some Docker Desktop users are running in constrained environments where they
> can't open ports on their machines. When choosing the language and framework for your backend, make sure it
> supports sockets connection.

{{< tabs group="lang" >}}
{{< tab name="Go" >}}

```go
package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func main() {
	var socketPath string
	flag.StringVar(&socketPath, "socket", "/run/guest/volumes-service.sock", "Unix domain socket to listen on")
	flag.Parse()

	os.RemoveAll(socketPath)

	logrus.New().Infof("Starting listening on %s\n", socketPath)
	router := echo.New()
	router.HideBanner = true

	startURL := ""

	ln, err := listen(socketPath)
	if err != nil {
		log.Fatal(err)
	}
	router.Listener = ln

	router.GET("/hello", hello)

	log.Fatal(router.Start(startURL))
}

func listen(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}

func hello(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, HTTPMessageBody{Message: "hello world"})
}

type HTTPMessageBody struct {
	Message string
}
```

{{< /tab >}}
{{< tab name="Node" >}}

> [!IMPORTANT]
>
> We don't have a working example for Node yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.25798127=Node)
> and let us know if you'd like a sample for Node.

{{< /tab >}}
{{< tab name="Python" >}}

> [!IMPORTANT]
>
> We don't have a working example for Python yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.25798127=Python)
> and let us know if you'd like a sample for Python.

{{< /tab >}}
{{< tab name="Java" >}}

> [!IMPORTANT]
>
> We don't have a working example for Java yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.25798127=Java)
> and let us know if you'd like a sample for Java.

{{< /tab >}}
{{< tab name=".NET" >}}

> [!IMPORTANT]
>
> We don't have a working example for .NET. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.25798127=.Net)
> and let us know if you'd like a sample for .NET.

{{< /tab >}}
{{< /tabs >}}

## Adapt the Dockerfile

> [!NOTE]
>
> When using the `docker extension init`, it creates a `Dockerfile` that already contains what is needed for a Go backend.

{{< tabs group="lang" >}}
{{< tab name="Go" >}}

To deploy your Go backend when installing the extension, you need first to configure the `Dockerfile`, so that it:
- Builds the backend application
- Copies the binary in the extension's container filesystem
- Starts the binary when the container starts listening on the extension socket

> [!TIP]
> 
> To ease version management, you can reuse the same image to build the frontend, build the
backend service, and package the extension.

```dockerfile
# syntax=docker/dockerfile:1
FROM node:17.7-alpine3.14 AS client-builder
# ... build frontend application

# Build the Go backend
FROM golang:1.17-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /backend
COPY vm/go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY vm/. .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o bin/service

FROM alpine:3.15
# ... add labels and copy the frontend application

COPY --from=builder /backend/bin/service /
CMD /service -socket /run/guest-services/extension-allthethings-extension.sock
```

{{< /tab >}}
{{< tab name="Node" >}}

> [!IMPORTANT]
>
> We don't have a working Dockerfile for Node yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.25798127=Node)
> and let us know if you'd like a Dockerfile for Node.

{{< /tab >}}
{{< tab name="Python" >}}

> [!IMPORTANT]
>
> We don't have a working Dockerfile for Python yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.25798127=Python)
> and let us know if you'd like a Dockerfile for Python.

{{< /tab >}}
{{< tab name="Java" >}}

> [!IMPORTANT]
>
> We don't have a working Dockerfile for Java yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.25798127=Java)
> and let us know if you'd like a Dockerfile for Java.

{{< /tab >}}
{{< tab name=".NET" >}}

> [!IMPORTANT]
>
> We don't have a working Dockerfile for .Net. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.25798127=.Net)
> and let us know if you'd like a Dockerfile for .Net.

{{< /tab >}}
{{< /tabs >}}

## Configure the metadata file

To start the backend service of your extension inside the VM of Docker Desktop, you have to configure the image name
in the `vm` section of the `metadata.json` file.

```json
{
  "vm": {
    "image": "${DESKTOP_PLUGIN_IMAGE}"
  },
  "icon": "docker.svg",
  "ui": {
    ...
  }
}
```

For more information on the `vm` section of the `metadata.json`, see [Metadata](../architecture/metadata.md).

> [!WARNING]
>
> Do not replace the `${DESKTOP_PLUGIN_IMAGE}` placeholder in the `metadata.json` file. The placeholder is replaced automatically with the correct image name when the extension is installed.

## Invoke the extension backend from your frontend

Using the [advanced frontend extension example](frontend-extension-tutorial.md), we can invoke our extension backend.

Use the Docker Desktop Client object and then invoke the `/hello` route from the backend service with `ddClient.
extension.vm.service.get` that returns the body of the response.

{{< tabs group="framework" >}}
{{< tab name="React" >}}

Replace the `ui/src/App.tsx` file with the following code:

```tsx

// ui/src/App.tsx
import React, { useEffect } from 'react';
import { createDockerDesktopClient } from "@docker/extension-api-client";

//obtain docker destkop extension client
const ddClient = createDockerDesktopClient();

export function App() {
  const ddClient = createDockerDesktopClient();
  const [hello, setHello] = useState<string>();

  useEffect(() => {
    const getHello = async () => {
      const result = await ddClient.extension.vm?.service?.get('/hello');
      setHello(JSON.stringify(result));
    }
    getHello()
  }, []);

  return (
    <Typography>{hello}</Typography>
  );
}

```

{{< /tab >}}
{{< tab name="Vue" >}}

> [!IMPORTANT]
>
> We don't have an example for Vue yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Vue)
> and let us know if you'd like a sample with Vue.

{{< /tab >}}
{{< tab name="Angular" >}}

> [!IMPORTANT]
>
> We don't have an example for Angular yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Angular)
> and let us know if you'd like a sample with Angular.

{{< /tab >}}
{{< tab name="Svelte" >}}

> [!IMPORTANT]
>
> We don't have an example for Svelte yet. [Fill out the form](https://docs.google.com/forms/d/e/1FAIpQLSdxJDGFJl5oJ06rG7uqtw1rsSBZpUhv_s9HHtw80cytkh2X-Q/viewform?usp=pp_url&entry.1333218187=Svelte)
> and let us know if you'd like a sample with Svelte.

{{< /tab >}}
{{< /tabs >}}

## Re-build the extension and update it

Since you have modified the configuration of the extension and added a stage in the Dockerfile, you must re-build the extension.

```bash
docker build --tag=awesome-inc/my-extension:latest .
```

Once built, you need to update it, or install it if you haven't already done so.

```bash
docker extension update awesome-inc/my-extension:latest
```

Now you can see the backend service running in the **Containers8* tab of the Docker Dashboard and watch the logs when you need to debug it.

> [!TIP]
>
> You may need to turn on the **Show system containers** option in **Settings** to see the backend container running.
> See [Show extension containers](../dev/test-debug.md#show-the-extension-containers) for more information.

Open Docker Dashboard and select the **Containers** tab. You should see the response from the backend service
call displayed.

## What's next?

- Learn how to [share and publish your extension](../extensions/index.md).
- Learn more about extensions [architecture](../architecture/index.md).
