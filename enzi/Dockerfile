FROM alpine:3.4

RUN apk --no-cache add ca-certificates

COPY enzi /usr/local/bin
COPY ui /ui

# Establish the version and upgrade matrix
ARG UCP_VERSION
LABEL com.docker.ucp.version="$UCP_VERSION"

USER nobody
ENTRYPOINT ["/usr/local/bin/enzi"]
