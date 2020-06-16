---
description: Receive real time events from containers.
keywords: fig, composition, compose, docker, orchestration, cli, events
title: docker-compose events
notoc: true
---

```
Usage: events [options] [SERVICE...]

Options:
    --json      Output events as a stream of json objects
```

Stream container events for every container in the project.

With the `--json` flag, a json object is printed one per line with the
format:

```json
{
    "time": "2015-11-20T18:01:03.615550",
    "type": "container",
    "action": "create",
    "id": "213cf7...5fc39a",
    "service": "web",
    "attributes": {
        "name": "application_web_1",
        "image": "alpine:edge"
    }
}
```

The events that can be received using this can be seen [here](../../engine/reference/commandline/events.md#object-types).
