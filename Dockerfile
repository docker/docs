# Start over with nginx:alpine
FROM nginx:alpine AS docs_base

# Set TARGET to the default location where Nginx serves files from
# (and where they are stored in all the archives)
ENV TARGET=/usr/share/nginx/html

# This is just so that docs-base is self-browseable

COPY index.html ${TARGET}/index.html

# For each of these blocks
# - Get an archive's HTML
# - Run fix-archives.sh, which calls scripts to fix relative links, minify
#   assets, create permanent links, and compress assets
#
# To make a new archive, copy the three lines and change the VER variable and
# and the tag in the --from flag (env vars are not allowed there yet)

ENV VER=v1.4
COPY --from=docs/docker.github.io:v1.4 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v1.5
COPY --from=docs/docker.github.io:v1.5 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v1.6
COPY --from=docs/docker.github.io:v1.6 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v1.7
COPY --from=docs/docker.github.io:v1.7 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v1.8
COPY --from=docs/docker.github.io:v1.8 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v1.9
COPY --from=docs/docker.github.io:v1.9 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v1.10
COPY --from=docs/docker.github.io:v1.10 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v1.11
COPY --from=docs/docker.github.io:v1.11 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v1.12
COPY --from=docs/docker.github.io:v1.12 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v1.13
COPY --from=docs/docker.github.io:v1.13 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v17.03
COPY --from=docs/docker.github.io:v17.03 ${TARGET}/${VER}/ ${TARGET}/${VER}/

ENV VER=v17.06
COPY --from=docs/docker.github.io:v17.06 ${TARGET}/${VER}/ ${TARGET}/${VER}/

# The static HTML is now ready to go

# Reset with nginx again, so we don't get scripts or extra apps in the final image
FROM nginx:alpine

# Reset TARGET since we lost it when we reset the image
ENV TARGET=/usr/share/nginx/html

# Copy the Nginx config from the docs-config image
COPY --from=docs/docker.github.io:docs-config /conf/nginx-overrides.conf /etc/nginx/conf.d/default.conf

# Copy the static HTML files from the docs_base build stage to where Nginx will serve them
COPY --from=docs_base ${TARGET}/ ${TARGET}/

# Serve the docs
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'
