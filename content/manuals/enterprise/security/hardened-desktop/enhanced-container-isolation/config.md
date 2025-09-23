---
title: Configure Docker socket exceptions and advanced settings
linkTitle: Configure advanced settings
description: Configure Docker socket exceptions and advanced settings for Enhanced Container Isolation
keywords: enhanced container isolation, docker socket, configuration, testcontainers, admin settings
aliases:
 - /desktop/hardened-desktop/enhanced-container-isolation/config/
 - /security/for-admins/hardened-desktop/enhanced-container-isolation/config/
weight: 20
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

This page shows you how to configure Docker socket exceptions and other advanced settings for Enhanced Container Isolation (ECI). These configurations enable trusted tools like Testcontainers to work with ECI while maintaining security.

## Docker socket mount permissions

By default, Enhanced Container Isolation blocks containers from mounting the Docker socket to prevent malicious access to Docker Engine. However, some tools require Docker socket access.

Common scenarios requiring Docker socket access include:

- Testing frameworks: Testcontainers, which manages test containers
- Build tools: Paketo buildpacks that create ephemeral build containers
- CI/CD tools: Tools that manage containers as part of deployment pipelines
- Development utilities: Docker CLI containers for container management

## Configure socket exceptions

Configure Docker socket exceptions using Settings Management:

{{< tabs >}}
{{< tab name="Admin Console" >}}

1. Sign in to [Docker Home](https://app.docker.com) and select your organization from the top-left account drop-down.
1. Go to **Admin Console** > **Desktop Settings Management**.
1. [Create or edit a setting policy](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md).
1. Find **Enhanced Container Isolation** settings.
1. Configure **Docker socket access control** with your trusted images and
command restrictions.

{{< /tab >}}
{{< tab name="JSON file" >}}

Create an [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md) and add:

```json
{
  "configurationFileVersion": 2,
  "enhancedContainerIsolation": {
    "locked": true,
    "value": true,
    "dockerSocketMount": {
      "imageList": {
        "images": [
          "docker.io/localstack/localstack:*",
          "docker.io/testcontainers/ryuk:*",
          "docker:cli"
        ],
        "allowDerivedImages": true
      },
      "commandList": {
        "type": "deny",
        "commands": ["push", "build"]
      }
    }
  }
}
```

{{< /tab >}}
{{< /tabs >}}

## Image allowlist configuration

The `imageList` defines which container images can mount the Docker socket.

### Image reference formats

| Format  | Description |
| :---------------------- | :---------- |
| `<image_name>[:<tag>]`  | Name of the image, with optional tag. If the tag is omitted, the `:latest` tag is used. If the tag is the wildcard `*`, then it means "any tag for that image." |
| `<image_name>@<digest>` | Name of the image, with a specific repository digest (e.g., as reported by `docker buildx imagetools inspect <image>`). This means only the image that matches that name and digest is allowed. |

### Example configurations

Basic allowlist for testing tools:

```json
"imageList": {
  "images": [
    "docker.io/testcontainers/ryuk:*",
    "docker:cli",
    "alpine:latest"
  ]
}
```

Wildcard allowlist (Docker Desktop 4.36 and later):

```json
"imageList": {
  "images": ["*"]
}
```

> [!WARNING]
>
> Using `"*"` allows all containers to mount the Docker socket, which reduces security. Only use this when explicitly listing allowed images isn't feasible.

### Security validation

Docker Desktop validates allowed images by:

1. Downloading image digests from registries for allowed images
1. Comparing container image digests against the allowlist when containers start
1. Blocking containers whose digests don't match allowed images

This prevents bypassing restrictions by re-tagging unauthorized images:

```console
$ docker tag malicious-image docker:cli
$ docker run -v /var/run/docker.sock:/var/run/docker.sock docker:cli
# This fails because the digest doesn't match the real docker:cli image
```

## Derived images support

For tools like Paketo buildpacks that create ephemeral local images, you can
allow images derived from trusted base images.

### Enable derived images

```json
"imageList": {
  "images": [
    "paketobuildpacks/builder:base"
  ],
  "allowDerivedImages": true
}
```

When `allowDerivedImages` is true, local images built from allowed base images (using `FROM` in Dockerfile) also gain Docker socket access.

### Derived images requirements

- Local images only: Derived images must not exist in remote registries
- Base image available: The parent image must be pulled locally first
- Performance impact: Adds up to 1 second to container startup for validation
- Version compatibility: Full wildcard support requires Docker Desktop 4.36+

## Command restrictions

### Deny list (recommended)

Blocks specified commands while allowing all others:

```json
"commandList": {
  "type": "deny",
  "commands": ["push", "build", "image*"]
}
```

### Allow list

Only allows specified commands while blocking all others:

```json
"commandList": {
  "type": "allow",
  "commands": ["ps", "container*", "volume*"]
}
```

### Command wildcards

| Wildcard | Blocks/allows |
| :---------------- | :---------- |
| `"container\*"`     | All "docker container ..." commands |
| `"image\*"`         | All "docker image ..." commands |
| `"volume\*"`        | All "docker volume ..." commands |
| `"network\*"`       | All "docker network ..." commands |
| `"build\*"`         | All "docker build ..." commands |
| `"system\*"`        | All "docker system ..." commands |

### Command blocking example

When a blocked command is executed:

```console
/ # docker push myimage
Error response from daemon: enhanced container isolation: docker command "/v1.43/images/myimage/push?tag=latest" is blocked; if you wish to allow it, configure the docker socket command list in the Docker Desktop settings.
```

## Common configuration examples

### Testcontainers setup

For Java/Python testing with Testcontainers:

```json
"dockerSocketMount": {
  "imageList": {
    "images": [
      "docker.io/testcontainers/ryuk:*",
      "testcontainers/*:*"
    ]
  },
  "commandList": {
    "type": "deny",
    "commands": ["push", "build"]
  }
}
```

### CI/CD pipeline tools

For controlled CI/CD container management:

```json
"dockerSocketMount": {
  "imageList": {
    "images": [
      "docker:cli",
      "your-registry.com/ci-tools/*:*"
    ]
  },
  "commandList": {
    "type": "allow",
    "commands": ["ps", "container*", "image*"]
  }
}
```

### Development environments

For local development with Docker-in-Docker:

```json
"dockerSocketMount": {
  "imageList": {
    "images": [
      "docker:dind",
      "docker:cli"
    ]
  },
  "commandList": {
    "type": "deny",
    "commands": ["system*"]
  }
}
```

## Security recommendations

### Image allowlist best practices

- Be restrictive: Only allow images you absolutely trust and need
- Use wildcards carefully: Tag wildcards (`*`) are convenient but less secure than specific tags
- Regular reviews: Periodically review and update your allowlist
- Digest pinning: Use digest references for maximum security in critical environments

### Command restrictions

- Default to deny: Start with a deny list blocking dangerous commands like `push` and `build`
- Principle of least privilege: Only allow commands your tools actually need
- Monitor usage: Track which commands are being blocked to refine your configuration

### Monitoring and maintenance

- Regular validation: Test your configuration after Docker Desktop updates, as image digests may change.
- Handle digest mismatches: If allowed images are unexpectedly blocked:
    ```console
    $ docker image rm <image>
    $ docker pull <image>
    ```

This resolves digest mismatches when upstream images are updated.

## Next steps

- Review [Enhanced Container Isolation limitations](/manuals/enterprise/security/hardened-desktop/enhanced-container-isolation/limitations.md).
- Review [Enhanced Container Isolation FAQs](/manuals/enterprise/security/hardened-desktop/enhanced-container-isolation/faq.md).
