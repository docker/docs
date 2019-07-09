---
title: Custom templates
description: Learn how to use a custom extension template
keywords: routing, proxy
---

# Using a custom extension template
A custom extension template can be
used if a needed option is not available in the extension configuration.

> Warning: This should be used with extreme caution as this completely bypasses
> the built-in extension template. Therefore, if you update the extension
> image in the future, you will not receive the updated template because you are
> using a custom one.

To use a custom template:

1. Create a Swarm configuration using a new template
2. Create a Swarm configuration object
3. Update the extension

## Creating a Swarm configuration using a new template
First, create a Swarm config using the new template, as shown in the following example. This example uses a custom Nginx configuration template, but you can use any extension configuration (for example, HAProxy).

The contents of the example `custom-template.conf` include:

{% raw %}
```
# CUSTOM INTERLOCK CONFIG
user {{ .ExtensionConfig.User  }};
worker_processes {{ .ExtensionConfig.WorkerProcesses  }};

error_log  {{ .ExtensionConfig.ErrorLogPath  }} warn;
pid        {{ .ExtensionConfig.PidPath  }};


events {
    worker_connections {{ .ExtensionConfig.MaxConnections  }};

}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    server_names_hash_bucket_size 128;

    # add custom HTTP options here, etc.

    log_format  main {{ .ExtensionConfig.MainLogFormat  }}

    log_format trace {{ .ExtensionConfig.TraceLogFormat  }}

    access_log  {{ .ExtensionConfig.AccessLogPath  }} main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  {{ .ExtensionConfig.KeepaliveTimeout  }};
    client_max_body_size {{ .ExtensionConfig.ClientMaxBodySize  }};
    client_body_buffer_size {{ .ExtensionConfig.ClientBodyBufferSize  }};
    client_header_buffer_size {{ .ExtensionConfig.ClientHeaderBufferSize  }};
    large_client_header_buffers {{ .ExtensionConfig.LargeClientHeaderBuffers  }};
    client_body_timeout {{ .ExtensionConfig.ClientBodyTimeout  }};
    underscores_in_headers {{ if .ExtensionConfig.UnderscoresInHeaders  }}on{{ else  }}off{{ end  }};

    add_header x-request-id $request_id;
    add_header x-proxy-id $hostname;
    add_header x-server-info "{{ .Version  }}";
    add_header x-upstream-addr $upstream_addr;
    add_header x-upstream-response-time $upstream_response_time;

    proxy_connect_timeout {{ .ExtensionConfig.ConnectTimeout  }};
    proxy_send_timeout {{ .ExtensionConfig.SendTimeout  }};
    proxy_read_timeout {{ .ExtensionConfig.ReadTimeout  }};
    proxy_set_header        X-Real-IP         $remote_addr;
    proxy_set_header        X-Forwarded-For   $proxy_add_x_forwarded_for;
    proxy_set_header        Host              $http_host;
    proxy_set_header x-request-id $request_id;
    send_timeout {{ .ExtensionConfig.SendTimeout  }};
    proxy_next_upstream error timeout invalid_header http_500 http_502 http_503 http_504;

    ssl_prefer_server_ciphers on;
    ssl_ciphers {{ .ExtensionConfig.SSLCiphers  }};
    ssl_protocols {{ .ExtensionConfig.SSLProtocols  }};
    {{ if (and (gt .ExtensionConfig.SSLDefaultDHParam 0) (ne .ExtensionConfig.SSLDefaultDHParamPath ""))  }}ssl_dhparam {{ .ExtensionConfig.SSLDefaultDHParamPath  }};{{ end  }}

    map $http_upgrade $connection_upgrade {
        default upgrade;
        ''      close;
    }

    {{ if not .HasDefaultBackend  }}
    # default host return 503
    server {
    listen {{ .Port  }} default_server;
    server_name _;

    root /usr/share/nginx/html;

    error_page   503 /503.html;
    location = /503.html {
        try_files /503.html @error;
        internal;
    }

    location @error {
        root /usr/share/nginx/html;
    }

    location / {
        return 503;

    }

    location /nginx_status {
        stub_status on;
        access_log off;
    }

    }
    {{ end  }}

    {{ range $host, $backends := .Hosts  }}
    {{ with $hostBackend := index $backends 0  }}
    {{ $sslBackend := index $.SSLBackends $host  }}
    upstream {{ backendName $host  }} {
        {{ if $hostBackend.IPHash  }}ip_hash; {{else}}zone {{ backendName $host  }}_backend 64k;{{ end  }}
    {{ if ne $hostBackend.StickySessionCookie ""  }}hash $cookie_{{ $hostBackend.StickySessionCookie  }} consistent; {{ end  }}
    {{ range $backend := $backends  }}
        {{ range $up := $backend.Targets  }}server {{ $up  }};
        {{ end  }}
        {{ end  }} {{/* end range backends */}}

    }
    {{ if not $sslBackend.Passthrough  }}
    server {
        listen {{ $.Port  }}{{ if $hostBackend.DefaultBackend  }} default_server{{ end  }};
    {{ if $hostBackend.DefaultBackend  }}server_name _;{{ else  }}server_name {{$host}};{{ end  }}

    {{ if (isRedirectHost $host $hostBackend.Redirects)  }}
    {{ range $redirect := $hostBackend.Redirects  }}
        {{ if isRedirectMatch $redirect.Source $host  }}return 302 {{ $redirect.Target  }}$request_uri;{{ end  }}
    {{ end  }}
    {{ else  }}

    {{ if eq ( len $hostBackend.ContextRoots  ) 0  }}
    {{ if not (isWebsocketRoot $hostBackend.WebsocketEndpoints)  }}
    location / {
            proxy_pass {{ if $hostBackend.SSLBackend  }}https://{{ else  }}http://{{ backendName $host  }};{{ end  }}
    }
    {{ end  }}

        {{ range $ws := $hostBackend.WebsocketEndpoints  }}
        location {{ $ws  }} {
            proxy_pass {{ if $hostBackend.SSLBackend  }}https://{{ else  }}http://{{ backendName $host  }};{{ end  }}
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection $connection_upgrade;
            proxy_set_header Origin '';
        }
        {{ end  }} {{/* end range WebsocketEndpoints */}}
    {{ else  }}

    {{ range $ctxroot := $hostBackend.ContextRoots  }}
    location {{ $ctxroot.Path  }} {
        {{ if $ctxroot.Rewrite  }}rewrite ^([^.]*[^/])$ $1/ permanent;
        rewrite  ^{{ $ctxroot.Path  }}/(.*)  /$1 break;{{ end  }}
        proxy_pass http://{{ backendName $host  }};
    }
    {{ end  }} {{/* end range contextroots */}}

    {{ end  }} {{/* end len $hostBackend.ContextRoots */}}
    location /nginx_status {
            stub_status on;
            access_log off;
    }
    {{ end  }}{{/* end isRedirectHost */}}

    }
    {{ end  }} {{/* end if not sslBackend.Passthrough */}}

    {{/* SSL */}}
    {{ if ne $hostBackend.SSLCert ""  }}
    {{ $sslBackend := index $.SSLBackends $host  }}
    server {
    listen 127.0.0.1:{{ $sslBackend.Port  }} ssl proxy_protocol;
        server_name {{ $host  }};
        ssl on;
        ssl_certificate /run/secrets/{{ $hostBackend.SSLCertTarget  }};
    {{ if ne $hostBackend.SSLKey ""  }}ssl_certificate_key /run/secrets/{{ $hostBackend.SSLKeyTarget  }};{{ end  }}
    set_real_ip_from 127.0.0.1/32;
    real_ip_header proxy_protocol;

    {{ if eq ( len $hostBackend.ContextRoots  ) 0  }}
    {{ if not (isWebsocketRoot $hostBackend.WebsocketEndpoints)  }}
    location / {
            proxy_pass {{ if $hostBackend.SSLBackend  }}https://{{ else  }}http://{{ backendName $host  }};{{ end  }}
    }
    {{ end  }}

        {{ range $ws := $hostBackend.WebsocketEndpoints  }}
        location {{ $ws  }} {
            proxy_pass {{ if $hostBackend.SSLBackend  }}https://{{ else  }}http://{{ backendName $host  }};{{ end  }}
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection $connection_upgrade;
        proxy_set_header Origin {{$host}};

        }
        {{ end  }} {{/* end range WebsocketEndpoints */}}
    {{ else  }}

    {{ range $ctxroot := $hostBackend.ContextRoots  }}
    location {{ $ctxroot.Path  }} {
        {{ if $ctxroot.Rewrite  }}rewrite ^([^.]*[^/])$ $1/ permanent;
        rewrite  ^{{ $ctxroot.Path  }}/(.*)  /$1 break;{{ end  }}
        proxy_pass http://{{ backendName $host  }};
    }
    {{ end  }} {{/* end range contextroots */}}

    {{ end  }} {{/* end len $hostBackend.ContextRoots */}}
    location /nginx_status {
            stub_status on;
            access_log off;
    }

    } {{ end  }} {{/* end $hostBackend.SSLCert */}}
    {{ end  }} {{/* end with hostBackend */}}

    {{ end  }} {{/* end range .Hosts */}}

    include       /etc/nginx/conf.d/*.conf;
}
stream {
    # main log compatible format
    log_format stream '$remote_addr - - [$time_local] "$ssl_preread_server_name -> $name ($protocol)" '
                          '$status $bytes_sent "" "" "" ';
                          map $ssl_preread_server_name $name {
    {{ range $host, $sslBackend := $.SSLBackends  }}
    {{ $sslBackend.Host  }} {{ if $sslBackend.Passthrough  }}pt-{{ backendName $host  }};{{ else  }}127.0.0.1:{{ $sslBackend.Port  }}; {{ end  }}
    {{ if $sslBackend.DefaultBackend  }}default {{ if $sslBackend.Passthrough  }}pt-{{ backendName $host  }};{{ else  }}127.0.0.1:{{ $sslBackend.Port  }}; {{ end  }}{{ end  }}
    {{ end  }}

                          }
    {{ range $host, $sslBackend := $.SSLBackends  }}
    upstream pt-{{ backendName $sslBackend.Host  }} {
    {{ $h := index $.Hosts $sslBackend.Host  }}{{ $hostBackend := index $h 0  }}
    {{ if $sslBackend.Passthrough  }}
    server 127.0.0.1:{{ $sslBackend.ProxyProtocolPort  }};
    {{ else  }}
    {{ range $up := $hostBackend.Targets  }}server {{ $up  }};
    {{ end  }} {{/* end range backend targets */}}
    {{ end  }} {{/* end range sslbackend */}}

    }{{ end  }} {{/* end range SSLBackends */}}

    {{ range $host, $sslBackend := $.SSLBackends  }}
    {{ $proxyProtocolPort := $sslBackend.ProxyProtocolPort  }}
    {{ $h := index $.Hosts $sslBackend.Host  }}{{ $hostBackend := index $h 0  }}
    {{ if ne $proxyProtocolPort 0  }}
    upstream proxy-{{ backendName $sslBackend.Host  }} {
    {{ range $up := $hostBackend.Targets  }}server {{ $up  }};
    {{ end  }} {{/* end range backend targets */}}

    }
    server {
    listen {{ $proxyProtocolPort  }} proxy_protocol;
    proxy_pass proxy-{{ backendName $sslBackend.Host  }};

    }
    {{ end  }} {{/* end if ne proxyProtocolPort 0 */}}
    {{ end  }} {{/* end range SSLBackends */}}

    server {
        listen {{ $.SSLPort  }};
        proxy_pass $name;
        proxy_protocol on;
        ssl_preread on;
        access_log {{ .ExtensionConfig.AccessLogPath  }} stream;
    }
}
```
{% endraw %}

## Creating a Swarm configuration object
To create a Swarm config object:

```
$> docker config create interlock-custom-template custom.conf
```

## Updating the extension
Now update the extension to use this new template:

```
$> docker service update --config-add source=interlock-custom-template,target=/etc/docker/extension-template.conf interlock-ext
```

This should trigger an update and a new proxy configuration will be generated.

## Removing the custom template
To remove the custom template and revert to using the built-in template:

```
$> docker service update --config-rm interlock-custom-template interlock-ext
```
