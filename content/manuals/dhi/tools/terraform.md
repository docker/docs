---
title: Use the DHI Terraform provider
linktitle: Terraform
description: Use the DHI Terraform provider to manage mirrors and customizations as infrastructure as code.
weight: 40
keywords: dhi terraform, docker hardened images terraform, infrastructure as code, dhi mirror terraform, dhi provider
---

The [DHI Terraform provider](https://registry.terraform.io/providers/docker-hardened-images/dhi/latest/docs)
lets you manage Docker Hardened Image mirrors and customizations as
infrastructure as code.

## Install and configure the provider

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    dhi = {
      source = "docker-hardened-images/dhi"
    }
  }
}

provider "dhi" {
  docker_hub_username = var.docker_username
  docker_hub_password = var.docker_password
  organization        = var.org_name
}
```

Instead of specifying credentials in the provider block, you can set environment
variables:

| Variable | Description |
|----------|-------------|
| `DOCKER_USERNAME` | Docker Hub username or organization namespace |
| `DOCKER_PASSWORD` | Docker Hub password or personal/organization access token |
| `DHI_ORG` | Target organization namespace |

You can authenticate using a personal access token (PAT) or an organization
access token (OAT) in place of a password. When using an OAT, permission scopes
apply:

- Read (pull) access is required to list mirrors.
- Push access is required to create or delete mirrors.

## Resources

### `dhi_mirror`

Manages a mirrored DHI repository in your organization. See [Mirror a Docker
Hardened Image repository](/dhi/how-to/mirror/) for task-based examples.

For the full list of resource attributes, see the [Terraform Registry
documentation](https://registry.terraform.io/providers/docker-hardened-images/dhi/latest/docs/resources/mirror).

### `dhi_customization`

Manages image customizations applied to a mirrored repository. See [Customize a
Docker Hardened Image](/dhi/how-to/customize/) for task-based examples.

For the full list of resource attributes, see the [Terraform Registry
documentation](https://registry.terraform.io/providers/docker-hardened-images/dhi/latest/docs/resources/customization).
