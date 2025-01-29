---
title: Using Docker with Zscaler
tags: [networking, admin]
summary: |
  This guide explains how to embed Zscaler’s root certificate into Docker
  images, allowing containers to operate securely with Zscaler proxies and
  avoid SSL errors.
params:
  time: 10 minutes
---

In many corporate environments, network traffic is intercepted and monitored
using HTTPS proxies, such as Zscaler. While Zscaler ensures security compliance
and network control, it can cause issues for developers using Docker,
particularly during build processes, where SSL certificate validation errors
might occur. This guide outlines how to configure Docker containers and builds
to properly handle Zscaler's custom certificates, ensuring smooth operation in
monitored environments.

## The role of certificates in Docker

When Docker builds or runs containers, it often needs to fetch resources from
the internet—whether it's pulling a base image from a registry, downloading
dependencies, or communicating with external services. In a proxied
environment, Zscaler intercepts HTTPS traffic and replaces the remote server's
certificate with its own. However, Docker doesn't trust this Zscaler
certificate by default, leading to SSL errors.

```plaintext
x509: certificate signed by unknown authority
```

These errors occur because Docker cannot verify the validity of the certificate
presented by Zscaler. To avoid this, you must configure Docker to trust
Zscaler's certificate.

## Configure Zscaler proxy for Docker Desktop

Depending on how Zscaler is deployed, you may need to configure Docker Desktop
proxy settings manually to use the Zscaler proxy.

If you're using Zscaler as a system-level proxy via the [Zscaler Client Connector](https://help.zscaler.com/zscaler-client-connector/what-is-zscaler-client-connector),
all traffic on the device is automatically routed through Zscaler, so Docker
Desktop uses the Zscaler proxy automatically with no additional configuration
necessary.

If you are not using Zscaler as a system-level proxy, manually configure proxy
settings in Docker Desktop. Set up proxy settings for all clients in the
organization using [Settings Management](/manuals/security/for-admins/hardened-desktop/settings-management/_index.md),
or edit proxy configuration in the Docker Desktop GUI under [**Settings > Resources > Proxies**](/manuals/desktop/settings-and-maintenance/settings.md#proxies).

## Install root certificates in Docker images

To enable containers to use and trust the Zscaler proxy, embed the certificate
in the image and configure the image's trust store. Installing certificates at
image build time is the preferred approach, as it removes the need for
configuration during startup and provides an auditable, consistent environment.

### Obtaining the root certificate

The easiest way to obtain the root certificate is to export it from a machine
where an administrator has already installed it. You can use either a web
browser or the system's certificate management service (for example, Windows
Certificate Store).

#### Example: Exporting the certificate using Google Chrome

1. In Google Chrome, navigate to `chrome://certificate-manager/`.
2. Under **Local certificates**, select **View imported certificates**.
3. Find the Zscaler root certificate, often labeled **Zscaler Root CA**.
4. Open the certificate details and select **Export**.
5. Save the certificate in ASCII PEM format.
6. Open the exported file in a text editor to confirm it includes `-----BEGIN CERTIFICATE-----` and `-----END CERTIFICATE-----`.

When you have obtained the certificate, store it in an accessible repository,
such as JFrog Artifactory or a Git repository. Alternatively, use generic
storage like AWS S3.

### Building with the certificate

To install these certificates when building images, copy the certificate into
the build container and update the trust store. An example Dockerfile looks
like this:

```dockerfile
FROM debian:bookworm
COPY zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates
```

Here, `zscaler-root-ca.crt` is the root certificate, located at the root of the
build context (often within the application's Git repository).

If you use an artifact repository, you can fetch the certificate directly using
the `ADD` instruction. You can also use the `--checksum` flag to verify that
the content digest of the certificate is correct.

```dockerfile
FROM debian:bookworm
ADD --checksum=sha256:24454f830cdb571e2c4ad15481119c43b3cafd48dd869a9b2945d1036d1dc68d \
    https://artifacts.example/certs/zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates
```

#### Using multi-stage builds

For multi-stage builds where certificates are needed in the final runtime
image, ensure the certificate installation occurs in the final stage.

```dockerfile
FROM debian:bookworm AS build
WORKDIR /build
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    curl \
    git
RUN --mount=target=. cmake -B output/

FROM debian:bookworm-slim AS final
ADD --checksum=sha256:24454f830cdb571e2c4ad15481119c43b3cafd48dd869a9b2945d1036d1dc68d \
    https://artifacts.example/certs/zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates
WORKDIR /app
COPY --from=build /build/output/bin .
ENTRYPOINT ["/app/bin"]
```

## Conclusion

Embedding the Zscaler root certificate directly into your Docker images ensures
that containers run smoothly within Zscaler-proxied environments. By using this
approach, you reduce potential runtime errors and create a consistent,
auditable configuration that allows for smooth Docker operations within a
monitored network.
