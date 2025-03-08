---
description: Configuration overview for Docker Trusted Registry
keywords: docker, documentation, about, technology, understanding, enterprise, hub, registry
redirect_from:
- /docker-trusted-registry/configure/configuration/
title: Configuration overview
---

When you first install Docker Trusted Registry, you need to configure it. Use
this overview to see what you can configure.

To start, navigate to the Trusted Registry user interface (UI) > Settings, to
view configuration options. Configuring is grouped by the following:

* [General settings](config-general.md) (ports, proxies, and Notary)
* [Security settings](config-security.md)
* [Storage settings](config-storage.md)
* [License](../install/license.md)
* Updates


Saving changes you've made to settings will restart various services, as follows:

 * General settings: full Docker Trusted Registry restart
 * License change: full Docker Trusted Registry restart
 * SSL change: Nginx reload
 * Storage config: only registries restart

## Docker daemon logs

Both the Trusted Registry and the Docker daemon collect and store log messages.
To limit duplication of the Docker daemon logs, add the following parameters in
a Trusted Registry CLI to the Docker daemon and then restart the daemon.

`docker daemon --log-opt max-size 100m max-file=1`


## See also

* [Monitor DTR](../monitor-troubleshoot/index.md)
* [Troubleshoot DTR](../monitor-troubleshoot/troubleshoot.md)