---
title: Image Access Management
description: Control which Docker Hub images developers can access with Image Access Management for enhanced supply chain security
keywords: image access management, docker official images, verified publisher, supply chain security, docker business
tags: [admin]
aliases:
 - /docker-hub/image-access-management/
 - /desktop/hardened-desktop/image-access-management/
 - /admin/organization/image-access/
 - /security/for-admins/image-access-management/
 - /security/for-admins/hardened-desktop/image-access-management/
weight: 40
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Image Access Management lets administrators control which types of images developers can pull from Docker Hub. This prevents developers from accidentally using untrusted community images that could pose security risks to your organization.

With Image Access Management, you can restrict access to:

- Docker Official Images: Curated images maintained by Docker
- Docker Verified Publisher Images: Images from trusted commercial publishers
- Organization images: Your organization's private repositories
- Community images: Public images from individual developers

## Who should use Image Access Management?

Image Access Management helps prevent supply chain attacks by ensuring developers only use trusted container images. For example, a developer building a new application might accidentally use a malicious community image as a component. Image Access Management prevents this by restricting access to only approved image types.

Common security scenarios include:

- Prevent use of unmaintained or malicious community images
- Ensure developers use only vetted, official base images
- Control access to commercial third-party images
- Maintain consistent security standards across development teams

## Prerequisites

Before configuring Image Access Management, you must:

- [Enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md) to ensure users authenticate with your organization
- Use [personal access tokens (PATs)](/manuals/security/access-tokens.md) for authentication (Organization access tokens aren't supported)
- Have a Docker Business subscription

> [!IMPORTANT]
>
> Image Access Management only takes effect when users are signed in to Docker Desktop with organization credentials.

## Configure image access

To configure Image Access Management:

1. Sign in to [Docker Home](https://app.docker.com) and select your organization from the top-left account drop-down.
1. Select **Admin Console**, then **Image access**.
1. Use the **toggle** to enable image access.
1. Select which image types to allow:
    - **Organization images**: Images from your organization (always allowed by default). These can be public or private images created by members within your organization.
    - **Community images**: Images contributed by various users that may pose security risks. This category includes Docker-Sponsored Open Source images and is turned off by default.
    - **Docker Verified Publisher Images**: Images from Docker partners in the Verified Publisher program, qualified for secure supply chains.
    - **Docker Official Images**: Curated Docker repositories that provide OS repositories, best practices for Dockerfiles, drop-in solutions, and timely security updates.

Once restrictions are applied, organization members can view the permissions page in read-only format.

> [!NOTE]
>
> Image Access Management is turned off by default. Organization owners have access to all images regardless of policy settings.

## Verify access restrictions

After configuring Image Access Management, test that restrictions work correctly.

When developers pull allowed image types:

```console
$ docker pull nginx  # Docker Official Image
# Pull succeeds if Docker Official Images are allowed
```

When developers pull blocked image types:

```console
$ docker pull someuser/custom-image  # Community image
Error response from daemon: image access denied: community images not allowed
```

Image access restrictions apply to all Docker Hub operations including pulls, builds using `FROM` instructions, and Docker Compose services.

## Security implementation

Start with the most restrictive policy and gradually expand based on legitimate business needs:

1. Start with: Docker Official Images and Organization images
2. Add if needed: Docker Verified Publisher Images for commercial tools
3. Carefully evaluate: Community images only for specific, vetted use cases

Other security recommendations include:

- Monitor usage patterns: Review which images developers are attempting to pull, identify legitimate requests for additional image types, regularly audit approved image categories for continued relevance, and use Docker Desktop analytics to monitor usage patterns.
- Layer security controls: Image Access Management works best with Registry Access Management to control which registries developers can access, Enhanced Container Isolation to secure containers at runtime, and Settings Management to control Docker Desktop configuration.

## Scope and bypass considerations

- Image Access Management only controls access to Docker Hub images. Images from other registries aren't affected by these policies. Use [Registry Access Management](/manuals/enterprise/security/hardened-desktop/registry-access-management.md) to control access to other registries.
- Users can potentially bypass Image Access Management by signing out of Docker Desktop (unless sign-in is enforced), using images from other registries that aren't restricted, or using registry mirrors or proxies. Enforce sign-in and combine with Registry Access Management for comprehensive control.
- Image restrictions apply to Dockerfile `FROM` instructions, Docker Compose services using restricted images will fail, multi-stage builds may be affected if intermediate images are restricted, and CI/CD pipelines using diverse image types may be impacted.
