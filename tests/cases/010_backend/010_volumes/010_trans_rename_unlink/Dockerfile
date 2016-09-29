FROM alpine:3.3

RUN apk --update add build-base
COPY trans_rename_unlink.c .
RUN gcc -o trans_rename_unlink trans_rename_unlink.c
