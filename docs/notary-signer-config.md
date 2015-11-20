<!--[metadata]>
+++
title = "Notary Signer Configuration File"
description = "Specifies the configuration file for Notary Signer"
keywords = ["docker, notary, notary-signer, configuration"]
[menu.main]
parent="mn_notary"
+++
<![end-metadata]-->

# Notary Signer Configuration File

An example (full) server configuration file.

```json
{
	"server": {
		"http_addr": ":4444",
		"grpc_addr": ":7899",
		"cert_file": "./fixtures/notary-signer.crt",
		"key_file": "./fixtures/notary-signer.key",
		"client_ca_file": "./fixtures/notary-server.crt"
	},
	"logging": {
		"level": 2
	},
	"storage": {
		"backend": "mysql",
		"db_url": "dockercondemo:dockercondemo@tcp(notarymysql:3306)/dockercondemo"
	}
}
```

## `server` section (required)

Example:

```json
"server": {
	"http_addr": ":4444",
	"grpc_addr": ":7899",
	"cert_file": "./fixtures/notary-signer.crt",
	"key_file": "./fixtures/notary-signer.key",
	"client_ca_file": "./fixtures/notary-server.crt"
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
		<td valign="top">The TCP address (IP and port) to listen for HTTP
			traffic on.  Examples:
			<ul>
			<li><code>":4444"</code> means listen on port 4444 on all IPs (and
				hence all interfaces, such as those listed when you run
				<code>ifconfig</code>)</li>
			<li><code>"127.0.0.1:4444"</code> means listen on port 4444 on
				localhost only.  That means that the server will not be
				acessible except locally (via SSH tunnel, or just on a local
				terminal)</li>
			</ul>
		</td>
	</tr>
	<tr>
		<td valign="top"><code>grpc_addr</code></td>
		<td valign="top">yes</td>
		<td valign="top">The TCP address (IP and port) to listen for GRPC
			traffic.  Examples:
			<ul>
			<li><code>":7899"</code> means listen on port 7899 on all IPs (and
				hence all interfaces, such as those listed when you run
				<code>ifconfig</code>)</li>
			<li><code>"127.0.0.1:7899"</code> means listen on port 7899 on
				localhost only.  That means that the server will not be
				acessible except locally (via SSH tunnel, or just on a local
				terminal)</li>
			</ul>
		</td>
	</tr>
	<tr>
		<td valign="top"><code>key_file</code></td>
		<td valign="top">yes</td>
		<td valign="top">The path to the private key to use for
			HTTPS. The path is relative to the directory where
			notary-signer is run.</td>
	</tr>
	<tr>
		<td valign="top"><code>cert_file</code></td>
		<td valign="top">yes</td>
		<td valign="top">The path to the certificate to use for
			HTTPS. The path is relative to the directory where
			notary-signer is run.</td>
	</tr>
	<tr>
		<td valign="top"><code>client_ca_file</code></td>
		<td valign="top">no</td>
		<td valign="top">The root certificate to trust for
			mutual authentication. If provided, any clients connecting to
			Notary Signer will have to have a client certificate signed by
			this root. If not provided, mutual authentication will not be
			required. The path is relative to the directory where
			notary-signer is run.</td>
	</tr>
</table>

## `storage` section (required)

We only support MySQL, currently, and it must be provided.

Example:

```json
"storage": {
	"backend": "mysql",
	"db_url": "dockercondemo:dockercondemo@tcp(notarymysql:3306)/dockercondemo"
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
		<td valign="top">Must be <code>"mysql"</code></td>
	</tr>
	<tr>
		<td valign="top"><code>db_url</code></td>
		<td valign="top">yes</td>
		<td valign="top">The URL used to access the DB, which includes both the
			endpoint the username/credentials</td>
	</tr>
</table>


## `logging` section (optional)

The logging section sets the log level of the server.  If not provided, or if
any part of this section is invalid, the server defaults to an ERROR logging
level.

Example:

```json
"logging": {
	"level": 2
}
```

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>level</code></td>
		<td valign="top">yes</td>
		<td valign="top">An integer between 0 and 5, representing values
			<code>"debug"</code> (5), <code>"info"</code> (4),
			<code>"warning"</code> (3), <code>"error"</code>(2),
			<code>"fatal"</code> (1), or <code>"panic"</code>(0)</td>
	</tr>
</table>
