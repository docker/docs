FROM starefossen/github-pages

RUN git clone https://www.github.com/docker/docker.github.io allv
RUN jekyll build -s allv -d allvbuild

RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.4
RUN mkdir allvbuild/v1.4
RUN jekyll build -s allv -d allvbuild/v1.4
RUN find allvbuild/v1.4 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.4/#g'
RUN find allvbuild/v1.4 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.4/#g'
RUN find allvbuild/v1.4 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.4/#g'

RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.5
RUN mkdir allvbuild/v1.5
RUN jekyll build -s allv -d allvbuild/v1.5
RUN find allvbuild/v1.5 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.5/#g'
RUN find allvbuild/v1.5 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.5/#g'
RUN find allvbuild/v1.5 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.5/#g'

RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.6
RUN mkdir allvbuild/v1.6
RUN jekyll build -s allv -d allvbuild/v1.6
RUN find allvbuild/v1.6 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.6/#g'
RUN find allvbuild/v1.6 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.6/#g'
RUN find allvbuild/v1.6 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.6/#g'

RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.7
RUN mkdir allvbuild/v1.7
RUN jekyll build -s allv -d allvbuild/v1.7
RUN find allvbuild/v1.7 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.7/#g'
RUN find allvbuild/v1.7 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.7/#g'
RUN find allvbuild/v1.7 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.7/#g'

RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.8
RUN mkdir allvbuild/v1.8
RUN jekyll build -s allv -d allvbuild/v1.8
RUN find allvbuild/v1.8 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.8/#g'
RUN find allvbuild/v1.8 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.8/#g'
RUN find allvbuild/v1.8 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.8/#g'

RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.9
RUN mkdir allvbuild/v1.9
RUN jekyll build -s allv -d allvbuild/v1.9
RUN find allvbuild/v1.9 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.9/#g'
RUN find allvbuild/v1.9 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.9/#g'
RUN find allvbuild/v1.9 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.9/#g'

RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.10
RUN mkdir allvbuild/v1.10
RUN jekyll build -s allv -d allvbuild/v1.10
RUN find allvbuild/v1.10 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.10/#g'
RUN find allvbuild/v1.10 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.10/#g'
RUN find allvbuild/v1.10 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.10/#g'

RUN git --git-dir=./allv/.git --work-tree=./allv checkout v1.11
RUN mkdir allvbuild/v1.11
RUN jekyll build -s allv -d allvbuild/v1.11
RUN find allvbuild/v1.11 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/v1.11/#g'
RUN find allvbuild/v1.11 -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/v1.11/#g'
RUN find allvbuild/v1.11 -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/v1.11/#g'

CMD jekyll serve -s /usr/src/app/allvbuild -d /_site --no-watch -H 0.0.0.0 -P 4000
