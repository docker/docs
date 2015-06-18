FROM golang

COPY . /go/src/github.com/docker/vetinari

RUN chmod 777 /tmp/

RUN GOPATH="/go/src/github.com/docker/vetinari/Godeps/_workspace:/go/" go install github.com/docker/vetinari/cmd/vetinari-server

EXPOSE 4443

WORKDIR /go/src/github.com/docker/vetinari

CMD vetinari-server -config cmd/vetinari-server/dev-config.json
