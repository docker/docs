---
title: Secrets top-level elements
description: Explore all the attributes the secrets top-level element can have.
keywords: compose, compose specification, secrets, compose file reference
---

Secrets are a flavor of [Configs](08-configs.md) focusing on sensitive data, with specific constraint for this usage. 

Services can only access secrets when explicitly granted by a [`secrets`](05-services.md#secrets) attribute within the `services` top-level element.

The top-level `secrets` declaration defines or references sensitive data that is granted to the services in your Compose
application. The source of the secret is either `file` or `environment`.

- `file`: The secret is created with the contents of the file at the specified path.
- `environment`: The secret is created with the value of an environment variable.
- `external`: If set to true, `external` specifies that this secret has already been created. Compose does
  not attempt to create it, and if it does not exist, an error occurs.
- `name`: The name of the secret object in Docker. This field can be used to
  reference secrets that contain special characters. The name is used as is
  and isn't scoped with the project name.

## Example 1

`server-certificate` secret is created as `<project_name>_server-certificate` when the application is deployed,
by registering content of the `server.cert` as a platform secret.

```yml
secrets:
  server-certificate:
    file: ./server.cert
```

## Example 2 

`token` secret  is created as `<project_name>_token` when the application is deployed,
by registering the content of the `OAUTH_TOKEN` environment variable as a platform secret.

```yml
secrets:
  token:
    environment: "OAUTH_TOKEN"
```

Alternatively, `server-certificate` can be declared as external. Compose looks up the `server-certificate` secret to expose to relevant services.

```yml
secrets:
  server-certificate:
    external: true
```

## Example 3

External secrets lookup can also use a distinct key by specifying a `name`. 

The following example modifies the previous example to look up a secret using the name `CERTIFICATE_KEY`. The actual lookup key is set at deployment time by the [interpolation](12-interpolation.md) of
variables, but exposed to containers as hard-coded ID `server-certificate`.

```yml
secrets:
  server-certificate:
    external: true
    name: "${CERTIFICATE_KEY}"
```

If `external` is set to `true`, all other attributes apart from `name` are irrelevant. If Compose detects any other attribute, it rejects the Compose file as invalid.

Your Compose file needs to explicitly grant access to the secrets to relevant services in your application.