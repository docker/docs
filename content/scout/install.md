---
title: Install Docker Scout
description: Installation instructions for the Docker Scout CLI plugin
keywords: scout, cli, install, download
---

The Docker Scout CLI plugin comes pre-installed with Docker Desktop.

If you run Docker Engine without Docker Desktop,
Docker Scout doesn't come pre-installed,
but you can install it as a standalone binary.

To install the latest version of the plugin, run the following commands:

```console
$ curl -fsSL https://raw.githubusercontent.com/docker/scout-cli/main/install.sh -o install-scout.sh
$ sh install-scout.sh
```

> **Note**
>
> Always examine scripts downloaded from the internet before running them
> locally. Before installing, make yourself familiar with potential risks and
> limitations of the convenience script.

If you want to install the plugin manually, you can find full instructions
and links to download in the [scout-cli repository](https://github.com/docker/scout-cli).

The Docker Scout CLI plugin is also available as [a container image](https://hub.docker.com/r/docker/scout-cli)
and as [a GitHub action](https://github.com/docker/scout-action).
