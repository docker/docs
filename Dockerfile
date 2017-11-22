# Get Jekyll build env and Nginx env
FROM docs/docker.github.io:docs-builder AS builder

# Build the docs from this branch
COPY . /source
RUN jekyll build --source /source --destination /site

# Reset to bare Nginx image
FROM nginx:alpine

# Copy the Nginx config
COPY --from=builder /conf/nginx-overrides.conf /etc/nginx/conf.d/default.conf

# Copy the HTML to the default location for Nginx and serve it
COPY --from=builder /site /usr/share/nginx/html
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'
