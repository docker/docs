FROM golang:1.4.2-cross

RUN apt-get update && apt-get install -y \
	libltdl-dev \
	--no-install-recommends \
	&& rm -rf /var/lib/apt/lists/*

RUN go get golang.org/x/tools/cmd/vet \
	&& go get golang.org/x/tools/cmd/cover

COPY . /go/src/github.com/docker/notary

ENV GOPATH /go/src/github.com/docker/notary/Godeps/_workspace:$GOPATH

WORKDIR /go/src/github.com/docker/notary
