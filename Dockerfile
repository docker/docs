# Get Jekyll build env
FROM docs/docker.github.io:docs-builder AS builder

# Set to the version for this archive
# You need to set this one more time below too
ENV VER=v1.13

# Build the docs from this branch
COPY . /source
RUN jekyll build --source /source --destination /site/${VER}

# Do post-processing on archive
COPY --from=docs/docker.github.io:docs-config /scripts/* /usr/bin/
RUN /usr/bin/fix-archives.sh /site/ ${VER}

# Make an index.html which will redirect / to /${VER}/
RUN VER=${VER} echo "<html><head><title>Redirect for $VER</title><meta http-equiv=\"refresh\" content=\"0;url='/$VER/'\" /></head><body><p>If you are not redirected automatically, click <a href=\"/$VER/\">here</a>.</p></body></html>" > /site/index.html

# Reset to bare Nginx image
FROM nginx:alpine

# Set to the version for this archive
ENV VER=v1.13

# Copy the Nginx config from the docs-config image
COPY --from=docs/docker.github.io:docs-config /conf/nginx-overrides.conf /etc/nginx/conf.d/default.conf

# Clean out any existing HTML files
RUN ls -a /usr/share/nginx/html/* && rm -rf /usr/share/nginx/html/*

# Copy the HTML from the builder stage to the default location for Nginx
COPY --from=builder /site/${VER}/ /usr/share/nginx/html/${VER}/
COPY --from=builder /site/index.html /usr/share/nginx/html/index.html
COPY --from=builder /site/index.html /usr/share/nginx/html/404.html

# Start Nginx to serve the archive at / (which will redirect to the version-specific dir)
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'
