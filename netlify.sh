#!/bin/bash

# Script to build master and edge docs on Netlify
# Replicates the non-archive functionality of the Dockerfile

# First build master
ENGINE_SVN_BRANCH="branches/17.03.x"
ENGINE_BRANCH="17.03.x"
DISTRIBUTION_SVN_BRANCH="branches/release/2.6"
DISTRIBUTION_BRANCH="release/2.6"

svn co https://github.com/docker/docker/$ENGINE_SVN_BRANCH/docs/extend engine/extend \
	&& wget -O engine/api/v1.18.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.18.md \
	&& wget -O engine/api/v1.19.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.19.md \
	&& wget -O engine/api/v1.20.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.20.md \
	&& wget -O engine/api/v1.21.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.21.md \
	&& wget -O engine/api/v1.22.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.22.md \
	&& wget -O engine/api/v1.23.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.23.md \
	&& wget -O engine/api/v1.24.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.24.md \
	&& wget -O engine/api/version-history.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/version-history.md \
	&& wget -O engine/reference/glossary.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/glossary.md \
	&& wget -O engine/reference/builder.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/builder.md \
	&& wget -O engine/reference/run.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/run.md \
	&& wget -O engine/reference/commandline/cli.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/commandline/cli.md \
	&& wget -O engine/deprecated.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/deprecated.md \
	&& svn co https://github.com/docker/distribution/$DISTRIBUTION_SVN_BRANCH/docs/spec registry/spec \
	&& rm registry/spec/api.md.tmpl \
	&& wget -O registry/configuration.md https://raw.githubusercontent.com/docker/distribution/$DISTRIBUTION_BRANCH/docs/configuration.md \
	&& rm -rf apidocs/cloud-api-source \
	&& rm -rf tests \
  && wget -O engine/api/v1.25/swagger.yaml https://raw.githubusercontent.com/docker/docker/v1.13.0/api/swagger.yaml \
  && wget -O engine/api/v1.26/swagger.yaml https://raw.githubusercontent.com/docker/docker/v17.03.0-ce/api/swagger.yaml \
  && wget -O engine/api/v1.27/swagger.yaml https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/api/swagger.yaml \
	&& jekyll build -d site --config _config.yml \
	&& rm -rf site/apidocs/layouts \
	&& find site -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g'

echo "Jekyll site in site/"

# Then build edge

ENGINE_SVN_BRANCH="branches/17.04.x"
ENGINE_BRANCH="17.04.x"
DISTRIBUTION_SVN_BRANCH="branches/release/2.6"
DISTRIBUTION_BRANCH="release/2.6"

svn co https://github.com/docker/docker/$ENGINE_SVN_BRANCH/docs/extend engine/extend \
  && wget -O engine/api/v1.18.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.18.md \
  && wget -O engine/api/v1.19.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.19.md \
  && wget -O engine/api/v1.20.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.20.md \
  && wget -O engine/api/v1.21.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.21.md \
  && wget -O engine/api/v1.22.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.22.md \
  && wget -O engine/api/v1.23.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.23.md \
  && wget -O engine/api/v1.24.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.24.md \
  && wget -O engine/api/version-history.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/version-history.md \
  && wget -O engine/reference/glossary.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/glossary.md \
  && wget -O engine/reference/builder.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/builder.md \
  && wget -O engine/reference/run.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/run.md \
  && wget -O engine/reference/commandline/cli.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/commandline/cli.md \
  && wget -O engine/deprecated.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/deprecated.md \
  && svn co https://github.com/docker/distribution/$DISTRIBUTION_SVN_BRANCH/docs/spec registry/spec \
  && rm registry/spec/api.md.tmpl \
  && wget -O registry/configuration.md https://raw.githubusercontent.com/docker/distribution/$DISTRIBUTION_BRANCH/docs/configuration.md \
  && rm -rf apidocs/cloud-api-source \
  && rm -rf tests \
  && wget -O engine/api/v1.25/swagger.yaml https://raw.githubusercontent.com/docker/docker/v1.13.0/api/swagger.yaml \
  && wget -O engine/api/v1.26/swagger.yaml https://raw.githubusercontent.com/docker/docker/v17.03.0-ce/api/swagger.yaml \
  && wget -O engine/api/v1.27/swagger.yaml https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/api/swagger.yaml \
  && jekyll build -d site/edge --config _config-edge.yml \
	&& echo "Jekyll site in site/edge"
  && rm -rf site/edge/apidocs/layouts \
  && find site/edge -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g'

  # Replace / rewrite some URLs so that links in the edge directory go to the correct
  # location. Note that the order in which these replacements are done is
  # important. Changing the order may result in replacements being done
  # multiple times.
  # First, remove the domain from URLs that include the domain
  VER="edge" \
  && BASEURL="$VER/" \
  && find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="http://docs-stage.docker.com/#href="/#g' \
  && find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="https://docs-stage.docker.com/#src="/#g' \
  && find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g' \
  && find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="https://docs.docker.com/#src="/#g' \
  && find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="http://docs.docker.com/#href="/#g' \
  && find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="http://docs.docker.com/#src="/#g'

  # Substitute https:// for schema-less resources (src="//analytics.google.com")
  # We're replacing them to prevent them being seen as absolute paths below
  find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="//#href="https://#g' \
  && find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="//#src="https://#g'

  # And some archive versions already have URLs starting with '/version/'
  find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/'"$BASEURL"'#href="/#g' \
  && find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/'"$BASEURL"'#src="/#g'

  # Archived versions 1.7 and under use some absolute links, and v1.10 uses
  # "relative" links to sources (href="./css/"). Remove those to make them
  # work :)
  find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="\./#href="/#g' \
  && find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="\./#src="/#g'

  # Create permalinks for archived versions
  find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/'"$BASEURL"'#g' \
  && find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/'"$BASEURL"'#g'

  # Fix 'Back to Stable Docs' URL
  find site/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#<li id="stable-cta"><a href="/edge/">Back to Stable docs</a></li>#<li id="stable-cta"><a href="/">Back to Stable docs</a></li>#g'

# Make sure we don't commit any changes to the actual source files
git reset --hard

if [ -d site/edge ]; then
	echo "Edge is where it is supposed to be"
fi
if [ -d _site ]; then
	rm -rf _site
fi

mv site _site
jekyll serve -d _site --skip-initial-build -D
