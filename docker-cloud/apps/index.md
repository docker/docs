---
description: Manage your Docker Cloud Applications
keywords: applications, reference, Cloud
title: Applications in Docker Cloud
notoc: true
---

Applications in Docker Cloud are usually several Services linked together using
the specifications from a [Stackfile](stacks.md) or a Compose file. You can also
create individual services using the Docker Cloud Services wizard, and you can
attach [Volumes](volumes.md) to use as long-lived storage for your services.

If you are using Docker Cloud's autobuild and autotest features, you can also
use [autoredeploy](auto-redeploy.md) to automatically redeploy the application
each time its underlying services are updated.

* [Deployment tags](deploy-tags.md)
* [Add a Deploy to Docker Cloud button](deploy-to-cloud-btn.md)
* [Manage service stacks](stacks.md)
    * [Stack YAML reference](stack-yaml-reference.md)
* [Publish and expose service or container ports](ports.md)
* [Redeploy running services](service-redeploy.md)
* [Scale your service](service-scaling.md)
* [Service API Roles](api-roles.md)
* [Service discovery and links](service-links.md)
* [Work with data volumes](volumes.md)
* [Create a proxy or load balancer](load-balance-hello-world.md)

### Automate your applications

Use the following features to automate specific actions on your Docker Cloud applications.

* [Automatic container destroy](auto-destroy.md)
* [Automatic container restart](autorestart.md)
* [Autoredeploy](auto-redeploy.md)
* [Use triggers](triggers.md)
