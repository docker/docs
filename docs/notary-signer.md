<!--[metadata]>
+++
title = "Notary Signer"
description = "Description of the Notary Signer"
keywords = ["docker, notary, notary-singer"]
[menu.main]
parent="mn_notary"
+++
<![end-metadata]-->

# Notary Signer

The Notary Signer is a remote store for private keys.  It will create and delete
keys, signs data, and return public key information on demand via its HTTP or
RPC api.

It is intended to be used as a remote RPC service for a
[Notary Server](notary-server.md)'s timestamp private keys.

### Authentication

Notary Signer supports mutual TLS authentication from clients (the only client
it supports so far is the [Notary Server](notary-server.md).

Note that when you generate client certificates to be used with Notary Signer,
please make sure that the certificates **are not CAs**.  Otherwise any server
that is compromised can sign any number of other client certs.

As an example, please see [this script](opensslCertGen.sh) to see how to
generate client SSL certs with basic constraints using OpenSSL.

### How to configure and run Notary Signer

A JSON configuration file is used to configure Notary Signer.  Please see the
[Notary Signer configuration document](notary-signer-config.md)
for more details about the format of the configuration file.

The parameters of the configuration file can also be overwritten using
environment variables of the form `NOTARY_SIGNER_var`, where `var` is the
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
to set would be `NOTARY_SIGNER_STORAGE_DB_URL`.

Note that you cannot override an intermediate level name.  Setting
`NOTARY_SIGNER_STORAGE=""` will not disable the MySQL storage.  Each leaf
parameter value must be set indepedently.

#### Running a Docker image

Get the official Docker image, which comes with some sane defaults.  You can
run it with your own MySQL DB:

```
$ docker pull docker.io/docker/notary-signer
$ docker run -p "4443:4443" \
	-e NOTARY_SIGNER_STORAGE_DB_URL="myuse@mypass:tcp(mydb)/dbNname" \
	notary-signer
```

Alternately, you can run with your own configuration file entirely.  The
docker image loads the config file from `/opt/notary-signer/config.json`, so
you can mount your config file at `/opt/notary-signer`:

```
$ docker run -p "4443:4443" -v /path/to/config/dir:/opt/notary-signer notary-signer
```

#### Running the binary
A JSON configuration file needs to be passed as a parameter/flag when starting
up the Notary Signer binary.  Environment variables can also be set in addition
to the configuration file, but the configuration file is required.

```
$ export NOTARY_SIGNER_STORAGE_DB_URL=myuser:mypass@tcp(my-db)/dbname
$ NOTARY_SIGNER_LOGGING_LEVEL=5 notary-signer -config /path/to/config.json
```

### What happens if the signer is compromised

If using a DB backend, then all the timestamp private keys stored on the signer
will be compromised.  The attacker cannot do anything with the timestamp keys
unless they also compromise the Notary Server, though.

The attacker can prevent Notary Signer from signing any Notary Server metadata,
and return invalid public key IDs when the Notary Server requests it, and this
execute a denial of service attack, which would prevent the Notary Server from
being able to update any metadata.

### Ops features

Notary signer provides the following endpoints for ops friendliness:

1. A [Bugsnag](https://bugsnag.com) hook for error logs, if a Bugsnag
	configuration is provided.
