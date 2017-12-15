# Start over with nginx:alpine
FROM nginx:alpine AS docs_base

# Set TARGET to the default location where Nginx serves files from
# (and where they are stored in all the archives)
ENV TARGET=/usr/share/nginx/html

# For each of these blocks
# - Get an archive's HTML
# - Run fix-archives.sh, which calls scripts to fix relative links, minify
#   assets, create permanent links, and compress assets
#
# To make a new archive, copy the last line and modify the tag
# in the --from flag (env vars are not allowed there yet)

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

# This is just so that docs-base is self-browseable and needs to go after all
# the archive copying stuff
COPY index.html ${TARGET}/index.html

# The static HTML is now ready to go

# Copy the Nginx config from the docs-config image
COPY --from=docs/docker.github.io:docs-config /conf/nginx-overrides.conf /etc/nginx/conf.d/default.conf

# Serve the docs
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'
