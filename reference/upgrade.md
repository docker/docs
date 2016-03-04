+++
title = "upgrade"
description = "Upgrade UCP controller"
keywords= ["upgrade, ucp"]
[menu.main]
identifier = "ucp_upgrade"
parent = "ucp_ref"
+++

# upgrade

Upgrade the UCP components on this Docker Engine.

```bash
docker run --rm -it \
      --name ucp \
      -v /var/run/docker.sock:/var/run/docker.sock \
      docker/ucp \
      upgrade [OPTIONS]
```

## Description

When upgrading UCP, you must run the `upgrade` command against every
Engine in your cluster.  You should upgrade your controller and replica
nodes first, followed by your compute nodes.  If you plan to upgrade your
Engine as well, upgrade the Engine first, before upgrading UCP on the node.

After upgrading each node, confirm the node is present in the UCP console
before proceeding to the next node.


## Options

| Option                 | Description                                                                      |
|:-----------------------|:---------------------------------------------------------------------------------|
| ` --debug`, `-D`       | Enable debug.                                                                    |
| ` --jsonlog`           | Produce json formatted output for easier parsing.                                |
| ` --interactive`, `-i` | Enable interactive mode.,You will be prompted to enter all required information. |
| ` --id`                | The ID of the UCP instance to upgrade.                                           |
| ` --pull "always"`     | Specify image pull behavior (`always`, when `missing`, or `never`).              |
