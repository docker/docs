# syntax=docker/dockerfile:1
# check=skip=InvalidBaseImagePlatform

ARG ALPINE_VERSION=3.23
ARG GO_VERSION=1.26
ARG HTMLTEST_VERSION=0.17.0
ARG VALE_VERSION=3.17.0
ARG HUGO_VERSION=0.163.0
ARG NODE_VERSION=24
ARG PAGEFIND_VERSION=1.5.2
ARG TARGETARCH

# base defines the generic base stage
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
RUN apk add --no-cache \
    git \
    nodejs \
    npm \
    gcompat \
    rsync

# npm downloads Node.js dependencies
FROM base AS npm
ENV NODE_ENV="production"
WORKDIR /out
RUN --mount=source=package.json,target=package.json \
    --mount=source=package-lock.json,target=package-lock.json \
    --mount=type=cache,target=/root/.npm \
    npm ci

# hugo downloads the Hugo binary
FROM base AS hugo
ARG TARGETARCH
ARG HUGO_VERSION
WORKDIR /out
ADD https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-${TARGETARCH}.tar.gz .
RUN tar xvf hugo_extended_${HUGO_VERSION}_linux-${TARGETARCH}.tar.gz

# build-base is the base stage used for building the site
FROM base AS build-base
WORKDIR /project
COPY --from=hugo /out/hugo /bin/hugo
COPY --from=npm /out/node_modules node_modules
COPY . .

# build creates production builds with Hugo
FROM build-base AS build
# HUGO_ENV sets the hugo.Environment (production, development, preview)
ARG HUGO_ENV="development"
# DOCS_URL sets the base URL for the site
ARG DOCS_URL="https://docs.docker.com"
ENV HUGO_CACHEDIR="/tmp/hugo_cache"
RUN --mount=type=cache,target=/tmp/hugo_cache \
    hugo \
      --gc \
      --minify \
      --panicOnWarning \
      --printPathWarnings \
      --printUnusedTemplates \
      -b $DOCS_URL \
      -e $HUGO_ENV
RUN node --test hack/test/flatten-and-resolve.mjs
RUN node hack/flatten-and-resolve.js public
RUN node hack/api-docs/verify-output.mjs public

# lint lints markdown files
FROM ghcr.io/rvben/rumdl:0.2.49-alpine AS lint
RUN --mount=type=bind,target=. rumdl check content

# vacuum downloads the prebuilt validator for the target architecture
FROM scratch AS vacuum-amd64
ADD --unpack --checksum=sha256:973b8ed30cd36533da4cdb76729e833324a3043bad48c74f61b6a1cc3df40b78 \
    https://github.com/daveshanley/vacuum/releases/download/v0.30.3/vacuum_0.30.3_linux_x86_64.tar.gz /out/

FROM scratch AS vacuum-arm64
ADD --unpack --checksum=sha256:04b604df0b1b570d5abd58b577c816122da162eb8ba12a7f70c3060aec12b3a9 \
    https://github.com/daveshanley/vacuum/releases/download/v0.30.3/vacuum_0.30.3_linux_arm64.tar.gz /out/

FROM vacuum-${TARGETARCH} AS vacuum

# validate-api-reference checks the vendored presentation data
FROM base AS validate-api-reference
RUN apk add --no-cache bash
WORKDIR /project
COPY --from=vacuum /out/vacuum /usr/local/bin/vacuum
COPY hack/api-docs ./hack/api-docs
COPY content/reference/api ./content/reference/api
COPY data/api-reference.json ./data/api-reference.json
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build <<"EOT"
set -eu
cp data/api-reference.json /tmp/api-reference-committed.json
./hack/api-docs/run.sh test
./hack/api-docs/run.sh generate
if ! cmp -s /tmp/api-reference-committed.json data/api-reference.json; then
  echo >&2 'ERROR: API reference data is stale. Run ./hack/api-docs/run.sh generate and commit data/api-reference.json.'
  exit 1
fi
EOT

# test validates HTML output and checks for broken links
FROM wjdp/htmltest:v${HTMLTEST_VERSION} AS test
WORKDIR /test
COPY --from=build /project/public ./public
ADD .htmltest.yml .htmltest.yml
RUN htmltest

# vale
FROM jdkato/vale:v${VALE_VERSION} AS vale-run
WORKDIR /src
ARG GITHUB_ACTIONS
RUN --mount=type=bind,target=.,rw <<EOT
  set -e
  mkdir /out
  args=""
  [ "$GITHUB_ACTIONS" = "true" ] && args="--output=.vale-rdjsonl.tmpl"
  set -x
  vale sync
  vale $args content/ | tee /out/vale.out
EOT

FROM scratch AS vale
COPY --from=vale-run /out/vale.out /

# update-modules downloads and vendors Hugo modules
FROM build-base AS update-modules
# MODULE is the Go module path and version of the module to update
ARG MODULE
RUN <<"EOT"
set -ex
if [ -n "$MODULE" ]; then
    hugo mod get ${MODULE}
else
    echo "no module set";
fi
EOT
RUN hugo mod vendor

# vendor is an empty stage with only vendored Hugo modules
FROM scratch AS vendor
COPY --from=update-modules /project/_vendor /_vendor
COPY --from=update-modules /project/go.* /

FROM base AS validate-vendor
RUN --mount=target=/context \
  --mount=type=bind,from=vendor,target=/out \
  --mount=target=.,type=tmpfs <<EOT
set -e
rsync -a /context/. .
git add -A
rm -rf _vendor
cp -rf /out/* .
if [ -n "$(git status --porcelain -- go.mod go.sum _vendor)" ]; then
  echo >&2 'ERROR: Vendor result differs. Please vendor your package with "make vendor"'
  git status --porcelain -- go.mod go.sum _vendor
  exit 1
fi
EOT

# build-upstream builds an upstream project with a replacement module
FROM build-base AS build-upstream
# UPSTREAM_MODULE_NAME is the canonical upstream repository name and namespace (e.g. moby/buildkit)
ARG UPSTREAM_MODULE_NAME
# UPSTREAM_REPO is the repository of the project to validate (e.g. dvdksn/buildkit)
ARG UPSTREAM_REPO
# UPSTREAM_COMMIT is the commit hash of the upstream project to validate
ARG UPSTREAM_COMMIT
# HUGO_MODULE_REPLACEMENTS is the replacement module for the upstream project
ENV HUGO_MODULE_REPLACEMENTS="github.com/${UPSTREAM_MODULE_NAME} -> github.com/${UPSTREAM_REPO} ${UPSTREAM_COMMIT}"
RUN hugo --ignoreVendorPaths "github.com/${UPSTREAM_MODULE_NAME}"

# validate-upstream validates HTML output for upstream builds
FROM wjdp/htmltest:v${HTMLTEST_VERSION} AS validate-upstream
WORKDIR /test
COPY --from=build-upstream /project/public ./public
ADD .htmltest.yml .htmltest.yml
RUN htmltest

# unused-media checks for unused graphics and other media
FROM alpine:${ALPINE_VERSION} AS unused-media
RUN apk add --no-cache fd ripgrep
WORKDIR /test
RUN --mount=type=bind,target=. ./hack/test/unused_media

# pagefind installs the Pagefind runtime
FROM base AS pagefind
ARG PAGEFIND_VERSION
COPY --from=build /project/public ./public
RUN --mount=type=bind,src=pagefind.yml,target=pagefind.yml \
    npx pagefind@v${PAGEFIND_VERSION} --output-path "/pagefind"

# index generates a Pagefind index
FROM scratch AS index
COPY --from=pagefind /pagefind .

# test-go-redirects checks that the /go/ redirects are valid
FROM alpine:${ALPINE_VERSION} AS test-go-redirects
WORKDIR /work
RUN apk add yq
COPY --from=build /project/public ./public
RUN --mount=type=bind,target=. <<"EOT"
set -ex
./hack/test/go_redirects
EOT

# release is an empty scratch image with only compiled assets
FROM scratch AS release
COPY --from=build /project/public /
COPY --from=pagefind /pagefind /pagefind
