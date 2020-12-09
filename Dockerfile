# This Dockerfile builds the docs for https://docs.docker.com/
# from the master branch of https://github.com/docker/docker.github.io
#
# Here is the sequence:
# 1.  Set up base stages for building and deploying
# 2.  Collect and build the reference documentation (from upstream resources)
# 3.  Build static HTML from the current branch
# 4.  Build the final image, combining the reference docs and current version
#     of the documentation
#
# When the image is run, it starts Nginx and serves the docs at port 4000

# Jekyll environment (development/production)
ARG JEKYLL_ENV=development

# Engine
# TODO change to 20.10 branch, once created
ARG ENGINE_BRANCH="master"

# Distribution
ARG DISTRIBUTION_BRANCH="release/2.7"

# Compose CLI
ARG COMPOSE_CLI_BRANCH="main"

###
# Set up base stages for building and deploying
###
FROM starefossen/github-pages:198 AS builderbase
ENV TARGET=/usr/share/nginx/html
WORKDIR /usr/src/app/md_source/

# Set vars used by fetch-upstream-resources.sh script as an environment variable,
# so that they are persisted in the image for use in later stages.
ARG ENGINE_BRANCH
ENV ENGINE_BRANCH=${ENGINE_BRANCH}

ARG DISTRIBUTION_BRANCH
ENV DISTRIBUTION_BRANCH=${DISTRIBUTION_BRANCH}

ARG COMPOSE_CLI_BRANCH
ENV COMPOSE_CLI_BRANCH=${COMPOSE_CLI_BRANCH}

# Fetch upstream resources (reference documentation)
# Only add the files that are needed to build these reference docs, so that these
# docs are only rebuilt if changes were made to ENGINE_BRANCH or DISTRIBUTION_BRANCH.
# Disable caching (docker build --no-cache) to force updating these docs.
FROM alpine AS upstream-resources
RUN apk add --no-cache subversion wget
WORKDIR /usr/src/app/md_source/
COPY ./_scripts/fetch-upstream-resources.sh ./_scripts/
ARG ENGINE_BRANCH
ARG DISTRIBUTION_BRANCH
ARG COMPOSE_CLI_BRANCH
RUN ./_scripts/fetch-upstream-resources.sh .


# Build the static HTML for the current docs.
# After building with jekyll, fix up some links
FROM builderbase AS current
COPY . .
COPY --from=upstream-resources /usr/src/app/md_source/. ./
# substitute the "{site.latest_engine_api_version}" in the title for the latest
# API docs, based on the latest_engine_api_version parameter in _config.yml
RUN ./_scripts/update-api-toc.sh
ARG JEKYLL_ENV
RUN echo "Building docs for ${JEKYLL_ENV} environment"
RUN set -eu; \
 if [ "${JEKYLL_ENV}" = "production" ]; then \
    jekyll build --profile -d ${TARGET} --config _config.yml,_config_production.yml; \
    sed -i 's#<loc>/#<loc>https://docs.docker.com/#' "${TARGET}/sitemap.xml"; \
 else \
    jekyll build --profile -d ${TARGET}; \
    echo '[]' > ${TARGET}/js/metadata.json; \
 fi; \
 find ${TARGET} -type f -name '*.html' | while read i; do sed -i 's#\(<a[^>]* href="\)https://docs.docker.com/#\1/#g' "$i"; done;


# This stage only contains the generated files. It can be used to host the
# documentation on a non-containerised service (e.g. to deploy to an s3 bucket).
# When using BuildKit, use the '--output' option to build the files and to copy
# them to your local filesystem.
#
# DOCKER_BUILDKIT=1 docker build --target=deploy-source --output=./_site .
FROM scratch AS deploy-source
COPY --from=current /usr/share/nginx/html /

# Final stage, which includes nginx, and the current docs.
#
# To build current docs:
# DOCKER_BUILDKIT=1 docker build -t docs .
FROM nginx:alpine AS deploy
ENV TARGET=/usr/share/nginx/html
WORKDIR $TARGET

COPY --from=current  /usr/share/nginx/html .

# Configure NGINX
COPY _deploy/nginx/default.conf /etc/nginx/conf.d/default.conf
ARG JEKYLL_ENV
ENV JEKYLL_ENV=${JEKYLL_ENV}
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000 (build target: ${JEKYLL_ENV})"; exec nginx -g 'daemon off;'
