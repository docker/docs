FROM starefossen/github-pages:112

# Install nginx

ENV NGINX_VERSION 1.11.9-1~jessie

RUN apt-key adv --keyserver hkp://pgp.mit.edu:80 --recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62 \
	&& echo "deb http://nginx.org/packages/mainline/debian/ jessie nginx" >> /etc/apt/sources.list \
	&& apt-get update \
	&& apt-get install --no-install-recommends --no-install-suggests -y \
						ca-certificates \
						nginx=${NGINX_VERSION} \
	&& rm -rf /var/lib/apt/lists/*

# Forward nginx request and error logs to docker log collector

RUN ln -sf /dev/stdout /var/log/nginx/access.log \
	&& ln -sf /dev/stderr /var/log/nginx/error.log

## Branch to pull from, per ref doc
ENV ENGINE_BRANCH="1.13.x"
ENV DISTRIBUTION_BRANCH="release/2.5"

## At the end of each layer, everything we need to pass on to the next layer
## should be in the "target" directory and we should have removed all temporary files

# Copy master into target directory (skipping files / folders in .dockerignore)
# These files represent the current docs
COPY . target

# Move built html into md_source directory so we can reuse the target directory
# to hold the static output.
# Pull reference docs from upstream locations, then build the master docs
# into static HTML in the "target" directory using Jekyll
# then nuke the md_source directory.

RUN mv target md_source \
	&& svn co https://github.com/docker/docker/branches/$ENGINE_BRANCH/docs/extend md_source/engine/extend \
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
	&& wget -O md_source/engine/api/v1.25/swagger.yaml https://raw.githubusercontent.com/docker/docker/$ENGINE_BRANCH/api/swagger.yaml \
	&& jekyll build -s md_source -d target \
	&& rm -rf target/apidocs/layouts \
	&& find target -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g' \
	&& rm -rf md_source

# Create archive; check out each version, create HTML under target/$VER, tweak links
# Nuke the archive_source directory. Only keep the target directory.

ENV VERSIONS="v1.4 v1.5 v1.6 v1.7 v1.8 v1.9 v1.10 v1.11 v1.12"

RUN git clone https://www.github.com/docker/docker.github.io archive_source; \
 for VER in $VERSIONS; do \
		git --git-dir=./archive_source/.git --work-tree=./archive_source checkout ${VER} \
		&& mkdir -p target/${VER} \
		&& jekyll build -s archive_source -d target/${VER} \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/'"$VER"'/#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/'"$VER"'/#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/'"$VER"'/#g'; \
	done; \
	rm -rf archive_source

# Serve the site (target), which is now all static HTML

CMD echo "Server running at http://0.0.0.0:4000" && exec nginx -c /usr/src/app/target/nginx.conf
