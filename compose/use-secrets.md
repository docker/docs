---
title: Using secrets in Compose
description: How to use secrets in Compose and their benefits
keywords: secrets, compose, security, environment variables
---

A secret is any piece of data, such as a password, certificates, or API key, that shouldn’t be transmitted over a network or stored unencrypted in a Dockerfile or in your application’s source code.

Docker Compose provides a way for you to use secrets without having to use environment variables to store information. If you’re injecting passwords and API keys as environment variables, you risk unintentional information exposure. Environment variables are often available to all processes, and it can be difficult to track access. They can also be printed in logs when debugging errors without your knowledge. Using secrets mitigates these risks.

Secrets are also useful for locally replicating a remote development setup with Kubernetes. Using secrets makes your application scalable. [TO_DO: EXPAND]

## Use secrets

Getting a secret into a container is a two-step process. First you need to define the secret, using the [top-level secrets attribute in your Compose file](compose-file/09-secrets.md). Then you update your service definitions to reference the secrets they require with the [secrets attribute](compose-file/05-services.md#secrets). Access to secrets is granted on a per-service basis.

Unlike the other methods, this guarantees that secrets are only available to the services that have been explicitly granted access, and secrets only exist in memory while that service is running.

## Examples

### Simple

```yaml
services:
  myapp:
    image: myapp:latest
    secrets:
      - my_secret
secrets:
  my_secret:
    file: ./my_secret.txt
```

In the example above, the frontend service is given access to the `my_secret` secret. The value of `my_secret` is set to the contents of the file `./my_secret.txt`.

### Advanced

```yaml
services:
   db:
     image: mysql:latest
     volumes:
       - db_data:/var/lib/mysql
     environment:
       MYSQL_ROOT_PASSWORD_FILE: /run/secrets/db_root_password
       MYSQL_DATABASE: wordpress
       MYSQL_USER: wordpress
       MYSQL_PASSWORD_FILE: /run/secrets/db_password
     secrets:
       - db_root_password
       - db_password

   wordpress:
     depends_on:
       - db
     image: wordpress:latest
     ports:
       - "8000:80"
     environment:
       WORDPRESS_DB_HOST: db:3306
       WORDPRESS_DB_USER: wordpress
       WORDPRESS_DB_PASSWORD_FILE: /run/secrets/db_password
     secrets:
       - db_password


secrets:
   db_password:
     file: db_password.txt
   db_root_password:
     file: db_root_password.txt

volumes:
    db_data:
```
In the advanced example above:

- The `secrets` attribute under each service defines the secrets you want to inject into the specific container.
- The top-level secrets section defines the variables `db_password` and `db_root_password` and provides the file that is used to populate their values.
- The deployment of each container means Docker creates a temporary filesystem mount under `/run/secrets/<secret_name>` with their specific values.

## Resources:

- [Secrets top-level element](compose-file/09-secrets.md)
- [Secrets attribute for services top-level element](compose-file/05-services.md#secrets)
