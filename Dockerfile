
# Build tdewolff/minify
# TODO move to separate image
FROM golang:1.9-alpine AS minifier
RUN apk add --no-cache git
RUN go get -d github.com/tdewolff/minify/cmd/minify \
 && go build -v -o /usr/bin/minify github.com/tdewolff/minify/cmd/minify

# TODO move to separate image
FROM scratch AS utilities
COPY --from=minifier /usr/bin/minify /
COPY ./scripts/* /

# Get archival docs from each canonical image
# When there is a new archive, add it here and also further down
# where we copy out of it
FROM docs/docker.github.io:v1.4 AS archive_v1.4
ENV  TARGET=/usr/share/nginx/html VER=v1.4
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v1.5 AS archive_v1.5
ENV  TARGET=/usr/share/nginx/html VER=v1.5
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v1.6 AS archive_v1.6
ENV  TARGET=/usr/share/nginx/html VER=v1.6
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v1.7 AS archive_v1.7
ENV  TARGET=/usr/share/nginx/html VER=v1.7
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v1.8 AS archive_v1.8
ENV  TARGET=/usr/share/nginx/html VER=v1.8
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v1.9 AS archive_v1.9
ENV  TARGET=/usr/share/nginx/html VER=v1.9
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v1.10 AS archive_v1.10
ENV  TARGET=/usr/share/nginx/html VER=v1.10
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v1.11 AS archive_v1.11
ENV  TARGET=/usr/share/nginx/html VER=v1.11
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v1.12 AS archive_v1.12
ENV  TARGET=/usr/share/nginx/html VER=v1.12
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v1.13 AS archive_v1.13
ENV  TARGET=/usr/share/nginx/html VER=v1.13
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v17.03 AS archive_v17.03
ENV  TARGET=/usr/share/nginx/html VER=v17.03
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

FROM docs/docker.github.io:v17.06 AS archive_v17.06
ENV  TARGET=/usr/share/nginx/html VER=v17.06
COPY --from=utilities /* /usr/bin/
RUN apk add -q --no-cache gzip \
 && normalize_links.sh ${TARGET} ${VER} \
 && minify_assets.sh ${TARGET} ${VER} \
 && create_permalinks.sh ${TARGET} ${VER} \
 && compress_assets.sh ${TARGET}

# Reset with nginx again, so we don't get scripts or extra apps in the final image
FROM nginx:alpine AS optimized

# Reset TARGET since we lost it when we reset the image
ENV TARGET=/usr/share/nginx/html

# Copy HTML from each stage above
COPY --from=archive_v1.4   ${TARGET} ${TARGET}/v1.4
COPY --from=archive_v1.5   ${TARGET} ${TARGET}/v1.5
COPY --from=archive_v1.6   ${TARGET} ${TARGET}/v1.6
COPY --from=archive_v1.7   ${TARGET} ${TARGET}/v1.7
COPY --from=archive_v1.8   ${TARGET} ${TARGET}/v1.8
COPY --from=archive_v1.9   ${TARGET} ${TARGET}/v1.9
COPY --from=archive_v1.10  ${TARGET} ${TARGET}/v1.10
COPY --from=archive_v1.11  ${TARGET} ${TARGET}/v1.11
COPY --from=archive_v1.12  ${TARGET} ${TARGET}/v1.12
COPY --from=archive_v1.13  ${TARGET} ${TARGET}/v1.13
COPY --from=archive_v17.03 ${TARGET} ${TARGET}/v17.03
COPY --from=archive_v17.06 ${TARGET} ${TARGET}/v17.06

## Copy the above and change the references to the version, to make a new archive

# Set `--build-arg REPORT_SIZE=1` to print the size-report during build
ARG REPORT_SIZE
COPY ./scripts/size_report.sh /usr/bin/
RUN \
  if [ -n "$REPORT_SIZE" ]; then \
      apk add -q --no-cache coreutils; \
      size_report.sh ${TARGET}; \
  fi;

# This index file gets overwritten, but it serves a sort-of useful purpose in
# making the docs/docs-base image browsable:
COPY index.html ${TARGET}/

# Copy the Nginx config
COPY ./default.conf  /etc/nginx/conf.d/default.conf

# Serve the docs
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'
