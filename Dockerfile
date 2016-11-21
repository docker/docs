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

# Get docker/docker ref docs from 1.12.x branch to be used in master builds
# Uses Github Subversion gateway to limit the checkout
RUN svn co https://github.com/docker/docker/branches/1.12.x/docs/reference allv/engine/reference \
 && svn co https://github.com/docker/docker/branches/1.12.x/docs/extend allv/engine/extend \
 && wget -O allv/engine/deprecated.md https://raw.githubusercontent.com/docker/docker/1.12.x/docs/deprecated.md \
 && jekyll build -s allv -d allvbuild \
 && rm -rf allv

# Serve the site, which is now all static HTML
CMD jekyll serve -s /usr/src/app/allvbuild -d /_site --no-watch -H 0.0.0.0 -P 4000
