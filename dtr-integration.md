<!--[metadata]>
+++
title = "Integrate with Trusted Registry"
description = "Integrate UCP with Docker Trusted Registry"
keywords = ["trusted, registry, integrate, UCP, DTR"]
[menu.main]
parent="mn_ucp"
+++
<![end-metadata]-->


# Integrate UCP with Docker Trusted Registry

This page explains how to integrate Universal Control Plane (UCP) with the
Docker Trusted Registry (DTR). Trusted Registry is a image storage and
management service that you can install within your company's private
infrastructure.

## Prerequisites

You must have already installed DTR on your infrastructure before performing
this procedure. The DTR server and the UCP controller must be able to
communicate over your network infrastructure.

The Universal Control Plane and Trusted Registry are both part of the Docker
Datacenter solution. This means the license you use for UCP works with DTR or,
if you have a DTR license, it also works with UCP.

## Step 1: (Optional) Prepare a cert script

If you are using a self-signed or third-party CA, you need to the prepare a
`cert_create.sh` script. You'll use this script to install the self-signed cert
on the nodes in your UCP cluster.

1. Create a file called `cert_create.sh` with your favorite editor.

2. Add the following to content to the file.

        DTR_HOST="<my-dtr-host-dns-name"
        mkdir -p ~/.docker/tls/${DTR_HOST}
        openssl s_client -host ${DTR_HOST} -port 443 </dev/null 2>/dev/null | openssl x509 -outform PEM | tee ~/.docker/tls/${DTR_HOST}/ca.crt

3.  Replace the `<my-dtr-host-dns-name>` value with the fully qualified DNS
    value for your DTR instance.

4. Save and close the `cert_create.sh` file.

4. Set execute permission on the file.

        $ chmod 755 cert_create.sh

## Step 2: Configure DTR and UCP

In this step, you configure DTR and UCP to communicate. To do this you need an admin level certificate bundle for UCP or terminal access to the UCP controller.  

1. Log into or connect to the UCP primary controller.

2. Generate the UCP certificates using the `ucp dump-certs` command.

    This command generates the certificates for the Swarm cluster.

        $ docker run --rm -it --name ucp -v /var/run/docker.sock:/var/run/docker.sock docker/ucp dump-certs  --cluster -ca > /tmp/cluster-root-chain.pem

3. Cat or edit the `cluster-root-chain.pem` file.

4. Copy the certificate.

    This example illustrates what you should copy, your installation certificate
    will be different.
    
        -----BEGIN CERTIFICATE-----
        MIIFGDCCAwCgAwIBAgIIIQjwMnZnj2gwDQYJKoZIhvcNAQENBQAwGDEWMBQGA1UE
        AxMNU3dhcm0gUm9vdCBDQTAeFw0xNjAyMTAxNzQzMDBaFw0yMTAyMDgxNzQzMDBa
        MBgxFjAUBgNVBAMTDVN3YXJtIFJvb3QgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4IC
        DwAwggIKAoICAQC5UtvO/xju7INdZkXA9TG7T6JYo1CIf5yZz9LZBDrexSAx7uPi
        7b5YmWGUA26VgBDvAFuLuQNRy/OlITNoFIEG0yovw6waLcqr597ox9d9jeaJ4ths
        ...<output snip>...
        2wDuqlzByRVTO0NL4BX0QV1J6LFtrlWU92WxTcOV8T7Zc4mzQNMHfiIZcHH/p3+7
        cRA7HVdljltI8UETcrEvTKb/h1BiPlhzpIfIHwMdA2UScGgJlaH7wA0LpeJGWtUc
        AKrb2kTIXNQq7phH
        -----END CERTIFICATE-----

5. Login to the Trusted Registry dashboard as a user.

6. Choose **Settings > General** page.

7. Locate the **Auth Bypass TLS Root CA** field.

8. Paste certificate you copied into the field.

9. (Optional) If you are using a self-signed or third-party CA, do the following
on each node in your UCP cluster:     

    a. Log into a UCP node using an account with `sudo` privileges.

    b. Copy the `cert_create.sh`to the node.

    c. Run the `cert_create.sh` on the node.

        $ sudo cert_create.sh

    d. Verify the cert was created.

        $ sudo cat ~/.docker/tls/${DTR_HOST}/ca.crt


## Step 2: Confirm the integration

The best way to confirm the integration is to push and pull an image from a UCP node.

1. Open a terminal session on a UCP node.

2. Pull Docker's `hello-world` image.

        $ docker pull hello-world
        Using default tag: latest
        latest: Pulling from library/hello-world
        03f4658f8b78: Pull complete
        a3ed95caeb02: Pull complete
        Digest: sha256:8be990ef2aeb16dbcb9271ddfe2610fa6658d13f6dfb8bc72074cc1ca36966a7
        Status: Downloaded newer image for hello-world:latest

3. Get the `IMAGE ID` of the `hello-world` image.

        $ docker images
        REPOSITORY          TAG                        IMAGE ID            CREATED             SIZE
        hello-world         latest                     690ed74de00f        4 months ago        960 B

5. Retag the `hello-world` image with a new tag.

    The syntax for tagging an image is:

        docker tag <ID> <username>/<image-name>:<tag>

    Make sure to replace `<username>` with your actual username and the <ID>
    with the ID of the `hello-world` image you pulled.

        $ docker tag 690ed74de00f username/hello-world:test

4. Login into the DTR instance from the command line.

    The example below uses `mydtr.company.com` as the URL for the DTR instance.
    Your's will be different.  

        $  docker login mydtr.company.com

    Provide your username password when prompted.

5. Push your newly tagged image to the DTR instance.

    The following is an example only, substitute your DTR URL and username when
    you run this command.            

        $ docker push  mydtr.company.com/username/hello-world:test


## Troubleshooting section

This section details common problems you can encounter when working with the DTR /
UCP integration.

### Unknown authority error on push

Example:

```
% docker push mydtr.acme.com/jdoe/myrepo:latest
The push refers to a repository [mydtr.acme.com/jdoe/myrepo]
unable to ping registry endpoint https://mydtr.acme.com/v0/
v2 ping attempt failed with error: Get https://mydtr.acme.com/v2/: x509: certificate signed by unknown authority
v1 ping attempt failed with error: Get https://mydtr.acme.com/v1/_ping: x509: certificate signed by unknown authority
```

Review the trust settings in DTR and make sure they are correct. Try repasting
the first PEM block from the `chain.pem` file.

### Authentication required

Example:

```
% docker push mydtr.acme.com/jdoe/myrepo:latest
The push refers to a repository [mydtr.acme.com/jdoe/myrepo]
5f70bf18a086: Preparing
2c84284818d1: Preparing
unauthorized: authentication required
```

You must login before you can push to DTR.
