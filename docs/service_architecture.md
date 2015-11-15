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

A Notary service consists of a Notary server, which stores and updates the
signed [TUF metadata files](
https://github.com/theupdateframework/tuf/blob/develop/docs/tuf-spec.txt#L348)
for multiple repositories in an associated database, and a Notary signer, which
stores private keys for and signs metadata for the Notary server. The following
diagram illustrates this architecture:

![Notary Service Architecture Diagram](images/service-architecture.svg)

Root, targets, and (sometimes) snapshot metadata are generated and signed by
clients, who upload the metadata to the Notary server. The server is
responsible for:

- ensuring that any uploaded metadata is valid, signed, and self-consistent
- generating the timestamp (and sometimes snapshot) metadata
- storing and serving to clients the latest valid metadata for any repository

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

5. Notary signer gets the requisite encrypted private keys from its database if
   it has them, decrypts the keys, and uses them to sign the metadata. If
   successful, it sends the signatures back to Notary server.

6. Notary server stores all the signed metadata, both client-uploaded and
   server-generated metadata, in the TUF database. The generated timestamp and
   snapshot metadata certify that the metadata files the client uploaded are the
   most recent for that repository.

   Finally, Notary server will notify the client that their upload was successful.

7. The client can now immediately download the latest metadata from the server,
   using the still-valid bearer token to connect. Notary server only needs to
   obtain the metadata from the database, since none of the metadata is expired
   yet.

   If the timestamp were expired, for example, Notary server would go through
   the entire sequence where it generates a new timestamp, request Notary signer
   for a signature, stores the newly signed timestamp in the database. It then
   sends this new timestamp, along with the rest of the stored metadata, to the
   requesting client.


## Threat model

Both the server and the signer are possible attack vectors. This section
discusses how the architecture responds when compromised.

## Notary server compromise

In the event of a Notary server compromise, an attacker would have access to
the metadata stored in the database as well as access to Notary signer to
 sign anything with any key the Signer holds.

- **Denial of Service** - An attacker could reject client requests and corrupt
 	or delete metadata from the database, thus preventing clients from being
    able to download or upload metadata.

- **Malicious Content** - An attacker can create, store, and serve arbitrary
    metadata content for one or more repositories. However, they do not have
    access to the original root, target, or (maybe) snapshots keys for
    existing repositories.

    Only clients who have never seen, and do not have any form of pinned trust
    for, the compromised repositories can be tricked into downloading and
    trusting the malicious content for these repositories.

    Clients who have pinned trust for the compromised repositories, either
    due to configuration or due to TOFU (trust on first use), will immediately
    detect that the content is malicious and would not trust any root, targets,
    or (maybe) snapshot metadata for these repositories.

 - **Rollback, Mix and Match** - The attacker can request that
    the Notary signer sign whatever timestamp (and maybe snapshot) metadata
    they want. They can create valid timestamp and snapshot metadata that
    certifies that the latest version of the repository contains
    strictly older, or even a mix of older and newer, data.

    Clients both with and without pinned trust would be vulnerable to these
    attacks.

 Note that the timestamp and snapshot keys are never compromised; once the
 Server compromise is mitigated, an attacker will not be able to generate
 valid timestamp or snapshot metadata and serve them on a malicious mirror, for
 example.

### Notary signer compromise

In the event of a Notary signer compromise, an attacker would have access to
all the private keys stored in a database. If the keys are stored in an HSM,
they would have the ability to interfere with the keys in the HSM, but not
to exfiltrate the private keys.

- **Denial of Service** - An attacker could reject all Notary server requests
  and corrupt or delete keys from the database (or even delete keys from an
  HSM), and thus prevent Notary servers from being able to sign generated
  timestamps or snapshots.

- **Key Compromise** - If the Notary signer uses a database as its backend,
  an attacker can exfiltrate all the private keys. This will let them set
  up a malicious mirror to perform rollback and mix and match attacks,
  for instance.

## Related information

* [Run a Notary service](running_a_service.md)
* [Notary configuration files](reference/index.md)
