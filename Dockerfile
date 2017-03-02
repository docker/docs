FROM starefossen/github-pages:112

# This is the source for docs/docs-base. Push to that location to ensure that
# the production image gets your update :)

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

COPY nginx.conf /etc/nginx/nginx.conf

## At the end of each layer, everything we need to pass on to the next layer
## should be in the "target" directory and we should have removed all temporary files

# Create archive; check out each version, create HTML under target/$VER, tweak links
# Nuke the archive_source directory. Only keep the target directory.

ENV VERSIONS="v1.4 v1.5 v1.6 v1.7 v1.8 v1.9 v1.10 v1.11 v1.12 v1.13"

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

# This index file gets overwritten, but it serves a sort-of useful purpose in
# making the docs/docs-base image browsable:

COPY index.html target

# Serve the site (target), which is now all static HTML

CMD echo "Docker docs are viewable at:" && echo "http://0.0.0.0:4000" && exec nginx
