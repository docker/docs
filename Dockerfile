FROM docs/base:latest
MAINTAINER Mary Anthony <mary@docker.com> (@moxiegirl)

# To get the git info for this repo
COPY . /src

COPY . /docs/content/docker-trusted-registry/
