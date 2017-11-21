# Build the docs
FROM starefossen/github-pages:147 AS build_docs
COPY . /source
RUN jekyll build --source /source --destination /site

# Reset with nginx, so we don't get docs source in the image
FROM nginx:alpine
# Override some nginx conf -- this gets added to default nginx conf
COPY nginx-overrides.conf /etc/nginx/conf.d/000-default.conf
COPY --from=build_docs /site /etc/nginx/html/
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'