# syntax=docker/dockerfile:1
# check=skip=InvalidBaseImagePlatform

# ALPINE_VERSION sets the Alpine Linux version for all Alpine stages
ARG ALPINE_VERSION=3.20
# GO_VERSION sets the Go version for the base stage
ARG GO_VERSION=1.23
# HTMLTEST_VERSION sets the wjdp/htmltest version for HTML testing
ARG HTMLTEST_VERSION=0.17.0
# HUGO_VERSION sets the version of Hugo to build the site with
ARG HUGO_VERSION=0.136.2

# build-base is the base stage used for building the site
FROM ghcr.io/gohugoio/hugo:v${HUGO_VERSION} AS build-base
USER root
ENV NODE_ENV="production"
RUN --mount=source=package.json,target=package.json \
    --mount=source=package-lock.json,target=package-lock.json \
    --mount=type=cache,target=/root/.npm \
    npm ci
COPY . .

# build creates production builds with Hugo
FROM build-base AS build
# HUGO_ENV sets the hugo.Environment (production, development, preview)
ARG HUGO_ENV="development"
# DOCS_URL sets the base URL for the site
ARG DOCS_URL="https://docs.docker.com"
RUN --mount=type=cache,target=/cache \
    hugo --gc --minify -d /out -e $HUGO_ENV -b $DOCS_URL

# lint lints markdown files
FROM davidanson/markdownlint-cli2:v0.14.0 AS lint
USER root
RUN --mount=type=bind,target=. \
    /usr/local/bin/markdownlint-cli2 \
    "content/**/*.md" \
    "#content/manuals/engine/release-notes/*.md" \
    "#content/manuals/desktop/previous-versions/*.md"

# test validates HTML output and checks for broken links
FROM wjdp/htmltest:v${HTMLTEST_VERSION} AS test
WORKDIR /test
COPY --from=build /out ./public
ADD .htmltest.yml .htmltest.yml
RUN htmltest

# update-modules downloads and vendors Hugo modules
FROM build-base AS update-modules
# MODULE is the Go module path and version of the module to update
ARG MODULE
RUN <<"EOT"
set -ex
if [ -n "$MODULE" ]; then
    hugo mod get ${MODULE}
    RESOLVED=$(cat go.mod | grep -m 1 "${MODULE/@*/}" | awk '{print $1 "@" $2}')
    go mod edit -replace "${MODULE/@*/}=${RESOLVED}";
else
    echo "no module set";
fi
EOT
RUN hugo mod vendor

# vendor is an empty stage with only vendored Hugo modules
FROM scratch AS vendor
COPY --from=update-modules /src/_vendor /_vendor
COPY --from=update-modules /src/go.* /

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
RUN hugo --ignoreVendorPaths "github.com/${UPSTREAM_MODULE_NAME}" -d /out

# validate-upstream validates HTML output for upstream builds
FROM wjdp/htmltest:v${HTMLTEST_VERSION} AS validate-upstream
WORKDIR /test
COPY --from=build-upstream /out ./public
ADD .htmltest.yml .htmltest.yml
RUN htmltest

# unused-media checks for unused graphics and other media
FROM alpine:${ALPINE_VERSION} AS unused-media
RUN apk add --no-cache fd ripgrep
WORKDIR /test
RUN --mount=type=bind,target=. <<"EOT"
set -ex
./scripts/test_unused_media.sh
EOT

# path-warnings checks for duplicate target paths
FROM build-base AS path-warnings
RUN hugo --printPathWarnings > /path-warnings.txt
RUN <<EOT
DUPLICATE_TARGETS=$(grep "Duplicate target paths" /path-warnings.txt)
if [ ! -z "$DUPLICATE_TARGETS" ]; then
    echo "$DUPLICATE_TARGETS"
    echo "You probably have a duplicate alias defined. Please check your aliases."
    exit 1
fi
EOT

# pagefind installs the Pagefind runtime
FROM node:alpine${ALPINE_VERSION} AS pagefind
ARG PAGEFIND_VERSION=1.1.1
COPY --from=build /out ./public
RUN --mount=type=bind,src=pagefind.yml,target=pagefind.yml \
    npx pagefind@v${PAGEFIND_VERSION} --output-path "/pagefind"

# index generates a Pagefind index
FROM scratch AS index
COPY --from=pagefind /pagefind .

# test-go-redirects checks that the /go/ redirects are valid
FROM alpine:${ALPINE_VERSION} AS test-go-redirects
WORKDIR /work
RUN apk add yq
COPY --from=build /out ./public
RUN --mount=type=bind,target=. <<"EOT"
set -ex
./scripts/test_go_redirects.sh
EOT

# release is an empty scratch image with only compiled assets
FROM scratch AS release
COPY --from=build /out /
COPY --from=pagefind /pagefind /pagefind
