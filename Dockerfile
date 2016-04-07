FROM golang:1.6.0

RUN apt-get update && apt-get install -y \
	curl \
	clang \
	libltdl-dev \
	libsqlite3-dev \
	patch \
	tar \
	xz-utils \
	--no-install-recommends \
	&& rm -rf /var/lib/apt/lists/*

RUN go get golang.org/x/tools/cmd/cover

# Configure the container for OSX cross compilation
ENV OSX_SDK MacOSX10.11.sdk
RUN set -x \
	&& export OSXCROSS_PATH="/osxcross" \
	&& git clone --depth 1 https://github.com/tpoechtrager/osxcross.git $OSXCROSS_PATH \
	&& curl -sSL https://s3.dockerproject.org/darwin/${OSX_SDK}.tar.xz -o "${OSXCROSS_PATH}/tarballs/${OSX_SDK}.tar.xz" \
	&& UNATTENDED=yes OSX_VERSION_MIN=10.6 ${OSXCROSS_PATH}/build.sh
ENV PATH /osxcross/target/bin:$PATH

ENV NOTARYDIR /go/src/github.com/docker/notary

COPY . ${NOTARYDIR}

ENV GOPATH ${NOTARYDIR}/Godeps/_workspace:$GOPATH

WORKDIR ${NOTARYDIR}

# Note this cannot use alpine because of the MacOSX Cross SDK: the cctools there uses sys/cdefs.h and that cannot be used in alpine: http://wiki.musl-libc.org/wiki/FAQ#Q:_I.27m_trying_to_compile_something_against_musl_and_I_get_error_messages_about_sys.2Fcdefs.h
