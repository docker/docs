---
title: Secure services with TLS
description: Learn how to configure your swarm services with TLS.
keywords: routing, proxy, tls
redirect_from:
  - /ee/ucp/interlock/usage/ssl/
---

After [deploying a layer 7 routing solution](../deploy/index.md), you have two options for securing your
services with TLS:

* [Let the proxy terminate the TLS connection.](#let-the-proxy-handle-tls) All traffic between end-users and
the proxy is encrypted, but the traffic going between the proxy and your swarm
service is not secured.
* [Let your swarm service terminate the TLS connection.](#let-your-service-handle-tls) The end-to-end traffic
is encrypted and the proxy service allows TLS traffic to passthrough unchanged.

Regardless of the option selected to secure swarm services, there are two steps required to
route traffic with TLS:

1. Create [Docker secrets](/engine/swarm/secrets.md) to manage from a central
place the private key and certificate used for TLS.
2. Add labels to your swarm service for UCP to reconfigure the proxy service.

## Let the proxy handle TLS
The following example deploys a swarm service and lets the proxy service handle
the TLS connection. All traffic between the proxy and the swarm service is
not secured, so use this option only if you trust that no one can
monitor traffic inside services running in your datacenter.

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

Notice that the demo service has labels specifying that the proxy service should
route `app.example.org` traffic to this service. All traffic between the
service and proxy takes place using the `demo-network` network. The service also
has labels specifying the Docker secrets to use on the proxy service for terminating
the TLS connection.

Because the private key and certificate are stored as Docker secrets, you can
easily scale the number of replicas used for running the proxy service. Docker
distributes the secrets to the replicas.

Set up your CLI client with a [UCP client bundle](../../user-access/cli.md),
and deploy the service:

```bash
docker stack deploy --compose-file docker-compose.yml demo
```

The service is now running. To test that everything is working correctly, update your `/etc/hosts` file to map `app.example.org` to the
IP address of a UCP node.

In a production deployment, you must create a DNS entry so that 
users can access the service using the domain name of your choice.
After creating the DNS entry, you can access your service:

```bash
https://<hostname>:<https-port>
```

FOr this example:
* `hostname` is the name you specified with the `com.docker.lb.hosts` label.
* `https-port` is the port you configured in the [UCP settings](../deploy/index.md).

![Browser screenshot](../../images/interlock-tls-2.png){: .with-border}

Because this example uses self-sign certificates, client tools like
browsers display a warning that the connection is insecure.

You can also test from the CLI:

```bash
curl --insecure \
  --resolve <hostname>:<https-port>:<ucp-ip-address> \
  https://<hostname>:<https-port>/ping
```

If everything is properly configured, you should get a JSON payload:

```json
{"instance":"f537436efb04","version":"0.1","request_id":"5a6a0488b20a73801aa89940b6f8c5d2"}
```

Because the proxy uses SNI to decide where to route traffic, make sure you are
using a version of `curl` that includes the SNI header with insecure requests.
Otherwise, `curl` displays an error saying that the SSL handshake
was aborted.

> **Note**: Currently there is no way to update expired certificates using this method.
> The proper way is to create a new secret then update the corresponding service. 

## Let your service handle TLS
The second option for securing with TLS involves encrypting traffic from end users to your swarm service.

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

The service is updated to start using the secrets with the
private key and certificate. The service is also labeled with
`com.docker.lb.ssl_passthrough: true`, signaling UCP to configure the proxy
service such that TLS traffic for `app.example.org` is passed to the service.

Since the connection is fully encrypted from end-to-end, the proxy service
cannot add metadata such as version information or request ID to the
response headers.
