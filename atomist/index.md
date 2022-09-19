---
title: Introduction to Atomist
description: Introduction to Atomist
keywords:
  atomist, software supply chain, indexing, sbom, vulnerabilities, automation
---

{% include atomist/disclaimer.md %}

Atomist is a data and automation platform for managing the software supply
chain. It will extract metadata from container images, evaluate the data, and
help you understand the state of the image.

Integrating Atomist into your systems and repositories grants you essential
information about the images you build, and the containers running in
production. Beyond collecting and visualizing information, Atomist can help you
further by giving you recommendations, notifications, validation, and more.

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
image, it will analyze and extract metadata about the image contents and any
base images used. The metadata is uploaded to an isolated partition in the
Atomist data plane where it's securely stored.

The Atomist data plane is a combination of metadata and a large knowledge graph
of public software and vulnerability data. Atomist determines the state of your
container by overlaying the image metadata with the knowledge graph.
