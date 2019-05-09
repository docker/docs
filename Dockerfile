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
ARG ENGINE_BRANCH="18.09.x"

# Distribution
ARG DISTRIBUTION_BRANCH="release/2.6"


###
# Set up base stages for building and deploying
###

# Get basic configs and Jekyll env
FROM docs/docker.github.io:docs-builder AS builderbase
ENV TARGET=/usr/share/nginx/html
WORKDIR /usr/src/app/md_source/

# Set vars used by fetch-upstream-resources.sh script
# Branch to pull from, per ref doc. To get master from svn the svn branch needs
# to be 'trunk'. To get a branch from svn it needs to be 'branches/branchname'
ARG ENGINE_BRANCH
ENV ENGINE_BRANCH=${ENGINE_BRANCH}
ENV ENGINE_SVN_BRANCH=branches/${ENGINE_BRANCH}

ARG DISTRIBUTION_BRANCH
ENV DISTRIBUTION_BRANCH=${DISTRIBUTION_BRANCH}
ENV DISTRIBUTION_SVN_BRANCH=branches/${DISTRIBUTION_BRANCH}


# Reset to alpine so we don't get any docs source or extra apps
FROM nginx:alpine AS deploybase
ENV TARGET=/usr/share/nginx/html

# Get the nginx config from the nginx-onbuild image
# This hardly ever changes so should usually be cached
COPY --from=docs/docker.github.io:nginx-onbuild /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf

# Set the default command to serve the static HTML site
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'


# Build the archived docs
# these docs barely change, so can be cached
FROM deploybase AS archives
# Get all the archive static HTML and put it into place. To add a new archive,
# add it here, and ALSO edit _data/docsarchives/archives.yaml to add it to the drop-down
COPY --from=docs/docker.github.io:v17.03 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v17.06 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v17.09 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v17.12 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v18.03 ${TARGET} ${TARGET}

# Fetch library samples (documentation from official images on Docker Hub)
# Only add the files that are needed to build these reference docs, so that
# these docs are only rebuilt if changes were made to the configuration.
# @todo find a way to build HTML in this stage, and still have them included in the navigation tree
FROM builderbase AS library-samples
COPY ./_scripts/fetch-library-samples.sh ./_scripts/
COPY ./_samples/boilerplate.txt ./_samples/
RUN bash ./_scripts/fetch-library-samples.sh

# Fetch upstream resources (reference documentation)
# Only add the files that are needed to build these reference docs, so that
# these docs are only rebuilt if changes were made to the configuration.
FROM builderbase AS upstream-resources
COPY ./_scripts/fetch-upstream-resources.sh ./_scripts/
COPY ./_config.yml .
COPY ./_data/toc.yaml ./_data/
RUN bash ./_scripts/fetch-upstream-resources.sh .


# Build the current docs from the checked out branch
FROM builderbase AS current
COPY . .
COPY --from=library-samples /usr/src/app/md_source/. ./
COPY --from=upstream-resources /usr/src/app/md_source/. ./

# Build the static HTML, now that everything is in place
RUN jekyll build -d ${TARGET}

# Fix up some links, don't touch the archives
RUN find ${TARGET} -type f -name '*.html' | grep -vE "v[0-9]+\." | while read i; do sed -i 's#href="https://docs.docker.com/#href="/#g' "$i"; done


# Docs with archives (for deploy)
FROM archives AS deploy

# Add the current version of the docs
COPY --from=current ${TARGET} ${TARGET}
