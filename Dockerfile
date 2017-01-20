FROM starefossen/github-pages

ENV VERSIONS="v1.4 v1.5 v1.6 v1.7 v1.8 v1.9 v1.10 v1.11"

# Create archive; check out each version, create HTML, tweak links
RUN git clone https://www.github.com/docker/docker.github.io temp; \
 for VER in $VERSIONS; do \
		git --git-dir=./temp/.git --work-tree=./temp checkout ${VER} \
		&& mkdir -p allvbuild/${VER} \
		&& jekyll build -s temp -d allvbuild/${VER} \
		&& find allvbuild/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/'"$VER"'/#g' \
		&& find allvbuild/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/'"$VER"'/#g' \
		&& find allvbuild/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/'"$VER"'/#g'; \
	done; \
	rm -rf temp

COPY . allv

## Branch to pull from, per ref doc
ENV ENGINE_BRANCH="1.13.x"
ENV DISTRIBUTION_BRANCH="release/2.5"

# The statements below pull reference docs from upstream locations,
# then build the whole site to static HTML using Jekyll

RUN svn co https://github.com/docker/docker/branches/$ENGINE_BRANCH/docs/extend allv/engine/extend \
 && wget -O allv/engine/api/v1.18.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.18.md \
 && wget -O allv/engine/api/v1.19.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.19.md \
 && wget -O allv/engine/api/v1.20.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.20.md \
 && wget -O allv/engine/api/v1.21.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.21.md \
 && wget -O allv/engine/api/v1.22.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.22.md \
 && wget -O allv/engine/api/v1.23.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.23.md \
 && wget -O allv/engine/api/v1.24.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/v1.24.md \
 && wget -O allv/engine/api/version-history.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/api/version-history.md \
 && wget -O allv/engine/reference/glossary.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/glossary.md \
 && wget -O allv/engine/reference/builder.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/builder.md \
 && wget -O allv/engine/reference/run.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/run.md \
 && wget -O allv/engine/reference/commandline/cli.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/reference/commandline/cli.md \
 && wget -O allv/engine/deprecated.md https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/docs/deprecated.md \
 && svn co https://github.com/docker/distribution/branches/$DISTRIBUTION_BRANCH/docs/spec allv/registry/spec \
 && wget -O allv/registry/configuration.md https://raw.githubusercontent.com/docker/distribution/$DISTRIBUTION_BRANCH/docs/configuration.md \
 && rm -rf allv/apidocs/cloud-api-source \
 && rm -rf allv/tests \
 && wget -O allv/engine/api/v1.25/swagger.yaml https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/api/swagger.yaml \
 && jekyll build -s allv -d allvbuild \
 && find allvbuild/engine/reference -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g' \
 && find allvbuild/engine/extend -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g' \
 && rm -rf allv

# Serve the site, which is now all static HTML
CMD jekyll serve -s /usr/src/app/allvbuild --no-watch -H 0.0.0.0 -P 4000
