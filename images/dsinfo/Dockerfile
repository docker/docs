FROM alpine:3.4

ENV DOCKER_VERSION 1.10.3
ENV TAG latest

RUN apk update && \
    apk fetch procps sysstat dmidecode curl && \
    apk add procps sysstat dmidecode curl && \
    curl -L https://get.docker.com/builds/Linux/x86_64/docker-${DOCKER_VERSION}  > /usr/local/bin/docker && \
    chmod +x /usr/local/bin/docker && \
    curl -L https://raw.githubusercontent.com/docker/docker/master/contrib/check-config.sh > /check-config.sh && \
    chmod +x /check-config.sh

ADD dsinfo.sh /dsinfo.sh
RUN chmod +x /dsinfo.sh


CMD /dsinfo.sh -b 2>/dev/null
