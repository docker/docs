---
title: Using Docker Compose with OCI artifacts
linkTitle: Use Compose OCI artifacts applications
weight: 20
description: How to start and publish Compose applications as OCI artifacts
keywords: cli, compose, oci
aliases:
- /compose/oci-artifact/
---

{{% include "compose/oci-artifact.md" %}}

## Starting an OCI artifact application

To start a Docker Compose application using an OCI artifact, you can use the `-f` (or `--file`) flag followed by the OCI artifact reference.
This allows you to specify a Compose file stored as an OCI artifact in a registry.  
The `oci://` prefix indicates that the Compose file should be pulled from an OCI-compliant registry rather than loaded from the local filesystem.

```bash
$ docker compose -f oci://docker.io/username/my-compose-app:latest up
```

To run the Compose application, use the `docker compose` command with the `-f` flag pointing to your OCI artifact:
```bash
$ docker compose -f oci://docker.io/username/my-compose-app:latest up
```

### Warnings/Messages displayed

When you run an application from an OCI artifact, Compose may display warning messages requiring your confirmation to limit risks of running a malicious application:
* Listing the interpolation variables used along with their values
* Listing all environment variables used by the application
* Let you know if your OCI artifact application is using another remote resources (via [`include`](/reference/compose-file/include/) for example)

```bash 
$ REGISTRY=myregistry.com docker compose -f oci://docker.io/username/my-compose-app:latest up

Found the following variables in configuration:
VARIABLE     VALUE                SOURCE        REQUIRED    DEFAULT
REGISTRY     myregistry.com      command-line   yes         
TAG          v1.0                environment    no          latest
DOCKERFILE   Dockerfile          default        no          Dockerfile
API_KEY      <unset>             none           no          

Do you want to proceed with these variables? [Y/n]:y

Warning: This Compose project includes files from remote sources:
- oci://registry.example.com/stack:latest
Remote includes could potentially be malicious. Make sure you trust the source.
Do you want to continue? [y/N]: 
```

If you agree to start the application, Compose will display the directory where all the resources from the OCI artifact have been downloaded.
```bash
...
Do you want to continue? [y/N]: y

Your compose stack "oci://registry.example.com/stack:latest" is stored in "~/Library/Caches/docker-compose/964e715660d6f6c3b384e05e7338613795f7dcd3613890cfa57e3540353b9d6d"
```
---

## Publishing your Compose application as an OCI artifact

To distribute your Compose application as an OCI artifact, you can **publish** it to an OCI-compliant registry. 
This allows others to deploy your application directly from the registry.

The publish function supports most of the composition capabilities of Compose, like overrides, extends or include, [with some limitations](#limitations-and-considerations)

### Steps

1. Navigate to Your Compose Application Directory  
   Ensure you're in the directory containing your `compose.yml` file or that you are specifying your Compose file with the `-f` flag.

2. Log in to Docker Hub
   Before publishing, make sure you're authenticated with Docker Hub:

   ```bash
   $ docker login
   ```

3. Publish the Compose application to Docker Hub  
   Use the `docker compose publish` command to push your application as an OCI artifact:

   ```bash
   $ docker compose publish username/my-compose-app:latest
   ```
   or passing multiple Compose files
   ```bash
   $ docker compose -f compose-base.yml -f compose-production.yml publish username/my-compose-app:latest
   ```
When publishing you can use options to specify the OCI version, whether to resolve image digests and if you want to include environment variables: 
* `--oci-version`: Specify the OCI version (default is automatically determined).
* `--resolve-image-digests`: Pin image tags to digests.
* `--with-env`: Include environment variables in the published OCI artifact.

Compose checks for you if there isn't any sensitive data in your configuration and displays your environment variables to confirm you want to publish them.

```bash
...
you are about to publish sensitive data within your OCI artifact.
please double check that you are not leaking sensitive data
AWS Client ID
"services.serviceA.environment.AWS_ACCESS_KEY_ID": xxxxxxxxxx
AWS Secret Key
"services.serviceA.environment.AWS_SECRET_ACCESS_KEY": aws"xxxx/xxxx+xxxx+"
Github authentication
"GITHUB_TOKEN": ghp_xxxxxxxxxx
JSON Web Token
"": xxxxxxx.xxxxxxxx.xxxxxxxx
Private Key
"": -----BEGIN DSA PRIVATE KEY-----
xxxxx
-----END DSA PRIVATE KEY-----
Are you ok to publish these sensitive data? [y/N]:y

you are about to publish environment variables within your OCI artifact.
please double check that you are not leaking sensitive data
Service/Config  serviceA
FOO=bar
Service/Config  serviceB
FOO=bar
QUIX=
BAR=baz
Are you ok to publish these environment variables? [y/N]: 
```

If you refuse the publish process will stop without sending anything to the registry.

---

## Limitations and considerations

There is limitations to publishing Compose applications as OCI artifacts:
* You can't publish Compose configuration with service(s) containing bind mounts
* You can't publish Compose configuration with service(s) containing only `build` section
* You can't publish Compose configuration using `include` of local files, publish them as well as remote `include` is supported
