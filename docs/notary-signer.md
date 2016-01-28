<!--[metadata]>
+++
title = "Signer"
description = "Description of the Notary Signer"
keywords = ["docker, notary, notary-singer"]
[menu.main]
parent="mn_notary"
+++
<![end-metadata]-->

# Notary Signer

The Notary Signer is a remote store for private keys.  It creates and delete
keys, signs data, and returns public key information on demand via its HTTP or
RPC api.

It is intended to be used as a remote RPC service for a
[Notary Server](notary-server.md)'s timestamp private keys.

### Authentication

Notary Signer supports mutual TLS authentication from
[Notary Server](notary-server.md).

Note that when you generate client certificates to be used with Notary Signer,
please make sure that the certificates **are not CAs**.  Otherwise any client
that is compromised can sign any number of other client certs.

As an example, please see [this script](https://github.com/docker/notary/blob/master/docs/opensslCertGen.sh) to see how to
generate client SSL certs with basic constraints using OpenSSL.

### Signer storage

Notary Signer uses MySQL as a backend for storing the encrypted private keys
that is responsible for.  The private keys [wrapped](
https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-31#section-4.4)
and [encrypted](
https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-31#section-4.8)
using [Javascript Object Signing and Encryption](
https://github.com/dvsekhvalnov/jose2go).

The passphrase used to encrypt the keys is passed as an environment variable,
the name of which [is specified by the confguration file](
notary-signer-config.md#storage-section-required).

### How to configure and run Notary Signer

A JSON configuration file is used to configure Notary Signer.  Please see the
[Notary Signer configuration document](notary-signer-config.md)
for more details about the format of the configuration file.

You can also override the parameters of the configuration by
setting environment variables of the form `NOTARY_SIGNER_var`.

`var` is the ALL-CAPS, `"_"`-delimited path of keys from the top level of the
configuration JSON.

For instance, if you wanted to override the storage URL of the Notary Signer
configuration:

```json
"storage": {
	"backend": "mysql",
	"db_url": "dockercondemo:dockercondemo@tcp(notary-mysql)/dockercondemo"
}
```

the full path of keys is `storage -> db_url`. So the environment variable you'd
need to set would be `NOTARY_SIGNER_STORAGE_DB_URL`.

Note that you cannot override a key whose value is another map.
For instance, setting `NOTARY_SIGNER_STORAGE=""` will not disable the
MySQL storage.  You can only override keys whose values are strings or numbers.

For example, if running the binary:

```
$ export NOTARY_SIGNER_STORAGE_DB_URL=myuser:mypass@tcp(my-db)/dbname?parseTime=true
$ NOTARY_SIGNER_LOGGING_LEVEL=info notary-signer -config /path/to/config.json
```

Note that you cannot override a key whose value is another map.
For instance, setting
`NOTARY_SIGNER_STORAGE='{"storage": {"backend": "memory"}}'` will not
set in-memory storage.  It just fails to parse.  You can only override keys
whose values are strings or numbers.

#### Running Notary Signer

Configuration options:

- `-config=<config file>` - The JSON configuration file.

- `-debug` - Passing this flag enables the debugging server on `localhost:8080`.
	The debugging server provides [pprof](https://golang.org/pkg/net/http/pprof/)
	and [expvar](https://golang.org/pkg/expvar/) endpoints.

Get the official Docker image, which comes with [some sane defaults](
https://github.com/docker/notary/blob/master/fixtures/signer-config-local.json),
which uses a local in-memory backing store (not recommended for production).

You can override the default configuration with environment variables.
For example, if you wanted to run it with your own DB
(recommended for production):

```
$ docker pull docker.io/docker/notary-signer
$ docker run -p "4444:4444" \
	-e NOTARY_SIGNER_STORAGE_DB_BACKEND="mysql" \
	-e NOTARY_SIGNER_STORAGE_DB_URL="myuser:mypass@tcp(my-db)/dbName"
	notary-signer
```

Alternately, you can run the image with your own configuration file entirely.
You just need to mount your configuration directory, and then pass the path to
that configuration file as an argument to the `docker run` command:

```
$ docker run -p "4444:4444" \
	-v /path/to/config/dir/on/host:/etc/docker/notary-signer \
	notary-signer -config=/etc/docker/notary-server/config.json
```

You can also pass the `-debug` flag to the container in addition to the
configuration file, but the debug server port is not exposed by the container.
In order to view the debug endpoints, you will have to `docker exec` into
your container.

### What happens if the signer is compromised

All the timestamp private keys stored on the signer will be compromised, and
an attacker can sign anything they wish with the timestamp key.

However, the attacker cannot do anything useful with the timestamp keys unless
they also [compromise the Notary Server](
notary-server.md#what-happens-if-the-server-is-compromised)

The attacker can prevent Notary Signer from signing timestamp metadata from
Notary Server and return invalid public key IDs when the Notary Server
requests it.  This means an attacker can execute a denial of service attack
that prevents the Notary Server from being able to update any metadata.

### Ops features

Notary Signer provides the following features for operational friendliness:

1. A [Bugsnag](https://bugsnag.com) hook for error logs, if a Bugsnag
	configuration is provided.
