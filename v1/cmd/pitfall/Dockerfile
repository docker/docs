FROM golang:alpine

ENV workdir /go/src/github.com/docker/pinata/v1/cmd/pitfall
ADD . ${workdir}
WORKDIR ${workdir}

RUN apk update && apk add git python
RUN go get -d -v
RUN GOOS=darwin go build -v 
RUN cp ./pitfall ./pitfall-darwin
RUN GOOS=linux go build -v 
RUN cp ./pitfall ./pitfall-linux
CMD ./pitfall -h
