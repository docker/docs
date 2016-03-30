<!--[metadata]>
+++
title = "Understand the service architecture"
description = "How the three requisite notary components interact"
keywords = ["docker, notary, notary-client, docker content trust, content trust, notary-server, notary server, notary-signer, notary signer, notary architecture"]
[menu.main]
parent="mn_notary"
weight=3
+++
<![end-metadata]-->


# Understand the Notary service architecture

On this page, you get an overview of the Notary service architecture. This
document assumes a prior understanding of [The Update
Framework](https://theupdateframework.github.io/).

## Architecture and components

Notary clients pull metadata from one or more (remote) Notary services.  Some
Notary clients will push metadata to one or more Notary services.

A Notary service consists of a Notary server, which stores and updates the
signed [TUF metadata files](
https://github.com/theupdateframework/tuf/blob/develop/docs/tuf-spec.txt#L348)
for multiple trusted collections in an associated database, and a Notary signer, which
stores private keys for and signs metadata for the Notary server. The following
diagram illustrates this architecture:

![Notary Service Architecture Diagram](images/service-architecture.svg)

Root, targets, and (sometimes) snapshot metadata are generated and signed by
clients, who upload the metadata to the Notary server. The server is
responsible for:

- ensuring that any uploaded metadata is valid, signed, and self-consistent
- generating the timestamp (and sometimes snapshot) metadata
- storing and serving to clients the latest valid metadata for any trusted collection

The Notary signer is responsible for:

- storing the private signing keys
[wrapped](
https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-31#section-4.4)
and [encrypted](
https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-31#section-4.8)
using [Javascript Object Signing and Encryption](
https://github.com/dvsekhvalnov/jose2go) in a database separate from the
Notary server database
- performing signing operations with these keys whenever the Notary server requests

## Example client-server-signer interaction

The following diagram illustrates the interactions between the Notary client,
sever, and signer:

![Notary Service Sequence Diagram](images/metadata-sequence.svg)

1. Notary server optionally supports authentication from clients using
   [JWT](http://jwt.io/) tokens. This requires an authorization server that
   manages access controls, and a cert bundle from this authorization server
   containing the public key it uses to sign tokens.

   If token authentication is enabled on Notary server, then any connecting
   client that does not have a token will be redirected to the authorization
   server.

   Please see the docs for [Docker Registry v2 authentication](
   https://github.com/docker/distribution/blob/master/docs/spec/auth/token.md)
   for more information.

2. The client will log in to the authorization server via basic auth over HTTPS,
   obtain a bearer token, and then present the token to Notary server on future
   requests.

3. When clients uploads new metadata files, Notary server checks them against
   any previous versions for conflicts, and verifies the signatures, checksums,
   and validity of the uploaded metadata.

4. Once all the uploaded metadata has been validated, Notary server
   generates the timestamp (and maybe snapshot) metadata. It sends this
   generated metadata to the Notary signer to be signed.

5. Notary signer retrieves the necessary encrypted private keys from its database
   if available, decrypts the keys, and uses them to sign the metadata. If
   successful, it sends the signatures back to Notary server.

6. Notary server is the source of truth for the state of a trusted collection of
   data, storing both client-uploaded and server-generated metadata in the TUF
   database. The generated timestamp and snapshot metadata certify that the
   metadata files the client uploaded are the most recent for that trusted collection.

   Finally, Notary server will notify the client that their upload was successful.

7. The client can now immediately download the latest metadata from the server,
   using the still-valid bearer token to connect. Notary server only needs to
   obtain the metadata from the database, since none of the metadata has expired.

   In the case that the timestamp has expired, Notary server would go through
   the entire sequence where it generates a new timestamp, request Notary signer
   for a signature, stores the newly signed timestamp in the database. It then
   sends this new timestamp, along with the rest of the stored metadata, to the
   requesting client.


## Threat model

Both the server and the signer are potential attack vectors. This section
discusses how our architecture is designed to deal with compromises.

### Notary server compromise

In the event of a Notary server compromise, an attacker would have direct access to
the metadata stored in the database as well as well as access to the credentials
used to communicate with Notary signer, and therefore, access to arbitrary signing
operations with any key the Signer holds.

- **Denial of Service** - An attacker could reject client requests and corrupt
    or delete metadata from the database, thus preventing clients from being
    able to download or upload metadata.

- **Malicious Content** - An attacker can create, store, and serve arbitrary
    metadata content for one or more trusted collections. However, they do not have
    access to any client-side keys, such as root, targets, and potentially the
    snapshot keys for the existing trusted collections.

    Only clients who have never seen the trusted collections, and who do not have any
    form of pinned trust, can be tricked into downloading and
    trusting the malicious content for these trusted collections.

    Clients that have previously interacted with any trusted collection, or that have
    their trust pinned to a specific certificate for the collections will immediately
    detect that the content is malicious and would not trust any root, targets,
    or (maybe) snapshot metadata for these collections.

- **Rollback, Freeze, Mix and Match** - The attacker can request that
    the Notary signer sign any arbitrary timestamp (and maybe snapshot) metadata
    they want. Attackers can launch a freeze attack, and, depending on whether
    the snapshot key is available, a mix-and-match attack up to the expiration
    of the targets file.

    Clients both with and without pinned trust would be vulnerable to these
    attacks, so long as the attacker ensures that the version number of their
    malicious metadata is higher than the version number of the most recent
    good metadata that any client may have.

 Note that the timestamp and snapshot keys cannot be compromised in a server-only
 compromise, so a key rotation would not be necessary.  Once the Server
 compromise is mitigated, an attacker will not be
 able to generate valid timestamp or snapshot metadata and serve them on a
 malicious mirror, for example.

### Notary signer compromise

In the event of a Notary signer compromise, an attacker would have access to
all the private keys stored in a database. If the keys are stored in an HSM,
they would have the ability to sign arbitrary content with, and to delete, the
keys in the HSM, but not to exfiltrate the private material.

- **Denial of Service** - An attacker could reject all Notary server requests
  and corrupt or delete keys from the database (or even delete keys from an
  HSM), and thus prevent Notary servers from being able to sign generated
  timestamps or snapshots.

- **Key Compromise** - If the Notary signer uses a database as its backend,
  an attacker can exfiltrate all the private material.  Note that the capabilities
  of an attacker are the same as of a Notary server compromise in terms of
  signing arbitrary metadata, with the important detail that in this particular
  case key rotations will be necessary to recover from the attack.

### Notary client keys and credentials compromise

The severity of the compromise of a trust collection owner/administrator's key
depends on which type of key in the key hierarchy is compromised:

<center><img src="images/key-hierarchy.svg" alt="TUF Key Hierarchy" style="max-width: 500px;"/></center>


Also, the severity depends upon whether a combination of keys were compromised (e.g.
the snapshot key and targets key, or just the targets key).

In general, with the right combination of compromised keys, an attacker would be
able to sign valid changes to the contents of that collection, and to any other
collections that use the same keys.

They can then set up a mirror to distribute this metadata, but they would not
be able to distribute it via a Notary service unless they also have write-capable
credentials for that service (e.g. the username/password of a user who could
push updates into that service).

Familiarity with [TUF (the update framework)](https://theupdateframework.github.io/)
would be helpful in understanding the different types of keys and roles
mentioned below.

Note that the descriptions below assume that the snapshot key is managed by the
Notary service - otherwise, in addition to the keys below, the snapshot key would
also have to be compromised in order to perform any of the attacks.

#### Delegation key compromise

A delegation key has the most limited capabilities of any client-managed key.
It is used to sign targets into [specific delegation roles, which may have path
restrictions](advanced_usage.md#work-with-delegation-roles), and can further
delegate trust to other delegation roles.

- **Limited Malicious Content, Rollback, Freeze, Mix and Match** - An attacker
    can add malicious content, remove legitimate content from a collection, and
    mix up the targets in a collection, but only within the particular delegation
    roles that the key can sign for.  Depending on the restrictions on that role,
    they may be restricted in what type of content they can modify.

- **Limited Denial of Service** - An attacker may add or remove the capabilities
    of other delegation keys with even less capabilities, but only those below it
    on the key hierarchy  (e.g. if `DelegationKey2` were compromised, it would
    only be able to modify the capabilityes of `DelegationKey4` and `DelegationKey5`),
    thus preventing holders of those keys from being able to modify content.

Mitigation:  if a compromise is detected, a higher level key (either the targets
key or another delegation key) holder must rotate the compromised key, and
push a clean set of targets using the new key.

#### Targets key compromise

A targets key, similarly to a delegation key, is used to sign targets into the
targets role, which is the ancestor role of any delegation role.  It delegates
trust to top level delegation roles.

- **Malicious Content, Rollback, Freeze, Mix and Match** - An attacker
    can add any malicious content, remove any legitimate content from a
    collection, and mix up the targets in a collection.

- **Limited Denial of Service** - An attacker may add or remove the capabilities
    of any delegation keys by removing the capabilities of the keys in the top
    level delegation roles, thus preventing holders of those keys from being
    able to modify content.

Mitigation:  if a compromise is detected, the root key holder must rotate the
compromised key and push a clean set of targets using the new key.

#### Root key compromise

A root key is the root of all trust.  It specifies the top keys used to
sign all the other top level metadata (the root, the timestamp, the snapshot,
and the targets keys).

- **Complete Key Compromise** An attacker can rotate all the top level keys,
    including the root key, giving themselves complete control over all keys in
    the repository.

- **Malicious Content, Rollback, Freeze, Mix and Match** - With their
    newly rotated keys, an attacker can add any malicious content, remove any
    legitimate content from a collection, and mix up the targets in a collection.

- **Denial of Service** - By rotating all keys including the root key,
    an attacker removes the capabilities for any other key to sign for new data.

Mitigation:  if a compromise is detected, the root key holder should contact
whomever runs the notary service to manually reverse any malicious changes to
the repository, and immediately rotate the root key.  This will create a fork
of the repository history, and thus break existing clients who have downloaded
any of the malicious changes.

## Related information

* [Run a Notary service](running_a_service.md)
* [Notary configuration files](reference/index.md)
