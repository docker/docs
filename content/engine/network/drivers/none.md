---
title: None network driver
description: How to isolate the networking stack of a container using the none driver
keywords: network, none, standalone
aliases:
- /network/none/
---

If you want to completely isolate the networking stack of a container, you can
use the `--network none` flag when starting the container. Within the container,
only the loopback device is created.

The following example shows the output of `ip link show` in an `alpine`
container using the `none` network driver.

```console
$ docker run --rm --network none alpine:latest ip link show
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
```

No IPv6 loopback address is configured for containers using the `none` driver.

```console
$ docker run --rm --network none --name no-net-alpine alpine:latest ip addr show
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
```

## Next steps

- Go through the [host networking tutorial](/engine/network/tutorials/host.md)
- Learn about [networking from the container's point of view](../_index.md)
- Learn about [bridge networks](bridge.md)
- Learn about [overlay networks](overlay.md)
- Learn about [Macvlan networks](macvlan.md)
