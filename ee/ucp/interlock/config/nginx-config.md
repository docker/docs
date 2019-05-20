---
title: Configure Nginx 
description: Learn how to configure an nginx extension
keywords: routing, proxy, interlock, load balancing
---

By default, nginx is used as a proxy, so the following configuration options are
available for the nginx extension:

| Option             | Type        | Description . | Defaults  |
|:------ |:------ |:------ |:------ |
| `User` | string | User to be used in the proxy | `nginx` |
| `PidPath` | string | Path to the pid file for the proxy service | `/var/run/proxy.pid` |
| `MaxConnections` | int | Maximum number of connections for proxy service | `1024` |
| `ConnectTimeout` | int | Timeout in seconds for clients to connect | `600` |
| `SendTimeout` | int | Timeout in seconds for the service to send a request to the proxied upstream | `600` |
| `ReadTimeout` | int | Timeout in seconds for the service to read a response from the proxied upstream | `600` |
| `SSLOpts` | string | Options to be passed when configuring SSL |  |
| `SSLDefaultDHParam` | int | Size of DH parameters | `1024` |
| `SSLDefaultDHParamPath` | string | Path to DH parameters file | |
| `SSLVerify` | string | SSL client verification | `required` |
| `WorkerProcesses` | string | Number of worker processes for the proxy service | `1` |
| `RLimitNoFile` | int | Number of maxiumum open files for the proxy service | `65535` |
| `SSLCiphers` | string | SSL ciphers to use for the proxy service | `HIGH:!aNULL:!MD5` |
| `SSLProtocols` | string | Enable the specified TLS protocols | `TLSv1.2` |
| `HideInfoHeaders`          | bool | Hide proxy-related response headers.                                                                  |
| `KeepaliveTimeout` | string | connection keepalive timeout | `75s` |
| `ClientMaxBodySize` | string | maximum allowed size of the client request body | `1m` |
| `ClientBodyBufferSize` | string | sets buffer size for reading client request body | `8k` |
| `ClientHeaderBufferSize` | string | sets buffer size for reading client request header | `1k` |
| `LargeClientHeaderBuffers` | string | sets the maximum number and size of buffers used for reading large client request header | `4 8k` |
| `ClientBodyTimeout` | string | timeout for reading client request body | `60s` |
| `UnderscoresInHeaders` | bool | enables or disables the use of underscores in client request header fields| `false` |
| `ServerNamesHashBucketSize` | int | sets the bucket size for the server names hash tables (in KB) | `128` |
| `UpstreamZoneSize` | int | size of the shared memory zone (in KB) | `64` |
| `GlobalOptions` | []string | list of options that are included in the global configuration | |
| `HTTPOptions` | []string | list of options that are included in the http configuration | |
| `TCPOptions` | []string | list of options that are included in the stream (TCP) configuration | |
| `AccessLogPath` | string | Path to use for access logs | `/dev/stdout` |
| `ErrorLogPath` | string | Path to use for error logs | `/dev/stdout` |
| `MainLogFormat` | string | [Format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format) to use for main logger | see default format |
| `TraceLogFormat` | string | [Format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format) to use for trace logger | see default format |

