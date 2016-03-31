<!--[metadata]>
+++
aliases = ["/ucp/plan-production-install/"]
title = "Plan a production installation"
description = "Learn about the Docker Universal Control Plane architecture, and the requirements to install it on production."
keywords = ["docker, ucp, install, checklist"]
[menu.main]
parent="mn_ucp_installation"
weight=10
+++
<![end-metadata]-->

# Plan a production installation

Docker Universal Control Plane can be installed on-premises, or
on a virtual private cloud. If you've never used Docker UCP before,
you should start by [installing it on a sandbox](../install-sandbox.md).

To secure your data, Docker UCP is automatically set up to use mutual TLS on
all communications. Before you install UCP, make sure you know:

* The fully qualified domain names (FQDN) of the hosts where you'll install UCP,
* Their Subject Alternative Names (SANs).

## Fully-qualified domain names

When installing Docker UCP, the installer tries to find the fully-qualified
domain names (FQDN) of your hosts.

If the installer can't detect this automatically, or if you want to use a
different FQDN or IP address,  use the `--host-address` option when installing.
This option allows you to specify the IP or hostname that UCP is going to use
to reach that host.

If you're installing UCP on a cloud provider such as AWS or Digital Ocean,
you might need to create a private network for you UCP installation. In that
case, make sure all nodes of the cluster can communicate using their private
IPs.


## Subject alternative names (SANs)

When joining new nodes to the cluster, UCP creates leaf certificates for that
node. Those certificates are then used by for communicating over mutual TLS
with other members of the cluster.

You can specify the subject alternative names (SANs) to be used in the
certificate. If you are installing UCP interactively you'll be prompted for
this. You can also use the `--san` option when installing and joining nodes
to the cluster.


## Where to go next

* [UCP system requirements](system-requirements.md)
* [Install UCP for production](install-production.md)
