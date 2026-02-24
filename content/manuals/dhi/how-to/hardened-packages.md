---
title: Use Hardened System Packages
linkTitle: Use hardened packages
weight: 30
keywords: hardened images, DHI, hardened packages, packages, alpine
description: Learn how to use and verify Docker's hardened system packages in your images.
---

Docker Hardened System Packages are built from source by Docker. This ensures
supply chain integrity throughout your entire image stack by eliminating risks
from potentially compromised public packages.

Access to hardened packages varies by subscription:

- **DHI Community**: Includes hardened packages in base images. Can configure the
  public package repository to access the same packages in custom images.
- **DHI Select**: Includes all Community packages, plus access to additional
  compliance-specific packages (such as FIPS variants) and Docker-patched
  packages through the image customization UI.
- **DHI Enterprise**: Includes all Select packages, plus the ability to configure
  the enterprise package repository directly in your own images for full access
  to compliance and security-patched packages.

## Built-in packages

Supported distributions of Docker Hardened Images (DHI) automatically include
hardened system packages. No additional configuration is required. Simply pull
and use the images as normal.

All packages in these images are built by Docker from source, maintaining
the same security standards as the base images themselves.

## Add hardened packages to your images

You can add hardened packages to your own images in the following two ways.

### Add packages through image customization {tier="DHI Select & DHI Enterprise"}

When customizing Docker Hardened Images with DHI Select or DHI Enterprise, you
can add hardened packages for Alpine-based images through the customization
interface. Follow the steps to [create an image
customization](./customize.md#create-an-image-customization) and select hardened
packages during the customization process.

### Configure the package manager

You can configure your package manager to pull from Docker's hardened package
repositories. This lets you install hardened packages in your own images.

#### Public repository

To use Docker's public hardened package repository in your own images, configure
the Alpine package manager in your Dockerfile.

The configuration process involves three steps:

1. Install the [signing key](https://github.com/docker-hardened-images/keyring)
2. Configure the package repository
3. Update and install packages

The following example shows how to configure the Alpine package manager in your
Dockerfile to use Docker's public hardened package repository:

```dockerfile
FROM alpine:3.23

# Install the signing key
RUN cd /etc/apk/keys && \
    wget https://dhi.io/keyring/dhi-apk@docker-0F81AD7700D99184.rsa.pub

# Replace the default repositories with the hardened package repository
RUN echo "https://dhi.io/apk/alpine/v3.23/main" > /etc/apk/repositories

# Update and install packages
RUN apk update && \
    apk add libpng
```

Replace `3.23` with your Alpine version in both the base image tag and repository URL.

To verify the configuration, build and run the image:

```console
$ docker build -t myapp:latest .
$ docker run -it myapp:latest sh
```

Inside the container, check the configured repositories:

```console
/ # cat /etc/apk/repositories
https://dhi.io/apk/alpine/v3.23/main
```

This ensures all packages are installed from Docker's hardened repository.

All packages installed from the Docker Hardened Images repository are built from
source by Docker and include full provenance.

#### Enterprise repository {tier="DHI Enterprise"}

With DHI Enterprise, you have access to an additional package
repository that includes hardened packages for compliance variants such as FIPS,
as well as additional security patches.

The configuration process involves five steps:

1. Install the [signing key](https://github.com/docker-hardened-images/keyring)
2. Configure the base package repository
3. Install the enterprise configuration package
4. Configure package installation with authentication
5. Build the image passing credentials as a secret using the DHI CLI

  > [!NOTE]
  >
  > You must have the Docker Hardened Images CLI installed and configured. For
  > more information, see [Use the DHI CLI](./cli.md).

The following example shows how to configure the Alpine package manager in your
Dockerfile to use Docker's enterprise hardened package repository:

```dockerfile
FROM alpine:3.23

# Install the signing key
RUN cd /etc/apk/keys && \
    wget https://dhi.io/keyring/dhi-apk@docker-0F81AD7700D99184.rsa.pub

# Replace the default repositories with the hardened package repository
RUN echo "https://dhi.io/apk/alpine/v3.23/main" > /etc/apk/repositories

# Update and install the enterprise configuration package to add the security repository
RUN apk update && \
    apk add dhi-enterprise-conf

# Install packages from the security repository with authentication
RUN --mount=type=secret,id=http_auth \
    HTTP_AUTH="$(cat /run/secrets/http_auth)" \
    apk update && \
    apk add openssl-fips
```

Build the image with authentication passed securely as a build secret:

```console
$ dhictl auth apk > http_auth.txt
$ docker build --secret id=http_auth,src=http_auth.txt -t myapp-enterprise:latest .
$ rm http_auth.txt
```

The `--secret` flag securely mounts the authentication credentials during build
without storing them in the image layers or metadata.

## Verify packages

Every hardened package is cryptographically signed and includes metadata that
proves its provenance and build integrity. You can verify the signatures and
view the metadata to ensure your packages come from Docker's trusted build
infrastructure.

### View package metadata

To view information about a hardened package, including its provenance:

```console
$ apk info -L <package-name>
```

This shows the files included in the package and its metadata.

### Verify package signatures

Hardened packages are cryptographically signed by Docker. When you install the
signing keys and configure your package manager as described previously, the
package manager automatically verifies signatures during installation.

If a package fails signature verification, the package manager will refuse to
install it, protecting you from tampered or compromised packages.

### Build provenance and cryptographic verification

Docker hardened packages are built by Docker's trusted infrastructure and include
verifiable metadata and cryptographic signatures.

To view this metadata for an installed package:

```console
$ apk info -a <package-name>
```

Or to view metadata for a package before installing:

```console
$ apk fetch --stdout <package-name> | tar -xzO .PKGINFO
```

The package signing keys ensure that packages haven't been tampered with after
being built. When you install the signing key and configure your package manager,
all packages are automatically verified before installation.

### Package attestations

Each hardened package includes its own attestations, similar to [image
attestations](./verify.md). These attestations provide provenance and build
information for individual packages, allowing you to trace the supply chain down
to the package level.

You can retrieve package attestations by first extracting package information
from the image's SLSA provenance, then using the package digest to access its
attestations.

#### Extract package information from image attestations

To get provenance information for a specific package from an image's SLSA
provenance attestation, you first need to retrieve the image's provenance and
then filter for the specific package you're interested in.

The SLSA provenance attestation includes a `materials` array that lists all
build inputs, including packages. You can use `jq` to filter this array for a
specific package:

```console
$ docker scout attest get dhi.io/golang:1.26-alpine3.23 \
    --predicate-type https://slsa.dev/provenance/v0.2 | \
    jq '.predicate.materials[] | select( .uri == "https://dhi.io/apk/alpine/v3.23/main/aarch64/golang-1.26-1.26.0-r0.apk" )'
```

Replace the package URI in the `select()` filter with the specific package
you're looking for. You can find available packages by first running the command
without the `select()` filter to see all materials.

This returns the package URI and its SHA-256 digest:

```json
{
  "uri": "https://dhi.io/apk/alpine/v3.23/main/aarch64/golang-1.26-1.26.0-r0.apk",
  "digest": {
    "sha256": "4082a2500abc2e7b8435f9398d3514d760044fa52ca3d10cf80015469124a838"
  }
}
```

#### List attestations for a package

Using the package digest from the previous section, you can list all available
attestations for that package:

```console
$ curl -s https://dhi.io/apk/alpine/v3.23/main/sha256:4082a2500abc2e7b8435f9398d3514d760044fa52ca3d10cf80015469124a838/attestations/list | jq .
```

This returns information about the package and its available attestations:

```json
{
  "subject": {
    "name": "pkg:apk/alpine/golang-1.26@1.26.0-r0?os_name=&os_version=",
    "digest": {
      "sha256": "4082a2500abc2e7b8435f9398d3514d760044fa52ca3d10cf80015469124a838"
    }
  },
  "attestations": [
    {
      "predicate_type": "https://slsa.dev/provenance/v1",
      "digest": {
        "sha256": "97c919cf0edb27087739bbabeea4c1ef88d069cd41791476ba64b69280d63a32"
      },
      "url": "https://dhi.io/apk/alpine/v3.23/main/sha256:4082a2500abc2e7b8435f9398d3514d760044fa52ca3d10cf80015469124a838/attestations/sha256:97c919cf0edb27087739bbabeea4c1ef88d069cd41791476ba64b69280d63a32"
    }
  ]
}
```

#### Retrieve package attestations

To retrieve the actual attestation content, use the URL provided in the
attestation list:

```console
$ curl -s https://dhi.io/apk/alpine/v3.23/main/sha256:4082a2500abc2e7b8435f9398d3514d760044fa52ca3d10cf80015469124a838/attestations/sha256:97c919cf0edb27087739bbabeea4c1ef88d069cd41791476ba64b69280d63a32 | jq .
```

This returns the full SLSA provenance attestation for the package, which
includes information about how the package was built, its dependencies, and
other build materials.

You can continue this process recursively to trace the supply chain all the way
down to the compiler and other build tools used to create the package.
