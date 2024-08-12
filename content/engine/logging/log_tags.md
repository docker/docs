---
description: Learn about how to format log output with Go templates
keywords: docker, logging, driver, syslog, Fluentd, gelf, journald
title: Customize log driver output
aliases:
  - /engine/reference/logging/log_tags/
  - /engine/admin/logging/log_tags/
  - /config/containers/logging/log_tags/
---

The `tag` log option specifies how to format a tag that identifies the
container's log messages. By default, the system uses the first 12 characters of
the container ID. To override this behavior, specify a `tag` option:

```console
$ docker run --log-driver=fluentd --log-opt fluentd-address=myhost.local:24224 --log-opt tag="mailer"
```

Docker supports some special template markup you can use when specifying a tag's value:

| Markup             | Description                                          |
| ------------------ | ---------------------------------------------------- |
| `{{.ID}}`          | The first 12 characters of the container ID.         |
| `{{.FullID}}`      | The full container ID.                               |
| `{{.Name}}`        | The container name.                                  |
| `{{.ImageID}}`     | The first 12 characters of the container's image ID. |
| `{{.ImageFullID}}` | The container's full image ID.                       |
| `{{.ImageName}}`   | The name of the image used by the container.         |
| `{{.DaemonName}}`  | The name of the Docker program (`docker`).           |

For example, specifying a `--log-opt tag="{{.ImageName}}/{{.Name}}/{{.ID}}"` value yields `syslog` log lines like:

```none
Aug  7 18:33:19 HOSTNAME hello-world/foobar/5790672ab6a0[9103]: Hello from Docker.
```

At startup time, the system sets the `container_name` field and `{{.Name}}` in
the tags. If you use `docker rename` to rename a container, the new name isn't
reflected in the log messages. Instead, these messages continue to use the
original container name.
