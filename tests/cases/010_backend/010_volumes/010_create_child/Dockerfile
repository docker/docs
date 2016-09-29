FROM alpine:3.3

RUN apk --update add build-base
COPY create_child.c .
RUN gcc -o create_child create_child.c
