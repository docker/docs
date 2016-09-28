# Highly Available Certificate Authorities

In v1.0 of UCP, the dual CAs (one for user certs, and one for the cluster)
are not highly available - they exist on only the initial node stood
up with the "install" command.  This document outlines the model for
improving that in 1.1.

Starting in 1.1, we run a pair of CAs on each controller.  These CAs
will use replicated root key-pairs so they are equivalent.

## Overview

* Run the CFSSL CA User and Cluster services on each controller
* Use the same named volume for each controller so they're ~interchangeable
* The user may provide root key material during replica join to pre-initialize the root CA on that node
* If no root key material is provided, bogus certs are set up for the root, and it will be ignored until populated with root material


## Upgrade

* New controllers will block old bootstrappers from joining by detecting the old join request payload
* **Initial Controller**
    * No substantive change, works as before
    * Configuration updated to new structure to expose a list of available CAs, pre-populated with just this initial root CA
* **Existing Replicas**
    * During upgrade, the bootstrapper will detect a pre 1.1 replica
    * If the user has pre-populated the root keys, the new CAs will come up and be usable
    * If no root keys are detected, using certs on the existing replica volumes, the bootstrapper will "side band" and go directly to the existing CAs, and make CSR requests for placeholder CA certs
    * The user may manually copy the key material to the node at any time
* **Existing non-Replicas**
    * Nothing changes, they'll still be wired up based on certs signed by the root, and trust

## Cert Replacement

There are various reasons a user might want to regenerate their certs.

1. They got SANs wrong and are having trouble connecting to the controller (e.g., through a LB they added after deploying)
2. Their certs expired (we set our expiration long, so this wont be a problem for many years)
3. A cluster node was compromised
    * In general, if a node in the cluster was compromised, you should probably regenerate the Root CAs and all certs for the cluster
4. A users bundle was compromised


* User Bundles
    * Easy - just download a new one, and remove the old one from the users list
* Non-replica Nodes:
    * Run `join --fresh-install` to re-join the node (no workloads will be impacted, and the certs will be regenerated)
* Controller Nodes:
    * Run `regen-certs --root-ca-only` first
    * Replicate the new Root CA material to all controllers with backup/restore
    * Run `regen-certs` without the CA flag to regenerate the local node (this will cause outages)
    * Once quorum is achieved etcd will recover, and the cluster will come back online
    * Restart all the controller daemons
    * Run through all non-replica nodes and re-join them

## Revocation

* If an infrastructure cert was compromised, for now, regenerate the Root CAs and all downstream certs.
* If only a user cert was compromised, have the user remove it from their list (we should add admin support for this)


* TODO: wire up OCSP (may not make 1.1 given schedule pressure, but we'll see)
* Store revoked certs info in the KV store
* URL on every controller routes back to the KV store so all of them can expose the revocation data
* More details TBD

