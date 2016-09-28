FROM alpine:3.4

MAINTAINER Daniel Hiltgen "<daniel.hiltgen@docker.com>"

ENV VERSION 2.3.6
EXPOSE 2380 7001 2379 4001
ENTRYPOINT ["/bin/etcd"]
CMD ["-data-dir=/data"]
ENV ETCD_COMMIT v2.3.6

# Establish the version and upgrade matrix
LABEL com.docker.ucp.version="$VERSION" com.docker.ucp.upgrades_from="2.2.0,0.0.0,2.2.3,2.2.4,2.2.5,2.3.6"

RUN apk update && \
    apk add ca-certificates go git gcc libc-dev libltdl libtool libgcc libpcap-dev && \
    export GOPATH=/go && \
    git clone https://github.com/coreos/etcd /go/src/github.com/coreos/etcd && \
    (cd /go/src/github.com/coreos/etcd && git checkout $ETCD_COMMIT && ./build && mv bin/etcd* /bin/ ) && \
    apk del go git gcc libc-dev libtool libgcc libpcap-dev && \
    rm -rf /var/cache/apk/* \
    rm -rf /go
USER nobody
