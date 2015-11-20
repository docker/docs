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

As an example, please see [this script](opensslCertGen.sh) to see how to
generate client SSL certs with basic constraints using OpenSSL.

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

#### Running a Docker image

Get the official Docker image, which comes with [some defaults](
https://github.com/docker/notary/blob/master/cmd/notary-signer/config.json).
You can override the default configuration with environment variables.
For example, if you wanted to run it with your own MySQL DB and a different
logging level:

```
$ docker pull docker.io/docker/notary-signer
$ docker run -p "4443:4443" \
	-e NOTARY_SIGNER_LOGGING_LEVEL=info \
	-e NOTARY_SIGNER_STORAGE_DB_URL="myuser:mypass@tcp(my-db)/dbName"
	notary-signer
```

Alternately, you can run the image with your own configuration file entirely.
The docker image loads the config file from `/opt/notary-signer/config.json`,
so you can mount a directory with your config file (named `config.json`)
at `/opt/notary-signer`:

```
$ docker run -p "4443:4443" -v /path/to/config/dir:/opt/notary-signer notary-signer
```

#### Running the binary
A JSON configuration file needs to be passed as a parameter/flag when starting
up the Notary Signer binary.  Environment variables can also be set in addition
to the configuration file, but the configuration file is required.  For example:

```
$ export NOTARY_SIGNER_STORAGE_DB_URL=myuser:mypass@tcp(my-db)/dbname
$ NOTARY_SIGNER_LOGGING_LEVEL=info notary-signer -config /path/to/config.json
```

### What happens if the signer is compromised

All the timestamp private keys stored on the signer will be compromised, and
an attacker can sign anything they wish with the timestamp key.

However, the attacker cannot do anything useful with the timestamp keys unless
they also [compromise the Notary Server](
notary-server.md#what-happens-if-the-server-is-compromised)

The attacker can prevent Notary Signer from signing timestap metadata from
Notary Server and return invalid public key IDs when the Notary Server
requests it.  This means an attacker can execute a denial of service attack
that prevents the Notary Server from being able to update any metadata.
