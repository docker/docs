---
title: HAProxy configuration
description: Learn how to configure an HAProxy extension
keywords: routing, proxy
---

# Configuring an HAProxy extension
The following configuration options are available:

| Option | Type | Description |
| --- | --- | --- |
| `PidPath` | string | path to the pid file for the proxy service |
| `MaxConnections` | int | maximum number of connections for proxy service |
| `ConnectTimeout` | int | timeout in seconds for clients to connect |
| `ClientTimeout` | int | timeout in seconds for the service to send a request to the proxied upstream |
| `ServerTimeout` | int | timeout in seconds for the service to read a response from the proxied upstream |
| `AdminUser` | string | username to be used with authenticated access to the proxy service |
| `AdminPass` | string | password to be used with authenticated access to the proxy service |
| `SSLOpts` | string | options to be passed when configuring SSL |
| `SSLDefaultDHParam` | int | size of DH parameters |
| `SSLVerify` | string | SSL client verification |
| `SSLCiphers` | string | SSL ciphers to use for the proxy service |
| `SSLProtocols` | string | enable the specified TLS protocols |
| `GlobalOptions` | []string | list of options that are included in the global configuration |
| `DefaultOptions` | []string | list of options that are included in the default configuration |

# Notes

When using SSL termination, the certificate and key must be combined into a single certificate (i.e. `cat cert.pem key.pem > combined.pem`).  The HAProxy extension only uses the certificate label to configure SSL.
