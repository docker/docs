FROM alpine
MAINTAINER Jessica Frazelle <jess@docker.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk update && apk add \
	ca-certificates \
	&& rm -rf /var/cache/apk/*

COPY . /go/src/github.com/docker/opensource

RUN buildDeps=' \
		go \
		git \
		gcc \
		libc-dev \
		libgcc \
	' \
	set -x \
	&& apk update \
	&& apk add $buildDeps \
	&& cd /go/src/github.com/docker/opensource \
	&& go get -d -v github.com/docker/opensource/maintainers \
	&& go generate ./maintainers \
	&& go build -o /usr/bin/maintainers ./maintainers \
	&& apk del $buildDeps \
	&& rm -rf /var/cache/apk/* \
	&& rm -rf /go \
	&& echo "Build complete."

WORKDIR /root/maintainers

ENTRYPOINT [ "maintainers" ]
