# This Dockerfile builds the docs for https://docs.docker.com/
# from the master branch of https://github.com/docker/docker.github.io
#
# Here is the sequence:
# 1.  Set up base stages for building and deploying
# 2.  Collect and build the archived documentation
# 3.  Collect and build the reference documentation (from upstream resources)
# 4.  Build static HTML from the current branch
# 5.  Build the final image, combining the archives, reference docs, and
#     current version of the documentation
#
# When the image is run, it starts Nginx and serves the docs at port 4000


# Engine
ARG ENGINE_BRANCH="19.03"

# Distribution
ARG DISTRIBUTION_BRANCH="release/2.7"

# Set to "false" to build the documentation without archives
ARG ENABLE_ARCHIVES=true

###
# Set up base stages for building and deploying
###

# Get basic configs and Jekyll env
FROM docs/docker.github.io:docs-builder AS builderbase
ENV TARGET=/usr/share/nginx/html
WORKDIR /usr/src/app/md_source/

# Set vars used by fetch-upstream-resources.sh script as an environment variable,
# so that they are persisted in the image for use in later stages.
ARG ENGINE_BRANCH
ENV ENGINE_BRANCH=${ENGINE_BRANCH}

ARG DISTRIBUTION_BRANCH
ENV DISTRIBUTION_BRANCH=${DISTRIBUTION_BRANCH}


# Empty stage if archives are disabled (ENABLE_ARCHIVES=false)
FROM scratch AS archives-false

# Stage with static HTML for all archives (ENABLE_ARCHIVES=true)
FROM scratch AS archives-true
ENV TARGET=/usr/share/nginx/html
# To add a new archive, add it here and ALSO edit _data/docsarchive/archives.yaml
# to add it to the drop-down
COPY --from=docs/docker.github.io:v17.06 ${TARGET} /
COPY --from=docs/docker.github.io:v18.03 ${TARGET} /
COPY --from=docs/docker.github.io:v18.09 ${TARGET} /

# Stage either with, or without archives, depending on ENABLE_ARCHIVES
FROM archives-${ENABLE_ARCHIVES} AS archives

# Fetch upstream resources (reference documentation)
# Only add the files that are needed to build these reference docs, so that these
# docs are only rebuilt if changes were made to ENGINE_BRANCH or DISTRIBUTION_BRANCH.
# Disable caching (docker build --no-cache) to force updating these docs.
FROM builderbase AS upstream-resources
COPY ./_scripts/fetch-upstream-resources.sh ./_scripts/
RUN ./_scripts/fetch-upstream-resources.sh .


# Build the static HTML for the current docs.
# After building with jekyll, fix up some links, but don't touch the archives
FROM builderbase AS current
COPY . .
COPY --from=upstream-resources /usr/src/app/md_source/. ./
# substitute the "{site.latest_engine_api_version}" in the title for the latest
# API docs, based on the latest_engine_api_version parameter in _config.yml
RUN ./_scripts/update-api-toc.sh
RUN jekyll build -d ${TARGET}
RUN find ${TARGET} -type f -name '*.html' | grep -vE "v[0-9]+\." | while read i; do sed -i 's#href="https://docs.docker.com/#href="/#g' "$i"; done


# This stage only contains the generated files. It can be used to host the
# documentation on a non-containerised service (e.g. to deploy to an s3 bucket).
# When using BuildKit, use the '--output' option to build the files and to copy
# them to your local filesystem.
#
# To build current docs, including archives:
# DOCKER_BUILDKIT=1 docker build --target=deploy-source --output=./_site .
#
# To build without archives:
# DOCKER_BUILDKIT=1 docker build --target=deploy-source --build-arg ENABLE_ARCHIVES=false --output=./_site .
FROM archives AS deploy-source
COPY --from=current /usr/share/nginx/html /

# Final stage, which includes nginx, and, depending on ENABLE_ARCHIVES, either
# current docs and archived versions (ENABLE_ARCHIVES=true), or only the current
# docs (ENABLE_ARCHIVES=false).
#
# To build current docs, including archives:
# DOCKER_BUILDKIT=1 docker build -t docs .
#
# To build without archives:
# DOCKER_BUILDKIT=1 docker build -t docs --build-arg ENABLE_ARCHIVES=false .
FROM nginx:alpine AS deploy
ENV TARGET=/usr/share/nginx/html
WORKDIR $TARGET

COPY --from=archives / .
COPY --from=current  /usr/share/nginx/html .

# Configure NGINX
COPY _deploy/nginx/default.conf /etc/nginx/conf.d/default.conf
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'
