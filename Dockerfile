# Get basic configs and Jekyll env
FROM docs/docker.github.io:docs-builder AS builder

# Set the target again
ENV TARGET=/usr/share/nginx/html

# Set the source directory to md_source
ENV SOURCE=md_source

# Put the archives into place
# This looks like it is going in the wrong direction, but the static
# HTML is in the "TARGET" location in the docs-base image, and we need to
# get it into our "SOURCE" location because it is part of our source.

COPY --from=docs/docker.github.io:docs-base ${TARGET} ${SOURCE}
# Get a colliding file out of the way
RUN rm ${SOURCE}/index.html

# Get the current docs from the checked out branch
# ${SOURCE} will contain a directory for each archive
COPY . ${SOURCE}

# ${SOURCE} now contains the static HTML for the archives in vrsion-specific
# directories, as well as the Markdown files for master.
# We still need to fetch upstream resources and then to build with Jekyll.


####### START UPSTREAM RESOURCES ########
# Set vars used by fetch-upstream-resources.sh script
## Branch to pull from, per ref doc
## To get master from svn the svn branch needs to be 'trunk'. To get a branch from svn it needs to be 'branches/branchname'

# Engine
ENV ENGINE_SVN_BRANCH="branches/17.06.x"
ENV ENGINE_BRANCH="17.06.x"

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

# Reset to alpine so we don't get any docs source or extra apps
FROM nginx:alpine

# Set the target again
ENV TARGET=/usr/share/nginx/html

# Get the built docs output from the previous step
COPY --from=builder ${TARGET} ${TARGET}

# Get the nginx config from the nginx-onbuild image
COPY --from=docs/docker.github.io:nginx-onbuild /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf

# Serve the site (target), which is now all static HTML
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'