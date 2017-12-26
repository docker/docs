# Set to the version for this archive
ARG VER=v1.5

# This image comes from the Dockerfile.builder.onbuild file in the publish-tools branch
# https://github.com/docker/docker.github.io/blob/publish-tools/Dockerfile.builder.onbuild
FROM docs/docker.github.io:docs-builder-onbuild AS builder

# Reset the docs/docker.github.io:nginx-onbuild image, which is based on nginx:alpine
# This image comes from the Dockerfule.nginx.onbuild in the publish-tools branch
# https://github.com/docker/docker.github.io/blob/publish-tools/Dockerfile.nginx.onbuild
FROM docs/docker.github.io:nginx-onbuild
