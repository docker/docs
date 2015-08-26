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


## User Supplied Certificates

Orca uses two separate root CAs for access control - one for Swarm,
and one for the Orca server itself.  The motivation for the dual root
certificates is to differentiate Docker remote API access to Orca
vs. Swarm.  Orca implements ACL and audit logging on a per-users basis
which are not offered in Swarm or the engines.  Swarm and the engine
proxies trust only the Swarm Root CA, while the Orca server trusts both
Root CAs.  Admins can access Orca, Swarm and the engines while normal
users are only granted access to Orca.

In Orca v1.0 we support user provided externally signed certificates
for the Orca server.  This cert is used by the main management web UI
(visible to your browser when you connect) as well as the Docker remote
API (visible to the Docker CLI and friends.)  The Swarm Root CA is
always manged by Orca itself in this release.  This external Orca Root
CA model supports customers managing their own CA, or purchasing certs
from a commercial CA.  When operating in this mode, Orca can not generate
regular user certificates, as those must be managed and signed externally,
however admin account certs can be generated as they are signed by the
internal Swarm Root CA.  Normal user accounts should be signed by the
same external Root CA (or a trusted intermediary), and the public keys
manually added through the UI.

To install Orca with an external Root CA, place the following files on the
engine host where you will install Orca **before** running the install:

* /var/lib/docker/orca\_ssl/orca\_ca.pem - Your Root CA Certificate chain (including any intermediaries)
* /var/lib/docker/orca\_ssl/orca\_controller.pem - Your signed Orca server cert
* /var/lib/docker/orca\_ssl/orca\_controller\_key.pem - Your Orca server private key

After setting up these files on the host, you can install with the "--external-orca-ca" flag.

```bash
docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name orca-bootstrap \
    dockerorca/orca-bootstrap \
    install --swarm-port 3376 -i --external-orca-ca
```


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
