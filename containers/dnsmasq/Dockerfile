FROM debian:jessie

MAINTAINER Chris Biscardi <chris@docker.com>

RUN apt-get update && apt-get install -y dnsmasq

EXPOSE 53/udp

ADD ./run /opt/run

CMD "/opt/run"
# docker run -d -p 53:53/udp --name docker-dnsmasq dnsmasq --address=/dev.docker.io/172.16.200.100
