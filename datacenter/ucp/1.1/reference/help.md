---
description: Shows a list of available commands for Docker Universal Control Plane.
keywords: help, ucp
title: docker/ucp help
---

Docker Universal Control Plane Tool

This tool has commands to 'install' the UCP initial controller and
'join' nodes to that controller. The tool can also 'uninstall' the product.
This tool must run as a container with a well-known name and with the
docker.sock volume mounted, which you can cut-and-paste from the usage
example below.

This tool will generate TLS certificates and will attempt to determine
your hostname and primary IP addresses. This may be overridden with the
'--host-address' option. The tool may not discover your
externally visible fully qualified hostname. For proper certificate
verification, you should pass one or more Subject Alternative Names with
'--san' during 'install' and 'join' that matches the fully qualified
hostname you intend to use to access the given system.

Many settings can be passed as flags or environment variables. When passing as
an environment variable, use the 'docker run -e VARIABLE_NAME ...' syntax to
pass the value from your shell, or 'docker run -e VARIABLE_NAME=value ...' to
specify the value explicitly on the command line.

Additional help is available for each command with the '--help' option.

## Usage

```bash
docker run --rm -it \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    command [arguments...]
```