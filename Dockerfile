FROM golang

COPY . /go/src/github.com/docker/vetinari

RUN GOPATH=/go/:/go/src/github.com/docker/vetinari/Godeps/_workspace go install github.com/docker/vetinari/cmd/vetinari-server

EXPOSE 4444

WORKDIR /go/src/github.com/docker/vetinari

CMD vetinari-server -config cmd/vetinari-server/dev-config.json
