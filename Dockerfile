FROM golang

COPY . /go/src/github.com/docker/notary

ENV GOPATH /go/src/github.com/docker/notary/Godeps/_workspace:$GOPATH

RUN go install github.com/docker/notary/cmd/notary-server

EXPOSE 4443

WORKDIR /go/src/github.com/docker/notary

CMD [ "notary-server", "-config", "cmd/notary-server/dev-config.json" ]
