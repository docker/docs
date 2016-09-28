FROM golang:1.5.0

ENV CONFIG_PATH /etc/garant/config.yml
ENV GARANT_DIR /go/src/github.com/docker/garant
ENV GO15VENDOREXPERIMENT 1

WORKDIR $GARANT_DIR
COPY . $GARANT_DIR
RUN make PREFIX=/go clean binary
RUN mkdir -pv "$(dirname $CONFIG_PATH)"
RUN cp -lv ./cmd/garant/config.yml ./cmd/garant/signing_key.json "$(dirname $CONFIG_PATH)"

EXPOSE 8080
CMD garant $CONFIG_PATH
