<!--[metadata]>
+++
title = "Client configuration"
description = "Configuring the Notary client, server and signer."
keywords = ["docker, notary, notary-client, notary-server, notary server, notary-signer, notary signer"]
[menu.main]
parent="mn_notary_config"
weight=4
+++
<![end-metadata]-->

# Notary client configuration file

This document is for power users of the [Notary client](../advanced_usage.md),
or for those who are [running their own service](../running_a_service.md) who
want to facilitate CLI interaction or specify custom options.

The configuration file for Notary client normally resides at `~/.notary/config.json`,
but the path to a different configuration file can be specified using the
`-c` or `--configFile` command line flag.

## Overview of the file

A configuration file consists of the following sections:

<table>
	<tr>
		<td><a href="#trust-dir-section-optional">trust_dir</a></td>
		<td>TUF key and metadata directory</td>
		<td>(optional)</td>
	</tr>
	<tr>
		<td><a href="#remote-server-section-optional">remote_server</a></td>
		<td>remote Notary Server configuration</td>
		<td>(optional)</td>
	</tr>
</table>

In addition please see the optional password [environment
variables](#environment-variables-optional) that the Notary client can take for
ease of use.

An example (full) client configuration file.

```json
{
  "trust_dir" : "~/.docker/trust",
  "remote_server": {
    "url": "https://my-notary-server.my-private-registry.com",
    "root-ca": "./fixtures/root-ca.crt",
    "tls_client_cert": "./fixtures/secure.example.com.crt",
    "tls_client_key": "./fixtures/secure.example.com.crt"
  }
}
```

## trust_dir section (optional)

The `trust_dir` specifies the location (as an absolute path or a path
relative to the directory of the configuration file) where the TUF metadata
and private keys will be stored.

This is normally defaults to `~/.notary`, but specifying `~/.docker/trust`
facilitates interoperability with content trust.

Note that this option can be overridden with the command line flag `--trustDir`.

## remote_server section (optional)

The `remote_server` specifies how to connect to a Notary server to download
metadata updates and publish metadata changes.

Remote server example:

```json
"remote_server": {
  "url": "https://my-notary-server.my-private-registry.com",
  "root-ca": "./fixtures/root-ca.crt",
  "tls_client_cert": "./fixtures/secure.example.com.crt",
  "tls_client_key": "./fixtures/secure.example.com.crt"
}
```

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>url</code></td>
		<td valign="top">no</td>
		<td valign="top">URL of the Notary server: defaults to https://notary.docker.io
			This configuration option can be overridden with the command line flag
			`-s` or `--server`.</td>
	</tr>
	<tr>
		<td valign="top"><code>root-ca</code></td>
		<td valign="top">no</td>
		<td valign="top"><p>The path to the file containing the root CA with which to verify
			the TLS certificate of the Notary server, for example if it is self-signed.
			The path is relative to the directory of the configuration file.</p>
			<p>This configuration option can overridden with the command line flag
			`--tlscacert`, which would specify a path relative to the current working
			directory where the Notary client is invoked.</p></td>
	</tr>
	<tr>
		<td valign="top"><code>tls_client_cert</code></td>
		<td valign="top">no</td>
		<td valign="top"><p>The path to the client certificate to use for mutual TLS with
			the Notary server.  Must be provided along with <code>tls_client_key</code>
			or not provided at all.  The path is relative to the directory of the
			configuration file.</p>
			<p>This configuration option can overridden with the command line flag
			`--tlscert`, which would specify a path relative to the current working
			directory where the Notary client is invoked.</p></td>
	</tr>
	<tr>
		<td valign="top"><code>tls_client_key</code></td>
		<td valign="top">no</td>
		<td valign="top"><p>The path to the client key to use for mutual TLS with
			the Notary server. Must be provided along with <code>tls_client_cert</code>
			or not provided at all.  The path is relative to the directory of the
			configuration file.</p>
			<p>This configuration option can overridden with the command line flag
			`--tlskey`, which would specify a path relative to the current working
			directory where the Notary client is invoked.</p></td>
	</tr>
</table>

## Environment variables (optional)

The following environment variables containing signing key passphrases can
be used to facilitate [Notary client CLI interaction](../advanced_usage.md).
If provided, these passwords will be used initially to sign TUF metadata.
If the passphrase is incorrect, you will be prompted to enter the correct
passphrase.


| Environment Variable        | Description                             |
| --------------------------- | --------------------------------------- |
|`NOTARY_ROOT_PASSPHRASE`     | The root/offline key passphrase         |
|`NOTARY_TARGETS_PASSPHRASE`  | The targets (an online) key passphrase  |
|`NOTARY_SNAPSHOT_PASSPHRASE` | The snapshot (an online) key passphrase |

## Related information

* [Notary Server Configuration File](server-config.md)
* [Notary Signer Configuration File](signer-config.md)
* [Configuration sections common to the Notary Server and Signer](common-configs.md)
