FROM alpine:3.3

RUN apk --update add gcc libc-dev
COPY chown_test.c chown_test.sh /
RUN gcc -o chown_test chown_test.c
