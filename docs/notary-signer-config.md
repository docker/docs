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
		"tls_cert_file": "./fixtures/notary-signer.crt",
		"tls_key_file": "./fixtures/notary-signer.key",
		"client_ca_file": "./fixtures/notary-server.crt"
	},
	"logging": {
		"level": 2
	},
	"storage": {
		"backend": "mysql",
		"db_url": "user:pass@tcp(notarymysql:3306)/databasename?parseTime=true",
		"default_alias": "password1"
	},
	"reporting": {
		"bugsnag": {
			"api_key": "c9d60ae4c7e70c4b6c4ebd3e8056d2b8",
			"release_stage": "production"
		}
	}
}
```

## `server` section (required)

"server" in this case refers to Notary Signer's HTTP/GRPC server, not
"Notary Server".

Example:

```json
"server": {
	"http_addr": ":4444",
	"grpc_addr": ":7899",
	"tls_cert_file": "./fixtures/notary-signer.crt",
	"tls_key_file": "./fixtures/notary-signer.key",
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
		<td valign="top"><code>tls_key_file</code></td>
		<td valign="top">yes</td>
		<td valign="top">The path to the private key to use for
			HTTPS. The path is relative to the directory of the
			configuration file.</td>
	</tr>
	<tr>
		<td valign="top"><code>tls_cert_file</code></td>
		<td valign="top">yes</td>
		<td valign="top">The path to the certificate to use for
			HTTPS. The path is relative to the directory of the
			configuration file.</td>
	</tr>
	<tr>
		<td valign="top"><code>client_ca_file</code></td>
		<td valign="top">no</td>
		<td valign="top">The root certificate to trust for
			mutual authentication. If provided, any clients connecting to
			Notary Signer will have to have a client certificate signed by
			this root. If not provided, mutual authentication will not be
			required. The path is relative to the directory of the
			configuration file.</td>
	</tr>
</table>

## `storage` section (required)

This is used to store encrypted priate keys.  We only support MySQL or an
in-memory store, currently.

Example:

```json
"storage": {
	"backend": "mysql",
	"db_url": "user:pass@tcp(notarymysql:3306)/databasename?parseTime=true",
	"default_alias": "password1"
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
			(note: please include "parseTime=true" as part of the the DSN)</td>
	</tr>
	<tr>
		<td valign="top"><code>default_alias</code></td>
		<td valign="top">yes if not <code>memory</code></td>
		<td valign="top">This parameter specifies the alias of the current
			password used to encrypt the private keys in the DB.  All new
			private keys will be encrypted using this password, which
			must also be provided as the environment variable
			<code>NOTARY_SIGNER_&lt;default_alias_value&gt;</code>.</td>
	</tr>
</table>

**Required environment variable:**
<code>NOTARY_SIGNER_&lt;default_alias_value&gt;</code>

**Optional environment variables:**
<code>NOTARY_SIGNER_&lt;old_alias_value&gt;</code> for as many old alias values
as needed.  This will ensure that older private keys, encrypted with older
passwords, can be decrypted.

## `logging` section (optional)

The logging section sets the log level of the server.  If it is not provided
or invalid, the signer defaults to an ERROR logging level.

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


## `reporting` section (optional)

The reporting section contains any configuration for useful for running the
service, such as reporting errors. Currently, we only support reporting errors
to [Bugsnag](https://bugsnag.com).

See [bugsnag-go](https://github.com/bugsnag/bugsnag-go/) for more information
about these configuration parameters.

```json
"reporting": {
	"bugsnag": {
		"api_key": "c9d60ae4c7e70c4b6c4ebd3e8056d2b8",
		"release_stage": "production"
	}
}
```

Note that this entire section is optional.  However, if you would like to
report errors to Bugsnag, then you need to include a `bugsnag` subsection,
along with the required parameters below, to configure it.

**Bugsnag reporting:**

<table>
	<tr>
		<th>Parameter</th>
		<th>Required</th>
		<th>Description</th>
	</tr>
	<tr>
		<td valign="top"><code>api_key</code></td>
		<td valign="top">yes</td>
		<td>The BugSnag API key to use to report errors.</td>
	</tr>
	<tr>
		<td valign="top"><code>release_stage</code></td>
		<td valign="top">yes</td>
		<td>The current release stage, such as "production".  You can
			use this value to filter errors in the Bugsnag dashboard.</td>
	</tr>
</table>
