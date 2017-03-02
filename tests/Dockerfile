FROM golang:1.7.3-alpine

COPY src /go/src
WORKDIR /go/src/validator

# when running the container, MOUNT docs repo in /docs
CMD ["go", "test", "-v", "-run", "FrontMatter"]
