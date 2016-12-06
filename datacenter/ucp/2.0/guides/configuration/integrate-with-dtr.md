---
description: Integrate UCP with Docker Trusted Registry
keywords: trusted, registry, integrate, UCP, DTR
title: Integrate with Docker Trusted Registry
---

Docker UCP integrates out of the box with Docker Trusted Registry (DTR). This
allows you to use a UCP client bundle to push and pull images to DTR, without
having to login directly into DTR.

If you've configured DTR to use certificates issued by a globally-trusted
certificate authority you can skip this and use UCP client bundles to push and
pull images from DTR.

If you're using the DTR default configurations or configured DTR to use
self-signed certificates, you need to configure all UCP nodes to trust
the certificate authority that signed the DTR certificates. Otherwise
UCP won't trust DTR, and when trying to pull or push images, you'll get an
error:

```none
x509: certificate signed by unknown authority
```

## Configure UCP Docker Engines

The configuration depends on your operating system.

### Ubuntu/ Debian

```bash
# Download the DTR CA certificate
$ sudo curl -k https://<dtr-domain-name>/ca -o /usr/local/share/ca-certificates/<dtr-domain-name>.crt

# Refresh the list of certificates to trust
$ sudo update-ca-certificates

# Restart the Docker daemon
$ sudo service docker restart
```

### RHEL/ CentOS

```bash
# Download the DTR CA certificate
$ sudo curl -k https://<dtr-domain-name>/ca -o /etc/pki/ca-trust/source/anchors/<dtr-domain-name>.crt

# Refresh the list of certificates to trust
$ sudo update-ca-trust

# Restart the Docker daemon
$ sudo /bin/systemctl restart docker.service
```

## Test the integration

The best way to confirm that everything is well configured, is to pull and push
images from a UCP node to a private DTR repository.

1. Create a test repository on DTR.

    Navigate to the **DTR web UI**, and create a new **hello-world** repository
    so that you can push and pull images. Set it as **private**, and save
    the changes.

    ![](../images/dtr-integration-1.png)

2. Use a [UCP client bundle](../access-ucp/cli-based-access.md) to run docker
commands in the UCP cluster.

3.  Pull an image from Docker Hub:

    ```bash
    $ docker pull hello-world
    ```

4.  Retag the image:

    ```bash
    $ docker tag hello-world:latest <dtr-domain-name>/<username>/hello-world:1
    ```

5.  Push the image from the UCP node to your private registry:

    ```bash
    $ docker push <dtr-domain-name>/<username>/hello-world:1
    ```

6.  Validate that your image is now stored on DTR.

    When successfully pushing the image you should see a result like:

    ```none
    The push refers to a repository [dtr/username/hello-world]
    5f70bf18a086: Pushed
    33e7801ac047: Pushed
    1: digest: sha256:7d9e482c0cc9e68c7f07bf76e0aafcb1869d32446547909200db990e7bc5461a size: 1930
    ```

    You can also check that the tag exists on the DTR web UI.

    ![](../images/dtr-integration-2.png)

## Troubleshooting

When one of the components is misconfigured, and doesn't trust the root CA
certificate of the other components, you'll get an error like:

```none
$ docker push dtr/username/hello-world:1

The push refers to a repository [dtr/username/hello-world]
Get https://dtr/v1/_ping: x509: certificate signed by unknown authority
```

## Where to go next

* [Monitor your cluster](../monitor/index.md)
* [Troubleshoot your cluster](../monitor/troubleshoot.md)
* [Run only signed images](../content-trust/index.md)
