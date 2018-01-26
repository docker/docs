---
description: Configure a Docker Universal Plane cluster to only allow running applications
  that use images you trust.
keywords: docker, ucp, backup, restore, recovery
title: Run only the images you trust
---

With Docker Universal Control Plane you can enforce applications to only use
Docker images signed by UCP users you trust. When a user tries to deploy an
application to the cluster, UCP checks if the application uses a Docker image
that is not trusted, and doesn't continue with the deployment if that's the case.

By signing and verifying the Docker images, you ensure that the images being
used in your cluster are the ones you trust and haven't been altered either
in the image registry or on their way from the image registry to your UCP
cluster.

## Configure UCP

To configure UCP to only allow running applications that use Docker images you
trust, go to the **UCP web UI**, navigate to the **Settings** page, and click
the **Content Trust** tab.

<!-- todo: add screenshot -->

Select the **Run only signed images** option to only allow deploying
applications if they use images you trust.

Then, in the **Require signature from** field, you can specify all the teams
that need sign the image, before it is trusted to run in the UCP cluster. If
you specify multiple teams, the image needs to be signed by a member of each
team, or someone that is a member of all those teams.
If you don't specify any team, the image will be trusted as long as it is signed
by any UCP user whose keys are trusted in a Notary delegation role.

## Set up the Docker Notary CLI client

After you configure UCP to only run applications that use Docker images you
trust, you need to specify which Docker images can be trusted using the Docker
Notary server that is built into Docker Trusted Registry.
You configure the Notary server to store signed metadata about the Docker
images you trust.

To interact with the Notary server, you need to
[install the Notary CLI client](https://github.com/docker/notary/releases).

Once you've installed the Notary client, you need to configure it to talk to
the Notary server that is built into Docker Trusted Registry. This can be done
using a [Notary configuration file](/notary/reference/client-config.md)
or by running:

```bash
# Create an alias to always have the notary client talking to the right server
$ alias notary="notary -s https://<dtr_url> -d ~/.docker/trust"
```

Where `-s` specifies the notary server to talk to, and `-d` specifies the
directory to store private keys and cache data.

If your Docker Trusted Registry is not using certificates signed by a globally
trusted certificate authority, you also need to configure notary to use the
certificate of the DTR CA:

```bash
$ alias notary="notary -s https://<dtr_url> -d ~/.docker/trust --tlscacert <dtr_ca.pem>"
```

## Set up a trusted image repository

Once your Docker Notary CLI client is configured, you can check if Notary has
information about a specific repository:

```bash
# <dtr_url>/<account>/<repository> is also known as a Globally Unique Name (GUN)
$ notary list <dtr_url>/<account>/<repository>
```

If notary has information about the repository it returns the list of
image tags it knows about, their expected digests, and the role of the private
key used to sign the metadata.

If Notary doesn't know yet about an image repository, run:

```bash
$ notary init -p <dtr_url>/<account>/<repository>
```

The Notary CLI client generates public and private key pairs, prompts you for
a passphrase to encrypt the private key, and stores the key pair in the
directory you've specified with the `notary -d` flag.
You should ensure you create backups for these keys, and that they are kept
securely and offline.
[Learn more about the keys used by Docker Notary.](/engine/security/trust/trust_key_mng.md)

## Sign and push an image

Now you can sign your images before pushing them to Docker Trusted Registry:

```bash
# Setting Docker content trust makes the Docker CLI client sign images before pushing them
$ export DOCKER_CONTENT_TRUST=1
# Push the image
$ docker push <dtr_url>/<account>/<repository>:<tag>
```

The Docker CLI client will prompt you for the passphrase you used to encrypt the
private keys, sign the image, and push it to the registry.


## Delegate image signing

Instead of signing the Docker images yourself, you can delegate that task
to other users.

Delegation roles simplify collaborator workflows in Notary trusted collections,
and also allow for fine-grained permissions within a collection's contents
across delegations.
Delegation roles act as signers in Notary that are managed by the targets key
and can be configured to use external signing keys. Keys can be dynamically
added to or removed from delegation roles as collaborators join and leave
trusted repositories.
[Learn more about Notary delegation roles.](/notary/advanced_usage.md)

Every change to the repository now needs to be signed by the snapshot key that
was generated with the `notary init` command.
To avoid having to distribute this key to other members so that they can also
sign images with this key, you can rotate the key and make it be managed by
the Notary server.

This operation only needs to be done once for the repository.

```bash
# This only needs to be done once for the repository
$ notary key rotate <dtr_url>/<account>/<repository> snapshot --server-managed
```

Then ask the users you want to delegate the image signing to share with you
the `cert.pem` files that are included in their client bundles. These files
should be shared using a trusted channel.

Then run the following command to create a new Notary delegation role, using the
user certificates:

```bash
$ notary delegation add -p <dtr_url>/<account>/<repository> targets/releases --all-paths user1.pem user2.pem
```

The above command adds the  the `targets/releases` delegation role to a trusted
repository.
This role is treated as an actual release branch for Docker Content Trust,
since `docker pull` commands with trust enabled will pull directly from this
role, if data exists.
All users that can release images should be added to this role.
[Learn more about the targets/releases role](/engine/security/trust/trust_delegation.md).

Notary has no limit on how many delegation roles can exist, so you can add more
delegation roles such as `targets/qa_team` or `targets/security_team` to the
trusted repository.

Valid delegation roles take the form of `targets/<delegation>`, where
`<delegation>` does not include further slashes.

You will need to add the key to at least one delegation in addition to the `targets/releases` delegation in order for UCP to honor the signed content:

```bash
$ notary delegation add -p <dtr_url>/<account>/<repository> targets/devops --all-paths user1.pem user2.pem
```

Before delegation role users can publish signed content with Notary or
Docker Content Trust, they must import the private key associated with the user certificate:

```bash
$ notary key import key.pem
```

## Where to go next

* [Manage trusted repositories](manage-trusted-repositories.md)
* [Get started with Notary](/notary/getting_started.md)
