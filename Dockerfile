FROM docs/docker.github.io:docs-base

# docs-base contains: GitHub Pages, nginx, and the docs archives, running on
# Debian Jesse. See the contents of docs-base at:
# https://github.com/docker/docker.github.io/tree/docs-base

# Copy master into target directory (skipping files / folders in .dockerignore)
# These files represent the current docs
COPY . md_source

# Move built html into md_source directory so we can reuse the target directory
# to hold the static output.
# Pull reference docs from upstream locations, then build the master docs
# into static HTML in the "target" directory using Jekyll
# then nuke the md_source directory.

## Branch to pull from, per ref doc
ENV ENGINE_BRANCH="17.03.x"
ENV DISTRIBUTION_BRANCH="release/2.5"

RUN svn co https://github.com/docker/docker/branches/$ENGINE_BRANCH/docs/extend md_source/engine/extend \
	&& wget -O md_source/engine/api/v1.18.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.18.md \
	&& wget -O md_source/engine/api/v1.19.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.19.md \
	&& wget -O md_source/engine/api/v1.20.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.20.md \
	&& wget -O md_source/engine/api/v1.21.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.21.md \
	&& wget -O md_source/engine/api/v1.22.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.22.md \
	&& wget -O md_source/engine/api/v1.23.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.23.md \
	&& wget -O md_source/engine/api/v1.24.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.24.md \
	&& wget -O md_source/engine/api/version-history.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/version-history.md \
	&& wget -O md_source/engine/reference/glossary.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/glossary.md \
	&& wget -O md_source/engine/reference/builder.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/builder.md \
	&& wget -O md_source/engine/reference/run.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/run.md \
	&& wget -O md_source/engine/reference/commandline/cli.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/commandline/cli.md \
	&& wget -O md_source/engine/deprecated.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/deprecated.md \
	&& svn co https://github.com/docker/distribution/branches/$DISTRIBUTION_BRANCH/docs/spec md_source/registry/spec \
	&& rm md_source/registry/spec/api.md.tmpl \
	&& wget -O md_source/registry/configuration.md https://raw.githubusercontent.com/docker/distribution/$DISTRIBUTION_BRANCH/docs/configuration.md \
	&& rm -rf md_source/apidocs/cloud-api-source \
	&& rm -rf md_source/tests \
  && wget -O md_source/engine/api/v1.25/swagger.yaml https://raw.githubusercontent.com/docker/docker/v1.13.0/api/swagger.yaml \
  && wget -O md_source/engine/api/v1.26/swagger.yaml https://raw.githubusercontent.com/docker/docker/v17.03.0-ce/api/swagger.yaml \
  && wget -O md_source/engine/api/v1.27/swagger.yaml https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/api/swagger.yaml \
	&& jekyll build -s md_source -d target \
	&& rm -rf target/apidocs/layouts \
	&& find target -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g' \
	&& rm -rf md_source
