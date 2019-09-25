---
title: Enable authentication using TLS client certificates
description: Learn how to enable user authentication via client certificates from your own public key infrastructure (PKI).
keywords: PKI, Client Certificates, Passwordless Authentication, Docker Enterprise, UCP, DTR, UCP PKI, DTR PKI
---

## Overview

In many organizations, authenticating to systems with a username and password combination is either restricted or outright prohibited. With Docker Enterprise 3.0, UCP's [CLI client certificate-based authentication](/ee/ucp/user-access/cli/) has been extended to the web user interface (web UI). DTR has also been enhanced to work with UCP's internally generated client bundles for client certificate-based authentication. If you have an external public key infrastructure (PKI) system, you can manage user authentication using a pool of X.509 client certificates in lieu of usernames and passwords.

## Benefits

The following table outlines existing and added capabilities when using client certificates — both internal to UCP and issued by an external certificate authority (CA) — for authentication.

| Operation                       | Benefit                                                                                                                                                                                                                                                          |
| ------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [UCP browser authentication](#ucp--dtr-browser-authentication) | Previously, UCP client bundles enabled communication between a local Docker client and UCP without the need of a username and password. Importing your client certificates into the browser extends this capability to the UCP web UI.                           |
| [DTR browser authentication](#ucp--dtr-browser-authentication) | You can bypass the login page for the DTR web UI when you use TLS client certificates as a DTR authentication method.                                                                                                                                            |
| [Image pulls and pushes to DTR](#image-pulls-and-pushes-to-dtr)        | You can update Docker engine with a client certificate for image pulls and pushes to DTR without the need for `docker login`. |
| [Image signing](#image-signing)                   | You can use client certificates to sign images that you push to DTR. Depending on which you configure to talk to DTR, the certificate files need to be located in certain directories. Alternatively, you can enable system-wide trust of your custom root certificates. |
| [DTR API access](#dtr-api-access)                  | You can use TLS client certificates in lieu of your user credentials to access the DTR API. |
| [Notary CLI operations with DTR](#notary-cli-operations-with-dtr)      | You can set your DTR as the remote trust server location and pass the certificate flags directly to the Notary CLI to access your DTR repositories. |

## Limitations

- The security of client certificates issued by your organization's PKI is outside of UCP’s control. UCP administrators are responsible for instructing their users on how to authenticate via client certificates.
- Username and password authentication cannot be disabled.
- If client certificates have been configured, they will be used for
all `docker push` and `docker pull` operations for _all users_ of the same
machine.
- Docker Enterprise 3.0 does not check certificate revocation lists (CRLs) or Online Certificate Status Protocol (OCSP) for revoked certificates.

## UCP / DTR browser authentication

The following instructions apply to UCP and DTR administrators. For non-admin users, contact your administrator for details on your PKI's client certificate configuration.

To bypass the browser login pages and hide the logout buttons for both UCP and DTR, follow the steps below.

1. Add your organization's root CA certificates [via the UCP web UI](/ee/ucp/admin/configure/use-your-own-tls-certificates/#configure-ucp-to-use-your-own-tls-certificates-and-keys) or [via the CLI installation command](https://success.docker.com/article/how-do-i-provide-an-externally-generated-security-certificate-during-the-ucp-command-line-installation).

    For testing purposes, you can download an [admin client bundle](/ee/ucp/user-access/cli/#download-client-certificates) from UCP and [convert the client certificates to `pkcs12`](#convert-your-client-certificates-to-a-PKCS12-file)

1. Download UCP's `ca.pem` from `https://<ucp-url>/ca` either in the browser or via `curl`. When using `curl`, redirect the response output to a file.
     `curl -sk https://<ucp-url>/ca -o ca.pem`

1. Enable client certificate authentication for DTR. If previously installed, reconfigure DTR with your UCP hostname's root CA certificate. This will be your organization's root certificate(s) appended to UCP's internal root CA certificates.

     ```
     docker run --rm -it docker/dtr:2.7.0 reconfigure --debug --ucp-url \
      <ucp-url> --ucp-username <ucp_admin_user> --ucp-password \ <ucp_admin_password> --enable-client-cert-auth
      --client-cert-auth-ca "$(cat ca.pem)"
     ```

     See [DTR installation](/reference/dtr/2.7/cli/install/) and [DTR reconfiguration](/reference/dtr/2.7/cli/reconfigure/) CLI reference pages for an explanation of the different options.

1. Import the PKCS12 file into [the browser](#pkcs12-file-browser-import) or [Keychain Access](https://www.digicert.com/ssl-support/p12-import-export-mac-mavericks-server.htm#import_certificate) if you're running macOS.

### Client certificate to PKCS12 file conversion

From the command line, switch to the directory of your client bundle and run the following command to convert the client bundle public and private key pair to a `.p12` file.

```bash
  openssl pkcs12 -export -out cert.p12 -inkey key.pem -in cert.pem
```

Create with a simple password, you will be prompted for it when you import the certificate into the browser or Mac's Keychain Access.

### PKCS12 file browser import

Instructions on how to import a certificate into a web browser vary according to your platform, OS, preferred browser and browser version. As a general rule, refer to one of the following how-to articles:
- ***Firefox***:
https://www.sslsupportdesk.com/how-to-import-a-certificate-into-firefox/
- ***Chrome***:
https://www.comodo.com/support/products/authentication_certs/setup/win_chrome.php
- ***Internet Explorer***:
https://www.comodo.com/support/products/authentication_certs/setup/ie7.php

## Image pulls and pushes to DTR

For pulling and pushing images to your DTR (with client certificate authentication method enabled) without performing a `docker login`, do the following:

1. Create a directory for your DTR public address or FQDN (Fully Qualified Domain Name) within your operating system's TLS certificate directory.

1. As a [superuser](https://en.wikipedia.org/wiki/Superuser), copy the private key (`client.pem`) and certificate (`client.cert`) to the machine you are using for pulling and pushing to DTR without doing a `docker login`. Note that the filenames must match.

1. Obtain the CA certificate from your DTR server, `ca.crt` from `https://<dtrurl>/ca`, and copy `ca.crt` to your operating system's TLS certificate directory so that your machine's Docker Engine will trust DTR. For Linux, this is `/etc/docker/certs.d/<dtrurl>/`. On Docker for Mac, this is `/<home_directory>/certs.d/<dtr_fqdn>/`.

    This is a convenient alternative to, for Ubuntu as an example, adding the DTR server certificate to `/etc/ca-certs` and running `update-ca-certificates`.
    ```curl
     curl -k https://<dtr>/ca -o ca.crt
    ```

    On Ubuntu
    ````bash
    cp ca.crt /etc/ca-certs
    ```

1. Restart the Docker daemon for the changes to take effect. See [Configure your host](/ee/dtr/user/access-dtr/#configure-your-host) for different ways to restart the Docker daemon.

### Add your DTR server CA certificate to system level

You have the option to add your DTR server CA certificate to your system's trusted root certificate pool. This is MacOS Keychain or `/etc/ca-certificates/` on Ubuntu. Note that you will have to remove the certificate if your DTR public address changes.

### Reference guides

- [Docker Engine](https://docs.docker.com/engine/security/certificates/)
- Docker Desktop
  - [Enterprise for Mac](/ee/desktop/user/mac-user/#add-tls-certificates)
  - [Enterprise for Windows](/ee/desktop/user/windows-user/#adding-tls-certificates)
  - [Community for Mac](/docker-for-mac/#add-tls-certificates)
  - [Community for Windows](/docker-for-windows/faqs/#certificates)

Note: The above configuration means that Docker Engine will use the same client certificate for all pulls and pushes to DTR for ***all users*** of the same machine.

## Image signing

DTR provides the Notary service for using Docker Content Trust (DCT) out of the box.

<table style="width:100%;">
<colgroup>
<col style="width: 35%" />
<col style="width: 30%" />
<col style="width: 35%" />
</colgroup>
<thead>
<tr class="night">
<th>Implementation</th>
<th>Component Pairing</th>
<th>Settings</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><a href="/engine/security/trust/content_trust/#signing-images-with-docker-content-trust">Sign with <code>docker trust sign</code></a></td>
<td><ul>
<li>Docker Engine - Enterprise 18.03 or higher</li>
<li>Docker Engine - Community 17.12 or higher</li>
</ul></td>
<td>Copy <code>ca.crt</code> from <code>https://&lt;dtr-external-url&gt;/ca</code> to:
<ul>
<li>Linux: <code>/etc/docker/certs.d/</code></li>
<li>Mac: <code>&lt;home_directory&gt;/.docker/certs.d/</code></li>
</ul></td>
</tr>
<tr class="even">
<td><a href="/engine/security/trust/content_trust/#runtime-enforcement-with-docker-content-trust">Enforce signature or hash verification on the Docker client</a></td>
<td><ul>
<li>Docker Engine - Enterprise 17.06 or higher</li>
<li>Docker Engine - Community 17.06 or higher</li>
</ul></td>
<td><code>export DOCKER_CONTENT_TRUST=1</code> to enable content trust on the Docker client. Copy <code>ca.crt</code> from <code>https://&lt;dtr-external-url&gt;/ca</code> to <code>/&lt;home_directory&gt;/.docker/tls/</code> on Linux and macOS. <code>docker push</code> will sign your images.</td>

</tr>
<tr class="odd">
<td><a href="/ee/dtr/user/manage-images/sign-images/">Sign images that UCP can trust</a></td>
<td><ul>
<li>Docker Engine - Enterprise 17.06 or higher</li>
<li>Docker UCP 2.2 or higher</li>
</ul></td>
<td>Configure UCP to <a href="/ee/ucp/admin/configure/run-only-the-images-you-trust/#configure-ucp">run only signed images</a>. See <a href="/ee/dtr/user/manage-images/">Sign an image</a> for detailed steps.</td>
</tr>
</tbody>
</table>

## DTR API access

With `curl`, you can interact with the DTR
API by passing a public certificate and private key pair instead of
your DTR username and password/authentication token.

```bash
curl --cert cert.pem --key key.pem  -X GET \
"https://<dtr-external-url>/api/v0/repositories?pageSize=10&count=false" \
-H "accept:application/json"
```

In the above example, `cert.pem` contains the public certificate and `key.pem`
contains the private key. For non-admin users, you can generate a client bundle from UCP or contact your administrator for your public and private key pair.

For Mac-specific quirks, see [curl on certain macOS versions](#curl-on-certain-macos-versions).

## Notary CLI operations with DTR

For establishing mutual trust between the Notary client and your trusted registry (DTR) using the Notary CLI, place your TLS client certificates in `<home_directory>/.docker/tls/<dtr-external-url>/` as `client.cert` and `client.key`. Note that the filenames must match. Pass the FQDN or publicly accessible IP address of your registry along with the TLS client certificate options to the Notary client. To get started, see [Use the Notary client for advanced users](/notary/advanced_usage/).

> ### Self-signed DTR server certificate
>
> Also place `ca.crt` in `<home_directory>/.docker/tls/<dtr-external-url>/` when you're using a self-signed server certificate for DTR.

## Troubleshooting tips

### DTR authentication via client Certificates

Hit your DTR's `basic_info` endpoint via `curl`:

```curl
curl --cert cert.pem --key key.pem -X GET "https://<dtr-external-url>/basic_info"
```

If successfully configured, you should see `TLSClientCertificate` listed as the `AuthnMethod` in the JSON response.

#### Example Response

```json
{
"CurrentVersion": "2.7.0",
"User": {
"name": "admin",
"id": "30f53dd2-763b-430d-bafb-dfa361279b9c",
"fullName": "",
"isOrg": false,
"isAdmin": true,
"isActive": true,
"isImported": false
},
"IsAdmin": true,
"AuthnMethod": "TLSClientCertificate"
}
```

### DTR as an insecure registry

Avoid adding DTR to Docker Engine's list of insecure registries as a workaround. This has the side effect of disabling the use of TLS certificates.

### DTR server certificate errors

#### Example Error

```bash
Error response from daemon: Get https://35.165.223.150/v2/: x509: certificate is valid for 172.17.0.1, not 35.165.223.150
```

- On the web UI, make sure to add the IP address or the FQDN associated with your custom TLS certificate under **System > General > Domains & Proxies**.

- From the command line interface, [reconfigure DTR](/reference/dtr/2.7/cli/reconfigure/) with the `--dtr-external-url` option and the associated PEM files for your certificate.

### Intermediate certificates

For chain of trust which includes intermediate certificates, you may optionally add those certificates when installing or reconfiguring DTR with `--enable-client-cert-auth` and `--client-cert-auth-ca`. You can do so by combining all of the certificates into a single PEM file.

### curl on certain macOS versions

Some versions of macOS include `curl` which only accepts `.p12` files and specifically requires a `./` prefix in front of the file name if running `curl` from the same directory as the `.p12` file:

```
curl --cert ./client.p12  -X GET \
"https://<dtr-external-url>/api/v0/repositories?pageSize=10&count=false" \
-H "accept:application/json"
```
