---
title: Manage trusted repositories
description: Learn how to use the Notary CLI client to manage trusted repositories
keywords: dtr, trust, notary, security
---

Once you
[configure the Notary CLI client](../../access-dtr/configure-your-notary-client.md),
you can use it to manage your private keys, list trust data from any repository
you have access to, authorize other team members to sign images, and rotate
keys if a private key has been compromised.

## List trust data

List the trust data for a repository by running:

```none
$ notary list <dtr_url>/<account>/<repository>
```

You can get one of the following errors, or a list with the images that have
been signed:

| Message                                     | Description                                                                                                      |
|:--------------------------------------------|:-----------------------------------------------------------------------------------------------------------------|
| `fatal: client is offline`                  | Either the repository server can't be reached, or your Notary CLI client is misconfigured                        |
| `fatal: <dtr_url> does not have trust data` | There's no trust data for the repository. Either run `notary init` or sign and push an image to that repository. |
| `No targets present in this repository`     | The repository has been initialized, but doesn't contain any signed images                                       |

## Initialize trust for a repository

There's two ways to initialize trust data for a repository. You can either
sign and push an image to that repository:

```none
export DOCKER_CONTENT_TRUST=1
docker push <dtr_url>/<account>/<repository>
```

or

```
notary init --publish <dtr_url>/<account>/<repository>
```

## Manage staged changes

The Notary CLI client stages changes before publishing them to the server.
You can manage the changes that are staged by running:

```bash
# Check what changes are staged
$ notary status <dtr_url>/<account>/<repository>

# Unstage a specific change
$ notary status <dtr_url>/<account>/<repository> --unstage 0

# Alternatively, unstage all changes
$ notary status <dtr_url>/<account>/<repository> --reset
```

When you're ready to publish your changes to the Notary server, run:

```bash
$ notary publish <dtr_url>/<account>/<repository>
```

## Delete trust data

Administrator users can remove all signatures from a trusted repository by
running:

```bash
$ notary delete <dtr_url>/<account>/<repository> --remote
```

If you don't include the `--remote` flag, Notary deletes local cached content
but will not delete data from the Notary server.


## Change the passphrase for a key

The Notary CLI client manages the keys used to sign the image metadata. To
list all the keys managed by the Notary CLI client, run:

```bash
$ notary key list
```

To change the passphrase used to encrypt one of the keys, run:

```bash
$ notary key passwd <key_id>
```

## Rotate keys

If one of the private keys is compromised you can rotate that key, so that
images that were signed with the key stop being trusted.

For keys that are kept offline and managed by the Notary CLI client, such the
keys with the root, targets, and snapshot roles, you can rotate them with:

```bash
$ notary key rotate <dtr_url>/<account>/<repository> <key_role>
```

The Notary CLI client generates a new key for the role you specified, and
prompts you for a passphrase to encrypt it.
Then you're prompted for the passphrase for the key you're rotating, and if it
is correct, the Notary CLI client contacts the Notary server to publish the
change.

You can also rotate keys that are stored in the Notary server, such as the keys
with the snapshot or timestamp role. For that, run:

```bash
$ notary key rotate <dtr_url>/<account>/<repository> <key_role> --server-managed
```

## Manage keys for delegation roles

To delegate image signing to other UCP users, get the `cert.pem` file that's
included in their client bundle and run:

```bash
$ notary delegation add -p <dtr_url>/<account>/<repository> targets/<role> --all-paths user1.pem user2.pem
```

You can also remove keys from a delegation role:

```bash
# Remove the given keys from a delegation role
$ notary delegation remove -p <dtr_url>/<account>/<repository> targets/<role> <keyID1> <keyID2>

# Alternatively, you can remove keys from all delegation roles
$ notary delegation purge <dtr_url>/<account>/<repository> --key <keyID1> --key <keyID2>
```

## Troubleshooting

Notary CLI has a `-D` flag that you can use to increase the logging level. You
can use this for troubleshooting.

Usually most problems are fixed by ensuring you're communicating with the
correct Notary server, using the `-s` flag, and that you're using the correct
directory where your private keys are stored, with the `-d` flag.

## Where to go next

* [Learn more about Notary](/notary/advanced_usage.md)
* [Notary architecture](/notary/service_architecture.md)
