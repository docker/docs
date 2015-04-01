FROM golang

COPY . /go/src/github.com/docker/vetinari

RUN go get github.com/docker/vetinari/cmd/vetinari-server

EXPOSE 4443

CMD vetinari-server -cert /go/src/github.com/docker/vetinari/fixtures/ca.pem -key /go/src/github.com/docker/vetinari/fixtures/ca-key.pem -debug
