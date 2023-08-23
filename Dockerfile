# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine as base
WORKDIR /src
RUN apk --update add nodejs npm git

FROM base as node
COPY package*.json .
RUN npm install && npm cache clean --force

FROM base as hugo
ARG HUGO_VERSION=0.116.1
ARG TARGETARCH
WORKDIR /bin
RUN go install github.com/gohugoio/hugo@v${HUGO_VERSION}

FROM base as build-base
COPY --from=hugo $GOPATH/bin/hugo /bin/hugo
COPY --from=node /src/node_modules /src/node_modules

FROM build-base as dev
COPY . .

FROM build-base as build
ARG HUGO_ENV
ARG DOCS_URL
COPY . .
RUN hugo --gc --minify -d /out -e $HUGO_ENV -b $DOCS_URL

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

FROM build-base as update-modules
ARG MODULE="-u"
WORKDIR /src
COPY . .
RUN hugo mod get ${MODULE}
RUN hugo mod vendor

FROM scratch as vendor
COPY --from=update-modules /src/_vendor /_vendor
COPY --from=update-modules /src/go.* /
