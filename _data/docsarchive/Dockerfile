FROM docs:base

# The branch name to check out.
# Override at build-time using --build-arg VERSION=<foo>
ARG BRANCH="master"

# Version. Only used informational currently
ARG VERSION="v1.12"

# Base URL to use (e.g. "v1.12/"). Must have a trailing slash if not empty
ARG BASEURL=""

WORKDIR /usr/share/nginx/

## Branch to pull from, per ref doc
ENV ENGINE_BRANCH="1.12.x"
ENV DISTRIBUTION_BRANCH="release/2.5"

# Get the branch contents without history
RUN echo "Started $VERSION, from branch $BRANCH, using base-url $BASEURL" \
    && git clone --depth 1 --branch $BRANCH https://github.com/docker/docker.github.io docs_src \
    \
    # Reference docs for "master" docs, as they are not 
    # part of the repository
    \
    && if [ "$BRANCH" == "master" ]; then \
      echo "Fetching reference docs" \
      && svn --non-interactive --trust-server-cert co https://github.com/docker/docker/branches/$ENGINE_BRANCH/docs/reference docs_src/engine/reference \
      && svn --non-interactive --trust-server-cert co https://github.com/docker/docker/branches/$ENGINE_BRANCH/docs/extend docs_src/engine/extend \
      && wget -O docs_src/engine/deprecated.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/deprecated.md \
      && svn --non-interactive --trust-server-cert co https://github.com/docker/distribution/branches/$DISTRIBUTION_BRANCH/docs/spec docs_src/registry/spec \
      && wget -O docs_src/registry/configuration.md https://raw.githubusercontent.com/docker/distribution/$DISTRIBUTION_BRANCH/docs/configuration.md; \
      \
      # remove api.md.tmp as it stresses out Liquid
      rm docs_src/registry/spec/api.md.tmpl; \
    fi \
    \
    && jekyll build -s docs_src -d html \
    && rm -rf docs_src \
    \
    # Replace / rewrite some URLs so that links in the archive go to the correct
    # location. Note that the order in which these replacements are done is
    # important. Changing the order may result in replacements being done
    # multiple times.
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="http://docs-stage.docker.com/#href="/#g' \
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="https://docs-stage.docker.com/#src="/#g' \
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g' \
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="https://docs.docker.com/#src="/#g' \
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="http://docs.docker.com/#href="/#g' \
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="http://docs.docker.com/#src="/#g' \
    \
    # Substitute https:// for schema-less resources (src="//analytics.google.com")
    # We're replacing them to prevent them being seen as absolute paths below
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="//#href="https://#g' \
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="//#src="https://#g' \
    \
    # And some archive versions already have URLs starting with '/version/'
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/'"$BASEURL"'#href="/#g' \
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/'"$BASEURL"'#src="/#g' \
    \
    # Archived versions 1.7 and under use some absolute links, and v1.10 uses
    # "relative" links to sources (href="./css/"). Remove those to make them
    # work :)
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="\./#href="/#g' \
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="\./#src="/#g' \
    \
    # Create permalinks for archived versions
    \
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/'"$BASEURL"'#g' \
    && find html -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/'"$BASEURL"'#g'

COPY nginx-docs.conf /etc/nginx/conf.d/default.conf
