package shared

const NginxTemplate = `
worker_processes auto;

events {
  worker_connections 1024;
  use epoll;
  multi_accept on;
}

http {
  tcp_nodelay on;

  # this is necessary for us to be able to disable request buffering in all cases
  proxy_http_version 1.1;

{{ if .HasLicense }}
  upstream image_storage {
    least_conn;
    server {{.StorageServerName}}:{{.StorageServerPort}};

    check interval=2000 rise=1 fall=1 timeout=5000 type=tcp;
  }
{{ end }}

  upstream auth {
    server {{.StorageServerName}}:{{.GarantPort}};
    check interval=2000 rise=1 fall=1 timeout=5000 type=tcp;
  }

  upstream rethink {
    server {{.RethinkFullName}}:{{.RethinkPort}};
    check interval=2000 rise=1 fall=1 timeout=5000 type=tcp;
  }

  upstream etcd {
    server {{.EtcdFullName}}:{{.EtcdPort1}};
    server {{.EtcdFullName}}:{{.EtcdPort2}};
    check interval=2000 rise=1 fall=1 timeout=5000 type=tcp;
  }

  upstream api {
    server {{.AdminServerFullName}}:{{.AdminPort}};
    check interval=2000 rise=1 fall=1 timeout=5000 type=tcp;
  }

  upstream notary {
    server {{.NotaryHost}}:{{.NotaryPort}};
    check interval=2000 rise=1 fall=1 timeout=5000 type=tcp;
  }

  server {
    listen 80;

    add_header X-Replica-ID XXXREPLICA_IDXXX always;

    # disable any limits to avoid HTTP 413 for large image uploads
    client_max_body_size 0;


    # admin UI access should always use HTTPS
    location / {
      return 301 https://$host{{.ColonHTTPSPort}}$request_uri;
    }

    # we should serve only the health api unencrypted
    location /health {
      proxy_pass http://api/health;
      # we don't really need all of this. It's just copied from the https route
      proxy_set_header Host $http_host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_hide_header {{.RegistryEventsHeaderName}};
      proxy_buffering off;
      proxy_request_buffering off;
    }

    location /v1/ {
      return 404;
    }

    location /v1/users/ {
{{ if .HasLicense }}
{{ else }}
      error_page 404 /no_license;
{{ end }}
      return 404;
    }

    location /v2/ {
{{ if .HasLicense }}
      proxy_pass http://image_storage/v2/;
      proxy_set_header Host $http_host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_buffering off;
      proxy_request_buffering off;
{{ else }}
      return 404;
      error_page 404 /no_license;
{{ end }}
    }

    location /{{.GarantSubroute}}/ {
{{ if .HasLicense }}
      proxy_pass http://auth/;
      proxy_set_header Host $http_host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_buffering off;
      proxy_request_buffering off;
{{ else }}
      error_page 404 /no_license;
      return 404;
{{ end }}
    }

    location = /no_license {
      root  /no_license/;
    }

    location /load_balancer_status {
      access_log off;
      check_status;
    }

    location /nginx_status {
        # Turn on nginx stats
        stub_status;
        access_log off;
    }
  }

  server {
    listen 443 default ssl;

    ssl_certificate             /config/server.pem;
    ssl_certificate_key         /config/server.pem;

    ssl_session_cache shared:SSL:20m;
    ssl_session_timeout 10m;

    ssl_dhparam                     /root/dhparams.pem;
    ssl_prefer_server_ciphers       on;
    ssl_protocols                   TLSv1.2;
    ssl_ciphers                     ECDH+AESGCM:DH+AESGCM:ECDH+AES256:DH+AES256:ECDH+AES128:DH+AES:ECDH+3DES:DH+3DES:RSA+AESGCM:RSA+AES:RSA+3DES:!aNULL:!MD5:!DSS;

    # This was removed because setting the HSTS header made it difficult to change CAs
    # add_header Strict-Transport-Security "max-age=31536000";
    add_header X-Replica-ID XXXREPLICA_IDXXX always;

    # disable any limits to avoid HTTP 413 for large image uploads
    client_max_body_size 0;

    {{ if .AuthBypassOU }}
        ssl_verify_client optional_no_ca;
        ssl_verify_depth 3;
        ssl_client_certificate /config/auth_bypass_ca.pem;
    {{ end }}

    location / {
      proxy_pass http://api/;
      proxy_set_header Host $http_host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_hide_header {{.RegistryEventsHeaderName}};
      proxy_buffering off;
      proxy_request_buffering off;
    }

    location /ws/ {
{{ if .HasLicense }}
      proxy_pass http://api/ws/;
      proxy_set_header Host $http_host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection "upgrade";
{{ else }}
        return 404;
{{ end }}
    }

    location /v1/search {
      proxy_pass http://api/api/v0/index/dockersearch;
      proxy_set_header Host $http_host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_hide_header {{.RegistryEventsHeaderName}};
      proxy_buffering off;
      proxy_request_buffering off;
    }

    location /v1/ {
      return 404;
    }

    location /v1/users/ {
{{ if .HasLicense }}
{{ else }}
      error_page 404 /no_license;
{{ end }}
      return 404;
    }

    location ~ ^/v2/(.*)/_trust/(.*) {
{{ if .HasLicense }}
        proxy_pass https://notary/v2/$1/_trust/$2;
        proxy_ssl_verify on;
        proxy_ssl_trusted_certificate {{.NotaryCACert}};
        proxy_ssl_name {{ .NotaryHost }};
        proxy_ssl_verify_depth 4;
        proxy_ssl_certificate {{.NotaryClientCert}};
        proxy_ssl_certificate_key {{.NotaryClientKey}};
        proxy_ssl_session_reuse on;
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_buffering off;
        proxy_request_buffering off;
{{ else }}
        error_page 404 /no_license;
        return 404;
{{ end }}
    }

    location /v2/ {
{{ if .HasLicense }}
      proxy_pass http://image_storage/v2/;
      proxy_set_header Host $http_host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_buffering off;
      proxy_request_buffering off;
{{ else }}
      error_page 404 /no_license;
      return 404;
{{ end }}
    }

    location = /no_license {
      root  /no_license/;
    }

    location /{{.GarantSubroute}}/ {
{{ if .HasLicense }}
      proxy_pass http://auth/{{.GarantSubroute}}/;
      proxy_set_header Host $http_host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;

      proxy_set_header X-Client-Cert $ssl_client_cert;
      proxy_set_header X-Client-Cert-Valid $ssl_client_verify;
      proxy_set_header X-Client-DN $ssl_client_s_dn;
      {{ if .AuthBypassOU }}
      set $valid_ou "false";
      if ($ssl_client_s_dn ~ "\bOU={{ .AuthBypassOU }}\b") {
          set $valid_ou "true";
      }
      proxy_set_header X-Client-OU-Valid $valid_ou;
      {{ end }}

      proxy_buffering off;
      proxy_request_buffering off;
{{ else }}
      error_page 404 /no_license;
      return 404;
{{ end }}
    }

    location /load_balancer_status {
      access_log off;
      check_status;
    }

    location /nginx_status {
        # Turn on nginx stats
        stub_status;
        access_log off;
    }
  }
}
`
