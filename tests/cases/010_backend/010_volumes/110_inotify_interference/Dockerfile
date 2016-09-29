FROM alpine:3.4

RUN apk update && apk add build-base
ADD interfere.c /
RUN gcc -Wall -Werror -o /interfere interfere.c
