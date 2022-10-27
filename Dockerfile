# syntax=docker/dockerfile:1

# This Dockerfile builds the docs for https://docs.docker.com/
# from the main branch of https://github.com/docker/docs

# Use same ruby version as the one in .ruby-version
# that is used by Netlify
ARG RUBY_VERSION=2.7.6
# Same as the one in Gemfile.lock
ARG BUNDLER_VERSION=2.3.13

ARG JEKYLL_ENV=development
ARG DOCS_URL=http://localhost:4000

# Base stage for building
FROM ruby:${RUBY_VERSION}-alpine AS base
WORKDIR /src
RUN apk add --no-cache bash build-base git

# Gem stage will install bundler used as dependency manager
# for our dependencies in Gemfile for Jekyll
FROM base AS gem
ARG BUNDLER_VERSION
COPY Gemfile* .
RUN gem uninstall -aIx bundler \
  && gem install bundler -v ${BUNDLER_VERSION} \
  && bundle install --jobs 4 --retry 3

# Vendor Gemfile for Jekyll
FROM gem AS vendored
ARG BUNDLER_VERSION
RUN bundle update \
  && mkdir /out \
  && cp Gemfile.lock /out

# Stage used to output the vendored Gemfile.lock:
# > make vendor
# or
# > docker buildx bake vendor
FROM scratch AS vendor
COPY --from=vendored /out /

# Build the static HTML for the current docs.
# After building with jekyll, fix up some links
FROM gem AS generate
ARG JEKYLL_ENV
ARG DOCS_URL
ENV TARGET=/out
RUN --mount=type=bind,target=.,rw \
    --mount=type=cache,target=/tmp/docker-docs-clone \
    --mount=type=cache,target=/src/.jekyll-cache <<EOT
  set -eu
  CONFIG_FILES="_config.yml"
  if [ "${JEKYLL_ENV}" = "production" ]; then
    CONFIG_FILES="${CONFIG_FILES},_config_production.yml"
  elif [ "${DOCS_URL}" = "https://docs-stage.docker.com" ]; then
    CONFIG_FILES="${CONFIG_FILES},_config_stage.yml"
  fi
  set -x
  bundle exec jekyll build --profile -d ${TARGET} --config ${CONFIG_FILES}
EOT

# htmlproofer checks for broken links
FROM gem AS htmlproofer-base
RUN --mount=type=bind,from=generate,source=/out,target=_site <<EOF
  htmlproofer ./_site \
    --disable-external \
    --internal-domains="docs.docker.com,docs-stage.docker.com,localhost:4000" \
    --file-ignore="/^./_site/engine/api/.*$/,./_site/registry/configuration/index.html" \
    --url-ignore="/^/docker-hub/api/latest/.*$/,/^/engine/api/v.+/#.*$/,/^/glossary/.*$/" > /results 2>&1
  rc=$?
  if [[ $rc -eq 0 ]]; then
    echo -n > /results
  fi
EOF

FROM htmlproofer-base as htmlproofer
RUN <<EOF
  cat /results
  [ ! -s /results ] || exit 1
EOF

FROM scratch as htmlproofer-output
COPY --from=htmlproofer-base /results /results

# mdl is a lint tool for markdown files
FROM gem AS mdl-base
ARG MDL_JSON
ARG MDL_STYLE
RUN --mount=type=bind,target=. <<EOF
  mdl --ignore-front-matter ${MDL_JSON:+'--json'} --style=${MDL_STYLE:-'.markdownlint.rb'} $( \
    find '.' -name '*.md' \
      -not -path './registry/*' \
      -not -path './desktop/extensions-sdk/*' \
  ) > /results
  rc=$?
  if [[ $rc -eq 0 ]]; then
    echo -n > /results
  fi
EOF

FROM mdl-base as mdl
RUN <<EOF
  cat /results
  [ ! -s /results ] || exit 1
EOF

FROM scratch as mdl-output
COPY --from=mdl-base /results /results

# Release the generated files in a scratch image
# Can be output to your host with:
# > make release
# or
# > docker buildx bake release
FROM scratch AS release
COPY --from=generate /out /

# Create a runnable nginx instance with generated HTML files.
# When the image is run, it starts Nginx and serves the docs at port 4000:
# > make deploy
# or
# > docker-compose up --build
FROM nginx:alpine AS deploy
COPY --from=release / /usr/share/nginx/html
COPY _deploy/nginx/default.conf /etc/nginx/conf.d/default.conf
ARG JEKYLL_ENV
ENV JEKYLL_ENV=${JEKYLL_ENV}
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000 (build target: ${JEKYLL_ENV})"; exec nginx -g 'daemon off;'

FROM deploy
