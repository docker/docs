---
title: Introduction to Atomist
description: Introduction to Atomist
keywords: atomist, software supply chain, vulnerability scanning, tutorial
---

{% include atomist/disclaimer.md %}

Atomist is a data and automation platform for managing the software supply
chain. It will extract metadata from container images, evaluate the data, and
help you understand the state of the image.

By integrating Atomist into your systems and repositories, you will obtain new,
essential information about your containers. Beyond collecting and presenting
data, Atomist can help you further by giving you recommendations, notifications,
validation, and more.

Example capabilities made possible with Atomist are:

- Stay up to date with advisory databases without having to re-analyze your
  images.
- Automatically open pull requests to update base images for improved product
  security.
- Check that your applications don't contain secrets, such as a password or API
  token, before they get deployed.
- Dissect Dockerfiles and see where vulnerabilities come from, line by line.

## How it works

Atomist monitors your container registry for new images. When it finds a new
image, it will analyze and extract metadata about the image contents as well as
any base images it uses. The metadata is uploaded to the Atomist data plane
where it's securely stored.

The Atomist data plane consists of a large knowledge graph of software and
vulnerability data. Atomist determines the state of your container by combining
the container metadata with the knowledge graph.

Atomist gets it's data about software vulnerabilities from public data
sources—security advisories. Information about your private container images is
never shared with anyone.
