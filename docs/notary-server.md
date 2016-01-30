<!--[metadata]>
+++
title = "Server Architecture"
description = "Description of the Notary Server"
keywords = ["docker, notary, notary-server"]
[menu.main]
parent="mn_notary"
+++
<![end-metadata]-->

# Notary Server Architecture

The Notary Server stores and updates the signed
[TUF metadata files](
https://github.com/theupdateframework/tuf/blob/develop/docs/tuf-spec.txt#L348)
for a repository.  The root, snapshot, and targets metadata files are generated
and signed by clients, and the timestamp metadata file is generated and signed
by the server.

The server creates and stores timestamp keys for each repository (preferably
using a remote key storage/signing service such as
[Notary Signer](notary-signer.md)).

When clients upload metadata files, the server checks them for conflicts and
verifies the signatures and key references in the files. If everything
checks out, the server then signs the timestamp metadata file for the
repository, which certifies that the files the client uploaded are the most
recent for that repository.

### Authentication

Notary Server supports authentication from clients using [JWT](http://jwt.io/)
tokens.  This requires an authorization server that manages access controls,
and a cert bundle from this authorization server containing the public key it
uses to sign tokens.

If token authentication is enabled on Notary Server, then any client that
does not have a token will be redirected to the authorization server.
The client will log in, obtain a token, and then present the token to
Notary Server on future requests.

Notary Server should be configured to trust signatures from that authorization
server.

Please see the docs for [Docker Registry v2 authentication](
https://github.com/docker/distribution/blob/master/docs/spec/auth/token.md)
for more information.

### Server storage

Notary Server uses MySQL as a backend for storing the timestamp
public keys and the TUF metadata for each repository.  It relies on a signing
service to store the private keys.

### Signing service

We recommend deploying Notary Server with a separate, remote signing
service: [Notary Signer](notary-signer.md).  This signing service generates
and stores the timestamp private keys and performs signing for the server.

By using remote a signing service, the private keys would never need to be
stored on the server itself.

Notary Signer supports mutual authentication - when you generate client
certificates for your deployment of Notary Server, please make
sure that the certificates **are not CAs**.  Otherwise if the server is
compromised, it can sign any number of other client certs.

As an example, please see [this script](https://github.com/docker/notary/blob/master/docs/opensslCertGen.sh) to see how to
generate client SSL certs with basic constraints using OpenSSL.

### How to configure and run Notary Server

A JSON configuration file is used to configure Notary Server.  Please see the
[Notary Server configuration document](notary-server-config.md)
for more details about the format of the configuration file.

You can also override the parameters of the configuration by
setting environment variables of the form `NOTARY_SERVER_var`.

`var` is the ALL-CAPS, `"_"`-delimited path of keys from the top level of the
configuration JSON.

For instance, if you wanted to override the storage URL of the Notary Server
configuration:

```json
"storage": {
	"backend": "mysql",
	"db_url": "dockercondemo:dockercondemo@tcp(notary-mysql)/dockercondemo"
}
```

the full path of keys is `storage -> db_url`. So the environment variable you'd
need to set would be `NOTARY_SERVER_STORAGE_DB_URL`.

For example, if running the binary:

```
$ export NOTARY_SERVER_STORAGE_DB_URL=myuser:mypass@tcp(my-db)/dbname?parseTime=true
$ NOTARY_SERVER_LOGGING_LEVEL=info notary-server -config /path/to/config.json
```

Note that you cannot override a key whose value is another map.
For instance, setting
`NOTARY_SERVER_STORAGE='{"storage": {"backend": "memory"}}'` will not
set in-memory storage.  It just fails to parse.  You can only override keys
whose values are strings or numbers.

#### Running Notary Server

Configuration options:

- `-config=<config file>` - The JSON configuration file.

- `-debug` - Passing this flag enables the debugging server on `localhost:8080`.
	The debugging server provides [pprof](https://golang.org/pkg/net/http/pprof/)
	and [expvar](https://golang.org/pkg/expvar/) endpoints.

- `-logf=<format>` - This flag sets the output format for the logs. Possible 
    formats are "json" and "logfmt"


Get the official Docker image, which comes with [some sane defaults](
https://github.com/docker/notary/blob/master/fixtures/server-config-local.json),
which include a remote trust service but local in-memory backing store.

You can override the default configuration with environment variables.
For example, if you wanted to run it with just a local signing service instead
(not recommended for production):

```
$ docker pull docker.io/docker/notary-server
$ docker run -p "4443:4443" \
	-e NOTARY_SERVER_TRUST_SERVICE_TYPE=local
	notary-server
```

Alternately, you can run the image with your own configuration file entirely.
You just need to mount your configuration directory, and then pass the path to
that configuration file as an argument to the `docker run` command:

```
$ docker run -p "4443:4443" \
	-v /path/to/config/dir/on/host:/etc/docker/notary-server/ \
	notary-server -config=/etc/docker/notary-server/config.json
```

You can also pass the `-debug` flag to the container in addition to the
configuration file, but the debug server port is not exposed by the container.
In order to view the debug endpoints, you will have to `docker exec` into
your container.

### What happens if the server is compromised

The server does not hold any keys for repositories, except the for timestamp
keys if you are using a local signing service, so the attacker cannot modify
the root, targets, or snapshots metadata.

If you are using a signer service, an attacker cannot get access to the
timestamp key either. They can use the server's credentials to get the signer
service to sign arbitrary data, such as an empty timestamp,
an invalid timestamp, or an old timestamp.

However, TOFU (trust on first use) would prevent the attacker from tricking
existing clients for existing repositories to download arbitrary data.
They would need the original root/target/snapshots keys to do that. The
attacker could only, by signing bad timestamps, prevent the such a user from
seeing any updated metadata.

The attacker can also make all new keys, and simply replace the repository
metadata with metadata signed with these new keys.  New clients who have not
seen the repository before will trust this bad data, but older clients will
know that something is wrong.

### Ops features

Notary Server provides the following features for operational friendliness:

1. A health endpoint at `/_notary_server/health` which returns 200 and a
	body of `{}` if the server is healthy, and a 500 with a map of
	failed services if the server cannot access its storage backend.

	If it cannot contact the signing service, an error will be logged but the
	service will still be considered healthy, because it can still serve
	existing metadata.  It cannot accept updates, so the service is degraded.

1. A [Bugsnag](https://bugsnag.com) hook for error logs, if a Bugsnag
	configuration is provided.

1. A [prometheus](http://prometheus.io/) endpoint at `/_notary_server/metrics`
	which provides HTTP stats.
