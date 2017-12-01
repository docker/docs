# Get archives
FROM docs/docker.github.io:docs-base AS docs_base

# Set the target to the default serving directory for Nginx
# which is where the archive HTML is
ENV TARGET=/usr/share/nginx/html

# Move the temporary index.html out of the way
RUN rm ${TARGET}/index.html

# Get basic configs and Jekyll env
FROM docs/docker.github.io:docs-builder AS builder

# Set the target again
ENV TARGET=/usr/share/nginx/html

# Get the current docs from the checked out branch
# md_source will contain a directory for each archive
COPY . md_source

# Move the archives into place
# md_source will still contain the archive directories but also the current docs
# at the root

COPY --from=docs_base ${TARGET} md_source

####### START UPSTREAM RESOURCES ########
# REMOVE FROM THE LINE ABOVE TO THE END UPSTREAM RESOURCES LINE FOR ARCHIVES #
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
RUN bash ./md_source/_scripts/fetch-upstream-resources.sh md_source
####### END UPSTREAM RESOURCES ########

# Build the static HTML
RUN jekyll build -s md_source -d ${TARGET} --config md_source/_config.yml

# Fix up some links, don't touch the archives
RUN find ${TARGET} -type f -name '*.html' | grep -vE "v[0-9]+\." | while read i; do sed -i 's#href="https://docs.docker.com/#href="/#g' "$i"; done

# Reset to alpine so we don't get any docs source or extra apps
FROM nginx:alpine

# Set the target again
ENV TARGET=/usr/share/nginx/html

# Get the built docs output from the previous step
COPY --from=builder ${TARGET} ${TARGET}

# Override some nginx conf -- this gets added to default nginx conf
COPY --from=docs/docker.github.io:docs-config /conf/nginx-overrides.conf /etc/nginx/conf.d/default.conf

# Serve the site (target), which is now all static HTML
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'