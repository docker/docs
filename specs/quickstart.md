+++
draft = "true"
+++

# UCP Quick Start Guide

For all the gory details about how our installer works, check out
[install spec](install_upgrade_spec.md) but you came here to get up and
running quickly, so lets dive right in!

# Prerequisites

* You'll need at leasts one docker engine (local or remote should work)
    * If you want to build a multi-node deployment, **all** the nodes must be able to see eachother, so make sure if you're using remote engines, they're all on the same remote network.
* UCP installs its own Swarm, so don't set up Swarm first, just install UCP directly on your engine(s)
* Ports - If you have firewalls configured, please make sure each node can access the following ports on the other nodes in the cluster.
    * UCP Server: 443
    * Swarm Manager: 2376 (user configurable) - Protected via mutual TLS
    * Engine Proxy: 12376 - Protected via mutual TLS
    * KV Store: 12379, 12380 - Protected by mutual TLS
    * CAs: 12381, 12382 - Protected by mutual TLS


# Initial Installation
(line wrapped for readability)
```bash
docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name ucp \
    docker/ucp \
    install --swarm-port 3376 -i
```

The above command will prompt you for some basic information to get UCP
stood up.  You can use "install --help" as the last line above to get
information about various options to the installer.

Important notes for first time users:
* We try to get the hostname/IPs right, but NAT can lead us astray.  Make sure you specify the **real** external hostname when prompted (or use --san for non-interactive mode)
* The last line of the installer output tells you where to go log in.
* If you didn't enter an admin password, the default login is "admin/orca"
* The first thing you probably want to do is download a cert bundle (upper right corner of UI, but subject to change)
    * With this, you can run docker CLI commands against UCP (and/or the swarm/engines if you're account is an admin account)
    * Take a look at the env.sh within the zip file for instructions (should be familiar if you've used machine)


## Data Persistence

UCP uses named volumes for persistence of user data.  By default,
the bootstrapper will create these using the default volume driver and
flags if they are not detected.  If you use a custom volume driver, you
can pre-create volumes prior to installing UCP.

* **orca-root-ca** - The certificate and key for the UCP Root CA
* **orca-swarm-root-ca** - The certificate and key for the Swarm Root CA
* **orca-server-certs** - The server certificates for the UCP web server
* **orca-swarm-node-certs** - The swarm certificates for the current node (repeated on every node in the cluster)
* **orca-swarm-kv-certs**   The Swarm KV client certificates for the current node (repeated on every node in the cluster)
* **orca-swarm-controller-certs**  The UCP Controller Swarm client certificates for the current node
* **orca-kv** - KV store persistence


## User Supplied Certificates

UCP uses two separate root CAs for access control - one for Swarm,
and one for the UCP server itself.  The motivation for the dual root
certificates is to differentiate Docker remote API access to UCP
vs. Swarm.  UCP implements ACL and audit logging on a per-users basis
which are not offered in Swarm or the engines.  Swarm and the engine
proxies trust only the Swarm Root CA, while the UCP server trusts both
Root CAs.  Admins can access UCP, Swarm and the engines while normal
users are only granted access to UCP.

In UCP v1.0 we support user provided externally signed certificates
for the UCP server.  This cert is used by the main management web UI
(visible to your browser when you connect) as well as the Docker remote
API (visible to the Docker CLI and friends.)  The Swarm Root CA is
always manged by UCP itself in this release.  This external UCP Root
CA model supports customers managing their own CA, or purchasing certs
from a commercial CA.  When operating in this mode, UCP can not generate
regular user certificates, as those must be managed and signed externally,
however admin account certs can be generated as they are signed by the
internal Swarm Root CA.  Normal user accounts should be signed by the
same external Root CA (or a trusted intermediary), and the public keys
manually added through the UI.

To install UCP with an external Root CA, create a named volume called **orca-server-certs**
on the engine host where you will install UCP **before** running the install, and ensure the following
files are present in the top-level directory of this volume:

* **ca.pem** - Your Root CA Certificate chain (including any intermediaries)
* **cert.pem** - Your signed UCP server cert
* **key.pem** - Your UCP server private key

After setting up these files on the host, you can install with the "--external-orca-ca" flag.

If you are creating your own storage volumes (for example, to take
advantage of a 3rd party storage driver) you can omit the **orca-root-ca**
volume as it will not be used when using an external UCP Root CA.

```bash
docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name ucp \
    docker/ucp \
    install --swarm-port 3376 -i --external-orca-ca
```


# Adding Nodes to the Cluster
To add capacity to your cluster, run the following on the engine you want to add (**not** the engine where you installed UCP above)
```bash
docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name ucp \
    docker/ucp \
    join -i
```

As with install, you can use "join --help" for more information.


# Cross-host Networking

See [networking.md](networking.md) for more details on the steps required
to enable cross-host networking

# Uninstalling
The installer can also uninstall the UCP software on either the primary
UCP node, as well as the secondary UCP nodes.  Run the following directly
against the engine you want to uninstall (**not** against UCP itself)

```bash
docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name ucp \
    docker/ucp \
    uninstall
```

As above, use "uninstall --help" to see what other flags are available to tune behavior.

# Upgrading

**Coming soon!** (before GA)

For now, just uninstall and re-install.  Your containers will remain on the engines and survive across the uninstall/install.
