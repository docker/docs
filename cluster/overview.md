---
description: Introduction and Overview of Docker Cluster
keywords: documentation, docs, docker, cluster, infrastructure, automation
title: Overview of Docker Cluster
---

Docker Cluster is a tool for lifecycle management of Docker clusters.
With Cluster, you use a YAML file to configure your provider's resources.
Then, with a single command, you provision and install all the resources
from your configuration.

Using Docker Cluster is a three-step process:

1. Ensure you have the credentials necessary to provision a cluster.

2. Define the resources that make up your cluster in `cluster.yml`

3. Run `docker cluster create` to have Cluster provision resources and install Docker Enterprise on the resources.

A `cluster.yml` file resembles the following example:

      variable:
        region: us-east-2
        ucp_password:
          type: prompt

      provider:
        aws:
          region: ${region}

      cluster:
        engine:
          version: "ee-stable-18.09.5"
        ucp:
          version: "docker/ucp:3.1.6"
          username: "admin"
          password: ${ucp_password}

      resource:
        aws_instance:
          managers:
            quantity: 1

For more information about Cluster files, refer to the
[Cluster file reference](cluster-file/index.md).

Docker Cluster has commands for managing the whole lifecycle of your cluster:

 * Create and destroy clusters
 * Scale up or Scale down clusters
 * Upgrade clusters
 * View the status of clusters
 * Backup and Restore clusters

## Cluster reference pages

- [Get started with Docker Cluster on AWS](aws.md)
- [Command line reference](/engine/reference/commandline/cluster/)
- [Cluster file reference](./cluster-file/index.md)

