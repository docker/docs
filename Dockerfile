# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21
ARG HTMLTEST_VERSION=0.17.0
ARG XX_VERSION=1.4.0
ARG ALPINE_VERSION=3.19

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as base
WORKDIR /src
RUN apk --update add nodejs npm git gcompat

FROM base as node
COPY package*.json .
ENV NODE_ENV=production
RUN npm install

FROM base as hugo
ARG HUGO_VERSION=0.122.0
ARG TARGETARCH
WORKDIR /tmp/hugo
RUN wget -O "hugo.tar.gz" "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-${TARGETARCH}.tar.gz"
RUN tar -xf "hugo.tar.gz" hugo

FROM base as build-base
COPY --from=hugo /tmp/hugo/hugo /bin/hugo
COPY --from=node /src/node_modules /src/node_modules
COPY . .

FROM build-base as dev

FROM build-base as build
ARG HUGO_ENV
ARG DOCS_URL
RUN hugo --gc --minify -d /out -e $HUGO_ENV -b $DOCS_URL

FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM alpine:${ALPINE_VERSION} as pagefind
ARG TARGETPLATFORM
ARG PAGEFIND_VERSION=1.0.4
COPY --from=xx / /
WORKDIR /tmp/pagefind
RUN wget -O "pagefind.tar.gz" \
    "https://github.com/CloudCannon/pagefind/releases/download/v${PAGEFIND_VERSION}/pagefind-v${PAGEFIND_VERSION}-$(xx-info alpine-arch)-unknown-linux-musl.tar.gz"
RUN tar -xf "pagefind.tar.gz" pagefind
COPY --from=build /out ./public
RUN ./pagefind --site "public" --output-path /pagefind

FROM scratch AS index
COPY --from=pagefind /pagefind .

FROM scratch as release
COPY --from=build /out /
COPY --from=pagefind /pagefind .

FROM davidanson/markdownlint-cli2:v0.12.1 as lint
USER root
RUN --mount=type=bind,target=. \
    /usr/local/bin/markdownlint-cli2 \
    "content/**/*.md" \
    "#content/engine/release-notes/*.md" \
    "#content/desktop/previous-versions/*.md"

FROM wjdp/htmltest:v${HTMLTEST_VERSION} as test
WORKDIR /test
COPY --from=build /out ./public
ADD .htmltest.yml .htmltest.yml
RUN htmltest

FROM build-base as update-modules
ARG MODULE="-u"
RUN hugo mod get ${MODULE}
RUN hugo mod vendor

FROM scratch as vendor
COPY --from=update-modules /src/_vendor /_vendor
COPY --from=update-modules /src/go.* /

FROM build-base as build-upstream
ARG UPSTREAM_MODULE_NAME
ARG UPSTREAM_REPO
ARG UPSTREAM_COMMIT
ENV HUGO_MODULE_REPLACEMENTS="github.com/${UPSTREAM_MODULE_NAME} -> github.com/${UPSTREAM_REPO} ${UPSTREAM_COMMIT}"
RUN hugo --ignoreVendorPaths "github.com/${UPSTREAM_MODULE_NAME}" -d /out

FROM wjdp/htmltest:v${HTMLTEST_VERSION} as validate-upstream
WORKDIR /test
COPY --from=build-upstream /out ./public
ADD .htmltest.yml .htmltest.yml
RUN htmltest

FROM dev
