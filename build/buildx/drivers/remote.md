---
title: "Remote driver"
keywords: build, buildx, driver, builder, remote
fetch_remote:
  line_start: 2
  line_end: -1
---

> Beta
>
> Remote driver is currently available as a beta feature. We recommend that you
> do not use this feature in production environments. You can [build Buildx from source](https://github.com/docker/buildx#building){: target="_blank" rel="noopener" class="_"}
> to test the remote driver or use the following command to download and
> install an edge release of Buildx:
> 
> ```console
> $ echo "FROM docker/buildx-bin:master" | docker buildx build --platform=local --output . -f - .
> $ mkdir -p ~/.docker/cli-plugins/
> $ mv buildx ~/.docker/cli-plugins/docker-buildx
> ```
{: .important}

