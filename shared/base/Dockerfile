# base Dockerfile for all of DTR's internal images:
# (bootstrapper, api, registry, nginx)
FROM alpine:3.3

# needed for the build system
ENTRYPOINT sh

RUN apk add -uU curl python3 pcre ca-certificates

RUN curl -L https://github.com/kelseyhightower/confd/releases/download/v0.11.0/confd-0.11.0-linux-amd64 > /usr/bin/confd && chmod +x /usr/bin/confd && mkdir -p /etc/confd

# put nginx on the path
RUN ln -s /usr/local/nginx/sbin/nginx /bin/nginx

# used by confd for writing out configs
RUN mkdir -p /config
# this is used by nginx to return a better error when no license is set
RUN mkdir -p /no_license
# used by adminserver
RUN mkdir -p /ui /swagger

# fuck the NSA
ADD dhparams.pem /root/dhparams.pem

EXPOSE 80 443 3000
