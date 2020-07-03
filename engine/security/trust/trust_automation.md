---
description: Automating content push pulls with trust
keywords: trust, security, docker, documentation, automation
title: Automation with content trust
---

It is very common for Docker Content Trust to be built into existing automation
systems. To allow tools to wrap Docker and push trusted content, there are 
environment variables that can be passed through to the client. 

This guide follows the steps as described 
[here](content_trust/#signing-images-with-docker-content-trust) so please read 
that and understand its prerequisites. 

When working directly with the Notary client, it uses its [own set of environment variables](../../../notary/reference/client-config.md#environment-variables-optional).

## Add a delegation private key

To automate importing a delegation private key to the local Docker trust store, we 
need to pass a passphrase for the new key. This passphrase will be required 
everytime that delegation signs a tag. 

```
$ export DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE="mypassphrase123"

$ docker trust key load delegation.key --name jeff
Loading key from "delegation.key"...
Successfully imported key from delegation.key
```

## Add a delegation public key

If you initialising a repository at the same time as adding a Delegation
public key, then you will need to use the local Notary Canonical Root Key's 
passphrase to create the repositories trust data. If the repository has already 
been initiated then you only need the repositories passphrase. 

```
# Export the Local Root Key Passphrase if required.
$ export DOCKER_CONTENT_TRUST_ROOT_PASSPHRASE="rootpassphrase123"

# Export the Repository Passphrase
$ export DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE="repopassphrase123"

# Initialise Repo and Push Delegation
$ docker trust signer add --key delegation.crt jeff registry.example.com/admin/demo
Adding signer "jeff" to registry.example.com/admin/demo...
Initializing signed repository for registry.example.com/admin/demo...
Successfully initialized "registry.example.com/admin/demo"
Successfully added signer: registry.example.com/admin/demo
```

## Sign an image

Finally when signing an image, we will need to export the passphrase of the 
signing key. This was created when the key was loaded into the local Docker 
trust store with `$ docker trust key load`.

```
$ export DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE="mypassphrase123"

$ docker trust sign registry.example.com/admin/demo:1
Signing and pushing trust data for local image registry.example.com/admin/demo:1, may overwrite remote trust data
The push refers to repository [registry.example.com/admin/demo]
428c97da766c: Layer already exists
2: digest: sha256:1a6fd470b9ce10849be79e99529a88371dff60c60aab424c077007f6979b4812 size: 524
Signing and pushing trust metadata
Successfully signed registry.example.com/admin/demo:1
```

## Build with content trust

You can also build with content trust. Before running the `docker build` command, 
you should set the environment variable `DOCKER_CONTENT_TRUST` either manually or 
in a scripted fashion. Consider the simple Dockerfile below.

```dockerfile
FROM docker/trusttest:latest
RUN echo
```

The `FROM` tag is pulling a signed image. You cannot build an image that has a
`FROM` that is not either present locally or signed. Given that content trust
data exists for the tag `latest`, the following build should succeed:

```bash
$  docker build -t docker/trusttest:testing .
Using default tag: latest
latest: Pulling from docker/trusttest

b3dbab3810fc: Pull complete
a9539b34a6ab: Pull complete
Digest: sha256:d149ab53f871
```

If content trust is enabled, building from a Dockerfile that relies on tag 
without trust data, causes the build command to fail:

```bash
$  docker build -t docker/trusttest:testing .
unable to process Dockerfile: No trust data for notrust
```

## Related information

* [Delegations for content trust](trust_delegation.md)
* [Content trust in Docker](content_trust.md)
* [Manage keys for content trust](trust_key_mng.md)
* [Play in a content trust sandbox](trust_sandbox.md)
