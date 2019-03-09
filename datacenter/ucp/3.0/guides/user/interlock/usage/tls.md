---
title: Applications with SSL
description: Learn how to configure your swarm services with TLS using the layer
  7 routing solution for UCP.
keywords: routing, proxy, tls
redirect_from:
  - /ee/ucp/interlock/usage/ssl/
---

Once the [layer 7 routing solution is enabled](../deploy/index.md), you can
start using it in your swarm services. You have two options for securing your
services with TLS:

* Let the proxy terminate the TLS connection. All traffic between end-users and
the proxy is encrypted, but the traffic going between the proxy and your swarm
service is not secured.
* Let your swarm service terminate the TLS connection. The end-to-end traffic
is encrypted and the proxy service allows TLS traffic to passthrough unchanged.

In this example we'll deploy a service that can be reached at `app.example.org`
using these two options.

No matter how you choose to secure your swarm services, there are two steps to
route traffic with TLS:

1. Create [Docker secrets](/engine/swarm/secrets.md) to manage from a central
place the private key and certificate used for TLS.
2. Add labels to your swarm service for UCP to reconfigure the proxy service.


## Let the proxy handle TLS

In this example we'll deploy a swarm service and let the proxy service handle
the TLS connection. All traffic between the proxy and the swarm service is
not secured, so you should only use this option if you trust that no one can
monitor traffic inside services running on your datacenter.

![TLS Termination](../../images/interlock-tls-1.png)

Start by getting a private key and certificate for the TLS connection. Make
sure the Common Name in the certificate matches the name where your service
is going to be available.

You can generate a self-signed certificate for `app.example.org` by running:

```bash
openssl req \
  -new \
  -newkey rsa:4096 \
  -days 3650 \
  -nodes \
  -x509 \
  -subj "/C=US/ST=CA/L=SF/O=Docker-demo/CN=app.example.org" \
  -keyout app.example.org.key \
  -out app.example.org.cert
```

Then, create a docker-compose.yml file with the following content:

```yml
version: "3.2"

services:
  demo:
    image: ehazlett/docker-demo
    deploy:
      replicas: 1
      labels:
        com.docker.lb.hosts: app.example.org
        com.docker.lb.network: demo-network
        com.docker.lb.port: 8080
        com.docker.lb.ssl_cert: demo_app.example.org.cert
        com.docker.lb.ssl_key: demo_app.example.org.key
    environment:
      METADATA: proxy-handles-tls
    networks:
      - demo-network

networks:
  demo-network:
    driver: overlay
secrets:
  app.example.org.cert:
    file: ./app.example.org.cert
  app.example.org.key:
    file: ./app.example.org.key
```

Notice that the demo service has labels describing that the proxy service should
route traffic to `app.example.org` to this service. All traffic between the
service and proxy takes place using the `demo-network` network. The service also
has labels describing the Docker secrets to use on the proxy service to terminate
the TLS connection.

Since the private key and certificate are stored as Docker secrets, you can
easily scale the number of replicas used for running the proxy service. Docker
takes care of distributing the secrets to the replicas.

Set up your CLI client with a [UCP client bundle](../../user-access/cli.md),
and deploy the service:

```bash
docker stack deploy --compose-file docker-compose.yml demo
```

The service is now running. To test that everything is working correctly you
first need to update your `/etc/hosts` file to map `app.example.org` to the
IP address of a UCP node.

In a production deployment, you'll have to create a DNS entry so that your
users can access the service using the domain name of your choice.
After doing that, you'll be able to access your service at:

```bash
https://<hostname>:<https-port>
```

Where:
* `hostname` is the name you used with the `com.docker.lb.hosts` label.
* `https-port` is the port you've configured in the [UCP settings](../deploy/index.md).

![Browser screenshot](../../images/interlock-tls-2.png){: .with-border}

Since we're using self-sign certificates in this example, client tools like
browsers display a warning that the connection is insecure.

You can also test from the CLI:

```bash
curl --insecure \
  --resolve <hostname>:<https-port>:<ucp-ip-address> \
  https://<hostname>:<https-port>/ping
```

If everything is properly configured you should get a JSON payload:

```json
{"instance":"f537436efb04","version":"0.1","request_id":"5a6a0488b20a73801aa89940b6f8c5d2"}
```

Since the proxy uses SNI to decide where to route traffic, make sure you're
using a version of curl that includes the SNI header with insecure requests.
If this doesn't happen, curl displays an error saying that the SSL handshake
was aborterd.


## Let your service handle TLS

You can also  encrypt the traffic from end-users to your swarm service.

![End-to-end encryption](../../images/interlock-tls-3.png)


To do that, deploy your swarm service using the following docker-compose.yml file:

```yml
version: "3.2"

services:
  demo:
    image: ehazlett/docker-demo
    command: --tls-cert=/run/secrets/cert.pem --tls-key=/run/secrets/key.pem
    deploy:
      replicas: 1
      labels:
        com.docker.lb.hosts: app.example.org
        com.docker.lb.network: demo-network
        com.docker.lb.port: 8080
        com.docker.lb.ssl_passthrough: "true"
    environment:
      METADATA: end-to-end-TLS
    networks:
      - demo-network
    secrets:
      - source: app.example.org.cert
        target: /run/secrets/cert.pem
      - source: app.example.org.key
        target: /run/secrets/key.pem

networks:
  demo-network:
    driver: overlay
secrets:
  app.example.org.cert:
    file: ./app.example.org.cert
  app.example.org.key:
    file: ./app.example.org.key
```

Notice that we've update the service to start using the secrets with the
private key and certificate. The service is also labeled with
`com.docker.lb.ssl_passthrough: true`, signaling UCP to configure the proxy
service such that TLS traffic for `app.example.org` is passed to the service.

Since the connection is fully encrypt from end-to-end, the proxy service
won't be able to add metadata such as version info or request ID to the
response headers.
