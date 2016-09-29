FROM alpine:latest
RUN apk update && apk add alpine-sdk
WORKDIR /src
COPY server.c common.h vm_sockets.h /src/
RUN cc -Wall -O3 -o server server.c
ENTRYPOINT ["/src/server"]
