FROM alpine
MAINTAINER Darin London <darin.london@duke.edu>

WORKDIR /usr/local/bin
RUN apk update \
  && apk add openssl \
  && wget https://releases.hashicorp.com/consul/0.6.4/consul_0.6.4_linux_amd64.zip \
  && unzip consul_0.6.4_linux_amd64.zip \
  && rm consul_0.6.4_linux_amd64.zip
RUN mkdir -p /etc/consul/ui
WORKDIR /etc/consul/ui
RUN wget https://releases.hashicorp.com/consul/0.6.4/consul_0.6.4_web_ui.zip \
  && unzip consul_0.6.4_web_ui.zip \
  && rm consul_0.6.4_web_ui.zip
EXPOSE 8300 8301 8301/udp 8302 8302/udp 8400 8500 53/udp
WORKDIR /root
