---
description: Learn how to use the Notary CLI client to manage trusted repositories
keywords: UCP, trust, notary, registry, security
title: Manage trusted repositories
---

Once you install the Notary CLI client, you can use it to manage your signing
keys, authorize other team members to sign images, and rotate the keys if
a private key has been compromised.

When using the Notary CLI client you need to specify where is Notary server
you want to communicate with, and where to store the private keys and cache for
the CLI client.

```bash
# Create an alias to always have the notary client talking to the right server
$ alias notary="notary -s https://<dtr_url> -d ~/.docker/trust"
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

To chance the passphrase used to encrypt one of the keys, run:

```bash
$ notary key passwd <key_id>
```

## Rotate keys

If one of the private keys is compromised you can rotate that key, so that
images that were signed with those keys stop being trusted.

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

* [Run only the images you trust](index.md)
* [Get started with Notary](/notary/getting_started.md)
