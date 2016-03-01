<!--[metadata]>
+++
title = "Server Configuration"
description = "Configuring the notary client, server and signer."
keywords = ["docker, notary, notary-client, notary-server, notary server, notary-signer, notary signer"]
[menu.main]
parent="mn_notary_config"
+++
<![end-metadata]-->

# Notary server configuration file

This document is for those who are [running their own Notary service](../running_a_service.md) who
want to specify custom options.

## Overview

A configuration file is required by Notary server, and the path to the
configuration file must be specified using the `-config` option on the command
line.

Here is a full server configuration file example; please click on the top level JSON keys to
learn more about the configuration section corresponding to that key:

<pre><code class="language-json">{
  <a href="#server-section-required">"server"</a>: {
    "http_addr": ":4443",
    "tls_key_file": "./fixtures/notary-server.key",
    "tls_cert_file": "./fixtures/notary-server.crt"
  },
  <a href="#trust-service-section-required">"trust_service"</a>: {
    "type": "remote",
    "hostname": "notarysigner",
    "port": "7899",
    "key_algorithm": "ecdsa",
    "tls_ca_file": "./fixtures/root-ca.crt",
    "tls_client_cert": "./fixtures/notary-server.crt",
    "tls_client_key": "./fixtures/notary-server.key"
  },
  <a href="#storage-section-required">"storage"</a>: {
    "backend": "mysql",
    "db_url": "user:pass@tcp(notarymysql:3306)/databasename?parseTime=true"
  },
  <a href="#auth-section-optional">"auth"</a>: {
    "type": "token",
    "options": {
      "realm": "https://auth.docker.io/token",
      "service": "notary-server",
      "issuer": "auth.docker.io",
      "rootcertbundle": "/path/to/auth.docker.io/cert"
    }
  },
  <a href="../common-configs/#logging-section-optional">"logging"</a>: {
    "level": "debug"
  },
  <a href="../common-configs/#reporting-section-optional">"reporting"</a>: {
    "bugsnag": {
      "api_key": "c9d60ae4c7e70c4b6c4ebd3e8056d2b8",
      "release_stage": "production"
    }
  }
}
</code></pre>

## server section (required)

Example:

```json
"server": {
  "http_addr": ":4443",
  "tls_key_file": "./fixtures/notary-server.key",
  "tls_cert_file": "./fixtures/notary-server.crt"
}
```

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>http_addr</code></td>
		<td valign="top">yes</td>
		<td valign="top">The TCP address (IP and port) to listen on.  Examples:
			<ul>
			<li><code>":4443"</code> means listen on port 4443 on all IPs (and
				hence all interfaces, such as those listed when you run
				<code>ifconfig</code>)</li>
			<li><code>"127.0.0.1:4443"</code> means listen on port 4443 on
				localhost only.  That means that the server will not be
				accessible except locally (via SSH tunnel, or just on a local
				terminal)</li>
			</ul>
		</td>
	</tr>
	<tr>
		<td valign="top"><code>tls_key_file</code></td>
		<td valign="top">no</td>
		<td valign="top">The path to the private key to use for
			HTTPS.  Must be provided together with <code>tls_cert_file</code>,
			or not at all. If neither are provided, the server will use HTTP
			instead of HTTPS. The path is relative to the directory of the
			configuration file.</td>
	</tr>
	<tr>
		<td valign="top"><code>tls_cert_file</code></td>
		<td valign="top">no</td>
		<td valign="top">The path to the certificate to use for HTTPS.
			Must be provided together with <code>tls_key_file</code>, or not
			at all. If neither are provided, the server will use HTTP instead
			of HTTPS. The path is relative to the directory of the
			configuration file.</td>
	</tr>
</table>


## trust_service section (required)

This section configures either a remote trust service, such as
[Notary signer](signer-config.md) or a local in-memory
ED25519 trust service.

Remote trust service example:

```json
"trust_service": {
  "type": "remote",
  "hostname": "notarysigner",
  "port": "7899",
  "key_algorithm": "ecdsa",
  "tls_ca_file": "./fixtures/root-ca.crt",
  "tls_client_cert": "./fixtures/notary-server.crt",
  "tls_client_key": "./fixtures/notary-server.key"
}
```

Local trust service example:

```json
"trust_service": {
  "type": "local"
}
```

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>type</code></td>
		<td valign="top">yes</td>
		<td valign="top">Must be <code>"remote"</code> or <code>"local"</code></td>
	</tr>
	<tr>
		<td valign="top"><code>hostname</code></td>
		<td valign="top">yes if remote</td>
		<td valign="top">The hostname of the remote trust service</td>
	</tr>
	<tr>
		<td valign="top"><code>port</code></td>
		<td valign="top">yes if remote</td>
		<td valign="top">The GRPC port of the remote trust service</td>
	</tr>
	<tr>
		<td valign="top"><code>key_algorithm</code></td>
		<td valign="top">yes if remote</td>
		<td valign="top">Algorithm to use to generate keys stored on the
			signing service.  Valid values are <code>"ecdsa"</code>,
			<code>"rsa"</code>, and <code>"ed25519"</code>.</td>
	</tr>
	<tr>
		<td valign="top"><code>tls_ca_file</code></td>
		<td valign="top">no</td>
		<td valign="top">The path to the root CA that signed the TLS
			certificate of the remote service. This parameter must be
			provided if said root CA is not in the system's default trust
			roots. The path is relative to the directory of the configuration file.</td>
	</tr>
	<tr>
		<td valign="top"><code>tls_client_key</code></td>
		<td valign="top">no</td>
		<td valign="top">The path to the private key to use for TLS mutual
			authentication. This must be provided together with
			<code>tls_client_cert</code> or not at all. The path is relative
			to the directory of the configuration file.</td>
	</tr>
	<tr>
		<td valign="top"><code>tls_client_cert</code></td>
		<td valign="top">no</td>
		<td valign="top">The path to the certificate to use for TLS mutual
			authentication. This must be provided together with
			<code>tls_client_key</code> or not at all. The path is relative
			to the directory of the configuration file.</td>
	</tr>
</table>

## storage section (required)

The storage section specifies which storage backend the server should use to
store TUF metadata.  Only MySQL or an in-memory store is supported.

DB storage example:

```json
"storage": {
  "backend": "mysql",
  "db_url": "user:pass@tcp(notarymysql:3306)/databasename?parseTime=true"
}
```

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>backend</code></td>
		<td valign="top">yes</td>
		<td valign="top">Must be <code>"mysql"</code> or <code>"memory"</code>.
			If <code>"memory"</code> is selected, the <code>db_url</code>
			is ignored.</td>
	</tr>
	<tr>
		<td valign="top"><code>db_url</code></td>
		<td valign="top">yes if not <code>memory</code></td>
		<td valign="top">The <a href="https://github.com/go-sql-driver/mysql">
			the Data Source Name used to access the DB.</a>
			(note: please include <code>parseTime=true</code> as part of the the DSN)</td>
	</tr>
</table>


## auth section (optional)

This sections specifies the authentication options for the server.
Currently, we only support token authentication.

Example:

```json
"auth": {
  "type": "token",
  "options": {
    "realm": "https://auth.docker.io",
    "service": "notary-server",
    "issuer": "auth.docker.io",
    "rootcertbundle": "/path/to/auth.docker.io/cert"
  }
}
```

Note that this entire section is optional.  However, if you would like
authentication for your server, then you need the required parameters below to
configure it.

**Token authentication:**

This is an implementation of the same authentication used by version 2 of the
[Docker registry](https://github.com/docker/distribution).  (JWT token-based
authentication post login.)

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>type</code></td>
		<td valign="top">yes</td>
		<td valign="top">Must be <code>"token"</code>; all other values will result in no
			authentication (and the rest of the parameters will be ignored)</td>
	</tr>
	<tr>
		<td valign="top"><code>options</code></td>
		<td valign="top">yes</td>
		<td valign="top">The options for token auth.  Please see
			<a href="https://github.com/docker/distribution/blob/master/docs/configuration.md#token">
			the registry token configuration documentation</a>
			for the parameter details.</td>
	</tr>
</table>

## Related information

* [Notary Signer Configuration File](signer-config.md)
* [Configuration sections common to the Notary Server and Signer](common-configs.md)
