FROM starefossen/github-pages

# Basic Git set-up for throwaway commits
RUN git config --global user.email "gordon@docker.com"
RUN git config --global user.name "Gordon"

# Clone the docs repo
RUN git clone https://www.github.com/docker/docker.github.io allv

# Get docker/docker ref docs from 1.12.x branch to be used in master builds
# Uses Github Subversion gateway to limit the checkout
RUN svn co https://github.com/docker/docker/branches/1.12.x/docs/reference allv/engine/reference
RUN svn co https://github.com/docker/docker/branches/1.12.x/docs/extend allv/engine/extend
# Can't use the svn trick to get a single file, use wget instead
RUN wget -O allv/engine/deprecated.md https://raw.githubusercontent.com/docker/docker/1.12.x/docs/deprecated.md
# Make a temporary commit for the files we added so we can check out other branches later
RUN git --git-dir=./allv/.git --work-tree=./allv commit -m "Temporary commit" -a

# Create HTML for master
RUN jekyll build -s allv -d allvbuild

# Check out 1.4 branch, create HTML, tweak links
RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.4
RUN mkdir allvbuild/v1.4
RUN jekyll build -s allv -d allvbuild/v1.4
RUN find allvbuild/v1.4 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.4/#g'
RUN find allvbuild/v1.4 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.4/#g'
RUN find allvbuild/v1.4 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.4/#g'

# Check out 1.5 branch, create HTML, tweak links
RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.5
RUN mkdir allvbuild/v1.5
RUN jekyll build -s allv -d allvbuild/v1.5
RUN find allvbuild/v1.5 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.5/#g'
RUN find allvbuild/v1.5 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.5/#g'
RUN find allvbuild/v1.5 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.5/#g'

# Check out 1.6, create HTML, tweak links
RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.6
RUN mkdir allvbuild/v1.6
RUN jekyll build -s allv -d allvbuild/v1.6
RUN find allvbuild/v1.6 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.6/#g'
RUN find allvbuild/v1.6 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.6/#g'
RUN find allvbuild/v1.6 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.6/#g'

# Check out 1.7, create HTML, tweak links
RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.7
RUN mkdir allvbuild/v1.7
RUN jekyll build -s allv -d allvbuild/v1.7
RUN find allvbuild/v1.7 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.7/#g'
RUN find allvbuild/v1.7 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.7/#g'
RUN find allvbuild/v1.7 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.7/#g'

# Check out 1.8, create HTML, tweak links
RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.8
RUN mkdir allvbuild/v1.8
RUN jekyll build -s allv -d allvbuild/v1.8
RUN find allvbuild/v1.8 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.8/#g'
RUN find allvbuild/v1.8 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.8/#g'
RUN find allvbuild/v1.8 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.8/#g'

# Check out 1.9, create HTML, tweak links
RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.9
RUN mkdir allvbuild/v1.9
RUN jekyll build -s allv -d allvbuild/v1.9
RUN find allvbuild/v1.9 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.9/#g'
RUN find allvbuild/v1.9 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.9/#g'
RUN find allvbuild/v1.9 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.9/#g'

# Check out 1.10, create HTML, tweak links
RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.10
RUN mkdir allvbuild/v1.10
RUN jekyll build -s allv -d allvbuild/v1.10
RUN find allvbuild/v1.10 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.10/#g'
RUN find allvbuild/v1.10 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.10/#g'
RUN find allvbuild/v1.10 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.10/#g'

# Check out 1.11, create HTML, tweak links
RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.11
RUN mkdir allvbuild/v1.11
RUN jekyll build -s allv -d allvbuild/v1.11
RUN find allvbuild/v1.11 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.11/#g'
RUN find allvbuild/v1.11 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.11/#g'
RUN find allvbuild/v1.11 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.11/#g'

# Serve the site, which is now all static HTML
CMD jekyll serve -s /usr/src/app/allvbuild -d /_site --no-watch -H 0.0.0.0 -P 4000
