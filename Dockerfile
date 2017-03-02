FROM starefossen/github-pages:onbuild

ONBUILD RUN git clone https://www.github.com/docker/docker.github.io docs

ONBUILD WORKDIR docs

ONBUILD RUN git checkout v1.13

ONBUILD COPY . /usr/src/app

CMD bundle update && bundle exec jekyll serve -d /_site --watch -H 0.0.0.0 -P 4000
