---
title: Introduction to Atomist
description: Introduction to Atomist
keywords: atomist, software supply chain, vulnerability scanning, tutorial
---

{% include atomist/disclaimer.md %}

Atomist is a data and automation platform for managing the software supply
chain. It will scan, evaluate, and visualize the state of container images and
the software that they contain.

By integrating Atomist into your systems and repositories, you will obtain new,
essential information about your containers. Beyond collecting and presenting
data, Atomist can help you further and provide recommendations, notifications,
validation, and more.

Example capabilities made possible with Atomist are:

- Automatically open pull requests to update base images for improved product
  security.
- Define policy rules that prevent applications from being deployed if they
  contain a secret, such as a password or API token.
- Dissect Dockerfiles and see where vulnerabilities are introduced, line by
  line.

## How it works

Atomist monitors your container registry for new images. When an image is found,
it is scanned, and metadata about the contents of the image is collected. The
metadata is uploaded to the Atomist data plane where it is securely stored.

The Atomist data plane consists of a large knowledge graph of software and
vulnerability data. Atomist determines the state of your container by combining
the container metadata with the knowledge graph.

The software and vulnerability data that is used by Atomist is derived from
public data sources. Metadata about your private container images is never
shared with anyone.
