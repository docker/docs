# syntax=docker/dockerfile:1

FROM alpine:3.17.1 as base
WORKDIR /src
RUN apk add --update nodejs npm git gcompat

FROM base as node
COPY package*.json .
RUN npm install

FROM base as hugo
WORKDIR /bin
ARG TARGETARCH
RUN wget https://github.com/gohugoio/hugo/releases/download/v0.110.0/hugo_extended_0.110.0_linux-${TARGETARCH}.tar.gz \
    && tar -xf hugo_extended_0.110.0_linux-${TARGETARCH}.tar.gz hugo

FROM base as build-base
COPY --from=hugo /bin/hugo /usr/local/bin/hugo
COPY --from=node /src/node_modules /src/node_modules

FROM build-base as build
ARG HUGO_ENV
COPY . .
RUN hugo --gc --minify -d /out -e "$HUGO_ENV"

FROM scratch as release
COPY --from=build /out /

FROM davidanson/markdownlint-cli2:v0.6.0 as lint
USER root
RUN --mount=type=bind,target=. \
    /usr/local/bin/markdownlint-cli2 content/**/*.md

FROM wjdp/htmltest:v0.17.0 as test
WORKDIR /test
COPY --from=build /out ./public
ADD .htmltest.yml .htmltest.yml
RUN htmltest
