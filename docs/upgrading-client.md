<!--[metadata]>
+++
title = "Upgrading Notary client"
description = "How to upgrade the Notary client"
keywords = ["docker, notary, notary-client"]
[menu.main]
parent="mn_notary"
+++
<![end-metadata]-->

# Upgrading the client

### From `< v0.2` to `v0.2`

Please download the v0.2 binary from [the Notary releases page](https://github.com/docker/notary/releases).

##### New features:

- Support for adding, editing, removing, and viewing information about [delegation roles](https://github.com/theupdateframework/tuf/blob/develop/docs/tuf-spec.txt#L252), and adding targets to delegation roles (as opposed to the targets role, which is currently the default).

	This will allow repository owners to add one or more delegated users as collaborators, allowing them to use their own keys to publish.  In Notary's previous workflow, a repository owner had to share the targets key with each collaborator in order to enable collaborators to publish.

- Support for rotating the snapshot key to one managed by the Notary server.  This is necessary in order to facilitate delegations without forcing repository owners to share the snapshot key with collaborators.

	A collaborator with their own user key may sign a delegation, but without server snapshot-signing, they would be unable to publish unless they also had the snapshot key.

	If repository owners want to reclaim control of their snapshot key, they can rotate the key again to a locally-generated key.

- Support for remote key rotation (forcing server-managed keys to rotate).  Forcing the server to rotate a server-managed key should be a very rare occurrence, necessarily only in the case of server compromise or remote key algorithm obsolescence.

- Drop-in support for user keys - users can just place their private key to be used for delegations in `~/.notary/private/tuf_keys` as `<keyID>.key`.  The PEM header "role" should be "user", although any invalid role is currently allowed.

##### Incompatible changes:

- Due to a bug in previous versions of Notary, if delegation keys are added to a repository, old clients (including, for instance, Docker versions < 1.10, which bundle older versions of Notary client) will not be able to view the repository.  Any attempt to download updated repository information with old client will result in an error.

	However, repositories created or edited by Notary Client v0.2 without delegations will be perfectly readable, downloadable, and publishable by old clients (with the caveat of server snapshot-signing, see below).

- If server snapshot-signing is enabled for a repository, old clients will be able to read from and download updates for the repository, but they may not publish.

- Due to a bug in previous versions of Notary, if new keys are created with Notary Client v0.2 and to disk, an old client (including, for instance, Docker versions < 1.10, which bundle older versions of Notary client) will not be able to use the same trust directory as the new client.  However, the new client will be able to and use read old-style keys.

	If you need to switch between old style and new style clients, please use different trust directories (for example, `~/.notaryv1/` and `~/.notaryv2`).

##### Bugfixes and improvements:

- The client no longer makes an extra request to the server when listing and getting targets.
- The client can use the local cache now if there are problems updating the new timestamp from the server.
- Improved error messages when there is a client failure
- Reduced the verbosity of logging during normal client operation
- The `-v`/`--verbose` command line option now display `ERROR` level logging (normal logging level is `FATAL`)
- Added the `-D`/`--debug` command line option to display `DEBUG` level logging
