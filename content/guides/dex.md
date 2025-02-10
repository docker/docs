---
title: Mocking OAuth services in testing with Dex
description: &desc Mocking OAuth services in testing with Dex
keywords: Dex, container-supported development
linktitle: Mocking OAuth services with Dex
summary: *desc
tags: [app-dev, distributed-systems]
languages: []
params:
  time: 10 minutes
---

Dex is an open-source OpenID Connect (OIDC) and OAuth 2.0 identity provider that can be configured to authenticate against various backend identity providers, such as LDAP, SAML, and OAuth. Running Dex in a Docker container allows developers to simulate an OAuth 2.0 server for testing and development purposes. This guide will walk you through setting up Dex as an OAuth mock server using Docker containers.

Nowadays OAuth is the preferred choice to authenticate in web services, the highest part of them give the possibility to access using popular OAuth services like Github, Google or Apple. Using OAuth guarantees a higher level of security and simplification since it is not necessary to create new profiles for each service. This means that, by allowing applications to access resources on behalf of users without sharing passwords, OAuth minimizes the risk of credential exposure.

In this guide, you'll learn how to:

- Use Docker to launch up a Dex container.
- Use mock OAuth in the local development without relying on an external OAuth provider.

## Using Dex with Docker

The official [Docker image for Dex](https://hub.docker.com/r/dexidp/dex/) provides a convenient way to deploy and manage Dex instances. Dex is available for various CPU architectures, including amd64, armv7, and arm64, ensuring compatibility with different devices and platforms. You can learn more about Dex standalone on the [Dex docs site](https://dexidp.io/docs/getting-started/).

### Prerequisites

[Docker Compose](https://docs.docker.com/compose/): Recommended for managing multi-container Docker applications.

### Setting Up Dex with Docker

Begin by creating a directory for your Dex project:

```bash
mkdir dex-mock-server
cd dex-mock-server
```
Organize your project with the following structure:

```bash
dex-mock-server/
├── config.yaml
└── docker-compose.yaml
```

Create the Dex Configuration File:

The config.yaml file defines Dex's settings, including connectors, clients, and storage. For a mock server setup, you can use the following minimal configuration:

```yaml
# config.yaml
issuer: http://localhost:5556/dex
storage:
  type: memory
web:
  http: 0.0.0.0:5556
staticClients:
  - id: example-app
    redirectURIs:
      - 'http://localhost:5555/callback'
    name: 'Example App'
    secret: ZXhhbXBsZS1hcHAtc2VjcmV0
enablePasswordDB: true
staticPasswords:
  - email: "admin@example.com"
    hash: "$2a$10$2b2cU8CPhOTaGrs1HRQuAueS7JTT5ZHsHSzYiFPm1leZck7Mc8T4W"
    username: "admin"
    userID: "1234"
```

Explanation:

- issuer: The public URL for Dex.

- storage: Using in-memory storage for simplicity.

- web: Dex will listen on port 5556.

- staticClients: Defines a client application (example-app) with its redirect URI and secret.

- enablePasswordDB: Enables static password authentication.

- staticPasswords: Defines a static user for authentication. The hash is a bcrypt hash of the password.

> Note: Ensure the hash is a valid bcrypt hash of your desired password. You can generate this using tools like [bcrypt-generator.com](https://bcrypt-generator.com/)
or use CLI tools like [htpasswd](https://httpd.apache.org/docs/2.4/programs/htpasswd.html) like in this following example:
```bash
echo password | htpasswd -BinC 10 admin | cut -d: -f2
```

Running Dex

With Docker Compose configured, start Dex:

```yaml
# docker-compose.yaml

services:
  dex:
    image: dexidp/dex:latest
    container_name: dex
    ports:
      - "5556:5556"
    volumes:
      - ./config.yaml:/etc/dex/config.yaml
    command: ["dex", "serve", "/etc/dex/config.yaml"]
```

Now it is possible to run the container using `docker compose` command.
```bash
docker compose up -d
```

This command will download the Dex Docker image (if not already available) and start the container in detached mode.


To Verify that Dex is running, check the logs to ensure Dex started successfully:

```bash
docker-compose logs -f dex
```
You should see output indicating that Dex is listening on the specified port.

Testing the OAuth Flow
Prepare a Test Application:

To test the OAuth flow, you'll need a client application configured to authenticate against Dex. Dex provides an example app that you can use for this purpose.

Clone the Dex Repository:

