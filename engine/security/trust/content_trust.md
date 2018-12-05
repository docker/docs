---
description: Enabling content trust in Docker
keywords: content, trust, security, docker, documentation
title: Content trust in Docker
---

When transferring data among networked systems, *trust* is a central concern. In
particular, when communicating over an untrusted medium such as the internet, it
is critical to ensure the integrity and the publisher of all the data a system
operates on. You use Docker Engine to push and pull images (data) to a public or private registry. Content trust
gives you the ability to verify both the integrity and the publisher of all the
data received from a registry over any channel.

## About trust in Docker

Docker Content Trust (DCT) allows operations with a remote Docker registry to enforce
client-side signing and verification of image tags. DCT provides the
ability to use digital signatures for data sent to and received from remote
Docker registries. These signatures allow client-side verification of the
integrity and publisher of specific image tags.

Once DCT is enabled, image publishers can sign their images. Image consumers 
can ensure that the images they use are signed. Publishers and consumers can 
either be individuals or organizations. DCT supports users and automated processes 
such as builds.

When you enable DCT, signing occurs on the client after push and verification 
happens on the client after pull if you use Docker CE. If you use UCP, and you 
have configured UCP to require images to be signed before deploying, signing is 
verified by UCP.

### Image tags and DCT

An individual image record has the following identifier:

```
[REGISTRY_HOST[:REGISTRY_PORT]/]REPOSITORY[:TAG]
```

A particular image `REPOSITORY` can have multiple tags. For example, `latest` and
 `3.1.2` are both tags on the `mongo` image. An image publisher can build an image
 and tag combination many times changing the image with each build.

DCT is associated with the `TAG` portion of an image. Each image
repository has a set of keys that image publishers use to sign an image tag.
Image publishers have discretion on which tags they sign.

An image repository can contain an image with one tag that is signed and another
tag that is not. For example, consider [the Mongo image
repository](https://hub.docker.com/r/library/mongo/tags/). The `latest`
tag could be unsigned while the `3.1.6` tag could be signed. It is the
responsibility of the image publisher to decide if an image tag is signed or
not. In this representation, some image tags are signed, others are not:

![Signed tags](images/tag_signing.png)

Publishers can choose to sign a specific tag or not. As a result, the content of
an unsigned tag and that of a signed tag with the same name may not match. For
example, a publisher can push a tagged image `someimage:latest` and sign it.
Later, the same publisher can push an unsigned `someimage:latest` image. This second
push replaces the last unsigned tag `latest` but does not affect the signed `latest` version.
The ability to choose which tags they can sign, allows publishers to iterate over
the unsigned version of an image before officially signing it.

Image consumers can enable DCT to ensure that images they use were
signed. If a consumer enables DCT, they can only pull, run, or build
with trusted images. Enabling DCT is like wearing a pair of
rose-colored glasses. Consumers "see" only signed image tags and the less
desirable, unsigned image tags are "invisible" to them.

![Trust view](images/trust_view.png)

To the consumer who has not enabled DCT, nothing about how they
work with Docker images changes. Every image is visible regardless of whether it
is signed or not.


### DCT operations and keys

When DCT is enabled, `docker` CLI commands that operate on tagged images must
either have content signatures or explicit content hashes. The commands that
operate with DCT are:

* `push`
* `build`
* `create`
* `pull`
* `run`

For example, with DCT enabled a `docker pull someimage:latest` only
succeeds if `someimage:latest` is signed. However, an operation with an explicit
content hash always succeeds as long as the hash exists:

```bash
$ docker pull someimage@sha256:d149ab53f8718e987c3a3024bb8aa0e2caadf6c0328f1d9d850b2a2a67f2819a
```

Trust for an image tag is managed through the use of signing keys. A key set is
created when an operation using DCT is first invoked. A key set consists
of the following classes of keys:

- an offline key that is the root of DCT for an image tag
- repository or tagging keys that sign tags
- server-managed keys such as the timestamp key, which provides freshness
	security guarantees for your repository

The following image depicts the various signing keys and their relationships:

![Content Trust components](images/trust_components.png)

>**WARNING**:
> Loss of the root key is **very difficult** to recover from.
>Correcting this loss requires intervention from [Docker
>Support](https://support.docker.com) to reset the repository state. This loss
>also requires **manual intervention** from every consumer that used a signed
>tag from this repository prior to the loss.
{:.warning}

You should backup the root key somewhere safe. Given that it is only required
to create new repositories, it is a good idea to store it offline in hardware.
For details on securing, and backing up your keys, make sure you
read how to [manage keys for DCT](trust_key_mng.md).

## Survey of typical DCT operations

This section surveys the typical trusted operations users perform with Docker
images. Specifically, we go through the following steps to help us exercise
these various trusted operations:

* Build and push an unsigned image
* Pull an unsigned image
* Build and push a signed image
* Pull the signed image pushed above
* Pull unsigned image pushed above

### Enabling DCT in Docker Engine Configuration

Engine Signature Verification prevents the following behaviors on an image:
* Running a container to build an image (the base image must be signed, or must be scratch)
* Creating a container from an image that is not signed

DCT does not verify that a running container’s filesystem has not been altered from what 
was in the image. For example,  it does not prevent a container from writing to the filesystem, nor 
the container’s filesystem from being altered on disk.

It will also pull and run signed images from registries, but will not prevent unsigned images from being 
imported, loaded, or created.

The image name, digest, or tag must be verified if DCT is enabled. The latest DCT metadata for 
an image must be downloaded from the trust server associated with the registry:
* If an image tag does not have a digest, the DCT metadata translates the name to an image digest
* If an image tag has an image digest, the DCT metadata verifies that the name matches the provided digest
* If an image digest does not have an image tag, the DCT metadata does a reverse lookup and provides the image tag as well as the digest.

The signature verification feature is configured in the Docker daemon configuration file 
`daemon.json`.

```
{
    ...
    “content-trust”: {
        “trust-pinning”: {
            “root-keys”: {
                “myregistry.com/myorg/*”: [“keyID1”, “keyID2”],
                “myregistry.com/otherorg/repo”: [“keyID3”]
            },
            “official-images”: true,
        },
        “mode”: “disabled” | “permissive” | “enforced”,
        “allow-expired-trust-cache”: true,
    }
}
```

| **Stanza**                     | **Description**                        |
|--------------------------------|----------------------------------------|
| `trust-pinning:root-keys`      | Root key IDs are canonical IDs that sign the root metadata of the image trust data. In Docker Certified Trust (DCT), the root keys are unique certificates tying the name of the image to the repo metadata.  The private key ID (the canonical key ID) corresponding to the certificate does not depend on the image name. If an image’s name matches more than one glob, then the most specific (longest) one is chosen.            |
| `trust-pinning:library-images` | This option pins the official libraries (`docker.io/library/*`) to the hard-coded Docker official images root key. DCT trusts the official images by default. This is in addition to whatever images are specified by `trust-pinning:root-keys`. If `trustpinning:root-keys` specifies a key mapping for `docker.io/library/*`, those keys will be preferred for trust pinning. Otherwise, if a more general `docker.io/*` or `*` are specified, the official images key will be preferred.                                                                |
| `allow-expired-trust-cache`    | Specifies whether cached locally expired metadata validates images if an external server is unreachable or does not have image trust metadata. This is necessary for machines which may be often offline, as may be the case for edge. This does not provide mitigations against freeze attacks, which is a necessary to provide availability in low-connectivity environments.                                                |
| `mode`                         | Specifies whether DCT is enabled and enforced. Valid modes are: `disabled`: Verification is not active and the remainder of the content-trust related metadata will be ignored. *NOTE* that this is the default configuration if `mode` is not specified. `permissive`: Verification will be performed, but only failures will only be logged and remain unenforced. This configuration is intended for testing of changes related to content-trust. `enforced`: DCT will be enforced and an image that cannot be verified successfully will not be pulled or run.  |

***Note:*** The DCT configuration defined here is agnostic of any policy defined in 
[UCP](https://docs.docker.com/v17.09/datacenter/ucp/2.0/guides/content-trust/#configure-ucp). 
Images that can be deployed by the UCP trust policy but are disallowed by the Docker Engine 
configuration will not successfully be deployed or run on that engine.

### Enable and disable DCT per-shell or per-invocation

Instead of enabling DCT through the system-wide configuration, DCT can be enabled or disabled 
on a per-shell or per-invocation basis.  

To enable on a per-shell basis, enable the `DOCKER_CONTENT_TRUST` environment variable. 
Enabling per-shell is useful because you can have one shell configured for trusted operations 
and another terminal shell for untrusted operations. You can also add this declaration to 
your shell profile to have it enabled by default.

To enable DCT in a `bash` shell enter the following command:

```bash
export DOCKER_CONTENT_TRUST=1
```

Once set, each of the "tag" operations requires a key for a trusted tag.

In an environment where `DOCKER_CONTENT_TRUST` is set, you can use the
`--disable-content-trust` flag to run individual operations on tagged images
without DCT on an as-needed basis.

Consider the following Dockerfile that uses an untrusted parent image:

```
$  cat Dockerfile
FROM docker/trusttest:latest
RUN echo
```

To build a container successfully using this Dockerfile, one can do:

```
$  docker build --disable-content-trust -t <username>/nottrusttest:latest .
Sending build context to Docker daemon 42.84 MB
...
Successfully built f21b872447dc
```

The same is true for all the other commands, such as `pull` and `push`:

```
$  docker pull --disable-content-trust docker/trusttest:latest
...
$  docker push --disable-content-trust <username>/nottrusttest:latest
...
```

To invoke a command with DCT enabled regardless of whether or how the `DOCKER_CONTENT_TRUST` variable is set:

```bash
$  docker build --disable-content-trust=false -t <username>/trusttest:testing .
```

All of the trusted operations support the `--disable-content-trust` flag.


### Push trusted content

To create signed content for a specific image tag, simply enable DCT
and push a tagged image. If this is the first time you have pushed an image
using DCT on your system, the session looks like this:

```bash
$ docker push <username>/trusttest:testing
The push refers to a repository [docker.io/<username>/trusttest] (len: 1)
9a61b6b1315e: Image already exists
902b87aaaec9: Image already exists
latest: digest: sha256:d02adacee0ac7a5be140adb94fa1dae64f4e71a68696e7f8e7cbf9db8dd49418 size: 3220
Signing and pushing trust metadata
You are about to create a new root signing key passphrase. This passphrase
will be used to protect the most sensitive key in your signing system. Please
choose a long, complex passphrase and be careful to keep the password and the
key file itself secure and backed up. It is highly recommended that you use a
password manager to generate the passphrase and keep it safe. There will be no
way to recover this key. You can find the key in your config directory.
Enter passphrase for new root key with id a1d96fb:
Repeat passphrase for new root key with id a1d96fb:
Enter passphrase for new repository key with id docker.io/<username>/trusttest (3a932f1):
Repeat passphrase for new repository key with id docker.io/<username>/trusttest (3a932f1):
Finished initializing "docker.io/<username>/trusttest"
```

When you push your first tagged image with DCT enabled, the `docker`
client recognizes this is your first push and:

 - alerts you that it is creating a new root key
 - requests a passphrase for the root key
 - generates a root key in the `~/.docker/trust` directory
 - requests a passphrase for the repository key
 - generates a repository key in the `~/.docker/trust` directory

The passphrase you chose for both the root key and your repository key-pair
should be randomly generated and stored in a *password manager*.

> **NOTE**: If you omit the `testing` tag, DCT is skipped. This is true
even if DCT is enabled and even if this is your first push.

```bash
$ docker push <username>/trusttest
The push refers to a repository [docker.io/<username>/trusttest] (len: 1)
9a61b6b1315e: Image successfully pushed
902b87aaaec9: Image successfully pushed
latest: digest: sha256:a9a9c4402604b703bed1c847f6d85faac97686e48c579bd9c3b0fa6694a398fc size: 3220
No tag specified, skipping trust metadata push
```

It is skipped because as the message states, you did not supply an image `TAG`
value. In DCT, signatures are associated with tags.

Once you have a root key on your system, subsequent images repositories
you create can use that same root key:

```bash
$ docker push docker.io/<username>/otherimage:latest
The push refers to a repository [docker.io/<username>/otherimage] (len: 1)
a9539b34a6ab: Image successfully pushed
b3dbab3810fc: Image successfully pushed
latest: digest: sha256:d2ba1e603661a59940bfad7072eba698b79a8b20ccbb4e3bfb6f9e367ea43939 size: 3346
Signing and pushing trust metadata
Enter key passphrase for root key with id a1d96fb:
Enter passphrase for new repository key with id docker.io/<username>/otherimage (bb045e3):
Repeat passphrase for new repository key with id docker.io/<username>/otherimage (bb045e3):
Finished initializing "docker.io/<username>/otherimage"
```

The new image has its own repository key and timestamp key. The `latest` tag is signed with both of
these.


### Pull image content

A common way to consume an image is to `pull` it. With DCT enabled, the Docker
client only allows `docker pull` to retrieve signed images. Let's try to pull the image
you signed and pushed earlier:

```
$  docker pull <username>/trusttest:testing
Pull (1 of 1): <username>/trusttest:testing@sha256:d149ab53f871
...
Tagging <username>/trusttest@sha256:d149ab53f871 as docker/trusttest:testing
```

In the following example, the command does not specify a tag, so the system uses
the `latest` tag by default again and the `docker/trusttest:latest` tag is not signed.

```bash
$ docker pull docker/trusttest
Using default tag: latest
no trust data available
```

Because the tag `docker/trusttest:latest` is not trusted, the `pull` fails.

## Related information

* [Manage keys for content trust](trust_key_mng.md)
* [Automation with content trust](trust_automation.md)
* [Delegations for content trust](trust_delegation.md)
* [Play in a content trust sandbox](trust_sandbox.md)
