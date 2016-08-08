<!--[metadata]>
+++
title = "Use externally-signed certificates"
description = "Learn how to configure Docker Universal Control Plane to use your own certificates."
keywords = ["Universal Control Plane, UCP, certificate, authentiation, tls"]
[menu.main]
parent="mn_ucp_configuration"
identifier="ucp_configure_certs"
weight=0
+++
<![end-metadata]-->

# Use externally-signed certificates

By default the UCP web UI is exposed using HTTPS, to ensure all
communications between clients and the cluster are encrypted. Since UCP
controllers use self-signed certificates for this, when a client accesses
UCP their browsers won't trust this certificate, so the browser displays a
warning message.

You can configure UCP to use your own certificates, so that it is automatically
trusted by your users' browser and client tools.

To ensure minimal impact to your business, you should plan for this change to
happen outside business peak hours. Your applications will continue
running normally, but UCP will be unresponsive while the controller containers
are restarted.

## Replace the server certificates

To configure UCP to use your own certificates and keys, go to the
**UCP web UI**, navigate to the **Admin Settings** page,
and click **Certificates**.

![](../images/use-externally-signed-certs-1.png)

Upload your certificates and keys:

* A ca.pem file with the root CA public certificate.
* A cert.pem file with the server certificate and any intermediate CA public
certificates. This certificate should also have SANs for all addresses used to
reach the UCP controller, including load balancers.
* A key.pem file with server private key.

Finally, click **Update** for the changes to take effect.

After replacing the certificates your users won't be able to authenticate
with their old client certificate bundles. Ask your users to go to the UCP
web UI and [get new client certificate bundles](../access-ucp/cli-based-access.md).

## Where to go next

* [Access UCP from the CLI](../access-ucp/cli-based-access.md)
