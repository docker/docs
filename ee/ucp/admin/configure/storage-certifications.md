---
title: Certified Swarm Storage Drivers
description: Deploy certified storage drivers for Swarm.
keywords: Universal Control Plane, UCP, storage, certification
---

## Driver Certification

Certification for 3rd party components is useful for Enterprises for a variety of reasons. Docker is actively engaged with the storage ecosystem to provide tested and validated storage solutions that are certified to work with Docker Enterprise. This removes the much of the burden from companies in ensuring that complex products from different companies naturally work together. Certification ensures the following things:

- The storage driver has been tested to install and operate with a specific set of versions of Docker Enterprise
- Docker's commercial technical support team has communication and case escalation paths with the certified storage partner to ensure that technical issues are worked collaboratively and that there is no finger pointing.
- Storage drivers are continuously tested and validated against new major versions of the driver and of Docker Enterprise


## Swarm Storage Driver Certifications

The following storage drivers are validated against Docker Swarm when running as a component of Docker Enterprise.

| Partner     | Storage Plugin     | Docker Enterprise Version| Notes     |
| ----------- | ----------- | ----------- | ----------- | 
| Pure Storage      | Pure Storage Docker Volume Plugin  | Docker Enterprise 2.0 |     |
| NetApp            | [Trident](https://store.docker.com/plugins/netapp-docker-volume-plugin-ndvp)                            | Docker Enterprise 2.0 | Only NFS protocol supported currently. No iSCSI support.|
| Veritas            | [Veritas Docker Volume Plugin V2](https://store.docker.com/plugins/veritas-docker-volume-plugin-v2) | Docker Enterprise 2.0 | |
| HPE            | 3PAR Volume Plugin for Docker        | Docker Enterprise 2.0 | |
| HPE            | [Nimble Storage Docker Volume Plugin](https://store.docker.com/plugins/nimble)   | Docker Enterprise 2.0 | |
| Nutanix        | [Nutanix Docker Volume Plugin](https://store.docker.com/plugins/nutanix-dvp-docker-volume-plug-in)   | Docker Enterprise 2.0 | |
| Nexenta        | [NexentaStor Docker NFS Volume Plug-In](https://store.docker.com/plugins/nexentastor-docker-nfs-volume-plug-in)   | Docker Enterprise 2.0 | |




