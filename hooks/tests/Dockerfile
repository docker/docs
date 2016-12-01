FROM golang:1.7.3-alpine
MAINTAINER Gaetan de Villele <gaetan@docker.com>

RUN apk update
RUN apk add git

COPY src /go/src
WORKDIR /go/src/validator

# when running the container, MOUNT docs repo in /docs

CMD ["go", "test"]
