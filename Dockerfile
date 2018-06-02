# This Dockerfile builds the docs for https://docs.docker.com/
# from the master branch of https://github.com/docker/docker.github.io
#
# Here is the sequence:
# 1.  Set up the build
# 2.  Fetch upstream resources
# 3.  Build static HTML from master
# 4.  Reset to clean tiny nginx image
# 5.  Copy Nginx config and archive HTML, which don't change often and can be cached
# 6.  Copy static HTML from previous build stage (step 3)
#
# When the image is run, it starts Nginx and serves the docs at port 4000

# Get basic configs and Jekyll env
FROM docs/docker.github.io:docs-builder AS builder

# Set the target again
ENV TARGET=/usr/share/nginx/html

# Set the source directory to md_source
ENV SOURCE=md_source

# Get the current docs from the checked out branch
# ${SOURCE} will contain a directory for each archive
COPY . ${SOURCE}

####### START UPSTREAM RESOURCES ########
# Set vars used by fetch-upstream-resources.sh script
## Branch to pull from, per ref doc
## To get master from svn the svn branch needs to be 'trunk'. To get a branch from svn it needs to be 'branches/branchname'

# Engine
ENV ENGINE_SVN_BRANCH="branches/17.09.x"
ENV ENGINE_BRANCH="17.09.x"

# Distribution
ENV DISTRIBUTION_SVN_BRANCH="branches/release/2.6"
ENV DISTRIBUTION_BRANCH="release/2.6"

# Fetch upstream resources
RUN bash ./${SOURCE}/_scripts/fetch-upstream-resources.sh ${SOURCE}
####### END UPSTREAM RESOURCES ########


# Build the static HTML, now that everything is in place

RUN jekyll build -s ${SOURCE} -d ${TARGET} --config ${SOURCE}/_config.yml

# Fix up some links, don't touch the archives
RUN find ${TARGET} -type f -name '*.html' | grep -vE "v[0-9]+\." | while read i; do sed -i 's#href="https://docs.docker.com/#href="/#g' "$i"; done

# BUILD OF MASTER DOCS IS NOW DONE!

# Reset to alpine so we don't get any docs source or extra apps
FROM nginx:alpine

# Set the target again
ENV TARGET=/usr/share/nginx/html

# Get the nginx config from the nginx-onbuild image
# This hardly ever changes so should usually be cached
COPY --from=docs/docker.github.io:nginx-onbuild /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf

# Get all the archive static HTML and put it into place
# Go oldest-to-newest to take advantage of the fact that we change older
# archives less often than new ones.
# To add a new archive, add it here
# AND ALSO edit _data/docsarchives/archives.yaml to add it to the drop-down
COPY --from=docs/docker.github.io:v1.4 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v1.5 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v1.6 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v1.7 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v1.8 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v1.9 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v1.10 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v1.11 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v1.12 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v1.13 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v17.03 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v17.06 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v17.09 ${TARGET} ${TARGET}
COPY --from=docs/docker.github.io:v17.12 ${TARGET} ${TARGET}

# Get the built docs output from the previous build stage
# This ordering means all previous layers can come from cache unless an archive
# changes

COPY --from=builder ${TARGET} ${TARGET}

# Serve the site (target), which is now all static HTML
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'
