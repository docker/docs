---
title: Use your own TLS certificates
description: Learn how to configure Docker Trusted Registry with your own TLS certificates.
keywords: dtr, tls, certificates, security
---

>{% include enterprise_label_shortform.md %}

Docker Trusted Registry (DTR) services are exposed using HTTPS by default. This
ensures encrypted communications between clients and your trusted registry. If
you do not pass a PEM-encoded TLS certificate during installation, DTR will
generate a [self-signed
certificate](https://en.wikipedia.org/wiki/Self-signed_certificate). This leads
to an insecure site warning when accessing DTR through a browser. Additionally,
DTR includes an [HSTS (HTTP Strict-Transport-Security)
header](https://en.wikipedia.org/wiki/HTTP_Strict_Transport_Security) in all
API responses which can further lead to your browser refusing to load DTR's web
interface.

You can configure DTR to use your own TLS certificates, so that it is
automatically trusted by your users' browser and client tools. As of v2.7, you
can also [enable user authentication via client
certificates](/ee/enable-client-certificate-authentication/) provided by your
organization's public key infrastructure (PKI).

## Replace the server certificates

You can upload your own TLS certificates and keys using the web interface, or pass them as CLI options when installing or reconfiguring your DTR instance.

### Web interface

Navigate to `https://<dtr-url>` and log in with your credentials. Select **System** from the left navigation pane, and scroll down to **Domain & Proxies**. 

![](/ee/dtr/images/use-your-certificates-1.png){: .with-border}

Enter your DTR domain name and upload or copy and paste the certificate details:

* ***Load balancer/public address.*** The domain name clients will use to access DTR.
* ***TLS private key.*** The server private key.
* ***TLS certificate chain.*** The server certificate and any intermediate public
certificates from your certificate authority (CA). This certificate needs to be valid for the DTR public address,
and have SANs for all addresses used to reach the DTR replicas, including load
balancers.
* ***TLS CA.*** The root CA public certificate.

Click **Save** to apply your changes.

If you've added certificates issued by a globally trusted CA,
any web browser or client tool should now trust DTR. If you're using an internal
CA, you will need to configure the client systems to trust that
CA.

### Command line interface

See [docker/dtr install](/reference/dtr/2.7/cli/install/) and [docker/dtr reconfigure](/reference/dtr/2.7/cli/reconfigure/) for TLS certificate options and usage. 

## Where to go next
- [Enable single sign-on](enable-single-sign-on)
- [Set up external storage](external-storage)
