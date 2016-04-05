FROM golang:1.6.0-alpine
MAINTAINER David Lawrence "david.lawrence@docker.com"

RUN apk add --update git gcc libc-dev && rm -rf /var/cache/apk/*

# Note that -ldflags "-extldflags -fno-PIC" and the "-extldflags -fno-PIC" below
# are to get around https://github.com/golang/go/issues/14851#issuecomment-200915770
# Once the golang alpine image updates, since it's been patched, we can
# remove the flags

# Install DB migration tool
RUN go get -ldflags "-extldflags -fno-PIC" github.com/mattes/migrate

ENV NOTARYPKG github.com/docker/notary

# Copy the local repo to the expected go path
COPY . /go/src/${NOTARYPKG}

WORKDIR /go/src/${NOTARYPKG}

ENV NOTARY_SIGNER_DEFAULT_ALIAS="timestamp_1"
ENV NOTARY_SIGNER_TIMESTAMP_1="testpassword"

EXPOSE 4444

# Install notary-signer
RUN go install \
    -tags pkcs11 \
    -ldflags "-w -X ${NOTARYPKG}/version.GitCommit=`git rev-parse --short HEAD` -X ${NOTARYPKG}/version.NotaryVersion=`cat NOTARY_VERSION` -extldflags -fno-PIC" \
    ${NOTARYPKG}/cmd/notary-signer && apk del git gcc libc-dev

ENTRYPOINT [ "notary-signer" ]
CMD [ "-config=fixtures/signer-config-local.json" ]
