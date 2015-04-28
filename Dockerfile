FROM golang

COPY . /go/src/github.com/docker/vetinari

RUN GOPATH=/go/:/go/src/github.com/docker/vetinari/Godeps/_workspace go install github.com/docker/vetinari/cmd/vetinari-server

EXPOSE 4443

CMD vetinari-server -cert /go/src/github.com/docker/vetinari/fixtures/vetinari.pem -key /go/src/github.com/docker/vetinari/fixtures/vetinary.key -debug
