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
after verifying the root, target, and snapshot signatures on the client update.

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

### How to configure and run notary server

A JSON configuration file is used to configure Notary Server.  Please see the
[Notary Server configuration document](notary-server-config.md)
for more details about the format of the configuration file.

The parameters of the configuration file can also be overwritten using
environment variables of the form `NOTARY_SERVER_var`, where `var` is the
full path from the top level of the configuration file to the variable you want
to override, in all caps.  A change in level is denoted with a `_`.

For instance, one part of the configuration file might look like:

```json
"storage": {
	"backend": "mysql",
	"db_url": "dockercondemo:dockercondemo@tcp(notary-mysql)/dockercondemo"
}
```

If you would like to specify a different `db_url`, the full path from the top
of the configuration tree is `storage -> db_url`, so the environment variable
to set would be `NOTARY_SERVER_STORAGE_DB_URL`.

Note that you cannot override an intermediate level name.  Setting
`NOTARY_SERVER_STORAGE=""` will not disable the MySQL storage.  Each leaf
parameter value must be set indepedently.

#### Running a Docker image

Get the official Docker image, which comes with some sane defaults.  You can
run it with your own signer service and mysql DB, or in the example below, with
just a local signing service and memory store:

```
$ docker pull docker.io/docker/notary-server
$ docker run -p "4443:4443" \
	-e NOTARY_SERVER_TRUST_SERVICE_TYPE=local \
	-e NOTARY_SERVER_STORAGE_BACKEND=""
	-e NOTARY_SERVER_STORAGE_DB_URL=""
	notary-server
```

Alternately, you can run with your own configuration file entirely.  The
docker image loads the config file from `/opt/notary-server/config.json`, so
you can mount your config file at `/opt/notary-server`:

```
$ docker run -p "4443:4443" -v /path/to/your/config/dir:/opt/notary-server
```

#### Running the binary
A JSON configuration file needs to be passed as a parameter/flag when starting
up the Notary Server binary.  Environment variables can also be set in addition
to the configuration file, but the configuration file is required.

```
$ export NOTARY_SERVER_STORAGE_DB_URL=myuser:mypass@tcp(my-db)/dbname
$ NOTARY_SERVER_LOGGING_LEVEL=info notary-server -config /path/to/config.json
```

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
