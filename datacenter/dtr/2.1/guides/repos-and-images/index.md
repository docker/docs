---
description: Learn how to configure your Docker Engine to push and pull images from
  Docker Trusted Registry.
keywords: docker, registry, TLS, certificates
title: Configure your Docker Engine
---

By default Docker Engine uses TLS when pushing and pulling images to an
image registry like Docker Trusted Registry.

If DTR is using the default configurations or was configured to use
self-signed certificates, you need to configure your Docker Engine to trust DTR.
Otherwise, when you try to login or push and pull images to DTR, you get an
error:

```none
$ docker login dtr.example.org

x509: certificate signed by unknown authority
```

The first step to make your Docker Engine trust the certificate authority used
by DTR is to get the DTR CA certificate. Then you configure your operating
system to trust that certificate.

## Configure your host

### macOS

In your browser navigate to `https://<dtr-url>/ca` to download the TLS
certificate used by DTR. Then
[add that certificate to the macOS trust store](https://support.apple.com/kb/PH18677?locale=en_US).

### Windows

In your browser navigate to `https://<dtr-url>/ca` to download the TLS
certificate used by DTR. Then
[add that certificate to the Windows trust store](https://technet.microsoft.com/en-us/library/cc754841(v=ws.11).aspx).


### Ubuntu/ Debian

```bash
# Download the DTR CA certificate
$ curl -k https://<dtr-domain-name>/ca -o /usr/local/share/ca-certificates/<dtr-domain-name>.crt
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

### Boot2Docker

1.  Login into the virtual machine with ssh:

    ```bash
    docker-machine ssh <machine-name>
    ```

2.  Create the `bootsync.sh` file, and make it executable:

    ```bash
    sudo touch /var/lib/boot2docker/bootsync.sh
    sudo chmod 755 /var/lib/boot2docker/bootsync.sh
    ```

3.  Add the following content to the `bootsync.sh` file. You can use nano or vi
    for this.

    ```bash
    #!/bin/sh

    cat /var/lib/boot2docker/server.pem >> /etc/ssl/certs/ca-certificates.crt
    ```

4.  Add the DTR CA certificate to the `server.pem` file:

    ```bash
    curl -k https://<dtr-domain-name>/ca | sudo tee -a /var/lib/boot2docker/server.pem
    ```

5.  Run `bootsync.sh` and restart the Docker daemon:

    ```bash
    sudo /var/lib/boot2docker/bootsync.sh
    sudo /etc/init.d/docker restart
    ```

## Login into DTR

To validate that your Docker daemon trusts DTR, trying authenticating against
DTR.

```bash
docker login dtr.example.org
```

## Where to go next

* [Pull an image from DTR](pull-an-image.md)
