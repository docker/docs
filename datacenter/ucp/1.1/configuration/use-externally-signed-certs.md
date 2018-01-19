---
description: Learn how to configure Docker Universal Control Plane to use your own
  certificates.
keywords: Universal Control Plane, UCP, certificate, authentication, tls
redirect_from:
- /ucp/configuration/use-externally-signed-certs/
title: Use externally-signed certificates
---

Docker Universal Control Plane uses TLS to encrypt the traffic between users
and your cluster. By default this is done using self-signed certificates.
Since self-signed certificates are not trusted by web browsers, when users
access the UCP web UI, their browsers display a security warning. To avoid this,
you can configure UCP to use externally signed certificates.

This can be done while
[installing the UCP cluster](../installation/install-production.md) by
providing the externally signed certificates during the installation.
If you install UCP without providing externally signed certificates, then
self-signed certificates are used by default. These certificates can be replaced
at any time.

Since client certificate bundles are signed and verified with the UCP server
certificates, if you replace the UCP server certificates, users have to
download new client certificate bundles to run Docker commands on
the cluster.

## Replace existing certificates

To replace the server certificates used by UCP, for each controller node:

1.  Login into the node with ssh.
2.  In the directory where you have the keys and certificates to use, run:

    ```none
    # Create a container that attaches to the same volume where certificates are stored
    $ docker create --name replace-certs -v ucp-controller-server-certs:/data busybox

    # Copy your keys and certificates to the container's volumes
    $ docker cp cert.pem replace-certs:/data/cert.pem
    $ docker cp ca.pem replace-certs:/data/ca.pem
    $ docker cp key.pem replace-certs:/data/key.pem

    # Remove the container, since you don't need it any longer
    $ docker rm replace-certs
    ```

3.  Restart the `ucp-controller` container.

    To avoid downtime, don't restart all the `ucp-controller` containers of
    your cluster at the same time.

    ```bash
    $ docker restart ucp-controller
    ```

4.  Let your users know.

    After replacing the certificates your users can't authenticate
    with their old client certificate bundles. Ask your users to go to the UCP
    web UI and [get new client certificate bundles](../access-ucp/cli-based-access.md).