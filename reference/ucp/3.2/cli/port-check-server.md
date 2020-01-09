---
title: docker/ucp port-check-server
description: Check the firewall ports for UCP
keywords: ucp, cli, images
---

>{% include enterprise_label_shortform.md %}

Checks the suitablility of the node for a UCP installation

## Usage

```
docker run --rm -it \
     -v /var/run/docker.sock:/var/run/docker.sock \
     docker/ucp \
     port-check-server [command options]
```

## Description

Checks the suitablility of the node for a UCP installation

## Options

| Option                        | Description                       |
|:------------------------------|:----------------------------------|
| --listen-address -l *value*   | Listen Address (default: ":2376") |
