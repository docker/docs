# Rsyslog daemon for receiving TLS syslogs on port 514
#
# docker run --rm -it \
#  -v /home/edgar/mysyslog:/etc/pki/rsyslog \
#  -p 515:514 \
#  --name rsyslog \
#  hinshun/rsyslog

FROM debian:jessie
MAINTAINER Edgar Lee <edgar.lee@docker.com>

# Install dependencies & clean up
RUN apt-get update && DEBIAN_FRONTEND=noninteractive \
    apt-get install -y --no-install-recommends \
    rsyslog \
    rsyslog-gnutls \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY rsyslog.conf /etc/rsyslog.conf

RUN mkdir -p /etc/pki/rsyslog && \
  touch /etc/pki/rsyslog/ca.pem \
  /etc/pki/rsyslog/cert.pem \
  /etc/pki/rsyslog/key.pem

COPY start.sh /start.sh

ENTRYPOINT ["/start.sh"]
