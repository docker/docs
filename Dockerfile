FROM starefossen/github-pages:177

# This is the source for docs/docs-base. Push to that location to ensure that
# the production image gets your update :)

# Install nginx as well as packages used by
# _scripts/fetch-upstream-resources.sh (master branch)

RUN apk --no-cache add \
	git \
	nginx \
	subversion \
	wget

# Forward nginx request and error logs to docker log collector

RUN ln -sf /dev/stdout /var/log/nginx/access.log \
    && ln -sf /dev/stderr /var/log/nginx/error.log

COPY nginx.conf /etc/nginx/nginx.conf

## At the end of each layer, everything we need to pass on to the next layer
## should be in the "target" directory and we should have removed all temporary files

# Create archive; check out each version, create HTML under target/$VER, tweak links
# Nuke the archive_source directory. Only keep the target directory.

ENV VERSIONS="v17.06 v17.03 v1.4 v1.5 v1.6 v1.7 v1.8 v1.9 v1.10 v1.11 v1.12 v1.13"

## Use shallow clone and shallow check-outs to only get the tip of each branch

RUN git clone --depth 1 --recursive https://www.github.com/docker/docker.github.io archive_source; \
  for VER in $VERSIONS; do \
    git --git-dir=./archive_source/.git --work-tree=./archive_source fetch origin ${VER}:${VER} --depth 1 \
    && git --git-dir=./archive_source/.git --work-tree=./archive_source checkout ${VER} \
    && mkdir -p target/${VER} \
    && jekyll build -s archive_source -d target/${VER} \
		# Replace / rewrite some URLs so that links in the archive go to the correct
	  # location. Note that the order in which these replacements are done is
	  # important. Changing the order may result in replacements being done
		# multiple times.
		# First, remove the domain from URLs that include the domain
		&& BASEURL="$VER/" \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="http://docs-stage.docker.com/#href="/#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="https://docs-stage.docker.com/#src="/#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="https://docs.docker.com/#src="/#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="http://docs.docker.com/#href="/#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="http://docs.docker.com/#src="/#g' \
		\
		# Substitute https:// for schema-less resources (src="//analytics.google.com")
		# We're replacing them to prevent them being seen as absolute paths below
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="//#href="https://#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="//#src="https://#g' \
		\
		# And some archive versions already have URLs starting with '/version/'
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/'"$BASEURL"'#href="/#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/'"$BASEURL"'#src="/#g' \
		\
		# Archived versions 1.7 and under use some absolute links, and v1.10 uses
		# "relative" links to sources (href="./css/"). Remove those to make them
		# work :)
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="\./#href="/#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="\./#src="/#g' \
		\
		# Create permalinks for archived versions
		\
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/'"$BASEURL"'#g' \
		&& find target/${VER} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/'"$BASEURL"'#g'; \
  done; \
  rm -rf archive_source

# This index file gets overwritten, but it serves a sort-of useful purpose in
# making the docs/docs-base image browsable:

COPY index.html target

# Serve the site (target), which is now all static HTML

CMD echo "Docker docs are viewable at:" && echo "http://0.0.0.0:4000" && nginx -g 'pid /tmp/nginx.pid;'
