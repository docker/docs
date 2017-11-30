---
description: Set up and configure content trust and signing policy for use with a continuous integration system
keywords: cup, trust, notary, security, continuous integration
title: Use trusted images for continuous integration
---

The document provides a minimal example on setting up Docker Content Trust (DCT) in
Universal Control Plane (UCP) for use with a Continuous Integration (CI) system. It
covers setting up the necessary accounts and trust delegations to restrict only those
images built by your CI system to be deployed to your UCP managed cluster.

## Set up UCP accounts and teams

The first step is to create a user account for your CI system. For the purposes of
this document we will assume you are using Jenkins as your CI system and will therefore
name the account "jenkins". As an admin user logged in to UCP, navigate to "User Management"
and select "Add User". Create a user with the name "jenkins" and set a strong password.

Next, create a team called "CI" and add the "jenkins" user to this team. All signing
policy is team based, so if we want only a single user to be able to sign images
destined to be deployed on the cluster, we must create a team for this one user.

## Set up the signing policy

While still logged in as an admin, navigate to "Admin Settings" and select the "Content Trust"
subsection. Select the checkbox to enable content trust and in the select box that appears,
select the "CI" team we have just created. Save the settings.

This policy will require that every image that referenced in a `docker image pull`,
`docker container run`, or `docker service create` must be signed by a key corresponding
to a member of the "CI" team. In this case, the only member is the "jenkins" user.

## Create keys for the Jenkins user

The signing policy implementation uses the certificates issued in user client bundles
to connect a signature to a user. Using an incognito browser window (or otherwise),
log in to the "jenkins" user account you created earlier. Download a client bundle for
this user. It is also recommended to change the description associated with the public
key stored in UCP such that you can identify in the future which key is being used for
signing.

Each time a user retrieves a new client bundle, a new keypair is generated. It is therefore
necessary to keep track of a specific bundle that a user chooses to designate as their signing bundle.

Once you have decompressed the client bundle, the only two files you need for the purposes
of signing are `cert.pem` and `key.pem`. These represent the public and private parts of
the user's signing identity respectively. We will load the `key.pem` file onto the Jenkins
servers, and use `cert.pem` to create delegations for the "jenkins" user in our
Trusted Collection.

## Prepare the Jenkins server

### Load `key.pem` on Jenkins

You will need to use the notary client to load keys onto your Jenkins server. Simply run
`notary -d /path/to/.docker/trust key import /path/to/key.pem`. You will be asked to set
a password to encrypt the key on disk. For automated signing, this password can be configured
into the environment under the variable name `DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE`. The `-d`
flag to the command specifies the path to the `trust` subdirectory within the server's `docker`
configuration directory. Typically this is found at `~/.docker/trust`.

### Enable content trust

There are two ways to enable content trust: globally, and per operation. To enabled content
trust globally, set the environment variable `DOCKER_CONTENT_TRUST=1`. To enable on a per
operation basis, wherever you run `docker image push` in your Jenkins scripts, add the flag
`--disable-content-trust=false`. You may wish to use this second option if you only want
to sign some images.

The Jenkins server is now prepared to sign images, but we need to create delegations referencing
the key to give it the necessary permissions.

## Initialize a repository

Any commands displayed in this section should _not_ be run from the Jenkins server. You
will most likely want to run them from your local system.

If this is a new repository, create it in Docker Trusted Registry (DTR) or Docker Hub,
depending on which you use to store your images, before proceeding further.

We will now initialize the trust data and create the delegation that provides the Jenkins
key with permissions to sign content. The following commands initialize the trust data and
rotate snapshotting responsibilities to the server. This is necessary to ensure human involvement
is not required to publish new content.

```
notary -s https://my_notary_server.com -d ~/.docker/trust init my_repository
notary -s https://my_notary_server.com -d ~/.docker/trust key rotate my_repository snapshot -r
notary -s https://my_notary_server.com -d ~/.docker/trust publish my_repository
```

The `-s` flag specifies the server hosting a notary service. If you are operating against
Docker Hub, this will be `https://notary.docker.io`. If you are operating against your own DTR
instance, this will be the same hostname you use in image names when running docker commands preceded
by the `https://` scheme. For example, if you would run `docker image push my_dtr:4443/me/an_image` the value
of the `-s` flag would be expected to be `https://my_dtr:4443`.

If you are using DTR, the name of the repository should be identical to the full name you use
in a `docker image push` command. If however you use Docker Hub, the name you use in a `docker image push`
must be preceded by `docker.io/`. i.e. if you ran `docker image push me/alpine`, you would
`notary init docker.io/me/alpine`.

For brevity, we will exclude the `-s` and `-d` flags from subsequent command, but be aware you
will still need to provide them for the commands to work correctly.

Now that the repository is initialized, we need to create the delegations for Jenkins. Docker
Content Trust treats a delegation role called `targets/releases` specially. It considers this
delegation to contain the canonical list of published images for the repository. It is therefore
generally desirable to add all users to this delegation with the following command:

```
notary delegation add my_repository targets/releases --all-paths /path/to/cert.pem
```

This solves a number of prioritization problems that would result from needing to determine
which delegation should ultimately be trusted for a specific image. However, because it
is anticipated that any user will be able to sign the `targets/releases` role it is not trusted
in determining if a signing policy has been met. Therefore it is also necessary to create a
delegation specifically for Jenkins:

```
notary delegation add my_repository targets/jenkins --all-paths /path/to/cert.pem
```

We will then publish both these updates (remember to add the correct `-s` and `-d` flags):

```
notary publish my_repository
```

Informational (Advanced): If we included the `targets/releases` role in determining if a signing policy
had been met, we would run into the situation of images being opportunistically deployed when
an appropriate user signs. In the scenario we have described so far, only images signed by
the "CI" team (containing only the "jenkins" user) should be deployable. If a user "Moby" could
also sign images but was not part of the "CI" team, they might sign and publish a new `targets/releases`
that contained their image. UCP would refuse to deploy this image because it was not signed
by the "CI" team. However, the next time Jenkins published an image, it would update and sign
the `targets/releases` role as whole, enabling "Moby" to deploy their image.

## Conclusion

With the Trusted Collection initialized, and delegations created, the Jenkins server will
now use the key we imported to sign any images we push to this repository.

Through either the Docker CLI, or the UCP browser interface, we will find that any images
that do not meet our signing policy cannot be used. The signing policy we set up requires
that the "CI" team must have signed any image we attempt to `docker image pull`, `docker container run`,
or `docker service create`, and the only member of that team is the "jenkins" user. This
restricts us to only running images that were published by our Jenkins CI system.
