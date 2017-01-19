FROM golang:1.7.3-alpine
MAINTAINER Adrien Duermael <adrien@docker.com>

COPY src /go/src
WORKDIR /go/src/validator

# when running the container, MOUNT docs repo in /docs
CMD ["go", "test"]
