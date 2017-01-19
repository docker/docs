FROM starefossen/github-pages

## Branch to pull from, per ref doc
ENV ENGINE_BRANCH="1.12.x"
ENV DISTRIBUTION_BRANCH="release/2.5"

# The statements below pull reference docs from upstream locations,
# then build the whole site to static HTML using Jekyll

RUN git clone --depth 1 --branch master https://github.com/docker/docker.github.io allv \
 && svn --non-interactive --trust-server-cert co https://github.com/docker/docker/branches/$ENGINE_BRANCH/docs/reference allv/engine/reference \
 && svn --non-interactive --trust-server-cert co https://github.com/docker/docker/branches/$ENGINE_BRANCH/docs/extend allv/engine/extend \
 && wget -O allv/engine/deprecated.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/deprecated.md \
 && svn --non-interactive --trust-server-cert co https://github.com/docker/distribution/branches/$DISTRIBUTION_BRANCH/docs/spec allv/registry/spec \
 && wget -O allv/registry/configuration.md https://raw.githubusercontent.com/docker/distribution/$DISTRIBUTION_BRANCH/docs/configuration.md \
 && rm -rf allv/apidocs/cloud-api-source \
 && rm -rf allv/tests \
 && jekyll build -s allv -d allvbuild \
 && find allvbuild/engine/reference -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g' \
 && find allvbuild/engine/extend -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g' \
 && rm -rf allv

# Serve the site, which is now all static HTML
CMD jekyll serve -s /usr/src/app/allvbuild --no-watch -H 0.0.0.0 -P 4000
