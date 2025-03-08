---
title: Deploy caches with TLS
description: Learn how to deploy and secure caches for Docker Trusted Registry, leveraging TLS
keywords: docker, dtr, tls
---

When running DTR caches on a production environment, you should secure them
with TLS. In this example we're going to deploy a DTR cache that uses TLS.

DTR caches use the same configuration file format used by Docker Registry.
You can learn more about the supported configuration in the
[Docker Registry documentation](/registry/configuration.md#tls).


## Get the TLS certificate and keys

Before deploying a DTR cache with TLS, you need to get a public key
certificate for the domain name used to deploy the cache. You also
need the public and private key files for that certificate.

Once you have then, transfer those files to the host used to deploy
the DTR cache.


## Create the cache configuration

Use SSH to log into the host used to deploy the DTR cache, and navigate to
the directory where you've stored the TLS certificate and keys.

Create the `config.yml` file with the following content:

```
version: 0.1
storage:
  delete:
    enabled: true
  filesystem:
    rootdirectory: /var/lib/registry
http:
  addr: :5000
  tls:
    certificate: /certs/dtr-cache-ca.pem
    key: /certs/dtr-cache-key.pem
middleware:
  registry:
      - name: downstream
        options:
          blobttl: 24h
          upstreams:
            - originhost: https://<dtr-url>
          cas:
            - /certs/dtr-ca.pem
```

The configuration file mentions:

* /certs/dtr-cache-ca.pem: this is the public key certificate the cache will use
* /certs/dtr-cache-key.pem: this is the TLS private key
* /certs/dtr-ca.pem is the CA certificate used by DTR

Run this command to download the CA certificate used by DTR:

```
curl -k https://<dtr-url>/ca > dtr-ca.pem
```

Now that we've got the cache configuration file and TLS certificates, we can
deploy the cache by running:

```none
docker run --detach --restart always \
  --name dtr-cache \
  --publish 5000:5000 \
  --volume $(pwd)/dtr-cache-ca.pem:/certs/dtr-cache-ca.pem \
  --volume $(pwd)/dtr-cache-key.pem:/certs/dtr-cache-key.pem \
  --volume $(pwd)/dtr-ca.pem:/certs/dtr-ca.pem \
  --volume $(pwd)/config.yml:/config.yml \
  docker/dtr-content-cache:<version> /config.yml
```

## Use Let's Encrypt

You can also use Let's Encrypt to automatically generate TLS certificates that
are trusted by most clients.

Learn more about [Let's Encrypt](https://letsencrypt.org/how-it-works/), and
how to
[create a configuration file that leverages it](/registry/configuration.md#letsencrypt).


## Where to go next

* [Chain multiple caches](chaining.md)
