---
title: Docker Hub APIs
description: >
  Information on working with Docker Hub programmatically (HTTP API, CLI)
keywords: hub, api, cli, distribution, dvp
---

<!-- prettier-ignore-start -->

The following table describes the three HTTP APIs that Docker Hub surfaces.

| Name                                        | Base URL                                             | Description                                                                                                        |
| ------------------------------------------- | ---------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| Registry API                                | `https://registry-1.docker.io/v2/`                   | Implements the API protocol defined in the [OCI distribution spec][1]{: target="blank" rel="noopener" class="_"}. |
| [Docker Hub API documentation](./latest.md) | `https://hub.docker.com/v2/`                         | Interface for integrating with Docker Hub platform services.                                                       |
| [DVP API](./dvp.md)                         | `https://hub.docker.com/api/publisher/analytics/v1/` | Verified publisher's API                                                                                           |

Docker also provides a
[Docker Hub CLI](https://github.com/docker/hub-tool#readme){: target="_blank"
rel="noopener" class="_"} tool (experimental) for interacting with Docker Hub
from the terminal. The CLI tool uses the Docker Hub API.

<!-- link to conformance when available: https://conformance.opencontainers.org/#hosted -->

## Authentication

The
[Registry token authentication specification](../../registry/spec/auth/index.md)
describes the implementation for API authentication in Docker Hub APIs.

[1]: https://github.com/opencontainers/distribution-spec/blob/v{{site.oci_distribution_spec_version}}/spec.md#api
<!-- prettier-ignore-end -->
