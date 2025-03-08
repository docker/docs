---
title: Enable IPv6 support
description: How to enable IPv6 support in the Docker daemon
keywords: daemon, network, networking, ipv6
redirect_from:
- /engine/userguide/networking/default_network/ipv6/
---

Before you can use IPv6 in Docker containers or swarm services, you need to
enable IPv6 support in the Docker daemon. Afterward, you can choose to use
either IPv4 or IPv6 (or both) with any container, service, or network.

> **Note**: IPv6 networking is only supported on Docker daemons running on Linux
> hosts.

1.  Edit `/etc/docker/daemon.json` and set the `ipv6` key to `true`.

    ```json
    {
      "ipv6": true
    }
    ```

    Save the file.

2.  Reload the Docker configuration file.

    ```bash
    $ systemctl reload docker
    ```

You can now create networks with the `--ipv6` flag and assign containers IPv6
addresses using the `--ip6` flag.

## Next steps

- [Networking overview](/network/index.md)
- [Container networking](/config/containers/container-networking.md)
