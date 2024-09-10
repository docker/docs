---
title: Use CA certificates with Docker
linkTitle: CA certificates
description: Learn how to install and use CA certificates on the Docker host and in Linux containers
keywords: docker, networking, ca, certs, host, container, proxy
---

> [!CAUTION]
> Best practices should be followed when using Man-in-the-Middle (MITM) CA
> certificates in production containers. If compromised, attackers could
> intercept sensitive data, spoof a trusted service, or perform
> man-in-the-middle attacks. Consult your security team before you proceed.

If your company uses a proxy that inspects HTTPS traffic, you might need to add
the required root certificates to your host machine and your Docker containers
or images. This is because Docker and its containers, when pulling images or
making network requests, need to trust the proxy’s certificates.

On the host, adding the root certificate ensures that any Docker commands (like
`docker pull`) work without issues. For containers, you'll need to add the root
certificate to the container's trust store either during the build process or
at runtime. This ensures that applications running inside the containers can
communicate through the proxy without encountering security warnings or
connection failures.

## Add CA certificate to the host

The following sections describe how to install CA certificates on your macOS or
Windows host. For Linux, refer to the documentation for your distribution.

### macOS

1. Download the CA certificate for your MITM proxy software.
2. Open the **Keychain Access** app.
3. In Keychain Access, select **System**, then switch to the **Certificates** tab.
4. Drag-and-drop the downloaded certificate into the list of certificates. Enter your password if prompted.
5. Find the newly added certificate, double-click it, and expand the **Trust** section.
6. Set **Always Trust** for the certificate. Enter your password if prompted.
7. Start Docker Desktop and verify that `docker pull` works, assuming Docker Desktop is configured to use the MITM proxy.

### Windows

Choose whether you want to install the certificate using the Microsoft
Management Console (MMC) or your web browser.

{{< tabs >}}
{{< tab name="MMC" >}}

1. Download CA certificate for the MITM proxy software.
2. Open the Microsoft Management Console (`mmc.exe`).
3. Add the **Certificates Snap-In** in the MMC.
   1. Select **File** → **Add/Remove Snap-in**, and then select **Certificates** → **Add >**.
   2. Select **Computer Account** and then **Next**.
   3. Select **Local computer** and then select **Finish**.
4. Import the CA certificate:
   1. From the MMC, expand **Certificates (Local Computer)**.
   2. Expand the **Trusted Root Certification Authorities** section.
   3. Right-click **Certificates** and select **All Tasks** and **Import…**.
   4. Follow the prompts to import your CA certificate.
5. Select **Finish** and then **Close**.
6. Start Docker Desktop and verify that `docker pull` succeeds (assuming Docker Desktop is already configured to use the MITM proxy server).

> [!NOTE]
> Depending on the SDK and/or runtime/framework in use, further steps may be
> required beyond adding the CA certificate to the operating system's trust
> store.

{{< /tab >}}
{{< tab name="Web browser" >}}

1. Download the CA certificate for your MITM proxy software.
2. Open your web browser, go to **Settings** and open **Manage certificates**
3. Select the **Trusted Root Certification Authorities** tab.
4. Select **Import**, then browse for the downloaded CA certificate.
5. Select **Open**, then choose **Place all certificates in the following store**.
6. Ensure **Trusted Root Certification Authorities** is selected and select **Next**.
7. Select **Finish** and then **Close**.
8. Start Docker Desktop and verify that `docker pull` succeeds (assuming Docker Desktop is already configured to use the MITM proxy server).

{{< /tab >}}
{{< /tabs >}}

## Add CA certificates to Linux images and containers

If you need to run containerized workloads that rely on internal or custom
certificates, such as in environments with corporate proxies or secure
services, you must ensure that the containers trust these certificates. Without
adding the necessary CA certificates, applications inside your containers may
encounter failed requests or security warnings when attempting to connect to
HTTPS endpoints.

By [adding CA certificates to images](#add-certificates-to-images) at build
time, you ensure that any containers started from the image will trust the
specified certificates. This is particularly important for applications that
require seamless access to internal APIs, databases, or other services during
production.

In cases where rebuilding the image isn't feasible, you can instead [add
certificates to containers](#add-certificates-to-containers) directly. However,
certificates added at runtime won’t persist if the container is destroyed or
recreated, so this method is typically used for temporary fixes or testing
scenarios.

## Add certificates to images

> [!NOTE]
> The following commands are for an Ubuntu base image. If your build uses a
> different Linux distribution, use equivalent commands for package management
> (`apt-get`, `update-ca-certificates`, and so on).

To add ca certificate to a container image when you're building it, add the
following instructions to your Dockerfile.

```dockerfile
# Install the ca-certificate package
RUN apt-get update && apt-get install -y ca-certificates
# Copy the CA certificate from the context to the build container
COPY your_certificate.crt /usr/local/share/ca-certificates/
# Update the CA certificates in the container
RUN update-ca-certificates
```

### Add certificates to containers

> [!NOTE]
> The following commands are for an Ubuntu-based container. If your container
> uses a different Linux distribution, use equivalent commands for package
> management (`apt-get`, `update-ca-certificates`, and so on).

To add a CA certificate to a running Linux container:

1. Download the CA certificate for your MITM proxy software.
2. If the certificate is in a format other than `.crt`, convert it to `.crt` format:

   ```console {title="Example command"}
   $ openssl x509 -in cacert.der -inform DER -out myca.crt
   ```

3. Copy the certificate into the running container:

    ```console
    $ docker cp myca.crt <containerid>:/tmp
    ```

4. Attach to the container:

    ```console
    $ docker exec -it <containerid> sh
    ```

5. Ensure the `ca-certificates` package is installed (required for updating certificates):

    ```console
    # apt-get update && apt-get install -y ca-certificates
    ```

6. Copy the certificate to the correct location for CA certificates:

    ```console
    # cp /tmp/myca.crt /usr/local/share/ca-certificates/root_cert.crt
    ```

7. Update the CA certificates:

    ```console
    # update-ca-certificates
    ```

    ```plaintext {title="Example output"}
    Updating certificates in /etc/ssl/certs...
    rehash: warning: skipping ca-certificates.crt, it does not contain exactly one certificate or CRL
    1 added, 0 removed; done.
    ```

8. Verify that the container can communicate via the MITM proxy:

    ```console
    # curl https://example.com
    ```

    ```plaintext {title="Example output"}
    <!doctype html>
    <html>
    <head>
        <title>Example Domain</title>
    ...
    ```
