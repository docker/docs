# Orca Quick Start Guide

For all the gory details about how our installer works, check out
[install spec](install_upgrade_spec.md) but you came here to get up and
running quickly, so lets dive right in!

# Prerequisites

* You'll need access to the dockerorca images on hub - if the link below gives an error, ask someone on the **#orca** slack channel to give you access.
    * [https://hub.docker.com/r/dockerorca/orca-bootstrap/](https://hub.docker.com/r/dockerorca/orca-bootstrap/)
* You'll need at leasts one docker engine (local or remote should work)
    * If you want to build a multi-node deployment, **all** the nodes must be able to see eachother, so make sure if you're using remote engines, they're all on the same remote network.
* Orca installs its own Swarm, so don't set up Swarm first, just install Orca directly on your engine(s)


# Initial Installation
(line wrapped for readability)
```bash
docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name orca-bootstrap \
    dockerorca/orca-bootstrap \
    install --swarm-port 3376 -i
```

The above command will prompt you for some basic information to get Orca
stood up.  You can use "install --help" as the last line above to get
information about various options to the installer.

Important notes for first time users:
* We try to get the hostname/IPs right, but NAT can lead us astray.  Make sure you specify the **real** external hostname when prompted (or use --san for non-interactive mode)
* The last line of the installer output tells you where to go log in.
* If you didn't enter an admin password, the default login is "admin/orca"
* The first thing you probably want to do is download a cert bundle (upper right corner of UI, but subject to change)
    * With this, you can run docker CLI commands against Orca (and/or the swarm/engines if you're account is an admin account)
    * Take a look at the env.sh within the zip file for instructions (should be familiar if you've used machine)


# Adding Nodes to the Cluster
To add capacity to your cluster, run the following on the engine you want to add (**not** the engine where you installed Orca above)
```bash
docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name orca-bootstrap \
    dockerorca/orca-bootstrap \
    join -i
```

As with install, you can use "join --help" for more information.


# Uninstalling
The installer can also uninstall the Orca software on either the primary
Orca node, as well as the secondary Orca nodes.  Run the following directly
against the engine you want to uninstall (**not** against Orca itself)

```bash
docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name orca-bootstrap \
    dockerorca/orca-bootstrap \
    uninstall
```

As above, use "uninstall --help" to see what other flags are available to tune behavior.

# Upgrading

**Coming soon!** (before GA)

For now, just uninstall and re-install.  Your containers will remain on the engines and survive across the uninstall/install.
