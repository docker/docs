<!--[metadata]>
+++
title = "Notary Server"
description = "Description of the Notary Server"
keywords = ["docker, notary, notary-server"]
[menu.main]
parent="mn_notary"
+++
<![end-metadata]-->

# Notary Server

The notary server is a remote store for, and coordinates updates to, the signed
metadata files for a repository (which are created by clients).  The server is
also responsible for creating and keeping track of timestamp keys for each repo,
and signing a timestamp file for each repo whenever a client sends updates,
after verifying the root/target/snapshot signatures on the client update.

### Authentication

Notary Server supports authentication from clients using [JWT](http://jwt.io/)
tokens.  This requires an authorization server that manages access controls,
and a cert bundle from this authorization server containing the public key it
uses to sign tokens.

The client will log into the server (it gets redirected by Notary Server if
authentication is configured), obtain a token, and present the token to Notary
Server, which should be configured to trust signatures from that authorization
server.

See the docs for [Docker Registry v2 authentication](
https://github.com/docker/distribution/blob/master/docs/spec/auth/token.md)
for more information.

### Server storage

Currently Notary Server uses MySQL as a backend for storing the timestamp
*public* keys and the TUF metadata for each repository.

### Signing service

The recommended usage of the server is with a separate signing service:
[Notary Signer](notary-signer.md).  The signing service actually
stores the timestamp *private* keys and performs signing for the server.

By using a signing service, the private keys then would never be stored on the
server itself.

Notary Signer supports mutual authentication - when you generate client
certificates for Notary Server to authenticate with Notary Signer, please make
sure that the certificates **are not CAs**.  Otherwise any server that is
compromised can sign any number of other client certs.

As an example, please see [this script](opensslGenCert.sh) to see how to
generate client SSL certs with basic constraints using OpenSSL.

### How to configure notary server

A JSON configuration file needs to be passed as a parameter/flag when starting
up Notary Server:

```
notary-server -config /path/to/configuration.json
```

Please see the [Notary Server configuration document](notary-server-config.md)
for more details about the format of the configuration file.

### What happens if the server is compromised

The server does not hold any keys for the repository except the timestamp key,
so the attacker cannot modify the root, targets, or snapshots metadata.

If using a signer service, an attacker cannot get access to the timestamp key.
They can use the server to make calls to the signer service to sign arbitrary
data, such as an empty timestamp, an invalid timestamp, or an old timestamp.

TOFU (trust on first use) would prevent the attacker from being able to make
existing clients for existing repositories download arbitrary data.  They would
need the original root/target/snapshots keys.  The attacker could, by signing
bad timestamps, prevent the user from seeing any updated metadata.

The attacker can also make all new keys, and simply replace the repository
metadata with metadata signed with these new keys.  New clients who have not
seen this repository before will trust this bad data, but older clients will
know that something is wrong.

### Ops features

Notary server provides the following endpoints for ops friendliness:

1. A health endpoint at `/_notary_server/health` which returns 200 and a
	body of `{}` if the server is healthy, and a 500 with a list of
	failed services if the server cannot access its storage backend.

	If it cannot contact the signing service (in which case service is degraded,
	but not down, since the server can still serve metadata, but not accept
	updates), an error will be logged, but the service will still be considered
	healthy.

1. A [Bugsnag](https://bugsnag.com) hook for error logs, if a Bugsnag
	configuration is provided.

1. A [prometheus](http://prometheus.io/) endpoint at `/_notary_server/metrics`
	which provides HTTP stats.
