FROM docs/docker.github.io:docs-base

# docs-base contains: GitHub Pages, nginx, wget, svn, and the docs archives,
# running on Alpine. See the contents of docs-base at:
# https://github.com/docker/docker.github.io/tree/docs-base

# First, build non-edge (all of this is duplicated later -- that is on purpose)

# Copy master into target directory (skipping files / folders in .dockerignore)
# These files represent the current docs
COPY . md_source

# Move built html into md_source directory so we can reuse the target directory
# to hold the static output.
# Pull reference docs from upstream locations, then build the master docs
# into static HTML in the "target" directory using Jekyll
# then nuke the md_source directory.

## Branch to pull from, per ref doc
## To get master from svn the svn branch needs to be 'trunk'. To get a branch from svn it needs to be 'branches/branchname'

# Engine
ENV ENGINE_SVN_BRANCH="branches/17.06.x"
ENV ENGINE_BRANCH="17.06.x"

# Distribution
ENV DISTRIBUTION_SVN_BRANCH="branches/release/2.6"
ENV DISTRIBUTION_BRANCH="release/2.6"

RUN sh md_source/_scripts/fetch-upstream-resources.sh \
	&& jekyll build -s md_source -d target --config md_source/_config.yml \
	&& rm -rf target/apidocs/layouts \
	&& find target -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g' \
	&& rm -rf md_source
