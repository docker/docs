---
description: Learn about the architecture of Docker Universal Control Plane.
keywords: docker, ucp, architecture
title: UCP architecture
---

Universal Control Plane is a containerized application that runs on the
Commercially Supported (CS) Docker Engine. It allows you to manage from a
centralized place a set of nodes that are part of the same swarm.

![](images/architecture-1.png)

## UCP components

The core component of UCP is a globally-scheduled service called `ucp-agent`.
When you install UCP on a node, or join a node to a swarm that is being managed
by UCP, the `ucp-agent` service starts running on that node.

Once this service is running, it deploys containers with other UCP components,
and ensures they keep running. The UCP components that are deployed
on a node depend on whether that node is a manager or a worker.
Manager nodes are responsible for maintaining the swarm state and scheduling
decisions. Worker nodes are responsible for executing workloads.

| Name                | Node            | Description                                                                                               |
|:--------------------|:----------------|:----------------------------------------------------------------------------------------------------------|
| ucp-auth-api        | Manager         | The centralized service for identity and authentication used by UCP and DTR                               |
| ucp-auth-store      | Manager         | Stores authentication configurations, and data for users, organizations and teams                         |
| ucp-auth-worker     | Manager         | Performs scheduled LDAP synchronizations and cleans authentication and authorization data                 |
| ucp-client-root-ca  | Manager         | A certificate authority to sign client bundles                                                            |
| ucp-cluster-root-ca | Manager         | A certificate authority used for TLS communication between UCP components                                 |
| ucp-controller      | Manager         | The UCP web server                                                                                        |
| ucp-kv              | Manager         | Used to store the UCP configurations. Don't use it in your applications, since it's for internal use only |
| ucp-proxy           | Manager, worker | A TLS proxy. It allows secure access to the local Docker Engine to UCP components                         |
| ucp-swarm-manager   | Manager         | Used to provide backwards-compatibility with Docker Swarm                                                 |

## Volumes

Docker UCP uses these named volumes to persist data:

| Volume name                 | Description                                                                              |
|:----------------------------|:-----------------------------------------------------------------------------------------|
| ucp-auth-api-certs          | Certificate and keys for the authentication and authorization service                    |
| ucp-auth-store-certs        | Certificate and keys for the authentication and authorization store                      |
| ucp-auth-store-data         | Data of the authentication and authorization store                                       |
| ucp-auth-worker-certs       | Certificate and keys for authentication worker                                           |
| ucp-auth-worker-data        | Data of the authentication worker                                                        |
| ucp-client-root-ca          | Root key material for the UCP root CA that issues client certificates                    |
| ucp-cluster-root-ca         | Root key material for the UCP root CA that issues certificates for swarm members         |
| ucp-controller-client-certs | Certificate and keys used by the UCP web server to communicate with other UCP components |
| ucp-controller-server-certs | Certificate and keys for the UCP web server running in the node                          |
| ucp-kv                      | UCP configuration data                                                                   |
| ucp-kv-certs                | Certificates and keys for the key-value store                                            |
| ucp-node-certs              | Certificate and keys for node communication                                              |

You can customize the volume driver used for these volumes, by creating
the volumes before installing UCP. During the installation, UCP checks which
volumes don't exist in the node, and creates them using the default volume
driver.

By default, the data for these volumes can be found at
`/var/lib/docker/volumes/<volume-name>/_data`.

## Where to go next

* [System requirements](installation/index.md)
* [Plan a production installation](installation/plan-production-install.md)
