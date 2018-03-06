# Build minifier utility
FROM golang:1.9-alpine AS minifier
RUN apk add --no-cache git
RUN go get -d github.com/tdewolff/minify/cmd/minify \
 && go build -v -o /minify github.com/tdewolff/minify/cmd/minify

# Set the version of Github Pages to use for each docs archive
FROM starefossen/github-pages:177

# Get some utilities we need for post-build steps
RUN apk add --no-cache bash wget subversion gzip

# Copy scripts used for static HTML post-processing.
COPY scripts /scripts
COPY --from=minifier /minify /scripts/minify

# Print out a message if someone tries to run this image on its own
CMD echo 'This image is only meant to be used as a base image for building docs.'
