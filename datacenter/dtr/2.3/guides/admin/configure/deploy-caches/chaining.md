---
title: Chain multiple caches
description: Learn how to deploy and chain multiple caches for Docker Trusted Registry, to cover multiple regions or offices
keywords: dtr, tls
---

If your users are distributed geographically, consider chaining multiple DTR
caches together for faster pulls.

![cache chaining](../../../images/chaining-1.svg)

Too many levels of chaining might slow down pulls, so you should try different
configurations and benchmark them, to find out the right configuration.

This example shows how to configure two caches. A dedicated cache for
the Asia region that pulls images directly from DTR, and a cache for China, that
pulls images from the Asia cache.

## Cache for the Asia region

This cache has TLS, and pulls images directly from DTR:

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
    certificate: /certs/asia-ca.pem
    key: /certs/asia-key.pem
middleware:
  registry:
      - name: downstream
        options:
          blobttl: 24h
          upstreams:
            - https://<dtr-url>
          cas:
            - /certs/dtr-ca.pem
```

## Cache for the China region

This cache has TLS, and pulls images from the Asia cache:

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
    certificate: /certs/china-ca.pem
    key: /certs/china-key.pem
middleware:
  registry:
      - name: downstream
        options:
          blobttl: 24h
          upstreams:
            - https://<asia-cache-url>
          cas:
            - /certs/asia-cache-ca.pem
```

Since the China cache doesn't need to communicate directly with DTR,
it only needs to trust the CA certificates for the next hop, in this case
the CA certificate used by the Asia cache.
