FROM docs-base:hugo-feature-work
MAINTAINER Mary Anthony <mary@docker.com> (@moxiegirl)

# to get the git info for this repo
COPY . /src

COPY . /docs/content/opensource/

RUN find /docs/content/opensource -type f -name "*.md" -exec sed -i.old  -e '/^<!.*metadata]>/g' -e '/^<!.*end-metadata.*>/g' {} \;