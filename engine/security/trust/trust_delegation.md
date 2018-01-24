---
description: Delegations for content trust
keywords: trust, security, delegations, keys, repository
title: Delegations for content trust
---

Docker Engine supports the usage of the `targets/releases` delegation as the
canonical source of a trusted image tag.

Using this delegation allows you to collaborate with other publishers without
sharing your repository key, which is a combination of your targets and snapshot keys.
See [Manage keys for content trust](trust_key_mng.md) for more information).
Collaborators can keep their own delegation keys private.

The `targets/releases` delegation is currently an optional feature - in order
to set up delegations, you must use the Notary CLI:

1. [Download the client](https://github.com/docker/notary/releases) and ensure that it is
available on your path

2. Create a configuration file at `~/.notary/config.json` with the following content:

	```
	{
	  "trust_dir" : "~/.docker/trust",
	  "remote_server": {
	    "url": "https://notary.docker.io"
	  }
	}
	```

	This tells Notary where the Docker Content Trust data is stored, and to use the
	Notary server used for images in Docker Hub.

For more detailed information about how to use Notary outside of the default
Docker Content Trust use cases, refer to the
[Notary CLI documentation](/notary/getting_started.md).

When publishing and listing delegation changes using the Notary client,
your Docker Hub credentials are required.

## Generating delegation keys

Your collaborator needs to generate a private key (either RSA or ECDSA)
and give you the public key so that you can add it to the `targets/releases`
delegation.

The easiest way for them to generate these keys is with OpenSSL.
Here is an example of how to generate a 2048-bit RSA portion key (all RSA keys
must be at least 2048 bits):

```
$ openssl genrsa -out delegation.key 2048
Generating RSA private key, 2048 bit long modulus
....................................................+++
............+++
e is 65537 (0x10001)

```

They should keep `delegation.key` private because it is used to sign tags.

Then they need to generate an x509 certificate containing the public key, which is
what you need from them. Here is the command to generate a CSR (certificate
signing request):

```
$ openssl req -new -sha256 -key delegation.key -out delegation.csr
```

Then they can send it to whichever CA you trust to sign certificates, or they
can self-sign the certificate (in this example, creating a certificate that is
valid for 1 year):

```
$ openssl x509 -req -sha256 -days 365 -in delegation.csr -signkey delegation.key -out delegation.crt
```

Then they need to give you `delegation.crt`, whether it is self-signed or signed
by a CA.

## Adding a delegation key to an existing repository

If your repository was created using a version of Docker Engine prior to 1.11,
then before adding any delegations, you should rotate the snapshot key to the server
so that collaborators don't need your snapshot key to sign and publish tags:

```
$ notary key rotate docker.io/<username>/<imagename> snapshot -r
```

This tells Notary to rotate a key for your particular image repository. The
`docker.io/` prefix is required. `snapshot -r` specifies that you want
to rotate the snapshot key and that you want the server to manage it (`-r`
stands for "remote").

When adding a delegation, your must acquire
[the PEM-encoded x509 certificate with the public key](#generating-delegation-keys)
of the collaborator you wish to delegate to.

Assuming you have the certificate `delegation.crt`, you can add a delegation
for this user and then publish the delegation change:

```
$ notary delegation add docker.io/<username>/<imagename> targets/releases delegation.crt --all-paths
$ notary publish docker.io/<username>/<imagename>
```

The preceding example illustrates a request to add the delegation
`targets/releases` to the image repository, if it doesn't exist. Be sure to use
`targets/releases` - Notary supports multiple delegation roles, so if you mistype
the delegation name, the Notary CLI does not error. However, Docker Engine
supports reading only from `targets/releases`.

It also adds the collaborator's public key to the delegation, enabling them to sign
the `targets/releases` delegation so long as they have the private key corresponding
to this public key. The `--all-paths` flag tells Notary not to restrict the tag
names that can be signed into `targets/releases`, which we highly recommend for
`targets/releases`.

Publishing the changes tells the server about the changes to the `targets/releases`
delegation.

After publishing, view the delegation information to ensure that you correctly added
the keys to `targets/releases`:

```
$ notary delegation list docker.io/<username>/<imagename>

      ROLE               PATHS                                   KEY IDS                                THRESHOLD
---------------------------------------------------------------------------------------------------------------
  targets/releases   "" <all paths>  729c7094a8210fd1e780e7b17b7bb55c9a28a48b871b07f65d97baf93898523a   1
```

You can see the `targets/releases` with its paths and the key ID you just added.

Notary currently does not map collaborators names to keys, so we recommend
that you add and list delegation keys one at a time, and keep a mapping of the key
IDs to collaborators yourself should you need to remove a collaborator.

## Removing a delegation key from an existing repository

To revoke a collaborator's ability to sign tags for your image repository, you
need to remove their keys from the `targets/releases` delegation. To do this,
you need the IDs of their keys.

```bash
$ notary delegation remove docker.io/<username>/<imagename> targets/releases 729c7094a8210fd1e780e7b17b7bb55c9a28a48b871b07f65d97baf93898523a

Removal of delegation role targets/releases with keys [729c7094a8210fd1e780e7b17b7bb55c9a28a48b871b07f65d97baf93898523a], to repository "docker.io/<username>/<imagename>" staged for next publish.
```

The revocation takes effect as soon as you publish:

```
$ notary publish docker.io/<username>/<imagename>
```

By removing all the keys from the `targets/releases` delegation, the
delegation (and any tags that are signed into it) is removed. That means that
these tags are all deleted, and you may end up with older, legacy tags that
were signed directly by the targets key.

## Removing the `targets/releases` delegation entirely from a repository

If you've decided that delegations aren't for you, you can delete the
`targets/releases` delegation entirely. This also removes all the tags that
are currently in `targets/releases`, however, and you may end up with older,
legacy tags that were signed directly by the targets key.

To delete the `targets/releases` delegation:

```
$ notary delegation remove docker.io/<username>/<imagename> targets/releases

Are you sure you want to remove all data for this delegation? (yes/no)
yes

Forced removal (including all keys and paths) of delegation role targets/releases to repository "docker.io/<username>/<imagename>" staged for next publish.

$ notary publish docker.io/<username>/<imagename>
```

## Pushing trusted data as a collaborator

As a collaborator with a private key that has been added to a repository's
`targets/releases` delegation, you need to import the private key that you
generated into Content Trust.

To do so, you can run:

```
$ notary key import delegation.key --role user
```

where `delegation.key` is the file containing your PEM-encoded private key.

After you have done so, running `docker push` on any repository that
includes your key in the `targets/releases` delegation automatically signs
tags using this imported key.

## `docker push` behavior

When running `docker push` with Docker Content Trust, Docker Engine
attempts to sign and push with the `targets/releases` delegation if it exists.
If it does not, the targets key is used to sign the tag, if the key is available.

## `docker pull` and `docker build` behavior

When running `docker pull` or `docker build` with Docker Content Trust, Docker
Engine pulls tags only signed by the `targets/releases` delegation role or
the legacy tags that were signed directly with the `targets` key.

## Related information

* [Content trust in Docker](content_trust.md)
* [Manage keys for content trust](trust_key_mng.md)
* [Automation with content trust](trust_automation.md)
* [Play in a content trust sandbox](trust_sandbox.md)
