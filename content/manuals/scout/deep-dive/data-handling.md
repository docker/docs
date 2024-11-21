---
description: How Docker Scout handles image metadata
keywords: |
  scout, scanning, supply chain, security, purl, sbom, provenance, environment,
  materials, config, ports, labels, os, registry, timestamp, digest, layers,
  architecture, license, dependencies, base image
title: Data collection and storage in Docker Scout
aliases:
  /scout/data-handling/
---

Docker Scout's image analysis works by collecting metadata from the container
images that you analyze. This metadata is stored on the Docker Scout platform.

## Data transmission

This section describes the data that Docker Scout collects and sends to the
platform.

### Image metadata

Docker Scout collects the following image metadata:

- Image creation timestamp
- Image digest
- Ports exposed by the image
- Environment variable names and values
- Name and value of image labels
- Order of image layers
- Hardware architecture
- Operating system type and version
- Registry URL and type

Image digests are created for each layer of an image when the image is built
and pushed to a registry. They are SHA256 digests of the contents of a layer.
Docker Scout doesn't create the digests; they're read from the image manifest.

The digests are matched against your own private images and Docker's database
of public images to identify images that share the same layers. The image that
shares most of the layers is considered a base image match for the image that's
currently being analyzed.

### SBOM metadata

Software Bill of Material (SBOM) metadata is used to match package types
and versions with vulnerability data to infer whether an image is affected.
When the Docker Scout platform receives information from security advisories
about new CVEs or other risk factors, such as leaked secrets, it cross-references
this information with the SBOM. If there's a match, Docker Scout displays the
results in the user interfaces where Docker Scout data is surfaced,
such as the Docker Scout Dashboard and in Docker Desktop.

Docker Scout collects the following SBOM metadata:

- Package URLs (PURL)
- Package author and description
- License IDs
- Package name and namespace
- Package scheme and size
- Package type and version
- Filepath within the image
- The type of direct dependency
- Total package count

The PURLs in Docker Scout follow the
[purl-spec](https://github.com/package-url/purl-spec) specification. Package
information is derived from the contents of image, including OS-level programs
and packages, and application-level packages such as maven, npm, and so on.

### Environment metadata

If you integrate Docker Scout with your runtime environment via the
[Sysdig integration](/manuals/scout/integrations/environment/sysdig.md),
Docker Scout collects the following data points about your deployments:

- Kubernetes namespace
- Workload name
- Workload type (for example, DaemonSet)

### Local analysis

For images analyzed locally on a developer's machine, Docker Scout only
transmits PURLs and layer digests. This data isn't persistently stored on the
Docker Scout platform; it's only used to run the analysis.

### Provenance

For images with [provenance attestations](/manuals/build/metadata/attestations/slsa-provenance.md),
Docker Scout stores the following data in addition to the SBOM:

- Materials
- Base image
- VCS information
- Dockerfile

## Data storage

For the purposes of providing the Docker Scout service, data is stored using:

- Amazon Web Services (AWS) on servers located in US East
- Google Cloud Platform (GCP) on servers located in US East

Data is used according to the processes described at
[docker.com/legal](https://www.docker.com/legal/) to provide the key
capabilities of Docker Scout.
