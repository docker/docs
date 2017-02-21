---
title: Configure the UCP and DTR servers for content trust
description: Configuration tasks for using content trust on UCP and DTR servers
---

These tasks allow an administrator to set up the UCP and DTR servers to require
content trust and to delegate the ability to sign images to UCP users. For an
overview of content trust in UCP, see [Run only the images you trust](index.md).

After completing these steps, continue to
[Client configuration for content trust in UCP](client_configuration.md).

## Prerequisites

Before completing these tasks, set up the teams for signing and add users to
them. For instance, if your business requirement is that images need to be
signed by `engineering`, `security`, and `quality` teams, set up those teams and
add the appropriate users to them. See
[Set up teams](/datacenter/ucp/2.1/guides/admin/manage-users/create-and-manage-teams.md)
and
[Create and manage users](/ucp/2.1/guides/admin/manage-users/create-and-manage-users.md).

## Overview

The administrator needs to complete the following tasks to configure UCP and DTR
for image signing.

2.  [Configure UCP](#configure-ucp) to only allow signed images to be used.

3.  [Set up the Docker Notary CLI client](#set-up-the-docker-notary-client)
    locally so that the administrator can initialize the trusted image
    repository.

4.  [Initialize the trusted image repository](#initialize-the-trusted-image-repository)
    which will store trusted images.

5.  [Delegate image signing](#delegate-image-signing) so that the appropriate
    users are able to sign images. This step is optional, and assumes that the
    UCP administrator will not be the only one signing images.

## Configure UCP

This step configures UCP to only allow deployment of signed images, as well as
the teams that must sign an image before it can be trusted. Set up the users
and groups before starting this step.

1.  Go to the **UCP web UI**, navigate to the **Admin Settings** page, and click
    the **Content Trust** menu item.

    ![Content Trust settings](/datacenter/images/ucp_content_trust_settings.png)

2.  Select the **Only run signed images** option.

    Click the **REQUIRE SIGNATURE FROM ALL OF THESE TEAMS** field and choose
    one or more teams. DDC will consider an image to be trustworthy only if it
    is signed by a member of _every team_ you select. The requirement can be
    fulfilled by a single user who is a member of all the teams, or by a signer
    in each team.

    > **Note**: If you don't specify any team (by leaving the field blank), an
    > image will be trusted as long as it is signed by any UCP user whose keys
    > are [configured in the Notary client](#set-up-the-docker-notary-client).

    The following screenshot shows a configuration that requires images to be
    signed by a member of the `engineering` team.

    ![Content Trust settings detail](/datacenter/images/ucp_content_trust_detail.png)

3.  Click **Update** to apply the changes.

UCP is now configured to only allow use of signed images, but you don't have the
ability to sign images yet. Next,
[set up the Docker Notary CLI client](#set-up-the-docker-notary-client).

## Set up the Docker Notary CLI client

After [configuring UCP](#configure-ucp), you need to specify which Docker images
can be trusted, using the Docker Notary server that is built into Docker Trusted
Registry (DTR). The following procedure configures the Notary server to store
signed metadata about the Docker images you trust. This set-up step only needs
to be done on the client of an administrator responsible for setting up
repositories and delegating the ability to sign images.

1.  **If you are on a Linux client**, install the Notary binary. If you use Docker
    for Mac or Docker for Windows, the Notary client is included in your
    installation.

    - [Download the latest client archive](https://github.com/docker/notary/releases).
    - Extract the archive.
    - Rename the binary to `notary` and set it to executable. Either move it to
      a location in your path or modify the examples below to include the full
      path to the binary.

2.  Configure the Notary client to communicate with the DTR server and store
    its metadata in the correct location. You can either use a
    [Notary configuration file](/notary/reference/client-config.md) or manually
    specify the following flags when you run the `notary` command.

    |Flag               |Purpose                                               |
    |-------------------|------------------------------------------------------|
    | `-s <dtr_url>`    | The hostname or IP address of the DTR server         |
    | `-d <trust_directory>` | The path to the local directory where trust metadata will be stored |
    | `--tlscacert <dtr_ca.pem>` | The path to the trust certificate for DTR. Only required if your DTR registry is not using certificates signed by a globally trusted certificate authority, such as self-signed certificates. Download the trust certificate from `https://<dtr_url>/ca` either from your browser or using `curl` or `wget`.|

    > **Tip**: If you don't want to provide the `-s`, `-d`, and `--tlscacert`
    > parameters each time you run a Nptary command, you can set up an alias in
    > Bash (Linux or macOS) or PowerShell (Windows) to save some typing. The
    > following examples do not include the `--tlscacert` flag, but you can add
    > it if necessary. All of the `notary` commands in the rest of this topic
    > assume that you have set up the alias.
    >
    >  - **Bash**: Type the following, or add it to your `~/.profile` file to make
    >    it permanent. Replace `<dtr_url>` with the hostname or IP address of your
    >    DTR instance.
    >
    >        alias notary="notary -s https://<dtr_url> -d ~/.docker/trust"
    >
    >  - **PowerShell**:  Type the following, or add it to your `profile.ps1` to
    >    make it permanent.
    >
    >        PS C:\> set-alias notary "notary -s https://<dtr_url> -d ~/.docker/trust"
    >
    >
    >  After setting up the alias, you only need to type `notary` and the server
    >  and destination directory will be included in the command automatically.
    {: id="notary_alias_config_note" }

3.  Find the _globally unique name (GUN)_ for your repository. The GUN
    is the `<registry/<account>/<repository>` string, such as
    `dtr-example.com/engineering/my-repo`. You can find the GUN for a repository
    by browsing to it in the DTR web UI and copying the part of the
    **Pull command** after `docker pull`.

Next, [initialize a trusted image repository](#set-up-a-trusted-image-repository).

## Initialize the trusted image repository

> **Tip - Yubikey integration**: Notary supports integration with Yubikey. If
> you have a Yubikey plugged in when you initialize a repository with Notary,
> the root key is stored on the Yubikey instead of in the trust directory. When
> you run any command that needs the `root` key, Notary looks on the Yubikey
> first, and uses the trust directory as a fallback.

This procedure needs to be done on the client of an administrator responsible for
setting up repositories. It needs to be done once per signing repository.

In these examples, add the `-s`, `-d`, and `--tlscacert` parameters you need
before `<GUN>` if you decided not to configure Notary using a
[configuration file](/notary/reference/client-config.md){: target="_blank" class="_" }
or a [terminal alias](#notary_alias_config_note").

1.  In the DTR web UI, create a new repository or browse to an existing
    repository that you want to reconfigure as a trusted image repository.
    Make a note of the GUN for the repository by copying the contents of the
    **Pull command** field after `docker pull`.

    ![Finding the GUN in the DTR UI](/datacenter/images/dtr_repo_find_gun.png)

2.  At the command line on your client, check whether Notary has information
    about the repository. Most likely if you are performing this task for the
    first time, the repository is not initialized.

    ```bash
    $ notary list <GUN>
    ```

    The response may be one of the following:

    - `fatal: client is offline`: Either the repository server can't be reached,
      or DTR is using certificates which are not signed by a globally trusted
      certificate authority, such as self-signed certificates. Run `notary list`
      again, adding the `--tlscacert` flag, with the path to the certificate
      authority for DTR. To get the certificate, download `https://<dtr_url>/ca`
      from your browser or using `curl` or `wget`. This certificate is
      different from the UCP trust certificate in the UCP client bundle.

    - `fatal: <dtr_url> does not have trust data for ddc-staging-dtr.qa.aws.dckr.io/engineering/redis`:
      The repository has not yet been initialized and you need to run
      `notary init`. Continue to step 2.

    - `No targets present in this repository.`: The repository has been
      initialized, but contains no signed images. You do not need to do step 2.

    - A list of signed image tags, their digests, and the role of the private
      key used to sign the metadata. This indicates that the repository is configured
      correctly and images have been signed and uploaded. You do not need to do
      step 2.

2.  To initialize the repository, run `notary init`, setting the `-p` flag to
    the GUN of the repository. You will be prompted to set passphrases for
    three different keys:

    - The `root` key is used to sign the `targets` and `snapshot` keys.
    - The `targets` key will be used to sign the keys of users authorized to sign
      images and designate them as trusted.
    - The `snapshot` key is used for snapshotting the repository, which is an
      optimization for updating the trust data.

    ```bash
    $ notary init <GUN>

    No root keys found. Generating a new root key...
    You are about to create a new root signing key passphrase. This passphrase
    will be used to protect the most sensitive key in your signing system. Please
    choose a long, complex passphrase and be careful to keep the password and the
    key file itself secure and backed up. It is highly recommended that you use a
    password manager to generate the passphrase and keep it safe. There will be no
    way to recover this key. You can find the key in your config directory.

    Enter passphrase for new root key with ID 717fa4b:
    Repeat passphrase for new root key with ID 717fa4b:
    Enter passphrase for new targets key with ID 776d924 (<GUN>):
    Repeat passphrase for new targets key with ID 776d924 (<GUN>):
    Enter passphrase for new snapshot key with ID d3cc399 (<GUN>):
    Repeat passphrase for new snapshot key with ID d3cc399 (<GUN>):
    Enter username: admin
    Enter password:
    ```

    As the help text in the command says, it's important to choose good
    passphrases and to save them in a secure location such as a password
    manager. The final username and password prompt are for the DTR login.

    Several important files are saved in the trust directory (the location you
    specified as the value of the `-d` flag. The following is an example listing
    of the trust directory:

    ```none
    ├── private
    │   ├── root_keys
    │   │   └── 92c11d487023de4447ef57747e84e7364cd7c62a4be28d8714ec05afe2f130f8.key
    │   └── tuf_keys
    │       └── dtr-example.com
    │           └── engineering
    │               └── testrepo
    │                   ├── 87a129ea47a4112fec6b989bde35f6ddea8325638450d41b3f44fabaf49dbe3d.key
    │                   └── aa4b236c610e3d951c930bf2a503861c41808ac28369bfbaf5c075b62cb3dd41.key
    └── tuf
        └── dtr-example.com
            └── engineering
                └── testrepo
                    ├── changelist
                    └── metadata
                        ├── root.json
                        ├── snapshot.json
                        └── targets.json
    ```

    The `tuf` directory contains metadata needed by Notary. The `private`
    directory contains the root key, target key, and snapshots key. It is
    important to protect these keys, especially the root key. If you are using
    the Yubikey integration feature, the root key is already stored on your
    Yubikey. You should back up the entire `private` subdirectory to secure
    offline storage and remove the `root_keys` subdirectory from the trust
    directory. If you do not use a Yubikey, back up the entire trust
    directory to secure offline storage, and bring it online only when you need
    to perform Notary operations.

3.  The metadata has been created but only exists on your client. To publish
    it to DTR, use `notary publish`.

    ```bash
    $ notary publish <GUN>

    Pushing changes to <GUN>
    Enter username: admin
    Enter password:
    Enter passphrase for targets key with ID 63c2d66:
    Enter passphrase for snapshot key with ID 6ac388d:
    Successfully published changes for repository <GUN>
    ```

    You will be prompted for the DTR login, the passphrase for the `targets` key,
    and the passphrase for the `snapshot` key.

Typically, the administrator is not part of the group which is authorized to
sign images. If you do attempt to sign images and you are not part of one of the
correct groups, the image will not be available to UCP.

Continue to [delegate image signing](#delegate-image-signing) to give the
appropriate users the ability to sign images.

You can also
[learn more about the keys used by Notary](/engine/security/trust/trust_key_mng.md).


## Delegate image signing

The administrator who manages Docker Trusted Registry is often not part of the
group which is allowed to sign images. This is where
[Notary delegation roles](/notary/advanced_usage.md) come in. Delegation roles
provide:

- Simple collaboration workflows
- Fine-grained permissions within a collection's contents across delegations
- Ability to dynamically add or remove keys from delegation roles when
  collaborators join or leave trusted repositories

When you [initialized the trusted repository](#nitialize-the-trusted-image-repository),
three keys were created:

- The `root` key signs the `targets` and `snapshot` keys.
- The `targets` key is used by Notary for delegation roles, which act as signers.
- Each change in the repository needs to be signed by the `snapshot` key.

To avoid the need to distribute the `snapshot` key to each person who will sign
images, you can configure the Notary server to manage it. In order to do this,
you need to also rotate the `snapshots` key, so that private keys do not need
to be transferred between the client and server.

1.  Rotate the key and configure the Notary server to manage it. This operation
    only needs to be done once for each trusted repository.

    ```bash
    $ notary key rotate <GUN> snapshot --server-managed
    ```

    You are prompted for the DTR credentials followed by the passphrase for the
    `root` key.

2.  For each user who should be able to sign images, ask that user to create a
    client bundle. They should:

    1.  Go to the UCP web UI.
    2.  Click your username at the top right. Click **Profile**.
    3.  Click **Create a Client Bundle**. A file is downloaded called
        `ucp-bundle-<username>.zip`.
    4.  Extract the zip file. The important file within the archive is the
        `cert.pem`, which is the user certificate.
    5.  Send you the `cert.pem` **through a secure, trusted channel**. If you
        plan to create more than one delegation, rename the `cert.pem` with the
        username or other identifying information.

3.  Run the following command to add the `targets/releases` delegation role for
    each user, using the `cert.pem` files. You can specify multiple `cert.pem`
    files at once.

    ```bash
    $ notary delegation add -p <GUN> targets/releases --all-paths user1.pem user2.pem
    ```

    You will be prompted for your DTR credentials and the passphrase for the
    `targets` key.

    > **Note**: You can also add arbitrary delegations, but `targets/releases`
    > is a special delegation, and is treated as an actual release branch for
    > Docker Content Trust. If a Docker client has content trust enabled, and
    > the client runs `docker pull`, this delegation is what signals that the
    > content is trusted.

    Each user who can release images should be added to the `targets/releases`
    role.

4.  Create at least one more delegation and add users to it, or UCP will not
    honor the signed content. This delegation indicates the team that is signing
    the release.

    Docker recommends adding one delegation per team. For instance, if you have
    an `engineering` team and a `qa` team, add a delegation for each of these.
    If a user is a member of both teams, that user will be able to indicate
    which team they are signing on behalf of. Notary has no limit on how many
    delegation roles can exist.

    Valid delegation roles take the form of `targets/<delegation>`. Do not include
    a trailing slash.

    The following command adds `user1` to the `targets/engineering` delegation:

    ```bash
    $ notary delegation add -p <GUN> targets/engineering --all-paths user1.pem
    ```

    You will be prompted for your DTR credentials followed by the passphrase
    for the `targets` key.

5.  Securely remove the `.pem` files of the users you added delegations to. If
    these keys are compromised, they could be used to sign images which should
    not be trusted.


## Next steps

The Notary server is now configured to allow users to sign images. Next, each
user needs to [configure their client](client_configuration.md)
and [sign some images](client_configuration.md#sign-and-push-images).

[Learn more about the targets/releases role](/engine/security/trust/trust_delegation.md).
