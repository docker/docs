#!/usr/bin/env bash

echo "Building utilities image"

# This is a separate image ("docs/utilities")
docker build -t docs/utilities . -f-<<-'EOF' >/dev/null
    FROM golang:1.9-alpine AS minifier
    RUN apk add -q --no-cache git
    RUN go get -d github.com/tdewolff/minify/cmd/minify \
     && go build -o /usr/bin/minify github.com/tdewolff/minify/cmd/minify

    FROM scratch
    COPY --from=minifier /usr/bin/minify /
    COPY ./scripts/* /
EOF


# Assume these are the dockerfiles for each archive branch
VERSIONS="v1.4 v1.5 v1.6 v1.7 v1.8 v1.9 v1.10 v1.11 v1.12 v1.13 v17.03 v17.06"

for VER in ${VERSIONS}; do

echo "Building docs-archive:${VER}"


docker build -t docs-archive:${VER} . -f-<<-EOF >/dev/null
    FROM docs/docker.github.io:${VER}
    ENV  TARGET=/usr/share/nginx/html VER=${VER}
    COPY --from=docs/utilities /* /usr/bin/
    RUN normalize_links.sh \${TARGET} \${VER} \
     && minify_assets.sh \${TARGET} \${VER}
EOF

done;
