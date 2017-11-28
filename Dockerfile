# Get Jekyll build env
FROM docs/docker.github.io:docs-builder AS builder

# Build the docs from this branch
COPY . /source
RUN jekyll build --source /source --destination /site

# Reset to bare Nginx image
FROM nginx:alpine

# Copy the Nginx config from the docs-config image
COPY --from=docs/docker.github.io:docs-config /conf/nginx-overrides.conf /etc/nginx/conf.d/default.conf

# Copy the HTML from the builder stage to the default location for Nginx
COPY --from=builder /site /usr/share/nginx/html

# Start Nginx to serve the archive at /
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'
