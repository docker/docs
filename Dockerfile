# Get the docs-builder image so we can get the nginx config from it later
FROM docs/docker.github.io:docs-builder AS builder

# Get archival docs from each canonical image
# When there is a new archive, add it here and also further down
# where we copy out of it
FROM docs/docker.github.io:v1.4 AS archive_v1.4
FROM docs/docker.github.io:v1.5 AS archive_v1.5
FROM docs/docker.github.io:v1.6 AS archive_v1.6
FROM docs/docker.github.io:v1.7 AS archive_v1.7
FROM docs/docker.github.io:v1.8 AS archive_v1.8
FROM docs/docker.github.io:v1.9 AS archive_v1.9
FROM docs/docker.github.io:v1.10 AS archive_v1.10
FROM docs/docker.github.io:v1.11 AS archive_v1.11
FROM docs/docker.github.io:v1.12 AS archive_v1.12
FROM docs/docker.github.io:v1.13 AS archive_v1.13
FROM docs/docker.github.io:v17.03 AS archive_v17.03
FROM docs/docker.github.io:v17.06 AS archive_v17.06

# Reset with nginx, so we don't get docs source in the image
FROM nginx:alpine AS docs_base

COPY fix_archives.sh /usr/bin/fix_archives.sh
ENV TARGET=/usr/share/nginx/html

# Copy HTML from each stage above and run script to fix relative links
COPY --from=archive_v1.4 ${TARGET} ${TARGET}/v1.4
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v1.4'

COPY --from=archive_v1.5 ${TARGET} ${TARGET}/v1.5
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v1.5'

COPY --from=archive_v1.6 ${TARGET} ${TARGET}/v1.6
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v1.6'

COPY --from=archive_v1.7 ${TARGET} ${TARGET}/v1.7
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v1.7'

COPY --from=archive_v1.8 ${TARGET} /${TARGET}/v1.8
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v1.8'

COPY --from=archive_v1.9 ${TARGET} ${TARGET}/v1.9
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v1.9'

COPY --from=archive_v1.10 ${TARGET} ${TARGET}/v1.10
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v1.10'

COPY --from=archive_v1.11 ${TARGET} ${TARGET}/v1.11
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v1.11'

COPY --from=archive_v1.12 ${TARGET} ${TARGET}/v1.12
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v1.12'

COPY --from=archive_v1.13 ${TARGET} ${TARGET}/v1.13
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v1.13'

COPY --from=archive_v17.03 ${TARGET} ${TARGET}/v17.03
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v17.03'

COPY --from=archive_v17.06 ${TARGET} ${TARGET}/v17.06
RUN sh /usr/bin/fix_archives.sh ${TARGET} 'v17.06'

## Copy the above two lines and change the three references
## to the version, to make a new archive

# This index file gets overwritten, but it serves a sort-of useful purpose in
# making the docs/docs-base image browsable:
COPY index.html ${TARGET}/

# Reset with nginx again, so we don't get scripts or extra apps in the final image
FROM nginx:alpine

# Copy the Nginx config
COPY --from=builder /conf/nginx-overrides.conf /etc/nginx/conf.d/default.conf

# Reset TARGET since we lost it when we reset the image
ENV TARGET=/usr/share/nginx/html

# Copy the static HTML files to where Nginx will serve them
COPY --from=docs_base ${TARGET} ${TARGET}

# Serve the docs
CMD echo -e "Docker docs are viewable at:\nhttp://0.0.0.0:4000"; exec nginx -g 'daemon off;'
