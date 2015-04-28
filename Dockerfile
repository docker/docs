FROM golang

COPY . /go/src/github.com/docker/vetinari

RUN GOPATH=/go/:/go/src/github.com/docker/vetinari/Godeps/_workspace go install github.com/docker/vetinari/cmd/vetinari-server

EXPOSE 4444

CMD vetinari-server -config /go/src/github.com/docker/vetinari/cmd/vetinari-server/dev-config.json
