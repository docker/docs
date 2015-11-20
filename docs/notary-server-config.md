<!--[metadata]>
+++
title = "Notary Serer Configuration File"
description = "Specifies the configuration file for Notary Server"
keywords = ["docker, notary, notary-server, configuration"]
[menu.main]
parent="mn_notary"
+++
<![end-metadata]-->

# Notary Server Configuration File

An example (full) server configuration file.

```json
{
	"server": {
		"addr": ":4443",
		"tls_key_file": "./fixtures/notary-server.key",
		"tls_cert_file": "./fixtures/notary-server.crt",
		"auth": {
			"type": "token",
			"options": {
				"realm": "https://auth.docker.io/token",
				"service": "notary-server",
				"issuer": "auth.docker.io",
				"rootcertbundle": "/path/to/auth.docker.io/cert"
			}
		}
	},
	"trust_service": {
		"type": "remote",
		"hostname": "notarysigner",
		"port": "7899",
		"key_algorithm": "ecdsa",
		"tls_ca_file": "./fixtures/root-ca.crt",
		"tls_client_cert": "./fixtures/notary-server.crt",
		"tls_client_key": "./fixtures/notary-server.key"
	},
	"logging": {
		"level": "debug"
	},
	"storage": {
		"backend": "mysql",
		"db_url": "dockercondemo:dockercondemo@tcp(notarymysql:3306)/dockercondemo"
	},
	"reporting": {
		"bugsnag": "yes",
		"bugsnag_api_key": "c9d60ae4c7e70c4b6c4ebd3e8056d2b8",
		"bugsnag_release_stage": "notary-server"
	}
}
```

## `server` section (required)

Example:

```json
"server": {
	"addr": ":4443",
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
		<td valign="top"><code>addr</code></td>
		<td valign="top">yes</td>
		<td valign="top">The TCP address (IP and port) to listen on.  Examples:
			<ul>
			<li><code>":4443"</code> means listen on port 4443 on all IPs (and
				hence all interfaces, such as those listed when you run
				<code>ifconfig</code>)</li>
			<li><code>"127.0.0.1:4443"</code> means listen on port 4443 on
				localhost only.  That means that the server will not be
				acessible except locally (via SSH tunnel, or just on a local
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
			instead of HTTPS. The path is relative to the directory where
			notary-server is run.</td>
	</tr>
	<tr>
		<td valign="top"><code>tls_cert_file</code></td>
		<td valign="top">no</td>
		<td valign="top">The path to the certificate to use for HTTPS.
			Must be provided together with <code>tls_key_file</code>, or not
			at all. If neither are provided, the server will use HTTP instead
			of HTTPS. The path is relative to the directory where notary-server
			is run.</td>
	</tr>
</table>

### `auth` subsection (optional)

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

This is an implementation of the same authentication used by
[docker registry](https://github.com/docker/distribution).  (JTW token-based
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
		<td valign="top">Must be `"token"`; all other values will result in no
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

## `trust service` section (optional but recommended)

This section is required to specify a remote trust service, such as
[Notary Signer](notary-signer.md).  If it is left out or invalid, a local
in-memory ED25519 trust service will be used instead.

Remote trust service example:

```json
"trust_service": {
	"type": "remote",
	"hostname": "notarysigner",
	"port": "7899",
	"key_algorithm": "ecdsa",
	"tls_ca_file": "./fixtures/root-ca.crt",
	"tls_client_key": "./fixtures/notary-server.key",
	"tls_client_cert": "./fixtures/notary-server.crt"
}
```

Note that this entire section is optional.  However, if you would like to use a
separate trust service (recommended), then you need the required parameters
below to configure it.

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>type</code></td>
		<td valign="top">yes</td>
		<td valign="top">Must be <code>"remote"</code>; all other values
			will result in a local trust service (and the rest of the
			parameters will be ignored)</td>
	</tr>
	<tr>
		<td valign="top"><code>hostname</code></td>
		<td valign="top">yes</td>
		<td valign="top">The hostname of the remote trust service</td>
	</tr>
	<tr>
		<td valign="top"><code>port</code></td>
		<td valign="top">yes</td>
		<td valign="top">The GRPC port of the remote trust service</td>
	</tr>
	<tr>
		<td valign="top"><code>key_algorithm</code></td>
		<td valign="top">yes</td>
		<td valign="top">Algorithm to use to generate keys stored on the
			signing service.  Valid values are <code>"ecdsa"</code>,
			<code>"rsa"</code>, and <code>"ed25519"</code>.</td>
	</tr>
	<tr>
		<td valign="top"><code>tls_ca_file</code></td>
		<td valign="top">no</td>
		<td valign="top">The path to the root CA that signed the TLS
			certificate of the remote service. This parameter if said root
			CA is not in the system's default trust roots. The path is
			relative to the directory where notary-server is run.</td>
	</tr>
	<tr>
		<td valign="top"><code>tls_client_key</code></td>
		<td valign="top">no</td>
		<td valign="top">The path to the private key to use for TLS mutual
			authentication. This must be provided together with
			<code>tls_client_cert</code> or not at all. The path is relative
			to the directory where notary-server is run.</td>
	</tr>
	<tr>
		<td valign="top"><code>tls_client_cert</code></td>
		<td valign="top">no</td>
		<td valign="top">The path to the certificate to use for TLS mutual
			authentication. This must be provided together with
			<code>tls_client_key</code> or not at all. The path is relative
			to the directory where notary-server is run.</td>
	</tr>
</table>


## `logging` section (optional)

The logging section sets the log level of the server.  If it is not provided
or invalid, the server defaults to an ERROR logging level.

Example:

```json
"logging": {
	"level": "debug"
}
```

Note that this entire section is optional.  However, if you would like to
specify a different log level, then you need the required parameters
below to configure it.

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>level</code></td>
		<td valign="top">yes</td>
		<td valign="top">One of <code>"debug"</code>, <code>"info"</code>,
			<code>"warning"</code>, <code>"error"</code>, <code>"fatal"</code>,
			or <code>"panic"</code></td>
	</tr>
</table>

## `storage` section (optional but recommended)

The storage section specifies which storage backend the server should use to
store TUF metadata.  Currently, we only support MySQL. If this
section is not provided or invalid, an in-memory store will be used instead.

DB storage example:

```json
"storage": {
	"backend": "mysql",
	"db_url": "dockercondemo:dockercondemo@tcp(notarymysql:3306)/dockercondemo"
}
```

Note that this entire section is optional.  However, if you would like to
use a database backend (recommended), then you need the required parameters
below to configure it.

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>backend</code></td>
		<td valign="top">yes</td>
		<td valign="top">Must be <code>"mysql"</code>; all other values will
			result in an in-memory store (and the rest of the parameters will
			be ignored)</td>
	</tr>
	<tr>
		<td valign="top"><code>db_url</code></td>
		<td valign="top">yes</td>
		<td valign="top">The URL used to access the DB, which includes both the
			endpoint the username/credentials</td>
	</tr>
</table>

## `reporting` section (optional)

The reporting section contains any configuration for useful for running the
service, such as reporting errors. Currently, we only support reporting errors
to [Bugsnag](https://bugsnag.com).

See [bugsnag-go](https://github.com/bugsnag/bugsnag-go/) for more information
about these configuration parameters.

```json
"reporting": {
	"bugsnag": "yes",
	"bugsnag_api_key": "c9d60ae4c7e70c4b6c4ebd3e8056d2b8",
	"bugsnag_release_stage": "notary-server"
}
```

Note that this entire section is optional.  However, if you would like to
report errors to Bugsnag, then you need the required parameters below to
configure it.

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>bugsnag</code></td>
		<td valign="top">yes</td>
		<td>Any string value. If this value is not set, no errors will be
			reported to Bugsnag (all other parameters will be ignored)</td>
	</tr>
	<tr>
		<td valign="top"><code>bugsnag_api_key</code></td>
		<td valign="top">yes</td>
		<td>The API key to use to report errors.</td>
	</tr>
	<tr>
		<td valign="top"><code>bugsnag_release_stage</code></td>
		<td valign="top">yes</td>
		<td>The current release stage, such as "production".  You can
			use this value to filter errors in the Bugsnag dashboard.</td>
	</tr>
</table>
